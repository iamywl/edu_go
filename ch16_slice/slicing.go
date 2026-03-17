package main

import "fmt"

func main() {
	// ============================================
	// 슬라이싱 연산
	// ============================================

	// 1. 기본 슬라이싱 문법: s[start:end]
	fmt.Println("=== 기본 슬라이싱 ===")
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println("원본:", s)

	fmt.Println("s[2:5] =", s[2:5]) // [2 3 4] - 인덱스 2부터 4까지
	fmt.Println("s[:3]  =", s[:3])  // [0 1 2] - 처음부터 인덱스 2까지
	fmt.Println("s[7:]  =", s[7:])  // [7 8 9] - 인덱스 7부터 끝까지
	fmt.Println("s[:]   =", s[:])   // 전체 (원본과 동일)

	// 2. 슬라이싱과 원본 배열 공유
	fmt.Println("\n=== 원본 배열 공유 ===")
	original := []int{10, 20, 30, 40, 50}
	sub := original[1:4] // [20 30 40]

	fmt.Println("원본:", original)
	fmt.Println("부분:", sub)

	sub[0] = 999 // 부분 슬라이스 수정
	fmt.Println("\n부분[0]을 999로 변경 후:")
	fmt.Println("원본:", original) // [10 999 30 40 50] - 원본도 변경!
	fmt.Println("부분:", sub)      // [999 30 40]

	// 3. 슬라이싱의 len과 cap
	fmt.Println("\n=== 슬라이싱의 len과 cap ===")
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s1 := arr[2:5]
	fmt.Printf("arr[2:5] = %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
	// len=3, cap=8 (인덱스 2부터 배열 끝까지)

	s2 := arr[5:8]
	fmt.Printf("arr[5:8] = %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))
	// len=3, cap=5 (인덱스 5부터 배열 끝까지)

	// 4. 3-인덱스 슬라이싱 (용량 제한)
	fmt.Println("\n=== 3-인덱스 슬라이싱 ===")
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// data[start:end:max] - cap = max - start
	s3 := data[2:5:7]
	fmt.Printf("data[2:5:7] = %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))
	// len=3, cap=5

	s4 := data[2:5:5]
	fmt.Printf("data[2:5:5] = %v, len=%d, cap=%d\n", s4, len(s4), cap(s4))
	// len=3, cap=3 -> append 시 새 배열 할당 보장

	// 5. 요소 삭제 (순서 유지)
	fmt.Println("\n=== 요소 삭제 (순서 유지) ===")
	nums := []int{10, 20, 30, 40, 50}
	fmt.Println("삭제 전:", nums)

	idx := 2 // 인덱스 2 (값 30)을 삭제
	nums = append(nums[:idx], nums[idx+1:]...)
	fmt.Println("인덱스 2 삭제 후:", nums) // [10 20 40 50]

	// 6. 요소 삽입
	fmt.Println("\n=== 요소 삽입 ===")
	values := []int{1, 2, 4, 5}
	fmt.Println("삽입 전:", values)

	// 인덱스 2에 3을 삽입
	insertIdx := 2
	values = append(values[:insertIdx+1], values[insertIdx:]...)
	values[insertIdx] = 3
	fmt.Println("인덱스 2에 3 삽입:", values) // [1 2 3 4 5]

	// 7. copy()로 독립적인 복사
	fmt.Println("\n=== copy()로 독립적인 복사 ===")
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	copied := copy(dst, src)

	fmt.Printf("복사된 요소 수: %d\n", copied)
	fmt.Println("원본:", src)
	fmt.Println("복사본:", dst)

	dst[0] = 999
	fmt.Println("\n복사본[0]을 999로 변경 후:")
	fmt.Println("원본:", src)  // [1 2 3 4 5] - 변경 없음!
	fmt.Println("복사본:", dst) // [999 2 3 4 5]

	// 8. 배열에서 슬라이스 만들기
	fmt.Println("\n=== 배열에서 슬라이스 만들기 ===")
	array := [5]int{10, 20, 30, 40, 50}
	slice := array[1:4] // 배열의 일부를 슬라이스로

	fmt.Println("배열:", array)
	fmt.Println("슬라이스:", slice)

	slice[0] = 999
	fmt.Println("슬라이스 수정 후 배열:", array) // [10 999 30 40 50]
}
