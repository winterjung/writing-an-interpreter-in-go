package parser

import (
	"strings"
	"testing"

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
		for _, stmt := range program.Statements {
			expStmt, ok := stmt.(*ast.ExpressionStatement)
			require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)
			identifier, ok := expStmt.Expression.(*ast.Identifier)
			require.Truef(t, ok, "expected: *ast.Identifier, got: %T", expStmt.Expression)
			require.Equal(t, "foobar", identifier.TokenLiteral())
			require.Equal(t, "foobar", identifier.Value)
		}
	})
	t.Run("integer expression", func(t *testing.T) {
		t.Parallel()

		input := `42;`

		program := parseProgram(t, input)
		require.Len(t, program.Statements, 1)
		for _, stmt := range program.Statements {
			expStmt, ok := stmt.(*ast.ExpressionStatement)
			require.Truef(t, ok, "expected: *ast.ExpressionStatement, got: %T", stmt)
			integer, ok := expStmt.Expression.(*ast.IntegerLiteral)
			require.Truef(t, ok, "expected: *ast.IntegerLiteral, got: %T", expStmt.Expression)
			require.Equal(t, "42", integer.TokenLiteral())
			require.Equal(t, int64(42), integer.Value)
		}
	})
}

func assertLetStatement(t *testing.T, stmt ast.Statement, name string) {
	letStmt, ok := stmt.(*ast.LetStatement)
	require.Truef(t, ok, "expected: *ast.LetStatement, got: %T", stmt)
	require.Equal(t, "let", letStmt.TokenLiteral())
	require.Equal(t, name, letStmt.Name.Value)
	require.Equal(t, name, letStmt.Name.TokenLiteral())
}

func parseProgram(t *testing.T, input string) *ast.Program {
	p := New(lexer.New(input))
	program := p.ParseProgram()
	require.NotNil(t, program)
	require.Nil(t, p.errs.ErrorOrNil())
	return program
}
