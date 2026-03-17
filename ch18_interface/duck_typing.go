package main

import "fmt"

// ============================================
// 덕 타이핑 (구조적 타이핑) 예제
// ============================================

// "오리처럼 걷고, 오리처럼 꽥꽥거리면, 그것은 오리다."
// Go에서는 인터페이스를 명시적으로 구현한다고 선언하지 않는다.
// 메서드만 맞으면 자동으로 인터페이스를 만족한다.

// Speaker 인터페이스 - 말할 수 있는 존재
type Speaker interface {
	Speak() string
}

// Mover 인터페이스 - 움직일 수 있는 존재
type Mover interface {
	Move() string
}

// SpeakingMover 인터페이스 - 말하고 움직일 수 있는 존재
type SpeakingMover interface {
	Speaker
	Mover
}

// --- 다양한 타입들 ---

// Dog 구조체
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return fmt.Sprintf("%s: 멍멍!", d.Name)
}

func (d Dog) Move() string {
	return fmt.Sprintf("%s이(가) 뛰어간다!", d.Name)
}

// Cat 구조체
type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return fmt.Sprintf("%s: 야옹~", c.Name)
}

func (c Cat) Move() string {
	return fmt.Sprintf("%s이(가) 살금살금 걸어간다.", c.Name)
}

// Robot 구조체 - 동물이 아니지만 같은 인터페이스를 구현
type Robot struct {
	Model string
}

func (r Robot) Speak() string {
	return fmt.Sprintf("[%s] 삐빕 삐빕!", r.Model)
}

func (r Robot) Move() string {
	return fmt.Sprintf("[%s] 위이잉~ 전진합니다.", r.Model)
}

// Parrot 구조체 - Speaker만 구현 (Mover는 구현하지 않음)
type Parrot struct {
	Name string
}

func (p Parrot) Speak() string {
	return fmt.Sprintf("%s: 안녕! 안녕!", p.Name)
}

// ============================================
// 인터페이스를 사용하는 함수들
// ============================================

// MakeSpeak은 Speaker 인터페이스를 구현한 모든 타입을 받을 수 있다
func MakeSpeak(s Speaker) {
	fmt.Println(" ", s.Speak())
}

// MakeMove는 Mover 인터페이스를 구현한 모든 타입을 받을 수 있다
func MakeMove(m Mover) {
	fmt.Println(" ", m.Move())
}

// Interact는 SpeakingMover 인터페이스를 구현한 타입만 받을 수 있다
func Interact(sm SpeakingMover) {
	fmt.Println(" ", sm.Speak())
	fmt.Println(" ", sm.Move())
}

func main() {
	// 1. 덕 타이핑 기본 예제
	fmt.Println("=== 덕 타이핑: Speaker 인터페이스 ===")
	fmt.Println("다양한 타입이 같은 인터페이스를 만족:")

	// Dog, Cat, Robot, Parrot 모두 Speak() 메서드가 있으므로 Speaker
	speakers := []Speaker{
		Dog{Name: "뽀삐"},
		Cat{Name: "나비"},
		Robot{Model: "RX-78"},
		Parrot{Name: "콩이"},
	}

	for _, s := range speakers {
		MakeSpeak(s)
	}

	// 2. 복합 인터페이스
	fmt.Println("\n=== 복합 인터페이스: SpeakingMover ===")
	fmt.Println("Speaker + Mover를 모두 구현한 타입만 사용 가능:")

	// Dog, Cat, Robot은 Speak()과 Move() 모두 구현
	movers := []SpeakingMover{
		Dog{Name: "뽀삐"},
		Cat{Name: "나비"},
		Robot{Model: "RX-78"},
		// Parrot은 Move()가 없어서 SpeakingMover를 만족하지 않음!
	}

	for _, sm := range movers {
		fmt.Printf("\n[%T]\n", sm)
		Interact(sm)
	}

	// 3. 인터페이스의 유연성 - 나중에 새 타입 추가
	fmt.Println("\n=== 인터페이스의 유연성 ===")
	fmt.Println("기존 코드를 수정하지 않고 새 타입을 추가할 수 있다:")

	// 새로운 타입을 만들어도 기존 함수(MakeSpeak)를 그대로 사용
	type Alien struct{}
	// 주의: 로컬 타입에는 메서드를 추가할 수 없으므로 위에서 정의된 타입을 사용

	// 다양한 타입을 하나의 인터페이스로 처리
	dog := Dog{Name: "뽀삐"}
	robot := Robot{Model: "T-800"}

	// Speaker 인터페이스를 통해 같은 방식으로 처리
	MakeSpeak(dog)   // Dog.Speak() 호출
	MakeSpeak(robot) // Robot.Speak() 호출

	// 4. 인터페이스 타입으로 변수 할당
	fmt.Println("\n=== 인터페이스 변수 동적 할당 ===")
	var speaker Speaker

	speaker = Dog{Name: "초코"}
	fmt.Printf("타입: %T, 값: %s\n", speaker, speaker.Speak())

	speaker = Cat{Name: "모모"}
	fmt.Printf("타입: %T, 값: %s\n", speaker, speaker.Speak())

	speaker = Robot{Model: "ASIMO"}
	fmt.Printf("타입: %T, 값: %s\n", speaker, speaker.Speak())
}
