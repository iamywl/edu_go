package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// 22.7 sync.WaitGroup과 sync.Once 사용법
// WaitGroup: 여러 고루틴의 완료를 대기한다.
// Once: 함수를 딱 한 번만 실행한다.
// ============================================================

// simulateTask - 작업을 시뮬레이션하는 함수
func simulateTask(id int, duration time.Duration) {
	fmt.Printf("  작업 %d 시작 (소요: %v)\n", id, duration)
	time.Sleep(duration)
	fmt.Printf("  작업 %d 완료!\n", id)
}

// Config - 설정 구조체 (싱글톤 패턴 예제)
type Config struct {
	DBHost string
	DBPort int
}

// 싱글톤 패턴을 위한 전역 변수
var (
	configInstance *Config
	configOnce     sync.Once
)

// GetConfig - sync.Once를 사용한 싱글톤 패턴
// 여러 고루틴에서 동시에 호출해도 초기화는 한 번만 수행됩니다.
func GetConfig() *Config {
	configOnce.Do(func() {
		fmt.Println("  [초기화] 설정을 로드합니다... (이 메시지는 한 번만 출력됩니다)")
		time.Sleep(100 * time.Millisecond) // 파일 읽기를 시뮬레이션
		configInstance = &Config{
			DBHost: "localhost",
			DBPort: 5432,
		}
	})
	return configInstance
}

func main() {
	fmt.Println("=== sync.WaitGroup 기본 사용법 ===")

	var wg sync.WaitGroup

	// 3개의 작업을 고루틴으로 실행
	tasks := []struct {
		id       int
		duration time.Duration
	}{
		{1, 300 * time.Millisecond},
		{2, 100 * time.Millisecond},
		{3, 200 * time.Millisecond},
	}

	for _, task := range tasks {
		wg.Add(1) // 카운터를 1 증가
		go func(id int, dur time.Duration) {
			defer wg.Done() // 작업 완료 시 카운터 1 감소
			simulateTask(id, dur)
		}(task.id, task.duration)
	}

	fmt.Println("모든 작업이 시작되었습니다. 완료를 기다립니다...")
	wg.Wait() // 카운터가 0이 될 때까지 대기
	fmt.Println("모든 작업이 완료되었습니다!")

	fmt.Println("\n=== WaitGroup으로 병렬 합계 계산 ===")

	var (
		wg2   sync.WaitGroup
		mu    sync.Mutex
		total int
	)

	// 1부터 100까지의 합을 10개의 고루틴으로 나누어 계산
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func(start int) {
			defer wg2.Done()

			// 각 고루틴은 10개의 숫자를 합산
			partialSum := 0
			for j := start; j < start+10; j++ {
				partialSum += j
			}

			// 뮤텍스로 보호하며 전체 합에 추가
			mu.Lock()
			total += partialSum
			mu.Unlock()

			fmt.Printf("  고루틴 %d: %d ~ %d 합계 = %d\n", start/10, start, start+9, partialSum)
		}(i*10 + 1)
	}

	wg2.Wait()
	fmt.Printf("1부터 100까지의 합: %d (검증: %d)\n", total, 100*101/2)

	fmt.Println("\n=== sync.Once 사용법 ===")

	var wg3 sync.WaitGroup

	// 10개의 고루틴이 동시에 GetConfig를 호출
	wg3.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg3.Done()
			config := GetConfig() // 내부 초기화는 한 번만 실행됨
			fmt.Printf("  고루틴 %d: DB=%s:%d\n", id, config.DBHost, config.DBPort)
		}(i)
	}

	wg3.Wait()
	fmt.Println("모든 고루틴이 같은 설정 인스턴스를 사용했습니다.")

	fmt.Println("\n=== sync.Once의 중요한 특성 ===")

	// Once.Do에 전달된 함수가 패닉을 발생시켜도
	// 이후 호출에서 다시 실행되지 않습니다.
	var once sync.Once
	for i := 0; i < 3; i++ {
		once.Do(func() {
			fmt.Println("  이 메시지는 딱 한 번만 출력됩니다!")
		})
	}

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
