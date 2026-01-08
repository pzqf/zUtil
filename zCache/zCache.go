package zCache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Item 缓存项结构体
type Item struct {
	Key        string
	Value      interface{}
	Expiration int64  // 过期时间戳（纳秒）
	Element    *list.Element
}

// Cache 缓存接口
type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Clear() error
	Len() int
	Keys() []string
}

// LRUCache 基于LRU的缓存实现
type LRUCache struct {
	mu            sync.Mutex
	cache         map[string]*Item
	list          *list.List
	capacity      int
	expiration    time.Duration  // 默认过期时间
	cleanupTicker *time.Ticker   // 清理过期项的定时器
	stopChan      chan struct{}  // 停止清理的通道
}

// NewLRUCache 创建新的LRU缓存
func NewLRUCache(capacity int, expiration time.Duration) *LRUCache {
	if capacity <= 0 {
		capacity = 100  // 默认容量
	}

	cache := &LRUCache{
		cache:      make(map[string]*Item),
		list:       list.New(),
		capacity:   capacity,
		expiration: expiration,
		stopChan:   make(chan struct{}),
	}

	// 启动定期清理过期项的goroutine
	if expiration > 0 {
		cache.cleanupTicker = time.NewTicker(expiration / 2)
		go cache.cleanupLoop()
	}

	return cache
}

// Set 设置缓存项
func (c *LRUCache) Set(key string, value interface{}, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果指定了过期时间，则使用指定的，否则使用默认的
	if expiration <= 0 {
		expiration = c.expiration
	}

	// 计算过期时间戳
	var expirationTime int64
	if expiration > 0 {
		expirationTime = time.Now().Add(expiration).UnixNano()
	}

	// 检查key是否已存在
	if item, ok := c.cache[key]; ok {
		// 更新值和过期时间
		item.Value = value
		item.Expiration = expirationTime
		// 移到链表头部
		c.list.MoveToFront(item.Element)
		return nil
	}

	// 如果缓存已满，移除最久未使用的项
	if len(c.cache) >= c.capacity {
		c.removeOldest()
	}

	// 创建新的缓存项
	item := &Item{
		Key:        key,
		Value:      value,
		Expiration: expirationTime,
	}

	// 添加到链表头部和map中
	item.Element = c.list.PushFront(item)
	c.cache[key] = item

	return nil
}

// Get 获取缓存项
func (c *LRUCache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查key是否存在
	item, ok := c.cache[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}

	// 检查是否过期
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		// 已过期，移除并返回错误
		c.remove(item)
		return nil, fmt.Errorf("key %s expired", key)
	}

	// 移到链表头部
	c.list.MoveToFront(item.Element)

	return item.Value, nil
}

// Delete 删除缓存项
func (c *LRUCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查key是否存在
	item, ok := c.cache[key]
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}

	// 移除项
	c.remove(item)

	return nil
}

// Clear 清空缓存
func (c *LRUCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清空map和链表
	c.cache = make(map[string]*Item)
	c.list.Init()

	return nil
}

// Len 获取缓存项数量
func (c *LRUCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.cache)
}

// Keys 获取所有缓存键
func (c *LRUCache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]string, 0, len(c.cache))
	for key := range c.cache {
		keys = append(keys, key)
	}

	return keys
}

// Stop 停止缓存的清理功能
func (c *LRUCache) Stop() {
	c.mu.Lock()
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
		c.cleanupTicker = nil
	}
	c.mu.Unlock()

	close(c.stopChan)
}

// removeOldest 移除最久未使用的项
func (c *LRUCache) removeOldest() {
	element := c.list.Back()
	if element != nil {
		item := element.Value.(*Item)
		c.remove(item)
	}
}

// remove 移除指定的缓存项
func (c *LRUCache) remove(item *Item) {
	c.list.Remove(item.Element)
	delete(c.cache, item.Key)
}

// cleanupLoop 定期清理过期项的循环
func (c *LRUCache) cleanupLoop() {
	for {
		select {
		case <-c.stopChan:
			return
		default:
			// 检查cleanupTicker是否为nil
			c.mu.Lock()
			ticker := c.cleanupTicker
			c.mu.Unlock()
			
			if ticker != nil {
				select {
				case <-c.stopChan:
					return
				case <-ticker.C:
					c.cleanup()
				}
			} else {
				// 如果ticker为nil，休眠一段时间后再次检查
				time.Sleep(time.Second)
			}
		}
	}
}

// cleanup 清理过期项
func (c *LRUCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for _, item := range c.cache {
		if item.Expiration > 0 && now > item.Expiration {
			c.remove(item)
		}
	}
}

// SimpleCache 简单缓存实现（无LRU，仅支持过期时间）
type SimpleCache struct {
	mu       sync.Mutex
	cache    map[string]*Item
	expiration time.Duration  // 默认过期时间
}

// NewSimpleCache 创建新的简单缓存
func NewSimpleCache(expiration time.Duration) *SimpleCache {
	return &SimpleCache{
		cache:      make(map[string]*Item),
		expiration: expiration,
	}
}

// Set 设置缓存项
func (c *SimpleCache) Set(key string, value interface{}, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果指定了过期时间，则使用指定的，否则使用默认的
	if expiration <= 0 {
		expiration = c.expiration
	}

	// 计算过期时间戳
	var expirationTime int64
	if expiration > 0 {
		expirationTime = time.Now().Add(expiration).UnixNano()
	}

	// 创建或更新缓存项
	c.cache[key] = &Item{
		Key:        key,
		Value:      value,
		Expiration: expirationTime,
	}

	return nil
}

// Get 获取缓存项
func (c *SimpleCache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查key是否存在
	item, ok := c.cache[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}

	// 检查是否过期
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		// 已过期，移除并返回错误
		delete(c.cache, key)
		return nil, fmt.Errorf("key %s expired", key)
	}

	return item.Value, nil
}

// Delete 删除缓存项
func (c *SimpleCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查key是否存在
	_, ok := c.cache[key]
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}

	// 移除项
	delete(c.cache, key)

	return nil
}

// Clear 清空缓存
func (c *SimpleCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清空map
	c.cache = make(map[string]*Item)

	return nil
}

// Len 获取缓存项数量
func (c *SimpleCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清理过期项
	now := time.Now().UnixNano()
	for key, item := range c.cache {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.cache, key)
		}
	}

	return len(c.cache)
}

// Keys 获取所有缓存键
func (c *SimpleCache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清理过期项
	now := time.Now().UnixNano()
	for key, item := range c.cache {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.cache, key)
		}
	}

	keys := make([]string, 0, len(c.cache))
	for key := range c.cache {
		keys = append(keys, key)
	}

	return keys
}

// TTLCache 基于TTL的缓存实现（自动清理过期项）
type TTLCache struct {
	*SimpleCache
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

// NewTTLCache 创建新的TTL缓存
func NewTTLCache(expiration time.Duration) *TTLCache {
	if expiration <= 0 {
		expiration = time.Minute  // 默认过期时间
	}

	cache := &TTLCache{
		SimpleCache:   NewSimpleCache(expiration),
		cleanupTicker: time.NewTicker(expiration / 2),
		stopChan:      make(chan struct{}),
	}

	// 启动定期清理过期项的goroutine
	go cache.cleanupLoop()

	return cache
}

// cleanupLoop 定期清理过期项的循环
func (c *TTLCache) cleanupLoop() {
	for {
		select {
		case <-c.stopChan:
			return
		case <-c.cleanupTicker.C:
			c.cleanup()
		}
	}
}

// cleanup 清理过期项
func (c *TTLCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range c.cache {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.cache, key)
		}
	}
}

// Stop 停止缓存的清理功能
func (c *TTLCache) Stop() {
	c.cleanupTicker.Stop()
	close(c.stopChan)
}

// DefaultCache 默认缓存实例（使用LRUCache）
var defaultCache Cache
var once sync.Once

// GetDefaultCache 获取默认缓存实例
func GetDefaultCache() Cache {
	once.Do(func() {
		defaultCache = NewLRUCache(1000, time.Hour)
	})
	return defaultCache
}

// SetDefault 设置默认缓存项
func SetDefault(key string, value interface{}, expiration time.Duration) error {
	return GetDefaultCache().Set(key, value, expiration)
}

// GetDefault 获取默认缓存项
func GetDefault(key string) (interface{}, error) {
	return GetDefaultCache().Get(key)
}

// DeleteDefault 删除默认缓存项
func DeleteDefault(key string) error {
	return GetDefaultCache().Delete(key)
}