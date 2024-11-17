package gsql

import (
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

const (
	tagKeyColumn = "column"
)

type Registry interface {
	Get(val any) (*model, error)
	Register(val any, opts ...ModelOption) (*model, error)
}

type ModelOption func(m *model) error

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

func (r *registry) Get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}
	m, err := r.Register(val)
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
//	m, er = r.Register(val)
//	if er != nil {
//		return nil, er
//	}
//	r.models[typ] = m
//
//	return m, nil
//}

func (r *registry) Register(entity any, opts ...ModelOption) (*model, error) {
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
		tags, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName, _ := tags[tagKeyColumn]
		if colName == "" {
			colName = fd.Name
		}
		fieldMap[fd.Name] = &field{
			colName: underscoreName(colName),
		}
	}

	tableName := ""
	if val, ok := entity.(TableName); ok {
		tableName = val.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(tye.Name())
	}

	res := &model{
		tableName: tableName,
		fields:    fieldMap,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := strings.TrimSpace(segs[0])
		val := strings.TrimSpace(segs[1])
		res[key] = val
	}
	return res, nil
}

func ModelWithTableName(tableName string) ModelOption {
	return func(m *model) error {
		m.tableName = tableName
		return nil
	}
}

func ModelWithColumnName(fieldName string, colName string) ModelOption {
	return func(m *model) error {
		fd, ok := m.fields[fieldName]
		if !ok {
			return errs.NewErrUnknownField(fieldName)
		}
		fd.colName = colName
		return nil
	}
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
