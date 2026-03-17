package main

import (
	"errors"
	"fmt"
)

// ============================================================
// 21.2 커스텀 에러 타입
// error 인터페이스를 구현하면 어떤 타입이든 에러로 사용할 수 있습니다.
// ============================================================

// ValidationError - 입력값 검증 에러
type ValidationError struct {
	Field   string // 에러가 발생한 필드 이름
	Message string // 에러 메시지
}

// Error 메서드를 구현하여 error 인터페이스를 충족한다.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("검증 실패 [%s]: %s", e.Field, e.Message)
}

// HttpError - HTTP 응답 에러
type HttpError struct {
	StatusCode int    // HTTP 상태 코드
	Message    string // 에러 메시지
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// validateUser - 사용자 입력을 검증하는 함수
func validateUser(name string, age int) error {
	if name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "이름은 비어있을 수 없습니다",
		}
	}
	if age < 0 || age > 150 {
		return &ValidationError{
			Field:   "age",
			Message: fmt.Sprintf("나이가 유효하지 않습니다: %d", age),
		}
	}
	return nil
}

// fetchData - HTTP 요청을 시뮬레이션하는 함수
func fetchData(url string) error {
	if url == "" {
		return &HttpError{
			StatusCode: 400,
			Message:    "잘못된 요청: URL이 비어있습니다",
		}
	}
	if url == "https://example.com/secret" {
		return &HttpError{
			StatusCode: 403,
			Message:    "접근이 금지되었습니다",
		}
	}
	if url == "https://example.com/missing" {
		return &HttpError{
			StatusCode: 404,
			Message:    "리소스를 찾을 수 없습니다",
		}
	}
	return nil
}

func main() {
	fmt.Println("=== 커스텀 에러 타입 사용 ===")

	// ValidationError 발생
	err := validateUser("", 25)
	if err != nil {
		fmt.Println("에러:", err)

		// errors.As를 사용하여 에러 타입 확인
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("  -> 필드: %s\n", valErr.Field)
			fmt.Printf("  -> 메시지: %s\n", valErr.Message)
		}
	}

	// 나이 검증 실패
	err = validateUser("홍길동", -5)
	if err != nil {
		fmt.Println("\n에러:", err)

		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("  -> 필드: %s\n", valErr.Field)
		}
	}

	fmt.Println("\n=== errors.As를 사용한 HTTP 에러 처리 ===")

	// 다양한 HTTP 에러 테스트
	urls := []string{
		"",
		"https://example.com/secret",
		"https://example.com/missing",
		"https://example.com/ok",
	}

	for _, url := range urls {
		err := fetchData(url)
		if err != nil {
			var httpErr *HttpError
			if errors.As(err, &httpErr) {
				// 상태 코드에 따라 다르게 처리
				switch {
				case httpErr.StatusCode == 400:
					fmt.Printf("[클라이언트 에러] %s\n", httpErr.Message)
				case httpErr.StatusCode == 403:
					fmt.Printf("[권한 에러] %s\n", httpErr.Message)
				case httpErr.StatusCode == 404:
					fmt.Printf("[미발견] %s\n", httpErr.Message)
				case httpErr.StatusCode >= 500:
					fmt.Printf("[서버 에러] %s\n", httpErr.Message)
				}
			}
			continue
		}
		fmt.Printf("[성공] %s 요청 완료\n", url)
	}

	fmt.Println("\n=== errors.Is 사용 예제 ===")

	// 센티널 에러 정의
	var ErrNotFound = errors.New("리소스를 찾을 수 없습니다")
	var ErrUnauthorized = errors.New("인증되지 않았습니다")

	// 에러 비교
	checkError := func(err error) {
		switch {
		case errors.Is(err, ErrNotFound):
			fmt.Println("-> 리소스를 찾을 수 없습니다. 다른 경로를 시도하세요.")
		case errors.Is(err, ErrUnauthorized):
			fmt.Println("-> 로그인이 필요합니다.")
		default:
			fmt.Println("-> 알 수 없는 에러:", err)
		}
	}

	checkError(ErrNotFound)
	checkError(ErrUnauthorized)
	checkError(errors.New("기타 에러"))
}
