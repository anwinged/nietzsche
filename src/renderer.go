package main

import "strings"

type Context map[string]interface{}

type ContextList []Context

// COMPILE

func Compile(template string) ([]Section, error) {
	var sections []Section
	tokens, err := Tokenize(template)
	if err != nil {
		return nil, err
	}
	var group bool = false
	var buffer []Section
	for _, token := range tokens {
		switch token.Type {
		case TextToken:
			if group {
				buffer = append(buffer, NewTextSection(token.Value))
			} else {
				sections = append(sections, NewTextSection(token.Value))
			}
		case ValueToken:
			if group {
				buffer = append(buffer, NewValueSection(token.Value))
			} else {
				sections = append(sections, NewValueSection(token.Value))
			}
		case OpenSectionToken:
			group = true
		case CloseSectionToken:
			group = false
			sections = append(sections, NewGroupSection(token.Value, buffer))
			buffer = []Section{}
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
