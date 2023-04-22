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
		require.NotNil(t, p.errs.ErrorOrNil())
		require.Equal(t, `2 errors occurred:
	* expected: IDENTIFIER, but got: INTEGER
	* expected: =, but got: INTEGER`, strings.TrimSpace(p.errs.Error()))
	})
	t.Run("let statements", func(t *testing.T) {
		t.Parallel()

		input := `
let x = 5;
let y = 10;
let foobar = 42;
`
		expectedIdentifiers := []string{"x", "y", "foobar"}
		program := parseProgram(t, input)
		require.Len(t, program.Statements, 3)
		for i, expected := range expectedIdentifiers {
			assertLetStatement(t, program.Statements[i], expected)
		}

	})
	t.Run("return statements", func(t *testing.T) {
		t.Parallel()

		input := `
return 5;
return 10;
return 42;
`
		program := parseProgram(t, input)
		require.Len(t, program.Statements, 3)
		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			require.Truef(t, ok, "expected: *ast.ReturnStatement, got: %T", stmt)
			require.Equal(t, "return", returnStmt.TokenLiteral())
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
		identifier, ok := expStmt.Expression.(*ast.Identifier)
		require.Truef(t, ok, "expected: *ast.Identifier, got: %T", expStmt.Expression)
		require.Equal(t, "foobar", identifier.TokenLiteral())
		require.Equal(t, "foobar", identifier.Value)
	})
	t.Run("integer expression", func(t *testing.T) {
		t.Parallel()

		input := `42;`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)
		integer, ok := expStmt.Expression.(*ast.IntegerLiteral)
		require.Truef(t, ok, "expected: *ast.IntegerLiteral, got: %T", expStmt.Expression)
		require.Equal(t, "42", integer.TokenLiteral())
		require.Equal(t, int64(42), integer.Value)
	})
	t.Run("prefix expression", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input           string
			expectedOp      string
			expectedInteger int64
		}{
			{
				input:           "!5;",
				expectedOp:      "!",
				expectedInteger: 5,
			},
			{
				input:           "-15;",
				expectedOp:      "-",
				expectedInteger: 15,
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
			assertIntegerLiteral(t, prefix.Right, tc.expectedInteger)
		}
	})
	t.Run("infix expression", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			input string
			left  int64
			op    string
			right int64
		}{
			{input: "5 + 5", left: 5, op: "+", right: 5},
			{input: "5 -5", left: 5, op: "-", right: 5},
			{input: "5* 5", left: 5, op: "*", right: 5},
			{input: "5 / 5;", left: 5, op: "/", right: 5},
			{input: "5 > 5", left: 5, op: ">", right: 5},
			{input: "5 < 5", left: 5, op: "<", right: 5},
			{input: "5 == 5", left: 5, op: "==", right: 5},
			{input: "5 != 5", left: 5, op: "!=", right: 5},
		}
		for _, tc := range cases {
			t.Run(tc.input, func(t *testing.T) {
				program := parseProgram(t, tc.input)
				require.Len(t, program.Statements, 1)

				stmt := program.Statements[0]
				expStmt, ok := stmt.(*ast.ExpressionStatement)
				require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)
				infix, ok := expStmt.Expression.(*ast.InfixExpression)
				require.Truef(t, ok, "expected: *ast.InfixExpression, got: %T", expStmt.Expression)
				assertIntegerLiteral(t, infix.Left, tc.left)
				require.Equal(t, tc.op, infix.Operator)
				assertIntegerLiteral(t, infix.Right, tc.right)
			})
		}
	})
	t.Run("operator precedence", func(t *testing.T) {
		cases := []struct {
			input    string
			expected string
		}{
			{input: "a+b", expected: "(a + b)"},
			{input: "a+-b", expected: "(a + (-b))"},
			{input: "-a * b", expected: "((-a) * b)"},
			{input: "!-a", expected: "(!(-a))"},
			{input: "a + b + c", expected: "((a + b) + c)"},
			{input: "a + b - c", expected: "((a + b) - c)"},
			{input: "a * b * c", expected: "((a * b) * c)"},
			{input: "a * b / c", expected: "((a * b) / c)"},
			{input: "a + b / c", expected: "(a + (b / c))"},
			{input: "a + b * c + d / e - f", expected: "(((a + (b * c)) + (d / e)) - f)"},
			{input: "3 + 4; -5 * 5", expected: "(3 + 4)((-5) * 5)"},
			{input: "5 > 4 == 3 < 4", expected: "((5 > 4) == (3 < 4))"},
			{input: "5 < 4 != 3 > 4", expected: "((5 < 4) != (3 > 4))"},
			{input: "3 + 4 * 5 == 3 * 1 + 4 * 5", expected: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		}
		for _, tc := range cases {
			program := parseProgram(t, tc.input)
			got := program.String()
			assert.Equal(t, tc.expected, got)
		}
	})
}

func parseProgram(t *testing.T, input string) *ast.Program {
	p := New(lexer.New(input))
	program := p.ParseProgram()
	require.NotNil(t, program)
	if p.errs.ErrorOrNil() != nil {
		t.Error(p.errs.Error())
		t.FailNow()
	}
	return program
}

func assertLetStatement(t *testing.T, stmt ast.Statement, name string) {
	letStmt, ok := stmt.(*ast.LetStatement)
	require.Truef(t, ok, "expected: *ast.LetStatement, got: %T", stmt)
	require.Equal(t, "let", letStmt.TokenLiteral())
	require.Equal(t, name, letStmt.Name.Value)
	require.Equal(t, name, letStmt.Name.TokenLiteral())
}

func assertIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	integer, ok := exp.(*ast.IntegerLiteral)
	require.Truef(t, ok, "expected: *ast.IntegerLiteral, got: %T", exp)
	require.Equal(t, value, integer.Value)
	require.Equal(t, strconv.FormatInt(value, 10), integer.TokenLiteral())
}
