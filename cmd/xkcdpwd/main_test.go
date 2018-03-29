package main

import "testing"

func Test_checkSeparatro(t *testing.T) {
	// valid
	valid := []string{"", " ", "-", ".", "_"}
	for _, sep := range valid {
		if !checkSeparator(sep) {
			t.Errorf("valid separator '%s' failed check", sep)
		}
	}
	invalid := []string{"ajfjke;ja;", "     ", "----", "....", "____", "a", "z", "A", "Z", "ü"}
	for _, sep := range invalid {
		if checkSeparator(sep) {
			t.Errorf("invalid separator '%s' passed check", sep)
		}
	}
}
