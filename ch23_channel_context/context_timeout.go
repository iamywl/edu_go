package main

import (
	"context"
	"fmt"
	"time"
)

// ============================================================
// 23.2 context.WithTimeout과 context.WithValue 사용법
// WithTimeout: 지정된 시간이 지나면 자동으로 취소됩니다.
// WithValue: 컨텍스트에 키-값 쌍을 저장한다.
// ============================================================

// contextKey - WithValue에서 사용할 키 타입
// 문자열 키 충돌을 방지하기 위해 별도의 타입을 정의한다.
type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

// simulateAPICall - 외부 API 호출을 시뮬레이션한다.
func simulateAPICall(ctx context.Context, apiName string, delay time.Duration) (string, error) {
	// 요청 ID를 컨텍스트에서 추출한다.
	reqID := "unknown"
	if v := ctx.Value(requestIDKey); v != nil {
		reqID = v.(string)
	}

	fmt.Printf("  [%s] API '%s' 호출 시작 (예상 소요: %v)\n", reqID, apiName, delay)

	// 작업 채널
	resultCh := make(chan string, 1)
	go func() {
		time.Sleep(delay) // API 호출 시뮬레이션
		resultCh <- fmt.Sprintf("API '%s' 응답 데이터", apiName)
	}()

	// 컨텍스트 취소 또는 결과 대기
	select {
	case result := <-resultCh:
		fmt.Printf("  [%s] API '%s' 응답 성공\n", reqID, apiName)
		return result, nil
	case <-ctx.Done():
		fmt.Printf("  [%s] API '%s' 취소됨: %v\n", reqID, apiName, ctx.Err())
		return "", ctx.Err()
	}
}

// processRequest - HTTP 요청 처리를 시뮬레이션한다.
func processRequest(ctx context.Context) error {
	userID := "unknown"
	if v := ctx.Value(userIDKey); v != nil {
		userID = v.(string)
	}
	fmt.Printf("  사용자 '%s'의 요청을 처리합니다.\n", userID)

	// 1단계: 사용자 정보 조회 (빠른 API)
	result, err := simulateAPICall(ctx, "사용자 조회", 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("사용자 조회 실패: %w", err)
	}
	fmt.Printf("  1단계 완료: %s\n", result)

	// 2단계: 주문 내역 조회 (느린 API)
	result, err = simulateAPICall(ctx, "주문 조회", 300*time.Millisecond)
	if err != nil {
		return fmt.Errorf("주문 조회 실패: %w", err)
	}
	fmt.Printf("  2단계 완료: %s\n", result)

	return nil
}

func main() {
	fmt.Println("=== context.WithTimeout 기본 사용법 ===")

	// 2초 타임아웃 컨텍스트
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 반드시 호출하여 리소스 누수를 방지한다.

	// 남은 시간 확인
	if deadline, ok := ctx.Deadline(); ok {
		fmt.Printf("만료 시점: %v (남은 시간: %v)\n", deadline.Format("15:04:05.000"), time.Until(deadline).Round(time.Millisecond))
	}

	// 빠른 API - 타임아웃 내에 완료
	result, err := simulateAPICall(ctx, "빠른API", 500*time.Millisecond)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Println("결과:", result)
	}

	fmt.Println("\n=== 타임아웃 발생 예제 ===")

	// 500ms 타임아웃 설정
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer shortCancel()

	// 느린 API - 타임아웃 발생!
	result, err = simulateAPICall(shortCtx, "느린API", 1*time.Second)
	if err != nil {
		fmt.Println("에러 발생:", err)
		if err == context.DeadlineExceeded {
			fmt.Println("-> 타임아웃으로 인한 취소입니다.")
		}
	}

	fmt.Println("\n=== context.WithValue 사용법 ===")

	// 기본 컨텍스트에 값을 추가한다.
	ctx = context.Background()
	ctx = context.WithValue(ctx, requestIDKey, "REQ-001")
	ctx = context.WithValue(ctx, userIDKey, "user-42")

	// 타임아웃도 함께 설정
	ctx, cancel = context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	fmt.Println("요청 처리 시작:")
	err = processRequest(ctx)
	if err != nil {
		fmt.Println("요청 처리 실패:", err)
	} else {
		fmt.Println("요청 처리 성공!")
	}

	fmt.Println("\n=== 타임아웃으로 인한 요청 실패 ===")

	// 짧은 타임아웃으로 두 번째 API 호출에서 실패
	ctx2 := context.WithValue(context.Background(), requestIDKey, "REQ-002")
	ctx2 = context.WithValue(ctx2, userIDKey, "user-99")
	ctx2, cancel2 := context.WithTimeout(ctx2, 250*time.Millisecond)
	defer cancel2()

	fmt.Println("요청 처리 시작 (짧은 타임아웃):")
	err = processRequest(ctx2)
	if err != nil {
		fmt.Println("요청 처리 실패:", err)
	}

	fmt.Println("\n=== context.WithDeadline 사용법 ===")

	// 특정 시점에 만료되는 컨텍스트
	deadline := time.Now().Add(300 * time.Millisecond)
	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), deadline)
	defer deadlineCancel()

	fmt.Printf("데드라인: %v\n", deadline.Format("15:04:05.000"))

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("작업 완료")
	case <-deadlineCtx.Done():
		fmt.Println("데드라인 초과:", deadlineCtx.Err())
	}

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
