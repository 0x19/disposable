// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// DisposableEmails -
type DisposableEmails struct {
	Emails []string
	Source string
}

// Load -
func (de *DisposableEmails) Load() error {
	log.Infof("[load] About to start loading disposable emails from (source: %s)", de.Source)

	data, err := ioutil.ReadFile(de.Source)
	if err != nil {
		log.Errorf("[load] Failed to load (source: %s) due to (err: %s)", de.Source, err)
		return err
	}

	de.Emails = strings.Split(string(data), "\n")
	return nil
}

// GetAll -
func (de *DisposableEmails) GetAll() []string {
	return de.Emails
}

// DomainExists -
func (de *DisposableEmails) DomainExists(domain string) bool {
	return StringInSlice(domain, de.Emails)
}

// IsOK -
func (de *DisposableEmails) IsOK(email string) bool {
	for _, domain := range de.Emails {
		if strings.Contains(email, domain) {
			return false
		}
	}

	return true
}

// Len -
func (de *DisposableEmails) Len() int {
	return len(de.Emails)
}

// DisposableEmails -
func NewDisposableEmails() (*DisposableEmails, error) {
	return &DisposableEmails{
		Source: OptionString("DISPOSABLE_EMAILS_SOURCE", "services/burner/emails.txt"),
	}, nil
}
