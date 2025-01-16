package time_func

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func astRewriteFile(path string, fn func(path string, n ast.Node) bool) (err error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("解析文件 %s 失败: %v", path, err)
	}

	ast.Inspect(file, func(node ast.Node) bool {
		return fn(path, node)
	})

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("打开文件 %s 用于写入失败: %v", path, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	err = format.Node(&buf, fset, file)
	if err != nil {
		return fmt.Errorf("格式化语法树 %s 失败: %v", path, err)
	}
	_, _ = io.Copy(f, &buf)
	return
}

// generateImport 生成 import
func generateImport(s string) *ast.ImportSpec {
	return &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s\"", s),
		},
		EndPos: token.NoPos,
	}
}

// generateImportSet 序列化import 包
func generateImportSet(impDecl *ast.GenDecl) map[string]bool {
	importSet := make(map[string]bool)
	for _, spec := range impDecl.Specs {
		if spec, ok := spec.(*ast.ImportSpec); ok {
			path := strings.Trim(spec.Path.Value, "\"")
			importSet[path] = true
		}
	}
	return importSet
}

func hasFunc(file *ast.File) bool {
	for _, decl := range file.Decls {
		if _, ok := decl.(*ast.FuncDecl); ok {
			return true
		}
	}
	return false
}

func hasImport(file *ast.File) (importDecl *ast.GenDecl, ok bool) {
	//有头部导入文件
	for _, decl := range file.Decls {
		impDecl, ok := decl.(*ast.GenDecl)
		if ok && impDecl.Tok == token.IMPORT {
			return impDecl, true
		}
	}
	return
}

func getAllDirs(dir ...string) (files []string, err error) {
	files = make([]string, 0)
	for _, s := range dir {
		err = filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".go") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return
		}
	}
	return
}
