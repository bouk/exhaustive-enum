// Exhaustive-enum is an exhaustive enum checker for Go.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"strings"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

type Violation struct {
	// Position is the position of the switch construct
	Position token.Position

	// Missing is a slice containing the missing constants
	Missing []*types.Const
}

func (v *Violation) String() string {
	var b bytes.Buffer
	dir, _ := os.Getwd()
	fmt.Fprintf(&b, "\t%s: ", formatPosition(dir, &v.Position))

	for i, c := range v.Missing {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(c.Name())
	}
	return b.String()
}

// formatPosition formats a token.Position with the given directory stripped
func formatPosition(dir string, pos *token.Position) string {
	s := pos.String()
	if strings.HasPrefix(s, dir) {
		// +1 to strip off path seperator
		s = s[len(dir)+1:]
	}

	return s
}

type Checker struct {
	// Where the output is logged
	Output io.Writer
}

func (c *Checker) Run(args ...string) error {
	var conf loader.Config
	conf.ParserMode = parser.ParseComments

	_, err := conf.FromArgs(gotool.ImportPaths(args), false)
	if err != nil {
		return err
	}

	program, err := conf.Load()
	if err != nil {
		return err
	}

	// TODO: parallelize
	var violations []*Violation
	for _, pkg := range program.InitialPackages() {
		v := visitor{pkg, program, nil}
		for _, f := range pkg.Files {
			ast.Walk(&v, f)
		}
		violations = append(violations, v.violations...)
	}

	if len(violations) > 0 {
		fmt.Fprintln(c.Output, "The following switch statements are missing constants:")
		for _, violation := range violations {
			fmt.Fprintln(c.Output, violation)
		}
	}

	return nil
}

func main() {
	c := Checker{
		Output: os.Stderr,
	}
	flag.Parse()
	if err := c.Run(flag.Args()...); err != nil {
		fmt.Println(err)
	}
}
