package main

import (
	"fmt"
	"runtime"
)

func main() {
	// === 기본 switch문 ===
	fmt.Println("=== 기본 switch문 ===")
	day := "수요일"

	switch day {
	case "월요일":
		fmt.Println("한 주의 시작!")
	case "수요일":
		fmt.Println("한 주의 중간!")
	case "금요일":
		fmt.Println("불금!")
	default:
		fmt.Println(day, "입니다")
	}

	// === 한 case에 여러 값 ===
	fmt.Println("\n=== 여러 값 매칭 ===")
	today := "토요일"

	switch today {
	case "토요일", "일요일":
		fmt.Println(today, "→ 주말입니다! 쉬세요~")
	case "월요일", "화요일", "수요일", "목요일", "금요일":
		fmt.Println(today, "→ 평일입니다. 열심히!")
	}

	// === 계절 판별 ===
	fmt.Println("\n=== 계절 판별 ===")
	month := 7

	switch month {
	case 3, 4, 5:
		fmt.Printf("%d월 → 봄\n", month)
	case 6, 7, 8:
		fmt.Printf("%d월 → 여름\n", month)
	case 9, 10, 11:
		fmt.Printf("%d월 → 가을\n", month)
	case 12, 1, 2:
		fmt.Printf("%d월 → 겨울\n", month)
	default:
		fmt.Println("잘못된 월입니다")
	}

	// === switch 초기문 ===
	fmt.Println("\n=== switch 초기문 ===")

	// 운영체제 확인
	switch os := runtime.GOOS; os {
	case "linux":
		fmt.Println("리눅스를 사용 중입니다")
	case "darwin":
		fmt.Println("macOS를 사용 중입니다")
	case "windows":
		fmt.Println("Windows를 사용 중입니다")
	default:
		fmt.Printf("기타 OS: %s\n", os)
	}

	// === 조건 없는 switch (if-else 대용) ===
	fmt.Println("\n=== 조건 없는 switch ===")
	score := 85

	switch {
	case score >= 90:
		fmt.Printf("%d점 → A학점\n", score)
	case score >= 80:
		fmt.Printf("%d점 → B학점\n", score)
	case score >= 70:
		fmt.Printf("%d점 → C학점\n", score)
	case score >= 60:
		fmt.Printf("%d점 → D학점\n", score)
	default:
		fmt.Printf("%d점 → F학점\n", score)
	}

	// === BMI 판별 (조건 없는 switch 활용) ===
	fmt.Println("\n=== BMI 판별 ===")
	height := 1.75 // 미터
	weight := 70.0 // kg
	bmi := weight / (height * height)

	fmt.Printf("키: %.2fm, 몸무게: %.1fkg, BMI: %.1f\n", height, weight, bmi)

	switch {
	case bmi < 18.5:
		fmt.Println("결과: 저체중")
	case bmi < 25.0:
		fmt.Println("결과: 정상")
	case bmi < 30.0:
		fmt.Println("결과: 과체중")
	default:
		fmt.Println("결과: 비만")
	}
}
