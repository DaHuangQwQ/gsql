package gsql

import (
	"database/sql"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	testCases := []struct {
		name string

		selector QueryBuilder

		wantQuery *Query
		wantErr   error
	}{
		{
			name:     "select",
			selector: &Selector[TestModel]{},
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "from",
			selector: (&Selector[TestModel]{}).From("test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "empty from",
			selector: (&Selector[TestModel]{}).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "from db",
			selector: (&Selector[TestModel]{}).From("test_db.test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`test_model`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "where",
			selector: (&Selector[TestModel]{}).Where(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel` WHERE `Age` = ?;",
				Args: []any{
					18,
				},
			},
			wantErr: nil,
		},
		{
			name:     "empty where",
			selector: (&Selector[TestModel]{}).Where(),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
			wantErr: nil,
		},
		{
			name:     "not",
			selector: (&Selector[TestModel]{}).Where(Not(C("Age").Eq(18))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel` WHERE  NOT (`Age` = ?);",
				Args: []any{
					18,
				},
			},
			wantErr: nil,
		},
		{
			name:     "and",
			selector: (&Selector[TestModel]{}).Where(C("Age").Eq(18).And(C("Name").Eq("dahuang"))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel` WHERE (`Age` = ?) AND (`Name` = ?);",
				Args: []any{
					18,
					"dahuang",
				},
			},
			wantErr: nil,
		},
		{
			name:     "or",
			selector: (&Selector[TestModel]{}).Where(C("Age").Eq(18).Or(C("Name").Eq("dahuang"))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel` WHERE (`Age` = ?) OR (`Name` = ?);",
				Args: []any{
					18,
					"dahuang",
				},
			},
			wantErr: nil,
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
