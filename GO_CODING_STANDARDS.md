# Go 코드 작성 표준 가이드

> Go 언어의 코딩 컨벤션과 모범 사례를 정리한 종합 가이드이다.
> Python의 PEP 8에 해당하는 Go 버전이라 할 수 있다.
> 모든 Go 개발자가 일관된 코드를 작성할 수 있도록 공식 스타일 가이드와
> 커뮤니티 모범 사례를 한데 모았다.

---

## 목차

1. [공식 스타일 가이드 소개](#1-공식-스타일-가이드-소개)
2. [코드 포매팅](#2-코드-포매팅)
3. [네이밍 컨벤션](#3-네이밍-컨벤션)
4. [패키지 설계](#4-패키지-설계)
5. [Import 스타일](#5-import-스타일)
6. [에러 처리](#6-에러-처리)
7. [함수 설계](#7-함수-설계)
8. [구조체와 메서드](#8-구조체와-메서드)
9. [인터페이스 설계](#9-인터페이스-설계)
10. [동시성 패턴](#10-동시성-concurrency-패턴)
11. [테스트 컨벤션](#11-테스트-컨벤션)
12. [주석과 문서화](#12-주석과-문서화)
13. [프로젝트 구조](#13-프로젝트-구조-표준-레이아웃)
14. [Go Proverbs](#14-go-proverbs-go-격언)
15. [자주 하는 실수와 안티패턴](#15-자주-하는-실수와-안티패턴)
16. [도구 생태계](#16-도구-생태계)

---

## 1. 공식 스타일 가이드 소개

Go 언어는 다른 언어와 달리 공식적인 스타일 가이드와 도구가 잘 갖추어져 있다.
코드 포매팅부터 네이밍, 설계 원칙까지 커뮤니티 전체가 합의한 표준이 존재한다.

### 1.1 Effective Go

Go 팀이 직접 작성한 공식 가이드 문서이다. Go 코드를 관용적(idiomatic)으로
작성하는 방법을 다루며, 언어의 설계 철학과 모범 사례를 설명한다.

- **URL**: https://go.dev/doc/effective_go
- **역할**: Go 언어의 관용적 사용법을 배우는 첫 번째 필독 문서이다.

### 1.2 Go Code Review Comments

Go 프로젝트의 코드 리뷰에서 자주 지적되는 사항을 정리한 문서이다.
Effective Go를 보완하는 실전 가이드이다.

- **URL**: https://go.dev/wiki/CodeReviewComments
- **역할**: 코드 리뷰 시 체크리스트로 활용한다.

### 1.3 Go Style Guide (Google)

Google 내부에서 사용하던 Go 스타일 가이드를 공개한 것이다.
스타일 결정(Style Decisions), 모범 사례(Best Practices), 스타일 가이드(Style Guide)
세 부분으로 구성된다.

- **URL**: https://google.github.io/styleguide/go/
- **역할**: 대규모 프로젝트에서의 일관성 유지 기준을 제공한다.

### 1.4 Uber Go Style Guide

Uber 엔지니어링 팀이 공개한 Go 스타일 가이드이다.
실무에서 자주 마주치는 패턴과 안티패턴을 구체적인 예시와 함께 다룬다.

- **URL**: https://github.com/uber-go/guide/blob/master/style.md
- **역할**: 실무 중심의 구체적인 코딩 가이드라인을 제공한다.

### 1.5 가이드 우선순위

스타일이 충돌할 경우 다음 우선순위를 따른다:

1. `gofmt` / `goimports` (자동 포매팅은 무조건 따른다)
2. Effective Go (공식 문서)
3. Go Code Review Comments (공식 보충)
4. Google Go Style Guide / Uber Go Style Guide (팀 합의에 따라 선택한다)

---

## 2. 코드 포매팅

### 2.1 gofmt 필수 사용

Go는 포매팅 논쟁이 존재하지 않는 언어이다. 모든 Go 코드는 반드시 `gofmt`를
통해 포매팅해야 한다. 이것은 선택이 아닌 필수이다.

```bash
# 파일 포매팅
gofmt -w main.go

# 디렉토리 전체 포매팅
gofmt -w .

# 차이점만 확인
gofmt -d main.go
```

> Rob Pike의 격언: "Gofmt's style is no one's favorite, yet gofmt is everyone's favorite."
> (gofmt의 스타일은 누구의 취향도 아니지만, gofmt 자체는 모두의 취향이다.)

### 2.2 goimports

`goimports`는 `gofmt`의 기능에 더해 import 구문을 자동으로 정리한다.
사용하지 않는 import를 제거하고, 필요한 import를 자동으로 추가한다.

```bash
# 설치
go install golang.org/x/tools/cmd/goimports@latest

# 사용
goimports -w main.go
```

에디터에 저장 시 자동 실행되도록 설정하는 것을 강력히 권장한다.

### 2.3 들여쓰기

Go는 들여쓰기에 **탭(tab)**을 사용한다. 스페이스가 아니다.
이것은 `gofmt`가 강제하는 규칙이므로 논쟁의 여지가 없다.

```go
// ✅ 올바른 예시 (탭 사용)
func main() {
	fmt.Println("Hello")  // 탭으로 들여쓰기
	if true {
		fmt.Println("World")  // 탭으로 들여쓰기
	}
}
```

### 2.4 줄 길이

Go는 공식적으로 엄격한 줄 길이 제한을 두지 않는다.
그러나 가독성을 위해 **80~100자**를 권장한다.
Google 스타일 가이드는 99자를 기준으로 제시한다.

```go
// ✅ 올바른 예시: 긴 함수 시그니처는 줄바꿈한다
func processUserRegistration(
	ctx context.Context,
	userID string,
	email string,
	options *RegistrationOptions,
) error {
	// ...
}

// ❌ 잘못된 예시: 한 줄에 모두 넣어 가독성이 떨어진다
func processUserRegistration(ctx context.Context, userID string, email string, options *RegistrationOptions) error {
	// ...
}
```

### 2.5 중괄호 위치

Go에서 여는 중괄호(`{`)는 **반드시 같은 줄**에 위치해야 한다.
이것은 스타일 선호의 문제가 아니라, Go의 세미콜론 자동 삽입 규칙 때문에
문법적으로 필수인 사항이다.

Go 컴파일러는 줄 끝의 특정 토큰 뒤에 자동으로 세미콜론을 삽입한다.
따라서 중괄호를 다음 줄에 놓으면 컴파일 에러가 발생한다.

```go
// ✅ 올바른 예시
func main() {
	if x > 0 {
		// ...
	}
}

// ❌ 잘못된 예시 (컴파일 에러 발생)
func main()
{
	if x > 0
	{
		// ...
	}
}
```

### 2.6 빈 줄 사용

논리적 단위를 구분하기 위해 빈 줄을 사용한다.
함수 내에서 관련된 코드끼리 그룹을 만들되, 과도하게 사용하지 않는다.

```go
// ✅ 올바른 예시
func processOrder(ctx context.Context, orderID string) error {
	// 주문 조회
	order, err := fetchOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("fetch order: %w", err)
	}

	// 재고 확인
	if err := checkInventory(ctx, order.Items); err != nil {
		return fmt.Errorf("check inventory: %w", err)
	}

	// 결제 처리
	if err := processPayment(ctx, order); err != nil {
		return fmt.Errorf("process payment: %w", err)
	}

	return nil
}
```

---

## 3. 네이밍 컨벤션

Go의 네이밍은 매우 중요하다. 이름이 곧 문서이며, 좋은 이름은 좋은 설계를 반영한다.

### 3.1 MixedCaps 규칙

Go는 **MixedCaps**(대문자 시작) 또는 **mixedCaps**(소문자 시작)를 사용한다.
언더스코어(`_`)를 포함한 이름은 사용하지 않는다.
유일한 예외는 테스트 함수명(`Test_xxx`)이다.

```go
// ✅ 올바른 예시
var maxRetryCount int
type UserProfile struct{}
func calculateTotalPrice() float64 {}

// ❌ 잘못된 예시
var max_retry_count int
type user_profile struct{}
func calculate_total_price() float64 {}
```

### 3.2 패키지 이름

패키지 이름은 **소문자 단일 단어**를 사용한다.
언더스코어나 대문자를 포함하지 않는다.

```go
// ✅ 올바른 예시
package http
package json
package user
package auth

// ❌ 잘못된 예시
package httpUtil      // 대문자 포함
package http_util     // 언더스코어 포함
package HttpUtil      // 대문자 시작
```

다음과 같은 범용적이고 의미 없는 패키지 이름은 금지한다:

```go
// ❌ 절대 사용하지 않는다
package util
package utils
package common
package base
package helpers
package misc
```

패키지 이름은 내용을 설명해야 한다. `util.ConvertToJSON()`보다
`json.Convert()`가 훨씬 명확하다.

패키지 이름과 내부 타입/함수 이름이 반복되지 않도록 한다:

```go
// ✅ 올바른 예시: http.Client (패키지명 + 타입명이 자연스럽다)
client := http.Client{}

// ❌ 잘못된 예시: http.HTTPClient (http가 중복된다)
client := http.HTTPClient{}
```

### 3.3 인터페이스 이름

메서드가 하나인 인터페이스는 **-er 접미사**를 붙인다:

```go
// ✅ 올바른 예시
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Stringer interface {
	String() string
}

type Closer interface {
	Close() error
}

type Formatter interface {
	Format(f fmt.State, verb rune)
}
```

메서드가 여러 개인 인터페이스는 의미를 정확히 전달하는 이름을 사용한다:

```go
// ✅ 올바른 예시
type ReadWriter interface {
	Reader
	Writer
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
```

### 3.4 Getter와 Setter

Go에서는 Getter에 `Get` 접두사를 붙이지 않는다.
Setter에는 `Set` 접두사를 붙인다.

```go
// ✅ 올바른 예시
type User struct {
	name string
}

func (u *User) Name() string {       // Getter: Get 접두사 없음
	return u.name
}

func (u *User) SetName(name string) { // Setter: Set 접두사 사용
	u.name = name
}

// ❌ 잘못된 예시
func (u *User) GetName() string {     // Go에서는 Get을 붙이지 않는다
	return u.name
}
```

### 3.5 약어 (Acronyms)

약어는 **전부 대문자** 또는 **전부 소문자**로 작성한다.
혼합하지 않는다.

```go
// ✅ 올바른 예시
var userID string       // ID 전부 대문자
var httpClient *Client  // http 전부 소문자
type URLParser struct{} // URL 전부 대문자
var xmlHTTPRequest      // XML, HTTP 각각 전부 대문자

func ServeHTTP(w http.ResponseWriter, r *http.Request) {}
func parseURL(rawURL string) (*URL, error) {}

// ❌ 잘못된 예시
var orderId string      // Id 혼합 (첫 글자만 대문자)
type UrlParser struct{} // Url 혼합
var userId string       // Id 혼합
func ServeHttp() {}     // Http 혼합
```

주요 약어 표기 목록:

| 약어 | 올바른 표기 | 잘못된 표기 |
|------|-------------|-------------|
| API  | `API` / `api` | `Api` |
| HTTP | `HTTP` / `http` | `Http` |
| ID   | `ID` / `id` | `Id` |
| JSON | `JSON` / `json` | `Json` |
| URL  | `URL` / `url` | `Url` |
| SQL  | `SQL` / `sql` | `Sql` |
| XML  | `XML` / `xml` | `Xml` |
| TCP  | `TCP` / `tcp` | `Tcp` |
| TLS  | `TLS` / `tls` | `Tls` |
| SSH  | `SSH` / `ssh` | `Ssh` |
| IP   | `IP` / `ip` | `Ip` |

### 3.6 지역 변수

Go에서는 **스코프가 좁은 변수일수록 짧은 이름**을 사용한다.
이것은 Go의 핵심 철학 중 하나이다.

```go
// ✅ 올바른 예시: 짧은 스코프에서 짧은 이름
for i, v := range items {
	fmt.Println(i, v)
}

// 관용적인 짧은 변수명
var (
	i   int             // 인덱스
	n   int             // 개수, 길이
	err error           // 에러
	ctx context.Context // 컨텍스트
	buf bytes.Buffer    // 버퍼
	mu  sync.Mutex      // 뮤텍스
	wg  sync.WaitGroup  // WaitGroup
	w   io.Writer       // Writer
	r   io.Reader       // Reader
	b   []byte          // 바이트 슬라이스
	s   string          // 문자열
	ok  bool            // 맵 존재 여부, 타입 단언
	c   *http.Client    // HTTP 클라이언트
)

// ❌ 잘못된 예시: 짧은 스코프에 불필요하게 긴 이름
for index, value := range items {
	fmt.Println(index, value)
}
```

단, **스코프가 넓거나 의미가 불분명한 경우**에는 설명적인 이름을 사용한다:

```go
// ✅ 넓은 스코프에서는 설명적 이름 사용
type Server struct {
	maxRetryCount    int           // 짧은 이름이면 의미를 알 수 없다
	requestTimeout   time.Duration
	shutdownTimeout  time.Duration
}
```

### 3.7 리시버 이름

메서드 리시버는 **1~2글자**의 짧은 이름을 사용한다.
타입 이름의 첫 글자(소문자)를 사용하는 것이 관용적이다.
`this`나 `self`는 절대 사용하지 않는다.

```go
// ✅ 올바른 예시
type Server struct{}

func (s *Server) Start() error { ... }
func (s *Server) Stop() error { ... }
func (s *Server) handleRequest(r *Request) { ... }

type DatabaseClient struct{}

func (dc *DatabaseClient) Query(q string) (*Result, error) { ... }

// ❌ 잘못된 예시
func (this *Server) Start() error { ... }   // this 금지
func (self *Server) Stop() error { ... }    // self 금지
func (server *Server) Start() error { ... } // 불필요하게 긴 이름
```

**중요**: 하나의 타입에서 리시버 이름은 반드시 일관되어야 한다.
한 메서드에서 `s`를 사용했으면 모든 메서드에서 `s`를 사용한다.

### 3.8 파일 이름

Go 파일 이름은 **snake_case**를 사용한다:

```
// ✅ 올바른 예시
user_repository.go
http_handler.go
string_utils.go
message_queue.go

// ❌ 잘못된 예시
userRepository.go    // camelCase 금지
UserRepository.go    // PascalCase 금지
user-repository.go   // kebab-case 금지
```

### 3.9 테스트 파일

테스트 파일은 반드시 `_test.go` 접미사를 붙인다:

```
user_repository.go       → user_repository_test.go
http_handler.go          → http_handler_test.go
```

### 3.10 Export 규칙

Go에서 이름의 **첫 글자가 대문자**이면 패키지 외부에서 접근 가능(exported)하고,
**첫 글자가 소문자**이면 패키지 내부에서만 접근 가능(unexported)하다.

```go
// ✅ 패키지 외부에서 사용할 것: 대문자 시작
type UserService struct{}        // exported
func NewUserService() *UserService {} // exported
const MaxRetries = 3             // exported

// ✅ 패키지 내부에서만 사용할 것: 소문자 시작
type userCache struct{}          // unexported
func validateEmail(e string) bool {} // unexported
const defaultTimeout = 30        // unexported
```

---

## 4. 패키지 설계

### 4.1 하나의 패키지 = 하나의 책임

각 패키지는 하나의 명확한 책임을 가져야 한다.
패키지 이름만 보고도 무엇을 하는지 알 수 있어야 한다.

```
// ✅ 올바른 예시: 각 패키지가 하나의 명확한 책임
auth/       → 인증 관련 기능
cache/      → 캐싱 관련 기능
email/      → 이메일 발송 기능
storage/    → 스토리지 관련 기능

// ❌ 잘못된 예시: 책임이 불명확하거나 너무 넓음
util/       → 무엇이든 들어갈 수 있다
common/     → 공통이라는 건 아무것도 아니다
helpers/    → 도우미라는 건 의미가 없다
```

### 4.2 패키지 이름으로 의미 전달

호출하는 쪽에서 패키지 이름과 함께 읽힐 때 자연스러워야 한다:

```go
// ✅ 올바른 예시: 패키지명 + 타입/함수명이 자연스럽다
client := http.Client{}
reader := bufio.NewReader(file)
color  := color.RGBA{R: 255}
ctx    := context.Background()
server := grpc.NewServer()

// ❌ 잘못된 예시: 패키지명이 타입/함수명에서 반복된다
client := http.HTTPClient{}        // http 반복
reader := bufio.NewBufioReader()   // bufio 반복
server := grpc.NewGRPCServer()     // grpc 반복
```

### 4.3 internal 패키지 활용

외부에 노출하지 않을 코드는 `internal` 패키지에 넣는다.
Go 컴파일러가 `internal` 패키지의 접근을 강제로 제한한다.

```
myproject/
├── internal/
│   ├── database/     # 외부 프로젝트에서 import 불가
│   ├── middleware/    # 외부 프로젝트에서 import 불가
│   └── validator/    # 외부 프로젝트에서 import 불가
├── pkg/
│   └── api/          # 외부 프로젝트에서 import 가능
└── cmd/
    └── server/
        └── main.go
```

```go
// ✅ 같은 프로젝트 내부에서는 접근 가능
import "myproject/internal/database"

// ❌ 외부 프로젝트에서는 컴파일 에러 발생
import "github.com/someone/myproject/internal/database"
// Error: use of internal package not allowed
```

### 4.4 순환 의존성 금지

Go는 순환 import를 허용하지 않는다. 컴파일 에러가 발생한다.
이것은 설계 결함의 신호이다.

```go
// ❌ 순환 의존성 (컴파일 에러)
// package a imports package b
// package b imports package a

// ✅ 해결 방법 1: 인터페이스를 사용한 의존성 역전
// package a: 인터페이스 정의
type UserStore interface {
	FindByID(id string) (*User, error)
}

// package b: 인터페이스 구현
type PostgresUserStore struct{}
func (s *PostgresUserStore) FindByID(id string) (*User, error) { ... }

// ✅ 해결 방법 2: 공통 타입을 별도 패키지로 분리
// package model: 공통 타입만 정의
type User struct {
	ID   string
	Name string
}
```

### 4.5 main 패키지는 최소한으로

`main` 패키지는 프로그램의 진입점이다.
최소한의 설정과 연결(wiring)만 수행하고, 실제 로직은 다른 패키지에 위임한다.

```go
// ✅ 올바른 예시: main은 조립만 담당한다
package main

import (
	"log"
	"os"

	"myproject/internal/config"
	"myproject/internal/server"
)

func main() {
	cfg, err := config.Load(os.Args[1:])
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("server create: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("server run: %v", err)
	}
}

// ❌ 잘못된 예시: main에 비즈니스 로직이 들어있다
package main

func main() {
	db, _ := sql.Open("postgres", "...")
	rows, _ := db.Query("SELECT * FROM users")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		// ... 100줄의 비즈니스 로직 ...
	}
}
```

---

## 5. Import 스타일

### 5.1 그룹 나누기

import 구문은 **세 그룹**으로 나누고, 빈 줄로 구분한다:

1. 표준 라이브러리
2. 외부(서드파티) 패키지
3. 내부(프로젝트) 패키지

```go
// ✅ 올바른 예시
import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"myproject/internal/auth"
	"myproject/internal/database"
)

// ❌ 잘못된 예시: 그룹 구분 없이 섞여 있다
import (
	"myproject/internal/auth"
	"fmt"
	"github.com/gorilla/mux"
	"context"
	"myproject/internal/database"
	"go.uber.org/zap"
	"net/http"
)
```

`goimports`를 사용하면 이 그룹 구분을 자동으로 처리한다.

### 5.2 dot import 금지

dot import(`. "패키지"`)는 패키지의 모든 식별자를 현재 스코프에 가져온다.
코드의 출처를 불명확하게 만들므로 사용하지 않는다.

```go
// ❌ 잘못된 예시
import . "fmt"

func main() {
	Println("어디서 온 함수인지 알 수 없다")
}

// ✅ 올바른 예시
import "fmt"

func main() {
	fmt.Println("fmt 패키지의 함수임이 명확하다")
}
```

**유일한 예외**: 테스트에서 순환 의존성을 피하기 위해 사용할 수 있다.
그러나 이 경우에도 가능하면 피하는 것이 좋다.

### 5.3 별칭(alias) 최소화

import 별칭은 **이름이 충돌할 때**만 사용한다:

```go
// ✅ 올바른 예시: 이름 충돌이 있을 때 별칭 사용
import (
	"crypto/rand"
	mrand "math/rand"
)

// ❌ 잘못된 예시: 불필요한 별칭
import (
	f "fmt"            // 이유 없는 별칭
	h "net/http"       // 이유 없는 별칭
)
```

### 5.4 blank import

blank import(`import _ "패키지"`)는 패키지의 `init()` 함수를 실행시키기 위해
사용한다. 주로 드라이버 등록에 활용한다.

```go
// ✅ 올바른 예시: 데이터베이스 드라이버 등록
import (
	"database/sql"

	_ "github.com/lib/pq"            // PostgreSQL 드라이버 등록
	_ "github.com/go-sql-driver/mysql" // MySQL 드라이버 등록
)

// ✅ 올바른 예시: 이미지 디코더 등록
import (
	"image"

	_ "image/png"  // PNG 디코더 등록
	_ "image/jpeg" // JPEG 디코더 등록
)
```

blank import를 사용할 때는 반드시 **주석으로 용도를 설명**한다.

---

## 6. 에러 처리

Go의 에러 처리는 언어의 핵심 철학이다.
"Errors are values" — 에러는 특별한 것이 아니라 값이다.

### 6.1 기본 패턴: if err != nil

Go에서 에러를 처리하는 기본 패턴이다:

```go
// ✅ 올바른 예시
result, err := doSomething()
if err != nil {
	return fmt.Errorf("do something: %w", err)
}
// result 사용
```

에러 검사는 즉시 수행하고, 성공 경로를 바깥쪽에 유지한다:

```go
// ✅ 올바른 예시: 에러 시 early return, 성공 경로가 바깥쪽
func processFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	return process(data)
}

// ❌ 잘못된 예시: else에 성공 로직이 중첩된다
func processFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	} else {
		data, err := io.ReadAll(f)
		if err != nil {
			return err
		} else {
			return process(data)
		}
	}
}
```

### 6.2 에러 래핑 (Error Wrapping)

`fmt.Errorf`와 `%w` 동사를 사용하여 에러에 컨텍스트를 추가한다.
이렇게 하면 에러 체인이 형성되어 `errors.Is`와 `errors.As`로 원인을 추적할 수 있다.

```go
// ✅ 올바른 예시: 컨텍스트와 함께 래핑
func getUserByID(ctx context.Context, id string) (*User, error) {
	row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user %s: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("get user %s: %w", id, err)
	}

	return &user, nil
}

// ❌ 잘못된 예시: 컨텍스트 없이 그대로 반환
func getUserByID(ctx context.Context, id string) (*User, error) {
	row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, err  // 어디서 발생한 에러인지 알 수 없다
	}

	return &user, nil
}
```

### 6.3 센티널 에러 (Sentinel Errors)

패키지 수준에서 공통적으로 사용할 에러를 정의한다.
변수명은 `Err` 접두사로 시작한다.

```go
// ✅ 올바른 예시
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal error")
)

// 사용
func GetUser(id string) (*User, error) {
	user, ok := users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

// 호출측에서 에러 종류 확인
user, err := GetUser("123")
if errors.Is(err, ErrNotFound) {
	// 404 응답
}
```

### 6.4 커스텀 에러 타입

더 많은 정보를 전달해야 할 때 커스텀 에러 타입을 사용한다:

```go
// ✅ 올바른 예시
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %s: %s", e.Field, e.Message)
}

// 사용
func validateAge(age int) error {
	if age < 0 || age > 150 {
		return &ValidationError{
			Field:   "age",
			Message: "must be between 0 and 150",
		}
	}
	return nil
}

// 호출측에서 에러 타입으로 분기
var ve *ValidationError
if errors.As(err, &ve) {
	fmt.Printf("필드 %s에서 검증 실패: %s\n", ve.Field, ve.Message)
}
```

### 6.5 errors.Is와 errors.As

Go 1.13부터 에러 체인을 탐색하는 표준 함수를 사용한다:

```go
// errors.Is: 특정 에러 값과 일치하는지 확인한다
if errors.Is(err, ErrNotFound) {
	// err 또는 그 체인 안에 ErrNotFound가 있다
}

// errors.As: 특정 에러 타입으로 변환한다
var pathErr *os.PathError
if errors.As(err, &pathErr) {
	fmt.Println("경로:", pathErr.Path)
}

// ❌ 잘못된 예시: 직접 비교 (래핑된 에러를 놓친다)
if err == ErrNotFound { ... }           // 래핑되면 실패
if _, ok := err.(*os.PathError); ok { } // 래핑되면 실패
```

### 6.6 에러 메시지 규칙

에러 메시지는 **소문자로 시작**하고 **마침표를 붙이지 않는다**.
이유는 에러가 래핑될 때 중간에 위치하게 되기 때문이다.

```go
// ✅ 올바른 예시
return fmt.Errorf("open database: %w", err)
return errors.New("invalid user ID")
return fmt.Errorf("parse config at %s: %w", path, err)

// 래핑 결과: "start server: open database: connection refused"
// → 자연스럽게 연결된다

// ❌ 잘못된 예시
return fmt.Errorf("Failed to open database: %w", err)  // 대문자 시작
return errors.New("Invalid user ID.")                   // 마침표 있음
return fmt.Errorf("Error: cannot parse config: %w", err) // "Error:" 접두사 불필요
```

### 6.7 panic 사용 제한

`panic`은 **정말로 복구 불가능한 프로그래밍 오류**에서만 사용한다.
일반적인 에러 상황에서는 절대 사용하지 않는다.

```go
// ✅ panic이 허용되는 경우
func MustCompileRegex(pattern string) *regexp.Regexp {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("invalid regex %q: %v", pattern, err))
	}
	return re
}

// 프로그램 초기화 시에만 사용 (Must 접두사로 명시)
var emailRegex = MustCompileRegex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ❌ panic을 사용하면 안 되는 경우
func GetUser(id string) *User {
	user, err := db.FindUser(id)
	if err != nil {
		panic(err) // ❌ 일반 에러에 panic을 사용하지 않는다
	}
	return user
}
```

### 6.8 에러 무시 금지

에러를 무시해야 할 경우에도 **명시적으로** 표현한다:

```go
// ✅ 올바른 예시: 에러를 명시적으로 무시하고 이유를 주석으로 남긴다
_ = writer.Close() // 이미 모든 데이터를 flush했으므로 close 에러는 무시한다

// ✅ 올바른 예시: fmt.Fprintf는 os.Stdout에 쓸 때 에러가 거의 없으므로 무시한다
fmt.Fprintf(os.Stdout, "result: %v\n", result)

// ❌ 잘못된 예시: 에러를 조용히 무시한다
result, _ := doSomething() // 왜 무시하는지 알 수 없다
doSomething()              // 반환값 자체를 받지 않는다
```

---

## 7. 함수 설계

### 7.1 함수 시그니처

인자 개수는 **3~4개 이하**를 유지한다.
그 이상이 필요하면 옵션 구조체를 사용한다.

```go
// ✅ 올바른 예시: 인자가 적다
func CreateUser(ctx context.Context, name, email string) (*User, error) {
	// ...
}

// ✅ 올바른 예시: 인자가 많으면 옵션 구조체 사용
type CreateUserOptions struct {
	Name        string
	Email       string
	Age         int
	Role        string
	Department  string
	ManagerID   string
}

func CreateUser(ctx context.Context, opts CreateUserOptions) (*User, error) {
	// ...
}

// ❌ 잘못된 예시: 인자가 너무 많다
func CreateUser(ctx context.Context, name, email string, age int, role, dept, managerID string) (*User, error) {
	// ...
}
```

### 7.2 에러는 마지막 반환값

에러를 반환하는 함수는 **항상 에러를 마지막 반환값**에 놓는다:

```go
// ✅ 올바른 예시
func FindUser(id string) (*User, error) { ... }
func ReadFile(path string) ([]byte, error) { ... }
func ParseConfig(data []byte) (*Config, bool, error) { ... }

// ❌ 잘못된 예시
func FindUser(id string) (error, *User) { ... }  // 에러가 첫 번째
```

### 7.3 Early Return 패턴

성공 경로는 가장 바깥쪽 들여쓰기에 유지한다.
에러나 예외 상황에서 먼저 반환한다.

```go
// ✅ 올바른 예시: 성공 경로가 바깥쪽에 있다
func authenticate(token string) (*User, error) {
	if token == "" {
		return nil, ErrEmptyToken
	}

	claims, err := parseToken(token)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if claims.Expired() {
		return nil, ErrTokenExpired
	}

	user, err := findUser(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

// ❌ 잘못된 예시: 중첩이 깊어진다
func authenticate(token string) (*User, error) {
	if token != "" {
		claims, err := parseToken(token)
		if err == nil {
			if !claims.Expired() {
				user, err := findUser(claims.UserID)
				if err == nil {
					return user, nil
				}
				return nil, err
			}
			return nil, ErrTokenExpired
		}
		return nil, err
	}
	return nil, ErrEmptyToken
}
```

### 7.4 함수 길이

함수는 **한 화면에 들어올 정도**(약 40~60줄)로 유지하는 것을 권장한다.
함수가 길어지면 더 작은 함수로 분리한다.

### 7.5 Named Return (명명된 반환값)

named return은 **짧은 함수**에서 godoc 문서화 목적으로만 사용한다.
긴 함수에서는 코드를 읽기 어렵게 만들므로 사용하지 않는다.
bare return(값 없는 return)은 사용하지 않는다.

```go
// ✅ 올바른 예시: 짧은 함수에서 문서화 목적
func divmod(a, b int) (quotient, remainder int) {
	quotient = a / b
	remainder = a % b
	return quotient, remainder  // 명시적으로 값을 반환한다
}

// ❌ 잘못된 예시: bare return 사용
func divmod(a, b int) (quotient, remainder int) {
	quotient = a / b
	remainder = a % b
	return  // 무엇을 반환하는지 알기 어렵다
}

// ❌ 잘못된 예시: 긴 함수에서 named return
func processOrder(ctx context.Context, id string) (order *Order, total float64, err error) {
	// ... 50줄의 코드 ...
	// 이 시점에서 order, total, err가 어디서 설정되었는지 추적하기 어렵다
	return
}
```

### 7.6 defer

리소스 정리에는 반드시 `defer`를 사용한다.
열었으면 닫고, 잠갔으면 풀어야 한다.

```go
// ✅ 올바른 예시
func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()  // 열자마자 바로 defer

	return io.ReadAll(f)
}

func updateCounter(mu *sync.Mutex, counter *int) {
	mu.Lock()
	defer mu.Unlock()  // 잠그자마자 바로 defer

	*counter++
}

// ✅ 올바른 예시: HTTP 응답 바디 닫기
resp, err := http.Get(url)
if err != nil {
	return err
}
defer resp.Body.Close()  // 반드시 닫는다
```

### 7.7 init() 함수

`init()` 함수는 최소화한다.
부작용(side effect)이 있으므로 테스트와 디버깅을 어렵게 만든다.

```go
// ✅ init()이 허용되는 경우: 간단한 변수 초기화
var defaultClient *http.Client

func init() {
	defaultClient = &http.Client{
		Timeout: 30 * time.Second,
	}
}

// ❌ init()을 사용하면 안 되는 경우
func init() {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)  // 테스트에서 문제가 된다
	}
	globalDB = db       // 전역 상태 변경
}

// ✅ 대안: 명시적 초기화 함수 사용
func NewDatabase(url string) (*sql.DB, error) {
	return sql.Open("postgres", url)
}
```

---

## 8. 구조체와 메서드

### 8.1 구조체 필드 정렬

필드는 **중요도순** 또는 **논리적 그룹**으로 정렬한다.
빈 줄로 그룹을 구분할 수 있다.

```go
// ✅ 올바른 예시: 논리적 그룹으로 정렬
type Server struct {
	// 필수 설정
	Addr    string
	Handler http.Handler

	// 타임아웃 설정
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// 내부 상태
	mu        sync.Mutex
	listeners []net.Listener
	shutdown  chan struct{}
}
```

### 8.2 생성자 패턴

구조체의 생성자는 `NewXxx()` 패턴을 사용한다.
패키지에 주요 타입이 하나뿐이면 `New()`를 사용한다.

```go
// ✅ 올바른 예시
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		shutdown:     make(chan struct{}),
	}
}

// 패키지에 주요 타입이 하나일 때
// package cache
func New(maxSize int) *Cache {
	return &Cache{
		maxSize: maxSize,
		items:   make(map[string]*entry),
	}
}
// 사용: cache.New(1000)

// ✅ 옵션 패턴 (Functional Options Pattern)
type Option func(*Server)

func WithReadTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.ReadTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.WriteTimeout = d
	}
}

func NewServer(addr string, opts ...Option) *Server {
	s := &Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,   // 기본값
		WriteTimeout: 10 * time.Second,  // 기본값
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// 사용
srv := NewServer(":8080",
	WithReadTimeout(10*time.Second),
	WithWriteTimeout(30*time.Second),
)
```

### 8.3 포인터 리시버 vs 값 리시버

다음 표를 기준으로 선택한다:

| 조건 | 리시버 타입 | 이유 |
|------|-------------|------|
| 메서드가 리시버를 수정한다 | **포인터** | 원본을 변경해야 하기 때문이다 |
| 리시버가 큰 구조체이다 | **포인터** | 복사 비용을 줄이기 위해서이다 |
| 리시버에 sync.Mutex 등이 있다 | **포인터** | 복사하면 동작이 깨진다 |
| 리시버가 map, func, chan이다 | **값** | 이미 참조 타입이기 때문이다 |
| 리시버가 작은 기본 타입이다 | **값** | 복사 비용이 없기 때문이다 |
| 리시버가 변경되지 않는 작은 구조체이다 | **값** | 안전하고 단순하기 때문이다 |
| 판단이 어려울 때 | **포인터** | 기본적으로 포인터를 선택한다 |

```go
// ✅ 포인터 리시버: 상태를 변경한다
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = true
	return nil
}

// ✅ 값 리시버: 상태를 읽기만 한다 (작은 구조체)
type Point struct {
	X, Y float64
}

func (p Point) Distance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}
```

**중요**: 한 타입의 메서드들은 리시버 타입을 **혼합하지 않는다**.
하나라도 포인터 리시버를 사용하면 모든 메서드를 포인터 리시버로 통일한다.

### 8.4 메서드 순서

파일 내에서 메서드는 다음 순서로 배치한다:

1. 생성자 (`NewXxx`)
2. 공개(exported) 메서드
3. 비공개(unexported) 메서드
4. 인터페이스 구현 메서드

```go
type UserService struct { ... }

// 1. 생성자
func NewUserService(repo UserRepository) *UserService { ... }

// 2. 공개 메서드
func (s *UserService) CreateUser(ctx context.Context, user *User) error { ... }
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) { ... }
func (s *UserService) DeleteUser(ctx context.Context, id string) error { ... }

// 3. 비공개 메서드
func (s *UserService) validateUser(user *User) error { ... }
func (s *UserService) notifyCreation(user *User) { ... }

// 4. 인터페이스 구현
func (s *UserService) String() string { ... }
```

### 8.5 제로값이 유용하도록 설계

Go의 제로값(zero value)이 바로 사용 가능하도록 구조체를 설계한다.
이것은 Go의 핵심 설계 원칙 중 하나이다.

```go
// ✅ 올바른 예시: 제로값이 바로 사용 가능하다
var buf bytes.Buffer
buf.WriteString("hello")  // 초기화 없이 바로 사용 가능

var mu sync.Mutex
mu.Lock()  // 초기화 없이 바로 사용 가능

// ✅ 올바른 예시: 제로값에 합리적인 기본 동작을 부여한다
type Logger struct {
	Output io.Writer
	Level  int
}

func (l *Logger) output() io.Writer {
	if l.Output == nil {
		return os.Stderr  // 제로값일 때 기본값 사용
	}
	return l.Output
}

// ❌ 잘못된 예시: 초기화 없이 사용하면 panic
type Cache struct {
	items map[string]interface{}
}

func (c *Cache) Set(key string, value interface{}) {
	c.items[key] = value  // nil map에 쓰기 → panic!
}

// ✅ 수정: lazy initialization으로 제로값 안전하게 처리
func (c *Cache) Set(key string, value interface{}) {
	if c.items == nil {
		c.items = make(map[string]interface{})
	}
	c.items[key] = value
}
```

### 8.6 구조체 비교

구조체의 모든 필드가 비교 가능하면 `==`로 비교할 수 있다.
그러나 `map`, `slice`, `func` 필드가 있으면 비교할 수 없다.

```go
// ✅ 비교 가능한 구조체
type Point struct {
	X, Y int
}
p1 := Point{1, 2}
p2 := Point{1, 2}
fmt.Println(p1 == p2) // true

// ❌ 비교 불가능한 구조체 (컴파일 에러)
type Config struct {
	Name   string
	Values []string  // 슬라이스는 비교할 수 없다
}
// c1 == c2 → 컴파일 에러

// ✅ 비교가 필요하면 reflect.DeepEqual 또는 커스텀 Equal 메서드 사용
func (c Config) Equal(other Config) bool {
	if c.Name != other.Name {
		return false
	}
	return slices.Equal(c.Values, other.Values)
}
```

---

## 9. 인터페이스 설계

### 9.1 작은 인터페이스 선호

Go의 표준 라이브러리를 보면 대부분의 인터페이스가 1~3개의 메서드만 가진다.
작은 인터페이스는 구현하기 쉽고, 조합하기 쉽고, 테스트하기 쉽다.

```go
// ✅ 올바른 예시: 작고 집중된 인터페이스
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}

// 작은 인터페이스를 조합하여 큰 인터페이스를 만든다
type ReadWriter interface {
	Reader
	Writer
}

type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// ❌ 잘못된 예시: 너무 큰 인터페이스
type FileSystem interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
	Remove(name string) error
	Rename(old, new string) error
	Stat(name string) (FileInfo, error)
	ReadDir(name string) ([]DirEntry, error)
	MkdirAll(path string, perm os.FileMode) error
	Chmod(name string, mode os.FileMode) error
	Chown(name string, uid, gid int) error
	// 메서드가 9개 — 이 인터페이스를 구현하려면 고통스럽다
}
```

> Go Proverb: "The bigger the interface, the weaker the abstraction."
> (인터페이스가 클수록 추상화는 약해진다.)

### 9.2 인터페이스는 사용하는 쪽에서 정의한다

인터페이스는 구현하는 쪽(producer)이 아니라 **사용하는 쪽(consumer)**에서 정의한다.
이것은 Go 인터페이스 설계의 핵심 원칙이다.

```go
// ✅ 올바른 예시: 사용하는 쪽에서 필요한 만큼만 인터페이스를 정의한다

// package order (사용하는 쪽)
type UserFinder interface {
	FindByID(ctx context.Context, id string) (*user.User, error)
}

type OrderService struct {
	users UserFinder  // 필요한 메서드만 요구한다
}

// package user (구현하는 쪽) — 인터페이스를 정의하지 않는다
type Repository struct {
	db *sql.DB
}

func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) { ... }
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) { ... }
func (r *Repository) Save(ctx context.Context, u *User) error { ... }
func (r *Repository) Delete(ctx context.Context, id string) error { ... }

// OrderService는 FindByID만 필요하므로 나머지 메서드에 의존하지 않는다

// ❌ 잘못된 예시: 구현 쪽에서 큰 인터페이스를 먼저 정의한다
// package user
type Repository interface {  // 구현 쪽에서 정의 — Go답지 않다
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string) error
}
```

### 9.3 빈 인터페이스(any) 남용 금지

`any`(또는 `interface{}`)는 타입 정보를 잃으므로 가능한 사용하지 않는다.
제네릭이 도입된 Go 1.18 이후로는 더욱 피해야 한다.

```go
// ❌ 잘못된 예시: any 남용
func Process(data any) any {
	// 타입을 알 수 없어서 타입 단언이 남발된다
	switch v := data.(type) {
	case string:
		return strings.ToUpper(v)
	case int:
		return v * 2
	default:
		return nil
	}
}

// ✅ 올바른 예시: 제네릭 사용 (Go 1.18+)
func Transform[T any](items []T, fn func(T) T) []T {
	result := make([]T, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

// ✅ 올바른 예시: 구체적 타입 사용
func ProcessString(data string) string {
	return strings.ToUpper(data)
}
```

`any`가 허용되는 경우:
- JSON 디코딩 시 스키마를 모를 때
- 로깅/디버깅 함수
- 리플렉션 기반 라이브러리 내부

### 9.4 인터페이스가 아닌 구체 타입을 반환한다

함수는 인터페이스가 아닌 **구체 타입을 반환**한다.
이것이 "Accept interfaces, return structs" 원칙이다.

```go
// ✅ 올바른 예시: 구체 타입을 반환한다
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// ❌ 잘못된 예시: 인터페이스를 반환한다
func NewUserService(repo UserRepository) UserServiceInterface {
	return &UserService{repo: repo}
}
```

인터페이스를 반환하면 호출자가 구체 타입의 추가 메서드에 접근할 수 없고,
불필요한 추상화 계층이 추가된다. 예외적으로 `error` 인터페이스를 반환하는 것은 괜찮다.

### 9.5 Accept Interfaces, Return Structs 원칙

함수의 **인자**에는 인터페이스를 사용하고, **반환값**에는 구체 타입을 사용한다.
이렇게 하면 함수의 유연성과 사용성이 모두 높아진다.

```go
// ✅ 올바른 예시
// 인자: 인터페이스 (유연하다 — 어떤 Reader든 받을 수 있다)
// 반환: 구체 타입 (명확하다 — 무엇을 받는지 정확히 안다)
func ParseConfig(r io.Reader) (*Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

// 다양한 소스에서 사용할 수 있다
cfg1, _ := ParseConfig(os.Stdin)                       // 표준 입력
cfg2, _ := ParseConfig(strings.NewReader(`{"key":"v"}`)) // 문자열
cfg3, _ := ParseConfig(resp.Body)                       // HTTP 응답
```

---

## 10. 동시성 (Concurrency) 패턴

### 10.1 핵심 원칙

> "Don't communicate by sharing memory; share memory by communicating."
> (메모리를 공유하여 통신하지 말고, 통신을 통해 메모리를 공유하라.)

이것은 Go 동시성의 근본 원칙이다.
가능하면 mutex 대신 채널을 사용하여 고루틴 간에 데이터를 전달한다.

### 10.2 고루틴 누수 방지

시작한 고루틴은 반드시 종료할 수 있어야 한다.
`context`나 done 채널을 사용하여 취소 신호를 전달한다.

```go
// ✅ 올바른 예시: context로 고루틴 취소 가능
func watch(ctx context.Context, ch <-chan Event) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()  // 외부에서 취소하면 종료한다
		case event, ok := <-ch:
			if !ok {
				return nil  // 채널이 닫히면 종료한다
			}
			process(event)
		}
	}
}

// 사용
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

go watch(ctx, eventCh)

// ❌ 잘못된 예시: 고루틴이 영원히 블로킹될 수 있다
func watch(ch <-chan Event) {
	for event := range ch {
		process(event)
	}
	// ch가 닫히지 않으면 영원히 대기한다 → 고루틴 누수
}
```

### 10.3 sync.Mutex 배치

뮤텍스는 **보호할 필드 바로 위에** 선언한다.
주석으로 어떤 필드를 보호하는지 명시한다.

```go
// ✅ 올바른 예시
type SafeCounter struct {
	mu     sync.Mutex // counts를 보호한다
	counts map[string]int
}

func (c *SafeCounter) Increment(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[key]++
}

func (c *SafeCounter) Get(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[key]
}

// 읽기가 많은 경우 RWMutex 사용
type Cache struct {
	mu    sync.RWMutex // items를 보호한다
	items map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()         // 읽기 잠금
	defer c.mu.RUnlock()
	v, ok := c.items[key]
	return v, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()          // 쓰기 잠금
	defer c.mu.Unlock()
	c.items[key] = value
}
```

### 10.4 sync.WaitGroup 패턴

여러 고루틴의 완료를 기다릴 때 `sync.WaitGroup`을 사용한다:

```go
// ✅ 올바른 예시
func processItems(items []Item) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(items))

	for _, item := range items {
		wg.Add(1)
		go func(it Item) {
			defer wg.Done()
			if err := process(it); err != nil {
				errs <- err
			}
		}(item)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err  // 첫 번째 에러 반환
		}
	}
	return nil
}
```

### 10.5 Channel vs Mutex 선택 기준

| 상황 | 선택 | 이유 |
|------|------|------|
| 소유권 전달 (한쪽이 다른 쪽에 데이터를 넘긴다) | **Channel** | 데이터의 소유권이 이동한다 |
| 파이프라인 (단계별 처리) | **Channel** | 생산자-소비자 패턴이다 |
| 타임아웃/취소가 필요하다 | **Channel + select** | select로 다중 이벤트를 처리한다 |
| 공유 상태를 보호한다 | **Mutex** | 여러 고루틴이 같은 데이터에 접근한다 |
| 카운터/캐시 | **Mutex** | 단순한 읽기/쓰기 보호이다 |
| 성능이 중요하다 | **Mutex** | 채널보다 오버헤드가 적다 |

```go
// ✅ Channel: 파이프라인 패턴
func pipeline(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n * n
		}
	}()
	return out
}

// ✅ Mutex: 공유 캐시 보호
type Cache struct {
	mu    sync.RWMutex
	items map[string]string
}
```

### 10.6 errgroup 활용

`golang.org/x/sync/errgroup`은 고루틴 그룹의 에러를 처리하는 표준 패턴이다:

```go
// ✅ 올바른 예시: errgroup으로 여러 작업을 병렬 실행
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	results := make([]string, len(urls))

	for i, url := range urls {
		i, url := i, url  // 루프 변수 캡처 (Go 1.22 이전)
		g.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("fetch %s: %w", url, err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("read %s: %w", url, err)
			}

			results[i] = string(body)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err  // 첫 번째 에러를 반환한다
	}

	return results, nil
}
```

`errgroup.WithContext`를 사용하면 한 고루틴이 실패했을 때
context가 취소되어 나머지 고루틴도 빠르게 종료할 수 있다.

---

## 11. 테스트 컨벤션

### 11.1 테이블 주도 테스트 (Table-Driven Tests)

Go에서 가장 널리 사용되는 테스트 패턴이다.
테스트 케이스를 데이터로 정의하고, 루프로 실행한다.

```go
// ✅ 올바른 예시
func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{name: "positive numbers", a: 2, b: 3, want: 5},
		{name: "negative numbers", a: -1, b: -2, want: -3},
		{name: "zero", a: 0, b: 0, want: 0},
		{name: "mixed signs", a: -1, b: 1, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// ❌ 잘못된 예시: 케이스마다 별도 테스트
func TestAddPositive(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Error("expected 5")
	}
}
func TestAddNegative(t *testing.T) {
	if Add(-1, -2) != -3 {
		t.Error("expected -3")
	}
}
```

### 11.2 테스트 함수 이름

테스트 함수명은 `TestXxx` 형식을 따른다.
서브 테스트는 `t.Run`으로 구분한다.

```go
// ✅ 올바른 예시
func TestUserService_CreateUser(t *testing.T) { ... }
func TestUserService_CreateUser_InvalidEmail(t *testing.T) { ... }
func TestParseConfig(t *testing.T) { ... }
func TestParseConfig_EmptyFile(t *testing.T) { ... }

// t.Run으로 서브 테스트
func TestUserService_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) { ... })
	t.Run("invalid email", func(t *testing.T) { ... })
	t.Run("duplicate user", func(t *testing.T) { ... })
}
```

### 11.3 testdata 디렉토리

테스트에서 사용하는 파일은 `testdata` 디렉토리에 넣는다.
Go 도구 체인은 `testdata` 디렉토리를 자동으로 무시한다.

```
mypackage/
├── parser.go
├── parser_test.go
└── testdata/
    ├── valid_config.json
    ├── invalid_config.json
    └── large_input.txt
```

```go
func TestParseConfig(t *testing.T) {
	data, err := os.ReadFile("testdata/valid_config.json")
	if err != nil {
		t.Fatalf("read test file: %v", err)
	}

	cfg, err := ParseConfig(data)
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}

	if cfg.Name != "test" {
		t.Errorf("Name = %q, want %q", cfg.Name, "test")
	}
}
```

### 11.4 테스트 헬퍼

테스트 헬퍼 함수는 반드시 `t.Helper()`를 호출한다.
이렇게 하면 테스트 실패 시 헬퍼 함수가 아닌 **호출한 줄**이 표시된다.

```go
// ✅ 올바른 예시
func assertNoError(t *testing.T, err error) {
	t.Helper()  // 반드시 호출한다
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func createTestUser(t *testing.T, name string) *User {
	t.Helper()
	user, err := NewUser(name, name+"@example.com")
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return user
}

// 사용
func TestSomething(t *testing.T) {
	user := createTestUser(t, "alice")
	err := doSomething(user)
	assertNoError(t, err)
	// 실패 시 이 줄이 표시된다 (헬퍼 내부가 아닌)
}
```

### 11.5 Golden File 테스트

복잡한 출력을 검증할 때 기대 결과를 파일로 저장하고 비교한다.
`-update` 플래그로 기대 파일을 업데이트한다.

```go
var update = flag.Bool("update", false, "update golden files")

func TestRender(t *testing.T) {
	result := Render(input)

	golden := filepath.Join("testdata", t.Name()+".golden")

	if *update {
		os.WriteFile(golden, []byte(result), 0644)
		return
	}

	expected, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}

	if result != string(expected) {
		t.Errorf("output mismatch.\ngot:\n%s\nwant:\n%s", result, expected)
	}
}
```

```bash
# golden 파일 업데이트
go test -run TestRender -update
```

### 11.6 벤치마크

벤치마크 함수는 `Benchmark` 접두사를 사용한다:

```go
func BenchmarkSort(b *testing.B) {
	data := generateTestData(10000)
	b.ResetTimer()  // 설정 시간 제외

	for i := 0; i < b.N; i++ {
		input := make([]int, len(data))
		copy(input, data)
		sort.Ints(input)
	}
}

func BenchmarkSort_ReportAllocs(b *testing.B) {
	b.ReportAllocs()  // 메모리 할당 보고

	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("hello %s", "world")
	}
}
```

```bash
# 벤치마크 실행
go test -bench=. -benchmem
```

### 11.7 Example 함수

`Example` 함수는 godoc에 실행 가능한 예시로 표시된다.
`// Output:` 주석으로 기대 출력을 명시하면 테스트로도 실행된다.

```go
func ExampleNewServer() {
	srv := NewServer(":8080")
	fmt.Println(srv.Addr)
	// Output: :8080
}

func ExampleServer_Start() {
	srv := NewServer(":8080")
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
```

---

## 12. 주석과 문서화

### 12.1 패키지 주석

모든 패키지는 **패키지 주석**을 가져야 한다.
`// Package xxx`로 시작하며, 패키지의 목적을 설명한다.

```go
// ✅ 올바른 예시
// Package auth provides authentication and authorization
// functionality for the application. It supports JWT tokens,
// OAuth2, and API key authentication.
package auth

// ✅ 올바른 예시: 간단한 패키지
// Package math provides basic mathematical functions.
package math

// ❌ 잘못된 예시: "Package"로 시작하지 않는다
// This package handles authentication.
package auth

// ❌ 잘못된 예시: 주석이 없다
package auth
```

큰 패키지의 경우 별도의 `doc.go` 파일에 패키지 주석을 작성한다:

```go
// doc.go

/*
Package http provides HTTP client and server implementations.

Get, Head, Post, and PostForm make HTTP (or HTTPS) requests:

	resp, err := http.Get("http://example.com/")
	...
*/
package http
```

### 12.2 함수/타입 주석

공개(exported) 함수와 타입은 반드시 주석을 달아야 한다.
주석은 **이름으로 시작**한다.

```go
// ✅ 올바른 예시: 이름으로 시작한다
// Server represents an HTTP server with graceful shutdown support.
type Server struct { ... }

// NewServer creates a new Server with the given address and options.
// If no options are provided, default values are used.
func NewServer(addr string, opts ...Option) *Server { ... }

// Start begins listening for incoming requests.
// It blocks until the server is shut down or an error occurs.
func (s *Server) Start() error { ... }

// ErrTimeout is returned when an operation exceeds its deadline.
var ErrTimeout = errors.New("operation timed out")

// MaxConnections defines the maximum number of concurrent connections.
const MaxConnections = 1000

// ❌ 잘못된 예시: 이름으로 시작하지 않는다
// This function creates a new server.
func NewServer(addr string) *Server { ... }

// Returns the current user.
func GetUser(id string) *User { ... }
```

### 12.3 godoc 규칙

godoc은 주석을 자동으로 문서화한다. 다음 규칙을 따른다:

- 주석의 첫 문장이 요약으로 사용된다
- 빈 줄로 단락을 구분한다
- 들여쓰기된 텍스트는 코드 블록으로 표시된다

```go
// ParseDuration parses a duration string.
//
// A duration string is a possibly signed sequence of decimal numbers,
// each with optional fraction and a unit suffix, such as "300ms",
// "-1.5h" or "2h45m".
//
// Valid time units are "ns", "us", "ms", "s", "m", "h".
//
// Example usage:
//
//	d, err := time.ParseDuration("1h30m")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ParseDuration(s string) (Duration, error) { ... }
```

### 12.4 TODO/FIXME

임시 코드나 개선이 필요한 부분에는 `TODO` 또는 `FIXME`를 사용한다.
작성자와 이슈 번호를 함께 기록한다.

```go
// TODO(username): Add rate limiting support. See #1234.

// FIXME(username): This doesn't handle Unicode properly. See #5678.
```

### 12.5 불필요한 주석 금지

코드 자체가 명확하면 주석을 달지 않는다.
주석은 "왜(why)"를 설명할 때 사용하고, "무엇을(what)"은 코드로 표현한다.

```go
// ❌ 잘못된 예시: 코드를 그대로 반복한다
// increment counter by 1
counter++

// check if user is nil
if user == nil {
	return ErrUserNotFound
}

// loop through items
for _, item := range items {
	process(item)
}

// ✅ 올바른 예시: "왜"를 설명한다
// We retry 3 times because the upstream service is occasionally flaky
// during deployments. See incident report #4521.
for i := 0; i < 3; i++ {
	if err := callUpstream(); err == nil {
		return nil
	}
}

// Use 32KB buffer size — benchmarks showed this is optimal for our
// typical payload sizes (see benchmark results in #3456).
buf := make([]byte, 32*1024)
```

### 12.6 Deprecated 표시

더 이상 사용하지 않는 함수/타입에는 `Deprecated` 주석을 달고
대안을 명시한다:

```go
// ✅ 올바른 예시
// ParseJSON parses JSON data into a map.
//
// Deprecated: Use encoding/json.Unmarshal instead.
func ParseJSON(data []byte) (map[string]any, error) { ... }

// ReadConfig reads configuration from the given path.
//
// Deprecated: Use LoadConfig instead, which supports multiple formats.
func ReadConfig(path string) (*Config, error) { ... }
```

---

## 13. 프로젝트 구조 (표준 레이아웃)

### 13.1 표준 디렉토리 구조

Go 프로젝트의 디렉토리 구조는 커뮤니티에서 널리 사용하는 관례를 따른다.
공식 표준은 아니지만 대부분의 대규모 프로젝트가 이 구조를 채택하고 있다.

```
myproject/
├── cmd/                    # 실행 파일 (main 패키지)
│   ├── server/
│   │   └── main.go         # 서버 실행 파일
│   └── cli/
│       └── main.go         # CLI 도구 실행 파일
│
├── internal/               # 비공개 패키지 (외부 import 불가)
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── database/
│   │   ├── postgres.go
│   │   └── postgres_test.go
│   ├── handler/
│   │   ├── user_handler.go
│   │   └── user_handler_test.go
│   ├── middleware/
│   │   └── logging.go
│   ├── model/
│   │   └── user.go
│   └── service/
│       ├── user_service.go
│       └── user_service_test.go
│
├── pkg/                    # 공개 라이브러리 (외부에서 import 가능, 선택적)
│   └── validator/
│       ├── validator.go
│       └── validator_test.go
│
├── api/                    # API 정의 파일
│   ├── proto/              # Protocol Buffers 정의
│   │   └── user.proto
│   └── openapi/            # OpenAPI/Swagger 정의
│       └── api.yaml
│
├── configs/                # 설정 파일 템플릿
│   ├── config.yaml
│   └── config.example.yaml
│
├── scripts/                # 빌드, 배포, CI/CD 스크립트
│   ├── build.sh
│   └── migrate.sh
│
├── migrations/             # 데이터베이스 마이그레이션
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
│
├── docs/                   # 추가 문서
│   └── architecture.md
│
├── testdata/               # 테스트 데이터
│   └── fixtures/
│
├── go.mod                  # 모듈 정의
├── go.sum                  # 의존성 체크섬
├── Makefile                # 빌드 자동화
├── Dockerfile              # 컨테이너 빌드
└── README.md               # 프로젝트 설명
```

### 13.2 각 디렉토리의 역할

#### cmd/

각 실행 파일마다 하위 디렉토리를 만든다.
`main.go`는 최소한의 코드만 포함하고, 실제 로직은 `internal/`에 위임한다.

```go
// cmd/server/main.go
package main

import (
	"log"
	"os"

	"myproject/internal/config"
	"myproject/internal/server"
)

func main() {
	cfg, err := config.Load(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := server.Run(cfg); err != nil {
		log.Fatalf("server: %v", err)
	}
}
```

#### internal/

프로젝트 내부에서만 사용하는 패키지를 넣는다.
Go 컴파일러가 외부 프로젝트에서의 import를 차단한다.
이것은 Go 도구 체인이 강제하는 접근 제어 메커니즘이다.

#### pkg/

외부 프로젝트에서도 사용할 수 있는 라이브러리 코드를 넣는다.
**주의**: `pkg/` 디렉토리 사용은 선택 사항이다.
프로젝트가 라이브러리가 아닌 애플리케이션이라면 `pkg/`를 생략하고
모든 코드를 `internal/`에 넣는 것이 더 안전하다.

#### api/

API 정의 파일(Protocol Buffers, OpenAPI/Swagger, GraphQL 스키마 등)을 넣는다.

### 13.3 간단한 프로젝트

소규모 프로젝트나 단일 바이너리 프로젝트는 간단한 구조를 사용한다.
과도한 디렉토리 구조는 오히려 복잡성을 증가시킨다.

```
# 작은 CLI 도구
mytool/
├── main.go
├── config.go
├── process.go
├── process_test.go
├── go.mod
└── go.sum

# 작은 라이브러리
mylib/
├── mylib.go
├── mylib_test.go
├── helper.go
├── go.mod
└── go.sum
```

---

## 14. Go Proverbs (Go 격언)

Rob Pike가 2015년 Gopherfest에서 발표한 Go 격언들이다.
각 격언은 Go의 설계 철학을 담고 있다.

**참조**: https://go-proverbs.github.io/

### 14.1 Don't communicate by sharing memory, share memory by communicating.

**메모리를 공유하여 통신하지 말고, 통신을 통해 메모리를 공유하라.**

전통적인 멀티스레드 프로그래밍에서는 공유 메모리와 잠금(lock)으로 통신한다.
Go에서는 채널을 통해 데이터를 전달하여 고루틴 간에 통신한다.

```go
// ✅ Go 방식: 채널로 통신한다
func producer(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i  // 데이터를 채널로 보낸다
	}
	close(ch)
}

func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)  // 채널에서 데이터를 받는다
	}
}
```

### 14.2 Concurrency is not parallelism.

**동시성은 병렬성이 아니다.**

동시성(concurrency)은 여러 작업을 **구조적으로 다루는 것**이다.
병렬성(parallelism)은 여러 작업을 **동시에 실행하는 것**이다.
동시성은 설계에 관한 것이고, 병렬성은 실행에 관한 것이다.

### 14.3 Channels orchestrate; mutexes serialize.

**채널은 조율하고, 뮤텍스는 직렬화한다.**

채널은 고루틴 간의 통신과 동기화를 **조율**하는 데 사용한다.
뮤텍스는 공유 리소스에 대한 접근을 **순차적으로 제한**하는 데 사용한다.

```go
// 채널: 조율 (orchestration)
done := make(chan struct{})
go func() {
	doWork()
	close(done)  // 작업 완료를 신호한다
}()
<-done  // 완료될 때까지 대기한다

// 뮤텍스: 직렬화 (serialization)
var mu sync.Mutex
mu.Lock()
counter++  // 한 번에 하나의 고루틴만 접근한다
mu.Unlock()
```

### 14.4 The bigger the interface, the weaker the abstraction.

**인터페이스가 클수록 추상화는 약해진다.**

`io.Reader`처럼 메서드가 하나인 인터페이스는 강력한 추상화를 제공한다.
파일, 네트워크 연결, 문자열, 버퍼 등 무엇이든 `Read` 메서드만 있으면 된다.
메서드가 많아질수록 구현할 수 있는 타입이 줄어들고 추상화의 힘이 약해진다.

### 14.5 Make the zero value useful.

**제로값을 유용하게 만들라.**

구조체의 제로값(모든 필드가 기본값인 상태)이 의미 있는 동작을 하도록 설계한다.

```go
// ✅ 제로값이 유용하다
var buf bytes.Buffer       // 초기화 없이 바로 사용 가능
buf.WriteString("hello")

var mu sync.Mutex          // 초기화 없이 바로 사용 가능
mu.Lock()
mu.Unlock()

var once sync.Once         // 초기화 없이 바로 사용 가능
once.Do(func() { ... })
```

### 14.6 interface{} says nothing.

**빈 인터페이스는 아무것도 말하지 않는다.**

`interface{}`(Go 1.18 이후 `any`)는 타입 정보를 전달하지 않으므로
가능한 사용을 피하고, 구체적인 타입이나 인터페이스를 사용한다.

### 14.7 Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.

**gofmt의 스타일은 누구의 취향도 아니지만, gofmt 자체는 모두의 취향이다.**

개인의 스타일 선호를 포기하는 대신, 모든 Go 코드가 일관된 형식을 가진다.
포매팅 논쟁에 시간을 낭비하지 않아도 된다.

### 14.8 A little copying is better than a little dependency.

**약간의 복사가 약간의 의존성보다 낫다.**

작은 유틸리티 함수를 위해 외부 패키지에 의존하느니,
그 함수를 직접 복사하는 것이 낫다. 의존성은 유지보수 비용을 수반한다.

```go
// ✅ 간단한 함수는 직접 작성한다
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ❌ 이 정도를 위해 외부 라이브러리를 가져올 필요 없다
import "github.com/someone/mathutil"
result := mathutil.Min(a, b)
```

### 14.9 Errors are values.

**에러는 값이다.**

에러는 예외(exception)가 아니라 일반적인 값이다.
프로그래밍으로 다룰 수 있고, 조합하고, 전달하고, 검사할 수 있다.

```go
// 에러를 값으로 다루는 패턴: errWriter
type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) write(buf []byte) {
	if ew.err != nil {
		return  // 이미 에러가 발생했으면 아무것도 하지 않는다
	}
	_, ew.err = ew.w.Write(buf)
}

// 사용: 에러 검사를 한 번만 하면 된다
ew := &errWriter{w: fd}
ew.write(header)
ew.write(body)
ew.write(footer)
if ew.err != nil {
	return ew.err
}
```

### 14.10 Don't just check errors, handle them gracefully.

**에러를 단순히 확인만 하지 말고, 우아하게 처리하라.**

에러를 받았을 때 단순히 위로 전달하는 것이 아니라,
컨텍스트를 추가하고 적절한 조치를 취해야 한다.

```go
// ❌ 단순히 확인만 한다
if err != nil {
	return err
}

// ✅ 우아하게 처리한다
if err != nil {
	return fmt.Errorf("save user %s: %w", user.ID, err)
}
```

### 14.11 Don't panic.

**panic하지 마라.**

`panic`은 정말로 복구 불가능한 프로그래밍 오류에서만 사용한다.
일반적인 에러 상황에서는 `error`를 반환한다.

### 14.12 Clear is better than clever.

**명확한 것이 영리한 것보다 낫다.**

트릭이나 기교를 부린 코드보다 읽기 쉬운 코드가 더 가치 있다.

```go
// ✅ 명확하다
func isEven(n int) bool {
	return n%2 == 0
}

// ❌ 영리하지만 불명확하다
func isEven(n int) bool {
	return n&1 == 0  // 비트 연산을 아는 사람만 이해한다
}
```

### 14.13 Reflection is never clear.

**리플렉션은 결코 명확하지 않다.**

`reflect` 패키지는 강력하지만, 코드를 이해하기 어렵게 만들고
컴파일 타임 타입 안전성을 잃게 한다.
반드시 필요한 경우(JSON 직렬화, ORM 등)에만 사용한다.

### 14.14 Don't just use concurrency because you can.

**할 수 있다고 해서 동시성을 사용하지 마라.**

고루틴이 가볍다고 해서 모든 곳에 사용하지 않는다.
동시성이 실제로 성능이나 설계를 개선하는 경우에만 사용한다.

### 14.15 Design the architecture, name the components, document the details.

**아키텍처를 설계하고, 컴포넌트에 이름을 붙이고, 세부 사항을 문서화하라.**

### 14.16 Documentation is for users.

**문서화는 사용자를 위한 것이다.**

주석과 문서는 코드를 작성한 자신이 아니라, 코드를 사용할 다른 사람을 위해 작성한다.

### 14.17 With the unsafe package there are no guarantees.

**unsafe 패키지를 사용하면 어떤 보장도 없다.**

`unsafe` 패키지는 Go의 타입 안전성과 메모리 안전성을 우회한다.
일반 애플리케이션에서는 사용하지 않는다.

---

## 15. 자주 하는 실수와 안티패턴

### 15.1 슬라이스와 맵의 nil vs empty 혼동

`nil` 슬라이스와 빈 슬라이스는 동작은 같지만 의미가 다르다.
JSON 직렬화에서 차이가 발생한다.

```go
// nil 슬라이스
var s []int          // s == nil, len(s) == 0
json.Marshal(s)      // "null"

// 빈 슬라이스
s := []int{}         // s != nil, len(s) == 0
json.Marshal(s)      // "[]"

s := make([]int, 0)  // s != nil, len(s) == 0
json.Marshal(s)      // "[]"

// ✅ API 응답에서 빈 배열이 필요하면 명시적으로 빈 슬라이스를 사용한다
func GetUsers() []User {
	users := fetchUsers()
	if users == nil {
		return []User{}  // JSON에서 null이 아닌 []를 반환한다
	}
	return users
}
```

맵도 마찬가지이다:

```go
var m map[string]int  // nil map — 읽기는 가능, 쓰기는 panic
m["key"]              // 0 반환 (panic 아님)
m["key"] = 1          // panic: assignment to entry in nil map

// ✅ 반드시 make로 초기화한다
m := make(map[string]int)
m["key"] = 1  // 정상 동작
```

### 15.2 고루틴에서 루프 변수 캡처 문제

Go 1.22 이전에는 루프 변수가 클로저에서 공유되어 문제가 발생했다.
Go 1.22부터 이 문제가 수정되었으나, 이전 버전에서는 주의해야 한다.

```go
// ❌ 잘못된 예시 (Go 1.21 이전)
for _, item := range items {
	go func() {
		process(item)  // 모든 고루틴이 마지막 item을 참조한다!
	}()
}

// ✅ 올바른 예시 (Go 1.21 이전): 변수를 로컬로 복사한다
for _, item := range items {
	item := item  // 루프 변수를 새 변수에 복사한다
	go func() {
		process(item)  // 각 고루틴이 자기만의 item을 가진다
	}()
}

// ✅ 올바른 예시 (Go 1.21 이전): 인자로 전달한다
for _, item := range items {
	go func(it Item) {
		process(it)
	}(item)
}

// ✅ Go 1.22 이후: 루프 변수가 반복마다 새로 생성되므로 문제없다
for _, item := range items {
	go func() {
		process(item)  // Go 1.22부터 안전하다
	}()
}
```

### 15.3 defer와 클로저 실수

`defer`는 함수가 반환될 때 실행되며, 클로저가 캡처한 변수의
**실행 시점**의 값을 사용한다.

```go
// ❌ 잘못된 예시: 항상 마지막 i 값(4)을 출력한다
func printNumbers() {
	for i := 0; i < 5; i++ {
		defer fmt.Println(i)  // 사실 이건 인자라서 올바르게 동작한다
	}
	// 출력: 4, 3, 2, 1, 0 (LIFO 순서)
}

// ❌ 진짜 문제: 클로저에서 변수 캡처
func openFiles(paths []string) {
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()  // 모든 파일이 함수 끝에서야 닫힌다!
	}
	// 많은 파일을 열면 파일 디스크립터가 고갈된다
}

// ✅ 올바른 예시: 별도 함수로 분리한다
func openFiles(paths []string) {
	for _, path := range paths {
		if err := processFile(path); err != nil {
			log.Println(err)
		}
	}
}

func processFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()  // 이 함수가 반환되면 즉시 닫힌다
	// ...
	return nil
}
```

### 15.4 인터페이스 nil 비교 함정

인터페이스 값은 (타입, 값) 쌍으로 구성된다.
타입이 설정되어 있으면 값이 nil이어도 인터페이스 자체는 nil이 아니다.

```go
// ❌ 함정: 인터페이스의 nil 비교
type MyError struct{}
func (e *MyError) Error() string { return "error" }

func getError() error {
	var err *MyError  // nil 포인터
	return err        // error 인터페이스로 반환 → nil이 아니다!
}

func main() {
	err := getError()
	fmt.Println(err == nil)  // false! (타입 정보가 있으므로)
}

// ✅ 올바른 예시: error 인터페이스를 직접 반환한다
func getError() error {
	var err *MyError
	if err == nil {
		return nil  // 명시적으로 nil을 반환한다
	}
	return err
}
```

### 15.5 불필요한 else 사용

`if` 블록에서 `return`하면 `else`가 불필요하다:

```go
// ❌ 잘못된 예시
func isPositive(n int) bool {
	if n > 0 {
		return true
	} else {
		return false
	}
}

// ✅ 올바른 예시
func isPositive(n int) bool {
	if n > 0 {
		return true
	}
	return false
}

// ✅ 더 좋은 예시
func isPositive(n int) bool {
	return n > 0
}
```

### 15.6 string과 []byte 불필요한 변환

`string`과 `[]byte` 간 변환은 메모리 할당을 수반한다.
루프 내에서 불필요한 변환을 반복하지 않는다.

```go
// ❌ 잘못된 예시: 불필요한 변환 반복
func processLines(data []byte) {
	lines := strings.Split(string(data), "\n")  // []byte → string
	for _, line := range lines {
		result := doSomething([]byte(line))  // string → []byte (매 반복마다)
		fmt.Println(string(result))          // []byte → string (매 반복마다)
	}
}

// ✅ 올바른 예시: bytes 패키지를 사용한다
func processLines(data []byte) {
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		result := doSomething(line)
		os.Stdout.Write(result)
	}
}
```

### 15.7 에러 무시

이미 6장에서 다루었지만, 가장 흔한 실수이므로 다시 강조한다.
에러를 무시하면 디버깅이 극도로 어려워진다.

```go
// ❌ 절대 하지 않는다
json.Unmarshal(data, &config)  // 에러 무시
http.ListenAndServe(":8080", nil)  // 에러 무시

// ✅ 반드시 에러를 처리한다
if err := json.Unmarshal(data, &config); err != nil {
	return fmt.Errorf("unmarshal config: %w", err)
}

if err := http.ListenAndServe(":8080", nil); err != nil {
	log.Fatalf("listen: %v", err)
}
```

### 15.8 과도한 추상화

Go는 실용적인 언어이다.
Java나 C# 스타일의 과도한 추상화 계층은 Go에 맞지 않는다.

```go
// ❌ 잘못된 예시: 과도한 추상화 (Java 스타일)
type UserRepositoryInterface interface { ... }
type UserRepositoryImpl struct { ... }
type UserServiceInterface interface { ... }
type UserServiceImpl struct { ... }
type UserControllerInterface interface { ... }
type UserControllerImpl struct { ... }
type UserDTOMapperInterface interface { ... }
type UserDTOMapperImpl struct { ... }

// ✅ 올바른 예시: Go답게 단순하게
type UserStore struct { db *sql.DB }
type UserService struct { store *UserStore }
type UserHandler struct { svc *UserService }
```

인터페이스는 **실제로 여러 구현이 필요할 때** 또는 **테스트에서 모킹이 필요할 때**만
도입한다. 미리 추상화하지 않는다.

---

## 16. 도구 생태계

### 16.1 go fmt / goimports

코드 포매팅 도구이다. 모든 Go 코드에 필수적으로 적용한다.

```bash
# go fmt: 표준 포매터
go fmt ./...

# goimports: go fmt + import 정리
goimports -w .
```

에디터에서 저장 시 자동 실행되도록 설정하는 것을 강력히 권장한다.

### 16.2 go vet

컴파일러가 잡지 못하는 의심스러운 코드를 검출한다.
CI/CD에서 필수적으로 실행해야 한다.

```bash
go vet ./...
```

검출하는 문제 예시:

```go
// go vet이 잡는 문제들
fmt.Printf("%d", "string")       // 형식 문자열 불일치
fmt.Printf("%s")                 // 인자 개수 부족
var mu sync.Mutex; _ = mu        // Mutex 복사
unreachableCode()                // 도달 불가능한 코드
```

### 16.3 staticcheck

Go를 위한 고급 정적 분석 도구이다.
`go vet`보다 더 많은 패턴을 검출한다.

```bash
# 설치
go install honnef.co/go/tools/cmd/staticcheck@latest

# 실행
staticcheck ./...
```

검출 예시:
- 사용되지 않는 코드
- 비효율적인 문자열 연결
- deprecated API 사용
- 불필요한 타입 변환

### 16.4 golangci-lint

여러 린터를 통합하여 실행하는 도구이다.
프로젝트에서 가장 널리 사용되는 린트 도구이다.

```bash
# 설치
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 실행
golangci-lint run ./...
```

권장하는 `.golangci.yml` 설정:

```yaml
linters:
  enable:
    - errcheck        # 에러 미처리 검출
    - govet           # go vet과 동일
    - staticcheck     # 정적 분석
    - unused          # 미사용 코드 검출
    - gosimple        # 코드 단순화 제안
    - ineffassign     # 무효한 할당 검출
    - typecheck       # 타입 에러 검출
    - gocritic        # 스타일 및 성능 제안
    - gofmt           # 포매팅 검사
    - goimports       # import 정리 검사
    - misspell        # 오타 검출
    - prealloc        # 슬라이스 사전 할당 제안
    - revive          # golint 대체

linters-settings:
  errcheck:
    check-type-assertions: true
  gocritic:
    enabled-tags:
      - diagnostic
      - style
      - performance

run:
  timeout: 5m

issues:
  max-issues-per-linter: 50
  max-same-issues: 10
```

### 16.5 go doc

패키지 문서를 터미널에서 확인하는 도구이다:

```bash
# 패키지 문서 확인
go doc fmt
go doc fmt.Println

# 소스 코드 포함
go doc -src fmt.Println

# 로컬 문서 서버 실행
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060
# 브라우저에서 http://localhost:6060 접속
```

### 16.6 go test -race

데이터 레이스를 검출하는 플래그이다.
테스트와 CI/CD에서 필수적으로 사용해야 한다.

```bash
# 레이스 감지 활성화
go test -race ./...

# 특정 테스트만
go test -race -run TestConcurrent ./...
```

레이스 감지기는 실행 시간에 동작하므로,
관련 코드 경로가 실행되는 테스트가 있어야 검출할 수 있다.

```go
// 이런 데이터 레이스를 검출한다
var counter int
var wg sync.WaitGroup

for i := 0; i < 1000; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		counter++  // 데이터 레이스! -race 플래그로 검출
	}()
}
wg.Wait()
```

### 16.7 go test -cover

코드 커버리지를 측정한다:

```bash
# 커버리지 요약
go test -cover ./...

# 커버리지 프로파일 생성
go test -coverprofile=coverage.out ./...

# HTML 보고서 생성
go tool cover -html=coverage.out -o coverage.html

# 함수별 커버리지 확인
go tool cover -func=coverage.out
```

### 16.8 go mod tidy

사용하지 않는 의존성을 제거하고 필요한 의존성을 추가한다.
커밋 전에 반드시 실행한다.

```bash
go mod tidy
```

```bash
# CI에서 go.mod이 정리되었는지 확인
go mod tidy
git diff --exit-code go.mod go.sum
```

### 16.9 go generate

코드 생성 명령을 실행한다.
소스 파일에 `//go:generate` 지시문으로 명령을 지정한다.

```go
//go:generate stringer -type=Color
//go:generate mockgen -source=repository.go -destination=mock_repository.go

type Color int

const (
	Red Color = iota
	Green
	Blue
)
```

```bash
# 모든 go:generate 지시문 실행
go generate ./...
```

### 16.10 도구 사용 요약

| 도구 | 용도 | 실행 시점 |
|------|------|-----------|
| `gofmt` / `goimports` | 코드 포매팅 | 저장 시 자동 실행 |
| `go vet` | 기본 정적 분석 | CI/CD 필수 |
| `staticcheck` | 고급 정적 분석 | CI/CD 권장 |
| `golangci-lint` | 통합 린터 | CI/CD 필수 |
| `go test -race` | 데이터 레이스 검출 | CI/CD 필수 |
| `go test -cover` | 커버리지 측정 | CI/CD 권장 |
| `go mod tidy` | 의존성 정리 | 커밋 전 필수 |
| `go generate` | 코드 생성 | 필요 시 |
| `go doc` | 문서 확인 | 개발 중 |

### 16.11 권장 CI/CD 파이프라인

```bash
#!/bin/bash
set -e

echo "=== Formatting Check ==="
goimports -l . | tee /dev/stderr | (! read)

echo "=== Vet ==="
go vet ./...

echo "=== Lint ==="
golangci-lint run ./...

echo "=== Test ==="
go test -race -cover -coverprofile=coverage.out ./...

echo "=== Module Tidy Check ==="
go mod tidy
git diff --exit-code go.mod go.sum

echo "=== Build ==="
go build ./...

echo "All checks passed!"
```

---

## 참고 자료

- **Effective Go**: https://go.dev/doc/effective_go
- **Go Code Review Comments**: https://go.dev/wiki/CodeReviewComments
- **Google Go Style Guide**: https://google.github.io/styleguide/go/
- **Uber Go Style Guide**: https://github.com/uber-go/guide/blob/master/style.md
- **Go Proverbs**: https://go-proverbs.github.io/
- **Go Blog**: https://go.dev/blog/
- **Go Standard Library**: https://pkg.go.dev/std
- **Go Specification**: https://go.dev/ref/spec

---

> 이 문서는 Go 커뮤니티의 공식 가이드와 모범 사례를 종합하여 정리한 것이다.
> 팀의 상황에 맞게 조정할 수 있으나, 핵심 원칙은 가능한 유지하는 것을 권장한다.
> Go는 **단순함**, **명확함**, **실용성**을 추구하는 언어이다.
> 코드를 작성할 때 항상 이 세 가지를 기억하라.
