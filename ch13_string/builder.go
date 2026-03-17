package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	const count = 100000

	// === 방법 1: + 연산자로 합치기 ===
	fmt.Println("=== 문자열 합치기 성능 비교 ===")
	fmt.Printf("'Go'를 %d번 합치기\n\n", count)

	start := time.Now()
	result1 := ""
	for i := 0; i < count; i++ {
		result1 += "Go" // 매번 새로운 문자열 생성!
	}
	elapsed1 := time.Since(start)
	fmt.Printf("+ 연산자:        %v (길이: %d)\n", elapsed1, len(result1))

	// === 방법 2: strings.Builder 사용 ===
	start = time.Now()
	var builder strings.Builder
	for i := 0; i < count; i++ {
		builder.WriteString("Go") // 내부 버퍼에 추가
	}
	result2 := builder.String()
	elapsed2 := time.Since(start)
	fmt.Printf("strings.Builder: %v (길이: %d)\n", elapsed2, len(result2))

	// === 방법 3: strings.Builder + Grow (최적) ===
	start = time.Now()
	var builder2 strings.Builder
	builder2.Grow(count * 2) // 미리 용량 확보 (2바이트 × count)
	for i := 0; i < count; i++ {
		builder2.WriteString("Go")
	}
	result3 := builder2.String()
	elapsed3 := time.Since(start)
	fmt.Printf("Builder + Grow:  %v (길이: %d)\n", elapsed3, len(result3))

	// === 방법 4: strings.Join ===
	start = time.Now()
	parts := make([]string, count)
	for i := 0; i < count; i++ {
		parts[i] = "Go"
	}
	result4 := strings.Join(parts, "")
	elapsed4 := time.Since(start)
	fmt.Printf("strings.Join:    %v (길이: %d)\n", elapsed4, len(result4))

	fmt.Println()

	// === 결과 비교 ===
	fmt.Println("=== 결과 요약 ===")
	fmt.Println("+ 연산자는 매번 새로운 문자열을 생성하므로 O(n^2) 복잡도입니다.")
	fmt.Println("strings.Builder는 내부 버퍼를 재사용하므로 O(n) 복잡도입니다.")
	fmt.Println("Grow()로 미리 용량을 확보하면 재할당을 방지하여 더 빠릅니다.")

	fmt.Println()

	// === strings.Builder 다양한 메서드 ===
	fmt.Println("=== strings.Builder 메서드들 ===")
	var sb strings.Builder
	sb.WriteString("이름: ") // 문자열 추가
	sb.WriteString("홍길동")
	sb.WriteByte('\n') // 바이트 추가
	sb.WriteString("나이: ")
	sb.WriteRune('2') // rune 추가
	sb.WriteRune('5')
	sb.WriteString("세")

	fmt.Println(sb.String())
	fmt.Printf("길이: %d 바이트\n", sb.Len())

	// Reset으로 초기화
	sb.Reset()
	fmt.Printf("Reset 후 길이: %d\n", sb.Len())

	fmt.Println()

	// === fmt.Sprintf 활용 ===
	fmt.Println("=== fmt.Sprintf 포맷 문자열 ===")
	name := "홍길동"
	age := 25
	score := 95.5
	formatted := fmt.Sprintf("%s님은 %d세이며, 점수는 %.1f점입니다.", name, age, score)
	fmt.Println(formatted)
}
