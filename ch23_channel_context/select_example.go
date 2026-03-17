package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ============================================================
// 23.1 select 문
// 여러 채널 연산을 동시에 대기하고, 준비된 것을 실행한다.
// ============================================================

// slowSearch - 느린 검색 엔진을 시뮬레이션한다.
func slowSearch(query string, delay time.Duration, result chan<- string) {
	time.Sleep(delay)
	result <- fmt.Sprintf("'%s' 검색 결과 (지연: %v)", query, delay)
}

func main() {
	fmt.Println("=== select 기본 사용법 ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// 두 고루틴이 서로 다른 시간에 결과를 보냄
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "채널 1의 메시지"
	}()
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "채널 2의 메시지"
	}()

	// 먼저 준비된 채널에서 받습니다.
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("ch1에서 받음:", msg)
		case msg := <-ch2:
			fmt.Println("ch2에서 받음:", msg)
		}
	}

	fmt.Println("\n=== select로 타임아웃 구현 ===")

	resultCh := make(chan string)

	go slowSearch("Go 언어", 500*time.Millisecond, resultCh)

	select {
	case result := <-resultCh:
		fmt.Println("검색 결과:", result)
	case <-time.After(300 * time.Millisecond):
		fmt.Println("타임아웃! 검색이 너무 오래 걸립니다.")
	}

	fmt.Println("\n=== select의 default 케이스 ===")

	ch := make(chan int, 1)

	// default가 있으면 채널이 준비되지 않았을 때 블로킹되지 않습니다.
	select {
	case v := <-ch:
		fmt.Println("받음:", v)
	default:
		fmt.Println("채널에 데이터가 없습니다 (비블로킹)")
	}

	// 채널에 데이터를 넣은 후 다시 시도
	ch <- 42
	select {
	case v := <-ch:
		fmt.Println("받음:", v)
	default:
		fmt.Println("채널에 데이터가 없습니다")
	}

	fmt.Println("\n=== select로 다중 검색 (먼저 응답한 결과 사용) ===")

	// 여러 검색 엔진에 동시에 요청하고 가장 빠른 결과를 사용한다.
	google := make(chan string, 1)
	bing := make(chan string, 1)
	naver := make(chan string, 1)

	// 각 검색 엔진의 응답 시간이 랜덤
	go slowSearch("Go", time.Duration(rand.Intn(300))*time.Millisecond, google)
	go slowSearch("Go", time.Duration(rand.Intn(300))*time.Millisecond, bing)
	go slowSearch("Go", time.Duration(rand.Intn(300))*time.Millisecond, naver)

	// 가장 먼저 응답한 결과만 사용
	select {
	case result := <-google:
		fmt.Println("Google:", result)
	case result := <-bing:
		fmt.Println("Bing:", result)
	case result := <-naver:
		fmt.Println("Naver:", result)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("모든 검색 엔진이 타임아웃!")
	}

	fmt.Println("\n=== select로 done 채널 패턴 ===")

	done := make(chan struct{}) // 종료 신호용 채널
	dataCh := make(chan int)

	// 데이터 생산자
	go func() {
		defer close(dataCh)
		for i := 1; ; i++ {
			select {
			case dataCh <- i:
				time.Sleep(50 * time.Millisecond)
			case <-done:
				fmt.Println("  생산자: 종료 신호를 받았습니다.")
				return
			}
		}
	}()

	// 5개만 받고 종료
	for i := 0; i < 5; i++ {
		fmt.Printf("  받은 데이터: %d\n", <-dataCh)
	}
	close(done)                        // 종료 신호 전송
	time.Sleep(100 * time.Millisecond) // 생산자가 종료될 시간

	fmt.Println("\n=== for-select 루프 패턴 ===")

	tick := time.Tick(100 * time.Millisecond)     // 주기적 이벤트
	timeout := time.After(550 * time.Millisecond) // 전체 타임아웃

	count := 0
loop:
	for {
		select {
		case <-tick:
			count++
			fmt.Printf("  틱 #%d\n", count)
		case <-timeout:
			fmt.Println("  타임아웃! 루프를 종료합니다.")
			break loop // select가 아닌 for 루프를 빠져나갑니다.
		}
	}

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
