# zUtil

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

zUtil 是一个 Go 语言通用工具集，提供并发安全的数据结构、加密、缓存、配置管理等基础功能，适用于游戏服务器及各类 Go 项目。

## 项目结构

```
zUtil/
├── zMap/            # 并发安全 Map - TypedMap/TypedShardedMap（泛型）
├── zCache/          # 缓存 - LRU/Simple/TTL 三种实现
├── zConcurrency/    # 并发控制 - WorkerPool/Semaphore/RateLimiter/AtomicCounter
├── zCrypto/         # 加密 - AES(CBC/GCM)/RSA/Hash/Base64
├── zKeyWordFilter/  # 敏感词过滤 - DFA 算法
├── zList/           # 线程安全双向链表
├── zQueue/          # 队列 - 链式队列 + 环形队列（自动扩容）
├── zStr/            # 字符串处理 - 中文支持、命名转换、编辑距离
├── zTime/           # 时间处理 - 时区转换、区间计算
├── zFile/           # 文件操作 - 读写/复制/移动/目录管理
├── zDataConv/       # 数据类型转换
├── zReflect/        # 反射工具 - 字段操作/深拷贝/方法调用
├── zRand/           # 随机数 - math/rand/v2 + crypto/rand
├── zStack/          # 固定大小栈
├── zHashtable/      # 哈希表 - 支持动态扩容
├── zTree/           # 树形结构生成
├── zGps/            # GPS 坐标系转换 - WGS84/GCJ02/BD09
├── zColor/          # 终端颜色输出
├── zError/          # 带错误码的错误类型
└── zUtils/          # 通用工具 - 目录获取/错误恢复
```

## 核心模块

### zMap - 并发安全 Map

基于 `sync.Map` 的扩展，提供类型安全的泛型实现。是 zUtil 中最核心、使用最广泛的模块。

#### 类型对比

| 类型 | 适用场景 | 说明 |
|------|----------|------|
| `TypedMap[K, V]` | 读多写少 | 基于 sync.Map，零类型断言 |
| `TypedShardedMap[K, V]` | 高并发写入 | 分片减少锁竞争 |
| `Map` | 需要动态类型 | 基于 sync.Map + atomic 计数 |
| `ShardedMap` | 高并发动态类型 | 分片版本 |

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zMap"
)

func main() {
    // TypedMap - 泛型类型安全 Map（推荐）
    tm := zMap.NewTypedMap[string, int]()
    tm.Store("player_001", 100)
    v, ok := tm.Load("player_001")  // v 是 int 类型，无需断言
    fmt.Println(v, ok)  // 100 true

    // TypedShardedMap - 高并发写入场景
    tsm := zMap.NewTypedShardedMap[string, int](32)
    tsm.Store("key1", 1)
    tsm.Store("key2", 2)
    fmt.Println(tsm.Len())  // 2

    // Map - 动态类型 Map
    m := zMap.NewMap()
    m.Store("name", "Alice")
    m.Store("age", 30)
    fmt.Println(m.Len())  // 2
    m.Range(func(key, value interface{}) bool {
        fmt.Printf("%v: %v\n", key, value)
        return true
    })

    // LoadOrStore - 原子性加载或存储
    existing, loaded := tm.LoadOrStore("player_001", 200)
    fmt.Println(existing, loaded)  // 100 true（已存在，不覆盖）

    // LoadAndDelete - 原子性加载并删除
    val, ok := tm.LoadAndDelete("player_001")
    fmt.Println(val, ok)  // 100 true
    fmt.Println(tm.Len())  // 0
}
```

#### 统一接口

所有 Map 类型提供一致的接口：

| 方法 | 说明 |
|------|------|
| `Load(key) (value, bool)` | 获取值 |
| `Store(key, value)` | 存储值 |
| `Delete(key)` | 删除值 |
| `Len() int64` | 元素数量 |
| `Range(f)` | 遍历所有元素 |
| `Clear()` | 清空所有元素 |
| `LoadOrStore(key, value) (value, bool)` | 加载或存储 |
| `LoadAndDelete(key) (value, bool)` | 加载并删除 |
| `CompareAndDelete(key, old) bool` | 比较并删除 |
| `CompareAndSwap(key, old, new) bool` | 比较并替换 |

---

### zCache - 缓存

提供 LRU、Simple、TTL 三种缓存实现。

| 类型 | 淘汰策略 | 过期清理 | 适用场景 |
|------|----------|----------|----------|
| `LRUCache` | LRU | 支持 | 需要淘汰策略的缓存 |
| `SimpleCache` | 无淘汰 | 支持 | 简单缓存 |
| `TTLCache` | 无淘汰 | 自动清理 | 需要自动过期的缓存 |

#### 使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/pzqf/zUtil/zCache"
)

func main() {
    // LRU 缓存 - 容量100，过期时间5分钟
    lru := zCache.NewLRUCache(100, 5*time.Minute)
    lru.Set("user:1001", "Alice")
    val, ok := lru.Get("user:1001")
    fmt.Println(val, ok)  // Alice true

    // TTL 缓存 - 自动过期清理
    ttl := zCache.NewTTLCache(10 * time.Second)
    ttl.Set("token", "abc123")
    time.Sleep(11 * time.Second)
    _, ok = ttl.Get("token")
    fmt.Println(ok)  // false（已过期）

    // 全局默认缓存
    zCache.SetDefault("key1", "value1")
    val, ok = zCache.GetDefault("key1")
    fmt.Println(val, ok)  // value1 true
}
```

---

### zConcurrency - 并发控制

提供工作池、信号量、限流器等并发控制工具。

| 类型 | 用途 |
|------|------|
| `WorkerPool` | 工作池，管理并发任务 |
| `Semaphore` | 信号量，限制并发访问 |
| `RateLimiter` | 令牌桶限流器 |
| `AtomicCounter` | 原子计数器 |
| `MutexWithTimeout` | 带超时的互斥锁 |
| `WaitGroupWithContext` | 带上下文的 WaitGroup |

#### 使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/pzqf/zUtil/zConcurrency"
)

func main() {
    // WorkerPool - 工作池
    pool := zConcurrency.NewWorkerPool(4, 100)  // 4个worker，队列100
    for i := 0; i < 10; i++ {
        task := zConcurrency.Task(func() error {
            time.Sleep(100 * time.Millisecond)
            return nil
        })
        pool.AddTask(task)
    }
    pool.Stop()

    // Semaphore - 信号量
    sem := zConcurrency.NewSemaphore(3)  // 最多3个并发
    sem.Acquire()
    // ... 执行受限操作
    sem.Release()

    // RateLimiter - 令牌桶限流
    limiter := zConcurrency.NewRateLimiter(100, 10)  // 100/s，桶容量10
    if limiter.Allow() {
        fmt.Println("请求通过")
    }

    // AtomicCounter - 原子计数器
    counter := zConcurrency.NewAtomicCounter(0)
    counter.Increment()
    counter.Add(10)
    fmt.Println(counter.Get())  // 11
}
```

---

### zCrypto - 加密工具

提供 AES/RSA/Hash/编码等加密功能。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zCrypto"
)

func main() {
    // AES-GCM 加密（推荐）
    key := zCrypto.GenerateAESKey(32)  // AES-256
    plaintext := []byte("Hello, World!")
    ciphertext, err := zCrypto.AESEncrypt(plaintext, key, zCrypto.AESModeGCM)
    if err != nil {
        panic(err)
    }
    decrypted, err := zCrypto.AESDecrypt(ciphertext, key, zCrypto.AESModeGCM)
    fmt.Println(string(decrypted))  // Hello, World!

    // RSA 密钥对
    privateKey, publicKey := zCrypto.GenerateRSAKeyPair(2048)
    encrypted, _ := zCrypto.RSAEncrypt(plaintext, publicKey)
    decrypted, _ = zCrypto.RSADecrypt(encrypted, privateKey)
    fmt.Println(string(decrypted))  // Hello, World!

    // 哈希
    hash := zCrypto.SHA256("password")
    fmt.Println(hash)

    // Base64
    encoded := zCrypto.Base64Encode(plaintext)
    decoded, _ := zCrypto.Base64Decode(encoded)
    fmt.Println(string(decoded))  // Hello, World!
}
```

---

### zKeyWordFilter - 敏感词过滤

基于 DFA（确定有限自动机）算法的高效敏感词过滤。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zKeyWordFilter"
)

func main() {
    // 创建过滤器
    filter := zKeyWordFilter.NewFilter()
    filter.AddWord("敏感词")
    filter.AddWord("违规")
    result := filter.Filter("这是一条包含敏感词的消息")
    fmt.Println(result)  // 这是一条包含***的消息

    // 使用全局默认过滤器
    zKeyWordFilter.InitDefaultFilter()
    zKeyWordFilter.AddWord("违禁")
    zKeyWordFilter.AddWord("禁止")
    result = zKeyWordFilter.Filter("违禁内容禁止传播")
    fmt.Println(result)  // **内容**传播

    // 从文件加载敏感词（每行一个词）
    zKeyWordFilter.ParseFromFile("sensitive_words.txt")
}
```

---

### zList - 线程安全双向链表

基于 `container/list` 的线程安全封装。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zList"
)

func main() {
    list := zList.New()
    list.PushBack("first")
    list.PushBack("second")
    list.PushFront("zero")

    fmt.Println(list.Len())  // 3

    // 遍历
    list.Range(func(e *list.Element, value any) bool {
        fmt.Println(value)
        return true
    })
    // 输出: zero first second

    // 获取首尾元素
    front := list.Front()
    fmt.Println(front.Value)  // zero
}
```

---

### zQueue - 队列

提供链式队列和环形队列两种实现。

| 类型 | 特点 | 适用场景 |
|------|------|----------|
| `Queue` | 链式队列，线程安全 | 通用队列 |
| `RingQueue` | 环形队列，自动扩容 | 高性能固定大小场景 |

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zQueue"
)

func main() {
    // 链式队列
    q := zQueue.NewQueue()
    q.Enqueue("first")
    q.Enqueue("second")
    q.Enqueue("third")
    fmt.Println(q.Length())  // 3

    val, _ := q.Dequeue()
    fmt.Println(val)  // first

    val, _ = q.Peek()
    fmt.Println(val)  // second

    // 环形队列
    rq := zQueue.NewRingQueue(1024)
    rq.Enqueue("item1")
    rq.Enqueue("item2")
    fmt.Println(rq.Len(), rq.Cap())  // 2 1024

    val, _ = rq.Dequeue()
    fmt.Println(val)  // item1
}
```

---

### zStr - 字符串处理

提供丰富的字符串处理函数，完整支持中文。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zStr"
)

func main() {
    // 中文子字符串
    s := "你好，世界！"
    fmt.Println(zStr.Substring(s, 0, 2))  // 你好

    // 命名转换
    fmt.Println(zStr.CamelCase("user_name"))   // userName
    fmt.Println(zStr.SnakeCase("UserName"))     // user_name
    fmt.Println(zStr.KebabCase("UserName"))     // user-name

    // 中文反转
    fmt.Println(zStr.Reverse("你好世界"))  // 界世好你

    // 编辑距离（支持中文）
    fmt.Println(zStr.Distance("你好", "你们好"))  // 1

    // 随机字符串
    fmt.Println(zStr.Random(16, ""))  // 随机16位字母数字
    fmt.Println(zStr.RandomDigit(6))  // 随机6位数字

    // 检查
    fmt.Println(zStr.IsEmpty("  "))       // true
    fmt.Println(zStr.IsAlpha("Hello"))    // true
    fmt.Println(zStr.IsNumeric("123"))    // true
    fmt.Println(zStr.IsAlphanumeric("abc123"))  // true

    // 首字母操作
    fmt.Println(zStr.Capitalize("hello"))    // Hello
    fmt.Println(zStr.Decapitalize("Hello"))  // hello

    // 填充
    fmt.Println(zStr.PadLeft("42", 5, "0"))  // 00042
    fmt.Println(zStr.PadRight("42", 5, "0")) // 42000

    // 截断
    fmt.Println(zStr.Truncate("你好世界再见", 3, "..."))  // 你好世...
}
```

---

### zTime - 时间处理

默认中国时区（CST, UTC+8），支持时间区间计算。

#### 使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/pzqf/zUtil/zTime"
)

func main() {
    // 当前时间
    now := zTime.Now()
    fmt.Println(now.Time())

    // 从 time.Time 创建
    t := time.Date(2026, 4, 14, 15, 30, 0, 0, time.UTC)
    zt := zTime.New(t)
    fmt.Println(zt.Time())

    // 从时间戳创建
    zt = zTime.NewFromSeconds(1713090600)

    // 从字符串创建
    zt = zTime.NewFromString("2006-01-02 15:04:05", "2026-04-14 15:30:00")

    // 时间区间
    fmt.Println(now.BeginOfDay())    // 当天 00:00:00
    fmt.Println(now.EndOfDay())      // 当天 23:59:59
    fmt.Println(now.BeginOfWeek())   // 本周一 00:00:00
    fmt.Println(now.EndOfWeek())     // 本周日 23:59:59
    fmt.Println(now.BeginOfMonth())  // 本月1日 00:00:00
    fmt.Println(now.EndOfMonth())    // 本月最后一天 23:59:59
    fmt.Println(now.BeginOfYear())   // 今年1月1日 00:00:00
    fmt.Println(now.EndOfYear())     // 今年12月31日 23:59:59

    // 时区设置
    zt.SetZone("EST", -5)  // 美东时间

    // 时间转换
    seconds := zTime.Time2Seconds(now.Time())
    str := zTime.Time2String(now.Time())
    fmt.Println(seconds, str)
}
```

---

### zFile - 文件操作

提供丰富的文件和目录操作函数。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zFile"
)

func main() {
    // 文件读写
    zFile.WriteFile("test.txt", []byte("Hello"))
    data, _ := zFile.ReadFile("test.txt")
    fmt.Println(string(data))  // Hello

    // 追加写入
    zFile.AppendFile("test.txt", []byte(" World"))

    // 目录操作
    zFile.CreateDirAll("./data/logs")
    files, _ := zFile.ListFiles("./data")

    // 文件信息
    size, _ := zFile.GetFileSize("test.txt")
    modTime, _ := zFile.GetFileModTime("test.txt")
    exists := zFile.Exists("test.txt")
    isDir := zFile.IsDir("./data")

    fmt.Println(size, modTime, exists, isDir)

    // 复制/移动
    zFile.CopyFile("test.txt", "test_copy.txt")
    zFile.MoveFile("test_copy.txt", "test_moved.txt")
}
```

---

### zDataConv - 数据类型转换

提供常用数据类型之间的转换函数。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zDataConv"
)

func main() {
    // 字符串转数字
    i := zDataConv.String2Int("123")
    i64 := zDataConv.String2Int64("9999999999")
    f32 := zDataConv.String2Float32("3.14")
    f64 := zDataConv.String2Float64("3.141592653589793")
    fmt.Println(i, i64, f32, f64)

    // 数字转字符串
    s := zDataConv.Int2String(123)
    s64 := zDataConv.Int642String(9999999999)
    fmt.Println(s, s64)

    // 字节和字符串互转（零拷贝）
    bytes := zDataConv.String2Byte("hello")
    str := zDataConv.Byte2String(bytes)
    fmt.Println(str)

    // Base64 编解码
    encoded := zDataConv.Base64Encode("hello")
    fmt.Println(encoded)

    // URL 编解码
    urlEncoded := zDataConv.UrlEncode("hello world")
    fmt.Println(urlEncoded)
}
```

---

### zReflect - 反射工具

提供反射相关的工具函数，简化反射操作。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zReflect"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func (u *User) Greet() string {
    return "Hello, " + u.Name
}

func main() {
    user := &User{Name: "Alice", Age: 30}

    // 获取字段值
    name, _ := zReflect.GetFieldValue(user, "Name")
    fmt.Println(name)  // Alice

    // 设置字段值
    zReflect.SetFieldValue(user, "Age", 31)

    // 获取结构体字段标签
    tags, _ := zReflect.GetStructFieldTags(user, "json")
    fmt.Println(tags)  // map[Name:name Age:age]

    // 调用方法
    result, _ := zReflect.CallMethod(user, "Greet")
    fmt.Println(result[0].String())  // Hello, Alice

    // 深拷贝
    copy, _ := zReflect.DeepCopy(user)
    fmt.Println(copy.(*User).Name)  // Alice

    // 转为 Map
    m, _ := zReflect.ToMapWithTags(user, "json")
    fmt.Println(m)  // map[name:Alice age:31]
}
```

---

### zRand - 随机数

提供基于 `math/rand/v2` 和 `crypto/rand` 的随机数生成。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zRand"
)

func main() {
    // 普通随机数
    n := zRand.RandN(100)          // [0, 100)
    interval := zRand.RandInterval(10, 20)  // [10, 20]
    fmt.Println(n, interval)

    // 随机字符串
    s := zRand.RandString(16)      // 16位随机字符串
    fmt.Println(s)

    // 自定义字符集
    custom := zRand.RandCustomString(8, "ABCDEF0123456789")
    fmt.Println(custom)  // 8位16进制字符串

    // 安全随机数（crypto/rand）
    secureInt, _ := zRand.SecureInt(100)
    secureStr, _ := zRand.SecureString(32)
    secureBytes, _ := zRand.SecureBytes(16)
    fmt.Println(secureInt, secureStr, secureBytes)
}
```

---

### zStack - 固定大小栈

固定大小的栈实现。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zStack"
)

func main() {
    stack := zStack.New(10)  // 容量10
    stack.Push("first")
    stack.Push("second")

    val, _ := stack.Pop()
    fmt.Println(val)  // second

    val, _ = stack.Peek()
    fmt.Println(val)  // first

    fmt.Println(stack.IsEmpty())  // false
}
```

---

### zHashtable - 哈希表

支持动态扩容的哈希表实现。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zHashtable"
)

func main() {
    ht := zHashtable.NewHashTable()
    ht.Add("key1", "value1")
    ht.Add("key2", "value2")

    val, _ := ht.Get("key1")
    fmt.Println(val)  // value1

    ht.Set("key1", "new_value")
    val, _ = ht.Get("key1")
    fmt.Println(val)  // new_value

    fmt.Println(ht.Size())  // 2
    ht.Remove("key2")
    fmt.Println(ht.Size())  // 1
}
```

---

### zTree - 树形结构

从扁平数据生成树形结构。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zTree"
)

type MenuItem struct {
    ID       int
    ParentID int
    Name     string
}

func (m MenuItem) GetId() int        { return m.ID }
func (m MenuItem) GetFatherId() int  { return m.ParentID }
func (m MenuItem) GetData() interface{} { return m }
func (m MenuItem) IsRoot() bool      { return m.ParentID == 0 }

func main() {
    items := []zTree.INode{
        MenuItem{ID: 1, ParentID: 0, Name: "Root"},
        MenuItem{ID: 2, ParentID: 1, Name: "Child1"},
        MenuItem{ID: 3, ParentID: 1, Name: "Child2"},
    }
    trees := zTree.GenerateTree(items)
    fmt.Println(len(trees))  // 1（一个根节点）
    fmt.Println(len(trees[0].Children))  // 2（两个子节点）
}
```

---

### zGps - GPS 坐标系转换

支持 WGS84/GCJ02/BD09 三种坐标系之间的转换。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zGps"
)

func main() {
    // WGS84 转 GCJ02（GPS坐标转国测局坐标）
    lng, lat := zGps.Gps84ToGcj02(116.397428, 39.90923)
    fmt.Printf("GCJ02: %.6f, %.6f\n", lng, lat)

    // GCJ02 转 BD09（国测局坐标转百度坐标）
    lng, lat = zGps.Gcj02ToBd09(lng, lat)
    fmt.Printf("BD09: %.6f, %.6f\n", lng, lat)

    // 计算两点距离（米）
    distance := zGps.GetDistance(39.90923, 116.397428, 39.90950, 116.39780)
    fmt.Printf("距离: %.2f 米\n", distance)
}
```

---

### zColor - 终端颜色输出

提供终端彩色文本输出。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zColor"
)

func main() {
    fmt.Println(zColor.Green("成功信息"))
    fmt.Println(zColor.Red("错误信息"))
    fmt.Println(zColor.Yellow("警告信息"))
    fmt.Println(zColor.Cyan("提示信息"))
    fmt.Println(zColor.Blue("调试信息"))
    fmt.Println(zColor.Purple("特殊信息"))
}
```

---

### zError - 带错误码的错误类型

提供带错误码的错误类型，便于错误分类和处理。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zError"
)

func main() {
    // 创建带错误码的错误
    err := zError.NewWithCode(404, "资源未找到")
    fmt.Println(err.GetCode())     // 404
    fmt.Println(err.GetMessage())  // 资源未找到
    fmt.Println(err.Error())       // 资源未找到

    // 格式化创建
    err2 := zError.Errorf("用户 %s 不存在", "Alice")
    fmt.Println(err2.Error())  // 用户 Alice 不存在
}
```

---

### zUtils - 通用工具

提供目录获取和错误恢复等通用工具。

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/pzqf/zUtil/zUtils"
)

func main() {
    // 获取当前工作目录
    dir, _ := zUtils.GetCurrentDirectory()
    fmt.Println(dir)

    // 错误恢复
    func() {
        defer func() {
            if err := zUtils.Recover(); err != nil {
                fmt.Println("恢复错误:", err)
            }
        }()
        panic("something went wrong")
    }()
}
```

## 依赖

| 依赖 | 版本 | 用途 |
|------|------|------|
| `gopkg.in/ini.v1` | v1.67.0 | INI 文件解析 |
| `gopkg.in/yaml.v3` | v3.0.1 | YAML 文件解析 |
| `github.com/stretchr/testify` | v1.11.1 | 测试断言 |

## 已修复的问题

| 版本 | 问题 | 修复方案 |
|------|------|----------|
| v0.0.2 | zMap.Map 计数竞态 | `atomic.Value` → `atomic.Int64`，使用 `Add()` 替代 read-modify-write |
| v0.0.2 | zTime.New 忽略参数 | `time.Now()` → 使用传入参数 `t` |
| v0.0.2 | zStr.Distance 不支持中文 | 按字节索引 → `[]rune` 按字符索引 |
| v0.0.2 | zStr.Random 使用已废弃 rand.Seed | `math/rand` + `rand.Seed` → `math/rand/v2` |
| v0.0.2 | zList.Range 未加锁 | 添加 `locker.Lock()/Unlock()` |
| v0.0.2 | zKeyWordFilter DefaultFilter 未初始化保护 | 添加 `ensureDefaultFilter()` 懒初始化 |
| v0.0.2 | zQueue.Queue.Length() O(n) | 添加 `count` 字段，Enqueue/Dequeue/Empty 维护 |

## 安装

```bash
go get github.com/pzqf/zUtil
```

## 许可证

MIT License

---

*最后更新: 2026-04-14*
