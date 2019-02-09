package parser

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	i, err := os.Open("__stubs__/input.eml")
	if err != nil {
		t.Fatal("Cannot read file: ", err)
	}

	o, err := ioutil.ReadFile("__stubs__/output.md")
	if err != nil {
		t.Fatal("Cannot read file: ", err)
	}

	result, err := Parse(i)
	if err != nil {
		t.Fatal("Cannot parse email: ", err)
	}

	if result != string(o) {
		t.Fatal("Expected: ", string(o), ", but found: ", result)
	}
}
