// scan_examples.go
// fmt 패키지의 표준 입력 함수 사용 예제
// 실행 방법: go run scan_examples.go

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("===== fmt.Scan 예제 =====")
	fmt.Println("공백이나 줄바꿈으로 구분하여 여러 값을 입력받습니다.")

	var name string
	var age int

	fmt.Print("이름과 나이를 입력하세요 (예: 홍길동 25): ")
	// Scan은 공백/줄바꿈으로 구분하여 값을 읽어옴
	// &는 변수의 메모리 주소를 전달 (포인터)
	n, err := fmt.Scan(&name, &age)
	if err != nil {
		fmt.Println("입력 에러:", err)
		return
	}
	fmt.Printf("입력받은 값 %d개: 이름=%s, 나이=%d\n", n, name, age)

	// 버퍼 비우기 (이전 입력의 나머지 줄바꿈 문자 제거)
	bufio.NewReader(os.Stdin).ReadString('\n')

	fmt.Println()
	fmt.Println("===== fmt.Scanln 예제 =====")
	fmt.Println("한 줄에서 공백으로 구분하여 입력받습니다.")
	fmt.Println("줄바꿈(Enter)이 나오면 입력을 중단합니다.")

	var city string
	var population int

	fmt.Print("도시명과 인구를 입력하세요 (예: 서울 9700000): ")
	// Scanln은 줄바꿈을 만나면 입력 중단
	fmt.Scanln(&city, &population)
	fmt.Printf("도시: %s, 인구: %d명\n", city, population)

	fmt.Println()
	fmt.Println("===== bufio.Scanner 예제 (권장) =====")
	fmt.Println("실무에서는 bufio.Scanner를 더 많이 사용합니다.")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("좋아하는 문장을 입력하세요: ")
	scanner.Scan()             // 한 줄을 읽어옴
	sentence := scanner.Text() // 읽은 텍스트 반환
	fmt.Printf("입력한 문장: %q\n", sentence)

	fmt.Println()
	fmt.Println("===== 여러 줄 입력 예제 =====")
	fmt.Println("'quit'을 입력하면 종료합니다.")

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "quit" {
			fmt.Println("프로그램을 종료합니다.")
			break
		}

		fmt.Printf("  입력값: %q (길이: %d)\n", input, len(input))
	}
}
