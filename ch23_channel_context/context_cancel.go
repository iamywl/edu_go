package main

import (
	"context"
	"fmt"
	"time"
)

// ============================================================
// 23.2 context.WithCancel 사용법
// 수동으로 취소 신호를 보낼 수 있는 컨텍스트를 생성한다.
// ============================================================

// doWork - 컨텍스트가 취소될 때까지 작업을 수행한다.
func doWork(ctx context.Context, id int, results chan<- string) {
	for i := 1; ; i++ {
		select {
		case <-ctx.Done():
			// 컨텍스트가 취소되면 정리 후 종료한다.
			fmt.Printf("  워커 %d: 취소됨 (이유: %v)\n", id, ctx.Err())
			return
		default:
			// 작업 수행
			result := fmt.Sprintf("워커 %d - 작업 %d 완료", id, i)
			select {
			case results <- result:
			case <-ctx.Done():
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// search - 데이터를 검색하고, 찾으면 결과를 반환한다.
func search(ctx context.Context, id int, target int, found chan<- int) {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("  검색자 %d: 취소됨 (검색 횟수: %d)\n", id, i)
			return
		default:
			// 무작위 검색 시뮬레이션
			if i == target {
				fmt.Printf("  검색자 %d: 찾았습니다! (값: %d)\n", id, i)
				found <- id
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	fmt.Println("=== context.WithCancel 기본 사용법 ===")

	// 취소 가능한 컨텍스트 생성
	ctx, cancel := context.WithCancel(context.Background())

	results := make(chan string, 10)

	// 3개의 워커 시작
	for i := 1; i <= 3; i++ {
		go doWork(ctx, i, results)
	}

	// 몇 개의 결과를 받은 후 취소
	for i := 0; i < 5; i++ {
		fmt.Println(<-results)
	}

	fmt.Println("\n모든 워커를 취소합니다...")
	cancel()                           // 모든 워커에게 취소 신호를 보냄
	time.Sleep(200 * time.Millisecond) // 워커들이 종료될 시간

	fmt.Println("\n=== 검색 후 나머지 취소 패턴 ===")
	fmt.Println("여러 고루틴 중 하나가 결과를 찾으면 나머지를 취소합니다.\n")

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	found := make(chan int, 1)

	// 여러 검색자를 실행한다. 각각 다른 target을 가집니다.
	go search(ctx2, 1, 50, found) // 50번째에 발견
	go search(ctx2, 2, 30, found) // 30번째에 발견
	go search(ctx2, 3, 80, found) // 80번째에 발견

	// 가장 먼저 찾은 검색자의 결과를 받습니다.
	winnerID := <-found
	fmt.Printf("\n검색자 %d이(가) 먼저 찾았습니다!\n", winnerID)

	// 나머지 검색자들을 취소한다.
	cancel2()
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=== 부모-자식 컨텍스트 계층 ===")
	fmt.Println("부모 컨텍스트가 취소되면 자식 컨텍스트도 함께 취소됩니다.\n")

	parentCtx, parentCancel := context.WithCancel(context.Background())

	// 자식 컨텍스트 생성
	childCtx1, childCancel1 := context.WithCancel(parentCtx)
	childCtx2, childCancel2 := context.WithCancel(parentCtx)
	defer childCancel1() // 리소스 누수 방지
	defer childCancel2()

	// 각 컨텍스트에서 고루틴 실행
	go func() {
		<-childCtx1.Done()
		fmt.Println("  자식 컨텍스트 1: 취소됨")
	}()
	go func() {
		<-childCtx2.Done()
		fmt.Println("  자식 컨텍스트 2: 취소됨")
	}()

	time.Sleep(100 * time.Millisecond)

	// 부모를 취소하면 모든 자식도 취소됩니다.
	fmt.Println("부모 컨텍스트를 취소합니다...")
	parentCancel()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
