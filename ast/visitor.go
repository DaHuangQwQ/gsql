package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)

type PrintVisitor struct {
}

func (p PrintVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		fmt.Println(nil)
		return p
	}
	typ := reflect.TypeOf(node)
	val := reflect.ValueOf(node)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	fmt.Printf("val: %+v, typ %s \n", val.Interface(), typ.Name())

	return p
}
