package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	port int
	bind string
	dir  string
)

var rootCmd = &cobra.Command{
	Use:   "okws",
	Short: "A simple HTTP server",
	Long:  "A simple HTTP server that serves files from a specified directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("error getting directory: %w", err)
		}

		if err := os.Chdir(absDir); err != nil {
			return fmt.Errorf("could not chdir: %w", err)
		}

		http.HandleFunc("/", directoryHandler(absDir))

		address := fmt.Sprintf("%s:%d", bind, port)
		fmt.Printf("Serving directory browser on http://%s from '%s'\n", address, absDir)

		return http.ListenAndServe(address, nil)
	},
}

// cobra stuff
func init() {
	rootCmd.PersistentFlags().IntVar(&port, "port", 8000, "bind to this port (default: 8000)")
	rootCmd.PersistentFlags().StringVar(&bind, "bind", "", "bind to this address (default: all interfaces)")
	rootCmd.PersistentFlags().StringVar(&dir, "directory", ".", "serve this directory (default: current directory)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// dirty way to embed html into go
// TODO: refactor using embedded html lib or something like that
var tmpl = template.Must(template.New("dir").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>{{ .Path }}</title>
<style>
img {
    max-width: 75vw;
    margin: 10px 0;
}
</style>
</head>
<body>

<h2>{{ .Path }}</h2>

<div>
{{ if .Parent }}
  <a href="{{ .Parent }}">&lt;-- parent directory</a><br><br>
{{ end }}

{{ range .Dirs }}
  <a href="{{ . }}/">{{ . }}/</a><br>
{{ end }}
</div>

{{ range .Images }}
  <img loading="lazy" src="{{ . }}" alt="" title="{{ . }}">
{{ end }}

{{ range .Files }}
  <a href="{{ . }}">{{ . }}</a><br>
{{ end }}

</body>
</html>
`))

func directoryHandler(base string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestedPath := filepath.Clean(r.URL.Path)

		fsPath := filepath.Join(base, requestedPath)

		// serve real files directly
		if fileExists(fsPath) && !isDir(fsPath) {
			http.ServeFile(w, r, fsPath)
			return
		}
		if fileExists(fsPath + "/index.html") {
			http.ServeFile(w, r, fsPath+"/index.html")
			return
		}
		if fileExists(fsPath + "/index.htm") {
			http.ServeFile(w, r, fsPath+"/index.htm")
			return
		}

		// render the directory otherwise
		entries, err := os.ReadDir(fsPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// we dont need a full dfs/bfs search because the url request
		// worries about that for us
		// btw its already sorted https://gohugo.io/functions/os/readdir/
		var dirs, images, files []string
		for _, e := range entries {
			name := e.Name()
			relativePath := filepath.Join(requestedPath, name)
			if e.IsDir() {
				dirs = append(dirs, relativePath)
			} else {
				ext := strings.ToLower(filepath.Ext(name))
				if isImage(ext) {
					images = append(images, relativePath)
				} else {
					files = append(files, relativePath)
				}
			}
		}

		// build parent directory link
		var parent string
		if requestedPath != "/" {
			parent = filepath.Dir(requestedPath)
			if parent == "." {
				parent = "/"
			}
		}

		data := struct {
			Path   string
			Parent string
			Dirs   []string
			Images []string
			Files  []string
		}{
			Path:   requestedPath,
			Parent: parent,
			Dirs:   dirs,
			Images: images,
			Files:  files,
		}

		tmpl.Execute(w, data)
	}
}

// helper funcs
func isImage(ext string) bool {
	switch ext {
	// NOTE: heic is a garbage file format but i keep it in here anyway in case
	//       your firefox extension supports rendering it
	// TODO: could add all sensible image formats here
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".exif", ".heic":
		return true
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}
