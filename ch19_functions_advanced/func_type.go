package main

import (
	"fmt"
	"strings"
)

// ============================================
// 함수 타입 변수 (함수를 값으로 사용)
// ============================================

// Operation은 두 정수를 받아 정수를 반환하는 함수 타입
type Operation func(int, int) int

// Predicate는 정수를 받아 불리언을 반환하는 함수 타입
type Predicate func(int) bool

// Transformer는 문자열을 변환하는 함수 타입
type Transformer func(string) string

// Apply는 Operation 함수를 적용한다
func Apply(op Operation, a, b int) int {
	return op(a, b)
}

// Filter는 조건에 맞는 요소만 반환한다
func Filter(nums []int, pred Predicate) []int {
	result := []int{}
	for _, n := range nums {
		if pred(n) {
			result = append(result, n)
		}
	}
	return result
}

// Map은 슬라이스의 모든 요소에 함수를 적용한다
func Map(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = fn(n)
	}
	return result
}

// TransformStrings는 문자열 슬라이스에 변환 함수를 적용한다
func TransformStrings(strs []string, fn Transformer) []string {
	result := make([]string, len(strs))
	for i, s := range strs {
		result[i] = fn(s)
	}
	return result
}

// Reduce는 슬라이스를 하나의 값으로 축소한다
func Reduce(nums []int, initial int, fn func(int, int) int) int {
	result := initial
	for _, n := range nums {
		result = fn(result, n)
	}
	return result
}

func main() {
	// 1. 함수 타입 변수
	fmt.Println("=== 함수 타입 변수 ===")
	var op Operation

	// 더하기 함수 할당
	op = func(a, b int) int { return a + b }
	fmt.Printf("더하기: 3 + 4 = %d\n", op(3, 4))

	// 곱하기 함수로 교체
	op = func(a, b int) int { return a * b }
	fmt.Printf("곱하기: 3 * 4 = %d\n", op(3, 4))

	// 2. 함수를 매개변수로 전달
	fmt.Println("\n=== 함수를 매개변수로 전달 ===")
	add := func(a, b int) int { return a + b }
	sub := func(a, b int) int { return a - b }
	mul := func(a, b int) int { return a * b }

	fmt.Println("Apply(add, 10, 3):", Apply(add, 10, 3))
	fmt.Println("Apply(sub, 10, 3):", Apply(sub, 10, 3))
	fmt.Println("Apply(mul, 10, 3):", Apply(mul, 10, 3))

	// 3. 함수 맵 (계산기 패턴)
	fmt.Println("\n=== 함수 맵 (계산기) ===")
	operations := map[string]Operation{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
	}

	a, b := 20, 4
	for symbol, op := range operations {
		fmt.Printf("  %d %s %d = %d\n", a, symbol, b, op(a, b))
	}

	// 4. Filter 함수 사용
	fmt.Println("\n=== Filter 함수 ===")
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 짝수만 필터링
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("짝수:", evens)

	// 5보다 큰 수만 필터링
	greaterThan5 := Filter(nums, func(n int) bool { return n > 5 })
	fmt.Println("5보다 큰 수:", greaterThan5)

	// 3의 배수만 필터링
	multiplesOf3 := Filter(nums, func(n int) bool { return n%3 == 0 })
	fmt.Println("3의 배수:", multiplesOf3)

	// 5. Map 함수 사용
	fmt.Println("\n=== Map 함수 ===")
	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Println("2배:", doubled)

	squared := Map(nums, func(n int) int { return n * n })
	fmt.Println("제곱:", squared)

	// 6. Reduce 함수 사용
	fmt.Println("\n=== Reduce 함수 ===")
	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Println("합계:", sum)

	product := Reduce([]int{1, 2, 3, 4, 5}, 1, func(acc, n int) int { return acc * n })
	fmt.Println("곱:", product)

	// 7. 문자열 변환
	fmt.Println("\n=== 문자열 변환 ===")
	words := []string{"hello", "world", "go", "language"}

	upper := TransformStrings(words, strings.ToUpper)
	fmt.Println("대문자:", upper)

	withPrefix := TransformStrings(words, func(s string) string {
		return ">> " + s
	})
	fmt.Println("접두사:", withPrefix)

	// 8. 함수를 반환하는 함수
	fmt.Println("\n=== 함수를 반환하는 함수 ===")
	greaterThan := func(threshold int) Predicate {
		return func(n int) bool {
			return n > threshold
		}
	}

	gt3 := greaterThan(3)
	gt7 := greaterThan(7)

	fmt.Println("3보다 큰 수:", Filter(nums, gt3))
	fmt.Println("7보다 큰 수:", Filter(nums, gt7))
}
