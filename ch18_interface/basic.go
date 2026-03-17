package main

import (
	"fmt"
	"math"
)

// ============================================
// 인터페이스 정의와 구현
// ============================================

// Shape 인터페이스 정의 - 도형이 가져야 할 메서드 목록
type Shape interface {
	Area() float64      // 넓이
	Perimeter() float64 // 둘레
	Name() string       // 도형 이름
}

// --- Rectangle 구조체 ---

type Rectangle struct {
	Width  float64
	Height float64
}

// Rectangle이 Shape 인터페이스를 구현 (암묵적)
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) Name() string {
	return "사각형"
}

// --- Circle 구조체 ---

type Circle struct {
	Radius float64
}

// Circle이 Shape 인터페이스를 구현
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) Name() string {
	return "원"
}

// --- Triangle 구조체 ---

type Triangle struct {
	Base   float64
	Height float64
	SideA  float64
	SideB  float64
	SideC  float64
}

// Triangle이 Shape 인터페이스를 구현
func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

func (t Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func (t Triangle) Name() string {
	return "삼각형"
}

// ============================================
// 인터페이스를 활용하는 함수들
// ============================================

// PrintShapeInfo는 Shape 인터페이스를 받아 도형 정보를 출력한다
// 어떤 도형이든 Shape를 구현하면 이 함수를 사용할 수 있다
func PrintShapeInfo(s Shape) {
	fmt.Printf("  %s - 넓이: %.2f, 둘레: %.2f\n",
		s.Name(), s.Area(), s.Perimeter())
}

// TotalArea는 여러 도형의 넓이 합계를 구한다
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// LargestShape는 넓이가 가장 큰 도형을 반환한다
func LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}
	largest := shapes[0]
	for _, s := range shapes[1:] {
		if s.Area() > largest.Area() {
			largest = s
		}
	}
	return largest
}

func main() {
	// 1. 인터페이스를 통한 다형성
	fmt.Println("=== 인터페이스를 통한 다형성 ===")

	// 다양한 도형 생성
	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	tri := Triangle{Base: 6, Height: 4, SideA: 5, SideB: 5, SideC: 6}

	// Shape 인터페이스 타입의 슬라이스에 모든 도형을 담을 수 있다
	shapes := []Shape{rect, circle, tri}

	// 같은 함수로 모든 도형을 처리
	for _, s := range shapes {
		PrintShapeInfo(s)
	}

	// 2. 인터페이스를 활용한 함수
	fmt.Println("\n=== 전체 넓이 합계 ===")
	total := TotalArea(shapes)
	fmt.Printf("모든 도형의 넓이 합계: %.2f\n", total)

	fmt.Println("\n=== 가장 큰 도형 ===")
	largest := LargestShape(shapes)
	fmt.Printf("가장 큰 도형: %s (넓이: %.2f)\n", largest.Name(), largest.Area())

	// 3. 인터페이스 변수 사용
	fmt.Println("\n=== 인터페이스 변수 ===")
	var s Shape // 인터페이스 변수 (nil)

	s = Rectangle{Width: 3, Height: 4}
	fmt.Printf("Shape = %s, 넓이 = %.2f\n", s.Name(), s.Area())

	s = Circle{Radius: 5} // 같은 변수에 다른 타입 할당
	fmt.Printf("Shape = %s, 넓이 = %.2f\n", s.Name(), s.Area())

	// 4. nil 인터페이스 체크
	fmt.Println("\n=== nil 인터페이스 ===")
	var nilShape Shape
	if nilShape == nil {
		fmt.Println("nilShape은 nil입니다")
	}
}
