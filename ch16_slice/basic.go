package main

import "fmt"

func main() {
	// ============================================
	// 슬라이스 생성과 기본 사용법
	// ============================================

	// 1. 리터럴로 슬라이스 생성
	fruits := []string{"사과", "바나나", "체리", "딸기"}
	fmt.Println("=== 리터럴로 생성 ===")
	fmt.Println("과일:", fruits)
	fmt.Println("길이:", len(fruits))
	fmt.Println("용량:", cap(fruits))

	// 2. make()로 슬라이스 생성
	fmt.Println("\n=== make()로 생성 ===")

	// make(타입, 길이) - 길이와 용량이 동일
	s1 := make([]int, 5)
	fmt.Printf("s1: %v (len=%d, cap=%d)\n", s1, len(s1), cap(s1))

	// make(타입, 길이, 용량) - 용량을 별도로 지정
	s2 := make([]int, 3, 10)
	fmt.Printf("s2: %v (len=%d, cap=%d)\n", s2, len(s2), cap(s2))

	// 3. var로 선언 (nil 슬라이스)
	var s3 []int
	fmt.Println("\n=== nil 슬라이스 ===")
	fmt.Printf("s3: %v (len=%d, cap=%d)\n", s3, len(s3), cap(s3))
	fmt.Println("s3 == nil:", s3 == nil) // true

	// 빈 슬라이스 (nil이 아님)
	s4 := []int{}
	fmt.Printf("s4: %v (len=%d, cap=%d)\n", s4, len(s4), cap(s4))
	fmt.Println("s4 == nil:", s4 == nil) // false

	// 4. 슬라이스에 요소 추가 (append)
	fmt.Println("\n=== append()로 요소 추가 ===")
	var numbers []int
	for i := 1; i <= 5; i++ {
		numbers = append(numbers, i)
		fmt.Printf("추가 %d: %v (len=%d, cap=%d)\n", i, numbers, len(numbers), cap(numbers))
	}

	// 5. 여러 요소 한번에 추가
	fmt.Println("\n=== 여러 요소 추가 ===")
	numbers = append(numbers, 6, 7, 8)
	fmt.Println("여러 요소 추가 후:", numbers)

	// 6. 슬라이스끼리 합치기
	fmt.Println("\n=== 슬라이스 합치기 ===")
	extra := []int{9, 10}
	numbers = append(numbers, extra...) // ...을 사용하여 슬라이스 풀어넣기
	fmt.Println("합친 후:", numbers)

	// 7. 슬라이스 순회
	fmt.Println("\n=== 슬라이스 순회 ===")
	colors := []string{"빨강", "초록", "파랑", "노랑"}

	// for range 사용
	for i, color := range colors {
		fmt.Printf("  colors[%d] = %s\n", i, color)
	}

	// 8. 배열과 슬라이스의 차이
	fmt.Println("\n=== 배열 vs 슬라이스 ===")

	// 배열: 값 복사
	arr1 := [3]int{1, 2, 3}
	arr2 := arr1
	arr2[0] = 100
	fmt.Println("배열 arr1:", arr1) // [1 2 3] - 변경 없음
	fmt.Println("배열 arr2:", arr2) // [100 2 3]

	// 슬라이스: 참조 복사 (같은 배열을 가리킴)
	slice1 := []int{1, 2, 3}
	slice2 := slice1
	slice2[0] = 100
	fmt.Println("슬라이스 slice1:", slice1) // [100 2 3] - 변경됨!
	fmt.Println("슬라이스 slice2:", slice2) // [100 2 3]
}
