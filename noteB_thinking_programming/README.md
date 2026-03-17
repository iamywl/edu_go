# 노트 B. 생각하는 프로그래밍

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run noteB_thinking_programming/oop_go.go
go run noteB_thinking_programming/constructor.go
go run noteB_thinking_programming/pointer_copy.go
go run noteB_thinking_programming/dependency_inversion.go
go run noteB_thinking_programming/gc_demo.go
```

> **Makefile 활용**: `make run CH=noteB_thinking_programming` 또는 `make run CH=noteB_thinking_programming FILE=oop_go.go`

---

Go 언어를 더 깊이 이해하기 위한 사고 실험과 설계 원칙을 다룬다.

---

## B.1 Go는 객체지향 언어인가?

### 공식 답변

Go 공식 FAQ에서는 이렇게 답한다:

> "Yes and no." (맞기도 하고 아니기도 하다.)

Go에는 클래스가 없고 상속도 없다. 하지만 객체지향 프로그래밍의 핵심 개념은 다른 방식으로 지원한다. 이 질문에 대한 답은 "객체지향"을 어떻게 정의하느냐에 달려 있다.

전통적 OOP의 네 가지 기둥과 Go의 대응 방식:
- **캡슐화**: 대문자/소문자 접근 제어로 지원한다.
- **추상화**: 인터페이스로 지원한다.
- **다형성**: 인터페이스의 암시적 구현으로 지원한다.
- **상속**: 지원하지 않는다. 대신 **조합(composition)**을 사용한다.

### 클래스 없는 OOP

Go는 **구조체 + 메서드 + 인터페이스**로 OOP를 구현한다:

```go
// 클래스 대신 구조체
type Dog struct {
    Name string
    Age  int
}

// 메서드를 구조체에 연결
func (d Dog) Bark() string {
    return d.Name + ": 멍멍!"
}
```

Go에서 메서드는 구조체뿐 아니라 모든 명명된 타입에 정의할 수 있다. 이것은 전통적 OOP 언어에서 클래스만 메서드를 가질 수 있는 것과 다른 점이다:

```go
type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
    return float64(c)*9/5 + 32
}
```

### 상속 대신 조합 (Composition over Inheritance)

Go는 상속 대신 **임베딩(embedding)**을 사용한다:

```go
// "부모 클래스" 역할
type Animal struct {
    Name string
}

func (a Animal) Breathe() string {
    return a.Name + "이(가) 숨을 쉰다"
}

// "자식 클래스" 역할 — 상속이 아닌 조합!
type Cat struct {
    Animal        // 임베딩: Animal의 필드와 메서드를 "승격"
    Indoor bool
}

func main() {
    cat := Cat{
        Animal: Animal{Name: "나비"},
        Indoor: true,
    }
    cat.Breathe()  // Animal의 메서드를 직접 호출 가능
    cat.Name       // Animal의 필드에 직접 접근 가능
}
```

**상속과의 차이점:**
- 임베딩은 "is-a" 관계가 아닌 "has-a" 관계이다.
- `Cat`은 `Animal` 타입이 아니다. `Animal`을 포함하고 있을 뿐이다.
- 메서드 오버라이딩이 가능하지만, 다형성은 인터페이스로 구현한다.

**임베딩의 메서드 승격 규칙:**

임베딩된 타입의 메서드는 외부 타입으로 "승격"된다. 하지만 외부 타입이 같은 이름의 메서드를 정의하면 외부 타입의 메서드가 우선한다:

```go
func (c Cat) Breathe() string {
    return c.Name + "이(가) 조용히 숨을 쉰다"  // Animal.Breathe를 가림
}

cat.Breathe()           // Cat의 Breathe 호출
cat.Animal.Breathe()    // Animal의 Breathe를 명시적으로 호출
```

**다중 임베딩도 가능하다:**

```go
type Logger struct{}
func (l Logger) Log(msg string) { fmt.Println(msg) }

type Validator struct{}
func (v Validator) Validate() bool { return true }

type Service struct {
    Logger
    Validator
}
// Service는 Log()와 Validate() 모두 사용 가능
```

### 다형성은 인터페이스로

```go
type Speaker interface {
    Speak() string
}

type Dog struct{ Name string }
func (d Dog) Speak() string { return d.Name + ": 멍멍!" }

type Cat struct{ Name string }
func (c Cat) Speak() string { return c.Name + ": 야옹~" }

// Speaker 인터페이스를 만족하는 모든 타입을 받을 수 있음
func MakeNoise(s Speaker) {
    fmt.Println(s.Speak())
}
```

Go의 인터페이스는 **암시적 구현(implicit implementation)**이다. `implements` 키워드 없이, 메서드만 구현하면 자동으로 인터페이스를 만족한다. 이 설계 덕분에:

- 서드파티 라이브러리의 타입도 내가 정의한 인터페이스를 만족시킬 수 있다.
- 구현과 인터페이스 사이의 결합도가 낮아진다.
- 인터페이스를 사후에 추가할 수 있다. (기존 코드를 수정할 필요가 없다.)

> 자세한 예제는 `oop_go.go`를 참고한다.

---

## B.2 구조체에 생성자를 둘 수 있나?

Go에는 생성자 문법이 없다. 대신 **`NewXxx` 함수 패턴**을 관례적으로 사용한다. 이것은 Go 표준 라이브러리 전체에서 일관되게 사용되는 패턴이다.

### NewXxx 패턴

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

// 생성자 역할의 함수
func NewServer(host string, port int) *Server {
    return &Server{
        host:    host,
        port:    port,
        timeout: 30 * time.Second, // 기본값 설정
    }
}
```

### 왜 이 패턴을 사용하는가?

1. **유효성 검증**: 잘못된 값으로 객체가 생성되는 것을 방지한다.

```go
func NewPort(port int) (*Port, error) {
    if port < 1 || port > 65535 {
        return nil, fmt.Errorf("유효하지 않은 포트: %d", port)
    }
    return &Port{number: port}, nil
}
```

2. **비공개 필드 초기화**: 소문자 필드는 외부에서 직접 설정 불가하다.

```go
type logger struct {
    prefix string    // 비공개 — 외부에서 접근 불가
    output io.Writer // 비공개
}

func NewLogger(prefix string) *logger {
    return &logger{
        prefix: prefix,
        output: os.Stdout,
    }
}
```

3. **기본값 제공**: 제로값이 의미 없는 경우 적절한 기본값을 설정한다.

4. **인터페이스 반환**: 구체 타입 대신 인터페이스를 반환하여 구현을 숨길 수 있다.

```go
type Store interface {
    Get(key string) (string, error)
    Set(key, value string) error
}

// 구체 타입이 아닌 인터페이스를 반환
func NewStore(driver string) (Store, error) {
    switch driver {
    case "memory":
        return &memoryStore{data: make(map[string]string)}, nil
    case "redis":
        return &redisStore{}, nil
    default:
        return nil, fmt.Errorf("지원하지 않는 드라이버: %s", driver)
    }
}
```

### 제로값이 유용한 설계

Go에서 좋은 설계란 구조체의 제로값이 곧바로 사용 가능한 유효한 상태인 것이다. 표준 라이브러리에서 이런 예를 많이 찾을 수 있다:

```go
var mu sync.Mutex      // 제로값이 곧 잠금 해제 상태
mu.Lock()              // NewMutex() 같은 생성자가 불필요

var buf bytes.Buffer   // 제로값이 곧 빈 버퍼
buf.WriteString("hi")  // NewBuffer() 없이 바로 사용 가능
```

제로값이 유효하도록 설계하면 `NewXxx` 함수가 필요 없어져 API가 간결해진다.

### Functional Options 패턴

선택적 설정이 많을 때 사용하는 고급 패턴이다. 이 패턴은 Rob Pike가 제안한 것으로, Go 커뮤니티에서 널리 사용된다:

```go
type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithMaxConns(n int) Option {
    return func(s *Server) {
        s.maxConns = n
    }
}

func NewServer(host string, port int, opts ...Option) *Server {
    s := &Server{
        host:     host,
        port:     port,
        timeout:  30 * time.Second,
        maxConns: 100,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// 사용 예
s := NewServer("localhost", 8080,
    WithTimeout(60 * time.Second),
    WithMaxConns(500),
)
```

**Functional Options 패턴의 장점:**
- 기본값이 명확하게 드러난다.
- 새로운 옵션을 추가해도 기존 코드가 깨지지 않는다 (하위 호환성).
- 옵션 이름이 자체 문서 역할을 한다.
- 옵션에 유효성 검증 로직을 포함할 수 있다.

> 자세한 예제는 `constructor.go`를 참고한다.

---

## B.3 포인터를 사용해도 복사가 일어나나?

**그렇다, 포인터 값 자체는 복사된다.** 하지만 포인터가 가리키는 데이터는 복사되지 않는다. 이 구분을 정확히 이해하는 것이 Go 프로그래밍의 핵심이다.

### 값 전달 vs 포인터 전달

```go
type BigData struct {
    data [1000000]int  // 약 8MB
}

// 값 전달: BigData 전체가 복사됨 (8MB 복사!)
func processValue(d BigData) {
    d.data[0] = 999  // 복사본을 수정 → 원본에 영향 없음
}

// 포인터 전달: 포인터만 복사됨 (8바이트 복사)
func processPointer(d *BigData) {
    d.data[0] = 999  // 원본을 수정!
}
```

### 포인터 복사의 의미

```go
func main() {
    original := &BigData{}
    copyOfPtr := original  // 포인터 값(메모리 주소)이 복사됨

    // original과 copyOfPtr는 같은 객체를 가리킴
    fmt.Println(original == copyOfPtr) // true

    // 하지만 포인터 변수 자체는 별개
    copyOfPtr = &BigData{}              // copyOfPtr만 변경
    fmt.Println(original == copyOfPtr)  // false — original은 그대로
}
```

### 그림으로 이해하기

```
포인터 전달 시:

  original ──────┐
                 ▼
              ┌──────────┐
              │ BigData  │  ← 실제 데이터 (하나만 존재)
              └──────────┘
                 ▲
  copyOfPtr ─────┘

포인터 값(주소)은 두 개지만, 가리키는 데이터는 하나이다.
```

### Go에서 모든 것은 값으로 전달된다

Go의 핵심 원칙 중 하나는 **모든 함수 인수가 값으로 전달(pass by value)**된다는 것이다. 포인터를 전달할 때도 포인터 값(메모리 주소)이 복사되는 것이지, 참조(reference)가 전달되는 것이 아니다. 이것은 C++의 참조 전달(pass by reference)과는 다른 개념이다.

```go
func changePointer(p *int) {
    val := 42
    p = &val  // 로컬 변수 p만 변경됨, 호출자의 포인터는 변경되지 않음
}

func main() {
    x := 10
    ptr := &x
    changePointer(ptr)
    fmt.Println(*ptr)  // 여전히 10
}
```

포인터가 가리키는 **값**은 변경할 수 있지만, 포인터 **자체**를 변경하는 것은 호출자에게 전파되지 않는다.

### 슬라이스도 비슷하다

슬라이스는 내부적으로 `(ptr, len, cap)` 세 필드를 가진 구조체이다. 함수에 전달하면 이 세 필드가 복사되지만, 내부 배열은 공유된다.

```go
func modify(s []int) {
    s[0] = 999  // 내부 배열 공유 → 원본 변경됨
    s = append(s, 100)  // 새로운 슬라이스 헤더 → 원본에 영향 없음
}
```

맵과 채널도 내부적으로 포인터이므로 함수에 전달하면 같은 데이터를 공유한다:

```go
func modifyMap(m map[string]int) {
    m["new"] = 42  // 원본 맵에도 반영됨
}
```

> 자세한 예제는 `pointer_copy.go`를 참고한다.

---

## B.4 값 타입을 쓸 것인가? 포인터를 쓸 것인가?

이 질문은 Go 개발에서 가장 자주 마주치는 설계 결정 중 하나이다. 명확한 기준을 갖고 있으면 일관된 코드를 작성할 수 있다.

### 값 타입을 사용하는 경우

1. **작은 구조체** (필드 2~3개 이하)
2. **불변 데이터** (한번 생성 후 변경하지 않음)
3. **맵의 키**로 사용하는 경우
4. **동시성 안전**이 필요한 경우 (복사본이므로 경합 없음)

```go
type Point struct {
    X, Y float64
}

type Color struct {
    R, G, B uint8
}
```

이런 작은 구조체는 값으로 전달해도 복사 비용이 미미하며, 복사본 덕분에 부작용(side effect)이 없는 코드를 작성할 수 있다.

### 포인터를 사용하는 경우

1. **큰 구조체** (복사 비용이 큰 경우)
2. **메서드에서 상태를 변경**해야 하는 경우
3. **nil이 의미 있는 값**인 경우 ("값 없음"을 표현)
4. **공유 상태**가 필요한 경우

```go
type DatabaseConnection struct {
    pool    *sql.DB
    config  Config
    metrics Metrics
}

// 상태를 변경하는 메서드 → 포인터 리시버
func (db *DatabaseConnection) Close() error {
    return db.pool.Close()
}
```

### "큰 구조체"의 기준

정확한 바이트 수 기준은 없지만, 일반적인 가이드라인은 다음과 같다:
- 필드가 4~5개 이하이고 모두 기본 타입이면 값 타입을 사용한다.
- 슬라이스, 맵, 문자열 등 참조 타입 필드가 있으면 이미 내부적으로 포인터이므로 값 타입으로 전달해도 큰 비용이 아니다.
- 확실하지 않으면 벤치마크로 측정한다. `go test -bench`가 항상 정답을 알려준다.

### 가이드라인 요약

| 상황 | 선택 | 이유 |
|------|------|------|
| 메서드가 상태를 변경 | `*T` (포인터) | 원본 수정 필요 |
| 구조체가 큰 경우 | `*T` (포인터) | 복사 비용 절감 |
| `nil`이 필요한 경우 | `*T` (포인터) | 값 타입은 nil 불가 |
| 동시성 안전 필요 | `T` (값) | 복사본으로 격리 |
| 맵 키로 사용 | `T` (값) | 포인터 키는 주소 비교 |
| 소규모 불변 데이터 | `T` (값) | 단순하고 안전 |
| `sync.Mutex` 포함 | `*T` (포인터) | Mutex는 복사 금지 |
| 인터페이스 구현 | `*T` (포인터) | 대부분의 경우 포인터 리시버 |

### 일관성 규칙

한 타입의 메서드 리시버는 **모두 값 리시버** 또는 **모두 포인터 리시버**로 통일하는 것이 좋다. 섞어 쓰면 인터페이스 만족 여부가 혼란스러워진다.

이유를 구체적으로 설명하면:
- 값 리시버 메서드는 `T`와 `*T` 모두에서 호출 가능하다.
- 포인터 리시버 메서드는 `*T`에서만 호출 가능하다.
- 인터페이스 변수에 `T` 값을 저장하면 포인터 리시버 메서드를 호출할 수 없다.

```go
type Counter struct{ n int }

func (c *Counter) Increment() { c.n++ }  // 포인터 리시버
func (c Counter) Value() int { return c.n }  // 값 리시버 — 혼합 사용

type Incrementer interface {
    Increment()
}

var inc Incrementer = Counter{}   // 컴파일 에러! *Counter만 가능
var inc Incrementer = &Counter{}  // OK
```

---

## B.5 구체화된 객체와 관계하라고?

이것은 SOLID 원칙 중 **의존성 역전 원칙(Dependency Inversion Principle, DIP)**과 관련된다.

> "추상에 의존하라. 구체에 의존하지 마라."

이 원칙의 핵심은 상위 수준 모듈이 하위 수준 모듈에 직접 의존해서는 안 되며, 둘 다 추상(인터페이스)에 의존해야 한다는 것이다.

### 나쁜 예: 구체 타입에 의존

```go
type MySQLDatabase struct {
    conn *sql.DB
}

func (db *MySQLDatabase) FindUser(id int) (*User, error) {
    // MySQL 전용 쿼리
}

type UserService struct {
    db *MySQLDatabase  // 구체 타입에 직접 의존!
}
```

이 설계의 문제점:
- `UserService`를 테스트하려면 실제 MySQL이 필요하다.
- PostgreSQL로 교체하려면 `UserService`를 수정해야 한다.
- `UserService`의 단위 테스트가 불가능하다 (통합 테스트만 가능하다).

### 좋은 예: 인터페이스에 의존

```go
// 인터페이스 정의 (추상)
type UserRepository interface {
    FindUser(id int) (*User, error)
    SaveUser(user *User) error
}

// 서비스는 인터페이스에 의존
type UserService struct {
    repo UserRepository  // 인터페이스에 의존!
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

이제 다양한 구현체를 주입할 수 있다:

```go
// 프로덕션: MySQL 사용
mysqlRepo := NewMySQLRepository(db)
service := NewUserService(mysqlRepo)

// 테스트: 가짜(mock) 사용
mockRepo := &MockRepository{users: map[int]*User{}}
testService := NewUserService(mockRepo)
```

### Go에서의 인터페이스 설계 원칙

1. **인터페이스는 사용하는 쪽에서 정의하라**

```go
// 나쁜 예: 구현 패키지에서 거대한 인터페이스 정의
package database
type Repository interface {
    FindUser(id int) (*User, error)
    SaveUser(user *User) error
    DeleteUser(id int) error
    ListUsers() ([]*User, error)
    // ... 20개 메서드
}

// 좋은 예: 사용하는 쪽에서 필요한 메서드만 정의
package userservice
type UserFinder interface {
    FindUser(id int) (*User, error)
}
```

이 원칙은 Go의 암시적 인터페이스 구현 덕분에 가능하다. Java나 C#에서는 구현 측에서 `implements`를 선언해야 하므로 인터페이스를 먼저 정의해야 하지만, Go에서는 사후에 인터페이스를 정의해도 기존 타입이 자동으로 만족한다.

2. **작은 인터페이스를 선호하라**

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 필요하면 조합
type ReadWriter interface {
    Reader
    Writer
}
```

Go 표준 라이브러리의 인터페이스는 대부분 메서드가 1~2개이다. 이를 **인터페이스 분리 원칙(ISP)**이라고도 한다. 작은 인터페이스가 좋은 이유는:
- 구현이 쉽다.
- 목(mock) 객체를 만들기 쉽다.
- 조합으로 큰 인터페이스를 만들 수 있다.
- 더 많은 타입이 만족할 수 있다.

3. **빈 인터페이스(`interface{}` / `any`)는 최소한으로 사용하라**

빈 인터페이스는 타입 안전성을 포기하는 것이다. 제네릭이 Go 1.18에 도입된 이후, 이전에 `interface{}`를 사용하던 많은 경우를 타입 파라미터로 대체할 수 있다.

> 자세한 예제는 `dependency_inversion.go`를 참고한다.

---

## B.6 Go 언어 가비지 컬렉터

Go는 자동 메모리 관리를 위해 가비지 컬렉터(GC)를 사용한다. Go의 GC는 **동시 삼색 마킹(Concurrent Tricolor Mark-and-Sweep)** 알고리즘을 사용한다. Go 팀은 처리량(throughput)보다 지연 시간(latency)을 우선시하는 설계 결정을 내렸으며, 이는 웹 서버와 같은 대화형 서비스에 적합한 선택이다.

### GC의 기본 원리

GC는 더 이상 접근할 수 없는 메모리를 자동으로 회수한다:

```go
func createData() {
    data := make([]byte, 1024*1024)  // 1MB 할당
    _ = data
}  // 함수 종료 후 data는 접근 불가 → GC가 회수
```

"접근할 수 없다"는 것은 프로그램의 어떤 변수에서도 해당 메모리에 도달하는 경로가 없다는 의미이다. 루트(스택 변수, 전역 변수)에서 시작하여 포인터를 따라가며 도달 가능한 모든 객체를 찾고, 도달할 수 없는 객체를 회수한다.

### 삼색 마킹 (Tricolor Marking)

GC는 객체를 세 가지 색으로 분류한다:

```
┌────────────────────────────────────────────┐
│                                            │
│   흰색 (White)                             │
│   → 아직 검사하지 않은 객체                   │
│   → GC 완료 후 흰색이면 회수 대상            │
│                                            │
│   회색 (Gray)                              │
│   → 검사 중인 객체                           │
│   → 이 객체가 참조하는 다른 객체를 아직 확인 안함│
│                                            │
│   검은색 (Black)                            │
│   → 검사 완료된 객체                         │
│   → 이 객체와 참조하는 모든 객체가 확인됨      │
│                                            │
└────────────────────────────────────────────┘
```

**마킹 과정:**

1. 처음에 모든 객체는 **흰색**이다.
2. 루트(스택, 전역 변수)에서 직접 접근 가능한 객체를 **회색**으로 표시한다.
3. 회색 객체를 꺼내서:
   - 이 객체가 참조하는 흰색 객체를 **회색**으로 변경한다.
   - 이 객체를 **검은색**으로 변경한다.
4. 회색 객체가 없을 때까지 3번을 반복한다.
5. 남은 **흰색** 객체를 모두 회수한다 (쓰레기).

```
단계 1: 루트에서 시작
  Root → [A] → [B] → [C]
              ↘ [D]
         [E] (루트에서 접근 불가)

단계 2: 마킹 후
  Root → [A검] → [B검] → [C검]
                 ↘ [D검]
         [E흰] ← 회수 대상!
```

### Go GC의 특징

1. **동시 실행 (Concurrent)**: GC가 프로그램과 동시에 실행된다. 대부분의 GC 작업이 프로그램을 멈추지 않고 진행된다.

2. **짧은 STW (Stop-The-World)**: 전체 프로그램을 멈추는 시간이 매우 짧다 (보통 수백 마이크로초 이하). Go 1.8 이후로 STW 시간은 100 마이크로초 미만을 목표로 한다.

3. **쓰기 장벽 (Write Barrier)**: 동시 마킹 중 프로그램이 참조를 변경해도 정확성을 보장한다. 프로그램이 포인터를 변경할 때 GC에 알리는 메커니즘으로, "검은색 객체가 흰색 객체를 직접 참조하지 않는다"는 불변 조건을 유지한다.

4. **세대 구분 없음**: Java의 G1 GC와 달리 Go GC는 세대별(generational) 구분을 하지 않는다. 모든 객체를 동일하게 처리한다.

### GC 튜닝

```go
import "runtime"

// GOGC 환경변수 또는 runtime에서 설정
// GOGC=100 (기본값): 힙이 이전 GC 후의 100% 만큼 커지면 GC 실행
// GOGC=200: 더 적게 GC 실행 (메모리 더 사용, CPU 덜 사용)
// GOGC=50:  더 자주 GC 실행 (메모리 덜 사용, CPU 더 사용)
// GOGC=off: GC 비활성화

// Go 1.19에 추가된 GOMEMLIMIT: 메모리 사용 상한을 설정
// GOMEMLIMIT=1GiB: 힙 메모리를 1GB로 제한

// 수동 GC 트리거
runtime.GC()

// 메모리 통계 확인
var stats runtime.MemStats
runtime.ReadMemStats(&stats)
fmt.Printf("힙 할당: %d MB\n", stats.HeapAlloc/1024/1024)
fmt.Printf("GC 횟수: %d\n", stats.NumGC)
fmt.Printf("마지막 GC 소요 시간: %v\n", time.Duration(stats.PauseNs[(stats.NumGC+255)%256]))
```

**GOGC와 GOMEMLIMIT의 관계:**
- `GOGC`는 GC 빈도를 조절한다. 값이 클수록 GC가 덜 자주 실행된다.
- `GOMEMLIMIT`은 메모리 상한을 설정한다. 이 한도에 가까워지면 `GOGC` 값과 관계없이 GC가 더 자주 실행된다.
- 컨테이너 환경에서는 `GOMEMLIMIT`을 설정하여 OOM(Out of Memory) Kill을 방지하는 것이 좋다.

### GC 친화적인 코드 작성법

1. **불필요한 힙 할당 줄이기**

```go
// 나쁜 예: 루프마다 새 슬라이스 할당
for i := 0; i < 1000; i++ {
    data := make([]byte, 1024)
    process(data)
}

// 좋은 예: 슬라이스를 재사용
data := make([]byte, 1024)
for i := 0; i < 1000; i++ {
    process(data)
}
```

2. **sync.Pool로 객체 재사용**

```go
var bufPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process() {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufPool.Put(buf)
    }()
    // buf 사용...
}
```

`sync.Pool`은 GC 사이클 사이에 객체를 재사용하는 메커니즘이다. GC가 실행되면 Pool의 객체가 회수될 수 있으므로, 캐시 용도로만 사용해야 하며 영구 저장소로 사용해서는 안 된다.

3. **포인터 남용 피하기**: 포인터가 많을수록 GC가 추적해야 할 참조가 많아진다.

```go
// GC 부담이 큰 구조
type Node struct {
    Value int
    Next  *Node  // 포인터 → GC가 추적
}

// GC 부담이 적은 구조 (슬라이스 기반)
type FlatList struct {
    Values []int  // 내부 배열은 포인터 하나
}
```

4. **이스케이프 분석 이해하기**: Go 컴파일러는 이스케이프 분석(escape analysis)을 통해 변수를 스택에 할당할지 힙에 할당할지 결정한다. 스택에 할당된 변수는 GC 대상이 아니므로 성능이 좋다.

```bash
go build -gcflags="-m" .  # 이스케이프 분석 결과 확인
```

변수가 힙으로 이스케이프하는 주요 원인:
- 함수가 로컬 변수의 포인터를 반환하는 경우
- 인터페이스 타입 변수에 값을 할당하는 경우
- 클로저가 로컬 변수를 캡처하는 경우
- 슬라이스의 용량이 런타임에 결정되는 경우

> 자세한 예제는 `gc_demo.go`를 참고한다.

---

## 연습문제

### 개념 문제

**1.** Go가 상속 대신 조합(composition)을 선택한 설계 이유를 설명하라. 상속이 가져올 수 있는 문제점(다이아몬드 문제, 취약한 기반 클래스 문제 등)과 조합이 이를 어떻게 해결하는지 서술하라.

**2.** Go의 암시적 인터페이스 구현이 Java나 C#의 명시적 구현(`implements` 키워드)과 비교하여 어떤 장단점이 있는지 설명하라.

**3.** 다음 코드에서 컴파일 에러가 발생하는 위치와 이유를 설명하라:

```go
type Printer interface {
    Print()
}

type Document struct{ content string }
func (d *Document) Print() { fmt.Println(d.content) }

func main() {
    var p Printer
    p = Document{content: "hello"}
    p.Print()
}
```

**4.** Functional Options 패턴과 Config 구조체 패턴(`NewServer(config Config)`)의 차이점을 비교하라. 각각 어떤 상황에서 더 적합한가?

**5.** "모든 것은 값으로 전달된다"는 Go의 원칙이 맵(map)에는 어떻게 적용되는가? 맵을 함수에 전달하면 원본이 수정되는 이유를 설명하라.

**6.** `sync.Mutex`를 포함하는 구조체를 값으로 복사하면 어떤 문제가 발생하는지 설명하라. `go vet`가 이를 어떻게 감지하는가?

**7.** Go GC의 삼색 마킹에서 "쓰기 장벽(write barrier)"이 없으면 어떤 문제가 발생하는지 구체적인 시나리오를 들어 설명하라.

**8.** 의존성 역전 원칙(DIP)을 적용하면 코드의 테스트 가능성(testability)이 어떻게 향상되는지, 구체적인 예를 들어 설명하라.

### 코딩 문제

**9.** 다음 인터페이스를 만족하는 `Stack` 구조체를 값 리시버와 포인터 리시버 중 적절한 것을 선택하여 구현하라. 선택 이유도 설명하라:

```go
type Stacker interface {
    Push(v int)
    Pop() (int, bool)
    Peek() (int, bool)
    Size() int
}
```

**10.** `GOGC` 값을 50, 100, 200으로 변경하면서 벤치마크를 실행하여 GC 빈도와 프로그램 성능 차이를 측정하는 테스트 코드를 작성하라. `runtime.ReadMemStats`를 사용하여 GC 횟수와 총 STW 시간을 출력한다.

---

## 구현 과제

**1. 동물원 시뮬레이터:** 다음 조건을 만족하는 프로그램을 작성하라:
- `Animal` 인터페이스(Speak, Eat, Sleep 메서드)를 정의한다.
- 최소 5종의 동물 구조체를 구현하되, 임베딩을 활용하여 공통 로직을 재사용한다.
- `Zoo` 구조체에 동물을 추가/제거하는 기능을 구현하고, `NewZoo` 생성자와 Functional Options 패턴을 적용한다.

**2. 의존성 역전 실습 — 알림 시스템:** 알림 전송 인터페이스를 정의하고, Email, SMS, Slack 세 가지 구현체를 만들어라. `NotificationService`는 인터페이스에만 의존하며, 테스트에서는 목(mock) 구현체를 사용하여 단위 테스트를 작성한다.

**3. 메모리 할당 분석기:** 다양한 데이터 구조(슬라이스, 맵, 연결 리스트, 이진 트리)를 생성하고, `runtime.ReadMemStats`를 사용하여 각 구조의 메모리 사용량과 GC 영향을 비교하는 프로그램을 작성하라. 결과를 표 형태로 출력한다.

**4. 값 vs 포인터 벤치마크:** 다양한 크기(16바이트, 64바이트, 256바이트, 1KB, 1MB)의 구조체를 값으로 전달하는 것과 포인터로 전달하는 것의 성능을 비교하는 벤치마크를 작성하라. `testing.B`를 사용하여 ns/op를 측정하고, 어느 시점에서 포인터 전달이 유리해지는지 분석한다.

**5. 인터페이스 기반 플러그인 시스템:** 인터페이스를 사용하여 간단한 텍스트 처리 플러그인 시스템을 구현하라:
- `Transformer` 인터페이스(`Transform(string) string` 메서드)를 정의한다.
- UpperCase, LowerCase, Reverse, CaesarCipher 등의 플러그인을 구현한다.
- 파이프라인으로 여러 Transformer를 연결하여 순차적으로 적용하는 기능을 구현한다.

---

## 프로젝트 과제

**1. 마이크로서비스 프레임워크:** 이번 노트에서 배운 모든 개념을 종합하여 간단한 마이크로서비스 프레임워크를 설계하고 구현하라:
- 인터페이스 기반 미들웨어 체인 (로깅, 인증, 속도 제한)
- Functional Options 패턴으로 서버 설정
- 의존성 주입 컨테이너 (인터페이스 → 구현체 매핑)
- 포인터와 값 타입을 적절히 사용한 요청/응답 구조체
- `sync.Pool`을 활용한 객체 재사용으로 GC 부담 최소화
- GC 메트릭 모니터링 엔드포인트 (`runtime.ReadMemStats` 활용)

**2. OOP 패턴 카탈로그:** Go 스타일로 다음 디자인 패턴을 구현하는 예제 모음을 만들어라:
- Strategy 패턴 (인터페이스 활용)
- Observer 패턴 (채널과 인터페이스 활용)
- Builder 패턴 (Functional Options 패턴과 비교)
- Decorator 패턴 (임베딩과 인터페이스 활용)
- 각 패턴에 대해 "전통적 OOP 구현"과 "Go 스타일 구현"을 비교하는 주석을 작성하고, Go에서 해당 패턴이 왜 다르게 표현되는지 설명한다.
