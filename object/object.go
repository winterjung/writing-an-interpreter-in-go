package object

import (
	"fmt"
	"strconv"
	"strings"

	"go-interpreter/ast"
)

type Type string

const (
	IntegerObject     Type = "int"
	BooleanObject     Type = "bool"
	StringObject      Type = "string"
	NullObject        Type = "null"
	ReturnValueObject Type = "return value"
	ErrorObject       Type = "error"
	FunctionObject    Type = "function"
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

type String struct {
	Value string
}

func (s *String) Type() Type {
	return StringObject
}

func (s *String) String() string {
	return s.Value
}

type Null struct{}

func (n *Null) Type() Type {
	return NullObject
}

func (n *Null) String() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (v *ReturnValue) Type() Type {
	return ReturnValueObject
}

func (v *ReturnValue) String() string {
	return v.Value.String()
}

type Error struct {
	// TODO: 렉서에 행과 열 추적기를 붙인 후 스택트레이스 추가
	Message string
}

func (e *Error) Type() Type {
	return ErrorObject
}

func (e *Error) String() string {
	return "Error: " + e.Message
}

type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Environment
}

func (f *Function) Type() Type {
	return FunctionObject
}

func (f *Function) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.String()
	}

	return fmt.Sprintf(
		"fn(%s) {\n%s\n}",
		strings.Join(params, ", "),
		f.Body,
	)
}
