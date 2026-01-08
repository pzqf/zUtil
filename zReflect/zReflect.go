package zReflect

import (
	"fmt"
	"reflect"
	"strings"
)

// GetFieldValue 获取结构体字段值
func GetFieldValue(obj interface{}, fieldName string) (interface{}, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj is not a struct, got %s", v.Kind())
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	if !field.CanInterface() {
		return nil, fmt.Errorf("cannot access field %s", fieldName)
	}

	return field.Interface(), nil
}

// SetFieldValue 设置结构体字段值
func SetFieldValue(obj interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("obj must be a pointer, got %s", v.Kind())
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("obj is not a struct, got %s", v.Kind())
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldName)
	}

	valueV := reflect.ValueOf(value)
	if !valueV.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("type mismatch: field type is %s, value type is %s", field.Type(), valueV.Type())
	}

	field.Set(valueV)
	return nil
}

// GetStructFields 获取结构体所有字段信息
func GetStructFields(obj interface{}) (map[string]reflect.StructField, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj is not a struct, got %s", v.Kind())
	}

	fields := make(map[string]reflect.StructField)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fields[field.Name] = field
	}

	return fields, nil
}

// GetStructFieldTags 获取结构体字段标签信息
func GetStructFieldTags(obj interface{}, tagName string) (map[string]string, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj is not a struct, got %s", v.Kind())
	}

	tags := make(map[string]string)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" {
			tags[field.Name] = tag
		}
	}

	return tags, nil
}

// CallMethod 调用结构体方法
func CallMethod(obj interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(obj)

	method := v.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found", methodName)
	}

	// 检查参数个数和类型
	if len(args) != method.Type().NumIn() {
		return nil, fmt.Errorf("method %s expects %d arguments, got %d", methodName, method.Type().NumIn(), len(args))
	}

	// 转换参数类型
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		argType := method.Type().In(i)
		argValue := reflect.ValueOf(arg)
		if !argValue.Type().AssignableTo(argType) {
			return nil, fmt.Errorf("argument %d type mismatch: expected %s, got %s", i, argType, argValue.Type())
		}
		in[i] = argValue
	}

	// 调用方法
	out := method.Call(in)

	// 转换返回值
	result := make([]interface{}, len(out))
	for i, val := range out {
		result[i] = val.Interface()
	}

	return result, nil
}

// DeepCopy 深度复制对象
func DeepCopy(src interface{}) (interface{}, error) {
	if src == nil {
		return nil, nil
	}

	v := reflect.ValueOf(src)
	return deepCopy(v), nil
}

func deepCopy(v reflect.Value) interface{} {
	if !v.IsValid() {
		return nil
	}

	t := v.Type()
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		newPtr := reflect.New(t.Elem())
		newVal := deepCopy(v.Elem())
		newPtr.Elem().Set(reflect.ValueOf(newVal))
		return newPtr.Interface()

	case reflect.Struct:
		newStruct := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			fieldValue := v.Field(i)
			newValue := deepCopy(fieldValue)
			newStruct.Field(i).Set(reflect.ValueOf(newValue))
		}
		return newStruct.Interface()

	case reflect.Slice:
		if v.IsNil() {
			return nil
		}
		newSlice := reflect.MakeSlice(t, v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			newValue := deepCopy(v.Index(i))
			newSlice.Index(i).Set(reflect.ValueOf(newValue))
		}
		return newSlice.Interface()

	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		newMap := reflect.MakeMap(t)
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()
			newKey := deepCopy(key)
			newValue := deepCopy(value)
			newMap.SetMapIndex(reflect.ValueOf(newKey), reflect.ValueOf(newValue))
		}
		return newMap.Interface()

	default:
		return v.Interface()
	}
}

// IsZeroValue 检查值是否为零值
func IsZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

// GetTypeName 获取类型名称
func GetTypeName(v interface{}) string {
	if v == nil {
		return "nil"
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}

// HasMethod 检查对象是否具有指定方法
func HasMethod(obj interface{}, methodName string) bool {
	v := reflect.ValueOf(obj)
	method := v.MethodByName(methodName)
	return method.IsValid()
}

// GetAllMethods 获取结构体所有方法
func GetAllMethods(obj interface{}) []string {
	v := reflect.ValueOf(obj)
	t := v.Type()
	methods := make([]string, 0, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methods = append(methods, method.Name)
	}

	return methods
}

// ToMap 将结构体转换为map
func ToMap(obj interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj is not a struct, got %s", v.Kind())
	}

	result := make(map[string]interface{})
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		if !fieldValue.CanInterface() {
			continue
		}
		result[field.Name] = fieldValue.Interface()
	}

	return result, nil
}

// ToMapWithTags 将结构体转换为map，使用指定标签作为key
func ToMapWithTags(obj interface{}, tagName string) (map[string]interface{}, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj is not a struct, got %s", v.Kind())
	}

	result := make(map[string]interface{})
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		if !fieldValue.CanInterface() {
			continue
		}

		tag := field.Tag.Get(tagName)
		if tag != "" {
			// 处理类似 `json:"name,omitempty"` 的标签
			if idx := strings.Index(tag, ","); idx != -1 {
				tag = tag[:idx]
			}
			result[tag] = fieldValue.Interface()
		} else {
			result[field.Name] = fieldValue.Interface()
		}
	}

	return result, nil
}

// DeepEqual 深度比较两个对象是否相等
func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}