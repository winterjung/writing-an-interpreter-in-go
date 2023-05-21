# writing-an-interpreter-in-go

- [《밑바닥부터 만드는 인터프리터 in Go》 책 소개](https://blog.insightbook.co.kr/2021/08/19/《밑바닥부터-만드는-인터프리터-in-go》/)

```shell
$ brew install go-task
$ task lint test format
# Run REPL
$ task run
```

## Chapters

- [x] 1장 렉싱
- [x] 2장 파싱: 프렛 파서(pratt parser)
- [x] 3장 평가: 트리 순회 인터프리터(tree-walking interpreter)
- [x] 4장 인터프리터 확장

## Showcase

```shell
>>> let age = 1;
>>> let name = "me"
>>> let result = 10 * (2 + 3)
>>> result
50
>>> print(name, age)
me 1
>>> let arr = [1, 2, 3]
>>> arr[1]
2
>>> arr[-1]
3
>>> {"name": "me", "age": 20}["name"]
me
>>> let adder = fn(n) { fn(x) { x + n } }
>>> adder(5)(10)
15
>>> let max = fn(x, y) { if (x > y) { x } else { y } }
>>> max(-1, 4)
4
```

## How to add syntax

각 단계는 필요하지 않으면 건너뛸 수 있음

1. [`token.go`](./token/token.go)에 렉싱할 토큰 추가
1. [`lexer_test.go`](./lexer/lexer_test.go)에 먼저 테스트를 추가하고 [`lexer.go`](./lexer/lexer.go)에 새로운 토큰 렉싱
1. 어떤 자료구조를 지닐지 생각하고 [`ast.go`](./ast/ast.go)에 노드 추가 (보통은 표현식을 추가하게 됨)
1. [`parser_test.go`](./parser/parser_test.go)에 테스트를 추가하고 [`parser.go`](./parser/parser.go)에 적절한 파싱함수를 등록
1. [`object.go`](./object/object.go)에 타입과 표현할 객체를 추가
1. [`evaluator_test.go`](./evaluator/evaluator_test.go)에 테스트를 추가하고 [`evaluator.go`](./evaluator/evaluator.go)의 eval switch-case에 추가
