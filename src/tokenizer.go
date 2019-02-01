package main

// Token type

type TokenType int

const (
	TextToken TokenType = iota
	ValueToken
)

// Token

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(tt TokenType, v string) *Token {
	t := new(Token)
	t.Type = tt
	t.Value = v
	return t
}

// Tokenizer

type Tokenizer struct {
	template string
	position int
	buffer   string
}

func NewTokenizer(template string) *Tokenizer {
	t := new(Tokenizer)
	t.template = template
	t.position = 0
	t.buffer = ""
	return t
}

func (t *Tokenizer) Next() (*Token, error) {
	if t.position >= len(t.template) {
		return nil, nil
	}

	var bracket int = 0
	var buffer string
	var pos int
	var char rune

	for pos, char = range t.template[t.position:] {
		if char == '{' {
			bracket += 1
			if bracket == 2 && buffer != "" {
				t.position += pos - 1
				return NewToken(TextToken, buffer), nil
			}
		} else if char == '}' {
			bracket -= 1
			if bracket == 0 && buffer != "" {
				t.position += pos + 1
				return NewToken(ValueToken, buffer), nil
			}
		} else {
			buffer += string(char)
		}
	}

	t.position += pos + 1
	return NewToken(TextToken, buffer), nil
}
