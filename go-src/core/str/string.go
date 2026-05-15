package str

import (
	"strings"
)

// TruncateStringBetween 截取两个字符串之间的内容
// start: 开始字符串
// end: 结束字符串
// str: 原始字符串
// 返回开始字符串和结束字符串之间的内容，不包含开始和结束字符串本身
func TruncateStringBetween(str, start, end string) string {
	// 查找开始字符串的位置
	startIndex := strings.Index(str, start)
	if startIndex == -1 {
		return ""
	}

	// 调整开始位置到开始字符串之后
	startIndex += len(start)

	// 从开始位置之后查找结束字符串
	endIndex := strings.Index(str[startIndex:], end)
	if endIndex == -1 {
		return ""
	}

	// 计算结束位置在原始字符串中的索引
	endIndex += startIndex

	// 返回截取的内容
	return str[startIndex:endIndex]
}

func GetCountSQL(sql string) string {
	temp := TruncateStringBetween(sql, "from", "order")
	countSQL := "select count(1) from  " + temp
	return countSQL
}

// TruncateStringAfter 截取指定字符串之后的内容
// delimiter: 分隔字符串
// str: 原始字符串
// 返回分隔字符串之后的内容，不包含分隔字符串本身
func TruncateStringAfter(str, delimiter string) string {
	index := strings.Index(str, delimiter)
	if index == -1 {
		return ""
	}

	return str[index+len(delimiter):]
}

// TruncateStringBefore 截取指定字符串之前的内容
// delimiter: 分隔字符串
// str: 原始字符串
// 返回分隔字符串之前的内容，不包含分隔字符串本身
func TruncateStringBefore(str, delimiter string) string {
	index := strings.Index(str, delimiter)
	if index == -1 {
		return str // 如果没找到分隔符，返回整个字符串
	}

	return str[:index]
}
