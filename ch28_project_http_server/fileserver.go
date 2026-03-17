// 28장: HTTP 웹 서버 만들기 - 정적 파일 서버 예제
// 실행: go run fileserver.go
// 접속: http://localhost:8083/
//
// 이 프로그램은 ./static 디렉토리의 파일을 서비스한다.
// 실행 전에 static 디렉토리와 샘플 파일을 생성한다.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// createSampleFiles - 예제용 정적 파일을 생성한다
func createSampleFiles() {
	// static 디렉토리 생성
	os.MkdirAll("static/css", 0755)
	os.MkdirAll("static/js", 0755)

	// HTML 파일 생성
	htmlContent := `<!DOCTYPE html>
<html lang="ko">
<head>
    <meta charset="UTF-8">
    <title>Go 파일 서버</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <h1>Go 정적 파일 서버</h1>
        <p>이 페이지는 Go의 http.FileServer로 서비스됩니다.</p>
        <ul>
            <li>HTML, CSS, JavaScript 파일을 서비스합니다</li>
            <li>http.StripPrefix로 URL 접두사를 제거합니다</li>
            <li>디렉토리 목록도 자동으로 표시됩니다</li>
        </ul>
        <p id="time"></p>
    </div>
    <script src="/static/js/app.js"></script>
</body>
</html>`

	// CSS 파일 생성
	cssContent := `/* Go 파일 서버 스타일시트 */
body {
    font-family: 'Noto Sans KR', sans-serif;
    background-color: #f5f5f5;
    margin: 0;
    padding: 20px;
}

.container {
    max-width: 800px;
    margin: 0 auto;
    background: white;
    padding: 30px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

h1 {
    color: #00ADD8; /* Go 파란색 */
}

ul {
    line-height: 1.8;
}
`

	// JavaScript 파일 생성
	jsContent := `// Go 파일 서버 JavaScript
document.addEventListener('DOMContentLoaded', function() {
    var timeEl = document.getElementById('time');
    if (timeEl) {
        timeEl.textContent = '페이지 로드 시간: ' + new Date().toLocaleString('ko-KR');
    }
    console.log('Go 파일 서버 예제가 로드되었습니다.');
});
`

	// 파일 작성
	os.WriteFile("static/index.html", []byte(htmlContent), 0644)
	os.WriteFile("static/css/style.css", []byte(cssContent), 0644)
	os.WriteFile("static/js/app.js", []byte(jsContent), 0644)

	fmt.Println("샘플 정적 파일이 생성되었습니다:")
	fmt.Println("  static/index.html")
	fmt.Println("  static/css/style.css")
	fmt.Println("  static/js/app.js")
}

// loggingMiddleware - 요청 로깅 미들웨어
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// 다음 핸들러 호출
		next.ServeHTTP(w, r)
		// 요청 로그 출력
		log.Printf("[%s] %s %s (%v)",
			r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func main() {
	// 샘플 파일 생성
	createSampleFiles()

	mux := http.NewServeMux()

	// === 정적 파일 서버 설정 ===

	// http.FileServer: 디렉토리의 파일을 HTTP로 서비스한다
	// http.Dir("./static"): 파일 시스템의 ./static 디렉토리를 지정한다
	fileServer := http.FileServer(http.Dir("./static"))

	// http.StripPrefix: URL에서 /static/ 접두사를 제거한다
	// 요청: /static/css/style.css -> 파일: ./static/css/style.css
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// 루트 경로: index.html 서비스
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})

	// API 엔드포인트 (파일 서버와 함께 사용)
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"running","time":"%s"}`,
			time.Now().Format(time.RFC3339))
	})

	// 로깅 미들웨어 적용
	handler := loggingMiddleware(mux)

	addr := ":8083"
	fmt.Printf("\n파일 서버 시작: http://localhost%s\n", addr)
	fmt.Printf("정적 파일: http://localhost%s/static/\n", addr)
	fmt.Printf("API: http://localhost%s/api/status\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
