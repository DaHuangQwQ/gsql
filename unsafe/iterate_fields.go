package unsafe

import "reflect"

func PrintFieldOffset(entity any) {
	typ := reflect.TypeOf(entity)
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		field := typ.Field(i)
		println(field.Offset)
	}
}
