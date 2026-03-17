package main

import "fmt"

// === 색상 열거값 ===
type Color int

const (
	Red Color = iota
	Green
	Blue
	Yellow
)

// Color의 문자열 표현
func (c Color) String() string {
	switch c {
	case Red:
		return "빨강"
	case Green:
		return "초록"
	case Blue:
		return "파랑"
	case Yellow:
		return "노랑"
	default:
		return "알 수 없는 색"
	}
}

// === HTTP 상태 그룹 열거값 ===
type StatusGroup int

const (
	StatusSuccess     StatusGroup = iota // 0: 2xx 성공
	StatusClientError                    // 1: 4xx 클라이언트 에러
	StatusServerError                    // 2: 5xx 서버 에러
)

// 상태 코드에서 그룹을 반환
func getStatusGroup(code int) StatusGroup {
	switch {
	case code >= 200 && code < 300:
		return StatusSuccess
	case code >= 400 && code < 500:
		return StatusClientError
	case code >= 500 && code < 600:
		return StatusServerError
	default:
		return -1
	}
}

// 상태 그룹에 대한 메시지
func statusMessage(group StatusGroup) string {
	switch group {
	case StatusSuccess:
		return "요청이 성공적으로 처리되었습니다"
	case StatusClientError:
		return "클라이언트 요청에 오류가 있습니다"
	case StatusServerError:
		return "서버에서 오류가 발생했습니다"
	default:
		return "알 수 없는 상태입니다"
	}
}

// === 방향 열거값 ===
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// 방향을 90도 시계방향으로 회전
func (d Direction) TurnRight() Direction {
	return (d + 1) % 4
}

func (d Direction) String() string {
	names := [...]string{"북", "동", "남", "서"}
	if d < North || d > West {
		return "?"
	}
	return names[d]
}

func main() {
	// === 색상 열거값 ===
	fmt.Println("=== 색상 열거값 ===")
	colors := []Color{Red, Green, Blue, Yellow}

	for _, c := range colors {
		fmt.Printf("Color(%d) → %s\n", c, c)
	}

	// === HTTP 상태 코드 ===
	fmt.Println("\n=== HTTP 상태 코드 판별 ===")
	codes := []int{200, 201, 404, 403, 500, 503}

	for _, code := range codes {
		group := getStatusGroup(code)
		msg := statusMessage(group)
		fmt.Printf("HTTP %d → %s\n", code, msg)
	}

	// === 방향 + switch ===
	fmt.Println("\n=== 방향 회전 ===")
	dir := North
	fmt.Println("시작 방향:", dir)

	// 4번 회전하면 원래 방향으로 돌아옴
	for i := 0; i < 4; i++ {
		dir = dir.TurnRight()
		fmt.Printf("오른쪽 회전 → %s\n", dir)
	}

	// === switch로 방향에 따른 이동 ===
	fmt.Println("\n=== 방향별 이동 ===")
	directions := []Direction{North, East, South, West}

	for _, d := range directions {
		switch d {
		case North:
			fmt.Println(d, "→ 위로 이동 (y+1)")
		case East:
			fmt.Println(d, "→ 오른쪽 이동 (x+1)")
		case South:
			fmt.Println(d, "→ 아래로 이동 (y-1)")
		case West:
			fmt.Println(d, "→ 왼쪽 이동 (x-1)")
		}
	}
}
