// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	disposable "github.com/0x19/disposable/protos"
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// ValidateEmail -
func ValidateEmail(email string) *disposable.DisposableResponse {
	if !govalidator.IsEmail(email) {
		return &disposable.DisposableResponse{
			Status:    false,
			RequestId: uuid.NewV4().String(),
			Error:     disposable.NewError(ErrorInvalidEmailAddress, TypeInvalidEmailAddress, nil),
		}
	}
	return nil
}
