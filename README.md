# okws

An ok webserver.

Whenever I'm unsatisfied with `python -m http.server`
I try to patch my wanted features into a Go server.

## extra features so far

- Use index.html if exists. Otherwise:
  - Generate a file list that lazy-loads all images as embedded
- Print best-effort outbound IP as default bind (as opposed to 0.0.0.0)

## installation

Although `okws` is, for convenience, a standalone executable
that does not need installation, you may want to
install to your `$PATH` anyway.

```bash
# Linux
go build
sudo install -Dm755 okws /usr/local/bin/okws
```
