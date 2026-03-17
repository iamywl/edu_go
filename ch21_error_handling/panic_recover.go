package main

import "fmt"

// ============================================================
// 21.3 패닉과 복구 (panic & recover)
// panic은 프로그램을 즉시 중단시키고,
// recover는 defer 내에서 panic을 잡아 프로그램 종료를 방지한다.
// ============================================================

// mustPositive - 양수가 아니면 패닉을 발생시킵니다.
// "Must" 접두사는 에러 대신 패닉을 발생시키는 함수의 관례이다.
func mustPositive(n int) int {
	if n <= 0 {
		panic(fmt.Sprintf("양수가 필요하지만 %d를 받았습니다", n))
	}
	return n
}

// safeCall - recover를 사용하여 패닉으로부터 안전하게 함수를 실행한다.
func safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("패닉 복구 완료:", r)
		}
	}()
	fn()
}

// safeIndex - 슬라이스의 안전한 인덱스 접근
func safeIndex(s []int, idx int) (result int, err error) {
	// defer + recover로 패닉을 에러로 변환한다.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("인덱스 접근 실패: %v", r)
		}
	}()

	return s[idx], nil
}

// demonstratePanicFlow - 패닉 발생 시의 실행 흐름을 보여줍니다.
func demonstratePanicFlow() {
	defer fmt.Println("3. defer: 이 함수의 defer는 실행됩니다")

	fmt.Println("1. 패닉 발생 전")
	panic("무언가 잘못되었습니다!")
	// 아래 코드는 실행되지 않습니다. (unreachable)
	// fmt.Println("이 줄은 실행되지 않습니다")
}

func main() {
	fmt.Println("=== panic 기본 동작 ===")

	// 정상적인 호출
	fmt.Println("mustPositive(5):", mustPositive(5))

	// 패닉을 safeCall로 감싸서 안전하게 처리
	safeCall(func() {
		mustPositive(-1) // 패닉 발생!
	})
	fmt.Println("패닉 이후에도 프로그램이 계속 실행됩니다")

	fmt.Println("\n=== 패닉 실행 흐름 ===")
	// defer + recover로 패닉을 잡습니다.
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("4. recover: 패닉을 복구했습니다:", r)
			}
		}()

		demonstratePanicFlow()
	}()
	fmt.Println("5. 메인 함수 계속 실행")

	fmt.Println("\n=== safeIndex: 패닉을 에러로 변환 ===")

	data := []int{10, 20, 30, 40, 50}

	// 정상적인 인덱스 접근
	val, err := safeIndex(data, 2)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Println("data[2] =", val)
	}

	// 범위를 벗어난 인덱스 접근 (패닉이 에러로 변환됨)
	val, err = safeIndex(data, 10)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Println("data[10] =", val)
	}

	fmt.Println("\n=== panic 사용이 적절한 경우 ===")

	// 1. 프로그램 초기화 시점의 필수 조건 검증
	// (예: 설정 파일 로드 실패, 필수 환경 변수 누락)
	initConfig := func(dbURL string) {
		if dbURL == "" {
			// 초기화 실패는 panic이 적절할 수 있습니다.
			panic("DB_URL 환경 변수가 설정되지 않았습니다")
		}
	}
	safeCall(func() {
		initConfig("")
	})

	// 2. 논리적으로 도달 불가능한 코드
	direction := "north"
	switch direction {
	case "north", "south", "east", "west":
		fmt.Println("유효한 방향:", direction)
	default:
		// 모든 케이스를 처리했으므로 여기 도달하면 버그입니다.
		panic(fmt.Sprintf("예상하지 못한 방향: %s", direction))
	}

	fmt.Println("\n프로그램이 정상적으로 종료되었습니다.")
}
