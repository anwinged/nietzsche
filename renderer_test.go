package nietzsche

import (
	"io/ioutil"
	"testing"

	yaml "gopkg.in/yaml.v2"
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

func TestPositiveBooleanGroupTag(t *testing.T) {
	cases := []interface{}{true, 1, float64(1.0), "yes"}
	for _, value := range cases {
		assertEquals(
			t,
			"Hello, Mike!",
			"Hello, {{#name}}Mike{{/name}}!",
			map[string]interface{}{
				"name": value,
			},
		)
	}
}

func TestNegativeBooleanGroupTag(t *testing.T) {
	cases := []interface{}{false, 0, float64(0.0), ""}
	for _, value := range cases {
		assertEquals(
			t,
			"Hello, Mike!",
			"Hello, {{^name}}Mike{{/name}}!",
			map[string]interface{}{
				"name": value,
			},
		)
	}
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

func TestDeepGroup(t *testing.T) {
	assertEquals(
		t,
		"Hello, Mike!",
		"Hello, {{#humans}}{{#persons}}{{name}}{{/persons}}{{/humans}}!",
		map[string]interface{}{
			"humans": map[string]interface{}{
				"persons": map[string]interface{}{
					"name": "Mike",
				},
			},
		},
	)
}

func TestInterpolation(t *testing.T) {
	testSpecSuite(t, "interpolation")
}

func TestSections(t *testing.T) {
	testSpecSuite(t, "sections")
}

func testSpecSuite(t *testing.T, suiteName string) {
	spec, err := ioutil.ReadFile("spec/specs/" + suiteName + ".yml")
	checkErrInTest(t, err)

	var data map[string]interface{}

	err = yaml.Unmarshal([]byte(spec), &data)
	checkErrInTest(t, err)

	suites := data["tests"].([]interface{})

	for _, test := range suites {
		s := test.(map[interface{}]interface{})
		name := s["name"]
		template := s["template"].(string)
		expected := s["expected"].(string)

		data := make(map[string]interface{})
		for k, v := range s["data"].(map[interface{}]interface{}) {
			data[k.(string)] = v
		}

		result, err := Render(template, data)

		checkErrInTest(t, err)

		if expected != result {
			t.Logf("%s\n", name)
			t.Logf("Render fails. Expect \"%s\", got \"%s\"", expected, result)
		}
	}
}

func checkErrInTest(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error: %s", err)
	}
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
