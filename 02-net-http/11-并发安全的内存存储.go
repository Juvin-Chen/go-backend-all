/*
lesson 11
并发安全的内存存储
*/

/*
1. 为什么需要并发安全？
在 Web 应用中，每个 HTTP 请求都会在一个独立的 goroutine 中处理。
如果多个请求同时读写同一个内存数据结构（比如 slice、map），就可能发生 数据竞争（data race），导致程序崩溃或数据错乱。
例如，一个简单的留言板，用 []Message 存储留言。两个请求同时添加留言，如果不对这个 slice 进行保护，可能会造成数据丢失或 panic。
Go 提供了 sync.Mutex 和 sync.RWMutex 来帮助我们在并发环境下安全地访问共享资源。

2. 使用 sync.Mutex 实现互斥锁
sync.Mutex 提供两个方法：
	Lock()：加锁，如果锁已被占用则阻塞等待
	Unlock()：解锁
我们可以在读写共享数据前后加锁，确保同一时间只有一个 goroutine 能够操作数据。
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// 简单的线程安全计数器
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func testCounter() {
	c := &Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	fmt.Println(c.Value()) // 输出 1000
}

/*
解释：
	mu sync.Mutex 是互斥锁。
	在 Inc 和 Value 方法中，先 Lock()，确保同一时间只有一个 goroutine 执行临界区代码。
	使用 defer c.mu.Unlock() 确保锁一定会被释放，即使发生 panic 也会释放。
*/

// 3.使用 sync.RWMutex 提高读性能
/*
sync.RWMutex 是读写锁，它允许多个读操作并发执行，但写操作是独占的。
( 读 和 写 用的是同一把锁)
方法：
	RLock() / RUnlock()：读锁
	Lock() / Unlock()：写锁

在“读多写少”的场景下，使用读写锁可以提高并发性能。
*/

type SafeSlice struct {
	mu   sync.RWMutex
	data []string
}

func (s *SafeSlice) Add(item string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, item)
}

func (s *SafeSlice) GetAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// 返回副本，避免外部修改
	result := make([]string, len(s.data))
	copy(result, s.data)
	return result
}

/*
注意：在 GetAll 中，我们返回了数据的副本，而不是直接返回 s.data。这是因为返回后外部可能修改这个切片，而我们不希望它影响内部数据。如果返回内部切片引用，读锁释放后外部可能并发修改，造成数据竞争。返回副本是安全的做法。
*/

// 4. 在 Web 应用中使用并发安全的内存存储
// 我们之前留言板项目中的 MemoryStore 就使用了 sync.RWMutex。
// 回顾一下核心代码：

type Message struct {
	ID        int
	Nickname  string
	Content   string
	CreatedAt time.Time
}

type MemoryStore struct {
	mu       sync.RWMutex
	messages []Message
	nextID   int
}

func (s *MemoryStore) Add(nickname, content string) Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	msg := Message{
		ID:        s.nextID,
		Nickname:  nickname,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.messages = append(s.messages, msg)
	s.nextID++
	return msg
}

func (s *MemoryStore) GetAll() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Message, len(s.messages))
	copy(result, s.messages)
	// 倒序（最新的在前）
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

/*
Add 使用写锁（Lock），因为要修改 messages 和 nextID。
GetAll 使用读锁（RLock），因为只读取数据。并且返回了数据的副本，避免外部修改影响内部存储。
这样，多个请求可以同时读取留言列表（并发读），而写入时会互斥，保证数据一致性
*/

// 5. 更复杂的数据结构：并发安全的 map
/*
Go 的 map 不是并发安全的，多个 goroutine 同时读写会 panic。
我们可以用 sync.RWMutex 保护 map，或者使用 sync.Map（适合特定场景）。
*/

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]string),
	}
}

func (sm *SafeMap) Set(key, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

func (sm *SafeMap) Get(key string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, ok := sm.m[key]
	return val, ok
}

/*
sync.Map 是官方提供的并发安全 map，但它的 API 是通用的，性能在某些场景下可能不如带锁的普通 map。一般我们先用 sync.RWMutex 保护 map，除非在大量读、少量写且 key 不重复的场景下才考虑 sync.Map。
*/

/*
常见问题
Q：为什么不用全局锁保护整个存储，而是用更细的锁？
A：全局锁会降低并发度，比如读操作也需要等待写锁。使用读写锁可以允许多个读并发。

Q：什么时候应该用 sync.Mutex，什么时候用 sync.RWMutex？
A：如果读操作远多于写操作，用 RWMutex 可以提升性能；如果读写操作比例接近或写操作多，用普通 Mutex 可能更简单。

Q：返回数据副本会不会有性能问题？
A：对于小数据量没问题。如果数据量很大，可以返回只读视图（比如封装成只读接口），但需要注意避免外部修改。在 Web 应用中，通常数据量不会太大，复制是安全的。
*/
