package main

import "fmt"

func main() {
	// === 기본 for문 (C 스타일) ===
	fmt.Println("=== 기본 for문 ===")
	for i := 0; i < 5; i++ {
		fmt.Printf("i = %d\n", i)
	}

	// === while 스타일 ===
	fmt.Println("\n=== while 스타일 ===")
	count := 1
	for count <= 5 {
		fmt.Printf("count = %d\n", count)
		count++
	}

	// === 무한 루프 + break ===
	fmt.Println("\n=== 무한 루프 + break ===")
	sum := 0
	n := 1
	for {
		sum += n
		if sum > 20 {
			fmt.Printf("합이 20을 초과! (sum=%d, n=%d)\n", sum, n)
			break
		}
		n++
	}

	// === 1부터 100까지 합 ===
	fmt.Println("\n=== 1~100 합계 ===")
	total := 0
	for i := 1; i <= 100; i++ {
		total += i
	}
	fmt.Printf("1부터 100까지 합: %d\n", total)

	// === continue 사용: 홀수만 출력 ===
	fmt.Println("\n=== continue (홀수만 출력) ===")
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			continue // 짝수는 건너뜀
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// === continue 사용: 3의 배수이면서 5의 배수가 아닌 수 ===
	fmt.Println("\n=== 3의 배수 & 5의 배수 아님 (1~30) ===")
	for i := 1; i <= 30; i++ {
		if i%3 != 0 { // 3의 배수가 아니면 건너뜀
			continue
		}
		if i%5 == 0 { // 5의 배수이면 건너뜀
			continue
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// === 역순 반복 ===
	fmt.Println("\n=== 역순 반복 ===")
	for i := 5; i >= 1; i-- {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// === 2칸씩 건너뛰기 ===
	fmt.Println("\n=== 2칸씩 건너뛰기 ===")
	for i := 0; i <= 10; i += 2 {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// === 별 찍기 ===
	fmt.Println("\n=== 별 찍기 ===")
	for i := 1; i <= 5; i++ {
		for j := 0; j < i; j++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}
