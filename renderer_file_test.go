package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	df "github.com/sergi/go-diff/diffmatchpatch"
)

func checkErrInTest(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func differ(diffs []df.Diff) bool {
	for _, d := range diffs {
		if d.Type != df.DiffEqual {
			return true
		}
	}
	return false
}

func testRenderFile(t *testing.T, testCase string) {
	templateFile := testCase + "/" + "template.mustache"
	dataFile := testCase + "/" + "data.json"
	resultFile := testCase + "/" + "result.txt"

	template, err := ioutil.ReadFile(templateFile)
	checkErrInTest(t, err)

	dataText, err := ioutil.ReadFile(dataFile)
	checkErrInTest(t, err)

	resultText, err := ioutil.ReadFile(resultFile)
	checkErrInTest(t, err)

	var data map[string]interface{}

	err = json.Unmarshal(dataText, &data)
	checkErrInTest(t, err)

	result, err := Render(string(template), data)
	checkErrInTest(t, err)

	dmp := df.New()
	diffs := dmp.DiffMain(string(resultText), result, false)

	if differ(diffs) {
		t.Logf("\n%s\n", data)
		t.Log("\n" + dmp.DiffPrettyText(diffs))
		t.Error("Render fails: result not match expected")
	}
}

func TestRenderFile(t *testing.T) {
	testRenderFile(t, "./share/test/simple")
}
