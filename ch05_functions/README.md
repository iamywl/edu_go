# Chapter 05: 함수

함수는 특정 작업을 수행하는 코드 묶음이다. Go에서 함수를 정의하고 활용하는 방법을 배운다. 함수는 프로그래밍에서 가장 중요한 추상화 도구 중 하나로, 코드를 논리적인 단위로 분리하여 재사용성, 가독성, 유지보수성을 높여 준다.

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch05_functions/basic.go
go run ch05_functions/multi_return.go
go run ch05_functions/recursion.go
```

> **Makefile 활용**: `make run CH=ch05_functions` 또는 `make run CH=ch05_functions FILE=basic.go`

---

## 5.1 함수 정의

### 기본 구조

```go
func 함수이름(매개변수 타입) 반환타입 {
    // 함수 본문
    return 반환값
}
```

Go의 함수 정의 문법은 C나 Java와 달리 반환 타입이 매개변수 뒤에 온다. 이 설계는 복잡한 함수 시그니처를 읽기 쉽게 만들어 준다.

### 다양한 함수 형태

#### 매개변수와 반환값이 없는 함수

```go
func greet() {
    fmt.Println("안녕하세요!")
}
```

반환 타입이 없는 함수는 부수 효과(side effect)를 위해 사용된다. 화면에 출력하거나, 파일에 쓰거나, 전역 상태를 변경하는 등의 작업을 수행한다.

#### 매개변수가 있는 함수

```go
func greetUser(name string) {
    fmt.Printf("안녕하세요, %s님!\n", name)
}
```

매개변수(parameter)는 함수가 외부로부터 받는 입력 값이다. 함수를 호출할 때 전달하는 실제 값은 인자(argument)라고 한다. 엄밀히 구분하면, 매개변수는 함수 정의에 선언된 변수이고, 인자는 호출 시 전달되는 값이다.

#### 매개변수와 반환값이 있는 함수

```go
func add(a int, b int) int {
    return a + b
}

// 같은 타입의 매개변수는 축약 가능하다
func add2(a, b int) int {
    return a + b
}
```

### 멀티 반환 (Multiple Return Values)

Go의 강력한 특징 중 하나이다. 함수가 **여러 값을 동시에 반환**할 수 있다. 이 기능은 Go의 에러 처리 패턴의 근간을 이룬다.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("0으로 나눌 수 없다")
    }
    return a / b, nil
}

// 호출
result, err := divide(10, 3)
if err != nil {
    fmt.Println("에러:", err)
} else {
    fmt.Printf("결과: %.2f\n", result)
}
```

Go에서는 함수가 에러를 반환할 때 `(결과, error)` 패턴을 사용하는 것이 관례이다. 에러가 없으면 `nil`을, 에러가 있으면 에러 값을 반환한다. 호출하는 쪽에서는 반드시 에러를 확인해야 한다. 이 패턴은 Go 코드 전반에 걸쳐 일관되게 사용되며, 예외(exception)를 사용하지 않는 Go만의 독특한 에러 처리 방식이다.

반환값 중 사용하지 않는 것은 블랭크 식별자(`_`)로 무시할 수 있다:

```go
result, _ := divide(10, 3)  // 에러를 무시한다 (권장하지 않는다)
```

### 이름이 붙은 반환값 (Named Return Values)

반환값에 이름을 붙여 코드의 가독성을 높일 수 있다:

```go
func rectangleInfo(width, height float64) (area, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return // 이름이 있으므로 값 생략 가능 (naked return)
}
```

이름이 붙은 반환값은 함수 시작 시 해당 타입의 기본값(zero value)으로 초기화된다. `naked return`(값 없는 return)은 짧은 함수에서는 유용하지만, 긴 함수에서는 가독성을 해칠 수 있으므로 주의해야 한다.

### 가변 인자 함수 (Variadic Function)

매개변수의 개수가 정해지지 않은 함수이다:

```go
func sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

// 호출
fmt.Println(sum(1, 2, 3))       // 6
fmt.Println(sum(1, 2, 3, 4, 5)) // 15
fmt.Println(sum())               // 0 (인자 없이 호출해도 된다)
```

가변 인자 매개변수는 함수 내부에서 슬라이스로 취급된다. 따라서 `range`를 사용하여 순회할 수 있다. 가변 인자는 반드시 매개변수 목록의 마지막에 위치해야 한다.

기존 슬라이스를 가변 인자 함수에 전달할 때는 `...` 연산자를 사용한다:

```go
nums := []int{1, 2, 3, 4, 5}
fmt.Println(sum(nums...))  // 15
```

---

## 5.2 함수를 호출하면 생기는 일

### 콜 스택 (Call Stack)

함수가 호출되면 **콜 스택(call stack)**에 새로운 **스택 프레임(stack frame)**이 쌓인다. 스택 프레임에는 함수의 매개변수, 지역 변수, 반환 주소(함수가 끝나면 돌아갈 위치) 등이 저장된다.

```
호출 순서: main() → greet() → fmt.Println()

콜 스택:
┌──────────────────┐
│ fmt.Println()    │ ← 현재 실행 중
├──────────────────┤
│ greet()          │
├──────────────────┤
│ main()           │
└──────────────────┘
```

1. `main()` 호출 → 스택에 main 프레임 추가
2. `greet()` 호출 → 스택에 greet 프레임 추가
3. `fmt.Println()` 호출 → 스택에 Println 프레임 추가
4. `fmt.Println()` 완료 → 프레임 제거
5. `greet()` 완료 → 프레임 제거
6. `main()` 완료 → 프레임 제거, 프로그램 종료

스택은 LIFO(Last In, First Out) 구조이다. 마지막에 호출된 함수가 가장 먼저 완료되고 제거된다.

### 매개변수 복사 (Call by Value)

Go에서 함수에 인자를 전달하면 **값이 복사**된다:

```go
func changeValue(x int) {
    x = 100  // 복사본을 변경한다 (원본에 영향 없다)
}

func main() {
    a := 10
    changeValue(a)
    fmt.Println(a) // 10 (변하지 않는다!)
}
```

이것이 Go의 "값에 의한 호출(Call by Value)" 방식이다. 함수에 전달되는 것은 원본 변수가 아니라 그 값의 복사본이므로, 함수 안에서 매개변수를 아무리 변경해도 호출한 쪽의 원본 변수는 영향을 받지 않는다.

원본을 변경하려면 **포인터**를 사용한다 (나중에 자세히 배운다):

```go
func changeValuePtr(x *int) {
    *x = 100  // 포인터를 통해 원본을 변경한다
}

func main() {
    a := 10
    changeValuePtr(&a)
    fmt.Println(a) // 100 (변경되었다!)
}
```

> **참고**: 슬라이스, 맵, 채널은 내부적으로 포인터를 포함하고 있어서, 값으로 전달해도 원본 데이터를 수정할 수 있다. 하지만 이것은 "참조에 의한 호출"이 아니라, "참조를 담고 있는 값의 복사"이다. 이 미묘한 차이는 나중에 슬라이스와 맵을 배울 때 더 자세히 다룬다.

---

## 5.3 함수는 왜 쓰나?

### 1. 코드 재사용

같은 작업을 반복하지 않고 함수로 한 번만 작성한다:

```go
// 함수 없이 (반복 코드)
fmt.Printf("%-10s %5d원\n", "사과", 1000)
fmt.Printf("%-10s %5d원\n", "바나나", 2000)
fmt.Printf("%-10s %5d원\n", "포도", 3000)

// 함수 사용 (재사용)
func printPrice(item string, price int) {
    fmt.Printf("%-10s %5d원\n", item, price)
}
printPrice("사과", 1000)
printPrice("바나나", 2000)
printPrice("포도", 3000)
```

DRY(Don't Repeat Yourself) 원칙에 따라, 동일한 코드가 두 번 이상 등장하면 함수로 분리하는 것이 좋다.

### 2. 가독성 향상

복잡한 로직을 의미 있는 이름의 함수로 분리하면 코드를 읽기 쉬워진다:

```go
// 가독성이 낮은 코드
if age >= 19 && hasID && !isBanned {
    // ...
}

// 가독성이 높은 코드
if canPurchase(age, hasID, isBanned) {
    // ...
}

func canPurchase(age int, hasID, isBanned bool) bool {
    return age >= 19 && hasID && !isBanned
}
```

좋은 함수 이름은 주석을 대체한다. `canPurchase`라는 이름만 보고도 이 조건이 무엇을 판별하는지 바로 이해할 수 있다.

### 3. 유지보수 용이

버그가 있을 때 함수 하나만 수정하면 해당 함수를 사용하는 모든 곳에 적용된다. 만약 함수 없이 같은 코드가 10곳에 복사되어 있다면, 10곳을 모두 찾아서 수정해야 하며, 하나라도 빠뜨리면 버그가 남게 된다.

### 4. 테스트 용이

함수 단위로 독립적인 테스트(단위 테스트)가 가능하다:

```go
func TestAdd(t *testing.T) {
    result := add(2, 3)
    if result != 5 {
        t.Errorf("add(2, 3) = %d, want 5", result)
    }
}
```

Go에는 테스트 프레임워크가 표준 라이브러리에 내장되어 있다(`testing` 패키지). 파일 이름을 `_test.go`로 끝나게 만들고, 함수 이름을 `Test`로 시작하면 `go test` 명령으로 자동 실행할 수 있다.

### 5. 추상화

함수는 구현 세부사항을 숨기고 인터페이스만 노출한다. 함수를 사용하는 쪽에서는 함수의 내부 구현을 알 필요 없이, 입력과 출력만 이해하면 된다:

```go
// 사용하는 쪽에서는 정렬 알고리즘의 세부 구현을 알 필요 없다
sort.Ints(numbers)
```

---

## 5.4 재귀 호출

함수가 자기 자신을 호출하는 것을 **재귀 호출(recursive call)**이라고 한다. 재귀는 문제를 같은 형태의 더 작은 문제로 분할하여 해결하는 기법이다.

### 팩토리얼

`n! = n * (n-1) * (n-2) * ... * 1`

```go
func factorial(n int) int {
    if n <= 1 {
        return 1  // 기저 조건 (base case): 재귀를 멈추는 조건이다
    }
    return n * factorial(n-1)  // 재귀 호출
}

// factorial(5) 실행 과정:
// factorial(5) = 5 * factorial(4)
//              = 5 * 4 * factorial(3)
//              = 5 * 4 * 3 * factorial(2)
//              = 5 * 4 * 3 * 2 * factorial(1)
//              = 5 * 4 * 3 * 2 * 1
//              = 120
```

### 피보나치 수열

`F(n) = F(n-1) + F(n-2)`, `F(0) = 0`, `F(1) = 1`

```go
func fibonacci(n int) int {
    if n <= 0 {
        return 0
    }
    if n == 1 {
        return 1
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// 0, 1, 1, 2, 3, 5, 8, 13, 21, 34, ...
```

위의 재귀 피보나치 구현은 이해하기 쉽지만, 성능이 매우 나쁘다. `fibonacci(5)`를 계산하기 위해 `fibonacci(3)`이 두 번, `fibonacci(2)`가 세 번 호출되는 등 중복 계산이 기하급수적으로 늘어난다. 이를 **지수적 시간 복잡도 O(2^n)**라고 한다.

### 재귀 호출 주의사항

1. **반드시 기저 조건(base case)이 있어야 한다** - 없으면 무한 재귀로 스택 오버플로우가 발생한다
2. **재귀 호출마다 문제의 크기가 줄어야 한다** - `factorial(n-1)`, `fibonacci(n-1)`처럼 입력이 점점 작아져야 한다
3. **깊은 재귀는 메모리를 많이 사용한다** - 각 호출마다 스택 프레임이 쌓이므로, Go의 기본 스택 크기를 초과할 수 있다 (Go는 고루틴 스택을 동적으로 확장하지만 한계가 있다)
4. **반복문으로 대체 가능하다** - 성능이 중요하면 반복문을 사용하는 것이 좋다

```go
// 반복문 버전 팩토리얼 (더 효율적이다)
func factorialLoop(n int) int {
    result := 1
    for i := 2; i <= n; i++ {
        result *= i
    }
    return result
}
```

재귀와 반복문 중 어떤 것을 선택할지는 상황에 따라 다르다. 트리 탐색, 분할 정복 등 문제 자체가 재귀적 구조를 가진 경우에는 재귀가 더 자연스럽고 읽기 쉬운 코드를 만든다. 단순 반복 작업에는 반복문이 더 효율적이다.

---

## 5.5 함수를 값으로 다루기

Go에서 함수는 **일급 시민(first-class citizen)**이다. 변수에 할당하고, 다른 함수에 인자로 전달하고, 함수에서 반환할 수 있다.

### 함수를 변수에 할당

```go
add := func(a, b int) int {
    return a + b
}
fmt.Println(add(3, 4))  // 7
```

이렇게 이름 없이 정의된 함수를 **익명 함수(anonymous function)** 또는 **함수 리터럴(function literal)**이라고 한다.

### 고차 함수 (Higher-Order Function)

함수를 매개변수로 받거나 함수를 반환하는 함수를 고차 함수라고 한다:

```go
func apply(a, b int, op func(int, int) int) int {
    return op(a, b)
}

add := func(a, b int) int { return a + b }
mul := func(a, b int) int { return a * b }

fmt.Println(apply(3, 4, add))  // 7
fmt.Println(apply(3, 4, mul))  // 12
```

### 클로저 (Closure)

익명 함수가 자신이 정의된 스코프의 변수를 캡처(참조)하는 것을 클로저라고 한다:

```go
func counter() func() int {
    count := 0
    return func() int {
        count++  // 외부 변수 count를 캡처한다
        return count
    }
}

c := counter()
fmt.Println(c())  // 1
fmt.Println(c())  // 2
fmt.Println(c())  // 3
```

클로저는 `count` 변수에 대한 참조를 유지하므로, 함수가 호출될 때마다 같은 `count` 변수를 증가시킨다. 이를 통해 함수가 상태를 유지할 수 있다.

---

## 핵심 요약

| 항목 | 설명 |
|------|------|
| `func` | 함수 정의 키워드이다 |
| 매개변수 | `func f(a int, b string)` 형태로 정의한다 |
| 반환값 | `func f() int { return 1 }` 형태로 반환한다 |
| 멀티 반환 | `func f() (int, error)` 형태로 여러 값을 반환한다 |
| 가변 인자 | `func f(nums ...int)` 형태로 가변 개수 인자를 받는다 |
| Call by Value | 인자가 복사되어 전달된다 |
| 재귀 호출 | 함수가 자기 자신을 호출하며, 반드시 기저 조건이 필요하다 |
| 일급 시민 | 함수를 변수에 할당하고, 인자로 전달하고, 반환할 수 있다 |

---

## 연습문제

1. 두 정수를 입력받아 큰 값을 반환하는 `max(a, b int) int` 함수를 작성하라.
2. 정수 슬라이스를 받아 합계와 평균을 동시에 반환하는 `sumAndAvg(nums ...int) (int, float64)` 함수를 작성하라.
3. 재귀 함수를 사용하여 정수의 각 자릿수를 더하는 `digitSum(n int) int` 함수를 작성하라 (예: `digitSum(123)` = 6).
4. 피보나치 수열을 반복문으로 구현하고, 재귀 버전과 실행 시간을 비교해 보라. (`time` 패키지의 `time.Now()`와 `time.Since()`를 활용한다.)
5. `func apply(a, b int, op func(int, int) int) int` 처럼 함수를 매개변수로 받는 고차 함수를 작성하고, 덧셈, 뺄셈, 곱셈, 나눗셈 함수를 각각 전달하여 호출해 보라.
6. 세 정수의 최대값을 반환하는 `max3(a, b, c int) int` 함수를 작성하라. 이미 만든 `max` 함수를 활용하여 구현하라.
7. 이름이 붙은 반환값(named return values)을 사용하여 원의 넓이와 둘레를 동시에 반환하는 `circleInfo(radius float64) (area, circumference float64)` 함수를 작성하라.
8. 클로저를 사용하여 호출할 때마다 누적 합계를 반환하는 `accumulator()` 함수를 작성하라. 예를 들어 `acc := accumulator()` 후 `acc(5)` → 5, `acc(3)` → 8, `acc(2)` → 10 형태로 동작해야 한다.
9. 문자열을 받아 뒤집은 문자열을 반환하는 `reverse(s string) string` 함수를 재귀로 구현하라. (힌트: 첫 문자를 떼어내고 나머지를 재귀적으로 뒤집은 후, 첫 문자를 뒤에 붙인다.)
10. 가변 인자를 받아 모든 인자 중 최솟값과 최댓값을 동시에 반환하는 `minMax(nums ...int) (int, int)` 함수를 작성하라. 인자가 없을 때는 어떻게 처리할지도 결정하라.

---

## 구현 과제

1. **재귀 거듭제곱 함수**: `power(base, exp int) int` 함수를 재귀로 구현하라. `exp`가 0이면 1을 반환하고, 그렇지 않으면 `base * power(base, exp-1)`을 반환한다. 추가로 분할 정복 방식(`exp`가 짝수면 `power(base, exp/2)`를 제곱)으로 최적화한 버전도 구현하여 비교하라.

2. **간이 계산기 (함수 버전)**: 사칙연산 각각을 별도의 함수(`add`, `subtract`, `multiply`, `divide`)로 구현하고, 연산자 문자열에 따라 적절한 함수를 호출하는 `calculate(a, b float64, op string) (float64, error)` 함수를 작성하라. 잘못된 연산자와 0으로 나누기에 대한 에러 처리를 포함하라.

3. **GCD와 LCM 계산기**: 유클리드 호제법을 재귀로 구현하여 최대공약수(GCD)를 구하는 함수를 작성하라. 이를 활용하여 최소공배수(LCM = a*b/GCD(a,b))를 구하는 함수도 작성하라. 가변 인자를 받아 여러 수의 GCD와 LCM을 계산하는 확장 버전도 만들어 보라.

4. **함수 조합기(Composer)**: 두 개의 함수 `f`와 `g`를 받아 합성 함수 `f(g(x))`를 반환하는 `compose` 함수를 작성하라. 예를 들어 `double := func(x int) int { return x * 2 }`와 `addOne := func(x int) int { return x + 1 }`을 합성하면 `doubleAndAddOne(x) = double(addOne(x))`가 되어야 한다.

5. **하노이의 탑**: 재귀를 사용하여 하노이의 탑 문제를 푸는 함수를 작성하라. 원반의 개수를 입력받고, 각 이동 단계를 "원반 1을 A에서 C로 이동" 형태로 출력하라. 총 이동 횟수도 함께 출력하라.

---

## 프로젝트 과제

1. **함수형 데이터 처리 파이프라인**: 정수 슬라이스에 대해 `filter`, `map`, `reduce` 함수를 각각 구현하라. `filter`는 조건 함수를 받아 조건을 만족하는 원소만 반환하고, `map`은 변환 함수를 받아 각 원소를 변환하고, `reduce`는 누적 함수를 받아 하나의 값으로 축약한다. 이 세 함수를 조합하여 "짝수만 골라서 각각 제곱한 뒤 모두 더하기" 같은 복합 연산을 수행하는 프로그램을 작성하라.

2. **재귀 미로 탐색기**: 2차원 배열로 표현된 미로에서 시작점에서 도착점까지의 경로를 재귀적으로 탐색하는 프로그램을 작성하라. 미로는 0(길)과 1(벽)로 구성되며, 상하좌우로 이동할 수 있다. 경로를 찾으면 미로 위에 경로를 표시하여 출력하라. 경로가 없는 경우도 처리해야 한다.
