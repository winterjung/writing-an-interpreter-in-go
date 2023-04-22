package parser

import (
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"go-interpreter/ast"
	"go-interpreter/lexer"
	"go-interpreter/token"
)

type (
	// 전위 파싱 함수
	prefixParseFn func() ast.Expression
	// 중위 파싱 함수
	infixParseFn func(ast.Expression) ast.Expression
	// TODO: ++ 같은 후위 파싱 지원
)
type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	// 만약 currToken 일 때 5;로 끝나는건지 5 + 10 처럼 처리해야하는지 알기 위해 필요함
	peekToken token.Token
	// 파싱 중 발생한 에러
	errs *multierror.Error

	// 현재 토큰 토큰에 따라 사용할 수 있는 파싱 함수
	prefixParseFnMap map[token.Type]prefixParseFn
	infixParseFnMap  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixParseFnMap = map[token.Type]prefixParseFn{
		token.IDENTIFIER: p.parseIdentifier,
		token.INTEGER:    p.parseIntegerLiteral,
	}

	// currToken, peekToken 세팅
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.currToken,
		Name:  nil,
		Value: nil,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: 일단 표현식은 무시하고 세미콜론까지 진행
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.currToken,
		Value: nil,
	}
	// return을 지나 표현식이 있는 곳으로 진행
	p.nextToken()

	// TODO: 일단 표현식은 무시하고 세미콜론까지 진행
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token:      p.currToken,
		Expression: p.parseExpression(LOWEST),
	}

	// REPL에서 "5 + 5"같은 표현식을 간편하게 사용하기 위해
	// 세미콜론을 선택적으로 검사
	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence opPrecedence) ast.Expression {
	prefix := p.prefixParseFnMap[p.currToken.Type]
	if prefix == nil {
		return nil
	}
	left := prefix()
	return left
}

func (p *Parser) parseIdentifier() ast.Expression {
	// nextToken()을 호출하지 않음
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	// nextToken()을 호출하지 않음
	i, err := strconv.ParseInt(p.currToken.Literal, 10, 64)
	if err != nil {
		p.errs = multierror.Append(p.errs, errors.Errorf("could not parse %q as integer", p.currToken.Literal))
		return nil
	}
	return &ast.IntegerLiteral{
		Token: p.currToken,
		Value: i,
	}
}

func (p *Parser) currentTokenIs(t token.Type) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	// 다음으로 오길 기대하는 토큰이 맞으면 해당 토큰을 소비하고 한단계 진행
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.markAsError(t)
	return false
}

// markAsError 메서드는 디버깅을 위해 파싱 과정에서 발생한 에러를 저장함
func (p *Parser) markAsError(expected token.Type) {
	p.errs = multierror.Append(p.errs, errors.Errorf("expected: %s, but got: %s", expected, p.peekToken.Type))
}
