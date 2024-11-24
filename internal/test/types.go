package test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/DaHuangQwQ/gutil"
)

// SimpleStruct 包含所有支持的类型
type SimpleStruct struct {
	Id      uint64
	Bool    bool
	BoolPtr *bool

	Int    int
	IntPtr *int

	Int8    int8
	Int8Ptr *int8

	Int16    int16
	Int16Ptr *int16

	Int32    int32
	Int32Ptr *int32

	Int64    int64
	Int64Ptr *int64

	Uint    uint
	UintPtr *uint

	Uint8    uint8
	Uint8Ptr *uint8

	Uint16    uint16
	Uint16Ptr *uint16

	Uint32    uint32
	Uint32Ptr *uint32

	Uint64    uint64
	Uint64Ptr *uint64

	Float32    float32
	Float32Ptr *float32

	Float64    float64
	Float64Ptr *float64

	Byte      byte
	BytePtr   *byte
	ByteArray []byte

	String string

	// 特殊类型
	NullStringPtr  *sql.NullString
	NullInt16Ptr   *sql.NullInt16
	NullInt32Ptr   *sql.NullInt32
	NullInt64Ptr   *sql.NullInt64
	NullBoolPtr    *sql.NullBool
	NullFloat64Ptr *sql.NullFloat64
	JsonColumn     *JsonColumn
}

// JsonColumn 是自定义的 JSON 类型字段
// Val 字段必须是结构体指针
type JsonColumn struct {
	Val   User
	Valid bool
}

type User struct {
	Name string
}

func (j *JsonColumn) Scan(src any) error {
	if src == nil {
		return nil
	}
	var bs []byte
	switch val := src.(type) {
	case string:
		bs = []byte(val)
	case []byte:
		bs = val
	case *[]byte:
		if val == nil {
			return nil
		}
		bs = *val
	default:
		return fmt.Errorf("不合法类型 %+v", src)
	}
	if len(bs) == 0 {
		return nil
	}
	err := json.Unmarshal(bs, &j.Val)
	if err != nil {
		return err
	}
	j.Valid = true
	return nil
}

// Value 参考 sql.NullXXX 类型定义的
func (j JsonColumn) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	bs, err := json.Marshal(j.Val)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func NewSimpleStruct(id uint64) *SimpleStruct {
	return &SimpleStruct{
		Id:            id,
		Bool:          true,
		BoolPtr:       gutil.ToPtr[bool](false),
		Int:           12,
		IntPtr:        gutil.ToPtr[int](13),
		Int8:          8,
		Int8Ptr:       gutil.ToPtr[int8](-8),
		Int16:         16,
		Int16Ptr:      gutil.ToPtr[int16](-16),
		Int32:         32,
		Int32Ptr:      gutil.ToPtr[int32](-32),
		Int64:         64,
		Int64Ptr:      gutil.ToPtr[int64](-64),
		Uint:          14,
		UintPtr:       gutil.ToPtr[uint](15),
		Uint8:         8,
		Uint8Ptr:      gutil.ToPtr[uint8](18),
		Uint16:        16,
		Uint16Ptr:     gutil.ToPtr[uint16](116),
		Uint32:        32,
		Uint32Ptr:     gutil.ToPtr[uint32](132),
		Uint64:        64,
		Uint64Ptr:     gutil.ToPtr[uint64](164),
		Float32:       3.2,
		Float32Ptr:    gutil.ToPtr[float32](-3.2),
		Float64:       6.4,
		Float64Ptr:    gutil.ToPtr[float64](-6.4),
		Byte:          byte(8),
		BytePtr:       gutil.ToPtr[byte](18),
		ByteArray:     []byte("hello"),
		String:        "world",
		NullStringPtr: &sql.NullString{String: "null string", Valid: true},
		NullInt16Ptr:  &sql.NullInt16{Int16: 16, Valid: true},
		NullInt32Ptr:  &sql.NullInt32{Int32: 32, Valid: true},
		NullInt64Ptr:  &sql.NullInt64{Int64: 64, Valid: true},
		NullBoolPtr:   &sql.NullBool{Bool: true, Valid: true},
		// NullTimePtr: &sql.NullTime{Time: t, Valid: true},
		NullFloat64Ptr: &sql.NullFloat64{Float64: 6.4, Valid: true},
		JsonColumn: &JsonColumn{
			Val:   User{Name: "Tom"},
			Valid: true,
		},
	}
}
