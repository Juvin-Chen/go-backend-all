// 用 go 的标准库 net/http 写第一个Go Web 服务器

package main

import (
	"fmt"
	"net/http"
)

// 1.处理器函数（Handler）
// w：用来给客户端写回信的笔 (Response)
// r：客户端寄过来的”信件包“ (Request)
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 从Request中获取用户访问的路径
	path := r.URL.Path
	fmt.Printf("收到一个请求，访问路径是：%s\n", path)

	// 使用 fmt.Fprintf 把字符串写入 ResponseWriter，发送给浏览器
	fmt.Fprintf(w, "Hello,Go Web! 你的访问路径是：%s", path)
}

func main() {
	// 2.注册路由
	// 当用户访问"/"开头的路径时，交给hellohandler去处理
	http.HandleFunc("/", helloHandler)

	// 3.启动服务器，监听 8080 端口
	fmt.Println("🚀 服务器已启动，请在浏览器访问 http://localhost:8080")

	// ListenAndServe 会一直阻塞运行，除非发生致命错误
	err := http.ListenAndServe(":8800", nil)
	if err != nil {
		fmt.Println("服务器启动失败：%v\n", err)
	}
}
