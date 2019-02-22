package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	df "github.com/sergi/go-diff/diffmatchpatch"
)

func assertEquals(t *testing.T, expected string, template string, params Context) {
	result, err := Render(template, params)
	if err != nil {
		t.Errorf("Render fails with error \"%s\"", err)
	}
	if expected != result {
		t.Errorf("Render fails. Expect \"%s\", got \"%s\"", expected, result)
	}
}

func assertError(t *testing.T, template string, params Context) {
	result, err := Render(template, params)
	if err == nil {
		t.Errorf("Render success, but error expected ('%s' => '%s')", template, result)
	}
}

func TestOneTag(t *testing.T) {
	assertEquals(
		t,
		"Bill Clinton",
		"{{name}} Clinton",
		Context{"name": "Bill"},
	)
	assertEquals(
		t,
		"Jack",
		"{{name}}",
		Context{"name": "Jack"},
	)
}

func TestTwoTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Great World!",
		"Hello, {{g}} {{w}}!",
		Context{"g": "Great", "w": "World"},
	)
}

func TestMissedTagValue(t *testing.T) {
	assertEquals(
		t,
		"Hello, Great !",
		"Hello, {{g}} {{w}}!",
		Context{"g": "Great"},
	)
}

func TestForgottenTag(t *testing.T) {
	assertError(
		t,
		"Hello, {{w",
		Context{"w": "World"},
	)
}

func TestBooleanGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike!",
		"Hello, {{#name}}Mike{{/name}}!",
		Context{
			"name": true,
		},
	)
}

func TestValueListGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, MikeMikeMike!",
		"Hello, {{#name}}Mike{{/name}}!",
		Context{
			"name": ValueList{1, 2, 3},
		},
	)
}

func TestContextGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike!",
		"Hello, {{#person}}{{name}}{{/person}}!",
		Context{
			"person": Context{"name": "Mike"},
		},
	)
}

func TestContextListGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike!",
		"Hello, {{#persons}}{{name}}{{/persons}}!",
		Context{
			"persons": ContextList{
				Context{"name": "Mike"},
			},
		},
	)
}

func TestContextStackInGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Moscow!",
		"Hello, {{#persons}}{{address}}{{/persons}}!",
		Context{
			"address": "Moscow",
			"persons": ContextList{
				Context{"name": "Mike"},
			},
		},
	)
}

func TestGroupTagList(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike, John, Kelly, !",
		"Hello, {{#persons}}{{name}}, {{/persons}}!",
		Context{
			"persons": ContextList{
				Context{"name": "Mike"},
				Context{"name": "John"},
				Context{"name": "Kelly"},
			},
		},
	)
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
	check(err)

	dataText, err := ioutil.ReadFile(dataFile)
	check(err)

	resultText, err := ioutil.ReadFile(resultFile)
	check(err)

	var data map[string]interface{}

	err = json.Unmarshal(dataText, &data)
	check(err)

	result, err := Render(string(template), Context(data))
	check(err)

	if err != nil {
		t.Errorf("Render fails: %s", err)
	}

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

// Benchmarking, func number - number of tokens

func BenchmarkRender1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render("{{x}}", Context{"x": "A"})
	}
}

func BenchmarkRender3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render("Hello, {{world}}!!", Context{"world": "World"})
	}
}

func BenchmarkRender10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render(
			`Hello, {{#persons}}{{fname}} {{lname}}, {{/persons}}!
			We are going to {{address}}.`,
			Context{
				"persons": ContextList{
					Context{"fname": "Mike", "lname": "Jackson"},
					Context{"fname": "John", "lname": "Rives"},
					Context{"fname": "Kelly", "lname": "Snow"},
				},
				"address": "Moscow",
			},
		)
	}
}
