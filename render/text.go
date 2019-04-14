/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package render

import (
	"net/http"
)

type Text struct {
	Body string
}

var (
	plainType = "text/plain"
)

/**
 * Text Write
 */
func (j Text) Write(w http.ResponseWriter) ([]byte, error) {
	data := []byte(j.Body)
	w.Header().Set("Content-Type", plainType)
	return data, nil
}
