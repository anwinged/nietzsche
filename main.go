package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()
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

	result, err := Render(string(template), Context(data))
	check(err)

	fmt.Println(result)
}
