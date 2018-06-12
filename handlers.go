// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	disposable "github.com/0x19/disposable/protos"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// Verify -
func (s *Service) Verify(c context.Context, req *disposable.DisposableRequest) (*disposable.DisposableResponse, error) {
	log.Infof("[verify] Starting email verification process (req: %v)", req)

	if err := ValidateEmail(req.Email); err != nil {
		log.Errorf("[verify] Could not validate email address due to (err: %s)", err)
		return err, nil
	}

	if ok := s.DisposableEmails.IsOK(req.Email); !ok {
		log.Errorf("[verify] Seems like provided (email: %s) is illegal. Returning error now...", req.Email)
		return &disposable.DisposableResponse{
			Status:    false,
			RequestId: GetUUID(),
			Error:     disposable.NewError(ErrorDomainNotPermitted, TypeDomainNotPermitted, nil),
		}, nil
	}

	log.Infof("[verify] Domain verification passed for (email: %s)", req.Email)

	return &disposable.DisposableResponse{
		Status:    true,
		RequestId: GetUUID(),
	}, nil
}
