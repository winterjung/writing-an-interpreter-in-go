package lexer

import (
	"unicode"

	"go-interpreter/token"
)

type Lexer struct {
	input string
	// 현재 위치 (현재 문자)
	position int
	// 현재 읽는 위치 (현재 문자의 다음)
	// 다음 문자를 미리 살펴봄과 동시에 현재 문자를 보존할 수 있어야 하기에 존재
	readPosition int
	// 현재 조사하고 있는 문자
	// TODO: rune 타입으로 바꾸고 읽는 방식을 바꿔 유니코드 지원
	ch byte
}

const (
	eof = 0
)

// TODO: io.Reader와 파일 이름으로 초기화 해 토큰에 파일 이름과 행 번호를 붙여,
// 렉싱과 파싱에서 생긴 에러를 더 쉽게 추적하도록 만들기
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = token.Token{
				Type:    token.EQ,
				Literal: "==",
			}
		default:
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = token.Token{
				Type:    token.NEQ,
				Literal: "!=",
			}
		default:
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case eof:
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	case '"':
		tok = token.Token{
			Type:    token.STRING,
			Literal: l.readString(),
		}
	// 앞에서 처리되지 않으면 식별자로 처리
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok = token.Token{
				Type:    token.LookupIdentifier(literal),
				Literal: literal,
			}
			// 이미 l.readIdentifier()에서 읽을만큼 읽었기에
			// l.readChar()를 수행하지 않고 반환함
			return tok
		}
		if unicode.IsDigit(rune(l.ch)) {
			return token.Token{
				Type:    token.INTEGER,
				Literal: l.readNumber(),
			}
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}

// 렉서가 현재 보고 있는 위치를 다음으로 이동하는 메서드
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = eof
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// 현재 위치를 바꾸지 않고 다음 문자만 살펴보는 메서드
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return eof
	}
	return l.input[l.readPosition]
}

// 문자가 아닐 때 까지 글자를 읽어 문자열을 반환함
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() string {
	start := l.position + 1 // " 다음부터
	for {
		l.readChar()
		if l.ch == '\\' && l.peekChar() == '"' {
			l.readChar()
			continue
		}
		if l.ch == '"' || l.ch == eof {
			break
		}
	}

	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

// 식별자의 허용 문자 범위를 결정하는 함수
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// TODO: 정수 이외에 실수형, 16진수 등도 지원
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '_' // 100_000 형태 지원
}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
