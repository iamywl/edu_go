// 29장: RESTful API 서버 만들기 - API 테스트
// 실행: go test -v
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// getTestRouter - 테스트용 라우터를 생성한다
// 각 테스트마다 새로운 저장소를 사용하여 테스트 간 격리를 보장한다
func getTestRouter() http.Handler {
	// 전역 store를 테스트용으로 재설정
	store = NewStudentStore()
	return setupRouter()
}

// TestGetAllStudents - 모든 학생 목록 조회 테스트
func TestGetAllStudents(t *testing.T) {
	router := getTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 상태 코드 확인
	if w.Code != http.StatusOK {
		t.Fatalf("상태 코드 = %d; 기대값 %d", w.Code, http.StatusOK)
	}

	// 응답 파싱
	var resp APIResponse
	body, _ := io.ReadAll(w.Body)
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("JSON 파싱 실패: %v", err)
	}

	if !resp.Success {
		t.Error("응답 success = false; 기대값 true")
	}

	// 초기 데이터 3명 확인
	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatal("Data가 배열이 아닙니다")
	}
	if len(data) != 3 {
		t.Errorf("학생 수 = %d; 기대값 3", len(data))
	}
}

// TestGetStudentByID - 특정 학생 조회 테스트
func TestGetStudentByID(t *testing.T) {
	router := getTestRouter()

	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantName   string
	}{
		{"존재하는 학생", "1", http.StatusOK, "김철수"},
		{"존재하는 학생 2", "2", http.StatusOK, "이영희"},
		{"존재하지 않는 학생", "999", http.StatusNotFound, ""},
		{"잘못된 ID", "abc", http.StatusBadRequest, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet,
				"/api/students/"+tt.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("상태 코드 = %d; 기대값 %d", w.Code, tt.wantStatus)
			}

			if tt.wantName != "" {
				var resp APIResponse
				body, _ := io.ReadAll(w.Body)
				json.Unmarshal(body, &resp)

				data, ok := resp.Data.(map[string]interface{})
				if !ok {
					t.Fatal("Data가 객체가 아닙니다")
				}
				if data["name"] != tt.wantName {
					t.Errorf("이름 = %v; 기대값 %s", data["name"], tt.wantName)
				}
			}
		})
	}
}

// TestCreateStudent - 학생 생성 테스트
func TestCreateStudent(t *testing.T) {
	router := getTestRouter()

	t.Run("정상 생성", func(t *testing.T) {
		newStudent := Student{
			Name:  "홍길동",
			Age:   22,
			Grade: 4,
			Email: "hong@school.com",
		}
		jsonBody, _ := json.Marshal(newStudent)

		req := httptest.NewRequest(http.MethodPost, "/api/students",
			bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 201 Created 확인
		if w.Code != http.StatusCreated {
			t.Fatalf("상태 코드 = %d; 기대값 %d", w.Code, http.StatusCreated)
		}

		// ID가 할당되었는지 확인
		var resp APIResponse
		body, _ := io.ReadAll(w.Body)
		json.Unmarshal(body, &resp)

		data := resp.Data.(map[string]interface{})
		if data["id"].(float64) == 0 {
			t.Error("ID가 할당되지 않았습니다")
		}
		if data["name"] != "홍길동" {
			t.Errorf("이름 = %v; 기대값 홍길동", data["name"])
		}
	})

	t.Run("이름 누락", func(t *testing.T) {
		jsonBody := `{"age":20,"grade":2}`
		req := httptest.NewRequest(http.MethodPost, "/api/students",
			bytes.NewReader([]byte(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("상태 코드 = %d; 기대값 %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("잘못된 JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/students",
			bytes.NewReader([]byte("잘못된 JSON")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("상태 코드 = %d; 기대값 %d", w.Code, http.StatusBadRequest)
		}
	})
}

// TestUpdateStudent - 학생 정보 수정 테스트
func TestUpdateStudent(t *testing.T) {
	router := getTestRouter()

	t.Run("정상 수정", func(t *testing.T) {
		updated := Student{
			Name:  "김철수(수정)",
			Age:   21,
			Grade: 3,
			Email: "kim_updated@school.com",
		}
		jsonBody, _ := json.Marshal(updated)

		req := httptest.NewRequest(http.MethodPut, "/api/students/1",
			bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("상태 코드 = %d; 기대값 %d", w.Code, http.StatusOK)
		}

		var resp APIResponse
		body, _ := io.ReadAll(w.Body)
		json.Unmarshal(body, &resp)

		data := resp.Data.(map[string]interface{})
		if data["name"] != "김철수(수정)" {
			t.Errorf("이름 = %v; 기대값 김철수(수정)", data["name"])
		}
	})

	t.Run("존재하지 않는 학생 수정", func(t *testing.T) {
		jsonBody := `{"name":"테스트"}`
		req := httptest.NewRequest(http.MethodPut, "/api/students/999",
			bytes.NewReader([]byte(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("상태 코드 = %d; 기대값 %d", w.Code, http.StatusNotFound)
		}
	})
}

// TestDeleteStudent - 학생 삭제 테스트
func TestDeleteStudent(t *testing.T) {
	router := getTestRouter()

	t.Run("정상 삭제", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/students/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("상태 코드 = %d; 기대값 %d", w.Code, http.StatusOK)
		}

		// 삭제 후 조회하면 404
		req2 := httptest.NewRequest(http.MethodGet, "/api/students/1", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusNotFound {
			t.Errorf("삭제 후 조회 상태 코드 = %d; 기대값 %d",
				w2.Code, http.StatusNotFound)
		}
	})

	t.Run("존재하지 않는 학생 삭제", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/students/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("상태 코드 = %d; 기대값 %d", w.Code, http.StatusNotFound)
		}
	})
}

// TestMethodNotAllowed - 허용되지 않은 메서드 테스트
func TestMethodNotAllowed(t *testing.T) {
	router := getTestRouter()

	req := httptest.NewRequest(http.MethodPatch, "/api/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("상태 코드 = %d; 기대값 %d",
			w.Code, http.StatusMethodNotAllowed)
	}
}

// TestCRUDFlow - CRUD 전체 흐름 통합 테스트
func TestCRUDFlow(t *testing.T) {
	router := getTestRouter()

	// 1. 학생 생성
	newStudent := `{"name":"테스트학생","age":22,"grade":4,"email":"test@school.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/students",
		bytes.NewReader([]byte(newStudent)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("생성 실패: 상태 코드 = %d", w.Code)
	}

	// 생성된 학생의 ID 추출
	var createResp APIResponse
	body, _ := io.ReadAll(w.Body)
	json.Unmarshal(body, &createResp)
	data := createResp.Data.(map[string]interface{})
	id := int(data["id"].(float64))
	t.Logf("생성된 학생 ID: %d", id)

	// 2. 생성된 학생 조회
	req = httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/api/students/%d", id), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("조회 실패: 상태 코드 = %d", w.Code)
	}

	// 3. 학생 정보 수정
	updateBody := `{"name":"수정된학생","age":23,"grade":4,"email":"updated@school.com"}`
	req = httptest.NewRequest(http.MethodPut,
		fmt.Sprintf("/api/students/%d", id),
		bytes.NewReader([]byte(updateBody)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("수정 실패: 상태 코드 = %d", w.Code)
	}

	// 4. 학생 삭제
	req = httptest.NewRequest(http.MethodDelete,
		fmt.Sprintf("/api/students/%d", id), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("삭제 실패: 상태 코드 = %d", w.Code)
	}

	// 5. 삭제 후 조회 확인 (404)
	req = httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/api/students/%d", id), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("삭제 후 조회: 상태 코드 = %d; 기대값 %d",
			w.Code, http.StatusNotFound)
	}

	t.Log("CRUD 전체 흐름 테스트 통과!")
}
