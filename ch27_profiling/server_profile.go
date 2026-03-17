// 27장: 프로파일링 - HTTP 서버에서 프로파일링 예제
// 실행: go run server_profile.go
// 브라우저에서 http://localhost:6060/debug/pprof/ 접속
//
// 프로파일 수집 및 분석:
//
//	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10
//	go tool pprof http://localhost:6060/debug/pprof/heap
//	go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap
package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof" // 빈 import: pprof HTTP 핸들러 자동 등록
	"sync"
)

// 메모리를 사용하는 전역 캐시 (프로파일링 대상)
var (
	cache   = make(map[string][]byte)
	cacheMu sync.RWMutex
)

// handleIndex - 인덱스 페이지 핸들러
func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>
<head><title>프로파일링 데모 서버</title></head>
<body>
<h1>프로파일링 데모 서버</h1>
<ul>
	<li><a href="/compute">/compute</a> - CPU 사용 작업</li>
	<li><a href="/alloc">/alloc</a> - 메모리 할당 작업</li>
	<li><a href="/cache?key=test&size=1024">/cache</a> - 캐시에 데이터 저장</li>
	<li><a href="/debug/pprof/">/debug/pprof/</a> - 프로파일링 대시보드</li>
</ul>
<h2>프로파일 분석 방법</h2>
<pre>
# CPU 프로파일 수집 (10초)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10

# 힙 메모리 프로파일
go tool pprof http://localhost:6060/debug/pprof/heap

# 고루틴 프로파일
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 웹 UI로 분석
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap
</pre>
</body>
</html>`)
}

// handleCompute - CPU를 많이 사용하는 핸들러 (소수 계산)
func handleCompute(w http.ResponseWriter, r *http.Request) {
	// 소수 계산으로 CPU를 사용한다
	count := 0
	for i := 2; i < 100000; i++ {
		if isPrimeNumber(i) {
			count++
		}
	}
	fmt.Fprintf(w, "100000 이하의 소수: %d개\n", count)
}

// isPrimeNumber - 소수 판별 함수
func isPrimeNumber(n int) bool {
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

// handleAlloc - 메모리를 많이 할당하는 핸들러
func handleAlloc(w http.ResponseWriter, r *http.Request) {
	// 대량의 문자열 슬라이스를 생성한다
	data := make([]string, 0, 10000)
	for i := 0; i < 10000; i++ {
		data = append(data, fmt.Sprintf("항목_%d: 이것은 메모리 할당 테스트 데이터입니다", i))
	}
	fmt.Fprintf(w, "메모리 할당 완료: %d개 항목 생성\n", len(data))
}

// handleCache - 캐시에 데이터를 저장하는 핸들러
func handleCache(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	sizeStr := r.URL.Query().Get("size")
	if key == "" {
		key = "default"
	}

	size := 1024 // 기본 1KB
	if sizeStr != "" {
		fmt.Sscanf(sizeStr, "%d", &size)
	}

	// 캐시에 데이터 저장
	cacheMu.Lock()
	cache[key] = make([]byte, size)
	cacheMu.Unlock()

	cacheMu.RLock()
	totalSize := 0
	for _, v := range cache {
		totalSize += len(v)
	}
	count := len(cache)
	cacheMu.RUnlock()

	fmt.Fprintf(w, "캐시 저장 완료: key=%s, size=%d bytes\n", key, size)
	fmt.Fprintf(w, "캐시 현황: %d개 항목, 총 %d bytes\n", count, totalSize)
}

func main() {
	// 핸들러 등록
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/compute", handleCompute)
	http.HandleFunc("/alloc", handleAlloc)
	http.HandleFunc("/cache", handleCache)

	// 서버 시작
	addr := ":6060"
	fmt.Printf("프로파일링 데모 서버 시작: http://localhost%s\n", addr)
	fmt.Printf("프로파일링 대시보드: http://localhost%s/debug/pprof/\n", addr)
	fmt.Println()
	fmt.Println("부하 테스트 예제:")
	fmt.Println("  for i in $(seq 1 100); do curl -s http://localhost:6060/compute > /dev/null; done")
	fmt.Println()
	fmt.Println("프로파일 수집:")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10")

	log.Fatal(http.ListenAndServe(addr, nil))
}
