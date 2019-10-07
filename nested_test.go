package vardump

import (
	"fmt"
	"testing"

	"github.com/sbreitf1/errors"
)

func defaultTestingNestedStringOptions() *NestedStringOptions {
	return &NestedStringOptions{
		Pointer:     "*",
		Indentation: "",
		BeginStruct: "{", EndStruct: "}", StructNameValueSeparator: ":", StructItemSeparator: ",",
		BeginArray: "[", EndArray: "]", ArraySeparator: ",",
		BreakEnumerationOnLen: 1000, BreakEnumerationItems: false,
		True: "true", False: "false",
		QuoteStructNames: true, QuoteStringValues: true,
		FallbackFormatString: "<%T>",
	}
}

func TestNestedString(t *testing.T) {
	obj := RawOuter{"foobar", RawNested{42, "l33t"}, []string{"bar", "foo"}}
	str, err := NestedString(obj, defaultTestingNestedStringOptions())
	errors.AssertNil(t, err)
	fmt.Println(str)
}
