package zConfig

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	data map[string]interface{}
	path string
}

// NewConfig 创建一个新的配置实例
func NewConfig() *Config {
	return &Config{
		data: make(map[string]interface{}),
	}
}

// LoadJSON 从JSON文件加载配置
func (c *Config) LoadJSON(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}
	
	c.data = data
	c.path = filePath
	return nil
}

// SaveJSON 保存配置到JSON文件
func (c *Config) SaveJSON(filePath string) error {
	content, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return err
	}
	
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return err
	}
	
	c.path = filePath
	return nil
}

// LoadYAML 从YAML文件加载配置
func (c *Config) LoadYAML(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	var data map[string]interface{}
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return err
	}
	
	c.data = data
	c.path = filePath
	return nil
}

// SaveYAML 保存配置到YAML文件
func (c *Config) SaveYAML(filePath string) error {
	content, err := yaml.Marshal(c.data)
	if err != nil {
		return err
	}
	
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return err
	}
	
	c.path = filePath
	return nil
}

// LoadINI 从INI文件加载配置
func (c *Config) LoadINI(filePath string) error {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return err
	}
	
	data := make(map[string]interface{})
	
	// 处理默认section
	defaultSection := cfg.Section("")
	if defaultSection != nil {
		sectionMap := make(map[string]interface{})
		for _, key := range defaultSection.Keys() {
			sectionMap[key.Name()] = convertINIValue(key.Value())
		}
		if len(sectionMap) > 0 {
			// 将默认section的内容直接放在根级别
			for k, v := range sectionMap {
				data[k] = v
			}
		}
	}
	
	// 处理其他section
	for _, section := range cfg.Sections() {
		if section.Name() == "" {
			continue // 跳过默认section
		}
		
		sectionMap := make(map[string]interface{})
		for _, key := range section.Keys() {
			sectionMap[key.Name()] = convertINIValue(key.Value())
		}
		
		data[section.Name()] = sectionMap
	}
	
	c.data = data
	c.path = filePath
	return nil
}

// SaveINI 保存配置到INI文件
func (c *Config) SaveINI(filePath string) error {
	cfg := ini.Empty()
	
	for key, value := range c.data {
		// 检查是否为嵌套的map（section）
		if sectionMap, ok := value.(map[string]interface{}); ok {
			section := cfg.Section(key)
			for k, v := range sectionMap {
				section.Key(k).SetValue(convertToINIString(v))
			}
		} else {
			// 根级别配置
			cfg.Section("").Key(key).SetValue(convertToINIString(value))
		}
	}
	
	err := cfg.SaveTo(filePath)
	if err != nil {
		return err
	}
	
	c.path = filePath
	return nil
}

// convertINIValue 将INI字符串值转换为合适的类型
func convertINIValue(value string) interface{} {
	// 尝试转换为布尔值
	if value == "true" || value == "True" || value == "TRUE" {
		return true
	}
	if value == "false" || value == "False" || value == "FALSE" {
		return false
	}
	
	// 尝试转换为整数
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	
	// 尝试转换为浮点数
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	
	// 默认返回字符串
	return value
}

// convertToINIString 将值转换为INI字符串
func convertToINIString(value interface{}) string {
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case string:
		return v
	default:
		return ""
	}
}

// Get 获取配置值
func (c *Config) Get(key string) (interface{}, error) {
	keys := strings.Split(key, ".")
	var value interface{}
	value = c.data
	
	for _, k := range keys {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid config path: " + key)
		}
		
		value, ok = m[k]
		if !ok {
			return nil, errors.New("key not found: " + key)
		}
	}
	
	return value, nil
}

// GetString 获取字符串类型的配置值
func (c *Config) GetString(key string) (string, error) {
	value, err := c.Get(key)
	if err != nil {
		return "", err
	}
	
	str, ok := value.(string)
	if !ok {
		return "", errors.New("value is not a string: " + key)
	}
	
	return str, nil
}

// GetInt 获取整数类型的配置值
func (c *Config) GetInt(key string) (int, error) {
	value, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	
	// 处理不同的数值类型
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case int64:
		return int(v), nil
	default:
		return 0, errors.New("value is not an integer: " + key)
	}
}

// GetFloat 获取浮点数类型的配置值
func (c *Config) GetFloat(key string) (float64, error) {
	value, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	
	// 处理不同的数值类型
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, errors.New("value is not a float: " + key)
	}
}

// GetBool 获取布尔类型的配置值
func (c *Config) GetBool(key string) (bool, error) {
	value, err := c.Get(key)
	if err != nil {
		return false, err
	}
	
	b, ok := value.(bool)
	if !ok {
		return false, errors.New("value is not a boolean: " + key)
	}
	
	return b, nil
}

// Set 设置配置值
func (c *Config) Set(key string, value interface{}) error {
	keys := strings.Split(key, ".")
	var current interface{}
	current = c.data
	
	for i, k := range keys {
		if i == len(keys)-1 {
			// 设置最终值
			m, ok := current.(map[string]interface{})
			if !ok {
				return errors.New("invalid config path: " + key)
			}
			m[k] = value
			return nil
		}
		
		// 创建嵌套的map
		m, ok := current.(map[string]interface{})
		if !ok {
			return errors.New("invalid config path: " + key)
		}
		
		next, ok := m[k]
		if !ok {
			// 如果键不存在，创建新的map
			newMap := make(map[string]interface{})
			m[k] = newMap
			next = newMap
		}
		
		// 检查是否是map类型
		_, ok = next.(map[string]interface{})
		if !ok {
			return errors.New("invalid config path: " + key)
		}
		
		current = next
	}
	
	return nil
}

// Unmarshal 将配置解析到结构体
func (c *Config) Unmarshal(v interface{}) error {
	content, err := json.Marshal(c.data)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(content, v)
}

// Marshal 将结构体解析到配置
func (c *Config) Marshal(v interface{}) error {
	content, err := json.Marshal(v)
	if err != nil {
		return err
	}
	
	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}
	
	c.data = data
	return nil
}

// Has 检查配置中是否存在某个键
func (c *Config) Has(key string) bool {
	_, err := c.Get(key)
	return err == nil
}

// Delete 删除配置中的某个键
func (c *Config) Delete(key string) error {
	keys := strings.Split(key, ".")
	var current interface{}
	current = c.data
	
	for i, k := range keys {
		if i == len(keys)-1 {
			// 删除最终值
			m, ok := current.(map[string]interface{})
			if !ok {
				return errors.New("invalid config path: " + key)
			}
			delete(m, k)
			return nil
		}
		
		// 检查下一级
		m, ok := current.(map[string]interface{})
		if !ok {
			return errors.New("invalid config path: " + key)
		}
		
		next, ok := m[k]
		if !ok {
			return errors.New("key not found: " + key)
		}
		
		// 检查是否是map类型
		_, ok = next.(map[string]interface{})
		if !ok {
			return errors.New("invalid config path: " + key)
		}
		
		current = next
	}
	
	return nil
}

// Clear 清空配置
func (c *Config) Clear() {
	c.data = make(map[string]interface{})
}

// Size 获取配置中的键值对数量
func (c *Config) Size() int {
	return len(c.data)
}

// IsEmpty 检查配置是否为空
func (c *Config) IsEmpty() bool {
	return len(c.data) == 0
}
