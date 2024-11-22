package gsql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DaHuangQwQ/gsql/internal/errs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		s         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "invalid column",
			s:       NewSelector[TestModel](db).Select(C("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "multiple columns",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("LastName")),
			wantQuery: &Query{
				SQL: "SELECT `first_name`,`last_name` FROM `test_model`;",
			},
		},
		{
			name: "columns alias",
			s:    NewSelector[TestModel](db).Select(C("FirstName").As("my_name"), C("LastName")),
			wantQuery: &Query{
				SQL: "SELECT `first_name` AS `my_name`,`last_name` FROM `test_model`;",
			},
		},
		{
			name: "avg",
			s:    NewSelector[TestModel](db).Select(Avg("Age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "avg alias",
			s:    NewSelector[TestModel](db).Select(Avg("Age").As("avg_age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`) AS `avg_age` FROM `test_model`;",
			},
		},
		{
			name: "sum",
			s:    NewSelector[TestModel](db).Select(Sum("Age")),
			wantQuery: &Query{
				SQL: "SELECT SUM(`age`) FROM `test_model`;",
			},
		},
		{
			name: "count",
			s:    NewSelector[TestModel](db).Select(Count("Age")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(`age`) FROM `test_model`;",
			},
		},
		{
			name: "max",
			s:    NewSelector[TestModel](db).Select(Max("Age")),
			wantQuery: &Query{
				SQL: "SELECT MAX(`age`) FROM `test_model`;",
			},
		},
		{
			name: "min",
			s:    NewSelector[TestModel](db).Select(Min("Age")),
			wantQuery: &Query{
				SQL: "SELECT MIN(`age`) FROM `test_model`;",
			},
		},
		{
			name:    "aggregate invalid columns",
			s:       NewSelector[TestModel](db).Select(Min("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "multiple aggregate",
			s:    NewSelector[TestModel](db).Select(Min("Age"), Max("Age")),
			wantQuery: &Query{
				SQL: "SELECT MIN(`age`),MAX(`age`) FROM `test_model`;",
			},
		},
		{
			name: "raw expression",
			s:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT `first_name`)")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.s.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

func TestSelector_Build(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	db, err := OpenDB(mockDB)
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
			name:     "from session",
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
		{
			name:     "raw expression as predicate",
			selector: NewSelector[TestModel](db).Where(Raw("`id`<?", 18).AsPredicate()),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`id`<?);",
				Args: []any{18},
			},
		},
		{
			name:     "raw expression used in predicate",
			selector: NewSelector[TestModel](db).Where(C("Id").Eq(Raw("`age`+?", 1))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = (`age`+?);",
				Args: []any{1},
			},
		},
		{
			name:     "columns alias in where",
			selector: NewSelector[TestModel](db).Where(C("Id").As("my_id").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = ?;",
				Args: []any{18},
			},
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

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("invalid query"))

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
	rows.AddRow("1", "da", "huang", "18")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	//rows = sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
	//rows.AddRow("xxx", "da", "huang", "18")
	//mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	testCases := []struct {
		name string
		s    *Selector[TestModel]

		wantErr error
		wantRes *TestModel
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("xxx").Eq(18)),
			wantErr: errs.NewErrUnknownField("xxx"),
		},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(18)),
			wantErr: errors.New("invalid query"),
		},
		{
			name:    "no rows",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(18)),
			wantErr: errs.ErrNoRows,
		},
		{
			name: "get row",
			s:    NewSelector[TestModel](db).Where(C("Id").Eq(18)),
			wantRes: &TestModel{
				Id:        1,
				FirstName: "da",
				LastName: &sql.NullString{
					String: "huang",
					Valid:  true,
				},
				Age: 18,
			},
		},
		//{
		//	name:    "get row: bad type",
		//	s:       NewSelector[TestModel](session).Where(C("Id").Eq(18)),
		//	wantErr: errs.NewErrUnknownColumn("id"),
		//},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, er := testCase.s.Get(context.Background())
			assert.Equal(t, testCase.wantErr, er)
			assert.Equal(t, testCase.wantRes, res)
		})
	}
}

func TestSelector_GetMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
	rows.AddRow("1", "da", "huang", "18")
	rows.AddRow("2", "xiao", "huang", "20")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	testCases := []struct {
		name string
		s    *Selector[TestModel]

		wantRes []*TestModel
		wantErr error
	}{
		{
			name: "get row",
			s:    NewSelector[TestModel](db).Where(C("Id").Eq(18)),
			wantRes: []*TestModel{
				{
					Id:        1,
					FirstName: "da",
					LastName: &sql.NullString{
						String: "huang",
						Valid:  true,
					},
					Age: 18,
				},
				{
					Id:        2,
					FirstName: "xiao",
					LastName: &sql.NullString{
						String: "huang",
						Valid:  true,
					},
					Age: 20,
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, er := testCase.s.GetMulti(context.Background())
			assert.Equal(t, testCase.wantErr, er)
			assert.Equal(t, testCase.wantRes, res)
		})
	}
}

func memoryDB(t *testing.T, opts ...DBOption) *DB {
	db, err := Open("sqlite3",
		"file:test.session?cache=shared&mode=memory",
		opts...)
	require.NoError(t, err)
	return db
}
