package main

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
)

// ============================================================
// 24.6 제네릭을 사용한 표준 라이브러리 패키지
// Go 1.21+에서 제공하는 slices, maps, cmp 패키지를 활용한다.
// ============================================================

// Person - 예제용 구조체
type Person struct {
	Name string
	Age  int
}

func main() {
	fmt.Println("=== slices 패키지 ===")

	// 정렬
	nums := []int{5, 3, 8, 1, 9, 2, 7, 4, 6}
	slices.Sort(nums)
	fmt.Println("Sort:", nums)

	// 역순 정렬
	slices.Reverse(nums)
	fmt.Println("Reverse:", nums)

	// 다시 정렬
	slices.Sort(nums)

	// 포함 여부 확인
	fmt.Println("Contains(5):", slices.Contains(nums, 5))
	fmt.Println("Contains(10):", slices.Contains(nums, 10))

	// 인덱스 검색
	fmt.Println("Index(7):", slices.Index(nums, 7))
	fmt.Println("Index(10):", slices.Index(nums, 10)) // -1

	// 최솟값, 최댓값
	fmt.Println("Min:", slices.Min(nums))
	fmt.Println("Max:", slices.Max(nums))

	// 이진 탐색 (정렬된 슬라이스에서)
	idx, found := slices.BinarySearch(nums, 7)
	fmt.Printf("BinarySearch(7): 인덱스=%d, 찾음=%v\n", idx, found)

	// 정렬 여부 확인
	fmt.Println("IsSorted:", slices.IsSorted(nums))

	fmt.Println("\n=== slices.SortFunc (커스텀 정렬) ===")

	people := []Person{
		{"홍길동", 30},
		{"김철수", 25},
		{"이영희", 35},
		{"박민수", 28},
	}

	// 나이순 정렬
	slices.SortFunc(people, func(a, b Person) int {
		return cmp.Compare(a.Age, b.Age)
	})
	fmt.Println("나이순:")
	for _, p := range people {
		fmt.Printf("  %s (%d세)\n", p.Name, p.Age)
	}

	// 이름순 정렬
	slices.SortFunc(people, func(a, b Person) int {
		return cmp.Compare(a.Name, b.Name)
	})
	fmt.Println("이름순:")
	for _, p := range people {
		fmt.Printf("  %s (%d세)\n", p.Name, p.Age)
	}

	fmt.Println("\n=== slices.Compact (연속 중복 제거) ===")

	dupes := []int{1, 1, 2, 3, 3, 3, 4, 4, 5}
	unique := slices.Compact(dupes)
	fmt.Println("원본:", dupes)
	fmt.Println("중복 제거:", unique)

	fmt.Println("\n=== maps 패키지 ===")

	// 맵 생성
	scores := map[string]int{
		"Go":     95,
		"Python": 88,
		"Rust":   92,
		"Java":   85,
	}

	// 맵 복제
	scoresCopy := maps.Clone(scores)
	fmt.Println("원본:", scores)
	fmt.Println("복제:", scoresCopy)

	// 맵 동일성 비교
	fmt.Println("Equal:", maps.Equal(scores, scoresCopy))

	// 키 목록
	keys := slices.Sorted(maps.Keys(scores))
	fmt.Println("키 목록 (정렬):", keys)

	// 값 목록
	values := slices.Sorted(maps.Values(scores))
	fmt.Println("값 목록 (정렬):", values)

	// 맵 삭제 (조건부)
	maps.DeleteFunc(scoresCopy, func(k string, v int) bool {
		return v < 90 // 90점 미만 삭제
	})
	fmt.Println("90점 이상만:", scoresCopy)

	fmt.Println("\n=== cmp 패키지 ===")

	// cmp.Compare: -1, 0, 1 반환
	fmt.Println("Compare(1, 2):", cmp.Compare(1, 2))             // -1
	fmt.Println("Compare(2, 2):", cmp.Compare(2, 2))             //  0
	fmt.Println("Compare(3, 2):", cmp.Compare(3, 2))             //  1
	fmt.Println("Compare(\"a\", \"b\"):", cmp.Compare("a", "b")) // -1

	// cmp.Or: 첫 번째 0이 아닌 값 반환
	fmt.Println("Or(0, 0, 3, 4):", cmp.Or(0, 0, 3, 4))             // 3
	fmt.Println("Or(\"\", \"\", \"기본값\"):", cmp.Or("", "", "기본값")) // 기본값

	// cmp.Or은 기본값 패턴에 유용하다.
	getPort := func(port int) int {
		return cmp.Or(port, 8080) // port가 0이면 8080 사용
	}
	fmt.Println("getPort(3000):", getPort(3000)) // 3000
	fmt.Println("getPort(0):", getPort(0))       // 8080

	fmt.Println("\n=== 실전 예제: 성적 처리 ===")

	students := map[string][]int{
		"홍길동": {85, 92, 78},
		"김철수": {90, 88, 95},
		"이영희": {75, 80, 70},
	}

	// 각 학생의 평균 계산
	averages := make(map[string]float64)
	for name, grades := range students {
		sum := 0
		for _, g := range grades {
			sum += g
		}
		averages[name] = float64(sum) / float64(len(grades))
	}

	// 이름 순으로 정렬하여 출력
	sortedNames := slices.Sorted(maps.Keys(averages))
	for _, name := range sortedNames {
		fmt.Printf("  %s: 평균 %.1f점\n", name, averages[name])
	}

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
