package object

type Environment struct {
	env map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{env: map[string]Object{}}
}

func (e *Environment) Get(name string) (Object, bool) {
	v, ok := e.env[name]
	return v, ok
}

func (e *Environment) Set(name string, v Object) Object {
	e.env[name] = v
	return v
}
