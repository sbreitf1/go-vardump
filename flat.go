package vardump

import (
	"fmt"
	"strings"
)

// FlatStringOptions defines options to format flat printing.
type FlatStringOptions struct {
	LinePrefix         string
	Pointer            string
	NameValueSeparator string
	FieldSeparator     string
	ArrayIndexFormat   string
	BaseTypeOptions    *BaseTypeOptions
}

// DefaultFlatStringOptions returns the default flat string options.
func DefaultFlatStringOptions() *FlatStringOptions {
	return &FlatStringOptions{
		LinePrefix:         "",
		Pointer:            "",
		NameValueSeparator: ": ",
		FieldSeparator:     ".",
		ArrayIndexFormat:   "[%d]",
		BaseTypeOptions:    DefaultBaseTypeOptions(),
	}
}

type flatStringVisitor struct {
	options   *FlatStringOptions
	hierarchy *stringStack
	sb        *stringBuilder
}

func (v *flatStringVisitor) Pointer() {
}
func (v *flatStringVisitor) Value(value interface{}, obscure bool) {
	for _, str := range v.hierarchy.Array() {
		v.sb.Append(str)
	}
	v.sb.Append(v.options.NameValueSeparator)
	v.sb.AppendValue(value, obscure)
	v.sb.AppendLine()
}
func (v *flatStringVisitor) BeginStruct(size int) {
	// a placeholder to swap in StructValueName()
	v.hierarchy.Push("")
}
func (v *flatStringVisitor) StructValueName(index int, name string) {
	if v.hierarchy.Count() > 1 {
		v.hierarchy.Swap("." + name)
	} else {
		v.hierarchy.Swap(name)
	}
}
func (v *flatStringVisitor) EndStruct() {
	v.hierarchy.Pop()
}
func (v *flatStringVisitor) BeginArray(size int) {
	// a placeholder to swap in ArrayValueIndex()
	v.hierarchy.Push("")
}
func (v *flatStringVisitor) ArrayValueIndex(index int) {
	v.hierarchy.Swap(fmt.Sprintf(v.options.ArrayIndexFormat, index))
}
func (v *flatStringVisitor) EndArray() {
	v.hierarchy.Pop()
}

func (v *flatStringVisitor) printLinePrefix(sb *strings.Builder) {
	sb.WriteString(v.options.LinePrefix)
}

func (v *flatStringVisitor) String() string {
	return v.sb.String()
}

// FlatString returns a nested, JSON-like representation with configurable formatting.
func FlatString(obj interface{}, options *FlatStringOptions) (string, error) {
	if options == nil {
		options = DefaultFlatStringOptions()
	}
	visitor := &flatStringVisitor{options: options, hierarchy: newStringStack()}
	visitor.sb = newStringBuilder(visitor.printLinePrefix, options.BaseTypeOptions)
	if err := visit(obj, visitor, ObscureByDefault); err != nil {
		return "", err
	}
	return visitor.String(), nil
}
