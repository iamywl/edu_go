package main

import (
	"fmt"
	"math"
	"strconv"
)

// 나눗셈 함수 (에러 반환 패턴)
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("0으로 나눌 수 없습니다")
	}
	return a / b, nil
}

func main() {
	// === if 초기문; 조건문 기본 ===
	fmt.Println("=== if 초기문; 조건문 ===")

	// 초기문에서 변수를 선언하고 바로 조건 검사
	if num := 10; num%2 == 0 {
		fmt.Printf("%d는 짝수입니다\n", num)
	} else {
		fmt.Printf("%d는 홀수입니다\n", num)
	}
	// 여기서 num은 사용 불가 (if 블록 스코프에만 존재)

	// === 문자열 길이 검사 ===
	fmt.Println("\n=== 문자열 길이 검사 ===")
	name := "Go 프로그래밍"

	if length := len(name); length >= 5 {
		fmt.Printf("'%s' (길이 %d) → 긴 문자열\n", name, length)
	} else {
		fmt.Printf("'%s' (길이 %d) → 짧은 문자열\n", name, length)
	}

	// === 함수 호출과 에러 처리 ===
	fmt.Println("\n=== 함수 호출 + 에러 처리 ===")

	// 정상적인 나눗셈
	if result, err := divide(10, 3); err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	// 0으로 나누기 (에러 발생)
	if result, err := divide(10, 0); err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("10 / 0 = %.2f\n", result)
	}

	// === 문자열 → 숫자 변환 ===
	fmt.Println("\n=== 문자열 변환 ===")

	// 성공하는 변환
	if n, err := strconv.Atoi("42"); err == nil {
		fmt.Printf("\"42\" → %d (성공)\n", n)
	}

	// 실패하는 변환
	if n, err := strconv.Atoi("abc"); err != nil {
		fmt.Printf("\"abc\" 변환 실패: %v\n", err)
	} else {
		fmt.Printf("\"abc\" → %d\n", n)
	}

	// === 수학 함수와 함께 ===
	fmt.Println("\n=== 수학 함수와 함께 ===")

	if abs := math.Abs(-7.5); abs > 5 {
		fmt.Printf("절댓값 %.1f은 5보다 큽니다\n", abs)
	}

	// === 실전 패턴: 맵에서 값 찾기 ===
	fmt.Println("\n=== 맵에서 값 찾기 ===")
	capitals := map[string]string{
		"한국": "서울",
		"일본": "도쿄",
		"미국": "워싱턴",
	}

	// if 초기문 패턴으로 맵 조회
	if capital, ok := capitals["한국"]; ok {
		fmt.Printf("한국의 수도: %s\n", capital)
	}

	if capital, ok := capitals["프랑스"]; ok {
		fmt.Printf("프랑스의 수도: %s\n", capital)
	} else {
		fmt.Println("프랑스의 수도 정보가 없습니다")
		_ = capital // 미사용 변수 경고 방지
	}
}
