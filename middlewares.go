// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	disposable "github.com/0x19/disposable/protos"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	"github.com/zang-cloud/micro-common/options"
)

// use provides a cleaner interface for chaining middleware for single routes.
// Middleware functions are simple HTTP handlers (w http.ResponseWriter, r *http.Request)
//
//  r.HandleFunc("/login", use(loginHandler, rateLimit, csrf))
//  r.HandleFunc("/form", use(formHandler, csrf))
//  r.HandleFunc("/about", aboutHandler)
func Use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

// HandleNotFound -
func HandleNotFound(w http.ResponseWriter, req *http.Request) {
	log.Infof("[http_handle_not_found] Received not-found (request: %+v)", req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	resp := &disposable.DisposableResponse{
		Status:    false,
		RequestId: uuid.NewV4().String(),
		Error:     disposable.NewError(ErrorPageNotFound, TypePageNotFound, nil),
	}

	j, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(j)
	return
}

// CapturePanic -
func CapturePanic(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("[capture_panic] Setting up middleware...")
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, true)]

				log.Errorf("[capture_panic] Got panic (err: %s) - (stack: %s)", err, stack)

				errresp := &disposable.DisposableResponse{
					Status:    false,
					RequestId: uuid.NewV4().String(),
					Error:     disposable.NewError(ErrorInternalServerError, TypeInternalServerError, errors.New(fmt.Sprintf("%v", err))),
				}

				j, err := json.MarshalIndent(errresp, "", " ")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Write(j)
				return
			}
		}()

		h.ServeHTTP(w, r)
	}
}

// BasicAuth -
func BasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("[basic_auth] Setting up middleware...")
		username := options.OptionString("HTTP_BASIC_USERNAME", "disposable")
		password := options.OptionString("HTTP_BASIC_PASSWORD", "disposable123")

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != username || pair[1] != password {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func CompressionHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("[compression_handler] Setting up middleware...")

	L:
		for _, enc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			switch strings.TrimSpace(enc) {
			case "gzip":
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Add("Vary", "Accept-Encoding")

				gw := gzip.NewWriter(w)
				defer gw.Close()

				w = &CompressResponseWriter{
					WriteResetter:  gw,
					ResponseWriter: w,
				}

				break L
			case "deflate":
				w.Header().Set("Content-Encoding", "deflate")
				w.Header().Add("Vary", "Accept-Encoding")

				fw, _ := flate.NewWriter(w, flate.DefaultCompression)
				defer fw.Close()

				w = &CompressResponseWriter{
					WriteResetter:  fw,
					ResponseWriter: w,
				}

				break L
			}
		}

		h.ServeHTTP(w, r)
	}
}
