// go-rs/rest-api-framework
// Copyright(c) 2019-2020 Roshan Gade. All rights reserved.
// MIT Licensed
package rest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type pattern struct {
	value  string
	regexp *regexp.Regexp
	keys   []string
}

const sep = "/"

// compile the pattern
func (p *pattern) compile() error {
	var err error
	pattern := ""
	p.keys = make([]string, 0)

	for _, val := range strings.Split(p.value, "/") {
		if val != "" {
			switch val[0] {
			case 94:
				pattern += "(?:/(.*))"

			case 42:
				pattern += "(?:/(.*))"
				p.keys = append(p.keys, "*")

			case 58:
				length := len(val)
				lastChar := val[length-1]
				if lastChar == 63 {
					pattern += "(?:/([^/]+?))?"
					p.keys = append(p.keys, val[1:(length-1)])
				} else {
					pattern += sep + "([^/]+?)"
					p.keys = append(p.keys, val[1:])
				}

			default:
				pattern += sep + val
			}
		}
	}

	p.regexp, err = regexp.Compile("^" + pattern + "/?$")
	return err
}

// match request URI with pattern
func (p *pattern) test(str string) bool {
	return p.regexp.MatchString(str)
}

// on URL path match, map every keys with pattern values
func (p *pattern) match(url string) map[string]string {
	if len(p.keys) == 0 {
		return nil
	}

	params := make(map[string]string)
	matches := p.regexp.FindAllSubmatch([]byte(url), -1)

	for i, k := range matches[0][1:] {
		params[p.keys[i]] = string(k)
		if p.keys[i] == "*" {
			params["*"] = sep + params["*"]
		}
	}
	fmt.Println(params)
	return params
}

// trim "/" from suffix
func trim(str string) string {
	if strings.HasSuffix(str, sep) {
		str = str[:len(str)-1]
	}
	return str
}

func jsonToBytes(data any) ([]byte, error) {
	if body, isString := data.(string); isString {
		return json.RawMessage(body).MarshalJSON()
	}

	//standard JSON as per RFC 7159
	return json.Marshal(data)
}

func xmlToBytes(data interface{}) ([]byte, error) {
	if body, isString := data.(string); isString {
		return []byte(body), nil
	}
	return xml.Marshal(data)
}
