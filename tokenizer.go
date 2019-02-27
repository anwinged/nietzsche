package main

import (
	"errors"
	"strings"
)

// Token type

type TokenType int

const (
	TextToken TokenType = iota
	ValueToken
	OpenGroupToken
	OpenInvertedGroupToken
	CloseGroupToken
)

// Token

type Token struct {
	Type  TokenType
	Value string
}

func Tokenize(template string) ([]Token, error) {
	var tokens []Token
	var buffer strings.Builder
	var bracket int = 0

	for _, char := range template {
		if char == '{' {
			bracket += 1
			if bracket == 2 && buffer.Len() != 0 {
				tokens = append(tokens, Token{Type: TextToken, Value: buffer.String()})
				buffer.Reset()
			}
		} else if char == '}' {
			bracket -= 1
			if bracket == 0 && buffer.Len() != 0 {
				tokens = append(tokens, createTagToken(buffer.String()))
				buffer.Reset()
			}
		} else {
			buffer.WriteRune(char)
		}
	}

	if bracket != 0 {
		return nil, errors.New("Unexpected bracket")
	}

	if buffer.Len() != 0 {
		tokens = append(tokens, Token{Type: TextToken, Value: buffer.String()})
	}

	return tokens, nil
}

func createTagToken(val string) Token {
	trimmed := strings.TrimSpace(val)
	var head byte = trimmed[0]
	var tail string = strings.TrimSpace(trimmed[1:])
	switch head {
	case '#':
		return Token{Type: OpenGroupToken, Value: tail}
	case '^':
		return Token{Type: OpenInvertedGroupToken, Value: tail}
	case '/':
		return Token{Type: CloseGroupToken, Value: tail}
	}
	return Token{Type: ValueToken, Value: trimmed}
}
