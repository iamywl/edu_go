package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// === byte vs rune 순회 비교 ===
	s := "Go 한글"

	fmt.Println("=== byte 단위 순회 ===")
	for i := 0; i < len(s); i++ {
		fmt.Printf("인덱스 %d: 바이트 0x%02x\n", i, s[i])
	}
	// 한글이 3바이트씩 출력됨 — 글자로 인식 불가

	fmt.Println()

	fmt.Println("=== rune 단위 순회 (for range) ===")
	for i, r := range s {
		fmt.Printf("바이트위치 %d: '%c' (U+%04X, 크기: %d바이트)\n",
			i, r, r, utf8.RuneLen(r))
	}
	// 한글은 3바이트씩 건너뜀 — 글자 단위로 올바르게 처리

	fmt.Println()

	// === 한글 문자열 처리 ===
	fmt.Println("=== 한글 문자열 처리 ===")
	korean := "대한민국"

	fmt.Printf("문자열: %s\n", korean)
	fmt.Printf("바이트 수: %d\n", len(korean))                     // 12
	fmt.Printf("글자 수:   %d\n", utf8.RuneCountInString(korean)) // 4

	// []rune으로 변환하면 인덱스로 개별 글자 접근 가능
	runes := []rune(korean)
	fmt.Printf("첫 글자: %c\n", runes[0])            // 대
	fmt.Printf("마지막:  %c\n", runes[len(runes)-1]) // 국

	fmt.Println()

	// === rune을 이용한 문자열 뒤집기 ===
	fmt.Println("=== 문자열 뒤집기 ===")
	original := "Hello 세계"
	reversed := reverseString(original)
	fmt.Printf("원본:   %s\n", original)
	fmt.Printf("뒤집기: %s\n", reversed)

	fmt.Println()

	// === 한글 자모 분리 확인 ===
	fmt.Println("=== 유니코드 코드포인트 확인 ===")
	sample := "가나다ABC"
	for _, r := range sample {
		fmt.Printf("'%c' → U+%04X (10진수: %d)\n", r, r, r)
	}

	fmt.Println()

	// === rune 타입 변환 ===
	fmt.Println("=== rune과 int32 ===")
	var r rune = '가'
	fmt.Printf("'%c'의 유니코드: %d (0x%X)\n", r, r, r)
	fmt.Printf("rune은 int32의 별칭: %T\n", r) // int32

	// 숫자를 rune으로
	nextChar := rune(r + 1)
	fmt.Printf("'가' 다음 유니코드: '%c' (U+%04X)\n", nextChar, nextChar)
}

// reverseString 은 rune 단위로 문자열을 뒤집는 함수
// 한글, 이모지 등 멀티바이트 문자도 올바르게 처리한다
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
