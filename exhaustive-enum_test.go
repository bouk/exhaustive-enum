package main

import (
	"bytes"
	"testing"
)

func TestExamples(t *testing.T) {
	for _, test := range []struct{ dir, expected string }{
		{"./examples/example", `The following switch statements are missing constants:
	examples/example/example.go:17:2: Sunday
`},
		{"./examples/example2", ``},
		{"./examples/multiple-decl", `The following switch statements are missing constants:
	examples/multiple-decl/example.go:30:2: Sunday
`},
	} {
		var b bytes.Buffer
		c := Checker{
			Output: &b,
		}
		err := c.Run(test.dir)
		if err != nil {
			t.Error(err)
			continue
		}

		output := b.String()
		if output != test.expected {
			t.Errorf("%q != %q", output, test.expected)
		}
	}
}
