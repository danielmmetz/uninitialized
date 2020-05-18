package main

import (
	"fmt"
	"go/ast"
	"sort"
	"strconv"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "uninitialized",
	Doc:      "check for uninitialized but required struct fields within composite literals",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.Preorder([]ast.Node{new(ast.CompositeLit)}, func(n ast.Node) {
		lit := n.(*ast.CompositeLit)
		required := requiredFields(lit)
		for _, element := range lit.Elts {
			kv, ok := element.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			ident, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			if !required[ident.Name] {
				continue
			}
			delete(required, ident.Name)
		}
		if len(required) > 0 {
			var typeName string
			if ident, ok := lit.Type.(*ast.Ident); ok {
				typeName = fmt.Sprintf("%s ", ident.Name)
			}
			pass.Reportf(n.Pos(), "%smissing required keys: %v", typeName, sortedKeys(required))
		}
	})

	return nil, nil
}

func requiredFields(lit *ast.CompositeLit) map[string]bool {
	ident, ok := lit.Type.(*ast.Ident)
	if !ok {
		return nil
	}
	if ident.Obj == nil {
		return nil
	}
	ts, ok := ident.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return nil
	}
	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		return nil
	}
	if st.Fields == nil {
		return nil
	}
	fields := make(map[string]bool)
	for _, field := range st.Fields.List {
		if field == nil {
			continue
		}
		if field.Tag == nil {
			continue
		}
		escapedTagValue := field.Tag.Value
		tags, err := structtag.Parse(unescapeTag(escapedTagValue))
		if err != nil {
			continue
		}
		tag, err := tags.Get("required")
		if err != nil {
			continue
		}
		value, _ := strconv.ParseBool(tag.Name)
		if !value {
			continue
		}
		if len(field.Names) == 0 {
			ident, ok := field.Type.(*ast.Ident)
			if !ok {
				continue
			}
			fields[ident.Name] = true
			continue
		}
		fields[field.Names[0].Name] = true
	}
	return fields
}

func sortedKeys(m map[string]bool) []string {
	r := []string{}
	for k := range m {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

func unescapeTag(escaped string) string {
	if len(escaped) < 2 {
		return escaped
	}
	return escaped[1 : len(escaped)-1]
}
