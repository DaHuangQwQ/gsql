package gsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInserter_SQLite_upsert(t *testing.T) {
	db := memoryDB(t, WithDialect(DialectSQLite))
	testCases := []struct {
		name string
		i    QueryBuilder

		wantErr error
		wantRes *Query
	}{
		{
			name: "upsert-update value",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Jerry", Valid: true},
			}).OnDuplicateKey().ConflictColumns("Id").Update(Assign("FirstName", "Huang"),
				Assign("Age", 19)),
			wantRes: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`last_name`,`age`) VALUES (?,?,?,?) " +
					"ON CONFLICT(`id`) DO UPDATE SET `first_name`=?,`age`=?;",
				Args: []any{int64(12), "Tom", &sql.NullString{String: "Jerry", Valid: true}, int8(18), "Huang", 19},
			},
		},
		{
			name: "upsert-update column",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Jerry", Valid: true},
			}, &TestModel{
				Id:        13,
				FirstName: "Da",
				Age:       19,
				LastName:  &sql.NullString{String: "Huang", Valid: true},
			}).OnDuplicateKey().ConflictColumns("FirstName", "LastName").Update(C("FirstName"), C("Age")),
			wantRes: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`last_name`,`age`) VALUES (?,?,?,?),(?,?,?,?) " +
					"ON CONFLICT(`first_name`,`last_name`) DO UPDATE SET `first_name`=excluded.`first_name`,`age`=excluded.`age`;",
				Args: []any{int64(12), "Tom", &sql.NullString{String: "Jerry", Valid: true}, int8(18),
					int64(13), "Da", &sql.NullString{String: "Huang", Valid: true}, int8(19)},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.i.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, q)
		})
	}
}

func TestInserter_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name string
		i    QueryBuilder

		wantRes *Query
		wantErr error
	}{
		{
			name: "single row",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        1,
				FirstName: "Da",
				LastName: &sql.NullString{
					String: "Huang",
					Valid:  true,
				},
				Age: 18,
			}),
			wantRes: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`,`last_name`,`age`) VALUES (?,?,?,?);",
				Args: []any{int64(1), "Da", &sql.NullString{String: "Huang", Valid: true}, int8(18)},
			},
		},
		{
			name: "multiple row",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        1,
				FirstName: "Da",
				LastName: &sql.NullString{
					String: "Huang",
					Valid:  true,
				},
				Age: 18,
			}, &TestModel{
				Id:        2,
				FirstName: "Xiao",
				LastName: &sql.NullString{
					String: "Huang",
					Valid:  true,
				},
				Age: 20,
			}),
			wantRes: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`last_name`,`age`) VALUES (?,?,?,?),(?,?,?,?);",
				Args: []any{int64(1), "Da", &sql.NullString{String: "Huang", Valid: true}, int8(18),
					int64(2), "Xiao", &sql.NullString{String: "Huang", Valid: true}, int8(20)},
			},
		},
		{
			name:    "no value in insert",
			i:       NewInserter[TestModel](db),
			wantErr: errs.ErrInsertZeroRow,
		},
		{
			// 插入多行、部分列
			name: "partial columns",
			i: NewInserter[TestModel](db).Columns("Id", "FirstName").Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Jerry", Valid: true},
			}, &TestModel{
				Id:        13,
				FirstName: "DaMing",
				Age:       19,
				LastName:  &sql.NullString{String: "Deng", Valid: true},
			}),
			wantRes: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`) VALUES (?,?),(?,?);",
				Args: []any{int64(12), "Tom", int64(13), "DaMing"},
			},
		},
		{
			name: "upsert-update value",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Jerry", Valid: true},
			}).OnDuplicateKey().Update(Assign("FirstName", "Deng"),
				Assign("Age", 19)),
			wantRes: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`last_name`,`age`) VALUES (?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `first_name`=?,`age`=?;",
				Args: []any{int64(12), "Tom", &sql.NullString{String: "Jerry", Valid: true}, int8(18), "Deng", 19},
			},
		},
		{
			name: "upsert-update column",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{String: "Jerry", Valid: true},
			}, &TestModel{
				Id:        13,
				FirstName: "DaMing",
				Age:       19,
				LastName:  &sql.NullString{String: "Deng", Valid: true},
			}).OnDuplicateKey().Update(C("FirstName"), C("Age")),
			wantRes: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`last_name`,`age`) VALUES (?,?,?,?),(?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `first_name`=VALUES(`first_name`),`age`=VALUES(`age`);",
				Args: []any{int64(12), "Tom", &sql.NullString{String: "Jerry", Valid: true}, int8(18),
					int64(13), "DaMing", &sql.NullString{String: "Deng", Valid: true}, int8(19)},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			q, err := testCase.i.Build()
			assert.Equal(t, testCase.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, testCase.wantRes, q)
		})
	}
}

func TestInserter_Exec(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)
	testCases := []struct {
		name     string
		i        *Inserter[TestModel]
		wantErr  error
		affected int64
	}{
		{
			name: "query error",
			i: func() *Inserter[TestModel] {
				return NewInserter[TestModel](db).Values(&TestModel{}).
					Columns("Invalid")
			}(),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "db error",
			i: func() *Inserter[TestModel] {
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(errors.New("db error"))
				return NewInserter[TestModel](db).Values(&TestModel{})
			}(),
			wantErr: errors.New("db error"),
		},
		{
			name: "exec",
			i: func() *Inserter[TestModel] {
				res := driver.RowsAffected(1)
				mock.ExpectExec("INSERT INTO .*").
					WillReturnResult(res)
				return NewInserter[TestModel](db).Values(&TestModel{})
			}(),
			affected: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			affected, err := res.RowsAffected()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.affected, affected)
		})
	}
}
