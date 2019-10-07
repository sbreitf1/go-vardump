package vardump

type RawOuter struct {
	RawValue      string
	Nested        RawNested
	ListOfStrings []string
}

type RawNested struct {
	Integer int
	StrData string
}
