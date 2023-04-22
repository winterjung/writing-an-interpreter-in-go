package parser

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go-interpreter/ast"
	"go-interpreter/lexer"
	"go-interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	// 만약 currToken 일 때 5;로 끝나는건지 5 + 10 처럼 처리해야하는지 알기 위해 필요함
	peekToken token.Token
	// 파싱 중 발생한 에러
	errs *multierror.Error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

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
	default:
		return nil
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
