# 14장. 패키지 (Package)

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행 (패키지 구조이므로 디렉토리 단위로 실행한다)
cd ch14_package && go run .
```

> **Makefile 활용**: `make run CH=ch14_package`

이 챕터는 `main.go`에서 하위 패키지(`mymath/`, `greeting/`)를 import하여 사용하는 구조이다. 개별 파일이 아닌 디렉토리 단위로 실행해야 한다.

---

**패키지(package)**는 Go에서 코드를 조직화하는 기본 단위이다. 관련 기능을 하나의 패키지로 묶어 재사용성과 유지보수성을 높인다. Go의 패키지 시스템은 단순하면서도 강력하며, 대규모 프로젝트에서 코드를 효과적으로 관리할 수 있게 해준다.

---

## 14.1 패키지 (코드 조직화)

모든 Go 파일은 반드시 하나의 패키지에 속해야 한다. 파일의 첫 번째 문장(주석 제외)은 반드시 `package` 선언이어야 한다.

```go
package main  // 이 파일은 main 패키지에 속함
```

### 패키지의 종류

| 종류 | 설명 | 예시 |
|------|------|------|
| **main 패키지** | 실행 가능한 프로그램의 진입점 | `package main` + `func main()` |
| **라이브러리 패키지** | 다른 패키지에서 import하여 사용 | `package fmt`, `package mymath` |

`main` 패키지는 특별하다. `package main`을 선언한 패키지에서 `func main()`을 정의하면, 이것이 프로그램의 시작점이 된다. 라이브러리 패키지는 직접 실행할 수 없으며, 다른 패키지에서 import하여 사용한다.

### 패키지 = 디렉토리

- 하나의 디렉토리 = 하나의 패키지이다.
- 같은 디렉토리의 모든 `.go` 파일은 같은 패키지에 속해야 한다.
- 디렉토리 이름과 패키지 이름이 반드시 같을 필요는 없지만, **관례적으로 같게** 하는 것이 좋다.

```
myproject/
├── go.mod
├── main.go          // package main
├── mymath/
│   └── mymath.go    // package mymath
└── greeting/
    └── greeting.go  // package greeting
```

### 같은 패키지 내 파일 간 접근

같은 패키지(같은 디렉토리)에 속한 파일들은 서로의 함수, 변수, 타입에 자유롭게 접근할 수 있다. 별도의 import가 필요 없다.

```
mymath/
├── add.go      // package mymath — func Add()
└── subtract.go // package mymath — func Subtract(), Add() 사용 가능
```

---

## 14.2 패키지 사용하기 (import)

### 기본 import

```go
import "fmt"
```

### 여러 패키지 import

```go
import (
    "fmt"
    "math"
    "strings"
)
```

괄호를 사용하여 여러 패키지를 묶어 import하는 것이 **권장 방식**이다.

### 로컬 패키지 import

```go
import (
    "myproject/mymath"      // 모듈 경로 기준
    "myproject/greeting"
)
```

로컬 패키지를 import할 때는 `go.mod`에 정의된 모듈 경로를 기준으로 한다. 상대 경로(`./mymath`)는 사용하지 않는다.

### import 별칭

```go
import (
    f "fmt"                 // 별칭 사용
    _ "database/sql/driver" // 사이드 이펙트만 (init 함수 실행)
    . "math"                // 패키지명 없이 사용 (비권장)
)

f.Println("별칭 사용")
Sqrt(16)  // math.Sqrt 대신 — . import 사용 시
```

| import 방식 | 설명 | 사용 상황 |
|-------------|------|----------|
| `import "fmt"` | 기본 import | 일반적인 경우 |
| `import f "fmt"` | 별칭 지정 | 패키지명 충돌 시 |
| `import _ "pkg"` | 빈 import | init() 함수만 실행할 때 |
| `import . "pkg"` | dot import | 테스트 등 특수한 경우 (비권장) |

### 사용하지 않는 import는 에러

Go에서는 import한 패키지를 사용하지 않으면 **컴파일 에러**가 발생한다. 이는 불필요한 의존성을 방지하기 위한 Go의 설계 철학이다.

```go
import "fmt"  // fmt를 사용하지 않으면 컴파일 에러!
```

개발 중 임시로 사용하지 않는 패키지가 있다면 `_`(빈 식별자)를 활용한다.

---

## 14.3 Go 모듈 (go mod)

### 모듈 초기화

```bash
go mod init myproject
```

이 명령은 `go.mod` 파일을 생성한다:

```
module myproject

go 1.21
```

### go.mod 파일의 역할

- **모듈 경로**: 이 프로젝트의 import 경로를 정의한다.
- **Go 버전**: 사용하는 Go 최소 버전을 명시한다.
- **의존성 관리**: 외부 패키지의 버전 정보를 기록한다.

### 외부 패키지 사용

```bash
go get github.com/some/package  # 외부 패키지 다운로드
go mod tidy                      # 사용하지 않는 의존성 정리
```

`go mod tidy`는 코드에서 실제로 사용하는 패키지만 `go.mod`에 남기고 나머지를 정리한다. 새로운 import를 추가한 후에는 `go mod tidy`를 실행하는 것이 좋다.

### go.sum 파일

`go.sum` 파일은 의존성 패키지의 **체크섬(hash)**을 저장한다. 이를 통해 의존성의 무결성을 검증한다. 이 파일은 자동으로 관리되므로 직접 수정할 필요가 없다.

### 주요 go mod 명령어

| 명령어 | 설명 |
|--------|------|
| `go mod init` | 모듈 초기화 |
| `go mod tidy` | 의존성 정리 |
| `go mod download` | 의존성 다운로드 |
| `go mod vendor` | vendor 디렉토리 생성 |
| `go mod graph` | 의존성 그래프 출력 |

---

## 14.4 패키지명과 패키지 외부 공개 (대문자 규칙)

Go에서는 **이름의 첫 글자**로 외부 공개 여부를 결정한다. 이것은 Go의 가장 독특한 특징 중 하나로, `public`/`private` 같은 별도의 키워드가 필요 없다.

### 대문자 = 공개 (Exported)

```go
package mymath

// Add 는 외부에서 사용 가능 (대문자)
func Add(a, b int) int {
    return a + b
}

// Pi 는 외부에서 접근 가능
const Pi = 3.14159
```

### 소문자 = 비공개 (Unexported)

```go
package mymath

// validate 는 패키지 내부에서만 사용 (소문자)
func validate(n int) bool {
    return n >= 0
}

// secret 은 외부에서 접근 불가
var secret = "내부 비밀"
```

### 외부에서 사용 시

```go
package main

import "myproject/mymath"

func main() {
    mymath.Add(1, 2)     // OK — 대문자
    // mymath.validate(1) // 컴파일 에러! — 소문자
}
```

### 규칙 요약

| 대상 | 대문자 시작 | 소문자 시작 |
|------|-----------|-----------|
| 함수 | 외부 공개 | 패키지 내부 전용 |
| 변수/상수 | 외부 공개 | 패키지 내부 전용 |
| 구조체 | 외부 공개 | 패키지 내부 전용 |
| 구조체 필드 | 외부 공개 | 패키지 내부 전용 |
| 인터페이스 | 외부 공개 | 패키지 내부 전용 |
| 메서드 | 외부 공개 | 패키지 내부 전용 |

### 구조체 필드의 공개/비공개

구조체 자체와 필드의 공개 여부는 독립적이다.

```go
type User struct {
    Name     string  // 외부 공개
    Email    string  // 외부 공개
    password string  // 패키지 내부 전용
}
```

외부에서는 `User.Name`과 `User.Email`에는 접근할 수 있지만, `User.password`에는 접근할 수 없다. 이를 통해 캡슐화를 구현한다.

---

## 14.5 패키지 초기화 (init 함수)

`init()` 함수는 패키지가 import될 때 **자동으로 실행**된다.

```go
package greeting

import "fmt"

var DefaultLang string

func init() {
    // 패키지가 로드될 때 자동 실행
    DefaultLang = "ko"
    fmt.Println("[greeting] init 실행: 기본 언어 설정 완료")
}
```

### init() 함수 특징

1. **매개변수와 반환값이 없다.**
2. **직접 호출할 수 없다** — import 시 자동 실행된다.
3. 한 파일에 **여러 init() 함수**를 정의할 수 있다. 같은 파일 내에서는 위에서 아래 순서로 실행된다.
4. 실행 순서: **import된 패키지의 init() → main 패키지의 init() → main()**

### 실행 순서

```
1. 의존 패키지의 전역 변수 초기화
2. 의존 패키지의 init() 실행
3. main 패키지의 전역 변수 초기화
4. main 패키지의 init() 실행
5. main() 함수 실행
```

### init() 함수 사용 예시

```go
// 데이터베이스 드라이버 등록
func init() {
    sql.Register("mydriver", &MyDriver{})
}

// 설정 파일 로드
func init() {
    config = loadConfig("config.json")
}

// 유효성 검증
func init() {
    if os.Getenv("API_KEY") == "" {
        log.Fatal("API_KEY 환경변수가 설정되지 않았다")
    }
}
```

### init() 사용 시 주의사항

- init()에서 무거운 작업(네트워크 요청, 대량 파일 I/O 등)을 수행하면 프로그램 시작이 느려진다.
- init() 간의 실행 순서에 의존하는 코드는 피해야 한다.
- 가능하다면 명시적인 초기화 함수를 사용하는 것이 테스트하기 더 쉽다.

---

## 14.6 패키지 설계 원칙

좋은 패키지를 설계하기 위한 원칙들을 정리한다.

### 패키지 이름 규칙

- **짧고 간결하게**: `stringutil`보다 `strutil`이 낫다.
- **소문자만 사용**: `myPackage`가 아닌 `mypackage`를 사용한다.
- **언더스코어/하이픈 지양**: `my_package`가 아닌 `mypackage`를 사용한다.
- **의미 있는 이름**: 패키지의 역할을 나타내는 이름을 사용한다.

### 순환 의존 금지

Go에서는 패키지 간 **순환 의존(circular dependency)**이 허용되지 않는다. 패키지 A가 B를 import하면, B는 A를 import할 수 없다. 이 제약은 의존 관계를 명확하게 유지하는 데 도움을 준다.

```
// 허용되지 않음
A → B → A  (순환 의존!)

// 해결 방법: 공통 부분을 별도 패키지로 분리
A → C
B → C
```

---

## 핵심 요약

1. 모든 Go 파일은 하나의 패키지에 속하며, 디렉토리 단위로 구분한다.
2. `import`로 다른 패키지를 가져와 사용한다.
3. `go mod init`으로 모듈을 초기화하고 `go.mod`로 의존성을 관리한다.
4. **대문자**로 시작하면 외부 공개, **소문자**로 시작하면 내부 전용이다.
5. `init()` 함수는 패키지 로드 시 자동 실행되며, 초기화 작업에 사용한다.

---

## 연습문제

### 문제 1: 패키지 만들기
`calculator` 패키지를 만들고 `Add`, `Subtract`, `Multiply`, `Divide` 함수를 구현하라. `Divide`는 0으로 나누면 에러를 반환하도록 하라.

### 문제 2: 공개/비공개
구조체 `User`를 만들되, `Name`은 공개, `password`는 비공개로 하라. 비공개 필드에 접근하는 `SetPassword`와 `CheckPassword` 함수를 공개로 만들라.

### 문제 3: init 함수
설정 값을 전역 변수에 저장하는 `init()` 함수를 작성하라. main에서 해당 전역 변수를 출력하여 init이 실행되었는지 확인하라.

### 문제 4: import 별칭
두 개의 서로 다른 패키지에 같은 이름의 함수가 있는 상황을 만들고, import 별칭을 사용하여 충돌을 해결하는 코드를 작성하라.

### 문제 5: 사용하지 않는 import
Go에서 사용하지 않는 import가 컴파일 에러를 발생시키는 이유를 설명하라. 개발 중에 임시로 import를 유지하는 방법(`_` 사용)을 예제 코드로 보여라.

### 문제 6: go mod 실습
새로운 디렉토리에서 `go mod init`으로 모듈을 초기화하고, `main.go`와 `mathutil/mathutil.go` 파일을 만들어 로컬 패키지를 import하는 프로젝트를 구성하라. 디렉토리 구조와 각 파일의 내용을 작성하라.

### 문제 7: 외부 패키지 사용
`go get`으로 외부 패키지(예: `golang.org/x/text`)를 다운로드하고, 해당 패키지를 사용하는 프로그램을 작성하라. `go.mod`와 `go.sum` 파일에 어떤 변화가 생기는지 확인하라.

### 문제 8: 패키지 문서화
GoDoc 형식에 맞는 주석을 작성하는 연습을 하라. 패키지 주석, 함수 주석, 타입 주석을 포함하는 `converter` 패키지를 작성하고, `go doc` 명령어로 문서를 확인하라.

### 문제 9: 순환 의존 해결
패키지 A가 B를, B가 A를 import해야 하는 상황을 가정하라. 이 순환 의존을 해결하기 위한 3가지 방법(공통 패키지 분리, 인터페이스 활용, 의존성 역전)을 각각 설명하고 코드 구조를 제시하라.

### 문제 10: 패키지 이름 규칙
다음 패키지 이름 중 Go 관례에 맞는 것과 맞지 않는 것을 구분하고, 맞지 않는 경우 올바른 이름을 제안하라.
- `myPackage`, `string_utils`, `httphandler`, `DB`, `io`, `net-utils`, `myMath`

---

## 구현 과제

### 과제 1: 유틸리티 패키지 모음
다음 3개의 유틸리티 패키지를 만들고, `main` 패키지에서 모두 사용하는 프로그램을 작성하라.
- `strutil`: 문자열 관련 유틸리티 (첫 글자 대문자 변환, 문자열 반복, 단어 수 세기)
- `mathutil`: 수학 관련 유틸리티 (최대공약수, 최소공배수, 소수 판별)
- `validate`: 검증 관련 유틸리티 (이메일 형식 검증, 비밀번호 강도 검사)
- 각 패키지는 공개/비공개 함수를 적절히 사용하고, 패키지 주석을 작성할 것

### 과제 2: 설정 관리 패키지
`config` 패키지를 만들어 프로그램 설정을 관리하는 시스템을 구현하라.
- `init()` 함수에서 기본 설정을 로드한다.
- `Get(key string) string`과 `Set(key, value string)` 함수를 제공한다.
- 설정값은 패키지 내부의 비공개 map에 저장한다.
- 존재하지 않는 키에 접근하면 기본값을 반환하는 `GetOrDefault(key, defaultVal string) string`을 구현하라.

### 과제 3: 로그 패키지
간단한 로그 패키지를 구현하라.
- 로그 레벨 지원: DEBUG, INFO, WARN, ERROR
- `init()` 함수에서 기본 로그 레벨을 설정한다.
- `SetLevel(level int)`, `Debug(msg string)`, `Info(msg string)`, `Warn(msg string)`, `Error(msg string)` 함수를 공개한다.
- 설정된 레벨 이상의 로그만 출력한다.
- 출력 형식: `[2024-01-15 10:30:00] [INFO] 메시지 내용`

### 과제 4: 단위 변환 패키지
`converter` 패키지를 만들어 다양한 단위 변환 기능을 제공하라.
- 온도 변환: 섭씨/화씨/켈빈
- 길이 변환: 미터/피트/인치/킬로미터/마일
- 무게 변환: 킬로그램/파운드/온스
- 각 변환 함수는 공개하고, 내부 변환 상수는 비공개로 유지하라.

---

## 프로젝트 과제

### 프로젝트 1: CLI 할일 관리 도구
여러 패키지로 구성된 CLI(Command Line Interface) 할일 관리 도구를 만들라.
- `main` 패키지: CLI 인터페이스와 사용자 입력 처리
- `todo` 패키지: Todo 구조체, 추가/삭제/완료/목록 기능
- `storage` 패키지: 파일 기반 데이터 저장/불러오기 (JSON 형식)
- `display` 패키지: 할일 목록 포맷팅 및 출력
- 각 패키지의 공개/비공개 요소를 적절히 설계하고, init() 함수를 활용하여 초기 데이터를 로드할 것

### 프로젝트 2: 모듈형 계산기
여러 패키지로 분리된 계산기 프로그램을 구현하라.
- `main` 패키지: 사용자 인터페이스 (수식 입력, 결과 출력)
- `parser` 패키지: 문자열 수식을 토큰으로 분리
- `evaluator` 패키지: 토큰을 평가하여 결과 계산 (사칙연산, 괄호 지원)
- `history` 패키지: 계산 이력 저장 및 조회
- 패키지 간 순환 의존이 발생하지 않도록 설계에 주의할 것
