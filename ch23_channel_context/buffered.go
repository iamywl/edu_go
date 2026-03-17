package main

import (
	"fmt"
	"time"
)

// ============================================================
// 23.1 버퍼드 채널 (Buffered Channel)
// 버퍼가 있는 채널은 버퍼가 가득 찰 때까지 블로킹되지 않습니다.
// ============================================================

// worker - 작업을 처리하는 워커
func worker(id int, jobs <-chan int, results chan<- string) {
	for job := range jobs {
		fmt.Printf("  워커 %d: 작업 %d 처리 중...\n", id, job)
		time.Sleep(100 * time.Millisecond) // 작업 시뮬레이션
		results <- fmt.Sprintf("워커 %d가 작업 %d 완료", id, job)
	}
}

func main() {
	fmt.Println("=== 버퍼드 채널 기본 ===")

	// 버퍼 크기 3인 채널 생성
	ch := make(chan int, 3)

	// 버퍼가 있으므로 받는 쪽 없이도 보낼 수 있습니다.
	ch <- 10
	ch <- 20
	ch <- 30
	fmt.Printf("채널 길이: %d, 용량: %d\n", len(ch), cap(ch))

	// 값 받기
	fmt.Println(<-ch) // 10
	fmt.Println(<-ch) // 20
	fmt.Println(<-ch) // 30
	fmt.Printf("채널 길이: %d, 용량: %d\n", len(ch), cap(ch))

	fmt.Println("\n=== 버퍼드 vs 언버퍼드 비교 ===")

	// 언버퍼드: 동기적 (핑퐁처럼 주고받기)
	fmt.Println("언버퍼드 채널: 보내기/받기가 동기화됨")
	unbuffered := make(chan string)
	go func() {
		unbuffered <- "메시지" // 받는 쪽이 준비될 때까지 대기
	}()
	fmt.Println("  받음:", <-unbuffered)

	// 버퍼드: 비동기적 (우편함처럼 넣어두기)
	fmt.Println("버퍼드 채널: 버퍼가 찰 때까지 비블로킹")
	buffered := make(chan string, 2)
	buffered <- "메시지1" // 블로킹 없음
	buffered <- "메시지2" // 블로킹 없음
	fmt.Println("  받음:", <-buffered)
	fmt.Println("  받음:", <-buffered)

	fmt.Println("\n=== 워커 풀 패턴 (Worker Pool) ===")
	fmt.Println("3개의 워커가 9개의 작업을 분배 처리합니다.\n")

	numJobs := 9
	numWorkers := 3

	// 작업 채널과 결과 채널 (버퍼드)
	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	// 워커 시작
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// 작업 전송
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // 모든 작업을 보낸 후 채널 닫기

	// 결과 수집
	fmt.Println("\n결과:")
	for i := 0; i < numJobs; i++ {
		result := <-results
		fmt.Println(" ", result)
	}

	fmt.Println("\n=== 버퍼드 채널을 세마포어로 사용 ===")

	// 동시에 실행할 수 있는 고루틴 수를 제한한다.
	maxConcurrent := 3
	semaphore := make(chan struct{}, maxConcurrent) // 빈 구조체로 메모리 절약

	for i := 1; i <= 7; i++ {
		semaphore <- struct{}{} // 슬롯 획득 (버퍼가 가득 차면 대기)
		go func(id int) {
			defer func() { <-semaphore }() // 슬롯 반환
			fmt.Printf("  작업 %d 실행 (최대 %d개 동시 실행)\n", id, maxConcurrent)
			time.Sleep(100 * time.Millisecond)
		}(i)
	}

	// 모든 고루틴이 완료될 때까지 대기
	for i := 0; i < maxConcurrent; i++ {
		semaphore <- struct{}{} // 모든 슬롯이 반환될 때까지 대기
	}

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
