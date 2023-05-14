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
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(node.Value)}
	// 표현식
	case *ast.PrefixExpression:
		return evalPrefix(node.Operator, Eval(node.Right))
	case *ast.InfixExpression:
		return evalInfix(node.Operator, Eval(node.Left), Eval(node.Right))
	case *ast.IfExpression:
		return evalIf(node)
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

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt)
		// 최종 리턴 값을 unwrap 해 반환함
		if v, ok := result.(*object.ReturnValue); ok {
			return v.Value
		}
	}
	return result
}

func evalBlockStatements(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt)
		// 맨 바깥에서 리턴을 처리하기 위해 unwrap 하지 않고 그대로 반환함
		if result != nil && result.Type() == object.ReturnValueObject {
			return result
		}
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

func evalIf(exp *ast.IfExpression) object.Object {
	cond := Eval(exp.Condition)
	// 정확히 true인 값을 따짐
	if cond == True {
		return Eval(exp.Consequence)
	}
	if exp.Alternative != nil {
		return Eval(exp.Alternative)
	}
	return Null
}
