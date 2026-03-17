# 27장 프로파일링으로 성능 개선하기

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch27_profiling/cpu_profile.go
go run ch27_profiling/mem_profile.go
go run ch27_profiling/server_profile.go

# 프로파일 분석 (cpu.prof, mem.prof 생성 후)
go tool pprof cpu.prof
go tool pprof mem.prof

# 웹 UI로 프로파일 분석
go tool pprof -http=:8080 cpu.prof
```

> **참고**: `go tool pprof`는 프로파일 파일을 분석하는 인터랙티브 도구이다. `top`, `list`, `web` 등의 명령어를 사용하여 병목 지점을 파악할 수 있다.

> **Makefile 활용**: `make run CH=ch27_profiling` 또는 `make run CH=ch27_profiling FILE=cpu_profile.go`

---

프로파일링은 프로그램이 어디에서 시간과 메모리를 가장 많이 소비하는지 분석하는 기법이다. Go는 `runtime/pprof`와 `net/http/pprof` 패키지를 통해 강력한 프로파일링 도구를 제공한다. 프로파일링은 성능 최적화의 첫 번째 단계이다. 추측이 아닌 데이터 기반으로 병목 지점을 찾아야 효과적인 최적화가 가능하다. Donald Knuth의 말처럼 "조기 최적화는 만악의 근원"이므로, 프로파일링 결과를 확인한 후에 최적화를 진행하는 것이 올바른 접근이다.

---

## 27.1 특정 구간 프로파일링

### pprof 패키지

`runtime/pprof` 패키지를 사용하면 프로그램의 특정 구간에서 CPU와 메모리 사용량을 프로파일링할 수 있다. 이 패키지는 Go 런타임과 긴밀하게 통합되어 있으며, 프로파일링에 따른 오버헤드가 매우 작아 프로덕션 환경에서도 사용할 수 있다.

### CPU 프로파일링

CPU 프로파일은 프로그램이 실행되는 동안 각 함수가 CPU 시간을 얼마나 사용하는지 기록한다. Go의 CPU 프로파일러는 기본적으로 100Hz(초당 100회)로 샘플링하며, 각 샘플에서 현재 실행 중인 고루틴의 스택 트레이스를 기록한다. 이 샘플을 모아 어떤 함수가 가장 많은 CPU 시간을 소비하는지 통계적으로 분석한다.

```go
import "runtime/pprof"

// CPU 프로파일 파일 생성
f, err := os.Create("cpu.prof")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

// CPU 프로파일링 시작
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()

// 프로파일링할 코드 실행
doWork()
```

`StartCPUProfile()`과 `StopCPUProfile()`은 프로그램 전체에서 한 번만 호출할 수 있다. 동시에 두 개의 CPU 프로파일을 수집하는 것은 불가능하다.

### 프로파일 분석

프로파일 파일을 생성한 후 `go tool pprof` 명령으로 분석한다:

```bash
# 프로그램 실행하여 프로파일 생성
go run cpu_profile.go

# 프로파일 분석 (인터랙티브 모드)
go tool pprof cpu.prof

# 자주 사용하는 pprof 명령어:
# top10     - CPU 사용량 상위 10개 함수
# top20     - CPU 사용량 상위 20개 함수
# list 함수명 - 특정 함수의 라인별 CPU 사용량
# web       - 호출 그래프를 웹 브라우저로 표시
# pdf       - PDF 형식으로 호출 그래프 저장
# peek 함수명 - 특정 함수의 호출자/피호출자 정보
```

### pprof 인터랙티브 모드의 주요 명령어

| 명령어 | 설명 |
|--------|------|
| `top [N]` | CPU/메모리 사용량 상위 N개 함수를 표시한다 |
| `list <함수명>` | 함수의 소스 코드와 라인별 비용을 표시한다 |
| `web` | 호출 그래프를 SVG로 생성하여 브라우저에서 표시한다 |
| `peek <함수명>` | 함수의 호출자와 피호출자를 표시한다 |
| `disasm <함수명>` | 함수의 어셈블리와 비용을 표시한다 |
| `cum` | 누적(cumulative) 비용 기준으로 정렬한다 |
| `flat` | 자체(flat) 비용 기준으로 정렬한다 |

**flat vs cum**: flat은 해당 함수 자체에서 소비한 시간이고, cum(cumulative)은 해당 함수와 그 함수가 호출한 모든 함수의 시간을 합산한 것이다. flat이 높은 함수는 직접적인 병목이고, cum이 높지만 flat이 낮은 함수는 호출하는 하위 함수에 병목이 있다.

### 메모리 프로파일링

메모리 프로파일은 힙 메모리 할당 패턴을 분석한다. 어떤 함수가 얼마나 많은 메모리를 할당하는지 파악하여 메모리 사용량을 최적화하는 데 활용한다:

```go
import "runtime/pprof"

// 메모리 프로파일 저장
f, err := os.Create("mem.prof")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

// 메모리 통계를 최신으로 갱신
runtime.GC()

// 힙 프로파일 기록
pprof.WriteHeapProfile(f)
```

`runtime.GC()`를 호출하면 가비지 컬렉션을 실행하여 더 이상 사용하지 않는 메모리를 회수한다. 이를 통해 현재 실제로 사용 중인 메모리만 프로파일에 포함시킬 수 있다.

메모리 프로파일 분석:

```bash
go tool pprof mem.prof

# 메모리 할당량 기준으로 분석 (프로그램 전체 수명 동안 할당된 총량)
go tool pprof -alloc_space mem.prof

# 현재 사용 중인 메모리 기준 (메모리 누수 탐지에 유용)
go tool pprof -inuse_space mem.prof

# 할당 횟수 기준 (GC 부하 분석에 유용)
go tool pprof -alloc_objects mem.prof
```

`-alloc_space`와 `-inuse_space`의 차이를 이해하는 것이 중요하다. `-alloc_space`는 프로그램이 실행되면서 총 할당한 메모리를 보여준다. 이미 해제된 메모리도 포함된다. `-inuse_space`는 프로파일 시점에 실제로 사용 중인 메모리만 보여주므로 메모리 누수를 찾는 데 유용하다.

### 벤치마크에서 프로파일링

`go test` 명령으로 벤치마크를 실행하면서 프로파일을 생성할 수도 있다. 이 방법은 별도의 프로파일링 코드를 작성할 필요가 없어 편리하다:

```bash
# CPU 프로파일 생성
go test -bench . -cpuprofile cpu.prof

# 메모리 프로파일 생성
go test -bench . -memprofile mem.prof

# 블록 프로파일 생성 (고루틴 블로킹 분석)
go test -bench . -blockprofile block.prof

# 뮤텍스 프로파일 생성 (뮤텍스 경합 분석)
go test -bench . -mutexprofile mutex.prof

# 분석
go tool pprof cpu.prof
```

---

## 27.2 서버에서 프로파일링

### net/http/pprof

실행 중인 웹 서버에서 프로파일링을 수행하려면 `net/http/pprof` 패키지를 import한다. 이 패키지는 import만 해도 자동으로 HTTP 핸들러가 `DefaultServeMux`에 등록된다. 이를 "빈 import" 또는 "side-effect import"라고 한다:

```go
import _ "net/http/pprof" // 빈 import: init() 함수만 실행
```

`_` (blank identifier)를 사용하면 패키지의 `init()` 함수만 실행되고, 패키지의 다른 심볼은 사용하지 않는다. `net/http/pprof`의 `init()` 함수는 `DefaultServeMux`에 프로파일링 관련 HTTP 핸들러를 등록한다.

### 사용 가능한 프로파일 엔드포인트

서버를 실행하면 다음 URL에서 프로파일 데이터에 접근할 수 있다:

| URL | 설명 |
|-----|------|
| `/debug/pprof/` | 프로파일 인덱스 페이지이다 |
| `/debug/pprof/profile` | CPU 프로파일이다 (30초 기본) |
| `/debug/pprof/heap` | 힙 메모리 프로파일이다 |
| `/debug/pprof/goroutine` | 고루틴 프로파일이다 |
| `/debug/pprof/block` | 블로킹 프로파일이다 |
| `/debug/pprof/threadcreate` | OS 스레드 생성 프로파일이다 |
| `/debug/pprof/mutex` | 뮤텍스 경합 프로파일이다 |
| `/debug/pprof/trace` | 실행 트레이스이다 |
| `/debug/pprof/allocs` | 메모리 할당 샘플링 데이터이다 |
| `/debug/pprof/cmdline` | 프로세스 커맨드라인 인수이다 |

### 서버 프로파일 분석

```bash
# 서버 실행 중에 CPU 프로파일 수집 (30초)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 힙 메모리 프로파일 수집
go tool pprof http://localhost:6060/debug/pprof/heap

# 고루틴 덤프
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 웹 UI로 분석 (추천)
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap

# 실행 트레이스 수집 (5초)
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
go tool trace trace.out
```

### 웹 UI 활용

`go tool pprof -http=:8080 프로파일` 명령을 사용하면 웹 브라우저에서 시각적으로 분석할 수 있다:

- **Graph**: 함수 호출 그래프이다. 노드 크기와 화살표 두께로 CPU/메모리 사용량을 시각적으로 파악할 수 있다.
- **Flame Graph**: 플레임 그래프이다. 가로축은 CPU 시간 비율을, 세로축은 호출 스택 깊이를 나타낸다. 넓은 막대가 병목 지점이다.
- **Top**: 가장 많은 리소스를 사용하는 함수 목록이다. flat과 cum 기준으로 정렬할 수 있다.
- **Source**: 소스 코드 레벨의 프로파일링 결과이다. 라인별로 CPU 시간이나 메모리 할당량을 확인할 수 있다.
- **Peek**: 특정 함수의 호출 관계를 표시한다.

### 프로덕션 환경에서의 프로파일링

프로덕션 환경에서 `net/http/pprof`를 사용할 때는 보안에 주의해야 한다:

```go
// 프로파일링용 서버를 별도 포트에서 실행
go func() {
    mux := http.NewServeMux()
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    mux.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
    // 내부 네트워크에서만 접근 가능한 포트
    http.ListenAndServe("localhost:6060", mux)
}()
```

### 프로파일링 팁

1. **프로덕션 환경 주의**: `net/http/pprof`는 디버깅용이므로 프로덕션에서는 별도 포트로 분리하고, 외부에서 접근할 수 없도록 한다.
2. **충분한 시간 수집**: CPU 프로파일은 최소 30초 이상 수집해야 의미 있는 결과를 얻는다. 샘플 수가 적으면 통계적 신뢰도가 낮다.
3. **GC 후 메모리 프로파일**: 메모리 프로파일 전에 `runtime.GC()`를 호출하면 더 정확한 결과를 얻는다.
4. **비교 분석**: 최적화 전후의 프로파일을 비교하면 개선 효과를 확인할 수 있다.
5. **반복 측정**: 한 번의 프로파일링으로 결론을 내지 말고, 여러 번 반복하여 일관된 결과를 확인한다.

```bash
# 두 프로파일 비교
go tool pprof -diff_base=before.prof after.prof
```

### Go 트레이스 도구

프로파일링과 함께 `go tool trace`도 유용한 도구이다. 트레이스는 프로파일링보다 더 세밀한 수준에서 고루틴 스케줄링, 시스템 콜, GC 활동 등을 시간순으로 기록한다:

```bash
# 트레이스 생성
go test -trace trace.out

# 트레이스 분석 (웹 UI)
go tool trace trace.out
```

트레이스 웹 UI에서는 타임라인 형태로 고루틴의 실행, 대기, GC 등을 시각적으로 확인할 수 있다.

---

## 핵심 요약

1. `runtime/pprof`로 프로그램의 특정 구간을 CPU/메모리 프로파일링할 수 있다.
2. `pprof.StartCPUProfile()`과 `pprof.StopCPUProfile()`로 CPU 프로파일을 수집한다.
3. `pprof.WriteHeapProfile()`로 메모리 프로파일을 수집한다.
4. `net/http/pprof`를 import하면 실행 중인 서버에서 프로파일링이 가능하다.
5. `go tool pprof`로 프로파일을 분석하고, `-http` 옵션으로 웹 UI를 활용한다.
6. 벤치마크에서도 `-cpuprofile`, `-memprofile` 플래그로 프로파일을 생성할 수 있다.
7. flat과 cum의 차이를 이해하면 병목 지점을 더 정확하게 파악할 수 있다.
8. `go tool trace`로 고루틴 스케줄링과 GC 활동을 시간순으로 분석할 수 있다.

---

## 연습문제

### 연습문제 1: CPU 프로파일링 실습
다음 함수의 CPU 프로파일을 생성하고 분석하라:
- 큰 슬라이스(10만 개 요소)를 버블 정렬하는 함수를 작성한다.
- `go tool pprof`의 `top` 명령으로 병목 구간을 찾는다.
- `list` 명령으로 해당 함수의 라인별 CPU 사용량을 확인한다.

### 연습문제 2: 메모리 프로파일링 실습
다음 두 가지 구현의 메모리 사용량을 프로파일링으로 비교하라:
- 슬라이스를 `append`로 키워가는 방식
- `make([]int, 0, capacity)`로 미리 용량을 할당하는 방식
- `-alloc_space`와 `-inuse_space` 결과의 차이를 설명하라.

### 연습문제 3: 서버 프로파일링
HTTP 서버를 만들고 `net/http/pprof`를 통해 다음을 확인하라:
- 부하 테스트 중 CPU 사용량 상위 함수
- 고루틴 수 변화
- 힙 메모리 사용 패턴

### 연습문제 4: flat vs cum 분석
다음과 같은 함수 호출 구조를 만들고 CPU 프로파일을 분석하라:
```
main() -> processAll() -> processItem() -> computeHash()
```
각 함수의 flat과 cum 값을 비교하고, 실제 병목이 어디인지 판단하라.

### 연습문제 5: 벤치마크 프로파일링
맵(map)에 10만 개의 키-값 쌍을 삽입하는 벤치마크를 작성하고, `-cpuprofile`과 `-memprofile`을 함께 생성하여 분석하라. 미리 크기를 지정한 맵(`make(map[string]int, 100000)`)과 지정하지 않은 맵의 차이를 비교하라.

### 연습문제 6: 플레임 그래프 해석
재귀 피보나치 함수와 반복문 피보나치 함수의 CPU 프로파일을 각각 생성하고, 웹 UI의 플레임 그래프를 비교하라. 재귀 버전의 플레임 그래프 형태가 특징적인 이유를 설명하라.

### 연습문제 7: 메모리 누수 탐지
의도적으로 메모리 누수를 발생시키는 프로그램을 작성하라 (예: 맵에 계속 데이터를 추가하되 삭제하지 않는 고루틴). `net/http/pprof`의 힙 프로파일을 시간 간격으로 여러 번 수집하여 메모리 증가 패턴을 확인하라.

### 연습문제 8: 트레이스 분석
여러 고루틴이 채널을 통해 데이터를 주고받는 파이프라인을 작성하고, `go test -trace`로 트레이스를 생성하라. `go tool trace`에서 고루틴 스케줄링 패턴과 채널 대기 시간을 분석하라.

### 연습문제 9: 프로파일 비교
문자열을 처리하는 함수의 두 가지 구현(바이트 슬라이스 기반 vs strings.Builder 기반)을 작성하고, `pprof -diff_base`를 사용하여 CPU 및 메모리 프로파일을 비교하라.

### 연습문제 10: GC 영향 분석
대량의 작은 객체를 생성하는 코드와 소수의 큰 객체를 생성하는 코드를 각각 작성하라. 메모리 프로파일의 `allocs` 항목을 비교하고, GC 부하 차이를 설명하라.

---

## 구현 과제

### 과제 1: 성능 병목 탐지 및 최적화
의도적으로 비효율적인 코드(비효율적인 정렬, 불필요한 메모리 할당, 과도한 문자열 연결 등)를 포함하는 프로그램을 작성하라. CPU와 메모리 프로파일링을 수행하여 병목 지점을 찾고, 단계별로 최적화하라. 각 최적화 단계 전후의 벤치마크 결과를 기록하라.

### 과제 2: HTTP 서버 성능 분석 도구
HTTP 서버에 대해 부하 테스트를 수행하고 자동으로 프로파일을 수집하는 스크립트를 작성하라:
- `go tool pprof`를 사용하여 CPU/메모리 프로파일을 수집한다.
- 고루틴 수와 힙 메모리를 주기적으로 기록한다.
- 결과를 요약하여 출력한다.

### 과제 3: 메모리 사용량 최적화
대량의 로그 데이터(100만 줄 이상)를 파싱하는 프로그램을 작성하라. 초기 구현 후 메모리 프로파일링을 수행하고, 다음 기법을 적용하여 메모리 사용량을 줄여라:
- 버퍼 크기 최적화
- 문자열 인턴(intern) 기법
- 구조체 필드 정렬로 패딩 최소화
- `sync.Pool` 활용

### 과제 4: 고루틴 누수 탐지기
`net/http/pprof`의 고루틴 프로파일을 주기적으로 수집하여 고루틴 누수를 감지하는 모니터링 도구를 만들어라:
- 일정 간격으로 고루틴 수를 기록한다.
- 고루틴 수가 지속적으로 증가하면 경고를 출력한다.
- 누수로 의심되는 고루틴의 스택 트레이스를 표시한다.

---

## 프로젝트 과제

### 프로젝트 1: 웹 애플리케이션 성능 대시보드
HTTP 서버의 실시간 성능 지표를 표시하는 대시보드를 만들어라:
- `runtime` 패키지를 사용하여 메모리 사용량, 고루틴 수, GC 통계 등을 수집한다.
- `/metrics` 엔드포인트에서 JSON 형태로 지표를 제공한다.
- 요청당 응답 시간 히스토그램을 구현한다.
- 프로파일링 엔드포인트를 별도 포트에서 안전하게 제공한다.
- HTML/JS 기반 간단한 대시보드 페이지를 제공한다.

### 프로젝트 2: 자동 성능 회귀 테스트 시스템
Git 커밋별로 벤치마크를 자동 실행하고 성능 변화를 추적하는 시스템을 만들어라:
- 지정된 벤치마크 함수를 실행하고 결과를 저장한다.
- 이전 커밋의 결과와 비교하여 성능 변화를 표시한다.
- 성능이 일정 비율 이상 악화되면 경고를 출력한다.
- 결과를 시간순으로 정리하여 성능 추이를 파악할 수 있게 한다.
