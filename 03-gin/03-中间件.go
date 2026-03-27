/*
Gin 中间件（Middleware）

1.Gin 中间件的本质
2.如何注册和使用中间件（全局、路由组、单个路由）
3.自定义中间件
4.中间件中的 Context 操作
5.常用内置中间件
6.执行顺序
7.实战：编写日志、认证、恢复中间件
*/

/*
1. Gin 中间件的本质
在 Gin 中，中间件本质上是一个函数，它接收 *gin.Context 并执行一些操作，然后可以选择继续执行后续的处理器（通过 c.Next()）或者中止（通过 c.Abort()）。
中间件的签名通常是：
	func MyMiddleware(c *gin.Context) {
		// 前置操作
		c.Next() // 调用后续处理器（如果有）
		// 后置操作
	}
Gin 把多个处理器（包括中间件和最终的业务处理器）组织成一个 处理器链，按照注册顺序依次执行。
这与 net/http 中通过 http.Handler 嵌套实现中间件本质相同，但 Gin 将其封装为更直观的 Use 方法和 c.Next()。

2. 注册中间件
Gin 提供了 r.Use() 方法，可以注册一个或多个中间件，它们将应用于该路由引擎或路由组。

2.1 全局中间件
r := gin.Default() // 默认已包含 Logger 和 Recovery 中间件
// 或者用 r := gin.New() 创建不带任何中间件的引擎

// 注册自定义全局中间件
r.Use(MyMiddleware)
r.GET("/test", func(c *gin.Context) {
    c.String(200, "OK")
})

2.2 路由组中间件
api := r.Group("/api")
api.Use(AuthMiddleware) // 该组内所有路由都应用 AuthMiddleware
{
    api.GET("/user", userHandler)
    api.GET("/profile", profileHandler)
}

2.3 单个路由中间件
r.GET("/admin", AuthMiddleware, adminHandler)
可以为一个路由注册多个中间件，按顺序执行。
*/

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 3. 编写自定义中间件

// 3.1 日志中间件（记录请求耗时）
func LoggerMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next() // 执行后续处理器
	latency := time.Since(start)
	log.Printf("Request %s %s took %v", c.Request.Method, c.Request.URL.Path, latency)
}

func testM() {
	r := gin.New()
	r.Use(LoggerMiddleware)
}

// 3.2 认证中间件
func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || token != "Bearer secret-token" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort() // 终止后续处理器
		return
	}
	c.Next()
}

// 3.3 恢复中间件（Gin 已内置，但可以自定义）
func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
		}
	}()
	c.Next()
}

// 4. 中间件中的上下文传递
// 在 Gin 中，可以使用 c.Set(key, value) 和 c.Get(key) 在中间件和处理器之间传递数据。
func AuthMiddleware_2(c *gin.Context) {
	// 假设验证通过，获取用户信息
	user := User{ID: 1, Name: "Alice"}
	c.Set("user", user)
	c.Next()
}

func UserHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
