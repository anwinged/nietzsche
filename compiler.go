package main

import "errors"

type SectionLayer struct {
	name     string
	negative bool
	items    []Section
}

type SectionGroupStack []SectionLayer

func NewStack() *SectionGroupStack {
	return new(SectionGroupStack)
}

func (s *SectionGroupStack) NewLayer(name string, negative bool) {
	layer := SectionLayer{name, negative, []Section{}}
	(*s) = append(*s, layer)
}

func (s *SectionGroupStack) CloseLayer(name string) ([]Section, bool, error) {
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

func (s *SectionGroupStack) AddSection(section Section) {
	length := len(*s)
	lastLayer := (*s)[length-1]
	lastLayer.items = append(lastLayer.items, section)
	(*s)[length-1] = lastLayer
}

func Compile(tokens []Token) ([]Section, error) {
	stack := NewStack()
	stack.NewLayer("__root__", false)
	for _, token := range tokens {
		switch token.Type {
		case TextToken:
			stack.AddSection(NewTextSection(token.Value))
		case ValueToken:
			stack.AddSection(NewValueSection(token.Value))
		case OpenSectionToken:
			stack.NewLayer(token.Value, false)
		case InvertedSectionToken:
			stack.NewLayer(token.Value, true)
		case CloseSectionToken:
			sections, neg, err := stack.CloseLayer(token.Value)
			if err != nil {
				return nil, err
			}
			if neg {
				stack.AddSection(NewNegativeSection(token.Value, sections))
			} else {
				stack.AddSection(NewGroupSection(token.Value, sections))
			}

		}
	}
	sections, _, err := stack.CloseLayer("__root__")
	if err != nil {
		return nil, err
	}

	return sections, nil
}
