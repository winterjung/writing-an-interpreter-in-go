package parser

import "go-interpreter/token"

// 연산자 우선순위
type opPrecedence int

const (
	LOWEST  opPrecedence = iota + 1
	EQ                   // ==
	LTGT                 // < or >
	SUM                  // +
	PRODUCT              // *
	PREFIX               // -x or !x
	CALL                 // x()
)

var (
	precedenceMap = map[token.Type]opPrecedence{
		token.EQ:       EQ,
		token.NEQ:      EQ,
		token.LT:       LTGT,
		token.GT:       LTGT,
		token.PLUS:     SUM,
		token.MINUS:    SUM,
		token.ASTERISK: PRODUCT,
		token.SLASH:    PRODUCT,
	}
)
