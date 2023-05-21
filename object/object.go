package object

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"

	"go-interpreter/ast"
)

type Type string

const (
	IntegerObject     Type = "int"
	BooleanObject     Type = "bool"
	StringObject      Type = "string"
	ArrayObject       Type = "array"
	HashObject        Type = "hash"
	NullObject        Type = "null"
	ReturnValueObject Type = "return value"
	ErrorObject       Type = "error"
	FunctionObject    Type = "function"
	BuiltinObject     Type = "builtin"
)

type Object interface {
	Type() Type
	String() string
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  Type
	Value uint64
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

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
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

func (b *Boolean) HashKey() HashKey {
	var v uint64
	if b.Value {
		v = 1
	}
	return HashKey{
		Type:  b.Type(),
		Value: v,
	}
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

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value)) // revive:disable-line
	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(), // TODO: hash collision 해결
	}
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() Type {
	return ArrayObject
}

func (a *Array) String() string {
	elems := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elems[i] = e.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() Type {
	return HashObject
}

func (h *Hash) String() string {
	pairs := make([]string, 0, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key, pair.Value))
	}
	sort.Strings(pairs)
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
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

type BuiltinFunc func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunc
}

func (b *Builtin) Type() Type {
	return BuiltinObject
}

func (b *Builtin) String() string {
	return "builtin function"
}
