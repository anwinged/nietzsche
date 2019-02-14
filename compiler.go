package main

func Compile(template string) ([]Section, error) {
	var sections []Section
	tokens, err := Tokenize(template)
	if err != nil {
		return nil, err
	}
	var group bool = false
	var buffer []Section
	for _, token := range tokens {
		switch token.Type {
		case TextToken:
			if group {
				buffer = append(buffer, NewTextSection(token.Value))
			} else {
				sections = append(sections, NewTextSection(token.Value))
			}
		case ValueToken:
			if group {
				buffer = append(buffer, NewValueSection(token.Value))
			} else {
				sections = append(sections, NewValueSection(token.Value))
			}
		case OpenSectionToken:
			group = true
		case CloseSectionToken:
			group = false
			sections = append(sections, NewGroupSection(token.Value, buffer))
			buffer = []Section{}
		}
	}
	return sections, nil
}
