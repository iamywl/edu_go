// scope.go
// 변수의 범위(스코프) 예제
// 실행 방법: go run scope.go

package main

import "fmt"

// 패키지 레벨 변수: 이 패키지의 모든 함수에서 접근 가능
var globalMessage = "안녕하세요, 저는 패키지 레벨 변수입니다"

func main() {
	fmt.Println("===== 패키지 레벨 변수 =====")
	fmt.Println(globalMessage)
	printGlobal() // 다른 함수에서도 접근 가능

	fmt.Println()
	fmt.Println("===== 함수 레벨 변수 =====")
	localVar := "저는 main 함수의 지역 변수입니다"
	fmt.Println(localVar)
	// localVar는 main 함수 안에서만 접근 가능

	fmt.Println()
	fmt.Println("===== 블록 레벨 변수 =====")
	x := 10
	fmt.Printf("블록 바깥 x = %d\n", x)

	{
		// 새로운 블록 시작
		y := 20 // 이 블록 안에서만 유효
		fmt.Printf("  블록 안쪽 x = %d (바깥 변수에 접근 가능)\n", x)
		fmt.Printf("  블록 안쪽 y = %d\n", y)
	}
	// fmt.Println(y)  // 컴파일 에러: y는 블록 안에서만 유효

	fmt.Printf("블록 바깥 x = %d (변함 없음)\n", x)

	fmt.Println()
	fmt.Println("===== if 블록의 스코프 =====")
	score := 85

	if score >= 80 {
		grade := "A" // if 블록 안에서만 유효
		fmt.Printf("점수: %d, 등급: %s\n", score, grade)
	}
	// fmt.Println(grade)  // 컴파일 에러: grade는 if 블록 안에서만 유효

	fmt.Println()
	fmt.Println("===== for 블록의 스코프 =====")

	for i := 0; i < 3; i++ {
		// i는 for 블록 안에서만 유효
		fmt.Printf("  i = %d\n", i)
	}
	// fmt.Println(i)  // 컴파일 에러: i는 for 블록 안에서만 유효

	fmt.Println()
	fmt.Println("===== 섀도잉 (Shadowing) =====")

	value := 100
	fmt.Printf("바깥 value = %d\n", value)

	{
		// 같은 이름으로 새 변수를 선언하면 바깥 변수를 "가림(shadow)"
		value := 200 // 새로운 변수! 바깥의 value와 다른 변수
		fmt.Printf("  안쪽 value = %d (새로운 변수)\n", value)

		value = 300 // 안쪽 value의 값 변경
		fmt.Printf("  안쪽 value = %d (값 변경)\n", value)
	}

	// 바깥의 value는 영향 받지 않음
	fmt.Printf("바깥 value = %d (변함 없음!)\n", value)

	fmt.Println()
	fmt.Println("===== 섀도잉 주의사항 =====")

	// 섀도잉은 버그의 원인이 될 수 있으므로 주의!
	count := 0
	fmt.Printf("초기 count = %d\n", count)

	for i := 0; i < 5; i++ {
		// 의도: count를 증가시키고 싶음
		count := count + 1 // 실수! 새로운 지역 변수 count를 만듦 (섀도잉)
		_ = count          // 사용하지 않으면 컴파일 에러이므로 임시로 사용 표시
	}
	fmt.Printf("반복 후 count = %d (여전히 0! 섀도잉 때문)\n", count)

	// 올바른 방법: = 사용 (기존 변수에 대입)
	for i := 0; i < 5; i++ {
		count = count + 1 // = 사용: 기존 변수에 값 대입
	}
	fmt.Printf("올바른 반복 후 count = %d\n", count)
}

// 다른 함수에서 패키지 레벨 변수에 접근
func printGlobal() {
	fmt.Printf("  printGlobal()에서: %s\n", globalMessage)
}
