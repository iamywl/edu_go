package main

import "fmt"

// === 사계절 열거값 ===
// iota는 const 블록마다 0으로 리셋됩니다.
const (
	Spring = iota // 0
	Summer        // 1
	Autumn        // 2
	Winter        // 3
)

// === 빈 식별자(_)로 0 건너뛰기 ===
const (
	_     = iota // 0은 건너뜀
	Red          // 1
	Green        // 2
	Blue         // 3
)

// === 비트 플래그 패턴 ===
// 파일 권한을 비트 연산으로 표현
const (
	ReadPerm  = 1 << iota // 1  (0001)
	WritePerm             // 2  (0010)
	ExecPerm              // 4  (0100)
)

// === 크기 단위 (KB, MB, GB) ===
const (
	_  = iota             // 0 건너뜀
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB                    // 1 << 20 = 1,048,576
	GB                    // 1 << 30 = 1,073,741,824
)

// === 커스텀 타입과 iota ===
type Weekday int

const (
	Sun Weekday = iota // 0
	Mon                // 1
	Tue                // 2
	Wed                // 3
	Thu                // 4
	Fri                // 5
	Sat                // 6
)

// String() 메서드로 출력 형태를 정의
func (d Weekday) String() string {
	names := [...]string{
		"일요일", "월요일", "화요일", "수요일",
		"목요일", "금요일", "토요일",
	}
	if d < Sun || d > Sat {
		return "알 수 없음"
	}
	return names[d]
}

func main() {
	// === 사계절 ===
	fmt.Println("=== 사계절 ===")
	fmt.Println("봄:", Spring)
	fmt.Println("여름:", Summer)
	fmt.Println("가을:", Autumn)
	fmt.Println("겨울:", Winter)

	// === 색상 (1부터) ===
	fmt.Println("\n=== 색상 ===")
	fmt.Println("Red:", Red)
	fmt.Println("Green:", Green)
	fmt.Println("Blue:", Blue)

	// === 비트 플래그 ===
	fmt.Println("\n=== 파일 권한 (비트 플래그) ===")
	fmt.Printf("읽기:  %d (%04b)\n", ReadPerm, ReadPerm)
	fmt.Printf("쓰기:  %d (%04b)\n", WritePerm, WritePerm)
	fmt.Printf("실행:  %d (%04b)\n", ExecPerm, ExecPerm)

	// OR 연산으로 권한 조합
	readExec := ReadPerm | ExecPerm
	fmt.Printf("읽기+실행: %d (%04b)\n", readExec, readExec)

	allPerm := ReadPerm | WritePerm | ExecPerm
	fmt.Printf("모든 권한: %d (%04b)\n", allPerm, allPerm)

	// AND 연산으로 권한 확인
	fmt.Println("\n=== 권한 확인 ===")
	myPerm := ReadPerm | ExecPerm // 읽기+실행 권한
	if myPerm&ReadPerm != 0 {
		fmt.Println("읽기 권한 있음")
	}
	if myPerm&WritePerm != 0 {
		fmt.Println("쓰기 권한 있음")
	} else {
		fmt.Println("쓰기 권한 없음")
	}

	// === 크기 단위 ===
	fmt.Println("\n=== 크기 단위 ===")
	fmt.Println("1 KB =", KB, "bytes")
	fmt.Println("1 MB =", MB, "bytes")
	fmt.Println("1 GB =", GB, "bytes")

	// === 커스텀 타입 + iota ===
	fmt.Println("\n=== 요일 (커스텀 타입) ===")
	today := Wed
	fmt.Println("오늘은", today) // String() 메서드 호출
	fmt.Println("내일은", today+1)
}
