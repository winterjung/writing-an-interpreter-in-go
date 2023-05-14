package evaluator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go-interpreter/lexer"
	"go-interpreter/object"
	"go-interpreter/parser"
)

func TestEvalInteger(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected int64
	}{
		{input: "0", expected: 0},
		{input: "5", expected: 5},
		{input: "-42", expected: -42},
		// TODO: --는 에러가 되어야함
		//{input: "--42", expected: 42},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertInteger(t, evaluated, tc.expected)
		})
	}
}

func TestEvalBoolean(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected bool
	}{
		{input: "false", expected: false},
		{input: "true", expected: true},
		{input: "!false", expected: true},
		{input: "!true", expected: false},
		{input: "!!false", expected: false},
		// TODO: !5, !null 문법은 지원하지 않음. 에러 검증 케이스 추가
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertBoolean(t, evaluated, tc.expected)
		})
	}
}

func evalFromString(t *testing.T, input string) object.Object {
	t.Helper()

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	require.NotNil(t, program)
	if p.Errs.ErrorOrNil() != nil {
		t.Error(p.Errs.Error())
		t.FailNow()
	}
	return Eval(program)
}

func assertInteger(t *testing.T, obj object.Object, expected int64) {
	t.Helper()

	i, ok := obj.(*object.Integer)
	require.Truef(t, ok, "expected: *object.Integer, got: %T", obj)
	require.Equal(t, expected, i.Value)
}

func assertBoolean(t *testing.T, obj object.Object, expected bool) {
	t.Helper()

	b, ok := obj.(*object.Boolean)
	require.Truef(t, ok, "expected: *object.Boolean, got: %T", obj)
	require.Equal(t, expected, b.Value)
}
