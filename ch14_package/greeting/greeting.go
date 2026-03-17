// Package greeting 은 인사말을 생성하는 패키지이다.
// init() 함수로 기본 언어를 초기화하는 예제를 포함한다.
package greeting

import "fmt"

// DefaultLang 은 기본 인사 언어이다. (대문자 → 외부 공개)
var DefaultLang string

// supportedLangs 는 지원하는 언어 목록이다. (소문자 → 비공개)
var supportedLangs = []string{"ko", "en", "ja"}

// init 함수 — 패키지가 import될 때 자동으로 실행됩니다.
// main() 보다 먼저 실행되어 패키지를 초기화한다.
func init() {
	DefaultLang = "ko" // 기본 언어를 한국어로 설정
	fmt.Println("[greeting] init 실행: 기본 언어를 'ko'로 설정")
}

// 두 번째 init 함수 — 한 파일에 여러 init을 정의할 수 있습니다.
func init() {
	fmt.Printf("[greeting] init 실행: 지원 언어 %v\n", supportedLangs)
}

// Hello 는 이름을 받아 현재 언어에 맞는 인사말을 반환한다.
func Hello(name string) string {
	switch DefaultLang {
	case "ko":
		return fmt.Sprintf("안녕하세요, %s님!", name)
	case "en":
		return fmt.Sprintf("Hello, %s!", name)
	case "ja":
		return fmt.Sprintf("こんにちは、%sさん!", name)
	default:
		return fmt.Sprintf("Hi, %s!", name)
	}
}

// SetLang 은 기본 언어를 변경한다.
func SetLang(lang string) {
	if !isSupported(lang) {
		fmt.Printf("지원하지 않는 언어: %s\n", lang)
		return
	}
	DefaultLang = lang
}

// isSupported 는 해당 언어가 지원되는지 확인한다. (소문자 → 비공개)
func isSupported(lang string) bool {
	for _, l := range supportedLangs {
		if l == lang {
			return true
		}
	}
	return false
}
