// 28장: HTTP 웹 서버 만들기 - 기본 HTTP 서버 예제
// 실행: go run main.go
// 테스트: curl http://localhost:8080/
//
//	curl http://localhost:8080/hello
//	curl http://localhost:8080/time
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// homeHandler - 홈 페이지 핸들러
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// r.URL.Path가 "/"인 경우만 처리 (다른 경로는 404)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// 응답 헤더 설정
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTML 응답 작성
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Go HTTP 서버</title></head>
<body>
<h1>Go HTTP 서버에 오신 것을 환영합니다!</h1>
<ul>
	<li><a href="/hello">인사하기</a></li>
	<li><a href="/hello?name=Go">이름으로 인사하기</a></li>
	<li><a href="/time">현재 시간</a></li>
	<li><a href="/info">요청 정보</a></li>
</ul>
</body>
</html>`)
}

// helloHandler - 인사 핸들러
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 쿼리 파라미터에서 이름을 가져옵니다
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "안녕하세요, %s님!\n", name)
}

// timeHandler - 현재 시간을 반환하는 핸들러
func timeHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "현재 시간: %s\n", now.Format("2006-01-02 15:04:05"))
}

// infoHandler - 요청 정보를 표시하는 핸들러
func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "=== 요청 정보 ===\n")
	fmt.Fprintf(w, "메서드: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %s\n", r.URL.String())
	fmt.Fprintf(w, "경로: %s\n", r.URL.Path)
	fmt.Fprintf(w, "프로토콜: %s\n", r.Proto)
	fmt.Fprintf(w, "호스트: %s\n", r.Host)
	fmt.Fprintf(w, "클라이언트: %s\n", r.RemoteAddr)

	fmt.Fprintf(w, "\n=== 요청 헤더 ===\n")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", name, value)
		}
	}
}

func main() {
	// ServeMux를 직접 생성한다 (DefaultServeMux 대신)
	mux := http.NewServeMux()

	// 핸들러 등록
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/time", timeHandler)
	mux.HandleFunc("/info", infoHandler)

	// 서버 설정
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("HTTP 서버를 시작합니다: http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
