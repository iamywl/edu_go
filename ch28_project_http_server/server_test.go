// 28장: HTTP 웹 서버 만들기 - httptest를 이용한 테스트 예제
// 실행: go test -v
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// === 핸들러 함수 (테스트 대상) ===

// testHelloHandler - 테스트용 인사 핸들러
func testHelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, " + name + "!"))
}

// testUserListHandler - 테스트용 사용자 목록 핸들러
func testUserListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	testUsers := []map[string]interface{}{
		{"id": 1, "name": "김철수"},
		{"id": 2, "name": "이영희"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testUsers)
}

// testCreateHandler - 테스트용 사용자 생성 핸들러
func testCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if body["name"] == nil || body["name"] == "" {
		http.Error(w, "이름은 필수입니다", http.StatusBadRequest)
		return
	}

	body["id"] = 100
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

// === 테스트 코드 ===

// TestHelloHandlerDefault - 기본 인사 핸들러 테스트
func TestHelloHandlerDefault(t *testing.T) {
	// httptest.NewRequest로 가짜 요청을 생성한다
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	// httptest.NewRecorder로 가짜 응답 기록기를 생성한다
	w := httptest.NewRecorder()

	// 핸들러를 직접 호출한다
	testHelloHandler(w, req)

	// 응답 결과를 확인한다
	resp := w.Result()
	defer resp.Body.Close()

	// 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		t.Errorf("상태 코드 = %d; 기대값 %d", resp.StatusCode, http.StatusOK)
	}

	// 응답 본문 확인
	body, _ := io.ReadAll(resp.Body)
	expected := "Hello, World!"
	if string(body) != expected {
		t.Errorf("응답 = %q; 기대값 %q", string(body), expected)
	}
}

// TestHelloHandlerWithName - 이름 파라미터가 있는 인사 핸들러 테스트
func TestHelloHandlerWithName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?name=Go", nil)
	w := httptest.NewRecorder()

	testHelloHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	expected := "Hello, Go!"
	if string(body) != expected {
		t.Errorf("응답 = %q; 기대값 %q", string(body), expected)
	}
}

// TestHelloHandlerTableDriven - 테이블 주도 테스트
func TestHelloHandlerTableDriven(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantBody   string
	}{
		{"이름 없음", "/hello", 200, "Hello, World!"},
		{"이름 있음", "/hello?name=Go", 200, "Hello, Go!"},
		{"한글 이름", "/hello?name=고퍼", 200, "Hello, 고퍼!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()

			testHelloHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("상태 코드 = %d; 기대값 %d", resp.StatusCode, tt.wantStatus)
			}

			body, _ := io.ReadAll(resp.Body)
			if string(body) != tt.wantBody {
				t.Errorf("응답 = %q; 기대값 %q", string(body), tt.wantBody)
			}
		})
	}
}

// TestUserListHandler - 사용자 목록 핸들러 테스트
func TestUserListHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	testUserListHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("상태 코드 = %d; 기대값 %d", resp.StatusCode, http.StatusOK)
	}

	// Content-Type 확인
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q; 기대값 %q", contentType, "application/json")
	}

	// JSON 파싱 확인
	var users []map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &users); err != nil {
		t.Fatalf("JSON 파싱 실패: %v", err)
	}

	// 사용자 수 확인
	if len(users) != 2 {
		t.Errorf("사용자 수 = %d; 기대값 2", len(users))
	}
}

// TestUserListHandlerMethodNotAllowed - 잘못된 메서드 테스트
func TestUserListHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	w := httptest.NewRecorder()

	testUserListHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("상태 코드 = %d; 기대값 %d",
			resp.StatusCode, http.StatusMethodNotAllowed)
	}
}

// TestCreateHandler - 사용자 생성 핸들러 테스트
func TestCreateHandler(t *testing.T) {
	// JSON 요청 본문 생성
	jsonBody := `{"name":"홍길동","email":"hong@test.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/users",
		strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testCreateHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// 201 Created 확인
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("상태 코드 = %d; 기대값 %d",
			resp.StatusCode, http.StatusCreated)
	}

	// 응답 JSON 확인
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("JSON 파싱 실패: %v", err)
	}

	if result["name"] != "홍길동" {
		t.Errorf("name = %v; 기대값 홍길동", result["name"])
	}
}

// TestCreateHandlerValidation - 유효성 검사 테스트
func TestCreateHandlerValidation(t *testing.T) {
	// 이름 없는 요청
	jsonBody := `{"email":"test@test.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/users",
		strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testCreateHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("상태 코드 = %d; 기대값 %d",
			resp.StatusCode, http.StatusBadRequest)
	}
}

// === httptest.NewServer를 사용한 통합 테스트 ===

// TestServerIntegration - 테스트 서버를 사용한 통합 테스트
func TestServerIntegration(t *testing.T) {
	// 테스트용 ServeMux 생성
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", testHelloHandler)
	mux.HandleFunc("/api/users", testUserListHandler)

	// 테스트 서버 시작 (자동으로 랜덤 포트 할당)
	ts := httptest.NewServer(mux)
	defer ts.Close() // 테스트 끝나면 서버 종료

	// 실제 HTTP 요청을 보냅니다
	t.Run("hello 엔드포인트", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/hello?name=테스트")
		if err != nil {
			t.Fatal("요청 실패:", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		expected := "Hello, 테스트!"
		if string(body) != expected {
			t.Errorf("응답 = %q; 기대값 %q", string(body), expected)
		}
	})

	t.Run("users 엔드포인트", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/users")
		if err != nil {
			t.Fatal("요청 실패:", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("상태 코드 = %d; 기대값 200", resp.StatusCode)
		}
	})
}
