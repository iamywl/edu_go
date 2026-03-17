# 28장 [Project] HTTP 웹 서버 만들기

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회, 포트 노출 필요: -p 8080:8080)
make shell

# 예제 실행
go run ch28_project_http_server/main.go
go run ch28_project_http_server/query.go
go run ch28_project_http_server/json_handler.go
go run ch28_project_http_server/fileserver.go

# 테스트 실행
go test ch28_project_http_server/ -v

# 서버 실행 후 다른 터미널에서 테스트
curl http://localhost:8080/
curl http://localhost:8080/hello?name=Go
```

> **참고**: HTTP 서버를 Docker 컨테이너 내부에서 실행할 때는 컨테이너의 포트를 호스트에 노출해야 외부에서 접근할 수 있다.

> **Makefile 활용**: `make run CH=ch28_project_http_server` 또는 `make run CH=ch28_project_http_server FILE=main.go`

---

Go의 `net/http` 패키지는 별도의 외부 라이브러리 없이도 강력한 HTTP 서버를 만들 수 있게 한다. 이 장에서는 기본 웹 서버부터 HTTPS까지 단계적으로 HTTP 서버를 구축한다. Go의 `net/http` 패키지는 프로덕션 수준의 HTTP 서버를 구현하기에 충분한 기능을 제공하며, 동시성 처리도 자동으로 지원하므로 별도의 설정 없이도 높은 성능을 달성할 수 있다.

---

## 28.1 HTTP 웹 서버 만들기

### net/http 패키지

Go의 `net/http` 패키지만으로 프로덕션 수준의 웹 서버를 만들 수 있다. 각 요청은 자동으로 별도의 고루틴에서 처리되므로 동시 요청을 효율적으로 다룬다:

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // 핸들러 함수 등록
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })

    // 서버 시작 (8080 포트)
    http.ListenAndServe(":8080", nil)
}
```

`http.ListenAndServe`의 두 번째 인자가 `nil`이면 `DefaultServeMux`를 사용한다. 이 함수는 서버가 종료될 때까지 블로킹되므로 일반적으로 `main()` 함수의 마지막에 위치한다. 에러가 발생하면 반환하므로 `log.Fatal(http.ListenAndServe(":8080", nil))`과 같이 에러 처리를 하는 것이 좋다.

### http.HandleFunc

`http.HandleFunc`는 URL 경로와 핸들러 함수를 연결한다:

```go
// 핸들러 함수의 시그니처
func handler(w http.ResponseWriter, r *http.Request) {
    // w: 응답을 작성하는 Writer
    // r: 요청 정보를 담고 있는 Request
}
```

`http.ResponseWriter`는 인터페이스로, `Header()`, `Write()`, `WriteHeader()` 세 가지 메서드를 제공한다. `Write()`를 호출하면 암묵적으로 `WriteHeader(http.StatusOK)`가 먼저 호출되므로, 상태 코드를 변경하려면 `Write()` 이전에 `WriteHeader()`를 호출해야 한다.

### http.Request의 주요 필드

| 필드 | 설명 |
|------|------|
| `r.Method` | HTTP 메서드이다 (GET, POST 등) |
| `r.URL` | 요청 URL 전체 정보이다 |
| `r.URL.Path` | URL 경로이다 |
| `r.URL.Query()` | 쿼리 파라미터를 `url.Values` 타입으로 반환한다 |
| `r.Header` | 요청 헤더이다 |
| `r.Body` | 요청 본문이다 (`io.ReadCloser` 타입) |
| `r.RemoteAddr` | 클라이언트 주소이다 |
| `r.ContentLength` | 요청 본문의 크기이다 |
| `r.Host` | 요청의 Host 헤더 값이다 |
| `r.Form` | `ParseForm()` 호출 후 사용 가능한 폼 데이터이다 |

---

## 28.2 HTTP 동작 원리

### Request와 Response

HTTP 통신은 클라이언트의 **Request**와 서버의 **Response**로 이루어진다. HTTP는 무상태(stateless) 프로토콜이므로 각 요청은 독립적이다:

```
클라이언트                    서버
   │                           │
   │   ─── HTTP Request ──►    │
   │   GET /hello HTTP/1.1     │
   │   Host: localhost:8080    │
   │                           │
   │   ◄── HTTP Response ───   │
   │   HTTP/1.1 200 OK        │
   │   Content-Type: text/html │
   │   <body>                  │
   │                           │
```

HTTP/1.1에서는 기본적으로 `Keep-Alive` 연결을 사용하여 하나의 TCP 연결로 여러 요청을 처리한다. Go의 `net/http`는 이를 자동으로 관리한다.

### HTTP 메서드

| 메서드 | 설명 | 안전 | 멱등 |
|--------|------|------|------|
| GET | 리소스 조회 | O | O |
| POST | 리소스 생성 | X | X |
| PUT | 리소스 전체 수정 | X | O |
| PATCH | 리소스 부분 수정 | X | X |
| DELETE | 리소스 삭제 | X | O |
| HEAD | GET과 동일하나 본문 없이 헤더만 반환 | O | O |
| OPTIONS | 서버가 지원하는 메서드 목록 반환 | O | O |

**안전(Safe)** 메서드는 서버의 상태를 변경하지 않는다. **멱등(Idempotent)** 메서드는 같은 요청을 여러 번 보내도 결과가 동일하다. 예를 들어 `DELETE /users/1`을 두 번 보내면 첫 번째는 삭제, 두 번째는 이미 없으므로 404이지만, 서버 상태는 동일하다.

### HTTP 상태 코드

| 코드 | 상수 | 설명 |
|------|------|------|
| 200 | `http.StatusOK` | 성공이다 |
| 201 | `http.StatusCreated` | 생성에 성공했다 |
| 204 | `http.StatusNoContent` | 성공했으나 반환할 본문이 없다 |
| 301 | `http.StatusMovedPermanently` | 영구적으로 리다이렉트한다 |
| 400 | `http.StatusBadRequest` | 잘못된 요청이다 |
| 401 | `http.StatusUnauthorized` | 인증이 필요하다 |
| 403 | `http.StatusForbidden` | 권한이 없다 |
| 404 | `http.StatusNotFound` | 리소스를 찾을 수 없다 |
| 405 | `http.StatusMethodNotAllowed` | 허용되지 않은 메서드이다 |
| 500 | `http.StatusInternalServerError` | 서버 내부 오류이다 |

상태 코드는 크게 다섯 범주로 나뉜다: 1xx(정보), 2xx(성공), 3xx(리다이렉션), 4xx(클라이언트 오류), 5xx(서버 오류).

---

## 28.3 HTTP 쿼리 인수 사용하기

URL에 `?key=value` 형태로 전달되는 쿼리 파라미터를 처리한다:

```go
func searchHandler(w http.ResponseWriter, r *http.Request) {
    // 쿼리 파라미터 가져오기
    query := r.URL.Query()

    keyword := query.Get("keyword")  // 단일 값
    page := query.Get("page")        // 단일 값

    // 여러 값을 가진 파라미터
    tags := query["tags"]            // []string

    fmt.Fprintf(w, "검색어: %s, 페이지: %s, 태그: %v",
        keyword, page, tags)
}
```

요청 예시: `GET /search?keyword=golang&page=1&tags=web&tags=server`

`query.Get()`은 키가 없으면 빈 문자열을 반환한다. 키의 존재 여부를 명시적으로 확인하려면 `query.Has("keyword")`를 사용하거나, `query["keyword"]`로 접근하여 `nil` 여부를 확인한다. 쿼리 파라미터 값은 항상 문자열이므로, 숫자로 사용하려면 `strconv.Atoi()` 등으로 변환해야 한다.

---

## 28.4 ServeMux 인스턴스 이용하기

기본 `http.DefaultServeMux` 대신 직접 `ServeMux`를 생성하면 더 안전하고 유연하다:

```go
func main() {
    // 새 ServeMux 생성
    mux := http.NewServeMux()

    // 핸들러 등록
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/hello", helloHandler)
    mux.HandleFunc("/api/users", usersHandler)

    // 서버 설정
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1MB
    }

    log.Fatal(server.ListenAndServe())
}
```

### ServeMux를 직접 만드는 이유

1. **보안**: 외부 라이브러리가 `DefaultServeMux`에 핸들러를 몰래 추가할 수 없다. 예를 들어 `net/http/pprof`를 import하면 `DefaultServeMux`에 디버그 핸들러가 등록되는데, 직접 만든 `ServeMux`에는 이런 일이 발생하지 않는다.
2. **테스트**: 테스트에서 별도의 `ServeMux`를 생성하여 격리된 테스트가 가능하다.
3. **서버 설정**: `http.Server` 구조체로 타임아웃 등 세밀한 설정이 가능하다.

### http.Server의 타임아웃 설정

타임아웃을 설정하지 않으면 슬로우로리스(Slowloris) 공격에 취약해질 수 있다. `ReadTimeout`은 요청을 읽는 최대 시간, `WriteTimeout`은 응답을 작성하는 최대 시간, `IdleTimeout`은 Keep-Alive 연결의 유휴 시간이다. 프로덕션 환경에서는 반드시 이 값들을 설정해야 한다.

---

## 28.5 파일 서버

정적 파일(HTML, CSS, JS, 이미지 등)을 제공하는 파일 서버를 만들 수 있다:

```go
// 특정 디렉토리의 파일을 제공
fs := http.FileServer(http.Dir("./static"))
mux.Handle("/static/", http.StripPrefix("/static/", fs))
```

`http.StripPrefix`는 URL 경로에서 접두사를 제거한다:
- 요청: `/static/css/style.css`
- `StripPrefix` 적용 후: `/css/style.css`
- 실제 파일: `./static/css/style.css`

`http.Dir`은 디렉토리 목록을 기본으로 노출하므로, 보안이 필요한 경우 디렉토리 리스팅을 비활성화하는 커스텀 파일 시스템을 구현해야 한다. `http.FS()`를 사용하면 `embed.FS`와 같은 `fs.FS` 인터페이스를 구현한 파일 시스템도 제공할 수 있다.

---

## 28.6 웹 서버 테스트 코드 만들기

### httptest 패키지

`net/http/httptest` 패키지를 사용하면 실제 서버를 시작하지 않고도 핸들러를 테스트할 수 있다. 네트워크 통신 없이 메모리 내에서 요청과 응답을 처리하므로 빠르고 안정적인 테스트가 가능하다:

```go
func TestHelloHandler(t *testing.T) {
    // 가짜 요청 생성
    req := httptest.NewRequest("GET", "/hello?name=Go", nil)

    // 가짜 응답 기록기 생성
    w := httptest.NewRecorder()

    // 핸들러 직접 호출
    helloHandler(w, req)

    // 결과 확인
    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("상태 코드 = %d; 기대값 200", resp.StatusCode)
    }

    body, _ := io.ReadAll(resp.Body)
    if !strings.Contains(string(body), "Go") {
        t.Errorf("응답에 'Go'가 포함되어야 한다: %s", body)
    }
}
```

`httptest.NewRequest`는 테스트용 `*http.Request`를 생성한다. POST 요청을 테스트할 때는 세 번째 인자로 요청 본문(`io.Reader`)을 전달한다. `httptest.NewRecorder`는 `http.ResponseWriter`를 구현하여 응답을 메모리에 기록한다.

### 테스트 서버

전체 서버를 테스트하려면 `httptest.NewServer`를 사용한다. 실제 TCP 포트를 열고 HTTP 통신을 수행하므로 통합 테스트에 적합하다:

```go
func TestServer(t *testing.T) {
    // 테스트용 서버 시작
    ts := httptest.NewServer(http.HandlerFunc(helloHandler))
    defer ts.Close()

    // 실제 HTTP 요청 보내기
    resp, err := http.Get(ts.URL + "/hello?name=Test")
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    // 응답 확인
    body, _ := io.ReadAll(resp.Body)
    t.Log("응답:", string(body))
}
```

`httptest.NewTLSServer`를 사용하면 HTTPS 테스트 서버도 생성할 수 있다.

---

## 28.7 JSON 데이터 전송

### JSON 응답 보내기

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    user := User{ID: 1, Name: "홍길동", Age: 30}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

`json.NewEncoder(w).Encode(user)`는 구조체를 JSON으로 직렬화하여 `ResponseWriter`에 직접 쓴다. `json.Marshal`을 사용하는 것보다 효율적인데, 중간 바이트 슬라이스를 생성하지 않고 스트리밍 방식으로 처리하기 때문이다.

### JSON 요청 받기

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "잘못된 JSON 형식이다", http.StatusBadRequest)
        return
    }

    // user 처리...
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

`json.NewDecoder`는 스트리밍 방식으로 JSON을 파싱하므로 큰 요청 본문을 다룰 때도 효율적이다. 단, 요청 본문의 크기를 제한하려면 `http.MaxBytesReader`를 사용하여 DoS 공격을 방지해야 한다:

```go
r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB 제한
```

### JSON 구조체 태그

| 태그 | 설명 |
|------|------|
| `` `json:"name"` `` | JSON 키 이름을 지정한다 |
| `` `json:"name,omitempty"` `` | 값이 zero value이면 JSON에서 제외한다 |
| `` `json:"-"` `` | JSON 직렬화/역직렬화에서 제외한다 |
| `` `json:",string"` `` | 숫자를 JSON 문자열로 인코딩한다 |

---

## 28.8 HTTPS 웹 서버 만들기

HTTPS는 TLS/SSL 인증서를 사용하여 통신을 암호화한다. 클라이언트와 서버 간의 모든 데이터가 암호화되므로 중간자 공격(MITM)을 방지할 수 있다:

```go
// 자체 서명 인증서 생성 (개발용)
// openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handler)

    server := &http.Server{
        Addr:    ":443",
        Handler: mux,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    }

    // HTTPS 서버 시작
    log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
```

### HTTP를 HTTPS로 리다이렉트

프로덕션 환경에서는 HTTP(80 포트)로 들어오는 요청을 HTTPS(443 포트)로 리다이렉트하는 것이 일반적이다:

```go
// HTTP(80)에서 HTTPS(443)로 리다이렉트
go http.ListenAndServe(":80", http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
        target := "https://" + r.Host + r.RequestURI
        http.Redirect(w, r, target, http.StatusMovedPermanently)
    },
))
```

`ListenAndServeTLS`는 인증서 파일(`cert.pem`)과 개인키 파일(`key.pem`)을 인자로 받는다. 프로덕션 환경에서는 Let's Encrypt 같은 무료 인증 기관에서 발급받은 인증서를 사용하는 것이 좋다. Go의 `golang.org/x/crypto/acme/autocert` 패키지를 사용하면 인증서 발급과 갱신을 자동으로 처리할 수 있다.

---

## 핵심 요약

1. `net/http` 패키지만으로 강력한 HTTP 서버를 만들 수 있다.
2. `http.HandleFunc`로 URL과 핸들러 함수를 연결한다.
3. `r.URL.Query()`로 쿼리 파라미터를 처리한다.
4. `http.NewServeMux()`로 직접 라우터를 생성하면 더 안전하다.
5. `http.FileServer`로 정적 파일을 서비스할 수 있다.
6. `httptest` 패키지로 핸들러를 쉽게 테스트할 수 있다.
7. `encoding/json`으로 JSON 요청/응답을 처리한다.
8. `ListenAndServeTLS`로 HTTPS 서버를 만들 수 있다.
9. `http.Server`의 타임아웃을 반드시 설정하여 보안을 강화해야 한다.

---

## 연습문제

### 연습문제 1: 방명록 서버
간단한 방명록 웹 서버를 만들어라:
- `GET /` - 방명록 목록을 HTML로 표시한다
- `POST /write` - 새 방명록을 등록한다
- 데이터는 메모리에 저장한다 (슬라이스 사용)

### 연습문제 2: 계산기 API
쿼리 파라미터로 계산을 수행하는 API를 만들어라:
- `GET /calc?op=add&a=10&b=20` -> `{"result": 30}`
- 지원 연산: add, sub, mul, div
- 잘못된 입력에 대한 에러 처리를 구현한다

### 연습문제 3: httptest 테스트
연습문제 2의 계산기 API에 대한 테스트 코드를 작성하라:
- 정상 동작 테스트
- 잘못된 연산자 테스트
- 0으로 나누기 테스트
- 테이블 주도 테스트 패턴을 사용한다

### 연습문제 4: HTTP 메서드 구분
하나의 URL 경로(`/api/items`)에서 HTTP 메서드에 따라 다른 동작을 수행하는 핸들러를 작성하라:
- GET: 목록 반환
- POST: 새 항목 추가
- 그 외 메서드: 405 Method Not Allowed 반환
- 각 메서드에 대한 테스트를 작성한다

### 연습문제 5: 미들웨어 체인
다음 미들웨어를 구현하고 체인으로 연결하라:
- 로깅 미들웨어: 요청 메서드, 경로, 처리 시간을 출력한다
- 인증 미들웨어: `Authorization` 헤더에 특정 토큰이 있는지 확인한다
- CORS 미들웨어: CORS 관련 헤더를 설정한다

### 연습문제 6: 파일 업로드 서버
파일을 업로드하고 다운로드할 수 있는 서버를 만들어라:
- `POST /upload` - `multipart/form-data`로 파일을 업로드한다
- `GET /files/{filename}` - 업로드된 파일을 다운로드한다
- `GET /files` - 업로드된 파일 목록을 JSON으로 반환한다
- 파일 크기 제한(10MB)을 설정한다

### 연습문제 7: Graceful Shutdown
서버가 `SIGINT`(Ctrl+C) 신호를 받으면 현재 처리 중인 요청을 완료한 후 안전하게 종료하는 코드를 작성하라. `http.Server.Shutdown(ctx)`를 사용한다.

### 연습문제 8: 요청 본문 크기 제한
`http.MaxBytesReader`를 사용하여 요청 본문의 크기를 제한하는 미들웨어를 작성하라. 제한을 초과하는 요청에 대해 413 Payload Too Large 상태 코드를 반환하라.

### 연습문제 9: HTML 템플릿 렌더링
`html/template` 패키지를 사용하여 동적 HTML 페이지를 렌더링하는 서버를 만들어라:
- 메인 레이아웃 템플릿과 페이지별 콘텐츠 템플릿으로 구성한다
- 템플릿에 데이터를 전달하여 동적 콘텐츠를 생성한다
- XSS 공격 방지를 위한 자동 이스케이핑이 어떻게 동작하는지 확인한다

### 연습문제 10: 상태 코드별 에러 핸들러
커스텀 에러 페이지를 제공하는 에러 핸들러를 구현하라. 404, 405, 500 등 상태 코드별로 다른 JSON 응답 또는 HTML 페이지를 반환하도록 한다.

---

## 구현 과제

### 과제 1: URL 단축 서비스
URL 단축 서비스를 구현하라:
- `POST /shorten` - 긴 URL을 받아 짧은 코드를 생성한다
- `GET /{code}` - 짧은 코드에 해당하는 원본 URL로 리다이렉트한다
- `GET /stats/{code}` - 해당 코드의 접속 횟수를 JSON으로 반환한다
- 데이터는 메모리에 저장하고, sync.RWMutex로 동시성을 처리한다
- 모든 엔드포인트에 대한 테스트를 작성한다

### 과제 2: JSON 설정 서버
HTTP API를 통해 설정 값을 관리하는 서버를 구현하라:
- `GET /config/{key}` - 특정 설정 값을 조회한다
- `PUT /config/{key}` - 설정 값을 저장한다
- `DELETE /config/{key}` - 설정 값을 삭제한다
- `GET /config` - 모든 설정을 JSON으로 반환한다
- 설정 변경 시 파일에도 저장하여 서버 재시작 후에도 유지되도록 한다

### 과제 3: 간이 프록시 서버
HTTP 요청을 다른 서버로 전달하는 리버스 프록시를 구현하라:
- `httputil.ReverseProxy`를 사용하지 않고 직접 구현한다
- 클라이언트 요청의 헤더와 본문을 대상 서버로 전달한다
- 대상 서버의 응답을 클라이언트에게 전달한다
- 요청/응답 로깅 기능을 추가한다

### 과제 4: 정적 사이트 생성기 서버
마크다운 파일을 HTML로 변환하여 제공하는 서버를 구현하라:
- 지정된 디렉토리의 `.md` 파일을 HTML로 변환하여 서비스한다
- CSS 파일로 스타일링을 적용한다
- 파일 변경 감지 기능을 추가한다 (개발 모드)

### 과제 5: 속도 제한(Rate Limiting) 미들웨어
IP 기반 속도 제한 미들웨어를 구현하라:
- 토큰 버킷(Token Bucket) 알고리즘을 사용한다
- IP별로 초당 최대 요청 수를 제한한다
- 제한 초과 시 429 Too Many Requests 상태 코드를 반환한다
- `X-RateLimit-Limit`, `X-RateLimit-Remaining` 헤더를 추가한다

---

## 프로젝트 과제

### 프로젝트 1: 개인 블로그 서버
Go `net/http`만으로 완전한 블로그 서버를 구현하라:
- 게시글 CRUD (Create, Read, Update, Delete)를 지원한다
- 마크다운으로 작성한 게시글을 HTML로 렌더링한다
- `html/template`를 사용하여 메인 페이지, 게시글 목록, 개별 게시글 페이지를 구현한다
- 정적 파일(CSS, 이미지)을 서비스한다
- 간단한 관리자 인증(Basic Auth)을 구현한다
- JSON 파일 기반 데이터 저장소를 사용한다
- `httptest`를 사용한 통합 테스트를 작성한다

### 프로젝트 2: 실시간 대시보드
서버 상태를 실시간으로 표시하는 대시보드를 만들어라:
- Server-Sent Events(SSE)를 사용하여 실시간 데이터를 클라이언트로 전송한다
- CPU 사용률, 메모리 사용량, 고루틴 수, 요청 수 등의 지표를 표시한다
- HTML/CSS/JS로 대시보드 UI를 구현한다
- `runtime` 패키지로 Go 런타임 지표를 수집한다
- 지표 히스토리를 메모리에 저장하여 그래프로 표시한다
