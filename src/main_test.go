package main

import "testing"

func assertEquals(t *testing.T, expected string, template string, params map[string]string) {
	result, err := Template(template, params)
	if err != nil {
		t.Errorf("Render fails with error \"%s\"", err)
	}
	if expected != result {
		t.Errorf("Render fails. Expect \"%s\", got \"%s\"", expected, result)
	}
}

func assertError(t *testing.T, template string, params map[string]string) {
	_, err := Template(template, params)
	if err == nil {
		t.Errorf("Render success, but error expected")
	}
}

func TestOneTagTemplate(t *testing.T) {
	assertEquals(
		t,
		"Bill Clinton",
		"{{name}} Clinton",
		map[string]string{"name": "Bill"},
	)
	assertEquals(
		t,
		"Jack",
		"{{name}}",
		map[string]string{"name": "Jack"},
	)
}

func TestTwoTagTemplate(t *testing.T) {
	assertEquals(
		t,
		"Hello, Great World!",
		"Hello, {{g}} {{w}}!",
		map[string]string{"g": "Great", "w": "World"},
	)
}

func TestMissedTagValueTemplate(t *testing.T) {
	assertEquals(
		t,
		"Hello, Great !",
		"Hello, {{g}} {{w}}!",
		map[string]string{"g": "Great"},
	)
}

func TestForgottenTag(t *testing.T) {
	assertError(
		t,
		"Hello, {{w",
		map[string]string{"w": "World"},
	)
}

func BenchmarkTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Template("Hello, {{world}}!", map[string]string{"world": "World"})
	}
}
