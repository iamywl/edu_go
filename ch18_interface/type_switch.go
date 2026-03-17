package main

import "fmt"

// ============================================
// 타입 단언과 타입 스위치
// ============================================

// Shape 인터페이스
type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

type Triangle struct {
	Base, Height float64
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

func main() {
	// 1. 타입 단언 (Type Assertion) 기본
	fmt.Println("=== 타입 단언 기본 ===")

	var val any = "안녕하세요"

	// 기본 타입 단언 - 성공
	str := val.(string)
	fmt.Println("문자열:", str)

	// 안전한 타입 단언 (ok 패턴) - 성공
	str2, ok := val.(string)
	fmt.Printf("string 타입? ok=%t, 값=%s\n", ok, str2)

	// 안전한 타입 단언 - 실패
	num, ok := val.(int)
	fmt.Printf("int 타입? ok=%t, 값=%d\n", ok, num) // ok=false, 값=0

	// 2. 타입 스위치 (Type Switch)
	fmt.Println("\n=== 타입 스위치 ===")

	values := []any{
		42,
		3.14,
		"Hello",
		true,
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		nil,
	}

	for _, v := range values {
		describeType(v)
	}

	// 3. 인터페이스에서 구체적인 타입 꺼내기
	fmt.Println("\n=== 인터페이스에서 구체적 타입 ===")

	shapes := []Shape{
		Rectangle{10, 5},
		Circle{7},
		Triangle{6, 4},
		Rectangle{3, 3},
	}

	for _, s := range shapes {
		describeShape(s)
	}

	// 4. 타입 단언으로 추가 메서드 접근
	fmt.Println("\n=== 타입 단언으로 추가 메서드 접근 ===")
	var s Shape = Rectangle{Width: 10, Height: 5}

	// Shape 인터페이스에는 Area()만 있지만
	// Rectangle으로 단언하면 Width, Height에 접근 가능
	if rect, ok := s.(Rectangle); ok {
		fmt.Printf("사각형 가로: %.1f, 세로: %.1f, 넓이: %.1f\n",
			rect.Width, rect.Height, rect.Area())
	}

	// 5. 실용적인 예제 - 다양한 타입의 값 합산
	fmt.Println("\n=== 실용 예제: 다양한 타입 합산 ===")
	mixed := []any{10, 20.5, "30", true, 40, 50.5, "무시됨"}
	total := sumValues(mixed)
	fmt.Printf("합계: %.1f\n", total)
}

// describeType은 any 타입의 값을 타입 스위치로 분석한다
func describeType(val any) {
	switch v := val.(type) {
	case int:
		fmt.Printf("  정수: %d\n", v)
	case float64:
		fmt.Printf("  실수: %.2f\n", v)
	case string:
		fmt.Printf("  문자열: %q (길이: %d)\n", v, len(v))
	case bool:
		fmt.Printf("  불리언: %t\n", v)
	case []int:
		fmt.Printf("  정수 슬라이스: %v (길이: %d)\n", v, len(v))
	case nil:
		fmt.Println("  nil 값")
	default:
		fmt.Printf("  기타 타입: %T = %v\n", v, v)
	}
}

// describeShape은 Shape 인터페이스의 구체적 타입을 분석한다
func describeShape(s Shape) {
	switch v := s.(type) {
	case Rectangle:
		fmt.Printf("  사각형 (%.1f x %.1f) -> 넓이: %.2f\n",
			v.Width, v.Height, v.Area())
	case Circle:
		fmt.Printf("  원 (반지름: %.1f) -> 넓이: %.2f\n",
			v.Radius, v.Area())
	case Triangle:
		fmt.Printf("  삼각형 (밑변: %.1f, 높이: %.1f) -> 넓이: %.2f\n",
			v.Base, v.Height, v.Area())
	default:
		fmt.Printf("  알 수 없는 도형: 넓이 %.2f\n", v.Area())
	}
}

// sumValues는 다양한 타입의 값을 숫자로 변환하여 합산한다
func sumValues(values []any) float64 {
	total := 0.0
	for _, v := range values {
		switch val := v.(type) {
		case int:
			total += float64(val)
		case float64:
			total += val
		case string:
			// 문자열은 숫자 변환 시도
			var n float64
			_, err := fmt.Sscanf(val, "%f", &n)
			if err == nil {
				total += n
				fmt.Printf("  문자열 %q -> 숫자 %.1f 변환 성공\n", val, n)
			} else {
				fmt.Printf("  문자열 %q -> 숫자 변환 실패, 건너뜀\n", val)
			}
		case bool:
			if val {
				total += 1
			}
		}
	}
	return total
}
