package vardump

type intStack struct {
	buffer []int
	count  int
}

func newIntStack() *intStack {
	return &intStack{make([]int, 10), 0}
}

func (s *intStack) Push(val int) {
	s.count++
	if s.count > len(s.buffer) {
		s.buffer = append(s.buffer, val)
	} else {
		s.buffer[s.count-1] = val
	}
}

func (s *intStack) Pop() (int, bool) {
	if s.count > 0 {
		s.count--
		return s.buffer[s.count], true
	}

	return 0, false
}

func (s *intStack) Count() int {
	return s.count
}
