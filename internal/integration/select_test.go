package integration

import (
	"context"
	"github.com/DaHuangQwQ/gsql"
	"github.com/DaHuangQwQ/gsql/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SelectSuite struct {
	Suite
}

func TestMySQLSelect(t *testing.T) {
	suite.Run(t, &SelectSuite{
		Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

func (s *InsertSuite) TearDownSuite() {
	gsql.RawQuery[test.SimpleStruct](s.db, "TRUNCATE TABLE `simple_struct`").Exec(context.Background())
}

func (s *SelectSuite) SetupSuite() {
	s.Suite.SetupSuite()
	res := gsql.NewInserter[test.SimpleStruct](s.db).Values(
		test.NewSimpleStruct(100)).Exec(context.Background())
	require.NoError(s.T(), res.Err())
}

func (s *SelectSuite) TestGet() {
	testCases := []struct {
		name string
		s    *gsql.Selector[test.SimpleStruct]

		wantRes *test.SimpleStruct
		wantErr error
	}{
		{
			name:    "get data",
			s:       gsql.NewSelector[test.SimpleStruct](s.db).Where(gsql.C("Id").Eq(100)),
			wantRes: test.NewSimpleStruct(100),
		},
		{
			name:    "no row",
			s:       gsql.NewSelector[test.SimpleStruct](s.db).Where(gsql.C("Id").Eq(200)),
			wantErr: gsql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			res, err := tc.s.Get(ctx)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
