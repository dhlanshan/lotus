package strutil

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Capitalize Convert the first character of the string to uppercase and the remaining characters to lowercase.
func Capitalize(s string) string {
	if s == "" {
		return s
	}

	firstRune, size := utf8.DecodeRuneInString(s)
	if !unicode.IsLetter(firstRune) {
		return s
	}

	return string(unicode.ToUpper(firstRune)) + strings.ToLower(s[size:])
}

// Center the string and fill it with the specified padding character (default is a space) to the specified total width
func Center(s string, width int, fillChar rune) string {
	if fillChar == 0 {
		fillChar = ' '
	}

	length := utf8.RuneCountInString(s)
	if width <= length {
		return s
	}

	totalPadding := width - length
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	var builder strings.Builder
	builder.Grow(width)

	builder.WriteString(strings.Repeat(string(fillChar), leftPadding))
	builder.WriteString(s)
	builder.WriteString(strings.Repeat(string(fillChar), rightPadding))

	return builder.String()
}

// Count the number of occurrences of the substring in the parent string
func Count(s, sub string, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > utf8.RuneCountInString(s) || end < 0 {
		end = utf8.RuneCountInString(s)
	}

	if len(sub) == 0 {
		return 0
	}

	var count int
	runes := []rune(s)
	subRunes := []rune(sub)

	for i := start; i+len(subRunes) <= end; i++ {
		if string(runes[i:i+len(subRunes)]) == sub {
			count++
		}
	}

	return count
}

// EndsWith Check if the string ends with the specified suffix, can accept a single suffix or a tuple (multiple suffixes).
func EndsWith(s string, suffixes ...string) bool {
	if len(suffixes) == 0 {
		return false
	}

	for _, suffix := range suffixes {
		if suffix == "" {
			continue
		}

		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	return false
}

// ExpandTabs Replace the tab (\ t) in the string with an appropriate number of spaces to ensure alignment.
func ExpandTabs(s string, tabsize int) string {
	if tabsize <= 0 {
		tabsize = 8
	}

	var builder strings.Builder
	column := 0

	for _, char := range s {
		if char == '\t' {
			if tabsize == 1 {
				continue
			}

			spaceCount := tabsize - (column % tabsize)
			builder.WriteString(strings.Repeat(" ", spaceCount))
			column += spaceCount
		} else {
			builder.WriteRune(char)
			if char == '\n' || char == '\r' {
				column = 0
			} else {
				column += utf8.RuneLen(char)
			}
		}
	}

	return builder.String()
}

// Find the position of the substring in the parent string.
// Return the first matching position index of the substring, and if not found, return -1.
// Meanwhile, it supports optional start and end parameters to limit the search scope.
func Find(s, sub string, start, end int) int {
	// Convert to []rune for Unicode-safe indexing
	runes := []rune(s)
	subRunes := []rune(sub)
	runesLen := len(runes)
	subLen := len(subRunes)

	// Normalize start and end
	if start < 0 {
		start = 0
	}
	if end > runesLen || end < 0 {
		end = runesLen
	}

	// Edge cases
	if subLen == 0 {
		return start
	}
	if subLen > (end - start) {
		return -1
	}

	// Search within the range
	for i := start; i <= end-subLen; i++ {
		if string(runes[i:i+subLen]) == sub {
			return i
		}
	}

	return -1
}

// Formatter 用于存储命名参数
type Formatter map[string]interface{}

// Format String formatting
func Format(template string, args ...interface{}) (string, error) {
	// 缓存 Formatter 类型的参数
	formatters := make(map[string]interface{})
	for _, arg := range args {
		if formatter, ok := arg.(Formatter); ok {
			for key, value := range formatter {
				formatters[key] = value
			}
		}
	}

	var builder strings.Builder
	builder.Grow(len(template)) // 预先分配足够的空间
	argIndex := 0
	length := len(template)

	for i := 0; i < length; i++ {
		char := template[i]

		if char == '{' {
			// 转义处理 "{{"
			if i+1 < length && template[i+1] == '{' {
				builder.WriteByte('{')
				i++
				continue
			}

			// 查找 "}" 的位置
			end := i + 1
			for end < length && template[end] != '}' {
				end++
			}
			if end >= length {
				return "", fmt.Errorf("unmatched '{' in format string")
			}

			// 获取占位符
			placeholder := template[i+1 : end]
			var replacement string

			if placeholder == "" { // 位置参数
				if argIndex >= len(args) {
					return "", fmt.Errorf("missing positional argument at index %d", argIndex)
				}
				replacement = fmt.Sprintf("%v", args[argIndex])
				argIndex++
			} else { // 命名参数
				if value, exists := formatters[placeholder]; exists {
					replacement = fmt.Sprintf("%v", value)
				} else {
					return "", fmt.Errorf("missing named argument: %s", placeholder)
				}
			}

			builder.WriteString(replacement)
			i = end
		} else if char == '}' {
			// 转义处理 "}}"
			if i+1 < length && template[i+1] == '}' {
				builder.WriteByte('}')
				i++
			} else {
				return "", fmt.Errorf("unmatched '}' in format string")
			}
		} else {
			builder.WriteByte(char)
		}
	}

	return builder.String(), nil
}

// Index Return the index position where the substring first appears
func Index(s, sub string, start, end int) int {
	// 转换为 []rune 以处理多字节字符
	runes := []rune(s)
	subRunes := []rune(sub)
	runesLen := len(runes)
	subLen := len(subRunes)

	// 处理负索引并规范化
	start = normalizeIndex(start, runesLen)
	end = normalizeIndex(end, runesLen)

	// 确保满足半开区间的边界条件
	if start < 0 {
		start = 0
	}
	if end > runesLen {
		end = runesLen
	}
	if start >= end {
		return -1
	}

	// 空子字符串总是匹配，返回起始索引
	if subLen == 0 {
		return start
	}

	// 在指定范围内查找子字符串
	for i := start; i+subLen <= end; i++ {
		if string(runes[i:i+subLen]) == string(subRunes) {
			return i
		}
	}

	return -1
}

// normalizeIndex 将负索引转换为正索引
func normalizeIndex(index, length int) int {
	if index < 0 {
		index += length
	}
	return index
}

// Isalnum checks if all characters in the string are alphanumeric and there is at least one character.
func Isalnum(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// IsAlpha checks if the string contains only alphabetic characters.
func IsAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

// IsAscii checks if the string contains only ASCII characters.
func IsAscii(s string) bool {
	if s == "" {
		return true
	}
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

// IsDecimal checks if the string contains only decimal digits (0-9).
func IsDecimal(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// IsDigit checks if the string contains only digit characters (0-9, full-width digits, etc.).
func IsDigit(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsLower checks if the string contains only lowercase letters.
func IsLower(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

// IsNumeric checks if the string contains only numeric characters (including full-width digits, fractions, etc.).
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

// IsPrintable checks if the string contains only printable characters (including space).
func IsPrintable(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// IsSpace checks if the string contains only whitespace characters.
func IsSpace(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

// IsTitle checks if the string follows title case rules.
func IsTitle(s string) bool {
	if s == "" {
		return false
	}

	newWord := true

	for _, r := range s {
		if unicode.IsLetter(r) {
			if newWord {
				if !unicode.IsUpper(r) {
					return false
				}
				newWord = false
			} else {
				if !unicode.IsLower(r) {
					return false
				}
			}
		} else {
			newWord = true
		}
	}

	return true
}

// IsUpper checks if all the alphabetic characters in the string are uppercase and there is at least one alphabetic character.
func IsUpper(s string) bool {
	hasLetter := false

	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}

	return hasLetter
}

// LJust returns the string left-justified in a string of length width.
// Padding is done using the specified fillChar, defaulting to a space.
func LJust(s string, width int, fillChar ...rune) (string, error) {
	padChar := ' '
	if len(fillChar) > 0 {
		if len(fillChar) != 1 {
			return "", errors.New("fillChar must be a single character")
		}
		padChar = fillChar[0]
	}

	padding := width - len(s)
	if padding <= 0 {
		return s, nil
	}

	return s + strings.Repeat(string(padChar), padding), nil
}

// Lower converts all uppercase characters in the string to lowercase.
func Lower(s string) string {
	return strings.ToLower(s)
}

// LStrip removes leading characters from the string. If chars is not provided, it defaults to whitespace.
func LStrip(s string, chars string) string {
	charSet := " \t\n\r\f\v"
	if chars != "" {
		charSet = chars
	}

	charMap := make(map[rune]struct{})
	for _, r := range charSet {
		charMap[r] = struct{}{}
	}

	start := len(s)
	for i, r := range s {
		if _, found := charMap[r]; !found {
			start = i
			break
		}
	}

	return s[start:]
}

// Partition splits the string into three parts using the first occurrence of sep.
// If sep is not found, it returns the string and two empty strings.
func Partition(s, sep string) (string, string, string, error) {
	if sep == "" {
		return "", "", "", errors.New("sep cannot be empty")
	}

	index := strings.Index(s, sep)
	if index == -1 {
		return s, "", "", nil
	}

	return s[:index], sep, s[index+len(sep):], nil
}

// RemovePrefix removes the specified prefix from the string.
// If the string does not start with the prefix, it returns the original string.
func RemovePrefix(s, prefix string) string {
	if prefix == "" || len(s) < len(prefix) {
		return s
	}

	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}

	return s
}

// RemoveSuffix removes the specified suffix from the string.
// If the string does not end with the suffix, it returns the original string.
func RemoveSuffix(s, suffix string) string {
	if suffix == "" || len(s) < len(suffix) {
		return s
	}

	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}

	return s
}

// Replace replaces occurrences of old with new in the string.
// If count is -1, all occurrences are replaced. If count >= 0, only the first `count` occurrences are replaced.
func Replace(s, old, new string, count int) string {
	if old == "" {
		return s // If old is empty, return the original string
	}
	if count == 0 {
		return s // If count is 0, return the original string
	}

	if count < 0 {
		// Replace all occurrences using strings.ReplaceAll
		return strings.ReplaceAll(s, old, new)
	}

	// Replace the first `count` occurrences
	return strings.Replace(s, old, new, count)
}

// RFind returns the highest index in the string where substring `sub` is found,
// such that `sub` is contained within s[start:end]. Returns -1 if `sub` is not found.
func RFind(s, sub string, start, end int) int {
	if sub == "" {
		if end > len(s) {
			end = len(s)
		}
		return end
	}

	// Normalize start and end
	if start < 0 {
		start = len(s) + start
	}
	if end < 0 {
		end = len(s) + end
	}
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return -1
	}

	lastIndex := strings.LastIndex(s[start:end], sub)
	if lastIndex == -1 {
		return -1
	}

	return start + lastIndex
}

// RIndex returns the highest index in the string where substring `sub` is found,
// such that `sub` is contained within s[start:end]. If `sub` is not found, it panics.
func RIndex(s, sub string, start, end int) (int, error) {
	index := RFind(s, sub, start, end) // Reuse RFind
	if index == -1 {
		return -1, errors.New("substring not found")
	}
	return index, nil
}

// RJust returns the string right-aligned in a string of length `width`.
// Padding is done using the specified `fillChar` (default is space).
func RJust(s string, width int, fillChar ...rune) string {
	if width <= len(s) {
		return s
	}

	fill := ' '
	if len(fillChar) > 0 {
		fill = fillChar[0]
	}

	padding := width - len(s)

	return strings.Repeat(string(fill), padding) + s
}

// RPartition splits the string into three parts using the last occurrence of `sep`.
// If `sep` is not found, it returns ("", "", s).
func RPartition(s, sep string) (string, string, string) {
	if sep == "" {
		return "", "", s
	}

	index := strings.LastIndex(s, sep)
	if index == -1 {
		return "", "", s
	}

	before := s[:index]
	after := s[index+len(sep):]
	return before, sep, after
}

// RSplit splits the string by the specified separator `sep` starting from the right.
// If `sep` is None (empty string), it splits on whitespace.
// maxsplit determines the maximum number of splits. Default is -1 (no limit).
func RSplit(s, sep string, maxsplit int) []string {
	if sep == "" {
		return rsplitWhitespace(s, maxsplit)
	}

	if maxsplit == 0 {
		return []string{s} // No splitting needed
	}

	if len(sep) == 1 {
		return rsplitSingleChar(s, rune(sep[0]), maxsplit)
	}

	// Use sep for splitting
	parts := strings.Split(s, sep)
	if maxsplit < 0 || maxsplit >= len(parts)-1 {
		return parts
	}

	// Combine parts to enforce maxsplit
	remaining := parts[:len(parts)-maxsplit]
	last := strings.Join(remaining, sep)
	return append([]string{last}, parts[len(parts)-maxsplit:]...)
}

// rsplitWhitespace splits a string based on whitespace from the right.
func rsplitWhitespace(s string, maxsplit int) []string {
	fields := strings.Fields(s) // Splits on all whitespace
	if maxsplit < 0 || maxsplit >= len(fields)-1 {
		return fields
	}

	// Combine parts to enforce maxsplit
	remaining := fields[:len(fields)-maxsplit]
	last := strings.Join(remaining, " ")
	return append([]string{last}, fields[len(fields)-maxsplit:]...)
}

// rsplitSingleChar splits the string by a single character separator starting from the right.
func rsplitSingleChar(s string, sep rune, maxsplit int) []string {
	// If the separator is a single character, use LastIndex for better performance
	parts := []string{}
	start := len(s)
	for count := 0; count < maxsplit && start > 0; {
		index := strings.LastIndexByte(s[:start], byte(sep))
		if index == -1 {
			break
		}
		parts = append([]string{s[index+1 : start]}, parts...)
		start = index
		count++
	}
	parts = append([]string{s[:start]}, parts...)
	return parts
}

// RStrip removes trailing characters specified in `chars` from the string `s`.
// If `chars` is not provided, it defaults to removing whitespace.
func RStrip(s string, chars ...string) string {
	var charSet map[rune]bool

	// Default to whitespace characters if `chars` is not provided
	if len(chars) == 0 {
		return rstripWhitespace(s)
	}

	// Create a set of characters from `chars`
	charSet = make(map[rune]bool)
	for _, r := range chars[0] {
		charSet[r] = true
	}

	// Find the index to slice the string
	end := len(s)
	for i := len(s) - 1; i >= 0; i-- {
		if !charSet[rune(s[i])] {
			break
		}
		end = i
	}
	return s[:end]
}

// rstripWhitespace removes trailing whitespace characters from `s`.
func rstripWhitespace(s string) string {
	end := len(s)
	for i := len(s) - 1; i >= 0; i-- {
		if !unicode.IsSpace(rune(s[i])) {
			break
		}
		end = i
	}
	return s[:end]
}

// Split splits the string `s` by the specified separator `sep`.
// If `sep` is `None`, it splits by whitespace. The `maxsplit` limits the number of splits.
func Split(s, sep string, maxsplit int) []string {
	if sep == "" {
		// Default to whitespace splitting if `sep` is empty (None in Python)
		return splitWhitespace(s, maxsplit)
	}
	// Split by the specified separator
	return splitBySeparator(s, sep, maxsplit)
}

// splitWhitespace splits the string `s` by whitespace characters.
func splitWhitespace(s string, maxsplit int) []string {
	// Using strings.Fields to split by whitespace
	fields := strings.Fields(s)
	if maxsplit < 0 || maxsplit >= len(fields)-1 {
		return fields
	}

	// Apply maxsplit if needed
	return append(fields[:len(fields)-maxsplit], strings.Join(fields[len(fields)-maxsplit:], " "))
}

// splitBySeparator splits the string `s` by the specified separator `sep` and applies `maxsplit`.
func splitBySeparator(s, sep string, maxsplit int) []string {
	// Split by separator using strings.Split
	parts := strings.Split(s, sep)
	if maxsplit < 0 || maxsplit >= len(parts)-1 {
		return parts
	}

	// Apply maxsplit if needed
	return append(parts[:len(parts)-maxsplit], strings.Join(parts[len(parts)-maxsplit:], sep))
}

// SplitLines splits a string into lines based on universal newline characters.
// If `keepends` is true, the line-ending characters are retained.
func SplitLines(s string, keepends bool) []string {
	var result []string
	var lineStart int

	for i, r := range s {
		if isNewline(r) {
			// Include the newline character(s) if keepends is true
			if keepends {
				result = append(result, s[lineStart:i+1])
			} else {
				result = append(result, s[lineStart:i])
			}
			// Handle CRLF ("\r\n")
			if r == '\r' && i+1 < len(s) && s[i+1] == '\n' {
				lineStart = i + 2
			} else {
				lineStart = i + 1
			}
		}
	}

	// Append the last line if there's no trailing newline
	if lineStart < len(s) {
		result = append(result, s[lineStart:])
	}
	return result
}

// isNewline checks if a rune is a universal newline character.
func isNewline(r rune) bool {
	switch r {
	case '\n', '\r', '\v', '\f', '\x1c', '\x1d', '\x1e', '\x85', '\u2028', '\u2029':
		return true
	}
	return false
}

// StartsWith checks if the string `s` starts with the prefix `prefix`.
// Optional `start` and `end` parameters define the substring to check.
func StartsWith(s, prefix string, start, end int) bool {
	// Handle negative start and end indices
	if start < 0 {
		start = len(s) + start
	}
	if end < 0 {
		end = len(s) + end
	}

	// Clamp start and end to valid ranges
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		return false
	}

	// Extract the substring
	substring := s[start:end]

	// Check if the substring starts with the prefix
	return len(substring) >= len(prefix) && substring[:len(prefix)] == prefix
}

// Strip removes characters specified in `chars` from both ends of the string `s`.
// If `chars` is empty, it removes whitespace characters by default.
func Strip(s string, chars string) string {
	if chars == "" {
		chars = " \t\n\r\v\f"
	}

	start := 0
	end := len(s)

	// Find the first character not in `chars` from the left
	for start < end && strings.ContainsRune(chars, rune(s[start])) {
		start++
	}

	// Find the first character not in `chars` from the right
	for end > start && strings.ContainsRune(chars, rune(s[end-1])) {
		end--
	}

	return s[start:end]
}

// SwapCase returns a copy of the string `s` with uppercase letters converted to lowercase
// and lowercase letters converted to uppercase.
func SwapCase(s string) string {
	// Create a rune slice to efficiently build the transformed string
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			runes[i] = unicode.ToLower(r)
		} else if unicode.IsLower(r) {
			runes[i] = unicode.ToUpper(r)
		}
	}
	return string(runes)
}

// Title returns a copy of the string `s` where the first letter of each word is capitalized
// and the rest of the letters are in lowercase.
func Title(s string) string {
	var result strings.Builder
	previousIsSpace := true

	for _, r := range s {
		if previousIsSpace && unicode.IsLetter(r) {
			result.WriteRune(unicode.ToUpper(r))
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
		previousIsSpace = unicode.IsSpace(r)
	}

	return result.String()
}

// Translate replaces characters in the string `s` based on the mapping `table`.
// If the value for a rune is `-1`, the character is removed.
func Translate(s string, table map[rune]rune) string {
	var result strings.Builder

	for _, r := range s {
		if replacement, ok := table[r]; ok {
			if replacement != -1 { // -1 indicates deletion
				result.WriteRune(replacement)
			}
		} else {
			result.WriteRune(r) // Keep the original character
		}
	}

	return result.String()
}

// Upper returns a copy of the string `s` with all lowercase letters converted to uppercase.
func Upper(s string) string {
	return strings.ToUpper(s)
}

// ZFill pads the string `s` with zeros ('0') on the left, to make its length equal to `width`.
// If `s` starts with '+' or '-', the padding zeros are inserted after the sign.
func ZFill(s string, width int) string {
	length := len(s)
	if length >= width {
		return s // No padding needed
	}

	padCount := width - length
	var result strings.Builder

	// Handle sign at the beginning
	if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
		result.WriteByte(s[0]) // Write the sign
		result.WriteString(strings.Repeat("0", padCount))
		result.WriteString(s[1:]) // Append the rest of the string
	} else {
		// No sign, pad with zeros at the front
		result.WriteString(strings.Repeat("0", padCount))
		result.WriteString(s)
	}

	return result.String()
}
