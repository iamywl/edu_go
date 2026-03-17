package main

import "fmt"

// =============================================
// Go에서의 OOP 패턴 예제
// 클래스 없이 구조체 + 메서드 + 인터페이스로 OOP 구현
// =============================================

// --- 인터페이스 정의: 다형성의 핵심 ---

// Speaker 인터페이스: Speak 메서드를 가진 모든 타입이 자동으로 구현
type Speaker interface {
	Speak() string
}

// Mover 인터페이스
type Mover interface {
	Move() string
}

// 인터페이스 조합: 작은 인터페이스를 합쳐서 큰 인터페이스 생성
type Animal interface {
	Speaker
	Mover
	Name() string
}

// --- 기본 구조체: "부모 클래스" 역할 ---

// BaseAnimal은 공통 필드와 메서드를 제공한다
type BaseAnimal struct {
	name string
	age  int
}

// Name 메서드: BaseAnimal을 임베드하는 모든 구조체에서 사용 가능
func (b BaseAnimal) Name() string {
	return b.name
}

// Info 메서드: 공통 정보 출력
func (b BaseAnimal) Info() string {
	return fmt.Sprintf("%s (나이: %d살)", b.name, b.age)
}

// --- Dog 구조체: BaseAnimal 임베딩 (조합) ---

type Dog struct {
	BaseAnimal // 임베딩: BaseAnimal의 필드와 메서드를 "승격"
	Breed      string
}

// Dog만의 Speak 구현
func (d Dog) Speak() string {
	return fmt.Sprintf("%s: 멍멍! 🐕", d.name)
}

// Dog만의 Move 구현
func (d Dog) Move() string {
	return fmt.Sprintf("%s이(가) 신나게 달립니다!", d.name)
}

// Dog에만 있는 메서드
func (d Dog) Fetch(item string) string {
	return fmt.Sprintf("%s이(가) %s을(를) 물어왔습니다!", d.name, item)
}

// --- Cat 구조체: BaseAnimal 임베딩 ---

type Cat struct {
	BaseAnimal
	Indoor bool
}

func (c Cat) Speak() string {
	return fmt.Sprintf("%s: 야옹~ 🐱", c.name)
}

func (c Cat) Move() string {
	if c.Indoor {
		return fmt.Sprintf("%s이(가) 집 안을 조용히 걸어다닙니다.", c.name)
	}
	return fmt.Sprintf("%s이(가) 담장 위를 걸어갑니다.", c.name)
}

// --- Robot 구조체: BaseAnimal 없이 독립적으로 인터페이스 구현 ---
// Go의 인터페이스는 암시적이므로, 메서드만 맞으면 어떤 구조체든 가능

type Robot struct {
	model string
}

func (r Robot) Name() string {
	return r.model
}

func (r Robot) Speak() string {
	return fmt.Sprintf("%s: 삐빅- 안녕하세요, 로봇입니다.", r.model)
}

func (r Robot) Move() string {
	return fmt.Sprintf("%s이(가) 바퀴로 이동합니다.", r.model)
}

// --- 다형성 활용 함수 ---

// introduceAnimal은 Animal 인터페이스를 만족하는 모든 타입을 받습니다
func introduceAnimal(a Animal) {
	fmt.Printf("이름: %s\n", a.Name())
	fmt.Printf("소리: %s\n", a.Speak())
	fmt.Printf("이동: %s\n", a.Move())
	fmt.Println()
}

// makeSpeech는 Speaker 인터페이스만 필요 (작은 인터페이스 활용)
func makeSpeech(speakers []Speaker) {
	fmt.Println("--- 모두 함께 말하기 ---")
	for _, s := range speakers {
		fmt.Printf("  %s\n", s.Speak())
	}
}

// --- 타입 단언(Type Assertion)과 타입 스위치 ---

func describeAnimal(a Animal) {
	fmt.Printf("[%s] ", a.Name())

	// 타입 스위치: 구체 타입에 따라 분기
	switch v := a.(type) {
	case Dog:
		fmt.Printf("견종: %s, %s\n", v.Breed, v.Fetch("공"))
	case Cat:
		if v.Indoor {
			fmt.Println("실내 고양이입니다.")
		} else {
			fmt.Println("실외 고양이입니다.")
		}
	case Robot:
		fmt.Printf("모델명: %s\n", v.model)
	default:
		fmt.Println("알 수 없는 동물입니다.")
	}
}

func main() {
	// --- 객체 생성 ---
	fmt.Println("========== Go OOP 패턴 ==========\n")

	dog := Dog{
		BaseAnimal: BaseAnimal{name: "바둑이", age: 3},
		Breed:      "진돗개",
	}

	cat := Cat{
		BaseAnimal: BaseAnimal{name: "나비", age: 5},
		Indoor:     true,
	}

	robot := Robot{model: "GoBot-3000"}

	// --- 임베딩으로 승격된 메서드 사용 ---
	fmt.Println("=== 임베딩으로 승격된 메서드 ===")
	fmt.Println("  " + dog.Info()) // BaseAnimal의 Info()를 직접 호출
	fmt.Println("  " + cat.Info())
	fmt.Println()

	// --- 다형성: 인터페이스를 통한 통합 처리 ---
	fmt.Println("=== 다형성: Animal 인터페이스 ===")
	animals := []Animal{dog, cat, robot}

	for _, a := range animals {
		introduceAnimal(a)
	}

	// --- 작은 인터페이스 활용 ---
	speakers := []Speaker{dog, cat, robot}
	makeSpeech(speakers)
	fmt.Println()

	// --- 타입 스위치 ---
	fmt.Println("=== 타입 스위치 ===")
	for _, a := range animals {
		describeAnimal(a)
	}

	// --- 인터페이스 비교: Go vs 전통 OOP ---
	fmt.Println("\n=== Go OOP 핵심 정리 ===")
	fmt.Println("  1. 클래스 대신 구조체(struct)를 사용합니다")
	fmt.Println("  2. 상속 대신 임베딩(embedding)으로 코드를 재사용합니다")
	fmt.Println("  3. 다형성은 인터페이스로 구현합니다")
	fmt.Println("  4. 인터페이스는 암시적으로 구현됩니다 (implements 키워드 없음)")
	fmt.Println("  5. 작은 인터페이스를 조합하여 큰 인터페이스를 만듭니다")
}
