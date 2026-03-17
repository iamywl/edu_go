// float_error.go
// 부동소수점 오차와 해결 방법 예제
// 실행 방법: go run float_error.go

package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("===== 부동소수점 오차 문제 =====")

	// 0.1 + 0.2가 0.3이 아닌 이유
	a := 0.1
	b := 0.2
	c := 0.3

	fmt.Printf("a = 0.1 → 실제 저장값: %.20f\n", a)
	fmt.Printf("b = 0.2 → 실제 저장값: %.20f\n", b)
	fmt.Printf("c = 0.3 → 실제 저장값: %.20f\n", c)
	fmt.Println()

	sum := a + b
	fmt.Printf("a + b = %.20f\n", sum)
	fmt.Printf("c     = %.20f\n", c)
	fmt.Printf("차이  = %.20f\n", sum-c)
	fmt.Println()

	// == 비교는 실패한다!
	fmt.Printf("a + b == c : %t (틀린 결과!)\n", a+b == c)

	fmt.Println()
	fmt.Println("===== 해결 방법 1: epsilon 비교 =====")

	// 두 값의 차이가 매우 작은 값(epsilon)보다 작으면 같다고 판단
	epsilon := 1e-10 // 0.0000000001

	diff := math.Abs(sum - c)
	fmt.Printf("차이의 절대값: %e\n", diff)
	fmt.Printf("epsilon:       %e\n", epsilon)
	fmt.Printf("차이 < epsilon: %t\n", diff < epsilon)

	if diff < epsilon {
		fmt.Println("→ epsilon 비교: 0.1 + 0.2는 0.3과 같다!")
	}

	fmt.Println()
	fmt.Println("===== 해결 방법 2: math.Nextafter =====")

	// Nextafter(x, y): x에서 y 방향으로 다음 표현 가능한 float64 값
	next := math.Nextafter(c, c+1)
	fmt.Printf("0.3 다음 표현 가능한 값: %.20f\n", next)
	fmt.Printf("0.3과의 차이:            %e\n", next-c)

	// Nextafter를 사용한 비교
	if sum <= math.Nextafter(c, c+1) && sum >= math.Nextafter(c, c-1) {
		fmt.Println("→ Nextafter 비교: 0.1 + 0.2는 0.3과 같다!")
	}

	fmt.Println()
	fmt.Println("===== 해결 방법 3: 정수로 변환 =====")

	// 금융 계산처럼 정확한 값이 필요하면 정수(센트 단위)로 계산
	// 0.1달러 = 10센트, 0.2달러 = 20센트
	cents1 := 10 // 0.1달러
	cents2 := 20 // 0.2달러
	cents3 := 30 // 0.3달러

	totalCents := cents1 + cents2
	fmt.Printf("10센트 + 20센트 = %d센트\n", totalCents)
	fmt.Printf("%d == %d : %t (정확한 비교!)\n", totalCents, cents3, totalCents == cents3)
	fmt.Printf("달러로 변환: $%.2f\n", float64(totalCents)/100)

	fmt.Println()
	fmt.Println("===== 다른 오차 사례들 =====")

	// 반복적인 덧셈에서 오차가 누적됨
	var total float64
	for i := 0; i < 10; i++ {
		total += 0.1
	}
	fmt.Printf("0.1을 10번 더하면: %.20f (1.0이어야 하지만...)\n", total)
	fmt.Printf("== 1.0 비교: %t\n", total == 1.0)

	// 큰 수와 작은 수의 덧셈에서 작은 수가 무시될 수 있음
	big := 1e16 // 10,000,000,000,000,000
	small := 1.0
	result := big + small - big
	fmt.Printf("\n1e16 + 1.0 - 1e16 = %f (1.0이어야 하지만...)\n", result)

	fmt.Println()
	fmt.Println("===== 실수 비교 함수 만들기 =====")

	// 실무에서 사용할 수 있는 실수 비교 함수
	fmt.Printf("almostEqual(0.1+0.2, 0.3) = %t\n", almostEqual(0.1+0.2, 0.3))
	fmt.Printf("almostEqual(1.0, 1.0) = %t\n", almostEqual(1.0, 1.0))
	fmt.Printf("almostEqual(1.0, 2.0) = %t\n", almostEqual(1.0, 2.0))
}

// almostEqual: 두 실수가 거의 같은지 비교하는 함수
func almostEqual(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}
