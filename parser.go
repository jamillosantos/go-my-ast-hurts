package myasthurts

import (
	"fmt"
	"go/ast"
)

func Parse(file *ast.File, definitions *Environment) {
	fmt.Println(file)
}
