package main

type Section interface {
	Render(context Context) string
}

// TEXT SECTION

type TextSection struct {
	Text string
}

func NewTextSection(text string) *TextSection {
	s := new(TextSection)
	s.Text = text
	return s
}

func (s *TextSection) Render(context Context) string {
	return s.Text
}

// VALUE SECTION

type ValueSection struct {
	Name string
}

func NewValueSection(name string) *ValueSection {
	s := new(ValueSection)
	s.Name = name
	return s
}

func (s *ValueSection) Render(context Context) string {
	return context[s.Name]
}
