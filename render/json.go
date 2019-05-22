/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
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

var (
	jsonType = "application/json"
)

var (
	invalidJson = errors.New("INVALID_JSON_RESPONSE")
)

/**
 * JSON Write
 */
func (j JSON) Write(w http.ResponseWriter) (data []byte, err error) {
	_type := reflect.TypeOf(j.Body).String()
	if _type == "int" || _type == "float64" || _type == "bool" {
		err = invalidJson
	} else if _type == "string" {
		data, err = json.RawMessage(j.Body.(string)).MarshalJSON()
	} else {
		data, err = json.Marshal(j.Body)
	}

	if err != nil {
		return
	}

	if !json.Valid(data) {
		err = invalidJson
		return
	}

	w.Header().Set("Content-Type", jsonType)
	return
}
