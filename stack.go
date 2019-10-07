package vardump

type stack struct {
	buffer []interface{}
	count  int
}

func newStack() *stack {
	return &stack{make([]interface{}, 10), 0}
}

func (s *stack) Push(val interface{}) {
	s.count++
	if s.count > len(s.buffer) {
		s.buffer = append(s.buffer, val)
	} else {
		s.buffer[s.count-1] = val
	}
}

func (s *stack) Pop() (interface{}, bool) {
	if s.count > 0 {
		s.count--
		return s.buffer[s.count], true
	}

	return 0, false
}

func (s *stack) Swap(val interface{}) bool {
	if s.count > 0 {
		s.buffer[s.count-1] = val
		return true
	}

	return false
}

func (s *stack) Array() []interface{} {
	return s.buffer[0:s.count]
}

func (s *stack) Count() int {
	return s.count
}

type intStack struct {
	s *stack
}

func newIntStack() *intStack {
	return &intStack{newStack()}
}

func (s *intStack) Push(val int) {
	s.s.Push(val)
}

func (s *intStack) Pop() (int, bool) {
	val, ok := s.s.Pop()
	if ok {
		return val.(int), true
	}

	return 0, false
}

func (s *intStack) Swap(val int) bool {
	return s.s.Swap(val)
}

func (s *intStack) Array() []int {
	arr := s.s.Array()
	intArr := make([]int, len(arr))
	for i := range arr {
		intArr[i] = arr[i].(int)
	}
	return intArr
}

func (s *intStack) Count() int {
	return s.s.Count()
}

type stringStack struct {
	s *stack
}

func newStringStack() *stringStack {
	return &stringStack{newStack()}
}

func (s *stringStack) Push(val string) {
	s.s.Push(val)
}

func (s *stringStack) Pop() (string, bool) {
	val, ok := s.s.Pop()
	if ok {
		return val.(string), true
	}

	return "", false
}

func (s *stringStack) Swap(val string) bool {
	return s.s.Swap(val)
}

func (s *stringStack) Array() []string {
	arr := s.s.Array()
	strArr := make([]string, len(arr))
	for i := range arr {
		strArr[i] = arr[i].(string)
	}
	return strArr
}

func (s *stringStack) Count() int {
	return s.s.Count()
}
