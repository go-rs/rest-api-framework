// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade.  All rights reserved.
// MIT Licensed

package render

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

//TODO: JSONP
type JSON struct {
	Body interface{}
}

const (
	jsonType = "application/json"
)

var (
	ErrInvalidJson = errors.New("INVALID_JSON_RESPONSE")
)

/**
 * JSON Write
 */
func (j JSON) Write(w http.ResponseWriter) (data []byte, err error) {
	_type := reflect.TypeOf(j.Body).String()
	if _type == "int" || _type == "float64" || _type == "bool" {
		err = ErrInvalidJson
	} else if _type == "string" {
		data, err = json.RawMessage(j.Body.(string)).MarshalJSON()
	} else {
		data, err = json.Marshal(j.Body)
	}

	if err != nil {
		return
	}

	if !json.Valid(data) {
		err = ErrInvalidJson
		return
	}

	w.Header().Set("Content-Type", jsonType)
	return
}
