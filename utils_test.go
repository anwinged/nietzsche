package main

import "testing"

func TestShortenLongText(t *testing.T) {
	text := "This is a long text"
	length := 5
	expected := "This ..."
	result := Shorten(text, length)
	if expected != result {
		t.Errorf("Shorten fails. Expect \"%s\", got \"%s\"", expected, result)
	}
}

func TestShortenShortText(t *testing.T) {
	text := "This is a short text"
	length := 100
	expected := "This is a short text"
	result := Shorten(text, length)
	if expected != result {
		t.Errorf("Shorten fails. Expect \"%s\", got \"%s\"", expected, result)
	}
}
