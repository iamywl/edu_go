# 8장. switch문

---

## 실행 방법

Docker 컨테이너 내부에서 다음과 같이 실행한다:

```bash
# 컨테이너 접속 (최초 1회)
make shell

# 예제 실행
go run ch08_switch/basic.go
go run ch08_switch/advanced.go
go run ch08_switch/enum_switch.go
```

> **Makefile 활용**: `make run CH=ch08_switch` 또는 `make run CH=ch08_switch FILE=basic.go`

---

## 8.1 switch문 동작 원리

`switch`문은 값에 따라 여러 갈래로 분기하는 제어문이다. `if-else if`를 여러 개 쓰는 것보다 코드가 깔끔하고, 분기 대상이 명확하게 드러난다.

```go
day := "화요일"

switch day {
case "월요일":
    fmt.Println("월요일입니다")
case "화요일":
    fmt.Println("화요일입니다")
case "수요일":
    fmt.Println("수요일입니다")
default:
    fmt.Println("다른 요일입니다")
}
```

Go의 `switch`는 C/Java와 다르게 동작한다:
- **`break`가 자동 적용**된다. 매칭된 case만 실행하고 빠져나간다. C/Java에서처럼 매 case마다 `break`를 쓸 필요가 없다.
- **`fallthrough`**를 명시해야 다음 case로 넘어간다 (C와 반대 동작이다).
- case의 값은 상수일 필요가 없다. 변수나 함수 호출 결과도 사용할 수 있다.

`switch`문의 실행 순서는 다음과 같다:
1. `switch` 키워드 뒤의 표현식을 평가한다.
2. 위에서부터 순서대로 `case` 값과 비교한다.
3. 일치하는 `case`를 찾으면 해당 블록을 실행하고 `switch`를 빠져나간다.
4. 일치하는 `case`가 없으면 `default` 블록을 실행한다 (있는 경우).

---

## 8.2 switch문을 언제 쓰는가?

### 하나의 값을 여러 값과 비교할 때

```go
// if-else if로 작성 (장황하다)
if cmd == "start" {
    // ...
} else if cmd == "stop" {
    // ...
} else if cmd == "restart" {
    // ...
}

// switch로 작성 (깔끔하다)
switch cmd {
case "start":
    // ...
case "stop":
    // ...
case "restart":
    // ...
}
```

동일한 변수를 여러 값과 반복 비교하는 상황에서는 `switch`가 `if-else if`보다 의도를 명확하게 전달한다. 비교 대상이 3개 이상이면 `switch`를 사용하는 것이 일반적이다.

### 열거값(enum)을 처리할 때

```go
switch status {
case StatusActive:
    // ...
case StatusInactive:
    // ...
case StatusDeleted:
    // ...
}
```

`iota`로 정의한 열거값을 처리할 때 `switch`는 필수적인 도구이다. 모든 열거값을 빠짐없이 처리했는지 코드에서 한눈에 확인할 수 있다.

### if-else 체인이 길어질 때

분기가 4개 이상이 되면 `if-else if` 체인의 가독성이 급격히 떨어진다. 이런 경우 `switch`로 전환하면 코드가 훨씬 읽기 쉬워진다.

---

## 8.3 다양한 switch문 형태

### 한 case에 여러 값

쉼표로 구분하여 하나의 `case`에 여러 값을 나열할 수 있다. 나열된 값 중 하나라도 일치하면 해당 블록이 실행된다.

```go
switch day {
case "토요일", "일요일":
    fmt.Println("주말입니다")
default:
    fmt.Println("평일입니다")
}
```

이 문법은 C/Java에서 `case`를 연속으로 나열하는 것(fall-through를 이용)과 같은 효과이지만, 훨씬 간결하다.

### 조건 없는 switch (if-else 대용)

`switch` 뒤에 비교 대상을 생략하면, 각 `case`에 조건식을 쓸 수 있다. 이 형태는 사실상 `if-else if` 체인과 동일하지만, 시각적으로 더 정돈된 느낌을 준다.

```go
score := 85

switch {
case score >= 90:
    fmt.Println("A")
case score >= 80:
    fmt.Println("B")
case score >= 70:
    fmt.Println("C")
default:
    fmt.Println("F")
}
```

조건 없는 `switch`는 `switch true`와 동일하다. 각 `case`의 조건이 `true`인 것을 찾는 것이다.

### switch 초기문

`if`문처럼 `switch`에도 초기문을 쓸 수 있다. 초기문에서 선언한 변수는 `switch` 블록 안에서만 유효하다.

```go
switch os := runtime.GOOS; os {
case "linux":
    fmt.Println("리눅스")
case "darwin":
    fmt.Println("macOS")
default:
    fmt.Println(os)
}
// 여기서 os는 사용 불가
```

### 타입 switch

인터페이스 값의 실제 타입을 검사할 때 사용하는 특수한 형태의 `switch`이다. 아직 인터페이스를 배우지 않았지만, 이런 문법이 있다는 것을 알아두면 좋다.

```go
switch v := value.(type) {
case int:
    fmt.Println("정수:", v)
case string:
    fmt.Println("문자열:", v)
default:
    fmt.Println("알 수 없는 타입")
}
```

---

## 8.4 const 열거값과 switch

`iota`로 정의한 상수를 `switch`로 처리하는 것은 Go에서 매우 흔한 패턴이다. 사용자 정의 타입을 만들어 열거값의 타입 안전성을 높이는 것이 좋은 습관이다.

```go
type Color int

const (
    Red Color = iota
    Green
    Blue
)

func colorName(c Color) string {
    switch c {
    case Red:
        return "빨강"
    case Green:
        return "초록"
    case Blue:
        return "파랑"
    default:
        return "알 수 없는 색"
    }
}
```

`default` case를 항상 포함하는 것이 좋은 습관이다. 나중에 새로운 열거값이 추가되었을 때, `default`가 없으면 아무 동작도 하지 않고 조용히 넘어갈 수 있기 때문이다.

### Stringer 인터페이스 패턴

열거값에 `String()` 메서드를 구현하면 `fmt.Println` 등에서 자동으로 문자열로 변환된다. 이 패턴은 Go에서 매우 자주 사용된다.

```go
func (c Color) String() string {
    switch c {
    case Red:
        return "빨강"
    case Green:
        return "초록"
    case Blue:
        return "파랑"
    default:
        return "알 수 없는 색"
    }
}

fmt.Println(Red)  // "빨강" 출력
```

---

## 8.5 break와 fallthrough 키워드

### break

Go의 `switch`는 매칭된 case만 실행하고 자동으로 빠져나간다. 하지만 case 안에서 조기 종료하고 싶을 때 `break`를 명시할 수 있다. 주로 `if`문과 함께 사용하여 case 블록의 나머지 부분을 건너뛸 때 유용하다.

```go
switch num {
case 1:
    if someCondition {
        break // 이 case 나머지를 건너뛰고 switch 종료
    }
    fmt.Println("조건이 아닐 때만 실행")
}
```

### fallthrough

`fallthrough`를 쓰면 다음 case의 조건을 검사하지 않고 **무조건** 실행한다. 다음 case의 조건이 거짓이더라도 실행된다는 점에 주의해야 한다.

```go
switch num := 3; {
case num >= 3:
    fmt.Println("3 이상")
    fallthrough
case num >= 2:
    fmt.Println("2 이상")
    fallthrough
case num >= 1:
    fmt.Println("1 이상")
}
// 출력:
// 3 이상
// 2 이상
// 1 이상
```

> **주의:** `fallthrough`는 Go에서 자주 사용되지 않는다. 대부분의 경우 한 case에 여러 값을 나열하거나, 조건 없는 switch를 사용하는 것이 더 깔끔하다. 필요한 경우에만 신중하게 사용해야 한다.

### fallthrough의 제약사항

- `fallthrough`는 case 블록의 마지막 문장이어야 한다.
- 타입 switch에서는 `fallthrough`를 사용할 수 없다.
- `fallthrough`는 다음 case의 조건을 무시하므로, 예상치 못한 동작을 유발할 수 있다.

---

## 핵심 요약

| 항목 | 설명 |
|------|------|
| `switch 값 { case: }` | 값에 따라 분기한다 |
| 여러 값 매칭 | `case "a", "b", "c":` 형태로 나열한다 |
| 조건 없는 switch | `switch { case 조건: }` 형태로 if-else 대용이다 |
| 자동 break | C와 달리 break가 불필요하다 |
| `fallthrough` | 다음 case를 무조건 실행한다 |
| `default` | 어떤 case에도 해당하지 않을 때 실행된다 |
| switch 초기문 | `switch 초기문; 값 { }` 형태로 변수 스코프를 제한한다 |

---

## 연습문제

1. 월(1~12)을 변수로 선언하고 해당 월의 계절을 출력하는 프로그램을 작성하라. (3~5월: 봄, 6~8월: 여름, 9~11월: 가을, 12~2월: 겨울) 한 case에 여러 값을 나열하는 방식을 사용하라.

2. 조건 없는 switch를 사용하여 BMI(체질량지수) 계산기를 만들어라. (18.5 미만: 저체중, 18.5~24.9: 정상, 25~29.9: 과체중, 30 이상: 비만)

3. `iota`로 HTTP 상태 코드 그룹(2xx 성공, 4xx 클라이언트 에러, 5xx 서버 에러)을 정의하고, `switch`로 각 그룹에 맞는 메시지를 출력하라.

4. `switch` 초기문을 사용하여 난수를 생성하고, 그 난수가 짝수인지 홀수인지 출력하는 코드를 작성하라.

5. 다음 코드의 출력 결과를 예측하고, 그 이유를 설명하라.
   ```go
   x := 5
   switch {
   case x > 3:
       fmt.Println("A")
       fallthrough
   case x > 10:
       fmt.Println("B")
   case x > 1:
       fmt.Println("C")
   }
   ```

6. 간단한 계산기를 `switch`문으로 구현하라. 연산자("+", "-", "*", "/")와 두 개의 피연산자를 변수로 선언하고, 연산자에 따라 결과를 출력하라. 0으로 나누는 경우와 잘못된 연산자 입력도 처리하라.

7. 년도와 월을 변수로 선언하고, `switch`를 사용하여 해당 월의 일수를 출력하는 프로그램을 작성하라. 2월의 경우 윤년 여부를 고려하라.

8. `fallthrough`를 사용하여, 입력된 등급(1~5)에 따라 누적 혜택을 출력하는 프로그램을 작성하라. (예: 등급 3이면 등급 3, 2, 1의 혜택을 모두 출력)

9. 타입에 따라 다른 동작을 하는 `switch`를 작성하라. `interface{}` 타입 변수에 `int`, `string`, `bool` 값을 넣어보고, 타입 switch로 각각 다른 메시지를 출력하라.

10. `switch`문과 `if`문의 사용 기준에 대해 자신만의 가이드라인을 정리하라. 어떤 상황에서 `switch`가 더 적절하고, 어떤 상황에서 `if`가 더 적절한지 예제와 함께 설명하라.

---

## 구현 과제

### 과제 1: 요일별 일정 안내기
요일(월~일)을 변수로 설정하고, `switch`를 사용하여 해당 요일의 일정을 출력하는 프로그램을 작성하라. 평일에는 각각 다른 수업 시간표를, 주말에는 자유 활동 안내를 출력하라. `default`로 잘못된 입력도 처리하라.

### 과제 2: 메뉴 선택 시스템
음식점 메뉴 시스템을 `switch`로 구현하라. 메뉴 번호(1~5)를 선택하면 음식 이름과 가격을 출력하고, 수량을 곱하여 총 금액을 계산하라. 잘못된 메뉴 번호, 0 이하의 수량 등 예외 상황도 처리하라.

### 과제 3: 학점 변환기
점수(0~100)를 입력받아 학점으로 변환하되, 조건 없는 switch를 사용하라. A+(95 이상), A(90 이상), B+(85 이상), B(80 이상), C+(75 이상), C(70 이상), D+(65 이상), D(60 이상), F(60 미만)로 세분화하라. 각 학점의 평점(4.5, 4.0, 3.5 등)도 함께 출력하라.

### 과제 4: 간이 명령어 처리기
문자열 명령어를 `switch`로 처리하는 프로그램을 작성하라. "help", "version", "list", "add", "delete", "quit" 명령어를 지원하고, 각 명령어에 대한 설명 메시지를 출력하라. "quit"일 때는 "종료합니다" 메시지를 출력하라. 인식할 수 없는 명령어에 대해서도 안내 메시지를 제공하라.

### 과제 5: 별자리 판별기
생일의 월과 일을 변수로 설정하고, `switch`를 활용하여 해당 날짜의 별자리를 출력하는 프로그램을 작성하라. 12개 별자리의 날짜 범위를 정확히 처리하라. (예: 양자리 3/21~4/19, 황소자리 4/20~5/20 등)

---

## 프로젝트 과제

### 프로젝트 1: 자판기 시뮬레이터
자판기 프로그램을 구현하라. 다음 기능을 `switch`문으로 처리하라:
- 음료 목록 표시 (이름, 가격을 상수로 정의)
- 음료 선택 (번호 입력)
- 투입 금액 확인 및 잔돈 계산
- 재고 관리 (각 음료의 재고를 변수로 관리)
- 품절 처리
- 여러 음료를 연속으로 구매할 수 있도록 반복문과 조합하라 (for문 사용 가능)

### 프로젝트 2: 텍스트 기반 RPG 전투 시스템
간단한 RPG 전투 시스템을 구현하라. 플레이어와 몬스터의 HP, 공격력을 변수로 관리하고, 매 턴마다 `switch`로 행동을 선택하라:
- "attack": 기본 공격 (공격력만큼 상대 HP 감소)
- "defend": 방어 (다음 턴 받는 피해 50% 감소)
- "heal": 회복 (HP 일부 회복, 최대 HP 초과 불가)
- "flee": 도망 (50% 확률로 성공)
- 몬스터는 매 턴 자동으로 공격한다
- HP가 0 이하가 되면 전투 종료하고 결과를 출력하라
