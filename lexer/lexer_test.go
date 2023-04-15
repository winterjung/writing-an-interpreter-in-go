package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go-interpreter/token"
)

func TestLexer_NextToken(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "simple tokens",
			input: "=+(){},;",
			expected: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SEMICOLON, Literal: ";"},
			},
		},
		{
			name: "simple use case",
			input: `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};
let result = add(five, ten);
`,
			expected: []token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENTIFIER, Literal: "five"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INTEGER, Literal: "5"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENTIFIER, Literal: "ten"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INTEGER, Literal: "10"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENTIFIER, Literal: "add"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.FUNCTION, Literal: "fn"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENTIFIER, Literal: "x"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENTIFIER, Literal: "y"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.IDENTIFIER, Literal: "x"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.IDENTIFIER, Literal: "y"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENTIFIER, Literal: "result"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.IDENTIFIER, Literal: "add"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENTIFIER, Literal: "five"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENTIFIER, Literal: "ten"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			lexer := New(tc.input)
			for i, expectedToken := range tc.expected {
				currentToken := lexer.NextToken()
				require.Equalf(t, expectedToken.Type, currentToken.Type, "input[%d] mismatched", i)
				require.Equalf(t, expectedToken.Literal, currentToken.Literal, "input[%d] mismatched", i)
			}
		})
	}
}
