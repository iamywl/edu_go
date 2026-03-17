package main

import "fmt"

func main() {
	// === 구구단 (세로 출력) ===
	fmt.Println("=== 구구단 (세로) ===")
	for i := 2; i <= 4; i++ { // 예시로 2~4단만
		fmt.Printf("--- %d단 ---\n", i)
		for j := 1; j <= 9; j++ {
			fmt.Printf("%d x %d = %2d\n", i, j, i*j)
		}
		fmt.Println()
	}

	// === 구구단 (가로 출력) ===
	fmt.Println("=== 구구단 (가로) ===")

	// 헤더 출력
	for i := 2; i <= 9; i++ {
		fmt.Printf("[ %d단 ]     ", i)
	}
	fmt.Println()

	// 각 행(1~9)을 순서대로 출력
	for j := 1; j <= 9; j++ {
		for i := 2; i <= 9; i++ {
			fmt.Printf("%d x %d = %2d  ", i, j, i*j)
		}
		fmt.Println()
	}
}
