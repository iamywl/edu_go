package main

import (
	"fmt"
	"time"
)

// ============================================
// defer 지연 실행
// ============================================

func main() {
	// 1. defer 기본 동작
	fmt.Println("=== defer 기본 동작 ===")
	fmt.Println("시작")
	defer fmt.Println("defer: 종료 (마지막에 실행됨)")
	fmt.Println("중간")
	fmt.Println("끝")
	// 출력 순서: 시작 -> 중간 -> 끝 -> defer: 종료

	// 2. 여러 defer의 실행 순서 (LIFO - 후입선출)
	fmt.Println("\n=== defer 실행 순서 (LIFO) ===")
	multiDefer()

	// 3. defer로 리소스 정리 시뮬레이션
	fmt.Println("\n=== defer로 리소스 정리 ===")
	processFile("data.txt")

	// 4. defer로 실행 시간 측정
	fmt.Println("\n=== defer로 실행 시간 측정 ===")
	heavyWork()

	// 5. defer와 인수 평가 시점
	fmt.Println("\n=== defer 인수 평가 시점 ===")
	deferArgEvaluation()

	// 6. defer와 반환값
	fmt.Println("\n=== defer와 반환값 ===")
	result := deferWithReturn()
	fmt.Println("반환값:", result)

	// 7. defer와 패닉 복구
	fmt.Println("\n=== defer와 패닉 복구 ===")
	safeDiv(10, 0)
	fmt.Println("프로그램 계속 실행 중...")

	// 8. 실용 예제: 데이터베이스 연결 시뮬레이션
	fmt.Println("\n=== 실용 예제: DB 시뮬레이션 ===")
	queryDB()
}

// multiDefer는 여러 defer가 LIFO 순서로 실행됨을 보여준다
func multiDefer() {
	for i := 1; i <= 5; i++ {
		defer fmt.Printf("  defer %d\n", i)
	}
	fmt.Println("  일반 코드 실행")
	// 출력: 일반 코드 -> defer 5 -> defer 4 -> ... -> defer 1
}

// processFile은 defer로 파일을 안전하게 닫는 패턴을 보여준다
func processFile(filename string) {
	fmt.Printf("  파일 '%s' 열기\n", filename)
	// 실제로는 os.Open(filename) 사용
	defer fmt.Printf("  파일 '%s' 닫기 (defer)\n", filename)

	fmt.Println("  파일 데이터 읽는 중...")
	fmt.Println("  파일 데이터 처리 중...")
	// 함수가 어떻게 종료되든 (정상, 에러, 패닉) defer가 실행됨
}

// heavyWork는 defer를 사용한 실행 시간 측정 패턴을 보여준다
func heavyWork() {
	defer elapsed("heavyWork")() // defer에서 함수를 반환하는 패턴

	// 무거운 작업 시뮬레이션
	fmt.Println("  무거운 작업 수행 중...")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("  작업 완료!")
}

// elapsed는 실행 시간 측정을 위한 헬퍼 함수
// 호출 시 시작 시간을 기록하고, 반환된 함수가 경과 시간을 출력
func elapsed(name string) func() {
	start := time.Now()
	fmt.Printf("  [%s] 시작\n", name)
	return func() {
		fmt.Printf("  [%s] 종료 (소요시간: %v)\n", name, time.Since(start))
	}
}

// deferArgEvaluation은 defer의 인수가 즉시 평가됨을 보여준다
func deferArgEvaluation() {
	x := 10
	defer fmt.Printf("  defer 시점의 x: %d (defer 등록 시 평가됨)\n", x)
	x = 20
	fmt.Printf("  현재 x: %d\n", x)
	x = 30
	fmt.Printf("  현재 x: %d\n", x)
	// defer는 x=10일 때의 값을 사용
}

// deferWithReturn은 defer가 반환값을 수정하는 예시
func deferWithReturn() (result int) {
	defer func() {
		result += 10 // 명명된 반환값을 defer에서 수정 가능
		fmt.Printf("  defer에서 result를 %d로 수정\n", result)
	}()

	result = 5
	fmt.Printf("  원래 result: %d\n", result)
	return result // defer 실행 후 result=15가 반환됨
}

// safeDiv는 defer + recover로 패닉을 안전하게 처리한다
func safeDiv(a, b int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  패닉 복구! 원인: %v\n", r)
		}
	}()

	fmt.Printf("  %d / %d = ", a, b)
	result := a / b // b가 0이면 패닉 발생!
	fmt.Println(result)
}

// queryDB는 데이터베이스 연결과 정리 패턴을 시뮬레이션한다
func queryDB() {
	// 데이터베이스 연결
	fmt.Println("  데이터베이스 연결...")
	defer fmt.Println("  데이터베이스 연결 해제 (defer)")

	// 트랜잭션 시작
	fmt.Println("  트랜잭션 시작...")
	defer fmt.Println("  트랜잭션 정리 (defer)")

	// 쿼리 실행
	fmt.Println("  SELECT * FROM users...")
	fmt.Println("  쿼리 결과 처리 중...")

	// 함수 종료 시 역순으로 정리:
	// 1. 트랜잭션 정리
	// 2. 데이터베이스 연결 해제
}
