package parser

import (
	"fmt"
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
	// 모든 파싱 함수는 아래 규약을 따름
	//   1. 현재 파싱 함수와 연관된 토큰 타입이 currToken인 상태로 진입하고
	//   2. 파싱하고자 하는 표현식 타입의 마지막 토큰이 currToken이 되도록 종료함
	p.prefixParseFnMap = map[token.Type]prefixParseFn{
		token.IDENTIFIER: p.parseIdentifier,
		token.INTEGER:    p.parseIntegerLiteral,
		token.BANG:       p.parsePrefixExpression,
		token.MINUS:      p.parsePrefixExpression,
		token.TRUE:       p.parseBoolean,
		token.FALSE:      p.parseBoolean,
		token.LPAREN:     p.parseGroupedExpression,
		token.IF:         p.parseIfExpression,
		token.FUNCTION:   p.parseFunctionLiteral,
	}
	p.infixParseFnMap = map[token.Type]infixParseFn{
		token.EQ:       p.parseInfixExpression,
		token.NEQ:      p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.LPAREN:   p.parseCallExpression,
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
		program.Statements = append(program.Statements, stmt)
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
	defer untrace(trace("선언문"))

	stmt := &ast.LetStatement{
		Token: p.currToken,
		Name:  nil,
		Value: nil,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = p.parseIdentifier().(*ast.Identifier)

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// = 토큰을 지나 표현식이 있는 곳으로 진행
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	defer untrace(trace("반환문"))

	stmt := &ast.ReturnStatement{
		Token: p.currToken,
		Value: nil,
	}
	// return을 지나 표현식이 있는 곳으로 진행
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("표현식 명령문"))

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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	defer untrace(trace("블록문"))

	block := &ast.BlockStatement{
		Token:      p.currToken,
		Statements: nil,
	}

	p.nextToken()

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}
	return block
}

func (p *Parser) parseExpression(precedence opPrecedence) ast.Expression {
	defer untrace(trace(fmt.Sprintf("표현식, LBP: %s, RBP: %s", precedence.String(), p.peekPrecedence())))

	prefix := p.prefixParseFnMap[p.currToken.Type]
	if prefix == nil {
		p.errs = multierror.Append(p.errs, errors.Errorf("no prefix parse function for %s", p.currToken.Type))
		return nil
	}
	left := prefix()

	// 현재 우선순위보다 낮은 우선순위의 토큰을 만날 때까지 반복
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFnMap[p.peekToken.Type]
		if infix == nil {
			return left
		}

		p.nextToken()
		left = infix(left)
	}
	return left
}

func (p *Parser) parseIdentifier() ast.Expression {
	defer untrace(trace("식별자"))

	// nextToken()을 호출하지 않음
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("정수"))

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

func (p *Parser) parseBoolean() ast.Expression {
	defer untrace(trace("불리언"))

	return &ast.Boolean{
		Token: p.currToken,
		Value: p.currentTokenIs(token.TRUE),
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("전위 표현식"))

	exp := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Right:    nil,
	}

	// "-"나 "!" 토큰만으론 전위 표현식을 파싱할 수 없기에
	// 토큰을 소모해 다음으로 진행시킴
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace(fmt.Sprintf("중위 표현식, left: %s", left)))

	exp := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
		Right:    nil,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	defer untrace(trace(fmt.Sprintf("함수 호출 표현식, fn: %s", fn)))

	exp := &ast.CallExpression{
		Token:     p.currToken,
		Function:  fn,
		Arguments: p.parseCallArguments(),
	}
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	// 빈 파라미터
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	p.nextToken()

	args := make([]ast.Expression, 0)
	for {
		args = append(args, p.parseExpression(LOWEST))
		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken() // token.COMMA
		p.nextToken() // 다음 표현식
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	defer untrace(trace("그룹 표현식"))

	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	defer untrace(trace("조건 표현식"))

	exp := &ast.IfExpression{
		Token: p.currToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	defer untrace(trace("함수"))

	l := &ast.FunctionLiteral{
		Token:  p.currToken,
		Params: nil,
		Body:   nil,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	l.Params = p.parseFunctionParams()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	l.Body = p.parseBlockStatement()
	return l
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	// 빈 파라미터
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	p.nextToken()

	ids := make([]*ast.Identifier, 0)
	for {
		ids = append(ids, p.parseIdentifier().(*ast.Identifier))
		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken() // token.COMMA
		p.nextToken() // 다음 식별자
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return ids
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

func (p *Parser) currPrecedence() opPrecedence {
	if p, ok := precedenceMap[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() opPrecedence {
	if p, ok := precedenceMap[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// markAsError 메서드는 디버깅을 위해 파싱 과정에서 발생한 에러를 저장함
func (p *Parser) markAsError(expected token.Type) {
	p.errs = multierror.Append(p.errs, errors.Errorf("expected: %s, but got: %s", expected, p.peekToken.Type))
}
