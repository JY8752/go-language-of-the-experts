package complexity

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func GetFuncNode(t *testing.T, code string) ast.Node {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			return fd
		}
	}
	t.Fatal("no function declear found")
	return nil
}

func TestComplexity(t *testing.T) {
	testcases := []struct {
		name       string
		code       string
		complexity int
	}{
		{
			name: "simple function",
			code: `
				package main

				import "fmt"

				func main() {
					fmt.Println("Hello World")
				}
			`,
			complexity: 1,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			//抽象構文木を取得する
			fn := GetFuncNode(t, testcase.code)

			//循環複雑度を取得する
			c := Count(fn)

			//アサーション
			if c != testcase.complexity {
				t.Errorf("expect %d, but %d\n", testcase.complexity, c)
			}
		})
	}
}
