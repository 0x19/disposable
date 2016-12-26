// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	graceful "gopkg.in/tylerb/graceful.v1"

	"github.com/gorilla/mux"
	"github.com/koding/cache"

	"github.com/0x19/disposable/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Service -
type Service struct {
	GRPCAddr, HTTPAddr string

	done             chan bool
	GRPCListener     net.Listener
	GRPC             *grpc.Server
	Cache            *cache.MemoryTTL
	Ctx              context.Context
	CtxCancel        context.CancelFunc
	DisposableEmails *DisposableEmails
}

// RegisterAndListenGrpcServer -
func (s *Service) RegisterAndListenGrpcServer() error {
	log.Infof("[register_grpc_server] Registering GRPC services...")
	var err error

	if s.GRPCListener, err = net.Listen("tcp", s.GRPCAddr); err != nil {
		return err
	}

	certfile := OptionString("GRPC_CA_FILE", "")

	if cert := OptionString("GRPC_CA", ""); cert != "" {
		cf, _ := ioutil.TempFile(os.TempDir(), "grpc_")
		cf.WriteString(cert)
		defer os.Remove(cf.Name())
		certfile = cf.Name()
	}

	keyfile := OptionString("GRPC_KEY_FILE", "")

	if key := OptionString("GRPC_KEY", ""); key != "" {
		kf, _ := ioutil.TempFile(os.TempDir(), "grpc_")
		kf.WriteString(key)
		defer os.Remove(kf.Name())
		keyfile = kf.Name()
	}

	creds, err := credentials.NewServerTLSFromFile(certfile, keyfile)
	if err != nil {
		log.Errorf("[register_grpc_server] Could not create new server from TLS due to (err: %s)", err)
		return err
	}

	s.GRPC = grpc.NewServer([]grpc.ServerOption{grpc.Creds(creds)}...)
	disposable.RegisterDisposableServiceServer(s.GRPC, s)

	return s.GRPC.Serve(s.GRPCListener)
}

// RegisterAndListenHTTPServer -
func (s *Service) RegisterAndListenHTTPServer() error {
	log.Infof("[register_http_server] Registering HTTP services...")

	httpmux := mux.NewRouter()
	httpmux.NotFoundHandler = http.HandlerFunc(HandleNotFound)

	httpmux.HandleFunc("/v1/verify", Use(func(w http.ResponseWriter, req *http.Request) {
		HandleVerifyEmail(s, w, req)
	}, BasicAuth, CapturePanic, CompressionHandler)).Methods("POST")

	srv := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    s.HTTPAddr,
			Handler: httpmux,
		},
	}

	return srv.ListenAndServe()
}

// HandleSigterm - Will basically wait for channel to close and than initiate
// service stop logic followed by actual exit.
func (s *Service) HandleSigterm(kill chan os.Signal) {
	<-kill
	s.CtxCancel()
}

// Start -
func (s *Service) Start() (err error) {
	log.Infof("[start] Starting up (grpc: %s) and (http: %s) services...", s.GRPCAddr, s.HTTPAddr)
	defer s.CtxCancel()

	errors := make(chan error)
	kill := make(chan os.Signal, 1)
	s.done = make(chan bool)

	signal.Notify(kill, os.Interrupt)
	signal.Notify(kill, syscall.SIGTERM)

	go s.HandleSigterm(kill)

	go func() { errors <- s.RegisterAndListenGrpcServer() }()
	go func() { errors <- s.RegisterAndListenHTTPServer() }()
	go func() { errors <- s.DisposableEmails.Load() }()

	select {
	case err := <-errors:
		log.Errorf("[start] Failed to serve services due to (err: %s)", err)
		s.CtxCancel()
		return err
	case <-s.Ctx.Done():
		log.Warn("[start] Received shutdown signal. Exiting service...")
		os.Exit(0)
	}

	return
}

// New -
func New() (*Service, error) {
	cache := cache.NewMemoryWithTTL(10 * time.Minute)
	cache.StartGC(1 * time.Second)

	emails, deerr := NewDisposableEmails()
	if deerr != nil {
		return nil, deerr
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		GRPCAddr:         OptionString("GRPC_ADDR", ":6874"),
		HTTPAddr:         OptionString("HTTP_ADDR", ":4432"),
		Cache:            cache,
		Ctx:              ctx,
		CtxCancel:        cancel,
		DisposableEmails: emails,
	}, nil
}
