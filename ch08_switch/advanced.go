package main

import "fmt"

func main() {
	// === fallthrough 키워드 ===
	fmt.Println("=== fallthrough ===")

	// fallthrough는 다음 case를 조건 검사 없이 무조건 실행
	num := 3
	fmt.Printf("num = %d일 때:\n", num)

	switch {
	case num >= 3:
		fmt.Println("  3 이상")
		fallthrough
	case num >= 2:
		fmt.Println("  2 이상")
		fallthrough
	case num >= 1:
		fmt.Println("  1 이상")
	}

	// === fallthrough 주의사항 ===
	fmt.Println("\n=== fallthrough 주의 ===")
	// fallthrough는 다음 case의 조건을 무시하고 실행
	value := 5
	switch {
	case value > 10:
		fmt.Println("10 초과")
		fallthrough
	case value > 0:
		fmt.Println("양수")
		fallthrough
	case value < 0: // fallthrough 때문에 value=5인데도 실행됨!
		fmt.Println("이것도 실행됨 (fallthrough 때문)")
	}

	// === case 안에서 break ===
	fmt.Println("\n=== case 안에서 break ===")
	grade := "B"
	extra := true

	switch grade {
	case "A":
		fmt.Println("장학금 대상")
	case "B":
		if !extra {
			break // 조건에 따라 case를 조기 종료
		}
		fmt.Println("B학점이지만 추가 활동으로 우수 학생!")
	case "C":
		fmt.Println("보통")
	}

	// === 조건 없는 switch로 범위 판별 ===
	fmt.Println("\n=== 온도 판별 ===")
	temperatures := []int{-5, 5, 15, 25, 35}

	for _, temp := range temperatures {
		fmt.Printf("%3d도 → ", temp)
		switch {
		case temp < 0:
			fmt.Println("영하! 매우 추움")
		case temp < 10:
			fmt.Println("추움, 겉옷 필요")
		case temp < 20:
			fmt.Println("선선함")
		case temp < 30:
			fmt.Println("따뜻함, 쾌적")
		default:
			fmt.Println("더움! 시원하게")
		}
	}

	// === 문자열 분류 ===
	fmt.Println("\n=== 문자 분류 ===")
	chars := []rune{'A', 'z', '5', '!', '가'}

	for _, ch := range chars {
		fmt.Printf("'%c' → ", ch)
		switch {
		case ch >= 'A' && ch <= 'Z':
			fmt.Println("영문 대문자")
		case ch >= 'a' && ch <= 'z':
			fmt.Println("영문 소문자")
		case ch >= '0' && ch <= '9':
			fmt.Println("숫자")
		case ch >= 0xAC00 && ch <= 0xD7A3:
			fmt.Println("한글")
		default:
			fmt.Println("기타 문자")
		}
	}
}
