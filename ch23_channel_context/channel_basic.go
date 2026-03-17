package main

import "fmt"

// ============================================================
// 23.1 채널 기본 사용법
// 채널은 고루틴 간에 데이터를 주고받는 통신 파이프이다.
// ============================================================

// producer - 데이터를 생성하여 채널로 보냅니다. (보내기 전용 채널)
func producer(out chan<- int, count int) {
	for i := 1; i <= count; i++ {
		fmt.Printf("  보내기: %d\n", i)
		out <- i // 채널에 값 보내기
	}
	close(out) // 모든 데이터를 보낸 후 채널 닫기
}

// consumer - 채널에서 데이터를 받아 처리한다. (받기 전용 채널)
func consumer(in <-chan int) {
	// range는 채널이 닫힐 때까지 반복한다.
	for value := range in {
		fmt.Printf("  받기: %d (제곱: %d)\n", value, value*value)
	}
}

// sum - 채널에서 받은 값들의 합계를 결과 채널로 보냅니다.
func sum(numbers []int, result chan<- int) {
	total := 0
	for _, n := range numbers {
		total += n
	}
	result <- total // 결과를 채널로 전송
}

func main() {
	fmt.Println("=== 언버퍼드 채널 기본 ===")

	// 언버퍼드 채널 생성
	ch := make(chan string)

	// 고루틴에서 채널에 값을 보냄
	go func() {
		ch <- "안녕하세요!" // 받는 쪽이 준비될 때까지 대기
	}()

	// 채널에서 값을 받음
	msg := <-ch // 보내는 쪽이 준비될 때까지 대기
	fmt.Println("받은 메시지:", msg)

	fmt.Println("\n=== 채널로 고루틴 간 통신 ===")

	dataCh := make(chan int)

	// 생산자와 소비자 패턴
	go producer(dataCh, 5)
	consumer(dataCh) // 메인 고루틴에서 소비

	fmt.Println("\n=== 채널로 결과 수집 ===")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	mid := len(numbers) / 2

	resultCh := make(chan int)

	// 두 고루틴이 각각 절반씩 합계를 계산
	go sum(numbers[:mid], resultCh) // 1~5
	go sum(numbers[mid:], resultCh) // 6~10

	// 두 결과를 받아서 합산
	sum1 := <-resultCh
	sum2 := <-resultCh
	fmt.Printf("부분합: %d + %d = %d\n", sum1, sum2, sum1+sum2)

	fmt.Println("\n=== 채널 닫기와 range ===")

	numCh := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			numCh <- i * 10
		}
		close(numCh) // 채널을 닫아야 range가 종료됩니다.
	}()

	// range로 채널의 모든 값을 받습니다.
	fmt.Print("받은 값: ")
	for v := range numCh {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	fmt.Println("\n=== 채널 닫힘 확인 ===")

	ch2 := make(chan int, 3)
	ch2 <- 100
	ch2 <- 200
	close(ch2)

	// 두 번째 반환값으로 채널이 닫혔는지 확인할 수 있습니다.
	v, ok := <-ch2
	fmt.Printf("값: %d, 열림: %v\n", v, ok) // 100, true

	v, ok = <-ch2
	fmt.Printf("값: %d, 열림: %v\n", v, ok) // 200, true

	v, ok = <-ch2
	fmt.Printf("값: %d, 열림: %v\n", v, ok) // 0, false (닫힌 채널)

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
