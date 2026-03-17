package main

import "fmt"

// ============================================================
// 24.4 제네릭 타입 (Generic Type)
// 구조체, 인터페이스 등에 타입 파라미터를 사용할 수 있습니다.
// ============================================================

// ---- 제네릭 스택 (Stack) ----

// Stack - LIFO(Last In, First Out) 자료구조
type Stack[T any] struct {
	items []T
}

// Push - 스택에 요소를 추가한다.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop - 스택에서 요소를 꺼냅니다. 비어있으면 제로값과 false를 반환한다.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	lastIdx := len(s.items) - 1
	item := s.items[lastIdx]
	s.items = s.items[:lastIdx]
	return item, true
}

// Peek - 스택의 맨 위 요소를 확인한다. (꺼내지 않음)
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len - 스택의 크기를 반환한다.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// IsEmpty - 스택이 비어있는지 확인한다.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// ---- 제네릭 큐 (Queue) ----

// Queue - FIFO(First In, First Out) 자료구조
type Queue[T any] struct {
	items []T
}

// Enqueue - 큐에 요소를 추가한다.
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue - 큐에서 요소를 꺼냅니다.
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Len - 큐의 크기를 반환한다.
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// ---- 제네릭 Result 타입 (Option/Result 패턴) ----

// Result - 성공 값 또는 에러를 담는 컨테이너
type Result[T any] struct {
	value T
	err   error
	ok    bool
}

// Ok - 성공 결과를 생성한다.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, ok: true}
}

// Err - 에러 결과를 생성한다.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, ok: false}
}

// Unwrap - 값을 반환한다. 에러이면 패닉을 발생시킵니다.
func (r Result[T]) Unwrap() T {
	if !r.ok {
		panic(fmt.Sprintf("Result.Unwrap 실패: %v", r.err))
	}
	return r.value
}

// UnwrapOr - 값을 반환하거나, 에러이면 기본값을 반환한다.
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if !r.ok {
		return defaultVal
	}
	return r.value
}

// IsOk - 성공인지 확인한다.
func (r Result[T]) IsOk() bool {
	return r.ok
}

// ---- 제네릭 Pair 타입 ----

// Pair - 두 값을 묶는 컨테이너
type Pair[F, S any] struct {
	First  F
	Second S
}

// NewPair - Pair를 생성한다.
func NewPair[F, S any](first F, second S) Pair[F, S] {
	return Pair[F, S]{First: first, Second: second}
}

func main() {
	fmt.Println("=== 제네릭 Stack ===")

	// int 스택
	intStack := &Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)
	fmt.Printf("스택 크기: %d\n", intStack.Len())

	if top, ok := intStack.Peek(); ok {
		fmt.Println("맨 위:", top) // 30
	}

	for !intStack.IsEmpty() {
		val, _ := intStack.Pop()
		fmt.Printf("  Pop: %d\n", val) // 30, 20, 10
	}

	// string 스택
	strStack := &Stack[string]{}
	strStack.Push("Go")
	strStack.Push("Python")
	strStack.Push("Rust")

	fmt.Println("\n문자열 스택:")
	for !strStack.IsEmpty() {
		val, _ := strStack.Pop()
		fmt.Printf("  Pop: %s\n", val)
	}

	fmt.Println("\n=== 제네릭 Queue ===")

	queue := &Queue[string]{}
	queue.Enqueue("첫 번째")
	queue.Enqueue("두 번째")
	queue.Enqueue("세 번째")

	fmt.Printf("큐 크기: %d\n", queue.Len())
	for queue.Len() > 0 {
		val, _ := queue.Dequeue()
		fmt.Printf("  Dequeue: %s\n", val) // FIFO 순서
	}

	fmt.Println("\n=== 제네릭 Result 타입 ===")

	// 성공 결과
	okResult := Ok(42)
	fmt.Printf("Ok(42): 성공=%v, 값=%d\n", okResult.IsOk(), okResult.Unwrap())

	// 에러 결과
	errResult := Err[int](fmt.Errorf("값을 찾을 수 없습니다"))
	fmt.Printf("Err: 성공=%v, 기본값=%d\n", errResult.IsOk(), errResult.UnwrapOr(-1))

	// Result를 반환하는 함수
	divide := func(a, b float64) Result[float64] {
		if b == 0 {
			return Err[float64](fmt.Errorf("0으로 나눌 수 없습니다"))
		}
		return Ok(a / b)
	}

	r1 := divide(10, 3)
	fmt.Printf("10 / 3 = %.2f\n", r1.Unwrap())

	r2 := divide(10, 0)
	fmt.Printf("10 / 0 = %.2f (기본값)\n", r2.UnwrapOr(0))

	fmt.Println("\n=== 제네릭 Pair 타입 ===")

	p1 := NewPair("이름", "홍길동")
	fmt.Printf("Pair: {%s: %s}\n", p1.First, p1.Second)

	p2 := NewPair(1, true)
	fmt.Printf("Pair: {%d: %v}\n", p2.First, p2.Second)

	// 다른 타입의 Pair
	p3 := NewPair("좌표", NewPair(37.5665, 126.9780))
	fmt.Printf("Pair: {%s: (%.4f, %.4f)}\n", p3.First, p3.Second.First, p3.Second.Second)

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
