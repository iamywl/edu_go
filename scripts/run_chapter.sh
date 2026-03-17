#!/bin/bash
# =============================================================================
# Go 챕터 실행 스크립트
# =============================================================================
# 사용법: ./scripts/run_chapter.sh <챕터_디렉토리> [파일명]
# 예시:
#   ./scripts/run_chapter.sh ch01_hello_go main.go
#   ./scripts/run_chapter.sh ch05_functions basic.go
#   ./scripts/run_chapter.sh ch01_hello_go          (기본: main.go)
# =============================================================================

set -e

# 색상 정의
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# 구분선
LINE="━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 인자 확인
if [ -z "$1" ]; then
    echo -e "${RED}오류: 챕터 디렉토리를 지정해야 한다.${NC}"
    echo ""
    echo "사용법: $0 <챕터_디렉토리> [파일명]"
    echo "예시:   $0 ch01_hello_go main.go"
    exit 1
fi

CHAPTER="$1"
FILE="${2:-main.go}"
FILEPATH="/workspace/${CHAPTER}/${FILE}"

# 파일 존재 확인
if [ ! -f "$FILEPATH" ]; then
    echo -e "${RED}오류: 파일을 찾을 수 없다: ${FILEPATH}${NC}"
    echo ""
    echo "사용 가능한 파일 목록:"
    ls /workspace/${CHAPTER}/*.go 2>/dev/null || echo "  (Go 파일 없음)"
    exit 1
fi

# 실행
echo ""
echo -e "${BLUE}${LINE}${NC}"
echo -e "${BOLD}  챕터: ${GREEN}${CHAPTER}${NC}"
echo -e "${BOLD}  파일: ${YELLOW}${FILE}${NC}"
echo -e "${BLUE}${LINE}${NC}"
echo ""

cd /workspace
go run "${CHAPTER}/${FILE}"

echo ""
echo -e "${BLUE}${LINE}${NC}"
echo -e "${GREEN}  실행 완료${NC}"
echo -e "${BLUE}${LINE}${NC}"
echo ""
