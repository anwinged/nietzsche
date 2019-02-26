package main

import (
	"strconv"
	"strings"
)

// SECTIONS

type Section interface {
	Type() string
	Desc() string
	Render(stack ContextStack) string
	Sections() []Section
}

// TEXT SECTION

type TextSection struct {
	Text string
}

func NewTextSection(text string) *TextSection {
	return &TextSection{Text: text}
}

func (s *TextSection) Type() string {
	return "TEXT"
}

func (s *TextSection) Desc() string {
	return "\"" + Shorten(s.Text, 15) + "\""
}

func (s *TextSection) Render(stack ContextStack) string {
	return s.Text
}

func (s *TextSection) Sections() []Section {
	return []Section{}
}

// VALUE SECTION

type ValueSection struct {
	Name string
}

func NewValueSection(name string) *ValueSection {
	return &ValueSection{Name: name}
}

func (s *ValueSection) Type() string {
	return "VALUE"
}

func (s *ValueSection) Desc() string {
	return s.Name
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

func (s *ValueSection) Sections() []Section {
	return []Section{}
}

// GROUP SECTION

type GroupSection struct {
	Name string
	Sns  []Section
}

func NewGroupSection(name string, sections []Section) *GroupSection {
	return &GroupSection{Name: name, Sns: sections}
}

func (s *GroupSection) Type() string {
	return "SECTION"
}

func (s *GroupSection) Desc() string {
	return s.Name
}

func (s *GroupSection) Sections() []Section {
	return s.Sns
}

func (s *GroupSection) Render(stack ContextStack) string {
	groupContextList := stack.FindValue(s.Name)
	if groupContextList == nil {
		return ""
	}
	switch groupContextList.(type) {
	case bool:
		return s.renderBool(stack, groupContextList.(bool))
	case int:
		return s.renderBool(stack, groupContextList.(int) != 0)
	case float64:
		return s.renderBool(stack, groupContextList.(float64) != 0.0)
	case string:
		return s.renderBool(stack, groupContextList.(string) != "")
	case map[string]interface{}:
		return s.renderContext(stack, groupContextList.(map[string]interface{}))
	case []interface{}:
		return s.renderValueList(stack, groupContextList.([]interface{}))
	default:
		return ""
	}
}

func (s *GroupSection) renderBool(stack ContextStack, condition bool) string {
	if !condition {
		return ""
	}
	var sb strings.Builder
	for _, section := range s.Sns {
		sb.WriteString(section.Render(stack))
	}
	return sb.String()
}

func (s *GroupSection) renderValueList(stack ContextStack, list []interface{}) string {
	var sb strings.Builder
	for _, el := range list {
		for _, section := range s.Sns {
			casted, ok := el.(map[string]interface{})
			if ok {
				newStack := stack.PushContext(casted)
				sb.WriteString(section.Render(newStack))
			} else {
				sb.WriteString(section.Render(stack))
			}
		}
	}
	return sb.String()
}

func (s *GroupSection) renderContext(stack ContextStack, context map[string]interface{}) string {
	var sb strings.Builder
	newStack := stack.PushContext(context)
	for _, section := range s.Sns {
		sb.WriteString(section.Render(newStack))
	}
	return sb.String()
}
