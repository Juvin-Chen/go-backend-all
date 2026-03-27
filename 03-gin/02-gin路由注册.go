/*
Gin 路由详解
*/
/*
1. 回顾：Gin 的路由注册
Gin 提供了丰富的 HTTP 方法对应的方法：
	r.GET()
	r.POST()
	r.PUT()
	r.DELETE()
	r.PATCH()
	r.HEAD()
	r.OPTIONS()`
也可以使用 r.Any() 匹配所有方法，或者 r.Handle() 指定方法。
r.GET("/users", listUsers)      // 仅 GET
r.POST("/users", createUser)     // 仅 POST
r.Any("/ping", ping)             // 所有方法

2. 路径参数
路径参数是 URL 的一部分，例如 /user/:id，其中 :id 是动态部分。
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")          // 获取路径参数
		c.String(http.StatusOK, "User ID: %s", id)
	})
测试：
	访问 /user/123 → 输出 User ID: 123
	访问 /user/abc → 输出 User ID: abc

多个参数：

下面的路径"/user/:id/post/:postID"中的post主要会起到一个分隔作用。

在路由 /user/:id/post/:postID 中：
	:id、:postID 是动态参数，值会变化；
	中间的 post 是固定路径单词，用于分隔两个参数，让 URL 语义清晰（表示 “某个用户的某篇帖子”），访问时必须原样写在地址里。

URL 必须严格遵循/user/[用户ID]/post/[帖子ID]格式，才能匹配成功：
✅ 合法地址：/user/100/post/520、/user/zhangsan/post/123
❌ 非法地址：/user/100/520（缺少固定词 post）、/user/100/article/520（post 写成 article）

r.GET("/user/:id/post/:postID", func(c *gin.Context) {
    userID := c.Param("id")
    postID := c.Param("postID")
    c.String(http.StatusOK, "User %s, Post %s", userID, postID)
})

3. 查询参数
查询参数是 URL 中 ? 之后的部分，如 /search?q=gin&page=1。
Gin 提供了 c.Query() 和 c.DefaultQuery() 方法。
	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")               // 获取 q 参数，若不存在返回空字符串
		page := c.DefaultQuery("page", "1") // 若不存在则返回默认值 "1"
		c.String(http.StatusOK, "Search: %s, Page: %s", q, page)
	})
也可以获取所有查询参数（返回 url.Values）：
	params := c.Request.URL.Query() // 是一个map[string][]string

4. 表单参数
处理 POST 表单时，使用 c.PostForm() 或 c.DefaultPostForm()。
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.DefaultPostForm("password", "default")
		c.String(http.StatusOK, "Username: %s, Password: %s", username, password)
	})
c.PostForm() 会自动解析请求体（表单格式），并返回对应字段的值。
注意：它默认只解析 application/x-www-form-urlencoded 格式；对于 JSON，需要绑定。
*/

package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 5. 参数绑定与验证（简介）
// Gin 提供了 c.Bind() 系列方法，可以将请求参数（查询参数、表单、JSON 等）自动绑定到结构体，并支持验证。
type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func testCbind() {
	r := gin.Default()
	r.POST("/login", func(ctx *gin.Context) {
		var req LoginRequest
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 处理登录
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}

/*
binding:"required" 表示字段必须存在且非空。
c.ShouldBind 会根据请求的 Content-Type 自动选择绑定方式（form, json, xml 等）。
我们会在后续课程专门讲解绑定与验证。
*/

// 6. 路由分组
// 当 API 有相同前缀时，可以用路由分组来组织代码，方便统一应用中间件。
func listUsers(c *gin.Context) {}

func testRouterGroup() {
	// 创建一个 /api 分组
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/users", listUsers)
		//api.GET("/users/:id", getUser)
		//api.POST("/users", createUser)
	}

	// 带注释是因为这些函数都没有定义，只是起到一个占位示例的作用

	// 带中间件的分组
	// authorized := r.Group("/admin")
	// authorized.Use(authMiddleware) // 对该分组内的路由使用认证中间件
	{
		//authorized.GET("/dashboard", dashboard)
		//authorized.POST("/settings", updateSettings)
	}
	// 分组可以嵌套，中间件会按顺序执行。
}

// 综合示例
// 我们编写一个小 API，包含用户列表、获取单个用户、创建用户，并展示路径参数、查询参数、表单/JSON 绑定。
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = []User{
	{ID: 1, Name: "zhangsan"},
	{ID: 2, Name: "lisi"},
}

func test2_example() {
	r := gin.Default()

	// 路由分组 / api
	api := r.Group("/api")
	{
		// GET /api/users 查询用户列表，支持 ?limit=N 分页
		api.GET("/users", func(ctx *gin.Context) {
			// 获取分页参数，默认 limit=10
			limitStr := ctx.DefaultQuery("limit", "10")
			limit, _ := strconv.Atoi(limitStr) // 字符串转数字
			// 真正的分页逻辑：截取前 N 条数据
			if limit > len(users) {
				limit = len(users)
			}
			pagedUsers := users[:limit]
			ctx.JSON(http.StatusOK, pagedUsers)
			// 简化处理：直接返回所有用户
			// ctx.JSON(http.StatusOK, users)
		})

		// GET /api/users/:id 获取单个用户
		api.GET("/user/:id", func(ctx *gin.Context) {
			idStr := ctx.Param("id")
			id, _ := strconv.Atoi(idStr)
			// 查找用户（简化）
			for _, u := range users {
				if u.ID == id {
					// gin.JSON 支持直接返回「结构体 / 切片」，不一定要用 gin.H！
					// 定义了结构体，并且写了 json 标签：
					ctx.JSON(http.StatusOK, u)
					return
				}
			}
			ctx.JSON(http.StatusNotFound, gin.H{"errors": "user not found"})
		})

		// POST /api/users 创建用户，接收 JSON
		api.POST("/users", func(ctx *gin.Context) {
			var newUser User
			if err := ctx.ShouldBindJSON(&newUser); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			newUser.ID = len(users) + 1
			users = append(users, newUser)
			ctx.JSON(http.StatusCreated, newUser)
		})
	}
	r.Run()
}
