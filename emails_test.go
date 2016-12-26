// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDisposableEmails(t *testing.T) {
	var emails *DisposableEmails
	var err error

	Convey("Disposable emails loaded successfully", t, func() {
		emails, err = NewDisposableEmails()
		So(emails, ShouldNotBeNil)
		So(err, ShouldBeNil)

		err := emails.Load()
		So(err, ShouldBeNil)

		So(emails.Len(), ShouldBeGreaterThan, 0)
		So(emails.DomainExists("wiki.8191.at"), ShouldBeTrue)
		So(emails.DomainExists("zang.io"), ShouldBeFalse)
	})
}
