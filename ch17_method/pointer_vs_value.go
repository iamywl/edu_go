package main

import "fmt"

// ============================================
// 포인터 리시버 vs 값 타입 리시버
// ============================================

// Player 구조체 - 게임 캐릭터
type Player struct {
	Name   string
	HP     int
	MaxHP  int
	Attack int
	Level  int
	Exp    int
}

// --- 값 타입 리시버 (읽기 전용) ---

// String은 플레이어 정보를 문자열로 반환한다
func (p Player) String() string {
	return fmt.Sprintf("[%s] Lv.%d HP:%d/%d ATK:%d EXP:%d",
		p.Name, p.Level, p.HP, p.MaxHP, p.Attack, p.Exp)
}

// IsAlive는 플레이어가 살아있는지 확인한다
func (p Player) IsAlive() bool {
	return p.HP > 0
}

// HPPercent는 HP 비율을 반환한다
func (p Player) HPPercent() float64 {
	return float64(p.HP) / float64(p.MaxHP) * 100
}

// --- 포인터 리시버 (값 변경) ---

// TakeDamage는 데미지를 받아 HP를 감소시킨다
func (p *Player) TakeDamage(damage int) {
	p.HP -= damage
	if p.HP < 0 {
		p.HP = 0
	}
	fmt.Printf("  %s이(가) %d의 데미지를 받았다! (HP: %d/%d)\n",
		p.Name, damage, p.HP, p.MaxHP)
}

// Heal은 HP를 회복시킨다
func (p *Player) Heal(amount int) {
	p.HP += amount
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	fmt.Printf("  %s이(가) %d만큼 회복했다! (HP: %d/%d)\n",
		p.Name, amount, p.HP, p.MaxHP)
}

// GainExp는 경험치를 획득하고 레벨업을 체크한다
func (p *Player) GainExp(exp int) {
	p.Exp += exp
	fmt.Printf("  %s이(가) 경험치 %d를 획득했다!\n", p.Name, exp)

	// 경험치 100 이상이면 레벨업
	for p.Exp >= 100 {
		p.Exp -= 100
		p.LevelUp()
	}
}

// LevelUp은 레벨을 올리고 능력치를 증가시킨다
func (p *Player) LevelUp() {
	p.Level++
	p.MaxHP += 20
	p.HP = p.MaxHP // 레벨업 시 풀 회복
	p.Attack += 5
	fmt.Printf("  ** %s 레벨 업! Lv.%d (HP:%d, ATK:%d) **\n",
		p.Name, p.Level, p.MaxHP, p.Attack)
}

// ============================================
// 값 타입 리시버로 변경을 시도하면?
// ============================================

// WrongHeal은 값 타입 리시버로 HP를 회복하려는 잘못된 예시
func (p Player) WrongHeal(amount int) {
	p.HP += amount // 복사본의 HP만 변경됨, 원본은 변경되지 않음!
	fmt.Printf("  (잘못된 회복) HP: %d/%d\n", p.HP, p.MaxHP)
}

func main() {
	// 1. 값 타입 리시버로는 원본을 변경할 수 없음
	fmt.Println("=== 값 타입 리시버의 한계 ===")
	hero := Player{
		Name:   "용사",
		HP:     100,
		MaxHP:  100,
		Attack: 20,
		Level:  1,
		Exp:    0,
	}
	fmt.Println("원본:", hero)

	hero.WrongHeal(50)                // 복사본만 변경됨
	fmt.Println("WrongHeal 후:", hero) // HP 변경 없음!
	fmt.Println("-> 값 타입 리시버는 복사본을 수정하므로 원본에 영향 없음")

	// 2. 포인터 리시버로 원본 변경
	fmt.Println("\n=== 포인터 리시버로 원본 변경 ===")
	fmt.Println("원본:", hero)

	hero.TakeDamage(30) // HP: 70/100
	fmt.Println("데미지 후:", hero)

	hero.Heal(50) // HP: 100/100 (최대 HP 초과 불가)
	fmt.Println("회복 후:", hero)

	// 3. 게임 시뮬레이션
	fmt.Println("\n=== 게임 시뮬레이션 ===")
	player := Player{
		Name:   "전사",
		HP:     80,
		MaxHP:  80,
		Attack: 15,
		Level:  1,
		Exp:    0,
	}
	fmt.Println("시작:", player)

	// 전투
	fmt.Println("\n-- 몬스터와 전투! --")
	player.TakeDamage(25)
	player.GainExp(40)

	fmt.Println("\n-- 보스와 전투! --")
	player.TakeDamage(50)
	player.GainExp(70) // 경험치 110 -> 레벨업!

	fmt.Println("\n최종 상태:", player)
	fmt.Printf("HP 비율: %.1f%%\n", player.HPPercent())
	fmt.Println("생존 여부:", player.IsAlive())

	// 4. 자동 변환 확인
	fmt.Println("\n=== 자동 변환 ===")

	// 값 타입 변수로 포인터 메서드 호출 가능
	val := Player{Name: "마법사", HP: 50, MaxHP: 50, Attack: 30, Level: 1}
	val.TakeDamage(10) // Go가 자동으로 (&val).TakeDamage(10)로 변환
	fmt.Println("값 변수로 포인터 메서드 호출:", val)

	// 포인터 변수로 값 타입 메서드 호출 가능
	ptr := &Player{Name: "궁수", HP: 60, MaxHP: 60, Attack: 25, Level: 1}
	fmt.Println("포인터로 값 메서드 호출:", ptr.IsAlive()) // (*ptr).IsAlive()로 자동 변환
}
