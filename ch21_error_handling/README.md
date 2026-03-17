# 21장 에러 핸들링

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch21_error_handling/basic.go
go run ch21_error_handling/custom_error.go
go run ch21_error_handling/panic_recover.go
go run ch21_error_handling/wrapping.go
```

> **Makefile 활용**: `make run CH=ch21_error_handling` 또는 `make run CH=ch21_error_handling FILE=basic.go`

---

Go 언어에서는 예외(exception)가 아닌 **에러 값(error value)**을 반환하는 방식으로 에러를 처리한다. 다른 언어(Java, Python, C++ 등)에서는 `try-catch` 기반의 예외 처리를 사용하지만, Go는 에러를 일반적인 값으로 취급하여 명시적으로 처리하도록 설계되어 있다. 이 장에서는 Go의 에러 처리 패턴을 체계적으로 학습한다.

---

## 21.1 에러 반환

### error 인터페이스

Go의 에러는 내장 인터페이스 `error`로 표현된다. 이 인터페이스는 단 하나의 메서드만 요구하므로, 어떤 타입이든 `Error() string` 메서드를 구현하면 에러로 사용할 수 있다.

```go
type error interface {
    Error() string
}
```

함수가 에러를 반환할 수 있으면, 관례적으로 **마지막 반환값**을 `error` 타입으로 선언한다. 이는 Go 커뮤니티 전체에서 따르는 강력한 관례이며, 표준 라이브러리도 이 패턴을 일관되게 사용한다.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("0으로 나눌 수 없습니다")
    }
    return a / b, nil
}
```

에러가 없는 경우에는 `nil`을 반환한다. 호출자는 반환된 에러가 `nil`인지 확인하여 성공 여부를 판단한다.

### errors.New

`errors` 패키지의 `New` 함수는 가장 간단한 에러 생성 방법이다. 정적인 에러 메시지를 가진 에러를 생성할 때 사용한다.

```go
import "errors"

var ErrNotFound = errors.New("찾을 수 없습니다")
```

패키지 수준에서 미리 정의해두는 에러를 **센티넬 에러(sentinel error)**라고 부른다. 표준 라이브러리에서도 `io.EOF`, `os.ErrNotExist` 등의 센티넬 에러를 제공한다. 센티넬 에러는 `var`로 선언하여 `errors.Is`로 비교할 수 있도록 하는 것이 관례이다.

### fmt.Errorf

동적인 정보를 포함한 에러 메시지를 만들 때 사용한다. `fmt.Sprintf`와 동일한 포맷 동사를 지원하며, 런타임 정보(파일명, 사용자 입력값 등)를 에러 메시지에 포함시킬 수 있다.

```go
import "fmt"

func openFile(name string) error {
    return fmt.Errorf("파일 '%s'을(를) 열 수 없습니다", name)
}
```

### 에러 처리 패턴

Go에서 가장 흔한 패턴은 에러를 즉시 확인하는 것이다. 이를 **guard clause** 패턴이라고도 부른다. 에러가 발생하면 즉시 처리하고 반환함으로써, 정상 경로(happy path)의 들여쓰기를 최소화한다.

```go
result, err := divide(10, 0)
if err != nil {
    fmt.Println("에러 발생:", err)
    return
}
fmt.Println("결과:", result)
```

> **주의:** `err`를 무시하지 않아야 한다! `_`로 에러를 버리는 것은 잠재적인 버그의 원인이다. 에러를 의도적으로 무시해야 하는 경우에는 그 이유를 주석으로 명시하는 것이 좋다.

---

## 21.2 에러 타입

### 커스텀 에러 타입

`error` 인터페이스를 구현하면 어떤 타입이든 에러로 사용할 수 있다. 커스텀 에러 타입을 사용하면 에러 메시지뿐만 아니라 에러의 종류, 발생 위치, 관련 데이터 등 **추가 정보**를 구조화하여 담을 수 있다.

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("검증 실패 [%s]: %s", e.Field, e.Message)
}
```

커스텀 에러 타입은 포인터 리시버로 `Error()` 메서드를 구현하는 것이 일반적이다. 이는 에러 비교 시 값 복사가 아닌 포인터 비교를 통해 의도치 않은 동등성 문제를 방지하기 위함이다.

### errors.Is

특정 에러 값과 일치하는지 확인한다. 에러 체인(wrapping)도 재귀적으로 탐색하므로, 래핑된 에러 내부에 원본 에러가 있는지 검사할 수 있다. 이전에는 `err == ErrPermission`처럼 직접 비교했지만, 에러 래핑이 도입된 이후로는 `errors.Is`를 사용하는 것이 올바른 방법이다.

```go
var ErrPermission = errors.New("권한이 없습니다")

if errors.Is(err, ErrPermission) {
    fmt.Println("권한 에러입니다")
}
```

### errors.As

에러를 특정 타입으로 변환(type assertion)한다. `errors.Is`가 특정 에러 **값**을 찾는 것이라면, `errors.As`는 특정 에러 **타입**을 찾는 것이다. 에러 체인을 순회하면서 해당 타입의 에러가 있으면 변환하여 반환한다.

```go
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println("필드:", valErr.Field)
    fmt.Println("메시지:", valErr.Message)
}
```

`errors.As`의 두 번째 인자는 반드시 포인터의 포인터여야 한다. 이는 함수 내부에서 역참조하여 값을 할당하기 위함이다.

### 에러 래핑 (Error Wrapping) - %w

`fmt.Errorf`에서 `%w` 동사를 사용하면 원본 에러를 감싸면서 맥락 정보를 추가할 수 있다. 래핑은 에러가 발생한 하위 레이어의 정보를 보존하면서, 상위 레이어에서 맥락을 덧붙이는 패턴이다.

```go
func readConfig(path string) error {
    _, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("설정 파일 읽기 실패: %w", err)
    }
    return nil
}
```

래핑된 에러는 `errors.Is`와 `errors.As`로 원본 에러를 추출할 수 있다.

```go
err := readConfig("config.json")
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("파일이 존재하지 않습니다")
}
```

> **`%v` vs `%w`:** `%v`는 단순히 문자열로 변환하여 에러 체인이 끊어지고, `%w`는 에러 체인을 유지하여 `errors.Is`/`errors.As`로 내부 에러를 탐색할 수 있게 한다. 의도적으로 내부 에러를 숨기고 싶을 때만 `%v`를 사용한다.

---

## 21.3 패닉 (panic & recover)

### panic

`panic`은 프로그램을 즉시 중단시키는 내장 함수이다. 현재 함수의 실행을 멈추고, defer된 함수들을 실행한 뒤, 호출 스택을 따라 올라가며 프로그램을 종료한다. panic이 발생하면 스택 트레이스가 출력되어 어디서 panic이 발생했는지 추적할 수 있다.

```go
func mustPositive(n int) {
    if n < 0 {
        panic("음수는 허용되지 않습니다")
    }
}
```

`panic`에는 `string`, `error`, 또는 어떤 타입의 값이든 인자로 전달할 수 있다. 그러나 관례적으로 문자열이나 `error` 타입을 전달한다.

### recover

`recover`는 `defer` 함수 내에서만 동작하며, 패닉을 잡아서 프로그램 종료를 방지한다. `recover`는 `panic`에 전달된 값을 반환하며, 패닉이 발생하지 않은 경우에는 `nil`을 반환한다.

```go
func safeFunction() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("패닉 복구:", r)
        }
    }()

    panic("심각한 오류!")
}
```

`recover`는 반드시 `defer`로 등록된 함수 내에서 직접 호출해야 한다. 중첩된 함수 호출에서는 동작하지 않는다.

### 언제 panic을 사용하나?

| 상황 | 사용 여부 |
|------|-----------|
| 프로그램 초기화 실패 (DB 연결 등) | O - panic 사용 |
| 논리적으로 도달 불가능한 코드 | O - panic 사용 |
| 일반적인 에러 (파일 없음, 입력 오류 등) | X - error 반환 |
| 라이브러리/패키지 코드 | X - error 반환 |

표준 라이브러리에서도 이 원칙을 따른다. 예를 들어 `regexp.MustCompile`은 정규식이 유효하지 않으면 panic을 발생시키는데, 이는 보통 컴파일 타임에 알 수 있는 정적 정규식에만 사용한다.

> **원칙:** 일반적인 에러는 `error`를 반환하고, 프로그램이 계속 실행될 수 없는 치명적 상황에서만 `panic`을 사용한다.

---

## 핵심 요약

| 개념 | 설명 |
|------|------|
| `error` 인터페이스 | Go의 에러 처리 기본 단위. `Error() string` 메서드를 가진다 |
| `errors.New` | 간단한 에러 메시지 생성 |
| `fmt.Errorf` | 포맷된 에러 메시지 생성. `%w`로 래핑 가능 |
| 센티넬 에러 | 패키지 수준에서 미리 정의하는 에러 값 |
| 커스텀 에러 타입 | `error` 인터페이스를 구현한 사용자 정의 타입 |
| `errors.Is` | 에러 체인에서 특정 에러 값 검색 |
| `errors.As` | 에러 체인에서 특정 에러 타입으로 변환 |
| `%w` | 에러 래핑 - 원본 에러를 유지하면서 맥락 추가 |
| `panic` | 프로그램 즉시 중단 (치명적 오류용) |
| `recover` | defer 내에서 panic 복구 |

---

## 연습문제

### 문제 1: 나이 검증 함수
사용자의 나이를 입력받아 검증하는 함수를 작성하라.
- 0 미만이면 "나이는 음수일 수 없습니다" 에러 반환
- 200 초과이면 "비현실적인 나이입니다" 에러 반환
- 정상이면 "성인" 또는 "미성년자" 문자열 반환 (기준: 18세)

### 문제 2: 커스텀 에러 타입
`HttpError` 구조체를 만들고, 상태 코드(StatusCode)와 메시지(Message)를 포함하도록 하라. `errors.As`를 사용하여 상태 코드에 따라 다른 처리를 하는 코드를 작성하라.

### 문제 3: 에러 래핑 체인
3개의 함수가 순차적으로 호출되면서 에러를 래핑하는 코드를 작성하라.
- `readDatabase() error` -> `processData() error` -> `handleRequest() error`
- 각 함수에서 `%w`로 에러를 래핑하고, 최종적으로 `errors.Is`로 원본 에러를 찾을 수 있어야 한다.

### 문제 4: panic과 recover
슬라이스의 인덱스 접근에서 발생할 수 있는 panic을 recover로 처리하는 `safeIndex` 함수를 작성하라.
- `func safeIndex(s []int, idx int) (int, error)`

### 문제 5: 센티넬 에러 정의
파일 처리 패키지를 위한 센티넬 에러를 설계하라.
- `ErrFileNotFound`, `ErrPermissionDenied`, `ErrFileTooLarge` 세 가지 에러를 정의하라.
- 파일 크기를 검사하는 `validateFile(path string, maxSize int64) error` 함수를 작성하되, 상황에 맞는 센티넬 에러를 반환하도록 하라.
- 호출자 측에서 `errors.Is`로 각 에러를 구분하여 처리하는 코드를 작성하라.

### 문제 6: 다중 에러 수집
여러 필드를 동시에 검증하고 발생한 모든 에러를 수집하는 함수를 작성하라.
- `type MultiError struct { Errors []error }`를 정의하고 `error` 인터페이스를 구현하라.
- `validateUser(name, email, age string) error` 함수를 만들어, 빈 이름, 잘못된 이메일 형식, 유효하지 않은 나이를 모두 검사하라.
- 에러가 하나도 없으면 `nil`을 반환하고, 하나 이상이면 `MultiError`를 반환하라.

### 문제 7: 에러 래핑과 Unwrap
`errors.Unwrap`을 사용하여 에러 체인을 순회하며 모든 에러 메시지를 출력하는 `printErrorChain(err error)` 함수를 작성하라. 3단계 이상 래핑된 에러를 만들어 테스트하라.

### 문제 8: fmt.Errorf의 %w와 %v 차이 확인
동일한 에러를 `%w`로 래핑한 경우와 `%v`로 래핑한 경우를 비교하는 프로그램을 작성하라.
- 두 경우 모두에서 `errors.Is`와 `errors.As`를 사용해보고, 결과 차이를 확인하라.
- 왜 `%v`는 에러 체인이 끊어지는지 주석으로 설명하라.

### 문제 9: recover를 활용한 안전한 함수 실행기
임의의 함수를 안전하게 실행하는 `safeRun(fn func()) error`를 작성하라.
- `fn`에서 panic이 발생하면 해당 panic 값을 에러로 변환하여 반환하라.
- panic이 발생하지 않으면 `nil`을 반환하라.
- 다양한 panic 상황(문자열, 에러, 정수 등)을 테스트하라.

### 문제 10: 에러 처리 리팩토링
아래 코드에서 에러 처리가 잘못된 부분을 모두 찾아 수정하라.
```go
func processFile(path string) {
    data, _ := os.ReadFile(path)
    var config Config
    json.Unmarshal(data, &config)
    result := compute(config)
    fmt.Println(result)
}
```

---

## 구현 과제

### 과제 1: 은행 계좌 에러 처리 시스템
은행 계좌를 관리하는 프로그램을 작성하라.
- `InsufficientFundsError`, `AccountNotFoundError`, `InvalidAmountError` 커스텀 에러 타입을 정의하라.
- `Deposit(accountID string, amount float64) error`와 `Withdraw(accountID string, amount float64) error` 함수를 구현하라.
- 호출자 측에서 `errors.As`를 사용하여 에러 종류에 따라 다른 에러 메시지를 출력하라.

### 과제 2: 설정 파일 로더
JSON 설정 파일을 읽어서 구조체로 파싱하는 프로그램을 작성하라.
- 파일 열기, 읽기, JSON 파싱 각 단계에서 발생하는 에러를 `%w`로 래핑하여 맥락을 추가하라.
- 필수 필드가 누락된 경우 `MissingFieldError` 커스텀 에러를 반환하라.
- 최종적으로 `errors.Is`와 `errors.As`를 사용하여 에러를 분류하고 사용자 친화적인 메시지를 출력하라.

### 과제 3: HTTP 요청 재시도 함수
에러 종류에 따라 재시도 여부를 결정하는 HTTP 요청 함수를 작성하라.
- `TemporaryError` 인터페이스(`Temporary() bool` 메서드를 가진)를 정의하라.
- 일시적 에러(네트워크 타임아웃 등)인 경우 최대 3회까지 재시도하라.
- 영구적 에러(404 Not Found 등)인 경우 즉시 반환하라.
- 각 재시도 시 에러를 래핑하여 시도 횟수 정보를 포함시켜라.

### 과제 4: 로그 파서와 에러 보고
로그 파일을 파싱하여 에러 통계를 보고하는 프로그램을 작성하라.
- 각 줄을 파싱할 때 발생하는 에러를 줄 번호 정보와 함께 수집하라 (`ParseError` 커스텀 타입 사용).
- 파싱 불가능한 줄은 건너뛰되, 모든 에러를 기록하라.
- 최종적으로 총 처리 줄 수, 성공 줄 수, 에러 줄 수, 에러 세부사항을 출력하라.

### 과제 5: panic-safe 미들웨어
HTTP 핸들러에서 발생하는 panic을 복구하는 미들웨어를 작성하라.
- `func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc` 형태로 구현하라.
- panic이 발생하면 500 Internal Server Error를 응답하고, panic 정보를 로그에 기록하라.
- 스택 트레이스를 `runtime/debug` 패키지로 캡처하여 로그에 포함시켜라.

---

## 프로젝트 과제

### 프로젝트 1: 데이터 검증 프레임워크
구조체 필드를 검증하는 범용 검증 프레임워크를 구현하라.
- `Validator` 인터페이스(`Validate() error`)를 정의하라.
- 다양한 검증 규칙(필수 값, 최소/최대 길이, 이메일 형식, 숫자 범위 등)을 커스텀 에러 타입으로 표현하라.
- 여러 검증 규칙을 체이닝하여 적용할 수 있도록 하라.
- 모든 검증 실패를 수집하여 한 번에 보고하는 기능을 제공하라.
- `User`, `Product` 등 2개 이상의 구조체에 대해 검증을 시연하라.

### 프로젝트 2: 에러 추적 시스템
에러 발생 시 호출 스택 정보를 자동으로 캡처하고, 구조화된 에러 보고서를 생성하는 시스템을 구현하라.
- `TrackedError` 타입을 만들어 에러 메시지, 발생 시각, 호출 스택, 래핑된 원본 에러를 포함시켜라.
- `runtime` 패키지를 활용하여 호출 스택을 캡처하라.
- 에러 보고서를 JSON 형태로 출력하는 기능을 구현하라.
- 3단계 이상의 함수 호출 체인에서 에러가 전파되는 시나리오를 시연하라.
