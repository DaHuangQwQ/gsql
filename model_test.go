package gsql

import (
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	testCases := []struct {
		name string

		entity any

		wantModel *model
		wantErr   error
		opts      []ModelOption
	}{
		{
			name:   "test model",
			entity: TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},
		{
			name:   "test pointer to model",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},
		{
			name: "another type(not struct)",
			entity: map[string]any{
				"1": 1,
			},
			wantErr: errs.ErrInvalidType,
		},
		{
			name:   "test pointer to model",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model_t",
				fields: map[string]*field{
					"Id": {
						colName: "id_t",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
			opts: []ModelOption{
				ModelWithTableName("test_model_t"),
				ModelWithColumnName("Id", "id_t"),
			},
		},
	}

	r := registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parseModel, err := r.Register(tc.entity, tc.opts...)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, parseModel)
		})
	}
}

func TestRegistry_Get(t *testing.T) {
	testCases := []struct {
		name string

		entity any

		wantModel *model
		wantErr   error
	}{
		{
			name:   "test model",
			entity: TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},
		{
			name: "test tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column = first_name_t"`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name_t",
					},
				},
			},
		},
		{
			name: "test tag: empty column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column = "`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
		{
			name: "test tag: only column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column "`
				}
				return &TagTable{}
			}(),
			wantErr: errs.NewErrInvalidTagContent("column"),
		},
		{
			name: "test tag: not my tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"abc = abc"`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
		{
			name:   "test custom table name",
			entity: &customTableName{},
			wantModel: &model{
				tableName: "custom_table_name_t",
				fields: map[string]*field{
					"Name": {
						colName: "name",
					},
				},
			},
		},
		{
			name:   "test custom empty table name",
			entity: &customEmptyTableName{},
			wantModel: &model{
				tableName: "custom_empty_table_name",
				fields: map[string]*field{
					"Name": {
						colName: "name",
					},
				},
			},
		},
	}

	r := newRegistry()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m, err := r.Get(testCase.entity)
			assert.Equal(t, testCase.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, testCase.wantModel, m)
			//assert.Equal(t, testCase.cacheSize, len(r.models))
			typ := reflect.TypeOf(testCase.entity)
			res, ok := r.models.Load(typ)
			m = res.(*model)
			assert.True(t, ok)
			assert.Equal(t, testCase.wantModel, m)
		})
	}
}

type customTableName struct {
	Name string
}

func (c customTableName) TableName() string {
	return "custom_table_name_t"
}

type customEmptyTableName struct {
	Name string
}

func (c customEmptyTableName) TableName() string {
	return ""
}
