// formatting.go
// fmt.Printf의 다양한 서식 지정자(format verbs) 예제
// 실행 방법: go run formatting.go

package main

import "fmt"

// Person 구조체 (서식 예제용)
type Person struct {
	Name string
	Age  int
}

func main() {
	fmt.Println("===== 정수 서식 =====")
	num := 42

	fmt.Printf("%%d  10진수:       %d\n", num)
	fmt.Printf("%%b  2진수:        %b\n", num)
	fmt.Printf("%%o  8진수:        %o\n", num)
	fmt.Printf("%%O  8진수(0o):    %O\n", num)
	fmt.Printf("%%x  16진수(소문자): %x\n", num)
	fmt.Printf("%%X  16진수(대문자): %X\n", num)
	fmt.Printf("%%c  유니코드 문자:  %c\n", num) // * (ASCII 42)
	fmt.Printf("%%U  유니코드 표기:  %U\n", num) // U+002A

	fmt.Println()
	fmt.Println("===== 실수 서식 =====")
	pi := 3.14159265358979

	fmt.Printf("%%f   기본:        %f\n", pi)
	fmt.Printf("%%.2f 소수점 2자리: %.2f\n", pi)
	fmt.Printf("%%.5f 소수점 5자리: %.5f\n", pi)
	fmt.Printf("%%e   지수(소문자): %e\n", pi)
	fmt.Printf("%%E   지수(대문자): %E\n", pi)
	fmt.Printf("%%g   짧은 표기:    %g\n", pi)

	fmt.Println()
	fmt.Println("===== 문자열 서식 =====")
	str := "Hello, Go!"

	fmt.Printf("%%s  문자열:        %s\n", str)
	fmt.Printf("%%q  따옴표 포함:   %q\n", str)
	fmt.Printf("%%x  16진수 바이트: %x\n", str)

	fmt.Println()
	fmt.Println("===== 불리언 서식 =====")
	fmt.Printf("%%t  true:  %t\n", true)
	fmt.Printf("%%t  false: %t\n", false)

	fmt.Println()
	fmt.Println("===== 일반 서식 (%v) =====")
	p := Person{Name: "홍길동", Age: 25}

	fmt.Printf("%%v   기본:        %v\n", p)
	fmt.Printf("%%+v  필드 이름:   %+v\n", p)
	fmt.Printf("%%#v  Go 문법:     %#v\n", p)
	fmt.Printf("%%T   타입:        %T\n", p)

	fmt.Println()
	fmt.Println("===== 너비(Width) 지정 =====")
	// 최소 너비를 지정하여 정렬할 수 있습니다

	// 정수 너비
	fmt.Printf("[%%10d]  오른쪽 정렬: [%10d]\n", 42)
	fmt.Printf("[%%-10d] 왼쪽 정렬:  [%-10d]\n", 42)
	fmt.Printf("[%%010d] 0으로 채움: [%010d]\n", 42)
	fmt.Printf("[%%+d]   부호 표시:  [%+d]\n", 42)
	fmt.Printf("[%%+d]   부호 표시:  [%+d]\n", -42)

	fmt.Println()

	// 문자열 너비
	fmt.Printf("[%%20s]  오른쪽 정렬: [%20s]\n", "Hello")
	fmt.Printf("[%%-20s] 왼쪽 정렬:  [%-20s]\n", "Hello")
	fmt.Printf("[%%.3s]  잘라내기:   [%.3s]\n", "Hello") // "Hel"

	fmt.Println()

	// 실수 너비와 정밀도
	fmt.Printf("[%%10.2f]  너비 10, 소수점 2: [%10.2f]\n", 3.14159)
	fmt.Printf("[%%-10.2f] 왼쪽 정렬:        [%-10.2f]\n", 3.14159)

	fmt.Println()
	fmt.Println("===== 표 형태 출력 =====")
	// 너비 지정을 활용한 깔끔한 표 출력

	fmt.Printf("%-10s %-5s %10s\n", "이름", "나이", "도시")
	fmt.Printf("%-10s %-5s %10s\n", "----------", "-----", "----------")
	fmt.Printf("%-10s %-5d %10s\n", "홍길동", 25, "서울")
	fmt.Printf("%-10s %-5d %10s\n", "김철수", 30, "부산")
	fmt.Printf("%-10s %-5d %10s\n", "이영희", 28, "대구")

	fmt.Println()
	fmt.Println("===== %% (퍼센트 기호 출력) =====")
	fmt.Printf("진행률: %d%%\n", 85) // 85%
}
