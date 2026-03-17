// main.go
// Go 언어의 첫 번째 프로그램: Hello Go World
// 실행 방법: go run main.go

package main // 실행 가능한 프로그램은 반드시 main 패키지여야 한다

import "fmt" // fmt 패키지: 텍스트 입출력을 위한 표준 라이브러리

// main 함수: 프로그램의 시작점(entry point)
func main() {
	// === 기본 Hello World ===
	fmt.Println("Hello, Go World!")

	// === 코드 뜯어보기 ===

	// 1. fmt.Println - 값을 출력하고 줄바꿈
	fmt.Println("1. Println은 출력 후 줄바꿈을 합니다")
	fmt.Println("2. 이렇게 각 줄이 분리됩니다")

	// 2. fmt.Print - 줄바꿈 없이 출력
	fmt.Print("3. Print는 ")
	fmt.Print("줄바꿈 없이 ")
	fmt.Print("이어서 출력합니다\n") // \n으로 직접 줄바꿈

	// 3. fmt.Printf - 서식(format)을 지정하여 출력
	name := "Go" // 짧은 변수 선언 (ch02에서 자세히 배웁니다)
	year := 2009
	fmt.Printf("4. %s 언어는 %d년에 탄생했습니다\n", name, year)

	// === 여러 값을 한 번에 출력 ===
	fmt.Println("5. 여러 값:", "Go", "언어", 2009, true)

	// === 이스케이프 시퀀스 ===
	fmt.Println("6. 줄바꿈: 첫째 줄\n   둘째 줄")
	fmt.Println("7. 탭: 이름\t나이\t도시")
	fmt.Println("8. 큰따옴표 출력: \"Go는 재미있다\"")
}
