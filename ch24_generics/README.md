# 24장 제네릭 프로그래밍

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch24_generics/basic.go
go run ch24_generics/constraints.go
go run ch24_generics/generic_type.go
go run ch24_generics/stdlib.go
```

> **Makefile 활용**: `make run CH=ch24_generics` 또는 `make run CH=ch24_generics FILE=basic.go`

---

Go 1.18부터 도입된 **제네릭(Generics)**을 사용하면 타입에 독립적인 코드를 작성할 수 있다. 제네릭은 Go 커뮤니티에서 가장 오랫동안 요청된 기능 중 하나로, 코드 중복을 줄이면서도 타입 안전성을 유지하는 것을 목표로 한다. 이 장에서는 제네릭의 기본 개념부터 실전 활용까지 학습한다.

---

## 24.1 제네릭 프로그래밍 소개

### 제네릭이 없던 시절의 문제

제네릭이 없으면 타입마다 같은 로직을 반복 작성해야 했다. 이는 코드 중복을 유발하고, 새로운 타입이 추가될 때마다 함수를 추가해야 하는 유지보수 부담을 만들었다.

```go
// 제네릭 없이 - 타입마다 별도 함수 필요
func MaxInt(a, b int) int {
    if a > b { return a }
    return b
}

func MaxFloat64(a, b float64) float64 {
    if a > b { return a }
    return b
}

func MaxString(a, b string) string {
    if a > b { return a }
    return b
}
```

또는 `interface{}`를 사용하면 타입 안전성을 잃었다. 잘못된 타입이 전달되어도 컴파일 시점에 잡히지 않고, 런타임에 panic이 발생할 수 있었다.

```go
// interface{} 사용 - 컴파일 시점에 타입 검사 불가
func Max(a, b interface{}) interface{} {
    // 런타임에 타입 단언 필요... 안전하지 않음
}
```

### 제네릭으로 해결

제네릭을 사용하면 하나의 함수로 모든 비교 가능한 타입을 처리할 수 있다. 타입 안전성은 컴파일 타임에 보장되며, 런타임 오버헤드도 없다.

```go
// 하나의 함수로 모든 비교 가능한 타입을 처리
func Max[T cmp.Ordered](a, b T) T {
    if a > b { return a }
    return b
}

Max(3, 5)           // int
Max(3.14, 2.71)     // float64
Max("abc", "xyz")   // string
```

Go의 제네릭은 **컴파일 타임에 단형화(monomorphization)**와 **딕셔너리 기반 디스패치**를 혼합하여 구현된다. 컴파일러가 상황에 따라 최적의 방식을 선택하므로, 개발자가 이를 신경 쓸 필요는 없다.

---

## 24.2 제네릭 함수

### 타입 파라미터 (Type Parameter)

함수 이름 뒤에 대괄호 `[]`로 타입 파라미터를 선언한다.

```go
func Print[T any](value T) {
    fmt.Println(value)
}
```

- `T`는 타입 파라미터의 이름이다 (관례적으로 대문자 한 글자를 사용한다: T, U, V, K, E 등).
- `any`는 모든 타입을 허용하는 제약 조건이다 (`interface{}`의 별칭).
- 타입 파라미터는 함수 시그니처와 본문에서 일반 타입처럼 사용할 수 있다.

### 타입 추론

대부분의 경우 Go 컴파일러가 인자의 타입에서 타입 파라미터를 추론한다. 명시적으로 타입을 지정할 수도 있지만, 보통은 불필요하다.

```go
Print[int](42)     // 명시적 타입 지정
Print(42)           // 타입 추론 (Go가 int로 추론)
Print("hello")      // 타입 추론 (Go가 string으로 추론)
```

타입 추론이 실패하는 경우(예: 반환 타입에만 타입 파라미터가 사용될 때)에는 명시적으로 지정해야 한다.

### 여러 타입 파라미터

쉼표로 구분하여 여러 타입 파라미터를 선언할 수 있다. 각 타입 파라미터는 독립적인 제약 조건을 가질 수 있다.

```go
func Pair[K comparable, V any](key K, value V) {
    fmt.Printf("키: %v, 값: %v\n", key, value)
}
```

---

## 24.3 제약 조건 (Constraints)

제약 조건은 타입 파라미터가 만족해야 하는 조건을 정의한다. 제약 조건이 없으면 타입 파라미터에 대해 어떤 연산도 수행할 수 없으므로, 적절한 제약 조건을 지정하는 것이 중요하다.

### 내장 제약 조건

| 제약 조건 | 설명 |
|-----------|------|
| `any` | 모든 타입 (`interface{}`의 별칭) |
| `comparable` | `==`, `!=` 비교가 가능한 타입 (맵의 키로 사용 가능) |
| `cmp.Ordered` | `<`, `>`, `<=`, `>=` 비교가 가능한 타입 (정수, 실수, 문자열) |

### 인터페이스를 사용한 제약 조건

Go에서는 인터페이스를 제약 조건으로 사용한다. 유니온(`|`) 연산자로 허용할 타입을 나열할 수 있다.

```go
// 숫자 타입만 허용하는 제약 조건
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~float32 | ~float64
}

func Sum[T Number](numbers []T) T {
    var total T
    for _, n := range numbers {
        total += n
    }
    return total
}
```

### `~` (틸다) 연산자

`~int`는 **기본 타입(underlying type)이 int인 모든 타입**을 포함한다. 이를 통해 사용자 정의 타입도 제약 조건에 포함시킬 수 있다.

```go
type MyInt int     // 기본 타입이 int
type Score int     // 기본 타입이 int

// ~int는 int, MyInt, Score 모두 허용
// int는 int만 허용
```

`~`가 없으면 정확히 해당 타입만 허용되므로, 사용자 정의 타입(`type MyInt int`)이 제외된다. 대부분의 경우 `~`를 사용하는 것이 더 유연하다.

### 메서드를 포함하는 제약 조건

제약 조건에 메서드 시그니처를 포함하면, 해당 메서드를 가진 타입만 허용된다. 타입 유니온과 메서드를 동시에 사용할 수도 있다.

```go
type Stringer interface {
    String() string
}

func PrintAll[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())
    }
}
```

---

## 24.4 제네릭 타입

### 제네릭 구조체

구조체에도 타입 파라미터를 적용하여 범용 자료구조를 만들 수 있다.

```go
// 제네릭 스택 구조
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}
```

제네릭 구조체를 사용할 때는 타입 인자를 명시해야 한다: `Stack[int]`, `Stack[string]` 등. 타입 추론은 함수에서만 동작하며, 구조체 생성 시에는 명시적 타입 지정이 필요하다.

> **참고:** `var zero T`는 타입 파라미터의 제로값을 얻는 관용적 방법이다. 숫자 타입이면 0, 문자열이면 "", 포인터면 nil이 된다.

### 제네릭 인터페이스

인터페이스에도 타입 파라미터를 적용할 수 있다.

```go
type Container[T any] interface {
    Add(item T)
    Get(index int) T
    Len() int
}
```

---

## 24.5 언제 제네릭을 사용해야 하는가?

### 사용하면 좋은 경우

1. **자료구조**: 스택, 큐, 리스트, 트리, 힙 등
2. **유틸리티 함수**: Map, Filter, Reduce, Contains, GroupBy 등
3. **타입 독립적 알고리즘**: 정렬, 검색, 비교 등
4. **컨테이너 타입**: Pair, Option, Result, Either 등

### 사용하지 않는 것이 좋은 경우

1. **단순히 `interface{}`를 대체하려는 경우**: `any`로 바꾸는 것만으로는 이점이 없다.
2. **구현이 타입마다 크게 다른 경우**: 일반 인터페이스가 더 적합하다.
3. **타입이 1~2개만 필요한 경우**: 오버엔지니어링이 될 수 있다.
4. **메서드에서 타입 파라미터를 추가하려는 경우**: Go에서는 메서드에 새 타입 파라미터를 선언할 수 없다 (구조체의 타입 파라미터만 사용 가능).

> **Go 팀의 가이드라인:** "코드를 작성할 때, 타입 파라미터가 아닌 인터페이스부터 시작하라. 공통 패턴이 보일 때 제네릭을 도입하라."

---

## 24.6 제네릭을 사용해 만든 유용한 기본 패키지

### slices 패키지 (Go 1.21+)

제네릭 도입 이후 만들어진 슬라이스 유틸리티 패키지이다. 이전에는 `sort.Ints`, `sort.Strings` 등 타입별 함수를 사용해야 했지만, 이제 하나의 `slices.Sort`로 모든 정렬 가능한 타입을 처리할 수 있다.

```go
import "slices"

nums := []int{3, 1, 4, 1, 5}
slices.Sort(nums)              // [1, 1, 3, 4, 5]
slices.Contains(nums, 3)       // true
idx := slices.Index(nums, 4)   // 3
slices.Reverse(nums)           // [5, 4, 3, 1, 1]
```

### maps 패키지 (Go 1.21+)

맵에 대한 유틸리티 함수를 제공한다.

```go
import "maps"

m1 := map[string]int{"a": 1, "b": 2}
m2 := maps.Clone(m1)           // 맵 복제
keys := slices.Collect(maps.Keys(m1))    // 모든 키
values := slices.Collect(maps.Values(m1)) // 모든 값
maps.Equal(m1, m2)              // true
```

### cmp 패키지

비교 관련 유틸리티를 제공한다.

```go
import "cmp"

cmp.Compare(1, 2)    // -1
cmp.Compare(2, 2)    //  0
cmp.Compare(3, 2)    //  1
cmp.Or(0, 0, 3)      //  3 (첫 번째 0이 아닌 값)
```

`cmp.Or`는 여러 값 중 첫 번째 제로값이 아닌 값을 반환하며, 기본값 설정에 유용하다.

---

## 핵심 요약

| 개념 | 설명 |
|------|------|
| 타입 파라미터 | `func F[T constraint](v T)` 형태로 선언 |
| `any` | 모든 타입을 허용하는 제약 조건 |
| `comparable` | `==`, `!=` 비교 가능 타입 |
| `cmp.Ordered` | 비교 연산자(`<`, `>` 등) 사용 가능 타입 |
| `~` (틸다) | 기본 타입이 같은 모든 타입 포함 |
| 유니온 타입 | `int \| float64` 형태로 허용 타입 나열 |
| 타입 추론 | 대부분의 경우 컴파일러가 타입을 자동 추론 |
| 제네릭 구조체 | `type Stack[T any] struct {}` |
| `slices` 패키지 | 슬라이스 유틸리티 (Sort, Contains, Index 등) |
| `maps` 패키지 | 맵 유틸리티 (Clone, Keys, Values 등) |

---

## 연습문제

### 문제 1: 제네릭 Filter 함수
조건을 만족하는 요소만 걸러내는 `Filter` 함수를 작성하라.
```go
func Filter[T any](items []T, fn func(T) bool) []T
```
- 짝수만 걸러내기, 양수만 걸러내기 등을 테스트하라.

### 문제 2: 제네릭 Map 함수
슬라이스의 각 요소를 변환하는 `Map` 함수를 작성하라.
```go
func Map[T, U any](items []T, fn func(T) U) []U
```
- `[]int`를 `[]string`으로 변환하는 예제를 작성하라.

### 문제 3: 제네릭 Queue
FIFO(First In, First Out) 큐를 제네릭으로 구현하라.
- `Enqueue(item T)`, `Dequeue() (T, bool)`, `Len() int` 메서드를 제공하라.

### 문제 4: 제네릭 Reduce 함수
슬라이스를 하나의 값으로 축약하는 `Reduce` 함수를 작성하라.
```go
func Reduce[T, U any](items []T, initial U, fn func(U, T) U) U
```
- 합계, 곱셈, 문자열 연결 등을 테스트하라.

### 문제 5: 제네릭 Contains 함수
슬라이스에 특정 값이 포함되어 있는지 확인하는 함수를 작성하라.
```go
func Contains[T comparable](items []T, target T) bool
```
- `comparable` 제약 조건이 왜 필요한지 주석으로 설명하라.

### 문제 6: 제네릭 GroupBy 함수
슬라이스의 요소를 특정 기준으로 그룹화하는 함수를 작성하라.
```go
func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T
```
- 사람 목록을 나이대별로 그룹화하는 예제를 작성하라.
- 문자열 목록을 첫 글자별로 그룹화하는 예제를 작성하라.

### 문제 7: 제네릭 Zip 함수
두 슬라이스를 쌍(pair)으로 묶는 함수를 작성하라.
```go
type Pair[A, B any] struct {
    First  A
    Second B
}

func Zip[A, B any](as []A, bs []B) []Pair[A, B]
```
- 두 슬라이스의 길이가 다른 경우 짧은 쪽에 맞추어라.

### 문제 8: 제네릭 제약 조건 직접 정의
다음 제약 조건을 직접 정의하고, 이를 활용하는 함수를 작성하라.
- `Numeric`: 모든 정수와 실수 타입을 포함하는 제약 조건
- `Addable`: `+` 연산이 가능한 타입(숫자와 문자열)을 포함하는 제약 조건
- 각 제약 조건을 사용하는 `Average[T Numeric]`와 `Concat[T Addable]` 함수를 구현하라.

### 문제 9: 제네릭 Optional 타입
값이 있거나 없을 수 있는 `Optional[T]` 타입을 구현하라.
- `Some(value T) Optional[T]`, `None[T]() Optional[T]` 생성자를 구현하라.
- `IsPresent() bool`, `Get() (T, bool)`, `OrElse(defaultVal T) T` 메서드를 구현하라.
- `Map[T, U any](opt Optional[T], fn func(T) U) Optional[U]` 함수를 구현하라.

### 문제 10: 제네릭 Set 타입
집합(Set) 자료구조를 제네릭으로 구현하라.
- `Add(item T)`, `Remove(item T)`, `Contains(item T) bool`, `Len() int` 메서드를 구현하라.
- `Union(other *Set[T]) *Set[T]`, `Intersection(other *Set[T]) *Set[T]` 메서드를 구현하라.
- `T`에 대해 `comparable` 제약 조건을 사용하라.

---

## 구현 과제

### 과제 1: 제네릭 정렬 라이브러리
다양한 정렬 알고리즘을 제네릭으로 구현하라.
- `BubbleSort[T cmp.Ordered](items []T)`, `InsertionSort`, `QuickSort`를 구현하라.
- 사용자 정의 비교 함수를 받는 버전도 구현하라: `SortFunc[T any](items []T, less func(T, T) bool)`.
- `int`, `float64`, `string` 슬라이스에 대해 각각 테스트하라.

### 과제 2: 제네릭 연결 리스트
이중 연결 리스트(Doubly Linked List)를 제네릭으로 구현하라.
- `PushFront(item T)`, `PushBack(item T)`, `PopFront() (T, bool)`, `PopBack() (T, bool)` 메서드를 구현하라.
- `Find(item T) bool` (comparable 제약 필요), `ForEach(fn func(T))` 메서드를 구현하라.
- `Len() int`, `ToSlice() []T` 메서드를 구현하라.

### 과제 3: 제네릭 결과 타입 (Result)
함수의 성공/실패를 표현하는 `Result[T]` 타입을 구현하라.
- `Ok(value T) Result[T]`, `Err[T](err error) Result[T]` 생성자를 구현하라.
- `IsOk() bool`, `IsErr() bool`, `Unwrap() T` (실패 시 panic), `UnwrapOr(defaultVal T) T` 메서드를 구현하라.
- `Map`, `FlatMap` 함수를 구현하여 Result를 체이닝할 수 있게 하라.

### 과제 4: 제네릭 캐시
LRU(Least Recently Used) 캐시를 제네릭으로 구현하라.
- `type LRUCache[K comparable, V any] struct`를 정의하라.
- `Get(key K) (V, bool)`, `Put(key K, value V)` 메서드를 구현하라.
- 캐시 용량을 초과하면 가장 오래 사용되지 않은 항목을 제거하라.
- 이중 연결 리스트와 맵을 조합하여 O(1) 시간 복잡도를 달성하라.

### 과제 5: 제네릭 이벤트 시스템
타입 안전한 이벤트 발행/구독 시스템을 구현하라.
- `type EventBus[T any] struct`를 정의하라.
- `Subscribe(handler func(T))`, `Publish(event T)`, `Unsubscribe(handler func(T))` 메서드를 구현하라.
- 여러 타입의 이벤트(`UserCreated`, `OrderPlaced` 등)에 대해 각각 `EventBus`를 생성하여 테스트하라.

---

## 프로젝트 과제

### 프로젝트 1: 제네릭 컬렉션 라이브러리
Go의 `slices`, `maps` 패키지를 참고하여 확장된 컬렉션 라이브러리를 구현하라.
- `Filter`, `Map`, `Reduce`, `FlatMap`, `GroupBy`, `Partition`, `Zip`, `Unzip` 함수를 구현하라.
- `SortBy`, `MinBy`, `MaxBy`, `Distinct`, `Take`, `Skip`, `Chunk` 함수를 구현하라.
- 맵에 대한 유틸리티: `MapKeys`, `MapValues`, `FilterMap`, `MergeMap` 함수를 구현하라.
- 각 함수에 대해 테스트 코드를 작성하라.
- 실제 데이터(학생 성적, 상품 목록 등)를 사용하는 종합 예제를 작성하라.

### 프로젝트 2: 제네릭 데이터 파이프라인
데이터 변환 파이프라인을 제네릭으로 구현하라.
- `Pipeline[In, Out]` 타입을 정의하고, 여러 변환 단계를 체이닝할 수 있게 하라.
- `Then[Next any](fn func(Out) Next) Pipeline[In, Next]` 메서드로 단계를 추가하라.
- `Execute(input In) Out` 메서드로 파이프라인을 실행하라.
- 에러 처리를 포함하는 `Pipeline[In, Result[Out]]` 변형을 구현하라.
- CSV 파싱 → 데이터 변환 → 필터링 → 집계의 파이프라인 예제를 작성하라.
