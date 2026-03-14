package payload

import (
	"reflect"
)

var PayloadSizeRegistry = make(map[reflect.Type]uintptr)

func GetRealSize(t reflect.Type) uintptr {
	if t.Kind() != reflect.Struct {
		return t.Size()
	}

	var totalSize uintptr
	for i := range t.NumField() {
		field := t.Field(i)

		if field.Type.Kind() == reflect.Struct {
			totalSize += GetRealSize(field.Type)

		} else {
			totalSize += field.Type.Size()
		}
	}

	return totalSize
}

func init() {
	register := func(p IPayload) {
		t := reflect.TypeOf(p).Elem()
		PayloadSizeRegistry[t] = GetRealSize(t)
	}

	register(&ManualExecResponsePayload{})
	register(&StatusResponsePayload{})
}

type IRequest interface {
	ToBytes() []byte
}

type NoOpToBytes struct{}

func (NoOpToBytes) ToBytes() []byte { return []byte{} }

type IResponse interface {
	Validate() error
}

type NoOpValidate struct{}

func (NoOpValidate) Validate() error { return nil }

type IPayload interface {
	IRequest
	IResponse
}

type EmptyPayload struct {
	NoOpToBytes
	NoOpValidate
}
