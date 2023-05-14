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
		{input: "1 + 2 + 3", expected: 6},
		{input: "-50 + 100 - 50", expected: 0},
		{input: "-50 + 100 + -50", expected: 0},
		{input: "4 * 4", expected: 16},
		{input: "0 / 42", expected: 0},
		{input: "4 * (2 + 3)", expected: 20},
		// TODO: division by zero 검증하기
		//{input: "1 / 0", expected: 0},
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
		{input: "1 < 2", expected: true},
		{input: "1 > 2", expected: false},
		{input: "1 == 1", expected: true},
		{input: "1 != 1", expected: false},
		{input: "1 == 2", expected: false},
		{input: "1 != 2", expected: true},
		{input: "true == true", expected: true},
		{input: "true == false", expected: false},
		{input: "true != false", expected: true},
		{input: "false == false", expected: true},
		{input: "(1 > 2) == false", expected: true},
		{input: "(1 == 2) == false", expected: true},
		// TODO: 아직 null은 직접 파싱하지 않음
		//{input: "null == null", expected: true},
		//{input: "null == true", expected: false},
		//{input: "null == false", expected: false},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertBoolean(t, evaluated, tc.expected)
		})
	}
}

func TestEvalIfElse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected any
	}{
		{input: "if (true) { 10 }", expected: 10},
		{input: "if (false) { 10 }", expected: nil},
		// TODO: 42 같은 truthy한 값은 true로 판단하지 않음
		// {input: "if (1) { 10 }", expected: 10},
		// {input: "if (null) { 10 }", expected: 10},
		{input: "if (1 == 1) { 10 }", expected: 10},
		{input: "if (1 > 2) { 10 }", expected: nil},
		{input: "if (1 > 2) { 10 } else { 42 }", expected: 42},
		{input: "if (1 < 2) { 10 } else { 42 }", expected: 10},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			i, ok := tc.expected.(int)
			if ok {
				assertInteger(t, evaluated, int64(i))
			} else {
				assertNull(t, evaluated)
			}
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

func assertNull(t *testing.T, obj object.Object) {
	t.Helper()

	require.Equal(t, Null, obj)
}
