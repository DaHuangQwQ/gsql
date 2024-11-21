package gsql

import (
	"database/sql"
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"github.com/stretchr/testify/assert"
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
