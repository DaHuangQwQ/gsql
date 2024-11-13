package gsql

import (
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"reflect"
	"regexp"
	"strings"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string
}

func ParseModel(entity any) (*model, error) {
	tye := reflect.TypeOf(entity)

	for tye.Kind() == reflect.Pointer {
		tye = tye.Elem()
	}

	if tye.Kind() != reflect.Struct {
		return nil, errs.ErrInvalidType
	}

	numFields := tye.NumField()

	fieldMap := make(map[string]*field, numFields)

	for i := 0; i < numFields; i++ {
		fd := tye.Field(i)
		fieldMap[fd.Name] = &field{
			colName: underscoreName(fd.Name),
		}
	}
	return &model{
		tableName: underscoreName(tye.Name()),
		fields:    fieldMap,
	}, nil
}

// underscoreName 使用正则表达式将驼峰命名转为下划线命名
func underscoreName(tableName string) string {
	// ID => id
	if tableName == strings.ToUpper(tableName) {
		return strings.ToLower(tableName)
	}

	re := regexp.MustCompile("([A-Z])")

	result := re.ReplaceAllStringFunc(tableName, func(s string) string {
		return "_" + strings.ToLower(s)
	})

	if len(result) > 0 && result[0] == '_' {
		result = result[1:]
	}

	return result
}
