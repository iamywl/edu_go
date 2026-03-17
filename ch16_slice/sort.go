package main

import (
	"fmt"
	"sort"
)

// Student 구조체 - 정렬 예제용
type Student struct {
	Name  string
	Score int
	Age   int
}

func main() {
	// ============================================
	// 슬라이스 정렬 (sort 패키지)
	// ============================================

	// 1. 기본 타입 정렬
	fmt.Println("=== 정수 정렬 ===")
	nums := []int{5, 3, 8, 1, 9, 2, 7, 4, 6}
	fmt.Println("정렬 전:", nums)
	sort.Ints(nums)
	fmt.Println("정렬 후:", nums)

	fmt.Println("\n=== 문자열 정렬 ===")
	words := []string{"바나나", "사과", "체리", "딸기", "포도"}
	fmt.Println("정렬 전:", words)
	sort.Strings(words)
	fmt.Println("정렬 후:", words)

	fmt.Println("\n=== 실수 정렬 ===")
	floats := []float64{3.14, 1.41, 2.71, 0.57, 1.73}
	fmt.Println("정렬 전:", floats)
	sort.Float64s(floats)
	fmt.Println("정렬 후:", floats)

	// 2. 정렬 확인
	fmt.Println("\n=== 정렬 확인 ===")
	sorted := []int{1, 2, 3, 4, 5}
	unsorted := []int{5, 3, 1, 4, 2}
	fmt.Println("sorted 정렬됨?:", sort.IntsAreSorted(sorted))     // true
	fmt.Println("unsorted 정렬됨?:", sort.IntsAreSorted(unsorted)) // false

	// 3. 역순 정렬
	fmt.Println("\n=== 역순 정렬 ===")
	desc := []int{3, 1, 4, 1, 5, 9, 2, 6}
	sort.Sort(sort.Reverse(sort.IntSlice(desc)))
	fmt.Println("내림차순:", desc)

	// 4. 구조체 슬라이스 정렬 (sort.Slice 사용)
	fmt.Println("\n=== 구조체 정렬 (sort.Slice) ===")
	students := []Student{
		{"김철수", 85, 20},
		{"이영희", 92, 19},
		{"박민수", 78, 21},
		{"최지영", 95, 20},
		{"정현우", 88, 22},
	}

	// 점수 기준 오름차순 정렬
	sort.Slice(students, func(i, j int) bool {
		return students[i].Score < students[j].Score
	})
	fmt.Println("점수 오름차순:")
	for _, s := range students {
		fmt.Printf("  %s: %d점 (나이: %d)\n", s.Name, s.Score, s.Age)
	}

	// 점수 기준 내림차순 정렬
	sort.Slice(students, func(i, j int) bool {
		return students[i].Score > students[j].Score
	})
	fmt.Println("\n점수 내림차순:")
	for _, s := range students {
		fmt.Printf("  %s: %d점 (나이: %d)\n", s.Name, s.Score, s.Age)
	}

	// 이름 기준 정렬
	sort.Slice(students, func(i, j int) bool {
		return students[i].Name < students[j].Name
	})
	fmt.Println("\n이름순:")
	for _, s := range students {
		fmt.Printf("  %s: %d점 (나이: %d)\n", s.Name, s.Score, s.Age)
	}

	// 5. 안정 정렬 (sort.SliceStable)
	// 동일한 키의 요소가 원래 순서를 유지
	fmt.Println("\n=== 안정 정렬 (SliceStable) ===")
	items := []Student{
		{"김철수", 85, 20},
		{"이영희", 85, 19},
		{"박민수", 90, 21},
		{"최지영", 90, 20},
	}

	// 나이순으로 먼저 정렬
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Age < items[j].Age
	})
	// 그 다음 점수순으로 안정 정렬 -> 같은 점수 내에서 나이순 유지
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Score < items[j].Score
	})
	fmt.Println("점수순 (같은 점수면 나이순):")
	for _, s := range items {
		fmt.Printf("  %s: %d점 (나이: %d)\n", s.Name, s.Score, s.Age)
	}

	// 6. 이진 탐색 (sort.SearchInts)
	fmt.Println("\n=== 이진 탐색 ===")
	sorted2 := []int{10, 20, 30, 40, 50, 60, 70, 80, 90}
	target := 50
	idx := sort.SearchInts(sorted2, target)
	fmt.Printf("정렬된 배열에서 %d의 위치: 인덱스 %d\n", target, idx)

	// 값이 없는 경우 - 삽입될 위치를 반환
	target2 := 35
	idx2 := sort.SearchInts(sorted2, target2)
	fmt.Printf("%d가 삽입될 위치: 인덱스 %d\n", target2, idx2)
}
