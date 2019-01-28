package main

import "errors"
import "fmt"

func Template(template string, params map[string]string) (string, error) {
	var result = ""
	var bracket = 0
	var name = ""
	for _, char := range template {
		if char == '{' {
			if bracket == 0 {
				bracket = 1
			} else if bracket == 1 {
				bracket = 2
			} else {
				return "", errors.New("Unexpected {")
			}
		} else if char == '}' {
			if bracket == 2 {
				bracket = 1
			} else if bracket == 1 {
				bracket = 0
			} else {
				return "", errors.New("Unexpected }")
			}
		} else {
			if bracket == 2 {
				name += string(char)
			} else {
				if name != "" {
					result += params[name]
					name = ""
				}
				result += string(char)
			}
		}
	}
	if name != "" {
		if bracket > 0 {
			return "", errors.New("Unexpected end of placeholder")
		}
		result += params[name]
	}

	return result, nil
}

func main() {
	fmt.Println("hello world")
}
