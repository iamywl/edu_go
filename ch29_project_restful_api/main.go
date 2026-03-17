// 29장: RESTful API 서버 만들기 - net/http를 이용한 학생 관리 API
// 실행: go run main.go
// 테스트:
//
//	curl http://localhost:8090/api/students
//	curl http://localhost:8090/api/students/1
//	curl -X POST -H "Content-Type: application/json" \
//	     -d '{"name":"홍길동","age":20,"grade":3,"email":"hong@test.com"}' \
//	     http://localhost:8090/api/students
//	curl -X PUT -H "Content-Type: application/json" \
//	     -d '{"name":"홍길동","age":21,"grade":4,"email":"hong@test.com"}' \
//	     http://localhost:8090/api/students/1
//	curl -X DELETE http://localhost:8090/api/students/1
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Student - 학생 정보 구조체
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade int    `json:"grade"` // 학년
	Email string `json:"email"`
}

// APIResponse - 표준 API 응답 구조체
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// StudentStore - 학생 데이터 저장소 (메모리 기반)
type StudentStore struct {
	mu       sync.RWMutex
	students map[int]Student
	nextID   int
}

// NewStudentStore - 새 저장소 생성 (초기 데이터 포함)
func NewStudentStore() *StudentStore {
	store := &StudentStore{
		students: make(map[int]Student),
		nextID:   1,
	}

	// 초기 데이터 추가
	store.Create(Student{Name: "김철수", Age: 20, Grade: 2, Email: "kim@school.com"})
	store.Create(Student{Name: "이영희", Age: 21, Grade: 3, Email: "lee@school.com"})
	store.Create(Student{Name: "박민수", Age: 19, Grade: 1, Email: "park@school.com"})

	return store
}

// GetAll - 모든 학생 조회
func (s *StudentStore) GetAll() []Student {
	s.mu.RLock()
	defer s.mu.RUnlock()

	students := make([]Student, 0, len(s.students))
	for _, student := range s.students {
		students = append(students, student)
	}
	return students
}

// GetByID - ID로 학생 조회
func (s *StudentStore) GetByID(id int) (Student, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	student, ok := s.students[id]
	return student, ok
}

// Create - 새 학생 등록
func (s *StudentStore) Create(student Student) Student {
	s.mu.Lock()
	defer s.mu.Unlock()

	student.ID = s.nextID
	s.nextID++
	s.students[student.ID] = student
	return student
}

// Update - 학생 정보 수정
func (s *StudentStore) Update(id int, student Student) (Student, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.students[id]; !ok {
		return Student{}, false
	}

	student.ID = id
	s.students[id] = student
	return student, true
}

// Delete - 학생 삭제
func (s *StudentStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.students[id]; !ok {
		return false
	}

	delete(s.students, id)
	return true
}

// === 핸들러 함수들 ===

var store = NewStudentStore()

// respondJSON - JSON 응답 헬퍼
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondSuccess - 성공 응답 헬퍼
func respondSuccess(w http.ResponseWriter, status int, data interface{}) {
	respondJSON(w, status, APIResponse{
		Success: true,
		Data:    data,
	})
}

// respondError - 에러 응답 헬퍼
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, APIResponse{
		Success: false,
		Error:   http.StatusText(status),
		Message: message,
	})
}

// extractID - URL 경로에서 ID를 추출한다
func extractID(path string) (int, error) {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) < 4 {
		return 0, fmt.Errorf("ID가 필요합니다")
	}
	return strconv.Atoi(parts[3])
}

// studentsHandler - /api/students 핸들러 (목록 조회, 생성)
func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 모든 학생 목록 조회
		students := store.GetAll()
		respondSuccess(w, http.StatusOK, students)

	case http.MethodPost:
		// 새 학생 등록
		var student Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			respondError(w, http.StatusBadRequest, "잘못된 JSON 형식입니다")
			return
		}

		// 유효성 검사
		if student.Name == "" {
			respondError(w, http.StatusBadRequest, "이름은 필수입니다")
			return
		}
		if student.Age < 1 || student.Age > 100 {
			respondError(w, http.StatusBadRequest, "나이는 1~100 사이여야 합니다")
			return
		}

		created := store.Create(student)
		respondSuccess(w, http.StatusCreated, created)

	default:
		respondError(w, http.StatusMethodNotAllowed, "GET 또는 POST만 허용됩니다")
	}
}

// studentByIDHandler - /api/students/{id} 핸들러 (조회, 수정, 삭제)
func studentByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "유효한 학생 ID가 필요합니다")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// 특정 학생 조회
		student, ok := store.GetByID(id)
		if !ok {
			respondError(w, http.StatusNotFound,
				fmt.Sprintf("ID %d인 학생을 찾을 수 없습니다", id))
			return
		}
		respondSuccess(w, http.StatusOK, student)

	case http.MethodPut:
		// 학생 정보 수정
		var student Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			respondError(w, http.StatusBadRequest, "잘못된 JSON 형식입니다")
			return
		}

		updated, ok := store.Update(id, student)
		if !ok {
			respondError(w, http.StatusNotFound,
				fmt.Sprintf("ID %d인 학생을 찾을 수 없습니다", id))
			return
		}
		respondSuccess(w, http.StatusOK, updated)

	case http.MethodDelete:
		// 학생 삭제
		if !store.Delete(id) {
			respondError(w, http.StatusNotFound,
				fmt.Sprintf("ID %d인 학생을 찾을 수 없습니다", id))
			return
		}
		respondSuccess(w, http.StatusOK, map[string]string{
			"message": fmt.Sprintf("ID %d 학생이 삭제되었습니다", id),
		})

	default:
		respondError(w, http.StatusMethodNotAllowed,
			"GET, PUT 또는 DELETE만 허용됩니다")
	}
}

// loggingMiddleware - 요청 로깅 미들웨어
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s (%v)", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// corsMiddleware - CORS 미들웨어
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// OPTIONS 요청 처리 (프리플라이트)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// setupRouter - 라우터 설정 (테스트에서도 사용)
func setupRouter() http.Handler {
	mux := http.NewServeMux()

	// API 핸들러 등록
	mux.HandleFunc("/api/students", studentsHandler)
	mux.HandleFunc("/api/students/", studentByIDHandler)

	// 인덱스 페이지
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<h1>학생 관리 RESTful API</h1>
<pre>
GET    /api/students      - 모든 학생 조회
GET    /api/students/{id} - 특정 학생 조회
POST   /api/students      - 새 학생 등록
PUT    /api/students/{id} - 학생 정보 수정
DELETE /api/students/{id} - 학생 삭제
</pre>`)
	})

	// 미들웨어 적용
	handler := loggingMiddleware(corsMiddleware(mux))
	return handler
}

func main() {
	handler := setupRouter()

	addr := ":8090"
	fmt.Printf("학생 관리 API 서버 시작: http://localhost%s\n", addr)
	fmt.Println("\n사용 예시:")
	fmt.Println("  curl http://localhost:8090/api/students")
	fmt.Println("  curl http://localhost:8090/api/students/1")
	fmt.Println(`  curl -X POST -H "Content-Type: application/json" \`)
	fmt.Println(`       -d '{"name":"홍길동","age":20,"grade":3,"email":"hong@test.com"}' \`)
	fmt.Println("       http://localhost:8090/api/students")
	fmt.Println("  curl -X DELETE http://localhost:8090/api/students/1")

	log.Fatal(http.ListenAndServe(addr, handler))
}
