package interpreter

type ReturnValue struct {
	Value Value
}

func (r ReturnValue) Error() string {
	return "return"
}