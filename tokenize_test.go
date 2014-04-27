package tokenize

import (
	"reflect"
	"strings"
	"testing"
)

var linetab = []struct {
	line   string
	tokens []string
}{
	{"word", []string{"word"}},
	{"a b", []string{"a", "b"}},
	{`"a" "b"`, []string{"a", "b"}},
	{`"a""" "b"`, []string{`a"`, "b"}},
}

func TestLine(t *testing.T) {
	for i, e := range linetab {
		v := Line(e.line)
		if !reflect.DeepEqual(v, e.tokens) {
			t.Errorf("test %d failed, expected %#v, saw %#v\n", i+1, e.tokens, v)
		}
	}
}

func TestReader(t *testing.T) {
	file := ""
	exp := make([][]string, 0)
	for i, e := range linetab {
		s := ""
		if i%3 == 0 {
			s = " " + s
		}
		s += e.line
		if i%4 == 0 {
			s += " "
		}
		if i%2 == 0 {
			s += "\r"
		}
		s += "\n"
		if i%2 == 0 {
			s += "\n" // empty line
		}
		if i%3 == 0 {
			s += "# comment\n"
		}
		if i%4 == 0 {
			s += " # comment not starting at start of line\n"
		}
		file += s
		exp = append(exp, e.tokens)
	}
	r := strings.NewReader(file)
	got, err := reader(r)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, exp) {
		t.Fatalf("expected %#v, saw %#v\n", exp, got)
	}
}
