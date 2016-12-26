// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

const (
	ErrorInternalServerError = "Internal server error happen."
	ErrorInvalidEmailAddress = "Invalid email address provided."
	ErrorPageNotFound        = "Requested page is not found."
	ErrorJSONParseError      = "Make sure to provide valid JSON. Could not parse JSON request body."
	ErrorDomainNotPermitted  = "Domain address is not permitted!"
)

const (
	TypeInternalServerError = "E_INTERNAL_SERVER_ERROR"
	TypeInvalidEmailAddress = "E_INVALID_EMAIL_ADDRESS"
	TypePageNotFound        = "E_PAGE_NOT_FOUND"
	TypeJSONParseError      = "E_JSON_PARSE_ERROR"
	TypeDomainNotPermitted  = "E_DOMAIN_NOT_PERMITTED"
)
