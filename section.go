package main

import (
	"strconv"
	"strings"
)

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
	case int:
		return strconv.Itoa(val.(int))
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
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
		return s.renderBool(stack, groupContextList.(bool))
	case ValueList:
		return s.renderValueList(stack, groupContextList.(ValueList))
	case Context:
		return s.renderContext(stack, groupContextList.(Context))
	case ContextList:
		return s.renderContextList(stack, groupContextList.(ContextList))
	default:
		return ""
	}
}

func (s *GroupSection) renderBool(stack ContextStack, condition bool) string {
	if !condition {
		return ""
	}
	var sb strings.Builder
	for _, section := range s.Sections {
		sb.WriteString(section.Render(stack))
	}
	return sb.String()
}

func (s *GroupSection) renderValueList(stack ContextStack, list ValueList) string {
	var sb strings.Builder
	for range list {
		for _, section := range s.Sections {
			sb.WriteString(section.Render(stack))
		}
	}
	return sb.String()
}

func (s *GroupSection) renderContextList(stack ContextStack, list ContextList) string {
	var sb strings.Builder
	for _, context := range list {
		newStack := stack.PushContext(context)
		for _, section := range s.Sections {
			sb.WriteString(section.Render(newStack))
		}
	}
	return sb.String()
}

func (s *GroupSection) renderContext(stack ContextStack, context Context) string {
	var sb strings.Builder
	newStack := stack.PushContext(context)
	for _, section := range s.Sections {
		sb.WriteString(section.Render(newStack))
	}
	return sb.String()
}
