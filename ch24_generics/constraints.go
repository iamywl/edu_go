package main

import (
	"cmp"
	"fmt"
	"strings"
)

// ============================================================
// 24.3 타입 제약 조건 (Type Constraints)
// 인터페이스를 사용하여 타입 파라미터에 제약을 둘 수 있습니다.
// ============================================================

// ---- 커스텀 제약 조건 정의 ----

// Number - 숫자 타입만 허용하는 제약 조건
// ~int는 기본 타입이 int인 모든 타입을 포함한다. (예: type MyInt int)
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Integer - 정수 타입만 허용
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Signed - 부호 있는 정수
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Stringer - String() 메서드를 가진 타입
type Stringer interface {
	String() string
}

// ---- 제약 조건을 사용하는 함수 ----

// Sum - 숫자 슬라이스의 합계를 구한다.
func Sum[T Number](numbers []T) T {
	var total T
	for _, n := range numbers {
		total += n
	}
	return total
}

// Average - 숫자 슬라이스의 평균을 구한다.
func Average[T Number](numbers []T) float64 {
	if len(numbers) == 0 {
		return 0
	}
	sum := Sum(numbers)
	return float64(sum) / float64(len(numbers))
}

// Abs - 절대값을 반환한다. (부호 있는 정수만)
func Abs[T Signed](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

// Filter - 조건을 만족하는 요소만 반환한다.
func Filter[T any](items []T, fn func(T) bool) []T {
	var result []T
	for _, item := range items {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map - 슬라이스의 각 요소를 변환한다.
func Map[T, U any](items []T, fn func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

// Reduce - 슬라이스를 하나의 값으로 축약한다.
func Reduce[T, U any](items []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, item := range items {
		result = fn(result, item)
	}
	return result
}

// SortSlice - cmp.Ordered 제약 조건을 사용한 정렬
func SortSlice[T cmp.Ordered](items []T) []T {
	// 간단한 버블 정렬 (교육 목적)
	result := make([]T, len(items))
	copy(result, items)

	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j] < result[i] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}

// ---- ~ (틸다) 연산자 데모를 위한 커스텀 타입 ----

// Celsius - 섭씨 온도 (기본 타입: float64)
type Celsius float64

// Fahrenheit - 화씨 온도 (기본 타입: float64)
type Fahrenheit float64

// Score - 점수 (기본 타입: int)
type Score int

func main() {
	fmt.Println("=== Number 제약 조건 ===")

	intSlice := []int{1, 2, 3, 4, 5}
	floatSlice := []float64{1.5, 2.5, 3.5}

	fmt.Println("Sum(int):", Sum(intSlice))         // 15
	fmt.Println("Sum(float64):", Sum(floatSlice))   // 7.5
	fmt.Println("Average(int):", Average(intSlice)) // 3.0
	fmt.Println("Average(float64):", Average(floatSlice))

	fmt.Println("\n=== ~ (틸다) 연산자 ===")

	// ~float64 덕분에 Celsius, Fahrenheit도 Sum에 사용 가능
	temps := []Celsius{36.5, 37.0, 38.2, 36.8}
	fmt.Printf("온도 합계: %.1f°C\n", Sum(temps))
	fmt.Printf("평균 온도: %.1f°C\n", Average(temps))

	scores := []Score{85, 92, 78, 95, 88}
	fmt.Printf("점수 합계: %d\n", Sum(scores))
	fmt.Printf("평균 점수: %.1f\n", Average(scores))

	fmt.Println("\n=== Abs (부호 있는 정수) ===")
	fmt.Println("Abs(-5):", Abs(-5))
	fmt.Println("Abs(3):", Abs(3))
	fmt.Println("Abs(-100):", Abs(int64(-100)))

	fmt.Println("\n=== Filter 함수 ===")

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 짝수만 필터링
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("짝수:", evens)

	// 5보다 큰 수만 필터링
	big := Filter(nums, func(n int) bool { return n > 5 })
	fmt.Println("5보다 큰 수:", big)

	// 문자열 필터링
	words := []string{"Go", "Python", "Rust", "Java", "Ruby"}
	short := Filter(words, func(s string) bool { return len(s) <= 3 })
	fmt.Println("짧은 단어:", short)

	fmt.Println("\n=== Map 함수 ===")

	// int -> string 변환
	strs := Map(nums, func(n int) string {
		return fmt.Sprintf("%d번", n)
	})
	fmt.Println("변환:", strs[:5])

	// 제곱
	squares := Map(nums, func(n int) int { return n * n })
	fmt.Println("제곱:", squares)

	// 대문자 변환
	uppers := Map(words, strings.ToUpper)
	fmt.Println("대문자:", uppers)

	fmt.Println("\n=== Reduce 함수 ===")

	// 합계
	total := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Println("합계:", total)

	// 곱셈
	product := Reduce([]int{1, 2, 3, 4, 5}, 1, func(acc, n int) int { return acc * n })
	fmt.Println("곱:", product)

	// 문자열 연결
	joined := Reduce(words, "", func(acc string, s string) string {
		if acc == "" {
			return s
		}
		return acc + ", " + s
	})
	fmt.Println("연결:", joined)

	fmt.Println("\n=== SortSlice 함수 ===")
	unsorted := []int{5, 3, 8, 1, 9, 2}
	sorted := SortSlice(unsorted)
	fmt.Println("정렬 전:", unsorted)
	fmt.Println("정렬 후:", sorted)

	sortedStrs := SortSlice([]string{"바나나", "사과", "귤", "포도"})
	fmt.Println("문자열 정렬:", sortedStrs)
}
