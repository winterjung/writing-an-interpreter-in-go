package object

import (
	"strconv"
)

type Type string

const (
	IntegerObject Type = "int"
	BooleanObject Type = "bool"
	NullObject    Type = "null"
)

type Object interface {
	Type() Type
	String() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return IntegerObject
}

func (i *Integer) String() string {
	return strconv.FormatInt(i.Value, 10)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BooleanObject
}

func (b *Boolean) String() string {
	return strconv.FormatBool(b.Value)
}

type Null struct{}

func (n *Null) Type() Type {
	return NullObject
}

func (n *Null) String() string {
	return "null"
}
