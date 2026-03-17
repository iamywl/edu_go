package main

import (
	"container/list"
	"fmt"
)

// ============================================
// 리스트 (container/list) - 이중 연결 리스트
// ============================================

func main() {
	// 1. 리스트 생성과 기본 연산
	fmt.Println("=== 리스트 생성과 기본 연산 ===")
	l := list.New()

	// 뒤에 추가 (PushBack)
	l.PushBack("바나나")
	l.PushBack("체리")
	l.PushBack("딸기")

	// 앞에 추가 (PushFront)
	l.PushFront("사과")

	fmt.Printf("리스트 길이: %d\n", l.Len())
	printList(l)

	// 2. 특정 위치에 삽입
	fmt.Println("\n=== 특정 위치에 삽입 ===")
	// "바나나" 요소를 찾아서 그 앞뒤에 삽입
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == "바나나" {
			l.InsertBefore("포도", e) // 바나나 앞에 포도 삽입
			l.InsertAfter("망고", e)  // 바나나 뒤에 망고 삽입
			break
		}
	}
	printList(l)

	// 3. 요소 제거
	fmt.Println("\n=== 요소 제거 ===")
	// "체리" 제거
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == "체리" {
			l.Remove(e)
			break
		}
	}
	fmt.Println("체리 제거 후:")
	printList(l)

	// 4. 순회 (앞에서 뒤로, 뒤에서 앞으로)
	fmt.Println("\n=== 순회 ===")

	fmt.Print("앞 -> 뒤: ")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value, " ")
	}
	fmt.Println()

	fmt.Print("뒤 -> 앞: ")
	for e := l.Back(); e != nil; e = e.Prev() {
		fmt.Print(e.Value, " ")
	}
	fmt.Println()

	// 5. 요소 이동
	fmt.Println("\n=== 요소 이동 ===")
	fmt.Println("이동 전:")
	printList(l)

	// 맨 뒤 요소를 맨 앞으로
	l.MoveToFront(l.Back())
	fmt.Println("맨 뒤 -> 맨 앞:")
	printList(l)

	// 맨 앞 요소를 맨 뒤로
	l.MoveToBack(l.Front())
	fmt.Println("맨 앞 -> 맨 뒤:")
	printList(l)

	// 6. 정수 리스트 - 정렬 삽입 예제
	fmt.Println("\n=== 정렬 삽입 ===")
	sortedList := list.New()
	values := []int{30, 10, 50, 20, 40, 15, 35}

	for _, v := range values {
		insertSorted(sortedList, v)
		fmt.Printf("  %d 삽입 후: ", v)
		printIntList(sortedList)
	}

	// 7. 리스트를 슬라이스로 변환
	fmt.Println("\n=== 리스트 -> 슬라이스 변환 ===")
	slice := listToSlice(sortedList)
	fmt.Println("슬라이스:", slice)

	// 8. 리스트를 큐(Queue)처럼 사용
	fmt.Println("\n=== 큐(Queue)로 사용 ===")
	queue := list.New()

	// 인큐 (Enqueue) - 뒤에 추가
	fmt.Println("인큐: 작업1, 작업2, 작업3")
	queue.PushBack("작업1")
	queue.PushBack("작업2")
	queue.PushBack("작업3")

	// 디큐 (Dequeue) - 앞에서 제거
	for queue.Len() > 0 {
		front := queue.Front()
		fmt.Printf("디큐: %v (남은 작업: %d)\n", front.Value, queue.Len()-1)
		queue.Remove(front)
	}

	// 9. 리스트를 스택(Stack)처럼 사용
	fmt.Println("\n=== 스택(Stack)으로 사용 ===")
	stack := list.New()

	// 푸시 (Push) - 뒤에 추가
	fmt.Println("푸시: A, B, C")
	stack.PushBack("A")
	stack.PushBack("B")
	stack.PushBack("C")

	// 팝 (Pop) - 뒤에서 제거
	for stack.Len() > 0 {
		back := stack.Back()
		fmt.Printf("팝: %v (남은 요소: %d)\n", back.Value, stack.Len()-1)
		stack.Remove(back)
	}
}

// printList는 리스트의 모든 요소를 출력한다
func printList(l *list.List) {
	fmt.Print("  [")
	for e := l.Front(); e != nil; e = e.Next() {
		if e != l.Front() {
			fmt.Print(" -> ")
		}
		fmt.Print(e.Value)
	}
	fmt.Println("]")
}

// printIntList는 정수 리스트를 출력한다
func printIntList(l *list.List) {
	fmt.Print("[")
	for e := l.Front(); e != nil; e = e.Next() {
		if e != l.Front() {
			fmt.Print(", ")
		}
		fmt.Print(e.Value)
	}
	fmt.Println("]")
}

// insertSorted는 정렬된 순서를 유지하며 값을 삽입한다
func insertSorted(l *list.List, val int) {
	// 빈 리스트이면 바로 추가
	if l.Len() == 0 {
		l.PushBack(val)
		return
	}

	// 적절한 위치를 찾아 삽입
	for e := l.Front(); e != nil; e = e.Next() {
		if val < e.Value.(int) {
			l.InsertBefore(val, e)
			return
		}
	}

	// 모든 요소보다 크면 맨 뒤에 추가
	l.PushBack(val)
}

// listToSlice는 정수 리스트를 슬라이스로 변환한다
func listToSlice(l *list.List) []int {
	result := make([]int, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(int))
	}
	return result
}
