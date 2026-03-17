package main

import (
	"errors"
	"fmt"
	"time"
)

// =============================================
// Go의 생성자 패턴 예제
// NewXxx 함수와 Functional Options 패턴
// =============================================

// --- 기본 NewXxx 패턴 ---

// Server 구조체: 비공개 필드를 가진 구조체
type Server struct {
	host     string
	port     int
	timeout  time.Duration
	maxConns int
	tls      bool
}

// NewServer는 Server의 "생성자" 역할을 한다
// 기본값을 설정하고 유효성을 검증한다
func NewServer(host string, port int) (*Server, error) {
	// 유효성 검증
	if host == "" {
		return nil, errors.New("호스트는 비어있을 수 없습니다")
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("유효하지 않은 포트 번호: %d (1-65535)", port)
	}

	return &Server{
		host:     host,
		port:     port,
		timeout:  30 * time.Second, // 기본 타임아웃
		maxConns: 100,              // 기본 최대 연결 수
		tls:      false,            // 기본 TLS 비활성화
	}, nil
}

// String 메서드: Server 정보를 문자열로 반환
func (s *Server) String() string {
	protocol := "http"
	if s.tls {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d (타임아웃: %v, 최대연결: %d)",
		protocol, s.host, s.port, s.timeout, s.maxConns)
}

// --- Functional Options 패턴 ---

// Option은 Server 설정을 변경하는 함수 타입
type Option func(*Server)

// WithTimeout은 타임아웃을 설정하는 옵션을 반환한다
func WithTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.timeout = d
	}
}

// WithMaxConns는 최대 연결 수를 설정하는 옵션을 반환한다
func WithMaxConns(n int) Option {
	return func(s *Server) {
		s.maxConns = n
	}
}

// WithTLS는 TLS를 활성화하는 옵션을 반환한다
func WithTLS() Option {
	return func(s *Server) {
		s.tls = true
	}
}

// NewServerWithOptions는 Functional Options 패턴을 사용하는 생성자이다
func NewServerWithOptions(host string, port int, opts ...Option) (*Server, error) {
	if host == "" {
		return nil, errors.New("호스트는 비어있을 수 없습니다")
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("유효하지 않은 포트 번호: %d", port)
	}

	// 기본값으로 서버 생성
	server := &Server{
		host:     host,
		port:     port,
		timeout:  30 * time.Second,
		maxConns: 100,
		tls:      false,
	}

	// 옵션 적용
	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

// --- 제로값이 유효한 구조체 ---

// Counter는 제로값 그 자체로 사용 가능한 구조체이다
// 별도의 생성자가 필요 없습니다!
type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) Value() int {
	return c.count
}

// --- 싱글톤 유사 패턴: 패키지 수준 변수 ---

// defaultLogger는 패키지 내부에서 관리하는 기본 로거
type Logger struct {
	prefix string
	level  string
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
		level:  "INFO",
	}
}

func (l *Logger) Log(msg string) {
	fmt.Printf("[%s] %s: %s\n", l.level, l.prefix, msg)
}

// --- Builder 패턴 (메서드 체이닝) ---

// QueryBuilder는 SQL 쿼리를 단계적으로 구성한다
type QueryBuilder struct {
	table      string
	conditions []string
	orderBy    string
	limit      int
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table: table,
		limit: -1,
	}
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	return qb // 자기 자신을 반환하여 메서드 체이닝 가능
}

func (qb *QueryBuilder) OrderBy(field string) *QueryBuilder {
	qb.orderBy = field
	return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
	qb.limit = n
	return qb
}

func (qb *QueryBuilder) Build() string {
	query := fmt.Sprintf("SELECT * FROM %s", qb.table)

	for i, cond := range qb.conditions {
		if i == 0 {
			query += " WHERE " + cond
		} else {
			query += " AND " + cond
		}
	}

	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	return query
}

func main() {
	// --- 기본 NewXxx 패턴 ---
	fmt.Println("========== 생성자 패턴 ==========\n")

	fmt.Println("=== 기본 NewXxx 패턴 ===")
	server, err := NewServer("localhost", 8080)
	if err != nil {
		fmt.Printf("에러: %v\n", err)
		return
	}
	fmt.Printf("  서버: %s\n", server)

	// 유효성 검증 실패 예
	_, err = NewServer("", 8080)
	fmt.Printf("  빈 호스트 에러: %v\n", err)

	_, err = NewServer("localhost", 99999)
	fmt.Printf("  잘못된 포트 에러: %v\n", err)

	// --- Functional Options 패턴 ---
	fmt.Println("\n=== Functional Options 패턴 ===")

	// 기본 옵션만
	s1, _ := NewServerWithOptions("api.example.com", 443)
	fmt.Printf("  기본 옵션: %s\n", s1)

	// 커스텀 옵션 적용
	s2, _ := NewServerWithOptions("api.example.com", 443,
		WithTimeout(60*time.Second),
		WithMaxConns(500),
		WithTLS(),
	)
	fmt.Printf("  커스텀 옵션: %s\n", s2)

	// --- 제로값이 유효한 구조체 ---
	fmt.Println("\n=== 제로값이 유효한 구조체 (생성자 불필요) ===")
	var counter Counter // 별도의 New 함수 없이 바로 사용 가능
	counter.Increment()
	counter.Increment()
	counter.Increment()
	fmt.Printf("  카운터 값: %d\n", counter.Value())

	// --- Logger 예제 ---
	fmt.Println("\n=== NewLogger 패턴 ===")
	logger := NewLogger("APP")
	logger.Log("서버가 시작되었습니다")
	logger.Log("요청을 처리 중입니다")

	// --- Builder 패턴 (메서드 체이닝) ---
	fmt.Println("\n=== Builder 패턴 (메서드 체이닝) ===")
	query := NewQueryBuilder("users").
		Where("age > 18").
		Where("active = true").
		OrderBy("name").
		Limit(10).
		Build()
	fmt.Printf("  생성된 쿼리: %s\n", query)

	// --- 정리 ---
	fmt.Println("\n=== 생성자 패턴 정리 ===")
	fmt.Println("  1. NewXxx: 가장 기본적인 패턴, 유효성 검증과 기본값 설정")
	fmt.Println("  2. Functional Options: 선택적 설정이 많을 때 유연한 패턴")
	fmt.Println("  3. 제로값 활용: 제로값이 유효하면 생성자가 불필요")
	fmt.Println("  4. Builder: 복잡한 객체를 단계적으로 구성할 때")
}
