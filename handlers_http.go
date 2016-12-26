// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"encoding/json"
	"net/http"
	"time"

	disposable "github.com/0x19/disposable/protos"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// HandleVerifyEmail -
func HandleVerifyEmail(s *Service, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)

	var vreq disposable.DisposableRequest

	if err := DecodeRequestBody(&vreq, req.Body); err != nil {
		j, derr := json.MarshalIndent(err, "", " ")
		if derr != nil {
			log.Errorf("[http_handle_verify_email] Unable to decode json into disposable request due to (err: %s)", derr)
			http.Error(w, derr.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	log.Infof("[http_handle_verify_email] Got new email verification request (req_vars: %+v) - (body: %+v)", vars, vreq)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	accresp, err := s.Verify(timeout, &vreq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.MarshalIndent(accresp, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !accresp.Status {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	w.Write(j)
	return
}
