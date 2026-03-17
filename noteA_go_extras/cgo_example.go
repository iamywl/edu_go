package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

// C 함수: 두 정수의 합을 반환
int add(int a, int b) {
    return a + b;
}

// C 함수: 문자열을 대문자로 변환 (제자리 변환)
void toUpperCase(char* s) {
    for (int i = 0; s[i]; i++) {
        if (s[i] >= 'a' && s[i] <= 'z') {
            s[i] -= 32;
        }
    }
}

// C 구조체 정의
typedef struct {
    double x;
    double y;
} Point;

// C 함수: 두 점 사이의 거리 계산
double distance(Point p1, Point p2) {
    double dx = p1.x - p2.x;
    double dy = p1.y - p2.y;
    return sqrt(dx*dx + dy*dy);
}
*/
import "C" // 이 줄은 위의 C 주석 바로 아래에 있어야 한다 (빈 줄 불가!)

import (
	"fmt"
	"unsafe"
)

func main() {
	// =============================================
	// cgo 기본 예제
	// 빌드: CGO_ENABLED=1 go run cgo_example.go
	// =============================================

	// --- C 함수 호출: 정수 덧셈 ---
	fmt.Println("=== C 함수 호출: 정수 덧셈 ===")
	result := C.add(C.int(3), C.int(7))
	fmt.Printf("  C.add(3, 7) = %d\n", int(result))

	// --- Go 문자열 → C 문자열 변환 ---
	fmt.Println("\n=== 문자열 변환 및 대문자 변환 ===")
	goStr := "hello, cgo world"

	// Go 문자열을 C 문자열로 변환
	// C.CString은 C 힙에 메모리를 할당하므로 반드시 free 해야 한다!
	cStr := C.CString(goStr)
	defer C.free(unsafe.Pointer(cStr)) // 메모리 누수 방지

	fmt.Printf("  변환 전: %s\n", goStr)

	// C 함수로 대문자 변환
	C.toUpperCase(cStr)

	// C 문자열을 다시 Go 문자열로 변환
	resultStr := C.GoString(cStr)
	fmt.Printf("  변환 후: %s\n", resultStr)

	// --- C 구조체 사용 ---
	fmt.Println("\n=== C 구조체 사용: 두 점 사이 거리 ===")

	// C 구조체를 Go에서 생성
	p1 := C.Point{x: C.double(0), y: C.double(0)}
	p2 := C.Point{x: C.double(3), y: C.double(4)}

	dist := C.distance(p1, p2)
	fmt.Printf("  점1(%.0f, %.0f)과 점2(%.0f, %.0f) 사이 거리: %.2f\n",
		float64(p1.x), float64(p1.y),
		float64(p2.x), float64(p2.y),
		float64(dist))

	// --- C 타입 크기 확인 ---
	fmt.Println("\n=== C 타입 크기 (바이트) ===")
	fmt.Printf("  C.int    크기: %d\n", unsafe.Sizeof(C.int(0)))
	fmt.Printf("  C.double 크기: %d\n", unsafe.Sizeof(C.double(0)))
	fmt.Printf("  C.char   크기: %d\n", unsafe.Sizeof(C.char(0)))
	fmt.Printf("  C.Point  크기: %d\n", unsafe.Sizeof(C.Point{}))

	// --- 주의사항 안내 ---
	fmt.Println("\n=== cgo 사용 시 주의사항 ===")
	fmt.Println("  1. C.CString()으로 할당한 메모리는 반드시 C.free()로 해제")
	fmt.Println("  2. import \"C\" 바로 위에 C 코드 주석이 있어야 함 (빈 줄 불가)")
	fmt.Println("  3. 크로스 컴파일 시 CGO_ENABLED=0이 기본값")
	fmt.Println("  4. cgo 호출은 일반 Go 함수 호출보다 느림 (약 100배)")
}
