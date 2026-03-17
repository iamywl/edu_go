// 28장: HTTP 웹 서버 만들기 - JSON 응답 핸들러 예제
// 실행: go run json_handler.go
// 테스트:
//
//	curl http://localhost:8082/api/users
//	curl http://localhost:8082/api/users/1
//	curl -X POST -H "Content-Type: application/json" \
//	     -d '{"name":"홍길동","email":"hong@test.com","age":30}' \
//	     http://localhost:8082/api/users
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// User - 사용자 정보 구조체
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// ErrorResponse - 에러 응답 구조체
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// 메모리 기반 사용자 저장소
var (
	users = []User{
		{ID: 1, Name: "김철수", Email: "kim@test.com", Age: 25},
		{ID: 2, Name: "이영희", Email: "lee@test.com", Age: 30},
		{ID: 3, Name: "박민수", Email: "park@test.com", Age: 28},
	}
	nextID  = 4
	usersMu sync.RWMutex
)

// respondJSON - JSON 응답을 보내는 헬퍼 함수
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON 인코딩 에러: %v", err)
	}
}

// respondError - 에러 JSON 응답을 보내는 헬퍼 함수
func respondError(w http.ResponseWriter, status int, errMsg, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   errMsg,
		Message: message,
	})
}

// usersHandler - 사용자 목록 조회 및 생성 핸들러
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 모든 사용자 조회
		usersMu.RLock()
		defer usersMu.RUnlock()
		respondJSON(w, http.StatusOK, users)

	case http.MethodPost:
		// 새 사용자 생성
		var newUser User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			respondError(w, http.StatusBadRequest,
				"invalid_json", "JSON 형식이 올바르지 않습니다")
			return
		}

		// 유효성 검사
		if newUser.Name == "" || newUser.Email == "" {
			respondError(w, http.StatusBadRequest,
				"validation_error", "이름과 이메일은 필수입니다")
			return
		}

		usersMu.Lock()
		newUser.ID = nextID
		nextID++
		users = append(users, newUser)
		usersMu.Unlock()

		respondJSON(w, http.StatusCreated, newUser)

	default:
		respondError(w, http.StatusMethodNotAllowed,
			"method_not_allowed", "GET 또는 POST만 허용됩니다")
	}
}

// userByIDHandler - 특정 사용자 조회 핸들러
func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	// URL에서 ID 추출: /api/users/1
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		respondError(w, http.StatusBadRequest,
			"invalid_id", "사용자 ID가 필요합니다")
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		respondError(w, http.StatusBadRequest,
			"invalid_id", "유효한 숫자 ID가 필요합니다")
		return
	}

	usersMu.RLock()
	defer usersMu.RUnlock()

	// 사용자 검색
	for _, user := range users {
		if user.ID == id {
			respondJSON(w, http.StatusOK, user)
			return
		}
	}

	respondError(w, http.StatusNotFound,
		"not_found", fmt.Sprintf("ID %d인 사용자를 찾을 수 없습니다", id))
}

// apiRouter - /api/users 경로를 분기하는 라우터
func apiRouter(w http.ResponseWriter, r *http.Request) {
	// /api/users 와 /api/users/{id} 구분
	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == "/api/users" {
		usersHandler(w, r)
	} else if strings.HasPrefix(path, "/api/users/") {
		userByIDHandler(w, r)
	} else {
		respondError(w, http.StatusNotFound,
			"not_found", "요청한 경로를 찾을 수 없습니다")
	}
}

func main() {
	mux := http.NewServeMux()

	// API 라우터 등록
	mux.HandleFunc("/api/users", apiRouter)
	mux.HandleFunc("/api/users/", apiRouter)

	// 인덱스 페이지
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<h1>JSON API 서버</h1>
<p>API 엔드포인트:</p>
<ul>
<li>GET <a href="/api/users">/api/users</a> - 사용자 목록</li>
<li>GET <a href="/api/users/1">/api/users/1</a> - 특정 사용자</li>
<li>POST /api/users - 사용자 생성 (JSON body 필요)</li>
</ul>`)
	})

	addr := ":8082"
	fmt.Printf("JSON API 서버 시작: http://localhost%s\n", addr)
	fmt.Println("사용 예시:")
	fmt.Println("  curl http://localhost:8082/api/users")
	fmt.Println("  curl http://localhost:8082/api/users/1")
	fmt.Println(`  curl -X POST -H "Content-Type: application/json" \`)
	fmt.Println(`       -d '{"name":"홍길동","email":"hong@test.com","age":30}' \`)
	fmt.Println("       http://localhost:8082/api/users")
	log.Fatal(http.ListenAndServe(addr, mux))
}
