# snai.pe/go-htutil

[![GoDoc](https://godoc.org/snai.pe/go-htutil?status.svg)](https://godoc.org/snai.pe/go-htutil)  

```
go get snai.pe/go-htutil
```

Go HTTP utilities with no dependencies.

This package provides the following utilities:

* an alternate implemenation of github.com/golang/gddo/httputil.NegotiateContentType.
  This package was in part motivated in providing a no-dependency package providing
  similar functionality.
* a wrapper over `net/url.URL` that implements `encoding.TextMarshaler` and `encoding.TextUnmarshaler`.
