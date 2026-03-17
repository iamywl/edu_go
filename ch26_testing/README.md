# 26장 테스트와 벤치마크하기

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 전체 테스트 실행
go test ch26_testing/ -v

# 특정 테스트 함수만 실행
go test ch26_testing/ -run TestAdd -v

# 벤치마크 실행
go test ch26_testing/ -bench . -benchmem

# 테스트 커버리지 확인
go test ch26_testing/ -cover

# 커버리지 HTML 리포트 생성
go test ch26_testing/ -coverprofile=coverage.out
go tool cover -html=coverage.out
```

> **참고**: 이 장은 `go run`이 아닌 `go test` 명령어를 사용한다. 테스트 파일(`calculator_test.go`, `calculator_bench_test.go`)은 `go test`로만 실행할 수 있다.

> **Makefile 활용**: `make run CH=ch26_testing` 또는 `make run CH=ch26_testing FILE=calculator.go`

---

Go 언어는 테스트를 언어 차원에서 강력하게 지원한다. 별도의 테스트 프레임워크 없이도 `testing` 패키지와 `go test` 명령어만으로 단위 테스트, 테이블 주도 테스트, 벤치마크까지 수행할 수 있다. 테스트는 소프트웨어 품질을 보장하는 가장 기본적이고 중요한 도구이며, Go는 이를 언어 설계 단계에서부터 고려하여 매우 자연스럽고 간결한 테스트 환경을 제공한다.

---

## 26.1 테스트 코드

### testing 패키지

Go의 표준 라이브러리인 `testing` 패키지는 테스트에 필요한 모든 기능을 제공한다. 외부 라이브러리 의존 없이 단위 테스트, 서브테스트, 벤치마크, 퍼즈 테스트 등 다양한 테스트 기법을 사용할 수 있다. `testing` 패키지는 Go 도구 체인과 긴밀하게 통합되어 있어, `go test` 명령 하나로 테스트 실행, 결과 확인, 커버리지 측정까지 수행할 수 있다.

### 테스트 파일 규칙

Go 테스트는 다음 규칙을 따른다:

1. **파일 이름**: `_test.go`로 끝나야 한다 (예: `calculator_test.go`)
2. **함수 이름**: `Test`로 시작해야 한다 (예: `TestAdd`)
3. **함수 시그니처**: `func TestXxx(t *testing.T)` 형태여야 한다

`_test.go` 파일은 `go build` 시에는 포함되지 않으며, 오직 `go test`를 실행할 때만 컴파일된다. 이 규칙 덕분에 테스트 코드와 프로덕션 코드를 같은 패키지에 두면서도 바이너리 크기에 영향을 주지 않는다.

```go
// calculator_test.go
package main

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

### go test 명령어

```bash
# 현재 패키지의 모든 테스트 실행
go test

# 자세한 출력과 함께 실행
go test -v

# 특정 테스트 함수만 실행
go test -run TestAdd

# 모든 하위 패키지 테스트 실행
go test ./...

# 테스트 커버리지 확인
go test -cover

# 커버리지 프로파일 생성 후 HTML 리포트
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

`-v` 플래그는 각 테스트 함수의 이름과 실행 결과를 상세히 출력한다. `-run` 플래그는 정규식을 지원하므로 `-run "TestAdd|TestSub"`처럼 여러 테스트를 동시에 지정할 수 있다. `-cover` 플래그는 코드 커버리지를 퍼센트로 표시하여 테스트가 소스 코드를 얼마나 다루는지 확인할 수 있게 한다.

### 테스트 함수에서 사용하는 주요 메서드

| 메서드 | 설명 |
|--------|------|
| `t.Error(args...)` | 에러 메시지를 출력하고 테스트를 실패로 표시한다 (이후 코드 계속 실행) |
| `t.Errorf(format, args...)` | 포맷된 에러 메시지를 출력한다 |
| `t.Fatal(args...)` | 에러 메시지를 출력하고 즉시 해당 테스트를 중단한다 |
| `t.Fatalf(format, args...)` | 포맷된 에러 메시지를 출력 후 즉시 중단한다 |
| `t.Log(args...)` | 테스트 로그를 출력한다 (`-v` 플래그와 함께 사용) |
| `t.Skip(args...)` | 테스트를 건너뛴다 (조건부 테스트에 유용하다) |
| `t.Helper()` | 해당 함수가 헬퍼 함수임을 표시한다 (에러 위치 보고를 개선한다) |
| `t.Cleanup(func())` | 테스트 종료 시 호출될 정리 함수를 등록한다 |
| `t.Parallel()` | 다른 테스트와 병렬로 실행하도록 지정한다 |

`t.Error`와 `t.Fatal`의 차이를 이해하는 것이 중요하다. `t.Error`는 실패를 기록하지만 나머지 테스트 코드를 계속 실행하므로, 하나의 테스트 함수에서 여러 조건을 검증할 때 사용한다. 반면 `t.Fatal`은 즉시 테스트를 중단하므로, 이후 테스트가 의미 없는 경우(예: 필수 초기화 실패)에 사용한다.

### 테이블 주도 테스트 (Table-Driven Tests)

Go에서 가장 권장되는 테스트 패턴이다. 여러 테스트 케이스를 구조체 슬라이스로 정의하고 반복 실행한다. 이 패턴은 새로운 테스트 케이스를 추가할 때 구조체 하나만 슬라이스에 추가하면 되므로 확장성이 뛰어나다:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"양수 더하기", 2, 3, 5},
        {"음수 더하기", -1, -2, -3},
        {"영 더하기", 0, 5, 5},
        {"큰 수 더하기", 1000000, 2000000, 3000000},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d",
                    tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

`t.Run()`을 사용하면 각 케이스가 **서브테스트**로 실행되어, 어떤 케이스가 실패했는지 명확히 알 수 있다. 또한 `-run "TestAdd/양수_더하기"`처럼 특정 서브테스트만 실행할 수도 있다. 서브테스트는 각각 독립적인 `*testing.T`를 가지므로 `t.Parallel()`을 호출하여 병렬로 실행하는 것도 가능하다.

### 테스트 헬퍼 함수

반복되는 테스트 로직은 헬퍼 함수로 추출하면 코드 중복을 줄일 수 있다:

```go
func assertEqual(t *testing.T, got, want int) {
    t.Helper() // 에러 발생 시 호출자 위치를 보고한다
    if got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}
```

`t.Helper()`를 호출하면 에러가 발생했을 때 헬퍼 함수 내부가 아닌 실제 테스트 코드의 위치를 보고하여 디버깅이 훨씬 수월해진다.

---

## 26.2 테스트 주도 개발 (TDD)

### TDD 사이클: Red-Green-Refactor

TDD는 다음 세 단계를 반복하는 개발 방법론이다. 테스트를 먼저 작성함으로써 요구사항을 명확히 정의하고, 최소한의 코드로 기능을 구현한 후, 리팩터링으로 코드 품질을 개선한다:

```
┌─────────────────────────────────────────┐
│                                         │
│   ① Red (실패하는 테스트 작성)            │
│        │                                │
│        ▼                                │
│   ② Green (테스트를 통과하는 최소 코드)    │
│        │                                │
│        ▼                                │
│   ③ Refactor (코드 개선)                 │
│        │                                │
│        └──────────► ①로 돌아감           │
│                                         │
└─────────────────────────────────────────┘
```

**Red 단계**에서는 아직 구현이 없으므로 테스트가 반드시 실패해야 한다. 만약 실패하지 않는다면 테스트가 잘못 작성된 것이다. **Green 단계**에서는 테스트를 통과시키기 위한 최소한의 코드만 작성한다. 완벽한 구현이 아닌, "동작하는" 코드가 목표이다. **Refactor 단계**에서는 기능은 유지하면서 코드 구조를 개선한다. 이때 테스트가 통과 상태를 유지하는지 지속적으로 확인한다.

### TDD 예제: Divide 함수 만들기

**1단계 (Red)**: 실패하는 테스트 먼저 작성

```go
func TestDivide(t *testing.T) {
    result, err := Divide(10, 2)
    if err != nil {
        t.Fatal(err)
    }
    if result != 5.0 {
        t.Errorf("Divide(10, 2) = %f; want 5.0", result)
    }
}

func TestDivideByZero(t *testing.T) {
    _, err := Divide(10, 0)
    if err == nil {
        t.Error("0으로 나눌 때 에러가 발생해야 한다")
    }
}
```

**2단계 (Green)**: 테스트를 통과하는 최소한의 코드 작성

```go
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("0으로 나눌 수 없다")
    }
    return a / b, nil
}
```

**3단계 (Refactor)**: 코드를 더 깔끔하게 개선한다. 에러 메시지를 상수로 분리하거나, 에러 타입을 정의하는 등의 개선이 가능하다.

### TDD의 장점

- 요구사항을 코드로 명확히 정의하므로 구현 방향이 분명해진다.
- 테스트가 항상 최신 상태로 유지된다.
- 과도한 구현을 방지하고 필요한 기능만 구현하게 된다.
- 리팩터링 시 안전망 역할을 한다.

---

## 26.3 벤치마크

### Benchmark 함수

벤치마크 함수는 코드의 성능을 측정하는 데 사용한다. 다음 규칙을 따른다:

1. `Benchmark`로 시작하는 함수 이름
2. `*testing.B` 매개변수를 받음
3. `b.N`번 반복 실행 (Go가 자동으로 적절한 횟수를 결정)

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(100, 200)
    }
}
```

Go의 벤치마크 프레임워크는 `b.N` 값을 자동으로 조정한다. 처음에는 작은 값으로 시작하여 결과가 안정적일 때까지 점점 증가시킨다. 이를 통해 통계적으로 신뢰할 수 있는 결과를 얻는다.

### go test -bench 명령어

```bash
# 모든 벤치마크 실행
go test -bench .

# 특정 벤치마크만 실행
go test -bench BenchmarkAdd

# 메모리 할당 정보도 함께 출력
go test -bench . -benchmem

# 10초 동안 벤치마크 실행
go test -bench . -benchtime 10s

# 벤치마크를 3회 반복
go test -bench . -count 3

# 일반 테스트를 건너뛰고 벤치마크만 실행
go test -bench . -run ^$
```

`-run ^$`는 정규식 `^$` (빈 문자열)에 매칭되는 테스트가 없으므로, 일반 테스트를 실행하지 않고 벤치마크만 수행한다.

### 벤치마크 결과 읽기

```
BenchmarkAdd-8    1000000000    0.2500 ns/op    0 B/op    0 allocs/op
```

| 항목 | 설명 |
|------|------|
| `BenchmarkAdd-8` | 함수 이름과 GOMAXPROCS 값이다 |
| `1000000000` | 실행 횟수 (b.N)이다 |
| `0.2500 ns/op` | 연산당 소요 시간이다 |
| `0 B/op` | 연산당 메모리 할당량이다 |
| `0 allocs/op` | 연산당 메모리 할당 횟수이다 |

`ns/op` 값이 클수록 느린 함수이다. `B/op`과 `allocs/op`는 메모리 효율을 나타낸다. 할당 횟수가 많을수록 가비지 컬렉터에 부담을 주므로 성능에 부정적인 영향을 미친다.

### 벤치마크 초기화 코드 제외

벤치마크에서 초기화 코드의 시간을 제외하려면 `b.ResetTimer()`를 사용한다:

```go
func BenchmarkComplexOperation(b *testing.B) {
    // 초기화 코드 (벤치마크에서 제외)
    data := prepareData()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        processData(data)
    }
}
```

반복문 내부에서 매번 초기화가 필요한 경우에는 `b.StopTimer()`와 `b.StartTimer()`를 사용하여 측정 구간을 더 세밀하게 제어할 수 있다:

```go
func BenchmarkWithSetup(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        data := prepareData() // 측정에서 제외
        b.StartTimer()
        processData(data)     // 이 부분만 측정
    }
}
```

### 서브 벤치마크

벤치마크에서도 `b.Run()`으로 서브 벤치마크를 만들 수 있다. 입력 크기별 성능 비교에 유용하다:

```go
func BenchmarkSort(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            data := generateData(size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                sort.Ints(data)
            }
        })
    }
}
```

---

## 핵심 요약

1. Go는 `testing` 패키지와 `go test` 명령어로 강력한 테스트를 지원한다.
2. 테스트 파일은 `_test.go`로 끝나고, 테스트 함수는 `Test`로 시작한다.
3. **테이블 주도 테스트**는 Go에서 가장 권장되는 테스트 패턴이다.
4. **TDD**는 Red-Green-Refactor 사이클을 반복하는 개발 방법론이다.
5. 벤치마크 함수는 `Benchmark`로 시작하고, `b.N`번 반복한다.
6. `go test -bench . -benchmem`으로 성능과 메모리 사용량을 측정한다.
7. `t.Helper()`를 사용하면 헬퍼 함수의 에러 위치 보고를 개선할 수 있다.
8. `-cover` 플래그로 테스트 커버리지를 확인하여 테스트 품질을 관리한다.

---

## 연습문제

### 연습문제 1: 문자열 유틸리티 테스트
다음 함수들을 테스트하는 테이블 주도 테스트를 작성하라:
- `Reverse(s string) string` - 문자열을 뒤집는 함수
- `IsPalindrome(s string) bool` - 회문 판별 함수
- 한글, 영문, 빈 문자열, 특수문자 등 다양한 케이스를 포함하라.

### 연습문제 2: TDD로 스택 구현
TDD 방식으로 정수형 스택(Stack)을 구현하라:
- `Push(val int)` - 값 추가
- `Pop() (int, error)` - 값 꺼내기
- `Peek() (int, error)` - 맨 위 값 확인
- `IsEmpty() bool` - 비어있는지 확인
- `Size() int` - 현재 스택 크기 반환

### 연습문제 3: 벤치마크 비교
문자열을 연결하는 세 가지 방법의 성능을 벤치마크로 비교하라:
- `+` 연산자
- `fmt.Sprintf`
- `strings.Builder`
- 문자열 10개, 100개, 1000개를 연결하는 서브 벤치마크를 작성하라.

### 연습문제 4: 테스트 커버리지 분석
간단한 계산기 패키지(Add, Sub, Mul, Div)를 만들고 테스트를 작성한 후, `go test -cover`로 커버리지를 확인하라. 100% 커버리지를 달성하기 위해 어떤 테스트 케이스를 추가해야 하는지 분석하라.

### 연습문제 5: t.Error vs t.Fatal
`t.Error`와 `t.Fatal`의 동작 차이를 실험하라. 하나의 테스트 함수에서 `t.Error`를 두 번 호출하는 경우와 `t.Fatal`을 호출하는 경우의 출력 차이를 확인하고, 각각 어떤 상황에서 사용해야 하는지 설명하라.

### 연습문제 6: 서브테스트와 병렬 실행
테이블 주도 테스트에서 `t.Parallel()`을 사용하여 서브테스트를 병렬로 실행하라. 병렬 실행 시 `tt` 변수 캡처에 주의해야 하는 이유를 설명하고, 올바른 코드를 작성하라.

### 연습문제 7: 테스트 헬퍼 함수 작성
JSON 응답을 검증하는 테스트 헬퍼 함수를 작성하라. `t.Helper()`를 호출하는 경우와 호출하지 않는 경우의 에러 메시지 출력 차이를 비교하라.

### 연습문제 8: 벤치마크 결과 해석
다음 벤치마크 결과를 해석하고, 어떤 구현이 더 좋은지 판단하라:
```
BenchmarkImplA-8    500000    3200 ns/op    256 B/op    4 allocs/op
BenchmarkImplB-8    300000    4100 ns/op     64 B/op    1 allocs/op
```

### 연습문제 9: Cleanup 함수 활용
임시 파일을 생성하여 테스트에 사용하고, `t.Cleanup()`을 이용하여 테스트 종료 후 자동으로 임시 파일을 삭제하는 테스트를 작성하라.

### 연습문제 10: 퍼즈 테스트
Go 1.18에서 도입된 퍼즈 테스트(Fuzz Testing)를 사용하여 `Reverse` 함수에 대한 퍼즈 테스트를 작성하라. `func FuzzReverse(f *testing.F)` 형태로 작성하고, `go test -fuzz .`로 실행하라.

---

## 구현 과제

### 과제 1: 수학 유틸리티 패키지
`mathutil` 패키지를 TDD 방식으로 구현하라. 다음 함수를 포함해야 한다:
- `GCD(a, b int) int` - 최대공약수
- `LCM(a, b int) int` - 최소공배수
- `IsPrime(n int) bool` - 소수 판별
- `Fibonacci(n int) int` - n번째 피보나치 수
- 모든 함수에 대해 테이블 주도 테스트를 작성하고, 100% 커버리지를 달성하라.

### 과제 2: 문자열 검증기
이메일, 전화번호, URL 등의 형식을 검증하는 `validator` 패키지를 구현하라:
- `IsValidEmail(email string) bool`
- `IsValidPhone(phone string) bool`
- `IsValidURL(url string) bool`
- 유효한 입력과 무효한 입력 모두에 대해 테이블 주도 테스트를 작성하라.

### 과제 3: 정렬 알고리즘 벤치마크
버블 정렬, 삽입 정렬, 퀵 정렬을 직접 구현하고, 입력 크기(100, 1000, 10000)별로 벤치마크를 수행하라. `sort.Ints`(표준 라이브러리)와도 비교하라. 결과를 표로 정리하라.

### 과제 4: 캐시 구현과 벤치마크
간단한 LRU 캐시를 구현하고, 다음 벤치마크를 수행하라:
- `Set` 성능
- `Get` (캐시 히트) 성능
- `Get` (캐시 미스) 성능
- 캐시 크기(100, 1000, 10000)에 따른 성능 변화

### 과제 5: 테스트 가능한 코드 리팩터링
다음과 같이 테스트하기 어려운 코드를 테스트 가능한 구조로 리팩터링하라:
- 현재 시간에 의존하는 함수 (인터페이스 주입 활용)
- 파일 I/O에 의존하는 함수 (io.Reader/Writer 활용)
- 외부 API에 의존하는 함수 (인터페이스 기반 Mock 활용)

---

## 프로젝트 과제

### 프로젝트 1: 자동 테스트 리포트 생성기
`go test -json` 출력을 파싱하여 테스트 결과를 요약하는 CLI 도구를 만들어라:
- 전체 테스트 수, 성공 수, 실패 수, 건너뛴 수를 집계한다.
- 실패한 테스트 목록과 에러 메시지를 보기 좋게 출력한다.
- 테스트 실행 시간이 가장 오래 걸린 상위 5개 테스트를 표시한다.
- 이 도구 자체에 대한 테스트도 작성한다.

### 프로젝트 2: 벤치마크 비교 도구
두 번의 벤치마크 실행 결과를 비교하여 성능 변화를 분석하는 프로그램을 만들어라:
- `go test -bench . -benchmem` 출력을 파싱한다.
- 각 벤치마크의 ns/op, B/op, allocs/op 변화를 비교한다.
- 성능이 개선된 항목과 악화된 항목을 구분하여 표시한다.
- `benchstat`과 유사한 기능을 직접 구현해 보라.
