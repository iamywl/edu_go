// arithmetic.go
// 산술 연산자 예제
// 실행 방법: go run arithmetic.go

package main

import "fmt"

func main() {
	fmt.Println("===== 기본 산술 연산자 =====")

	a, b := 10, 3

	fmt.Printf("%d + %d = %d\n", a, b, a+b)  // 덧셈: 13
	fmt.Printf("%d - %d = %d\n", a, b, a-b)  // 뺄셈: 7
	fmt.Printf("%d * %d = %d\n", a, b, a*b)  // 곱셈: 30
	fmt.Printf("%d / %d = %d\n", a, b, a/b)  // 나눗셈: 3 (정수 나눗셈!)
	fmt.Printf("%d %% %d = %d\n", a, b, a%b) // 나머지: 1

	fmt.Println()
	fmt.Println("===== 정수 나눗셈 vs 실수 나눗셈 =====")

	// 정수 나눗셈: 소수점 이하가 버려짐 (주의!)
	fmt.Printf("정수: 10 / 3 = %d\n", 10/3) // 3
	fmt.Printf("정수: 7 / 2 = %d\n", 7/2)   // 3
	fmt.Printf("정수: -7 / 2 = %d\n", -7/2) // -3 (0을 향해 자름)

	// 실수 나눗셈: 소수점까지 계산
	fmt.Printf("실수: 10.0 / 3.0 = %f\n", 10.0/3.0) // 3.333333
	fmt.Printf("실수: 7.0 / 2.0 = %f\n", 7.0/2.0)   // 3.500000

	// 정수를 실수로 변환하여 나눗셈
	x, y := 10, 3
	result := float64(x) / float64(y)
	fmt.Printf("변환: float64(%d) / float64(%d) = %f\n", x, y, result)

	fmt.Println()
	fmt.Println("===== 나머지 연산자 (%) =====")

	// 나머지 연산의 활용
	fmt.Println("짝수/홀수 판별:")
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Printf("  %d: 짝수\n", i)
		} else {
			fmt.Printf("  %d: 홀수\n", i)
		}
	}

	fmt.Println()
	fmt.Println("===== 증감 연산자 (++, --) =====")

	count := 0
	fmt.Printf("초기값: %d\n", count)

	count++                            // count = count + 1
	fmt.Printf("count++: %d\n", count) // 1

	count++                            // count = count + 1
	fmt.Printf("count++: %d\n", count) // 2

	count--                            // count = count - 1
	fmt.Printf("count--: %d\n", count) // 1

	// Go에서 주의할 점:
	// 1. 전위 증감(++count)은 없음 → 컴파일 에러
	// 2. 식(expression)으로 사용 불가 → j := count++ 불가
	// 3. 오직 문장(statement)으로만 사용 가능

	fmt.Println()
	fmt.Println("===== 문자열 연결 (+) =====")

	firstName := "홍"
	lastName := "길동"
	fullName := firstName + lastName
	fmt.Printf("이름: %s\n", fullName) // 홍길동

	// += 으로 문자열 이어 붙이기
	greeting := "Hello"
	greeting += ", "
	greeting += "Go!"
	fmt.Println(greeting) // Hello, Go!

	fmt.Println()
	fmt.Println("===== 대입 연산자 =====")

	v := 100
	fmt.Printf("초기값:   %d\n", v)

	v += 10                         // v = v + 10
	fmt.Printf("v += 10:  %d\n", v) // 110

	v -= 20                         // v = v - 20
	fmt.Printf("v -= 20:  %d\n", v) // 90

	v *= 2                          // v = v * 2
	fmt.Printf("v *= 2:   %d\n", v) // 180

	v /= 3                          // v = v / 3
	fmt.Printf("v /= 3:   %d\n", v) // 60

	v %= 7                           // v = v % 7
	fmt.Printf("v %%= 7:   %d\n", v) // 4
}
