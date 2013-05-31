/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package twilio is a library for interacting with the Twilio VoIP and SMS service.
//
// It provides helpers for writing succinct applications that serve TwiML responses.
package twilio

import (
	"bytes"
	"fmt"
	"net/http"
)

type Context interface {
	// Value returns the request form value for the specified key,
	// or an empty string if the key is not present in the request.
	Value(key string) string

	// Response sends a TwiML response to the Twilio service.
	Response(s string)

	// Responsef is like Response but takes a format string and arguments.
	// See the fmt package documentation for the format string syntax.
	Responsef(format string, args ...interface{})

	// Hangup is a convenience method that sends a <Hangup/> response.
	Hangup()
}

// HandlerFunc is a twilio handler function. It implements http.Handler.
type HandlerFunc func(Context)

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn(&context{w, r})
}

// Handle is a convenience function that registers the specified handler
// function under the given path using the net/http package's DefaultServeMux.
func Handle(path string, fn HandlerFunc) {
	http.Handle(path, fn)
}

type context struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *context) Value(key string) string {
	return c.r.FormValue(key)
}

func (c *context) Response(s string) {
	buf := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?><Response>`)
	buf.WriteString(s)
	buf.WriteString(`</Response>`)
	buf.WriteTo(c.w)
}

func (c *context) Responsef(format string, args ...interface{}) {
	c.Response(fmt.Sprintf(format, args...))
}

func (c *context) Hangup() {
	c.Response("<Hangup/>")
}
