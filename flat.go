package vardump

import (
	"strings"
)

// FlatStringOptions defines options to format flat printing.
type FlatStringOptions struct {
	LinePrefix       string
	Pointer          string
	FieldSeparator   string
	ArrayIndexFormat string
}

// DefaultFlatStringOptions returns the default flat string options.
func DefaultFlatStringOptions() *FlatStringOptions {
	return &FlatStringOptions{
		LinePrefix:       "",
		Pointer:          "",
		FieldSeparator:   ".",
		ArrayIndexFormat: "[%d]",
	}
}

type flatStringVisitor struct {
	options   *FlatStringOptions
	enumLens  *intStack
	hierarchy *stack
	sb        strings.Builder
}

func (v *flatStringVisitor) Pointer() {
}
func (v *flatStringVisitor) Value(value interface{}) {
}
func (v *flatStringVisitor) BeginStruct(size int) {
}
func (v *flatStringVisitor) StructValueName(index int, name string) {
}
func (v *flatStringVisitor) EndStruct() {
}
func (v *flatStringVisitor) BeginArray(size int) {
}
func (v *flatStringVisitor) ArrayValueIndex(index int) {
}
func (v *flatStringVisitor) EndArray() {
}

func (v *flatStringVisitor) String() string {
	return strings.TrimRight(v.sb.String(), "\n")
}

// FlatString returns a nested, JSON-like representation with configurable formatting.
func FlatString(obj interface{}, options *FlatStringOptions) (string, error) {
	if options == nil {
		options = DefaultFlatStringOptions()
	}
	visitor := &flatStringVisitor{options: options, enumLens: newIntStack(), hierarchy: newStack()}
	if err := visit(obj, visitor); err != nil {
		return "", err
	}
	return visitor.String(), nil
}
