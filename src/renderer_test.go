package main

import "testing"

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

func TestGroupTag(t *testing.T) {
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

func BenchmarkRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Render("Hello, {{world}}!", Context{"world": "World"})
	}
}
