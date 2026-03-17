// basic.go
// 함수의 기본 정의와 호출 예제
// 실행 방법: go run basic.go

package main

import "fmt"

// === 매개변수와 반환값이 없는 함수 ===
func greet() {
	fmt.Println("안녕하세요!")
}

// === 매개변수가 있는 함수 ===
func greetUser(name string) {
	fmt.Printf("안녕하세요, %s님!\n", name)
}

// === 매개변수와 반환값이 있는 함수 ===
func add(a int, b int) int {
	return a + b
}

// 같은 타입의 매개변수는 타입을 한 번만 쓸 수 있음
func multiply(a, b int) int {
	return a * b
}

// === 이름이 붙은 반환값 (Named Return Values) ===
func rectangleInfo(width, height float64) (area, perimeter float64) {
	area = width * height
	perimeter = 2 * (width + height)
	return // 이름이 있으므로 반환값 생략 가능 (naked return)
}

// === 가변 인자 함수 (Variadic Function) ===
func sum(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// === Call by Value 확인 ===
func tryChange(x int) {
	fmt.Printf("  함수 안 (변경 전): x = %d\n", x)
	x = 999
	fmt.Printf("  함수 안 (변경 후): x = %d\n", x)
}

// === 함수를 매개변수로 받는 고차 함수 ===
func apply(a, b int, op func(int, int) int) int {
	return op(a, b)
}

// === 클로저 (Closure) ===
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	fmt.Println("===== 매개변수/반환값 없는 함수 =====")
	greet()

	fmt.Println()
	fmt.Println("===== 매개변수가 있는 함수 =====")
	greetUser("홍길동")
	greetUser("Go")

	fmt.Println()
	fmt.Println("===== 매개변수와 반환값이 있는 함수 =====")
	result := add(10, 20)
	fmt.Printf("add(10, 20) = %d\n", result)
	fmt.Printf("multiply(5, 6) = %d\n", multiply(5, 6))

	// 반환값을 바로 사용
	fmt.Printf("add(3, 4) + multiply(2, 5) = %d\n", add(3, 4)+multiply(2, 5))

	fmt.Println()
	fmt.Println("===== 이름이 붙은 반환값 =====")
	area, perimeter := rectangleInfo(5.0, 3.0)
	fmt.Printf("가로 5, 세로 3 → 넓이: %.1f, 둘레: %.1f\n", area, perimeter)

	fmt.Println()
	fmt.Println("===== 가변 인자 함수 =====")
	fmt.Printf("sum(1, 2, 3) = %d\n", sum(1, 2, 3))
	fmt.Printf("sum(1, 2, 3, 4, 5) = %d\n", sum(1, 2, 3, 4, 5))
	fmt.Printf("sum() = %d\n", sum()) // 인자 없이도 호출 가능

	// 슬라이스를 가변 인자로 전달 (... 사용)
	nums := []int{10, 20, 30, 40}
	fmt.Printf("sum(10, 20, 30, 40) = %d\n", sum(nums...))

	fmt.Println()
	fmt.Println("===== Call by Value (값에 의한 호출) =====")
	myValue := 42
	fmt.Printf("함수 호출 전: myValue = %d\n", myValue)
	tryChange(myValue)
	fmt.Printf("함수 호출 후: myValue = %d (변하지 않음!)\n", myValue)

	fmt.Println()
	fmt.Println("===== 함수를 값으로 사용 =====")

	// 함수를 변수에 저장
	myFunc := add
	fmt.Printf("myFunc(100, 200) = %d\n", myFunc(100, 200))

	// 익명 함수 (anonymous function)
	double := func(x int) int {
		return x * 2
	}
	fmt.Printf("double(15) = %d\n", double(15))

	// 고차 함수 사용
	fmt.Printf("apply(10, 3, add) = %d\n", apply(10, 3, add))
	fmt.Printf("apply(10, 3, multiply) = %d\n", apply(10, 3, multiply))

	// 익명 함수를 직접 전달
	result2 := apply(10, 3, func(a, b int) int {
		return a - b
	})
	fmt.Printf("apply(10, 3, subtract) = %d\n", result2)

	fmt.Println()
	fmt.Println("===== 클로저 (Closure) =====")
	// 클로저: 함수가 자신이 생성된 환경의 변수를 기억
	counter := makeCounter()
	fmt.Printf("counter() = %d\n", counter()) // 1
	fmt.Printf("counter() = %d\n", counter()) // 2
	fmt.Printf("counter() = %d\n", counter()) // 3

	// 새로운 카운터는 독립적
	counter2 := makeCounter()
	fmt.Printf("counter2() = %d\n", counter2()) // 1 (독립적인 count)
}
