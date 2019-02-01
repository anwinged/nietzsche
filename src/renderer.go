package main

import "strings"

type Context map[string]string

// COMPILE

func compile(template string) ([]Section, error) {
	var sections []Section
	tokenizer := NewTokenizer(template)
	for {
		token, err := tokenizer.Next()
		if err != nil {
			return []Section{}, err
		}
		if token == nil {
			break
		}
		switch token.Type {
		case TextToken:
			sections = append(sections, NewTextSection(token.Value))
		case ValueToken:
			sections = append(sections, NewTagSection(token.Value))
		}
	}
	return sections, nil
}

func Render(template string, params Context) (string, error) {
	sections, err := compile(template)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, s := range sections {
		sb.WriteString(s.Render(params))
	}
	return sb.String(), nil
}
