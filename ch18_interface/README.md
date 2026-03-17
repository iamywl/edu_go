# Chapter 18: 인터페이스

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch18_interface/basic.go
go run ch18_interface/duck_typing.go
go run ch18_interface/type_switch.go
go run ch18_interface/empty_interface.go
```

> **Makefile 활용**: `make run CH=ch18_interface` 또는 `make run CH=ch18_interface FILE=basic.go`

---

Go 언어의 **인터페이스(Interface)**는 메서드의 집합을 정의하는 타입이다. 인터페이스를 통해 다형성(polymorphism)을 구현하고, 코드 간 결합도를 낮출 수 있다. Go의 인터페이스는 암묵적으로 구현되기 때문에, 다른 언어와 비교했을 때 매우 유연하고 강력한 추상화 도구이다.

---

## 18.1 인터페이스 (메서드 집합 정의)

### 인터페이스란?

인터페이스는 **메서드의 목록**을 정의한다. 어떤 타입이 인터페이스에 정의된 모든 메서드를 구현하면, 그 타입은 해당 인터페이스를 구현한 것이다. 인터페이스는 "무엇을 할 수 있는가"를 정의하지, "어떻게 하는가"를 정의하지 않는다. 이것이 인터페이스의 핵심 철학이다.

### 기본 문법

```go
type 인터페이스이름 interface {
    메서드이름1(매개변수) 반환타입
    메서드이름2(매개변수) 반환타입
}
```

인터페이스 이름은 관례적으로 메서드 이름에 `-er` 접미사를 붙인다. 예를 들어 `Read()` 메서드를 가진 인터페이스는 `Reader`, `Write()` 메서드를 가진 인터페이스는 `Writer`라고 이름 짓는다.

### 예시

```go
// Shape 인터페이스 정의
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle이 Shape를 구현 (명시적 선언 없이!)
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// Circle도 Shape를 구현
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Shape 인터페이스 타입으로 사용
func PrintInfo(s Shape) {
    fmt.Printf("넓이: %.2f, 둘레: %.2f\n", s.Area(), s.Perimeter())
}
```

### 인터페이스의 특징

- **암묵적 구현**: `implements` 키워드가 없다. Java나 C#과 달리 타입이 인터페이스를 구현한다고 선언할 필요가 없다.
- 인터페이스에 정의된 **모든 메서드를 구현**하면 자동으로 그 인터페이스를 만족한다.
- 인터페이스는 **타입**이므로 변수, 매개변수, 반환값으로 사용할 수 있다.
- 인터페이스 변수는 내부적으로 두 개의 포인터를 가진다: 실제 값을 가리키는 포인터와 타입 정보를 가리키는 포인터이다.

---

## 18.2 인터페이스 왜 쓰나? (다형성, 디커플링)

### 다형성 (Polymorphism)

같은 인터페이스 타입으로 **다양한 구현체**를 사용할 수 있다. 이것이 다형성이다. 호출하는 쪽에서는 구체적인 타입을 알 필요 없이, 인터페이스에 정의된 메서드만 사용하면 된다.

```go
shapes := []Shape{
    Rectangle{10, 5},
    Circle{7},
    Rectangle{3, 4},
}

for _, s := range shapes {
    PrintInfo(s)  // 각 타입에 맞는 메서드가 호출됨
}
```

위 코드에서 `PrintInfo()`는 `Shape` 인터페이스만 알고 있으며, 실제로 전달되는 것이 `Rectangle`인지 `Circle`인지 알 필요가 없다. 런타임에 실제 타입의 메서드가 호출된다.

### 디커플링 (Decoupling)

인터페이스를 사용하면 **구현에 의존하지 않고 동작에 의존**할 수 있다. 이를 "의존성 역전 원칙(Dependency Inversion Principle)"이라 하며, 좋은 소프트웨어 설계의 핵심이다.

```go
// 데이터베이스 인터페이스
type Database interface {
    Get(key string) (string, error)
    Set(key string, value string) error
}

// 실제 DB, 메모리 DB, 테스트용 DB 등 다양한 구현 가능
// 사용하는 코드는 Database 인터페이스만 알면 됨
func ProcessData(db Database) {
    // db가 어떤 구현체인지 몰라도 됨
    val, _ := db.Get("key")
    fmt.Println(val)
}
```

이 패턴의 장점은 `ProcessData()` 함수를 수정하지 않고도, 새로운 데이터베이스 구현체를 추가하여 사용할 수 있다는 것이다.

### 테스트 용이성

인터페이스를 활용하면 테스트 시 실제 구현체 대신 **가짜(mock) 구현체**를 사용할 수 있다. 이를 통해 외부 의존성(데이터베이스, 네트워크 등) 없이 단위 테스트를 수행할 수 있다.

```go
// 테스트할 때 가짜(mock) 구현체를 사용할 수 있음
type MockDB struct {
    data map[string]string
}
func (m MockDB) Get(key string) (string, error) { ... }
func (m MockDB) Set(key, val string) error       { ... }
```

### Go에서의 인터페이스 설계 원칙

Go 커뮤니티에서는 **작은 인터페이스**를 선호한다. 메서드가 하나인 인터페이스가 가장 유용하고 재사용성이 높다. 표준 라이브러리의 `io.Reader`, `io.Writer`, `fmt.Stringer` 등이 좋은 예이다. 큰 인터페이스가 필요하면 작은 인터페이스를 조합하여 만드는 것이 바람직하다.

---

## 18.3 덕 타이핑 (구조적 타이핑, 암묵적 구현)

### 덕 타이핑이란?

> "오리처럼 걷고, 오리처럼 꽥꽥거리면, 그것은 오리다."

Go에서는 타입이 인터페이스의 **메서드를 모두 구현하기만 하면** 자동으로 해당 인터페이스를 만족한다. 명시적인 선언이 필요 없다. 이를 "구조적 타이핑(structural typing)"이라 하며, 덕 타이핑의 정적 타입 버전이라고 볼 수 있다.

```go
type Stringer interface {
    String() string
}

type Dog struct {
    Name string
}

// Dog은 String() 메서드를 가지고 있으므로 Stringer 인터페이스를 자동으로 구현
func (d Dog) String() string {
    return fmt.Sprintf("강아지: %s", d.Name)
}
```

### 구조적 타이핑 vs 명목적 타이핑

| 구분 | 구조적 타이핑 (Go) | 명목적 타이핑 (Java, C#) |
|------|-------------------|------------------------|
| 구현 방식 | 메서드만 맞으면 자동 구현 | `implements` 키워드 필요 |
| 유연성 | 높음 | 낮음 |
| 기존 타입 확장 | 쉬움 (새 인터페이스 정의만) | 소스 코드 수정 필요 |
| 컴파일 타임 검증 | 사용 시점에 검증 | 선언 시점에 검증 |

### 덕 타이핑의 장점

1. **기존 코드 수정 없이** 새로운 인터페이스에 맞출 수 있다.
2. 패키지 간 **느슨한 결합**을 유지할 수 있다.
3. 서드파티 라이브러리의 타입도 인터페이스를 통해 추상화할 수 있다.
4. **인터페이스를 소비자(사용자) 측에서 정의**할 수 있다. 구현체를 제공하는 패키지가 아니라, 그것을 사용하는 패키지에서 필요한 인터페이스를 정의하는 것이 Go의 관용적 패턴이다.

### 인터페이스 구현 확인

컴파일 타임에 특정 타입이 인터페이스를 구현하는지 확인하고 싶다면, 다음과 같은 패턴을 사용한다:

```go
// 컴파일 타임에 Dog이 Stringer를 구현하는지 확인
var _ Stringer = Dog{}
var _ Stringer = (*Dog)(nil)
```

이 변수는 실제로 사용되지 않지만, Dog이 Stringer를 구현하지 않으면 컴파일 에러가 발생하여 즉시 문제를 발견할 수 있다.

---

## 18.4 인터페이스 기능 더 알기

### 빈 인터페이스 (interface{} / any)

메서드가 하나도 없는 인터페이스이다. **모든 타입**이 빈 인터페이스를 만족한다. 모든 타입은 "메서드가 0개인 집합"을 당연히 구현하기 때문이다.

```go
// Go 1.18 이전
var val interface{}

// Go 1.18 이후 (any는 interface{}의 별칭)
var val any

val = 42         // 정수
val = "hello"    // 문자열
val = true       // 불리언
val = []int{1,2} // 슬라이스
```

### 빈 인터페이스의 활용

```go
// 어떤 타입이든 받을 수 있는 함수
func PrintAnything(val any) {
    fmt.Println(val)
}

// 여러 타입을 담을 수 있는 슬라이스
mixed := []any{42, "hello", true, 3.14}
```

빈 인터페이스는 편리하지만 남용하면 안 된다. 타입 안전성을 잃기 때문이다. 가능하면 구체적인 인터페이스나 제네릭(Go 1.18+)을 사용하는 것이 좋다.

### 타입 단언 (Type Assertion)

인터페이스 값에서 **구체적인 타입을 꺼내는** 방법이다. 인터페이스 변수에 저장된 실제 값의 타입을 확인하고, 해당 타입으로 변환할 때 사용한다.

```go
var val any = "hello"

// 기본 타입 단언 (실패 시 패닉)
str := val.(string)
fmt.Println(str) // "hello"

// 안전한 타입 단언 (ok 패턴)
str, ok := val.(string)
if ok {
    fmt.Println("문자열:", str)
} else {
    fmt.Println("문자열이 아니다")
}

// 잘못된 타입 단언
num, ok := val.(int)
fmt.Println(num, ok) // 0, false
```

> **주의**: ok 패턴 없이 타입 단언을 사용하면, 타입이 일치하지 않을 때 **패닉**이 발생한다. 따라서 확실하지 않은 경우 반드시 ok 패턴을 사용해야 한다.

---

## 18.5 인터페이스 변환하기 (타입 스위치)

### 타입 스위치 (Type Switch)

인터페이스 값의 **실제 타입에 따라 분기**하는 방법이다. 여러 타입을 검사해야 할 때 타입 단언을 반복하는 것보다 깔끔하다.

```go
func Describe(val any) {
    switch v := val.(type) {
    case int:
        fmt.Printf("정수: %d\n", v)
    case string:
        fmt.Printf("문자열: %s (길이: %d)\n", v, len(v))
    case bool:
        fmt.Printf("불리언: %t\n", v)
    case []int:
        fmt.Printf("정수 슬라이스: %v (길이: %d)\n", v, len(v))
    default:
        fmt.Printf("알 수 없는 타입: %T = %v\n", v, v)
    }
}
```

타입 스위치에서 `v`는 각 case 블록 내에서 해당 타입으로 자동 변환된다. `case int:` 블록 안에서 `v`는 `int` 타입이고, `case string:` 블록 안에서 `v`는 `string` 타입이다.

### 인터페이스 조합 (임베딩)

작은 인터페이스들을 조합하여 큰 인터페이스를 만들 수 있다. 이는 Go의 "구성(composition)" 철학을 잘 보여주는 패턴이다.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 인터페이스 임베딩으로 조합
type ReadWriter interface {
    Reader
    Writer
}
```

`ReadWriter` 인터페이스를 구현하려면 `Read()`와 `Write()` 메서드를 모두 구현해야 한다. 이렇게 작은 인터페이스를 조합하면 코드 재사용성이 높아지고, 각 인터페이스를 독립적으로 사용할 수도 있다.

### 인터페이스 변환

더 큰 인터페이스에서 더 작은 인터페이스로 변환할 수 있다. 이는 메서드 집합이 상위 집합인 타입은 하위 집합의 인터페이스를 자동으로 만족하기 때문이다.

```go
var rw ReadWriter = ... // ReadWriter 구현체
var r Reader = rw       // ReadWriter -> Reader (자동 변환)
```

반대로, 작은 인터페이스에서 큰 인터페이스로의 변환은 자동으로 되지 않으며 타입 단언이 필요하다.

### nil 인터페이스 주의사항

인터페이스의 nil에는 주의할 점이 있다. 인터페이스 변수가 nil인 것과, 인터페이스가 nil 포인터를 가리키는 것은 다르다.

```go
var s Shape           // nil 인터페이스 (타입 정보도 없음)
fmt.Println(s == nil) // true

var r *Rectangle      // nil 포인터
var s2 Shape = r      // nil이 아님! (타입 정보는 있음)
fmt.Println(s2 == nil) // false
```

이 차이를 이해하지 못하면 버그가 발생할 수 있으므로 주의해야 한다.

---

## 핵심 요약

| 개념 | 설명 |
|------|------|
| 인터페이스 | 메서드의 집합을 정의하는 타입이다 |
| 암묵적 구현 | 메서드만 맞으면 자동으로 인터페이스를 구현한다 (덕 타이핑) |
| 다형성 | 같은 인터페이스로 다양한 타입을 사용할 수 있다 |
| 빈 인터페이스 | `interface{}` / `any`로 모든 타입을 수용한다 |
| 타입 단언 | `val.(Type)`으로 구체적 타입을 추출한다 |
| 타입 스위치 | `switch v := val.(type)`으로 타입별 분기한다 |
| 인터페이스 조합 | 인터페이스 임베딩으로 큰 인터페이스를 구성한다 |

---

## 연습문제

### 문제 1: Shape 인터페이스
`Shape` 인터페이스를 정의하고, `Rectangle`, `Circle`, `Triangle` 구조체에서 구현하라. 모든 도형의 넓이 합계를 구하는 함수 `TotalArea(shapes []Shape) float64`를 작성하라.

### 문제 2: Stringer 인터페이스
`fmt.Stringer` 인터페이스(`String() string`)를 구현하는 `Student` 구조체를 만들고, `fmt.Println()`으로 출력해보라. `Stringer`를 구현하면 `fmt` 패키지의 출력 함수들이 자동으로 `String()` 메서드를 호출한다는 점을 확인하라.

### 문제 3: 타입 스위치
`any` 타입의 슬라이스를 받아 각 요소의 타입과 값을 출력하는 함수를 작성하라. 숫자 타입(`int`, `float64`)의 합계도 계산하라.

### 문제 4: 인터페이스 설계
`Logger` 인터페이스를 정의하고(`Log(message string)`), `ConsoleLogger`(콘솔 출력)와 `FileLogger`(파일명 출력) 두 가지 구현체를 만들라. 같은 코드에서 두 로거를 교체하며 사용해보라.

### 문제 5: 덕 타이핑 활용
`Closer` 인터페이스(`Close() error`)를 정의하고, 이를 구현하는 `DBConnection`, `FileHandle`, `NetworkSocket` 타입을 만들어 모두 닫는 함수 `CloseAll(closers []Closer) []error`를 작성하라.

### 문제 6: 빈 인터페이스와 타입 단언
`any` 타입의 값을 받아 정수면 2를 곱하고, 문자열이면 대문자로 변환하고, 불리언이면 반전시키는 함수 `Transform(val any) any`를 작성하라. 지원하지 않는 타입이면 원래 값을 그대로 반환하라.

### 문제 7: 인터페이스 조합
`Reader` 인터페이스(`Read() string`)와 `Writer` 인터페이스(`Write(data string)`)를 정의하고, 이를 조합한 `ReadWriter` 인터페이스를 만들라. `MemoryBuffer` 구조체가 `ReadWriter`를 구현하도록 하고, `Reader`나 `Writer`로도 사용할 수 있음을 보여라.

### 문제 8: 정렬 인터페이스 구현
`sort.Interface`(Len, Less, Swap)를 구현하는 `ByAge` 타입을 만들어 `Person` 구조체 슬라이스를 나이순으로 정렬하라. `sort.Sort()`와 `sort.Slice()` 두 가지 방식의 차이를 비교하라.

### 문제 9: nil 인터페이스 이해
nil 인터페이스와 nil 포인터를 가진 인터페이스의 차이를 보여주는 코드를 작성하라. 함수가 `error` 인터페이스를 반환할 때 nil 포인터 문제가 발생하는 예시를 만들고, 이를 올바르게 수정하라.

### 문제 10: 다중 인터페이스 구현
하나의 구조체가 여러 인터페이스를 동시에 구현하는 예시를 작성하라. `Serializer`(`Serialize() string`), `Validator`(`Validate() error`), `Stringer`(`String() string`) 인터페이스를 정의하고, `User` 구조체가 세 인터페이스를 모두 구현하도록 하라.

---

## 구현 과제

### 과제 1: 플러그인 시스템
`Plugin` 인터페이스(`Name() string`, `Execute(input string) string`)를 정의하고, 여러 플러그인(대문자 변환, 단어 수 세기, 문자열 뒤집기 등)을 구현하라. `PluginManager` 구조체를 만들어 플러그인 등록, 이름으로 검색, 실행 기능을 구현하라.

### 과제 2: 결제 시스템
`PaymentProcessor` 인터페이스(`ProcessPayment(amount float64) error`, `Refund(amount float64) error`, `Name() string`)를 정의하고, `CreditCard`, `BankTransfer`, `DigitalWallet` 구현체를 만들라. `PaymentService` 구조체에서 여러 결제 수단을 관리하고, 결제 처리 및 로깅 기능을 구현하라.

### 과제 3: 알림 시스템
`Notifier` 인터페이스(`Send(to string, message string) error`)를 정의하고, `EmailNotifier`, `SMSNotifier`, `SlackNotifier` 구현체를 만들라. 하나의 알림을 여러 채널로 동시에 보내는 `MultiNotifier`도 구현하라(데코레이터 패턴).

### 과제 4: 데이터 저장소 추상화
`Store` 인터페이스(`Get(key string) (string, error)`, `Set(key, value string) error`, `Delete(key string) error`, `Keys() []string`)를 정의하고, `MemoryStore`(맵 기반)와 `FileStore`(파일 기반, 시뮬레이션 가능)를 구현하라. 어떤 Store를 사용하든 동일하게 동작하는 `CacheService`를 만들라.

### 과제 5: 미들웨어 체인
`Handler` 인터페이스(`Handle(request string) string`)를 정의하고, 미들웨어 패턴을 구현하라. `LoggingMiddleware`, `AuthMiddleware`, `TimingMiddleware` 등의 미들웨어가 요청을 가로채어 처리한 후 다음 핸들러에게 전달하는 체인 구조를 만들라.

---

## 프로젝트 과제

### 프로젝트 1: 파일 변환기 프레임워크
`Converter` 인터페이스(`Convert(input []byte) ([]byte, error)`, `InputFormat() string`, `OutputFormat() string`)를 정의하여 파일 형식 변환 프레임워크를 구현하라. CSV->JSON, JSON->CSV, 대문자 변환, Base64 인코딩/디코딩 등 다양한 변환기를 구현하고, 변환기를 체이닝하여 여러 변환을 순차적으로 적용할 수 있도록 하라. 변환기 레지스트리를 만들어 입력/출력 형식으로 적절한 변환기를 자동으로 찾는 기능도 추가하라.

### 프로젝트 2: 간단한 의존성 주입 컨테이너
인터페이스를 활용하여 간단한 DI(Dependency Injection) 컨테이너를 구현하라. 컨테이너에 인터페이스와 구현체를 등록하고, 필요할 때 구현체를 꺼내서 사용할 수 있도록 하라. 예를 들어 `UserService`가 `Database` 인터페이스에 의존할 때, 프로덕션에서는 `PostgresDB`를, 테스트에서는 `MockDB`를 주입하는 시나리오를 구현하라. 타입 단언과 빈 인터페이스를 적절히 활용하라.
