// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"os"
	"strconv"
)

// OptionString - Will return ENV value back based on provided name casted as string
func OptionString(name, def string) string {
	if res := os.Getenv(name); res != "" {
		return string(res)
	}

	return def
}

// OptionBool - Will return ENV value back based on provided name casted as bool
func OptionBool(name string, def bool) bool {
	var res string

	if res = os.Getenv(name); res == "" {
		return def
	}

	b, err := strconv.ParseBool(res)
	if err != nil {
		return def
	}

	return b
}

// OptionInt - Will return ENV value back based on provided name casted as int64
func OptionInt(name string, def int) int {
	var res string

	if res = os.Getenv(name); res == "" {
		return def
	}

	b, err := strconv.Atoi(res)
	if err != nil {
		return def
	}

	return b
}
