package main

import "fmt"

// ============================================
// 빈 인터페이스 (interface{} / any) 활용
// ============================================

// Pair는 키-값 쌍을 저장하는 구조체 (any 타입 사용)
type Pair struct {
	Key   string
	Value any // interface{}와 동일
}

// SimpleMap은 간단한 키-값 저장소
type SimpleMap struct {
	items []Pair
}

// Set은 키-값 쌍을 추가한다
func (m *SimpleMap) Set(key string, value any) {
	// 기존 키가 있으면 업데이트
	for i, item := range m.items {
		if item.Key == key {
			m.items[i].Value = value
			return
		}
	}
	// 없으면 새로 추가
	m.items = append(m.items, Pair{Key: key, Value: value})
}

// Get은 키에 해당하는 값을 반환한다
func (m *SimpleMap) Get(key string) (any, bool) {
	for _, item := range m.items {
		if item.Key == key {
			return item.Value, true
		}
	}
	return nil, false
}

// PrintAll은 모든 키-값 쌍을 출력한다
func (m *SimpleMap) PrintAll() {
	for _, item := range m.items {
		fmt.Printf("  %s: %v (%T)\n", item.Key, item.Value, item.Value)
	}
}

func main() {
	// 1. any (interface{}) 기본 사용
	fmt.Println("=== any 타입 기본 ===")

	var val any

	val = 42
	fmt.Printf("정수: %v (타입: %T)\n", val, val)

	val = "Go 언어"
	fmt.Printf("문자열: %v (타입: %T)\n", val, val)

	val = true
	fmt.Printf("불리언: %v (타입: %T)\n", val, val)

	val = []int{1, 2, 3}
	fmt.Printf("슬라이스: %v (타입: %T)\n", val, val)

	val = struct{ Name string }{"Kim"}
	fmt.Printf("구조체: %v (타입: %T)\n", val, val)

	// 2. any 슬라이스 - 여러 타입을 하나의 슬라이스에 저장
	fmt.Println("\n=== any 슬라이스 ===")
	mixed := []any{
		42,
		"안녕하세요",
		3.14,
		true,
		[]string{"Go", "Python", "Java"},
		map[string]int{"나이": 25},
	}

	for i, v := range mixed {
		fmt.Printf("  [%d] %v (%T)\n", i, v, v)
	}

	// 3. any를 매개변수로 받는 함수
	fmt.Println("\n=== any 매개변수 함수 ===")
	PrintWithType(42)
	PrintWithType("Hello")
	PrintWithType(3.14)
	PrintWithType([]int{1, 2, 3})

	// 4. 가변 인수와 any 조합
	fmt.Println("\n=== 가변 인수 + any ===")
	PrintAll("이름:", "김철수", "나이:", 25, "활성:", true)

	// 5. SimpleMap 활용 (any로 다양한 타입 저장)
	fmt.Println("\n=== SimpleMap (any 활용) ===")
	m := &SimpleMap{}
	m.Set("이름", "이영희")
	m.Set("나이", 30)
	m.Set("키", 165.5)
	m.Set("취미", []string{"독서", "등산"})
	m.Set("기혼", false)

	fmt.Println("저장된 데이터:")
	m.PrintAll()

	// 값 가져오기와 타입 단언
	fmt.Println("\n특정 값 가져오기:")
	if name, ok := m.Get("이름"); ok {
		// 타입 단언으로 실제 타입으로 변환
		if nameStr, ok := name.(string); ok {
			fmt.Println("  이름:", nameStr)
		}
	}

	if age, ok := m.Get("나이"); ok {
		if ageInt, ok := age.(int); ok {
			fmt.Println("  나이:", ageInt)
		}
	}

	// 6. fmt.Println이 any를 사용하는 원리
	fmt.Println("\n=== fmt.Println의 원리 ===")
	fmt.Println("fmt.Println은 다음과 같이 정의되어 있습니다:")
	fmt.Println("  func Println(a ...any) (n int, err error)")
	fmt.Println("-> any 가변 인수를 받으므로 어떤 타입이든 출력 가능!")

	// 7. any 사용 시 주의사항
	fmt.Println("\n=== any 사용 시 주의사항 ===")
	fmt.Println("1. any는 타입 안전성을 잃으므로 꼭 필요할 때만 사용")
	fmt.Println("2. 가능하면 구체적인 타입이나 인터페이스를 사용")
	fmt.Println("3. any에서 값을 꺼낼 때는 반드시 타입 단언/스위치 사용")
	fmt.Println("4. Go 1.18+에서는 제네릭이 any보다 나은 대안일 수 있음")
}

// PrintWithType은 any 타입의 값과 그 타입을 출력한다
func PrintWithType(val any) {
	fmt.Printf("  값: %-15v 타입: %T\n", val, val)
}

// PrintAll은 가변 인수를 모두 출력한다
// fmt.Println과 유사한 시그니처
func PrintAll(args ...any) {
	for _, arg := range args {
		fmt.Print(arg, " ")
	}
	fmt.Println()
}
