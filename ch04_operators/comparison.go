// comparison.go
// 비교 연산자 예제
// 실행 방법: go run comparison.go

package main

import "fmt"

func main() {
	fmt.Println("===== 정수 비교 =====")

	a, b := 10, 20

	fmt.Printf("%d == %d : %t\n", a, b, a == b) // false
	fmt.Printf("%d != %d : %t\n", a, b, a != b) // true
	fmt.Printf("%d < %d  : %t\n", a, b, a < b)  // true
	fmt.Printf("%d > %d  : %t\n", a, b, a > b)  // false
	fmt.Printf("%d <= %d : %t\n", a, b, a <= b) // true
	fmt.Printf("%d >= %d : %t\n", a, b, a >= b) // false

	fmt.Println()
	fmt.Println("===== 같은 값 비교 =====")

	x, y := 10, 10
	fmt.Printf("%d == %d : %t\n", x, y, x == y) // true
	fmt.Printf("%d <= %d : %t\n", x, y, x <= y) // true
	fmt.Printf("%d >= %d : %t\n", x, y, x >= y) // true

	fmt.Println()
	fmt.Println("===== 문자열 비교 =====")

	// 문자열은 사전순(lexicographic)으로 비교
	fmt.Printf("%q == %q : %t\n", "Go", "Go", "Go" == "Go")              // true
	fmt.Printf("%q == %q : %t\n", "Go", "go", "Go" == "go")              // false (대소문자 구분)
	fmt.Printf("%q < %q  : %t\n", "apple", "banana", "apple" < "banana") // true
	fmt.Printf("%q < %q  : %t\n", "ABC", "abc", "ABC" < "abc")           // true (대문자가 더 작음)

	fmt.Println()
	fmt.Println("===== bool 비교 =====")

	t, f := true, false
	fmt.Printf("%t == %t : %t\n", t, t, t == t) // true
	fmt.Printf("%t == %t : %t\n", t, f, t == f) // false
	fmt.Printf("%t != %t : %t\n", t, f, t != f) // true

	fmt.Println()
	fmt.Println("===== 비교 연산자 활용 예제 =====")

	// 점수에 따른 등급 판정
	score := 85

	fmt.Printf("점수: %d\n", score)
	fmt.Printf("90점 이상인가? %t\n", score >= 90)
	fmt.Printf("80점 이상인가? %t\n", score >= 80)
	fmt.Printf("70점 이상인가? %t\n", score >= 70)

	// 비밀번호 확인
	password := "go1234"
	input := "go1234"

	if password == input {
		fmt.Println("\n비밀번호가 일치합니다!")
	} else {
		fmt.Println("\n비밀번호가 틀렸습니다.")
	}

	fmt.Println()
	fmt.Println("===== 주의: 서로 다른 타입은 비교 불가 =====")

	// Go는 서로 다른 타입끼리 비교할 수 없음
	var i32 int32 = 10
	var i64 int64 = 10

	// fmt.Println(i32 == i64) // 컴파일 에러! int32와 int64는 비교 불가
	fmt.Printf("int32(%d) == int64(%d) : %t\n", i32, i64, int64(i32) == i64) // 타입 변환 후 비교

	// 정수와 실수 비교도 불가
	var intVal int = 10
	var floatVal float64 = 10.0

	// fmt.Println(intVal == floatVal) // 컴파일 에러!
	fmt.Printf("int(%d) == float64(%.1f) : %t\n",
		intVal, floatVal, float64(intVal) == floatVal) // 타입 변환 후 비교
}
