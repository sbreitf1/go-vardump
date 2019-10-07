package vardump

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbreitf1/errors"
)

func testingFlatStringOptions() *FlatStringOptions {
	return &FlatStringOptions{
		LinePrefix:         "",
		Pointer:            "",
		NameValueSeparator: ":",
		FieldSeparator:     ".",
		ArrayIndexFormat:   "[%d]",
		BaseTypeOptions:    DefaultBaseTypeOptions(),
	}
}

func TestFlatString(t *testing.T) {
	obj := RawOuter{"foobar", RawNested{42, "l33t"}, []string{"bar", "foo"}}
	str, err := FlatString(obj, testingFlatStringOptions())
	errors.AssertNil(t, err)
	assert.Contains(t, str, "RawValue:\"foobar\"")
	assert.Contains(t, str, "Nested.Integer:42")
	assert.Contains(t, str, "Nested.StrData:\"l33t\"")
	assert.Contains(t, str, "ListOfStrings[0]:\"bar\"")
	assert.Contains(t, str, "ListOfStrings[1]:\"foo\"")
	assert.Contains(t, str, "\n")
}

type AnnotedOuter struct {
	Data1 AnnotedInner `vardump:"-"`
	Data2 AnnotedInner `vardump:"MainData"`
}

type AnnotedInner struct {
	UserName string
	Password string `vardump:",obscure"`
}

func TestFlatStringAnnotated(t *testing.T) {
	obj := AnnotedOuter{AnnotedInner{"admin", "classified"}, AnnotedInner{"guest", "secret"}}
	str, err := FlatString(obj, testingFlatStringOptions())
	errors.AssertNil(t, err)
	assert.NotContains(t, str, "Data1.UserName:")
	assert.NotContains(t, str, "admin")
	assert.NotContains(t, str, "Data1.Password:")
	assert.NotContains(t, str, "classified")
	assert.Contains(t, str, "MainData.UserName:\"guest\"")
	assert.Contains(t, str, "MainData.Password:")
	assert.NotContains(t, str, "secret")
	assert.Contains(t, str, "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b") // hex sha256 of "secret"
}
