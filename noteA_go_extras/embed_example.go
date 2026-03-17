package main

import (
	"embed"
	"fmt"
	"io/fs"
)

// =============================================
// embed 패키지 예제
// //go:embed 디렉티브로 파일을 바이너리에 포함
// =============================================

// 파일 내용을 string으로 임베드
// hello.txt 파일이 같은 디렉토리에 있어야 한다
//
//go:embed hello.txt
var helloMessage string

// 파일 내용을 []byte로 임베드 (바이너리 파일에 적합)
//
//go:embed hello.txt
var helloBytes []byte

// embed.FS로 임베드 (여러 파일, 디렉토리 지원)
//
//go:embed hello.txt
var helloFS embed.FS

func main() {
	// --- string으로 임베드된 파일 읽기 ---
	fmt.Println("=== string으로 임베드 ===")
	fmt.Printf("  내용: %s\n", helloMessage)
	fmt.Printf("  길이: %d 바이트\n", len(helloMessage))

	// --- []byte로 임베드된 파일 읽기 ---
	fmt.Println("\n=== []byte로 임베드 ===")
	fmt.Printf("  내용: %s\n", string(helloBytes))
	fmt.Printf("  길이: %d 바이트\n", len(helloBytes))
	fmt.Printf("  첫 5바이트 (hex): ")
	limit := 5
	if len(helloBytes) < limit {
		limit = len(helloBytes)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("0x%02x ", helloBytes[i])
	}
	fmt.Println()

	// --- embed.FS로 임베드된 파일 읽기 ---
	fmt.Println("\n=== embed.FS로 임베드 ===")

	// ReadFile로 개별 파일 읽기
	data, err := helloFS.ReadFile("hello.txt")
	if err != nil {
		fmt.Printf("  에러: %v\n", err)
	} else {
		fmt.Printf("  ReadFile 결과: %s\n", string(data))
	}

	// ReadDir로 디렉토리 내용 나열
	fmt.Println("\n=== embed.FS 디렉토리 탐색 ===")
	entries, err := helloFS.ReadDir(".")
	if err != nil {
		fmt.Printf("  에러: %v\n", err)
	} else {
		for _, entry := range entries {
			info, _ := entry.Info()
			fmt.Printf("  파일: %-15s 크기: %d 바이트\n", entry.Name(), info.Size())
		}
	}

	// --- fs.WalkDir로 재귀 탐색 ---
	fmt.Println("\n=== fs.WalkDir로 재귀 탐색 ===")
	fs.WalkDir(helloFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			fmt.Printf("  [디렉토리] %s\n", path)
		} else {
			info, _ := d.Info()
			fmt.Printf("  [파일]     %s (%d 바이트)\n", path, info.Size())
		}
		return nil
	})

	// --- 활용 팁 ---
	fmt.Println("\n=== embed 활용 팁 ===")
	fmt.Println("  1. 웹 서버의 정적 파일(HTML, CSS, JS)을 바이너리에 포함")
	fmt.Println("  2. SQL 마이그레이션 파일을 바이너리에 포함")
	fmt.Println("  3. 설정 파일 템플릿을 바이너리에 포함")
	fmt.Println("  4. 버전 정보 파일을 바이너리에 포함")
	fmt.Println()
	fmt.Println("  주의: 큰 파일을 임베드하면 바이너리 크기가 커집니다!")
}
