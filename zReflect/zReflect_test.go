package zReflect

import (
	"testing"
)

type TestStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Score float64
}

func (t *TestStruct) GetName() string {
	return t.Name
}

func (t *TestStruct) SetAge(age int) {
	t.Age = age
}

func (t *TestStruct) Add(a, b int) int {
	return a + b
}

func TestGetFieldValue(t *testing.T) {
	testObj := TestStruct{
		Name:  "test",
		Age:   18,
		Score: 90.5,
	}

	// 测试获取字段值
	value, err := GetFieldValue(testObj, "Name")
	if err != nil {
		t.Errorf("GetFieldValue failed: %v", err)
	}
	if value != "test" {
		t.Errorf("expected 'test', got %v", value)
	}

	// 测试获取不存在的字段
	_, err = GetFieldValue(testObj, "UnknownField")
	if err == nil {
		t.Error("expected error for unknown field, got nil")
	}

	// 测试非结构体类型
	_, err = GetFieldValue(123, "Name")
	if err == nil {
		t.Error("expected error for non-struct type, got nil")
	}
}

func TestSetFieldValue(t *testing.T) {
	testObj := &TestStruct{
		Name:  "test",
		Age:   18,
		Score: 90.5,
	}

	// 测试设置字段值
	err := SetFieldValue(testObj, "Age", 20)
	if err != nil {
		t.Errorf("SetFieldValue failed: %v", err)
	}
	if testObj.Age != 20 {
		t.Errorf("expected age 20, got %d", testObj.Age)
	}

	// 测试设置不存在的字段
	err = SetFieldValue(testObj, "UnknownField", "value")
	if err == nil {
		t.Error("expected error for unknown field, got nil")
	}

	// 测试类型不匹配
	err = SetFieldValue(testObj, "Age", "20")
	if err == nil {
		t.Error("expected error for type mismatch, got nil")
	}
}

func TestGetStructFields(t *testing.T) {
	testObj := TestStruct{
		Name:  "test",
		Age:   18,
		Score: 90.5,
	}

	fields, err := GetStructFields(testObj)
	if err != nil {
		t.Errorf("GetStructFields failed: %v", err)
	}

	if len(fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fields))
	}

	if _, ok := fields["Name"]; !ok {
		t.Error("expected field 'Name' not found")
	}
}

func TestGetStructFieldTags(t *testing.T) {
	testObj := TestStruct{
		Name:  "test",
		Age:   18,
		Score: 90.5,
	}

	tags, err := GetStructFieldTags(testObj, "json")
	if err != nil {
		t.Errorf("GetStructFieldTags failed: %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}

	if tags["Name"] != "name" {
		t.Errorf("expected tag 'name' for field 'Name', got %s", tags["Name"])
	}

	if tags["Age"] != "age" {
		t.Errorf("expected tag 'age' for field 'Age', got %s", tags["Age"])
	}
}

func TestCallMethod(t *testing.T) {
	testObj := &TestStruct{
		Name:  "test",
		Age:   18,
		Score: 90.5,
	}

	// 测试无参数方法
	result, err := CallMethod(testObj, "GetName")
	if err != nil {
		t.Errorf("CallMethod failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
	if result[0] != "test" {
		t.Errorf("expected 'test', got %v", result[0])
	}

	// 测试有参数无返回值方法
	_, err = CallMethod(testObj, "SetAge", 25)
	if err != nil {
		t.Errorf("CallMethod failed: %v", err)
	}
	if testObj.Age != 25 {
		t.Errorf("expected age 25, got %d", testObj.Age)
	}

	// 测试有参数有返回值方法
	result, err = CallMethod(testObj, "Add", 10, 20)
	if err != nil {
		t.Errorf("CallMethod failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
	if result[0] != 30 {
		t.Errorf("expected 30, got %v", result[0])
	}

	// 测试不存在的方法
	_, err = CallMethod(testObj, "UnknownMethod")
	if err == nil {
		t.Error("expected error for unknown method, got nil")
	}
}

func TestDeepCopy(t *testing.T) {
	original := &TestStruct{
		Name:  "test",
		Age:   18,
		Score: 90.5,
	}

	// 测试深度复制
	copyObj, err := DeepCopy(original)
	if err != nil {
		t.Errorf("DeepCopy failed: %v", err)
	}

	copyTest, ok := copyObj.(*TestStruct)
	if !ok {
		t.Error("DeepCopy returned wrong type")
	}

	// 验证复制的对象与原对象内容相同
	if copyTest.Name != original.Name || copyTest.Age != original.Age || copyTest.Score != original.Score {
		t.Error("DeepCopy failed to copy values correctly")
	}

	// 修改复制后的对象，验证原对象不受影响
	copyTest.Name = "modified"
	if original.Name == "modified" {
		t.Error("original object was modified after copying")
	}
}

func TestIsZeroValue(t *testing.T) {
	// 测试零值
	if !IsZeroValue(0) {
		t.Error("expected 0 to be zero value")
	}
	if !IsZeroValue("") {
		t.Error("expected empty string to be zero value")
	}
	if !IsZeroValue(TestStruct{}) {
		t.Error("expected empty struct to be zero value")
	}

	// 测试非零值
	if IsZeroValue(1) {
		t.Error("expected 1 to be non-zero value")
	}
	if IsZeroValue("test") {
		t.Error("expected 'test' to be non-zero value")
	}
}

func TestGetTypeName(t *testing.T) {
	// 测试基本类型
	if GetTypeName(1) != "int" {
		t.Errorf("expected 'int', got %s", GetTypeName(1))
	}
	if GetTypeName("test") != "string" {
		t.Errorf("expected 'string', got %s", GetTypeName("test"))
	}

	// 测试结构体类型
	testObj := TestStruct{}
	if GetTypeName(testObj) != "TestStruct" {
		t.Errorf("expected 'TestStruct', got %s", GetTypeName(testObj))
	}

	// 测试指针类型
	if GetTypeName(&testObj) != "*TestStruct" {
		t.Errorf("expected '*TestStruct', got %s", GetTypeName(&testObj))
	}

	// 测试nil
	if GetTypeName(nil) != "nil" {
		t.Errorf("expected 'nil', got %s", GetTypeName(nil))
	}
}

func TestHasMethod(t *testing.T) {
	testObj := &TestStruct{}

	// 测试存在的方法
	if !HasMethod(testObj, "GetName") {
		t.Error("expected method 'GetName' to exist")
	}

	// 测试不存在的方法
	if HasMethod(testObj, "UnknownMethod") {
		t.Error("expected method 'UnknownMethod' to not exist")
	}
}

func TestGetAllMethods(t *testing.T) {
	testObj := &TestStruct{}

	methods := GetAllMethods(testObj)
	expectedMethods := []string{"Add", "GetName", "SetAge"}

	if len(methods) != len(expectedMethods) {
		t.Errorf("expected %d methods, got %d", len(expectedMethods), len(methods))
	}

	// 简单检查是否包含所有预期方法
	methodMap := make(map[string]bool)
	for _, m := range methods {
		methodMap[m] = true
	}

	for _, m := range expectedMethods {
		if !methodMap[m] {
		t.Errorf("expected method '%s' not found", m)
		}
	}
}

func TestToMap(t *testing.T) {
	testObj := TestStruct{
		Name:  "test",
		Age:   18,
		Score: 90.5,
	}

	resultMap, err := ToMap(testObj)
	if err != nil {
		t.Errorf("ToMap failed: %v", err)
	}

	if len(resultMap) != 3 {
		t.Errorf("expected 3 fields, got %d", len(resultMap))
	}

	if resultMap["Name"] != "test" {
		t.Errorf("expected Name='test', got %v", resultMap["Name"])
	}
	if resultMap["Age"] != 18 {
		t.Errorf("expected Age=18, got %v", resultMap["Age"])
	}
	if resultMap["Score"] != 90.5 {
		t.Errorf("expected Score=90.5, got %v", resultMap["Score"])
	}
}

func TestToMapWithTags(t *testing.T) {
	testObj := TestStruct{
		Name:  "test",
		Age:   18,
		Score: 90.5,
	}

	resultMap, err := ToMapWithTags(testObj, "json")
	if err != nil {
		t.Errorf("ToMapWithTags failed: %v", err)
	}

	if len(resultMap) != 3 {
		t.Errorf("expected 3 fields, got %d", len(resultMap))
	}

	if resultMap["name"] != "test" {
		t.Errorf("expected name='test', got %v", resultMap["name"])
	}
	if resultMap["age"] != 18 {
		t.Errorf("expected age=18, got %v", resultMap["age"])
	}
	if resultMap["Score"] != 90.5 {
		t.Errorf("expected Score=90.5, got %v", resultMap["Score"])
	}
}

func TestDeepEqual(t *testing.T) {
	// 测试相等的对象
	obj1 := TestStruct{Name: "test", Age: 18}
	obj2 := TestStruct{Name: "test", Age: 18}

	if !DeepEqual(obj1, obj2) {
		t.Error("expected objects to be equal")
	}

	// 测试不相等的对象
	obj3 := TestStruct{Name: "test2", Age: 20}
	if DeepEqual(obj1, obj3) {
		t.Error("expected objects to be not equal")
	}

	// 测试nil
	if !DeepEqual(nil, nil) {
		t.Error("expected nil to be equal to nil")
	}
	if DeepEqual(nil, obj1) {
		t.Error("expected nil to be not equal to object")
	}
}