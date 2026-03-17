# =============================================================================
# Go 언어 학습 환경 Dockerfile
# =============================================================================
# 이 Dockerfile은 Go 언어 학습에 필요한 모든 도구를 포함한 컨테이너 이미지를 생성한다.
# golang:1.22-alpine 기반으로 디버거, 린터, gRPC 도구 등을 함께 설치한다.
# =============================================================================

FROM golang:1.22-alpine

# -----------------------------------------------------------------------------
# 시스템 패키지 설치
# -----------------------------------------------------------------------------
RUN apk add --no-cache \
    git \
    curl \
    vim \
    bash \
    make \
    protobuf \
    protobuf-dev \
    tree \
    gcc \
    musl-dev

# -----------------------------------------------------------------------------
# Go 개발 도구 설치
# -----------------------------------------------------------------------------
# dlv: Go 디버거
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# gopls: Go Language Server (VS Code 등 에디터 연동)
RUN go install golang.org/x/tools/gopls@latest

# staticcheck: 정적 분석 도구
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

# golangci-lint: 통합 린터
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2

# -----------------------------------------------------------------------------
# gRPC / Protocol Buffers 도구 설치 (ch30_project_grpc_chat 용)
# -----------------------------------------------------------------------------
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# -----------------------------------------------------------------------------
# 작업 디렉토리 설정
# -----------------------------------------------------------------------------
WORKDIR /workspace

# 프로젝트 파일 복사
COPY . .

# Go 모듈 초기화 (go.mod가 없는 경우)
RUN if [ ! -f go.mod ]; then \
        go mod init edu_go && \
        go mod tidy; \
    fi

# -----------------------------------------------------------------------------
# 기본 쉘: bash
# -----------------------------------------------------------------------------
CMD ["bash"]
