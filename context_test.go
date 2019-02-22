package main

import "testing"

func makeContextStack() ContextStack {
	contexts := []map[string]interface{}{
		{"a": 1, "b": 2},
		{"c": 1, "d": 2},
		{"e": 1, "f": 2},
	}

	stack := ContextStack{}

	for _, c := range contexts {
		stack = stack.PushContext(c)
	}

	return stack
}

func TestFindValueOnFirstLevel(t *testing.T) {
	stack := makeContextStack()
	val := stack.FindValue("a")
	if val == nil {
		t.Errorf("Value expexted, nothing found")
	}
}

func TestFindValueOnLastLevel(t *testing.T) {
	stack := makeContextStack()
	val := stack.FindValue("f")
	if val == nil {
		t.Errorf("Value expexted, nothing found")
	}
}

func TestFindAbscentValue(t *testing.T) {
	stack := makeContextStack()
	val := stack.FindValue("zz")
	if val != nil {
		t.Errorf("Nothing expexted, value found")
	}
}
