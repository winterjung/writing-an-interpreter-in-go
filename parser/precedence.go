package parser

// 연산자 우선순위
type opPrecedence int

const (
	LOWEST      opPrecedence = iota + 1
	EQ                       // ==
	LESSGREATER              // < or >
	SUM                      // +
	PRODUCT                  // *
	PREFIX                   // -x or !x
	CALL                     // x()
)
