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

func visitStruct(obj interface{}, visitor Visitor, obscure bool) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	fieldCount := t.NumField()
	//TODO estimate correct field count -> ignored struct fields
	visitor.BeginStruct(fieldCount)
	for i := 0; i < fieldCount; i++ {
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

		visitor.StructValueName(i, fieldName)
		if err := visit(v.Field(i).Interface(), visitor, childObscure); err != nil {
			return err
		}
	}
	visitor.EndStruct()
	return nil
}
