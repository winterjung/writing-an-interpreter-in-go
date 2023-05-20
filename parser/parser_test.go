package parser

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-interpreter/ast"
	"go-interpreter/lexer"
)

func TestParser_ParseProgram(t *testing.T) {
	t.Parallel()

	t.Run("parsing errors", func(t *testing.T) {
		t.Parallel()

		input := `
let 5;
let x 5;
`
		p := New(lexer.New(input))
		program := p.ParseProgram()
		require.NotNil(t, program)
		require.NotNil(t, p.Errs.ErrorOrNil())
		require.Equal(t, `2 errors occurred:
	* expected: IDENTIFIER, but got: INTEGER
	* expected: =, but got: INTEGER`, strings.TrimSpace(p.Errs.Error()))
	})
	t.Run("let statements", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input string
			id    string
			v     any
		}{
			{input: "let x= 5;", id: "x", v: 5},
			{input: "let y = true", id: "y", v: true},
			{input: "let z = y;", id: "z", v: "y"},
		}

		for _, tc := range cases {
			program := parseProgram(t, tc.input)
			require.Len(t, program.Statements, 1)

			stmt := program.Statements[0]
			assertLetStatement(t, stmt, tc.id)
			v := stmt.(*ast.LetStatement).Value
			assertLiteralExpression(t, v, tc.v)
		}

	})
	t.Run("return statements", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input    string
			expected any
		}{
			{input: "return 5;", expected: 5},
			{input: "return true", expected: true},
			{input: "return y;", expected: "y"},
		}

		for _, tc := range cases {
			program := parseProgram(t, tc.input)
			require.Len(t, program.Statements, 1)

			stmt := program.Statements[0]
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			require.Truef(t, ok, "expected: *ast.ReturnStatement, got: %T", stmt)
			assert.Equal(t, "return", returnStmt.TokenLiteral())
			assertLiteralExpression(t, returnStmt.Value, tc.expected)
		}

	})
	t.Run("identifier expression", func(t *testing.T) {
		t.Parallel()

		input := `foobar;`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		assertLiteralExpression(t, expStmt.Expression, "foobar")
	})
	t.Run("integer expression", func(t *testing.T) {
		t.Parallel()

		input := `42;`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		assertLiteralExpression(t, expStmt.Expression, 42)
	})
	t.Run("string expression", func(t *testing.T) {
		t.Parallel()

		input := `"hello world";`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		assertStringLiteral(t, expStmt.Expression, "hello world")
	})
	t.Run("boolean expression", func(t *testing.T) {
		t.Parallel()

		input := `true;`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		assertLiteralExpression(t, expStmt.Expression, true)
	})
	t.Run("array expression", func(t *testing.T) {
		t.Parallel()

		input := `[1, 2*2, "3"]`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		array, ok := expStmt.Expression.(*ast.ArrayLiteral)
		require.Truef(t, ok, "expected: *ast.ArrayLiteral, got: %T", expStmt.Expression)
		require.Len(t, array.Elements, 3)
		assertIntegerLiteral(t, array.Elements[0], 1)
		assertInfixExpression(t, array.Elements[1], 2, "*", 2)
		assertStringLiteral(t, array.Elements[2], "3")
	})
	t.Run("prefix expression", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input      string
			expectedOp string
			expected   any
		}{
			{
				input:      "!5;",
				expectedOp: "!",
				expected:   5,
			},
			{
				input:      "-15;",
				expectedOp: "-",
				expected:   15,
			},
			{
				input:      "!true;",
				expectedOp: "!",
				expected:   true,
			},
			{
				input:      "!false;",
				expectedOp: "!",
				expected:   false,
			},
		}
		for _, tc := range cases {
			program := parseProgram(t, tc.input)
			require.Len(t, program.Statements, 1)

			stmt := program.Statements[0]
			expStmt, ok := stmt.(*ast.ExpressionStatement)
			require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)
			prefix, ok := expStmt.Expression.(*ast.PrefixExpression)
			require.Truef(t, ok, "expected: *ast.PrefixExpression, got: %T", expStmt.Expression)
			require.Equal(t, tc.expectedOp, prefix.Operator)
			assertLiteralExpression(t, prefix.Right, tc.expected)
		}
	})
	t.Run("infix expression", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input string
			left  any
			op    string
			right any
		}{
			{input: "5 + 5", left: 5, op: "+", right: 5},
			{input: "5 -5", left: 5, op: "-", right: 5},
			{input: "5* 5", left: 5, op: "*", right: 5},
			{input: "5 / 5;", left: 5, op: "/", right: 5},
			{input: "5 > 5", left: 5, op: ">", right: 5},
			{input: "5 < 5", left: 5, op: "<", right: 5},
			{input: "5 == 5", left: 5, op: "==", right: 5},
			{input: "5 != 5", left: 5, op: "!=", right: 5},
			{input: "true == true", left: true, op: "==", right: true},
			{input: "true != false", left: true, op: "!=", right: false},
			{input: "false == false", left: false, op: "==", right: false},
		}
		for _, tc := range cases {
			t.Run(tc.input, func(t *testing.T) {
				program := parseProgram(t, tc.input)
				require.Len(t, program.Statements, 1)

				stmt := program.Statements[0]
				expStmt, ok := stmt.(*ast.ExpressionStatement)
				require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

				assertInfixExpression(t, expStmt.Expression, tc.left, tc.op, tc.right)
			})
		}
	})
	t.Run("operator precedence", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input    string
			expected string
		}{
			{input: "a+b", expected: "(a + b)"},
			{input: "a+-b", expected: "(a + (-b))"},
			{input: "-a * b", expected: "((-a) * b)"},
			{input: "-(a * b)", expected: "(-(a * b))"},
			{input: "!-a", expected: "(!(-a))"},
			{input: "1 + 2 * 3", expected: "(1 + (2 * 3))"},
			{input: "(1 + 2) * 3", expected: "((1 + 2) * 3)"},
			{input: "a + b + c", expected: "((a + b) + c)"},
			{input: "a + b - c", expected: "((a + b) - c)"},
			{input: "a * b * c", expected: "((a * b) * c)"},
			{input: "a * b / c", expected: "((a * b) / c)"},
			{input: "a * (b / c)", expected: "(a * (b / c))"},
			{input: "a + b / c", expected: "(a + (b / c))"},
			{input: "a + b * c + d / e - f", expected: "(((a + (b * c)) + (d / e)) - f)"},
			{input: "3 + 4; -5 * 5", expected: "(3 + 4)((-5) * 5)"},
			{input: "5 > 4 == 3 < 4", expected: "((5 > 4) == (3 < 4))"},
			{input: "5 < 4 != 3 > 4", expected: "((5 < 4) != (3 > 4))"},
			{input: "3 + 4 * 5 == 3 * 1 + 4 * 5", expected: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
			{input: "true", expected: "true"},
			{input: "false;", expected: "false"},
			{input: "!(true == true)", expected: "(!(true == true))"},
			{input: "3 > 5 == false", expected: "((3 > 5) == false)"},
			{input: "3 < 5 == true", expected: "((3 < 5) == true)"},
			{input: "a + add(b * c) + d", expected: "((a + add((b * c))) + d)"},
			{input: "add(1, add(2, 3 * 4))", expected: "add(1, add(2, (3 * 4)))"},
		}
		for _, tc := range cases {
			t.Run(tc.input, func(t *testing.T) {
				program := parseProgram(t, tc.input)
				got := program.String()
				assert.Equal(t, tc.expected, got)
			})
		}
	})
	t.Run("if expression", func(t *testing.T) {
		t.Parallel()

		input := `if (x < y) { x }`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		ifExp, ok := expStmt.Expression.(*ast.IfExpression)
		require.Truef(t, ok, "expected: *ast.IfExpression, got: %T", expStmt.Expression)

		assertInfixExpression(t, ifExp.Condition, "x", "<", "y")
		require.Len(t, ifExp.Consequence.Statements, 1)
		require.Nil(t, ifExp.Alternative)

		consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", ifExp.Consequence.Statements[0])
		assertLiteralExpression(t, consequence.Expression, "x")
	})
	t.Run("if else expression", func(t *testing.T) {
		t.Parallel()

		input := `if (x < y) { x } else { y }`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		ifExp, ok := expStmt.Expression.(*ast.IfExpression)
		require.Truef(t, ok, "expected: *ast.IfExpression, got: %T", expStmt.Expression)

		assertInfixExpression(t, ifExp.Condition, "x", "<", "y")
		require.Len(t, ifExp.Consequence.Statements, 1)
		require.Len(t, ifExp.Alternative.Statements, 1)

		consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", ifExp.Consequence.Statements[0])
		assertLiteralExpression(t, consequence.Expression, "x")

		alternative, ok := ifExp.Alternative.Statements[0].(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", ifExp.Alternative.Statements[0])
		assertLiteralExpression(t, alternative.Expression, "y")
	})
	t.Run("function literal", func(t *testing.T) {
		t.Parallel()

		input := `fn(x, y) { x + y; }`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		fn, ok := expStmt.Expression.(*ast.FunctionLiteral)
		require.Truef(t, ok, "expected: *ast.FunctionLiteral, got: %T", expStmt.Expression)

		require.Len(t, fn.Params, 2)
		assertLiteralExpression(t, fn.Params[0], "x")
		assertLiteralExpression(t, fn.Params[1], "y")

		require.Len(t, fn.Body.Statements, 1)
		bodyStmt := fn.Body.Statements[0].(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", fn.Body.Statements[0])
		assertInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
	})
	t.Run("function params", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input    string
			expected []string
		}{
			{input: "fn() {}", expected: []string{}},
			{input: "fn(x) {}", expected: []string{"x"}},
			{input: "fn(x, y, z) {}", expected: []string{"x", "y", "z"}},
		}

		for _, tc := range cases {
			t.Run(tc.input, func(t *testing.T) {
				program := parseProgram(t, tc.input)
				require.Len(t, program.Statements, 1)

				fn := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral)
				require.Len(t, fn.Params, len(tc.expected))
				for i, expected := range tc.expected {
					assertLiteralExpression(t, fn.Params[i], expected)
				}
			})
		}
	})
	t.Run("call expression", func(t *testing.T) {
		t.Parallel()

		input := `add(1, 2*3, 4 + 5)`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)

		call, ok := expStmt.Expression.(*ast.CallExpression)
		require.Truef(t, ok, "expected: *ast.CallExpression, got: %T", expStmt.Expression)

		assertIdentifier(t, call.Function, "add")
		require.Len(t, call.Arguments, 3)
		assertLiteralExpression(t, call.Arguments[0], 1)
		assertInfixExpression(t, call.Arguments[1], 2, "*", 3)
		assertInfixExpression(t, call.Arguments[2], 4, "+", 5)
	})
	t.Run("call arguments", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input    string
			expected []string
		}{
			{input: "add()", expected: []string{}},
			{input: "add(1)", expected: []string{"1"}},
			{input: "add(1, 2*3, 4 + 5)", expected: []string{"1", "(2 * 3)", "(4 + 5)"}},
		}

		for _, tc := range cases {
			t.Run(tc.input, func(t *testing.T) {
				program := parseProgram(t, tc.input)
				require.Len(t, program.Statements, 1)

				call := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.CallExpression)
				require.Len(t, call.Arguments, len(tc.expected))
				for i, expected := range tc.expected {
					assert.Equal(t, call.Arguments[i].String(), expected)
				}
			})
		}
	})
}

func parseProgram(t *testing.T, input string) *ast.Program {
	t.Helper()

	p := New(lexer.New(input))
	program := p.ParseProgram()
	require.NotNil(t, program)
	if p.Errs.ErrorOrNil() != nil {
		t.Error(p.Errs.Error())
		t.FailNow()
	}
	return program
}

func assertLetStatement(t *testing.T, stmt ast.Statement, name string) {
	t.Helper()

	letStmt, ok := stmt.(*ast.LetStatement)
	require.Truef(t, ok, "expected: *ast.LetStatement, got: %T", stmt)
	require.Equal(t, "let", letStmt.TokenLiteral())
	require.Equal(t, name, letStmt.Name.Value)
	require.Equal(t, name, letStmt.Name.TokenLiteral())
}

func assertIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	t.Helper()

	integer, ok := exp.(*ast.IntegerLiteral)
	require.Truef(t, ok, "expected: *ast.IntegerLiteral, got: %T", exp)
	require.Equal(t, value, integer.Value)
	require.Equal(t, strconv.FormatInt(value, 10), integer.TokenLiteral())
}

func assertStringLiteral(t *testing.T, exp ast.Expression, value string) {
	t.Helper()

	s, ok := exp.(*ast.StringLiteral)
	require.Truef(t, ok, "expected: *ast.StringLiteral, got: %T", exp)
	require.Equal(t, value, s.Value)
	require.Equal(t, value, s.TokenLiteral())
}

func assertIdentifier(t *testing.T, exp ast.Expression, value string) {
	t.Helper()

	identifier, ok := exp.(*ast.Identifier)
	require.Truef(t, ok, "expected: *ast.Identifier, got: %T", exp)
	require.Equal(t, value, identifier.Value)
	require.Equal(t, value, identifier.TokenLiteral())
}

func assertBoolean(t *testing.T, exp ast.Expression, value bool) {
	t.Helper()

	identifier, ok := exp.(*ast.Boolean)
	require.Truef(t, ok, "expected: *ast.Boolean, got: %T", exp)
	require.Equal(t, value, identifier.Value)
	require.Equal(t, strconv.FormatBool(value), identifier.TokenLiteral())
}

func assertLiteralExpression(t *testing.T, exp ast.Expression, expected any) {
	t.Helper()

	switch x := expected.(type) {
	case int:
		assertIntegerLiteral(t, exp, int64(x))
	case int64:
		assertIntegerLiteral(t, exp, x)
	case string:
		assertIdentifier(t, exp, x)
	case bool:
		assertBoolean(t, exp, x)
	default:
		t.Errorf("unknown expression, got: %T", exp)
	}
}

func assertInfixExpression(t *testing.T, exp ast.Expression, left any, op string, right any) {
	t.Helper()

	infix, ok := exp.(*ast.InfixExpression)
	require.Truef(t, ok, "expected: *ast.InfixExpression, got: %T", exp)

	assertLiteralExpression(t, infix.Left, left)
	require.Equal(t, op, infix.Operator)
	assertLiteralExpression(t, infix.Right, right)
}
