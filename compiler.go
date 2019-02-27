package main

import "errors"

type NodeLayer struct {
	name     string
	negative bool
	items    []Node
}

type NodeGroupStack []NodeLayer

func NewStack() *NodeGroupStack {
	return new(NodeGroupStack)
}

func (s *NodeGroupStack) NewLayer(name string, negative bool) {
	layer := NodeLayer{name, negative, []Node{}}
	(*s) = append(*s, layer)
}

func (s *NodeGroupStack) CloseLayer(name string) ([]Node, bool, error) {
	length := len(*s)
	lastLayer := (*s)[length-1]
	if lastLayer.name != name {
		return nil, false, errors.New("Unexpected close token " + name)
	}
	result := lastLayer.items
	neg := lastLayer.negative
	(*s) = (*s)[:length-1]
	return result, neg, nil
}

func (s *NodeGroupStack) AddNode(section Node) {
	length := len(*s)
	lastLayer := (*s)[length-1]
	lastLayer.items = append(lastLayer.items, section)
	(*s)[length-1] = lastLayer
}

func Compile(tokens []Token) ([]Node, error) {
	stack := NewStack()
	stack.NewLayer("__root__", false)
	for _, token := range tokens {
		switch token.Type {
		case TextToken:
			stack.AddNode(NewTextNode(token.Value))
		case ValueToken:
			stack.AddNode(NewValueNode(token.Value))
		case OpenNodeToken:
			stack.NewLayer(token.Value, false)
		case InvertedNodeToken:
			stack.NewLayer(token.Value, true)
		case CloseNodeToken:
			sections, neg, err := stack.CloseLayer(token.Value)
			if err != nil {
				return nil, err
			}
			if neg {
				stack.AddNode(NewNegativeNode(token.Value, sections))
			} else {
				stack.AddNode(NewGroupNode(token.Value, sections))
			}

		}
	}
	sections, _, err := stack.CloseLayer("__root__")
	if err != nil {
		return nil, err
	}

	return sections, nil
}
