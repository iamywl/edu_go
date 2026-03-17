// 26장: 테스트와 벤치마크하기 - 벤치마크 예제
// go test -bench . -benchmem 명령으로 실행한다.
package main

import (
	"fmt"
	"strings"
	"testing"
)

// === 기본 벤치마크 ===

// BenchmarkAdd - Add 함수의 벤치마크
func BenchmarkAdd(b *testing.B) {
	// b.N번 반복한다 (Go가 자동으로 적절한 횟수를 결정)
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

// BenchmarkSubtract - Subtract 함수의 벤치마크
func BenchmarkSubtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Subtract(100, 200)
	}
}

// BenchmarkMultiply - Multiply 함수의 벤치마크
func BenchmarkMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Multiply(100, 200)
	}
}

// BenchmarkDivide - Divide 함수의 벤치마크
func BenchmarkDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Divide(100, 200) //nolint:errcheck
	}
}

// === 피보나치 벤치마크 - 입력 크기별 비교 ===

func BenchmarkFibonacci10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(10)
	}
}

func BenchmarkFibonacci20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

func BenchmarkFibonacci40(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(40)
	}
}

// === 팩토리얼 벤치마크 ===

func BenchmarkFactorial10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Factorial(10) //nolint:errcheck
	}
}

func BenchmarkFactorial20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Factorial(20) //nolint:errcheck
	}
}

// === 문자열 연결 방식 비교 벤치마크 ===
// 세 가지 방식의 성능 차이를 확인한다

// BenchmarkStringConcat - + 연산자로 문자열 연결
func BenchmarkStringConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 100; j++ {
			s += "hello"
		}
	}
}

// BenchmarkStringSprintf - fmt.Sprintf로 문자열 연결
func BenchmarkStringSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 100; j++ {
			s = fmt.Sprintf("%s%s", s, "hello")
		}
	}
}

// BenchmarkStringBuilder - strings.Builder로 문자열 연결 (가장 빠름)
func BenchmarkStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		for j := 0; j < 100; j++ {
			builder.WriteString("hello")
		}
		_ = builder.String()
	}
}

// === ResetTimer 사용 예제 ===

// BenchmarkWithResetTimer - 초기화 시간을 벤치마크에서 제외
func BenchmarkWithResetTimer(b *testing.B) {
	// 초기화 코드: 큰 슬라이스를 준비한다
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	// 타이머를 리셋하여 초기화 시간을 제외한다
	b.ResetTimer()

	// 실제 벤치마크 대상 코드
	for i := 0; i < b.N; i++ {
		sum := 0
		for _, v := range data {
			sum += v
		}
	}
}

// === 서브 벤치마크 ===

// BenchmarkPower - 다양한 입력에 대한 서브 벤치마크
func BenchmarkPower(b *testing.B) {
	// 서브 벤치마크: 다양한 지수 크기를 테스트한다
	benchmarks := []struct {
		name string
		base float64
		exp  float64
	}{
		{"2^10", 2, 10},
		{"2^100", 2, 100},
		{"2^1000", 2, 1000},
		{"10^10", 10, 10},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Power(bm.base, bm.exp)
			}
		})
	}
}
