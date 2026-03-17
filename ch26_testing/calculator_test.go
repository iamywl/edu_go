// 26장: 테스트와 벤치마크하기 - 테스트 예제
// go test -v 명령으로 실행한다.
package main

import (
	"testing"
)

// === 기본 테스트 ===

// TestAdd - Add 함수의 기본 테스트
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; 기대값 5", result)
	}
}

// TestSubtract - Subtract 함수의 기본 테스트
func TestSubtract(t *testing.T) {
	result := Subtract(10, 4)
	if result != 6 {
		t.Errorf("Subtract(10, 4) = %d; 기대값 6", result)
	}
}

// === 테이블 주도 테스트 (Table-Driven Tests) ===

// TestAddTableDriven - 테이블 주도 방식의 Add 테스트
func TestAddTableDriven(t *testing.T) {
	// 테스트 케이스를 구조체 슬라이스로 정의한다
	tests := []struct {
		name string // 테스트 케이스 이름
		a, b int    // 입력값
		want int    // 기대 결과
	}{
		{"양수 + 양수", 2, 3, 5},
		{"음수 + 음수", -1, -2, -3},
		{"양수 + 음수", 5, -3, 2},
		{"영 + 양수", 0, 5, 5},
		{"영 + 영", 0, 0, 0},
		{"큰 수", 1000000, 2000000, 3000000},
	}

	// 각 테스트 케이스를 서브테스트로 실행한다
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d, 기대값 %d",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestMultiplyTableDriven - 테이블 주도 방식의 Multiply 테스트
func TestMultiplyTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"양수끼리", 3, 7, 21},
		{"음수끼리", -3, -7, 21},
		{"양수와 음수", 3, -7, -21},
		{"영 곱하기", 0, 100, 0},
		{"1 곱하기", 1, 42, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Multiply(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Multiply(%d, %d) = %d, 기대값 %d",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestDivide - Divide 함수의 정상 동작과 에러 테스트
func TestDivide(t *testing.T) {
	// 정상 동작 테스트
	t.Run("정상 나누기", func(t *testing.T) {
		tests := []struct {
			name string
			a, b float64
			want float64
		}{
			{"10 / 2", 10, 2, 5},
			{"7 / 2", 7, 2, 3.5},
			{"1 / 3", 1, 3, 0.3333333333333333},
			{"-10 / 2", -10, 2, -5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Divide(tt.a, tt.b)
				if err != nil {
					t.Fatalf("예상치 못한 에러: %v", err)
				}
				if got != tt.want {
					t.Errorf("Divide(%.1f, %.1f) = %.10f, 기대값 %.10f",
						tt.a, tt.b, got, tt.want)
				}
			})
		}
	})

	// 0으로 나누기 에러 테스트
	t.Run("0으로 나누기", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Error("0으로 나눌 때 에러가 발생해야 합니다")
		}
	})
}

// TestAbs - 절대값 함수 테스트
func TestAbs(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-1, 1},
	}

	for _, tt := range tests {
		got := Abs(tt.input)
		if got != tt.want {
			t.Errorf("Abs(%d) = %d, 기대값 %d", tt.input, got, tt.want)
		}
	}
}

// TestFactorial - 팩토리얼 함수 테스트
func TestFactorial(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    int
		wantErr bool
	}{
		{"0의 팩토리얼", 0, 1, false},
		{"1의 팩토리얼", 1, 1, false},
		{"5의 팩토리얼", 5, 120, false},
		{"10의 팩토리얼", 10, 3628800, false},
		{"음수 에러", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Factorial(tt.input)
			// 에러 발생 여부 확인
			if (err != nil) != tt.wantErr {
				t.Fatalf("Factorial(%d) 에러 = %v, 에러 기대 = %v",
					tt.input, err, tt.wantErr)
			}
			// 에러가 없을 때만 결과 확인
			if !tt.wantErr && got != tt.want {
				t.Errorf("Factorial(%d) = %d, 기대값 %d",
					tt.input, got, tt.want)
			}
		})
	}
}

// TestFibonacci - 피보나치 함수 테스트
func TestFibonacci(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{10, 55},
		{20, 6765},
	}

	for _, tt := range tests {
		got := Fibonacci(tt.n)
		if got != tt.want {
			t.Errorf("Fibonacci(%d) = %d, 기대값 %d", tt.n, got, tt.want)
		}
	}
}

// === 헬퍼 함수 사용 예제 ===

// assertEqual - 테스트 헬퍼 함수
// t.Helper()를 호출하면 에러 발생 시 이 함수가 아닌 호출한 곳의 위치가 표시됩니다
func assertEqual(t *testing.T, got, want int) {
	t.Helper() // 헬퍼 함수로 표시
	if got != want {
		t.Errorf("결과 = %d, 기대값 %d", got, want)
	}
}

func TestWithHelper(t *testing.T) {
	assertEqual(t, Add(1, 2), 3)
	assertEqual(t, Add(-1, 1), 0)
	assertEqual(t, Multiply(3, 4), 12)
}
