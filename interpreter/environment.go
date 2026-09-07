package interpreter

import "fmt"

type Environment struct {
    values map[string]Value
	enclosing *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]Value),
	}
}

func NewEnclosedEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]Value),
		enclosing: enclosing,
	}
}

func (e *Environment) Define(name string, value Value) {
    e.values[name] = value
}

func (e *Environment) Get(name string) (Value, bool) {
	value, exists := e.values[name]

	if exists {
		return value, true
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return Value{}, false
}

func (e *Environment) Assign(name string, value Value) error {
	currentValue, exists := e.values[name]

	if exists {
		if currentValue.Type != value.Type {
			return fmt.Errorf(
				"não é possível atribuir %s a uma variável do tipo %s",
				value.Type,
				currentValue.Type,
			)
		}

		e.values[name] = value

		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return fmt.Errorf(
		"variável '%s' não foi declarada",
		name,
	)
}