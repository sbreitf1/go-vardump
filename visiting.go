package vardump

import (
	"reflect"
)

// Visitor describes functionality for traversing an object.
type Visitor interface {
	Pointer()
	Value(value interface{})
	BeginStruct(size int)
	StructValueName(index int, name string)
	EndStruct()
	BeginArray(size int)
	ArrayValueIndex(index int)
	EndArray()
}

func visit(obj interface{}, visitor Visitor) error {
	t := reflect.TypeOf(obj)

	switch t.Kind() {
	case reflect.Ptr:
		return visitPointer(obj, visitor)

	case reflect.Array:
		return visitSlice(obj, visitor)
	case reflect.Slice:
		return visitSlice(obj, visitor)

	case reflect.Struct:
		return visitStruct(obj, visitor)

	default:
		visitor.Value(obj)
		return nil
	}
}

func visitPointer(obj interface{}, visitor Visitor) error {
	visitor.Pointer()
	return visit(reflect.ValueOf(obj).Elem().Interface(), visitor)
}

func visitSlice(obj interface{}, visitor Visitor) error {
	v := reflect.ValueOf(obj)
	count := v.Len()
	visitor.BeginArray(count)
	for i := 0; i < count; i++ {
		visitor.ArrayValueIndex(i)
		if err := visit(v.Index(i).Interface(), visitor); err != nil {
			return err
		}
	}
	visitor.EndArray()
	return nil
}

func visitStruct(obj interface{}, visitor Visitor) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	fieldCount := t.NumField()
	visitor.BeginStruct(fieldCount)
	for i := 0; i < fieldCount; i++ {
		f := t.Field(i)
		visitor.StructValueName(i, f.Name)
		if err := visit(v.Field(i).Interface(), visitor); err != nil {
			return err
		}
	}
	visitor.EndStruct()
	return nil
}
