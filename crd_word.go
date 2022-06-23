package crd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var wordMap map[string]int

func buildWordMap(rootPath string) {
	wordMap = make(map[string]int)

	inspectFunc := func(node ast.Node) bool {
		// ident, ok := node.(*ast.Ident)
		// if ok {
		// 	spew.Dump(ident.Name)
		// }
		// func
		funcDecl, ok := node.(*ast.FuncDecl)
		if ok {
			printIdent("F", funcDecl.Name)
			// func params
			for _, param := range funcDecl.Type.Params.List {
				printNames("fi", param.Names)
			}
			// return true
		}
		funcLit, ok := node.(*ast.FuncLit)
		if ok {
			for _, param := range funcLit.Type.Params.List {
				printNames("fi", param.Names)
			}
		}

		//type
		structSpec, ok := node.(*ast.TypeSpec)
		if ok {
			printIdent("T", structSpec.Name)

		}
		// st
		structType, ok := node.(*ast.StructType)
		if ok {
			for _, field := range structType.Fields.List {
				printNames("f", field.Names)
			}
		}
		// it
		interfaceType, ok := node.(*ast.InterfaceType)
		if ok {
			for _, method := range interfaceType.Methods.List {
				printNames("im", method.Names)
			}
		}
		// var
		valueSpec, ok := node.(*ast.ValueSpec)
		if ok {
			printNames("v", valueSpec.Names)
		}
		// assign
		assignStmt, ok := node.(*ast.AssignStmt)
		if ok && assignStmt.Tok == token.DEFINE {
			printAssignStmt(assignStmt)
		}
		// for
		// forStmt, ok := node.(*ast.ForStmt)
		// if ok {
		// 	assignStmt, ok := forStmt.Init.(*ast.AssignStmt)
		// 	if ok && assignStmt.Tok == token.DEFINE {
		// 		printAssignStmt(assignStmt)
		// 	}
		// }
		// range
		rangeStmt, ok := node.(*ast.RangeStmt)
		if ok {
			if rangeStmt.Tok == token.DEFINE {
				printIdent("va", rangeStmt.Key.(*ast.Ident))
				printIdent("va", rangeStmt.Value.(*ast.Ident))
			}
		}

		return true
	}

	// 遍历文件夹
	walkDir(rootPath, inspectFunc)
	saveWordMap()

	wordMap = nil
}

func walkDir(dirPath string, inspectFunc func(node ast.Node) bool) {
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, fileInfo := range fileInfos {
		fileName := dirPath + "/" + fileInfo.Name()
		if fileInfo.IsDir() {
			walkDir(fileName, inspectFunc)
		} else if filepath.Ext(fileName) == ".go" {
			inspectFile(fileName, inspectFunc)
		}
	}
}

func inspectFile(filename string, inspectFunc func(node ast.Node) bool) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)

	if err != nil {
		fmt.Printf("err = %s", err)
	}

	ast.Inspect(f, inspectFunc)
}

func printAssignStmt(stmt *ast.AssignStmt) {
	for _, expr := range stmt.Lhs {
		printIdent("va", expr.(*ast.Ident))
	}
}

func printNames(mark string, names []*ast.Ident) {
	for _, name := range names {
		printIdent(mark, name)
	}
}

func printIdent(mark string, ident *ast.Ident) {
	spew.Dump(mark + ":" + ident.Name)
	splitWord(ident.Name)
}

func splitWord(str string) {
	var beginIdx = 0
	var endIdx = 0
	for {
		var allUpperCase bool
		beginIdx = endIdx
		// find beginIdx
		for beginIdx < len(str) {
			var beginCh byte = str[beginIdx]
			if isWordBegin(beginCh) {
				allUpperCase = beginCh >= 'A' && beginCh <= 'Z'
				break
			}
			beginIdx++
			endIdx = beginIdx + 1

		}
		// no word find
		if beginIdx >= len(str) {
			break
		}

		// search end
		for endIdx < len(str) {
			var endCh byte = str[endIdx]
			if !isSameWord(endCh, &allUpperCase) {
				break
			}
			endIdx++
		}

		// save word
		var word = strings.ToLower(str[beginIdx:endIdx])
		if val, ok := wordMap[word]; ok {
			wordMap[word] = val + 1
		} else {
			wordMap[word] = 1
		}
	}
}
func isWordBegin(ch byte) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') {
		return true
	}
	return false
}
func isSameWord(ch byte, allUpperCase *bool) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
		*allUpperCase = false
		return true
	}
	if (*allUpperCase) && ch >= 'A' && ch <= 'Z' {
		return true
	}
	return false
}

func saveWordMap() {
	keyList := make([]string, len(wordMap))
	keyIdx := 0
	for key, _ := range wordMap {
		keyList[keyIdx] = key
		keyIdx++
	}
	sort.Strings(keyList)

	// save word map
	path, _ := os.Getwd()
	crdPath := path + "./.crd/"
	err := os.MkdirAll(crdPath, 0644)
	if err != nil {
		return
	}
	filename := crdPath + "word_list.txt"
	os.Remove(filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	for _, word := range keyList {
		file.WriteString(word + " " + fmt.Sprintf("%d", wordMap[word]) + "\n")
	}
}
