package main

import "errors"
import "strings"

type Context map[string]string

type Section interface {
	Render(context Context) string
}

// TEXT SECTION

type TextSection struct {
	Text string
}

func NewTextSection(text string) *TextSection {
	s := new(TextSection)
	s.Text = text
	return s
}

func (s *TextSection) Render(context Context) string {
	return s.Text
}

// TAG SECTION

type TagSection struct {
	Name string
}

func NewTagSection(name string) *TagSection {
	s := new(TagSection)
	s.Name = name
	return s
}

func (s *TagSection) Render(context Context) string {
	return context[s.Name]
}

// COMPILE

func compile(template string) ([]Section, error) {
	var sections []Section
	var bracket = 0
	var name = ""
	var chunk = ""
	for _, char := range template {
		if char == '{' {
			if bracket == 0 {
				bracket = 1
			} else if bracket == 1 {
				bracket = 2
			} else {
				return nil, errors.New("Unexpected {")
			}
		} else if char == '}' {
			if bracket == 2 {
				bracket = 1
			} else if bracket == 1 {
				bracket = 0
			} else {
				return nil, errors.New("Unexpected }")
			}
		} else {
			if bracket == 2 {
				if chunk != "" {
					sections = append(sections, NewTextSection(chunk))
					chunk = ""
				}
				name += string(char)
			} else {
				if name != "" {
					sections = append(sections, NewTagSection(name))
					name = ""
				}
				chunk += string(char)
			}
		}
	}
	if name != "" {
		if bracket > 0 {
			return nil, errors.New("Unexpected end of placeholder")
		}
		sections = append(sections, NewTagSection(name))
	}
	if chunk != "" {
		sections = append(sections, NewTextSection(chunk))
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