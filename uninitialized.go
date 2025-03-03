package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"maps"
	"sort"
	"strconv"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:      "uninitialized",
	Doc:       "check for uninitialized but required struct fields within composite literals",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{new(hasRequiredFields)},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.Preorder([]ast.Node{new(ast.TypeSpec)}, func(n ast.Node) {
		ts := n.(*ast.TypeSpec)
		fields := requiredFieldsFromType(ts)
		if len(fields) > 0 {
			obj, ok := pass.TypesInfo.Defs[ts.Name]
			if !ok {
				return
			}
			pass.ExportObjectFact(obj, &hasRequiredFields{requiredFields: fields})
		}
	})
	inspect.Preorder([]ast.Node{new(ast.CompositeLit)}, func(n ast.Node) {
		lit := n.(*ast.CompositeLit)
		required := requiredFieldsForCompositLit(pass, lit)
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
			switch t := lit.Type.(type) {
			case *ast.Ident:
				typeName = fmt.Sprintf("%s ", t.Name)
			case *ast.SelectorExpr:
				typeName = fmt.Sprintf("%s ", t.Sel.Name)
			}
			pass.Reportf(n.Pos(), "%smissing required keys: %v", typeName, sortedKeys(required))
		}
	})

	return nil, nil
}

func requiredFieldsForCompositLit(pass *analysis.Pass, lit *ast.CompositeLit) map[string]bool {
	var obj types.Object
	var ok bool
	switch t := lit.Type.(type) {
	case *ast.Ident:
		obj, ok = pass.TypesInfo.Uses[t]
		if !ok {
			return nil
		}
		if t.Obj == nil {
			break
		}
		if ts, ok := t.Obj.Decl.(*ast.TypeSpec); ok {
			return requiredFieldsFromType(ts)
		}
	case *ast.SelectorExpr:
		obj, ok = pass.TypesInfo.Uses[t.Sel]
		if !ok {
			return nil
		}
	default:
		return nil
	}
	var f hasRequiredFields
	if ok := pass.ImportObjectFact(obj, &f); ok {
		return maps.Clone(f.requiredFields)
	}
	return nil
}

type hasRequiredFields struct {
	requiredFields map[string]bool
}

func (f *hasRequiredFields) AFact() {}

func (f *hasRequiredFields) String() string {
	keys := make([]string, 0, len(f.requiredFields))
	for k := range f.requiredFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return fmt.Sprint(keys)
}

func requiredFieldsFromType(lit *ast.TypeSpec) map[string]bool {
	st, ok := lit.Type.(*ast.StructType)
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
