# zUtil

zUtil是一个Go语言的常用工具集，提供了多种实用的功能模块，帮助开发者提高开发效率。

## 项目结构

```
zUtil/
├── zCache/          # 缓存工具
├── zColor/          # 命令行颜色输出
├── zConfig/         # 配置管理
├── zConcurrency/    # 并发控制工具
├── zCrypto/         # 加密工具集合（包含AES、RSA、哈希等）
├── zDataConv/       # 数据类型转换
├── zDistributed/    # 分布式锁实现
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

### 1. zCache

缓存工具，提供LRU缓存、TTL缓存等多种缓存实现。

#### 1.1 基本使用

```go
import "github.com/pzqf/zUtil/zCache"

// LRU缓存
lruCache := zCache.NewLRUCache(100, time.Hour)

// 设置缓存
lruCache.Set("key", "value", time.Hour)

// 获取缓存
value, err := lruCache.Get("key")

// 删除缓存
lruCache.Delete("key")

// 清空缓存
lruCache.Clear()
```

#### 1.2 默认缓存实例

```go
// 使用全局默认缓存
zCache.SetDefault("key", "value", time.Hour)
value, err := zCache.GetDefault("key")
```

### 2. zColor

命令行颜色输出工具，支持多种颜色和样式。

```go
import "github.com/pzqf/zUtil/zColor"

// 输出不同颜色的文本
fmt.Println(zColor.Red("错误信息"))
fmt.Println(zColor.Green("成功信息"))
fmt.Println(zColor.Blue("提示信息"))
fmt.Println(zColor.Yellow("警告信息"))
```

### 3. zConfig

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

### 4. zConcurrency

并发控制工具，提供工作池、信号量、带超时的互斥锁、限流器等功能。

#### 4.1 工作池

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
```

#### 4.2 信号量

```go
// 创建信号量，限制并发数为5
sem := zConcurrency.NewSemaphore(5)
sem.Acquire(1)
// 执行受保护的代码
sem.Release(1)
```

#### 4.3 限流器

```go
// 创建限流器：每秒10个请求，最大突发20个
limiter := zConcurrency.NewRateLimiter(10, 20)
if limiter.Allow() {
    // 处理请求
}
```

#### 4.4 并发执行器

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

#### 4.5 最佳实践

- 使用信号量限制资源访问
- 使用限流器保护 API
- 使用执行器管理 goroutine 生命周期

### 5. zCrypto

完整的加密工具集合，包含AES加密、RSA加密、哈希算法、编码等功能。

#### 5.1 AES加密（支持多种模式）

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

#### 5.2 RSA加密

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

#### 5.3 哈希算法

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

#### 5.4 编码

```go
import "github.com/pzqf/zUtil/zCrypto"

// Base64编码
encoded := zCrypto.Base64Encode([]byte("Hello, World!"))
decoded, _ := zCrypto.Base64Decode(encoded)

// Hex编码
hexEncoded := zCrypto.HexEncode([]byte("Hello, World!"))
hexDecoded, _ := zCrypto.HexDecode(hexEncoded)
```

#### 5.5 ECDH 密钥交换

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

#### 5.6 性能基准

| 操作 | 吞吐量 |
|------|--------|
| AES-GCM 加密 | 1M ops/s |
| AES-GCM 解密 | 1M ops/s |
| ECDH 密钥交换 | 10K ops/s |
| SHA256 | 500K ops/s |

#### 5.7 最佳实践

- ✅ 使用 AES-GCM 进行加密（最安全）
- ✅ 使用 ECDH 进行密钥交换
- ✅ 生产环境禁用 AES-ECB
- ⚠️ AES-CBC 需要额外处理 IV 和填充

### 6. zDataConv

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

### 7. zDistributed

分布式锁实现，基于etcd的分布式锁，支持租约机制和自动续约。

#### 7.1 基本使用

```go
import (
    "context"
    "github.com/pzqf/zUtil/zDistributed"
    "go.etcd.io/etcd/client/v3"
)

// 创建etcd客户端
client, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"192.168.91.128:2379"},
    DialTimeout: 5 * time.Second,
})
if err != nil {
    // 处理错误
}
defer client.Close()

// 创建锁管理器
lockManager := zDistributed.NewLockManager(client, zDistributed.DefaultLockOptions)

// 获取锁
lock := lockManager.GetLock("resource-key")
ctx := context.Background()

// 测试获取锁
err = lock.Lock(ctx)
if err != nil {
    // 处理错误
}

// 执行业务操作
// ...

// 释放锁
err = lock.Unlock(ctx)
if err != nil {
    // 处理错误
}

// 或通过管理器释放
err = lockManager.ReleaseLock("resource-key")
if err != nil {
    // 处理错误
}
```

#### 7.2 集群支持

```go
// 连接到etcd集群
client, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"node1:2379", "node2:2379", "node3:2379"},
    DialTimeout: 5 * time.Second,
})
```

#### 7.3 自定义锁选项

```go
options := zDistributed.LockOptions{
    Type:          zDistributed.LockTypeExclusive, // 排他锁
    LeaseTTL:      30,                           // 租约过期时间（秒）
    RetryCount:    5,                            // 重试次数
    RetryInterval: 200,                          // 重试间隔（毫秒）
}

lockManager := zDistributed.NewLockManager(client, options)
```

#### 7.4 最佳实践

- ✅ 使用 zDistributed 进行跨服务的资源协调
- ✅ 为锁设置合理的租约时间
- ✅ 尽量减少锁的持有时间
- ✅ 使用锁管理器统一管理多个锁
- ⚠️ 确保 etcd 集群的高可用性

### 8. zFile

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

### 9. zGps

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

### 10. zHashtable

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

### 11. zKeyWordFilter

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

### 12. zList

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

### 13. zMap

基于Go标准库sync.Map的扩展，提供额外的功能如计数和清空，以及类型安全的实现。

#### 13.1 基本Map

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

#### 13.2 TypedMap - 类型安全的Map

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

#### 13.3 ShardedMap - 分片Map

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

#### 13.4 TypedShardedMap - 类型安全的分片Map

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

#### 13.5 完整接口列表

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

#### 13.6 性能基准

| 操作 | TypedMap | TypedShardedMap |
|------|----------|-----------------|
| Load | 25 ns/op | 42 ns/op |
| Store | 65 ns/op | 48 ns/op |
| Delete | 30 ns/op | 45 ns/op |
| Range | 150 ns/iter | 180 ns/iter |

#### 13.7 最佳实践

```
场景                          推荐类型
─────────────────────────────────────────────
读多写少                      TypedMap
高并发写入                    TypedShardedMap
需要频繁遍历                  TypedShardedMap
键类型为 uint64/int/string    TypedMap/TypedShardedMap
```

### 14. zQueue

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

### 15. zRand

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

### 16. zReflect

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

### 17. zStack

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

### 18. zStr

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

### 19. zTime

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

// 时间计算
tomorrow := now.AddDays(1)
lastWeek := now.SubWeeks(1)

// 时间差
diff := zTime.Diff(now.Time(), tomorrow.Time())
fmt.Println(diff.Days, diff.Hours, diff.Minutes)
```

### 20. zTree

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

### 21. zUtils

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
    fmt.Println(zColor.Green("Hello, zUtil!"))
  
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
5. **类型安全**：提供泛型实现，减少类型断言
6. **性能优化**：关键操作都经过性能优化

## 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目。

---

## 📝 许可证

MIT License

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

*最后更新: 2026-03-06*