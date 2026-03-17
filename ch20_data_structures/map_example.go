package main

import (
	"fmt"
	"sort"
	"strings"
)

// ============================================
// 맵 (map) - 해시 테이블 기반 키-값 저장소
// ============================================

// Student 구조체
type Student struct {
	Name  string
	Score int
	Grade string
}

func main() {
	// 1. 맵 생성과 초기화
	fmt.Println("=== 맵 생성과 초기화 ===")

	// 리터럴로 생성
	fruits := map[string]int{
		"사과":  1000,
		"바나나": 2000,
		"딸기":  3000,
	}
	fmt.Println("과일 가격:", fruits)

	// make()로 생성
	ages := make(map[string]int)
	ages["김철수"] = 25
	ages["이영희"] = 30
	fmt.Println("나이:", ages)

	// 2. CRUD 연산
	fmt.Println("\n=== CRUD 연산 ===")

	// Create (생성)
	scores := make(map[string]int)
	scores["김철수"] = 85
	scores["이영희"] = 92
	scores["박민수"] = 78
	fmt.Println("생성:", scores)

	// Read (읽기)
	fmt.Println("김철수 점수:", scores["김철수"])

	// 존재하지 않는 키 읽기 -> 제로값 반환
	fmt.Println("최지영 점수:", scores["최지영"]) // 0

	// 존재 여부 확인 (comma ok 패턴)
	if val, ok := scores["이영희"]; ok {
		fmt.Println("이영희 존재, 점수:", val)
	}
	if _, ok := scores["최지영"]; !ok {
		fmt.Println("최지영은 존재하지 않음")
	}

	// Update (수정)
	scores["김철수"] = 90
	fmt.Println("수정 후:", scores)

	// Delete (삭제)
	delete(scores, "박민수")
	fmt.Println("삭제 후:", scores)
	fmt.Println("맵 크기:", len(scores))

	// 3. 맵 순회
	fmt.Println("\n=== 맵 순회 ===")
	colors := map[string]string{
		"빨강": "#FF0000",
		"초록": "#00FF00",
		"파랑": "#0000FF",
		"노랑": "#FFFF00",
	}

	// 키-값 순회 (순서 보장 안 됨!)
	for key, val := range colors {
		fmt.Printf("  %s: %s\n", key, val)
	}

	// 정렬된 순서로 순회하려면 키를 정렬
	fmt.Println("\n정렬된 순서로 순회:")
	keys := make([]string, 0, len(colors))
	for k := range colors {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %s\n", k, colors[k])
	}

	// 4. 값으로 구조체 사용
	fmt.Println("\n=== 구조체 값 맵 ===")
	students := map[int]Student{
		1: {Name: "김철수", Score: 85, Grade: "B"},
		2: {Name: "이영희", Score: 92, Grade: "A"},
		3: {Name: "박민수", Score: 78, Grade: "C"},
	}

	for id, s := range students {
		fmt.Printf("  [%d] %s: %d점 (%s)\n", id, s.Name, s.Score, s.Grade)
	}

	// 5. 값으로 슬라이스 사용 (그래프, 인접 리스트)
	fmt.Println("\n=== 맵의 값으로 슬라이스 (그래프) ===")
	graph := map[string][]string{
		"서울": {"부산", "대전", "인천"},
		"부산": {"서울", "대구", "울산"},
		"대전": {"서울", "대구"},
		"대구": {"부산", "대전"},
		"인천": {"서울"},
		"울산": {"부산"},
	}

	for city, neighbors := range graph {
		fmt.Printf("  %s -> %v\n", city, neighbors)
	}

	// 서울에서 직접 연결된 도시
	fmt.Println("\n서울에서 직접 연결:", graph["서울"])

	// 6. 단어 빈도수 세기
	fmt.Println("\n=== 단어 빈도수 ===")
	text := "go go go python java go python rust go java go"
	wordCount := countWords(text)

	// 빈도수 기준 정렬하여 출력
	type wordFreq struct {
		word  string
		count int
	}
	var freqs []wordFreq
	for w, c := range wordCount {
		freqs = append(freqs, wordFreq{w, c})
	}
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].count > freqs[j].count
	})
	for _, f := range freqs {
		fmt.Printf("  %s: %d회\n", f.word, f.count)
	}

	// 7. 맵을 집합(Set)처럼 사용
	fmt.Println("\n=== 맵을 집합(Set)으로 사용 ===")
	set := make(map[string]bool)

	// 추가
	set["Go"] = true
	set["Python"] = true
	set["Java"] = true

	// 존재 확인
	fmt.Println("Go 존재?:", set["Go"])     // true
	fmt.Println("Rust 존재?:", set["Rust"]) // false

	// 삭제
	delete(set, "Java")

	// 집합 순회
	fmt.Print("집합 요소: ")
	for item := range set {
		fmt.Print(item, " ")
	}
	fmt.Println()

	// 8. 맵은 참조 타입
	fmt.Println("\n=== 맵은 참조 타입 ===")
	original := map[string]int{"a": 1, "b": 2}
	copied := original // 참조 복사 (같은 맵을 가리킴)
	copied["c"] = 3    // 원본에도 영향

	fmt.Println("원본:", original) // map[a:1 b:2 c:3]
	fmt.Println("복사:", copied)   // map[a:1 b:2 c:3]

	// 독립적인 복사가 필요하면 수동 복사
	independent := make(map[string]int)
	for k, v := range original {
		independent[k] = v
	}
	independent["d"] = 4
	fmt.Println("독립 복사:", independent)
	fmt.Println("원본 (변경 없음):", original)

	// 9. 맵의 키 타입
	fmt.Println("\n=== 다양한 키 타입 ===")

	// 정수 키
	intMap := map[int]string{1: "일", 2: "이", 3: "삼"}
	fmt.Println("정수 키:", intMap)

	// 구조체 키
	type Point struct {
		X, Y int
	}
	pointMap := map[Point]string{
		{0, 0}: "원점",
		{1, 0}: "오른쪽",
		{0, 1}: "위쪽",
	}
	fmt.Println("구조체 키:", pointMap)
	fmt.Println("원점:", pointMap[Point{0, 0}])
}

// countWords는 텍스트에서 단어 빈도수를 세어 반환한다
func countWords(text string) map[string]int {
	counts := make(map[string]int)
	words := strings.Fields(text) // 공백으로 분리
	for _, word := range words {
		counts[word]++
	}
	return counts
}
