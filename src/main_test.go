package main

import "testing"

func TestTemplate(t *testing.T) {
	result := Template()
	if result != "World" {
		t.Errorf("Incorrect render")
	}
}
