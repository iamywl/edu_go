# 30장 [Project] TCP와 gRPC로 채팅 앱 만들기

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회, 포트 노출 필요: -p 9000:9000 -p 50051:50051)
make shell

# === TCP Echo 서버/클라이언트 ===
# 터미널 1: 서버 실행
go run ch30_project_grpc_chat/echo_server.go

# 터미널 2: 클라이언트 실행
go run ch30_project_grpc_chat/echo_client.go

# === TCP 채팅 서버/클라이언트 ===
# 터미널 1: 채팅 서버 실행
go run ch30_project_grpc_chat/chat_server.go

# 터미널 2, 3: 채팅 클라이언트 실행 (여러 개 접속 가능)
go run ch30_project_grpc_chat/chat_client.go

# === gRPC 채팅 서버/클라이언트 ===
# 터미널 1: gRPC 서버 실행
go run ch30_project_grpc_chat/grpc_server.go

# 터미널 2, 3: gRPC 클라이언트 실행
go run ch30_project_grpc_chat/grpc_client.go
```

> **참고**: TCP 채팅과 gRPC 채팅 모두 서버와 클라이언트를 **별도의 터미널**에서 실행해야 한다. Docker 컨테이너에 여러 터미널로 접속하려면 `docker exec -it <container> /bin/sh`를 사용한다.

> **Makefile 활용**: `make run CH=ch30_project_grpc_chat` 또는 `make run CH=ch30_project_grpc_chat FILE=echo_server.go`

---

이 장에서는 네트워크 프로그래밍의 기초인 TCP 소켓부터 시작하여, gRPC를 이용한 채팅 프로그램까지 단계적으로 만들어 본다. TCP 소켓은 네트워크 통신의 가장 기본적인 단위이며, 이를 이해하면 HTTP, gRPC 등 상위 프로토콜의 동작 원리를 명확히 파악할 수 있다. gRPC는 마이크로서비스 아키텍처에서 널리 사용되는 고성능 통신 프레임워크이다.

---

## 30.1 TCP를 이용한 Echo 서버 제작

### net 패키지로 TCP 서버 만들기

Go의 `net` 패키지는 TCP, UDP 등 저수준 네트워크 프로그래밍을 지원한다. TCP(Transmission Control Protocol)는 신뢰성 있는 연결 지향 프로토콜로, 데이터의 순서와 무결성을 보장한다. `net.Listen()`은 지정된 주소에서 연결을 대기하는 리스너를 생성하고, `listener.Accept()`는 새로운 클라이언트 연결이 들어올 때까지 블로킹된다:

```go
// TCP 리스너 생성
listener, err := net.Listen("tcp", ":9000")
if err != nil {
    log.Fatal(err)
}
defer listener.Close()

// 클라이언트 연결 대기
for {
    conn, err := listener.Accept()
    if err != nil {
        log.Println("연결 에러:", err)
        continue
    }
    // 각 연결을 고루틴으로 처리
    go handleConnection(conn)
}
```

각 클라이언트 연결을 별도의 고루틴에서 처리하므로 동시에 여러 클라이언트를 수용할 수 있다. Go의 고루틴은 경량이므로 수천 개의 동시 연결도 효율적으로 처리한다. `net.Conn` 인터페이스는 `Read()`, `Write()`, `Close()` 메서드를 제공하며 `io.Reader`와 `io.Writer`를 모두 구현한다.

### Echo 서버

Echo 서버는 클라이언트가 보낸 메시지를 그대로 돌려보내는 가장 단순한 형태의 서버이다. 네트워크 프로그래밍의 "Hello, World"에 해당한다:

```go
func handleConnection(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            return
        }
        conn.Write(buf[:n]) // 받은 데이터를 그대로 돌려보냄
    }
}
```

`conn.Read()`는 데이터가 도착할 때까지 블로킹된다. 클라이언트가 연결을 종료하면 `io.EOF` 에러를 반환한다. 버퍼 크기(1024)보다 큰 메시지는 여러 번의 Read 호출에 걸쳐 수신될 수 있으므로, 실제 프로토콜에서는 메시지 경계를 처리하는 로직이 필요하다. 일반적으로 길이 접두사(length-prefix) 또는 구분자(delimiter) 방식을 사용한다.

### TCP 연결의 생명주기

```
클라이언트                서버
    │ ── SYN ──────►  │   (3-way handshake)
    │ ◄── SYN+ACK ──  │
    │ ── ACK ──────►  │
    │                   │   연결 수립
    │ ◄──── 데이터 ───► │   데이터 교환
    │                   │
    │ ── FIN ──────►  │   (4-way handshake)
    │ ◄── ACK ────── │
    │ ◄── FIN ────── │
    │ ── ACK ──────► │   연결 종료
```

---

## 30.2 클라이언트 제작

### TCP 클라이언트

`net.Dial()`을 사용하여 서버에 연결한다. `Dial`은 TCP 3-way handshake를 수행하고, 연결이 성공하면 `net.Conn`을 반환한다:

```go
// 서버에 연결
conn, err := net.Dial("tcp", "localhost:9000")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// 메시지 전송
conn.Write([]byte("Hello, Server!"))

// 응답 수신
buf := make([]byte, 1024)
n, _ := conn.Read(buf)
fmt.Println("응답:", string(buf[:n]))
```

`net.DialTimeout()`을 사용하면 연결 타임아웃을 설정할 수 있다. 네트워크 상태가 불안정한 환경에서는 반드시 타임아웃을 설정해야 한다:

```go
conn, err := net.DialTimeout("tcp", "localhost:9000", 5*time.Second)
```

### 고루틴을 활용한 비동기 통신

클라이언트에서 메시지 송신과 수신을 동시에 처리한다. 수신은 별도 고루틴에서 처리하고, 송신은 메인 고루틴에서 사용자 입력을 읽어 처리한다:

```go
// 수신 고루틴
go func() {
    buf := make([]byte, 1024)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            return
        }
        fmt.Println("받음:", string(buf[:n]))
    }
}()

// 송신 (메인 고루틴)
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    conn.Write([]byte(scanner.Text()))
}
```

이 패턴은 채팅 클라이언트의 기본 구조이다. 수신 고루틴이 종료되면 메인 고루틴에 알려야 하는데, 이를 위해 `done` 채널을 사용할 수 있다.

---

## 30.3 채팅 서버 제작

### 다중 클라이언트 채팅

채팅 서버는 여러 클라이언트의 연결을 관리하고, 한 클라이언트의 메시지를 모든 클라이언트에게 전달한다. 이를 브로드캐스트(broadcast)라고 한다.

핵심 구조:

```go
type ChatServer struct {
    clients    map[net.Conn]string    // 연결 -> 닉네임
    broadcast  chan string            // 브로드캐스트 메시지 채널
    register   chan net.Conn          // 새 클라이언트 등록
    unregister chan net.Conn          // 클라이언트 해제
    mu         sync.RWMutex
}
```

`ChatServer`는 채널을 사용하여 클라이언트 등록, 해제, 메시지 브로드캐스트를 처리한다. 채널을 사용하면 뮤텍스 없이도 안전하게 공유 상태를 관리할 수 있다. `register`와 `unregister` 채널은 클라이언트 맵의 동시 접근을 직렬화하고, `broadcast` 채널은 메시지 전달을 직렬화한다.

### 브로드캐스트 패턴

채널을 사용하여 메시지를 모든 클라이언트에게 전달한다:

```go
func (s *ChatServer) broadcastLoop() {
    for msg := range s.broadcast {
        s.mu.RLock()
        for conn := range s.clients {
            conn.Write([]byte(msg + "\n"))
        }
        s.mu.RUnlock()
    }
}
```

`RLock()`(읽기 잠금)을 사용하므로 여러 고루틴이 동시에 클라이언트 목록을 읽을 수 있다. 클라이언트가 추가/제거될 때만 `Lock()`(쓰기 잠금)을 사용한다. `Write` 호출이 실패하면 해당 클라이언트의 연결이 끊어진 것이므로, 에러 처리 후 `unregister` 채널로 알려야 한다.

### 이벤트 루프 패턴

select 문을 사용하여 여러 채널의 이벤트를 단일 고루틴에서 처리하는 패턴이다:

```go
func (s *ChatServer) run() {
    for {
        select {
        case conn := <-s.register:
            s.mu.Lock()
            s.clients[conn] = ""
            s.mu.Unlock()
        case conn := <-s.unregister:
            s.mu.Lock()
            delete(s.clients, conn)
            conn.Close()
            s.mu.Unlock()
        case msg := <-s.broadcast:
            s.mu.RLock()
            for conn := range s.clients {
                conn.Write([]byte(msg + "\n"))
            }
            s.mu.RUnlock()
        }
    }
}
```

---

## 30.4 gRPC란?

### gRPC 개요

gRPC는 Google이 개발한 고성능 RPC(Remote Procedure Call) 프레임워크이다. RPC는 원격 서버의 함수를 마치 로컬 함수를 호출하듯 사용할 수 있게 하는 기술이다. gRPC는 HTTP/2를 전송 프로토콜로, Protocol Buffers를 직렬화 형식으로 사용하여 높은 성능과 강한 타입 안전성을 제공한다:

```
┌──────────────┐         ┌──────────────┐
│   클라이언트   │ ──────► │    서버      │
│              │ HTTP/2  │              │
│  Stub        │ ◄────── │  Service     │
│  (자동 생성)  │ Protobuf│  (구현)      │
└──────────────┘         └──────────────┘
```

클라이언트의 Stub(스텁)은 `.proto` 파일에서 자동 생성된 코드이다. 개발자는 스텁의 메서드를 호출하기만 하면 되고, 직렬화, 네트워크 통신, 역직렬화는 gRPC 프레임워크가 자동으로 처리한다.

### gRPC vs REST

| 항목 | gRPC | REST |
|------|------|------|
| 프로토콜 | HTTP/2 | HTTP/1.1 (주로) |
| 직렬화 | Protocol Buffers (바이너리) | JSON (텍스트) |
| 스트리밍 | 양방향 스트리밍 지원 | 제한적 (SSE, WebSocket) |
| 코드 생성 | .proto에서 자동 생성 | 수동 작성 |
| 타입 안전성 | 컴파일 타임 검증 | 런타임 검증 |
| 성능 | 바이너리 직렬화로 빠르다 | 텍스트 직렬화로 상대적으로 느리다 |
| 브라우저 지원 | 제한적 (gRPC-Web 필요) | 완전 지원 |
| API 디버깅 | 바이너리라 직접 읽기 어렵다 | JSON이라 읽기 쉽다 |
| 적합한 용도 | 마이크로서비스 간 통신 | 공개 API, 웹 클라이언트 |

gRPC는 마이크로서비스 간의 내부 통신에 적합하고, REST는 외부 클라이언트(웹 브라우저, 모바일 앱)에 적합하다. 많은 시스템에서 두 가지를 함께 사용하여, 내부 통신은 gRPC로, 외부 API는 REST로 제공한다.

### Protocol Buffers (protobuf)

gRPC는 Protocol Buffers를 사용하여 서비스와 메시지를 정의한다. protobuf는 Google이 개발한 언어 중립적인 직렬화 형식으로, JSON보다 빠르고 작다:

```protobuf
syntax = "proto3";

package chat;

option go_package = "./chatpb";

// 서비스 정의
service ChatService {
    // 단순 RPC
    rpc SendMessage(ChatMessage) returns (ChatResponse);

    // 서버 스트리밍 RPC
    rpc Subscribe(SubscribeRequest) returns (stream ChatMessage);

    // 양방향 스트리밍 RPC
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}

// 메시지 정의
message ChatMessage {
    string user = 1;      // 필드 번호 1
    string content = 2;   // 필드 번호 2
    int64 timestamp = 3;  // 필드 번호 3
}

message ChatResponse {
    bool success = 1;
    string message = 2;
}

message SubscribeRequest {
    string user = 1;
}
```

각 필드에 할당된 번호(1, 2, 3)는 바이너리 형식에서 필드를 식별하는 데 사용한다. 한번 할당된 번호는 변경하면 안 된다. 필드를 제거할 때는 `reserved` 키워드를 사용하여 해당 번호가 재사용되지 않도록 해야 한다.

### gRPC 통신 유형

1. **단순 RPC (Unary)**: 요청 1개 -> 응답 1개. 일반적인 함수 호출과 유사하다.
2. **서버 스트리밍**: 요청 1개 -> 응답 여러 개. 서버가 데이터를 순차적으로 전송한다. 실시간 피드, 대용량 데이터 조회에 적합하다.
3. **클라이언트 스트리밍**: 요청 여러 개 -> 응답 1개. 클라이언트가 데이터를 순차적으로 전송한다. 파일 업로드, 로그 수집에 적합하다.
4. **양방향 스트리밍**: 요청/응답 여러 개. 양쪽이 독립적으로 데이터를 주고받는다. 채팅, 실시간 게임에 적합하다.

### 설치 및 코드 생성

```bash
# Protocol Buffers 컴파일러 설치
brew install protobuf

# Go용 gRPC 플러그인 설치
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# PATH에 $GOPATH/bin 추가 (zsh 기준)
export PATH="$PATH:$(go env GOPATH)/bin"

# .proto 파일에서 Go 코드 생성
protoc --go_out=. --go-grpc_out=. proto/chat.proto
```

`protoc` 명령은 `.proto` 파일을 입력으로 받아 두 가지 Go 파일을 생성한다:
- `chat.pb.go`: 메시지 타입의 직렬화/역직렬화 코드
- `chat_grpc.pb.go`: 서비스의 클라이언트 스텁과 서버 인터페이스

---

## 30.5 gRPC를 이용한 채팅 프로그램

### 서버 측 구현

```go
type chatServer struct {
    pb.UnimplementedChatServiceServer
    mu      sync.RWMutex
    clients map[string]pb.ChatService_ChatServer
}

// Chat - 양방향 스트리밍 RPC 구현
func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
    // 첫 메시지에서 사용자 이름을 추출하여 등록
    msg, err := stream.Recv()
    if err != nil {
        return err
    }
    username := msg.User

    s.mu.Lock()
    s.clients[username] = stream
    s.mu.Unlock()

    defer func() {
        s.mu.Lock()
        delete(s.clients, username)
        s.mu.Unlock()
    }()

    // 스트림에서 메시지를 수신하고 모든 클라이언트에게 전달
    for {
        msg, err := stream.Recv()
        if err != nil {
            return err
        }
        s.broadcast(msg)
    }
}

func (s *chatServer) broadcast(msg *pb.ChatMessage) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, stream := range s.clients {
        stream.Send(msg)
    }
}
```

`UnimplementedChatServiceServer`를 임베딩하면 서비스 인터페이스의 모든 메서드에 대한 기본 구현이 제공된다. 이를 통해 새로운 RPC 메서드가 추가되어도 기존 서버 코드가 컴파일 에러 없이 동작한다. 구현하지 않은 메서드를 호출하면 "Unimplemented" 에러를 반환한다.

### gRPC 서버 시작

```go
func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("리스닝 실패: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterChatServiceServer(grpcServer, &chatServer{
        clients: make(map[string]pb.ChatService_ChatServer),
    })

    log.Println("gRPC 서버 시작: :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("서버 실행 실패: %v", err)
    }
}
```

### 클라이언트 측 구현

```go
// 서버에 연결
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := pb.NewChatServiceClient(conn)

// 양방향 스트리밍 연결
stream, err := client.Chat(context.Background())
if err != nil {
    log.Fatal(err)
}

// 수신 고루틴
go func() {
    for {
        msg, err := stream.Recv()
        if err != nil {
            return
        }
        fmt.Printf("[%s] %s\n", msg.User, msg.Content)
    }
}()

// 송신
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    stream.Send(&pb.ChatMessage{
        User:      "사용자",
        Content:   scanner.Text(),
        Timestamp: time.Now().Unix(),
    })
}
```

`grpc.WithInsecure()`는 TLS 없이 연결한다는 의미이다. 프로덕션 환경에서는 반드시 TLS를 사용해야 하며, `grpc.WithTransportCredentials()`로 인증서를 설정한다.

### 전체 흐름

```
1. chat.proto 정의
       │
       ▼
2. protoc으로 Go 코드 생성 (chat.pb.go, chat_grpc.pb.go)
       │
       ▼
3. 서버: ChatServiceServer 인터페이스 구현
       │
       ▼
4. 클라이언트: 생성된 Stub으로 RPC 호출
```

### gRPC 인터셉터

gRPC에서 미들웨어 역할을 하는 것이 인터셉터(Interceptor)이다. HTTP 미들웨어와 유사하게 로깅, 인증, 메트릭 수집 등에 사용한다:

```go
// 단순 RPC용 인터셉터
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    log.Printf("RPC: %s, 소요시간: %v, 에러: %v",
        info.FullMethod, time.Since(start), err)
    return resp, err
}

// 인터셉터 적용
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(loggingInterceptor),
)
```

---

## 핵심 요약

1. `net` 패키지의 `net.Listen()`과 `net.Dial()`로 TCP 서버/클라이언트를 만든다.
2. 고루틴과 채널을 활용하면 다중 클라이언트 채팅 서버를 구현할 수 있다.
3. **gRPC**는 HTTP/2 기반의 고성능 RPC 프레임워크이다.
4. **Protocol Buffers**로 서비스와 메시지를 정의하고, 코드를 자동 생성한다.
5. gRPC의 **양방향 스트리밍**은 실시간 채팅에 적합하다.
6. `protoc` 명령으로 `.proto` 파일에서 Go 코드를 생성한다.
7. TCP 프로토콜에서는 메시지 경계 처리가 중요하다.
8. gRPC 인터셉터로 로깅, 인증 등 횡단 관심사를 처리한다.

---

## 연습문제

### 연습문제 1: 에코 서버 확장
에코 서버를 다음과 같이 확장하라:
- 메시지를 대문자로 변환하여 반환한다
- 클라이언트별 메시지 카운트를 추적한다
- 타임아웃 처리를 추가한다 (30초 동안 무응답 시 연결 종료)

### 연습문제 2: 채팅방 기능
채팅 서버에 채팅방 기능을 추가하라:
- `/join 방이름` - 채팅방에 참여한다
- `/leave` - 채팅방을 나간다
- `/list` - 현재 방 목록을 표시한다
- 같은 방에 있는 사용자에게만 메시지를 전달한다

### 연습문제 3: gRPC 기능 확장
gRPC 채팅 프로그램에 다음 기능을 추가하라:
- 접속 중인 사용자 목록 조회 (단순 RPC)
- 최근 메시지 히스토리 조회 (서버 스트리밍)
- 귓속말 기능 (특정 사용자에게만 메시지 전달)

### 연습문제 4: 메시지 프로토콜 설계
TCP 에코 서버에 길이 접두사(length-prefix) 프로토콜을 구현하라:
- 메시지 앞에 4바이트 길이 정보를 추가한다
- 길이 정보를 먼저 읽은 후, 해당 길이만큼 메시지를 읽는다
- 버퍼 크기보다 큰 메시지도 올바르게 처리되는지 테스트한다

### 연습문제 5: 연결 상태 관리
TCP 채팅 서버에 연결 상태 관리를 추가하라:
- 클라이언트가 접속하면 "XXX님이 입장했다" 메시지를 브로드캐스트한다
- 클라이언트가 퇴장하면 "XXX님이 퇴장했다" 메시지를 브로드캐스트한다
- `/who` 명령으로 현재 접속자 목록을 확인할 수 있게 한다
- 하트비트(heartbeat)를 구현하여 비정상 종료를 감지한다

### 연습문제 6: gRPC 에러 처리
gRPC 채팅 서버에 적절한 에러 처리를 추가하라:
- `google.golang.org/grpc/status` 패키지를 사용하여 gRPC 상태 코드를 반환한다
- 닉네임 중복 시 `AlreadyExists` 에러를 반환한다
- 존재하지 않는 사용자에게 귓속말 시 `NotFound` 에러를 반환한다
- 클라이언트에서 에러를 적절히 처리한다

### 연습문제 7: protobuf 메시지 설계
다음 시나리오에 적합한 protobuf 메시지를 설계하라:
- 채팅 메시지 (텍스트, 이미지, 파일 등 여러 타입)
- `oneof` 키워드를 사용하여 메시지 타입에 따라 다른 필드를 갖게 한다
- `enum`을 사용하여 메시지 상태(전송 중, 전송 완료, 읽음)를 정의한다

### 연습문제 8: TCP 서버 벤치마크
에코 서버의 처리량(throughput)을 측정하라:
- 동시 클라이언트 수(1, 10, 100, 1000)별 처리량을 측정한다
- 메시지 크기(100B, 1KB, 10KB)별 처리량을 측정한다
- `testing.B`를 사용한 벤치마크를 작성한다

### 연습문제 9: gRPC 인터셉터 구현
다음 gRPC 인터셉터를 구현하라:
- 요청 로깅 인터셉터: RPC 이름, 소요 시간, 에러 여부를 로그에 기록한다
- 인증 인터셉터: 메타데이터에서 토큰을 검증한다
- 스트리밍 인터셉터: 스트림 메시지 수를 카운트한다

### 연습문제 10: 연결 복구
gRPC 클라이언트에 자동 재연결 기능을 구현하라:
- 서버 연결이 끊어지면 자동으로 재연결을 시도한다
- 지수 백오프(exponential backoff) 전략을 사용한다
- 최대 재시도 횟수를 설정한다
- 재연결 성공 시 채팅 스트림을 다시 설정한다

---

## 구현 과제

### 과제 1: 멀티룸 TCP 채팅 서버
여러 채팅방을 지원하는 TCP 채팅 서버를 완전히 구현하라:
- 채팅방 생성, 참여, 나가기, 삭제 기능
- 방별 사용자 목록 확인
- 방 내 브로드캐스트와 귓속말 기능
- 닉네임 변경 기능 (`/nick 새닉네임`)
- 서버 상태를 모니터링하는 관리자 명령어

### 과제 2: gRPC 파일 전송 서비스
gRPC를 이용한 파일 전송 서비스를 구현하라:
- 클라이언트 스트리밍으로 파일을 업로드한다 (청크 단위 전송)
- 서버 스트리밍으로 파일을 다운로드한다
- 단순 RPC로 파일 목록을 조회한다
- 전송 진행률을 표시한다
- 대용량 파일(100MB 이상)도 메모리 효율적으로 처리한다

### 과제 3: 프로토콜 변환 게이트웨이
gRPC 서비스를 REST API로 노출하는 게이트웨이를 구현하라:
- HTTP 요청을 받아 gRPC 호출로 변환한다
- gRPC 응답을 JSON으로 변환하여 HTTP 응답으로 반환한다
- `POST /api/chat/send` -> `ChatService.SendMessage` RPC 호출
- `GET /api/chat/messages` -> `ChatService.Subscribe` RPC 호출 (SSE로 스트리밍)

### 과제 4: TCP 채팅 클라이언트 TUI
`tcell` 또는 `bubbletea` 라이브러리를 사용하여 터미널 기반 채팅 클라이언트 UI를 구현하라:
- 화면을 메시지 영역과 입력 영역으로 분리한다
- 사용자 목록을 사이드바에 표시한다
- 색상으로 사용자별 메시지를 구분한다
- 스크롤 기능을 구현한다

### 과제 5: gRPC 헬스 체크 서비스
gRPC 서버에 표준 헬스 체크 서비스를 추가하라:
- `grpc_health_v1.Health` 서비스를 구현한다
- 서버 상태(SERVING, NOT_SERVING)를 관리한다
- 클라이언트에서 헬스 체크를 수행하는 유틸리티를 작성한다
- 워치(Watch) 스트림을 사용하여 상태 변화를 실시간으로 감지한다

---

## 프로젝트 과제

### 프로젝트 1: 실시간 협업 도구
gRPC 양방향 스트리밍을 사용한 실시간 협업 도구를 만들어라:
- **공유 메모장**: 여러 사용자가 동시에 텍스트를 편집한다. 각 사용자의 변경 사항이 실시간으로 다른 사용자에게 반영된다.
- **채팅 기능**: 편집 중에 사용자 간 채팅이 가능하다.
- **사용자 커서 위치 공유**: 다른 사용자의 커서 위치를 표시한다.
- **변경 히스토리**: 최근 변경 내역을 조회할 수 있다.
- protobuf 메시지를 적절히 설계하고, 충돌 처리 전략을 구현한다.

### 프로젝트 2: 마이크로서비스 채팅 시스템
채팅 시스템을 여러 마이크로서비스로 분리하여 구현하라:
- **인증 서비스**: 사용자 등록, 로그인, 토큰 발급 (gRPC)
- **채팅 서비스**: 메시지 송수신, 채팅방 관리 (gRPC 양방향 스트리밍)
- **히스토리 서비스**: 메시지 저장 및 조회 (gRPC)
- **API 게이트웨이**: REST API를 gRPC 호출로 변환 (HTTP 서버)
- 각 서비스 간 통신은 gRPC를 사용한다.
- 서비스 간 인증은 gRPC 메타데이터와 인터셉터를 사용한다.
- 각 서비스에 대한 테스트를 작성한다.
