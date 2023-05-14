package evaluator

import (
	"go-interpreter/ast"
	"go-interpreter/object"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// 명령문
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// 표현식
	case *ast.PrefixExpression:
		return evalPrefix(node.Operator, Eval(node.Right))
	case *ast.InfixExpression:
		return evalInfix(node.Operator, Eval(node.Left), Eval(node.Right))
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return toBooleanObject(node.Value)
	}
	return nil
}

func toBooleanObject(b bool) *object.Boolean {
	if b {
		return True
	}
	return False
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalPrefix(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBang(right)
	case "-":
		return evalMinus(right)
	default:
		return Null
	}
}

func evalBang(right object.Object) object.Object {
	if right == True {
		return False
	}
	return True
}

func evalMinus(right object.Object) object.Object {
	if right.Type() != object.IntegerObject {
		return Null // TODO: 에러가 돼야함
	}

	return &object.Integer{Value: -right.(*object.Integer).Value}
}

func evalInfix(op string, left, right object.Object) object.Object {
	if left.Type() == object.IntegerObject && right.Type() == object.IntegerObject {
		return evalInfixInteger(op, left, right)
	}
	switch op {
	case "==":
		return toBooleanObject(left == right)
	case "!=":
		return toBooleanObject(left != right)
	}
	return Null
}

func evalInfixInteger(op string, left, right object.Object) object.Object {
	l, r := left.(*object.Integer).Value, right.(*object.Integer).Value
	switch op {
	// 정수는 포인터 비교로 동등성을 따질 수 없기에 불 연산자보다 먼저 와야함
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		return &object.Integer{Value: l / r}
	case "<":
		return toBooleanObject(l < r)
	case ">":
		return toBooleanObject(l > r)
	case "==":
		return toBooleanObject(l == r)
	case "!=":
		return toBooleanObject(l != r)
	default:
		return Null
	}
}
