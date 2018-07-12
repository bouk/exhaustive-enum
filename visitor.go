package main

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/loader"
)

type visitor struct {
	// pkg is the package that we are currently inspecting
	pkg *loader.PackageInfo

	// program is the whole program containing all the type info
	program *loader.Program

	// violations is a slice containing all the found violations
	violations []*Violation
}

// isExhaustiveEnum checks whether the given type is an exhaustive enum.
// Exhaustive enums have a comment `//exhaustive-enum` applied to them
// TODO(bouk): add support for supplying your own types to be accepted as an exhaustive enum
// TODO(bouk): add caching
func (v *visitor) isExhaustiveEnum(typ *types.Named) bool {
	obj := typ.Obj()
	pkg := obj.Pkg()

	// Search through ASTs to find comments associated with type
	for _, f := range v.program.AllPackages[pkg].Files {
		for _, d := range f.Decls {
			// Only process type statements
			decl, ok := d.(*ast.GenDecl)
			if !(ok && decl.Tok == token.TYPE) {
				continue
			}

			// Check if this is the type statement we're looking for
			for _, spec := range decl.Specs {
				spec := spec.(*ast.TypeSpec)
				if spec.Name.Name != obj.Name() {
					continue
				}

				cmap := ast.NewCommentMap(v.program.Fset, f, f.Comments)

				var comments []*ast.CommentGroup
				// If there's multiple specs we only want the comments that are specific to the one we are looking at
				if len(decl.Specs) == 1 {
					comments = cmap[decl]
				} else {
					comments = cmap[spec]
				}

				for _, comment := range comments {
					for _, s := range comment.List {
						if s.Text == "//exhaustive-enum" {
							return true
						}
					}
				}

				return false
			}
		}
	}

	return false
}

func (v *visitor) missingConstants(switc *ast.SwitchStmt) (missing []*types.Const) {
	t, ok := v.pkg.Info.TypeOf(switc.Tag).(*types.Named)
	if !ok || !v.isExhaustiveEnum(t) {
		return nil
	}

	// Collect all the constants that are defined for the type
	var constants []*types.Const
	scope := t.Obj().Pkg().Scope()
	for _, name := range scope.Names() {
		c, ok := scope.Lookup(name).(*types.Const)
		if !(ok && c.Type() == t) {
			continue
		}

		constants = append(constants, c)
	}

	// Collect all the constants that are used
	usedConstants := make(map[constant.Value]struct{})
	for _, s := range switc.Body.List {
		cas, ok := s.(*ast.CaseClause)
		// A default: case has an empty List
		if !ok || len(cas.List) == 0 {
			// In this case the switch is always exhaustive
			return nil
		}

		for _, e := range cas.List {
			e = getInnerFromParen(e)
			switch e := e.(type) {
			case *ast.SelectorExpr:
				v := v.pkg.Types[e].Value
				if v == nil {
					// This is not a constant value, so it can be anything, so abort
					return nil
				}
				usedConstants[v] = struct{}{}
			case *ast.Ident:
				v := v.pkg.Types[e].Value
				if v == nil {
					// This is not a constant value, so it can be anything, so abort
					return nil
				}
				usedConstants[v] = struct{}{}
			case *ast.BasicLit:
				v := constant.MakeFromLiteral(e.Value, e.Kind, 0)
				usedConstants[v] = struct{}{}
			default:
				// This is not an enum switch so abort
				return nil
			}
		}
	}

	// Compare used and defined constants
	for _, c := range constants {
		val := c.Val()
		if _, ok := usedConstants[val]; !ok {
			missing = append(missing, c)
		}
	}

	return missing
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	s, ok := node.(*ast.SwitchStmt)
	if !ok {
		return v
	}

	if s.Init != nil {
		ast.Walk(v, s.Init)
	}

	if s.Tag != nil {
		missing := v.missingConstants(s)
		if len(missing) > 0 {
			v.violations = append(v.violations, &Violation{
				Position: v.program.Fset.Position(s.Pos()),
				Missing:  missing,
			})
		}

		ast.Walk(v, s.Tag)
	}

	// TODO(bouk): We need to keep track of which constants were not used, and propegate them through into default: case
	ast.Walk(v, s.Body)

	return nil
}
