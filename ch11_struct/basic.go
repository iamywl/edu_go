package main

import "fmt"

// Student 구조체 정의
type Student struct {
	Name  string
	Age   int
	Grade string
}

func main() {
	// === 구조체 선언 및 기본 사용 ===
	fmt.Println("=== 구조체 기본 사용 ===")

	// 방법 1: 제로값으로 생성 후 필드 설정
	var s1 Student
	s1.Name = "홍길동"
	s1.Age = 20
	s1.Grade = "A"
	fmt.Printf("s1: %+v\n", s1)

	// 방법 2: 필드명을 지정하여 초기화 (권장)
	s2 := Student{
		Name:  "김철수",
		Age:   22,
		Grade: "B",
	}
	fmt.Printf("s2: %+v\n", s2)

	// 방법 3: 순서대로 초기화 (비권장)
	s3 := Student{"이영희", 21, "A"}
	fmt.Printf("s3: %+v\n", s3)

	// 방법 4: 일부 필드만 초기화 (나머지는 제로값)
	s4 := Student{
		Name: "박민수",
		// Age와 Grade는 제로값 (0, "")
	}
	fmt.Printf("s4: %+v\n", s4)

	// === 구조체 필드 접근 ===
	fmt.Println("\n=== 필드 접근 ===")
	fmt.Printf("%s 학생은 %d세이며, 성적은 %s입니다.\n", s2.Name, s2.Age, s2.Grade)

	// === 구조체 비교 ===
	fmt.Println("\n=== 구조체 비교 ===")
	s5 := Student{Name: "홍길동", Age: 20, Grade: "A"}
	fmt.Printf("s1 == s5: %v\n", s1 == s5) // 모든 필드가 같으면 true

	// === 함수에 구조체 전달 (값 복사) ===
	fmt.Println("\n=== 함수에 구조체 전달 ===")
	printStudent(s2)

	// 값 복사이므로 원본에 영향 없음
	changeGrade(s2, "A+")
	fmt.Printf("변경 시도 후 s2.Grade: %s (변경 안 됨)\n", s2.Grade)
}

// printStudent 는 학생 정보를 출력하는 함수
func printStudent(s Student) {
	fmt.Printf("[정보] 이름: %s, 나이: %d, 성적: %s\n", s.Name, s.Age, s.Grade)
}

// changeGrade 는 성적을 변경하려 하지만, 값 복사이므로 원본에는 영향이 없다
func changeGrade(s Student, newGrade string) {
	s.Grade = newGrade
	fmt.Printf("[함수 내부] %s의 성적을 %s로 변경함\n", s.Name, newGrade)
}
