# 29장 [Project] RESTful API 서버 만들기

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회, 포트 노출 필요: -p 8080:8080)
make shell

# net/http 기반 서버 실행
go run ch29_project_restful_api/main.go

# Gin 기반 서버 실행
go run ch29_project_restful_api/gin_server.go

# 테스트 실행
go test ch29_project_restful_api/ -v

# 서버 실행 후 다른 터미널에서 API 테스트
curl http://localhost:8080/api/students
curl -X POST http://localhost:8080/api/students \
  -H "Content-Type: application/json" \
  -d '{"name":"홍길동","age":20,"grade":3,"email":"hong@example.com"}'
```

> **참고**: HTTP 서버를 Docker 컨테이너 내부에서 실행할 때는 컨테이너의 포트를 호스트에 노출해야 외부에서 접근할 수 있다.

> **Makefile 활용**: `make run CH=ch29_project_restful_api` 또는 `make run CH=ch29_project_restful_api FILE=main.go`

---

이 장에서는 학생 관리 시스템을 RESTful API로 설계하고 구현한다. 먼저 `net/http`로 직접 구현한 후, Gin 프레임워크로 같은 기능을 더 간결하게 만들어 본다. REST 아키텍처는 웹 API의 사실상 표준이며, 이를 이해하고 올바르게 구현하는 것은 백엔드 개발의 핵심 역량이다.

---

## 29.1 해법: 학생 관리 API 설계

### 요구사항

학생 정보를 CRUD(생성, 조회, 수정, 삭제)할 수 있는 API를 만든다. API는 JSON 형식으로 요청과 응답을 처리하며, 적절한 HTTP 상태 코드를 반환해야 한다.

### API 설계

| 메서드 | 경로 | 설명 | 성공 코드 |
|--------|------|------|-----------|
| GET | `/api/students` | 모든 학생 목록 조회 | 200 |
| GET | `/api/students/{id}` | 특정 학생 조회 | 200 |
| POST | `/api/students` | 새 학생 등록 | 201 |
| PUT | `/api/students/{id}` | 학생 정보 수정 | 200 |
| DELETE | `/api/students/{id}` | 학생 삭제 | 204 |

### 데이터 모델

```go
type Student struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Grade int    `json:"grade"`
    Email string `json:"email"`
}
```

구조체 태그(`json:"..."`)는 JSON 직렬화/역직렬화 시 사용되는 키 이름을 지정한다. Go의 관례인 PascalCase 필드 이름을 JSON의 관례인 camelCase 또는 snake_case로 매핑할 수 있다.

### 응답 형식 표준화

일관된 API 응답을 위해 공통 응답 구조를 정의하는 것이 좋다:

```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

---

## 29.2 사전 지식: RESTful API

### REST 원칙

REST(Representational State Transfer)는 다음 원칙을 따른다:

1. **클라이언트-서버 분리**: 클라이언트와 서버가 독립적으로 발전한다. 클라이언트는 서버의 데이터 저장 방식을 알 필요가 없고, 서버는 클라이언트의 UI를 알 필요가 없다.
2. **무상태(Stateless)**: 각 요청은 독립적이며, 서버는 클라이언트의 상태를 저장하지 않는다. 요청에 필요한 모든 정보는 요청 자체에 포함되어야 한다.
3. **균일한 인터페이스**: URL로 리소스를 식별하고, HTTP 메서드로 행위를 표현한다. 모든 리소스에 대해 일관된 접근 방법을 사용한다.
4. **계층 구조**: 클라이언트는 서버의 내부 구조를 알 필요 없다. 로드 밸런서, 캐시, 프록시 등 중간 계층을 투명하게 추가할 수 있다.
5. **캐시 가능**: 서버 응답에 캐시 가능 여부를 명시하여 클라이언트가 응답을 재사용할 수 있게 한다.

### HTTP 메서드와 CRUD

| CRUD | HTTP 메서드 | 예시 | 설명 |
|------|-------------|------|------|
| Create (생성) | POST | `POST /api/students` | 요청 본문에 데이터를 포함한다 |
| Read (조회) | GET | `GET /api/students/1` | URL로 리소스를 식별한다 |
| Update (수정) | PUT | `PUT /api/students/1` | 전체 리소스를 교체한다 |
| Update (부분 수정) | PATCH | `PATCH /api/students/1` | 일부 필드만 수정한다 |
| Delete (삭제) | DELETE | `DELETE /api/students/1` | 리소스를 제거한다 |

PUT과 PATCH의 차이를 이해하는 것이 중요하다. PUT은 리소스 전체를 교체하므로 모든 필드를 보내야 한다. PATCH는 변경하고 싶은 필드만 보내면 된다. 예를 들어 학생의 이메일만 변경하려면 PATCH가 적절하다.

### 좋은 URL 설계

```
좋은 예:
  GET    /api/students       # 목록 조회
  GET    /api/students/1     # 상세 조회
  POST   /api/students       # 생성
  PUT    /api/students/1     # 수정
  DELETE /api/students/1     # 삭제

나쁜 예:
  GET    /api/getStudents
  POST   /api/createStudent
  POST   /api/deleteStudent?id=1
```

좋은 URL 설계의 핵심 원칙은 다음과 같다:
- **명사를 사용한다**: URL은 리소스(명사)를 나타내고, 행위(동사)는 HTTP 메서드로 표현한다.
- **복수형을 사용한다**: `/api/students` (O), `/api/student` (X)
- **계층 관계를 표현한다**: `/api/students/1/courses` (1번 학생의 수강 과목)
- **소문자와 하이픈을 사용한다**: `/api/course-registrations` (O), `/api/courseRegistrations` (X)

---

## 29.3 RESTful API 서버 만들기

### net/http로 구현

`net/http` 패키지만으로 RESTful API를 구현한다. 표준 라이브러리만 사용하므로 외부 의존성이 없다:

```go
func studentsHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getStudents(w, r)
    case http.MethodPost:
        createStudent(w, r)
    default:
        http.Error(w, "Method Not Allowed",
            http.StatusMethodNotAllowed)
    }
}
```

`net/http`의 `ServeMux`는 경로만 매칭하고 HTTP 메서드를 구분하지 않으므로, 핸들러 내부에서 `r.Method`를 검사하여 분기해야 한다. 이는 코드가 복잡해지는 원인이며, 웹 프레임워크를 사용하는 주된 이유 중 하나이다.

### 미들웨어 패턴

요청 전후에 공통 로직을 실행하는 미들웨어를 만들 수 있다. 미들웨어는 핸들러를 감싸는 함수로, 횡단 관심사(cross-cutting concerns)를 처리하는 데 사용한다:

```go
func jsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next(w, r)
    }
}

// 로깅 미들웨어
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    }
}

// 미들웨어 체인: 여러 미들웨어를 순서대로 적용
mux.HandleFunc("/api/students",
    loggingMiddleware(jsonMiddleware(studentsHandler)))
```

---

## 29.4 테스트 코드 작성하기

`httptest` 패키지로 API를 테스트한다. API 테스트에서는 요청 생성, 핸들러 호출, 응답 검증의 세 단계를 따른다:

```go
func TestGetStudents(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/students", nil)
    w := httptest.NewRecorder()

    handler := setupRouter()
    handler.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("상태 코드 = %d; 기대값 200", w.Code)
    }

    // JSON 응답 파싱하여 검증
    var students []Student
    if err := json.NewDecoder(w.Body).Decode(&students); err != nil {
        t.Fatalf("JSON 파싱 실패: %v", err)
    }
}

func TestCreateStudent(t *testing.T) {
    body := `{"name":"홍길동","age":20,"grade":3,"email":"hong@example.com"}`
    req := httptest.NewRequest("POST", "/api/students",
        strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    handler := setupRouter()
    handler.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Errorf("상태 코드 = %d; 기대값 201", w.Code)
    }
}
```

테스트에서 `setupRouter()`를 별도 함수로 분리하면 테스트와 메인 코드에서 동일한 라우터 구성을 재사용할 수 있다.

---

## 29.5 특정 학생 데이터 반환하기

URL 경로에서 학생 ID를 추출하여 특정 학생 데이터를 반환한다:

```go
// /api/students/1 에서 ID 추출
func extractID(path string) (int, error) {
    parts := strings.Split(path, "/")
    if len(parts) < 4 {
        return 0, errors.New("잘못된 경로이다")
    }
    return strconv.Atoi(parts[3])
}
```

이 방식은 URL 구조에 강하게 결합되어 있다는 단점이 있다. 경로가 변경되면 `extractID` 함수도 수정해야 한다. Go 1.22부터는 `ServeMux`에서 경로 변수를 지원하므로 `mux.HandleFunc("GET /api/students/{id}", handler)`와 같이 작성하고 `r.PathValue("id")`로 값을 추출할 수 있다.

---

## 29.6 학생 데이터 추가/삭제하기

### POST: 학생 추가

```go
func createStudent(w http.ResponseWriter, r *http.Request) {
    var student Student
    if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
        respondError(w, http.StatusBadRequest, "잘못된 요청이다")
        return
    }

    // 유효성 검증
    if student.Name == "" {
        respondError(w, http.StatusBadRequest, "이름은 필수이다")
        return
    }

    // 저장 로직...
    respondJSON(w, http.StatusCreated, student)
}
```

### DELETE: 학생 삭제

```go
func deleteStudent(w http.ResponseWriter, r *http.Request) {
    id, err := extractID(r.URL.Path)
    if err != nil {
        respondError(w, http.StatusBadRequest, "잘못된 ID이다")
        return
    }
    // 삭제 로직...
    w.WriteHeader(http.StatusNoContent)
}
```

### 유틸리티 함수

JSON 응답과 에러 응답을 보내는 유틸리티 함수를 만들면 코드 중복을 줄일 수 있다:

```go
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{"error": message})
}
```

---

## 29.7 RESTful API로의 발전

표준 `net/http`로 만든 API의 한계는 다음과 같다:

1. **URL 경로 파라미터 추출이 번거롭다**: `/students/{id}`에서 ID를 추출하려면 문자열을 수동으로 파싱해야 한다.
2. **메서드별 라우팅을 수동으로 해야 한다**: 핸들러 내에서 `switch r.Method`로 분기해야 한다.
3. **미들웨어 체인 관리가 복잡하다**: 미들웨어를 여러 개 적용하려면 중첩 호출이 깊어진다.
4. **입력 유효성 검증을 직접 구현해야 한다**: 필드별 검증 로직을 일일이 작성해야 한다.
5. **그룹 라우팅이 불가능하다**: `/api/v1/...`과 같은 접두사를 공유하는 경로 그룹을 만들기 어렵다.

이러한 한계를 극복하기 위해 웹 프레임워크를 사용할 수 있다. 단, Go 1.22 이후에는 표준 `ServeMux`의 기능이 크게 향상되어 간단한 API는 프레임워크 없이도 충분히 구현할 수 있다.

---

## 29.8 Gin으로 서버 만들기

### Gin 프레임워크

[Gin](https://github.com/gin-gonic/gin)은 Go에서 가장 인기 있는 웹 프레임워크이다. `httprouter` 기반의 빠른 라우팅을 제공하며, 미들웨어, JSON 바인딩, 유효성 검증 등 API 개발에 필요한 기능을 포괄적으로 지원한다:

```bash
go get -u github.com/gin-gonic/gin
```

### Gin의 장점

1. **경로 파라미터**: `/students/:id`로 간편하게 파라미터를 추출한다
2. **메서드별 라우팅**: `r.GET()`, `r.POST()` 등 직관적인 API를 제공한다
3. **미들웨어**: `r.Use()`로 간편하게 미들웨어를 적용한다
4. **바인딩**: `c.ShouldBindJSON()`으로 자동 JSON 바인딩을 수행한다
5. **검증**: 구조체 태그로 자동 유효성 검사를 수행한다
6. **그룹 라우팅**: `r.Group()`으로 경로 그룹을 만들 수 있다

### Gin 예제

```go
func main() {
    r := gin.Default() // Logger + Recovery 미들웨어 포함

    // 라우트 그룹
    api := r.Group("/api")
    {
        api.GET("/students", getStudents)
        api.GET("/students/:id", getStudent)
        api.POST("/students", createStudent)
        api.PUT("/students/:id", updateStudent)
        api.DELETE("/students/:id", deleteStudent)
    }

    r.Run(":8080")
}

func getStudent(c *gin.Context) {
    id := c.Param("id")  // 경로 파라미터 추출
    // ...
    c.JSON(http.StatusOK, student)
}

func createStudent(c *gin.Context) {
    var student Student
    if err := c.ShouldBindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // ...
    c.JSON(http.StatusCreated, student)
}
```

### Gin의 유효성 검증

Gin은 `go-playground/validator`를 내장하고 있어 구조체 태그로 유효성 검증을 수행한다:

```go
type Student struct {
    ID    int    `json:"id"`
    Name  string `json:"name" binding:"required,min=2,max=50"`
    Age   int    `json:"age" binding:"required,gte=1,lte=150"`
    Grade int    `json:"grade" binding:"required,gte=1,lte=6"`
    Email string `json:"email" binding:"required,email"`
}
```

`binding` 태그를 사용하면 `ShouldBindJSON()` 호출 시 자동으로 유효성을 검사하며, 검증에 실패하면 상세한 에러 메시지를 반환한다.

---

## 핵심 요약

1. RESTful API는 URL로 리소스를 식별하고 HTTP 메서드로 행위를 표현한다.
2. `net/http`만으로도 RESTful API를 구현할 수 있지만, 경로 파라미터 처리가 번거롭다.
3. `httptest` 패키지로 API 핸들러를 쉽게 테스트할 수 있다.
4. Gin 프레임워크를 사용하면 경로 파라미터, 메서드 라우팅, JSON 바인딩 등을 간편하게 처리할 수 있다.
5. 미들웨어 패턴으로 로깅, 인증, CORS 등 공통 로직을 분리할 수 있다.
6. PUT은 전체 리소스 교체, PATCH는 부분 수정에 사용한다.
7. 일관된 응답 형식과 적절한 상태 코드를 사용하는 것이 중요하다.

---

## 연습문제

### 연습문제 1: 도서 관리 API
도서 관리 RESTful API를 설계하고 구현하라:
- `Book` 구조체: ID, Title, Author, ISBN, Price
- 전체 CRUD 엔드포인트를 구현한다
- 제목으로 검색 기능을 추가한다: `GET /api/books?title=Go`

### 연습문제 2: 인증 미들웨어
간단한 API 키 인증 미들웨어를 구현하라:
- `Authorization: Bearer <api-key>` 헤더를 확인한다
- 유효하지 않은 키: 401 Unauthorized 응답을 반환한다
- 미들웨어를 특정 경로에만 적용한다

### 연습문제 3: 페이징 처리
학생 목록 API에 페이징을 추가하라:
- `GET /api/students?page=1&size=10`
- 응답에 총 개수, 총 페이지 수, 현재 페이지를 포함한다

### 연습문제 4: 입력 유효성 검증
`net/http` 기반 API에 입력 유효성 검증을 추가하라:
- 이름: 필수, 2~50자
- 나이: 1~150 범위
- 이메일: 이메일 형식 확인
- 유효성 검증 실패 시 어떤 필드가 잘못되었는지 상세한 에러 메시지를 반환한다

### 연습문제 5: PATCH 메서드 구현
학생 정보의 부분 수정(PATCH)을 구현하라. PUT과 달리 전송된 필드만 업데이트해야 한다. JSON에서 전송되지 않은 필드와 zero value를 구분하기 위해 포인터 타입을 사용하는 방법을 적용하라.

### 연습문제 6: 정렬과 필터링
학생 목록 API에 정렬과 필터링 기능을 추가하라:
- 정렬: `GET /api/students?sort=name&order=asc`
- 필터: `GET /api/students?grade=3&age_min=18&age_max=25`
- 여러 조건을 조합할 수 있어야 한다

### 연습문제 7: 에러 응답 표준화
일관된 에러 응답 형식을 설계하고 구현하라:
- 에러 코드, 메시지, 상세 정보를 포함한다
- 유효성 검증 에러는 필드별 에러 목록을 반환한다
- 404, 400, 500 등 상태 코드별로 적절한 에러 응답을 반환한다

### 연습문제 8: API 버전 관리
API 버전 관리를 구현하라:
- URL 경로 방식: `/api/v1/students`, `/api/v2/students`
- v1과 v2에서 서로 다른 응답 형식을 제공한다
- Gin의 그룹 라우팅을 활용한다

### 연습문제 9: CORS 미들웨어
Cross-Origin Resource Sharing(CORS)을 처리하는 미들웨어를 직접 구현하라:
- `Access-Control-Allow-Origin` 헤더를 설정한다
- `OPTIONS` 요청(Preflight)에 적절히 응답한다
- 허용할 출처, 메서드, 헤더를 설정 가능하게 한다

### 연습문제 10: API 문서화
학생 관리 API의 각 엔드포인트에 대해 요청/응답 예시를 포함한 테이블 주도 테스트를 작성하라. 테스트 코드 자체가 API 사용 문서 역할을 하도록 상세한 테스트 이름과 주석을 작성한다.

---

## 구현 과제

### 과제 1: 할일(Todo) API 서버
할일 관리 RESTful API를 완전하게 구현하라:
- `Todo` 모델: ID, Title, Description, Completed, CreatedAt, UpdatedAt
- CRUD 엔드포인트 전체 구현
- 완료/미완료 필터링, 생성일 기준 정렬
- 페이징 지원
- 모든 엔드포인트에 대한 테스트 작성
- `net/http`와 Gin 두 가지 버전으로 구현하여 코드량과 가독성을 비교한다

### 과제 2: 사용자 인증 시스템
JWT(JSON Web Token) 기반 인증 시스템을 구현하라:
- `POST /auth/register` - 사용자 등록 (비밀번호 해시화)
- `POST /auth/login` - 로그인 후 JWT 토큰 발급
- 인증이 필요한 엔드포인트에 JWT 미들웨어 적용
- 토큰 만료 처리
- `golang.org/x/crypto/bcrypt` 패키지로 비밀번호를 해시화한다

### 과제 3: 관계형 리소스 API
학생과 과목(Course)의 다대다 관계를 가진 API를 설계하고 구현하라:
- `GET /api/students/:id/courses` - 특정 학생의 수강 과목 목록
- `POST /api/students/:id/courses` - 수강 신청
- `DELETE /api/students/:id/courses/:courseId` - 수강 취소
- 데이터는 메모리에 맵과 슬라이스로 저장한다

### 과제 4: API 게이트웨이
여러 마이크로서비스를 하나의 진입점으로 통합하는 간이 API 게이트웨이를 구현하라:
- `/api/users/*` -> 사용자 서비스로 프록시
- `/api/products/*` -> 상품 서비스로 프록시
- 공통 미들웨어(로깅, 인증, 속도 제한)를 적용한다
- 대상 서비스의 헬스 체크를 수행한다

### 과제 5: 벌크 작업 API
여러 학생을 한 번에 생성/수정/삭제하는 벌크 API를 구현하라:
- `POST /api/students/bulk` - 여러 학생을 한 번에 생성한다
- `PUT /api/students/bulk` - 여러 학생을 한 번에 수정한다
- `DELETE /api/students/bulk` - 여러 학생을 한 번에 삭제한다
- 각 항목별 성공/실패 결과를 상세히 반환한다

---

## 프로젝트 과제

### 프로젝트 1: 온라인 서점 API
온라인 서점의 백엔드 API를 완전히 설계하고 구현하라:
- **도서 관리**: CRUD, 카테고리별 조회, 검색, 페이징
- **사용자 관리**: 회원가입, 로그인, 프로필 수정
- **장바구니**: 도서 추가/삭제, 수량 변경
- **주문**: 장바구니에서 주문 생성, 주문 내역 조회
- Gin 프레임워크를 사용하고, 인증 미들웨어를 적용한다
- 모든 엔드포인트에 대한 테스트를 작성한다
- 데이터는 메모리(맵/슬라이스)에 저장한다

### 프로젝트 2: REST API 테스트 프레임워크
RESTful API를 쉽게 테스트할 수 있는 테스트 유틸리티 라이브러리를 만들어라:
- 요청을 빌더 패턴으로 생성한다: `NewRequest().GET("/api/students").WithHeader("Authorization", token).Do()`
- 응답 검증을 메서드 체인으로 수행한다: `.ExpectStatus(200).ExpectJSON("name", "홍길동")`
- 테이블 주도 테스트에 쉽게 통합할 수 있어야 한다
- 이 라이브러리를 사용하여 실제 API 테스트를 작성하고, 코드의 가독성이 얼마나 개선되는지 비교한다
