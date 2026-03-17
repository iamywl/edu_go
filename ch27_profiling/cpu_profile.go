// 27장: 프로파일링 - CPU 프로파일링 예제
// 실행: go run cpu_profile.go
// 분석: go tool pprof cpu.prof
//
// pprof 명령어:
//
//	top10  - CPU 사용량 상위 10개 함수
//	list 함수명 - 특정 함수의 라인별 CPU 사용량
//	web    - 호출 그래프 시각화 (Graphviz 필요)
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"time"
)

// bubbleSort - 버블 정렬 (CPU를 많이 사용하는 함수)
func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

// isPrime - 소수 판별 함수
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	limit := int(math.Sqrt(float64(n)))
	for i := 2; i <= limit; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// findPrimes - 주어진 범위에서 소수를 찾습니다
func findPrimes(max int) []int {
	var primes []int
	for i := 2; i <= max; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

// generateRandomSlice - 의사 랜덤 슬라이스를 생성한다
func generateRandomSlice(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		// 간단한 의사 난수 생성 (프로파일링 용도)
		slice[i] = (i*17 + 31) % size
	}
	return slice
}

func main() {
	// === CPU 프로파일 파일 생성 ===
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("CPU 프로파일 파일 생성 실패:", err)
	}
	defer f.Close()

	// === CPU 프로파일링 시작 ===
	fmt.Println("CPU 프로파일링을 시작합니다...")
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("CPU 프로파일링 시작 실패:", err)
	}
	defer pprof.StopCPUProfile()

	// === 프로파일링할 작업 1: 버블 정렬 ===
	fmt.Println("\n[작업 1] 버블 정렬 실행 중...")
	start := time.Now()

	data := generateRandomSlice(10000)
	bubbleSort(data)

	fmt.Printf("버블 정렬 완료: %v\n", time.Since(start))

	// === 프로파일링할 작업 2: 표준 정렬과 비교 ===
	fmt.Println("\n[작업 2] 표준 라이브러리 정렬 실행 중...")
	start = time.Now()

	data2 := generateRandomSlice(10000)
	sort.Ints(data2)

	fmt.Printf("표준 정렬 완료: %v\n", time.Since(start))

	// === 프로파일링할 작업 3: 소수 찾기 ===
	fmt.Println("\n[작업 3] 소수 찾기 실행 중...")
	start = time.Now()

	primes := findPrimes(100000)

	fmt.Printf("소수 찾기 완료: %d개 소수 발견 (%v)\n",
		len(primes), time.Since(start))

	// === 프로파일링 완료 ===
	fmt.Println("\n프로파일링이 완료되었습니다.")
	fmt.Println("다음 명령으로 분석하세요:")
	fmt.Println("  go tool pprof cpu.prof")
	fmt.Println("  (pprof) top10")
	fmt.Println("  (pprof) list bubbleSort")
	fmt.Println("  (pprof) web  (Graphviz 필요)")
}
