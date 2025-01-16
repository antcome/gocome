package time_func

import (
	"embed"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

//go:embed demo/demo.go
var demoFS embed.FS
var demoFile string

func init() {
	data, _ := demoFS.ReadFile("demo/demo.go")
	demoFile = string(data)
}

const (
	demoFuncCalTime            = "calTime"
	demoFuncCalTimeFirstAssign = "startTime"
)

var (
	demoImports          *ast.GenDecl
	demoFuncCallTimeBody []ast.Stmt
	demoImportSet        = make(map[string]bool)
)

func init() {
	var file, _ = parser.ParseFile(token.NewFileSet(), "", strings.NewReader(demoFile), 0)
	for _, decl := range file.Decls {
		if impDecl, ok := decl.(*ast.GenDecl); ok && impDecl.Tok == token.IMPORT {
			demoImports = impDecl
		}
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == demoFuncCalTime {
				demoFuncCallTimeBody = funcDecl.Body.List
			}
		}
	}
	demoImportSet = generateImportSet(demoImports)
}

func importDemoFile(n ast.Node) {
	file, ok := n.(*ast.File)
	if !ok {
		return
	}

	if !hasFunc(file) {
		return
	}

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if ok {
			addDemoFuncCalTime(funcDecl)
		}
	}

	if importDecl, ok := hasImport(file); ok {
		importSet := generateImportSet(importDecl)
		for s := range demoImportSet {
			if !importSet[s] {
				importDecl.Specs = append(importDecl.Specs, generateImport(s))
			}
		}
	} else {
		file.Decls = append([]ast.Decl{demoImports}, file.Decls...)
	}
}

func addDemoFuncCalTime(fn *ast.FuncDecl) {
	if len(fn.Body.List) == 0 {
		return
	}
	block, ok := fn.Body.List[0].(*ast.BlockStmt)
	if ok {
		assign2, ok := block.List[0].(*ast.AssignStmt)
		if ok {
			if expr, ok := assign2.Lhs[0].(*ast.Ident); ok {
				if expr.Name == demoFuncCalTimeFirstAssign {
					return
				}
			}
		}
	}
	fn.Body.List = append(demoFuncCallTimeBody, fn.Body.List...)
}
