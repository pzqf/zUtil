# zUtils

zUtils是一个Go语言的常用工具集，提供了多种实用的功能模块，帮助开发者提高开发效率。

## 项目结构

```
zUtil/
├── contenttype/     # MIME类型管理
├── zCache/          # 缓存工具
├── zColor/          # 命令行颜色输出
├── zConfig/         # 配置管理
├── zConcurrency/    # 并发控制工具
├── zCrypto/         # 加密工具集合（包含AES、RSA、哈希等）
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

### 2. zCrypto

完整的加密工具集合，包含AES加密、RSA加密、哈希算法、编码等功能。

#### 2.1 AES加密（支持多种模式）

```go
import "github.com/pzqf/zUtil/zCrypto"

// AES-CBC加密（兼容zAes接口）
key := []byte("1234567890123456")
origData := []byte("Hello, World!")
encrypted := zCrypto.EncryptCBC(origData, key)

// AES-CBC解密（兼容zAes接口）
decrypted := zCrypto.DecryptCBC(encrypted, key)

// AES-CFB加密（返回错误）
encryptedCFB, err := zCrypto.EncryptCFB(origData, key)
if err != nil {
    // 处理错误
}

// AES-CFB解密（返回错误）
decryptedCFB, err := zCrypto.DecryptCFB(encryptedCFB, key)
if err != nil {
    // 处理错误
}

// AES-ECB加密（不安全，仅用于兼容）
encryptedECB := zCrypto.EncryptECB(origData, key)
decryptedECB := zCrypto.DecryptECB(encryptedECB, key)

// 推荐使用：AES-GCM加密（支持认证加密，更安全）
encryptedGCM, _ := zCrypto.AESEncrypt(origData, key, nil, zCrypto.AESModeGCM)
decryptedGCM, _ := zCrypto.AESDecrypt(encryptedGCM, key, nil, zCrypto.AESModeGCM)
```

#### 2.2 RSA加密

```go
import "github.com/pzqf/zUtil/zCrypto"

// 生成RSA密钥对
privateKey, publicKey, err := zCrypto.GenerateRSAKeyPair(2048)
if err != nil {
    // 处理错误
}

// RSA加密
encrypted, err := zCrypto.RSAEncrypt([]byte("Hello, World!"), publicKey)
if err != nil {
    // 处理错误
}

// RSA解密
decrypted, err := zCrypto.RSADecrypt(encrypted, privateKey)
if err != nil {
    // 处理错误
}

// RSA签名
signature, err := zCrypto.RSASign([]byte("Hello, World!"), privateKey)
if err != nil {
    // 处理错误
}

// RSA验证签名
err = zCrypto.RSAVerify([]byte("Hello, World!"), signature, publicKey)
if err != nil {
    // 签名无效
}
```

#### 2.3 哈希算法

```go
import "github.com/pzqf/zUtil/zCrypto"

// MD5哈希
md5Hash := zCrypto.MD5("Hello, World!")

// SHA1哈希
sha1Hash := zCrypto.SHA1("Hello, World!")

// SHA256哈希
sha256Hash := zCrypto.SHA256("Hello, World!")

// SHA512哈希
sha512Hash := zCrypto.SHA512("Hello, World!")
```

#### 2.4 编码

```go
import "github.com/pzqf/zUtil/zCrypto"

// Base64编码
encoded := zCrypto.Base64Encode([]byte("Hello, World!"))
decoded, _ := zCrypto.Base64Decode(encoded)

// Hex编码
hexEncoded := zCrypto.HexEncode([]byte("Hello, World!"))
hexDecoded, _ := zCrypto.HexDecode(hexEncoded)
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
value, exists := m.Load("key1")

// 获取长度
length := m.Len()

// 清空Map
m.Clear()
```

#### TypedMap - 类型安全的Map

```go
// 创建类型安全的Map
m := zMap.NewTypedMap[string, int]()

// 存储元素（类型安全）
m.Store("key1", 123)

// 获取元素（类型安全，无需类型断言）
value, exists := m.Load("key1")
// value 是 int 类型，可以直接使用

// 获取长度
length := m.Len()

// 清空Map
m.Clear()
```

#### ShardedMap - 分片Map

```go
// 创建分片Map（默认32个分片）
m := zMap.NewShardedMap(32)

// 或使用快捷方法创建32分片Map
m := zMap.NewShardedMap32()

// 存储元素
m.Store("key1", "value1")

// 获取元素
value, exists := m.Load("key1")

// 获取长度
length := m.Len()
```

#### TypedShardedMap - 类型安全的分片Map

```go
// 创建类型安全的分片Map
m := zMap.NewTypedShardedMap[string, int](32)

// 或使用快捷方法创建32分片的类型安全Map
m := zMap.NewTypedShardedMap32[string, int]()

// 存储元素（类型安全）
m.Store("key1", 123)

// 获取元素（类型安全）
value, exists := m.Load("key1")
```

#### 完整接口列表

```go
// 核心接口
Load(key interface{}) (interface{}, bool)        // 获取元素
Store(key, value interface{})                     // 存储元素
Delete(key interface{})                           // 删除元素
Len() int64                                       // 获取元素数量
Range(f func(key, value interface{}) bool)        // 遍历元素
Clear()                                           // 清空Map
LoadOrStore(key, value interface{}) (interface{}, bool)   // 加载或存储
LoadAndDelete(key interface{}) (interface{}, bool)        // 加载并删除
CompareAndDelete(key, oldValue interface{}) bool          // 比较并删除
CompareAndSwap(key, oldValue, newValue interface{}) bool  // 比较并替换
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

---

## 🆕 最新特性

### 1. zMap 类型安全并发集合

zMap 模块提供了多种高性能、类型安全的并发 Map 实现。

#### Map 类型一览

| 类型 | 适用场景 | 特点 |
|------|---------|------|
| `Map` | 通用并发场景 | sync.Map 包装，提供 Len() 方法 |
| `TypedMap[K, V]` | 读多写少 | 泛型类型安全，零类型断言 |
| `ShardedMap` | 高并发写入 | 分片设计，减少锁竞争 |
| `TypedShardedMap[K, V]` | 高并发+类型安全 | 分片+泛型，最佳选择 |

#### TypedMap 使用示例

```go
import "github.com/pzqf/zUtil/zMap"

// 创建类型安全的 Map
m := zMap.NewTypedMap[string, int]()

// 基本操作
m.Store("key1", 100)
value, ok := m.Load("key1")  // value 是 int 类型，无需断言

// 原子操作
actual, loaded := m.LoadOrStore("key2", 200)
deleted := m.CompareAndDelete("key1", 100)
swapped := m.CompareAndSwap("key1", 100, 200)

// 遍历
m.Range(func(key string, value int) bool {
    fmt.Printf("%s: %d\n", key, value)
    return true
})
```

#### TypedShardedMap 使用示例

```go
// 创建分片 Map（默认32分片）
m := zMap.NewTypedShardedMap[uint64, *User]()

// 创建自定义分片数
m := zMap.NewTypedShardedMapWithShardCount[uint64, *User](64)

// 高并发场景下性能更佳
for i := 0; i < 1000000; i++ {
    m.Store(uint64(i), &User{ID: uint64(i)})
}
```

#### 性能对比

```
BenchmarkTypedMap_Load-8         50000000    25.5 ns/op
BenchmarkTypedMap_Store-8        20000000    65.2 ns/op
BenchmarkTypedShardedMap_Load-8  30000000    42.1 ns/op
BenchmarkTypedShardedMap_Store-8 30000000    48.3 ns/op
```

### 2. zCrypto 完整加密方案

zCrypto 模块提供了完整的加密解决方案，包含对称加密、非对称加密、哈希、编码等功能。

#### AES 加密模式

```go
import "github.com/pzqf/zUtil/zCrypto"

key := []byte("1234567890123456")
data := []byte("Hello, World!")

// 推荐: AES-GCM (认证加密，最安全)
encrypted, _ := zCrypto.AESEncrypt(data, key, nil, zCrypto.AESModeGCM)
decrypted, _ := zCrypto.AESDecrypt(encrypted, key, nil, zCrypto.AESModeGCM)

// AES-CBC (需要IV)
iv := []byte("1234567890123456")
encrypted := zCrypto.EncryptCBC(data, key)
decrypted := zCrypto.DecryptCBC(encrypted, key)
```

#### 哈希计算

```go
// MD5
md5Hash := zCrypto.MD5("data")

// SHA256
sha256Hash := zCrypto.SHA256("data")

// SHA1
sha1Hash := zCrypto.SHA1("data")

// 自定义哈希
hash := zCrypto.Hash("data", crypto.SHA512)
```

#### RSA 加密

```go
// 生成密钥对
privateKey, publicKey, _ := zCrypto.GenerateRSAKeyPair(2048)

// 加密/解密
encrypted, _ := zCrypto.RSAEncrypt([]byte("data"), publicKey)
decrypted, _ := zCrypto.RSADecrypt(encrypted, privateKey)

// 签名/验签
signature, _ := zCrypto.RSASign([]byte("data"), privateKey)
valid := zCrypto.RSAVerify([]byte("data"), signature, publicKey)
```

#### 编码解码

```go
// Base64
encoded := zCrypto.Base64Encode([]byte("data"))
decoded, _ := zCrypto.Base64Decode(encoded)

// Base64 URL Safe
encoded := zCrypto.Base64UrlEncode([]byte("data"))
decoded, _ := zCrypto.Base64UrlDecode(encoded)

// Hex
hexStr := zCrypto.Bytes2Hex([]byte("data"))
bytes := zCrypto.Hex2Bytes(hexStr)
```

### 3. ECDH 密钥交换

支持 Elliptic Curve Diffie-Hellman 密钥交换算法。

```go
// 服务器端
serverDH, _ := zCrypto.NewDHKeyExchange()
serverPubKey := serverDH.GetPublicKey()  // 64字节公钥

// 客户端
clientDH, _ := zCrypto.NewDHKeyExchange()
clientPubKey := clientDH.GetPublicKey()

// 交换公钥后计算共享密钥
serverSharedKey, _ := serverDH.ComputeSharedSecret(clientPubKey)
clientSharedKey, _ := clientDH.ComputeSharedSecret(serverPubKey)

// serverSharedKey == clientSharedKey
```

### 4. zConcurrency 并发工具

提供多种并发控制工具。

#### 信号量

```go
import "github.com/pzqf/zUtil/zConcurrency"

// 创建信号量，限制并发数为10
sem := zConcurrency.NewSemaphore(10)

// 获取许可
sem.Acquire()
defer sem.Release()

// 执行受限操作
doWork()
```

#### 限流器

```go
// 创建限流器：每秒100次请求
limiter := zConcurrency.NewRateLimiter(100, time.Second)

// 检查是否允许
if limiter.Allow() {
    processRequest()
}
```

#### 并发执行器

```go
// 创建执行器，限制最大并发数
executor := zConcurrency.NewConcurrentExecutor(10)

// 提交任务
for i := 0; i < 100; i++ {
    executor.Submit(func() {
        doWork()
    })
}

// 等待所有任务完成
executor.Wait()
```

### 5. zCache 缓存工具

提供多种缓存实现。

#### LRU 缓存

```go
import "github.com/pzqf/zUtil/zCache"

// 创建 LRU 缓存
cache := zCache.NewLRUCache(1000)

// 设置缓存（带过期时间）
cache.Set("key", "value", time.Hour)

// 获取缓存
value, found := cache.Get("key")

// 删除缓存
cache.Delete("key")

// 清空缓存
cache.Clear()
```

#### 默认缓存实例

```go
// 使用全局默认缓存
zCache.SetDefault("key", "value", time.Hour)
value, err := zCache.GetDefault("key")
```

### 6. zTime 时间处理

增强的时间处理工具。

```go
import "github.com/pzqf/zUtil/zTime"

// 创建时间对象
now := zTime.Now()

// 时间计算
tomorrow := now.AddDays(1)
lastWeek := now.SubWeeks(1)

// 时间范围
dayStart := now.BeginOfDay()
dayEnd := now.EndOfDay()
weekStart := now.BeginOfWeek()
monthStart := now.BeginOfMonth()

// 格式化
str := zTime.Time2String(now.Time())
t, _ := zTime.String2Time(str)

// 时间差
diff := zTime.Diff(now.Time(), tomorrow.Time())
fmt.Println(diff.Days, diff.Hours, diff.Minutes)
```

---

## 📊 性能基准

### zMap 性能

| 操作 | TypedMap | TypedShardedMap |
|------|----------|-----------------|
| Load | 25 ns/op | 42 ns/op |
| Store | 65 ns/op | 48 ns/op |
| Delete | 30 ns/op | 45 ns/op |
| Range | 150 ns/iter | 180 ns/iter |

### zCrypto 性能

| 操作 | 吞吐量 |
|------|--------|
| AES-GCM 加密 | 1M ops/s |
| AES-GCM 解密 | 1M ops/s |
| ECDH 密钥交换 | 10K ops/s |
| SHA256 | 500K ops/s |

---

## 🔧 最佳实践

### 1. Map 选择指南

```
场景                          推荐类型
─────────────────────────────────────────────
读多写少                      TypedMap
高并发写入                    TypedShardedMap
需要频繁遍历                  TypedShardedMap
键类型为 uint64/int/string    TypedMap/TypedShardedMap
```

### 2. 加密使用建议

- ✅ 使用 AES-GCM 进行加密（最安全）
- ✅ 使用 ECDH 进行密钥交换
- ✅ 生产环境禁用 AES-ECB
- ⚠️ AES-CBC 需要额外处理 IV 和填充

### 3. 并发控制

- 使用信号量限制资源访问
- 使用限流器保护 API
- 使用执行器管理 goroutine 生命周期

---

## 📝 许可证

MIT License

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

*最后更新: 2026-02*
