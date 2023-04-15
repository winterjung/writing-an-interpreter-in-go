package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = "ILLEGAL" // 알 수 없는 토큰
	EOF               = "EOF"     // 파일의 끝

	// 식별자 + 리터럴
	IDENTIFIER = "IDENTIFIER" // 변수 이름
	INTEGER    = "INTEGER"

	// 연산자
	ASSIGN = "="
	PLUS   = "+"

	// 구분자
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// 예약어
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
