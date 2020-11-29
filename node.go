package nietzsche

import (
	"strconv"
	"strings"
)

// SECTIONS

type Node interface {
	Type() string
	Desc() string
	Render(stack ContextStack) string
	Nodes() []Node
}

// TEXT SECTION

type TextNode struct {
	Text string
}

func NewTextNode(text string) *TextNode {
	return &TextNode{Text: text}
}

func (s *TextNode) Type() string {
	return "TEXT"
}

func (s *TextNode) Desc() string {
	return "\"" + Shorten(s.Text, 15) + "\""
}

func (s *TextNode) Render(stack ContextStack) string {
	return s.Text
}

func (s *TextNode) Nodes() []Node {
	return []Node{}
}

// VALUE SECTION

type ValueNode struct {
	Name string
}

func NewValueNode(name string) *ValueNode {
	return &ValueNode{Name: name}
}

func (s *ValueNode) Type() string {
	return "VALUE"
}

func (s *ValueNode) Desc() string {
	return s.Name
}

func (s *ValueNode) Render(stack ContextStack) string {
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

func (s *ValueNode) Nodes() []Node {
	return []Node{}
}

// GROUP SECTION

type GroupNode struct {
	Name string
	Sns  []Node
}

func NewGroupNode(name string, sections []Node) *GroupNode {
	return &GroupNode{Name: name, Sns: sections}
}

func (s *GroupNode) Type() string {
	return "GROUP"
}

func (s *GroupNode) Desc() string {
	return s.Name
}

func (s *GroupNode) Nodes() []Node {
	return s.Sns
}

func (s *GroupNode) Render(stack ContextStack) string {
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

func (s *GroupNode) renderBool(stack ContextStack, condition bool) string {
	if !condition {
		return ""
	}
	var sb strings.Builder
	for _, section := range s.Sns {
		sb.WriteString(section.Render(stack))
	}
	return sb.String()
}

func (s *GroupNode) renderValueList(stack ContextStack, list []interface{}) string {
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

func (s *GroupNode) renderContext(stack ContextStack, context map[string]interface{}) string {
	var sb strings.Builder
	newStack := stack.PushContext(context)
	for _, section := range s.Sns {
		sb.WriteString(section.Render(newStack))
	}
	return sb.String()
}

// NEGATIVE SECTION

type NegativeNode struct {
	Name string
	Sns  []Node
}

func NewNegativeNode(name string, sections []Node) *NegativeNode {
	return &NegativeNode{Name: name, Sns: sections}
}

func (s *NegativeNode) Type() string {
	return "NEGATIVE"
}

func (s *NegativeNode) Desc() string {
	return s.Name
}

func (s *NegativeNode) Nodes() []Node {
	return s.Sns
}

func (s *NegativeNode) Render(stack ContextStack) string {
	groupContextList := stack.FindValue(s.Name)
	if groupContextList == nil {
		return ""
	}
	switch groupContextList.(type) {
	case bool:
		return s.renderNegativeBool(stack, groupContextList.(bool))
	case int:
		return s.renderNegativeBool(stack, groupContextList.(int) != 0)
	case float64:
		return s.renderNegativeBool(stack, groupContextList.(float64) != 0.0)
	case string:
		return s.renderNegativeBool(stack, groupContextList.(string) != "")
	case map[string]interface{}:
		return s.renderNegativeBool(stack, len(groupContextList.(map[string]interface{})) != 0)
	case []interface{}:
		return s.renderNegativeBool(stack, len(groupContextList.([]interface{})) != 0)
	default:
		return ""
	}
}

func (s *NegativeNode) renderNegativeBool(stack ContextStack, condition bool) string {
	if condition {
		return ""
	}
	var sb strings.Builder
	for _, section := range s.Sns {
		sb.WriteString(section.Render(stack))
	}
	return sb.String()
}
