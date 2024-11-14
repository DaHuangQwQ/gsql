package gsql

import (
	"database/sql"
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	db, err := NewDB()
	require.NoError(t, err)
	testCases := []struct {
		name string

		selector QueryBuilder

		wantQuery *Query
		wantErr   error
	}{
		{
			name:     "select",
			selector: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "from",
			selector: NewSelector[TestModel](db).From("TestModel"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "empty from",
			selector: NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "from db",
			selector: NewSelector[TestModel](db).From("test_db.test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`test_model`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "where",
			selector: NewSelector[TestModel](db).Where(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE `age` = ?;",
				Args: []any{
					18,
				},
			},
			wantErr: nil,
		},
		{
			name:     "empty where",
			selector: NewSelector[TestModel](db).Where(),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "not",
			selector: NewSelector[TestModel](db).Where(Not(C("Age").Eq(18))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{
					18,
				},
			},
			wantErr: nil,
		},
		{
			name:     "and",
			selector: NewSelector[TestModel](db).Where(C("Age").Eq(18).And(C("FirstName").Eq("dahuang"))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`first_name` = ?);",
				Args: []any{
					18,
					"dahuang",
				},
			},
			wantErr: nil,
		},
		{
			name:     "or",
			selector: NewSelector[TestModel](db).Where(C("Age").Eq(18).Or(C("FirstName").Eq("dahuang"))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE (`age` = ?) OR (`first_name` = ?);",
				Args: []any{
					18,
					"dahuang",
				},
			},
			wantErr: nil,
		},
		{
			name:     "unknown field",
			selector: NewSelector[TestModel](db).Where(C("Age").Eq(18).Or(C("XX").Eq("dahuang"))),
			wantErr:  errs.NewErrUnknownField("XX"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			q, err := testCase.selector.Build()
			assert.Equal(t, testCase.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, testCase.wantQuery, q)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	LastName  *sql.NullString
	Age       int8
}
