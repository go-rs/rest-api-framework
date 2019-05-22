/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package utils

import (
	"regexp"
	"testing"
)

func validateRegExp(t *testing.T, from string, regex *regexp.Regexp, _regex *regexp.Regexp, keys int, expKeys int, err error) {
	t.Log("case: ", from)
	if err != nil {
		t.Error("caught an error - ", err)
	}
	if keys != expKeys {
		t.Error("wrong number of keys - ", keys)
	}
	if regex.String() != _regex.String() {
		t.Error("regexp does not match - ", regex.String())
	}
}

// wildcard
func TestCompile1(t *testing.T) {
	var regex *regexp.Regexp
	var _regex *regexp.Regexp
	var keys []string
	var err error

	// case 1
	regex, keys, err = Compile("/*")
	_regex, _ = regexp.Compile("^(?:/(.*))/?$")
	validateRegExp(t, "wildcard 1", regex, _regex, len(keys), 1, err)

	// case 2
	regex, keys, err = Compile("/foo/:bar/*")
	_regex, _ = regexp.Compile("^/foo/([^/]+?)(?:/(.*))/?$")
	validateRegExp(t, "wildcard 2", regex, _regex, len(keys), 2, err)

	// case 3 - wrong example
	regex, keys, err = Compile("/:foo/:bar?/*")
	_regex, _ = regexp.Compile("^/([^/]+?)(?:/([^/]+?))?(?:/(.*))/?$")
	validateRegExp(t, "wildcard 3", regex, _regex, len(keys), 3, err)

}

// params
func TestCompile2(t *testing.T) {
	var regex *regexp.Regexp
	var _regex *regexp.Regexp
	var keys []string
	var err error

	// case 1
	regex, keys, err = Compile("/:foo/:bar")
	_regex, _ = regexp.Compile("^/([^/]+?)/([^/]+?)/?$")
	validateRegExp(t, "params 1", regex, _regex, len(keys), 2, err)

	// case 2
	regex, keys, err = Compile("/foo/:bar")
	_regex, _ = regexp.Compile("^/foo/([^/]+?)/?$")
	validateRegExp(t, "params 2", regex, _regex, len(keys), 1, err)

	// case 3
	regex, keys, err = Compile("/:foo/bar")
	_regex, _ = regexp.Compile("^/([^/]+?)/bar/?$")
	validateRegExp(t, "params 3", regex, _regex, len(keys), 1, err)

	// case 4
	regex, keys, err = Compile("/::foo")
	_regex, _ = regexp.Compile("^/([^/]+?)/?$")
	validateRegExp(t, "params 4", regex, _regex, len(keys), 1, err)

}

// plain text
func TestCompile3(t *testing.T) {
	var regex *regexp.Regexp
	var _regex *regexp.Regexp
	var keys []string
	var err error

	// case 1
	regex, keys, err = Compile("/foo/bar")
	_regex, _ = regexp.Compile("^/foo/bar/?$")
	validateRegExp(t, "plain text 1", regex, _regex, len(keys), 0, err)

	// case 2
	regex, keys, err = Compile("/foo*")
	_regex, _ = regexp.Compile("^/foo*/?$")
	validateRegExp(t, "plain text 2", regex, _regex, len(keys), 0, err)

	// case 2
	regex, keys, err = Compile("/")
	_regex, _ = regexp.Compile("^/?$")
	validateRegExp(t, "plain text 3", regex, _regex, len(keys), 0, err)

}

// optional param
func TestCompile4(t *testing.T) {
	var regex *regexp.Regexp
	var _regex *regexp.Regexp
	var keys []string
	var err error

	// case 1
	regex, keys, err = Compile("/:foo/:bar?")
	_regex, _ = regexp.Compile("^/([^/]+?)(?:/([^/]+?))?/?$")
	validateRegExp(t, "optional param 1", regex, _regex, len(keys), 2, err)

	// case 2 - wrong example (handle error)
	regex, keys, err = Compile("/:foo/:?")
	_regex, _ = regexp.Compile("^/([^/]+?)(?:/([^/]+?))?/?$")
	validateRegExp(t, "optional param 2", regex, _regex, len(keys), 2, err)
}

// wildcard character
var regex1, keys1, _ = Compile("/*")
var str1 = []byte("/foo/bar")

// wildcard
func TestExec1(t *testing.T) {
	var x map[string]string
	x = Exec(regex1, keys1, str1)

	if x["*"] != "/foo/bar" {
		t.Error("Error in asterisk")
	}
}

func BenchmarkExec1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Exec(regex1, keys1, str1)
	}
}

// params
var regex2, keys2, _ = Compile("/:foo/:bar")
var str2 = []byte("/foo/bar")

func TestExec2(t *testing.T) {
	var x map[string]string
	x = Exec(regex2, keys2, str2)

	if x["foo"] != "foo" || x["bar"] != "bar" {
		t.Error("Error in params")
	}
}

func BenchmarkExec2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Exec(regex2, keys2, str2)
	}
}

// plain text
var regex3, keys3, _ = Compile("/foo/bar")
var str3 = []byte("/foo/bar")

func TestExec3(t *testing.T) {
	var x map[string]string
	x = Exec(regex3, keys3, str3)

	if len(x) != 0 {
		t.Error("Error in plain text")
	}
}

func BenchmarkExec3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Exec(regex3, keys3, str3)
	}
}

// optional param
var regex4, keys4, _ = Compile("/foo/:bar?")
var str4 = []byte("/foo/bar")

func TestExec4(t *testing.T) {
	var x map[string]string
	x = Exec(regex4, keys4, str4)

	if x["bar"] != "bar" {
		t.Error("Error in optional")
	}

	x = Exec(regex4, keys4, []byte("/foo"))

	if x["bar"] != "" {
		t.Error("Error in optional")
	}
}

func BenchmarkExec4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Exec(regex4, keys4, str4)
	}
}
