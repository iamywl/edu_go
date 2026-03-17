package main

import (
	"fmt"
	"sync"
)

// ============================================================
// 22.4 경쟁 상태 (Race Condition) 데모
// 여러 고루틴이 동시에 공유 데이터에 접근하면 예상치 못한 결과가 발생한다.
// 이 프로그램을 `go run -race race_condition.go`로 실행해 보세요.
// ============================================================

func main() {
	fmt.Println("=== 경쟁 상태 데모 ===")
	fmt.Println("이 프로그램을 'go run -race race_condition.go'로 실행하면")
	fmt.Println("경쟁 상태를 탐지할 수 있습니다.\n")

	// 공유 변수
	counter := 0
	var wg sync.WaitGroup

	// 1000개의 고루틴이 동시에 counter를 증가시킵니다.
	numGoroutines := 1000
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			// 경쟁 상태 발생!
			// counter++는 실제로 3단계 연산이다:
			// 1. 메모리에서 counter 값을 읽음 (예: 42)
			// 2. 값을 1 증가시킴 (43)
			// 3. 메모리에 다시 씀 (43)
			//
			// 두 고루틴이 동시에 실행되면:
			// 고루틴 A: 읽기(42) → 증가(43) → 쓰기(43)
			// 고루틴 B: 읽기(42) → 증가(43) → 쓰기(43)
			// 결과: 44가 되어야 하지만 43이 됨!
			counter++
		}()
	}

	wg.Wait()

	fmt.Printf("기대값: %d\n", numGoroutines)
	fmt.Printf("실제값: %d\n", counter)

	if counter != numGoroutines {
		fmt.Println("-> 경쟁 상태로 인해 값이 손실되었습니다!")
	} else {
		fmt.Println("-> 이번에는 운좋게 정확하지만, 보장되지 않습니다!")
	}

	fmt.Println("\n=== 맵 경쟁 상태 데모 ===")

	// 맵도 동시 접근 시 경쟁 상태가 발생한다.
	// (심한 경우 런타임 패닉이 발생한다)
	data := make(map[string]int)
	var wg2 sync.WaitGroup

	wg2.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg2.Done()
			key := fmt.Sprintf("key_%d", id%3) // 의도적으로 같은 키에 접근
			data[key] = id                     // 경쟁 상태!
		}(i)
	}
	wg2.Wait()

	fmt.Println("맵 내용:", data)
	fmt.Println("-> 맵의 동시 쓰기는 매우 위험합니다!")
	fmt.Println("   'go run -race'로 실행하면 경고가 표시됩니다.")
}
