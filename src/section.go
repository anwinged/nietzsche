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

// TAG SECTION

type TagSection struct {
	Name string
}

func NewTagSection(name string) *TagSection {
	s := new(TagSection)
	s.Name = name
	return s
}

func (s *TagSection) Render(context Context) string {
	return context[s.Name]
}
