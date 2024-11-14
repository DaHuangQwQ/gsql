package gsql

import (
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string
}

// registry 元数据的注册中心
type registry struct {
	models sync.Map
	//lock   sync.RWMutex
}

func newRegistry() *registry {
	return &registry{
		models: sync.Map{},
	}
}

func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}
	m, err := r.ParseModel(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*model), nil
}

// get RWMutex double check
//func (r *registry) getV1(val any) (*model, error) {
//	typ := reflect.TypeOf(val)
//
//	r.lock.RLock()
//	m, ok := r.models[typ]
//	r.lock.RUnlock()
//	if ok {
//		return m, nil
//	}
//
//	r.lock.Lock()
//	defer r.lock.Unlock()
//	m, ok = r.models[typ]
//	if ok {
//		return m, nil
//	}
//
//	var er error
//	m, er = r.ParseModel(val)
//	if er != nil {
//		return nil, er
//	}
//	r.models[typ] = m
//
//	return m, nil
//}

func (r *registry) ParseModel(entity any) (*model, error) {
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
