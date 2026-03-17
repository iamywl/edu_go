package main

import "fmt"

func main() {
	// === 배열 선언과 초기화 ===
	fmt.Println("=== 배열 선언과 초기화 ===")

	// 기본 선언 (zero value로 초기화)
	var a [5]int
	fmt.Println("기본 선언:", a) // [0 0 0 0 0]

	// 선언과 동시에 초기화
	b := [3]string{"사과", "바나나", "체리"}
	fmt.Println("초기화:", b)

	// 크기 자동 계산
	c := [...]int{10, 20, 30, 40, 50}
	fmt.Println("자동 크기:", c, "/ 길이:", len(c))

	// 특정 인덱스만 초기화
	d := [5]int{1: 100, 3: 300}
	fmt.Println("부분 초기화:", d) // [0 100 0 300 0]

	// === 요소 접근과 수정 ===
	fmt.Println("\n=== 요소 접근과 수정 ===")
	fruits := [4]string{"딸기", "포도", "수박", "참외"}

	fmt.Println("첫 번째:", fruits[0])
	fmt.Println("마지막:", fruits[len(fruits)-1])

	fruits[1] = "블루베리" // 수정
	fmt.Println("수정 후:", fruits)

	// === for문으로 순회 ===
	fmt.Println("\n=== 배열 순회 ===")
	scores := [5]int{90, 85, 78, 92, 88}

	// 인덱스 기반 순회
	fmt.Println("인덱스 순회:")
	for i := 0; i < len(scores); i++ {
		fmt.Printf("  scores[%d] = %d\n", i, scores[i])
	}

	// range 순회
	fmt.Println("range 순회:")
	for i, v := range scores {
		fmt.Printf("  scores[%d] = %d\n", i, v)
	}

	// === 합계와 평균 ===
	fmt.Println("\n=== 합계와 평균 ===")
	sum := 0
	for _, v := range scores {
		sum += v
	}
	avg := float64(sum) / float64(len(scores))
	fmt.Printf("점수: %v\n", scores)
	fmt.Printf("합계: %d, 평균: %.1f\n", sum, avg)

	// === 최댓값, 최솟값 찾기 ===
	fmt.Println("\n=== 최댓값, 최솟값 ===")
	nums := [7]int{23, 45, 12, 67, 34, 89, 56}

	maxVal := nums[0]
	minVal := nums[0]

	for _, v := range nums {
		if v > maxVal {
			maxVal = v
		}
		if v < minVal {
			minVal = v
		}
	}
	fmt.Printf("배열: %v\n", nums)
	fmt.Printf("최댓값: %d, 최솟값: %d\n", maxVal, minVal)

	// === 배열은 값 타입 (복사) ===
	fmt.Println("\n=== 배열은 값 타입 ===")
	original := [3]int{1, 2, 3}
	copied := original // 전체 복사
	copied[0] = 999

	fmt.Println("원본:", original) // [1 2 3] (변경 안 됨)
	fmt.Println("복사:", copied)   // [999 2 3]

	// === 배열 비교 ===
	fmt.Println("\n=== 배열 비교 ===")
	x := [3]int{1, 2, 3}
	y := [3]int{1, 2, 3}
	z := [3]int{1, 2, 4}

	fmt.Println("x == y:", x == y) // true
	fmt.Println("x == z:", x == z) // false
}
