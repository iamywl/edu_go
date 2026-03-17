package main

import (
	"fmt"
	"math"
)

func main() {
	// === 레이블 없이 break (안쪽 for만 종료) ===
	fmt.Println("=== break (레이블 없음) ===")
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				break // 안쪽 for만 종료
			}
			fmt.Printf("(%d, %d) ", i, j)
		}
	}
	fmt.Println()

	// === 레이블 + break (바깥 for까지 종료) ===
	fmt.Println("\n=== break OuterLoop ===")
OuterBreak:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Println("→ (1,1)에서 전체 종료!")
				break OuterBreak // 바깥 for까지 종료
			}
			fmt.Printf("(%d, %d) ", i, j)
		}
	}

	// === 레이블 + continue (바깥 for의 다음 반복으로) ===
	fmt.Println("\n=== continue OuterLoop ===")
OuterContinue:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				continue OuterContinue // 바깥 for의 다음 i로
			}
			fmt.Printf("(%d, %d) ", i, j)
		}
	}
	fmt.Println()

	// === 실전 예제: 소수(prime) 찾기 ===
	fmt.Println("\n=== 1~50 소수 찾기 ===")
	fmt.Print("소수: ")

	for num := 2; num <= 50; num++ {
		isPrime := true
		// num의 제곱근까지만 확인하면 충분
		limit := int(math.Sqrt(float64(num)))

		for div := 2; div <= limit; div++ {
			if num%div == 0 {
				isPrime = false
				break // 약수를 찾았으면 더 볼 필요 없음
			}
		}

		if isPrime {
			fmt.Printf("%d ", num)
		}
	}
	fmt.Println()

	// === 실전 예제: 레이블로 소수 찾기 ===
	fmt.Println("\n=== 소수 찾기 (레이블 활용) ===")
	fmt.Print("소수: ")

	// 레이블 없이 시도하면? → 잘못된 결과가 나온다
	// continue가 안쪽 for문에만 적용되어, 나누어떨어지는 수도 소수로 출력된다
	for num := 2; num <= 50; num++ {
		isPrime := true
		limit := int(math.Sqrt(float64(num)))
		for div := 2; div <= limit; div++ {
			if num%div == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			fmt.Printf("%d ", num)
		}
	}

	// 레이블을 활용한 올바른 구현
	fmt.Print("\n소수(레이블): ")
NextNum:
	for num := 2; num <= 50; num++ {
		limit := int(math.Sqrt(float64(num)))
		for div := 2; div <= limit; div++ {
			if num%div == 0 {
				continue NextNum // 바깥 for의 다음 num으로
			}
		}
		// 여기에 도달 = 어떤 수로도 나누어지지 않음 = 소수
		fmt.Printf("%d ", num)
	}
	fmt.Println()

	// === 실전 예제: 2차원 탐색에서 값 찾기 ===
	fmt.Println("\n=== 2차원 배열에서 값 찾기 ===")
	matrix := [3][4]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
	}
	target := 7

Search:
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if matrix[i][j] == target {
				fmt.Printf("값 %d을(를) [%d][%d] 위치에서 찾았습니다!\n", target, i, j)
				break Search // 찾았으면 전체 탐색 종료
			}
		}
	}
}
