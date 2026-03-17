package main

import (
	"fmt"
	"unsafe"
)

// 비효율적 구조체: 패딩이 많이 발생하는 필드 배치
type BadLayout struct {
	A bool    // 1바이트 + 패딩 7바이트 (다음 float64 정렬을 위해)
	B float64 // 8바이트
	C int32   // 4바이트 + 패딩 4바이트 (전체 구조체 정렬을 위해)
}

// 효율적 구조체: 패딩을 최소화한 필드 배치
type GoodLayout struct {
	B float64 // 8바이트
	C int32   // 4바이트
	A bool    // 1바이트 + 패딩 3바이트
}

// 다양한 타입의 크기를 확인하기 위한 구조체
type Example struct {
	A bool   // 1바이트
	B int8   // 1바이트
	C int16  // 2바이트
	D int32  // 4바이트
	E int64  // 8바이트
	F string // 16바이트 (포인터 8 + 길이 8)
}

func main() {
	// === 기본 타입 크기 확인 ===
	fmt.Println("=== 기본 타입 크기 ===")
	fmt.Printf("bool:    %d 바이트\n", unsafe.Sizeof(true))
	fmt.Printf("int8:    %d 바이트\n", unsafe.Sizeof(int8(0)))
	fmt.Printf("int16:   %d 바이트\n", unsafe.Sizeof(int16(0)))
	fmt.Printf("int32:   %d 바이트\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("int64:   %d 바이트\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("float64: %d 바이트\n", unsafe.Sizeof(float64(0)))
	fmt.Printf("string:  %d 바이트\n", unsafe.Sizeof(""))

	fmt.Println()

	// === 구조체 크기 비교 ===
	fmt.Println("=== 구조체 크기 비교 (패딩 효과) ===")

	bad := BadLayout{}
	good := GoodLayout{}

	fmt.Printf("BadLayout  크기: %d 바이트 (bool, float64, int32 순서)\n", unsafe.Sizeof(bad))
	fmt.Printf("GoodLayout 크기: %d 바이트 (float64, int32, bool 순서)\n", unsafe.Sizeof(good))
	fmt.Printf("절약된 메모리: %d 바이트\n", unsafe.Sizeof(bad)-unsafe.Sizeof(good))

	fmt.Println()

	// === 왜 패딩이 발생하는가? ===
	fmt.Println("=== 메모리 정렬 설명 ===")
	fmt.Println("CPU는 특정 크기의 데이터를 특정 주소 경계에서 읽어야 효율적입니다.")
	fmt.Println("예: 8바이트 float64는 8의 배수 주소에 위치해야 합니다.")
	fmt.Println()

	// 필드별 오프셋 확인
	fmt.Println("=== BadLayout 필드 오프셋 ===")
	fmt.Printf("A (bool)    오프셋: %d\n", unsafe.Offsetof(bad.A))
	fmt.Printf("B (float64) 오프셋: %d\n", unsafe.Offsetof(bad.B))
	fmt.Printf("C (int32)   오프셋: %d\n", unsafe.Offsetof(bad.C))

	fmt.Println()

	fmt.Println("=== GoodLayout 필드 오프셋 ===")
	fmt.Printf("B (float64) 오프셋: %d\n", unsafe.Offsetof(good.B))
	fmt.Printf("C (int32)   오프셋: %d\n", unsafe.Offsetof(good.C))
	fmt.Printf("A (bool)    오프셋: %d\n", unsafe.Offsetof(good.A))

	fmt.Println()

	// === 최적화 팁 ===
	fmt.Println("=== 최적화 팁 ===")
	fmt.Println("규칙: 큰 타입(8바이트) → 중간 타입(4바이트) → 작은 타입(1바이트) 순서로 배치")
	fmt.Println("대량의 구조체를 다룰 때 메모리 절약 효과가 큽니다.")

	// 예: 100만 개 생성 시 차이
	count := 1_000_000
	badTotal := int(unsafe.Sizeof(bad)) * count
	goodTotal := int(unsafe.Sizeof(good)) * count
	fmt.Printf("\n구조체 %d개 생성 시:\n", count)
	fmt.Printf("BadLayout:  %d 바이트 (약 %.1f MB)\n", badTotal, float64(badTotal)/1024/1024)
	fmt.Printf("GoodLayout: %d 바이트 (약 %.1f MB)\n", goodTotal, float64(goodTotal)/1024/1024)
	fmt.Printf("절약량:     %d 바이트 (약 %.1f MB)\n", badTotal-goodTotal, float64(badTotal-goodTotal)/1024/1024)
}
