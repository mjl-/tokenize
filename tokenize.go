// Package tokenize helps parsing lines and entire files that contain (quoted) strings.
//
// Example of parsing a single line with Line():
//	these are words "but this is a quoted string"
// Results in:
//	[`these`, `are`, `words`, `but this is a quoted string`]
//
// Escaping:
//	"a single double quote ("") can be had by putting two subsequent double quotes in a qouted string"
// Results in:
//	[`a single double quote (") can be had by putting two subsequent double quotes in a quoted string`]
//
package tokenize

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// s starts with doublequote, read until the ending doublequote.
// double chars escape a single one.
// no error is raised if the ending quote is missing.
func parsequotedstring(s string) (string, string) {
	r := make([]rune, 0)
	prevquote := false
	s = s[1:]
	for i, c := range s {
		switch c {
		case '"':
			if prevquote {
				r = append(r, c)
				prevquote = false
			} else {
				prevquote = true
			}
		case ' ', '\t':
			if prevquote {
				return string(r), s[i:]
			}
			r = append(r, c)
		default:
			r = append(r, c)
		}
	}
	return string(r), ""
}

// Line splits line into space/tab separate tokens, parsing double quoted strings.
// Escape an explicit double quote with two double quotes.
func Line(line string) []string {
	t := make([]string, 0)
	var tok string
	for line != "" {
		if strings.HasPrefix(line, `"`) {
			tok, line = parsequotedstring(line)
		} else {
			o := strings.IndexAny(line, " \t")
			if o >= 0 {
				tok, line = line[:o], line[o:]
			} else {
				tok, line = line, ""
			}
		}
		line = strings.TrimLeft(line, " \t")
		t = append(t, tok)
	}
	return t
}

func reader(rd io.Reader) ([][]string, error) {
	b := bufio.NewScanner(rd)
	r := make([][]string, 0)
	for b.Scan() {
		s := b.Text()
		s = strings.Trim(s, " \t\r")
		if strings.HasPrefix(s, "#") || s == "" {
			continue
		}
		r = append(r, Line(s))
	}
	return r, b.Err()
}

// File splits all lines in path into tokens, by calling Line().
// Leading and trailing whitespace is stripped from each line.
// Empty lines and lines starting with a # are ignored.
//
// Example:
//	# this is a comment
//	# empty lines are ignored
//	this line is non-empty and parsed as a list of words
//	"this is a single string"
func File(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return reader(f)
}
