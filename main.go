// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import log "github.com/Sirupsen/logrus"

var (
	service *Service
)

func main() {
	var err error

	if service, err = New(); err != nil {
		log.Fatalf("Failed to initiate new service (err: %s)", err)
	}

	if err := service.Start(); err != nil {
		log.Fatalf("Failed to start service (err: %s)", err)
	}
}
