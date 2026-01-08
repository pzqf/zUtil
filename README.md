# zUtils

zUtils是一个Go语言的常用工具集，提供了多种实用的功能模块，帮助开发者提高开发效率。

## 项目结构

```
zUtil/
├── contenttype/     # MIME类型管理
├── zAes/            # AES加密算法
├── zCache/          # 缓存工具
├── zColor/          # 命令行颜色输出
├── zConfig/         # 配置管理
├── zConcurrency/    # 并发控制工具
├── zCrypto/         # 加密工具集合
├── zDataConv/       # 数据类型转换
├── zFile/           # 文件操作
├── zGps/            # GPS坐标转换
├── zHashtable/      # 哈希表实现
├── zKeyWordFilter/  # 关键词过滤
├── zList/           # 线程安全链表
├── zMap/            # 扩展sync.Map
├── zQueue/          # 队列实现
├── zRand/           # 随机数生成
├── zReflect/        # 反射工具
├── zStack/          # 栈实现
├── zStr/            # 字符串处理
├── zTime/           # 时间处理
├── zTree/           # 树形结构生成
└── zUtils/          # 通用工具函数
```

## 功能模块

### 1. contenttype

提供MIME类型管理功能，根据文件扩展名获取对应的MIME类型。

```go
import "github.com/pzqf/zUtil/contenttype"

// 获取文件的MIME类型
mime := contenttype.GetFileContentType("test.jpg")
// 输出: image/jpeg
```

### 2. zAes

AES加密算法实现，支持CBC、ECB、CFB三种模式。

```go
import "github.com/pzqf/zUtil/zAes"

// AES-CBC加密
key := []byte("1234567890123456")
origData := []byte("Hello, World!")
encrypted := zAes.EncryptCBC(origData, key)

// AES-CBC解密
decrypted := zAes.DecryptCBC(encrypted, key)

// AES-CFB加密（现在返回错误而不是panic）
encryptedCFB, err := zAes.EncryptCFB(origData, key)
if err != nil {
    // 处理错误
}

// AES-CFB解密
decryptedCFB, err := zAes.DecryptCFB(encryptedCFB, key)
if err != nil {
    // 处理错误
}
```

### 3. zColor

命令行颜色输出工具，支持多种颜色和样式。

```go
import "github.com/pzqf/zUtil/zColor"

// 输出不同颜色的文本
fmt.Println(zColor.Red("错误信息"))
fmt.Println(zColor.Green("成功信息"))
fmt.Println(zColor.Blue("提示信息"))
fmt.Println(zColor.Yellow("警告信息"))
```

### 4. zDataConv

数据类型转换工具，提供多种类型之间的转换函数。

```go
import "github.com/pzqf/zUtil/zDataConv"

// 字符串转整数
num, err := zDataConv.String2Int("123")

// 整数转字符串
str := zDataConv.Int2String(123)

// 字符串转浮点数
f, err := zDataConv.String2Float64("3.14")
```

### 5. zGps

GPS坐标转换工具，支持WGS-84、GCJ-02、BD-09坐标系之间的转换，以及距离计算。

```go
import "github.com/pzqf/zUtil/zGps"

// WGS-84转GCJ-02
lat, lon := zGps.Gps84ToGcj02(39.908823, 116.397470)

// GCJ-02转BD-09
lat, lon = zGps.Gcj02ToBd09(lat, lon)

// 计算两点之间的距离
distance := zGps.GetDistance(39.908823, 116.397470, 39.910000, 116.400000)
```

### 6. zHashtable

哈希表实现，支持添加、获取、设置和删除操作。现在支持动态扩容以减少哈希冲突。

```go
import "github.com/pzqf/zUtil/zHashtable"

// 创建哈希表
ht := zHashtable.NewHashTable()

// 添加元素（当负载因子超过0.75时会自动扩容）
ht.Add("key1", "value1")
ht.Add("key2", 123)

// 获取元素
value, exists := ht.Get("key1")

// 设置元素
ht.Set("key1", "new value")

// 删除元素
ht.Remove("key2")
```

### 7. zKeyWordFilter

关键词过滤工具，使用DFA算法实现敏感词过滤。

```go
import "github.com/pzqf/zUtil/zKeyWordFilter"

// 创建过滤器
filter := zKeyWordFilter.NewFilter()

// 添加关键词
filter.AddWord("敏感词")
filter.AddWord("不良内容")

// 过滤文本
result := filter.Filter("这是一段包含敏感词的文本")
// 输出: 这是一段包含***的文本
```

### 8. zList

基于Go标准库container/list的线程安全链表。

```go
import "github.com/pzqf/zUtil/zList"

// 创建链表
list := zList.New()

// 添加元素
list.PushBack("item1")
list.PushFront("item0")

// 遍历链表
list.Range(func(e *list.Element, value any) bool {
    fmt.Println(value)
    return true
})
```

### 9. zMap

基于Go标准库sync.Map的扩展，提供额外的功能如计数和清空。

```go
import "github.com/pzqf/zUtil/zMap"

// 创建Map
m := zMap.NewMap()

// 存储元素
m.Store("key1", "value1")

// 获取元素
value, exists := m.Get("key1")

// 获取长度
length := m.Len()

// 清空Map
m.Clear()
```

### 10. zQueue

队列实现，支持线程安全的入队和出队操作。Dequeue操作已优化为O(1)时间复杂度。

```go
import "github.com/pzqf/zUtil/zQueue"

// 创建队列
queue := zQueue.NewQueue()

// 入队
queue.Enqueue("item1")
queue.Enqueue("item2")

// 出队（现在是O(1)时间复杂度）
item, exists := queue.Dequeue()

// 获取队列长度
length := queue.Length()

// 查看队首元素
peek, exists := queue.Peek()

// 检查队列是否为空
isEmpty := queue.IsEmpty()
```

### 11. zRand

随机数生成工具，包括随机整数、字符串等。

```go
import "github.com/pzqf/zUtil/zRand"

// 生成0-9之间的随机整数
randNum := zRand.RandN(10)

// 生成指定区间的随机整数
randInt := zRand.RandInterval(1, 100)

// 生成随机字符串
randStr := zRand.RandString(10)

// 生成自定义字符集的随机字符串
customStr := zRand.RandCustomString(10, "0123456789")
```

### 12. zStack

栈实现，支持固定大小的栈操作。

```go
import "github.com/pzqf/zUtil/zStack"

// 创建栈
stack := zStack.New(10)

// 入栈
stack.Push("item1")
stack.Push("item2")

// 出栈
item, err := stack.Pop()

// 查看栈顶元素
peek, err := stack.Peek()
```

### 13. zTime

时间处理工具，支持时间转换、格式化、区间计算等。

```go
import "github.com/pzqf/zUtil/zTime"

// 获取当前时间
now := zTime.Now()

// 时间转换
seconds := zTime.Time2Seconds(now.Time())
t := zTime.Seconds2Time(seconds)

// 格式化时间
str := zTime.Time2String(now.Time())

// 获取当天开始和结束时间
beginOfDay := now.BeginOfDay()
endOfDay := now.EndOfDay()

// 获取当周开始和结束时间
beginOfWeek := now.BeginOfWeek()
endOfWeek := now.EndOfWeek()
```

### 14. zTree

树形结构生成工具，基于父子关系生成树形结构。

```go
import "github.com/pzqf/zUtil/zTree"

// 实现INode接口
type Node struct {
    ID       int
    FatherID int
    Data     string
}

func (n Node) GetId() int {
    return n.ID
}

func (n Node) GetFatherId() int {
    return n.FatherID
}

func (n Node) GetData() interface{} {
    return n.Data
}

func (n Node) IsRoot() bool {
    return n.FatherID == 0
}

// 生成树形结构
nodes := []zTree.INode{
    Node{ID: 1, FatherID: 0, Data: "根节点"},
    Node{ID: 2, FatherID: 1, Data: "子节点1"},
    Node{ID: 3, FatherID: 1, Data: "子节点2"},
}

trees := zTree.GenerateTree(nodes)
```

### 15. zUtils

通用工具函数，包括目录操作、错误恢复等。

```go
import "github.com/pzqf/zUtil/zUtils"

// 获取当前目录
dir, err := zUtils.GetCurrentDirectory()

// 错误恢复
defer func() {
    if err := zUtils.Recover(); err != nil {
        fmt.Println("发生错误:", err)
    }
}()
```

### 16. zFile

文件操作工具，提供文件读写、目录管理、文件复制等功能。

```go
import "github.com/pzqf/zUtil/zFile"

// 读取文件内容
content, err := zFile.ReadFile("/path/to/file.txt")

// 写入文件内容
err = zFile.WriteFile("/path/to/file.txt", []byte("内容"), 0644)

// 追加文件内容
err = zFile.AppendFile("/path/to/file.txt", []byte("追加内容"))

// 复制文件
err = zFile.CopyFile("/path/to/src.txt", "/path/to/dst.txt")

// 移动文件
err = zFile.MoveFile("/path/to/src.txt", "/path/to/dst.txt")

// 递归列出目录中的文件
files, err := zFile.ListFilesRecursive("/path/to/dir")
```

### 17. zConfig

配置管理工具，支持JSON、YAML和INI格式的配置文件加载和解析。

```go
import "github.com/pzqf/zUtil/zConfig"

// 创建配置实例
config := zConfig.NewConfig()

// 加载JSON配置文件
err := config.LoadJSON("/path/to/config.json")

// 加载YAML配置文件
err := config.LoadYAML("/path/to/config.yaml")

// 加载INI配置文件
err := config.LoadINI("/path/to/config.ini")

// 获取字符串配置
value, err := config.GetString("server.port")

// 获取整数配置
port, err := config.GetInt("server.port")

// 获取布尔配置
enable, err := config.GetBool("server.enable")

// 设置配置
config.Set("server.timeout", 30)

// 保存JSON配置文件
err := config.SaveJSON("/path/to/config.json")

// 保存YAML配置文件
err := config.SaveYAML("/path/to/config.yaml")

// 保存INI配置文件
err := config.SaveINI("/path/to/config.ini")

// 将配置解析为结构体
type ServerConfig struct {
    Port   int  `json:"port"`
    Enable bool `json:"enable"`
}
var serverConfig ServerConfig
err = config.Unmarshal(&serverConfig)
```



### 19. zStr

字符串处理工具，支持中文字符处理、命名转换、编辑距离计算等功能。

```go
import "github.com/pzqf/zUtil/zStr"

// 中文支持的子字符串提取
substr := zStr.Substring("这是一段中文文本", 2, 4)

// 命名转换
camelCase := zStr.CamelCase("snake_case_string")
snakeCase := zStr.SnakeCase("CamelCaseString")
kebabCase := zStr.KebabCase("CamelCaseString")

// 编辑距离计算
distance := zStr.Distance("kitten", "sitting")

// 随机字符串生成
randStr := zStr.RandomString(10)
```

### 20. zReflect

反射工具，提供结构体字段检查、方法调用、深拷贝比较等功能。

```go
import "github.com/pzqf/zUtil/zReflect"

// 获取结构体字段值
type User struct {
    Name string
    Age  int
}
user := User{Name: "张三", Age: 25}
name, err := zReflect.GetFieldValue(user, "Name")

// 设置结构体字段值
zReflect.SetFieldValue(&user, "Age", 26)

// 调用结构体方法
result, err := zReflect.CallMethod(user, "MethodName", arg1, arg2)

// 深拷贝
newUser, err := zReflect.DeepCopy(user)
```

### 21. zCrypto

加密工具集合，支持哈希计算、编码解码、对称加密、非对称加密等功能。

```go
import "github.com/pzqf/zUtil/zCrypto"

// 哈希计算
md5Hash := zCrypto.MD5("data")
sha256Hash := zCrypto.SHA256("data")

// 编码解码
base64Encoded := zCrypto.Base64Encode([]byte("data"))
base64Decoded, err := zCrypto.Base64Decode(base64Encoded)

// AES加密解密
key := []byte("1234567890123456")
iv := []byte("1234567890123456")
encrypted, err := zCrypto.AESEncrypt([]byte("data"), key, iv, zCrypto.AESModeCBC)
decrypted, err := zCrypto.AESDecrypt(encrypted, key, iv, zCrypto.AESModeCBC)

// RSA加密解密
pubKey, privKey, err := zCrypto.RSAGenerateKey(2048)
rsaEncrypted, err := zCrypto.RSAEncrypt([]byte("data"), pubKey)
rsaDecrypted, err := zCrypto.RSADecrypt(rsaEncrypted, privKey)
```

### 22. zConcurrency

并发控制工具，提供工作池、信号量、带超时的互斥锁、限流器等功能。

```go
import "github.com/pzqf/zUtil/zConcurrency"

// 工作池
pool := zConcurrency.NewWorkerPool(10, 100)
pool.Start()

for i := 0; i < 100; i++ {
    pool.Submit(func() error {
        // 执行任务
        return nil
    })
}

// 信号量
sem := zConcurrency.NewSemaphore(5)
sem.Acquire(1)
// 执行受保护的代码
sem.Release(1)

// 限流器
limiter := zConcurrency.NewRateLimiter(10, 20) // 每秒10个请求，最大突发20个
if limiter.Allow() {
    // 处理请求
}
```

### 23. zCache

缓存工具，提供LRU缓存、TTL缓存等多种缓存实现。

```go
import "github.com/pzqf/zUtil/zCache"

// LRU缓存
lruCache := zCache.NewLRUCache(100, time.Hour)

// 设置缓存
lruCache.Set("key", "value", time.Hour)

// 获取缓存
value, err := lruCache.Get("key")

// 默认缓存实例
zCache.SetDefault("key", "value", time.Hour)
value, err := zCache.GetDefault("key")
```

## 安装

```bash
go get -u github.com/pzqf/zUtil
```

## 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zColor"
    "github.com/pzqf/zUtil/zTime"
)

func main() {
    // 输出带颜色的文本
    fmt.Println(zColor.Green("Hello, zUtils!"))
    
    // 获取当前时间
    now := zTime.Now()
    fmt.Println("当前时间:", zTime.Time2String(now.Time()))
    fmt.Println("本周开始时间:", zTime.Time2String(now.BeginOfWeek()))
}
```

## 项目特点

1. **模块化设计**：每个功能独立成模块，方便使用和维护
2. **线程安全**：大部分容器类型都实现了线程安全
3. **易于扩展**：简单的接口设计，方便扩展新功能
4. **全面的测试**：每个模块都有对应的测试文件

## 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目。

## 许可证

本项目采用MIT许可证。
