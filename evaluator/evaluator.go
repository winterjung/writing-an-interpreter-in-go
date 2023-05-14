package evaluator

import (
	"fmt"

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
		v := Eval(node.Value)
		if isError(v) {
			return v
		}
		return &object.ReturnValue{Value: v}
	// 표현식
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefix(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfix(node.Operator, left, right)
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
		switch result := result.(type) {
		// 최종 리턴 값을 unwrap 해 반환함
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatements(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt)
		if result != nil {
			switch result.Type() {
			// 맨 바깥에서 리턴을 처리하기 위해 unwrap 하지 않고 그대로 반환함
			case object.ReturnValueObject:
				return result
			case object.ErrorObject:
				return result
			}
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
		return makeError("unsupported operator: %s'%s'", op, right.Type())
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
		return makeError("unsupported operator: -'%s'", right.Type())
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
	return makeError("unsupported operator: '%s' %s '%s'", left.Type(), op, right.Type())
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
		return makeError("unsupported operator: '%s' %s '%s'", left.Type(), op, right.Type())
	}
}

func evalIf(exp *ast.IfExpression) object.Object {
	cond := Eval(exp.Condition)
	if isError(cond) {
		return cond
	}
	// 정확히 true인 값을 따짐
	if cond == True {
		return Eval(exp.Consequence)
	}
	if exp.Alternative != nil {
		return Eval(exp.Alternative)
	}
	return Null
}

func makeError(format string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == object.ErrorObject
}
