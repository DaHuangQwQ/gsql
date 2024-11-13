package reflect

import (
	"github.com/DaHuangQwQ/gweb/reflect/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFunc(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantRes map[string]FuncInfo
		wantErr error
	}{
		{
			name:   "struct",
			entity: types.NewUser("Tom", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",
					// 下标 0 的指向接收器
					InputTypes:  []reflect.Type{reflect.TypeOf(types.User{})},
					OutputTypes: []reflect.Type{reflect.TypeOf(0)},
					Result:      []any{18},
				},
				// "ChangeName": {
				// 	Name: "ChangeName",
				// 	InputTypes: []reflect.Type{reflect.TypeOf("")},
				// 	// OutputTypes: []reflect.Type{reflect.TypeOf(0)},
				// 	// Result: []any{18},
				// },
			},
		},
		{
			name:   "pointer",
			entity: types.NewUserPtr("Tom", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",
					// 下标 0 的指向接收器
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{})},
					OutputTypes: []reflect.Type{reflect.TypeOf(0)},
					Result:      []any{18},
				},
				"ChangeName": {
					Name:        "ChangeName",
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{}), reflect.TypeOf("")},
					OutputTypes: []reflect.Type{},
					Result:      []any{},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFunc(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
