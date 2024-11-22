package valuer

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	gsql "github.com/DaHuangQwQ/gsql/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_reflectValue_SetColumns(t *testing.T) {
	testSetColumns(t, NewReflectValue)
}

func testSetColumns(t *testing.T, creator Creator) {
	testCases := []struct {
		name string
		// 一定是指针
		entity any
		rows   func() *sqlmock.Rows

		wantErr    error
		wantEntity any
	}{
		{
			name:   "set columns",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow("1", "Tom", "18", "Jerry")
				return rows
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},

		{
			// 测试列的不同顺序
			name:   "order",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "last_name", "first_name", "age"})
				rows.AddRow("1", "Jerry", "Tom", "18")
				return rows
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},

		{
			// 测试列的不同顺序
			name:   "partial columns",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "last_name"})
				rows.AddRow("1", "Jerry")
				return rows
			},
			wantEntity: &TestModel{
				Id:       1,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
	}

	r := gsql.NewRegistry()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 构造 rows
			mockRows := tc.rows()
			mock.ExpectQuery("SELECT XX").WillReturnRows(mockRows)
			rows, err := mockDB.Query("SELECT XX")
			require.NoError(t, err)

			rows.Next()

			m, err := r.Get(tc.entity)
			require.NoError(t, err)
			val := creator(m, tc.entity)
			err = val.SetColumns(rows)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			// 比较一下 tc.entity 究竟有没有被设置好数据
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}

type TestModel struct {
	Id int64
	// ""
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
