package main

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFileVisitor_Get(t *testing.T) {
	testCases := []struct {
		src  string
		want File
	}{
		{
			src: `
package orm_gen
import (
	"fmt"
    "database/sql"
) 

import (
	dri "database/sql/driver"
)
type (
	StructType struct {
		// Public is a field
		// @type string
		Public string
        Ptr *sql.NullString
		Struct sql.NullInt64
		Age *int8
		Slice []byte
	}
)
`,
			want: File{
				Package: "orm_gen",
				Imports: []string{`"fmt"`, `"database/sql"`, `dri "database/sql/driver"`},
				Types: []Type{
					{
						Name: "StructType",
						Fields: []Field{
							{
								Name: "Public",
								Type: "string",
							},
							{
								Name: "Ptr",
								Type: "*sql.NullString",
							},
							{
								Name: "Struct",
								Type: "sql.NullInt64",
							},
							{
								Name: "Age",
								Type: "*int8",
							},
							{
								Name: "Slice",
								Type: "[]byte",
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "src.go", tc.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		tv := &SingleFileEntryVisitor{}
		ast.Walk(tv, f)
		file := tv.Get()
		assert.Equal(t, tc.want, file)
	}
}
