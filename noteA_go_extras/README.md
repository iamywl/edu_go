# 노트 A. Go 문법 보충 수업

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run noteA_go_extras/for_range.go
go run noteA_go_extras/bufio_example.go
go run noteA_go_extras/embed_example.go

# cgo 예제 실행 (CGO_ENABLED=1 필요)
CGO_ENABLED=1 go run noteA_go_extras/cgo_example.go
```

> **참고**: `cgo_example.go`는 C 컴파일러가 필요하며, `CGO_ENABLED=1` 환경변수를 명시적으로 설정해야 한다. Docker 이미지에 C 컴파일러(gcc)가 설치되어 있어야 한다.

> **Makefile 활용**: `make run CH=noteA_go_extras` 또는 `make run CH=noteA_go_extras FILE=for_range.go`

---

본문에서 다루지 못한 Go 문법과 도구들을 보충한다.

---

## A.1 배열과 슬라이스

### 배열 (Array)

배열은 **고정 길이**의 동일 타입 요소 모음이다. 길이가 타입의 일부이므로 `[3]int`와 `[5]int`는 서로 다른 타입이다. 배열의 길이는 컴파일 타임에 결정되며, 런타임에 변경할 수 없다.

```go
var a [5]int            // 길이 5인 int 배열 (0으로 초기화)
b := [3]string{"Go", "Python", "Rust"}
c := [...]int{1, 2, 3}  // 컴파일러가 길이를 추론 → [3]int
d := [5]int{0: 10, 4: 50}  // 인덱스를 지정하여 초기화: [10, 0, 0, 0, 50]
```

**배열의 특성:**
- 값 타입이므로 함수에 전달하면 **전체가 복사**된다. 배열의 크기가 클수록 복사 비용이 크므로 주의해야 한다.
- 길이가 타입에 포함되므로 `[3]int`를 받는 함수에 `[5]int`를 전달할 수 없다.
- 비교 연산(`==`, `!=`)이 가능하다. 두 배열의 길이와 요소가 모두 같아야 `==`가 `true`이다.
- 배열은 연속된 메모리 공간에 저장되므로 CPU 캐시 친화적이며, 인덱스 접근 시 O(1)의 시간 복잡도를 가진다.

**다차원 배열:**

Go에서 다차원 배열도 선언할 수 있다:

```go
var matrix [3][3]int  // 3x3 정수 행렬
matrix[0][0] = 1
matrix[1][1] = 1
matrix[2][2] = 1      // 단위 행렬 생성
```

### 슬라이스 (Slice)

슬라이스는 배열 위에 만들어진 **가변 길이** 자료구조이다. 내부적으로 세 필드를 가진다:

```
┌──────────┬──────┬──────────┐
│ 포인터    │ 길이  │ 용량      │
│ (ptr)    │ (len)│ (cap)    │
└──────────┴──────┴──────────┘
      │
      ▼
  ┌───┬───┬───┬───┬───┐
  │ 0 │ 1 │ 2 │ 3 │ 4 │  ← 내부 배열 (backing array)
  └───┴───┴───┴───┴───┘
```

- **포인터(ptr)**: 슬라이스가 참조하는 내부 배열의 시작 위치이다.
- **길이(len)**: 현재 슬라이스에 포함된 요소 수이다. `len()` 함수로 확인한다.
- **용량(cap)**: 내부 배열에서 사용 가능한 최대 요소 수이다. `cap()` 함수로 확인한다.

```go
s := make([]int, 3, 5) // len=3, cap=5
s = append(s, 10)      // len=4, cap=5 (용량 내에서 추가)
s = append(s, 20, 30)  // len=6, cap=10 (용량 초과 → 새 배열 할당, 2배 성장)
```

**슬라이스 용량 증가 전략:**

`append`로 용량을 초과하면 Go 런타임은 새로운 내부 배열을 할당한다. Go 1.18 이후의 증가 전략은 다음과 같다:
- 현재 용량이 256 미만이면 대략 2배로 증가한다.
- 256 이상이면 약 1.25배 + 192 만큼 증가한다.
- 정확한 증가량은 요청한 크기와 메모리 정렬에 따라 달라진다.

**슬라이싱 연산:**

슬라이스에서 부분 슬라이스를 만들 때 세 번째 인덱스로 용량을 제한할 수 있다:

```go
s := []int{1, 2, 3, 4, 5}
sub := s[1:3:4]  // [2, 3], len=2, cap=3 (용량을 4-1=3으로 제한)
```

### 배열 vs 슬라이스 비교표

| 구분 | 배열 | 슬라이스 |
|------|------|----------|
| 길이 | 고정 (컴파일 타임) | 가변 (런타임) |
| 타입 | `[N]T` (길이 포함) | `[]T` |
| 전달 방식 | 값 복사 | 헤더 복사 (내부 배열 공유) |
| 비교 | `==` 가능 | `==` 불가 (`slices.Equal` 사용) |
| 생성 | `[3]int{1,2,3}` | `make([]int, 3)` 또는 `[]int{1,2,3}` |
| 메모리 위치 | 스택 또는 힙 | 헤더는 스택, 내부 배열은 힙 |

### 슬라이스 주의사항: 내부 배열 공유

```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:3] // [2, 3] — 같은 내부 배열을 공유!

sub[0] = 99
fmt.Println(original) // [1, 99, 3, 4, 5] ← original도 변경됨!
```

독립적인 복사본이 필요하면 `copy()`를 사용한다:

```go
independent := make([]int, len(sub))
copy(independent, sub)
```

**슬라이스 메모리 누수 주의:**

큰 슬라이스에서 작은 부분만 잘라서 사용하면, 큰 내부 배열이 GC에 의해 회수되지 않을 수 있다:

```go
// 메모리 누수 가능성 있는 코드
func getFirstTwo(data []int) []int {
    return data[:2]  // 원본 내부 배열 전체가 유지됨
}

// 안전한 코드
func getFirstTwo(data []int) []int {
    result := make([]int, 2)
    copy(result, data[:2])
    return result  // 새로운 작은 배열만 유지됨
}
```

---

## A.2 for range

Go에서 `for range`는 슬라이스, 배열, 맵, 문자열, 채널 등을 순회할 때 사용한다. 반복 가능한 자료구조를 간결하게 순회하는 Go의 핵심 문법이다.

### 기본 패턴

```go
// 1. 인덱스와 값 모두 사용
for i, v := range slice {
    fmt.Printf("index=%d, value=%d\n", i, v)
}

// 2. 인덱스만 사용
for i := range slice {
    fmt.Println(i)
}

// 3. 값만 사용 (인덱스 무시)
for _, v := range slice {
    fmt.Println(v)
}

// 4. 순회만 (값 불필요)
for range slice {
    fmt.Println("한 번 실행")
}
```

**중요한 특성:** `for range`에서 반환되는 값 `v`는 원본 요소의 **복사본**이다. 따라서 `v`를 수정해도 원본 슬라이스에는 영향이 없다. 원본을 수정하려면 인덱스를 통해 직접 접근해야 한다:

```go
nums := []int{1, 2, 3}
for i := range nums {
    nums[i] *= 2  // 인덱스를 사용하여 원본 수정
}
// nums: [2, 4, 6]
```

### 맵 순회

```go
m := map[string]int{"Go": 2009, "Rust": 2010, "Python": 1991}
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}
// 주의: 맵 순회 순서는 매번 다를 수 있다!
```

맵의 순회 순서가 비결정적인 것은 의도된 설계이다. Go 런타임이 의도적으로 순서를 무작위화하여, 개발자가 특정 순서에 의존하는 코드를 작성하지 않도록 유도한다. 정렬된 순서로 순회하려면 키를 별도로 정렬해야 한다.

### 문자열 순회 — 바이트 vs 룬

```go
s := "Go 언어"

// range는 UTF-8 룬(rune) 단위로 순회
for i, r := range s {
    fmt.Printf("byte offset=%d, rune=%c, code=%U\n", i, r, r)
}

// 바이트 단위 순회가 필요하면
for i := 0; i < len(s); i++ {
    fmt.Printf("byte[%d] = %x\n", i, s[i])
}
```

`range`로 문자열을 순회하면 UTF-8 멀티바이트 문자를 올바르게 처리한다. 한글 한 글자는 UTF-8에서 3바이트를 차지하므로, `len("언어")`는 6이지만 룬 수는 2이다.

### 채널 순회

```go
ch := make(chan int)
go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch) // 채널을 닫아야 range가 종료됨
}()

for v := range ch {
    fmt.Println(v)
}
```

채널에 대한 `for range`는 채널이 닫힐 때까지 값을 계속 수신한다. 채널을 닫지 않으면 영원히 블로킹되므로 반드시 `close()`를 호출해야 한다.

### Go 1.22+ 정수 range

Go 1.22부터 정수에 대한 range가 가능하다:

```go
for i := range 5 {
    fmt.Println(i) // 0, 1, 2, 3, 4
}
```

이 문법은 `for i := 0; i < 5; i++`와 동일한 동작을 하며, 더 간결한 표현이다.

> 자세한 예제는 `for_range.go`를 참고한다.

---

## A.3 입출력 처리

### os 패키지 — 파일 기본 조작

`os` 패키지는 운영체제와 상호작용하는 기본 기능을 제공한다. 파일 읽기/쓰기, 환경변수 접근, 프로세스 관리 등을 처리한다.

```go
// 파일 읽기
data, err := os.ReadFile("input.txt")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))

// 파일 쓰기
err = os.WriteFile("output.txt", []byte("Hello Go"), 0644)

// 표준 입출력
os.Stdin   // 표준 입력
os.Stdout  // 표준 출력
os.Stderr  // 표준 에러
```

**파일 권한 비트(`0644`)의 의미:**
- 첫 번째 자리(0): 8진수 표기법 접두사이다.
- 6 (owner): 읽기(4) + 쓰기(2) = 6이다.
- 4 (group): 읽기(4)이다.
- 4 (others): 읽기(4)이다.

### bufio 패키지 — 버퍼링된 입출력

`bufio`는 입출력에 버퍼를 추가하여 성능을 높인다. 특히 줄 단위 읽기에 유용하다. 버퍼링이란 작은 읽기/쓰기 요청을 모아서 한 번에 처리하는 것으로, 시스템 콜 횟수를 줄여 성능을 크게 향상시킨다.

```go
// Scanner로 줄 단위 읽기
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println("입력:", line)
}

// Scanner 에러 확인 (항상 체크해야 한다)
if err := scanner.Err(); err != nil {
    log.Fatal(err)
}

// Writer로 버퍼링된 쓰기
writer := bufio.NewWriter(os.Stdout)
fmt.Fprintln(writer, "버퍼링된 출력")
writer.Flush() // 반드시 Flush 호출!
```

**Scanner의 기본 버퍼 크기는 64KB이다.** 한 줄이 이보다 길면 `scanner.Buffer()`로 버퍼 크기를 조정해야 한다:

```go
scanner := bufio.NewScanner(file)
scanner.Buffer(make([]byte, 1024*1024), 1024*1024)  // 1MB 버퍼
```

### io 패키지 — Reader/Writer 인터페이스

Go의 I/O는 두 가지 핵심 인터페이스를 중심으로 설계되어 있다:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

이 인터페이스를 구현하는 타입은 모두 호환된다. 이것이 Go I/O 설계의 핵심이다. 파일, 네트워크 연결, 메모리 버퍼, 압축 스트림 등 모든 것이 동일한 인터페이스로 처리된다:

```go
// 파일 → 표준 출력으로 복사
f, _ := os.Open("input.txt")
defer f.Close()
io.Copy(os.Stdout, f) // os.Stdout은 Writer, f는 Reader

// 여러 Reader를 연결
r := io.MultiReader(strings.NewReader("Hello "), strings.NewReader("World"))
io.Copy(os.Stdout, r) // "Hello World"

// io.TeeReader: 읽으면서 동시에 다른 곳에 기록
var buf bytes.Buffer
tee := io.TeeReader(f, &buf)
io.Copy(os.Stdout, tee)  // 출력하면서 buf에도 저장
```

> 자세한 예제는 `bufio_example.go`를 참고한다.

---

## A.4 알아두면 유용한 go 명령어

### go run — 컴파일 + 실행을 한 번에

```bash
go run main.go           # 단일 파일 실행
go run .                 # 현재 디렉토리의 패키지 실행
go run ./cmd/server      # 특정 패키지 실행
```

임시로 컴파일하고 즉시 실행한다. 바이너리 파일이 남지 않는다. 개발 중 빠른 피드백을 위해 주로 사용한다.

### go build — 바이너리 빌드

```bash
go build -o myapp .      # myapp 이름으로 빌드
GOOS=linux GOARCH=amd64 go build -o myapp-linux .  # 크로스 컴파일
go build -ldflags="-s -w" -o myapp .  # 디버그 정보 제거하여 바이너리 크기 축소
```

`GOOS`와 `GOARCH` 환경변수를 조합하면 다양한 플랫폼용 바이너리를 생성할 수 있다. 예를 들어 macOS에서 Linux용 바이너리를, Windows에서 ARM용 바이너리를 만들 수 있다.

### go fmt — 코드 포맷팅

```bash
go fmt ./...             # 모든 패키지의 코드를 표준 스타일로 정리
gofmt -d .               # 변경될 내용을 diff로 확인
```

Go는 **공식 코드 스타일이 하나**이다. 팀 내 스타일 논쟁이 없다. 들여쓰기는 탭을 사용하며, 이는 `go fmt`가 강제한다.

### go vet — 정적 분석

```bash
go vet ./...             # 코드에서 의심스러운 패턴 검출
```

컴파일은 되지만 버그일 가능성이 높은 코드를 찾아준다:
- `fmt.Printf`의 형식 문자열 불일치
- 도달할 수 없는 코드
- 잘못된 구조체 태그
- 루프 변수 캡처 문제
- 잠금(lock) 복사 문제

`go vet`는 CI/CD 파이프라인에 반드시 포함시켜야 하는 도구이다.

### go install — 바이너리 설치

```bash
go install golang.org/x/tools/gopls@latest   # 도구 설치
go install .                                   # 현재 패키지를 $GOPATH/bin에 설치
```

설치된 바이너리는 `$GOPATH/bin` (기본값 `~/go/bin`)에 위치한다. 이 경로가 `$PATH`에 포함되어 있어야 어디서든 실행할 수 있다.

### go get — 의존성 추가/업데이트

```bash
go get github.com/gin-gonic/gin@latest        # 최신 버전 추가
go get github.com/gin-gonic/gin@v1.9.1        # 특정 버전 지정
go get -u ./...                                # 모든 의존성 업데이트
```

### go mod tidy — 모듈 정리

```bash
go mod tidy              # 사용하지 않는 의존성 제거, 필요한 의존성 추가
```

`go.mod`와 `go.sum`을 깔끔하게 유지하는 필수 명령어이다. 코드에서 import한 패키지가 `go.mod`에 없으면 자동으로 추가하고, import하지 않는 패키지는 제거한다.

### go test — 테스트 실행

```bash
go test ./...                    # 모든 패키지 테스트
go test -v ./...                 # 상세 출력
go test -run TestMyFunc ./...    # 특정 테스트만 실행
go test -cover ./...             # 커버리지 확인
go test -bench=. ./...           # 벤치마크 실행
go test -race ./...              # 데이터 레이스 감지
```

### 정리표

| 명령어 | 용도 |
|--------|------|
| `go run` | 컴파일 없이 바로 실행 |
| `go build` | 바이너리 생성 |
| `go fmt` | 코드 포맷팅 |
| `go vet` | 정적 분석 |
| `go install` | 바이너리 설치 |
| `go get` | 의존성 관리 |
| `go mod tidy` | 모듈 정리 |
| `go test` | 테스트 실행 |
| `go doc` | 문서 확인 |
| `go env` | Go 환경변수 확인 |
| `go generate` | 코드 생성 도구 실행 |

---

## A.5 cgo로 C 언어 호출하기

cgo는 Go 코드에서 C 라이브러리를 직접 호출할 수 있게 해주는 기능이다. C로 작성된 방대한 기존 라이브러리 생태계를 Go에서 활용할 수 있게 해주는 중요한 브릿지이다.

### 기본 사용법

```go
package main

/*
#include <stdio.h>
#include <stdlib.h>

void sayHello(const char* name) {
    printf("Hello from C, %s!\n", name);
}
*/
import "C"           // ← 바로 위의 주석이 C 코드로 해석됨
import "unsafe"

func main() {
    name := C.CString("Gopher")
    defer C.free(unsafe.Pointer(name))  // C 메모리는 수동 해제 필수!
    C.sayHello(name)
}
```

**핵심 규칙:**
- `import "C"`는 반드시 C 코드 주석 **바로 아래**에 위치해야 한다. 빈 줄이 있으면 안 된다.
- C에서 할당한 메모리(`C.CString` 등)는 반드시 `C.free`로 해제해야 한다. Go의 가비지 컬렉터는 C 힙 메모리를 관리하지 않는다.
- `import "C"`는 다른 import 문과 그룹으로 묶을 수 없다. 반드시 단독 import 문이어야 한다.

### CGO_ENABLED 환경변수

```bash
CGO_ENABLED=1 go build .   # cgo 활성화 (기본값, 네이티브 빌드 시)
CGO_ENABLED=0 go build .   # cgo 비활성화 (순수 Go 빌드)
```

크로스 컴파일 시에는 `CGO_ENABLED`가 자동으로 0이 된다. cgo를 사용하는 패키지를 크로스 컴파일하려면 해당 플랫폼의 C 컴파일러(크로스 컴파일러)가 필요하다. 이것이 cgo 사용 시 가장 큰 실용적 제약이다.

### 타입 매핑

| C 타입 | Go에서의 표현 |
|--------|--------------|
| `int` | `C.int` |
| `char` | `C.char` |
| `char*` | `*C.char` |
| `void*` | `unsafe.Pointer` |
| `size_t` | `C.size_t` |
| `long` | `C.long` |
| `double` | `C.double` |
| `float` | `C.float` |

### 문자열 변환

Go 문자열과 C 문자열 사이의 변환은 항상 복사가 발생한다:

```go
// Go → C: 메모리 할당 + 복사 (반드시 free 필요)
cStr := C.CString("hello")
defer C.free(unsafe.Pointer(cStr))

// C → Go: 복사만 발생 (free 불필요)
goStr := C.GoString(cStr)

// C → Go (길이 지정): 널 종료 문자가 없는 경우에 사용
goStr2 := C.GoStringN(cStr, 5)
```

### cgo의 장단점

**장점:**
- 기존 C 라이브러리(SQLite, OpenSSL 등)를 바로 활용 가능하다.
- 성능이 중요한 부분만 C로 작성 가능하다.
- 시스템 레벨 API에 직접 접근 가능하다.

**단점:**
- 크로스 컴파일이 어려워진다.
- 빌드 시간이 증가한다.
- Go의 가비지 컬렉터가 C 메모리를 관리하지 않는다.
- Go와 C 사이의 함수 호출 오버헤드가 존재한다 (일반 Go 함수 호출의 수십~수백 배).
- 디버깅이 복잡해진다 (Go 디버거와 C 디버거를 함께 사용해야 할 수 있다).
- 바이너리가 정적 링크 대신 동적 링크될 수 있어 배포가 복잡해진다.

> 자세한 예제는 `cgo_example.go`를 참고한다.

---

## A.6 go doc

Go는 코드 자체가 문서가 되는 철학을 따른다. 별도의 문서화 도구 없이 주석만으로 문서를 생성할 수 있다. 이 접근 방식은 코드와 문서의 동기화 문제를 근본적으로 해결한다.

### 문서화 규칙

1. **패키지 주석**: 패키지 선언 바로 위에 작성한다.

```go
// Package mathutil은 수학 관련 유틸리티 함수를 제공한다.
package mathutil
```

2. **함수/타입 주석**: 이름으로 시작하는 주석을 작성한다.

```go
// Add는 두 정수를 더한 결과를 반환한다.
// 오버플로가 발생하면 결과는 정의되지 않는다.
func Add(a, b int) int {
    return a + b
}

// Calculator는 기본 산술 연산을 수행하는 구조체이다.
type Calculator struct {
    // Precision은 소수점 이하 자릿수를 지정한다.
    Precision int
}
```

3. **Deprecated 표시**: 더 이상 사용하지 않는 함수를 표시한다.

```go
// OldFunc는 이전 버전의 처리 함수이다.
//
// Deprecated: NewFunc를 대신 사용한다.
func OldFunc() {}
```

4. **예제 함수**: 테스트 파일에 `Example` 접두사로 작성한다.

```go
func ExampleAdd() {
    fmt.Println(Add(2, 3))
    // Output: 5
}
```

예제 함수는 `go test` 실행 시 자동으로 검증되므로, 항상 최신 상태의 실행 가능한 문서 역할을 한다.

### go doc 명령어 사용법

```bash
go doc fmt              # fmt 패키지 문서
go doc fmt.Println      # 특정 함수 문서
go doc -all fmt         # 모든 심볼 포함
go doc -src fmt.Println # 소스 코드까지 표시
```

### godoc 웹 서버

```bash
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060       # http://localhost:6060 에서 문서 열람
```

로컬에서 자신의 프로젝트 문서를 웹 브라우저로 확인할 수 있으며, 표준 라이브러리 문서도 오프라인으로 열람 가능하다.

---

## A.7 Embed

Go 1.16에서 추가된 `embed` 패키지를 사용하면 파일을 바이너리에 직접 포함시킬 수 있다. 이 기능을 통해 설정 파일, HTML 템플릿, 정적 자산 등을 단일 바이너리에 포함하여 배포를 단순화할 수 있다.

### 기본 사용법

```go
package main

import (
    _ "embed"
    "fmt"
)

//go:embed hello.txt
var message string  // 파일 내용이 문자열로 포함됨

func main() {
    fmt.Println(message)
}
```

### 다양한 임베드 방식

```go
import "embed"

//go:embed version.txt
var version string          // string: 텍스트 파일

//go:embed logo.png
var logo []byte             // []byte: 바이너리 파일

//go:embed templates/*
var templates embed.FS      // embed.FS: 디렉토리 전체
```

각 방식의 사용 시나리오:
- `string`: 설정 파일, SQL 쿼리, 버전 정보 등 텍스트 데이터에 적합하다.
- `[]byte`: 이미지, 폰트, 인증서 등 바이너리 데이터에 적합하다.
- `embed.FS`: 여러 파일을 포함하는 디렉토리(웹 정적 파일, 템플릿 모음 등)에 적합하다.

### embed.FS 활용

```go
//go:embed static/*
var staticFiles embed.FS

func main() {
    // 개별 파일 읽기
    data, err := staticFiles.ReadFile("static/index.html")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(data))

    // http.FileServer와 함께 사용
    http.Handle("/", http.FileServer(http.FS(staticFiles)))
}
```

### 여러 패턴 지정

하나의 변수에 여러 `//go:embed` 디렉티브를 사용할 수 있다:

```go
//go:embed static/*
//go:embed templates/*
//go:embed config.json
var content embed.FS
```

### 주의사항

- `//go:embed` 디렉티브와 `//` 사이에 공백이 있으면 안 된다 (일반 주석이 된다).
- 임베드 경로는 현재 패키지 디렉토리 기준 상대 경로이다.
- `.`이나 `_`로 시작하는 파일은 기본적으로 제외된다 (`all:` 접두사로 포함 가능하다).
- 임베드된 파일은 빌드 시점에 바이너리에 포함되므로, 바이너리 크기가 커질 수 있다.
- 임베드 변수는 반드시 패키지 수준에서 선언해야 한다 (함수 내 지역 변수 불가).

```go
//go:embed all:templates  // 숨김 파일도 포함
var allTemplates embed.FS
```

> 자세한 예제는 `embed_example.go`를 참고한다.

---

## 연습문제

### 개념 문제

**1.** `[3]int`와 `[]int`의 차이점을 세 가지 이상 서술하라.

**2.** 다음 코드의 출력을 예측하고, 그 이유를 설명하라:

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]
b = append(b, 99)
fmt.Println(a)
fmt.Println(b)
```

**3.** `for range`에서 반환되는 값이 원본의 복사본인 이유를 Go의 설계 철학과 연관지어 설명하라.

**4.** `bufio.Scanner`와 `bufio.Reader`의 차이점은 무엇이며, 각각 어떤 상황에서 사용하는 것이 적합한가?

**5.** `io.Reader`와 `io.Writer` 인터페이스가 Go I/O 설계의 핵심인 이유를 설명하라. 이 인터페이스를 사용함으로써 얻는 이점을 세 가지 이상 서술하라.

**6.** cgo를 사용할 때 `C.CString()`으로 만든 문자열을 반드시 `C.free()`로 해제해야 하는 이유를 Go의 메모리 관리 모델과 연관지어 설명하라.

**7.** `go vet`와 `go fmt`의 역할 차이를 설명하고, 두 도구를 모두 CI/CD에 포함시켜야 하는 이유를 서술하라.

**8.** `embed.FS`로 파일을 바이너리에 포함시키는 방식의 장점과 단점을 각각 세 가지씩 서술하라.

### 코딩 문제

**9.** 슬라이스 `s := []int{10, 20, 30, 40, 50}`에서 인덱스 2의 요소를 삭제하는 함수를 작성하라. 원본 슬라이스의 순서를 유지해야 한다.

**10.** `for range`를 사용하여 문자열에 포함된 한글 글자 수만 세는 함수를 작성하라. (힌트: 유니코드 범위 `0xAC00` ~ `0xD7A3`)

---

## 구현 과제

**1. 줄 번호 출력기:** `bufio.Scanner`를 사용하여 파일을 읽고, 각 줄 앞에 줄 번호를 붙여 출력하는 프로그램을 작성하라. 파일명은 명령줄 인수로 받는다. 줄 번호는 오른쪽 정렬하여 출력한다.

**2. 슬라이스 유틸리티:** 제네릭을 사용하여 다음 슬라이스 유틸리티 함수들을 구현하라:
- `Filter[T any](s []T, pred func(T) bool) []T`
- `Map[T, U any](s []T, fn func(T) U) []U`
- `Reduce[T, U any](s []T, init U, fn func(U, T) U) U`

**3. 간단한 C 라이브러리 래퍼:** cgo를 사용하여 C의 `math.h`에서 `sqrt`, `pow`, `sin`, `cos` 함수를 래핑하는 Go 패키지를 만들어라. 각 함수에 대해 Go 스타일의 래퍼 함수와 문서 주석을 작성한다.

**4. 파일 임베드 웹 서버:** `embed` 패키지를 사용하여 정적 HTML/CSS 파일을 내장한 웹 서버를 만들어라. `static/` 디렉토리에 `index.html`과 `style.css`를 두고, 이를 바이너리에 포함시켜 단일 실행 파일로 배포 가능하게 한다.

**5. go doc 준수 패키지:** 수학 유틸리티 패키지(`mathutil`)를 만들고, 모든 공개 함수와 타입에 `go doc` 규칙을 준수하는 주석을 작성하라. `Example` 테스트 함수도 포함하여 `godoc`으로 문서를 확인할 수 있도록 한다.

---

## 프로젝트 과제

**1. 다중 포맷 파일 변환기:** 명령줄에서 입력 파일과 출력 포맷을 받아 파일을 변환하는 도구를 만들어라. `io.Reader`/`io.Writer` 인터페이스를 활용하여 JSON, CSV, YAML 포맷 간 변환을 지원한다. `embed`로 기본 변환 규칙 설정 파일을 내장하고, `go doc` 규칙에 맞게 모든 공개 API를 문서화한다. `bufio`를 사용하여 대용량 파일도 효율적으로 처리한다.

**2. Go 프로젝트 분석 도구:** 주어진 Go 프로젝트 디렉토리를 분석하여 다음 정보를 리포트하는 CLI 도구를 만들어라:
- 각 패키지의 파일 수, 함수 수, 테스트 커버리지
- 사용된 `go` 명령어 히스토리 기반 빌드 설정 분석
- 의존성 목록과 각 의존성의 라이선스 정보
- 결과를 HTML 리포트로 출력하되, 템플릿은 `embed`로 바이너리에 포함한다.
