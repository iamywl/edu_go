package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================
// 25장 단어 검색 프로그램 (순차 버전)
// 지정된 경로에서 파일을 재귀적으로 탐색하고,
// 각 파일에서 검색어가 포함된 줄을 출력한다.
//
// 사용법: go run main.go <검색어> <경로>
// 예시:   go run main.go "fmt" .
// ============================================================

// SearchResult - 검색 결과를 담는 구조체
type SearchResult struct {
	FilePath string // 파일 경로
	LineNum  int    // 줄 번호
	Line     string // 해당 줄의 내용
}

// parseArgs - 실행 인수를 파싱한다.
// 검색어와 경로를 반환한다.
func parseArgs() (string, string, error) {
	if len(os.Args) < 3 {
		return "", "", fmt.Errorf(
			"사용법: %s <검색어> <경로>\n"+
				"예시: %s \"fmt\" .",
			os.Args[0], os.Args[0],
		)
	}
	keyword := os.Args[1]
	root := os.Args[2]

	// 경로가 존재하는지 확인
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return "", "", fmt.Errorf("경로가 존재하지 않습니다: %s", root)
	}

	return keyword, root, nil
}

// getFileList - 지정된 경로에서 파일 목록을 재귀적으로 수집한다.
// 숨김 디렉토리(.git 등)와 바이너리 파일은 건너뜁니다.
func getFileList(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 접근 권한이 없는 파일/디렉토리는 건너뜁니다.
			fmt.Fprintf(os.Stderr, "경고: %s 접근 실패: %v\n", path, err)
			return nil
		}

		// 숨김 디렉토리 건너뛰기 (.git, .vscode 등)
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
			return filepath.SkipDir
		}

		// 디렉토리는 건너뛰기
		if info.IsDir() {
			return nil
		}

		// 텍스트 파일만 포함 (확장자로 판단)
		ext := strings.ToLower(filepath.Ext(path))
		textExts := map[string]bool{
			".go": true, ".txt": true, ".md": true, ".json": true,
			".yaml": true, ".yml": true, ".toml": true, ".xml": true,
			".html": true, ".css": true, ".js": true, ".ts": true,
			".py": true, ".rb": true, ".java": true, ".c": true,
			".h": true, ".cpp": true, ".rs": true, ".sh": true,
			".sql": true, ".csv": true, ".log": true, ".cfg": true,
			".ini": true, ".env": true, ".mod": true, ".sum": true,
		}

		// 확장자가 없는 파일이나 알려진 텍스트 파일만 포함
		if ext == "" || textExts[ext] {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// searchFile - 개별 파일에서 검색어를 찾습니다.
// 검색어가 포함된 줄의 정보를 SearchResult 슬라이스로 반환한다.
func searchFile(filePath, keyword string) ([]SearchResult, error) {
	// 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	var results []SearchResult
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// 한 줄씩 읽으며 검색
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 검색어가 포함되어 있는지 확인
		if strings.Contains(line, keyword) {
			results = append(results, SearchResult{
				FilePath: filePath,
				LineNum:  lineNum,
				Line:     line,
			})
		}
	}

	// 스캐너 에러 확인
	if err := scanner.Err(); err != nil {
		return results, fmt.Errorf("파일 읽기 에러: %w", err)
	}

	return results, nil
}

// searchFiles - 모든 파일에서 검색어를 순차적으로 찾습니다.
func searchFiles(files []string, keyword string) int {
	totalMatches := 0
	filesWithMatches := 0

	for _, file := range files {
		results, err := searchFile(file, keyword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "경고: %v\n", err)
			continue
		}

		if len(results) > 0 {
			filesWithMatches++
		}

		for _, r := range results {
			// 줄 내용에서 검색어를 강조하여 출력
			fmt.Printf("\033[35m%s\033[0m:\033[32m%d\033[0m: %s\n",
				r.FilePath, r.LineNum, highlightKeyword(r.Line, keyword))
			totalMatches++
		}
	}

	return totalMatches
}

// highlightKeyword - 검색어를 강조한다. (ANSI 색상 코드 사용)
func highlightKeyword(line, keyword string) string {
	// 검색어를 빨간색으로 강조
	highlighted := strings.ReplaceAll(line, keyword,
		"\033[1;31m"+keyword+"\033[0m")
	return highlighted
}

func main() {
	// 1단계: 실행 인수 파싱
	keyword, root, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "에러:", err)
		os.Exit(1)
	}

	fmt.Printf("검색어: \"%s\"\n", keyword)
	fmt.Printf("검색 경로: %s\n", root)
	fmt.Println(strings.Repeat("-", 50))

	// 2단계: 파일 목록 수집
	files, err := getFileList(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "파일 목록 수집 실패:", err)
		os.Exit(1)
	}
	fmt.Printf("검색 대상 파일: %d개\n\n", len(files))

	// 3단계: 파일 검색 실행
	totalMatches := searchFiles(files, keyword)

	// 4단계: 결과 요약
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("총 %d개의 파일에서 %d개의 결과를 찾았습니다.\n", len(files), totalMatches)
}
