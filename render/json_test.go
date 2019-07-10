// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade. All rights reserved.
// MIT Licensed

package render

import (
	"net/http/httptest"
	"testing"
)

func TestJSON_Write1(t *testing.T) {
	json := JSON{
		Body: "Hello World",
	}
	w := httptest.NewRecorder()
	_, err := json.ToBytes(w)

	if err != errInvalidJson {
		t.Error("Should not render text")
	}

	if w.Header().Get("Content-Type") != "" {
		t.Error("Content-Type Header is not set.")
	}
}

func TestJSON_Write2(t *testing.T) {
	json := JSON{
		Body: "{\"Message\":\"Hello World\"}",
	}
	w := httptest.NewRecorder()
	_, err := json.ToBytes(w)

	if err != nil {
		t.Error("Should not throw an error")
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type Header is not set.")
	}
}

func TestJSON_Write3(t *testing.T) {
	json := JSON{
		Body: make(map[string]string),
	}
	w := httptest.NewRecorder()
	_, err := json.ToBytes(w)

	if err != nil {
		t.Error("Should not throw an error")
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type Header is not set.")
	}
}

func TestJSON_Write4(t *testing.T) {
	json := JSON{
		Body: true,
	}
	w := httptest.NewRecorder()
	_, err := json.ToBytes(w)

	if err != errInvalidJson {
		t.Error("Should not throw an error")
	}

	if w.Header().Get("Content-Type") != "" {
		t.Error("Content-Type Header is not set.")
	}
}
