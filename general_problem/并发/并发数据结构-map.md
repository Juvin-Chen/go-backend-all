# 并发安全 Map 与 sync.RWMutex 详解

> 更多了解也可以参考: 02-net-http / 11-并发安全的内存存储

适用场景：多 goroutine 同时读写 map，避免程序崩溃

## 一、核心结论
1. Go 内置的 `map` **不是线程安全**的，也**不是原子操作**。
2. 多个 goroutine 同时对 map 进行**读+写/写+写**，程序会直接 panic。
3. 必须通过 `sync.Mutex` 或 `sync.RWMutex` 加锁保护，才能安全并发使用。

## 二、为什么 map 不安全？
- map 底层结构复杂，赋值、查找、删除都需要多步操作（计算哈希、定位桶、扩容、迁移数据）。
- 执行过程中会被其他 goroutine 打断，导致数据结构混乱。
- 典型报错：
  ```
  fatal error: concurrent map read and map write
  ```

## 三、什么是原子操作？
原子操作：**一步完成、不可被中断**的操作。
例如简单数值赋值、加减等。
map 操作由多步组成，**不属于原子操作**，因此无法靠自身保证并发安全。

## 四、sync.RWMutex 读写锁
比普通互斥锁 `sync.Mutex` 更高效，适合**读多写少**场景：
- `RLock()` / `RUnlock()`：读锁，**多人可同时读**，但不能写。
- `Lock()` / `Unlock()`：写锁，**独占**，读写都阻塞。

规则：
- 读与读可以并发
- 读与写互斥
- 写与写互斥

## 五、SafeMap 完整实现与解析
```go
type SafeMap struct {
	mu sync.RWMutex // 读写锁：保护内部 map
	m  map[string]string // 实际存储数据的 map
}

// NewSafeMap 创建一个安全的 map 实例
func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]string),
	}
}

// Set 写操作：加写锁
func (sm *SafeMap) Set(key, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

// Get 读操作：加读锁
func (sm *SafeMap) Get(key string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, ok := sm.m[key]
	return val, ok
}
```

### 关键点
- `defer Unlock`：确保函数退出一定会释放锁，避免死锁。
- 读用 `RLock`，写用 `Lock`，提升并发读性能。

## 六、总结
1. map 非线程安全，并发读写必加锁。
2. `RWMutex` 适合读多写少，比普通锁效率更高。
3. 封装 SafeMap 后，外部调用无需关心锁细节，直接安全使用。

