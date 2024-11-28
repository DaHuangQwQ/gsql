package gsql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleter_Build(t *testing.T) {
	db := memoryDB(t, WithDialect(DialectSQLite))
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no where",
			builder: (NewDeleter[TestModel](db)).From("test_model"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "where",
			builder: (NewDeleter[TestModel](db)).Where(C("Id").Eq(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id` = ?;",
				Args: []any{16},
			},
		},
		{
			name:    "from",
			builder: (NewDeleter[TestModel](db)).From("test_model").Where(C("Id").Eq(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id` = ?;",
				Args: []any{16},
			},
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			query, err := c.builder.Build()
			assert.Equal(t, c.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
