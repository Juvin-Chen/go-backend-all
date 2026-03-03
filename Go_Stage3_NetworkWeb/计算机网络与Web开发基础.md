# 计算机网络与 Web 开发基础（Go 语言版）

## 第一章：计算机网络与通信协议基础

### 1.1 什么是计算机网络与通信协议？

**计算机网络**是指将地理位置不同的具有独立功能的计算机及其外部设备，通过通信线路连接起来，在网络操作系统、管理软件及**网络通信协议**的协调下，实现**资源共享和信息传递**的系统。

- **组成**：硬件（服务器、路由器、光纤等）、软件（操作系统、管理软件、通信协议）。

**网络通信协议**是计算机之间通信时事先达成的一种“约定”。因为网络设备来自不同厂商、有着不同的 OS，必须遵循相同的规则（如速率、传输代码、错误控制等）才能通信。

### 1.2 OSI 七层模型与 TCP/IP 模型

为了让网络标准不再复杂，国际标准化组织（ISO）提出了著名的 **OSI 七层参考模型**，而实际互联网中最常用的是 **TCP/IP 协议族（四层/五层模型）**。

- **应用层**（OSI 第 7 层 / TCP/IP 第 4 层）：为应用程序提供服务。相关协议：HTTP, FTP, DNS, SNMP。
- **表示层、会话层**（OSI 第 6, 5 层）：数据格式转化、加密，建立和维护会话。在 TCP/IP 中被合并入应用层。
- **传输层**（OSI 第 4 层 / TCP/IP 第 3 层）：建立端到端的连接。提供进程间的逻辑通信。相关协议：TCP, UDP。
- **网络层**（OSI 第 3 层 / TCP/IP 第 2 层）：IP 选址及路由选择，使数据包在网络中路由。相关协议：IP, ICMP, ARP。
- **数据链路层、物理层**（网络接口层）：提供介质访问、转化比特流等底层物理通信。

### 1.3 数据封装与解封流程

数据在网络中传输，就像打包快递：

- **数据封装（发送方，从上到下）**：应用层数据 -> 传输层加上 TCP 头变成**段（Segment）** -> 网络层加上 IP 头变成**包（Packet）** -> 数据链路层加上 MAC 头尾变成**帧（Frame）** -> 物理层转化为比特流发送。
- **数据解封（接收方，从下到上）**：物理层接收比特流 -> 数据链路层剥离 MAC 头 -> 网络层剥离 IP 头 -> 传输层剥离 TCP 头 -> 交给应用层。
- *核心概念：对等层通信，接收方的每一层只处理发送方同等层封装的数据。*

### 1.4 广域网 (WAN) 与 局域网 (LAN)

- **局域网（LAN）**：“小范围的私有网络”。覆盖单一空间（如家庭、公司），自己管理，速度快、延迟极低（<10ms）、极其安全。设备以交换机、家用路由器为主。
- **广域网（WAN）**：“大范围的公共网络”。跨城市/国家，由运营商维护，速度受限、延迟较高、需加密保障安全。核心用于实现不同 LAN 之间的互联（如你连家里 WiFi 逛淘宝，数据通过运营商 WAN 传到淘宝服务器）。

------

## 第二章：网络核心元素：IP、域名、端口与 URL

### 2.1 IP 地址

IP 地址用来标识网络中唯一的通信实体（如计算机、路由器）。主流为 IPv4（32位整数，已分配完毕）和 IPv6（128位整数，海量地址）。

- **公有 IP**：由互联网信息中心分配，直接访问互联网。
- **私有 IP**：专门为组织机构内部（局域网）使用。如常用的 `192.168.x.x`、`10.x.x.x`、`172.16.x.x`。`127.0.0.1` 为本机回环地址。

> **Q：为什么百度有 IP 地址？** 百度是一家公司，提供服务需要成千上万台服务器。百度的 IP 就像公司的“总机号码”，用户拨打总机，系统会自动分配给内部空闲的服务器（分机）来为你服务。

### 2.2 域名与 DNS 解析

**域名**是 IP 的别名（如 `www.baidu.com`），便于人类记忆。**DNS（域名系统）**就像通讯录，负责将域名“翻译”为 IP 地址。

- **访问全过程**：浏览器缓存 -> 系统缓存 -> 路由器缓存 -> ISP（运营商）DNS 服务器 -> 递归查询（根域名 -> 顶级域 -> 权威域） -> 拿到 IP 建立连接。

### 2.3 端口 (Port)

如果 IP 是门牌号，端口就是**房间号**。端口号用于识别计算机中进行通信的具体**应用程序**，范围 0~65535。

- **公认端口 (0-1023)**：紧密绑定特定服务，如 HTTP 默认 80，HTTPS 默认 443，FTP 默认 21。
- **注册端口 (1024-65535)**：松散绑定，多用于动态分配。
- *注意：不能两个应用同时监听同一个端口（会冲突），但一个应用可以监听多个端口。*

### 2.4 URI、URL 与 URN

- **URI (统一资源标识符)**：是一个纯粹的语法结构，用于标识互联网唯一资源的字符串。
- **URL (统一资源定位符)**：URI 的子集，相当于资源的身份证号+家庭住址。组成：`协议://服务器(域名/IP):端口/请求路径?传递的数据(key=value)`。
- **URN (统一资源名称)**：URI 的子集，只标识名字不提供访问方式。

------

## 第三章：传输层核心：TCP 与 UDP 协议详解

TCP 和 UDP 是位于传输层的两种截然不同的协议，它们都是通过 **Socket（套接字）** 供应用层调用的。

### 3.1 TCP 协议 (面向连接、可靠)

类似于“打电话”，必须建立专属虚拟连接。具备流量控制和拥塞控制。适用于对可靠性要求高的场景（如文件传输、网页浏览）。

**核心机制：三次握手与四次挥手**

- 建立连接（三次握手）：
  1. 客户端发送 SYN（同步包）到服务器，进入 SYN-SEND 状态。
  2. 服务器收到，返回 SYN + ACK（确认包），进入 SYN-RECV 状态。
  3. 客户端收到，返回 ACK 给服务器，双方进入 ESTABLISHED（已连接）状态。
- 断开连接（四次挥手）：
  1. 客户端发起 FIN 包，表示数据发完了，进入 FIN_WAIT_1。
  2. 服务器收到 FIN，返回 ACK 确认，进入 CLOSE_WAIT 状态（此时服务器可能还有数据没发完，客户端仍需接收）。
  3. 服务器数据发完后，发送 FIN 包给客户端，进入 LAST_ACK。
  4. 客户端收到，返回 ACK，进入 TIME_WAIT（等一小会儿确保断开），服务器收到后彻底 CLOSED。

### 3.2 UDP 协议 (无连接、不可靠)

类似于“发短信”，不需要建立连接，只管发，不保证对方收到。优点是速度快、无连接开销。适用于对实时性要求高的场景（视频会议、直播、游戏）。

------

## 第四章：应用层核心：HTTP 协议全解

**HTTP (超文本传输协议)** 是万维网通信的基础，基于 TCP/IP 运行。

### 4.1 HTTP 协议特性

1. **简单快速、灵活**：只需传送请求方法和路径，可传输任意类型（通过 `Content-Type` 标记）。
2. **无连接 (短连接)**：每次请求响应结束后断开。但 HTTP 1.1 后默认开启 `Connection: keep-alive` 实现长连接。
3. **单向性、无状态**：永远是客户端主动请求；服务器不记忆上下文（需要 Cookie 和 Session 技术维持状态）。

### 4.2 HTTP 发展版本

- **HTTP 1.0**：短连接，每次请求都要重新建立 TCP 连接，无状态。
- **HTTP 1.1**：广泛使用，支持**长连接**（一次连接多次请求），但请求需要排队发送阻塞。
- HTTP 2.0 ：
  - **多路复用**：同一个连接并发处理多个请求（引入帧和流的概念），解决队头阻塞。
  - **二进制传输**：取代文本传输，解析更高效。
  - **首部压缩**：压缩 Header 减少数据量。
  - **服务端推送**：服务器可主动向客户端推资源。

### 4.3 HTTP 请求 (Request) 与 GET/POST 的区别

请求包含三部分：**请求行、请求头、请求体**（请求头和请求体之间有空行）。

- **请求行**：`方法 URL 协议版本` (例如 `GET /login.html HTTP/1.1`)。
- **请求头**：`Key: Value`。常见有 `Host`、`User-Agent`、`Accept`、`Cookie` 等。
- **请求体**：POST 方法传给服务器的数据。

**GET 与 POST 的核心区别（面试必考）**：

1. GET 参数在 URL 中（不安全，有长度限制）；POST 参数在请求体中（相对安全，无限制，支持字节流）。
2. GET 请求可被浏览器缓存/收藏/保存历史记录，POST 不行。
3. GET 刷新/回退无害，POST 回退会重新提交表单。

### 4.4 HTTP 响应 (Response) 与 MIME 类型

包含三部分：**响应行、响应头、响应体**。

- **响应行**：`协议版本 状态码 描述` (如 `HTTP/1.1 200 OK`)。
- 常见状态码：
  - `200`：成功。
  - `301/302`：永久/临时重定向。
  - `400`：客户端请求语法错误。
  - `401/403`：未授权 / 拒绝访问。
  - `404`：资源未找到。
  - `500`：服务端异常。
- **MIME 类型**：在响应头 `Content-Type` 中指定，告诉浏览器返回的数据是什么类型（如 `text/html`, `application/json`, `image/jpeg`），浏览器据此决定如何解析。

------

## 第五章：Go 语言网络编程实战 (对标 Java)

在 Java 中，网络编程依赖 `java.net` 包下的 `InetAddress`、`URL`、`Socket` 等类。**在 Go 语言中，我们主要使用 `net` 和 `net/url` 标准库。**

> **关于 Socket 的深层理解：** Socket 是应用层和传输层之间的桥梁。通信必须有两端：`Socket(IP, Port, 协议)` 组成的三元组代表一个端点。 在服务端，`ServerSocket` 就像公司的**总机接线员**（只负责监听端口），当有连接进来时（`accept()`），分配一个新的 `Socket`（相当于**员工分机**）去和客户端专门通信。

### 5.1 解析 URL (对应 Java 的 `URL` 类)

```
package main

import (
	"fmt"
	"net/url"
)

func main() {
	// 对应 Java 中的 new URL(...)
	rawUrl := "https://www.itbaizhan.com/search.html?kw=java"
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}

	fmt.Println("协议 (Protocol):", parsedUrl.Scheme)       // https
	fmt.Println("主机名 (Host):", parsedUrl.Host)           // www.itbaizhan.com
	fmt.Println("路径 (Path):", parsedUrl.Path)           // /search.html
	fmt.Println("参数部分 (Query):", parsedUrl.RawQuery)     // kw=java

    // 获取具体的参数值
	queryParams := parsedUrl.Query()
	fmt.Println("kw 的值:", queryParams.Get("kw"))        // java
}
```

### 5.2 TCP 编程：服务端与客户端 (对应 Java `ServerSocket` 与 `Socket`)

**Go 服务端实现 (Server)** 在 Go 中，我们不需要像 Java 那样通过手动分配多线程（`extends Thread`）来处理多客户端并发。Go 原生提供轻量级的 `goroutine` 进行高并发处理，极其简洁。

```
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// 1. 对应 Java: ServerSocket serverSocket = new ServerSocket(8888);
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("监听失败:", err)
		return
	}
	defer listener.Close()
	fmt.Println("服务端启动，等待监听 8888 端口...")

	for {
		// 2. 对应 Java: Socket socket = serverSocket.accept(); (阻塞等待)
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接收连接失败:", err)
			continue
		}
		fmt.Println("有客户端连接了:", conn.RemoteAddr())

		// 3. 启动 Goroutine 处理该客户端的读写（替代 Java 的多线程机制）
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// 4. 对应 Java: br.readLine() 获取客户端消息
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("客户端断开连接:", conn.RemoteAddr())
			return
		}
		fmt.Print("客户端说: ", msg)

		// 回复客户端 (对应 Java: pw.println(str); pw.flush();)
		reply := fmt.Sprintf("服务器已收到: %s", msg)
		conn.Write([]byte(reply))
	}
}
```

**Go 客户端实现 (Client)**

```
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 1. 对应 Java: Socket socket = new Socket("127.0.0.1", 8888);
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接服务端失败:", err)
		return
	}
	defer conn.Close()
	fmt.Println("客户端启动，连接服务端成功！")

	inputReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Print("请输入发送内容 (exit退出): ")
		// 读取键盘输入
		msg, _ := inputReader.ReadString('\n')

		// 对应 Java: pw.println(msg); pw.flush();
		_, err = conn.Write([]byte(msg))
		if msg == "exit\n" || msg == "exit\r\n" {
			break
		}

		// 等待接收服务端回复
		reply, _ := serverReader.ReadString('\n')
		fmt.Print("服务端返回: ", reply)
	}
}
```

### 5.3 UDP 编程：基本数据类型与对象的传输 (对标 Java `DatagramSocket`)

UDP 不需要提前建连，在 Go 中使用 `net.ListenUDP` 和 `net.DialUDP`。

**传输自定义对象 (使用 JSON 序列化代替 Java 的 Serializable)** 在 Java 中，传对象必须要实现 `Serializable` 接口并使用 `ObjectOutputStream`。 在现代 Web 尤其是 Go 语言开发中，**传递结构体（对象）最标准、跨语言的做法是将其序列化为 JSON 字节数组。**

**UDP 客户端 (发送方)**

```
package main

import (
	"encoding/json"
	"fmt"
	"net"
)

// 对应 Java 中的 Person 类，注意 Go 中的字段需要大写才能被 json 包导出
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// UDP 服务端地址
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	// 本地客户端分配随机 UDP 端口
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("建立UDP连接失败:", err)
		return
	}
	defer conn.Close()

	// 1. 实例化对象并进行 JSON 序列化
	p := Person{Name: "Oldlu", Age: 18}
	// json.Marshal 替代了 Java 的 ObjectOutputStream.writeObject()
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("序列化失败:", err)
		return
	}

	// 2. 发送 UDP 数据报文 (对应 Java的 datagramSocket.send(dp))
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("发送失败:", err)
	} else {
		fmt.Println("对象发送成功!")
	}
}
```

**UDP 服务端 (接收方)**

```go
package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// 1. 监听本地 9999 端口
	addr, _ := net.ResolveUDPAddr("udp", ":9999")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("监听失败:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP服务端启动，等待数据...")

	// 2. 对应 Java 的 byte[] b = new byte; DatagramPacket dp = new DatagramPacket...
	buf := make([]byte, 1024)

	for {
		// 阻塞接收数据包
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("接收数据错误:", err)
			continue
		}

		// 3. JSON 反序列化 (替代 Java 的 ObjectInputStream.readObject())
		var p Person
		err = json.Unmarshal(buf[:n], &p)
		if err != nil {
			fmt.Println("反序列化失败:", err)
			continue
		}

		fmt.Printf("收到来自 %v 的对象数据: Name=%s, Age=%d\n", clientAddr, p.Name, p.Age)
	}
}
```

