package evaluator

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

func TestEvalString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{input: `"hello world!"`, expected: "hello world!"},
		// TODO: \n 제대로 지원하기
		{input: `"hello\nworld!"`, expected: "hello\\nworld!"},
		{input: `"hello" + " " + "world!"`, expected: "hello world!"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertString(t, evaluated, tc.expected)
		})
	}
}

func TestEvalArray(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected []int
	}{
		{input: "[]", expected: []int{}},
		{input: "[1]", expected: []int{1}},
		{input: "[1, 2 + 3]", expected: []int{1, 5}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertArray(t, evaluated, tc.expected)
		})
	}
}

func TestEvalArrayIndex(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected any
	}{
		{input: "[1, 2][0]", expected: 1},
		{input: "[1, 2][1]", expected: 2},
		{input: "[1, 2][-0]", expected: 1},
		{input: "[1, 2, 3, 4][-1]", expected: 4},
		{input: "[1, 2, 3, 4][-2]", expected: 3},
		{input: "let i = 0; [1][i]", expected: 1},
		{input: "[1, 2][0+1]", expected: 2},
		{input: "let a = [1]; a[0];", expected: 1},
		{input: "let a = [1, 2]; a[0] + a[1];", expected: 3},
		{input: "[1, 2][2]", expected: errors.New("list index out of range")},
		{input: "[1, 2][-3]", expected: errors.New("list index out of range")},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)

			switch expected := tc.expected.(type) {
			case int:
				assertInteger(t, evaluated, int64(expected))
			case error:
				assertError(t, evaluated, expected.Error())
			}
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

func TestEvalReturn(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected int64
	}{
		{input: "return 10;", expected: 10},
		{input: "return 10; 9;", expected: 10},
		{input: "return 2 * 5; 9;", expected: 10},
		{input: "9; return 10; 9;", expected: 10},
		{input: "if (true) { if (true) { return 10; }} return 1;", expected: 10},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertInteger(t, evaluated, tc.expected)
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "5 + true",
			expected: "unsupported operator: 'int' + 'bool'",
		},
		{
			// 평가가 중단돼야함
			input:    "5 + true; 42;",
			expected: "unsupported operator: 'int' + 'bool'",
		},
		{
			input:    "false + true",
			expected: "unsupported operator: 'bool' + 'bool'",
		},
		{
			input:    "-true",
			expected: "unsupported operator: -'bool'",
		},
		{
			// 평가가 중단돼야함
			input:    "-true; 42;",
			expected: "unsupported operator: -'bool'",
		},
		{
			input:    "if (true) { true * false }",
			expected: "unsupported operator: 'bool' * 'bool'",
		},
		{
			input:    "foobar",
			expected: "undefined name: 'foobar'",
		},
		{
			input:    `"hello" - "world"`,
			expected: "unsupported operator: 'string' - 'string'",
		},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertError(t, evaluated, tc.expected)
		})
	}
}

func TestEvalLet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected int64
	}{
		{input: "let a = 5; a;", expected: 5},
		{input: "let a = 5 * 5; a;", expected: 25},
		{input: "let a = 5; let b = a; b;", expected: 5},
		{input: "let a = 5; let b = a; a + b + 5;", expected: 15},
		// TODO: 파서에서 assign 지원하기
		//{input: "let a = 5; let b = a; a = 25; b;", expected: 5},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertInteger(t, evaluated, tc.expected)
		})
	}
}

func TestEvalFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := evalFromString(t, input)
	fn, ok := evaluated.(*object.Function)
	require.Truef(t, ok, "expected: *object.Function, got: %T", evaluated)
	assert.Equal(t, 1, len(fn.Params))
	assert.Equal(t, "x", fn.Params[0].String())
	assert.Equal(t, "(x + 2)", fn.Body.String())
}

func TestEvalFunctionCall(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "명시적인 return",
			input: `
let identity = fn(x) { return x; };
identity(42)
`,
			expected: 42,
		},
		{
			name: "암묵적인 return",
			input: `
let identity = fn(x) { x };
identity(42)
`,
			expected: 42,
		},
		{
			name: "다중 파라미터",
			input: `
let add = fn(a, b, c) { return a + b + c }
add(1, 2, 3)
`,
			expected: 6,
		},
		{
			name: "if-else",
			input: `
let max = fn(x, y) {
  if (x > y) { x } else { y }
};
max(-10, 10);
`,
			expected: 10,
		},
		{
			name: "함수 전달 전 인수 평가",
			input: `
let max = fn(x, y) {
  if (x > y) { x } else { y }
};
max(-10 + 10, max(-42, 10));
`,
			expected: 10,
		},
		{
			name: "익명 함수",
			input: `
fn(x) { x }(2 * 10)
`,
			expected: 20,
		},
		{
			name: "고차 함수: 함수도 인자로",
			input: `
let addThree = fn(x) { return x + 3; }
let twice = fn(x, func) { func(func(x)) }
twice(3, addThree)
`,
			expected: 9,
		},
		{
			name: "클로져",
			input: `
let adder = fn(x) {
  return fn(y) { x + y };
}
adder(2)(5)
`,
			expected: 7,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)
			assertInteger(t, evaluated, tc.expected)
		})
	}
}

func TestEvalBuiltinFunctions(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected any
	}{
		{input: `len("")`, expected: 0},
		{input: `len("four")`, expected: 4},
		{input: `len(1)`, expected: errors.New("unsupported argument type of len(): 'int'")},
		{input: `len("one", "two")`, expected: errors.New("len() takes exactly one argument: 2 given")},
		{input: `len()`, expected: errors.New("len() takes exactly one argument: 0 given")},
		{input: `len([])`, expected: 0},
		{input: `len([1, 2])`, expected: 2},
		{input: `len([1, 2], [3, 4])`, expected: errors.New("len() takes exactly one argument: 2 given")},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			evaluated := evalFromString(t, tc.input)

			switch expected := tc.expected.(type) {
			case int:
				assertInteger(t, evaluated, int64(expected))
			case error:
				assertError(t, evaluated, expected.Error())
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
	env := object.NewEnvironment()
	return Eval(program, env)
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

func assertString(t *testing.T, obj object.Object, expected string) {
	t.Helper()

	s, ok := obj.(*object.String)
	require.Truef(t, ok, "expected: *object.String, got: %T", obj)
	require.Equal(t, expected, s.Value)
}

func assertArray(t *testing.T, obj object.Object, expected []int) {
	t.Helper()

	array, ok := obj.(*object.Array)
	require.Truef(t, ok, "expected: *object.Array, got: %T", obj)
	require.Len(t, array.Elements, len(expected))
	for i, elem := range expected {
		assertInteger(t, array.Elements[i], int64(elem))
	}
}

func assertNull(t *testing.T, obj object.Object) {
	t.Helper()

	require.Equal(t, Null, obj)
}

func assertError(t *testing.T, obj object.Object, expected string) {
	t.Helper()

	err, ok := obj.(*object.Error)
	require.Truef(t, ok, "expected: *object.Error, got: %T", obj)
	require.Equal(t, expected, err.Message)
}
