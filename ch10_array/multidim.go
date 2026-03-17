package main

import "fmt"

func main() {
	// === 2차원 배열 선언과 초기화 ===
	fmt.Println("=== 2차원 배열 기본 ===")

	// 기본 선언 (zero value)
	var grid [2][3]int
	fmt.Println("기본 선언:", grid)

	// 초기화와 함께
	matrix := [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	// === 요소 접근 ===
	fmt.Println("\n=== 요소 접근 ===")
	fmt.Println("matrix[0][0] =", matrix[0][0]) // 1
	fmt.Println("matrix[1][2] =", matrix[1][2]) // 6
	fmt.Println("matrix[2][1] =", matrix[2][1]) // 8

	// === 2차원 배열 순회 ===
	fmt.Println("\n=== 행렬 출력 ===")
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Printf("%3d", matrix[i][j])
		}
		fmt.Println()
	}

	// range 사용
	fmt.Println("\n=== range로 순회 ===")
	for i, row := range matrix {
		for j, val := range row {
			fmt.Printf("matrix[%d][%d]=%d  ", i, j, val)
		}
		fmt.Println()
	}

	// === 행렬 덧셈 ===
	fmt.Println("\n=== 행렬 덧셈 ===")
	a := [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	b := [3][3]int{
		{9, 8, 7},
		{6, 5, 4},
		{3, 2, 1},
	}
	var result [3][3]int

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result[i][j] = a[i][j] + b[i][j]
		}
	}

	fmt.Println("A:")
	printMatrix(a)
	fmt.Println("+")
	fmt.Println("B:")
	printMatrix(b)
	fmt.Println("=")
	fmt.Println("결과:")
	printMatrix(result)

	// === 단위 행렬 (Identity Matrix) ===
	fmt.Println("=== 5x5 단위 행렬 ===")
	var identity [5][5]int
	for i := 0; i < 5; i++ {
		identity[i][i] = 1 // 대각선만 1
	}
	for _, row := range identity {
		for _, val := range row {
			fmt.Printf("%2d", val)
		}
		fmt.Println()
	}

	// === 2차원 배열 활용: 좌석 배치도 ===
	fmt.Println("\n=== 좌석 배치도 ===")
	// 0: 빈 좌석, 1: 예약됨
	seats := [4][6]int{
		{1, 1, 0, 0, 1, 1},
		{1, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 0, 0},
	}

	fmt.Println("     1열 2열 3열 4열 5열 6열")
	for i, row := range seats {
		fmt.Printf("%d행: ", i+1)
		for _, val := range row {
			if val == 1 {
				fmt.Print(" [X]") // 예약됨
			} else {
				fmt.Print(" [ ]") // 빈 좌석
			}
		}
		fmt.Println()
	}

	// 빈 좌석 수 세기
	empty := 0
	for _, row := range seats {
		for _, val := range row {
			if val == 0 {
				empty++
			}
		}
	}
	fmt.Printf("빈 좌석: %d / 전체: %d\n", empty, 4*6)
}

// 3x3 행렬을 출력하는 헬퍼 함수
func printMatrix(m [3][3]int) {
	for _, row := range m {
		for _, val := range row {
			fmt.Printf("%3d", val)
		}
		fmt.Println()
	}
}
