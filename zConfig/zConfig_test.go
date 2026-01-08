package zConfig

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	// 创建测试配置文件路径
	filePath := "test_config.json"
	
	// 创建配置实例
	config := NewConfig()
	
	// 设置配置值
	err := config.Set("server.host", "localhost")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("server.port", 8080)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("server.debug", true)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("database.url", "mysql://user:pass@localhost:3306/db")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	// 保存配置到文件
	err = config.SaveJSON(filePath)
	if err != nil {
		t.Fatalf("SaveJSON failed: %v", err)
	}
	
	// 创建新的配置实例
	newConfig := NewConfig()
	
	// 从文件加载配置
	err = newConfig.LoadJSON(filePath)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}
	
	// 测试获取配置值
	host, err := newConfig.GetString("server.host")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	
	if host != "localhost" {
		t.Errorf("Host mismatch: expected 'localhost', got '%s'", host)
	}
	
	port, err := newConfig.GetInt("server.port")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	
	if port != 8080 {
		t.Errorf("Port mismatch: expected 8080, got %d", port)
	}
	
	debug, err := newConfig.GetBool("server.debug")
	if err != nil {
		t.Fatalf("GetBool failed: %v", err)
	}
	
	if !debug {
		t.Errorf("Debug mismatch: expected true, got false")
	}
	
	dbURL, err := newConfig.GetString("database.url")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	
	if dbURL != "mysql://user:pass@localhost:3306/db" {
		t.Errorf("DB URL mismatch: expected 'mysql://user:pass@localhost:3306/db', got '%s'", dbURL)
	}
	
	// 测试Has方法
	if !newConfig.Has("server.host") {
		t.Errorf("Has failed: expected true for 'server.host'")
	}
	
	if newConfig.Has("nonexistent.key") {
		t.Errorf("Has failed: expected false for 'nonexistent.key'")
	}
	
	// 测试Delete方法
	err = newConfig.Delete("server.debug")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	if newConfig.Has("server.debug") {
		t.Errorf("Delete failed: 'server.debug' still exists")
	}
	
	// 测试Unmarshal方法
	type ServerConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	
	type DatabaseConfig struct {
		URL string `json:"url"`
	}
	
	type AppConfig struct {
		Server   ServerConfig   `json:"server"`
		Database DatabaseConfig `json:"database"`
	}
	
	var appConfig AppConfig
	err = newConfig.Unmarshal(&appConfig)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	
	if appConfig.Server.Host != "localhost" {
		t.Errorf("Unmarshal: Host mismatch: expected 'localhost', got '%s'", appConfig.Server.Host)
	}
	
	if appConfig.Server.Port != 8080 {
		t.Errorf("Unmarshal: Port mismatch: expected 8080, got %d", appConfig.Server.Port)
	}
	
	if appConfig.Database.URL != "mysql://user:pass@localhost:3306/db" {
		t.Errorf("Unmarshal: DB URL mismatch: expected 'mysql://user:pass@localhost:3306/db', got '%s'", appConfig.Database.URL)
	}
	
	// 测试Marshal方法
	newAppConfig := AppConfig{
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: 9090,
		},
		Database: DatabaseConfig{
			URL: "postgres://user:pass@localhost:5432/db",
		},
	}
	
	err = config.Marshal(newAppConfig)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	
	// 验证Marshal后的配置
	newHost, err := config.GetString("server.host")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	
	if newHost != "127.0.0.1" {
		t.Errorf("Marshal: Host mismatch: expected '127.0.0.1', got '%s'", newHost)
	}
	
	// 测试Clear方法
	config.Clear()
	if !config.IsEmpty() {
		t.Errorf("Clear failed: config is not empty")
	}
	
	// 清理测试文件
	os.Remove(filePath)
}

func TestConfigYAML(t *testing.T) {
	// 创建测试配置文件路径
	filePath := "test_config.yaml"
	
	// 创建配置实例
	config := NewConfig()
	
	// 设置配置值
	err := config.Set("server.host", "localhost")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("server.port", 8080)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("server.debug", true)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("database.url", "mysql://user:pass@localhost:3306/db")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	// 保存配置到YAML文件
	err = config.SaveYAML(filePath)
	if err != nil {
		t.Fatalf("SaveYAML failed: %v", err)
	}
	
	// 创建新的配置实例
	newConfig := NewConfig()
	
	// 从YAML文件加载配置
	err = newConfig.LoadYAML(filePath)
	if err != nil {
		t.Fatalf("LoadYAML failed: %v", err)
	}
	
	// 测试获取配置值
	host, err := newConfig.GetString("server.host")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	
	if host != "localhost" {
		t.Errorf("Host mismatch: expected 'localhost', got '%s'", host)
	}
	
	port, err := newConfig.GetInt("server.port")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	
	if port != 8080 {
		t.Errorf("Port mismatch: expected 8080, got %d", port)
	}
	
	// 清理测试文件
	os.Remove(filePath)
}

func TestConfigINI(t *testing.T) {
	// 创建测试配置文件路径
	filePath := "test_config.ini"
	
	// 创建配置实例
	config := NewConfig()
	
	// 设置配置值（INI格式的简单结构）
	err := config.Set("host", "localhost")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("port", 8080)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("debug", true)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	// 设置section
	err = config.Set("database.url", "mysql://user:pass@localhost:3306/db")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	err = config.Set("database.max_connections", 100)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	
	// 保存配置到INI文件
	err = config.SaveINI(filePath)
	if err != nil {
		t.Fatalf("SaveINI failed: %v", err)
	}
	
	// 创建新的配置实例
	newConfig := NewConfig()
	
	// 从INI文件加载配置
	err = newConfig.LoadINI(filePath)
	if err != nil {
		t.Fatalf("LoadINI failed: %v", err)
	}
	
	// 测试获取根级别配置值
	host, err := newConfig.GetString("host")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	
	if host != "localhost" {
		t.Errorf("Host mismatch: expected 'localhost', got '%s'", host)
	}
	
	port, err := newConfig.GetInt("port")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	
	if port != 8080 {
		t.Errorf("Port mismatch: expected 8080, got %d", port)
	}
	
	// 测试获取section配置值
	dbURL, err := newConfig.GetString("database.url")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	
	if dbURL != "mysql://user:pass@localhost:3306/db" {
		t.Errorf("DB URL mismatch: expected 'mysql://user:pass@localhost:3306/db', got '%s'", dbURL)
	}
	
	maxConns, err := newConfig.GetInt("database.max_connections")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	
	if maxConns != 100 {
		t.Errorf("Max connections mismatch: expected 100, got %d", maxConns)
	}
	
	// 清理测试文件
	os.Remove(filePath)
}
