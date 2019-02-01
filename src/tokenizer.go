package main

// Token type

type TokenType int

const (
	Text TokenType = iota
	Value
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
}

func NewTokenizer(template string) *Tokenizer {
	t := new(Tokenizer)
	t.template = template
	t.position = 0
	return t
}

func (t *Tokenizer) Next() (*Token, error) {
	return nil, nil
}
