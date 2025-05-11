package infix_to_postfix

type Stack struct {
	Top  *Element
	Size int
}

type Element struct {
	Value interface{}
	Next  *Element
}

func (s *Stack) Empty() bool {
	return s.Size == 0
}

func (s *Stack) TopFunc() interface{} {
	return s.Top.Value
}

func (s *Stack) Push(value interface{}) {
	s.Top = &Element{value, s.Top}
	s.Size++
}

func (s *Stack) Pop() (value interface{}) {
	if s.Size > 0 {
		value, s.Top = s.Top.Value, s.Top.Next
		s.Size--
		return
	}
	return nil
}
