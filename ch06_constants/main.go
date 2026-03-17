package main

import "fmt"

// 단일 상수 선언
const Pi = 3.14159

// 여러 상수를 그룹으로 선언
const (
	AppName    = "Go 학습"
	AppVersion = "1.0.0"
	MaxRetry   = 3
)

// 타입 있는 상수
const TypedMax int32 = 100

func main() {
	// === 기본 상수 사용 ===
	fmt.Println("=== 기본 상수 ===")
	fmt.Println("Pi:", Pi)
	fmt.Println("앱 이름:", AppName)
	fmt.Println("앱 버전:", AppVersion)
	fmt.Println("최대 재시도:", MaxRetry)

	// === 상수로 원의 넓이 계산 ===
	fmt.Println("\n=== 원의 넓이 ===")
	radius := 5.0
	area := Pi * radius * radius
	fmt.Printf("반지름 %.1f인 원의 넓이: %.2f\n", radius, area)

	// === 타입 없는 상수의 유연함 ===
	fmt.Println("\n=== 타입 없는 상수 ===")
	const Big = 1000000

	var a int32 = Big   // int32로 사용
	var b int64 = Big   // int64로 사용
	var c float64 = Big // float64로 사용

	fmt.Printf("int32:   %d\n", a)
	fmt.Printf("int64:   %d\n", b)
	fmt.Printf("float64: %.1f\n", c)

	// === iota를 사용한 열거값 ===
	fmt.Println("\n=== iota 열거값 (요일) ===")
	const (
		Sunday    = iota // 0
		Monday           // 1
		Tuesday          // 2
		Wednesday        // 3
		Thursday         // 4
		Friday           // 5
		Saturday         // 6
	)
	fmt.Println("Sunday:", Sunday)
	fmt.Println("Monday:", Monday)
	fmt.Println("Friday:", Friday)
	fmt.Println("Saturday:", Saturday)

	// === iota + 1 (1부터 시작) ===
	fmt.Println("\n=== iota + 1 (월) ===")
	const (
		January  = iota + 1 // 1
		February            // 2
		March               // 3
	)
	fmt.Println("January:", January)
	fmt.Println("February:", February)
	fmt.Println("March:", March)
}
