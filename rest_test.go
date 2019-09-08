package rest

import (
	"testing"
)

func TestNew(t *testing.T) {
	api := New("/")

	if api == nil {
		t.Error("New function should return API reference.")
	}
}
