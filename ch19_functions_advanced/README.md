# Chapter 19: 함수 고급편

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch19_functions_advanced/variadic.go
go run ch19_functions_advanced/defer_example.go
go run ch19_functions_advanced/func_type.go
go run ch19_functions_advanced/closure.go
```

> **Makefile 활용**: `make run CH=ch19_functions_advanced` 또는 `make run CH=ch19_functions_advanced FILE=variadic.go`

---

Go 언어의 함수는 단순한 코드 블록 이상의 강력한 기능을 제공한다. 가변 인수, defer, 함수 타입 변수, 클로저 등 고급 기능을 통해 유연하고 표현력 높은 코드를 작성할 수 있다. 이 장에서 다루는 내용은 Go 프로그래밍에서 매우 자주 사용되는 핵심 패턴들이다.

---

## 19.1 가변 인수 함수 (... 문법)

### 가변 인수란?

함수가 **정해지지 않은 수의 인수**를 받을 수 있는 기능이다. 매개변수 타입 앞에 `...`을 붙여 선언한다. 가변 인수는 함수 내부에서 슬라이스로 취급되므로, `range`로 순회하거나 `len()`으로 길이를 확인할 수 있다.

### 기본 문법

```go
func 함수이름(매개변수 ...타입) {
    // 매개변수는 슬라이스로 전달됨
}
```

### 예시

```go
func Sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 다양한 수의 인수로 호출 가능
Sum(1, 2)           // 3
Sum(1, 2, 3, 4, 5)  // 15
Sum()               // 0
```

인수가 하나도 전달되지 않으면 가변 인수는 길이 0인 nil 슬라이스가 된다. 따라서 `Sum()`처럼 인수 없이 호출해도 안전하게 동작한다.

### 슬라이스를 가변 인수로 전달

기존 슬라이스를 가변 인수 함수에 전달하려면 `...` 연산자를 사용하여 슬라이스를 풀어야 한다.

```go
nums := []int{1, 2, 3, 4, 5}
Sum(nums...)  // 슬라이스를 풀어서 전달
```

이때 슬라이스가 복사되는 것이 아니라, 슬라이스 자체가 가변 인수로 사용된다. 따라서 함수 내에서 슬라이스의 요소를 수정하면 원본도 영향을 받을 수 있다.

### 일반 매개변수와 함께 사용

```go
// 가변 인수는 반드시 마지막에 위치해야 한다
func Printf(format string, args ...any) {
    // format은 일반 매개변수, args는 가변 인수
}

// 잘못된 예: 가변 인수 뒤에 일반 매개변수를 둘 수 없다
// func Wrong(nums ...int, name string) {}  // 컴파일 에러
```

### 가변 인수의 내부 동작

가변 인수는 컴파일러에 의해 슬라이스로 변환된다. `Sum(1, 2, 3)`은 내부적으로 `Sum([]int{1, 2, 3})`과 유사하게 처리된다. 다만, 인수를 직접 전달할 때와 `...`로 슬라이스를 전달할 때의 동작에는 미묘한 차이가 있다. 직접 전달하면 새 슬라이스가 생성되지만, `...`로 전달하면 기존 슬라이스가 그대로 사용된다.

---

## 19.2 defer 지연 실행

### defer란?

`defer` 키워드를 사용하면 함수가 **종료되기 직전에** 실행할 코드를 예약할 수 있다. 주로 리소스 정리(파일 닫기, 잠금 해제 등)에 사용한다. `defer`는 함수가 정상적으로 반환되든, 패닉이 발생하든 **반드시 실행**된다는 점이 핵심이다.

### 기본 문법

```go
func ReadFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()  // 함수 종료 시 자동으로 파일 닫기

    // 파일 처리 로직...
    return nil
}
```

리소스를 획득한 직후에 `defer`로 해제 코드를 작성하면, 이후 어떤 경로로 함수가 종료되더라도 리소스가 해제됨을 보장할 수 있다.

### defer의 실행 순서 (LIFO)

여러 defer가 있으면 **후입선출(LIFO)** 순서로 실행된다. 마지막에 등록된 defer가 먼저 실행된다. 이는 스택 자료구조와 같은 원리이다.

```go
func main() {
    defer fmt.Println("1번째 defer")
    defer fmt.Println("2번째 defer")
    defer fmt.Println("3번째 defer")
    fmt.Println("일반 코드")
}
// 출력:
// 일반 코드
// 3번째 defer
// 2번째 defer
// 1번째 defer
```

LIFO 순서인 이유는, 나중에 획득한 리소스가 먼저 해제되어야 의존성 문제가 발생하지 않기 때문이다. 예를 들어 파일을 열고 그 파일에 대한 잠금을 획득한 경우, 잠금을 먼저 해제한 후 파일을 닫아야 한다.

### defer의 활용 사례

1. **파일 닫기**: `defer f.Close()`
2. **뮤텍스 잠금 해제**: `defer mu.Unlock()`
3. **데이터베이스 연결 해제**: `defer db.Close()`
4. **패닉 복구**: `defer func() { recover() }()`
5. **실행 시간 측정**: 시작 시간 기록 후 defer로 경과 시간 출력
6. **임시 파일/디렉토리 정리**: `defer os.Remove(tmpFile)`
7. **HTTP 응답 본문 닫기**: `defer resp.Body.Close()`

### 실행 시간 측정 패턴

```go
func TrackTime(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s 실행 시간: %v\n", name, time.Since(start))
    }
}

func SomeFunction() {
    defer TrackTime("SomeFunction")()
    // ... 함수 본문
}
```

### 주의사항

```go
// 1. defer는 함수 호출 시점에 인수가 평가된다
x := 10
defer fmt.Println(x)  // 10이 출력됨 (나중에 x가 변경되어도)
x = 20

// 2. 클로저를 사용하면 최신 값을 참조할 수 있다
y := 10
defer func() {
    fmt.Println(y)  // 20이 출력됨 (y를 참조)
}()
y = 20

// 3. 루프 내에서 defer 사용 시 주의
// 함수가 끝날 때까지 리소스가 해제되지 않으므로,
// 루프에서 많은 파일을 열면 리소스 누수가 발생할 수 있다
for _, filename := range files {
    f, _ := os.Open(filename)
    defer f.Close()  // 위험! 모든 파일이 함수 종료 시까지 열려 있음
}

// 올바른 방법: 별도 함수로 분리
for _, filename := range files {
    processFile(filename)  // 각 호출마다 defer가 실행됨
}
```

---

## 19.3 함수 타입 변수 (함수를 변수에 저장)

### 함수도 값이다

Go에서 함수는 **일급 시민(first-class citizen)**이다. 변수에 저장하고, 매개변수로 전달하고, 반환값으로 사용할 수 있다. 이는 함수형 프로그래밍 스타일을 가능하게 하는 핵심 기능이다.

### 함수 타입 변수

```go
// 함수 타입 선언
var add func(int, int) int

// 함수 할당
add = func(a, b int) int {
    return a + b
}

fmt.Println(add(3, 4))  // 7
```

함수 타입 변수의 제로 값은 `nil`이다. `nil` 함수 변수를 호출하면 패닉이 발생하므로, 호출 전에 nil 체크를 하거나 반드시 할당 후 사용해야 한다.

### 함수 타입 정의

함수 시그니처가 반복되면 `type` 키워드로 함수 타입에 이름을 부여할 수 있다. 이렇게 하면 코드 가독성이 향상된다.

```go
// 타입 별칭으로 함수 타입 정의
type Operation func(int, int) int

func Apply(op Operation, a, b int) int {
    return op(a, b)
}

// 사용
result := Apply(func(a, b int) int { return a + b }, 3, 4)
```

### 함수를 매개변수로 전달 (콜백)

함수를 매개변수로 전달하면 동작을 외부에서 주입할 수 있다. 이를 **콜백(callback)** 패턴이라 한다.

```go
func Filter(nums []int, predicate func(int) bool) []int {
    result := []int{}
    for _, n := range nums {
        if predicate(n) {
            result = append(result, n)
        }
    }
    return result
}

// 짝수만 필터링
evens := Filter([]int{1,2,3,4,5}, func(n int) bool {
    return n%2 == 0
})
// [2, 4]

// 양수만 필터링
positives := Filter([]int{-3, -1, 0, 2, 5}, func(n int) bool {
    return n > 0
})
// [2, 5]
```

같은 `Filter()` 함수에 다른 조건 함수를 전달하여 다양한 필터링을 수행할 수 있다. 이것이 함수를 값으로 다루는 것의 강력함이다.

### 함수를 반환값으로 사용

함수가 함수를 반환하는 패턴은 "팩토리 함수"라고도 불린다. 설정을 캡처하여 특화된 함수를 생성할 수 있다.

```go
func Multiplier(factor int) func(int) int {
    return func(n int) int {
        return n * factor
    }
}

double := Multiplier(2)
triple := Multiplier(3)
fmt.Println(double(5))  // 10
fmt.Println(triple(5))  // 15
```

---

## 19.4 함수 리터럴 (익명 함수, 클로저)

### 익명 함수 (Anonymous Function)

이름이 없는 함수이다. 변수에 할당하거나 즉시 실행할 수 있다. 한 번만 사용되는 짧은 함수를 정의할 때 유용하다.

```go
// 변수에 할당
greet := func(name string) {
    fmt.Printf("안녕하세요, %s님!\n", name)
}
greet("김철수")

// 즉시 실행 (IIFE - Immediately Invoked Function Expression)
func(msg string) {
    fmt.Println(msg)
}("즉시 실행!")
```

IIFE 패턴은 새로운 스코프를 만들어야 할 때나, 고루틴에서 변수를 캡처할 때 유용하다.

### 클로저 (Closure)

클로저는 **외부 변수를 캡처**하는 함수이다. 함수가 정의된 환경(렉시컬 스코프)의 변수에 접근할 수 있으며, 해당 변수는 클로저가 살아있는 동안 메모리에서 해제되지 않는다.

```go
func Counter() func() int {
    count := 0  // 외부 변수
    return func() int {
        count++      // 외부 변수를 캡처하여 사용
        return count
    }
}

counter := Counter()
fmt.Println(counter())  // 1
fmt.Println(counter())  // 2
fmt.Println(counter())  // 3
// count 변수가 클로저에 의해 유지됨
```

`Counter()` 함수가 반환된 후에도 `count` 변수는 사라지지 않는다. 클로저가 이 변수를 참조하고 있기 때문에 Go의 가비지 컬렉터가 해제하지 않는다. 이를 "변수가 힙으로 탈출(escape to heap)했다"고 한다.

### 클로저의 활용

```go
// 1. 설정을 기억하는 함수
func Logger(prefix string) func(string) {
    return func(msg string) {
        fmt.Printf("[%s] %s\n", prefix, msg)
    }
}

info := Logger("INFO")
errLog := Logger("ERROR")
info("서버 시작")       // [INFO] 서버 시작
errLog("연결 실패")     // [ERROR] 연결 실패

// 2. 누적 계산
func Accumulator(initial int) func(int) int {
    sum := initial
    return func(n int) int {
        sum += n
        return sum
    }
}

acc := Accumulator(100)
fmt.Println(acc(10))  // 110
fmt.Println(acc(20))  // 130
fmt.Println(acc(30))  // 160
```

클로저는 상태를 캡슐화하는 간결한 방법이다. 구조체와 메서드를 사용하는 것의 경량 대안이라고 볼 수 있다.

### 주의: 클로저와 반복문

클로저가 반복문의 변수를 캡처할 때 흔히 발생하는 실수가 있다. 클로저는 변수의 값이 아니라 **변수 자체(참조)**를 캡처하기 때문이다.

```go
// 잘못된 예 (Go 1.22 이전)
funcs := make([]func(), 5)
for i := 0; i < 5; i++ {
    funcs[i] = func() {
        fmt.Println(i)  // 모두 5를 출력! (i를 참조)
    }
}

// 올바른 예
for i := 0; i < 5; i++ {
    i := i  // 새 변수로 캡처 (shadowing)
    funcs[i] = func() {
        fmt.Println(i)  // 0, 1, 2, 3, 4 각각 출력
    }
}
```

> **참고**: Go 1.22부터는 `for` 루프 변수가 각 반복마다 새로운 변수로 생성되므로, 위의 "잘못된 예"도 올바르게 동작한다. 그러나 하위 호환성을 위해 shadowing 패턴을 알아두는 것이 좋다.

---

## 핵심 요약

| 개념 | 설명 |
|------|------|
| 가변 인수 | `...타입`으로 선언하며, 슬라이스로 전달된다 |
| defer | 함수 종료 직전에 실행되며, LIFO 순서이다. 리소스 정리에 사용한다 |
| 함수 타입 변수 | 함수를 변수에 저장하고, 매개변수/반환값으로 사용할 수 있다 |
| 익명 함수 | 이름 없는 함수이다. 즉시 실행이 가능하다 |
| 클로저 | 외부 변수를 캡처하는 함수이다. 상태 유지에 활용한다 |

---

## 연습문제

### 문제 1: 가변 인수
정수를 가변 인수로 받아 최대값, 최소값, 평균을 반환하는 함수 `Stats(nums ...int) (int, int, float64)`를 작성하라. 인수가 없으면 적절한 기본값을 반환하라.

### 문제 2: defer 활용
파일을 열고, 내용을 읽고, 닫는 시뮬레이션 함수를 작성하라. `defer`를 사용하여 파일이 반드시 닫히도록 하라. 에러 발생 시에도 파일이 닫히는지 확인하라.

### 문제 3: 함수 타입
`Map([]int, func(int) int) []int`와 `Filter([]int, func(int) bool) []int` 함수를 구현하라. 이를 사용하여 정수 슬라이스에서 짝수만 골라 제곱한 결과를 구하라.

### 문제 4: 클로저
피보나치 수열을 생성하는 클로저를 만들라. `Fibonacci() func() int` 함수를 호출할 때마다 다음 피보나치 수를 반환해야 한다.

### 문제 5: 종합
문자열 슬라이스를 받아 각 문자열에 변환 함수를 적용하는 `Transform([]string, func(string) string) []string` 함수를 만들고, 대문자 변환, 접두사 추가 등 다양한 변환을 적용해보라.

### 문제 6: defer 순서 예측
다음 코드의 출력 결과를 예측하고, 실제로 실행하여 확인하라:
```go
func mystery() {
    for i := 0; i < 4; i++ {
        defer fmt.Printf("%d ", i)
    }
    fmt.Println("시작")
}
```

### 문제 7: 함수 맵
`map[string]func(int, int) int` 타입의 맵을 만들어 사칙연산("+", "-", "*", "/")을 등록하고, 문자열 연산자를 입력받아 계산을 수행하는 간단한 계산기를 구현하라.

### 문제 8: Reduce 함수
`Reduce([]int, func(acc, val int) int, initial int) int` 함수를 구현하라. 이를 사용하여 슬라이스의 합계, 곱, 최대값을 각각 구하라.

### 문제 9: 함수 합성(Composition)
두 함수를 합성하여 새 함수를 반환하는 `Compose(f, g func(int) int) func(int) int`를 구현하라. `Compose(f, g)(x)`는 `f(g(x))`와 같아야 한다. 세 개 이상의 함수를 합성하는 `ComposeMany(funcs ...func(int) int) func(int) int`도 구현하라.

### 문제 10: 메모이제이션
계산 비용이 큰 함수의 결과를 캐싱하는 `Memoize(f func(int) int) func(int) int` 함수를 구현하라. 클로저와 맵을 사용하여, 같은 인수로 호출되면 캐시된 결과를 반환하도록 하라. 피보나치 함수에 적용하여 성능 차이를 확인하라.

---

## 구현 과제

### 과제 1: 파이프라인 빌더
정수 슬라이스를 처리하는 파이프라인을 구축하는 `Pipeline` 구조체를 구현하라. `AddStep(func([]int) []int)` 메서드로 처리 단계를 추가하고, `Execute(input []int) []int` 메서드로 모든 단계를 순서대로 실행하라. 필터링, 변환, 정렬 등 다양한 단계를 체이닝할 수 있어야 한다.

### 과제 2: 이벤트 시스템
콜백 기반의 이벤트 시스템을 구현하라. `EventEmitter` 구조체에 `On(event string, handler func(data any))` 메서드로 이벤트 핸들러를 등록하고, `Emit(event string, data any)` 메서드로 이벤트를 발생시켜 등록된 모든 핸들러를 호출하라. `Off(event string)` 메서드로 핸들러를 제거하는 기능도 추가하라.

### 과제 3: 재시도(Retry) 유틸리티
실패할 수 있는 함수를 자동으로 재시도하는 `Retry(maxAttempts int, delay time.Duration, fn func() error) error` 함수를 구현하라. 지수 백오프(exponential backoff) 전략을 적용하여 재시도 간격을 점진적으로 늘리는 `RetryWithBackoff` 버전도 구현하라.

### 과제 4: 미들웨어 패턴
HTTP 핸들러를 시뮬레이션하는 미들웨어 패턴을 구현하라. `type Handler func(request string) string` 타입을 정의하고, `type Middleware func(Handler) Handler` 타입의 미들웨어를 만들라. 로깅, 인증 체크, 실행 시간 측정 미들웨어를 구현하고, 이들을 체이닝하여 적용하라.

### 과제 5: 지연 평가 시퀀스
클로저를 활용하여 지연 평가(lazy evaluation) 시퀀스를 구현하라. `type Sequence func() (int, Sequence)` 타입을 정의하여 호출할 때마다 다음 값과 다음 시퀀스를 반환하도록 하라. 자연수 시퀀스, 짝수 시퀀스, 피보나치 시퀀스를 생성하는 함수를 만들고, `Take(seq Sequence, n int) []int`, `Map(seq Sequence, f func(int) int) Sequence`, `Filter(seq Sequence, f func(int) bool) Sequence` 등의 유틸리티 함수도 구현하라.

---

## 프로젝트 과제

### 프로젝트 1: 함수형 데이터 처리 라이브러리
Go에서 함수형 프로그래밍 스타일의 데이터 처리 라이브러리를 구현하라. 다음 기능을 포함해야 한다:
- `Map`, `Filter`, `Reduce`, `ForEach` 함수
- `Compose`, `Pipe` 함수 합성 유틸리티
- `Curry`(커링) 함수: `func(a, b int) int`를 `func(a int) func(b int) int`로 변환
- `Partial`(부분 적용): 일부 인수를 미리 바인딩
- `Memoize`: 결과 캐싱
- 정수, 문자열, `any` 타입 각각에 대한 버전을 구현하라

이 라이브러리를 사용하여 실제 데이터(학생 성적, 상품 목록 등)를 처리하는 예제를 작성하라.

### 프로젝트 2: 작업 스케줄러
`defer`와 클로저를 활용하여 간단한 작업 스케줄러를 구현하라. 다음 기능을 포함해야 한다:
- 작업 등록: `Schedule(name string, task func()) error`
- 지연 실행: `ScheduleAfter(name string, delay time.Duration, task func())`
- 반복 실행: `ScheduleEvery(name string, interval time.Duration, task func())`
- 작업 취소: `Cancel(name string)`
- 작업 상태 조회: `Status(name string) string`
- 모든 작업 정리: `Shutdown()` (defer 패턴 활용)
- 작업 실행 결과 로깅: 각 작업의 시작/종료 시간, 소요 시간, 성공/실패 기록

클로저를 사용하여 각 작업의 상태를 캡슐화하고, defer를 사용하여 리소스 정리를 보장하라.
