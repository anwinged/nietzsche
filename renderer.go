package main

import "strings"

// RENDER

func Render(template string, context Context) (string, error) {
	sections, err := Compile(template)
	if err != nil {
		return "", err
	}
	stack := ContextStack{context}
	var sb strings.Builder
	for _, s := range sections {
		sb.WriteString(s.Render(stack))
	}
	return sb.String(), nil
}