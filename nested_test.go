package vardump

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbreitf1/errors"
)

type SingleArray struct {
	ListOfValues []interface{}
}

type SingleString struct {
	StrValue string
}

type SingleInt struct {
	IntValue int
}

func singleLineTestingNestedStringOptions() *NestedStringOptions {
	return &NestedStringOptions{
		LinePrefix:  "",
		Pointer:     "*",
		Indentation: "",
		BeginStruct: "{", EndStruct: "}", StructNameFormat: "%q:", StructItemSeparator: ",",
		BeginArray: "[", EndArray: "]", ArrayIndexFormat: "", ArrayElementSeparator: ",",
		BreakEnumerationOnLen: 1000, BreakEnumerationItems: false,
		QuoteStructNames: true,
		BaseTypeOptions:  DefaultBaseTypeOptions(),
	}
}

func multiLineTestingNestedStringOptions() *NestedStringOptions {
	return &NestedStringOptions{
		LinePrefix:  "",
		Pointer:     "*",
		Indentation: "  ",
		BeginStruct: "{", EndStruct: "}", StructNameFormat: "%q:", StructItemSeparator: ",",
		BeginArray: "[", EndArray: "]", ArrayIndexFormat: "", ArrayElementSeparator: ",",
		BreakEnumerationOnLen: 1, BreakEnumerationItems: true,
		QuoteStructNames: true,
		BaseTypeOptions:  DefaultBaseTypeOptions(),
	}
}

func TestNestedString(t *testing.T) {
	obj := RawOuter{"foobar", RawNested{42, "l33t"}, []string{"bar", "foo"}}
	str, err := NestedString(obj, singleLineTestingNestedStringOptions())
	errors.AssertNil(t, err)
	assert.Equal(t, '{', rune(str[0]))
	assert.Contains(t, str, `"RawValue":"foobar"`)
	assert.Contains(t, str, `"Nested":{`)
	assert.Contains(t, str, `"Integer":42`)
	assert.Contains(t, str, `"StrData":"l33t"`)
	assert.Contains(t, str, `"ListOfStrings":["bar","foo"]`)
	assert.NotContains(t, str, "\n")
	assert.Equal(t, '}', rune(str[len(str)-1]))
}

func TestNestedStringIndentation(t *testing.T) {
	obj := SingleArray{ListOfValues: []interface{}{SingleString{"foobar"}, SingleInt{42}}}
	str, err := NestedString(obj, multiLineTestingNestedStringOptions())
	errors.AssertNil(t, err)
	assert.Equal(t, "{\n  \"ListOfValues\":[\n    {\n      \"StrValue\":\"foobar\"\n    },\n    {\n      \"IntValue\":42\n    }\n  ]\n}", str)
}

func TestNestedStringBreakOnLen(t *testing.T) {
	obj := SingleArray{ListOfValues: []interface{}{SingleString{"foobar"}, SingleInt{42}}}
	options := multiLineTestingNestedStringOptions()
	options.BreakEnumerationOnLen = 2
	str, err := NestedString(obj, options)
	errors.AssertNil(t, err)
	assert.Equal(t, "{\"ListOfValues\":[\n    {\"StrValue\":\"foobar\"},\n    {\"IntValue\":42}\n  ]}", str)
}
