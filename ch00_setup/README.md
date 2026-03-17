# Chapter 00: 개발 환경 구축

Go 언어를 배우기 위한 첫 번째 단계는 개발 환경을 구축하는 것이다. 이 프로젝트는 Docker 기반으로 구성되어 있어, Docker만 설치하면 Go를 별도로 설치하지 않아도 바로 학습을 시작할 수 있다.

이 장에서는 Docker를 이용한 개발 환경 구축부터 첫 번째 Go 프로그램 실행까지의 전체 과정을 다룬다.

---

## 0.1 Docker Desktop 설치 (macOS)

### Homebrew를 이용한 설치 (권장)

```bash
# Homebrew가 설치되어 있지 않다면 먼저 설치한다
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Docker Desktop 설치
brew install --cask docker
```

### 수동 설치

1. [Docker Desktop 공식 사이트](https://www.docker.com/products/docker-desktop/)에서 macOS용 설치 파일을 다운로드한다.
2. 다운로드한 `.dmg` 파일을 열고, Docker 아이콘을 Applications 폴더로 드래그한다.
3. Applications에서 Docker를 실행한다.

### 설치 확인

Docker Desktop을 실행한 후, 터미널에서 다음 명령어로 설치를 확인한다:

```bash
docker --version
# Docker version 24.x.x, build xxxxxxx

docker compose version
# Docker Compose version v2.x.x
```

상단 메뉴바에 고래 모양 아이콘이 나타나고 "Docker Desktop is running" 상태이면 정상적으로 설치된 것이다.

---

## 0.2 Docker 기본 개념

Docker를 사용하기 전에 핵심 개념 세 가지를 이해해야 한다.

### 이미지 (Image)

이미지는 컨테이너를 만들기 위한 템플릿이다. 운영체제, 프로그래밍 언어, 도구 등이 포함된 일종의 스냅샷이라고 생각하면 된다. 이 프로젝트에서는 `golang:1.22-alpine` 이미지를 기반으로 Go 개발에 필요한 도구를 추가한 커스텀 이미지를 사용한다.

### 컨테이너 (Container)

컨테이너는 이미지를 실행한 인스턴스이다. 가상 머신과 비슷하지만 훨씬 가볍고 빠르다. 컨테이너 내부에서 Go 코드를 작성하고 실행할 수 있다.

### 볼륨 (Volume)

볼륨은 호스트(내 Mac)와 컨테이너 사이에 파일을 공유하는 방법이다. 이 프로젝트에서는 프로젝트 디렉토리를 컨테이너의 `/workspace`에 마운트한다. 따라서 Mac에서 코드를 수정하면 컨테이너 내부에 즉시 반영된다.

```
┌──────────────────┐          ┌──────────────────────┐
│   macOS (호스트)   │          │   Docker 컨테이너      │
│                    │  볼륨    │                        │
│  ~/edu_go/  ──────────────── /workspace/             │
│  (코드 편집)       │  마운트   │  (코드 실행)            │
└──────────────────┘          └──────────────────────┘
```

---

## 0.3 프로젝트 Docker 환경 빌드 및 실행

### Step 1: 프로젝트 디렉토리로 이동

```bash
cd edu_go
```

### Step 2: Docker 이미지 빌드

```bash
make build
# 또는
docker compose build
```

첫 빌드는 Go 도구(디버거, 린터 등)를 설치하므로 5~10분 정도 소요될 수 있다. 이후 빌드는 캐시를 사용하므로 빠르다.

### Step 3: 컨테이너 시작

```bash
make up
# 또는
docker compose up -d
```

`-d` 플래그는 백그라운드에서 실행하는 옵션이다.

### Step 4: 컨테이너 접속

```bash
make shell
# 또는
docker exec -it go-edu-container bash
```

접속에 성공하면 프롬프트가 다음과 같이 변경된다:

```
bash-5.2#
```

이제 컨테이너 내부에서 Go 명령어를 실행할 수 있다.

### Step 5: 컨테이너 종료

학습을 마치면 컨테이너를 종료한다:

```bash
# 컨테이너 내부에서 나가기
exit

# 컨테이너 중지 및 제거
make down
```

---

## 0.4 컨테이너 내부에서 Go 확인

컨테이너에 접속한 후 Go가 정상적으로 설치되어 있는지 확인한다.

### Go 버전 확인

```bash
go version
# go version go1.22.x linux/amd64 (또는 linux/arm64)
```

### Go 환경 변수 확인

```bash
go env
```

주요 환경 변수는 다음과 같다:

| 환경 변수 | 설명 | 컨테이너 내 값 |
|-----------|------|----------------|
| `GOROOT` | Go가 설치된 경로 | `/usr/local/go` |
| `GOPATH` | Go 워크스페이스 경로 | `/go` |
| `GOBIN` | `go install`로 설치한 바이너리 경로 | `/go/bin` |
| `GOOS` | 운영체제 | `linux` |
| `GOARCH` | CPU 아키텍처 | `amd64` 또는 `arm64` |

### 설치된 도구 확인

```bash
# 디버거
dlv version

# Language Server
gopls version

# 정적 분석 도구
staticcheck --version

# 린터
golangci-lint --version

# Protocol Buffers 컴파일러
protoc --version
```

---

## 0.5 VS Code + Dev Containers 연동 (선택 사항)

VS Code에서 Docker 컨테이너에 직접 연결하여 개발할 수 있다. 이 방법을 사용하면 코드 자동 완성, 디버깅 등 VS Code의 모든 기능을 컨테이너 환경에서 사용할 수 있다.

### 설정 방법

1. VS Code에서 **Dev Containers** 확장을 설치한다 (제작자: Microsoft).
2. `Cmd + Shift + P`를 눌러 명령 팔레트를 연다.
3. "Dev Containers: Attach to Running Container"를 선택한다.
4. `go-edu-container`를 선택한다.
5. 새 VS Code 창이 열리면 `/workspace` 폴더를 연다.

### 추천 VS Code 확장 (컨테이너 내부)

컨테이너에 연결된 VS Code에서 다음 확장을 설치하면 편리하다:

- **Go** (Go Team at Google): Go 언어 지원
- **Go Test Explorer**: 테스트 탐색기
- **Error Lens**: 에러를 코드 옆에 바로 표시

### 추천 설정 (settings.json)

```json
{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "editor.formatOnSave": true,
    "[go]": {
        "editor.defaultFormatter": "golang.go",
        "editor.codeActionsOnSave": {
            "source.organizeImports": "explicit"
        }
    }
}
```

---

## 0.6 첫 번째 Go 프로그램 실행

컨테이너에 접속한 상태에서 첫 번째 프로그램을 실행해 본다.

### 설치 확인 프로그램 실행

```bash
go run ch00_setup/check_install.go
```

### 예상 출력

```
Go 설치가 정상적으로 완료되었습니다!
Go 버전: go1.22.x
운영체제: linux
아키텍처: amd64
```

컨테이너 내부에서 실행하므로 운영체제가 `linux`으로 표시되는 것이 정상이다.

### 직접 프로그램 작성해 보기

호스트(Mac)에서 아무 텍스트 에디터로 새 파일을 만들어도 되고, 컨테이너 내부에서 `vim`으로 작성해도 된다.

```bash
# 컨테이너 내부에서 vim으로 작성
vim ch00_setup/hello.go
```

```go
package main

import "fmt"

func main() {
    fmt.Println("Docker 환경에서 Go를 실행하고 있다!")
}
```

```bash
go run ch00_setup/hello.go
# 출력: Docker 환경에서 Go를 실행하고 있다!
```

---

## 0.7 Go 모듈 이해

### Go Modules란

Go Modules는 Go의 공식 의존성 관리 시스템이다. 프로젝트에서 사용하는 외부 패키지의 버전을 관리하고, 재현 가능한 빌드를 보장한다.

이 프로젝트의 Docker 이미지에는 이미 `go.mod`가 초기화되어 있으므로 별도의 설정이 필요 없다.

### go.mod 파일 확인

```bash
cat go.mod
```

```
module edu_go

go 1.22
```

- `module`: 모듈 이름이다
- `go`: 사용하는 Go 버전이다

외부 패키지를 사용하는 챕터(ch29의 Gin, ch30의 gRPC 등)에서는 `go mod tidy` 명령으로 의존성을 자동으로 다운로드한다.

---

## 유용한 명령어 모음

### Docker 관련

| 명령어 | 설명 |
|--------|------|
| `make build` | Docker 이미지를 빌드한다 |
| `make up` | 컨테이너를 시작한다 |
| `make shell` | 컨테이너에 접속한다 |
| `make down` | 컨테이너를 중지한다 |
| `make clean` | 컨테이너와 이미지를 모두 제거한다 |

### Go 관련

| 명령어 | 설명 |
|--------|------|
| `go run` | 소스 코드를 컴파일하고 바로 실행한다 |
| `go build` | 소스 코드를 컴파일하여 실행 파일을 생성한다 |
| `go fmt` | 코드 포맷팅을 수행한다 |
| `go vet` | 코드의 잠재적 오류를 검사한다 |
| `go mod init` | 새 모듈을 초기화한다 |
| `go mod tidy` | 의존성을 정리한다 |
| `go test` | 테스트를 실행한다 |
| `go doc` | 패키지/함수 문서를 조회한다 |

---

## 부록: Docker 없이 로컬 설치하기

Docker를 사용하지 않고 직접 Go를 설치하려면 다음 방법을 따른다.

### macOS (Homebrew)

```bash
brew install go
```

### macOS (공식 패키지)

1. [Go 공식 다운로드 페이지](https://go.dev/dl/)에서 macOS용 `.pkg` 파일을 다운로드한다.
2. 패키지를 실행하여 설치한다.

### Windows

1. [Go 공식 다운로드 페이지](https://go.dev/dl/)에서 Windows용 `.msi` 파일을 다운로드한다.
2. 설치 파일을 실행하고 안내에 따라 설치한다.
3. 기본 설치 경로는 `C:\Program Files\Go`이다.

### Linux

```bash
# 최신 버전 다운로드
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

# 기존 설치 제거 후 압축 해제
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# PATH 설정 (~/.bashrc 또는 ~/.zshrc에 추가)
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin
```

### 설치 확인

```bash
go version
go env GOROOT GOPATH
```

### GOPATH 설정

| 환경 변수 | 설명 | 기본값 |
|-----------|------|--------|
| `GOROOT` | Go 설치 경로 | `/usr/local/go` |
| `GOPATH` | Go 워크스페이스 경로 | `$HOME/go` |
| `GOBIN` | 바이너리 설치 경로 | `$GOPATH/bin` |

`GOPATH/bin`을 시스템 PATH에 추가해 두면 `go install`로 설치한 도구를 어디서든 실행할 수 있다:

```bash
# ~/.zshrc에 추가
export PATH=$PATH:$(go env GOPATH)/bin
```

### VS Code 설정 (로컬)

1. [VS Code](https://code.visualstudio.com/)를 설치한다.
2. 확장 탭에서 **Go** 확장을 설치한다 (Go Team at Google).
3. VS Code를 재시작하면 Go 도구 설치를 묻는 알림이 나타난다. **Install All**을 클릭한다.

---

## 연습문제

1. Docker Desktop을 설치하고, `docker --version`과 `docker compose version` 명령으로 설치된 버전을 확인하라.
2. `make build` 명령으로 Docker 이미지를 빌드하고, `docker images` 명령으로 생성된 이미지를 확인하라. 이미지 크기가 얼마인지 확인하라.
3. 컨테이너에 접속하여 `go version`과 `go env`를 실행하고, `GOROOT`, `GOPATH`, `GOOS`, `GOARCH`의 값을 확인하라. 호스트(macOS)에서의 값과 어떻게 다른지 비교하라.
4. 컨테이너 내부에서 `which dlv`, `which gopls`, `which staticcheck` 명령을 실행하여 Go 개발 도구가 어디에 설치되어 있는지 확인하라.
5. 호스트(Mac)에서 아무 텍스트 에디터로 `ch00_setup/` 디렉토리에 새 `.go` 파일을 만들고, 컨테이너 내부에서 해당 파일이 즉시 보이는지 확인하라. 이것이 가능한 이유를 Docker 볼륨 마운트 개념으로 설명하라.
6. `docker ps`, `docker images`, `docker logs go-edu-container` 명령을 각각 실행하고, 출력 결과의 각 컬럼이 무엇을 의미하는지 조사하라.
7. `go doc fmt Println` 명령을 컨테이너 내부에서 실행하여 `fmt.Println` 함수의 공식 문서를 읽어 보라. 반환값이 무엇인지 확인하라.
8. `go fmt`와 `go vet`의 차이점을 조사하고, 각각 어떤 상황에서 사용하는지 정리하라. 의도적으로 포맷이 엉망인 코드를 작성하여 `go fmt`로 정리해 보라.
9. `tree -L 1 /workspace` 명령을 컨테이너 내부에서 실행하여 프로젝트 전체 구조를 확인하라. 각 디렉토리가 어떤 주제를 다루는지 README.md를 참고하여 정리하라.
10. `docker compose down`과 `docker compose down --rmi all`의 차이점을 조사하라. 각각 실행 후 `docker ps -a`와 `docker images`로 상태를 확인하여 차이를 설명하라.

---

## 구현 과제

1. **시스템 정보 출력 프로그램**: `runtime` 패키지를 활용하여 Go 버전, 운영체제, 아키텍처, CPU 개수(`runtime.NumCPU()`), 고루틴 개수(`runtime.NumGoroutine()`)를 모두 출력하는 프로그램을 작성하라. Docker 컨테이너 내부에서 실행하여 결과를 확인하라.

2. **환경 변수 비교 프로그램**: `os` 패키지의 `os.Getenv()` 함수를 사용하여 `GOPATH`, `GOROOT`, `HOME`, `PATH` 환경 변수를 읽어 출력하는 프로그램을 작성하라. 같은 코드를 호스트(로컬 Go 설치 시)와 Docker 컨테이너에서 각각 실행하여 값이 어떻게 다른지 비교하라.

3. **Docker 환경 감지 프로그램**: Go 프로그램이 Docker 컨테이너 내부에서 실행되고 있는지 감지하는 프로그램을 작성하라. `/.dockerenv` 파일의 존재 여부를 확인하거나, `os.Hostname()`의 결과를 활용할 수 있다.

4. **개발 도구 자동 검증 프로그램**: `os/exec` 패키지를 사용하여 `go`, `dlv`, `gopls`, `staticcheck`, `golangci-lint`, `protoc` 명령이 모두 사용 가능한지 자동으로 검증하는 프로그램을 작성하라. 각 도구의 버전 정보도 함께 출력하라.

5. **크로스 컴파일 실험**: 컨테이너 내부에서 `GOOS`와 `GOARCH` 환경 변수를 설정하여 다른 운영체제용 바이너리를 빌드해 보라. 예를 들어 `GOOS=darwin GOARCH=arm64 go build` 명령으로 macOS용 바이너리를 생성하고, `file` 명령으로 바이너리의 형식을 확인하라.

---

## 프로젝트 과제

1. **개발 환경 헬스체크 대시보드**: Go 개발 환경의 전체 상태를 한눈에 보여주는 프로그램을 작성하라. Go 버전 확인, 필수 도구 설치 여부, 디스크 공간, 메모리 사용량, 네트워크 연결 상태 등을 점검하고, 결과를 보기 좋은 표 형식으로 출력해야 한다. 모든 항목이 통과하면 "환경 준비 완료", 실패 항목이 있으면 해결 방법을 함께 안내하라.

2. **Go 프로젝트 템플릿 생성기**: 프로젝트 이름을 입력받아 디렉토리를 생성하고, `go mod init`을 실행하며, 기본 `main.go`, `.gitignore`, `Dockerfile`, `Makefile`을 자동으로 생성하는 프로그램을 작성하라. 생성된 프로젝트가 바로 `go run`으로 실행 가능해야 한다.

---

다음 장에서는 Go 언어의 역사와 첫 번째 "Hello World" 프로그램을 작성해 본다.
