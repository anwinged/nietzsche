package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	df "github.com/sergi/go-diff/diffmatchpatch"
)

func assertEquals(t *testing.T, expected string, template string, params map[string]interface{}) {
	result, err := Render(template, params)
	if err != nil {
		t.Errorf("Render fails with error \"%s\"", err)
	}
	if expected != result {
		t.Errorf("Render fails. Expect \"%s\", got \"%s\"", expected, result)
	}
}

func assertError(t *testing.T, template string, params map[string]interface{}) {
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
		map[string]interface{}{"name": "Bill"},
	)
	assertEquals(
		t,
		"Jack",
		"{{name}}",
		map[string]interface{}{"name": "Jack"},
	)
}

func TestTwoTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Great World!",
		"Hello, {{g}} {{w}}!",
		map[string]interface{}{"g": "Great", "w": "World"},
	)
}

func TestMissedTagValue(t *testing.T) {
	assertEquals(
		t,
		"Hello, Great !",
		"Hello, {{g}} {{w}}!",
		map[string]interface{}{"g": "Great"},
	)
}

func TestForgottenTag(t *testing.T) {
	assertError(
		t,
		"Hello, {{w",
		map[string]interface{}{"w": "World"},
	)
}

func TestBooleanGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike!",
		"Hello, {{#name}}Mike{{/name}}!",
		map[string]interface{}{
			"name": true,
		},
	)
}

func TestValueListGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, MikeMikeMike!",
		"Hello, {{#name}}Mike{{/name}}!",
		map[string]interface{}{
			"name": []interface{}{1, 2, 3},
		},
	)
}

func TestContextGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike!",
		"Hello, {{#person}}{{name}}{{/person}}!",
		map[string]interface{}{
			"person": map[string]interface{}{"name": "Mike"},
		},
	)
}

func TestContextListGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike!",
		"Hello, {{#persons}}{{name}}{{/persons}}!",
		map[string]interface{}{
			"persons": []interface{}{
				map[string]interface{}{"name": "Mike"},
			},
		},
	)
}

func TestContextStackInGroupTag(t *testing.T) {
	assertEquals(
		t,
		"Hello, Moscow!",
		"Hello, {{#persons}}{{address}}{{/persons}}!",
		map[string]interface{}{
			"address": "Moscow",
			"persons": []interface{}{
				map[string]interface{}{"name": "Mike"},
			},
		},
	)
}

func TestGroupTagList(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike, John, Kelly, !",
		"Hello, {{#persons}}{{name}}, {{/persons}}!",
		map[string]interface{}{
			"persons": []interface{}{
				map[string]interface{}{"name": "Mike"},
				map[string]interface{}{"name": "John"},
				map[string]interface{}{"name": "Kelly"},
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

	result, err := Render(string(template), data)
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
		Render("{{x}}", map[string]interface{}{"x": "A"})
	}
}

func BenchmarkRender3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render("Hello, {{world}}!!", map[string]interface{}{"world": "World"})
	}
}

func BenchmarkRender10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render(
			`Hello, {{#persons}}{{fname}} {{lname}}, {{/persons}}!
			We are going to {{address}}.`,
			map[string]interface{}{
				"persons": []interface{}{
					map[string]interface{}{"fname": "Mike", "lname": "Jackson"},
					map[string]interface{}{"fname": "John", "lname": "Rives"},
					map[string]interface{}{"fname": "Kelly", "lname": "Snow"},
				},
				"address": "Moscow",
			},
		)
	}
}
