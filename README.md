# Go 언어 완벽 가이드 - 기초부터 서버 개발까지

> Go 언어의 기본 문법부터 동시성 프로그래밍, HTTP 서버, gRPC까지 체계적으로 학습하는 프로젝트이다.

---

## 1. 프로젝트 소개

이 프로젝트는 Go 언어를 처음 배우는 개발자를 위한 체계적인 학습 자료이다. 총 31개의 챕터와 2개의 부록으로 구성되어 있으며, 기초 문법부터 실전 서버 개발까지 단계별로 학습할 수 있다.

각 챕터는 다음과 같은 구조를 가진다:
- **README.md**: 개념 설명과 예제 코드
- **.go 파일**: 직접 실행할 수 있는 예시 코드
- **연습문제**: 개념 확인 문제
- **구현과제**: 직접 코드를 작성하는 과제
- **프로젝트과제**: 종합적인 미니 프로젝트

### 참고 문서

- **[GO_CODING_STANDARDS.md](GO_CODING_STANDARDS.md)** - Go 코드 작성 표준 가이드 (Python의 PEP 8에 해당). Effective Go, Go Code Review Comments, Go Proverbs 등 Go 공식 스타일 가이드를 한국어로 정리한 문서이다. 코드를 작성하기 전에 반드시 읽어보는 것을 권장한다.
- **[STUDY_PLAN.md](STUDY_PLAN.md)** - 학습 계획표. 하루 2시간 기준 8주 완성 코스, 속성 4주 코스, 느긋한 12주 코스를 제공한다. Day 1부터 Day 56까지 일별 세부 학습 내용과 시간 배분이 포함되어 있다.

모든 학습 환경은 Docker로 구성되어 있어, Go를 별도로 설치할 필요 없이 Docker만 있으면 바로 학습을 시작할 수 있다.

---

## 2. 사전 준비

### Docker Desktop 설치 (macOS)

```bash
# Homebrew가 설치되어 있지 않다면 먼저 설치한다
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Docker Desktop 설치
brew install --cask docker
```

설치 후 Docker Desktop 앱을 실행하고, 상단 메뉴바에 Docker 아이콘(고래 모양)이 나타나면 준비가 완료된 것이다.

### 요구 사항

| 항목 | 최소 요구 |
|------|-----------|
| Docker Desktop | 4.0 이상 |
| 디스크 공간 | 약 2GB (이미지 포함) |
| 메모리 | 4GB 이상 권장 |

---

## 3. 빠른 시작 (Quick Start)

### Step 1: 프로젝트 클론/다운로드

```bash
git clone <repository-url>
cd edu_go
```

### Step 2: Docker 이미지 빌드

```bash
make build
# 또는
docker compose build
```

### Step 3: 컨테이너 실행

```bash
make up
# 또는
docker compose up -d
```

### Step 4: 컨테이너 접속

```bash
make shell
# 또는
docker exec -it go-edu-container bash
```

### Step 5: 첫 번째 예제 실행

```bash
# 컨테이너 내부에서 실행한다
go run ch01_hello_go/main.go
```

---

## 4. 학습 방법

각 챕터를 다음 순서로 학습하는 것을 권장한다:

```
1. README.md 읽기     →  개념과 문법을 이해한다
2. 예시 코드 실행      →  코드를 직접 실행하고 결과를 확인한다
3. 연습문제 풀기       →  개념을 제대로 이해했는지 점검한다
4. 구현과제 수행       →  직접 코드를 작성하며 실력을 키운다
5. 프로젝트과제 도전   →  종합적인 프로그램을 만들어 본다
```

챕터 순서대로 진행하되, 이미 알고 있는 내용은 건너뛰어도 무방하다. 단, ch15, ch25, ch28~ch30은 이전 챕터의 내용을 종합하는 프로젝트 챕터이므로 순서를 지키는 것이 좋다.

---

## 5. 챕터별 실행 가이드

### 전체 챕터 목록

| 챕터 | 제목 | 난이도 | 주요 파일 | 실행 명령어 |
|------|------|--------|-----------|-------------|
| ch00 | 개발 환경 구축 | ⭐ | `check_install.go` | `go run ch00_setup/check_install.go` |
| ch01 | Hello Go | ⭐ | `main.go` | `go run ch01_hello_go/main.go` |
| ch02 | 변수 | ⭐ | `main.go`, `scope.go`, `type_conversion.go` | `go run ch02_variables/main.go` |
| ch03 | fmt 패키지 | ⭐ | `print_examples.go`, `formatting.go`, `scan_examples.go` | `go run ch03_fmt/print_examples.go` |
| ch04 | 연산자 | ⭐ | `arithmetic.go`, `comparison.go`, `float_error.go` | `go run ch04_operators/arithmetic.go` |
| ch05 | 함수 | ⭐ | `basic.go`, `multi_return.go`, `recursion.go` | `go run ch05_functions/basic.go` |
| ch06 | 상수 | ⭐ | `main.go`, `iota_example.go` | `go run ch06_constants/main.go` |
| ch07 | if 조건문 | ⭐ | `basic.go`, `short_statement.go` | `go run ch07_if/basic.go` |
| ch08 | switch 문 | ⭐ | `basic.go`, `advanced.go`, `enum_switch.go` | `go run ch08_switch/basic.go` |
| ch09 | for 반복문 | ⭐ | `basic.go`, `gugudan.go`, `label.go` | `go run ch09_for/basic.go` |
| ch10 | 배열 | ⭐⭐ | `basic.go`, `multidim.go` | `go run ch10_array/basic.go` |
| ch11 | 구조체 | ⭐⭐ | `basic.go`, `embedded.go`, `padding.go` | `go run ch11_struct/basic.go` |
| ch12 | 포인터 | ⭐⭐ | `basic.go`, `instance.go`, `why_pointer.go` | `go run ch12_pointer/basic.go` |
| ch13 | 문자열 | ⭐⭐ | `basic.go`, `rune.go`, `builder.go` | `go run ch13_string/basic.go` |
| ch14 | 패키지 | ⭐⭐ | `main.go` | `go run ch14_package/main.go` |
| ch15 | 프로젝트: 숫자 맞추기 게임 | ⭐⭐ | `main.go` | `go run ch15_project_guessing_game/main.go` |
| ch16 | 슬라이스 | ⭐⭐ | `basic.go`, `internal.go`, `slicing.go` | `go run ch16_slice/basic.go` |
| ch17 | 메서드 | ⭐⭐ | `basic.go`, `pointer_vs_value.go` | `go run ch17_method/basic.go` |
| ch18 | 인터페이스 | ⭐⭐⭐ | `basic.go`, `duck_typing.go`, `empty_interface.go` | `go run ch18_interface/basic.go` |
| ch19 | 고급 함수 | ⭐⭐⭐ | `closure.go`, `defer_example.go`, `func_type.go` | `go run ch19_functions_advanced/closure.go` |
| ch20 | 자료구조 | ⭐⭐⭐ | `map_example.go`, `list_example.go`, `ring_example.go` | `go run ch20_data_structures/map_example.go` |
| ch21 | 에러 처리 | ⭐⭐⭐ | `basic.go`, `custom_error.go`, `panic_recover.go` | `go run ch21_error_handling/basic.go` |
| ch22 | 고루틴 | ⭐⭐⭐ | `basic.go`, `mutex.go`, `deadlock.go` | `go run ch22_goroutine/basic.go` |
| ch23 | 채널과 컨텍스트 | ⭐⭐⭐ | `channel_basic.go`, `buffered.go`, `context_cancel.go` | `go run ch23_channel_context/channel_basic.go` |
| ch24 | 제네릭 | ⭐⭐⭐ | `basic.go`, `constraints.go`, `generic_type.go` | `go run ch24_generics/basic.go` |
| ch25 | 프로젝트: 단어 검색기 | ⭐⭐⭐ | `main.go`, `main_concurrent.go` | `go run ch25_project_word_search/main.go` |
| ch26 | 테스팅 | ⭐⭐⭐ | `calculator.go`, `calculator_test.go` | `go test -v ./ch26_testing/...` |
| ch27 | 프로파일링 | ⭐⭐⭐⭐ | `cpu_profile.go`, `mem_profile.go`, `server_profile.go` | `go run ch27_profiling/cpu_profile.go` |
| ch28 | 프로젝트: HTTP 서버 | ⭐⭐⭐⭐ | `main.go`, `json_handler.go`, `fileserver.go` | `go run ch28_project_http_server/main.go` |
| ch29 | 프로젝트: RESTful API | ⭐⭐⭐⭐ | `main.go`, `gin_server.go` | `go run ch29_project_restful_api/main.go` |
| ch30 | 프로젝트: gRPC 채팅 | ⭐⭐⭐⭐⭐ | `chat_server.go`, `chat_client.go`, `echo_client.go` | `go run ch30_project_grpc_chat/chat_server.go` |
| noteA | Go 추가 기능 | ⭐⭐ | - | - |
| noteB | 프로그래밍 사고방식 | ⭐ | - | - |

### 챕터 실행 방법 (컨테이너 내부)

```bash
# 직접 실행
go run ch01_hello_go/main.go

# Makefile 사용 (호스트에서 실행)
make run CH=ch01_hello_go FILE=main.go

# 스크립트 사용 (컨테이너 내부)
bash scripts/run_chapter.sh ch01_hello_go main.go
```

---

## 6. 유용한 Docker 명령어

```bash
# 이미지 빌드
make build

# 컨테이너 시작
make up

# 컨테이너 접속
make shell

# 컨테이너 중지
make down

# 컨테이너와 이미지 모두 제거
make clean

# 컨테이너 상태 확인
docker ps

# 컨테이너 로그 확인
docker logs go-edu-container

# 이미지 목록 확인
docker images
```

---

## 7. 유용한 Go 명령어

```bash
# 소스 코드 실행
go run main.go

# 바이너리 빌드
go build -o myapp main.go

# 테스트 실행
go test -v ./...

# 코드 포맷팅
go fmt ./...

# 정적 분석
go vet ./...

# 모듈 의존성 정리
go mod tidy

# 패키지 문서 확인
go doc fmt Println

# Go 환경 변수 확인
go env

# 외부 패키지 설치
go get github.com/example/package

# 벤치마크 실행
go test -bench=. ./ch26_testing/
```

---

## 8. 디렉토리 구조

```
edu_go/
├── Dockerfile                  # Docker 이미지 정의
├── docker-compose.yml          # Docker Compose 설정
├── Makefile                    # 빌드/실행 명령어 모음
├── README.md                   # 프로젝트 안내 (현재 파일)
├── .dockerignore               # Docker 빌드 제외 파일
├── scripts/
│   └── run_chapter.sh          # 챕터 실행 헬퍼 스크립트
│
├── ch00_setup/                 # 개발 환경 구축
├── ch01_hello_go/              # Hello Go
├── ch02_variables/             # 변수
├── ch03_fmt/                   # fmt 패키지
├── ch04_operators/             # 연산자
├── ch05_functions/             # 함수
├── ch06_constants/             # 상수
├── ch07_if/                    # if 조건문
├── ch08_switch/                # switch 문
├── ch09_for/                   # for 반복문
├── ch10_array/                 # 배열
├── ch11_struct/                # 구조체
├── ch12_pointer/               # 포인터
├── ch13_string/                # 문자열
├── ch14_package/               # 패키지
├── ch15_project_guessing_game/ # 프로젝트: 숫자 맞추기 게임
├── ch16_slice/                 # 슬라이스
├── ch17_method/                # 메서드
├── ch18_interface/             # 인터페이스
├── ch19_functions_advanced/    # 고급 함수
├── ch20_data_structures/       # 자료구조 (맵, 리스트, 링)
├── ch21_error_handling/        # 에러 처리
├── ch22_goroutine/             # 고루틴
├── ch23_channel_context/       # 채널과 컨텍스트
├── ch24_generics/              # 제네릭
├── ch25_project_word_search/   # 프로젝트: 단어 검색기
├── ch26_testing/               # 테스팅
├── ch27_profiling/             # 프로파일링
├── ch28_project_http_server/   # 프로젝트: HTTP 서버
├── ch29_project_restful_api/   # 프로젝트: RESTful API
├── ch30_project_grpc_chat/     # 프로젝트: gRPC 채팅
├── noteA_go_extras/            # 부록 A: Go 추가 기능
└── noteB_thinking_programming/ # 부록 B: 프로그래밍 사고방식
```

---

## 9. 커리큘럼 로드맵

```
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  [1단계: 기초]  ch00 ~ ch15                                         │
│  ─────────────────────────────                                      │
│  환경설정 → 변수 → 출력 → 연산자 → 함수 → 상수                      │
│  → 조건문 → 반복문 → 배열 → 구조체 → 포인터 → 문자열                │
│  → 패키지 → 🎮 숫자 맞추기 게임                                     │
│                          │                                          │
│                          ▼                                          │
│  [2단계: 고급]  ch16 ~ ch27                                         │
│  ─────────────────────────────                                      │
│  슬라이스 → 메서드 → 인터페이스 → 고급 함수 → 자료구조               │
│  → 에러 처리 → 고루틴 → 채널/컨텍스트 → 제네릭                      │
│  → 🔍 단어 검색기 → 테스팅 → 프로파일링                             │
│                          │                                          │
│                          ▼                                          │
│  [3단계: 서버 개발]  ch28 ~ ch30                                    │
│  ──────────────────────────────                                     │
│  🌐 HTTP 서버 → 🔄 RESTful API → 💬 gRPC 채팅                      │
│                                                                     │
│  [부록]                                                             │
│  noteA: Go 추가 기능  │  noteB: 프로그래밍 사고방식                  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 10. FAQ / 문제 해결

### Q: `make build`에서 오류가 발생한다

Docker Desktop이 실행 중인지 확인한다:
```bash
docker info
```
Docker Desktop 앱이 실행되고 있어야 한다. 메뉴바에서 고래 아이콘이 보이지 않으면 Docker Desktop을 실행한다.

### Q: `make shell`에서 "No such container" 오류가 발생한다

컨테이너가 실행 중인지 확인한다:
```bash
docker ps
```
컨테이너가 목록에 없으면 `make up`으로 먼저 시작한다.

### Q: 포트가 이미 사용 중이라는 오류가 발생한다

해당 포트를 사용하는 프로세스를 확인하고 종료한다:
```bash
# macOS
lsof -i :8080
kill -9 <PID>
```

또는 `docker-compose.yml`에서 포트 번호를 변경한다.

### Q: 컨테이너 내부에서 코드 변경이 반영되지 않는다

`docker-compose.yml`에서 볼륨 마운트(`.:/workspace`)가 설정되어 있으므로, 호스트에서 코드를 수정하면 컨테이너 내부에 즉시 반영된다. 반영되지 않는 경우 컨테이너를 재시작한다:
```bash
make down && make up
```

### Q: `go mod tidy`에서 오류가 발생한다

컨테이너 내부에서 실행해야 한다:
```bash
make shell
# 컨테이너 내부
go mod tidy
```

### Q: Apple Silicon (M1/M2/M3) Mac에서 빌드가 느리다

Docker Desktop의 설정에서 "Use Rosetta for x86/amd64 emulation on Apple Silicon"을 활성화하면 성능이 개선될 수 있다. golang:1.22-alpine 이미지는 ARM64를 기본 지원하므로 대부분의 경우 문제가 없다.

### Q: 디스크 공간이 부족하다

Docker 캐시를 정리한다:
```bash
docker system prune -a
```

---

## 라이선스

이 프로젝트는 학습 목적으로 작성되었다.
