package main

import "fmt"

// Address 구조체 — 주소 정보
type Address struct {
	City   string
	Street string
	Zip    string
}

// Company 구조체 — 회사 정보
type Company struct {
	Name    string
	Address // 내장(임베딩) 구조체 — 필드명 생략!
}

// Employee 구조체 — 직원 정보
type Employee struct {
	Name    string
	Age     int
	Company // Company를 임베딩
}

// ContactInfo 구조체 — 일반 포함 방식 비교용
type ContactInfo struct {
	Name    string
	Address Address // 필드명을 명시한 일반 포함
}

func main() {
	// === 내장(Embedded) 구조체 ===
	fmt.Println("=== 내장(Embedded) 구조체 ===")

	emp := Employee{
		Name: "홍길동",
		Age:  30,
		Company: Company{
			Name: "Go 주식회사",
			Address: Address{
				City:   "서울",
				Street: "테헤란로 123",
				Zip:    "06100",
			},
		},
	}

	// 임베딩 덕분에 중첩된 필드에 바로 접근 가능!
	fmt.Printf("직원 이름: %s\n", emp.Name)         // Employee.Name
	fmt.Printf("회사 이름: %s\n", emp.Company.Name) // Company.Name (이름 충돌 시 명시)
	fmt.Printf("도시: %s\n", emp.City)            // Address.City에 바로 접근!
	fmt.Printf("거리: %s\n", emp.Street)          // Address.Street에 바로 접근!
	fmt.Printf("우편번호: %s\n", emp.Zip)           // Address.Zip에 바로 접근!

	fmt.Println()

	// === 일반 포함 방식과 비교 ===
	fmt.Println("=== 일반 포함 (필드명 명시) ===")

	contact := ContactInfo{
		Name: "김철수",
		Address: Address{
			City:   "부산",
			Street: "해운대로 456",
			Zip:    "48000",
		},
	}

	// 일반 포함 방식은 반드시 필드명을 거쳐야 함
	fmt.Printf("이름: %s\n", contact.Name)
	fmt.Printf("도시: %s\n", contact.Address.City) // contact.City 불가!
	fmt.Printf("거리: %s\n", contact.Address.Street)

	fmt.Println()

	// === 필드명 충돌 시 동작 ===
	fmt.Println("=== 필드명 충돌 예시 ===")
	// Employee.Name과 Company.Name이 둘 다 있다
	// 이 경우 바깥쪽(Employee)의 Name이 우선한다
	fmt.Printf("emp.Name → Employee의 Name: %s\n", emp.Name)
	fmt.Printf("emp.Company.Name → Company의 Name: %s\n", emp.Company.Name)
}
