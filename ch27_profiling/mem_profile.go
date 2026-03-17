// 27장: 프로파일링 - 메모리 프로파일링 예제
// 실행: go run mem_profile.go
// 분석: go tool pprof mem.prof
//
// pprof 명령어:
//
//	top             - 메모리 할당 상위 함수
//	list 함수명      - 특정 함수의 라인별 메모리 사용량
//	go tool pprof -alloc_space mem.prof  - 총 할당량 기준
//	go tool pprof -inuse_space mem.prof  - 현재 사용량 기준
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

// allocateWithAppend - append로 슬라이스를 키우는 방식 (메모리 재할당 발생)
func allocateWithAppend() []int {
	var result []int
	for i := 0; i < 100000; i++ {
		result = append(result, i)
	}
	return result
}

// allocateWithCapacity - 미리 용량을 할당하는 방식 (효율적)
func allocateWithCapacity() []int {
	result := make([]int, 0, 100000)
	for i := 0; i < 100000; i++ {
		result = append(result, i)
	}
	return result
}

// createManyStrings - 많은 문자열을 생성하는 함수
func createManyStrings() []string {
	result := make([]string, 0, 10000)
	for i := 0; i < 10000; i++ {
		// 매번 새로운 문자열을 할당한다
		s := fmt.Sprintf("문자열_%d_데이터_%d", i, i*i)
		result = append(result, s)
	}
	return result
}

// createManyMaps - 많은 맵을 생성하는 함수
func createManyMaps() []map[string]int {
	result := make([]map[string]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		m := make(map[string]int)
		for j := 0; j < 100; j++ {
			key := fmt.Sprintf("key_%d", j)
			m[key] = j * i
		}
		result = append(result, m)
	}
	return result
}

// printMemStats - 현재 메모리 사용량을 출력한다
func printMemStats(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("[%s]\n", label)
	fmt.Printf("  Alloc = %d KB (현재 힙에 할당된 메모리)\n", m.Alloc/1024)
	fmt.Printf("  TotalAlloc = %d KB (총 할당된 메모리)\n", m.TotalAlloc/1024)
	fmt.Printf("  Sys = %d KB (OS로부터 할당받은 메모리)\n", m.Sys/1024)
	fmt.Printf("  NumGC = %d (GC 실행 횟수)\n", m.NumGC)
	fmt.Println()
}

func main() {
	// 시작 시 메모리 상태 출력
	printMemStats("시작")

	// === 메모리를 많이 사용하는 작업 수행 ===

	fmt.Println("=== append로 슬라이스 생성 ===")
	data1 := allocateWithAppend()
	printMemStats("append 후")
	_ = data1

	fmt.Println("=== 용량 미리 할당 후 슬라이스 생성 ===")
	data2 := allocateWithCapacity()
	printMemStats("capacity 할당 후")
	_ = data2

	fmt.Println("=== 문자열 대량 생성 ===")
	strings := createManyStrings()
	printMemStats("문자열 생성 후")
	_ = strings

	fmt.Println("=== 맵 대량 생성 ===")
	maps := createManyMaps()
	printMemStats("맵 생성 후")
	_ = maps

	// === 메모리 프로파일 저장 ===

	// GC를 실행하여 사용하지 않는 메모리를 해제한다
	runtime.GC()
	printMemStats("GC 후")

	// 메모리 프로파일 파일 생성
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("메모리 프로파일 파일 생성 실패:", err)
	}
	defer f.Close()

	// 힙 프로파일 기록
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("메모리 프로파일 기록 실패:", err)
	}

	fmt.Println("메모리 프로파일이 mem.prof에 저장되었습니다.")
	fmt.Println("다음 명령으로 분석하세요:")
	fmt.Println("  go tool pprof mem.prof")
	fmt.Println("  go tool pprof -alloc_space mem.prof  (총 할당량 기준)")
	fmt.Println("  go tool pprof -inuse_space mem.prof  (현재 사용량 기준)")
}
