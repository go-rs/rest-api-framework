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

const (
	plainType = "text/plain;charset=UTF-8"
)

/**
 * Text Write
 */
func (j Text) Write(w http.ResponseWriter) (data []byte, err error) {
	data = []byte(j.Body)
	w.Header().Set("Content-Type", plainType)
	return
}
