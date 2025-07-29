package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	Run: func(cmd *cobra.Command, args []string) {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			fmt.Println("Error getting absolute directory:", err)
			return
		}
		os.Chdir(absDir)

		handler := http.FileServer(http.Dir(absDir))
		address := fmt.Sprintf("%s:%d", bind, port)

		fmt.Printf("Serving on %s from directory '%s'...\n", address, absDir)
		if err := http.ListenAndServe(address, handler); err != nil {
			fmt.Println("Error starting server:", err)
		}
	},
}

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
