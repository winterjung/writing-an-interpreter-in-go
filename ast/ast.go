package ast

import (
	"bytes"
	"fmt"
	"strings"

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

type BlockStatement struct {
	Token      token.Token // token.LBRACE 토큰
	Statements []Statement
}

func (s *BlockStatement) statementNode() {}

func (s *BlockStatement) TokenLiteral() string { return s.Token.Literal }

func (s *BlockStatement) String() string {
	var out bytes.Buffer
	for _, stmt := range s.Statements {
		_, _ = out.WriteString(stmt.String())
	}
	return out.String()
}

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

type StringLiteral struct {
	Token token.Token // token.STRING 토큰
	Value string
}

func (l *StringLiteral) expressionNode() {}

func (l *StringLiteral) TokenLiteral() string { return l.Token.Literal }

func (l *StringLiteral) String() string { return l.Token.Literal }

// <prefix operator><expression>
type PrefixExpression struct {
	Token    token.Token // 전위 연산자 토큰 (e.g. -, !)
	Operator string
	Right    Expression
}

// [<comma separated expressions>]
type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET 토큰
	Elements []Expression
}

func (l *ArrayLiteral) expressionNode() {}

func (l *ArrayLiteral) TokenLiteral() string { return l.Token.Literal }

func (l *ArrayLiteral) String() string {
	elems := make([]string, len(l.Elements))
	for i, e := range l.Elements {
		elems[i] = e.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
}

func (exp *PrefixExpression) expressionNode() {}

func (exp *PrefixExpression) TokenLiteral() string { return exp.Token.Literal }

func (exp *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", exp.Operator, exp.Right)
}

// <expression> <infix operator> <expression>
type InfixExpression struct {
	Token    token.Token // 중위 연산자 토큰 (e.g. +, -)
	Operator string
	Left     Expression
	Right    Expression
}

func (exp *InfixExpression) expressionNode() {}

func (exp *InfixExpression) TokenLiteral() string { return exp.Token.Literal }

func (exp *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", exp.Left, exp.Operator, exp.Right)
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

func (b *Boolean) String() string { return b.Token.Literal }

// if (<condition>) <consequence> else <alternative>
type IfExpression struct {
	Token       token.Token // token.IF 토큰
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (exp *IfExpression) expressionNode() {}

func (exp *IfExpression) TokenLiteral() string { return exp.Token.Literal }

func (exp *IfExpression) String() string {
	s := fmt.Sprintf("if %s %s", exp.Condition, exp.Consequence)
	if exp.Alternative != nil {
		s += fmt.Sprintf(" else %s", exp.Alternative)
	}
	return s
}

// fn <parameters> <block statement>
type FunctionLiteral struct {
	Token  token.Token // token.FUNCTION 토큰
	Params []*Identifier
	Body   *BlockStatement
}

func (l *FunctionLiteral) expressionNode() {}

func (l *FunctionLiteral) TokenLiteral() string { return l.Token.Literal }

func (l *FunctionLiteral) String() string {
	params := make([]string, len(l.Params))
	for i, p := range l.Params {
		params[i] = p.String()
	}

	return fmt.Sprintf(
		"%s(%s) %s",
		l.TokenLiteral(),
		strings.Join(params, ", "),
		l.Body,
	)
}

// <expression>(<comma separated expressions>)
type CallExpression struct {
	Token     token.Token // token.LPAREN 토큰
	Function  Expression  // 식별자(e.g. add(1, 2))거나 함수 리터럴(e.g. fn(x) { x; }(42))
	Arguments []Expression
}

func (exp *CallExpression) expressionNode() {}

func (exp *CallExpression) TokenLiteral() string { return exp.Token.Literal }

func (exp *CallExpression) String() string {
	args := make([]string, len(exp.Arguments))
	for i, arg := range exp.Arguments {
		args[i] = arg.String()
	}

	return fmt.Sprintf(
		"%s(%s)",
		exp.Function,
		strings.Join(args, ", "),
	)
}

// <expression>[<expression>]
type IndexExpression struct {
	Token token.Token // token.LBRACKET 토큰
	Left  Expression
	Index Expression
}

func (exp *IndexExpression) expressionNode() {}

func (exp *IndexExpression) TokenLiteral() string { return exp.Token.Literal }

func (exp *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", exp.Left, exp.Index)
}
