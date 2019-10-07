package vardump

import (
	"fmt"
	"strconv"
	"strings"
)

// NestedStringOptions defines options to format nested printing.
type NestedStringOptions struct {
	Pointer                                                               string
	Indentation                                                           string
	BeginStruct, EndStruct, StructNameValueSeparator, StructItemSeparator string
	BeginArray, EndArray, ArraySeparator                                  string
	BreakEnumerationOnLen                                                 int
	BreakEnumerationItems                                                 bool
	True, False                                                           string
	QuoteStructNames, QuoteStringValues                                   bool
	FallbackFormatString                                                  string
}

// DefaultNestedStringOptions returns the default nested string options for a JSON-like representation.
func DefaultNestedStringOptions() *NestedStringOptions {
	return &NestedStringOptions{
		Pointer:     "",
		Indentation: "  ",
		BeginStruct: "{", EndStruct: "}", StructNameValueSeparator: ": ", StructItemSeparator: ",",
		BeginArray: "[", EndArray: "]", ArraySeparator: ",",
		BreakEnumerationOnLen: 1, BreakEnumerationItems: true,
		True: "true", False: "false",
		QuoteStructNames: true, QuoteStringValues: true,
		FallbackFormatString: "<%T>",
	}
}

type nestedStringVisitor struct {
	options  *NestedStringOptions
	sb       strings.Builder
	enumLens *intStack
}

func (v *nestedStringVisitor) Pointer() {
	v.append(v.options.Pointer)
}
func (v *nestedStringVisitor) Value(value interface{}) {
	switch val := value.(type) {
	case bool:
		if val {
			v.append(v.options.True)
		} else {
			v.append(v.options.False)
		}

	case byte:
		v.append(strconv.Itoa(int(val)))
	case int:
		v.append(strconv.Itoa(val))
	case int64:
		v.append(strconv.FormatInt(val, 10))

	case string:
		v.append(condQuote(val, v.options.QuoteStringValues))
	case fmt.Stringer:
		v.append(condQuote(val.String(), v.options.QuoteStringValues))

	default:
		v.append(fmt.Sprintf(v.options.FallbackFormatString, val))
	}
}
func (v *nestedStringVisitor) BeginStruct(size int) {
	v.append(v.options.BeginStruct)
	v.enumLens.Push(size)
	if size >= v.options.BreakEnumerationOnLen {
		v.appendLine()
	}
}
func (v *nestedStringVisitor) StructValueName(index int, name string) {
	if index > 0 {
		v.append(v.options.StructItemSeparator)
		if v.options.BreakEnumerationItems {
			v.appendLine()
		}
	}
	v.append(condQuote(name, v.options.QuoteStructNames))
	v.append(v.options.StructNameValueSeparator)
}
func (v *nestedStringVisitor) EndStruct() {
	size, _ := v.enumLens.Pop()
	if size >= v.options.BreakEnumerationOnLen {
		v.appendLine()
	}
	v.append(v.options.EndStruct)
}
func (v *nestedStringVisitor) BeginArray(size int) {
	v.append(v.options.BeginArray)
	v.enumLens.Push(size)
	if size >= v.options.BreakEnumerationOnLen {
		v.appendLine()
	}
}
func (v *nestedStringVisitor) ArrayValueIndex(index int) {
	if index > 0 {
		v.append(v.options.ArraySeparator)
		if v.options.BreakEnumerationItems {
			v.appendLine()
		}
	}
}
func (v *nestedStringVisitor) EndArray() {
	size, _ := v.enumLens.Pop()
	if size >= v.options.BreakEnumerationOnLen {
		v.appendLine()
	}
	v.append(v.options.EndArray)
}

func condQuote(str string, quote bool) string {
	if quote {
		return fmt.Sprintf("%q", str)
	}
	return str
}

func (v *nestedStringVisitor) append(str string) {
	v.sb.WriteString(str)
}
func (v *nestedStringVisitor) appendLine() {
	v.sb.WriteString("\n")
	for i := 0; i < v.enumLens.Count(); i++ {
		v.sb.WriteString(v.options.Indentation)
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
	visitor := &nestedStringVisitor{options: options, enumLens: newIntStack()}
	if err := visit(obj, visitor); err != nil {
		return "", err
	}
	return visitor.String(), nil
}
