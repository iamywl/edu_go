// 30장: TCP 채팅 앱 - 채팅 서버
// 실행: go run chat_server.go
// 클라이언트 접속: go run chat_client.go
//
// 여러 클라이언트가 접속하여 메시지를 주고받을 수 있는 채팅 서버이다.
// 한 클라이언트가 보낸 메시지를 모든 클라이언트에게 브로드캐스트한다.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// ChatServer - 채팅 서버 구조체
type ChatServer struct {
	clients    map[net.Conn]string // 연결 -> 닉네임 매핑
	broadcast  chan string         // 브로드캐스트 메시지 채널
	register   chan net.Conn       // 새 클라이언트 등록 채널
	unregister chan net.Conn       // 클라이언트 해제 채널
	mu         sync.RWMutex        // 클라이언트 맵 보호용 뮤텍스
}

// NewChatServer - 새 채팅 서버를 생성한다
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:    make(map[net.Conn]string),
		broadcast:  make(chan string, 100), // 버퍼 있는 채널
		register:   make(chan net.Conn, 10),
		unregister: make(chan net.Conn, 10),
	}
}

// Run - 채팅 서버의 이벤트 루프를 실행한다
func (s *ChatServer) Run() {
	for {
		select {
		case conn := <-s.register:
			// 새 클라이언트 등록
			s.mu.Lock()
			s.clients[conn] = conn.RemoteAddr().String()
			s.mu.Unlock()

			count := len(s.clients)
			msg := fmt.Sprintf("[시스템] %s님이 입장했습니다. (현재 %d명)",
				conn.RemoteAddr().String(), count)
			fmt.Println(msg)
			s.broadcast <- msg

		case conn := <-s.unregister:
			// 클라이언트 해제
			s.mu.Lock()
			nickname := s.clients[conn]
			delete(s.clients, conn)
			s.mu.Unlock()
			conn.Close()

			count := len(s.clients)
			msg := fmt.Sprintf("[시스템] %s님이 퇴장했습니다. (현재 %d명)",
				nickname, count)
			fmt.Println(msg)
			s.broadcast <- msg

		case msg := <-s.broadcast:
			// 모든 클라이언트에게 메시지 전송
			s.mu.RLock()
			for conn := range s.clients {
				_, err := conn.Write([]byte(msg + "\n"))
				if err != nil {
					// 전송 실패한 클라이언트는 나중에 정리됩니다
					log.Printf("전송 실패 (%s): %v",
						conn.RemoteAddr().String(), err)
				}
			}
			s.mu.RUnlock()
		}
	}
}

// HandleClient - 개별 클라이언트의 메시지를 처리한다
func (s *ChatServer) HandleClient(conn net.Conn) {
	// 클라이언트 등록
	s.register <- conn

	// 닉네임 설정 안내
	conn.Write([]byte("[시스템] 채팅 서버에 오신 것을 환영합니다!\n"))
	conn.Write([]byte("[시스템] 닉네임을 입력하세요: "))

	scanner := bufio.NewScanner(conn)

	// 닉네임 수신
	var nickname string
	if scanner.Scan() {
		nickname = strings.TrimSpace(scanner.Text())
		if nickname == "" {
			nickname = conn.RemoteAddr().String()
		}
	}

	// 닉네임 등록
	s.mu.Lock()
	s.clients[conn] = nickname
	s.mu.Unlock()

	joinMsg := fmt.Sprintf("[시스템] '%s'님이 닉네임을 설정했습니다.", nickname)
	fmt.Println(joinMsg)
	s.broadcast <- joinMsg

	// 메시지 수신 루프
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		// 명령어 처리
		if strings.HasPrefix(text, "/") {
			s.handleCommand(conn, nickname, text)
			continue
		}

		// 일반 메시지 브로드캐스트
		timestamp := time.Now().Format("15:04:05")
		msg := fmt.Sprintf("[%s] %s: %s", timestamp, nickname, text)
		fmt.Println(msg)
		s.broadcast <- msg
	}

	// 클라이언트 해제
	s.unregister <- conn
}

// handleCommand - 채팅 명령어를 처리한다
func (s *ChatServer) handleCommand(conn net.Conn, nickname, cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "/help":
		help := `[시스템] 사용 가능한 명령어:
  /help   - 도움말 표시
  /users  - 접속 중인 사용자 목록
  /quit   - 채팅 종료
`
		conn.Write([]byte(help))

	case "/users":
		s.mu.RLock()
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("[시스템] 접속 중인 사용자 (%d명):\n", len(s.clients)))
		for _, name := range s.clients {
			sb.WriteString(fmt.Sprintf("  - %s\n", name))
		}
		s.mu.RUnlock()
		conn.Write([]byte(sb.String()))

	case "/quit":
		conn.Write([]byte("[시스템] 연결을 종료합니다.\n"))
		s.unregister <- conn

	default:
		conn.Write([]byte(fmt.Sprintf("[시스템] 알 수 없는 명령어: %s (/help로 도움말 확인)\n", parts[0])))
	}
}

func main() {
	// TCP 리스너 생성
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("리스너 생성 실패:", err)
	}
	defer listener.Close()

	fmt.Println("=== TCP 채팅 서버 ===")
	fmt.Println("주소: localhost:9001")
	fmt.Println("클라이언트 접속: go run chat_client.go")
	fmt.Println("또는: nc localhost 9001")
	fmt.Println("종료: Ctrl+C")
	fmt.Println()

	// 채팅 서버 생성 및 이벤트 루프 시작
	server := NewChatServer()
	go server.Run()

	// 클라이언트 연결 대기 루프
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("연결 수락 에러:", err)
			continue
		}
		go server.HandleClient(conn)
	}
}
