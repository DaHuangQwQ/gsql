package integration

import (
	"context"
	"github.com/DaHuangQwQ/gsql"
	"github.com/DaHuangQwQ/gsql/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type InsertSuite struct {
	Suite
}

func TestMySQLInsert(t *testing.T) {
	suite.Run(t, &InsertSuite{
		Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

func (s *InsertSuite) TearDownTest() {
	gsql.RawQuery[test.SimpleStruct](s.db, "TRUNCATE TABLE `simple_struct`").Exec(context.Background())
}

func (i *InsertSuite) TestInsert() {
	db := i.db
	t := i.T()
	testCases := []struct {
		name         string
		i            *gsql.Inserter[test.SimpleStruct]
		wantAffected int64 // 插入行数
	}{
		{
			name:         "insert one",
			i:            gsql.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(12)),
			wantAffected: 1,
		},
		{
			name: "insert multiple",
			i: gsql.NewInserter[test.SimpleStruct](db).Values(
				test.NewSimpleStruct(13),
				test.NewSimpleStruct(14)),
			wantAffected: 2,
		},
		{
			name:         "insert id",
			i:            gsql.NewInserter[test.SimpleStruct](db).Values(&test.SimpleStruct{Id: 15}),
			wantAffected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			res := tc.i.Exec(ctx)
			affected, err := res.RowsAffected()
			assert.NoError(t, err)
			assert.Equal(t, tc.wantAffected, affected)
		})
	}
}
