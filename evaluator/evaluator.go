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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// 명령문
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.ReturnStatement:
		v := Eval(node.Value, env)
		if isError(v) {
			return v
		}
		return &object.ReturnValue{Value: v}
	case *ast.LetStatement:
		v := Eval(node.Value, env)
		if isError(v) {
			return v
		}
		env.Set(node.Name.Value, v)
	// 표현식
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefix(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfix(node.Operator, left, right)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndex(left, index)
	case *ast.IfExpression:
		return evalIf(node, env)
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		// 인자 평가 도중 에러가 발생했다면 항상 에러만 반환됨
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return toBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elems := evalExpressions(node.Elements, env)
		// 평가 도중 에러가 발생했다면 항상 에러만 반환됨
		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{
			Params: node.Params,
			Body:   node.Body,
			Env:    env,
		}
	}
	return nil
}

func toBooleanObject(b bool) *object.Boolean {
	if b {
		return True
	}
	return False
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
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

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
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
	if left.Type() == object.StringObject && right.Type() == object.StringObject {
		return evalInfixString(op, left, right)
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
func evalInfixString(op string, left, right object.Object) object.Object {
	l, r := left.(*object.String).Value, right.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: l + r}
	default:
		return makeError("unsupported operator: '%s' %s '%s'", left.Type(), op, right.Type())
	}
}

func evalIndex(left, index object.Object) object.Object {
	if left.Type() == object.ArrayObject && index.Type() == object.IntegerObject {
		return evalArrayIndex(left, index)
	}
	return makeError("unsupported index: '%s'", left.Type())
}

func evalArrayIndex(left, index object.Object) object.Object {
	array := left.(*object.Array)
	idx := index.(*object.Integer).Value
	max := len(array.Elements) - 1
	if idx < 0 || int(idx) > max {
		return makeError("list index out of range")
	}
	return array.Elements[idx]
}

func evalIf(exp *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(exp.Condition, env)
	if isError(cond) {
		return cond
	}
	// 정확히 true인 값을 따짐
	if cond == True {
		return Eval(exp.Consequence, env)
	}
	if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
	}
	return Null
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if v, ok := env.Get(node.Value); ok {
		return v
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return makeError("undefined name: '%s'", node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	results := make([]object.Object, 0, len(exps))
	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		results = append(results, evaluated)
	}
	return results
}

func applyFunction(obj object.Object, args []object.Object) object.Object {
	switch fn := obj.(type) {
	case *object.Function:
		env := fn.Env.Extend()
		for i, p := range fn.Params {
			env.Set(p.Value, args[i])
		}

		evaluated := Eval(fn.Body, env)
		// unwrap
		if v, ok := evaluated.(*object.ReturnValue); ok {
			return v.Value
		}
		return evaluated
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return makeError("not a function: %s", obj.Type())
	}
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
