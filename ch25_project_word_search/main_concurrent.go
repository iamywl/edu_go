package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ============================================================
// 25장 단어 검색 프로그램 (병렬 버전)
// 고루틴과 채널을 사용하여 여러 파일을 동시에 검색한다.
//
// 사용법: go run main_concurrent.go <검색어> <경로>
// 예시:   go run main_concurrent.go "fmt" .
// ============================================================

// SearchResultC - 검색 결과를 담는 구조체
type SearchResultC struct {
	FilePath string // 파일 경로
	LineNum  int    // 줄 번호
	Line     string // 해당 줄의 내용
}

// parseArgsC - 실행 인수를 파싱한다.
func parseArgsC() (string, string, error) {
	if len(os.Args) < 3 {
		return "", "", fmt.Errorf(
			"사용법: %s <검색어> <경로>\n"+
				"예시: %s \"fmt\" .",
			os.Args[0], os.Args[0],
		)
	}
	keyword := os.Args[1]
	root := os.Args[2]

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return "", "", fmt.Errorf("경로가 존재하지 않습니다: %s", root)
	}

	return keyword, root, nil
}

// getFileListC - 파일 목록을 재귀적으로 수집한다.
func getFileListC(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 접근 실패한 파일은 건너뜀
		}

		// 숨김 디렉토리 건너뛰기
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		// 텍스트 파일만 포함
		ext := strings.ToLower(filepath.Ext(path))
		textExts := map[string]bool{
			".go": true, ".txt": true, ".md": true, ".json": true,
			".yaml": true, ".yml": true, ".toml": true, ".xml": true,
			".html": true, ".css": true, ".js": true, ".ts": true,
			".py": true, ".rb": true, ".java": true, ".c": true,
			".h": true, ".cpp": true, ".rs": true, ".sh": true,
			".sql": true, ".csv": true, ".log": true, ".mod": true,
			".sum": true,
		}

		if ext == "" || textExts[ext] {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// searchFileC - 개별 파일에서 검색어를 찾습니다.
func searchFileC(filePath, keyword string) ([]SearchResultC, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []SearchResultC
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if strings.Contains(line, keyword) {
			results = append(results, SearchResultC{
				FilePath: filePath,
				LineNum:  lineNum,
				Line:     line,
			})
		}
	}

	return results, scanner.Err()
}

// searchFilesConcurrent - 고루틴을 사용하여 파일을 병렬로 검색한다.
func searchFilesConcurrent(files []string, keyword string) int {
	// 결과를 전달할 채널
	results := make(chan SearchResultC, 100)

	// WaitGroup으로 모든 고루틴의 완료를 추적한다.
	var wg sync.WaitGroup

	// 세마포어로 동시 실행 고루틴 수를 제한한다.
	// CPU 코어 수의 2배로 설정한다.
	maxWorkers := runtime.NumCPU() * 2
	sem := make(chan struct{}, maxWorkers)

	// 각 파일에 대해 고루틴을 시작한다.
	for _, file := range files {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			// 세마포어 슬롯 획득 (동시 실행 수 제한)
			sem <- struct{}{}
			defer func() { <-sem }()

			// 파일 검색 수행
			found, err := searchFileC(path, keyword)
			if err != nil {
				// 에러는 stderr로 출력
				fmt.Fprintf(os.Stderr, "경고: %s 읽기 실패: %v\n", path, err)
				return
			}

			// 결과를 채널로 전송
			for _, r := range found {
				results <- r
			}
		}(file)
	}

	// 모든 고루틴이 완료되면 결과 채널을 닫습니다.
	go func() {
		wg.Wait()
		close(results)
	}()

	// 결과를 수집하고 출력한다.
	totalMatches := 0
	for r := range results {
		// 검색어를 강조하여 출력
		highlighted := strings.ReplaceAll(r.Line, keyword,
			"\033[1;31m"+keyword+"\033[0m")
		fmt.Printf("\033[35m%s\033[0m:\033[32m%d\033[0m: %s\n",
			r.FilePath, r.LineNum, highlighted)
		totalMatches++
	}

	return totalMatches
}

func main() {
	// 1단계: 실행 인수 파싱
	keyword, root, err := parseArgsC()
	if err != nil {
		fmt.Fprintln(os.Stderr, "에러:", err)
		os.Exit(1)
	}

	fmt.Printf("검색어: \"%s\"\n", keyword)
	fmt.Printf("검색 경로: %s\n", root)
	fmt.Printf("워커 수: %d (CPU 코어의 2배)\n", runtime.NumCPU()*2)
	fmt.Println(strings.Repeat("-", 50))

	// 2단계: 파일 목록 수집
	files, err := getFileListC(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "파일 목록 수집 실패:", err)
		os.Exit(1)
	}
	fmt.Printf("검색 대상 파일: %d개\n\n", len(files))

	// 3단계: 병렬 검색 실행 (시간 측정)
	start := time.Now()
	totalMatches := searchFilesConcurrent(files, keyword)
	elapsed := time.Since(start)

	// 4단계: 결과 요약
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("총 %d개의 파일에서 %d개의 결과를 찾았습니다.\n", len(files), totalMatches)
	fmt.Printf("검색 소요 시간: %v\n", elapsed)
}
