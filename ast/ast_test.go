package ast

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go-interpreter/token"
)

func TestProgram_String(t *testing.T) {
	p := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "x"},
					Value: "x",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "y"},
					Value: "y",
				},
			},
		},
	}
	require.Equal(t, "let x = y;", p.String())
}
