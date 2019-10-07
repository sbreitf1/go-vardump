package vardump

import (
	"fmt"
	"strconv"
	"strings"
)

type printHandler func(sb *strings.Builder)

type stringBuilder struct {
	options           *BaseTypeOptions
	builder           strings.Builder
	requireNewLine    bool
	linePrefixPrinter printHandler
}

func newStringBuilder(linePrefixPrinter printHandler, options *BaseTypeOptions) *stringBuilder {
	if options == nil {
		options = DefaultBaseTypeOptions()
	}
	return &stringBuilder{options, strings.Builder{}, false, linePrefixPrinter}
}

// BaseTypeOptions defines options to format ordinal types.
type BaseTypeOptions struct {
	Pointer              string
	True, False          string
	QuoteStringValues    bool
	FallbackFormatString string
}

// DefaultBaseTypeOptions returns default values to format base types.
func DefaultBaseTypeOptions() *BaseTypeOptions {
	return &BaseTypeOptions{
		Pointer: "*",
		True:    "true", False: "false",
		QuoteStringValues:    true,
		FallbackFormatString: "<%T>",
	}
}

func (s *stringBuilder) AppendValue(value interface{}) {
	switch val := value.(type) {
	case bool:
		if val {
			s.Append(s.options.True)
		} else {
			s.Append(s.options.False)
		}

	case byte:
		s.Append(strconv.Itoa(int(val)))
	case int:
		s.Append(strconv.Itoa(val))
	case int64:
		s.Append(strconv.FormatInt(val, 10))

	case string:
		s.Append(condQuote(val, s.options.QuoteStringValues))
	case fmt.Stringer:
		s.Append(condQuote(val.String(), s.options.QuoteStringValues))

	default:
		s.Append(fmt.Sprintf(s.options.FallbackFormatString, val))
	}
}

func condQuote(str string, quote bool) string {
	if quote {
		return fmt.Sprintf("%q", str)
	}
	return str
}

func (s *stringBuilder) Append(str string) {
	if s.builder.Len() > 0 && s.requireNewLine {
		s.requireNewLine = false
		s.builder.WriteString("\n")
		if s.linePrefixPrinter != nil {
			s.linePrefixPrinter(&s.builder)
		}
	} else if s.builder.Len() == 0 {
		if s.linePrefixPrinter != nil {
			s.linePrefixPrinter(&s.builder)
		}
	}
	s.builder.WriteString(str)
}
func (s *stringBuilder) AppendLine() {
	// do not print directly but only require -> prevent empty lines
	s.requireNewLine = true
}

func (s *stringBuilder) String() string {
	return s.builder.String()
}
