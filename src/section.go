package main

import "strings"

// SECTIONS

type Section interface {
	Render(stack ContextStack) string
}

// TEXT SECTION

type TextSection struct {
	Text string
}

func NewTextSection(text string) *TextSection {
	return &TextSection{Text: text}
}

func (s *TextSection) Render(stack ContextStack) string {
	return s.Text
}

// VALUE SECTION

type ValueSection struct {
	Name string
}

func NewValueSection(name string) *ValueSection {
	return &ValueSection{Name: name}
}

func (s *ValueSection) Render(stack ContextStack) string {
	val := stack.FindValue(s.Name)
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
	return &GroupSection{Name: name, Sections: sections}
}

func (s *GroupSection) Render(stack ContextStack) string {
	groupContextList := stack.FindValue(s.Name)
	if groupContextList == nil {
		return ""
	}
	switch groupContextList.(type) {
	case bool:
		condition := groupContextList.(bool)
		if !condition {
			return ""
		}
		var sb strings.Builder
		for _, section := range s.Sections {
			sb.WriteString(section.Render(stack))
		}
		return sb.String()
	case ContextList:
		var sb strings.Builder
		for _, context := range groupContextList.(ContextList) {
			newStack := stack.PushContext(context)
			for _, section := range s.Sections {
				sb.WriteString(section.Render(newStack))
			}
		}
		return sb.String()
	default:
		return ""
	}
}
