package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// =============================================
	// bufio 및 io 패키지 활용 예제
	// =============================================

	// --- bufio.Scanner: 줄 단위 읽기 ---
	fmt.Println("=== Scanner로 문자열 줄 단위 읽기 ===")
	input := "첫 번째 줄\n두 번째 줄\n세 번째 줄"
	scanner := bufio.NewScanner(strings.NewReader(input))

	lineNum := 1
	for scanner.Scan() {
		fmt.Printf("  %d: %s\n", lineNum, scanner.Text())
		lineNum++
	}
	// 스캔 중 에러가 발생했는지 확인
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "읽기 에러: %v\n", err)
	}

	// --- bufio.Scanner: 단어 단위 읽기 ---
	fmt.Println("\n=== Scanner로 단어 단위 읽기 ===")
	wordInput := "Go 언어는 정말 재미있습니다"
	wordScanner := bufio.NewScanner(strings.NewReader(wordInput))
	wordScanner.Split(bufio.ScanWords) // 기본은 ScanLines, 여기서는 ScanWords 사용

	fmt.Print("  단어들: ")
	for wordScanner.Scan() {
		fmt.Printf("[%s] ", wordScanner.Text())
	}
	fmt.Println()

	// --- bufio.Writer: 버퍼링된 쓰기 ---
	fmt.Println("\n=== bufio.Writer 버퍼링된 쓰기 ===")
	writer := bufio.NewWriter(os.Stdout)

	// 버퍼에 데이터를 씁니다 (아직 출력되지 않을 수 있음)
	fmt.Fprint(writer, "  버퍼링된 ")
	fmt.Fprint(writer, "출력입니다. ")
	fmt.Fprintln(writer, "Flush 후 화면에 나타납니다.")

	// Flush를 호출해야 버퍼의 내용이 실제로 출력됩니다
	writer.Flush()

	// --- bufio.ReadWriter: 읽기/쓰기 결합 ---
	fmt.Println("\n=== bufio.NewReadWriter ===")
	var buf strings.Builder
	r := bufio.NewReader(strings.NewReader("Hello, Go!\n"))
	w := bufio.NewWriter(&buf)
	rw := bufio.NewReadWriter(r, w)

	// ReadWriter로 읽기
	line, _ := rw.ReadString('\n')
	fmt.Printf("  읽은 내용: %s", line)

	// ReadWriter로 쓰기
	rw.WriteString("새로운 내용")
	rw.Flush()
	fmt.Printf("  쓴 내용: %s\n", buf.String())

	// --- io.Copy: Reader에서 Writer로 복사 ---
	fmt.Println("\n=== io.Copy: Reader → Writer ===")
	src := strings.NewReader("  io.Copy로 복사된 텍스트입니다.\n")
	// src(Reader) → os.Stdout(Writer)로 복사
	io.Copy(os.Stdout, src)

	// --- io.MultiReader: 여러 Reader 연결 ---
	fmt.Println("\n=== io.MultiReader: 여러 소스 연결 ===")
	r1 := strings.NewReader("  안녕하세요, ")
	r2 := strings.NewReader("Go ")
	r3 := strings.NewReader("세계!\n")
	multi := io.MultiReader(r1, r2, r3)
	io.Copy(os.Stdout, multi)

	// --- io.TeeReader: 읽으면서 동시에 기록 ---
	fmt.Println("\n=== io.TeeReader: 읽으면서 기록 ===")
	original := strings.NewReader("TeeReader 테스트 데이터")
	var captured strings.Builder

	// original에서 읽을 때마다 captured에도 동시에 기록
	tee := io.TeeReader(original, &captured)
	data, _ := io.ReadAll(tee)

	fmt.Printf("  읽은 데이터: %s\n", string(data))
	fmt.Printf("  캡처된 데이터: %s\n", captured.String())

	// --- io.LimitReader: 읽기 제한 ---
	fmt.Println("\n=== io.LimitReader: 읽기 크기 제한 ===")
	longText := strings.NewReader("이것은 매우 긴 텍스트입니다. 일부만 읽겠습니다.")
	limited := io.LimitReader(longText, 30) // 최대 30바이트만 읽기
	limitedData, _ := io.ReadAll(limited)
	fmt.Printf("  제한된 읽기 (30바이트): %s\n", string(limitedData))

	fmt.Println("\n프로그램 종료")
}
