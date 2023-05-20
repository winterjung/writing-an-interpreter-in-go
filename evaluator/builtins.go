package evaluator

import "go-interpreter/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return makeError("len() takes exactly one argument: %d given", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return makeError("unsupported argument type of len(): '%s'", arg.Type())
			}
		},
	},
}
