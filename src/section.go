package main

import "strings"

// SECTIONS

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

// VALUE SECTION

type ValueSection struct {
	Name string
}

func NewValueSection(name string) *ValueSection {
	s := new(ValueSection)
	s.Name = name
	return s
}

func (s *ValueSection) Render(context Context) string {
	val := context[s.Name]
	switch val.(type) {
	case nil:
		return ""
	case string:
		return val.(string)
	default:
		return ""
	}
}

// GROUP SECTION

type GroupSection struct {
	Name     string
	Sections []Section
}

func NewGroupSection(name string, sections []Section) *GroupSection {
	s := new(GroupSection)
	s.Name = name
	s.Sections = sections
	return s
}

func (s *GroupSection) Render(context Context) string {
	var sb strings.Builder
	gctx, ok := context[s.Name]
	if !ok {
		return ""
	}
	list, ok := gctx.(ContextList)
	if !ok {
		return ""
	}
	for _, c := range list {
		for _, section := range s.Sections {
			sb.WriteString(section.Render(c))
		}
	}
	return sb.String()
}
