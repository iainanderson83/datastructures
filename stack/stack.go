package stack

// Stack is an interface for a classic stack data structure
type Stack interface {
	Push(interface{})
	Pop() interface{}
	Peek() interface{}
}

var _ Stack = &ArrStack{}

type ArrStack struct {
	s []interface{}
}

func (a *ArrStack) Push(v interface{}) {
	if v == nil {
		return
	}
	a.s = append(a.s, v)
}

func (a *ArrStack) Peek() interface{} {
	if len(a.s) == 0 {
		return nil
	}
	return a.s[len(a.s)-1]
}

func (a *ArrStack) Pop() interface{} {
	v := a.Peek()
	if len(a.s) == 0 {
		return v
	}
	a.s = a.s[:len(a.s)-1]
	return v
}
