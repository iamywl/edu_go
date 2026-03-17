package main

import (
	"fmt"
	"strings"
	"time"
)

// =============================================
// 의존성 역전 원칙 (Dependency Inversion Principle)
// 인터페이스를 활용한 느슨한 결합
// =============================================

// --- 도메인 모델 ---

// User는 사용자 정보를 나타냅니다
type User struct {
	ID    int
	Name  string
	Email string
}

// --- 인터페이스 정의 (추상) ---
// 인터페이스는 사용하는 쪽에서 정의하는 것이 Go의 관례이다

// UserRepository는 사용자 데이터 저장소의 추상이다
type UserRepository interface {
	FindByID(id int) (*User, error)
	Save(user *User) error
	FindAll() ([]*User, error)
}

// Notifier는 알림 발송의 추상이다
type Notifier interface {
	Notify(to, message string) error
}

// Logger는 로깅의 추상이다
type Logger interface {
	Log(message string)
}

// --- 구현체 1: 메모리 기반 저장소 (테스트용) ---

// MemoryRepository는 메모리에 데이터를 저장한다
type MemoryRepository struct {
	users  map[int]*User
	nextID int
}

// NewMemoryRepository는 메모리 저장소를 생성한다
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (r *MemoryRepository) FindByID(id int) (*User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("사용자를 찾을 수 없습니다: ID=%d", id)
	}
	return user, nil
}

func (r *MemoryRepository) Save(user *User) error {
	if user.ID == 0 {
		user.ID = r.nextID
		r.nextID++
	}
	r.users[user.ID] = user
	return nil
}

func (r *MemoryRepository) FindAll() ([]*User, error) {
	users := make([]*User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}

// --- 구현체 2: 콘솔 알림 ---

// ConsoleNotifier는 콘솔에 알림을 출력한다
type ConsoleNotifier struct{}

func (n *ConsoleNotifier) Notify(to, message string) error {
	fmt.Printf("  [알림 → %s] %s\n", to, message)
	return nil
}

// --- 구현체 3: 이메일 알림 (가상) ---

// EmailNotifier는 이메일로 알림을 보냅니다 (시뮬레이션)
type EmailNotifier struct {
	smtpServer string
}

func NewEmailNotifier(smtp string) *EmailNotifier {
	return &EmailNotifier{smtpServer: smtp}
}

func (n *EmailNotifier) Notify(to, message string) error {
	fmt.Printf("  [이메일 → %s via %s] %s\n", to, n.smtpServer, message)
	return nil
}

// --- 구현체 4: 콘솔 로거 ---

// ConsoleLogger는 콘솔에 로그를 출력한다
type ConsoleLogger struct {
	prefix string
}

func NewConsoleLogger(prefix string) *ConsoleLogger {
	return &ConsoleLogger{prefix: prefix}
}

func (l *ConsoleLogger) Log(message string) {
	now := time.Now().Format("15:04:05")
	fmt.Printf("  [%s %s] %s\n", now, l.prefix, message)
}

// --- 서비스: 인터페이스에 의존 ---

// UserService는 사용자 관련 비즈니스 로직을 담당한다
// 구체 타입이 아닌 인터페이스에 의존한다!
type UserService struct {
	repo     UserRepository // 인터페이스에 의존
	notifier Notifier       // 인터페이스에 의존
	logger   Logger         // 인터페이스에 의존
}

// NewUserService는 의존성을 주입받는 생성자이다
func NewUserService(repo UserRepository, notifier Notifier, logger Logger) *UserService {
	return &UserService{
		repo:     repo,
		notifier: notifier,
		logger:   logger,
	}
}

// RegisterUser는 새 사용자를 등록한다
func (s *UserService) RegisterUser(name, email string) (*User, error) {
	s.logger.Log(fmt.Sprintf("사용자 등록 시작: %s", name))

	user := &User{Name: name, Email: email}

	// 저장소에 저장 (어떤 저장소인지 모름 — 추상에 의존)
	if err := s.repo.Save(user); err != nil {
		return nil, fmt.Errorf("사용자 저장 실패: %w", err)
	}

	// 알림 발송 (어떤 알림인지 모름 — 추상에 의존)
	s.notifier.Notify(email, "가입을 환영합니다, "+name+"님!")

	s.logger.Log(fmt.Sprintf("사용자 등록 완료: ID=%d, Name=%s", user.ID, user.Name))
	return user, nil
}

// GetUser는 사용자를 조회한다
func (s *UserService) GetUser(id int) (*User, error) {
	return s.repo.FindByID(id)
}

// ListUsers는 모든 사용자를 조회한다
func (s *UserService) ListUsers() ([]*User, error) {
	return s.repo.FindAll()
}

// --- Mock 구현: 테스트에 유용 ---

// MockNotifier는 테스트용 알림 구현체이다
type MockNotifier struct {
	Calls []struct {
		To      string
		Message string
	}
}

func (m *MockNotifier) Notify(to, message string) error {
	m.Calls = append(m.Calls, struct {
		To      string
		Message string
	}{to, message})
	return nil
}

func main() {
	fmt.Println("========== 의존성 역전 원칙 ==========\n")

	// --- 시나리오 1: 개발 환경 (메모리 저장소 + 콘솔 알림) ---
	fmt.Println("=== 시나리오 1: 개발 환경 ===")

	memRepo := NewMemoryRepository()
	consoleNotifier := &ConsoleNotifier{}
	logger := NewConsoleLogger("DEV")

	// 의존성 주입: 구체 타입을 인터페이스로 전달
	devService := NewUserService(memRepo, consoleNotifier, logger)

	devService.RegisterUser("김고퍼", "gopher@go.dev")
	devService.RegisterUser("이러스트", "rust@example.com")

	users, _ := devService.ListUsers()
	fmt.Printf("\n  등록된 사용자 수: %d\n", len(users))
	for _, u := range users {
		fmt.Printf("    - ID=%d, Name=%s, Email=%s\n", u.ID, u.Name, u.Email)
	}

	// --- 시나리오 2: 프로덕션 환경 (같은 서비스, 다른 구현체) ---
	fmt.Println("\n=== 시나리오 2: 프로덕션 환경 (이메일 알림으로 교체) ===")

	emailNotifier := NewEmailNotifier("smtp.example.com")
	prodLogger := NewConsoleLogger("PROD")

	// 같은 UserService이지만 다른 구현체를 주입
	prodService := NewUserService(memRepo, emailNotifier, prodLogger)
	prodService.RegisterUser("박파이썬", "python@example.com")

	// --- 시나리오 3: 테스트 환경 (Mock 사용) ---
	fmt.Println("\n=== 시나리오 3: 테스트 (Mock으로 검증) ===")

	testRepo := NewMemoryRepository()
	mockNotifier := &MockNotifier{}
	testLogger := NewConsoleLogger("TEST")

	testService := NewUserService(testRepo, mockNotifier, testLogger)
	testService.RegisterUser("테스트유저", "test@test.com")

	// Mock을 통해 알림이 정상적으로 호출되었는지 검증
	if len(mockNotifier.Calls) == 1 {
		call := mockNotifier.Calls[0]
		fmt.Printf("  [검증 성공] 알림 발송됨: To=%s, Message=%s\n",
			call.To, call.Message)
	}

	// --- 핵심 정리 ---
	fmt.Println("\n=== 의존성 역전 원칙 핵심 정리 ===")
	printSummary()
}

func printSummary() {
	summary := []string{
		"1. 구체 타입이 아닌 인터페이스에 의존하라",
		"2. 인터페이스는 사용하는 쪽에서 정의하라 (Go 관례)",
		"3. 생성자를 통해 의존성을 주입하라 (DI)",
		"4. 작은 인터페이스를 선호하라 (1~3개 메서드)",
		"5. 테스트 시 Mock 구현체를 주입하여 격리된 테스트 가능",
	}

	for _, s := range summary {
		fmt.Println("  " + s)
	}

	fmt.Println()
	fmt.Println("  나쁜 예:")
	fmt.Println("    type Service struct { db *MySQLDatabase } // 구체 타입에 의존!")
	fmt.Println()
	fmt.Println("  좋은 예:")
	fmt.Println("    type Service struct { repo UserRepository } // 인터페이스에 의존!")

	_ = strings.Join(nil, "")
}
