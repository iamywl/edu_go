package main

import "fmt"

// =============================================
// 포인터와 복사에 대한 이해
// 포인터를 사용해도 포인터 값 자체는 복사됩니다
// =============================================

// LargeStruct는 큰 구조체 예제이다
type LargeStruct struct {
	Name string
	Data [100]int // 비교적 큰 데이터
}

// SmallStruct는 작은 구조체 예제이다
type SmallStruct struct {
	X, Y int
}

// --- 값으로 전달: 전체 구조체가 복사됨 ---
func modifyByValue(s LargeStruct) {
	s.Name = "변경됨(값)"
	s.Data[0] = 999
	fmt.Printf("  함수 내부 (값): Name=%s, Data[0]=%d\n", s.Name, s.Data[0])
}

// --- 포인터로 전달: 포인터(주소값)만 복사됨 ---
func modifyByPointer(s *LargeStruct) {
	s.Name = "변경됨(포인터)"
	s.Data[0] = 999
	fmt.Printf("  함수 내부 (포인터): Name=%s, Data[0]=%d\n", s.Name, s.Data[0])
}

func main() {
	fmt.Println("========== 포인터와 복사 ==========\n")

	// --- 1. 값 전달 vs 포인터 전달 ---
	fmt.Println("=== 1. 값 전달 vs 포인터 전달 ===")

	original := LargeStruct{Name: "원본", Data: [100]int{1, 2, 3}}

	fmt.Printf("  호출 전: Name=%s, Data[0]=%d\n", original.Name, original.Data[0])

	modifyByValue(original) // 복사본이 전달됨
	fmt.Printf("  값 전달 후: Name=%s, Data[0]=%d (변경 안 됨!)\n",
		original.Name, original.Data[0])

	modifyByPointer(&original) // 주소가 전달됨
	fmt.Printf("  포인터 전달 후: Name=%s, Data[0]=%d (변경됨!)\n",
		original.Name, original.Data[0])

	// --- 2. 포인터 값 자체의 복사 ---
	fmt.Println("\n=== 2. 포인터 값 자체의 복사 ===")

	data := &LargeStruct{Name: "공유 데이터"}
	copyOfPtr := data // 포인터 값(메모리 주소)이 복사됨

	fmt.Printf("  data 주소:      %p\n", data)
	fmt.Printf("  copyOfPtr 주소: %p\n", copyOfPtr)
	fmt.Printf("  같은 객체? %v\n", data == copyOfPtr) // true

	// copyOfPtr를 통해 수정하면 data도 영향 받음
	copyOfPtr.Name = "copyOfPtr가 변경함"
	fmt.Printf("  data.Name:      %s\n", data.Name) // "copyOfPtr가 변경함"
	fmt.Printf("  copyOfPtr.Name: %s\n", copyOfPtr.Name)

	// 하지만 포인터 변수 자체를 교체하면 원래 포인터에 영향 없음
	copyOfPtr = &LargeStruct{Name: "새로운 객체"}
	fmt.Printf("\n  포인터 교체 후:")
	fmt.Printf("\n  data.Name:      %s (변경 안 됨!)\n", data.Name)
	fmt.Printf("  copyOfPtr.Name: %s (새 객체)\n", copyOfPtr.Name)
	fmt.Printf("  같은 객체? %v\n", data == copyOfPtr) // false

	// --- 3. 슬라이스의 복사 동작 ---
	fmt.Println("\n=== 3. 슬라이스의 복사 동작 ===")
	fmt.Println("  슬라이스 = (포인터, 길이, 용량) 헤더")

	nums := []int{10, 20, 30, 40, 50}
	sub := nums[1:3] // [20, 30] — 같은 내부 배열 공유

	fmt.Printf("  nums: %v\n", nums)
	fmt.Printf("  sub:  %v (nums[1:3])\n", sub)

	sub[0] = 999 // 내부 배열이 공유되므로 nums도 변경됨!
	fmt.Printf("\n  sub[0] = 999 후:")
	fmt.Printf("\n  nums: %v (영향 받음!)\n", nums)
	fmt.Printf("  sub:  %v\n", sub)

	// 독립적인 복사본 만들기
	independent := make([]int, len(sub))
	copy(independent, sub)
	independent[0] = 111

	fmt.Printf("\n  copy 후 independent[0] = 111:")
	fmt.Printf("\n  nums:        %v (영향 없음)\n", nums)
	fmt.Printf("  independent: %v (독립적)\n", independent)

	// --- 4. 함수에 슬라이스를 전달할 때 ---
	fmt.Println("\n=== 4. 함수에 슬라이스 전달 ===")

	scores := []int{100, 200, 300}
	fmt.Printf("  전달 전: %v\n", scores)

	modifySlice(scores)
	fmt.Printf("  수정 후: %v (요소 변경은 반영됨)\n", scores)

	appendToSlice(scores)
	fmt.Printf("  append 후: %v (append는 반영 안 됨!)\n", scores)

	// append 결과를 반영하려면 반환값을 받아야 함
	scores = appendAndReturn(scores)
	fmt.Printf("  반환값 받은 후: %v (반영됨)\n", scores)

	// --- 5. 포인터를 사용해야 하는 경우 vs 값을 사용해야 하는 경우 ---
	fmt.Println("\n=== 5. 값 vs 포인터 선택 가이드 ===")
	fmt.Println("  값 타입 사용:")
	fmt.Println("    - 작은 구조체 (Point, Color 등)")
	fmt.Println("    - 불변 데이터")
	fmt.Println("    - 동시성에서 안전하게 사용할 때")
	fmt.Println()
	fmt.Println("  포인터 사용:")
	fmt.Println("    - 큰 구조체 (복사 비용 절감)")
	fmt.Println("    - 메서드에서 상태를 변경할 때")
	fmt.Println("    - nil이 의미 있는 값일 때")
	fmt.Println("    - 여러 곳에서 같은 객체를 공유할 때")
}

// 슬라이스의 요소를 변경 — 원본에 반영됨
func modifySlice(s []int) {
	if len(s) > 0 {
		s[0] = 1 // 내부 배열 공유 → 원본 변경됨
	}
}

// 슬라이스에 append — 원본에 반영 안 됨
func appendToSlice(s []int) {
	s = append(s, 400) // 새 슬라이스 헤더 생성 → 원본 영향 없음
	fmt.Printf("  함수 내부 append: %v\n", s)
}

// append 결과를 반환하여 반영
func appendAndReturn(s []int) []int {
	return append(s, 400)
}
