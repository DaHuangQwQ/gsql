package gsql

import (
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseModel(t *testing.T) {
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
	}

	r := registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parseModel, err := r.ParseModel(tc.entity)
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
		cacheSize int
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
			cacheSize: 1,
		},
	}

	r := newRegistry()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m, err := r.get(testCase.entity)
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
