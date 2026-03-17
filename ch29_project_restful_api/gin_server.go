// 29장: RESTful API 서버 만들기 - Gin 프레임워크 버전
// 설치: go get -u github.com/gin-gonic/gin
// 실행: go run gin_server.go
//
// 이 파일은 main.go와 같은 기능을 Gin 프레임워크로 구현한 예제이다.
// main.go와 함께 빌드할 수 없으므로 (main 함수 중복),
// 단독으로 실행하려면: go run gin_server.go
//
//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// GinStudent - 학생 정보 구조체 (Gin용 바인딩 태그 포함)
type GinStudent struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`     // 필수 필드
	Age   int    `json:"age" binding:"min=1,max=100"` // 범위 검증
	Grade int    `json:"grade" binding:"min=1,max=6"`
	Email string `json:"email" binding:"required,email"` // 이메일 형식 검증
}

// GinStudentStore - Gin용 학생 저장소
type GinStudentStore struct {
	mu       sync.RWMutex
	students map[int]GinStudent
	nextID   int
}

// NewGinStudentStore - 새 저장소 생성
func NewGinStudentStore() *GinStudentStore {
	s := &GinStudentStore{
		students: make(map[int]GinStudent),
		nextID:   1,
	}
	// 초기 데이터
	s.students[1] = GinStudent{ID: 1, Name: "김철수", Age: 20, Grade: 2, Email: "kim@school.com"}
	s.students[2] = GinStudent{ID: 2, Name: "이영희", Age: 21, Grade: 3, Email: "lee@school.com"}
	s.students[3] = GinStudent{ID: 3, Name: "박민수", Age: 19, Grade: 1, Email: "park@school.com"}
	s.nextID = 4
	return s
}

func main() {
	ginStore := NewGinStudentStore()

	// Gin 라우터 생성 (기본 미들웨어: Logger + Recovery)
	r := gin.Default()

	// === API 그룹 ===
	api := r.Group("/api")
	{
		// GET /api/students - 모든 학생 조회
		api.GET("/students", func(c *gin.Context) {
			ginStore.mu.RLock()
			defer ginStore.mu.RUnlock()

			students := make([]GinStudent, 0, len(ginStore.students))
			for _, s := range ginStore.students {
				students = append(students, s)
			}

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    students,
				"count":   len(students),
			})
		})

		// GET /api/students/:id - 특정 학생 조회
		api.GET("/students/:id", func(c *gin.Context) {
			// Gin은 경로 파라미터를 c.Param()으로 쉽게 추출한다
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "유효한 숫자 ID가 필요합니다",
				})
				return
			}

			ginStore.mu.RLock()
			student, ok := ginStore.students[id]
			ginStore.mu.RUnlock()

			if !ok {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"error":   fmt.Sprintf("ID %d인 학생을 찾을 수 없습니다", id),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    student,
			})
		})

		// POST /api/students - 새 학생 등록
		api.POST("/students", func(c *gin.Context) {
			var student GinStudent

			// ShouldBindJSON: JSON 바인딩 + 유효성 검사
			if err := c.ShouldBindJSON(&student); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "유효성 검사 실패",
					"details": err.Error(),
				})
				return
			}

			ginStore.mu.Lock()
			student.ID = ginStore.nextID
			ginStore.nextID++
			ginStore.students[student.ID] = student
			ginStore.mu.Unlock()

			c.JSON(http.StatusCreated, gin.H{
				"success": true,
				"data":    student,
			})
		})

		// PUT /api/students/:id - 학생 정보 수정
		api.PUT("/students/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "유효한 숫자 ID가 필요합니다",
				})
				return
			}

			var student GinStudent
			if err := c.ShouldBindJSON(&student); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			ginStore.mu.Lock()
			defer ginStore.mu.Unlock()

			if _, ok := ginStore.students[id]; !ok {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"error":   fmt.Sprintf("ID %d인 학생을 찾을 수 없습니다", id),
				})
				return
			}

			student.ID = id
			ginStore.students[id] = student

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    student,
			})
		})

		// DELETE /api/students/:id - 학생 삭제
		api.DELETE("/students/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "유효한 숫자 ID가 필요합니다",
				})
				return
			}

			ginStore.mu.Lock()
			defer ginStore.mu.Unlock()

			if _, ok := ginStore.students[id]; !ok {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"error":   fmt.Sprintf("ID %d인 학생을 찾을 수 없습니다", id),
				})
				return
			}

			delete(ginStore.students, id)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": fmt.Sprintf("ID %d 학생이 삭제되었습니다", id),
			})
		})
	}

	// 서버 시작
	fmt.Println("Gin 학생 관리 API 서버 시작: http://localhost:8091")
	fmt.Println("\nGin의 장점:")
	fmt.Println("  - c.Param('id')로 경로 파라미터를 쉽게 추출")
	fmt.Println("  - c.ShouldBindJSON()으로 자동 JSON 바인딩 + 유효성 검사")
	fmt.Println("  - r.GET(), r.POST() 등 메서드별 라우팅")
	fmt.Println("  - gin.H{}로 간편한 JSON 응답")
	fmt.Println("  - 미들웨어 자동 적용 (Logger, Recovery)")
	r.Run(":8091")
}
