package main

import (
	"fmt"
	"sync"
)

// ============================================================
// 22.5 뮤텍스를 이용한 동시성 문제 해결
// sync.Mutex와 sync.RWMutex를 사용하여 공유 데이터를 안전하게 보호한다.
// ============================================================

// SafeCounter - 뮤텍스로 보호되는 카운터
type SafeCounter struct {
	mu    sync.Mutex // 뮤텍스
	value int        // 보호할 데이터
}

// Increment - 카운터를 안전하게 증가시킵니다.
func (c *SafeCounter) Increment() {
	c.mu.Lock()         // 잠금 획득 (다른 고루틴은 대기)
	defer c.mu.Unlock() // 함수 종료 시 잠금 해제
	c.value++
}

// Value - 카운터 값을 안전하게 읽습니다.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// SafeMap - RWMutex를 사용하는 안전한 맵
type SafeMap struct {
	mu   sync.RWMutex // 읽기/쓰기 구분 뮤텍스
	data map[string]string
}

// NewSafeMap - SafeMap 생성자
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]string),
	}
}

// Set - 값을 안전하게 저장한다. (쓰기 잠금)
func (m *SafeMap) Set(key, value string) {
	m.mu.Lock() // 쓰기 잠금 (독점 접근)
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get - 값을 안전하게 읽습니다. (읽기 잠금)
func (m *SafeMap) Get(key string) (string, bool) {
	m.mu.RLock() // 읽기 잠금 (여러 고루틴 동시 읽기 가능)
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

// Len - 맵의 크기를 안전하게 반환한다.
func (m *SafeMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

func main() {
	fmt.Println("=== sync.Mutex로 카운터 보호 ===")

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 1000개의 고루틴이 동시에 카운터를 증가
	numGoroutines := 1000
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Increment() // 뮤텍스로 보호되므로 안전!
		}()
	}

	wg.Wait()

	fmt.Printf("기대값: %d\n", numGoroutines)
	fmt.Printf("실제값: %d\n", counter.Value())
	fmt.Println("-> 뮤텍스 덕분에 항상 정확한 값입니다!")

	fmt.Println("\n=== sync.RWMutex로 맵 보호 ===")

	safeMap := NewSafeMap()
	var wg2 sync.WaitGroup

	// 여러 고루틴에서 동시에 쓰기
	writers := 10
	wg2.Add(writers)
	for i := 0; i < writers; i++ {
		go func(id int) {
			defer wg2.Done()
			key := fmt.Sprintf("user_%d", id)
			value := fmt.Sprintf("이름_%d", id)
			safeMap.Set(key, value)
		}(i)
	}
	wg2.Wait()

	// 여러 고루틴에서 동시에 읽기
	readers := 20
	wg2.Add(readers)
	for i := 0; i < readers; i++ {
		go func(id int) {
			defer wg2.Done()
			key := fmt.Sprintf("user_%d", id%writers)
			if val, ok := safeMap.Get(key); ok {
				_ = val // 값 사용 (출력하면 너무 많아지므로 생략)
			}
		}(i)
	}
	wg2.Wait()

	fmt.Printf("맵 크기: %d (기대값: %d)\n", safeMap.Len(), writers)

	// 몇 가지 값 확인
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("user_%d", i)
		if val, ok := safeMap.Get(key); ok {
			fmt.Printf("  %s = %s\n", key, val)
		}
	}

	fmt.Println("\n=== defer를 사용한 안전한 잠금 해제 ===")

	// defer를 사용하면 패닉이 발생해도 잠금이 해제됩니다.
	var mu sync.Mutex
	safeOperation := func() (result string, err error) {
		mu.Lock()
		defer mu.Unlock() // 패닉이 발생해도 반드시 해제됨

		// 여기서 작업 수행...
		return "성공", nil
	}

	result, _ := safeOperation()
	fmt.Println("작업 결과:", result)

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
