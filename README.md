# zUtil

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

zUtil是一个Go语言的常用工具集，提供了多种实用的功能模块，帮助开发者提高开发效率。它设计为模块化、线程安全且性能优化的工具库，适用于各种Go项目。

## 项目概述

zUtil的设计理念是"让Go开发更简单"，通过提供一系列实用的工具函数和数据结构，帮助开发者快速构建高质量的应用。无论你是开发游戏服务器、Web应用还是其他类型的Go项目，zUtil都能为你提供所需的工具支持。

### 为什么选择zUtil？

- **模块化设计**：每个功能独立成模块，方便使用和维护
- **线程安全**：大部分容器类型都实现了线程安全
- **类型安全**：提供泛型实现，减少类型断言
- **性能优化**：关键操作都经过性能优化
- **全面的测试**：每个模块都有对应的测试文件
- **易于扩展**：简单的接口设计，方便扩展新功能

## 项目结构

```
zUtil/
├── zCache/          # 缓存工具 - LRU、TTL等多种缓存实现
├── zColor/          # 命令行颜色输出 - 支持多种颜色和样式
├── zConfig/         # 配置管理 - 支持JSON、YAML、INI格式
├── zConcurrency/    # 并发控制工具 - 工作池、信号量、限流器等
├── zCrypto/         # 加密工具集合 - AES、RSA、哈希等
├── zDataConv/       # 数据类型转换 - 各种类型之间的转换
├── zFile/           # 文件操作 - 文件读写、目录管理等
├── zGps/            # GPS坐标转换 - 支持多种坐标系转换
├── zHashtable/      # 哈希表实现 - 支持动态扩容
├── zKeyWordFilter/  # 关键词过滤 - DFA算法实现敏感词过滤
├── zList/           # 线程安全链表 - 基于标准库的线程安全实现
├── zMap/            # 扩展sync.Map - 提供类型安全的Map实现
├── zQueue/          # 队列实现 - 线程安全的队列
├── zRand/           # 随机数生成 - 各种随机数生成功能
├── zReflect/        # 反射工具 - 结构体操作、方法调用等
├── zStack/          # 栈实现 - 固定大小的栈操作
├── zStr/            # 字符串处理 - 中文字符支持、命名转换等
├── zTime/           # 时间处理 - 时间转换、格式化、区间计算等
├── zTree/           # 树形结构生成 - 基于父子关系生成树形结构
└── zUtils/          # 通用工具函数 - 目录操作、错误恢复等
```

## 核心功能模块

### 1. zCache - 缓存工具

提供多种缓存实现，包括LRU缓存、TTL缓存等，支持缓存过期和自动清理。

#### 主要特性
- 支持LRU（最近最少使用）缓存策略
- 支持TTL（生存时间）缓存
- 线程安全的实现
- 可自定义缓存大小和过期时间
- 提供全局默认缓存实例

#### 使用示例

```go
import "github.com/pzqf/zUtil/zCache"

// 创建LRU缓存，容量100，默认过期时间1小时
lruCache := zCache.NewLRUCache(100, time.Hour)

// 设置缓存，指定过期时间
lruCache.Set("key", "value", time.Hour)

// 获取缓存
value, err := lruCache.Get("key")

// 删除缓存
lruCache.Delete("key")

// 清空缓存
lruCache.Clear()

// 使用全局默认缓存
zCache.SetDefault("key", "value", time.Hour)
value, err := zCache.GetDefault("key")
```

### 2. zColor - 命令行颜色输出

命令行颜色输出工具，支持多种颜色和样式，使终端输出更加美观。

#### 主要特性
- 支持多种颜色（红、绿、蓝、黄等）
- 支持粗体、下划线等样式
- 跨平台兼容
- 简单易用的API

#### 使用示例

```go
import "github.com/pzqf/zUtil/zColor"

// 输出不同颜色的文本
fmt.Println(zColor.Red("错误信息"))
fmt.Println(zColor.Green("成功信息"))
fmt.Println(zColor.Blue("提示信息"))
fmt.Println(zColor.Yellow("警告信息"))

// 输出带样式的文本
fmt.Println(zColor.Bold(zColor.Green("加粗成功信息")))
```

### 3. zConfig - 配置管理

配置管理工具，支持JSON、YAML和INI格式的配置文件加载和解析。

#### 主要特性
- 支持多种配置格式（JSON、YAML、INI）
- 支持配置文件加载和保存
- 支持嵌套配置访问
- 支持配置解析为结构体
- 支持默认值设置

#### 使用示例

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

// 将配置解析为结构体
type ServerConfig struct {
    Port   int  `json:"port"`
    Enable bool `json:"enable"`
}
var serverConfig ServerConfig
err = config.Unmarshal(&serverConfig)
```

### 4. zConcurrency - 并发控制工具

并发控制工具，提供工作池、信号量、带超时的互斥锁、限流器等功能。

#### 主要特性
- 工作池：管理并发任务执行
- 信号量：限制并发访问
- 限流器：控制请求速率
- 并发执行器：管理goroutine生命周期
- 带超时的互斥锁：避免死锁

#### 使用示例

```go
import "github.com/pzqf/zUtil/zConcurrency"

// 创建工作池
pool := zConcurrency.NewWorkerPool(10, 100)
pool.Start()

for i := 0; i < 100; i++ {
    pool.Submit(func() error {
        // 执行任务
        return nil
    })
}

// 停止工作池
pool.Stop()

// 创建信号量，限制并发数为5
sem := zConcurrency.NewSemaphore(5)
sem.Acquire(1)
// 执行受保护的代码
sem.Release(1)

// 创建限流器：每秒10个请求，最大突发20个
limiter := zConcurrency.NewRateLimiter(10, 20)
if limiter.Allow() {
    // 处理请求
}

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

### 5. zCrypto - 加密工具

完整的加密工具集合，包含AES加密、RSA加密、哈希算法、编码等功能。

#### 主要特性
- AES加密：支持CBC、CFB、ECB、GCM等模式
- RSA加密：支持密钥生成、加密、解密、签名
- 哈希算法：支持MD5、SHA1、SHA256、SHA512
- 编码：支持Base64、Hex编码
- ECDH密钥交换：支持椭圆曲线Diffie-Hellman密钥交换

#### 使用示例

```go
import "github.com/pzqf/zUtil/zCrypto"

// AES-GCM加密（推荐，更安全）
key := []byte("1234567890123456")
origData := []byte("Hello, World!")
encrypted, _ := zCrypto.AESEncrypt(origData, key, nil, zCrypto.AESModeGCM)
decrypted, _ := zCrypto.AESDecrypt(encrypted, key, nil, zCrypto.AESModeGCM)

// 生成RSA密钥对
privateKey, publicKey, err := zCrypto.GenerateRSAKeyPair(2048)

// RSA加密
encrypted, err := zCrypto.RSAEncrypt([]byte("Hello, World!"), publicKey)

// RSA解密
decrypted, err := zCrypto.RSADecrypt(encrypted, privateKey)

// SHA256哈希
sha256Hash := zCrypto.SHA256("Hello, World!")

// Base64编码
encoded := zCrypto.Base64Encode([]byte("Hello, World!"))
decoded, _ := zCrypto.Base64Decode(encoded)

// ECDH密钥交换
serverDH, _ := zCrypto.NewDHKeyExchange()
serverPubKey := serverDH.GetPublicKey()

clientDH, _ := zCrypto.NewDHKeyExchange()
clientPubKey := clientDH.GetPublicKey()

serverSharedKey, _ := serverDH.ComputeSharedSecret(clientPubKey)
clientSharedKey, _ := clientDH.ComputeSharedSecret(serverPubKey)
```

### 6. zDataConv - 数据类型转换

数据类型转换工具，提供多种类型之间的转换函数，减少类型转换的代码复杂度。

#### 主要特性
- 字符串与整数、浮点数之间的转换
- 支持错误处理
- 提供简洁的API

#### 使用示例

```go
import "github.com/pzqf/zUtil/zDataConv"

// 字符串转整数
num, err := zDataConv.String2Int("123")

// 整数转字符串
str := zDataConv.Int2String(123)

// 字符串转浮点数
f, err := zDataConv.String2Float64("3.14")
```

### 7. zFile - 文件操作

文件操作工具，提供文件读写、目录管理、文件复制等功能。

#### 主要特性
- 文件读写：支持读取、写入、追加文件内容
- 文件操作：支持复制、移动、删除文件
- 目录操作：支持创建、删除目录，列出目录内容
- 递归操作：支持递归列出目录中的文件

#### 使用示例

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

### 8. zGps - GPS坐标转换

GPS坐标转换工具，支持WGS-84、GCJ-02、BD-09坐标系之间的转换，以及距离计算。

#### 主要特性
- 支持多种坐标系转换
- 支持距离计算
- 高精度实现

#### 使用示例

```go
import "github.com/pzqf/zUtil/zGps"

// WGS-84转GCJ-02
lat, lon := zGps.Gps84ToGcj02(39.908823, 116.397470)

// GCJ-02转BD-09
lat, lon = zGps.Gcj02ToBd09(lat, lon)

// 计算两点之间的距离
distance := zGps.GetDistance(39.908823, 116.397470, 39.910000, 116.400000)
```

### 9. zHashtable - 哈希表实现

哈希表实现，支持添加、获取、设置和删除操作，支持动态扩容以减少哈希冲突。

#### 主要特性
- 支持动态扩容
- 支持基本的哈希表操作
- 简单易用的API

#### 使用示例

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

### 10. zKeyWordFilter - 关键词过滤

关键词过滤工具，使用DFA（确定有限自动机）算法实现敏感词过滤，高效且准确。

#### 主要特性
- 使用DFA算法，高效过滤
- 支持添加和移除关键词
- 支持自定义替换字符

#### 使用示例

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

### 11. zList - 线程安全链表

基于Go标准库container/list的线程安全链表，支持并发操作。

#### 主要特性
- 线程安全
- 支持基本的链表操作
- 兼容标准库list接口

#### 使用示例

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

### 12. zMap - 扩展Map

基于Go标准库sync.Map的扩展，提供额外的功能如计数和清空，以及类型安全的实现。

#### 主要特性
- 类型安全的Map实现
- 分片Map，减少锁竞争
- 支持基本的Map操作
- 支持获取长度和清空操作

#### 使用示例

```go
import "github.com/pzqf/zUtil/zMap"

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

// 创建分片Map（默认32个分片）
shardedMap := zMap.NewShardedMap(32)

// 创建类型安全的分片Map
typedShardedMap := zMap.NewTypedShardedMap[string, int](32)
```

### 13. zQueue - 队列实现

队列实现，支持线程安全的入队和出队操作，Dequeue操作已优化为O(1)时间复杂度。

#### 主要特性
- 线程安全
- O(1)时间复杂度的出队操作
- 支持基本的队列操作

#### 使用示例

```go
import "github.com/pzqf/zUtil/zQueue"

// 创建队列
queue := zQueue.NewQueue()

// 入队
queue.Enqueue("item1")
queue.Enqueue("item2")

// 出队（O(1)时间复杂度）
item, exists := queue.Dequeue()

// 获取队列长度
length := queue.Length()

// 查看队首元素
peek, exists := queue.Peek()

// 检查队列是否为空
isEmpty := queue.IsEmpty()
```

### 14. zRand - 随机数生成

随机数生成工具，包括随机整数、字符串等，提供多种随机数生成功能。

#### 主要特性
- 支持生成随机整数
- 支持生成随机字符串
- 支持自定义字符集
- 线程安全

#### 使用示例

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

### 15. zReflect - 反射工具

反射工具，提供结构体字段检查、方法调用、深拷贝比较等功能，简化反射操作。

#### 主要特性
- 支持获取和设置结构体字段值
- 支持调用结构体方法
- 支持深拷贝
- 支持结构体比较

#### 使用示例

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

### 16. zStack - 栈实现

栈实现，支持固定大小的栈操作，提供基本的栈功能。

#### 主要特性
- 固定大小的栈
- 支持基本的栈操作
- 线程安全

#### 使用示例

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

### 17. zStr - 字符串处理

字符串处理工具，支持中文字符处理、命名转换、编辑距离计算等功能。

#### 主要特性
- 支持中文字符处理
- 支持命名转换（驼峰、蛇形、短横线等）
- 支持编辑距离计算
- 支持随机字符串生成

#### 使用示例

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

### 18. zTime - 时间处理

时间处理工具，支持时间转换、格式化、区间计算等，提供丰富的时间操作功能。

#### 主要特性
- 支持时间转换（时间戳、字符串等）
- 支持时间格式化
- 支持时间区间计算（天、周、月等）
- 支持时间差计算

#### 使用示例

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

// 时间计算
tomorrow := now.AddDays(1)
lastWeek := now.SubWeeks(1)

// 时间差
diff := zTime.Diff(now.Time(), tomorrow.Time())
fmt.Println(diff.Days, diff.Hours, diff.Minutes)
```

### 19. zTree - 树形结构生成

树形结构生成工具，基于父子关系生成树形结构，方便处理层级数据。

#### 主要特性
- 基于父子关系生成树形结构
- 支持自定义节点类型
- 提供树形结构遍历

#### 使用示例

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

### 20. zUtils - 通用工具函数

通用工具函数，包括目录操作、错误恢复等，提供各种实用的辅助函数。

#### 主要特性
- 提供当前目录获取
- 提供错误恢复
- 提供各种辅助函数

#### 使用示例

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

## 安装

```bash
go get -u github.com/pzqf/zUtil
```

## 快速开始

### 基本使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zColor"
    "github.com/pzqf/zUtil/zTime"
)

func main() {
    // 输出带颜色的文本
    fmt.Println(zColor.Green("Hello, zUtil!"))
  
    // 获取当前时间
    now := zTime.Now()
    fmt.Println("当前时间:", zTime.Time2String(now.Time()))
    fmt.Println("本周开始时间:", zTime.Time2String(now.BeginOfWeek()))
}
```

### 高级使用示例

#### 缓存使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/pzqf/zUtil/zCache"
)

func main() {
    // 创建LRU缓存
    cache := zCache.NewLRUCache(100, time.Hour)
    
    // 设置缓存
    cache.Set("user:1", map[string]interface{}{
        "id":   1,
        "name": "张三",
        "age":  25,
    }, time.Hour)
    
    // 获取缓存
    user, err := cache.Get("user:1")
    if err == nil {
        fmt.Println("用户信息:", user)
    }
    
    // 删除缓存
    cache.Delete("user:1")
}
```

#### 并发控制示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/pzqf/zUtil/zConcurrency"
)

func main() {
    // 创建工作池
    pool := zConcurrency.NewWorkerPool(5, 100)
    pool.Start()
    
    // 提交任务
    for i := 0; i < 20; i++ {
        taskID := i
        pool.Submit(func() error {
            fmt.Printf("执行任务 %d\n", taskID)
            time.Sleep(100 * time.Millisecond)
            return nil
        })
    }
    
    // 等待所有任务完成
    time.Sleep(2 * time.Second)
    
    // 停止工作池
    pool.Stop()
    fmt.Println("所有任务执行完成")
}
```

## 性能基准

| 模块 | 操作 | 性能 |
|------|------|------|
| zCrypto | AES-GCM 加密 | 1M ops/s |
| zCrypto | AES-GCM 解密 | 1M ops/s |
| zCrypto | ECDH 密钥交换 | 10K ops/s |
| zCrypto | SHA256 | 500K ops/s |
| zMap | TypedMap Load | 25 ns/op |
| zMap | TypedMap Store | 65 ns/op |
| zMap | TypedShardedMap Load | 42 ns/op |
| zMap | TypedShardedMap Store | 48 ns/op |
| zQueue | Enqueue | 30 ns/op |
| zQueue | Dequeue | 25 ns/op |

## 最佳实践

### zMap 使用建议

| 场景 | 推荐类型 |
|------|----------|
| 读多写少 | TypedMap |
| 高并发写入 | TypedShardedMap |
| 需要频繁遍历 | TypedShardedMap |
| 键类型为 uint64/int/string | TypedMap/TypedShardedMap |

### zCrypto 使用建议

- ✅ 使用 AES-GCM 进行加密（最安全）
- ✅ 使用 ECDH 进行密钥交换
- ✅ 生产环境禁用 AES-ECB
- ⚠️ AES-CBC 需要额外处理 IV 和填充

### zConcurrency 使用建议

- ✅ 使用信号量限制资源访问
- ✅ 使用限流器保护 API
- ✅ 使用执行器管理 goroutine 生命周期
- ⚠️ 合理设置工作池大小，避免过度创建 goroutine

## 项目特点

1. **模块化设计**：每个功能独立成模块，方便使用和维护
2. **线程安全**：大部分容器类型都实现了线程安全
3. **易于扩展**：简单的接口设计，方便扩展新功能
4. **全面的测试**：每个模块都有对应的测试文件
5. **类型安全**：提供泛型实现，减少类型断言
6. **性能优化**：关键操作都经过性能优化
7. **跨平台兼容**：支持Windows、Linux、macOS等平台
8. **丰富的文档**：详细的使用说明和示例

## 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目。

### 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

MIT License

## 联系方式

- 项目主页：https://github.com/pzqf/zUtil
- 问题反馈：https://github.com/pzqf/zUtil/issues

---

*最后更新: 2026-03-31*