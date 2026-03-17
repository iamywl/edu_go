package main

import (
	"errors"
	"fmt"
)

// ============================================================
// 21.1 에러 반환 패턴
// Go에서는 함수의 마지막 반환값으로 error를 반환하는 것이 관례이다.
// ============================================================

// errors.New를 사용한 센티널 에러 정의
var ErrDivisionByZero = errors.New("0으로 나눌 수 없습니다")

// divide - 나눗셈 함수. 에러를 마지막 반환값으로 반환한다.
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// parseAge - fmt.Errorf를 사용하여 동적 에러 메시지를 생성한다.
func parseAge(age int) (string, error) {
	if age < 0 {
		return "", fmt.Errorf("유효하지 않은 나이: %d (음수)", age)
	}
	if age > 150 {
		return "", fmt.Errorf("유효하지 않은 나이: %d (너무 큼)", age)
	}
	if age >= 18 {
		return "성인", nil
	}
	return "미성년자", nil
}

func main() {
	fmt.Println("=== errors.New 사용 예제 ===")

	// 정상적인 나눗셈
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	// 0으로 나누기 시도
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("10 / 0 = %.2f\n", result)
	}

	// 센티널 에러 비교
	if errors.Is(err, ErrDivisionByZero) {
		fmt.Println("-> 이것은 0으로 나누기 에러입니다")
	}

	fmt.Println("\n=== fmt.Errorf 사용 예제 ===")

	// 정상적인 나이
	category, err := parseAge(25)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Println("25세:", category)
	}

	// 음수 나이
	category, err = parseAge(-5)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Println("-5세:", category)
	}

	// 너무 큰 나이
	category, err = parseAge(200)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Println("200세:", category)
	}

	fmt.Println("\n=== 에러 무시의 위험성 ===")

	// 나쁜 예: 에러를 무시하면 안 됩니다!
	// result, _ := divide(10, 0)  // 이렇게 하면 에러를 놓칠 수 있습니다

	// 좋은 예: 항상 에러를 확인하세요
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("에러를 확인하고 적절히 처리합니다:", err)
	}
	_ = result // 에러가 발생했으므로 result를 사용하지 않습니다
}
