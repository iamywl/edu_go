# Chapter 17: 메서드

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch17_method/basic.go
go run ch17_method/pointer_vs_value.go
```

> **Makefile 활용**: `make run CH=ch17_method` 또는 `make run CH=ch17_method FILE=basic.go`

---

Go 언어에서 **메서드(Method)**는 특정 타입에 속한 함수이다. 구조체와 메서드를 결합하면 데이터와 기능을 하나로 묶어 객체 지향적인 프로그래밍을 할 수 있다. Go는 전통적인 클래스 기반 상속 대신 메서드와 인터페이스를 통해 다형성과 캡슐화를 구현한다.

---

## 17.1 메서드 선언 (리시버 함수)

### 메서드란?

메서드는 **리시버(receiver)**가 있는 함수이다. 리시버는 메서드가 어떤 타입에 속하는지를 나타낸다. 일반 함수와 달리 메서드는 특정 타입에 "소속"되어 있으며, 해당 타입의 값을 통해 호출한다. 이를 통해 데이터와 동작을 자연스럽게 연결할 수 있다.

### 기본 문법

```go
func (리시버변수 리시버타입) 메서드이름(매개변수) 반환타입 {
    // 메서드 본문
}
```

리시버는 함수 이름 앞에 괄호로 감싸서 선언한다. 이것이 일반 함수와 메서드를 구분짓는 유일한 문법적 차이이다.

### 예시

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Rectangle 타입의 메서드
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// 사용
rect := Rectangle{Width: 10, Height: 5}
fmt.Println(rect.Area())       // 50
fmt.Println(rect.Perimeter())  // 30
```

### 리시버 이름 관례

- 리시버 변수 이름은 타입의 **첫 글자를 소문자**로 사용한다.
- `Rectangle` -> `r`, `Student` -> `s`, `Account` -> `a`
- `this`나 `self`를 사용하지 않는 것이 Go의 관례이다. 이는 Go 커뮤니티에서 강하게 지키는 관례이므로 반드시 따르는 것이 좋다.
- 같은 타입의 모든 메서드에서 리시버 이름을 통일해야 한다.

### 어떤 타입에 메서드를 붙일 수 있나?

메서드를 정의할 수 있는 타입에는 제한이 있다. 같은 패키지에 정의된 타입에만 메서드를 추가할 수 있다.

```go
// 구조체에 메서드 추가
type Circle struct {
    Radius float64
}
func (c Circle) Area() float64 { ... }

// 사용자 정의 타입에도 가능
type MyInt int
func (m MyInt) Double() MyInt { return m * 2 }

// 기본 타입(int, string 등)에는 직접 추가 불가
// 같은 패키지에 정의된 타입에만 메서드 추가 가능
```

외부 패키지에 정의된 타입에 메서드를 추가하고 싶다면, 해당 타입을 래핑(wrapping)하는 새로운 타입을 정의하면 된다. 예를 들어 `type MyString string`과 같이 선언한 후 `MyString`에 메서드를 추가할 수 있다.

---

## 17.2 메서드는 왜 필요한가? (데이터와 기능 결합)

### 함수 vs 메서드

```go
// 일반 함수 방식
func CalculateArea(r Rectangle) float64 {
    return r.Width * r.Height
}
area := CalculateArea(rect)

// 메서드 방식
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
area := rect.Area()
```

두 방식 모두 동일한 결과를 반환하지만, 메서드 방식은 "이 동작이 Rectangle에 속한다"는 의미를 명확히 전달한다.

### 메서드의 장점

1. **코드 가독성**: `rect.Area()`가 `CalculateArea(rect)`보다 읽기 쉽다. 주어-동사 형태로 읽히기 때문이다.
2. **데이터와 기능의 결합**: 관련된 동작을 타입에 묶어 관리할 수 있다. 타입의 정의와 동작을 한 곳에서 파악할 수 있어 코드 이해도가 높아진다.
3. **네임스페이스 분리**: 다른 타입에 같은 이름의 메서드를 정의할 수 있다. `Rectangle.Area()`와 `Circle.Area()`는 이름이 같지만 충돌하지 않는다.
4. **인터페이스 구현**: 메서드를 통해 인터페이스를 구현한다 (18장에서 학습). 이것이 Go의 다형성을 가능하게 하는 핵심 메커니즘이다.

```go
// 같은 이름의 메서드를 다른 타입에 정의 가능
func (r Rectangle) Area() float64 { return r.Width * r.Height }
func (c Circle) Area() float64    { return math.Pi * c.Radius * c.Radius }
```

### 메서드를 활용한 설계

메서드를 잘 설계하면 타입이 자체적으로 자신의 상태를 관리하도록 할 수 있다. 외부에서는 메서드만 호출하면 되므로 내부 구현을 숨길 수 있다.

```go
type Account struct {
    Owner   string
    Balance int
}

func (a *Account) Deposit(amount int) {
    a.Balance += amount
}

func (a *Account) Withdraw(amount int) error {
    if a.Balance < amount {
        return fmt.Errorf("잔액 부족: %d원 필요, %d원 보유", amount, a.Balance)
    }
    a.Balance -= amount
    return nil
}

func (a Account) String() string {
    return fmt.Sprintf("[%s] 잔액: %d원", a.Owner, a.Balance)
}
```

위 예시에서 `Deposit()`과 `Withdraw()`는 잔액을 변경해야 하므로 포인터 리시버를 사용하고, `String()`은 읽기만 하므로 값 리시버를 사용한다. 그러나 실무에서는 일관성을 위해 한 타입의 모든 메서드를 포인터 리시버로 통일하는 경우가 많다.

---

## 17.3 포인터 메서드 vs 값 타입 메서드

### 값 타입 리시버

```go
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
```

- 리시버의 **복사본**이 전달된다.
- 메서드 내에서 원본을 변경할 수 **없다**.
- 읽기 전용 작업에 적합하다.
- 구조체가 작고(필드 몇 개 이하) 변경이 불필요할 때 사용한다.

### 포인터 리시버

```go
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}
```

- 리시버의 **포인터**가 전달된다.
- 메서드 내에서 원본을 **직접 변경**할 수 있다.
- 구조체가 클 때 복사 비용을 줄일 수 있다.
- 포인터 리시버는 8바이트(64비트 시스템)만 복사되므로, 큰 구조체에서 성능상 유리하다.

### 비교 정리

| 구분 | 값 타입 리시버 | 포인터 리시버 |
|------|---------------|--------------|
| 문법 | `func (r Rect)` | `func (r *Rect)` |
| 원본 변경 | 불가 | 가능 |
| 복사 | 전체 복사 | 포인터만 복사 (8바이트) |
| 사용 시점 | 읽기 전용, 작은 구조체 | 값 수정 필요시, 큰 구조체 |
| nil 호출 | 불가능 | 가능 (nil 체크 필요) |

### 자동 변환

Go는 메서드 호출 시 **자동으로 포인터와 값을 변환**한다. 이 기능 덕분에 프로그래머는 리시버 타입을 크게 신경 쓰지 않고 메서드를 호출할 수 있다:

```go
rect := Rectangle{Width: 10, Height: 5}

// 값 타입 변수로 포인터 메서드 호출 가능
rect.Scale(2)  // Go가 자동으로 (&rect).Scale(2)로 변환

// 포인터로 값 타입 메서드 호출도 가능
p := &rect
fmt.Println(p.Area())  // Go가 자동으로 (*p).Area()로 변환
```

단, 이 자동 변환은 **주소를 얻을 수 있는(addressable) 값**에서만 동작한다. 맵의 값이나 함수 반환값처럼 주소를 얻을 수 없는 경우에는 포인터 리시버 메서드를 직접 호출할 수 없다.

### 포인터 리시버를 사용해야 하는 경우

1. 메서드가 리시버의 **값을 변경**해야 할 때
2. 리시버가 **큰 구조체**일 때 (복사 비용 절감)
3. **일관성**: 한 타입의 메서드 중 하나라도 포인터 리시버를 사용하면, 나머지도 포인터 리시버로 통일하는 것이 좋다
4. 리시버가 `sync.Mutex` 등 **복사하면 안 되는 필드**를 포함할 때

---

## 핵심 요약

| 개념 | 설명 |
|------|------|
| 메서드 | 리시버가 있는 함수이다. `func (r Type) Name()` 형태로 선언한다 |
| 리시버 | 메서드가 속하는 타입을 지정한다. 타입 첫 글자를 소문자로 사용한다 |
| 값 리시버 | 복사본으로 동작하며, 원본 변경이 불가하다 |
| 포인터 리시버 | 원본을 참조하며, 값 변경이 가능하고 큰 구조체에 효율적이다 |
| 자동 변환 | Go가 값/포인터 간 자동 변환을 처리한다 |

---

## 연습문제

### 문제 1: 기본 메서드
`Circle` 구조체에 `Area()`와 `Circumference()`(둘레) 메서드를 구현하라. `math.Pi`를 사용하라.

### 문제 2: 포인터 리시버
`Counter` 구조체에 `Increment()`, `Decrement()`, `Reset()`, `Value()` 메서드를 구현하라. 어떤 메서드가 포인터 리시버여야 하는지 생각해보라.

### 문제 3: 은행 계좌
`BankAccount` 구조체에 입금(`Deposit`), 출금(`Withdraw`), 잔액조회(`Balance`), 거래내역(`String`) 메서드를 구현하라. 출금 시 잔액이 부족하면 에러를 반환하라.

### 문제 4: 사용자 정의 타입
`type Celsius float64`와 `type Fahrenheit float64` 타입을 만들고, 각각 다른 온도 단위로 변환하는 메서드를 구현하라. (`ToFahrenheit()`, `ToCelsius()`)

### 문제 5: 문자열 빌더
`type StringBuilder struct { data []byte }` 타입을 정의하고, `Append(s string)`, `Prepend(s string)`, `String() string`, `Len() int`, `Clear()` 메서드를 구현하라.

### 문제 6: 값 리시버 vs 포인터 리시버 관찰
`Point` 구조체(`X, Y float64`)를 만들고, 값 리시버 메서드 `MoveByValue(dx, dy float64)`와 포인터 리시버 메서드 `MoveByPointer(dx, dy float64)`를 각각 구현하라. 두 메서드를 호출한 후 원본 값이 어떻게 변하는지(또는 변하지 않는지) 확인하는 코드를 작성하라.

### 문제 7: 메서드 체이닝
`Query` 구조체를 만들고, `Select(fields string)`, `From(table string)`, `Where(condition string)` 메서드를 포인터 리시버로 구현하여 메서드 체이닝이 가능하도록 하라. 각 메서드가 `*Query`를 반환하여 `q.Select("*").From("users").Where("age > 20")` 형태로 사용할 수 있어야 한다.

### 문제 8: 컬렉션 타입 메서드
`type IntSlice []int` 타입을 정의하고, `Sum() int`, `Avg() float64`, `Max() int`, `Min() int`, `Contains(val int) bool`, `Filter(f func(int) bool) IntSlice` 메서드를 구현하라.

### 문제 9: 시간 관련 타입
`Duration` 구조체(시, 분, 초 필드)를 만들고, `TotalSeconds() int`, `Add(other Duration) Duration`, `String() string` 메서드를 구현하라. `String()`은 `"2h 30m 15s"` 형식으로 반환해야 한다.

### 문제 10: nil 리시버 처리
포인터 리시버 메서드에서 리시버가 `nil`인 경우를 안전하게 처리하는 방법을 구현하라. `SafeList` 구조체를 만들고, `nil`인 `*SafeList`에서도 `Len()`, `String()` 메서드가 패닉 없이 동작하도록 하라.

---

## 구현 과제

### 과제 1: 벡터(Vector) 타입
`Vector2D` 구조체(`X, Y float64`)를 구현하고, 다음 메서드를 추가하라:
- `Add(other Vector2D) Vector2D` : 벡터 덧셈
- `Sub(other Vector2D) Vector2D` : 벡터 뺄셈
- `Scale(factor float64) Vector2D` : 스칼라 곱
- `Magnitude() float64` : 크기(길이)
- `Normalize() Vector2D` : 단위 벡터
- `Dot(other Vector2D) float64` : 내적
- `Distance(other Vector2D) float64` : 두 벡터 사이의 거리

### 과제 2: 링크드 리스트 구현
`Node`와 `LinkedList` 구조체를 사용하여 단방향 연결 리스트를 직접 구현하라. `Append(val int)`, `Prepend(val int)`, `Delete(val int) bool`, `Find(val int) bool`, `Len() int`, `String() string` 메서드를 포함하라.

### 과제 3: 날짜(Date) 타입
`Date` 구조체(Year, Month, Day)를 만들고, 다음 메서드를 구현하라:
- `IsValid() bool` : 유효한 날짜인지 확인
- `DaysInMonth() int` : 해당 월의 일수 반환
- `IsLeapYear() bool` : 윤년 여부 확인
- `AddDays(n int) Date` : n일 후의 날짜 반환
- `DaysBetween(other Date) int` : 두 날짜 사이의 일수 반환
- `String() string` : "2024-01-15" 형식으로 반환

### 과제 4: 행렬(Matrix) 타입
`Matrix` 구조체를 만들고, 다음 메서드를 구현하라:
- `NewMatrix(rows, cols int) *Matrix` : 생성자 함수
- `Set(row, col int, val float64)` : 값 설정
- `Get(row, col int) float64` : 값 읽기
- `Add(other *Matrix) *Matrix` : 행렬 덧셈
- `Transpose() *Matrix` : 전치 행렬
- `String() string` : 보기 좋은 형식으로 출력

### 과제 5: 간단한 JSON 빌더
`JSONBuilder` 구조체를 만들어 JSON 문자열을 단계적으로 구축하는 메서드 체이닝 패턴을 구현하라. `StartObject()`, `EndObject()`, `AddString(key, val string)`, `AddInt(key string, val int)`, `AddBool(key string, val bool)`, `Build() string` 메서드를 구현하라.

---

## 프로젝트 과제

### 프로젝트 1: 도형 계산기
다양한 도형 타입(`Rectangle`, `Circle`, `Triangle`, `Trapezoid` 등)을 정의하고, 각 도형에 `Area()`, `Perimeter()`, `String()`, `Scale(factor float64)` 메서드를 구현하라. 사용자로부터 도형 종류와 치수를 입력받아 넓이와 둘레를 계산하고, 여러 도형의 넓이 합계와 가장 큰 도형을 찾는 프로그램을 작성하라.

### 프로젝트 2: 은행 시뮬레이션
`Bank` 구조체와 `Account` 구조체를 설계하여 은행 시뮬레이션 프로그램을 작성하라. 다음 기능을 포함해야 한다:
- 계좌 개설, 폐쇄
- 입금, 출금, 계좌 이체
- 거래 내역 조회 (각 거래를 `Transaction` 구조체로 기록)
- 전체 고객 목록 조회
- 잔액 기준 정렬
- 총 자산 통계

각 구조체에 적절한 메서드를 정의하고, 값 리시버와 포인터 리시버를 상황에 맞게 사용하라.
