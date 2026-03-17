package main

import "fmt"

func main() {
	// === 포인터 기본 개념 ===
	fmt.Println("=== 포인터 기본 ===")

	a := 42
	p := &a // p는 a의 메모리 주소를 저장

	fmt.Printf("a의 값:    %d\n", a)
	fmt.Printf("a의 주소:  %p\n", &a)
	fmt.Printf("p의 값:    %p (a의 주소와 동일)\n", p)
	fmt.Printf("*p의 값:   %d (p가 가리키는 곳의 값)\n", *p)

	fmt.Println()

	// === 포인터로 값 변경 ===
	fmt.Println("=== 포인터로 값 변경 ===")
	fmt.Printf("변경 전 a: %d\n", a)
	*p = 100 // 포인터를 통해 a의 값을 변경
	fmt.Printf("변경 후 a: %d\n", a)

	fmt.Println()

	// === 포인터 타입 ===
	fmt.Println("=== 포인터 타입 ===")
	var intPtr *int    // int 포인터, 제로값은 nil
	var strPtr *string // string 포인터

	fmt.Printf("intPtr: %v (nil 여부: %v)\n", intPtr, intPtr == nil)
	fmt.Printf("strPtr: %v (nil 여부: %v)\n", strPtr, strPtr == nil)

	// nil 포인터에 값을 할당하려면 먼저 주소를 지정해야 함
	num := 77
	intPtr = &num
	fmt.Printf("intPtr가 가리키는 값: %d\n", *intPtr)

	fmt.Println()

	// === 포인터와 함수 ===
	fmt.Println("=== 포인터와 함수 ===")
	x, y := 10, 20
	fmt.Printf("swap 전: x=%d, y=%d\n", x, y)
	swap(&x, &y)
	fmt.Printf("swap 후: x=%d, y=%d\n", x, y)

	fmt.Println()

	// === nil 포인터 안전하게 다루기 ===
	fmt.Println("=== nil 체크 ===")
	var safePtr *int
	if safePtr != nil {
		fmt.Println(*safePtr)
	} else {
		fmt.Println("포인터가 nil입니다. 역참조하면 패닉이 발생합니다!")
	}
}

// swap 은 두 정수의 값을 교환하는 함수 (포인터 사용)
func swap(a, b *int) {
	*a, *b = *b, *a
}
