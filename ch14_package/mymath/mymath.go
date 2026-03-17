// Package mymath 는 기본적인 수학 연산을 제공하는 패키지이다.
// 대문자로 시작하는 함수/변수는 외부에 공개(Exported)됩니다.
package mymath

import (
	"errors"
	"fmt"
)

// Pi 는 원주율 상수 (대문자 → 외부 공개)
const Pi = 3.141592653589793

// maxValue 는 내부에서만 사용하는 상수 (소문자 → 비공개)
const maxValue = 1000000

// init 함수 — 패키지가 import될 때 자동 실행
func init() {
	fmt.Println("[mymath] init 실행: 수학 패키지 로드 완료")
}

// Add 는 두 정수의 합을 반환한다. (대문자 → 외부 공개)
func Add(a, b int) int {
	return a + b
}

// Subtract 는 두 정수의 차를 반환한다. (대문자 → 외부 공개)
func Subtract(a, b int) int {
	return a - b
}

// Multiply 는 두 정수의 곱을 반환한다. (대문자 → 외부 공개)
func Multiply(a, b int) int {
	return a * b
}

// Divide 는 두 정수의 나눗셈 결과를 반환한다.
// 0으로 나누면 에러를 반환한다. (대문자 → 외부 공개)
func Divide(a, b int) (float64, error) {
	if b == 0 {
		return 0, errors.New("0으로 나눌 수 없습니다")
	}
	if !validate(a) || !validate(b) {
		return 0, errors.New("값이 허용 범위를 초과합니다")
	}
	return float64(a) / float64(b), nil
}

// validate 는 값이 허용 범위 내인지 확인한다. (소문자 → 비공개)
// 이 함수는 패키지 외부에서 호출할 수 없습니다.
func validate(n int) bool {
	if n < 0 {
		n = -n
	}
	return n <= maxValue
}
