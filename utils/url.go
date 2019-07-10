// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade. All rights reserved.
// MIT Licensed

package utils

import (
	"regexp"
	"strings"
)

const sep = "/"

// Pre-compile pattern/string to avoid runtime or every request execution.
// Also, it returns keys, which are used in pattern, if any.
func Compile(str string) (regex *regexp.Regexp, keys []string, err error) {
	pattern := ""
	keys = make([]string, 0)

	for _, val := range strings.Split(str, "/") {
		if val != "" {
			switch val[0] {
			case 42:
				pattern += "(?:/(.*))"
				keys = append(keys, "*")

			case 58:
				length := len(val)
				lastChar := val[length-1]
				if lastChar == 63 {
					pattern += "(?:/([^/]+?))?"
					keys = append(keys, val[1:(length-1)])
				} else {
					pattern += sep + "([^/]+?)"
					keys = append(keys, val[1:])
				}

			default:
				pattern += sep + val
			}
		}
	}

	// if len(keys) == 0 {
	// 	pattern += "(?:/)?"
	// }

	regex, err = regexp.Compile("^" + pattern + "/?$")
	return
}

// On URL path match, map every keys with pattern values
func Exec(regex *regexp.Regexp, keys []string, uri []byte) (params map[string]string) {
	params = make(map[string]string)

	if len(keys) == 0 {
		return
	}

	matches := regex.FindAllSubmatch(uri, -1)

	for _, val := range matches {
		for i, k := range val[1:] {
			params[keys[i]] = string(k)
			if keys[i] == "*" {
				params["*"] = sep + params["*"]
			}
		}
	}

	return
}
