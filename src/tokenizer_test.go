package main

import "testing"
import "reflect"
import "strings"

func assertTokenValues(t *testing.T, template string, expected []Token) {
	tokens, err := Tokenize(template)
	if err != nil {
		t.Errorf("Error not expexted, template '%s'", template)
		return
	}
	if !reflect.DeepEqual(expected, tokens) {
		t.Errorf(
			"Unexpected tokens, template '%s', expect %s, got %s",
			template, strTokens(expected), strTokens(tokens),
		)
	}
}

func strTokens(tokens []Token) string {
	var sb strings.Builder
	sb.WriteString("[")
	for pos, t := range tokens {
		sb.WriteString("`" + t.Value + "`")
		if pos < len(tokens)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func assertErrorTokenization(t *testing.T, template string) {
	_, err := Tokenize(template)
	if err == nil {
		t.Errorf("Error expexted, but tokenization successful")
	}
}

func TestOneTextToken(t *testing.T) {
	assertTokenValues(
		t,
		"hello",
		[]Token{
			{Type: TextToken, Value: "hello"},
		},
	)
}

func TestOneValueToken(t *testing.T) {
	assertTokenValues(
		t,
		"{{hello}}",
		[]Token{
			{Type: ValueToken, Value: "hello"},
		},
	)
}

func TestSectionTokens(t *testing.T) {
	assertTokenValues(
		t,
		"{{#persons}}{{/persons}}",
		[]Token{
			{Type: OpenSectionToken, Value: "persons"},
			{Type: CloseSectionToken, Value: "persons"},
		},
	)
}

func TestInvertedTokens(t *testing.T) {
	assertTokenValues(
		t,
		"{{^persons}}{{/persons}}",
		[]Token{
			{Type: InvertedSectionToken, Value: "persons"},
			{Type: CloseSectionToken, Value: "persons"},
		},
	)
}

func TestComplexTokens(t *testing.T) {
	assertTokenValues(
		t,
		"Hi, {{name}}, we are {{#persons}}{{ name }}{{/ persons}}!",
		[]Token{
			{Type: TextToken, Value: "Hi, "},
			{Type: ValueToken, Value: "name"},
			{Type: TextToken, Value: ", we are "},
			{Type: OpenSectionToken, Value: "persons"},
			{Type: ValueToken, Value: "name"},
			{Type: CloseSectionToken, Value: "persons"},
			{Type: TextToken, Value: "!"},
		},
	)
}

func TestMissedOpenBrackets(t *testing.T) {
	assertErrorTokenization(
		t,
		"Hello, w}}orld",
	)
}

func TestMissedClosetBrackets(t *testing.T) {
	assertErrorTokenization(
		t,
		"Hello, {{world",
	)
}

func TestTooManyPairBrackets(t *testing.T) {
	assertErrorTokenization(
		t,
		"Hello, {{world}}}}",
	)
}
