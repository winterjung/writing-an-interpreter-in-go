package object

type Environment struct {
	env map[string]Object
	// 함수 안에서 참조할 바깥 환경
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		env:   map[string]Object{},
		outer: nil,
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	v, ok := e.env[name]
	if !ok && e.outer != nil {
		v, ok = e.outer.Get(name)
	}
	return v, ok
}

func (e *Environment) Set(name string, v Object) Object {
	e.env[name] = v
	return v
}

func (e *Environment) Extend() *Environment {
	env := NewEnvironment()
	env.outer = e
	return env
}
