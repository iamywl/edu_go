// main.go
// 변수의 다양한 선언 방법과 기본 타입을 학습한다.
// 실행 방법: go run main.go

package main

import "fmt"

// 패키지 레벨 변수 (함수 바깥에서 선언)
// 짧은 선언(:=)은 여기서 사용할 수 없습니다.
var packageVar string = "패키지 레벨 변수"

// 여러 변수를 한 번에 선언하는 var 블록
var (
	appName    string  = "Go 학습"
	appVersion float64 = 1.0
)

// 상수 선언
const Pi = 3.14159265358979

// iota를 사용한 상수 열거
const (
	Spring = iota // 0
	Summer        // 1
	Autumn        // 2
	Winter        // 3
)

func main() {
	fmt.Println("===== 2.2 변수 선언 (var 키워드) =====")

	// 방법 1: 타입을 명시하여 선언 (기본값으로 초기화)
	var age int
	var name string
	var isStudent bool
	fmt.Printf("age = %d (int의 기본값)\n", age)
	fmt.Printf("name = %q (string의 기본값)\n", name)
	fmt.Printf("isStudent = %t (bool의 기본값)\n", isStudent)

	// 방법 2: 선언과 동시에 초기화
	var height float64 = 175.5
	fmt.Printf("height = %.1f\n", height)

	// 방법 3: 타입 추론 (값에서 타입을 자동으로 결정)
	var city = "서울" // string으로 추론
	fmt.Printf("city = %s (타입: %T)\n", city, city)

	fmt.Println()
	fmt.Println("===== 2.3 기본값 (Zero Value) =====")

	var (
		zeroInt     int
		zeroFloat   float64
		zeroBool    bool
		zeroString  string
		zeroPointer *int
	)
	fmt.Printf("int 기본값:     %d\n", zeroInt)
	fmt.Printf("float64 기본값: %f\n", zeroFloat)
	fmt.Printf("bool 기본값:    %t\n", zeroBool)
	fmt.Printf("string 기본값:  %q\n", zeroString)
	fmt.Printf("pointer 기본값: %v\n", zeroPointer)

	fmt.Println()
	fmt.Println("===== 2.4 짧은 선언 (:=) =====")

	// 짧은 선언: var 키워드와 타입 생략
	// 함수 안에서만 사용 가능!
	userName := "홍길동"
	userAge := 25
	userHeight := 175.5
	isActive := true
	fmt.Printf("이름: %s, 나이: %d, 키: %.1f, 활성: %t\n",
		userName, userAge, userHeight, isActive)

	// 여러 변수를 한 줄에 선언
	x, y, z := 1, "hello", true
	fmt.Printf("x=%d, y=%s, z=%t\n", x, y, z)

	// 같은 타입의 여러 변수
	var a, b, c int = 10, 20, 30
	fmt.Printf("a=%d, b=%d, c=%d\n", a, b, c)

	fmt.Println()
	fmt.Println("===== 타입 종류 확인 =====")

	// 정수 타입
	var i8 int8 = 127     // -128 ~ 127
	var i16 int16 = 32767 // -32768 ~ 32767
	var i32 int32 = 2147483647
	var i64 int64 = 9223372036854775807

	fmt.Printf("int8:  %d (타입: %T)\n", i8, i8)
	fmt.Printf("int16: %d (타입: %T)\n", i16, i16)
	fmt.Printf("int32: %d (타입: %T)\n", i32, i32)
	fmt.Printf("int64: %d (타입: %T)\n", i64, i64)

	// 부호 없는 정수
	var u8 uint8 = 255
	fmt.Printf("uint8: %d (타입: %T)\n", u8, u8)

	// byte와 rune
	var b1 byte = 'A' // uint8의 별칭
	var r1 rune = '가' // int32의 별칭 (유니코드 코드 포인트)
	fmt.Printf("byte: %c (%d)\n", b1, b1)
	fmt.Printf("rune: %c (%d)\n", r1, r1)

	fmt.Println()
	fmt.Println("===== 상수와 iota =====")

	fmt.Printf("Pi = %f\n", Pi)
	fmt.Printf("계절: 봄=%d, 여름=%d, 가을=%d, 겨울=%d\n",
		Spring, Summer, Autumn, Winter)

	fmt.Println()
	fmt.Println("===== 2.7 숫자 표현 =====")

	// 다양한 진법 표현
	decimal := 42      // 10진수
	binary := 0b101010 // 2진수
	octal := 0o52      // 8진수
	hex := 0x2A        // 16진수

	fmt.Printf("10진수: %d\n", decimal)
	fmt.Printf("2진수:  0b%b = %d\n", binary, binary)
	fmt.Printf("8진수:  0o%o = %d\n", octal, octal)
	fmt.Printf("16진수: 0x%X = %d\n", hex, hex)

	// 가독성을 위한 밑줄 사용
	billion := 1_000_000_000
	fmt.Printf("10억: %d\n", billion)
}
