// check_install.go
// Go 설치 확인 프로그램
// 실행 방법: go run check_install.go

package main

import (
	"fmt"
	"runtime"
)

func main() {
	// Go 설치 확인 메시지 출력
	fmt.Println("========================================")
	fmt.Println("  Go 설치 확인 프로그램")
	fmt.Println("========================================")

	// runtime 패키지를 사용하여 Go 환경 정보 출력
	fmt.Printf("Go 버전    : %s\n", runtime.Version())
	fmt.Printf("운영체제   : %s\n", runtime.GOOS)
	fmt.Printf("아키텍처   : %s\n", runtime.GOARCH)
	fmt.Printf("CPU 코어 수: %d\n", runtime.NumCPU())
	fmt.Printf("GOROOT     : %s\n", runtime.GOROOT())

	fmt.Println("========================================")
	fmt.Println("  Go 설치가 정상적으로 완료되었습니다!")
	fmt.Println("========================================")
}
