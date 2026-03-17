package main

import "fmt"

// ============================================
// 가변 인수 함수 (Variadic Functions)
// ============================================

// Sum은 정수를 가변 인수로 받아 합계를 반환한다
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Max는 정수를 가변 인수로 받아 최대값을 반환한다
// 인수가 없으면 0을 반환한다
func Max(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	max := nums[0]
	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}
	return max
}

// JoinStrings는 구분자와 문자열 가변 인수를 받아 결합한다
// 일반 매개변수와 가변 인수를 함께 사용하는 예시
func JoinStrings(separator string, strs ...string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += separator
		}
		result += s
	}
	return result
}

// PrintFormatted는 접두사와 다양한 타입의 값을 출력한다
// any 타입의 가변 인수 (fmt.Println과 유사)
func PrintFormatted(prefix string, values ...any) {
	fmt.Printf("[%s] ", prefix)
	for i, v := range values {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%v(%T)", v, v)
	}
	fmt.Println()
}

// Stats는 가변 인수로 받은 정수의 합계, 평균, 최대, 최소를 반환한다
func Stats(nums ...int) (sum int, avg float64, max int, min int) {
	if len(nums) == 0 {
		return 0, 0, 0, 0
	}

	max = nums[0]
	min = nums[0]

	for _, n := range nums {
		sum += n
		if n > max {
			max = n
		}
		if n < min {
			min = n
		}
	}
	avg = float64(sum) / float64(len(nums))
	return
}

func main() {
	// 1. 기본 가변 인수 사용
	fmt.Println("=== 기본 가변 인수 ===")
	fmt.Println("Sum():", Sum())                           // 0
	fmt.Println("Sum(1, 2):", Sum(1, 2))                   // 3
	fmt.Println("Sum(1, 2, 3, 4, 5):", Sum(1, 2, 3, 4, 5)) // 15

	// 2. 슬라이스를 가변 인수로 전달 (... 사용)
	fmt.Println("\n=== 슬라이스 -> 가변 인수 ===")
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Println("슬라이스:", numbers)
	fmt.Println("Sum(numbers...):", Sum(numbers...)) // ... 으로 풀어서 전달

	// 3. Max 함수
	fmt.Println("\n=== Max 함수 ===")
	fmt.Println("Max(3, 7, 2, 9, 4):", Max(3, 7, 2, 9, 4))
	fmt.Println("Max(100):", Max(100))

	// 4. 일반 매개변수 + 가변 인수
	fmt.Println("\n=== 일반 매개변수 + 가변 인수 ===")
	result := JoinStrings(", ", "Go", "Python", "Java", "Rust")
	fmt.Println("JoinStrings:", result)

	result2 := JoinStrings(" -> ", "시작", "처리", "완료")
	fmt.Println("JoinStrings:", result2)

	// 5. any 타입 가변 인수
	fmt.Println("\n=== any 타입 가변 인수 ===")
	PrintFormatted("정보", "이름", "김철수", "나이", 25)
	PrintFormatted("데이터", 42, 3.14, true, "hello")

	// 6. Stats 함수 (다중 반환값 + 가변 인수)
	fmt.Println("\n=== Stats 함수 ===")
	scores := []int{85, 92, 78, 95, 88, 76, 90}
	sum, avg, max, min := Stats(scores...)
	fmt.Printf("점수: %v\n", scores)
	fmt.Printf("합계: %d, 평균: %.1f, 최대: %d, 최소: %d\n", sum, avg, max, min)

	// 7. 가변 인수 전달 (함수 간 전달)
	fmt.Println("\n=== 가변 인수 전달 ===")
	wrapper := func(nums ...int) {
		// 가변 인수를 다른 함수에 그대로 전달
		fmt.Printf("  전달 받은 값: %v, 합계: %d\n", nums, Sum(nums...))
	}
	wrapper(1, 2, 3)
	wrapper(10, 20, 30, 40)
}
