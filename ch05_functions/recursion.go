// recursion.go
// 재귀 호출(Recursion) 예제
// 실행 방법: go run recursion.go

package main

import "fmt"

// === 팩토리얼 (재귀 버전) ===
// n! = n * (n-1) * (n-2) * ... * 1
// 0! = 1, 1! = 1
func factorial(n int) int {
	// 기저 조건 (base case): 재귀를 멈추는 조건
	if n <= 1 {
		return 1
	}
	// 재귀 호출: 문제의 크기를 줄여가며 자기 자신을 호출
	return n * factorial(n-1)
}

// === 팩토리얼 (반복문 버전) ===
func factorialLoop(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// === 피보나치 수열 (재귀 버전) ===
// F(0)=0, F(1)=1, F(n)=F(n-1)+F(n-2)
func fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// === 피보나치 수열 (반복문 버전 - 더 효율적) ===
func fibonacciLoop(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	prev, curr := 0, 1
	for i := 2; i <= n; i++ {
		prev, curr = curr, prev+curr
	}
	return curr
}

// === 거듭제곱 (재귀 버전) ===
// base^exp
func power(base, exp int) int {
	if exp == 0 {
		return 1
	}
	return base * power(base, exp-1)
}

// === 자릿수 합 (재귀 버전) ===
// 예: digitSum(123) = 1 + 2 + 3 = 6
func digitSum(n int) int {
	if n < 0 {
		n = -n // 음수 처리
	}
	if n < 10 {
		return n // 한 자릿수이면 그대로 반환
	}
	return n%10 + digitSum(n/10)
}

// === 하노이의 탑 ===
// n개의 원판을 from에서 to로 이동 (via를 경유)
func hanoi(n int, from, to, via string) {
	if n == 1 {
		fmt.Printf("  원판 %d: %s → %s\n", n, from, to)
		return
	}
	// 1. n-1개의 원판을 from에서 via로 이동
	hanoi(n-1, from, via, to)
	// 2. 가장 큰 원판을 from에서 to로 이동
	fmt.Printf("  원판 %d: %s → %s\n", n, from, to)
	// 3. n-1개의 원판을 via에서 to로 이동
	hanoi(n-1, via, to, from)
}

// === 재귀 호출 과정을 시각적으로 보여주는 팩토리얼 ===
func factorialVerbose(n int, depth int) int {
	// 들여쓰기로 재귀 깊이 표시
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	fmt.Printf("%sfactorial(%d) 호출\n", indent, n)

	if n <= 1 {
		fmt.Printf("%sfactorial(%d) = 1 반환\n", indent, n)
		return 1
	}

	result := n * factorialVerbose(n-1, depth+1)
	fmt.Printf("%sfactorial(%d) = %d * factorial(%d) = %d 반환\n",
		indent, n, n, n-1, result)
	return result
}

func main() {
	fmt.Println("===== 팩토리얼 =====")
	for i := 0; i <= 10; i++ {
		fmt.Printf("%d! = %d\n", i, factorial(i))
	}

	fmt.Println()
	fmt.Println("===== 팩토리얼 실행 과정 (시각화) =====")
	fmt.Println("factorial(5):")
	result := factorialVerbose(5, 0)
	fmt.Printf("최종 결과: %d\n", result)

	fmt.Println()
	fmt.Println("===== 피보나치 수열 =====")
	fmt.Print("피보나치 수열: ")
	for i := 0; i <= 15; i++ {
		fmt.Printf("%d ", fibonacciLoop(i))
	}
	fmt.Println()

	// 재귀 vs 반복문 비교 (작은 수로 비교)
	fmt.Println()
	fmt.Println("===== 재귀 vs 반복문 비교 =====")
	n := 10
	fmt.Printf("factorial(%d): 재귀=%d, 반복문=%d\n",
		n, factorial(n), factorialLoop(n))
	fmt.Printf("fibonacci(%d): 재귀=%d, 반복문=%d\n",
		n, fibonacci(n), fibonacciLoop(n))

	// 참고: fibonacci(40) 이상은 재귀 버전이 매우 느림!
	// 재귀: O(2^n), 반복문: O(n)

	fmt.Println()
	fmt.Println("===== 거듭제곱 =====")
	fmt.Printf("2^10 = %d\n", power(2, 10))
	fmt.Printf("3^5 = %d\n", power(3, 5))

	fmt.Println()
	fmt.Println("===== 자릿수 합 =====")
	fmt.Printf("digitSum(123) = %d\n", digitSum(123))
	fmt.Printf("digitSum(9999) = %d\n", digitSum(9999))
	fmt.Printf("digitSum(0) = %d\n", digitSum(0))

	fmt.Println()
	fmt.Println("===== 하노이의 탑 (원판 3개) =====")
	fmt.Println("A(시작), B(경유), C(도착)")
	hanoi(3, "A", "C", "B")
}
