package vardump

import (
	"fmt"
	"strings"
)

// NestedStringOptions defines options to format nested printing.
type NestedStringOptions struct {
	LinePrefix                                                    string
	Pointer                                                       string
	Indentation                                                   string
	BeginStruct, EndStruct, StructNameFormat, StructItemSeparator string
	QuoteStructNames                                              bool
	BeginArray, EndArray, ArrayIndexFormat, ArrayElementSeparator string
	BreakEnumerationOnLen                                         int
	BreakEnumerationItems                                         bool
	BaseTypeOptions                                               *BaseTypeOptions
}

// DefaultNestedStringOptions returns the default nested string options for a JSON-like representation.
func DefaultNestedStringOptions() *NestedStringOptions {
	return &NestedStringOptions{
		LinePrefix:  "",
		Pointer:     "",
		Indentation: "  ",
		BeginStruct: "{", EndStruct: "}", StructNameFormat: "%q: ", StructItemSeparator: ",",
		QuoteStructNames: true,
		BeginArray:       "[", EndArray: "]", ArrayIndexFormat: "", ArrayElementSeparator: ",",
		BreakEnumerationOnLen: 1, BreakEnumerationItems: true,
		BaseTypeOptions: DefaultBaseTypeOptions(),
	}
}

type nestedStringVisitor struct {
	options        *NestedStringOptions
	enumLens       *intStack
	hierarchy      *stack
	sb             *stringBuilder
	requireNewLine bool
}

func (v *nestedStringVisitor) Pointer() {
	v.sb.Append(v.options.Pointer)
}
func (v *nestedStringVisitor) Value(value interface{}) {
	v.sb.AppendValue(value)
}
func (v *nestedStringVisitor) BeginStruct(size int) {
	v.sb.Append(v.options.BeginStruct)
	v.enterEnum(size)
}
func (v *nestedStringVisitor) StructValueName(index int, name string) {
	if index > 0 {
		v.sb.Append(v.options.StructItemSeparator)
		if v.options.BreakEnumerationItems {
			v.sb.AppendLine()
		}
	}
	v.sb.Append(fmt.Sprintf(v.options.StructNameFormat, name))
}
func (v *nestedStringVisitor) EndStruct() {
	v.leaveEnum()
	v.sb.Append(v.options.EndStruct)
}
func (v *nestedStringVisitor) BeginArray(size int) {
	v.sb.Append(v.options.BeginArray)
	v.enterEnum(size)
}
func (v *nestedStringVisitor) ArrayValueIndex(index int) {
	if index > 0 {
		v.sb.Append(v.options.ArrayElementSeparator)
		if v.options.BreakEnumerationItems {
			v.sb.AppendLine()
		}
	}
}
func (v *nestedStringVisitor) EndArray() {
	v.leaveEnum()
	v.sb.Append(v.options.EndArray)
}

func (v *nestedStringVisitor) enterEnum(size int) {
	v.enumLens.Push(size)
	if size >= v.options.BreakEnumerationOnLen {
		v.sb.AppendLine()
	}
}
func (v *nestedStringVisitor) leaveEnum() {
	size, _ := v.enumLens.Pop()
	if size >= v.options.BreakEnumerationOnLen {
		v.sb.AppendLine()
	}
}

func (v *nestedStringVisitor) printLinePrefix(sb *strings.Builder) {
	sb.WriteString(v.options.LinePrefix)
	for i := 0; i < v.enumLens.Count(); i++ {
		sb.WriteString(v.options.Indentation)
	}
}

func (v *nestedStringVisitor) String() string {
	return strings.TrimRight(v.sb.String(), "\n")
}

// NestedString returns a nested, JSON-like representation with configurable formatting.
func NestedString(obj interface{}, options *NestedStringOptions) (string, error) {
	if options == nil {
		options = DefaultNestedStringOptions()
	}
	visitor := &nestedStringVisitor{options: options, enumLens: newIntStack(), hierarchy: newStack()}
	visitor.sb = newStringBuilder(visitor.printLinePrefix, options.BaseTypeOptions)
	if err := visit(obj, visitor); err != nil {
		return "", err
	}
	return visitor.String(), nil
}
