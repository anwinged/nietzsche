package main

import "strings"

type Context map[string]string

// COMPILE

func Compile(template string) ([]Section, error) {
	var sections []Section
	tokens, err := Tokenize(template)
	if err != nil {
		return nil, err
	}
	for _, token := range tokens {
		switch token.Type {
		case TextToken:
			sections = append(sections, NewTextSection(token.Value))
		case ValueToken:
			sections = append(sections, NewValueSection(token.Value))
		}
	}
	return sections, nil
}

// RENDER

func Render(template string, params Context) (string, error) {
	sections, err := Compile(template)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, s := range sections {
		sb.WriteString(s.Render(params))
	}
	return sb.String(), nil
}
