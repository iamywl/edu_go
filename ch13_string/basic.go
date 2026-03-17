package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	// === 문자열 기본 ===
	fmt.Println("=== 문자열 기본 ===")

	s := "안녕하세요, Go!"
	fmt.Printf("문자열: %s\n", s)
	fmt.Printf("바이트 수 (len): %d\n", len(s))                         // 19
	fmt.Printf("글자 수 (RuneCount): %d\n", utf8.RuneCountInString(s)) // 9

	fmt.Println()

	// === 일반 문자열 vs Raw 문자열 ===
	fmt.Println("=== 일반 vs Raw 문자열 ===")

	normal := "줄바꿈: \n탭: \t끝"
	raw := `줄바꿈: \n탭: \t끝`

	fmt.Println("일반 문자열:")
	fmt.Println(normal)
	fmt.Println()
	fmt.Println("Raw 문자열:")
	fmt.Println(raw)

	fmt.Println()

	// === 문자열 인덱싱 ===
	fmt.Println("=== 문자열 인덱싱 (바이트 접근) ===")
	eng := "Hello"
	fmt.Printf("eng[0] = %c (바이트: %d)\n", eng[0], eng[0]) // H
	fmt.Printf("eng[4] = %c (바이트: %d)\n", eng[4], eng[4]) // o

	// 한글은 인덱싱으로 접근하면 깨짐
	kor := "한글"
	fmt.Printf("kor[0] = %d (한글의 첫 번째 바이트, 글자가 아님!)\n", kor[0])

	fmt.Println()

	// === 유용한 strings 패키지 함수들 ===
	fmt.Println("=== strings 패키지 ===")

	text := "Hello, Go World! Go is great."

	fmt.Printf("Contains(\"Go\"): %v\n", strings.Contains(text, "Go"))
	fmt.Printf("Count(\"Go\"): %d\n", strings.Count(text, "Go"))
	fmt.Printf("HasPrefix(\"Hello\"): %v\n", strings.HasPrefix(text, "Hello"))
	fmt.Printf("HasSuffix(\"great.\"): %v\n", strings.HasSuffix(text, "great."))
	fmt.Printf("Index(\"Go\"): %d\n", strings.Index(text, "Go"))
	fmt.Printf("ToUpper: %s\n", strings.ToUpper(text))
	fmt.Printf("ToLower: %s\n", strings.ToLower(text))
	fmt.Printf("Replace: %s\n", strings.Replace(text, "Go", "Golang", -1))
	fmt.Printf("TrimSpace: [%s]\n", strings.TrimSpace("  공백 제거  "))

	// Split과 Join
	csv := "사과,바나나,포도,딸기"
	fruits := strings.Split(csv, ",")
	fmt.Printf("Split: %v\n", fruits)
	fmt.Printf("Join:  %s\n", strings.Join(fruits, " | "))
}
