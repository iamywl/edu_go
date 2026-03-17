package main

import (
	"fmt"
	"runtime"
	"time"
)

// =============================================
// Go 가비지 컬렉터(GC) 동작 데모
// 메모리 할당, GC 트리거, 통계 확인
// =============================================

// printMemStats는 현재 메모리 통계를 출력한다
func printMemStats(label string) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fmt.Printf("  [%s]\n", label)
	fmt.Printf("    힙 할당 중 (HeapAlloc):  %6d KB\n", stats.HeapAlloc/1024)
	fmt.Printf("    힙 시스템 (HeapSys):     %6d KB\n", stats.HeapSys/1024)
	fmt.Printf("    힙 객체 수 (HeapObjects): %6d 개\n", stats.HeapObjects)
	fmt.Printf("    GC 실행 횟수 (NumGC):    %6d 회\n", stats.NumGC)
	fmt.Printf("    총 할당량 (TotalAlloc):  %6d KB\n", stats.TotalAlloc/1024)
	fmt.Println()
}

// allocateMemory는 지정된 크기의 메모리를 할당한다
func allocateMemory(sizeMB int) []byte {
	data := make([]byte, sizeMB*1024*1024)
	// 실제로 메모리를 사용해야 OS가 할당한다
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}

// createGarbage는 금방 버려질 임시 객체들을 대량 생성한다
func createGarbage(count int) {
	for i := 0; i < count; i++ {
		// 각 반복에서 생성된 슬라이스는 다음 반복에서 참조를 잃음 → GC 대상
		data := make([]byte, 1024) // 1KB씩 할당
		_ = data
	}
}

// --- 탈출 분석 예제 ---

// noEscape: 로컬 변수가 함수 밖으로 나가지 않음 → 스택에 할당
//
//go:noinline
func noEscape() int {
	x := 42 // 스택에 할당 (GC 대상이 아님)
	return x
}

// doesEscape: 포인터를 반환하면 힙에 할당됨 → GC 대상
//
//go:noinline
func doesEscape() *int {
	x := 42   // 힙에 할당 (함수 밖에서 사용되므로)
	return &x // 포인터를 반환 → 탈출(escape)
}

func main() {
	fmt.Println("========== Go 가비지 컬렉터 데모 ==========\n")

	// --- 1. 초기 메모리 상태 ---
	fmt.Println("=== 1. 초기 메모리 상태 ===")
	runtime.GC() // 깨끗한 상태에서 시작
	printMemStats("초기 상태")

	// --- 2. 메모리 할당 후 ---
	fmt.Println("=== 2. 메모리 할당 (5MB) ===")
	data := allocateMemory(5) // 5MB 할당
	printMemStats("5MB 할당 후")

	// --- 3. 참조 제거 + GC 실행 ---
	fmt.Println("=== 3. 참조 제거 + GC 실행 ===")
	_ = data     // 컴파일러 경고 방지
	data = nil   // 참조 제거 → GC 회수 가능
	runtime.GC() // 명시적 GC 트리거
	printMemStats("참조 제거 + GC 후")

	// --- 4. 가비지 대량 생성 ---
	fmt.Println("=== 4. 가비지 대량 생성 (10,000개 객체) ===")
	printMemStats("가비지 생성 전")

	createGarbage(10000) // 10,000개의 임시 객체 생성
	printMemStats("가비지 생성 직후")

	runtime.GC()
	printMemStats("GC 실행 후")

	// --- 5. GC 동작 관찰 ---
	fmt.Println("=== 5. GC 자동 트리거 관찰 ===")
	fmt.Println("  GOGC=100(기본값): 힙이 이전 GC 후의 100% 성장하면 GC 실행")
	fmt.Println()

	var gcCount uint32
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	gcCount = stats.NumGC

	// 반복적으로 메모리를 할당하며 GC 자동 트리거 관찰
	for i := 0; i < 20; i++ {
		_ = allocateMemory(1) // 1MB 할당 후 즉시 참조 해제

		runtime.ReadMemStats(&stats)
		if stats.NumGC > gcCount {
			fmt.Printf("  [반복 %2d] GC가 자동 실행됨! (총 %d회)\n", i, stats.NumGC)
			gcCount = stats.NumGC
		}
	}

	// --- 6. GC 일시 정지 시간 측정 ---
	fmt.Println("\n=== 6. GC 일시 정지(STW) 시간 ===")
	runtime.GC()
	runtime.ReadMemStats(&stats)

	// 최근 GC 일시 정지 시간 확인
	// PauseNs는 최근 256번의 GC pause 시간을 나노초로 저장
	recentPauses := stats.PauseNs[:5] // 최근 5번
	fmt.Println("  최근 GC 일시 정지 시간:")
	for i, pause := range recentPauses {
		if pause > 0 {
			fmt.Printf("    GC[%d]: %v\n", i, time.Duration(pause))
		}
	}

	// --- 7. 탈출 분석 설명 ---
	fmt.Println("\n=== 7. 탈출 분석 (Escape Analysis) ===")
	fmt.Println("  Go 컴파일러는 변수가 함수 밖으로 '탈출'하는지 분석합니다.")
	fmt.Println()

	v1 := noEscape()   // 스택 할당 → GC 부담 없음
	v2 := doesEscape() // 힙 할당 → GC 대상

	fmt.Printf("  noEscape() = %d   (스택에 할당, GC 대상 아님)\n", v1)
	fmt.Printf("  doesEscape() = %d (힙에 할당, GC가 관리)\n", *v2)
	fmt.Println()
	fmt.Println("  탈출 분석 확인: go build -gcflags='-m' 명령으로 확인 가능")

	// --- 8. GC 친화적인 코드 팁 ---
	fmt.Println("\n=== 8. GC 친화적인 코드 작성법 ===")
	fmt.Println("  1. 슬라이스를 미리 할당하라: make([]T, 0, expectedSize)")
	fmt.Println("  2. 루프 안에서 불필요한 할당을 피하라")
	fmt.Println("  3. sync.Pool로 자주 사용하는 객체를 재활용하라")
	fmt.Println("  4. 작은 구조체는 값 타입으로 사용하라 (포인터 추적 감소)")
	fmt.Println("  5. 문자열 연결 시 strings.Builder를 사용하라")
	fmt.Println()

	// --- 미리 할당 예제 ---
	fmt.Println("  미리 할당 예제:")
	fmt.Println("    나쁜 예: s := []int{}; for ... { s = append(s, v) }")
	fmt.Println("    좋은 예: s := make([]int, 0, 1000); for ... { s = append(s, v) }")

	// --- GOGC 설정 안내 ---
	fmt.Println("\n=== GOGC 환경변수 ===")
	fmt.Println("  GOGC=100 (기본값): 힙이 100% 성장하면 GC 실행")
	fmt.Println("  GOGC=200: GC 덜 실행 (메모리 더 사용, CPU 절약)")
	fmt.Println("  GOGC=50:  GC 더 자주 (메모리 절약, CPU 더 사용)")
	fmt.Println("  GOGC=off: GC 비활성화 (주의!)")
	fmt.Println()
	fmt.Println("  GOMEMLIMIT: Go 1.19+ 메모리 상한 설정")
	fmt.Println("  예: GOMEMLIMIT=1GiB ./myapp")
}
