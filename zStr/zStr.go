package zStr

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"
	"unicode"
)

// Substring 截取字符串（支持中文）
func Substring(str string, start, length int) string {
	if start < 0 {
		start = 0
	}
	
	chars := []rune(str)
	charLen := len(chars)
	
	if start >= charLen {
		return ""
	}
	
	end := start + length
	if end > charLen {
		end = charLen
	}
	
	return string(chars[start:end])
}

// Truncate 截断字符串并添加省略号（支持中文）
func Truncate(str string, maxLength int, suffix string) string {
	if maxLength <= 0 {
		return ""
	}
	
	chars := []rune(str)
	if len(chars) <= maxLength {
		return str
	}
	
	if suffix == "" {
		suffix = "..."
	}
	
	return string(chars[:maxLength]) + suffix
}

// Contains 检查字符串是否包含子串（支持中文）
func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// ContainsAny 检查字符串是否包含任一子串
func ContainsAny(str string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(str, substr) {
			return true
		}
	}
	return false
}

// ContainsAll 检查字符串是否包含所有子串
func ContainsAll(str string, substrs ...string) bool {
	for _, substr := range substrs {
		if !strings.Contains(str, substr) {
			return false
		}
	}
	return true
}

// Replace 替换字符串中的子串（支持中文）
func Replace(str, old, new string, n int) string {
	return strings.Replace(str, old, new, n)
}

// ReplaceAll 替换字符串中的所有子串（支持中文）
func ReplaceAll(str, old, new string) string {
	return strings.ReplaceAll(str, old, new)
}

// ReplaceRegexp 使用正则表达式替换字符串
func ReplaceRegexp(str, pattern, repl string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(str, repl), nil
}

// Split 分割字符串
func Split(str, sep string) []string {
	return strings.Split(str, sep)
}

// SplitN 分割字符串，最多分割n次
func SplitN(str, sep string, n int) []string {
	return strings.SplitN(str, sep, n)
}

// Join 连接字符串
func Join(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

// Trim 去除字符串首尾的空白字符（支持中文）
func Trim(str string, cutset string) string {
	if cutset == "" {
		return strings.TrimSpace(str)
	}
	return strings.Trim(str, cutset)
}

// TrimLeft 去除字符串左侧的空白字符（支持中文）
func TrimLeft(str string, cutset string) string {
	if cutset == "" {
		return strings.TrimLeftFunc(str, unicode.IsSpace)
	}
	return strings.TrimLeft(str, cutset)
}

// TrimRight 去除字符串右侧的空白字符（支持中文）
func TrimRight(str string, cutset string) string {
	if cutset == "" {
		return strings.TrimRightFunc(str, unicode.IsSpace)
	}
	return strings.TrimRight(str, cutset)
}

// PadLeft 在字符串左侧填充指定字符
func PadLeft(str string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	
	if len(str) >= length {
		return str
	}
	
	padding := strings.Repeat(pad, length-len(str))
	return padding + str
}

// PadRight 在字符串右侧填充指定字符
func PadRight(str string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	
	if len(str) >= length {
		return str
	}
	
	padding := strings.Repeat(pad, length-len(str))
	return str + padding
}

// Upper 转换为大写
func Upper(str string) string {
	return strings.ToUpper(str)
}

// Lower 转换为小写
func Lower(str string) string {
	return strings.ToLower(str)
}

// Title 转换为首字母大写
func Title(str string) string {
	return strings.Title(str)
}

// CamelCase 转换为驼峰命名
func CamelCase(str string) string {
	parts := strings.Split(str, "_")
	for i, part := range parts {
		if i == 0 {
			continue
		}
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

// SnakeCase 转换为蛇形命名
func SnakeCase(str string) string {
	var result bytes.Buffer
	for i, char := range str {
		if unicode.IsUpper(char) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// KebabCase 转换为短横线命名
func KebabCase(str string) string {
	var result bytes.Buffer
	for i, char := range str {
		if unicode.IsUpper(char) {
			if i > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// Reverse 反转字符串（支持中文）
func Reverse(str string) string {
	chars := []rune(str)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

// Repeat 重复字符串
func Repeat(str string, count int) string {
	return strings.Repeat(str, count)
}

// Count 统计子串出现的次数（支持中文）
func Count(str, substr string) int {
	return strings.Count(str, substr)
}

// Index 查找子串首次出现的位置（支持中文）
func Index(str, substr string) int {
	return strings.Index(str, substr)
}

// LastIndex 查找子串最后出现的位置（支持中文）
func LastIndex(str, substr string) int {
	return strings.LastIndex(str, substr)
}

// IsEmpty 检查字符串是否为空
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsNotEmpty 检查字符串是否非空
func IsNotEmpty(str string) bool {
	return strings.TrimSpace(str) != ""
}

// IsAlpha 检查字符串是否只包含字母
func IsAlpha(str string) bool {
	for _, char := range str {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return len(str) > 0
}

// IsNumeric 检查字符串是否只包含数字
func IsNumeric(str string) bool {
	for _, char := range str {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return len(str) > 0
}

// IsAlphanumeric 检查字符串是否只包含字母和数字
func IsAlphanumeric(str string) bool {
	for _, char := range str {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return len(str) > 0
}

// IsASCII 检查字符串是否只包含ASCII字符
func IsASCII(str string) bool {
	for _, char := range str {
		if char > 127 {
			return false
		}
	}
	return true
}

// Random 生成随机字符串
func Random(length int, charset string) string {
	if length <= 0 {
		return ""
	}

	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	var result bytes.Buffer
	for i := 0; i < length; i++ {
		result.WriteByte(charset[rand.IntN(len(charset))])
	}
	return result.String()
}

// RandomLetter 生成随机字母字符串
func RandomLetter(length int) string {
	return Random(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// RandomDigit 生成随机数字字符串
func RandomDigit(length int) string {
	return Random(length, "0123456789")
}

// RandomAlphanumeric 生成随机字母和数字字符串
func RandomAlphanumeric(length int) string {
	return Random(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

// Wrap 换行字符串
func Wrap(str string, width int) string {
	if width <= 0 {
		return str
	}
	
	var result bytes.Buffer
	words := strings.Fields(str)
	currentLine := ""
	
	for _, word := range words {
		if len(currentLine) == 0 {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			result.WriteString(currentLine + "\n")
			currentLine = word
		}
	}
	
	if currentLine != "" {
		result.WriteString(currentLine + "\n")
	}
	
	return result.String()
}

// Unwrap 去除字符串中的换行符
func Unwrap(str string) string {
	// 替换各种换行符为空格
	str = strings.ReplaceAll(str, "\r\n", " ")
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "\r", " ")
	
	// 合并连续的空格
	return regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")
}

// Lines 分割字符串为行
func Lines(str string) []string {
	return strings.Split(strings.ReplaceAll(str, "\r\n", "\n"), "\n")
}

// JoinLines 连接行为字符串
func JoinLines(lines []string, separator string) string {
	return strings.Join(lines, separator)
}

// Indent 缩进字符串
func Indent(str, indent string) string {
	lines := Lines(str)
	for i, line := range lines {
		lines[i] = indent + line
	}
	return JoinLines(lines, "\n")
}

// Unindent 去除字符串的缩进
func Unindent(str string) string {
	lines := Lines(str)
	if len(lines) == 0 {
		return ""
	}
	
	// 找到最小缩进
	minIndent := -1
	for _, line := range lines {
		if IsEmpty(line) {
			continue
		}
		
		indent := 0
		for _, char := range line {
			if char == ' ' || char == '\t' {
				indent++
			} else {
				break
			}
		}
		
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}
	
	// 去除缩进
	if minIndent > 0 {
		for i, line := range lines {
			if IsEmpty(line) {
				continue
			}
			lines[i] = Substring(line, minIndent, len(line))
		}
	}
	
	return JoinLines(lines, "\n")
}

// Quote 给字符串添加引号
func Quote(str string, quote string) string {
	if quote == "" {
		quote = "\""
	}
	return quote + str + quote
}

// Unquote 去除字符串的引号
func Unquote(str string) string {
	if len(str) < 2 {
		return str
	}
	
	// 检查是否有引号
	if (str[0] == '"' && str[len(str)-1] == '"') || (str[0] == '\'' && str[len(str)-1] == '\'') {
		return str[1 : len(str)-1]
	}
	
	return str
}

// Format 格式化字符串
func Format(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Capitalize 首字母大写
func Capitalize(str string) string {
	chars := []rune(str)
	if len(chars) == 0 {
		return ""
	}
	chars[0] = unicode.ToUpper(chars[0])
	return string(chars)
}

// Decapitalize 首字母小写
func Decapitalize(str string) string {
	chars := []rune(str)
	if len(chars) == 0 {
		return ""
	}
	chars[0] = unicode.ToLower(chars[0])
	return string(chars)
}

// Distance 计算编辑距离（Levenshtein距离，支持中文）
func Distance(str1, str2 string) int {
	r1 := []rune(str1)
	r2 := []rune(str2)
	m := len(r1)
	n := len(r2)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			dp[i][j] = min(
				dp[i-1][j]+1,
				dp[i][j-1]+1,
				dp[i-1][j-1]+cost,
			)
		}
	}

	return dp[m][n]
}

// min 辅助函数，返回较小值
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
