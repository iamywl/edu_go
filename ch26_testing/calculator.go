// 26장: 테스트와 벤치마크하기 - 간단한 계산기 함수
// 이 파일은 테스트 대상이 되는 계산기 함수들을 정의한다.
package main

import (
	"errors"
	"fmt"
	"math"
)

// Add - 두 정수를 더한다
func Add(a, b int) int {
	return a + b
}

// Subtract - 두 정수를 뺍니다
func Subtract(a, b int) int {
	return a - b
}

// Multiply - 두 정수를 곱한다
func Multiply(a, b int) int {
	return a * b
}

// Divide - 두 실수를 나눕니다 (0으로 나누기 에러 처리 포함)
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("0으로 나눌 수 없습니다")
	}
	return a / b, nil
}

// Abs - 절대값을 반환한다
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Power - a의 b제곱을 반환한다
func Power(a, b float64) float64 {
	return math.Pow(a, b)
}

// Factorial - n의 팩토리얼을 계산한다 (재귀)
func Factorial(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("음수의 팩토리얼은 정의되지 않습니다")
	}
	if n == 0 || n == 1 {
		return 1, nil
	}
	result, err := Factorial(n - 1)
	if err != nil {
		return 0, err
	}
	return n * result, nil
}

// Fibonacci - n번째 피보나치 수를 반환한다
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func main() {
	// 계산기 함수 사용 예제
	fmt.Println("=== 간단한 계산기 ===")
	fmt.Printf("Add(2, 3) = %d\n", Add(2, 3))
	fmt.Printf("Subtract(10, 4) = %d\n", Subtract(10, 4))
	fmt.Printf("Multiply(3, 7) = %d\n", Multiply(3, 7))

	result, err := Divide(10, 3)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("Divide(10, 3) = %.4f\n", result)
	}

	_, err = Divide(10, 0)
	if err != nil {
		fmt.Println("Divide(10, 0) 에러:", err)
	}

	fmt.Printf("Abs(-42) = %d\n", Abs(-42))
	fmt.Printf("Power(2, 10) = %.0f\n", Power(2, 10))

	fact, _ := Factorial(10)
	fmt.Printf("Factorial(10) = %d\n", fact)

	fmt.Printf("Fibonacci(10) = %d\n", Fibonacci(10))
}
