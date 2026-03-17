// 28장: HTTP 웹 서버 만들기 - 쿼리 파라미터 처리 예제
// 실행: go run query.go
// 테스트:
//
//	curl "http://localhost:8081/search?keyword=golang&page=1"
//	curl "http://localhost:8081/calc?op=add&a=10&b=20"
//	curl "http://localhost:8081/filter?tags=web&tags=go&tags=server"
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// searchHandler - 검색 쿼리 파라미터 처리
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// URL 쿼리 파라미터를 파싱한다
	query := r.URL.Query()

	keyword := query.Get("keyword") // 단일 값 가져오기
	page := query.Get("page")
	limit := query.Get("limit")

	// 기본값 설정
	if keyword == "" {
		http.Error(w, "keyword 파라미터가 필요합니다", http.StatusBadRequest)
		return
	}
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "검색 결과\n")
	fmt.Fprintf(w, "검색어: %s\n", keyword)
	fmt.Fprintf(w, "페이지: %s\n", page)
	fmt.Fprintf(w, "표시 개수: %s\n", limit)
}

// calcHandler - 계산기 쿼리 파라미터 처리
func calcHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	op := query.Get("op")
	aStr := query.Get("a")
	bStr := query.Get("b")

	// 파라미터 유효성 검사
	if op == "" || aStr == "" || bStr == "" {
		http.Error(w, "op, a, b 파라미터가 모두 필요합니다\n예: /calc?op=add&a=10&b=20",
			http.StatusBadRequest)
		return
	}

	// 문자열을 숫자로 변환
	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		http.Error(w, "a 값이 유효한 숫자가 아닙니다", http.StatusBadRequest)
		return
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		http.Error(w, "b 값이 유효한 숫자가 아닙니다", http.StatusBadRequest)
		return
	}

	// 연산 수행
	var result float64
	switch op {
	case "add":
		result = a + b
	case "sub":
		result = a - b
	case "mul":
		result = a * b
	case "div":
		if b == 0 {
			http.Error(w, "0으로 나눌 수 없습니다", http.StatusBadRequest)
			return
		}
		result = a / b
	default:
		http.Error(w, "지원하지 않는 연산입니다. 사용 가능: add, sub, mul, div",
			http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%.2f %s %.2f = %.2f\n", a, op, b, result)
}

// filterHandler - 다중 값 쿼리 파라미터 처리
func filterHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// 단일 값 파라미터
	category := query.Get("category")

	// 다중 값 파라미터 (같은 키로 여러 값)
	// 예: /filter?tags=web&tags=go&tags=server
	tags := query["tags"] // []string 반환

	// 값 존재 여부 확인
	_, hasMinPrice := query["min_price"]
	_, hasMaxPrice := query["max_price"]

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "=== 필터 조건 ===\n")
	fmt.Fprintf(w, "카테고리: %s\n", category)
	fmt.Fprintf(w, "태그: %v (%d개)\n", tags, len(tags))
	fmt.Fprintf(w, "최소가격 설정: %v\n", hasMinPrice)
	fmt.Fprintf(w, "최대가격 설정: %v\n", hasMaxPrice)
}

func main() {
	mux := http.NewServeMux()

	// 핸들러 등록
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/calc", calcHandler)
	mux.HandleFunc("/filter", filterHandler)

	// 인덱스 페이지
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<h1>쿼리 파라미터 예제</h1>
<ul>
<li><a href="/search?keyword=golang&page=1&limit=20">검색: golang</a></li>
<li><a href="/calc?op=add&a=10&b=20">계산: 10 + 20</a></li>
<li><a href="/calc?op=mul&a=7&b=8">계산: 7 * 8</a></li>
<li><a href="/filter?category=books&tags=go&tags=programming&tags=backend">필터: 다중 태그</a></li>
</ul>`)
	})

	addr := ":8081"
	fmt.Printf("쿼리 파라미터 예제 서버 시작: http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
