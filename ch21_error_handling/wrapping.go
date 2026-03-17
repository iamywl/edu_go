package main

import (
	"errors"
	"fmt"
	"os"
)

// ============================================================
// 21.2 에러 래핑 (Error Wrapping)
// %w 동사를 사용하여 에러를 감싸면서 맥락 정보를 추가할 수 있습니다.
// 래핑된 에러는 errors.Is, errors.As로 원본 에러를 추출할 수 있습니다.
// ============================================================

// 센티널 에러 정의 - 패키지 수준의 재사용 가능한 에러
var (
	ErrDatabase   = errors.New("데이터베이스 에러")
	ErrConnection = errors.New("연결 실패")
)

// DatabaseError - 커스텀 에러 타입
type DatabaseError struct {
	Query   string
	Message string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("DB 에러 [쿼리: %s]: %s", e.Query, e.Message)
}

// ---- 3단계 함수 호출 체인 (에러 래핑 데모) ----

// connectDB - 데이터베이스 연결을 시뮬레이션한다.
func connectDB(host string) error {
	if host == "" {
		// 원본 에러를 %w로 래핑한다.
		return fmt.Errorf("DB 연결 실패: %w", ErrConnection)
	}
	return nil
}

// queryDB - 데이터베이스 쿼리를 시뮬레이션한다.
func queryDB(query string) error {
	err := connectDB("") // 빈 호스트로 연결 시도
	if err != nil {
		// 에러를 다시 래핑하여 맥락을 추가한다.
		return fmt.Errorf("쿼리 '%s' 실행 중: %w", query, err)
	}
	return nil
}

// handleRequest - 요청 처리를 시뮬레이션한다.
func handleRequest(endpoint string) error {
	err := queryDB("SELECT * FROM users")
	if err != nil {
		// 한 번 더 래핑한다.
		return fmt.Errorf("엔드포인트 '%s' 처리 중: %w", endpoint, err)
	}
	return nil
}

// ---- 파일 관련 에러 래핑 예제 ----

// readFile - 파일 읽기를 시도한다.
func readFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// os 에러를 래핑한다.
		return nil, fmt.Errorf("파일 '%s' 읽기 실패: %w", path, err)
	}
	return data, nil
}

// loadConfig - 설정 파일을 로드한다.
func loadConfig(path string) error {
	_, err := readFile(path)
	if err != nil {
		return fmt.Errorf("설정 로드 실패: %w", err)
	}
	return nil
}

func main() {
	fmt.Println("=== 에러 래핑 체인 ===")

	// 3단계 래핑된 에러 발생
	err := handleRequest("/api/users")
	if err != nil {
		// 전체 에러 메시지 (래핑된 맥락이 모두 포함됨)
		fmt.Println("전체 에러 메시지:")
		fmt.Println(" ", err)
		fmt.Println()

		// errors.Is로 원본 에러를 찾을 수 있습니다.
		if errors.Is(err, ErrConnection) {
			fmt.Println("-> errors.Is: 연결 실패 에러를 감지했습니다!")
		}
	}

	fmt.Println("\n=== 파일 에러 래핑 ===")

	// 존재하지 않는 설정 파일 로드
	err = loadConfig("/존재하지않는/config.json")
	if err != nil {
		fmt.Println("에러:", err)
		fmt.Println()

		// os.ErrNotExist와 비교 (래핑된 에러에서도 동작)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("-> errors.Is: 파일이 존재하지 않습니다!")
		}
	}

	fmt.Println("\n=== errors.Unwrap으로 에러 체인 풀기 ===")

	// 에러 체인을 하나씩 풀어봅니다.
	err = handleRequest("/api/users")
	for err != nil {
		fmt.Printf("에러: %s\n", err.Error())
		err = errors.Unwrap(err)
	}

	fmt.Println("\n=== %v vs %w 비교 ===")

	originalErr := errors.New("원본 에러")

	// %v: 문자열로만 변환, 에러 체인 끊김
	wrappedV := fmt.Errorf("맥락 추가: %v", originalErr)
	// %w: 에러 체인 유지
	wrappedW := fmt.Errorf("맥락 추가: %w", originalErr)

	fmt.Printf("%%v로 래핑: errors.Is 결과 = %v\n", errors.Is(wrappedV, originalErr))
	fmt.Printf("%%w로 래핑: errors.Is 결과 = %v\n", errors.Is(wrappedW, originalErr))

	fmt.Println("\n=== 커스텀 에러 타입과 래핑 ===")

	// 커스텀 에러를 래핑한다.
	dbErr := &DatabaseError{
		Query:   "INSERT INTO users",
		Message: "중복 키 제약 조건 위반",
	}
	wrappedErr := fmt.Errorf("사용자 생성 실패: %w", dbErr)

	fmt.Println("래핑된 에러:", wrappedErr)

	// errors.As로 원본 타입을 추출한다.
	var extractedErr *DatabaseError
	if errors.As(wrappedErr, &extractedErr) {
		fmt.Printf("  -> 쿼리: %s\n", extractedErr.Query)
		fmt.Printf("  -> 메시지: %s\n", extractedErr.Message)
	}
}
