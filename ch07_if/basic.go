package main

import "fmt"

func main() {
	// === 기본 if문 ===
	fmt.Println("=== 기본 if문 ===")
	age := 20

	if age >= 18 {
		fmt.Println("성인입니다")
	}

	// === if-else ===
	fmt.Println("\n=== if-else ===")
	score := 75

	if score >= 60 {
		fmt.Println(score, "점: 합격!")
	} else {
		fmt.Println(score, "점: 불합격")
	}

	// === if-else if-else (학점 계산) ===
	fmt.Println("\n=== 학점 계산 ===")
	scores := []int{95, 82, 73, 64, 55}

	for _, s := range scores {
		var grade string
		if s >= 90 {
			grade = "A"
		} else if s >= 80 {
			grade = "B"
		} else if s >= 70 {
			grade = "C"
		} else if s >= 60 {
			grade = "D"
		} else {
			grade = "F"
		}
		fmt.Printf("%d점 → %s\n", s, grade)
	}

	// === 논리 연산자 &&, ||, ! ===
	fmt.Println("\n=== 논리 연산자 ===")
	temp := 25
	humid := 60

	// AND: 모든 조건이 참이어야 참
	if temp >= 20 && temp <= 28 && humid < 70 {
		fmt.Printf("온도 %d도, 습도 %d%% → 쾌적합니다\n", temp, humid)
	}

	// OR: 하나라도 참이면 참
	day := "토요일"
	if day == "토요일" || day == "일요일" {
		fmt.Println(day, "→ 주말입니다!")
	}

	// NOT: 조건 반전
	isRaining := false
	if !isRaining {
		fmt.Println("비가 오지 않습니다. 우산 불필요!")
	}

	// === 중첩 if ===
	fmt.Println("\n=== 중첩 if ===")
	userAge := 20
	hasTicket := true

	if userAge >= 18 {
		if hasTicket {
			fmt.Println("입장 가능합니다")
		} else {
			fmt.Println("티켓을 먼저 구매하세요")
		}
	} else {
		fmt.Println("미성년자는 입장할 수 없습니다")
	}

	// === 윤년 판별 ===
	fmt.Println("\n=== 윤년 판별 ===")
	year := 2024

	// 윤년 조건: 4의 배수이면서 100의 배수가 아닌 해, 또는 400의 배수인 해
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		fmt.Printf("%d년은 윤년입니다\n", year)
	} else {
		fmt.Printf("%d년은 평년입니다\n", year)
	}
}
