package zStr

import (
	"testing"
)

func TestSubstring(t *testing.T) {
	// 测试英文
	result := Substring("Hello, World!", 0, 5)
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
	
	// 测试中文
	result = Substring("你好，世界！", 0, 2)
	if result != "你好" {
		t.Errorf("Expected '你好', got '%s'", result)
	}
	
	// 测试超出范围
	result = Substring("Hello", 10, 5)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
	
	// 测试从中间开始
	result = Substring("Hello, World!", 7, 5)
	if result != "World" {
		t.Errorf("Expected 'World', got '%s'", result)
	}
}

func TestTruncate(t *testing.T) {
	// 测试英文
	result := Truncate("Hello, World!", 5, "...")
	if result != "Hello..." {
		t.Errorf("Expected 'Hello...', got '%s'", result)
	}
	
	// 测试中文
	result = Truncate("你好，世界！", 2, "...")
	if result != "你好..." {
		t.Errorf("Expected '你好...', got '%s'", result)
	}
	
	// 测试不需要截断
	result = Truncate("Hello", 10, "...")
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
}

func TestContains(t *testing.T) {
	// 测试英文
	result := Contains("Hello, World!", "World")
	if !result {
		t.Errorf("Expected true, got false")
	}
	
	// 测试中文
	result = Contains("你好，世界！", "世界")
	if !result {
		t.Errorf("Expected true, got false")
	}
	
	// 测试不包含
	result = Contains("Hello", "World")
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestContainsAny(t *testing.T) {
	result := ContainsAny("Hello, World!", "World", "Go")
	if !result {
		t.Errorf("Expected true, got false")
	}
	
	result = ContainsAny("Hello", "World", "Go")
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestContainsAll(t *testing.T) {
	result := ContainsAll("Hello, World!", "Hello", "World")
	if !result {
		t.Errorf("Expected true, got false")
	}
	
	result = ContainsAll("Hello, World!", "Hello", "Go")
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestReplace(t *testing.T) {
	// 测试英文
	result := Replace("Hello, World!", "World", "Go", 1)
	if result != "Hello, Go!" {
		t.Errorf("Expected 'Hello, Go!', got '%s'", result)
	}
	
	// 测试中文
	result = Replace("你好，世界！", "世界", "中国", 1)
	if result != "你好，中国！" {
		t.Errorf("Expected '你好，中国！', got '%s'", result)
	}
	
	// 测试替换所有
	result = ReplaceAll("Hello, World! Hello, World!", "World", "Go")
	if result != "Hello, Go! Hello, Go!" {
		t.Errorf("Expected 'Hello, Go! Hello, Go!', got '%s'", result)
	}
}

func TestReplaceRegexp(t *testing.T) {
	result, err := ReplaceRegexp("Hello, 123 World!", `\d+`, "456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if result != "Hello, 456 World!" {
		t.Errorf("Expected 'Hello, 456 World!', got '%s'", result)
	}
}

func TestTrim(t *testing.T) {
	result := Trim("  Hello, World!  ", "")
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", result)
	}
	
	result = Trim("xxHello, World!xx", "x")
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", result)
	}
}

func TestPad(t *testing.T) {
	result := PadLeft("Hello", 10, " ")
	if result != "     Hello" {
		t.Errorf("Expected '     Hello', got '%s'", result)
	}
	
	result = PadRight("Hello", 10, " ")
	if result != "Hello     " {
		t.Errorf("Expected 'Hello     ', got '%s'", result)
	}
	
	result = PadLeft("Hello", 5, " ")
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
}

func TestCase(t *testing.T) {
	result := Upper("Hello")
	if result != "HELLO" {
		t.Errorf("Expected 'HELLO', got '%s'", result)
	}
	
	result = Lower("HELLO")
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
	
	result = Title("hello world")
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result)
	}
}

func TestNamingConventions(t *testing.T) {
	result := CamelCase("hello_world")
	if result != "helloWorld" {
		t.Errorf("Expected 'helloWorld', got '%s'", result)
	}
	
	result = SnakeCase("helloWorld")
	if result != "hello_world" {
		t.Errorf("Expected 'hello_world', got '%s'", result)
	}
	
	result = KebabCase("helloWorld")
	if result != "hello-world" {
		t.Errorf("Expected 'hello-world', got '%s'", result)
	}
}

func TestReverse(t *testing.T) {
	// 测试英文
	result := Reverse("Hello")
	if result != "olleH" {
		t.Errorf("Expected 'olleH', got '%s'", result)
	}
	
	// 测试中文
	result = Reverse("你好世界")
	if result != "界世好你" {
		t.Errorf("Expected '界世好你', got '%s'", result)
	}
}

func TestRandom(t *testing.T) {
	result := RandomLetter(10)
	if len(result) != 10 {
		t.Errorf("Expected length 10, got %d", len(result))
	}
	
	result = RandomDigit(10)
	if len(result) != 10 {
		t.Errorf("Expected length 10, got %d", len(result))
	}
	
	result = RandomAlphanumeric(10)
	if len(result) != 10 {
		t.Errorf("Expected length 10, got %d", len(result))
	}
}

func TestWrap(t *testing.T) {
	result := Wrap("This is a long text that should be wrapped at 20 characters.", 20)
	lines := Lines(result)
	
	for _, line := range lines {
		if len(line) > 20 && line != "" {
			t.Errorf("Line '%s' is longer than 20 characters", line)
		}
	}
}

func TestUnwrap(t *testing.T) {
	result := Unwrap("This is a\nlong text\r\nwith multiple\rnewlines.")
	if result != "This is a long text with multiple newlines." {
		t.Errorf("Expected 'This is a long text with multiple newlines.', got '%s'", result)
	}
}

func TestIndent(t *testing.T) {
	result := Indent("Hello\nWorld", "  ")
	if result != "  Hello\n  World" {
		t.Errorf("Expected '  Hello\n  World', got '%s'", result)
	}
}

func TestUnindent(t *testing.T) {
	result := Unindent("  Hello\n  World\n    Foo")
	if result != "Hello\nWorld\n  Foo" {
		t.Errorf("Expected 'Hello\nWorld\n  Foo', got '%s'", result)
	}
}

func TestQuote(t *testing.T) {
	result := Quote("Hello", "")
	if result != `"Hello"` {
		t.Errorf("Expected '\"Hello\"', got '%s'", result)
	}
	
	result = Quote("Hello", "'")
	if result != "'Hello'" {
		t.Errorf("Expected %q, got %q", "'Hello'", result)
	}
}

func TestUnquote(t *testing.T) {
	result := Unquote(`"Hello"`)
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
	
	result = Unquote("'Hello'")
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
	
	result = Unquote("Hello")
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
}

func TestDistance(t *testing.T) {
	result := Distance("hello", "world")
	if result != 4 {
		t.Errorf("Expected 4, got %d", result)
	}
	
	result = Distance("kitten", "sitting")
	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}
	
	result = Distance("hello", "hello")
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestIsEmpty(t *testing.T) {
	if !IsEmpty("") {
		t.Errorf("Expected true for empty string")
	}
	
	if !IsEmpty("   ") {
		t.Errorf("Expected true for whitespace string")
	}
	
	if IsEmpty("Hello") {
		t.Errorf("Expected false for non-empty string")
	}
}

func TestIsNotEmpty(t *testing.T) {
	if IsNotEmpty("") {
		t.Errorf("Expected false for empty string")
	}
	
	if IsNotEmpty("   ") {
		t.Errorf("Expected false for whitespace string")
	}
	
	if !IsNotEmpty("Hello") {
		t.Errorf("Expected true for non-empty string")
	}
}

func TestIsAlpha(t *testing.T) {
	if !IsAlpha("Hello") {
		t.Errorf("Expected true for 'Hello'")
	}
	
	if IsAlpha("Hello123") {
		t.Errorf("Expected false for 'Hello123'")
	}
}

func TestIsNumeric(t *testing.T) {
	if !IsNumeric("123") {
		t.Errorf("Expected true for '123'")
	}
	
	if IsNumeric("123abc") {
		t.Errorf("Expected false for '123abc'")
	}
}

func TestIsAlphanumeric(t *testing.T) {
	if !IsAlphanumeric("Hello123") {
		t.Errorf("Expected true for 'Hello123'")
	}
	
	if IsAlphanumeric("Hello123!") {
		t.Errorf("Expected false for 'Hello123!'")
	}
}

func TestCapitalize(t *testing.T) {
	result := Capitalize("hello")
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
}

func TestDecapitalize(t *testing.T) {
	result := Decapitalize("Hello")
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}
