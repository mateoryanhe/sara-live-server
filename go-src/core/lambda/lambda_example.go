package lambda

import (
	"fmt"
	"strings"
)

func main() {
	// 示例1: 处理整数切片
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 使用Filter筛选偶数
	evens := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println("偶数:", evens) // 输出: 偶数: [2 4 6 8 10]

	// 使用Map将整数转换为其平方
	squares := Map(evens, func(n int) int {
		return n * n
	})
	fmt.Println("平方:", squares) // 输出: 平方: [4 16 36 64 100]

	// 使用Reduce计算总和
	sum := Reduce(squares, func(acc int, n int) int {
		return acc + n
	}, 0)
	fmt.Println("总和:", sum) // 输出: 总和: 120

	// 示例2: 处理字符串切片
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}

	// 筛选长度大于5的单词
	longWords := Filter(words, func(s string) bool {
		return len(s) > 5
	})
	fmt.Println("长单词:", longWords) // 输出: 长单词: [banana cherry elderberry]

	// 将单词转换为大写
	upperWords := Map(longWords, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Println("大写单词:", upperWords) // 输出: 大写单词: [BANANA CHERRY ELDERBERRY]

	// 使用ForEach打印每个单词
	fmt.Print("遍历输出: ")
	ForEach(upperWords, func(s string) {
		fmt.Print(s + " ")
	})
	// 输出: 遍历输出: BANANA CHERRY ELDERBERRY

	// 检查是否有以"C"开头的单词
	hasCWord := AnyMatch(upperWords, func(s string) bool {
		return strings.HasPrefix(s, "C")
	})
	fmt.Println("\n有以C开头的单词:", hasCWord) // 输出: 有以C开头的单词: true
}
