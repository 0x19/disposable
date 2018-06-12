// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	disposable "github.com/0x19/disposable/protos"
	"github.com/asaskevich/govalidator"
)

// ValidateEmail -
func ValidateEmail(email string) *disposable.DisposableResponse {
	if !govalidator.IsEmail(email) {
		return &disposable.DisposableResponse{
			Status:    false,
			RequestId: GetUUID(),
			Error:     disposable.NewError(ErrorInvalidEmailAddress, TypeInvalidEmailAddress, nil),
		}
	}
	return nil
}
