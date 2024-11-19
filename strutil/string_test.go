package strutil

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// 常规情况
		{"hello", "Hello"},
		{"HELLO", "Hello"},
		{"hELLo", "Hello"},
		{"world", "World"},
		// 空字符串
		{"", ""},
		// 单字符
		{"a", "A"},
		{"A", "A"},
		// 非字母字符开头
		{"123abc", "123abc"},
		{"!hello", "!hello"},
		// Unicode 字符
		{"éxample", "Éxample"},
		{"ÉXAMPLE", "Éxample"},
		{"你好", "你好"}, // 不变，因为首字符是中文
		// 混合大小写和特殊字符
		{"gO-lang", "Go-lang"},
		{"🚀rocket", "🚀rocket"}, // Emoji 不变
	}

	for _, test := range tests {
		result := Capitalize(test.input)
		if result != test.expected {
			t.Errorf("Capitalize(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func BenchmarkCapitalize(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{"ShortASCII", "hello"},
		{"MixedCaseASCII", "GoLaNg"},
		{"LongASCII", strings.Repeat("hello ", 1000)}, // 长字符串
		{"UnicodeShort", "éxample"},
		{"UnicodeLong", strings.Repeat("你好世界", 1000)}, // 长 Unicode 字符串
		{"SpecialChars", "123!@#Hello"},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Capitalize(test.input)
			}
		})
	}
}

func TestCenter(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		fillchar rune
		expected string
	}{
		// 常规情况
		{"hello", 10, '*', "**hello***"},
		{"hello", 11, '-', "---hello---"},
		{"hello", 5, '#', "hello"}, // 宽度小于等于字符串长度
		// Unicode 字符
		{"你好", 6, '-', "--你好--"},
		{"你好", 7, '*', "**你好***"},
		// 填充字符为 0（使用默认空格）
		{"world", 10, 0, "  world   "},
		// 空字符串
		{"", 5, '+', "+++++"},
		{"", 0, '!', ""},
		// 边界情况：宽度小于零
		{"test", -1, '*', "test"},
		// 特殊字符
		{"123", 7, '=', "==123=="},
		{"🚀", 5, '-', "--🚀--"},
	}

	for _, test := range tests {
		result := Center(test.input, test.width, test.fillchar)
		if result != test.expected {
			t.Errorf("Center(%q, %d, %q) = %q; want %q", test.input, test.width, test.fillchar, result, test.expected)
		}
	}
}

func BenchmarkCenter(b *testing.B) {
	tests := []struct {
		name     string
		input    string
		width    int
		fillchar rune
	}{
		{"ShortASCII", "hello", 20, '*'},
		{"LongASCII", "hello", 1000, '-'},
		{"UnicodeShort", "你好", 20, ' '},
		{"UnicodeLong", "你好", 1000, '='},
		{"EmptyString", "", 50, '+'},
		{"SpecialChars", "🚀123", 50, '#'},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Center(test.input, test.width, test.fillchar)
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		input    string
		sub      string
		start    int
		end      int
		expected int
	}{
		// 常规测试
		{"hello world", "o", 0, 11, 2},                 // "o" 在 "hello world" 中出现两次
		{"hello world", "l", 0, 11, 3},                 // "l" 在 "hello world" 中出现三次
		{"hello world", "z", 0, 11, 0},                 // "z" 不在 "hello world" 中
		{"hello world", "world", 0, 11, 1},             // "world" 出现一次
		{"", "o", 0, 0, 0},                             // 空字符串，没有匹配
		{"hello", "", 0, 5, 0},                         // 空子串，返回0次
		{"hello world hello world", "hello", 0, 23, 2}, // "hello" 出现两次

		// Unicode 字符串
		{"你好，世界你好，世界", "世界", 0, 18, 2}, // "世界" 在字符串中出现两次

		// 边界情况
		{"你好，世界你好，世界", "你好", 0, -1, 2},  // end 为负数，表示整个字符串
		{"你好，世界你好，世界", "世界", 10, 18, 0}, // 只查找子串 "世界" 在字符串的部分
	}

	for _, test := range tests {
		result := Count(test.input, test.sub, test.start, test.end)
		if result != test.expected {
			t.Errorf("Count(%q, %q, %d, %d) = %d; want %d", test.input, test.sub, test.start, test.end, result, test.expected)
		}
	}
}

func BenchmarkCount(b *testing.B) {
	tests := []struct {
		name     string
		input    string
		sub      string
		start    int
		end      int
		expected int
	}{
		{"ShortString", "hello world", "o", 0, 11, 2},
		{"LongString", "hello world hello world hello world hello world", "hello", 0, 50, 4},
		{"UnicodeString", "你好，世界你好，世界你好，世界", "你好", 0, 18, 3},
		{"EmptyString", "", "a", 0, 0, 0},
		{"SpecialChars", "123$%$^&*123$%$^&*123", "$%$^", 0, 30, 3},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Count(test.input, test.sub, test.start, test.end)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	tests := []struct {
		input    string
		suffixes []string
		expected bool
	}{
		{"hello world", []string{"world"}, true},
		{"hello world", []string{"world", "hello"}, true},
		{"hello world", []string{"goodbye"}, false},
		{"你好，世界", []string{"世界"}, true},
		{"你好，世界", []string{"你好"}, false}, // 不能匹配
		{"", []string{"world"}, false},   // 空字符串，不能匹配任何后缀
		{"hello", []string{""}, false},   // 空后缀不匹配
		{"hello🚀", []string{"🚀"}, true},  // 空后缀不匹配
	}

	for _, test := range tests {
		result := EndsWith(test.input, test.suffixes...)
		if result != test.expected {
			t.Errorf("EndsWith(%q, %v) = %v; want %v", test.input, test.suffixes, result, test.expected)
		}
	}
}

func BenchmarkEndsWith(b *testing.B) {
	tests := []struct {
		name     string
		input    string
		suffixes []string
	}{
		{"ShortString", "hello world", []string{"world"}},
		{"LongString", "hello world hello world hello world", []string{"world", "hello"}},
		{"UnicodeString", "你好，世界你好，世界你好，世界", []string{"世界", "你好"}},
		{"EmptyString", "", []string{"world"}},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				EndsWith(test.input, test.suffixes...)
			}
		})
	}
}

func TestExpandTabs(t *testing.T) {
	tests := []struct {
		input    string
		tabsize  int
		expected string
	}{
		{"hello\tworld", 4, "hello   world"},
		{"\t\thello", 4, "        hello"},
		{"hello\tworld\nhi\tthere", 4, "hello   world\nhi  there"},
		{"", 4, ""},
		{"\t", 4, "    "},
		{"abc\tdef", 0, "abc     def"},    // 默认 tabsize = 8
		{"abc\tdef", -5, "abc     def"},   // 负数处理为默认值
		{"你好\t世界", 4, "你好  世界"},           // Unicode 测试
		{"hello\tworld", 1, "helloworld"}, // tabsize = 1, 无空格直接对齐
	}

	for _, test := range tests {
		result := ExpandTabs(test.input, test.tabsize)
		if result != test.expected {
			t.Errorf("ExpandTabs(%q, %d) = %q; want %q", test.input, test.tabsize, result, test.expected)
		}
	}
}

func BenchmarkExpandTabs(b *testing.B) {
	tests := []struct {
		input   string
		tabsize int
	}{
		// 短字符串，tabsize = 4
		{"hello\tworld", 4},
		// 包含多制表符
		{"hello\tworld\tGo\tis\tawesome", 4},
		// 中文字符与制表符混合
		{"你好\t世界\t欢迎\t光临", 4},
		// 特殊情况：tabsize = 1（制表符被跳过）
		{"hello\tworld\tGo\tis\tawesome", 1},
		// 特殊情况：tabsize = 0（默认 8）
		{"hello\tworld\tGo\tis\tawesome", 0},
	}

	for _, test := range tests {
		b.Run(test.input[:10]+"...", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ExpandTabs(test.input, test.tabsize)
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		input    string
		sub      string
		start    int
		end      int
		expected int
	}{
		{"hello world", "world", 0, -1, 6},  // Basic match
		{"hello world", "hello", 0, -1, 0},  // Match at the start
		{"hello world", "o", 0, -1, 4},      // Match single character
		{"hello world", "z", 0, -1, -1},     // No match
		{"hello world", "world", 7, -1, -1}, // No match in restricted range
		{"hello world", "world", 0, 5, -1},  // No match in restricted range
		{"你好世界", "世", 0, -1, 2},             // Unicode match
		{"你好世界", "你好", 0, -1, 0},            // Unicode match at the start
		{"你好世界", "界", 1, 3, -1},             // Unicode match in range
		{"hello world", "", 0, -1, 0},       // Empty substring
		{"hello world", "world", -5, 20, 6}, // Start and end normalization
		{"hello world", "world", 6, -1, 6},  // Start matches
		{"hello", "hello world", 0, -1, -1}, // Substring longer than input
	}

	for _, test := range tests {
		t.Run(test.input+"_"+test.sub, func(t *testing.T) {
			result := Find(test.input, test.sub, test.start, test.end)
			if result != test.expected {
				t.Errorf("Find(%q, %q, %d, %d) = %d; want %d",
					test.input, test.sub, test.start, test.end, result, test.expected)
			}
		})
	}
}

func BenchmarkFind(b *testing.B) {
	tests := []struct {
		input string
		sub   string
		start int
		end   int
	}{
		{"hello world", "world", 0, -1},
		{"hello world", "o", 0, -1},
		{"你好世界", "世", 0, -1},
		{"你好世界", "界", 1, 3},
		{strings.Repeat("a", 10000) + "b", "b", 0, -1}, // Large input, match at the end
		{strings.Repeat("a", 10000), "b", 0, -1},       // Large input, no match
	}

	for _, test := range tests {
		b.Run(test.input[:10]+"_"+test.sub, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Find(test.input, test.sub, test.start, test.end)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		template string
		args     []interface{}
		expected string
		err      bool
	}{
		// 测试位置参数
		{"Hello, {}!", []interface{}{"World"}, "Hello, World!", false},
		{"{} + {} = {}", []interface{}{1, 2, 3}, "1 + 2 = 3", false},
		{"Empty: {}", []interface{}{}, "", true},

		// 测试命名参数
		{"Hello, {name}!", []interface{}{Formatter{"name": "Alice"}}, "Hello, Alice!", false},
		{"Missing {key}", []interface{}{Formatter{"other": 123}}, "", true},

		// 测试转义符
		{"{{}}", nil, "{}", false},
		{"{{name}}", nil, "{name}", false},
		{"{{Hello}}", nil, "{Hello}", false},

		// 测试混合参数 (不支持)
		//{"{greet}, {name}!", []interface{}{"Hi", Formatter{"name": "Alice"}}, "Hi, Alice!", false},

		// 边界测试
		{"No placeholders", nil, "No placeholders", false},
		{"Unmatched {", nil, "", true},
		{"Unmatched }", nil, "", true},
	}

	for _, tt := range tests {
		result, err := Format(tt.template, tt.args...)
		if (err != nil) != tt.err {
			t.Errorf("Format(%q, %v) error = %v, wantErr %v", tt.template, tt.args, err, tt.err)
		}
		if result != tt.expected {
			t.Errorf("Format(%q, %v) = %q, want %q", tt.template, tt.args, result, tt.expected)
		}
	}
}

func BenchmarkFormat(b *testing.B) {
	// 示例模板和参数
	template := "Hello, {name}! Today is {day}. You have {count} new messages."
	args := []interface{}{Formatter{
		"name":  "Alice",
		"day":   "Monday",
		"count": 5,
	}}

	for i := 0; i < b.N; i++ {
		_, err := Format(template, args...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestIndex(t *testing.T) {
	tests := []struct {
		s        string
		sub      string
		start    int
		end      int
		expected int
	}{
		{"hello world", "world", 0, len("hello world"), 6},
		{"hello world", "world", 0, 5, -1},
		{"hello world", "o", 4, 10, 4},
		{"你好，世界", "世界", 0, -1, -1},
		{"你好，世界", "好", -4, -1, 1},
		{"hello world", "world", -5, -1, -1},
		{"hello world", "", 3, 8, 3},
		{"hello world", "z", 0, -1, -1},
	}

	for _, tt := range tests {
		result := Index(tt.s, tt.sub, tt.start, tt.end)
		if result != tt.expected {
			t.Errorf("Index(%q, %q, %d, %d) = %d; want %d",
				tt.s, tt.sub, tt.start, tt.end, result, tt.expected)
		}
	}
}

func BenchmarkIndex(b *testing.B) {
	tests := []struct {
		s     string
		sub   string
		start int
		end   int
	}{
		{"hello world", "world", 0, len("hello world")}, // 普通字符串查找
		{"hello world", "hello", 0, len("hello world")}, // 字符串开头查找
		{"hello world", "world", 0, 5},                  // 查找范围限制
		{"你好，世界", "世界", 0, len("你好，世界")},                // 中文字符查找
		{"你好，世界", "好", -4, -1},                          // 使用负索引查找
		{"", "hello", 0, 0},                             // 空字符串
	}

	for _, tt := range tests {
		// 使用 b.N 来控制基准测试的循环次数
		b.Run(tt.s, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// 测试每种情况
				Index(tt.s, tt.sub, tt.start, tt.end)
			}
		})
	}
}

func TestIsalnum(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"1234567890", true},      // 全数字
		{"abcdefg", true},         // 全字母
		{"123abc456", true},       // 字母数字混合
		{"HelloWorld123", true},   // 字母数字混合，带大写
		{"Hello@World123", false}, // 含非字母数字字符
		{"你好世界", true},            // 中文字符
		{"", false},               // 空字符串
		{"123!@#", false},         // 含有特殊字符
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result := Isalnum(tt.s)
			if result != tt.expected {
				t.Errorf("Isalnum(%q) = %v; want %v", tt.s, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsalnum(b *testing.B) {
	tests := []string{
		"1234567890",     // 全数字
		"abcdefg",        // 全字母
		"123abc456",      // 字母数字混合
		"HelloWorld123",  // 字母数字混合，带大写
		"Hello@World123", // 含非字母数字字符
		"你好世界",           // 中文字符
	}

	for _, test := range tests {
		b.Run(test, func(b *testing.B) {
			// 重复执行多次以获取准确的基准值
			for i := 0; i < b.N; i++ {
				Isalnum(test)
			}
		})
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"hello", true},        // 只包含字母
		{"hello123", false},    // 包含数字
		{"你好", true},           // Unicode 字符，中文
		{"hello world", false}, // 包含空格
		{"", false},            // 空字符串
		{"12345", false},       // 只包含数字
		{"@#$%", false},        // 包含特殊字符
		{"ABCDE", true},        // 只包含大写字母
		{"abcdef", true},       // 只包含小写字母
		{"ABcdEf", true},       // 混合大小写字母
		{"你好，世界", false},       // 包含中文和标点符号，应该是 false
		{"hello_123", false},   // 包含下划线和数字
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result := IsAlpha(tt.s)
			if result != tt.expected {
				t.Errorf("IsAlpha(%q) = %v; want %v", tt.s, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsAlpha(b *testing.B) {
	tests := []struct {
		s string
	}{
		{"hello"},
		{"hello123"},
		{"你好"},
		{"hello world"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"},
		{"12345"},
		{"@#$%"},
		{"!@#abcABC"},
	}

	for _, tt := range tests {
		b.Run(tt.s, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IsAlpha(tt.s)
			}
		})
	}
}

func TestIsAscii(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"hello", true},
		{"你好", false}, // 非 ASCII 字符
		{"world!", true},
		{"\x80", false}, // 非 ASCII 字符
		{"", true},      // 空字符串也认为是 ASCII
		{"Hello123", true},
		{"你好，世界", false}, // 非 ASCII 字符
	}

	for _, tt := range tests {
		result := IsAscii(tt.s)
		if result != tt.expected {
			t.Errorf("IsAscii(%q) = %v; want %v", tt.s, result, tt.expected)
		}
	}
}

func BenchmarkIsAscii(b *testing.B) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},           // 简单的ASCII字符串
		{"hello 世界", false},       // 包含非ASCII字符的字符串
		{"", true},                // 空字符串
		{"ASCII123", true},        // 仅包含ASCII字符
		{"你好", false},             // 仅包含非ASCII字符
		{"hello world 123", true}, // 仅包含ASCII字符
	}

	for _, tt := range tests {
		b.Run(tt.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := IsAscii(tt.input)
				if result != tt.expected {
					b.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestIsDecimal(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"12345", true},      // 纯数字
		{"12345abc", false},  // 包含非数字字符
		{"", false},          // 空字符串
		{"9876543210", true}, // 纯数字
		{"1234 5678", false}, // 包含空格
		{"0123456789", true}, // 纯数字，含0开头
		{"-12345", false},    // 包含负号
		{"123.45", false},    // 包含小数点
		{"1234567890", true}, // 纯数字
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsDecimal(tt.input)
			if result != tt.expected {
				t.Errorf("IsDecimal(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsDecimal(b *testing.B) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"12345", true},      // 纯数字
		{"12345abc", false},  // 包含非数字字符
		{"", false},          // 空字符串
		{"9876543210", true}, // 纯数字
		{"1234 5678", false}, // 包含空格
		{"0123456789", true}, // 纯数字，含0开头
	}

	for _, tt := range tests {
		b.Run(tt.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := IsDecimal(tt.input)
				if result != tt.expected {
					b.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"12345", true},      // 纯数字
		{"12345abc", false},  // 包含非数字字符
		{"", false},          // 空字符串
		{"9876543210", true}, // 纯数字
		{"１２３４５", true},      // 全角数字
		{"1234 5678", false}, // 包含空格
		{"-12345", false},    // 包含负号
		{"123.45", false},    // 包含小数点
		{"٠١٢٣٤٥٦٧٨٩", true}, // 阿拉伯数字
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsDigit(tt.input)
			if result != tt.expected {
				t.Errorf("IsDigit(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsDigit(b *testing.B) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"12345", true},      // 纯数字
		{"12345abc", false},  // 包含非数字字符
		{"", false},          // 空字符串
		{"9876543210", true}, // 纯数字
		{"１２３４５", true},      // 全角数字
		{"1234 5678", false}, // 包含空格
		{"-12345", false},    // 包含负号
		{"123.45", false},    // 包含小数点
		{"٠١٢٣٤٥٦٧٨٩", true}, // 阿拉伯数字
	}

	for _, tt := range tests {
		b.Run(tt.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := IsDigit(tt.input)
				if result != tt.expected {
					b.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestIsLower(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},         // 纯小写字母
		{"HELLO", false},        // 纯大写字母
		{"Hello", false},        // 混合大小写字母
		{"", false},             // 空字符串
		{"hello123", false},     // 包含数字
		{"hello!", false},       // 包含符号
		{"你好", false},           // 包含中文
		{"lowercase", true},     // 纯小写字母
		{"lowercase123", false}, // 纯小写字母 + 数字
		{"lower_case", false},   // 包含下划线
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsLower(tt.input)
			if result != tt.expected {
				t.Errorf("IsLower(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsLower(b *testing.B) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},         // 纯小写字母
		{"HELLO", false},        // 纯大写字母
		{"Hello", false},        // 混合大小写字母
		{"", false},             // 空字符串
		{"hello123", false},     // 包含数字
		{"hello!", false},       // 包含符号
		{"你好", false},           // 包含中文
		{"lowercase", true},     // 纯小写字母
		{"lowercase123", false}, // 纯小写字母 + 数字
		{"lower_case", false},   // 包含下划线
	}

	for _, tt := range tests {
		b.Run(tt.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := IsLower(tt.input)
				if result != tt.expected {
					b.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123456", true},     // 纯数字字符
		{"１２３４５６", true},     // 全角数字字符
		{"123.45", false},    // 包含小数点
		{"你好", false},        // 中文字符
		{"123456abc", false}, // 包含字母
		{"", false},          // 空字符串
		{"⅔", true},          // 分数字符
		{"٠١٢٣٤٥٦", true},    // 阿拉伯数字
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("IsNumeric(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
func BenchmarkIsNumeric(b *testing.B) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123456", true},     // 纯数字字符
		{"１２３４５６", true},     // 全角数字字符
		{"123.45", false},    // 包含小数点
		{"你好", false},        // 中文字符
		{"123456abc", false}, // 包含字母
		{"", false},          // 空字符串
		{"⅔", true},          // 分数字符
		{"٠١٢٣٤٥٦", true},    // 阿拉伯数字
	}

	for _, tt := range tests {
		b.Run(tt.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := IsNumeric(tt.input)
				if result != tt.expected {
					b.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello world", true},   // 纯可打印字符
		{"hello\nworld", false}, // 包含换行符
		{"hello\tworld", false}, // 包含制表符
		{"你好，世界", true},         // 中文字符
		{"hello世界", true},       // 混合字符
		{"", false},             // 空字符串
		{" ", true},             // 空格是可打印字符
		{"\x01\x02\x03", false}, // 非打印字符
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsPrintable(tt.input)
			if result != tt.expected {
				t.Errorf("IsPrintable(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsPrintable(b *testing.B) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello world", true},   // 纯可打印字符
		{"hello\nworld", false}, // 包含换行符
		{"hello\tworld", false}, // 包含制表符
		{"你好，世界", true},         // 中文字符
		{"hello世界", true},       // 混合字符
		{"", false},             // 空字符串
		{" ", true},             // 空格是可打印字符
		{"\x01\x02\x03", false}, // 非打印字符
	}

	for _, tt := range tests {
		b.Run(tt.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := IsPrintable(tt.input)
				if result != tt.expected {
					b.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestIsSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"    ", true},          // 纯空白字符
		{"\t\n ", true},         // 含有空格、换行、制表符
		{"hello", false},        // 包含非空白字符
		{"\t\n", true},          // 仅有制表符和换行符
		{"", false},             // 空字符串
		{" hello ", false},      // 包含空格，但非全空白
		{"\x01\x02\x03", false}, // 非空白字符
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsSpace(tt.input)
			if result != tt.expected {
				t.Errorf("IsSpace(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsSpace(b *testing.B) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"    ", true},          // 纯空白字符
		{"\t\n ", true},         // 含有空格、换行、制表符
		{"hello", false},        // 包含非空白字符
		{"\t\n", true},          // 仅有制表符和换行符
		{"", false},             // 空字符串
		{" hello ", false},      // 包含空格，但非全空白
		{"\x01\x02\x03", false}, // 非空白字符
	}

	for _, tt := range tests {
		b.Run(tt.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := IsSpace(tt.input)
				if result != tt.expected {
					b.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestIsTitle(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"The Quick Brown Fox", true},  // 正确的标题
		{"the quick brown fox", false}, // 错误的标题，首字母小写
		{"The Quick brown Fox", false}, // 错误的标题，"brown"不小写
		{"", false},                    // 空字符串
		{"Hello", true},                // 单个单词的标题
		{"Hello World", true},          // 多个单词的标题
		{"hello world", false},         // 错误的标题，单词首字母小写
		{"HELLO world", false},         // 错误的标题，首字母大写但其他字母大写
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result := IsTitle(tt.s)
			if result != tt.expected {
				t.Errorf("IsTitle(%q) = %v; want %v", tt.s, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsTitle(b *testing.B) {
	// 定义一些测试用例
	tests := []struct {
		s string
		o string
	}{
		{"The Quick Brown Fox", ""}, // 正确的标题
		{"the quick brown fox", ""}, // 错误的标题，首字母小写
		{"The Quick brown Fox", ""}, // 错误的标题，"brown"不小写
		{"Hello World", ""},         // 单个单词的标题
		{"HELLO world", ""},         // 错误的标题，首字母大写但其他字母大写
		{"", ""},                    // 空字符串
		{"A Quick Brown Fox Jumped Over The Lazy Dog", ""}, // 较长的标题
	}

	for _, tt := range tests {
		b.Run(tt.s, func(b *testing.B) {
			// 运行基准测试
			for i := 0; i < b.N; i++ {
				_ = IsTitle(tt.s) // 调用 IsTitle
			}
		})
	}
}

func TestIsUpper(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"HELLO", true},       // All uppercase letters
		{"HELLO WORLD", true}, // All uppercase letters with space
		{"hello", false},      // All lowercase letters
		{"Hello", false},      // Mixed case
		{"1234", false},       // No letters
		{"", false},           // Empty string
		{"HELLO123", true},    // Uppercase letters with digits
		{"123ABC", true},      // Digits before uppercase letters
		{"aBcD", false},       // Mixed case letters
		{"H E L L O", true},   // Uppercase letters with spaces
		{"!@#^&*", false},     // Non-alphabetic characters
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result := IsUpper(tt.s)
			if result != tt.expected {
				t.Errorf("IsUpper(%q) = %v; want %v", tt.s, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsUpper(b *testing.B) {
	tests := []struct {
		s string
	}{
		{"HELLO WORLD"},             // All uppercase letters
		{"hello world"},             // All lowercase letters
		{"Hello World"},             // Mixed case letters
		{"HELLO1234567890!@#^&*()"}, // Uppercase with non-alphabetic characters
		{""},                        // Empty string
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ"}, // Large uppercase
		{"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"}, // Large lowercase
	}

	for _, tt := range tests {
		b.Run(tt.s, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = IsUpper(tt.s)
			}
		})
	}
}

func TestLJust(t *testing.T) {
	tests := []struct {
		s        string
		width    int
		fillchar []rune
		expected string
		err      bool
	}{
		{"hello", 10, []rune{'*'}, "hello*****", false}, // 使用星号填充
		{"hello", 5, []rune{}, "hello", false},          // 不需要填充
		{"hi", 5, []rune{'-'}, "hi---", false},          // 使用短横线填充
		{"test", 3, []rune{}, "test", false},            // 字符串宽度小于实际长度，无需填充
		{"abc", 6, []rune{'#'}, "abc###", false},        // 使用 # 填充
		{"", 5, []rune{'$'}, "$$$$$", false},            // 空字符串填充
		{"hello", 10, []rune{'*', '+'}, "", true},       // 填充字符不合法
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result, err := LJust(tt.s, tt.width, tt.fillchar...)
			if (err != nil) != tt.err {
				t.Errorf("LJust(%q, %d, %q) = %v, want error? %v", tt.s, tt.width, tt.fillchar, err, tt.err)
			}
			if result != tt.expected {
				t.Errorf("LJust(%q, %d, %q) = %q; want %q", tt.s, tt.width, tt.fillchar, result, tt.expected)
			}
		})
	}
}

func BenchmarkLJust(b *testing.B) {
	tests := []struct {
		s     string
		width int
	}{
		{"hello", 100},
		{"short", 1000},
		{"this is a much longer string", 10000},
	}

	for _, test := range tests {
		b.Run(test.s, func(b *testing.B) {
			// 使用默认填充字符进行基准测试
			for i := 0; i < b.N; i++ {
				LJust(test.s, test.width)
			}
		})
	}
}

func TestLower(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{"hello", "hello"},
		{"HELLO", "hello"},
		{"HeLLo", "hello"},
		{"123", "123"},       // 纯数字不受影响
		{"", ""},             // 空字符串
		{"你好", "你好"},         // 非 ASCII 字符（Unicode）不受影响
		{"123!@#", "123!@#"}, // 非字母字符不受影响
		{"A B C", "a b c"},   // 字符间有空格，确保空格不受影响
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result := Lower(tt.s)
			if result != tt.expected {
				t.Errorf("Lower(%q) = %q; want %q", tt.s, result, tt.expected)
			}
		})
	}
}

func BenchmarkLower(b *testing.B) {
	tests := []struct {
		s string
	}{
		{"hello world"},
		{"HELLO WORLD"},
		{"a long string with many characters to test performance"},
		{"1234567890"},
		{"你好，世界"},
	}

	for _, test := range tests {
		b.Run(test.s, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Lower(test.s)
			}
		})
	}
}

func TestLStrip(t *testing.T) {
	tests := []struct {
		s        string
		chars    string
		expected string
	}{
		{"  hello", "", "hello"},
		{"\t\n hello", "", "hello"},
		{"abca123", "abc", "123"},
		{"###hello###", "#", "hello###"},
		{"  \t", "", ""},
		{"", "", ""},
		{"你好世界", "你", "好世界"},
		{"\t\n中文字符", "", "中文字符"},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result := LStrip(tt.s, tt.chars)
			if result != tt.expected {
				t.Errorf("LStrip(%q, %q) = %q; want %q", tt.s, tt.chars, result, tt.expected)
			}
		})
	}
}

func BenchmarkLStrip(b *testing.B) {
	// 定义测试数据
	tests := []struct {
		s     string
		chars string
	}{
		{"  \t\n  Hello, World!  \t", ""},   // 默认去除空白字符
		{"aaaaaHello, World!", "a"},         // 去除指定字符
		{"#######Hello, World!######", "#"}, // 去除重复的前导字符
		{"你好你好世界", "你"},                     // 去除中文字符
		{"\t\n\r\f\vHello, World!", ""},     // 特殊空白字符
		{"", ""},                            // 空字符串
		{"      ", ""},                      // 全空白字符
		{"Hello, World!", ""},               // 无需去除
	}

	// 逐个测试
	for _, tt := range tests {
		b.Run(strings.TrimSpace(tt.s), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				LStrip(tt.s, tt.chars)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		s        string
		sep      string
		expected [3]string
		err      error
	}{
		{"hello world", " ", [3]string{"hello", " ", "world"}, nil},
		{"hello world", "world", [3]string{"hello ", "world", ""}, nil},
		{"hello world", "hello", [3]string{"", "hello", " world"}, nil},
		{"hello world", "x", [3]string{"hello world", "", ""}, nil},
		{"hello world", "", [3]string{"", "", ""}, errors.New("sep cannot be empty")},
		{"", "x", [3]string{"", "", ""}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.sep, func(t *testing.T) {
			part1, part2, part3, err := Partition(tt.s, tt.sep)

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}

			if [3]string{part1, part2, part3} != tt.expected {
				t.Errorf("Partition(%q, %q) = (%q, %q, %q); want (%q, %q, %q)",
					tt.s, tt.sep, part1, part2, part3, tt.expected[0], tt.expected[1], tt.expected[2])
			}
		})
	}
}

func BenchmarkPartition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Partition("this is a benchmark test string", "benchmark")
	}
}

func TestRemovePrefix(t *testing.T) {
	tests := []struct {
		s        string
		prefix   string
		expected string
	}{
		{"hello world", "hello", " world"},
		{"hello world", "world", "hello world"},
		{"你好，世界", "你好，", "世界"},
		{"你好，世界", "世界", "你好，世界"},
		{"", "prefix", ""},
		{"prefix", "", "prefix"},
		{"你好你好", "你好", "你好"},
	}

	for _, tt := range tests {
		result := RemovePrefix(tt.s, tt.prefix)
		if result != tt.expected {
			t.Errorf("RemovePrefix(%q, %q) = %q; want %q",
				tt.s, tt.prefix, result, tt.expected)
		}
	}
}

func BenchmarkRemovePrefix(b *testing.B) {
	tests := []struct {
		s      string
		prefix string
	}{
		{"hello world", "hello"},
		{"hello world", "world"},
		{"你好，世界", "你好，"},
		{"你好，世界", "世界"},
		{"prefixprefixprefix", "prefix"},
		{"你好你好你好", "你好"},
	}

	for _, tt := range tests {
		b.Run(tt.s+"_"+tt.prefix, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = RemovePrefix(tt.s, tt.prefix)
			}
		})
	}
}

func TestRemoveSuffix(t *testing.T) {
	tests := []struct {
		s        string
		suffix   string
		expected string
	}{
		{"hello world", "world", "hello "},
		{"hello world", "hello", "hello world"},
		{"你好，世界", "世界", "你好，"},
		{"你好，世界", "你好", "你好，世界"},
		{"golang", "lang", "go"},
		{"golang", "python", "golang"},
		{"", "", ""},
		{"abc", "", "abc"},
		{"abc", "abc", ""},
		{"中文测试", "测试", "中文"},
		{"中文测试", "中文", "中文测试"},
	}

	for _, tt := range tests {
		result := RemoveSuffix(tt.s, tt.suffix)
		if result != tt.expected {
			t.Errorf("RemoveSuffix(%q, %q) = %q; want %q", tt.s, tt.suffix, result, tt.expected)
		}
	}
}

func BenchmarkRemoveSuffix(b *testing.B) {
	tests := []struct {
		s      string
		suffix string
	}{
		{"hello world", "world"},
		{"hello world", "hello"},
		{"你好，世界", "世界"},
		{"你好，世界", "你好"},
		{"golang programming", "programming"},
		{"golang programming", "golang"},
		{"abcabcabcabc", "abc"},
		{"中文字符串测试", "测试"},
	}

	for _, tt := range tests {
		b.Run(tt.s+"_"+tt.suffix, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = RemoveSuffix(tt.s, tt.suffix)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		s        string
		old      string
		new      string
		count    int
		expected string
	}{
		{"hello world", "world", "Go", -1, "hello Go"},
		{"hello world world", "world", "Go", 1, "hello Go world"},
		{"hello world world", "world", "Go", -1, "hello Go Go"},
		{"你好，世界，世界", "世界", "Go", 1, "你好，Go，世界"},
		{"你好，世界，世界", "世界", "Go", -1, "你好，Go，Go"},
		{"abcdabcdabcd", "abc", "123", 2, "123d123dabcd"},
		{"abcdabcdabcd", "abc", "123", -1, "123d123d123d"},
		{"", "abc", "123", -1, ""},
		{"abcd", "", "123", -1, "abcd"}, // old is empty
		{"abcd", "abcd", "", -1, ""},    // new is empty
	}

	for _, tt := range tests {
		result := Replace(tt.s, tt.old, tt.new, tt.count)
		if result != tt.expected {
			t.Errorf("Replace(%q, %q, %q, %d) = %q; want %q", tt.s, tt.old, tt.new, tt.count, result, tt.expected)
		}
	}
}

func BenchmarkReplace(b *testing.B) {
	tests := []struct {
		s     string
		old   string
		new   string
		count int
	}{
		{"hello world", "world", "Go", -1},
		{"hello world world", "world", "Go", 1},
		{"hello world world", "world", "Go", -1},
		{"你好，世界，世界", "世界", "Go", 1},
		{"abcdabcdabcd", "abc", "123", 2},
		{"abcdabcdabcd", "abc", "123", -1},
	}

	for _, tt := range tests {
		b.Run(tt.s+"_"+tt.old, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Replace(tt.s, tt.old, tt.new, tt.count)
			}
		})
	}
}

func TestRFind(t *testing.T) {
	tests := []struct {
		s, sub     string
		start, end int
		expected   int
	}{
		// 标准测试
		{"hello world", "world", 0, len("hello world"), len("hello world") - len("world")},
		{"hello world", "hello", 0, len("hello world"), 0},                  // "hello" 从头开始
		{"hello world", "o", 0, len("hello world"), len("hello world") - 4}, // 最后一个 'o'
		{"hello world", "z", 0, len("hello world"), -1},                     // 不存在的子串

		// 边界情况
		{"", "hello", 0, 0, -1},                      // 空字符串
		{"hello", "", 0, len("hello"), len("hello")}, // 空子串
		{"hello", "lo", 0, 3, -1},                    // 子串在start-end范围内
		{"hello", "lo", 0, 10, 3},                    // end 大于字符串长度

		// 负向start, end 测试
		{"hello world", "world", -5, -1, -1}, // 使用负索引
		{"hello world", "o", -5, -1, 7},      // 从倒数第5个字符开始搜索
	}

	for _, tt := range tests {
		t.Run(tt.s+tt.sub, func(t *testing.T) {
			got := RFind(tt.s, tt.sub, tt.start, tt.end)
			if got != tt.expected {
				t.Errorf("RFind(%q, %q, %d, %d) = %d; want %d", tt.s, tt.sub, tt.start, tt.end, got, tt.expected)
			}
		})
	}
}

func TestRIndex(t *testing.T) {
	tests := []struct {
		s, sub     string
		start, end int
		expected   int
		expectErr  bool
	}{
		{"hello world", "world", 0, len("hello world"), len("hello world") - len("world"), false},
		{"hello world", "hello", 0, len("hello world"), 0, false},
		{"hello world", "o", 0, len("hello world"), 7, false},
		{"hello world", "z", 0, len("hello world"), -1, true}, // 不存在的子串

		// 边界情况
		{"", "hello", 0, 0, -1, true},            // 空字符串
		{"hello", "", 0, len("hello"), 5, false}, // 空子串
		{"hello", "lo", 0, 3, -1, true},          // 子串在start-end范围内
		{"hello", "lo", 0, 10, 3, false},         // end 大于字符串长度

		// 负向start, end 测试
		{"hello world", "world", -5, -1, -1, true},
		{"hello world", "o", -5, -1, 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.s+tt.sub, func(t *testing.T) {
			got, err := RIndex(tt.s, tt.sub, tt.start, tt.end)
			if (err != nil) != tt.expectErr {
				t.Errorf("RIndex(%q, %q, %d, %d) unexpected error: %v", tt.s, tt.sub, tt.start, tt.end, err)
			}
			if got != tt.expected {
				t.Errorf("RIndex(%q, %q, %d, %d) = %d; want %d", tt.s, tt.sub, tt.start, tt.end, got, tt.expected)
			}
		})
	}
}

func TestRJust(t *testing.T) {
	tests := []struct {
		s, expected string
		width       int
		fillChar    []rune
	}{
		{"hello", "     hello", 10, nil},         // Default padding with space
		{"hello", "*****hello", 10, []rune{'*'}}, // Custom padding with asterisk
		{"hello", "hello", 5, nil},               // No padding needed
		{"hello", " hello", 6, nil},              // Single space padding
		{"hello", "  hello", 7, nil},             // Padding with space
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := RJust(tt.s, tt.width, tt.fillChar...)
			if got != tt.expected {
				t.Errorf("RJust(%q, %d, %v) = %q; want %q", tt.s, tt.width, tt.fillChar, got, tt.expected)
			}
		})
	}
}

func BenchmarkRJust(b *testing.B) {
	// 测试用的大字符串
	longStr := "a" + string(make([]byte, 1024*1024)) // 1MB的字符串，开头是 'a'

	// 基准测试
	b.Run("Right justify with space", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RJust(longStr, len(longStr)+10, ' ')
		}
	})

	b.Run("Right justify with asterisk", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RJust(longStr, len(longStr)+10, '*')
		}
	})
}

func TestRPartition(t *testing.T) {
	tests := []struct {
		s, sep    string
		expected1 string
		expected2 string
		expected3 string
	}{
		{"a/b/c", "/", "a/b", "/", "c"},
		{"a/b/c", "b", "a/", "b", "/c"},
		{"hello world", " ", "hello", " ", "world"},
		{"hello world", "x", "", "", "hello world"},
		{"", "x", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got1, got2, got3 := RPartition(tt.s, tt.sep)
			if got1 != tt.expected1 || got2 != tt.expected2 || got3 != tt.expected3 {
				t.Errorf("RPartition(%q, %q) = (%q, %q, %q); want (%q, %q, %q)", tt.s, tt.sep, got1, got2, got3, tt.expected1, tt.expected2, tt.expected3)
			}
		})
	}
}

func BenchmarkRPartition(b *testing.B) {
	// 测试用大字符串
	longStr := "a" + string(make([]byte, 1024*1024)) // 1MB的字符串，开头是 'a'

	// 基准测试
	b.Run("Partition with separator", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RPartition(longStr, "a")
		}
	})

	b.Run("Partition without separator", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RPartition(longStr, "z")
		}
	})
}

func TestRSplit(t *testing.T) {
	tests := []struct {
		s        string
		sep      string
		maxsplit int
		expected []string
	}{
		{"a/b/c/d/e", "/", 2, []string{"a/b/c", "d", "e"}},
		{"a b c d e", " ", 2, []string{"a b c", "d", "e"}},
		{"a b c d e", " ", -1, []string{"a b c d e"}},
		{"a a a d e", "", 2, []string{"a a a", "d", "e"}},
		{"hello world", "", 1, []string{"hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := RSplit(tt.s, tt.sep, tt.maxsplit)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RSplit(%v, %v, %v) = %v; want %v", tt.s, tt.sep, tt.maxsplit, got, tt.expected)
			}
		})
	}
}