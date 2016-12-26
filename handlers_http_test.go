// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	disposable "github.com/0x19/disposable/protos"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHandleVerifyEmail(t *testing.T) {
	var service *Service
	var err error

	grpcPort := Random(6200, 6400)
	httpPort := grpcPort + 1

	go func(service *Service, err error) {
		os.Setenv("GRPC_ADDR", fmt.Sprintf(":%d", grpcPort))
		os.Setenv("HTTP_ADDR", fmt.Sprintf(":%d", httpPort))
		service, err = NewTestServer()
	}(service, err)

	time.Sleep(2 * time.Second)

	Convey("Disposable service started successfully", t, func() {
		So(service, ShouldHaveSameTypeAs, &Service{})
		So(err, ShouldBeNil)
	})

	Convey("Valid JSON is required", t, func() {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("http://localhost:%d/v1/verify", httpPort),
			nil,
		)
		So(err, ShouldBeNil)
		So(req, ShouldHaveSameTypeAs, &http.Request{})

		req.SetBasicAuth("disposable", "aa")

		client := &http.Client{}
		resp, err := client.Do(req)
		So(err, ShouldBeNil)
		So(resp, ShouldHaveSameTypeAs, &http.Response{})
		defer resp.Body.Close()

		So(resp.StatusCode, ShouldEqual, 401)
	})

	Convey("Valid JSON is required", t, func() {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("http://localhost:%d/v1/verify", httpPort),
			nil,
		)
		So(err, ShouldBeNil)
		So(req, ShouldHaveSameTypeAs, &http.Request{})

		req.SetBasicAuth("disposable", "disposable123")

		client := &http.Client{}
		resp, err := client.Do(req)
		So(err, ShouldBeNil)
		So(resp, ShouldHaveSameTypeAs, &http.Response{})
		defer resp.Body.Close()

		So(resp.StatusCode, ShouldEqual, 400)

		So(resp.Header["Content-Type"], ShouldNotBeNil)
		So(resp.Header["Content-Type"], ShouldResemble, []string{"application/json"})

		var dr disposable.DisposableResponse
		jrerr := DecodeJSONBody(&dr, resp.Body)
		So(jrerr, ShouldBeNil)

		So(dr.Status, ShouldBeFalse)
		So(dr.Error, ShouldHaveSameTypeAs, &disposable.Error{})
		So(dr.Error.Message, ShouldEqual, ErrorJSONParseError)
		So(dr.Error.Type, ShouldEqual, TypeJSONParseError)
	})

	Convey("Email is required", t, func() {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("http://localhost:%d/v1/verify", httpPort),
			bytes.NewBuffer([]byte(`{}`)),
		)
		So(err, ShouldBeNil)
		So(req, ShouldHaveSameTypeAs, &http.Request{})

		req.SetBasicAuth("disposable", "disposable123")

		client := &http.Client{}
		resp, err := client.Do(req)
		So(err, ShouldBeNil)
		So(resp, ShouldHaveSameTypeAs, &http.Response{})
		defer resp.Body.Close()

		So(resp.StatusCode, ShouldEqual, 400)

		So(resp.Header["Content-Type"], ShouldNotBeNil)
		So(resp.Header["Content-Type"], ShouldResemble, []string{"application/json"})

		var dr disposable.DisposableResponse
		jrerr := DecodeJSONBody(&dr, resp.Body)
		So(jrerr, ShouldBeNil)

		So(dr.Status, ShouldBeFalse)
		So(dr.Error, ShouldHaveSameTypeAs, &disposable.Error{})
		So(dr.Error.Message, ShouldEqual, ErrorInvalidEmailAddress)
		So(dr.Error.Type, ShouldEqual, TypeInvalidEmailAddress)
	})

	Convey("Valid email is required", t, func() {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("http://localhost:%d/v1/verify", httpPort),
			bytes.NewBuffer([]byte(`{"email": "brr"}`)),
		)
		So(err, ShouldBeNil)
		So(req, ShouldHaveSameTypeAs, &http.Request{})

		req.SetBasicAuth("disposable", "disposable123")

		client := &http.Client{}
		resp, err := client.Do(req)
		So(err, ShouldBeNil)
		So(resp, ShouldHaveSameTypeAs, &http.Response{})
		defer resp.Body.Close()

		So(resp.StatusCode, ShouldEqual, 400)

		So(resp.Header["Content-Type"], ShouldNotBeNil)
		So(resp.Header["Content-Type"], ShouldResemble, []string{"application/json"})

		var dr disposable.DisposableResponse
		jrerr := DecodeJSONBody(&dr, resp.Body)
		So(jrerr, ShouldBeNil)

		So(dr.Status, ShouldBeFalse)
		So(dr.Error, ShouldHaveSameTypeAs, &disposable.Error{})
		So(dr.Error.Message, ShouldEqual, ErrorInvalidEmailAddress)
		So(dr.Error.Type, ShouldEqual, TypeInvalidEmailAddress)
	})

	Convey("Valid email is required", t, func() {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("http://localhost:%d/v1/verify", httpPort),
			bytes.NewBuffer([]byte(`{"email": "buuf@wiki.8191.at"}`)),
		)
		So(err, ShouldBeNil)
		So(req, ShouldHaveSameTypeAs, &http.Request{})

		req.SetBasicAuth("disposable", "disposable123")

		client := &http.Client{}
		resp, err := client.Do(req)
		So(err, ShouldBeNil)
		So(resp, ShouldHaveSameTypeAs, &http.Response{})
		defer resp.Body.Close()

		So(resp.StatusCode, ShouldEqual, 400)

		So(resp.Header["Content-Type"], ShouldNotBeNil)
		So(resp.Header["Content-Type"], ShouldResemble, []string{"application/json"})

		var dr disposable.DisposableResponse
		jrerr := DecodeJSONBody(&dr, resp.Body)
		So(jrerr, ShouldBeNil)

		So(dr.Status, ShouldBeFalse)
		So(dr.Error, ShouldHaveSameTypeAs, &disposable.Error{})
		So(dr.Error.Message, ShouldEqual, ErrorDomainNotPermitted)
		So(dr.Error.Type, ShouldEqual, TypeDomainNotPermitted)
	})

	Convey("Passed email address is not blacklisted.", t, func() {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("http://localhost:%d/v1/verify", httpPort),
			bytes.NewBuffer([]byte(`{"email": "nevio.vesic@gmail.com"}`)),
		)
		So(err, ShouldBeNil)
		So(req, ShouldHaveSameTypeAs, &http.Request{})

		req.SetBasicAuth("disposable", "disposable123")

		client := &http.Client{}
		resp, err := client.Do(req)
		So(err, ShouldBeNil)
		So(resp, ShouldHaveSameTypeAs, &http.Response{})
		defer resp.Body.Close()

		So(resp.StatusCode, ShouldEqual, 200)

		So(resp.Header["Content-Type"], ShouldNotBeNil)
		So(resp.Header["Content-Type"], ShouldResemble, []string{"application/json"})

		var dr disposable.DisposableResponse
		jrerr := DecodeJSONBody(&dr, resp.Body)
		So(jrerr, ShouldBeNil)

		So(dr.Status, ShouldBeTrue)
		So(dr.Error, ShouldBeNil)
	})
}
