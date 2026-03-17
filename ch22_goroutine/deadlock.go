package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// 22.6 데드락 (Deadlock) 예제
// 두 개 이상의 고루틴이 서로가 가진 잠금을 기다리며
// 영원히 대기하는 상태를 데드락이라 한다.
// ============================================================

// ---- 데드락 발생 시나리오 (시뮬레이션) ----

// 두 개의 뮤텍스를 교차로 잠그는 데드락 시나리오
// 실제로 실행하면 프로그램이 멈추므로 타임아웃으로 감지한다.
func demonstrateDeadlockScenario() {
	var mu1, mu2 sync.Mutex
	done := make(chan bool, 2)

	// 고루틴 A: mu1 -> mu2 순서로 잠금
	go func() {
		mu1.Lock()
		fmt.Println("  고루틴 A: mu1 잠금 획득")
		time.Sleep(100 * time.Millisecond) // 고루틴 B가 mu2를 잠글 시간을 줌

		fmt.Println("  고루틴 A: mu2 잠금 시도...")
		mu2.Lock() // 고루틴 B가 mu2를 들고 있으므로 대기 (데드락!)
		fmt.Println("  고루틴 A: mu2 잠금 획득")
		mu2.Unlock()
		mu1.Unlock()
		done <- true
	}()

	// 고루틴 B: mu2 -> mu1 순서로 잠금 (역순!)
	go func() {
		mu2.Lock()
		fmt.Println("  고루틴 B: mu2 잠금 획득")
		time.Sleep(100 * time.Millisecond) // 고루틴 A가 mu1을 잠글 시간을 줌

		fmt.Println("  고루틴 B: mu1 잠금 시도...")
		mu1.Lock() // 고루틴 A가 mu1을 들고 있으므로 대기 (데드락!)
		fmt.Println("  고루틴 B: mu1 잠금 획득")
		mu1.Unlock()
		mu2.Unlock()
		done <- true
	}()

	// 타임아웃으로 데드락을 감지한다.
	select {
	case <-done:
		fmt.Println("  작업 완료 (데드락 없음)")
	case <-time.After(1 * time.Second):
		fmt.Println("  [타임아웃] 데드락이 발생했습니다!")
		fmt.Println("  고루틴 A는 mu2를 기다리고, 고루틴 B는 mu1을 기다립니다.")
	}
}

// ---- 데드락 해결: 잠금 순서 일관성 ----

func demonstrateDeadlockFix() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup

	// 해결 방법: 항상 같은 순서로 잠금을 획득한다.
	// 두 고루틴 모두 mu1 -> mu2 순서로 잠급니다.

	wg.Add(2)

	// 고루틴 A: mu1 -> mu2 순서
	go func() {
		defer wg.Done()
		mu1.Lock()
		fmt.Println("  고루틴 A: mu1 잠금 획득")
		time.Sleep(50 * time.Millisecond)

		mu2.Lock()
		fmt.Println("  고루틴 A: mu2 잠금 획득")

		// 작업 수행
		fmt.Println("  고루틴 A: 작업 완료")

		mu2.Unlock()
		mu1.Unlock()
	}()

	// 고루틴 B: mu1 -> mu2 순서 (동일한 순서!)
	go func() {
		defer wg.Done()
		mu1.Lock()
		fmt.Println("  고루틴 B: mu1 잠금 획득")
		time.Sleep(50 * time.Millisecond)

		mu2.Lock()
		fmt.Println("  고루틴 B: mu2 잠금 획득")

		// 작업 수행
		fmt.Println("  고루틴 B: 작업 완료")

		mu2.Unlock()
		mu1.Unlock()
	}()

	wg.Wait()
}

// ---- 잠금 범위 최소화 예제 ----

type BankAccount struct {
	mu      sync.Mutex
	balance int
}

func (a *BankAccount) Deposit(amount int) {
	// 좋은 예: 잠금 범위를 최소화
	a.mu.Lock()
	a.balance += amount
	a.mu.Unlock()
	// 잠금 해제 후 로깅 등 비임계 작업 수행
	fmt.Printf("  %d원 입금 완료 (잔액: %d원)\n", amount, a.balance)
}

func (a *BankAccount) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func main() {
	fmt.Println("=== 데드락 시나리오 (타임아웃으로 감지) ===")
	fmt.Println("두 고루틴이 서로의 잠금을 기다리는 상황:")
	demonstrateDeadlockScenario()

	fmt.Println("\n=== 데드락 해결: 잠금 순서 일관성 ===")
	fmt.Println("두 고루틴이 같은 순서로 잠금을 획득:")
	demonstrateDeadlockFix()
	fmt.Println("-> 데드락 없이 완료!")

	fmt.Println("\n=== 데드락 방지 팁 ===")
	fmt.Println("1. 여러 뮤텍스를 사용할 때 항상 같은 순서로 잠금을 획득하세요.")
	fmt.Println("2. 잠금 범위를 최소화하세요 (필요한 코드만 보호).")
	fmt.Println("3. 가능하면 하나의 뮤텍스만 사용하세요.")
	fmt.Println("4. 채널 기반 설계를 고려하세요.")
	fmt.Println("5. context.WithTimeout으로 타임아웃을 설정하세요.")

	fmt.Println("\n=== 잠금 범위 최소화 예제 ===")
	account := &BankAccount{}
	var wg sync.WaitGroup

	wg.Add(5)
	for i := 1; i <= 5; i++ {
		go func(amount int) {
			defer wg.Done()
			account.Deposit(amount * 1000)
		}(i)
	}
	wg.Wait()

	fmt.Printf("최종 잔액: %d원\n", account.Balance())
	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
