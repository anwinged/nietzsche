package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	showTokens := flag.Bool("tokens", false, "Show tokens of template")

	flag.Parse()

	if flag.NArg() == 1 && *showTokens {
		templateFile := flag.Arg(0)
		printSections(templateFile)
		os.Exit(0)
	}

	if flag.NArg() < 2 {
		os.Exit(1)
	}

	templateFile := flag.Arg(0)
	dataFile := flag.Arg(1)

	template, err := ioutil.ReadFile(templateFile)
	check(err)

	dataText, err := ioutil.ReadFile(dataFile)
	check(err)

	var data map[string]interface{}

	err = json.Unmarshal(dataText, &data)
	check(err)

	result, err := Render(string(template), data)
	check(err)

	fmt.Println(result)
}

func printSections(templateFile string) {
	template, err := ioutil.ReadFile(templateFile)
	check(err)

	tokens, err := Tokenize(string(template))
	check(err)

	sections, err := Compile(tokens)
	check(err)

	printSectionsRecursive(sections, 0)
}

func printSectionsRecursive(sections []Section, level int) {
	for _, s := range sections {
		fmt.Printf("%s%-8s %s\n", strings.Repeat("    ", level), s.Type(), s.Desc())
		subs := s.Sections()
		if len(subs) > 0 {
			printSectionsRecursive(subs, 1)
		}
	}
}
