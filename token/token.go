package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // 알 수 없는 토큰
	EOF     = "EOF"     // 파일의 끝

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

var keywords = map[string]Type{
	"fn":  FUNCTION,
	"let": LET,
}

// 주어진 식별자가 예약어인지 아닌지 판단
func LookupIdentifier(s string) Type {
	if tokenType, ok := keywords[s]; ok {
		return tokenType
	}
	return IDENTIFIER
}
