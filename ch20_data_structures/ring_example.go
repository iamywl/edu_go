package main

import (
	"container/ring"
	"fmt"
	"time"
)

// ============================================
// 링 (container/ring) - 원형 연결 리스트
// ============================================

func main() {
	// 1. 링 생성과 기본 사용
	fmt.Println("=== 링 생성과 기본 사용 ===")
	r := ring.New(5) // 크기 5인 링 생성
	fmt.Println("링 크기:", r.Len())

	// 값 설정 (순회하며 설정)
	for i := 0; i < r.Len(); i++ {
		r.Value = fmt.Sprintf("항목_%d", i+1)
		r = r.Next()
	}

	// 2. 링 순회 (Do 메서드)
	fmt.Println("\n=== Do()로 순회 ===")
	r.Do(func(val any) {
		fmt.Printf("  %v\n", val)
	})

	// 3. Next()와 Prev()로 이동
	fmt.Println("\n=== Next/Prev로 이동 ===")
	fmt.Println("현재:", r.Value)
	fmt.Println("다음:", r.Next().Value)
	fmt.Println("이전:", r.Prev().Value)

	// 3칸 이동 (음수면 반대 방향)
	moved := r.Move(3)
	fmt.Println("3칸 이동:", moved.Value)
	moved = r.Move(-2)
	fmt.Println("-2칸 이동:", moved.Value)

	// 4. 원형 특성 확인 - 계속 순회하면 처음으로 돌아옴
	fmt.Println("\n=== 원형 특성 (2바퀴 순회) ===")
	current := r
	for i := 0; i < r.Len()*2; i++ {
		fmt.Printf("  [%d] %v\n", i, current.Value)
		current = current.Next()
	}

	// 5. 실용 예제: 최근 N개 로그 유지
	fmt.Println("\n=== 최근 로그 유지 (크기 5) ===")
	logRing := ring.New(5)

	// 로그 메시지 추가 (8개 추가하면 처음 3개는 덮어씀)
	logs := []string{
		"서버 시작",
		"DB 연결 성공",
		"요청 수신: GET /api",
		"요청 처리 완료",
		"캐시 갱신",
		"요청 수신: POST /api", // 여기부터 "서버 시작"을 덮어씀
		"인증 성공",
		"응답 전송",
	}

	for _, log := range logs {
		logRing.Value = log
		logRing = logRing.Next()
		fmt.Printf("  추가: %s\n", log)
	}

	fmt.Println("\n최근 5개 로그:")
	logRing.Do(func(val any) {
		if val != nil {
			fmt.Printf("  - %v\n", val)
		}
	})

	// 6. 실용 예제: 라운드 로빈 스케줄러
	fmt.Println("\n=== 라운드 로빈 스케줄러 ===")
	servers := ring.New(3)

	// 서버 설정
	serverNames := []string{"서버A (8080)", "서버B (8081)", "서버C (8082)"}
	for i := 0; i < servers.Len(); i++ {
		servers.Value = serverNames[i]
		servers = servers.Next()
	}

	// 요청을 라운드 로빈으로 분배
	for i := 1; i <= 7; i++ {
		fmt.Printf("  요청 #%d -> %v\n", i, servers.Value)
		servers = servers.Next() // 다음 서버로 이동
	}

	// 7. 실용 예제: 시계 (12시간)
	fmt.Println("\n=== 시계 시뮬레이션 ===")
	clock := ring.New(12)
	for i := 0; i < 12; i++ {
		clock.Value = i + 1 // 1~12시
		clock = clock.Next()
	}

	// 현재 시각부터 5시간 후까지 표시
	now := time.Now().Hour() % 12
	if now == 0 {
		now = 12
	}
	// 현재 시각으로 이동
	clock = clock.Move(now - 1)
	fmt.Printf("현재: %d시\n", clock.Value)
	for i := 1; i <= 5; i++ {
		clock = clock.Next()
		fmt.Printf("  +%d시간: %d시\n", i, clock.Value)
	}

	// 8. 링 연결 (Link)과 분리 (Unlink)
	fmt.Println("\n=== 링 연결과 분리 ===")
	r1 := ring.New(3)
	r2 := ring.New(3)

	// r1 값 설정
	for i := 0; i < 3; i++ {
		r1.Value = fmt.Sprintf("A%d", i+1)
		r1 = r1.Next()
	}

	// r2 값 설정
	for i := 0; i < 3; i++ {
		r2.Value = fmt.Sprintf("B%d", i+1)
		r2 = r2.Next()
	}

	fmt.Print("r1: ")
	r1.Do(func(v any) { fmt.Print(v, " ") })
	fmt.Println()

	fmt.Print("r2: ")
	r2.Do(func(v any) { fmt.Print(v, " ") })
	fmt.Println()

	// 두 링 연결
	r1.Link(r2)
	fmt.Print("연결 후: ")
	r1.Do(func(v any) { fmt.Print(v, " ") })
	fmt.Println()
	fmt.Println("연결 후 크기:", r1.Len())
}
