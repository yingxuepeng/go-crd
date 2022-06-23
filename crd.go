package crd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func BuildDigest(rootPath string) {
	buildWordMap(rootPath)
}

func PrintAst(filePath string, astPath string) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)

	if err != nil {
		fmt.Printf("err = %s", err)
	}
	of, _ := os.OpenFile(astPath, os.O_CREATE|os.O_RDWR, 0)
	ast.Fprint(of, fset, f, nil)
}
