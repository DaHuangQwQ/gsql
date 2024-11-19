package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type UnsafeAccessor struct {
	fields  map[string]FieldMeta
	address unsafe.Pointer
}

func NewUnsafeAccessor(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity)
	typ = typ.Elem()
	numField := typ.NumField()
	fields := make(map[string]FieldMeta, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fields[fd.Name] = FieldMeta{
			Offset: fd.Offset,
			typ:    fd.Type,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnsafeAccessor{
		fields:  fields,
		address: val.UnsafePointer(),
	}
}

func (a *UnsafeAccessor) Field(field string) (any, error) {
	// 起始地址 + 字段偏移量
	fd, ok := a.fields[field]
	if !ok {
		return nil, errors.New("非法字段")
	}
	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)
	// 如果知道类型，就这么读
	// return *(*int)(fdAddress), nil

	// 不知道确切类型
	return reflect.NewAt(fd.typ, fdAddress).Elem().Interface(), nil
}

func (a *UnsafeAccessor) SetField(field string, val any) error {
	// 起始地址 + 字段偏移量
	fd, ok := a.fields[field]
	if !ok {
		return errors.New("非法字段")
	}
	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)

	// 你知道确切类型就这么写
	// *(*int)(fdAddress) = val.(int)

	// 你不知道确切类型
	reflect.NewAt(fd.typ, fdAddress).Elem().Set(reflect.ValueOf(val))
	return nil
}

type FieldMeta struct {
	Offset uintptr
	typ    reflect.Type
}
