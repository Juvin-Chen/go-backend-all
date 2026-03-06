# English to Chinese Reference
这里是中英文对照参考。

| 英文单词                           | 中文释义         | 在代码 / Go 标准库中的具体含义                               | 用法场景（结合你的代码）                                     |
| :--------------------------------- | :--------------- | :----------------------------------------------------------- | :----------------------------------------------------------- |
| `parse`（动词）                    | 解析、拆分       | 把字符串形式的 URL 拆分成 Go 能识别的结构化对象              | `url.Parse(rawUrl)`：解析原始 URL 字符串，返回`*url.URL`对象 |
| `panic`（动词 / 函数）             | 恐慌、终止程序   | Go 内置函数，用于处理**无法恢复的严重错误**，调用后立即终止程序并打印错误 + 调用栈 | `panic(err)`：URL 解析失败时，终止程序（因为后续逻辑无意义） |
| `scheme`（名词）                   | 协议、方案       | `url.URL`结构体字段，对应 URL 的协议部分（http/https/ftp 等） | `parseUrl.Scheme`：获取 URL 的协议（比如代码中输出`https`）  |
| `host`（名词）                     | 主机、主机名     | `url.URL`结构体字段，对应 URL 的域名 + 端口（如果有）        | `parseUrl.Host`：获取 URL 的主机名（比如代码中输出`www.baidu.com`） |
| `path`（名词）                     | 路径             | `url.URL`结构体字段，对应 URL 中域名后、参数前的部分         | `parseUrl.Path`：获取 URL 的路径（比如代码中输出`/s`）       |
| `query`（名词 / 动词）             | 查询、参数       | 1. 名词：URL 中`?`后的查询参数；2. 动词：解析查询参数        | `parseUrl.Query()`：解析 URL 的查询参数，返回`url.Values`（键值对）；`queryParams.Get("wd")`：获取指定查询参数的值 |
| `err`（名词，error 的缩写）        | 错误             | Go 中表示错误的变量名（约定俗成），存储函数返回的错误信息    | `parseUrl, err := url.Parse(rawUrl)`：接收 URL 解析的错误信息 |
| `raw`（形容词）                    | 原始的、未处理的 | 修饰 URL，表示未解析的字符串形式                             | `rawUrl`：变量名，指原始的 URL 字符串                        |
| `params`（名词，parameter 的复数） | 参数             | 查询参数的统称                                               | `queryParams`：变量名，指解析后的查询参数集合                |
| `get`（动词）                      | 获取             | `url.Values`的方法，获取指定键的参数值                       | `queryParams.Get("wd")`：获取`wd`参数的值                    |

------

### 额外补充：`net/url`包中高频单词（写 Gin 博客会用到）


| 英文单词   | 中文释义 | 用法                                                         |
| :--------- | :------- | :----------------------------------------------------------- |
| `port`     | 端口     | `parseUrl.Port()`：获取主机名中的端口（比如`www.baidu.com:8080`的`8080`） |
| `fragment` | 锚点     | `parseUrl.Fragment`：获取 URL 中`#`后的锚点（比如`https://xxx#top`的`top`） |
| `encode`   | 编码     | `url.Encode("go语言")`：把中文转成 URL 编码（`go%E8%AF%AD%E8%A8%80`） |
| `decode`   | 解码     | `url.QueryUnescape("go%E8%AF%AD%E8%A8%80")`：把 URL 编码转回中文 |
| `values`   | 值集合   | `url.Values`：存储查询参数的键值对类型（本质是`map[string][]string`） |