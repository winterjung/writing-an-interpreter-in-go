package ast

import (
	"bytes"
	"fmt"

	"go-interpreter/token"
)

type Node interface {
	// TokenLiteral 메서드는 토큰에 대응하는 리터럴 값을 반환하며
	// 디버깅, 테스트 용도로만 사용함
	TokenLiteral() string
	String() string
}

// Statement 인터페이스는 5, return 5; 같은 명령문을 의미함
type Statement interface {
	Node
	// 명령문, 표현식 혼용을 방지하기 위한 더미 메서드
	statementNode()
}

// Expression 인터페이스는 add(5, 5) 같은 표현식을 의미함
type Expression interface {
	Node
	// 명령문, 표현식 혼용을 방지하기 위한 더미 메서드
	expressionNode()
}

// Program 노드는 파서가 생산하는 모든 AST의 루트 노드
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.Statements {
		_, _ = out.WriteString(stmt.String())
	}
	return out.String()
}

// let <identifier> = <expression>;
type LetStatement struct {
	Token token.Token // token.LET 토큰
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) statementNode() {}

func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }

func (s *LetStatement) String() string {
	return fmt.Sprintf("%s %s = %s;", s.TokenLiteral(), s.Name, s.Value)
}

type Identifier struct {
	Token token.Token // token.IDENTIFIER 토큰
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string {
	return i.Value
}

// return <expression>;
type ReturnStatement struct {
	Token token.Token // token.RETURN 토큰
	Value Expression
}

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }

func (s *ReturnStatement) String() string {
	return fmt.Sprintf("%s %s;", s.TokenLiteral(), s.Value)
}
