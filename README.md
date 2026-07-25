# zUtil Go 通用工具集

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Zero Dependencies](https://img.shields.io/badge/Dependencies-Zero-green.svg)]()

zUtil 是一个 Go 语言通用工具集，提供并发安全的数据结构、加密、缓存、字符串处理等基础功能。所有模块零外部依赖（仅标准库），模块间完全独立，可按需引入。

## 目录

- [模块总览](#模块总览)
- [zMap 并发安全 Map](#zmap-并发安全-map)
- [zConcurrency 并发控制](#zconcurrency-并发控制)
- [zCrypto 加密工具](#zcrypto-加密工具)
- [zCache 缓存](#zcache-缓存)
- [zKeyWordFilter 敏感词过滤](#zkeywordfilter-敏感词过滤)
- [zList 线程安全链表](#zlist-线程安全链表)
- [zQueue 队列](#zqueue-队列)
- [zStack 栈](#zstack-栈)
- [zTree 树形结构](#ztree-树形结构)
- [zStr 字符串处理](#zstr-字符串处理)
- [zTime 时间处理](#ztime-时间处理)
- [zRand 随机数](#zrand-随机数)
- [zReflect 反射工具](#zreflect-反射工具)
- [zFile 文件操作](#zfile-文件操作)
- [zHashtable 哈希表](#zhashtable-哈希表)
- [zDataConv 数据类型转换](#zdataconv-数据类型转换)
- [zError 带错误码的错误类型](#zerror-带错误码的错误类型)
- [zUtils 通用工具](#zutils-通用工具)
- [zColor 终端颜色](#zcolor-终端颜色)
- [zGps GPS 坐标转换](#zgps-gps-坐标转换)

## 模块总览

| 模块 | 包路径 | 功能 | 线程安全 | 测试覆盖 | 代码行数 |
|------|--------|------|----------|----------|----------|
| zMap | zUtil/zMap | 并发安全 Map（分片/泛型） | ✅ | 非常充分 | ~400 |
| zConcurrency | zUtil/zConcurrency | 工作池/信号量/限流器/原子计数器 | ✅ | 较充分 | ~350 |
| zCrypto | zUtil/zCrypto | AES/RSA/哈希/编码 | ✅ | 充分 | ~500 |
| zCache | zUtil/zCache | LRU/Simple/TTL 缓存 | ✅ | 充分 | ~460 |
| zKeyWordFilter | zUtil/zKeyWordFilter | DFA 敏感词过滤 | ✅ | 充分 | ~180 |
| zList | zUtil/zList | 线程安全双向链表 | ✅ | 充分 | ~115 |
| zQueue | zUtil/zQueue | 环形队列 | ✅ | 不足 | ~130 |
| zStack | zUtil/zStack | 固定大小栈/并发安全栈 | ✅ | 不足 | ~115 |
| zTree | zUtil/zTree | 树形结构生成 | ❌ | 基本 | ~80 |
| zStr | zUtil/zStr | 字符串处理（40+函数） | ❌ | 充分 | ~500 |
| zTime | zUtil/zTime | 时间包装器/转换 | ❌ | 不足 | ~150 |
| zRand | zUtil/zRand | 随机数/安全随机 | ✅ | 不足 | ~60 |
| zReflect | zUtil/zReflect | 反射工具/深拷贝 | ❌ | 充分 | ~200 |
| zFile | zUtil/zFile | 文件读写/目录操作 | ❌ | 充分 | ~200 |
| zHashtable | zUtil/zHashtable | 哈希表（自动扩容） | ✅ | 充分 | ~250 |
| zDataConv | zUtil/zDataConv | 数据类型转换 | ❌ | 不足 | ~100 |
| zError | zUtil/zError | 带错误码的错误类型 | ❌ | 充分 | ~50 |
| zUtils | zUtil/zUtils | panic 恢复/工作目录 | ❌ | 充分 | ~45 |
| zColor | zUtil/zColor | ANSI 终端颜色 | ❌ | 不足 | ~80 |
| zGps | zUtil/zGps | GPS 坐标系转换 | ❌ | 不足 | ~150 |

## 设计原则

1. **零外部依赖**：所有模块仅使用 Go 标准库，无第三方依赖
2. **模块独立**：每个模块完全独立，可按需引入，不存在模块间依赖
3. **接口简洁**：API 设计直观易用，与标准库风格保持一致
4. **并发安全优先**：高频使用的数据结构均提供并发安全版本
5. **泛型支持**：zMap 提供 Go 1.18+ 泛型封装，消除类型断言

## 适用范围

**适合**：任何 Go 项目里想要"拿来即用、不引入第三方依赖"的基础件——并发安全容器
（zMap/zList/zQueue/zStack/zHashtable）、加密与编码（zCrypto）、缓存（zCache）、限流与原子计数
（zConcurrency）、字符串/时间/随机/反射/文件小工具等。每个模块相互独立，可单独 `import`，
不会牵入其它模块。

**与 [zEngine](https://github.com/pzqf/zEngine) 的关系**：zUtil 是最底层、零第三方依赖的工具库；
zEngine（游戏服务器引擎）依赖它作为基础件。你可以只用 zUtil 而不碰 zEngine。

**注意**：并发安全性**按模块而定**——见上方「模块总览」的「线程安全」列（标 ❌ 的需调用方自行加锁）。
版本 `0.0.x`，成熟前不发 `1.0`。

## 快速开始

```bash
go get github.com/pzqf/zUtil
```

```go
import (
    "github.com/pzqf/zUtil/zMap"
    "github.com/pzqf/zUtil/zConcurrency"
    "github.com/pzqf/zUtil/zCache"
)
```

***

## zMap 并发安全 Map

zMap 提供四种并发安全的 Map 实现，从基础到高性能逐步递进。

### 实现方式

#### Map（基础版）

基于 `sync.Map` 封装，额外维护 `atomic.Int64` 计数器：

- `Store` 时先 `Load` 检查 key 是否存在，不存在则 `count.Add(1)`
- `Delete` 使用 `LoadAndDelete` 原子操作，删除成功则 `count.Add(-1)`
- `Clear` 先 Range 收集所有 key，再逐个 Delete，最后 `count.Store(0)`

#### ShardedMap（分片版）

将数据分散到多个 `Map` 分片中：

- 使用 FNV-1a 64 位哈希计算 key 的分片索引
- 支持 int64/int32/int/string/[]byte 类型 key，其他类型哈希值为 0
- 减少锁竞争，适合高并发场景

#### TypedMap（泛型版）

泛型包装 `Map`，在 Load 等返回时做类型断言 `value.(V)`，消除外部类型断言。

#### TypedShardedMap（泛型分片版）

泛型包装 `ShardedMap`，结合分片性能和类型安全。

### 使用说明

#### 基础 Map

```go
m := zMap.NewMap()

m.Store("key", "value")
val, ok := m.Load("key")       // val = "value", ok = true
count := m.Len()                // count = 1

m.Delete("key")
m.Clear()                       // 清空所有

// 原子操作
actual, loaded := m.LoadOrStore("key", "default")
val, loaded := m.LoadAndDelete("key")
swapped := m.CompareAndSwap("key", "old", "new")
deleted := m.CompareAndDelete("key", "old")

// 遍历
m.Range(func(key, value interface{}) bool {
    fmt.Println(key, value)
    return true
})
```

#### 分片 Map

```go
sm := zMap.NewShardedMap32()
sm.Store("key", "value")
val, ok := sm.Load("key")
sm.Delete("key")
```

#### 泛型 Map

```go
tm := zMap.NewTypedMap[string, int]()
tm.Store("hp", 100)
hp, ok := tm.Load("hp")

tsm := zMap.NewTypedShardedMap32[int64, *Player]()
tsm.Add(1001, &Player{Name: "test"})
player, ok := tsm.Get(1001)
```

### 注意事项

1. **分片数量**：建议为 2 的幂（16/32/64），分片越多并发性能越好，但内存开销略增
2. **计数精度**：`Len()` 返回的是近似值，并发操作时可能有微小误差
3. **Clear 性能**：`Clear()` 需要先遍历收集 key 再逐个删除，大数据量时较慢
4. **key 类型**：`ShardedMap` 的哈希函数仅支持 int64/int32/int/string/[]byte，其他类型都映射到分片 0
5. **Range 一致性**：`Range` 不保证遍历期间数据不变，可能看到或遗漏中间状态

***

## zConcurrency 并发控制

zConcurrency 提供六种并发控制工具，覆盖工作池、信号量、限流、互斥锁、计数器等场景。

### 实现方式

| 工具 | 底层实现 | 适用场景 |
|------|----------|----------|
| WorkerPool | context.Cancel + chan Task | 消息处理、任务调度 |
| Semaphore | sync.Mutex + chan struct{} | 资源限流、连接池 |
| RateLimiter | 令牌桶 + atomic.Int64 | API 限流、流量控制 |
| MutexWithTimeout | sync.Mutex + chan struct{} | 超时加锁、防死锁 |
| AtomicCounter | atomic.Int64 | 计数统计、并发计数 |
| WaitGroupWithContext | sync.WaitGroup + []error | 并发等待、错误收集 |

### 使用说明

#### WorkerPool

```go
pool := zConcurrency.NewWorkerPool(4, 100)
pool.Start()

err := pool.Submit(func() error {
    return nil
})

err = pool.SubmitWithContext(ctx, func() error {
    return nil
})

pool.Stop()
pool.Wait()
errors := pool.Errors()
```

#### Semaphore

```go
sem := zConcurrency.NewSemaphore(5)
sem.Acquire()
defer sem.Release()

err := sem.AcquireWithContext(ctx, 1)
available := sem.AvailablePermits()
```

#### RateLimiter

```go
limiter := zConcurrency.NewRateLimiter(100, 10)

if limiter.Allow() { /* 可以执行 */ }
if limiter.AllowN(5) { /* 消耗 5 个令牌 */ }
limiter.Wait(ctx)
limiter.Reset(200, 20)
```

#### MutexWithTimeout

```go
mu := zConcurrency.NewMutexWithTimeout()
mu.Lock()
mu.Unlock()

ok := mu.LockWithTimeout(time.Second * 5)
ok = mu.LockWithContext(ctx)
```

#### AtomicCounter

```go
counter := zConcurrency.NewAtomicCounter(0)
counter.Increment()
counter.Decrement()
counter.Add(10)
val := counter.Get()
swapped := counter.CompareAndSwap(100, 200)
```

#### WaitGroupWithContext

```go
wg := zConcurrency.NewWaitGroupWithContext()
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := doWork(); err != nil {
            wg.RecordError(err)
        }
    }()
}
err := wg.WaitWithContext(ctx)
errors := wg.Errors()
```

### 注意事项

1. **WorkerPool 必须启动**：`Submit` 前必须调用 `Start()`
2. **WorkerPool 队列满**：`Submit` 会阻塞，使用 `SubmitWithContext` 避免无限等待
3. **Semaphore 死锁**：`Acquire(n)` 中 n 不能超过初始 permits
4. **RateLimiter 轮询**：`Wait` 使用 1ms 轮询，高频场景考虑优化
5. **MutexWithTimeout panic**：未 Lock 时调用 Unlock 会 panic

***

## zCrypto 加密工具

zCrypto 提供完整的加密解密工具集，涵盖对称加密、非对称加密、哈希、编码等。

### 功能矩阵

| 功能 | 方法 | 说明 |
|------|------|------|
| AES-GCM | AESEncrypt/AESDecrypt | 推荐模式，认证加密 |
| AES-CBC | AESEncrypt/AESDecrypt | PKCS#7 填充 |
| RSA | RSAEncrypt/RSADecrypt | PKCS1v15 填充 |
| RSA 签名 | RSASign/RSAVerify | SHA256 哈希 |
| 哈希 | MD5/SHA1/SHA256/SHA512 | 十六进制输出 |
| 编码 | Base64Encode/HexEncode | 双向编解码 |
| 随机数 | GenerateRandomString/Bytes | crypto/rand |

### 使用说明

```go
// AES-GCM 加密（推荐）
key, _ := zCrypto.GenerateAESKey(zCrypto.AESKeySize32)
encrypted, _ := zCrypto.AESEncrypt(plaintext, key, nil, zCrypto.AESModeGCM)
decrypted, _ := zCrypto.AESDecrypt(encrypted, key, nil, zCrypto.AESModeGCM)

// RSA 加密
privateKey, publicKey, _ := zCrypto.GenerateRSAKeyPair(2048)
encrypted, _ := zCrypto.RSAEncrypt(plaintext, publicKey)
decrypted, _ := zCrypto.RSADecrypt(encrypted, privateKey)

// 哈希
md5Hash := zCrypto.MD5("hello")
sha256Hash := zCrypto.SHA256("hello")

// 安全随机
randomStr, _ := zCrypto.GenerateRandomString(32)
```

### 注意事项

1. **GCM vs CBC**：推荐使用 GCM 模式，提供认证加密
2. **IV 唯一性**：GCM 模式下 IV 绝不能重复使用相同密钥加密
3. **RSA 数据大小**：RSA 加密数据不能超过密钥长度减去填充开销
4. **ECB 不安全**：ECB 模式仅用于兼容旧系统，不推荐新代码使用

***

## zCache 缓存

zCache 提供三种缓存实现：LRU 缓存、简单缓存和 TTL 自动过期缓存。

### 实现对比

| 特性 | LRUCache | SimpleCache | TTLCache |
|------|----------|-------------|----------|
| 淘汰策略 | LRU | 无 | 无 |
| 过期支持 | ✅ | ✅ | ✅ |
| 后台清理 | ✅ | ❌ | ✅ |
| 容量限制 | ✅ | ❌ | ❌ |
| 性能 | 中 | 高 | 高 |

### 使用说明

```go
// LRU 缓存
lru := zCache.NewLRUCache(1000, time.Hour)
lru.Set("key", "value", time.Hour)
val, err := lru.Get("key")
lru.Delete("key")
lru.Stop()

// 简单缓存
simple := zCache.NewSimpleCache(time.Hour)
simple.Set("key", "value", time.Hour)

// TTL 缓存
ttl := zCache.NewTTLCache(time.Minute * 30)
ttl.Set("key", "value", time.Minute*30)
ttl.Stop()

// 全局默认缓存
zCache.SetDefault("key", "value", time.Hour)
val, err := zCache.GetDefault("key")
```

### 注意事项

1. **清理 goroutine**：LRUCache 和 TTLCache 会启动后台 goroutine，必须调用 `Stop()` 释放
2. **LRU 容量**：容量满时淘汰最久未访问的项
3. **SimpleCache 性能**：`Len()/Keys()` 每次调用都会清理过期项

***

## zKeyWordFilter 敏感词过滤

zKeyWordFilter 基于 DFA（确定有限自动机）算法实现高效敏感词过滤。

### 使用说明

```go
filter := zKeyWordFilter.NewFilter()
filter.AddWord("敏感词")
filter.AddWord("违禁")

result := filter.Filter("这是一条包含敏感词的文本")
// result = "这是一条包含***的文本"

zKeyWordFilter.InitDefaultFilter()
zKeyWordFilter.AddWord("敏感词")
result = zKeyWordFilter.Filter("包含敏感词的文本")

err := zKeyWordFilter.ParseFromFile("sensitive_words.txt")
```

### 注意事项

1. **线程安全**：`AddWord` 和 `Filter` 均为线程安全，使用 `sync.RWMutex` 保护，支持并发读取
2. **性能**：DFA 算法时间复杂度 O(n)，与敏感词数量无关
3. **中文支持**：按 rune 处理，天然支持中文、日文等多字节字符

***

## zList 线程安全链表

zList 基于 `container/list` 封装线程安全双向链表，支持安全遍历删除。

```go
l := zList.New()
l.PushFront("head")
l.PushBack("tail")
front := l.Front()
back := l.Back()

// 遍历（回调中不能调用 Remove，会导致死锁）
l.Range(func(e *list.Element, value any) bool {
    fmt.Println(value)
    return true
})

// 安全遍历并删除（推荐）
l.RangeWithDelete(func(e *list.Element, value any) (bool, bool) {
    if value.(int) > 3 {
        return true, true  // 继续遍历，删除当前元素
    }
    return true, false     // 继续遍历，保留当前元素
})

l.Remove(front)
l.Clear()  // 清空链表
```

***

## zQueue 队列

zQueue 提供环形队列实现，满时自动 2 倍扩容。

```go
rq := zQueue.NewRingQueue(100)
rq.Enqueue("item1")
rq.Enqueue("item2")
val, ok := rq.Dequeue()
val, ok = rq.Peek()
length := rq.Len()
capacity := rq.Cap()
rq.Clear()
```

***

## zStack 栈

zStack 提供固定大小的栈实现和线程安全的并发栈。

```go
// 非线程安全栈
s := zStack.New(100)
s.Push("item1")
s.Push("item2")
val, err := s.Pop()
val, err = s.Peek()
all := s.Get()
s.Empty()

// 线程安全并发栈
cs := zStack.NewConcurrent(100)
cs.Push("item1")
val, err := cs.Pop()
val, err = cs.Peek()
all := cs.Get()
length := cs.Len()
cs.Empty()
```

***

## zTree 树形结构

zTree 提供从扁平数据生成树形结构的功能。

```go
type Department struct {
    ID       int
    ParentID int
    Name     string
}

func (d *Department) GetId() interface{}      { return d.ID }
func (d *Department) GetFatherId() interface{} { return d.ParentID }
func (d *Department) GetData() interface{}     { return d }
func (d *Department) IsRoot() bool             { return d.ParentID == 0 }

nodes := []zTree.INode{
    &Department{ID: 1, ParentID: 0, Name: "总部"},
    &Department{ID: 2, ParentID: 1, Name: "技术部"},
}
trees := zTree.GenerateTree(nodes)
```

***

## zStr 字符串处理

zStr 提供 40+ 个字符串处理函数，全面支持中文。

```go
// 截取
sub := zStr.Substring("你好世界", 0, 2)     // "你好"
truncated := zStr.Truncate("hello world", 5) // "hello..."

// 包含检查
zStr.Contains("hello", "ell")
zStr.ContainsAny("hello", "a,b,c")
zStr.ContainsAll("hello world", "hello,world")

// 替换
zStr.Replace("hello", "l", "L", 1)
zStr.ReplaceAll("hello", "l", "L")
zStr.ReplaceRegexp("hello123", `\d`, "X")

// 分割/连接
parts := zStr.Split("a,b,c", ",")
joined := zStr.Join([]string{"a", "b"}, ",")

// 修剪
zStr.Trim("  hello  ")
zStr.TrimLeft("hello", "he")

// 大小写
zStr.CamelCase("hello_world")
zStr.SnakeCase("HelloWorld")
zStr.KebabCase("HelloWorld")

// 其他
zStr.IsEmail("test@example.com")
zStr.IsURL("https://example.com")
zStr.Reverse("hello")
zStr.PadLeft("5", 3, "0") // "005"
```

***

## zTime 时间处理

zTime 提供时间包装器和转换工具。

```go
now := zTime.Now()
timestamp := zTime.NowUnix()
dateTime := zTime.FormatDateTime(now)
dateOnly := zTime.FormatDate(now)

// 时间转换
t, _ := zTime.ParseDateTime("2024-01-01 12:00:00")
t, _ = zTime.ParseDate("2024-01-01")
duration := zTime.ParseDuration("1h30m")

// 时间计算
tomorrow := zTime.AddDay(now, 1)
startOfDay := zTime.StartOfDay(now)
endOfDay := zTime.EndOfDay(now)
```

***

## zRand 随机数

zRand 提供随机数和安全随机数生成。

```go
// 普通随机
val := zRand.Int(100)
val64 := zRand.Int64(1000)
str := zRand.String(16)

// 范围随机
val := zRand.Range(10, 100)

// 安全随机（crypto/rand）
secureStr := zRand.SecureString(32)
secureBytes := zRand.SecureBytes(16)
```

***

## zReflect 反射工具

zReflect 提供反射工具和深拷贝功能。

```go
// 深拷贝
copied := zReflect.DeepCopy(original)

// 类型判断
isStruct := zReflect.IsStruct(obj)
isSlice := zReflect.IsSlice(obj)

// 字段操作
fields := zReflect.GetStructFields(obj)
val := zReflect.GetFieldValue(obj, "Name")
zReflect.SetFieldValue(obj, "Name", "new_name")
```

***

## zFile 文件操作

zFile 提供文件读写和目录操作工具。

```go
// 文件读写
content, err := zFile.ReadFile("config.json")
err := zFile.WriteFile("output.json", data)

// 目录操作
exists := zFile.DirExists("logs")
err := zFile.Mkdir("logs")
files, err := zFile.ListFiles("logs")

// 路径操作
absPath := zFile.AbsPath("config.json")
ext := zFile.Ext("config.json")
base := zFile.BaseName("config.json")
```

***

## zHashtable 哈希表

zHashtable 提供自动扩容的哈希表实现。

```go
ht := zHashtable.NewHashtable()
ht.Put("key", "value")
val, ok := ht.Get("key")
ht.Remove("key")
size := ht.Size()
ht.Clear()
```

***

## zDataConv 数据类型转换

zDataConv 提供安全的数据类型转换函数。

```go
// 数值转换
i := zDataConv.ToInt("123")
i64 := zDataConv.ToInt64("1234567890")
f64 := zDataConv.ToFloat64("3.14")
b := zDataConv.ToBool("true")

// 安全转换（带默认值）
i := zDataConv.ToIntWithDefault("abc", 0)
str := zDataConv.ToString(123)
```

***

## zError 带错误码的错误类型

zError 提供带错误码的错误类型，实现标准 `error` 接口。

```go
err := zError.New("player not found")
code := err.GetCode()       // 0
msg := err.GetMessage()     // "player not found"

err = zError.NewWithCode(1001, "player not found")
code = err.GetCode()        // 1001
msg = err.GetMessage()      // "player not found"

err = zError.Errorf("player %d not found", 1001)
```

***

## zUtils 通用工具

zUtils 提供 panic 恢复和工作目录等通用工具。

```go
// Panic 恢复（推荐用法）
var err error
defer zUtils.RecoverToErr(&err)

// 工作目录
dir, err := zUtils.GetCurrentDirectory()
```

***

## zColor 终端颜色

zColor 提供 ANSI 终端颜色输出。

```go
fmt.Println(zColor.Red("Error message"))
fmt.Println(zColor.Green("Success message"))
fmt.Println(zColor.Yellow("Warning message"))
fmt.Println(zColor.Blue("Info message"))
fmt.Println(zColor.Bold("Bold text"))
```

***

## zGps GPS 坐标转换

zGps 提供 GPS 坐标系转换工具（WGS84/GCJ02/BD09）。

```go
// WGS84 转 GCJ02（火星坐标系）
gcjLng, gcjLat := zGps.WGS84ToGCJ02(116.3912, 39.9073)

// GCJ02 转 BD09（百度坐标系）
bdLng, bdLat := zGps.GCJ02ToBD09(116.3974, 39.9089)

// BD09 转 GCJ02
gcjLng, gcjLat = zGps.BD09ToGCJ02(116.4038, 39.9155)

// WGS84 转 BD09
bdLng, bdLat = zGps.WGS84ToBD09(116.3912, 39.9073)
```

## 代码质量评估

### 整体评分

| 维度 | 评分 | 说明 |
|------|------|------|
| **功能完整性** | ⭐⭐⭐⭐ | 20 个模块覆盖常用工具场景，部分模块功能偏简单 |
| **代码质量** | ⭐⭐⭐⭐ | 整体代码规范，命名清晰，部分模块存在并发 Bug |
| **测试覆盖** | ⭐⭐⭐ | 18 个测试文件/3464 行，核心模块覆盖充分，部分模块缺失或不足 |
| **API 设计** | ⭐⭐⭐⭐ | 接口简洁直观，与标准库风格一致，泛型支持良好 |
| **文档完善度** | ⭐⭐⭐⭐ | README 详细，每个模块有使用说明和注意事项 |

### 代码统计

| 指标 | 数值 |
|------|------|
| 源码文件 | 29 |
| 源码行数 | 3,537 |
| 测试文件 | 18 |
| 测试行数 | 3,464 |
| 测试/代码比 | 97.9% |
| 外部依赖 | 0（仅标准库） |

### 各模块质量详情

| 模块 | 行数 | 测试 | 测试状态 | 质量评级 | 备注 |
|------|------|------|----------|----------|------|
| zMap | 357 | ✅ 充分 | PASS | A | 核心模块，4种Map实现，泛型支持 |
| zConcurrency | 383 | ✅ 较充分 | PASS | A | 6种并发工具，功能完整 |
| zCrypto | 387 | ✅ 充分 | PASS | A | AES/RSA/哈希/编码，功能完整 |
| zCache | 380 | ✅ 充分 | PASS | A | LRU/Simple/TTL 三种缓存 |
| zStr | 469 | ✅ 充分 | PASS | A | 40+ 函数，中文支持 |
| zReflect | 263 | ✅ 充分 | PASS | A | 深拷贝/字段操作 |
| zFile | 184 | ✅ 充分 | PASS | A | 文件读写/目录操作 |
| zHashtable | 129 | ✅ 充分 | PASS | B+ | 自动扩容哈希表 |
| zKeyWordFilter | 149 | ✅ 充分 | PASS | A- | DFA 算法正确，线程安全 |
| zQueue | 163 | ⚠️ 不足 | PASS | B | 环形队列，功能偏简单 |
| zTime | 125 | ⚠️ 不足 | PASS | B | 时间包装器，缺少时区处理 |
| zGps | 137 | ⚠️ 不足 | PASS | B | 坐标转换正确 |
| zDataConv | 59 | ⚠️ 不足 | PASS | B- | 类型转换，功能偏简单 |
| zRand | 76 | ⚠️ 不足 | PASS | B | 随机数，缺少加权随机等 |
| zColor | 63 | ⚠️ 不足 | PASS | B- | 终端颜色，功能简单 |
| zTree | 46 | ⚠️ 基本 | PASS | B- | 树形生成，仅支持 INode 接口 |
| zStack | 115 | ⚠️ 不足 | PASS | B | 固定大小栈 + 并发安全栈 |
| zList | 115 | ✅ 充分 | PASS | A- | 线程安全链表，支持安全遍历删除 |
| zError | 50 | ✅ 充分 | PASS | B+ | 带错误码错误类型 |
| zUtils | 45 | ✅ 充分 | PASS | B+ | panic 恢复/工作目录 |

### 已知问题

1. **zQueue**：环形队列缺少并发安全版本
2. **zTree**：仅支持 INode 接口，缺少泛型版本
3. **zTime**：缺少时区处理功能
4. **zRand**：缺少加权随机等高级随机功能

### 更新日志

#### v0.2.0 (2026-04)

- **zList**：修复 Range 中 Remove 死锁 Bug，新增 `RangeWithDelete` 安全遍历删除方法和 `Clear` 方法
- **zKeyWordFilter**：添加 `sync.RWMutex` 线程安全保护，支持并发 AddWord 和 Filter
- **zStack**：新增 `ConcurrentStack` 线程安全并发栈
- **zError**：新增单元测试
- **zUtils**：新增 `RecoverToErr` 方法（正确支持 defer 直接调用），新增单元测试
- **zList**：重写测试，使用 testify 断言，增加并发测试

### 改进计划

#### P1 - 质量提升

- [ ] zQueue 添加并发安全版本
- [ ] zTree 添加泛型版本

#### P2 - 功能增强

- [ ] zRand 添加加权随机（WeightedRandom）、洗牌（Shuffle）等
- [ ] zTime 添加时区处理、工作日计算等
- [ ] zQueue 添加阻塞队列（BlockingQueue）、优先级队列
- [ ] zDataConv 添加更多类型转换（uint 系列、[]byte ↔ hex 等）
- [ ] zCache 添加统计接口（命中率、淘汰数等）

#### P3 - 长期优化

- [ ] 考虑添加泛型优先队列（zPriorityQueue）
- [ ] 考虑添加布隆过滤器（zBloomFilter）
- [ ] 考虑添加跳表（zSkipList）
- [ ] 考虑添加限流器增强版（滑动窗口限流、漏桶限流）

## 许可证

MIT License
