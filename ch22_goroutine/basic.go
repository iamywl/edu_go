package main

import (
	"fmt"
	"time"
)

// ============================================================
// 22.2 고루틴 기본 사용법
// go 키워드를 사용하여 함수를 고루틴으로 실행한다.
// ============================================================

// printNumbers - 숫자를 출력하는 함수
func printNumbers(prefix string, count int) {
	for i := 1; i <= count; i++ {
		fmt.Printf("%s: %d\n", prefix, i)
		time.Sleep(100 * time.Millisecond) // 실행 순서를 관찰하기 위한 지연
	}
}

// greet - 인사를 출력하는 함수
func greet(name string) {
	fmt.Printf("안녕하세요, %s님!\n", name)
}

func main() {
	fmt.Println("=== 고루틴 기본 사용 ===")

	// 일반 함수 호출 (순차 실행)
	fmt.Println("\n--- 순차 실행 ---")
	start := time.Now()
	printNumbers("순차A", 3)
	printNumbers("순차B", 3)
	fmt.Printf("순차 실행 소요 시간: %v\n", time.Since(start))

	// 고루틴으로 실행 (동시 실행)
	fmt.Println("\n--- 동시 실행 (고루틴) ---")
	start = time.Now()
	go printNumbers("고루틴A", 3)         // 고루틴으로 실행
	go printNumbers("고루틴B", 3)         // 고루틴으로 실행
	time.Sleep(500 * time.Millisecond) // 고루틴이 완료될 때까지 대기
	fmt.Printf("동시 실행 소요 시간: %v\n", time.Since(start))

	fmt.Println("\n=== 익명 함수 고루틴 ===")

	// 익명 함수를 고루틴으로 실행
	go func() {
		fmt.Println("익명 함수 고루틴 1에서 실행!")
	}()

	// 인자를 전달하는 익명 함수 고루틴
	go func(msg string) {
		fmt.Println("메시지:", msg)
	}("고루틴에서 보낸 메시지")

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== 여러 고루틴 동시 실행 ===")

	// 5개의 고루틴을 동시에 실행
	names := []string{"김철수", "이영희", "박민수", "정수진", "최동현"}
	for _, name := range names {
		go greet(name)
	}

	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=== 주의: 클로저와 고루틴 ===")

	// 잘못된 예: 루프 변수를 직접 캡처하면 예상치 못한 결과가 발생할 수 있습니다.
	fmt.Println("--- 잘못된 예 (루프 변수 캡처) ---")
	for i := 0; i < 5; i++ {
		go func() {
			// Go 1.22부터는 루프 변수가 각 반복에서 새로 생성되므로
			// 이 문제가 해결되었지만, 이전 버전에서는 모두 같은 값이 출력될 수 있습니다.
			fmt.Printf("  i = %d\n", i)
		}()
	}
	time.Sleep(100 * time.Millisecond)

	// 올바른 예: 인자로 전달
	fmt.Println("--- 올바른 예 (인자로 전달) ---")
	for i := 0; i < 5; i++ {
		go func(n int) {
			fmt.Printf("  n = %d\n", n)
		}(i) // 현재 i 값을 인자로 전달
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n프로그램 종료")
	// 참고: 실제 코드에서는 time.Sleep 대신 sync.WaitGroup이나 채널을 사용하세요!
}
