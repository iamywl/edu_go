package main

import "fmt"

// ============================================
// 함수 리터럴 (익명 함수)과 클로저
// ============================================

// Counter는 카운터 클로저를 반환한다
// 호출할 때마다 1씩 증가하는 값을 반환
func Counter() func() int {
	count := 0 // 클로저에 의해 캡처되는 외부 변수
	return func() int {
		count++
		return count
	}
}

// CounterWithStep은 지정된 간격으로 증가하는 카운터를 반환한다
func CounterWithStep(start, step int) func() int {
	current := start - step
	return func() int {
		current += step
		return current
	}
}

// Fibonacci는 피보나치 수열 생성기를 반환한다
func Fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		result := a
		a, b = b, a+b
		return result
	}
}

// Logger는 접두사를 기억하는 로깅 함수를 반환한다
func Logger(prefix string) func(string) {
	count := 0
	return func(msg string) {
		count++
		fmt.Printf("[%s #%d] %s\n", prefix, count, msg)
	}
}

// Accumulator는 누적 합계를 관리하는 클로저를 반환한다
func Accumulator(initial int) (add func(int) int, getTotal func() int) {
	total := initial
	add = func(n int) int {
		total += n
		return total
	}
	getTotal = func() int {
		return total
	}
	return
}

// Cache는 간단한 메모이제이션 캐시를 구현한다
func Cache(fn func(int) int) func(int) int {
	cache := make(map[int]int)
	return func(n int) int {
		if val, ok := cache[n]; ok {
			fmt.Printf("    캐시 히트: f(%d) = %d\n", n, val)
			return val
		}
		result := fn(n)
		cache[n] = result
		fmt.Printf("    캐시 미스: f(%d) = %d (계산됨)\n", n, result)
		return result
	}
}

func main() {
	// 1. 익명 함수 기본
	fmt.Println("=== 익명 함수 기본 ===")

	// 변수에 익명 함수 할당
	greet := func(name string) string {
		return fmt.Sprintf("안녕하세요, %s님!", name)
	}
	fmt.Println(greet("김철수"))
	fmt.Println(greet("이영희"))

	// 즉시 실행 함수 (IIFE)
	result := func(a, b int) int {
		return a + b
	}(3, 4)
	fmt.Println("즉시 실행 결과:", result)

	// 2. 클로저 - Counter
	fmt.Println("\n=== 클로저: Counter ===")
	counter1 := Counter()
	counter2 := Counter() // 독립적인 카운터

	fmt.Println("counter1:", counter1()) // 1
	fmt.Println("counter1:", counter1()) // 2
	fmt.Println("counter1:", counter1()) // 3
	fmt.Println("counter2:", counter2()) // 1 (독립적)
	fmt.Println("counter2:", counter2()) // 2

	// 3. 클로저 - 커스텀 카운터
	fmt.Println("\n=== 커스텀 카운터 ===")
	evens := CounterWithStep(0, 2)  // 짝수 생성기
	odds := CounterWithStep(1, 2)   // 홀수 생성기
	tens := CounterWithStep(10, 10) // 10씩 증가

	fmt.Print("짝수: ")
	for i := 0; i < 5; i++ {
		fmt.Print(evens(), " ")
	}
	fmt.Println()

	fmt.Print("홀수: ")
	for i := 0; i < 5; i++ {
		fmt.Print(odds(), " ")
	}
	fmt.Println()

	fmt.Print("10의 배수: ")
	for i := 0; i < 5; i++ {
		fmt.Print(tens(), " ")
	}
	fmt.Println()

	// 4. 클로저 - 피보나치 수열
	fmt.Println("\n=== 피보나치 수열 생성기 ===")
	fib := Fibonacci()
	fmt.Print("피보나치: ")
	for i := 0; i < 10; i++ {
		fmt.Print(fib(), " ")
	}
	fmt.Println()

	// 5. 클로저 - Logger
	fmt.Println("\n=== Logger 클로저 ===")
	infoLog := Logger("INFO")
	errorLog := Logger("ERROR")

	infoLog("서버가 시작되었습니다")
	infoLog("요청을 처리 중입니다")
	errorLog("데이터베이스 연결 실패")
	infoLog("요청 처리 완료")
	errorLog("타임아웃 발생")

	// 6. 클로저 - 누적기 (다중 반환)
	fmt.Println("\n=== Accumulator 클로저 ===")
	add, getTotal := Accumulator(100) // 초기값 100
	fmt.Println("초기값:", getTotal())
	fmt.Println("+ 50 =", add(50))
	fmt.Println("+ 30 =", add(30))
	fmt.Println("- 20 =", add(-20))
	fmt.Println("현재 합계:", getTotal())

	// 7. 클로저 - 캐시 (메모이제이션)
	fmt.Println("\n=== 캐시 (메모이제이션) ===")
	// 무거운 계산을 시뮬레이션
	expensiveSquare := Cache(func(n int) int {
		return n * n
	})

	fmt.Println("  결과:", expensiveSquare(5)) // 캐시 미스
	fmt.Println("  결과:", expensiveSquare(3)) // 캐시 미스
	fmt.Println("  결과:", expensiveSquare(5)) // 캐시 히트!
	fmt.Println("  결과:", expensiveSquare(3)) // 캐시 히트!
	fmt.Println("  결과:", expensiveSquare(7)) // 캐시 미스

	// 8. 클로저와 반복문 주의사항
	fmt.Println("\n=== 클로저와 반복문 주의 ===")

	// 올바른 방법: 반복 변수를 새 변수로 캡처
	funcs := make([]func(), 5)
	for i := 0; i < 5; i++ {
		i := i // 새 변수로 캡처 (Go에서는 이 패턴이 관용적)
		funcs[i] = func() {
			fmt.Printf("  함수 %d 호출됨\n", i)
		}
	}

	for _, f := range funcs {
		f() // 0, 1, 2, 3, 4 각각 출력
	}
}
