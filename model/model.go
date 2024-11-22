package model

import (
	"github.com/DaHuangQwQ/gsql/internal/errs"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

const (
	tagKeyColumn = "column"
)

type Model struct {
	TableName string

	Fields []*Field
	// FieldMap 字段名到字段定义的映射
	FieldMap map[string]*Field
	// ColumnMap 列名到字段定义的映射
	ColumnMap map[string]*Field
}

type Field struct {
	GoName  string
	ColName string
	Typ     reflect.Type
	Offset  uintptr
}

// registry 元数据的注册中心
type registry struct {
	models sync.Map
	//lock   sync.RWMutex
}

func NewRegistry() *registry {
	return &registry{
		models: sync.Map{},
	}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	m, err := r.Register(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), nil
}

func (r *registry) Register(entity any, opts ...ModelOption) (*Model, error) {
	tye := reflect.TypeOf(entity)

	for tye.Kind() == reflect.Pointer {
		tye = tye.Elem()
	}

	if tye.Kind() != reflect.Struct {
		return nil, errs.ErrInvalidType
	}

	numFields := tye.NumField()

	fieldMap := make(map[string]*Field, numFields)
	columnMap := make(map[string]*Field, numFields)
	fields := make([]*Field, 0, numFields)

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

		fdMeta := &Field{
			GoName:  fd.Name,
			ColName: underscoreName(colName),
			Typ:     fd.Type,
			Offset:  fd.Offset,
		}

		fieldMap[fd.Name] = fdMeta
		columnMap[underscoreName(colName)] = fdMeta
		fields = append(fields, fdMeta)
	}

	tableName := ""
	if val, ok := entity.(TableName); ok {
		tableName = val.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(tye.Name())
	}

	res := &Model{
		TableName: tableName,
		Fields:    fields,
		FieldMap:  fieldMap,
		ColumnMap: columnMap,
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

func WithTableName(tableName string) ModelOption {
	return func(m *Model) error {
		m.TableName = tableName
		return nil
	}
}

func WithColumnName(fieldName string, colName string) ModelOption {
	return func(m *Model) error {
		fd, ok := m.FieldMap[fieldName]
		if !ok {
			return errs.NewErrUnknownField(fieldName)
		}
		fd.ColName = colName
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
