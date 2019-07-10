// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade.  All rights reserved.
// MIT Licensed

package render

import (
	"net/http/httptest"
	"testing"
)

func TestText_Write(t *testing.T) {
	txt := Text{
		Body: "Hello World",
	}
	w := httptest.NewRecorder()
	body, err := txt.Write(w)

	if err != nil || string(body) != txt.Body {
		t.Error("Render text is not valid")
	}

	if w.Header().Get("Content-Type") != "text/plain" {
		t.Error("Content-Type Header is not set.")
	}
}
