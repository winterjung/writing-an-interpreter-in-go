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

// <expression>;
// "x + 10;"처럼 표현식 하나로만 구성되는 명령문
type ExpressionStatement struct {
	Token      token.Token // 표현식의 첫 번째 토큰
	Expression Expression
}

func (s *ExpressionStatement) statementNode() {}

func (s *ExpressionStatement) TokenLiteral() string { return s.Token.Literal }

func (s *ExpressionStatement) String() string { return s.Expression.String() }

type Identifier struct {
	Token token.Token // token.IDENTIFIER 토큰
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }

type IntegerLiteral struct {
	Token token.Token // token.INTEGER 토큰
	Value int64
}

func (l *IntegerLiteral) expressionNode() {}

func (l *IntegerLiteral) TokenLiteral() string { return l.Token.Literal }

func (l *IntegerLiteral) String() string { return l.Token.Literal }

// <prefix operator><expression>;
type PrefixExpression struct {
	Token    token.Token // 전위 연산자 토큰 (e.g. -, !)
	Operator string
	Right    Expression
}

func (exp *PrefixExpression) expressionNode() {}

func (exp *PrefixExpression) TokenLiteral() string { return exp.Token.Literal }

func (exp *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", exp.Operator, exp.Right)
}
