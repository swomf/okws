# okws

An ok webserver.

Whenever I'm unsatisfied with `python -m http.server`
I try to patch my wanted features into a Go server.

## extra features so far

- Lazy-loaded image embedding

## installation

```bash
go build
sudo install -Dm755 okws /usr/local/bin/okws
```
