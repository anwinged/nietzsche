package main

import "errors"

// Token type

type TokenType int

const (
	TextToken TokenType = iota
	ValueToken
	OpenSectionToken
	CloseSectionToken
)

// Token

type Token struct {
	Type  TokenType
	Value string
}

func Tokenize(template string) ([]Token, error) {
	var tokens []Token
	var buffer string
	var bracket int = 0

	for _, char := range template {
		if char == '{' {
			bracket += 1
			if bracket == 2 && buffer != "" {
				tokens = append(tokens, Token{Type: TextToken, Value: buffer})
				buffer = ""
			}
		} else if char == '}' {
			bracket -= 1
			if bracket == 0 && buffer != "" {
				tokens = append(tokens, createTagToken(buffer))
				buffer = ""
			}
		} else {
			buffer += string(char)
		}
	}

	if bracket != 0 {
		return nil, errors.New("Unexpected bracket")
	}

	if buffer != "" {
		tokens = append(tokens, Token{Type: TextToken, Value: buffer})
	}

	return tokens, nil
}

func createTagToken(val string) Token {
	var head byte = val[0]
	var tail string = val[1:]
	switch head {
	case '#':
		return Token{Type: OpenSectionToken, Value: tail}
	case '/':
		return Token{Type: CloseSectionToken, Value: tail}
	}
	return Token{Type: ValueToken, Value: val}
}
