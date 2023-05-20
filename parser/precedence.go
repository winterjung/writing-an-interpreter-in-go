package parser

import "go-interpreter/token"

// 연산자 우선순위
type opPrecedence int

func (p opPrecedence) String() string {
	switch p {
	case LOWEST:
		return "LOWEST(1)"
	case EQ:
		return "EQ(2)"
	case LTGT:
		return "LTGT(3)"
	case SUM:
		return "SUM(4)"
	case PRODUCT:
		return "PRODUCT(5)"
	case PREFIX:
		return "PREFIX(6)"
	case CALL:
		return "CALL(7)"
	case INDEX:
		return "INDEX(8)"
	default:
		return "UNKNOWN(0)"
	}
}

const (
	LOWEST  opPrecedence = iota + 1
	EQ                   // ==
	LTGT                 // < or >
	SUM                  // +
	PRODUCT              // *
	PREFIX               // -x or !x
	CALL                 // x()
	INDEX                // x[index]
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
		token.LPAREN:   CALL,
		token.LBRACKET: INDEX,
	}
)
