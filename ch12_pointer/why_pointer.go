package main

import "fmt"

// Student 구조체
type Student struct {
	Name  string
	Age   int
	Grade string
}

func main() {
	// === 값 복사 vs 포인터 전달 비교 ===

	s := Student{Name: "홍길동", Age: 20, Grade: "B"}

	// --- 값 복사로 전달 ---
	fmt.Println("=== 값 복사로 전달 ===")
	fmt.Printf("호출 전: %+v\n", s)
	tryPromoteByValue(s) // 복사본이 전달됨
	fmt.Printf("호출 후: %+v (변경 안 됨!)\n", s)

	fmt.Println()

	// --- 포인터로 전달 ---
	fmt.Println("=== 포인터로 전달 ===")
	fmt.Printf("호출 전: %+v\n", s)
	promoteByPointer(&s) // 주소가 전달됨
	fmt.Printf("호출 후: %+v (변경됨!)\n", s)

	fmt.Println()

	// === 슬라이스의 학생들을 포인터로 수정 ===
	fmt.Println("=== 여러 학생 성적 일괄 변경 ===")
	students := []Student{
		{Name: "김철수", Age: 21, Grade: "C"},
		{Name: "이영희", Age: 22, Grade: "C"},
		{Name: "박민수", Age: 20, Grade: "C"},
	}

	fmt.Println("변경 전:")
	for _, st := range students {
		fmt.Printf("  %s: %s\n", st.Name, st.Grade)
	}

	// 인덱스를 사용하여 원본에 접근
	for i := range students {
		promoteByPointer(&students[i])
	}

	fmt.Println("변경 후:")
	for _, st := range students {
		fmt.Printf("  %s: %s\n", st.Name, st.Grade)
	}
}

// tryPromoteByValue 는 값 복사로 전달받아 원본을 변경할 수 없다
func tryPromoteByValue(s Student) {
	s.Grade = "A+" // 복사본만 변경됨
	fmt.Printf("  [함수 내부] %s의 성적: %s\n", s.Name, s.Grade)
}

// promoteByPointer 는 포인터로 전달받아 원본을 직접 수정한다
func promoteByPointer(s *Student) {
	s.Grade = "A+" // 원본이 변경됨!
	// s.Grade는 (*s).Grade의 축약형 — Go가 자동 역참조해줌
	fmt.Printf("  [함수 내부] %s의 성적: %s\n", s.Name, s.Grade)
}
