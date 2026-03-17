# 25장 [Project] 단어 검색 프로그램 만들기

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 순차 버전 실행 (검색어와 경로를 인수로 전달한다)
go run ch25_project_word_search/main.go <검색어> <경로>
# 예시: go run ch25_project_word_search/main.go fmt .

# 병렬 버전 실행
go run ch25_project_word_search/main_concurrent.go <검색어> <경로>
# 예시: go run ch25_project_word_search/main_concurrent.go error .
```

> **Makefile 활용**: `make run CH=ch25_project_word_search` 또는 `make run CH=ch25_project_word_search FILE=main.go`

---

난이도: ★★☆☆

이 프로젝트에서는 지금까지 배운 내용을 종합하여 **파일에서 특정 단어를 검색하는 프로그램**을 만든다. Unix의 `grep` 명령어와 유사한 기능을 Go로 구현하며, 순차 버전과 병렬 버전을 모두 작성하여 성능 차이를 체감할 수 있다.

---

## 25.1 해법 (설계)

### 프로그램 기능

```
$ go run main.go <검색어> <경로>
```

- 지정된 경로에서 파일을 재귀적으로 탐색한다.
- 각 파일을 열어서 한 줄씩 읽으며 검색어를 찾는다.
- 검색어가 포함된 줄을 파일명, 줄번호와 함께 출력한다.

프로그램의 핵심은 크게 세 단계로 나뉜다: (1) 명령줄 인수 파싱, (2) 파일 목록 수집, (3) 파일 내용 검색. 각 단계를 독립적인 함수로 분리하여 테스트 가능하고 재사용 가능한 구조로 설계한다.

### 프로그램 구조

```
main()
├── parseArgs()         // 실행 인수 파싱
├── getFileList()       // 파일 목록 수집
└── searchFiles()       // 파일 검색 실행
    └── searchFile()    // 개별 파일 검색
```

### 출력 형식

```
파일경로:줄번호: 내용
```

예시:
```
main.go:15: func main() {
main.go:23:     fmt.Println("Hello, World!")
```

이 형식은 `grep`의 기본 출력 형식과 동일하며, 많은 텍스트 편집기와 IDE가 이 형식을 인식하여 해당 위치로 바로 이동할 수 있도록 지원한다.

---

## 25.2 사전 지식

### os.Args

프로그램의 실행 인수를 읽는다. `os.Args`는 문자열 슬라이스이며, 첫 번째 요소(`os.Args[0]`)는 실행 파일의 경로이다.

```go
import "os"

// 실행: go run main.go hello world
// os.Args[0] = 실행 파일 경로
// os.Args[1] = "hello"
// os.Args[2] = "world"
```

더 복잡한 명령줄 인수 처리가 필요하면 `flag` 패키지를 사용할 수 있다. 하지만 이 프로젝트에서는 `os.Args`로 충분하다.

### filepath 패키지

파일 경로를 다루는 유틸리티를 제공한다. `filepath.Walk`는 디렉토리를 재귀적으로 탐색하며, 각 파일/디렉토리에 대해 콜백 함수를 호출한다.

```go
import "path/filepath"

// 디렉토리를 재귀적으로 탐색
filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
        fmt.Println(path)
    }
    return nil
})
```

> **참고:** Go 1.16부터는 `filepath.WalkDir`이 추가되었다. `Walk`보다 효율적이며, `os.DirEntry` 인터페이스를 사용하여 불필요한 `Stat` 호출을 줄인다.

### bufio 패키지

버퍼링된 I/O를 제공한다. `bufio.Scanner`는 파일을 한 줄씩 읽을 때 매우 유용하며, 대용량 파일도 메모리를 효율적으로 사용하면서 읽을 수 있다.

```go
import "bufio"

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text() // 한 줄씩 읽기
    fmt.Println(line)
}
if err := scanner.Err(); err != nil {
    fmt.Println("읽기 에러:", err)
}
```

`Scanner`의 기본 버퍼 크기는 64KB이다. 한 줄이 이보다 길면 에러가 발생하므로, 필요에 따라 `scanner.Buffer()`로 버퍼 크기를 조절할 수 있다.

### strings 패키지

문자열 관련 유틸리티를 제공한다. 이 프로젝트에서는 `strings.Contains`를 사용하여 줄에 검색어가 포함되어 있는지 확인한다.

```go
import "strings"

strings.Contains("Hello, World!", "World")  // true
strings.Contains("Hello, World!", "world")  // false (대소문자 구분)
strings.ToLower("Hello")                     // "hello"
```

---

## 25.3 실행 인수 읽고 파일 목록 가져오기

### 1단계: 인수 파싱

```go
func parseArgs() (string, string, error) {
    if len(os.Args) < 3 {
        return "", "", fmt.Errorf("사용법: %s <검색어> <경로>", os.Args[0])
    }
    return os.Args[1], os.Args[2], nil
}
```

이 함수는 최소 2개의 인수(검색어, 경로)가 필요하다. 인수가 부족하면 사용법을 안내하는 에러를 반환한다.

### 2단계: 파일 목록 수집

```go
func getFileList(root string) ([]string, error) {
    var files []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}
```

이 함수는 지정된 디렉토리를 재귀적으로 탐색하여 모든 파일의 경로를 슬라이스로 반환한다. 디렉토리는 제외하고 파일만 수집한다.

---

## 25.4 파일을 열어서 라인 읽기

### 개별 파일 검색

검색 결과를 구조체로 정의하면, 결과를 구조화하여 다루기 편리하다.

```go
type SearchResult struct {
    FilePath string
    LineNum  int
    Line     string
}

func searchFile(filePath, keyword string) ([]SearchResult, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var results []SearchResult
    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        if strings.Contains(line, keyword) {
            results = append(results, SearchResult{
                FilePath: filePath,
                LineNum:  lineNum,
                Line:     line,
            })
        }
    }
    return results, scanner.Err()
}
```

`defer file.Close()`는 함수가 어떤 경로로 반환되든 파일이 반드시 닫히도록 보장한다. 파일을 닫지 않으면 파일 디스크립터가 누수되어, 대량의 파일을 처리할 때 "too many open files" 에러가 발생할 수 있다.

---

## 25.5 파일 검색 프로그램 완성하기

모든 파일에 대해 `searchFile`을 호출하고 결과를 출력한다.

```go
func searchFiles(files []string, keyword string) {
    totalMatches := 0
    for _, file := range files {
        results, err := searchFile(file, keyword)
        if err != nil {
            fmt.Fprintf(os.Stderr, "경고: %s 읽기 실패: %v\n", file, err)
            continue
        }
        for _, r := range results {
            fmt.Printf("%s:%d: %s\n", r.FilePath, r.LineNum, r.Line)
            totalMatches++
        }
    }
    fmt.Printf("\n총 %d개의 결과를 찾았습니다.\n", totalMatches)
}
```

에러가 발생한 파일은 건너뛰고 다음 파일을 계속 처리한다. 에러 메시지는 `os.Stderr`로 출력하여, 검색 결과와 에러 메시지가 섞이지 않도록 한다. 이는 파이프라인에서 결과만 다른 프로그램으로 전달할 때 유용하다.

완성된 코드는 `main.go`를 참고하라.

---

## 25.6 개선하기 (고루틴으로 병렬 검색)

### 병렬 검색의 장점

- 대량의 파일을 동시에 검색하여 성능을 향상시킨다.
- I/O 대기 시간(파일 읽기)을 고루틴 간에 중첩시켜 전체 처리 시간을 단축한다.
- 고루틴과 채널을 사용하여 자연스럽게 구현할 수 있다.

### 병렬 버전 설계

```
main()
├── parseArgs()
├── getFileList()
└── searchFilesConcurrent()
    ├── 고루틴 1: searchFile("file1.go") ─┐
    ├── 고루틴 2: searchFile("file2.go") ─┼─→ 결과 채널 → 출력
    ├── 고루틴 3: searchFile("file3.go") ─┘
    └── WaitGroup으로 완료 대기
```

각 파일을 별도의 고루틴에서 검색하되, 동시에 열리는 파일 수를 제한하여 시스템 리소스 고갈을 방지한다. 이를 **세마포어 패턴**이라 한다.

### 핵심 구현

```go
func searchFilesConcurrent(files []string, keyword string) {
    results := make(chan SearchResult)
    var wg sync.WaitGroup

    // 최대 동시 실행 수 제한 (세마포어)
    sem := make(chan struct{}, 10)

    for _, file := range files {
        wg.Add(1)
        go func(path string) {
            defer wg.Done()
            sem <- struct{}{}        // 슬롯 획득
            defer func() { <-sem }() // 슬롯 반환

            found, err := searchFile(path, keyword)
            if err != nil { return }
            for _, r := range found {
                results <- r
            }
        }(file)
    }

    // 모든 고루틴이 완료되면 채널을 닫음
    go func() {
        wg.Wait()
        close(results)
    }()

    // 결과 출력
    for r := range results {
        fmt.Printf("%s:%d: %s\n", r.FilePath, r.LineNum, r.Line)
    }
}
```

이 구현에서 세마포어 크기(10)는 동시에 열 수 있는 최대 파일 수를 의미한다. 이 값은 시스템의 파일 디스크립터 제한과 디스크 I/O 특성에 따라 조절할 수 있다.

완성된 코드는 `main_concurrent.go`를 참고하라.

### 실행 방법

```bash
# 순차 버전
go run main.go "fmt" .

# 병렬 버전
go run main_concurrent.go "fmt" .

# 특정 디렉토리에서 검색
go run main.go "error" /path/to/project
```

---

## 학습 포인트 정리

이 프로젝트에서 활용한 개념들:

| 개념 | 사용 위치 |
|------|----------|
| `os.Args` | 실행 인수 읽기 |
| `filepath.Walk` | 디렉토리 재귀 탐색 |
| `bufio.Scanner` | 파일을 한 줄씩 읽기 |
| `strings.Contains` | 문자열 검색 |
| 에러 처리 | 파일 열기/읽기 에러 처리 |
| 구조체 | SearchResult로 결과 캡슐화 |
| 고루틴 | 파일 병렬 검색 |
| 채널 | 검색 결과 전달 |
| WaitGroup | 고루틴 완료 대기 |
| 세마포어 패턴 | 동시 실행 수 제한 |

---

## 연습문제

### 문제 1: 대소문자 무시 검색
검색 시 대소문자를 구분하지 않는 옵션(`-i` 플래그)을 추가하라.
- `strings.ToLower`를 활용하여 검색어와 줄 모두를 소문자로 변환한 뒤 비교하라.

### 문제 2: 정규식 검색
검색어 대신 정규표현식으로 검색하는 기능을 구현하라.
- `regexp` 패키지를 사용하라.
- 정규식이 유효하지 않으면 에러 메시지를 출력하고 종료하라.

### 문제 3: 파일 확장자 필터
특정 확장자의 파일만 검색하는 옵션(`-ext .go,.txt`)을 추가하라.
- `filepath.Ext` 함수를 사용하여 파일 확장자를 확인하라.
- 여러 확장자를 쉼표로 구분하여 지정할 수 있도록 하라.

### 문제 4: 줄 번호 주변 컨텍스트
검색 결과에서 일치하는 줄의 전후 N줄을 함께 출력하는 기능을 구현하라.
- `-C 3` 옵션으로 전후 3줄을 출력하라 (`grep`의 `-C` 옵션과 동일).
- 연속된 결과가 겹치면 하나로 합쳐서 출력하라.

### 문제 5: 검색 결과 통계
검색 완료 후 다음 통계를 출력하라.
- 검색한 총 파일 수
- 일치하는 파일 수
- 총 일치하는 줄 수
- 소요 시간 (`time.Since` 사용)

### 문제 6: 바이너리 파일 건너뛰기
바이너리 파일(텍스트가 아닌 파일)을 자동으로 건너뛰는 기능을 구현하라.
- 파일의 처음 512바이트를 읽어 NUL 바이트(`\x00`)가 포함되어 있으면 바이너리로 판단하라.
- 바이너리 파일은 건너뛰고 경고 메시지를 출력하라.

### 문제 7: 검색 결과 색상 강조
터미널에서 검색어가 포함된 부분을 색상으로 강조하여 출력하라.
- ANSI 이스케이프 코드를 사용하라 (예: `\033[31m`은 빨간색).
- 검색어 부분만 빨간색으로 표시하라.

### 문제 8: 순차 vs 병렬 성능 비교
순차 버전과 병렬 버전의 성능을 비교하는 벤치마크를 작성하라.
- 큰 디렉토리(예: Go 표준 라이브러리 소스 코드)에서 동일한 검색어로 테스트하라.
- 세마포어 크기(1, 5, 10, 20, 50)에 따른 성능 변화를 측정하라.
- 결과를 표 형태로 출력하라.

### 문제 9: 제외 패턴 지원
특정 디렉토리나 파일을 검색에서 제외하는 기능을 추가하라.
- `--exclude-dir .git,node_modules` 옵션으로 디렉토리를 제외하라.
- `--exclude *.log,*.tmp` 옵션으로 특정 패턴의 파일을 제외하라.
- `.gitignore` 파일을 읽어 자동으로 제외 패턴을 적용하는 기능을 추가하라 (보너스).

### 문제 10: 검색 결과 정렬
병렬 검색의 경우 결과 순서가 비결정적이다. 결과를 파일 경로와 줄 번호 기준으로 정렬하여 출력하는 옵션(`--sort`)을 추가하라.
- 모든 결과를 먼저 수집한 뒤 정렬하여 출력하라.
- `slices.SortFunc`를 사용하라.

---

## 구현 과제

### 과제 1: flag 패키지로 CLI 개선
현재 `os.Args`로 처리하는 명령줄 인수를 `flag` 패키지를 사용하도록 개선하라.
- `-i`: 대소문자 무시
- `-r`: 재귀 검색 (기본: 활성화)
- `-n`: 줄 번호 출력 (기본: 활성화)
- `-c`: 일치 횟수만 출력
- `-l`: 일치하는 파일명만 출력
- `-w`: 단어 단위 일치 (검색어가 독립된 단어로 존재할 때만 일치)
- `--help`로 사용법을 출력하라.

### 과제 2: 검색 결과를 JSON으로 출력
검색 결과를 JSON 형식으로 출력하는 옵션(`--json`)을 추가하라.
- 각 결과를 `{"file": "...", "line": N, "content": "..."}` 형태로 출력하라.
- 전체 결과를 JSON 배열로 출력하라.
- `encoding/json` 패키지를 사용하라.

### 과제 3: 검색 결과 캐싱
동일한 파일을 반복 검색할 때 속도를 높이기 위한 캐싱 메커니즘을 구현하라.
- 파일의 수정 시간(`ModTime`)과 크기를 기준으로 캐시 유효성을 검사하라.
- 캐시된 파일 내용은 메모리에 저장하되, 최대 캐시 크기를 제한하라.
- 캐시 적중률을 출력하라.

### 과제 4: 대체(Replace) 기능
검색한 결과를 다른 문자열로 대체하는 기능을 구현하라.
- `--replace <대체문자열>` 옵션을 추가하라.
- 실제 파일을 수정하기 전에 변경 내용을 미리보기로 출력하라.
- `--dry-run` 옵션을 추가하여 실제 수정 없이 미리보기만 하는 기능을 제공하라.
- 변경 전 백업 파일을 생성하라 (`.bak` 확장자).

### 과제 5: context를 활용한 타임아웃과 취소
검색에 타임아웃을 설정하고, 사용자가 Ctrl+C로 취소할 수 있도록 개선하라.
- `--timeout 10s` 옵션으로 전체 검색 시간을 제한하라.
- `os/signal` 패키지로 SIGINT 신호를 잡아 검색을 취소하라.
- 취소 시 현재까지의 결과를 출력하고 종료하라.
- `context.WithTimeout`과 `context.WithCancel`을 적절히 활용하라.

---

## 프로젝트 과제

### 프로젝트 1: 미니 grep 도구 (mygrep)
실제 `grep` 명령어의 주요 기능을 구현한 완전한 CLI 도구를 만들어라.
- 지원 기능: 대소문자 무시(`-i`), 줄 번호(`-n`), 일치 횟수(`-c`), 반전 검색(`-v`), 파일명만 출력(`-l`), 재귀 검색(`-r`), 전후 문맥(`-C`), 정규식 지원, 색상 강조
- `flag` 패키지로 명령줄 인수를 처리하라.
- 병렬 검색을 기본으로 하되, 순차 검색 옵션도 제공하라.
- 파이프 입력(`stdin`)도 지원하라: `cat file.txt | mygrep "pattern"`
- 바이너리 파일을 자동으로 건너뛰라.
- 사용법(`--help`)과 버전 정보(`--version`)를 출력하라.

### 프로젝트 2: 코드 검색 엔진
프로젝트 디렉토리에서 코드를 검색하는 특화된 검색 도구를 만들어라.
- `.gitignore` 규칙을 자동으로 적용하여 불필요한 파일을 건너뛰라.
- 프로그래밍 언어별 구문 인식: 함수 정의 검색(`--func`), 구조체/클래스 검색(`--type`), 주석 검색(`--comment`)
- 검색 결과에 파일 타입(Go, Python, JavaScript 등)을 표시하라.
- 프로젝트 전체의 통계를 출력하라: 파일 수, 코드 줄 수, 언어별 분포.
- 고루틴으로 병렬 검색하되, `context.WithCancel`로 취소 기능을 구현하라.
- 자주 검색하는 프로젝트의 파일 목록을 인덱싱하여 재검색 시 속도를 향상시켜라.
