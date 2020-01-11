package rest

import "testing"

func TestPattern_compile1(t *testing.T) {
	p := pattern{
		value: "/",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.regexp.String() != "^/?$" {
		t.Error("Pattern is not compiled successfully.", p.regexp.String())
	}

	if len(p.keys) > 0 {
		t.Error("Pattern should not return any key.", p.keys)
	}
}

func TestPattern_compile2(t *testing.T) {
	p := pattern{
		value: "/page/:number",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.regexp.String() != "^/page/([^/]+?)/?$" {
		t.Error("Pattern is not compiled successfully.", p.regexp.String())
	}

	if len(p.keys) != 1 && p.keys[0] == "number" {
		t.Error("Pattern should not return any key.", p.keys)
	}
}

func TestPattern_compile3(t *testing.T) {
	p := pattern{
		value: "/page/:number/:limit?",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.regexp.String() != "^/page/([^/]+?)(?:/([^/]+?))?/?$" {
		t.Error("Pattern is not compiled successfully.", p.regexp.String())
	}

	if len(p.keys) != 2 && p.keys[0] == "number" && p.keys[1] == "limit" {
		t.Error("Pattern should not return any key.", p.keys)
	}
}

func TestPattern_compile4(t *testing.T) {
	p := pattern{
		value: "/page/:number/type/:category",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.regexp.String() != "^/page/([^/]+?)/type/([^/]+?)/?$" {
		t.Error("Pattern is not compiled successfully.", p.regexp.String())
	}

	if len(p.keys) != 2 && p.keys[0] == "number" && p.keys[1] == "category" {
		t.Error("Pattern should not return any key.", p.keys)
	}
}

func TestPattern_compile5(t *testing.T) {
	p := pattern{
		value: "/page/:number/*",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.regexp.String() != "^/page/([^/]+?)(?:/(.*))/?$" {
		t.Error("Pattern is not compiled successfully.", p.regexp.String())
	}

	if len(p.keys) != 2 && p.keys[0] == "number" && p.keys[1] == "*" {
		t.Error("Pattern should not return any key.", p.keys)
	}
}

func TestPattern_test1(t *testing.T) {
	p := pattern{
		value: "/",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.test("/") != true {
		t.Error("Unable to test regular expression for '/'")
	}
}

func TestPattern_test2(t *testing.T) {
	p := pattern{
		value: "/page/:number",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.test("/page/1") != true {
		t.Error("Unable to test regular expression for '/page/:number'")
	}
}

func TestPattern_test3(t *testing.T) {
	p := pattern{
		value: "/page/:number/:limit?",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.test("/page/1") != true {
		t.Error("Unable to test regular expression for '/page/:number/:limit?'")
	}

	if p.test("/page/1/10") != true {
		t.Error("Unable to test regular expression for '/page/:number/:limit?'")
	}
}

func TestPattern_test4(t *testing.T) {
	p := pattern{
		value: "/page/:number/type/:category",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.test("/page/1/type/go") != true {
		t.Error("Unable to test regular expression for '/page/:number/type/:category'")
	}
}

func TestPattern_test5(t *testing.T) {
	p := pattern{
		value: "/page/:number/*",
	}

	err := p.compile()

	if err != nil {
		t.Error("Failed to compile pattern", err)
	}

	if p.test("/page/1/posts") != true {
		t.Error("Unable to test regular expression for '/page/:number'")
	}
}
