// print_examples.go
// fmt 패키지의 다양한 출력 함수 사용 예제
// 실행 방법: go run print_examples.go

package main

import "fmt"

func main() {
	fmt.Println("===== fmt.Println =====")
	// Println: 값을 출력하고 줄바꿈
	fmt.Println("Hello, Go!")
	fmt.Println(42)
	fmt.Println(3.14)
	fmt.Println(true)
	fmt.Println("이름:", "홍길동", "나이:", 25) // 여러 값은 공백으로 구분
	fmt.Println()                        // 빈 줄 출력

	fmt.Println("===== fmt.Print =====")
	// Print: 줄바꿈 없이 출력
	fmt.Print("Hello, ")
	fmt.Print("Go!")
	fmt.Print("\n") // 직접 줄바꿈

	// Print에 여러 값을 전달하면 공백 없이 이어붙임
	fmt.Print("A", "B", "C", "\n") // ABC

	fmt.Println()
	fmt.Println("===== fmt.Printf =====")

	name := "홍길동"
	age := 25
	height := 175.5
	isStudent := true

	// 서식 지정자를 사용한 출력
	fmt.Printf("이름: %s\n", name)
	fmt.Printf("나이: %d세\n", age)
	fmt.Printf("키: %.1fcm\n", height)
	fmt.Printf("학생 여부: %t\n", isStudent)

	// 한 줄에 여러 값
	fmt.Printf("이름: %s, 나이: %d세, 키: %.1fcm\n", name, age, height)

	fmt.Println()
	fmt.Println("===== %v (기본 형식) =====")
	// %v는 값의 타입에 맞는 기본 형식으로 출력
	fmt.Printf("문자열: %v\n", "Go")
	fmt.Printf("정수:   %v\n", 42)
	fmt.Printf("실수:   %v\n", 3.14)
	fmt.Printf("불리언: %v\n", true)

	fmt.Println()
	fmt.Println("===== %T (타입 출력) =====")
	fmt.Printf("%v의 타입: %T\n", "Go", "Go")
	fmt.Printf("%v의 타입: %T\n", 42, 42)
	fmt.Printf("%v의 타입: %T\n", 3.14, 3.14)
	fmt.Printf("%v의 타입: %T\n", true, true)

	fmt.Println()
	fmt.Println("===== Sprintf (문자열 생성) =====")
	// Sprintf: 화면에 출력하지 않고 문자열로 반환
	msg := fmt.Sprintf("이름: %s, 나이: %d", name, age)
	fmt.Println("생성된 문자열:", msg)

	// 날짜 형식 문자열 만들기
	year, month, day := 2024, 1, 1
	dateStr := fmt.Sprintf("%04d년 %02d월 %02d일", year, month, day)
	fmt.Println("날짜:", dateStr)

	fmt.Println()
	fmt.Println("===== 이스케이프 시퀀스 =====")
	fmt.Println("줄바꿈: 첫째 줄\n둘째 줄")
	fmt.Println("탭:     이름\t나이\t도시")
	fmt.Println("따옴표: \"큰따옴표\"")
	fmt.Println("역슬래시: \\")
	fmt.Println("퍼센트: fmt.Printf(\"100%%\")")
}
