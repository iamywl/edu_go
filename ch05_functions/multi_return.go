// multi_return.go
// 멀티 반환(Multiple Return Values) 예제
// 실행 방법: go run multi_return.go

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// divide: 나눗셈 결과와 에러를 동시에 반환
// Go에서 가장 흔한 멀티 반환 패턴: (결과값, error)
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("0으로 나눌 수 없습니다")
	}
	return a / b, nil // nil = 에러 없음
}

// minMax: 슬라이스에서 최솟값과 최댓값을 반환
func minMax(numbers []int) (min, max int) {
	if len(numbers) == 0 {
		return 0, 0
	}

	min = numbers[0]
	max = numbers[0]

	for _, n := range numbers[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return // 이름이 있는 반환값이므로 naked return 가능
}

// sumAndAvg: 합계와 평균을 동시에 반환
func sumAndAvg(numbers ...int) (sum int, avg float64) {
	if len(numbers) == 0 {
		return 0, 0
	}

	for _, n := range numbers {
		sum += n
	}
	avg = float64(sum) / float64(len(numbers))
	return
}

// circleInfo: 원의 넓이와 둘레를 반환
func circleInfo(radius float64) (area, circumference float64) {
	area = math.Pi * radius * radius
	circumference = 2 * math.Pi * radius
	return
}

// parseName: 전체 이름을 성과 이름으로 분리
func parseName(fullName string) (firstName, lastName string, err error) {
	parts := strings.Fields(fullName)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("성과 이름을 모두 입력해주세요: %q", fullName)
	}
	return parts[0], parts[1], nil
}

func main() {
	fmt.Println("===== 에러 처리 패턴 (가장 흔한 패턴) =====")

	// 성공 케이스
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("10 / 3 = %.4f\n", result)
	}

	// 실패 케이스
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("10 / 0 = %.4f\n", result)
	}

	fmt.Println()
	fmt.Println("===== 반환값 무시 (_) =====")

	// 필요 없는 반환값은 _ (빈 식별자)로 무시
	onlyResult, _ := divide(20, 4) // 에러를 무시 (확실한 경우에만!)
	fmt.Printf("20 / 4 = %.1f (에러 무시)\n", onlyResult)

	fmt.Println()
	fmt.Println("===== 최솟값/최댓값 =====")

	nums := []int{38, 12, 95, 7, 63, 41}
	min, max := minMax(nums)
	fmt.Printf("숫자: %v\n", nums)
	fmt.Printf("최솟값: %d, 최댓값: %d\n", min, max)

	fmt.Println()
	fmt.Println("===== 합계와 평균 =====")

	total, average := sumAndAvg(80, 90, 75, 88, 92)
	fmt.Printf("합계: %d, 평균: %.1f\n", total, average)

	fmt.Println()
	fmt.Println("===== 원의 정보 =====")

	area, circumference := circleInfo(5.0)
	fmt.Printf("반지름 5인 원 → 넓이: %.2f, 둘레: %.2f\n", area, circumference)

	fmt.Println()
	fmt.Println("===== 이름 파싱 =====")

	first, last, err := parseName("홍 길동")
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("성: %s, 이름: %s\n", first, last)
	}

	_, _, err = parseName("홍길동") // 공백이 없는 경우
	if err != nil {
		fmt.Println("에러:", err)
	}

	fmt.Println()
	fmt.Println("===== 표준 라이브러리의 멀티 반환 예제 =====")

	// strconv.Atoi: 문자열 → 정수 변환
	num, err := strconv.Atoi("123")
	if err != nil {
		fmt.Println("변환 에러:", err)
	} else {
		fmt.Printf("\"123\" → %d\n", num)
	}

	num, err = strconv.Atoi("abc")
	if err != nil {
		fmt.Println("변환 에러:", err)
	}

	// strings.Cut: 구분자로 문자열 분리
	before, after, found := strings.Cut("name=홍길동", "=")
	if found {
		fmt.Printf("키: %s, 값: %s\n", before, after)
	}

	fmt.Println()
	fmt.Println("===== 멀티 반환으로 swap 구현 =====")

	x, y := 10, 20
	fmt.Printf("교환 전: x=%d, y=%d\n", x, y)

	// Go에서는 멀티 반환을 이용해 간단하게 값 교환 가능
	x, y = y, x
	fmt.Printf("교환 후: x=%d, y=%d\n", x, y)
}
