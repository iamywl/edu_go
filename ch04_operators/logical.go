// logical.go
// 논리 연산자 예제
// 실행 방법: go run logical.go

package main

import "fmt"

func main() {
	fmt.Println("===== 논리 연산자 기본 =====")

	// && (AND): 둘 다 참이어야 참
	fmt.Println("--- AND (&&) ---")
	fmt.Printf("true  && true  = %t\n", true && true)
	fmt.Printf("true  && false = %t\n", true && false)
	fmt.Printf("false && true  = %t\n", false && true)
	fmt.Printf("false && false = %t\n", false && false)

	fmt.Println()

	// || (OR): 하나라도 참이면 참
	fmt.Println("--- OR (||) ---")
	fmt.Printf("true  || true  = %t\n", true || true)
	fmt.Printf("true  || false = %t\n", true || false)
	fmt.Printf("false || true  = %t\n", false || true)
	fmt.Printf("false || false = %t\n", false || false)

	fmt.Println()

	// ! (NOT): 참과 거짓을 반전
	fmt.Println("--- NOT (!) ---")
	fmt.Printf("!true  = %t\n", !true)
	fmt.Printf("!false = %t\n", !false)

	fmt.Println()
	fmt.Println("===== 비교 연산자와 조합 =====")

	age := 25
	hasLicense := true
	hasCar := false

	// 20세 이상이고 면허가 있어야 운전 가능
	canDrive := age >= 20 && hasLicense
	fmt.Printf("나이: %d, 면허: %t → 운전 가능: %t\n", age, hasLicense, canDrive)

	// 차가 있거나 면허가 있으면 교통 관련 정보 표시
	showTraffic := hasCar || hasLicense
	fmt.Printf("차: %t, 면허: %t → 교통 정보 표시: %t\n", hasCar, hasLicense, showTraffic)

	// 미성년자가 아닌 경우
	isMinor := age < 19
	isAdult := !isMinor
	fmt.Printf("나이: %d, 미성년: %t, 성인: %t\n", age, isMinor, isAdult)

	fmt.Println()
	fmt.Println("===== 단축 평가 (Short-circuit Evaluation) =====")

	// && : 왼쪽이 false이면 오른쪽은 평가하지 않음
	fmt.Println("--- && 단축 평가 ---")
	fmt.Print("false && check(): ")
	result := false && check("이 함수는 호출되지 않음")
	fmt.Printf("결과: %t\n", result)

	fmt.Print("true && check():  ")
	result = true && check("이 함수는 호출됨")
	fmt.Printf("결과: %t\n", result)

	fmt.Println()

	// || : 왼쪽이 true이면 오른쪽은 평가하지 않음
	fmt.Println("--- || 단축 평가 ---")
	fmt.Print("true || check():  ")
	result = true || check("이 함수는 호출되지 않음")
	fmt.Printf("결과: %t\n", result)

	fmt.Print("false || check(): ")
	result = false || check("이 함수는 호출됨")
	fmt.Printf("결과: %t\n", result)

	fmt.Println()
	fmt.Println("===== 실용 예제: 로그인 시스템 =====")

	username := "admin"
	password := "1234"
	isActive := true

	// 사용자 이름과 비밀번호가 맞고, 계정이 활성화 상태여야 로그인 가능
	inputUser := "admin"
	inputPass := "1234"

	isValidUser := inputUser == username && inputPass == password
	canLogin := isValidUser && isActive

	fmt.Printf("입력 - 이름: %s, 비밀번호: %s\n", inputUser, inputPass)
	fmt.Printf("인증 성공: %t\n", isValidUser)
	fmt.Printf("로그인 가능: %t\n", canLogin)

	fmt.Println()
	fmt.Println("===== 실용 예제: 윤년 판별 =====")

	year := 2024

	// 윤년 조건:
	// 4로 나누어 떨어지고, 100으로 나누어 떨어지지 않음
	// 또는 400으로 나누어 떨어짐
	isLeapYear := (year%4 == 0 && year%100 != 0) || (year%400 == 0)

	fmt.Printf("%d년은 윤년인가? %t\n", year, isLeapYear)

	// 여러 해 확인
	years := []int{2000, 1900, 2024, 2023, 2100}
	for _, y := range years {
		leap := (y%4 == 0 && y%100 != 0) || (y%400 == 0)
		fmt.Printf("  %d년: 윤년=%t\n", y, leap)
	}
}

// check: 단축 평가를 확인하기 위한 함수
// 호출되면 메시지를 출력하고 true를 반환
func check(msg string) bool {
	fmt.Printf("[%s] ", msg)
	return true
}
