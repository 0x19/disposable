// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"io"
	"net/http"
)

type WriteResetter interface {
	io.Writer
	Reset(io.Writer)
}

type CompressResponseWriter struct {
	WriteResetter
	http.ResponseWriter
}

func (w *CompressResponseWriter) WriteHeader(c int) {
	w.ResponseWriter.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(c)
}

func (w *CompressResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *CompressResponseWriter) Write(b []byte) (int, error) {
	h := w.ResponseWriter.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(b))
	}
	h.Del("Content-Length")

	return w.WriteResetter.Write(b)
}
