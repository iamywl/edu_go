package main

import "fmt"

func main() {
	// =============================================
	// for range 다양한 패턴 예제
	// =============================================

	// --- 슬라이스 순회 ---
	fmt.Println("=== 슬라이스 순회 ===")
	fruits := []string{"사과", "바나나", "포도", "딸기"}

	// 패턴 1: 인덱스와 값 모두 사용
	fmt.Println("인덱스 + 값:")
	for i, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}

	// 패턴 2: 인덱스만 사용
	fmt.Println("\n인덱스만:")
	for i := range fruits {
		fmt.Printf("  인덱스: %d\n", i)
	}

	// 패턴 3: 값만 사용 (인덱스 무시)
	fmt.Println("\n값만:")
	for _, fruit := range fruits {
		fmt.Printf("  과일: %s\n", fruit)
	}

	// --- 맵 순회 ---
	fmt.Println("\n=== 맵 순회 ===")
	scores := map[string]int{
		"수학": 95,
		"영어": 88,
		"과학": 92,
	}

	for subject, score := range scores {
		fmt.Printf("  %s: %d점\n", subject, score)
	}
	// 주의: 맵의 순회 순서는 보장되지 않습니다!

	// --- 문자열 순회 (룬 단위) ---
	fmt.Println("\n=== 문자열 순회 (룬 단위) ===")
	greeting := "Go 언어!"

	// range는 UTF-8 룬 단위로 순회한다
	for i, r := range greeting {
		fmt.Printf("  바이트 오프셋=%d, 문자='%c', 유니코드=%U\n", i, r, r)
	}

	// --- 문자열 바이트 단위 순회 ---
	fmt.Println("\n=== 문자열 바이트 단위 순회 ===")
	for i := 0; i < len(greeting); i++ {
		fmt.Printf("  byte[%d] = 0x%02x\n", i, greeting[i])
	}

	// --- 채널 순회 ---
	fmt.Println("\n=== 채널 순회 ===")
	ch := make(chan string)

	go func() {
		messages := []string{"첫 번째", "두 번째", "세 번째"}
		for _, msg := range messages {
			ch <- msg
		}
		close(ch) // 채널을 닫아야 range 루프가 종료됩니다
	}()

	for msg := range ch {
		fmt.Printf("  수신: %s\n", msg)
	}

	// --- 2차원 슬라이스 순회 ---
	fmt.Println("\n=== 2차원 슬라이스 순회 ===")
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	for row, cols := range matrix {
		for col, val := range cols {
			fmt.Printf("  [%d][%d] = %d", row, col, val)
		}
		fmt.Println()
	}

	// --- range에서 값은 복사본 ---
	fmt.Println("\n=== range 값은 복사본 ===")
	numbers := []int{10, 20, 30}

	// 이렇게 하면 원본이 변경되지 않습니다!
	for _, v := range numbers {
		v *= 2 // v는 복사본이므로 원본에 영향 없음
		_ = v
	}
	fmt.Printf("  변경 시도 후: %v (변경되지 않음)\n", numbers)

	// 원본을 변경하려면 인덱스를 사용해야 한다
	for i := range numbers {
		numbers[i] *= 2
	}
	fmt.Printf("  인덱스로 변경 후: %v (변경됨)\n", numbers)

	// --- Go 1.22+ 정수 range ---
	fmt.Println("\n=== 정수 range (Go 1.22+) ===")
	fmt.Print("  0부터 4까지: ")
	for i := range 5 {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
}
