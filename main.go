package main

import (
	"fmt"
	"go/token"
	"go/parser"
	"go/ast"
	"os"
	"github.com/davecgh/go-spew/spew"
	"go/printer"
)

type SourceFile struct {
	assignments []*ast.AssignStmt // :=
	values []*ast.ValueSpec // consts, =
}


func main() {
	mySource := SourceFile{}
	// make some arbitrary integer with a bin expr value (e.g.: 1+5,7*900,x<<1, etc..)
	var five = 1 + 2 + 1 + 1
	six := 6
	fmt.Printf("Five equals %d\n", five)
	fmt.Printf("Six equals %d\n", six)
	var helloStr = "Hello World!"
	fmt.Println(helloStr)
	// new fset object for parsing (lotta constants)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "main.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// printing the AST tree
	ast.Fprint(os.Stdout, fset, node, nil)

	// the first parameter is the AST
	// the second parameter specifies a function we will execute PER node
	// visited during the breadth first search of the AST
	ast.Inspect(node, func (n ast.Node) bool{
		/*

		OLD CODE

		var t, s string
		switch x := n.(type) {
		case *ast.BasicLit:
			// if it's a literal value (int, str, etc..)
			t = "Literal"
			s = x.Value
		case *ast.Ident:
			// if it's an Identity (?)
			t = "Ident"
			s = x.Name
		}
		if s != "" {
			fmt.Printf("%s:\t%s\t%s\n", fset.Position(n.Pos()), t, s)
		}
		*/

		// print out assignments
		// x := 1
		// t := myFunc()
		aStmt, ok := n.(*ast.AssignStmt)
		if ok {
			fmt.Printf("Found an assign token: %s\n", aStmt.Tok)
			fmt.Printf("lhs: %s\n", aStmt.Lhs)
			fmt.Printf("rhs: %s\n", aStmt.Rhs)
			mySource.assignments = append(mySource.assignments, aStmt)
			return true
		}

		// print out declarations
		// var y int = 5
		defs, ok := n.(*ast.ValueSpec)
		if ok {
			fmt.Printf("Found a decl name: %s\n", defs.Names)
			fmt.Printf("Type : %s\n", defs.Type)
			fmt.Printf("Values : %s\n", defs.Values)
			mySource.values = append(mySource.values, defs)
			return true
		}

		// print out all literals "asdf", 5, 0xDEADBEEF, etc...
		//literals, ok := n.(*ast.BasicLit)
		//if ok {
		//	fmt.Printf("Found a basic lit: %s\n", literals.Value)
		//	//fmt.Printf("token CONST_KIND : %s\n", literals.Kind)
		//	return true
		//}
		return true
	})
	fmt.Printf("\n\nOkay, now we've traversed the AST, let's identify the function we want to change\n")
	fmt.Printf("Let's walk through each of our variables and find the one we want")
	for i := range mySource.values{
		if len(mySource.values[i].Values) > 1 {
			fmt.Printf("Found a variable with multiple values?? Inspect variable number %d", i)
		}
		binExpr, ok := mySource.values[i].Values[0].(*ast.BinaryExpr) // check and see if the value is a binary expression
		name := mySource.values[i].Names[0]
		if ok {
			if name.Name == "five" { // hard coded for now, let's just change one
				fmt.Printf("\n\nFound it!\n")
				spew.Dump(binExpr)
				testing := binExpr.Y.(*ast.BasicLit)
				// changing a value
				testing.Value = "500"
				fmt.Printf("\n\nChanged a value in five, printing source code from AST!\n\n\n\n\n")
			}

		}
	}
	printer.Fprint(os.Stdout, fset, node)
}