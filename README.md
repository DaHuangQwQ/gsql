# gsql
an go easy orm framework

```shell
go get github.com/DaHuangQwQ/gsql
```

## 示例
```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DaHuangQwQ/gsql"
	_ "github.com/go-sql-driver/mysql"
)

type TestModel struct {
	Id        int64
	FirstName string
	LastName  *sql.NullString
	Age       int8
}

func main() {
	db, err := gsql.Open("mysql", "root:root@tcp(localhost:13306)/test_model")
	if err != nil {
		panic(err)
	}
	
	gsql.RawQuery[TestModel](db, "TRUNCATE TABLE `test_model`.test_model").Exec(context.Background())

	res := gsql.NewInserter[TestModel](db).Values(&TestModel{
		Id:        1,
		FirstName: "Da",
		LastName: &sql.NullString{
			String: "Huang",
			Valid:  true,
		},
		Age: 18,
	}).Exec(context.Background())
	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("affected:", affected)

	get, err := gsql.NewSelector[TestModel](db).Where(gsql.C("Age").Eq(18)).Get(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("get: %v", get)
}
```

## aop
- cache
- nodelete
- opentelemetry
- prometheus
- querylog
- safedml
- slowquery

## gen
ast + template 代码生成
