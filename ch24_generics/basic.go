package main

import (
	"cmp"
	"fmt"
)

// ============================================================
// 24.2 제네릭 함수 기본
// 타입 파라미터를 사용하여 타입에 독립적인 함수를 작성한다.
// ============================================================

// Print - 어떤 타입이든 출력할 수 있는 함수
// [T any]는 "T는 어떤 타입이든 될 수 있다"는 의미이다.
func Print[T any](value T) {
	fmt.Printf("값: %v (타입: %T)\n", value, value)
}

// Max - 두 값 중 큰 값을 반환한다.
// cmp.Ordered는 비교 연산자(<, >, <=, >=)를 지원하는 타입으로 제한한다.
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min - 두 값 중 작은 값을 반환한다.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Contains - 슬라이스에 특정 값이 포함되어 있는지 확인한다.
// comparable은 ==, != 비교가 가능한 타입으로 제한한다.
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// IndexOf - 슬라이스에서 특정 값의 인덱스를 반환한다. 없으면 -1을 반환한다.
func IndexOf[T comparable](slice []T, target T) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

// Swap - 두 값을 교환한다.
func Swap[T any](a, b T) (T, T) {
	return b, a
}

// Pair - 여러 타입 파라미터를 사용하는 함수
func Pair[K comparable, V any](key K, value V) string {
	return fmt.Sprintf("{%v: %v}", key, value)
}

func main() {
	fmt.Println("=== 제네릭 Print 함수 ===")
	Print(42)       // 타입 추론: int
	Print(3.14)     // 타입 추론: float64
	Print("안녕하세요")  // 타입 추론: string
	Print(true)     // 타입 추론: bool
	Print[int](100) // 명시적 타입 지정

	fmt.Println("\n=== 제네릭 Max / Min 함수 ===")
	fmt.Println("Max(3, 7):", Max(3, 7))
	fmt.Println("Max(3.14, 2.71):", Max(3.14, 2.71))
	fmt.Println("Max(\"abc\", \"xyz\"):", Max("abc", "xyz"))
	fmt.Println("Min(10, 5):", Min(10, 5))
	fmt.Println("Min(\"가\", \"나\"):", Min("가", "나"))

	fmt.Println("\n=== 제네릭 Contains 함수 ===")
	intSlice := []int{1, 2, 3, 4, 5}
	fmt.Println("Contains(intSlice, 3):", Contains(intSlice, 3))
	fmt.Println("Contains(intSlice, 9):", Contains(intSlice, 9))

	strSlice := []string{"Go", "Python", "Rust"}
	fmt.Println("Contains(strSlice, \"Go\"):", Contains(strSlice, "Go"))
	fmt.Println("Contains(strSlice, \"Java\"):", Contains(strSlice, "Java"))

	fmt.Println("\n=== 제네릭 IndexOf 함수 ===")
	fmt.Println("IndexOf(intSlice, 4):", IndexOf(intSlice, 4))
	fmt.Println("IndexOf(intSlice, 9):", IndexOf(intSlice, 9))
	fmt.Println("IndexOf(strSlice, \"Rust\"):", IndexOf(strSlice, "Rust"))

	fmt.Println("\n=== 제네릭 Swap 함수 ===")
	a, b := Swap(10, 20)
	fmt.Printf("Swap(10, 20) = (%d, %d)\n", a, b)

	x, y := Swap("hello", "world")
	fmt.Printf("Swap(\"hello\", \"world\") = (\"%s\", \"%s\")\n", x, y)

	fmt.Println("\n=== 여러 타입 파라미터 ===")
	fmt.Println(Pair("name", "홍길동"))
	fmt.Println(Pair(1, true))
	fmt.Println(Pair("age", 25))

	fmt.Println("\n=== 제네릭이 없었다면? ===")
	fmt.Println("제네릭 없이 Contains를 구현하려면:")
	fmt.Println("  - ContainsInt, ContainsString, ContainsFloat64...")
	fmt.Println("  - 또는 interface{}와 reflect를 사용 (느리고 안전하지 않음)")
	fmt.Println("제네릭 덕분에 하나의 함수로 모든 타입을 지원합니다!")
}
