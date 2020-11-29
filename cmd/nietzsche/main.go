package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	nz "github.com/anwinged/nietzsche"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	showTokens := flag.Bool("tree", false, "Show template structure")

	flag.Parse()

	if flag.NArg() == 2 {
		templateFile := flag.Arg(0)
		dataFile := flag.Arg(1)
		processTemplateFromFile(templateFile, dataFile)
		os.Exit(0)
	}

	if flag.NArg() == 1 && *showTokens {
		templateFile := flag.Arg(0)
		printTemplateStructure(templateFile)
		os.Exit(0)
	}

	if flag.NArg() == 1 {
		template := captureInput()
		dataFile := flag.Arg(0)
		processTemplateWithDataFile(template, dataFile)
		os.Exit(0)
	}

	os.Exit(1)
}

func captureInput() string {
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		check(err)
	}
	return string(text)
}

func processTemplateFromFile(templateFile, dataFile string) {
	template, err := ioutil.ReadFile(templateFile)
	check(err)

	processTemplateWithDataFile(string(template), dataFile)
}

func processTemplateWithDataFile(template, dataFile string) {
	dataText, err := ioutil.ReadFile(dataFile)
	check(err)

	var data map[string]interface{}

	err = json.Unmarshal(dataText, &data)
	check(err)

	result, err := nz.Render(template, data)
	check(err)

	fmt.Println(result)
}

func printTemplateStructure(templateFile string) {
	template, err := ioutil.ReadFile(templateFile)
	check(err)

	tokens, err := nz.Tokenize(string(template))
	check(err)

	sections, err := nz.Compile(tokens)
	check(err)

	printNodesRecursive(sections, 0)
}

func printNodesRecursive(sections []nz.Node, level int) {
	for _, s := range sections {
		fmt.Printf("%s%-8s %s\n", strings.Repeat("    ", level), s.Type(), s.Desc())
		subs := s.Nodes()
		if len(subs) > 0 {
			printNodesRecursive(subs, 1)
		}
	}
}
