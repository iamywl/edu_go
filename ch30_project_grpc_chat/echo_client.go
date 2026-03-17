// 30장: TCP 채팅 앱 - Echo 클라이언트
// 실행: go run echo_client.go
// 먼저 echo_server.go를 실행해야 한다.
//
// 사용자가 입력한 메시지를 서버에 보내고, 에코 응답을 받아 출력한다.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// 서버에 TCP 연결
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal("서버 연결 실패:", err)
	}
	defer conn.Close()

	fmt.Println("=== TCP Echo 클라이언트 ===")
	fmt.Println("서버에 연결되었습니다: localhost:9000")
	fmt.Println("메시지를 입력하면 서버에서 에코 응답을 보냅니다.")
	fmt.Println("종료: Ctrl+C 또는 'quit' 입력")
	fmt.Println()

	// 수신 고루틴: 서버로부터 응답을 비동기로 받습니다
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("\n서버 연결이 끊어졌습니다.")
				os.Exit(0)
			}
			fmt.Printf("[에코] %s", string(buf[:n]))
		}
	}()

	// 송신: 사용자 입력을 서버로 보냅니다
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()

		// 종료 명령
		if text == "quit" || text == "exit" {
			fmt.Println("연결을 종료합니다.")
			break
		}

		// 서버로 메시지 전송 (개행 문자 포함)
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			log.Println("전송 에러:", err)
			break
		}
	}
}
