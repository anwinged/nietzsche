package main

import "fmt"

func main() {
	result, err := Render(
		"{{h}}, {{w}}!",
		Context{"h": "Hello", "w": "World"},
	)

	if err != nil {
		fmt.Println("Error")
	} else {
		fmt.Println(result)
	}
}
