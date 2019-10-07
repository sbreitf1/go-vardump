package vardump

import (
	"reflect"
	"strings"
)

var (
	// ObscureByDefault denotes whether a value requires the "obscure" tag to be obscured. When set to true, only values with "safe" tag are printed in plain text.
	ObscureByDefault = false
)

// Visitor describes functionality for traversing an object.
type Visitor interface {
	Pointer()
	Value(value interface{}, obscure bool)
	BeginStruct(size int)
	StructValueName(index int, name string)
	EndStruct()
	BeginArray(size int)
	ArrayValueIndex(index int)
	EndArray()
}

func visit(obj interface{}, visitor Visitor, obscure bool) error {
	t := reflect.TypeOf(obj)

	switch t.Kind() {
	case reflect.Ptr:
		return visitPointer(obj, visitor, obscure)

	case reflect.Array:
		return visitSlice(obj, visitor, obscure)
	case reflect.Slice:
		return visitSlice(obj, visitor, obscure)

	case reflect.Struct:
		return visitStruct(obj, visitor, obscure)

	default:
		visitor.Value(obj, obscure)
		return nil
	}
}

func visitPointer(obj interface{}, visitor Visitor, obscure bool) error {
	visitor.Pointer()
	return visit(reflect.ValueOf(obj).Elem().Interface(), visitor, obscure)
}

func visitSlice(obj interface{}, visitor Visitor, obscure bool) error {
	v := reflect.ValueOf(obj)
	count := v.Len()
	visitor.BeginArray(count)
	for i := 0; i < count; i++ {
		visitor.ArrayValueIndex(i)
		if err := visit(v.Index(i).Interface(), visitor, obscure); err != nil {
			return err
		}
	}
	visitor.EndArray()
	return nil
}

type fieldMetaData struct {
	Index     int
	FieldName string
	Obscure   bool
}

func visitStruct(obj interface{}, visitor Visitor, obscure bool) error {
	fields := make([]fieldMetaData, 0)

	t := reflect.TypeOf(obj)
	totalFieldCount := t.NumField()
	for i := 0; i < totalFieldCount; i++ {
		f := t.Field(i)
		fieldName := f.Name
		childObscure := obscure

		tag := f.Tag.Get("vardump")
		if len(tag) > 0 {
			if tag == "-" {
				// skip this field
				continue
			}

			parts := strings.Split(tag, ",")
			if len(parts[0]) > 0 {
				fieldName = parts[0]
			}
			if len(parts) > 1 && parts[1] == "obscure" {
				childObscure = true
			}
		}

		fields = append(fields, fieldMetaData{i, fieldName, childObscure})
	}

	v := reflect.ValueOf(obj)
	visitor.BeginStruct(len(fields))
	for i, f := range fields {
		visitor.StructValueName(i, f.FieldName)
		if err := visit(v.Field(f.Index).Interface(), visitor, f.Obscure); err != nil {
			return err
		}
	}
	visitor.EndStruct()
	return nil
}
