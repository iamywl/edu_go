package main

import (
	"fmt"

	// 로컬 패키지 import (모듈 경로 기준)
	"ch14_package/greeting"
	"ch14_package/mymath"
)

// main 패키지의 init 함수 — greeting과 mymath의 init 이후에 실행됨
func init() {
	fmt.Println("[main] init 실행")
}

func main() {
	fmt.Println("\n=== main() 함수 시작 ===")
	fmt.Println()

	// === mymath 패키지 사용 ===
	fmt.Println("=== mymath 패키지 ===")
	fmt.Printf("Add(10, 3)      = %d\n", mymath.Add(10, 3))
	fmt.Printf("Subtract(10, 3) = %d\n", mymath.Subtract(10, 3))
	fmt.Printf("Multiply(10, 3) = %d\n", mymath.Multiply(10, 3))

	result, err := mymath.Divide(10, 3)
	if err != nil {
		fmt.Printf("Divide 에러: %v\n", err)
	} else {
		fmt.Printf("Divide(10, 3)   = %.2f\n", result)
	}

	// 0으로 나누기 시도
	_, err = mymath.Divide(10, 0)
	if err != nil {
		fmt.Printf("Divide(10, 0)   = 에러: %v\n", err)
	}

	// 공개 상수 사용
	fmt.Printf("Pi              = %f\n", mymath.Pi)

	// 비공개 함수는 접근 불가
	// mymath.validate(5)  // 컴파일 에러!

	fmt.Println()

	// === greeting 패키지 사용 ===
	fmt.Println("=== greeting 패키지 ===")
	fmt.Printf("기본 언어: %s\n", greeting.DefaultLang)
	fmt.Println(greeting.Hello("홍길동"))
	fmt.Println(greeting.Hello("Kim"))

	greeting.SetLang("en")
	fmt.Printf("언어 변경 후: %s\n", greeting.DefaultLang)
	fmt.Println(greeting.Hello("Kim"))
}
