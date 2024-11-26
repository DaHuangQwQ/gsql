package main

import (
	"fmt"
	"go/ast"
)

type SingleFileEntryVisitor struct {
	file *fileVisitor
}

func (s *SingleFileEntryVisitor) Get() File {
	if s.file != nil {
		return s.file.Get()
	}
	return File{}
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.File)
	if ok {
		s.file = &fileVisitor{
			pkg: n.Name.String(),
		}
		return s.file
	}
	return s
}

type fileVisitor struct {
	types   []*typeVisitor
	imports []string
	pkg     string
}

func (f *fileVisitor) Get() File {
	types := make([]Type, 0, len(f.types))
	for _, t := range f.types {
		types = append(types, t.Get())
	}
	return File{
		Package: f.pkg,
		Imports: f.imports,
		Types:   types,
	}
}

func (f *fileVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		res := &typeVisitor{
			name:   n.Name.String(),
			fields: make([]Field, 0, 0),
		}
		f.types = append(f.types, res)
		return res
	case *ast.ImportSpec:
		path := n.Path.Value
		if n.Name != nil && n.Name.String() != "" {
			path = n.Name.String() + " " + path
		}
		f.imports = append(f.imports, path)
	}
	return f
}

type File struct {
	Types []Type
	// 就是直接用户写出来的那种样子
	Imports []string

	Package string
}

type typeVisitor struct {
	name   string
	fields []Field
}

func (t *typeVisitor) Get() Type {
	return Type{
		Name:   t.name,
		Fields: t.fields,
	}
}

func (t *typeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fd, ok := node.(*ast.Field)
	if ok {
		var typName string
		switch fdType := fd.Type.(type) {
		case *ast.Ident:
			typName = fdType.String()
		case *ast.StarExpr:
			switch expr := fdType.X.(type) {
			case *ast.Ident:
				typName = fmt.Sprintf("*%s", expr.String())
			case *ast.SelectorExpr:
				x := expr.X.(*ast.Ident).String()
				name := expr.Sel.String()
				typName = fmt.Sprintf("*%s.%s", x, name)
			}
		case *ast.SelectorExpr:
			x := fdType.X.(*ast.Ident).String()
			name := fdType.Sel.String()
			typName = fmt.Sprintf("%s.%s", x, name)
		case *ast.ArrayType:
			switch ele := fdType.Elt.(type) {
			case *ast.Ident:
				typName = fmt.Sprintf("[]%s", ele.String())
			}
		}
		t.fields = append(t.fields, Field{
			Name: fd.Names[0].String(),
			Type: typName,
		})
		return nil
	}
	return t
}

type Type struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}
