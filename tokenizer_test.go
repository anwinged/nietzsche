package main

import (
	"reflect"
	"strings"
	"testing"
)

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

func TestNodeTokens(t *testing.T) {
	assertTokenValues(
		t,
		"{{#persons}}{{/persons}}",
		[]Token{
			{Type: OpenGroupToken, Value: "persons"},
			{Type: CloseGroupToken, Value: "persons"},
		},
	)
}

func TestInvertedTokens(t *testing.T) {
	assertTokenValues(
		t,
		"{{^persons}}{{/persons}}",
		[]Token{
			{Type: OpenInvertedGroupToken, Value: "persons"},
			{Type: CloseGroupToken, Value: "persons"},
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
			{Type: OpenGroupToken, Value: "persons"},
			{Type: ValueToken, Value: "name"},
			{Type: CloseGroupToken, Value: "persons"},
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

// Benchmarking, func number - number of tokens

func BenchmarkTokenize1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tokenize("{{x}}")
	}
}

func BenchmarkTokenize3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tokenize("Hello, {{world}}!!")
	}
}

func BenchmarkTokenize10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tokenize(
			`Hello, {{#persons}}{{fname}} {{lname}}, {{/persons}}!
			We are going to {{address}}.`,
		)
	}
}
