// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	disposable "github.com/0x19/disposable/protos"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TestVerifyEmail(t *testing.T) {
	var service *Service
	var err error

	grpcPort := Random(6200, 6400)
	httpPort := grpcPort + 1
	caCert := OptionString("GRPC_CA_FILE", "")
	caHost := OptionString("GRPC_HOST", "localhost")

	caCertt, _ := ioutil.ReadFile(caCert)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertt)

	go func(service *Service, err error) {
		os.Setenv("GRPC_ADDR", fmt.Sprintf(":%d", grpcPort))
		os.Setenv("HTTP_ADDR", fmt.Sprintf(":%d", httpPort))
		service, err = NewTestServer()
	}(service, err)

	time.Sleep(1 * time.Second)

	Convey("Account service started successfully", t, func() {
		So(service, ShouldHaveSameTypeAs, &Service{})
		So(err, ShouldBeNil)
	})

	Convey("Email address is required", t, func() {
		opts := []grpc.DialOption{}
		creds := credentials.NewClientTLSFromCert(caCertPool, caHost)
		opts = append(opts, grpc.WithTransportCredentials(creds))
		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(fmt.Sprintf(":%d", grpcPort), opts...)
		defer conn.Close()

		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		client := disposable.NewDisposableServiceClient(conn)
		resp, err := client.Verify(timeout, &disposable.DisposableRequest{})

		So(resp, ShouldHaveSameTypeAs, &disposable.DisposableResponse{})
		So(err, ShouldBeNil)
		So(resp.Status, ShouldBeFalse)

		// Ensure proper UUID is passed to us
		_, uuiderr := uuid.FromString(resp.RequestId)
		So(uuiderr, ShouldBeNil)

		So(resp.Error, ShouldHaveSameTypeAs, &disposable.Error{})
		So(resp.Error.Type, ShouldEqual, TypeInvalidEmailAddress)
		So(resp.Error.Message, ShouldEqual, ErrorInvalidEmailAddress)
	})

	Convey("Valid email address is required", t, func() {
		opts := []grpc.DialOption{}
		creds := credentials.NewClientTLSFromCert(caCertPool, caHost)
		opts = append(opts, grpc.WithTransportCredentials(creds))
		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(fmt.Sprintf(":%d", grpcPort), opts...)
		defer conn.Close()

		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		client := disposable.NewDisposableServiceClient(conn)
		resp, err := client.Verify(timeout, &disposable.DisposableRequest{
			Email: "buuf",
		})

		So(resp, ShouldHaveSameTypeAs, &disposable.DisposableResponse{})
		So(err, ShouldBeNil)
		So(resp.Status, ShouldBeFalse)

		// Ensure proper UUID is passed to us
		_, uuiderr := uuid.FromString(resp.RequestId)
		So(uuiderr, ShouldBeNil)

		So(resp.Error, ShouldHaveSameTypeAs, &disposable.Error{})
		So(resp.Error.Type, ShouldEqual, TypeInvalidEmailAddress)
		So(resp.Error.Message, ShouldEqual, ErrorInvalidEmailAddress)
	})

	Convey("Valid email address but under illegal domain", t, func() {
		opts := []grpc.DialOption{}
		creds := credentials.NewClientTLSFromCert(caCertPool, caHost)
		opts = append(opts, grpc.WithTransportCredentials(creds))
		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(fmt.Sprintf(":%d", grpcPort), opts...)
		defer conn.Close()

		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		client := disposable.NewDisposableServiceClient(conn)
		resp, err := client.Verify(timeout, &disposable.DisposableRequest{
			Email: "buuf@wiki.8191.at",
		})

		So(resp, ShouldHaveSameTypeAs, &disposable.DisposableResponse{})
		So(err, ShouldBeNil)
		So(resp.Status, ShouldBeFalse)

		// Ensure proper UUID is passed to us
		_, uuiderr := uuid.FromString(resp.RequestId)
		So(uuiderr, ShouldBeNil)

		So(resp.Error, ShouldHaveSameTypeAs, &disposable.Error{})
		So(resp.Error.Type, ShouldEqual, TypeDomainNotPermitted)
		So(resp.Error.Message, ShouldEqual, ErrorDomainNotPermitted)
	})

	Convey("Valid email address and valid domain", t, func() {
		opts := []grpc.DialOption{}
		creds := credentials.NewClientTLSFromCert(caCertPool, caHost)
		opts = append(opts, grpc.WithTransportCredentials(creds))
		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(fmt.Sprintf(":%d", grpcPort), opts...)
		defer conn.Close()

		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		client := disposable.NewDisposableServiceClient(conn)
		resp, err := client.Verify(timeout, &disposable.DisposableRequest{
			Email: "nevio.vesic@gmail.com",
		})

		So(resp, ShouldHaveSameTypeAs, &disposable.DisposableResponse{})
		So(err, ShouldBeNil)
		So(resp.Status, ShouldBeTrue)

		// Ensure proper UUID is passed to us
		_, uuiderr := uuid.FromString(resp.RequestId)
		So(uuiderr, ShouldBeNil)

		So(resp.Error, ShouldBeNil)
	})
}
