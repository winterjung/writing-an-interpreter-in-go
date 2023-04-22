package parser

import (
	"github.com/stretchr/testify/require"
	"go-interpreter/ast"
	"go-interpreter/lexer"
	"strings"
	"testing"
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
		p := New(lexer.New(input))
		program := p.ParseProgram()
		require.NotNil(t, program)
		require.Nil(t, p.errs.ErrorOrNil())
		require.Len(t, program.Statements, 3)
		for i, expected := range expectedIdentifiers {
			assertLetStatement(t, program.Statements[i], expected)
		}

	})
}

func assertLetStatement(t *testing.T, stmt ast.Statement, name string) {
	require.Equal(t, "let", stmt.TokenLiteral())
	letStmt, ok := stmt.(*ast.LetStatement)
	require.Truef(t, ok, "expected: *ast.LetStatement, got: %T", stmt)
	require.Equal(t, name, letStmt.Name.Value)
	require.Equal(t, name, letStmt.Name.TokenLiteral())
}
