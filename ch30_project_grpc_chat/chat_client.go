// 30장: TCP 채팅 앱 - 채팅 클라이언트
// 실행: go run chat_client.go
// 먼저 chat_server.go를 실행해야 한다.
//
// 채팅 서버에 접속하여 다른 사용자들과 메시지를 주고받습니다.
// 명령어: /help, /users, /quit
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// 채팅 서버에 TCP 연결
	conn, err := net.Dial("tcp", "localhost:9001")
	if err != nil {
		log.Fatal("서버 연결 실패:", err)
	}
	defer conn.Close()

	fmt.Println("=== TCP 채팅 클라이언트 ===")
	fmt.Println("서버에 연결되었습니다: localhost:9001")
	fmt.Println("명령어: /help, /users, /quit")
	fmt.Println()

	// 수신 고루틴: 서버로부터 메시지를 비동기로 받습니다
	done := make(chan struct{})
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			msg := scanner.Text()
			// 커서를 줄 앞으로 이동하고 메시지 출력 후 프롬프트 다시 표시
			fmt.Printf("\r%s\n> ", msg)
		}
		fmt.Println("\n서버 연결이 끊어졌습니다.")
	}()

	// 송신: 사용자 입력을 서버로 보냅니다
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		// 메시지 전송
		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			log.Println("전송 에러:", err)
			break
		}

		// /quit 명령이면 종료
		if text == "/quit" {
			fmt.Println("채팅을 종료합니다.")
			break
		}
	}

	// 수신 고루틴이 끝날 때까지 대기
	<-done
}
