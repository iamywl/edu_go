package main

import (
	"fmt"
	"math"
)

// ============================================
// 메서드 선언과 기본 사용법
// ============================================

// Rectangle 구조체 정의
type Rectangle struct {
	Width  float64
	Height float64
}

// Circle 구조체 정의
type Circle struct {
	Radius float64
}

// --- Rectangle 메서드들 ---

// Area는 사각형의 넓이를 반환한다 (값 타입 리시버)
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter는 사각형의 둘레를 반환한다
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// String은 사각형 정보를 문자열로 반환한다
func (r Rectangle) String() string {
	return fmt.Sprintf("사각형(가로=%.1f, 세로=%.1f)", r.Width, r.Height)
}

// --- Circle 메서드들 ---

// Area는 원의 넓이를 반환한다
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Circumference는 원의 둘레를 반환한다
func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

// String은 원 정보를 문자열로 반환한다
func (c Circle) String() string {
	return fmt.Sprintf("원(반지름=%.1f)", c.Radius)
}

// --- 사용자 정의 타입에 메서드 추가 ---

// MyString은 사용자 정의 문자열 타입
type MyString string

// Reverse는 문자열을 뒤집어 반환한다
func (s MyString) Reverse() MyString {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return MyString(runes)
}

// Len은 문자열의 길이를 반환한다
func (s MyString) Len() int {
	return len([]rune(s))
}

func main() {
	// 1. Rectangle 메서드 사용
	fmt.Println("=== Rectangle 메서드 ===")
	rect := Rectangle{Width: 10, Height: 5}
	fmt.Println(rect)                          // String() 메서드 자동 호출
	fmt.Printf("넓이: %.1f\n", rect.Area())      // 50.0
	fmt.Printf("둘레: %.1f\n", rect.Perimeter()) // 30.0

	// 2. Circle 메서드 사용
	fmt.Println("\n=== Circle 메서드 ===")
	circle := Circle{Radius: 7}
	fmt.Println(circle)                              // String() 자동 호출
	fmt.Printf("넓이: %.2f\n", circle.Area())          // 153.94
	fmt.Printf("둘레: %.2f\n", circle.Circumference()) // 43.98

	// 3. 같은 이름의 메서드 - 타입에 따라 다르게 동작
	fmt.Println("\n=== 같은 이름, 다른 타입 ===")
	fmt.Printf("사각형 넓이: %.1f\n", rect.Area())
	fmt.Printf("원 넓이: %.2f\n", circle.Area())
	// 둘 다 Area()이지만 타입에 맞는 메서드가 호출됨

	// 4. 사용자 정의 타입의 메서드
	fmt.Println("\n=== 사용자 정의 타입 메서드 ===")
	greeting := MyString("안녕하세요")
	fmt.Println("원본:", greeting)
	fmt.Println("뒤집기:", greeting.Reverse())
	fmt.Println("길이:", greeting.Len())

	// 5. 메서드 vs 함수 비교
	fmt.Println("\n=== 메서드 vs 함수 ===")
	r := Rectangle{Width: 8, Height: 3}

	// 메서드 방식 (더 직관적)
	fmt.Printf("메서드: %s의 넓이 = %.1f\n", r, r.Area())

	// 일반 함수 방식 (비교용)
	fmt.Printf("함수:   넓이 = %.1f\n", calculateArea(r))
}

// 일반 함수 (비교용)
func calculateArea(r Rectangle) float64 {
	return r.Width * r.Height
}
