# 23장 채널과 컨텍스트

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch23_channel_context/channel_basic.go
go run ch23_channel_context/buffered.go
go run ch23_channel_context/select_example.go
go run ch23_channel_context/context_cancel.go
go run ch23_channel_context/context_timeout.go
```

> **Makefile 활용**: `make run CH=ch23_channel_context` 또는 `make run CH=ch23_channel_context FILE=channel_basic.go`

---

Go의 동시성 프로그래밍에서 가장 중요한 두 가지 도구인 **채널(Channel)**과 **컨텍스트(Context)**를 학습한다. Go의 철학인 "메모리를 공유하여 통신하지 말고, 통신으로 메모리를 공유하라(Don't communicate by sharing memory; share memory by communicating)"를 실현하는 핵심 도구이다. 채널은 고루틴 간의 데이터 전달과 동기화를 담당하고, 컨텍스트는 고루틴의 생명주기 관리와 취소 전파를 담당한다.

---

## 23.1 채널 사용하기

### 채널이란?

채널은 고루틴 간에 데이터를 주고받는 **통신 파이프**이다. 채널은 타입이 지정되어 있어, 선언된 타입의 값만 주고받을 수 있다. 채널 연산(`<-`)은 고루틴을 블로킹할 수 있으므로, 동기화 도구로도 사용된다.

```go
ch := make(chan int)    // int 타입 채널 생성
ch <- 42               // 채널에 값 보내기
value := <-ch           // 채널에서 값 받기
```

채널의 제로값은 `nil`이다. `nil` 채널에 대한 송수신은 영원히 블로킹되므로, 채널을 사용하기 전에 반드시 `make`로 초기화해야 한다.

### 언버퍼드 채널 (Unbuffered Channel)

버퍼가 없는 채널은 **동기적**이다. 보내는 쪽과 받는 쪽이 모두 준비될 때까지 블로킹된다. 이는 두 고루틴 간의 동기화 지점(synchronization point)을 만드는 효과가 있다.

```go
ch := make(chan int)    // 버퍼 크기 0 (언버퍼드)

go func() {
    ch <- 42            // 받는 쪽이 준비될 때까지 대기
}()

value := <-ch           // 보내는 쪽이 준비될 때까지 대기
```

언버퍼드 채널은 두 고루틴이 반드시 만나야(rendezvous) 하므로, 값의 전달과 동시에 동기화가 보장된다.

### 버퍼드 채널 (Buffered Channel)

버퍼가 있는 채널은 버퍼가 가득 찰 때까지 블로킹되지 않는다. 생산자와 소비자의 속도가 다를 때 버퍼를 완충 장치로 사용할 수 있다.

```go
ch := make(chan int, 3) // 버퍼 크기 3

ch <- 1  // 블로킹 없음
ch <- 2  // 블로킹 없음
ch <- 3  // 블로킹 없음
// ch <- 4  // 버퍼가 가득 차서 블로킹!
```

`len(ch)`로 현재 버퍼에 들어있는 요소 수를, `cap(ch)`로 버퍼 크기를 확인할 수 있다. 버퍼 크기는 성능 튜닝의 대상이 되며, 너무 크면 메모리가 낭비되고 너무 작으면 블로킹이 자주 발생한다.

### 채널 닫기

`close(ch)`로 채널을 닫으면 더 이상 값을 보낼 수 없다. 닫힌 채널에 값을 보내면 panic이 발생한다. 받는 쪽은 남은 값을 모두 받은 후 제로값을 받으며, 두 번째 반환값으로 채널이 닫혔는지 확인할 수 있다.

```go
close(ch)

// range로 채널이 닫힐 때까지 읽기
for value := range ch {
    fmt.Println(value)
}

// 또는 두 번째 반환값으로 확인
value, ok := <-ch  // ok가 false이면 채널이 닫힌 것이다
```

> **원칙:** 채널은 보내는 쪽에서 닫아야 한다. 받는 쪽에서 닫으면 보내는 쪽이 닫힌 채널에 값을 보내 panic이 발생할 수 있다. 여러 송신자가 있는 경우에는 `sync.WaitGroup` 등을 활용하여 모든 송신이 완료된 후 닫는다.

### 채널 방향 지정

함수 매개변수에서 채널의 방향을 제한할 수 있다. 이를 통해 함수가 채널을 잘못 사용하는 것을 컴파일 타임에 방지할 수 있다.

```go
func producer(out chan<- int) { ... }  // 보내기 전용
func consumer(in <-chan int)  { ... }  // 받기 전용
```

양방향 채널은 방향이 제한된 채널에 자동으로 변환되지만, 반대 방향으로의 변환은 허용되지 않는다.

### select 문

여러 채널 연산을 동시에 대기하고, 준비된 것을 실행한다. `switch`와 비슷한 형태이지만, 채널 연산에 특화되어 있다.

```go
select {
case msg := <-ch1:
    fmt.Println("ch1에서 받음:", msg)
case msg := <-ch2:
    fmt.Println("ch2에서 받음:", msg)
case ch3 <- "hello":
    fmt.Println("ch3에 보냄")
default:
    fmt.Println("아무 채널도 준비되지 않음")
}
```

**select의 특징:**
- 여러 case가 동시에 준비되면 **무작위**로 하나를 선택한다. 이는 특정 채널이 기아(starvation) 상태에 빠지는 것을 방지한다.
- `default`가 있으면 블로킹되지 않는다 (비블로킹 채널 연산).
- `default`가 없으면 적어도 하나의 case가 준비될 때까지 블로킹된다.
- 타임아웃 구현에 자주 사용된다.

### 타임아웃 패턴

`select`와 `time.After`를 조합하면 타임아웃 패턴을 구현할 수 있다.

```go
select {
case result := <-ch:
    fmt.Println("결과:", result)
case <-time.After(3 * time.Second):
    fmt.Println("타임아웃!")
}
```

---

## 23.2 컨텍스트 사용하기

### context 패키지

`context` 패키지는 고루틴의 **취소 신호**, **타임아웃**, **값 전달**을 위한 도구이다. 서버 프로그래밍에서 요청의 생명주기를 관리하는 핵심 도구이며, 표준 라이브러리의 많은 함수가 첫 번째 매개변수로 `context.Context`를 받는다.

모든 컨텍스트는 루트 컨텍스트에서 파생된다. `context.Background()`는 최상위 루트 컨텍스트이며, `context.TODO()`는 어떤 컨텍스트를 사용할지 아직 결정하지 못했을 때 사용하는 플레이스홀더이다.

### context.WithCancel

수동으로 취소 신호를 보낼 수 있는 컨텍스트를 생성한다. 반환된 `cancel` 함수를 호출하면 해당 컨텍스트와 그로부터 파생된 모든 자식 컨텍스트에 취소 신호가 전파된다.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // 리소스 누수 방지를 위해 반드시 호출

go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("취소됨:", ctx.Err())
            return
        default:
            // 작업 수행
        }
    }
}(ctx)

cancel() // 취소 신호 전송
```

`ctx.Done()`은 채널을 반환하며, 컨텍스트가 취소되면 이 채널이 닫힌다. `ctx.Err()`는 취소 이유를 반환한다 (`context.Canceled` 또는 `context.DeadlineExceeded`).

### context.WithTimeout

지정된 시간이 지나면 자동으로 취소되는 컨텍스트를 생성한다. 내부적으로 `context.WithDeadline`을 사용하여 현재 시각에 duration을 더한 시점을 데드라인으로 설정한다.

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

select {
case result := <-doWork(ctx):
    fmt.Println("결과:", result)
case <-ctx.Done():
    fmt.Println("타임아웃:", ctx.Err())
}
```

### context.WithDeadline

특정 시점에 만료되는 컨텍스트를 생성한다. `WithTimeout`과 달리 절대적인 시점을 지정한다.

```go
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

`ctx.Deadline()` 메서드로 설정된 데드라인을 조회할 수 있다. 데드라인이 없는 컨텍스트는 `ok` 값이 `false`로 반환된다.

### context.WithValue

컨텍스트에 키-값 쌍을 저장한다. 요청 범위의 데이터 전달에 사용한다.

```go
type contextKey string

ctx := context.WithValue(parentCtx, contextKey("userID"), "user-123")
userID := ctx.Value(contextKey("userID")).(string)
```

> **주의:** `WithValue`는 요청 범위의 메타데이터(요청 ID, 인증 정보 등)에만 사용하고, 함수의 매개변수를 대체하지 않아야 한다. 키 타입은 충돌을 방지하기 위해 비공개 타입(`unexported type`)을 사용하는 것이 관례이다.

### 컨텍스트 사용 규칙

1. **첫 번째 매개변수**로 전달한다: `func DoSomething(ctx context.Context, ...)`
2. **nil을 전달하지 않는다**: 어떤 컨텍스트를 사용할지 모르면 `context.TODO()`를 사용한다.
3. **구조체에 저장하지 않는다**: 함수 매개변수로 전달한다.
4. **cancel을 반드시 호출**한다: 리소스 누수를 방지한다. `defer cancel()`을 사용하는 것이 가장 안전하다.
5. **컨텍스트는 불변(immutable)**이다: `WithValue`, `WithCancel` 등은 새로운 컨텍스트를 반환하며, 원본을 수정하지 않는다.

---

## 23.3 채널 활용 패턴

### 파이프라인 패턴

여러 단계의 처리를 채널로 연결하여 데이터 처리 파이프라인을 구성하는 패턴이다.

```go
// 생성 → 변환 → 출력
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}
```

### Fan-Out / Fan-In 패턴

하나의 채널에서 여러 고루틴이 읽는 것을 Fan-Out, 여러 채널의 결과를 하나의 채널로 합치는 것을 Fan-In이라 한다. CPU 집약적 작업을 병렬화할 때 유용하다.

### Done 채널 패턴

고루틴에게 종료 신호를 보내는 패턴이다. 컨텍스트가 도입되기 전에 많이 사용되었으며, 현재는 `context.WithCancel`이 이를 대체한다.

```go
done := make(chan struct{})
go func() {
    for {
        select {
        case <-done:
            return
        default:
            // 작업 수행
        }
    }
}()
close(done) // 종료 신호
```

---

## 핵심 요약

| 개념 | 설명 |
|------|------|
| 채널 (Channel) | 고루틴 간 데이터 통신 파이프 |
| 언버퍼드 채널 | `make(chan T)` - 동기적, 양쪽 준비될 때까지 블로킹 |
| 버퍼드 채널 | `make(chan T, n)` - 버퍼가 찰 때까지 비블로킹 |
| `close(ch)` | 채널 닫기. `range`와 함께 사용 |
| `chan<-` / `<-chan` | 보내기 전용 / 받기 전용 채널 |
| `select` | 여러 채널 연산 중 준비된 것 실행 |
| `nil` 채널 | 송수신 모두 영원히 블로킹 |
| `context.WithCancel` | 수동 취소 가능한 컨텍스트 |
| `context.WithTimeout` | 타임아웃 자동 취소 컨텍스트 |
| `context.WithDeadline` | 특정 시점에 만료되는 컨텍스트 |
| `context.WithValue` | 요청 범위 값 전달 |

---

## 연습문제

### 문제 1: 파이프라인 패턴
세 단계의 파이프라인을 채널로 구현하라.
1. `generator`: 1~10 숫자를 채널로 보냄
2. `square`: 받은 숫자를 제곱하여 다른 채널로 보냄
3. `printer`: 최종 결과를 출력

### 문제 2: Fan-Out / Fan-In
하나의 작업을 여러 고루틴으로 분배(Fan-Out)하고, 결과를 하나의 채널로 모으는(Fan-In) 패턴을 구현하라.

### 문제 3: 타임아웃 처리
외부 API 호출을 시뮬레이션하고, `context.WithTimeout`으로 3초 타임아웃을 구현하라.
- API 응답이 타임아웃 내에 오면 결과를 출력하라.
- 타임아웃이 되면 에러 메시지를 출력하라.

### 문제 4: 작업 취소
여러 고루틴이 동시에 검색을 수행하다가, 하나의 고루틴이 결과를 찾으면 나머지를 `context.WithCancel`로 취소하는 프로그램을 작성하라.

### 문제 5: 버퍼드 vs 언버퍼드 채널 비교
동일한 작업을 언버퍼드 채널과 버퍼드 채널(크기 1, 10, 100)로 각각 구현하고, 성능 차이를 `time.Since`로 측정하라.
- 생산자가 1000개의 값을 보내고, 소비자가 받는 시나리오를 사용하라.
- 버퍼 크기가 성능에 미치는 영향을 관찰하고 분석하라.

### 문제 6: select를 이용한 다중 채널 처리
3개의 채널에서 동시에 데이터를 받는 프로그램을 작성하라.
- 각 채널은 서로 다른 간격(100ms, 200ms, 500ms)으로 데이터를 보낸다.
- `select`로 어떤 채널에서든 데이터가 오면 즉시 처리하라.
- 5초 후에 모든 채널을 닫고 프로그램을 종료하라.
- 각 채널에서 받은 데이터 수를 출력하라.

### 문제 7: 컨텍스트 체인
부모-자식 컨텍스트 체인을 3단계로 구성하라.
- 최상위 컨텍스트: `WithTimeout(5초)`
- 중간 컨텍스트: `WithValue("requestID", "req-123")`
- 최하위 컨텍스트: `WithCancel`
- 각 단계의 고루틴에서 `ctx.Done()`을 모니터링하고, 부모 컨텍스트가 취소되면 자식도 취소되는 것을 확인하라.

### 문제 8: 채널을 이용한 세마포어
채널을 사용하여 최대 동시 실행 수를 제한하는 세마포어를 구현하라.
- `type Semaphore chan struct{}`를 정의하라.
- `Acquire()`, `Release()`, `TryAcquire(timeout time.Duration) bool` 메서드를 구현하라.
- 20개의 고루틴이 동시에 실행되지만, 최대 5개만 동시에 작업하는 것을 확인하라.

### 문제 9: 채널을 이용한 제한된 재시도
실패할 수 있는 작업을 채널과 컨텍스트를 사용하여 제한된 횟수만큼 재시도하는 함수를 구현하라.
- `func retry(ctx context.Context, maxAttempts int, fn func() error) error`
- 각 재시도 사이에 지수 백오프(exponential backoff)를 적용하라 (1초, 2초, 4초...).
- 컨텍스트가 취소되면 즉시 중단하라.

### 문제 10: 생산자-소비자 패턴
여러 생산자와 여러 소비자가 하나의 버퍼드 채널을 공유하는 패턴을 구현하라.
- 3개의 생산자가 각각 100개의 아이템을 생산하라.
- 5개의 소비자가 채널에서 아이템을 가져와 처리하라.
- 모든 생산이 완료되면 채널을 닫고, 모든 소비자가 종료될 때까지 대기하라.
- 각 소비자가 처리한 아이템 수를 출력하라.

---

## 구현 과제

### 과제 1: 채팅 서버 (채널 기반)
채널을 사용한 간단한 채팅 서버를 구현하라.
- 각 클라이언트(고루틴)는 메시지를 보낼 수 있고, 모든 클라이언트에게 브로드캐스트된다.
- 허브(hub) 고루틴이 채널을 통해 메시지를 수집하고 분배하라.
- 클라이언트 접속/퇴장을 관리하라.
- `context.WithCancel`로 서버 종료 시 모든 클라이언트를 정리하라.

### 과제 2: 파이프라인 데이터 처리기
CSV 데이터를 읽고 변환하고 출력하는 3단계 파이프라인을 구현하라.
- 1단계(reader): 파일에서 CSV 행을 읽어 채널로 보낸다.
- 2단계(transformer): 데이터를 변환(필터링, 매핑)하여 다른 채널로 보낸다.
- 3단계(writer): 최종 결과를 출력하거나 파일에 쓴다.
- `context.WithTimeout`으로 전체 파이프라인의 타임아웃을 설정하라.

### 과제 3: 동시 API 요청기
여러 외부 API를 동시에 호출하고 결과를 집계하는 프로그램을 구현하라.
- 각 API 호출을 고루틴으로 실행하고, 결과를 채널로 수집하라.
- 개별 API에 대한 타임아웃과 전체 타임아웃을 모두 설정하라.
- 일부 API가 실패해도 성공한 결과는 반환하라.
- `context.WithCancel`로 첫 번째 에러 발생 시 나머지 요청을 취소하는 옵션을 제공하라.

### 과제 4: 속도 제한기 (Rate Limiter)
토큰 버킷 알고리즘을 채널로 구현한 속도 제한기를 작성하라.
- `type RateLimiter struct`를 정의하고, `Allow() bool`, `Wait(ctx context.Context) error` 메서드를 구현하라.
- 초당 N개의 요청만 허용하도록 제한하라.
- 버스트(burst)를 지원하라 (버퍼드 채널의 크기로 조절).
- 여러 고루틴에서 동시에 사용해도 안전하게 동작해야 한다.

### 과제 5: 작업 스케줄러
주기적 작업과 일회성 작업을 관리하는 스케줄러를 채널과 컨텍스트로 구현하라.
- `Schedule(interval time.Duration, task func())`: 주기적으로 실행하는 작업 등록
- `RunOnce(delay time.Duration, task func())`: 지연 후 한 번 실행하는 작업 등록
- `Shutdown(ctx context.Context)`: 진행 중인 작업이 완료되면 종료
- 모든 작업은 컨텍스트를 통해 취소할 수 있어야 한다.

---

## 프로젝트 과제

### 프로젝트 1: 분산 작업 큐 시스템
생산자-소비자 패턴을 확장한 분산 작업 큐 시스템을 구현하라.
- 작업 큐(버퍼드 채널)에 작업을 제출하는 여러 생산자를 구현하라.
- 워커 풀(고정된 수의 소비자 고루틴)이 작업을 처리하라.
- 작업 우선순위를 지원하라 (여러 채널과 select를 활용).
- 작업 결과를 결과 채널로 수집하라.
- 진행 상황을 실시간으로 보고하는 모니터 고루틴을 구현하라.
- `context.WithTimeout`으로 전체 처리 시간을 제한하고, 타임아웃 시 진행 중인 작업을 정리하라.
- 작업 실패 시 재시도 로직을 포함하라.

### 프로젝트 2: 실시간 로그 수집 및 분석 시스템
여러 로그 소스에서 실시간으로 로그를 수집하고 분석하는 시스템을 구현하라.
- 3개 이상의 로그 생산자가 서로 다른 형식의 로그를 생성하라 (INFO, WARN, ERROR 레벨).
- 채널 파이프라인으로 로그를 수집 → 파싱 → 필터링 → 집계하라.
- ERROR 레벨 로그가 10초 내에 5개 이상 발생하면 경고를 출력하라.
- 30초마다 로그 통계(레벨별 개수, 초당 처리량)를 출력하라.
- `context.WithCancel`로 전체 시스템을 깔끔하게 종료하라.
- 각 단계에서 발생하는 에러를 적절히 처리하라.
