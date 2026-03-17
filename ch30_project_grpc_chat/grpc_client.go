// 30장: gRPC 채팅 앱 - gRPC 클라이언트
// 이 파일은 gRPC를 이용한 채팅 클라이언트의 구현 예제이다.
//
// 실행 전 준비:
//   1. protoc으로 Go 코드 생성
//   2. grpc_server.go를 먼저 실행
//
// 실행: go run grpc_client.go
//
// 참고: 이 파일은 proto에서 생성된 코드가 필요하다.
//       생성된 코드 없이는 컴파일되지 않으며, 구조를 이해하기 위한 예제이다.

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// 생성된 protobuf 코드를 import한다
	// pb "your_module/proto/chatpb"
)

func main() {
	// === gRPC 서버에 연결 ===

	// 다이얼 옵션 설정 (개발용: TLS 없이 연결)
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("gRPC 서버 연결 실패:", err)
	}
	defer conn.Close()

	// 클라이언트 스텁 생성
	// 실제로는: client := pb.NewChatServiceClient(conn)
	_ = conn

	fmt.Println("=== gRPC 채팅 클라이언트 ===")
	fmt.Println("서버에 연결되었습니다: localhost:50051")
	fmt.Println()

	// 닉네임 입력
	fmt.Print("닉네임을 입력하세요: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())
	if username == "" {
		username = "익명"
	}

	fmt.Printf("\n%s님으로 채팅을 시작합니다.\n", username)
	fmt.Println("명령어: /users (사용자 목록), /quit (종료)")
	fmt.Println()

	// === 양방향 스트리밍 채팅 예시 (실제 gRPC 코드) ===
	//
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// // 양방향 스트리밍 연결
	// stream, err := client.Chat(ctx)
	// if err != nil {
	//     log.Fatal("스트림 생성 실패:", err)
	// }
	//
	// // 입장 메시지 전송
	// stream.Send(&pb.ChatMessage{
	//     User:      username,
	//     Content:   fmt.Sprintf("%s님이 입장했습니다.", username),
	//     Timestamp: time.Now().UnixMilli(),
	//     Type:      pb.MessageType_JOIN,
	// })
	//
	// // 수신 고루틴
	// go func() {
	//     for {
	//         msg, err := stream.Recv()
	//         if err != nil {
	//             fmt.Println("\n서버 연결이 끊어졌습니다.")
	//             cancel()
	//             return
	//         }
	//         timestamp := time.UnixMilli(msg.Timestamp).Format("15:04:05")
	//
	//         switch msg.Type {
	//         case pb.MessageType_NORMAL:
	//             fmt.Printf("\r[%s] %s: %s\n> ", timestamp, msg.User, msg.Content)
	//         case pb.MessageType_JOIN:
	//             fmt.Printf("\r[%s] %s\n> ", timestamp, msg.Content)
	//         case pb.MessageType_LEAVE:
	//             fmt.Printf("\r[%s] %s\n> ", timestamp, msg.Content)
	//         case pb.MessageType_SYSTEM:
	//             fmt.Printf("\r[시스템] %s\n> ", msg.Content)
	//         }
	//     }
	// }()
	//
	// // 송신 루프
	// for scanner.Scan() {
	//     text := strings.TrimSpace(scanner.Text())
	//     if text == "" {
	//         fmt.Print("> ")
	//         continue
	//     }
	//
	//     switch {
	//     case text == "/quit":
	//         stream.Send(&pb.ChatMessage{
	//             User:      username,
	//             Content:   fmt.Sprintf("%s님이 퇴장했습니다.", username),
	//             Timestamp: time.Now().UnixMilli(),
	//             Type:      pb.MessageType_LEAVE,
	//         })
	//         fmt.Println("채팅을 종료한다.")
	//         return
	//
	//     case text == "/users":
	//         // 단순 RPC로 사용자 목록 조회
	//         resp, err := client.GetUsers(context.Background(),
	//             &pb.UserListRequest{})
	//         if err != nil {
	//             fmt.Println("사용자 목록 조회 실패:", err)
	//         } else {
	//             fmt.Printf("접속 중인 사용자 (%d명):\n", resp.Count)
	//             for _, user := range resp.Users {
	//                 fmt.Printf("  - %s\n", user)
	//             }
	//         }
	//
	//     default:
	//         // 일반 메시지 전송
	//         stream.Send(&pb.ChatMessage{
	//             User:      username,
	//             Content:   text,
	//             Timestamp: time.Now().UnixMilli(),
	//             Type:      pb.MessageType_NORMAL,
	//         })
	//     }
	//     fmt.Print("> ")
	// }

	// === 데모용 시뮬레이션 (protobuf 코드 없이 동작 확인) ===
	fmt.Println()
	fmt.Println("[데모 모드] protobuf 생성 코드 없이 시뮬레이션합니다.")
	fmt.Println("[데모 모드] 실제 gRPC 연결은 위 주석의 코드를 활성화하세요.")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 데모: 사용자 입력을 받아 출력
	go func() {
		<-ctx.Done()
	}()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		if text == "/quit" {
			fmt.Println("채팅을 종료합니다.")
			break
		}

		if text == "/users" {
			fmt.Printf("[데모] 접속 중인 사용자: %s\n", username)
			continue
		}

		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("[%s] %s: %s\n", timestamp, username, text)
		fmt.Println("[데모] 실제 gRPC에서는 서버를 통해 모든 클라이언트에게 전달됩니다.")
	}
}
