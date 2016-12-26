// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const ()

// NewTestServer -
func NewTestServer() (*Service, error) {
	var service *Service
	var err error

	if service, err = New(); err != nil {
		return nil, err
	}

	if err := service.Start(); err != nil {
		return nil, err
	}

	return service, nil
}

func TestService(t *testing.T) {
	var service *Service
	var err error

	go func(service *Service, err error) {
		os.Setenv("GRPC_ADDR", ":5432")
		os.Setenv("HTTP_ADDR", ":6432")

		Convey("goroutine block - ignore", t, func() {
			service, err = NewTestServer()
			defer service.CtxCancel()
			So(service, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

	}(service, err)

	time.Sleep(5 * time.Second)

	Convey("Disposable service started successfully", t, func() {
		So(service, ShouldHaveSameTypeAs, &Service{})
		So(err, ShouldBeNil)
	})
}
