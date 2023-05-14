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
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return True
		}
		return False
	}
	return nil
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
