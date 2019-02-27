package main

import "strings"

// RENDER

func Render(template string, context map[string]interface{}) (string, error) {
	tokens, err := Tokenize(template)
	if err != nil {
		return "", err
	}

	sections, err := Compile(tokens)
	if err != nil {
		return "", err
	}

	return RenderAST(sections, context)
}

func RenderAST(sections []Node, context map[string]interface{}) (string, error) {
	stack := ContextStack{context}
	var sb strings.Builder
	for _, s := range sections {
		sb.WriteString(s.Render(stack))
	}
	return sb.String(), nil
}
