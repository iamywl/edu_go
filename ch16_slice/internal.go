package main

import (
	"fmt"
	"unsafe"
)

// SliceHeader 구조를 시각적으로 보여주는 예제

func main() {
	// ============================================
	// 슬라이스 내부 구조 (SliceHeader) 이해하기
	// ============================================

	// 1. len과 cap의 차이
	fmt.Println("=== len vs cap ===")
	s := make([]int, 3, 5) // 길이 3, 용량 5
	fmt.Printf("슬라이스: %v\n", s)
	fmt.Printf("길이(len): %d - 현재 사용 중인 요소 수\n", len(s))
	fmt.Printf("용량(cap): %d - 내부 배열의 전체 크기\n", cap(s))

	// 2. append()에 따른 len, cap 변화 관찰
	fmt.Println("\n=== append()에 따른 변화 ===")
	var nums []int
	prevCap := cap(nums)

	for i := 0; i < 20; i++ {
		nums = append(nums, i)
		if cap(nums) != prevCap {
			fmt.Printf("  len=%2d, cap=%2d (용량 증가! %d -> %d)\n",
				len(nums), cap(nums), prevCap, cap(nums))
			prevCap = cap(nums)
		}
	}
	fmt.Println("  -> 용량이 부족하면 보통 2배로 증가합니다")

	// 3. 슬라이스 대입은 참조 복사
	fmt.Println("\n=== 참조 복사 확인 ===")
	original := []int{10, 20, 30}
	copied := original // SliceHeader만 복사됨

	fmt.Println("원본:", original)
	fmt.Println("복사본:", copied)

	copied[1] = 999 // 복사본 수정
	fmt.Println("\n복사본[1]을 999로 변경 후:")
	fmt.Println("원본:", original) // [10 999 30] - 원본도 변경됨!
	fmt.Println("복사본:", copied)  // [10 999 30]
	fmt.Println("-> 같은 내부 배열을 공유하기 때문!")

	// 4. append()가 새 배열을 만드는 경우
	fmt.Println("\n=== append()와 새 배열 할당 ===")
	a := make([]int, 3, 3) // 용량이 꽉 참
	a[0], a[1], a[2] = 1, 2, 3

	b := a // 같은 배열을 가리킴
	fmt.Printf("a의 주소: %p\n", a)
	fmt.Printf("b의 주소: %p (같음)\n", b)

	// a에 append하면 용량 초과 -> 새 배열 할당
	a = append(a, 4)
	fmt.Printf("\nappend 후 a의 주소: %p (변경됨 - 새 배열!)\n", a)
	fmt.Printf("         b의 주소: %p (그대로)\n", b)

	a[0] = 100
	fmt.Println("\na[0]을 100으로 변경:")
	fmt.Println("a:", a) // [100 2 3 4]
	fmt.Println("b:", b) // [1 2 3] - 영향 없음! 다른 배열이므로

	// 5. 슬라이스 헤더의 크기
	fmt.Println("\n=== SliceHeader 크기 ===")
	var slice []int
	fmt.Printf("슬라이스 헤더 크기: %d 바이트\n", unsafe.Sizeof(slice))
	fmt.Println("  - Data (포인터): 8 바이트")
	fmt.Println("  - Len (길이):    8 바이트")
	fmt.Println("  - Cap (용량):    8 바이트")
	fmt.Println("  -> 슬라이스를 함수에 전달해도 24바이트만 복사됩니다")

	// 6. 함수에 슬라이스 전달 (참조처럼 동작)
	fmt.Println("\n=== 함수에 슬라이스 전달 ===")
	data := []int{1, 2, 3, 4, 5}
	fmt.Println("함수 호출 전:", data)
	doubleValues(data)
	fmt.Println("함수 호출 후:", data) // 값이 2배로 변경됨
}

// doubleValues는 슬라이스의 모든 요소를 2배로 만든다
// 슬라이스 헤더(24바이트)만 복사되며, 내부 배열은 공유됨
func doubleValues(s []int) {
	for i := range s {
		s[i] *= 2
	}
}
