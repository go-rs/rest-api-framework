/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package render

import (
	"encoding/json"
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

/**
 * JSON Write
 */
func (j JSON) Write(w http.ResponseWriter) ([]byte, error) {
	var data []byte
	var err error
	if reflect.TypeOf(j.Body).String() == "string" {
		rawIn := json.RawMessage(j.Body.(string))
		data, err = rawIn.MarshalJSON()
	} else {
		data, err = json.Marshal(j.Body)
	}
	if err != nil {
		return nil, err
	}
	w.Header().Set("Content-Type", jsonType)
	return data, nil
}
