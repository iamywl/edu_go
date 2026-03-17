package main

import "fmt"

// Config 구조체 — 설정 정보
type Config struct {
	Host  string
	Port  int
	Debug bool
}

func main() {
	// === new() 함수로 인스턴스 생성 ===
	fmt.Println("=== new() 함수 ===")

	// new()는 제로값으로 초기화된 인스턴스의 포인터를 반환
	cfg1 := new(Config)
	fmt.Printf("타입: %T\n", cfg1)    // *main.Config
	fmt.Printf("값:   %+v\n", *cfg1) // 제로값: {Host: Port:0 Debug:false}

	// 필드를 나중에 설정
	cfg1.Host = "localhost"
	cfg1.Port = 8080
	cfg1.Debug = true
	fmt.Printf("설정 후: %+v\n", *cfg1)

	fmt.Println()

	// === &struct{} 방식으로 인스턴스 생성 (더 자주 사용) ===
	fmt.Println("=== &struct{} 방식 ===")

	// 선언과 동시에 초기값 지정 가능!
	cfg2 := &Config{
		Host:  "127.0.0.1",
		Port:  3000,
		Debug: false,
	}
	fmt.Printf("타입: %T\n", cfg2) // *main.Config
	fmt.Printf("값:   %+v\n", *cfg2)

	fmt.Println()

	// === 팩토리 함수 패턴 ===
	fmt.Println("=== 팩토리 함수 패턴 ===")

	cfg3 := NewConfig("api.example.com", 443)
	fmt.Printf("기본 설정: %+v\n", *cfg3)

	cfg4 := NewDebugConfig("localhost", 8080)
	fmt.Printf("디버그 설정: %+v\n", *cfg4)

	fmt.Println()

	// === 같은 인스턴스를 여러 포인터가 가리킬 수 있다 ===
	fmt.Println("=== 여러 포인터가 같은 인스턴스 참조 ===")

	original := &Config{Host: "original.com", Port: 80, Debug: false}
	alias := original // 같은 인스턴스를 가리킴 (복사가 아님!)

	fmt.Printf("original: %+v\n", *original)
	fmt.Printf("alias:    %+v\n", *alias)

	// alias를 통해 수정하면 original도 변경됨
	alias.Port = 9999
	fmt.Printf("\nalias.Port = 9999 후:\n")
	fmt.Printf("original.Port: %d\n", original.Port) // 9999
	fmt.Printf("alias.Port:    %d\n", alias.Port)    // 9999
	fmt.Printf("같은 주소?     %v\n", original == alias)
}

// NewConfig 는 Config 인스턴스를 생성하는 팩토리 함수
func NewConfig(host string, port int) *Config {
	return &Config{
		Host:  host,
		Port:  port,
		Debug: false,
	}
}

// NewDebugConfig 는 디버그 모드가 활성화된 Config를 생성하는 팩토리 함수
func NewDebugConfig(host string, port int) *Config {
	return &Config{
		Host:  host,
		Port:  port,
		Debug: true,
	}
}
