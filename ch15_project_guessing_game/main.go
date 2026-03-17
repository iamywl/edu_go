package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// === 게임 설정 ===
	const (
		minNum     = 1   // 최소 숫자
		maxNum     = 100 // 최대 숫자
		maxAttempt = 10  // 최대 시도 횟수
	)

	// === 랜덤 숫자 생성 (1~100) ===
	answer := rand.Intn(maxNum-minNum+1) + minNum

	// === 게임 시작 안내 ===
	fmt.Println("===================================")
	fmt.Println("  숫자 맞추기 게임에 오신 것을 환영합니다!")
	fmt.Println("===================================")
	fmt.Printf("%d~%d 사이의 숫자를 맞춰보세요.\n", minNum, maxNum)
	fmt.Printf("기회는 %d번입니다.\n\n", maxAttempt)

	// === 게임 루프 ===
	for attempt := 1; attempt <= maxAttempt; attempt++ {
		// 현재 시도 횟수 표시 및 입력 받기
		var guess int
		fmt.Printf("[시도 %d/%d] 숫자를 입력하세요: ", attempt, maxAttempt)

		// 사용자 입력 받기
		_, err := fmt.Scan(&guess)
		if err != nil {
			// 잘못된 입력 처리 (문자 등)
			fmt.Println("올바른 숫자를 입력해주세요!")
			// 입력 버퍼 비우기
			var discard string
			fmt.Scanln(&discard)
			attempt-- // 잘못된 입력은 시도 횟수에 포함하지 않음
			continue
		}

		// 범위 검증
		if guess < minNum || guess > maxNum {
			fmt.Printf("%d~%d 사이의 숫자를 입력해주세요!\n\n", minNum, maxNum)
			attempt-- // 범위 밖 입력은 시도 횟수에 포함하지 않음
			continue
		}

		// === 숫자 비교 ===
		if guess > answer {
			fmt.Println("더 작은 수입니다!")
			// 5회 이상 실패 시 범위 힌트 제공
			if attempt >= 5 {
				printHint(answer, guess)
			}
		} else if guess < answer {
			fmt.Println("더 큰 수입니다!")
			// 5회 이상 실패 시 범위 힌트 제공
			if attempt >= 5 {
				printHint(answer, guess)
			}
		} else {
			// 정답!
			fmt.Println()
			fmt.Println("*****************************")
			fmt.Printf("  정답입니다! %d번 만에 맞추셨습니다!\n", attempt)
			fmt.Println("*****************************")

			// 시도 횟수에 따른 평가
			switch {
			case attempt <= 3:
				fmt.Println("  대단합니다! 천재시군요!")
			case attempt <= 5:
				fmt.Println("  훌륭합니다! 감이 좋으시네요!")
			case attempt <= 7:
				fmt.Println("  좋습니다! 잘 하셨어요!")
			default:
				fmt.Println("  아슬아슬했지만 성공!")
			}
			return // 게임 종료
		}
		fmt.Println()
	}

	// === 모든 기회를 소진한 경우 ===
	fmt.Println()
	fmt.Println("===================================")
	fmt.Printf("  아쉽습니다! 정답은 %d이었습니다.\n", answer)
	fmt.Println("  다음에 다시 도전해보세요!")
	fmt.Println("===================================")
}

// printHint 는 정답에 가까워지도록 범위 힌트를 제공하는 함수
func printHint(answer, guess int) {
	diff := answer - guess
	if diff < 0 {
		diff = -diff
	}

	switch {
	case diff <= 5:
		fmt.Println("  [힌트] 아주 가까워요! (차이: 5 이내)")
	case diff <= 15:
		fmt.Println("  [힌트] 거의 다 왔어요! (차이: 15 이내)")
	case diff <= 30:
		fmt.Println("  [힌트] 조금 더 노력해보세요. (차이: 30 이내)")
	default:
		fmt.Println("  [힌트] 아직 멀었어요...")
	}
}
