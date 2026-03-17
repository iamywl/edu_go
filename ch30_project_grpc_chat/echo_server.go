// 30장: TCP 채팅 앱 - Echo 서버
// 실행: go run echo_server.go
// 테스트: nc localhost 9000 또는 echo_client.go 실행
//
// 클라이언트가 보낸 메시지를 그대로 돌려보내는 Echo 서버이다.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// handleEchoConnection - 클라이언트 연결을 처리한다
func handleEchoConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("[접속] 클라이언트 연결됨: %s\n", clientAddr)

	// 읽기/쓰기 버퍼
	buf := make([]byte, 4096)

	for {
		// 읽기 타임아웃 설정 (60초)
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// 클라이언트로부터 데이터 읽기
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("[해제] 클라이언트 정상 종료: %s\n", clientAddr)
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("[타임아웃] 클라이언트 타임아웃: %s\n", clientAddr)
			} else {
				fmt.Printf("[에러] 읽기 에러 (%s): %v\n", clientAddr, err)
			}
			return
		}

		// 수신한 메시지 출력
		message := string(buf[:n])
		fmt.Printf("[수신] %s: %s", clientAddr, message)

		// Echo: 받은 데이터를 그대로 클라이언트에게 돌려보냄
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Printf("[에러] 쓰기 에러 (%s): %v\n", clientAddr, err)
			return
		}
	}
}

func main() {
	// TCP 리스너 생성 (9000번 포트)
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("리스너 생성 실패:", err)
	}
	defer listener.Close()

	fmt.Println("=== TCP Echo 서버 ===")
	fmt.Println("주소: localhost:9000")
	fmt.Println("테스트: nc localhost 9000")
	fmt.Println("종료: Ctrl+C")
	fmt.Println()

	// 클라이언트 연결 대기 루프
	for {
		// 새 클라이언트 연결을 기다립니다
		conn, err := listener.Accept()
		if err != nil {
			log.Println("연결 수락 에러:", err)
			continue
		}

		// 각 클라이언트를 별도의 고루틴으로 처리한다
		go handleEchoConnection(conn)
	}
}
