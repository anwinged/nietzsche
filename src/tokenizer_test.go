package main

import "testing"

func assertTokenValues(t *testing.T, template string, expectedTokens []string) {
	tokenizer := NewTokenizer(template)
	index := 0
	for {
		token, err := tokenizer.Next()
		if token == nil {
			break
		}
		if err != nil {
			t.Errorf("Error not expexted, template '%s'", template)
			return
		}
		if index >= len(expectedTokens) {
			t.Errorf(
				"Unexpeted token, template '%s', expect nothing got '%s'",
				template, token.Value)
			return
		}
		expected := expectedTokens[index]
		if token.Value != expected {
			t.Errorf(
				"Unexpected token, template '%s', expected '%s', got '%s'",
				template, expected, token.Value)
			return
		}
		index++
	}
	length := len(expectedTokens)
	if index != length {
		t.Errorf(
			"Unexpected token count, template '%s', expected %d, got %d",
			template, length, index)
	}
}

func TestOneTextToken(t *testing.T) {
	assertTokenValues(
		t,
		"hello",
		[]string{"hello"},
	)
}

func TestOneValueTonen(t *testing.T) {
	assertTokenValues(
		t,
		"{{hello}}",
		[]string{"hello"},
	)
}
