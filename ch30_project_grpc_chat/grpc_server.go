// 30장: gRPC 채팅 앱 - gRPC 서버
// 이 파일은 gRPC를 이용한 채팅 서버의 구현 예제이다.
//
// 실행 전 준비:
//   1. protoc으로 Go 코드 생성:
//      protoc --go_out=. --go-grpc_out=. proto/chat.proto
//   2. 필요한 패키지 설치:
//      go get google.golang.org/grpc
//      go get google.golang.org/protobuf
//
// 실행: go run grpc_server.go
// 클라이언트: go run grpc_client.go
//
// 참고: 이 파일은 proto에서 생성된 코드가 필요하다.
//       생성된 코드 없이는 컴파일되지 않으며, 구조를 이해하기 위한 예제이다.

package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	// 생성된 protobuf 코드를 import한다
	// pb "your_module/proto/chatpb"
)

// === gRPC 채팅 서버 구현 ===

// chatServiceServer - ChatService 서버 구현체
// pb.UnimplementedChatServiceServer를 임베딩하여 forward compatibility를 보장한다
type chatServiceServer struct {
	// pb.UnimplementedChatServiceServer // 자동 생성된 코드의 기본 구현

	mu          sync.RWMutex
	clients     map[string]chan *ChatMsg // 사용자별 메시지 채널
	chatStreams map[string]interface{}   // 양방향 스트리밍 연결
}

// ChatMsg - 내부용 채팅 메시지 (protobuf 대신 사용)
type ChatMsg struct {
	User      string
	Content   string
	Timestamp int64
	MsgType   int // 0: NORMAL, 1: JOIN, 2: LEAVE, 3: SYSTEM
}

// newChatServiceServer - 새 gRPC 채팅 서버 생성
func newChatServiceServer() *chatServiceServer {
	return &chatServiceServer{
		clients:     make(map[string]chan *ChatMsg),
		chatStreams: make(map[string]interface{}),
	}
}

// Join - 채팅방 입장 (단순 RPC)
// 클라이언트가 채팅방에 입장할 때 호출됩니다
func (s *chatServiceServer) Join(user, room string) (bool, string, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 이미 접속한 사용자인지 확인
	if _, exists := s.clients[user]; exists {
		return false, "이미 접속 중인 사용자입니다", len(s.clients)
	}

	// 메시지 채널 생성 (버퍼: 100개)
	s.clients[user] = make(chan *ChatMsg, 100)

	// 입장 알림 브로드캐스트
	joinMsg := &ChatMsg{
		User:      "시스템",
		Content:   fmt.Sprintf("%s님이 입장했습니다.", user),
		Timestamp: time.Now().UnixMilli(),
		MsgType:   1, // JOIN
	}
	s.broadcastMessage(joinMsg)

	return true, fmt.Sprintf("%s 방에 입장했습니다", room), len(s.clients)
}

// Leave - 채팅방 퇴장 (단순 RPC)
func (s *chatServiceServer) Leave(user string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, exists := s.clients[user]; exists {
		close(ch)
		delete(s.clients, user)

		// 퇴장 알림 브로드캐스트
		leaveMsg := &ChatMsg{
			User:      "시스템",
			Content:   fmt.Sprintf("%s님이 퇴장했습니다.", user),
			Timestamp: time.Now().UnixMilli(),
			MsgType:   2, // LEAVE
		}
		s.broadcastMessage(leaveMsg)
	}
}

// GetUsers - 접속 중인 사용자 목록 조회 (단순 RPC)
func (s *chatServiceServer) GetUsers() ([]string, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]string, 0, len(s.clients))
	for user := range s.clients {
		users = append(users, user)
	}
	return users, len(users)
}

// SendMessage - 메시지 전송 (단순 RPC)
func (s *chatServiceServer) SendMessage(msg *ChatMsg) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.broadcastMessage(msg)
}

// broadcastMessage - 모든 클라이언트에게 메시지 전송 (내부 함수)
// 주의: 호출 시 mu가 이미 잠겨있어야 한다
func (s *chatServiceServer) broadcastMessage(msg *ChatMsg) {
	for _, ch := range s.clients {
		// 비차단 전송 (채널이 가득 차면 건너뜀)
		select {
		case ch <- msg:
		default:
			// 채널이 가득 찼으면 메시지 드롭
			log.Println("메시지 드롭: 클라이언트 채널이 가득 찼습니다")
		}
	}
}

// Subscribe - 메시지 구독 (서버 스트리밍 RPC)
// 실제 gRPC에서는 stream으로 메시지를 전달한다:
//
//	func (s *chatServiceServer) Subscribe(req *pb.JoinRequest,
//	    stream pb.ChatService_SubscribeServer) error {
//	    ch := s.clients[req.User]
//	    for msg := range ch {
//	        if err := stream.Send(msg); err != nil {
//	            return err
//	        }
//	    }
//	    return nil
//	}

// Chat - 양방향 스트리밍 (양방향 스트리밍 RPC)
// 실제 gRPC에서는 다음과 같이 구현한다:
//
//	func (s *chatServiceServer) Chat(
//	    stream pb.ChatService_ChatServer) error {
//
//	    // 첫 메시지에서 사용자 이름을 추출한다
//	    firstMsg, err := stream.Recv()
//	    if err != nil {
//	        return err
//	    }
//	    user := firstMsg.User
//
//	    // 사용자 등록
//	    s.mu.Lock()
//	    s.clients[user] = make(chan *ChatMsg, 100)
//	    s.mu.Unlock()
//
//	    // 수신 고루틴: 클라이언트 -> 서버
//	    go func() {
//	        for {
//	            msg, err := stream.Recv()
//	            if err != nil {
//	                s.Leave(user)
//	                return
//	            }
//	            s.SendMessage(msg)
//	        }
//	    }()
//
//	    // 송신: 서버 -> 클라이언트 (채널에서 메시지를 읽어 전송)
//	    ch := s.clients[user]
//	    for msg := range ch {
//	        if err := stream.Send(msg); err != nil {
//	            return err
//	        }
//	    }
//	    return nil
//	}

func main() {
	// TCP 리스너 생성 (50051은 gRPC의 기본 포트)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("리스너 생성 실패:", err)
	}

	// gRPC 서버 생성
	grpcServer := grpc.NewServer()

	// 채팅 서비스 등록
	server := newChatServiceServer()
	_ = server // 실제로는 pb.RegisterChatServiceServer(grpcServer, server)로 등록

	// 서버에 서비스를 등록한다 (protobuf 생성 코드 필요):
	// pb.RegisterChatServiceServer(grpcServer, server)

	fmt.Println("=== gRPC 채팅 서버 ===")
	fmt.Println("주소: localhost:50051")
	fmt.Println()
	fmt.Println("이 서버를 실행하려면 다음 단계가 필요합니다:")
	fmt.Println("  1. protoc --go_out=. --go-grpc_out=. proto/chat.proto")
	fmt.Println("  2. go mod init 및 go mod tidy")
	fmt.Println("  3. import 경로를 실제 생성된 패키지로 수정")
	fmt.Println()
	fmt.Println("gRPC 서버를 시작합니다...")

	// gRPC 서버 시작
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("gRPC 서버 시작 실패:", err)
	}
}
