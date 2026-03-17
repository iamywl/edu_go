// type_conversion.go
// 타입 변환(Type Conversion) 예제
// 실행 방법: go run type_conversion.go

package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("===== 숫자 타입 간 변환 =====")

	// int → float64
	var intVal int = 42
	var floatVal float64 = float64(intVal)
	fmt.Printf("int(%d) → float64(%f)\n", intVal, floatVal)

	// float64 → int (소수점 이하 버림)
	var pi float64 = 3.14159
	var piInt int = int(pi)
	fmt.Printf("float64(%f) → int(%d) [소수점 이하 버림!]\n", pi, piInt)

	// int → int8 (범위 초과 주의!)
	var bigNum int = 300
	var smallNum int8 = int8(bigNum)
	fmt.Printf("int(%d) → int8(%d) [오버플로우 발생!]\n", bigNum, smallNum)

	// int32 → int64 (Go에서는 같은 int 계열도 명시적 변환 필요)
	var i32 int32 = 100
	var i64 int64 = int64(i32)
	fmt.Printf("int32(%d) → int64(%d)\n", i32, i64)

	fmt.Println()
	fmt.Println("===== 정수와 실수의 연산 =====")

	// 서로 다른 타입끼리 직접 연산 불가
	a := 10  // int
	b := 3.0 // float64

	// result := a / b  // 컴파일 에러! int와 float64는 직접 연산 불가
	result := float64(a) / b // int를 float64로 변환 후 연산
	fmt.Printf("%d / %f = %f\n", a, b, result)

	// 정수 나눗셈 vs 실수 나눗셈
	fmt.Printf("정수 나눗셈: 10 / 3 = %d\n", 10/3)         // 3 (소수점 버림)
	fmt.Printf("실수 나눗셈: 10.0 / 3.0 = %f\n", 10.0/3.0) // 3.333333

	fmt.Println()
	fmt.Println("===== 문자열 ↔ 숫자 변환 (strconv 패키지) =====")

	// 정수 → 문자열
	num := 42
	str := strconv.Itoa(num) // Itoa = Integer to ASCII
	fmt.Printf("int(%d) → string(%q)\n", num, str)

	// 문자열 → 정수
	str2 := "123"
	num2, err := strconv.Atoi(str2) // Atoi = ASCII to Integer
	if err != nil {
		fmt.Println("변환 에러:", err)
	} else {
		fmt.Printf("string(%q) → int(%d)\n", str2, num2)
	}

	// 변환 실패 예시
	str3 := "hello"
	num3, err := strconv.Atoi(str3)
	if err != nil {
		fmt.Printf("string(%q) → int 변환 실패: %v\n", str3, err)
	} else {
		fmt.Printf("num3 = %d\n", num3)
	}

	// 실수 → 문자열
	f := 3.14159
	fStr := strconv.FormatFloat(f, 'f', 2, 64) // 소수점 2자리
	fmt.Printf("float64(%f) → string(%q)\n", f, fStr)

	// 문자열 → 실수
	fStr2 := "2.718"
	f2, err := strconv.ParseFloat(fStr2, 64)
	if err != nil {
		fmt.Println("변환 에러:", err)
	} else {
		fmt.Printf("string(%q) → float64(%f)\n", fStr2, f2)
	}

	// bool → 문자열
	bStr := strconv.FormatBool(true)
	fmt.Printf("bool(true) → string(%q)\n", bStr)

	// 문자열 → bool
	bVal, _ := strconv.ParseBool("true")
	fmt.Printf("string(\"true\") → bool(%t)\n", bVal)

	fmt.Println()
	fmt.Println("===== byte와 rune 변환 =====")

	// 문자 → 숫자
	var ch byte = 'A'
	fmt.Printf("문자 %c의 ASCII 코드: %d\n", ch, ch)

	// 숫자 → 문자
	var code byte = 65
	fmt.Printf("ASCII 코드 %d의 문자: %c\n", code, code)

	// 한글과 rune
	var hangul rune = '가'
	fmt.Printf("문자 %c의 유니코드: U+%04X (%d)\n", hangul, hangul, hangul)
}
