# Lesson 9 静态文件服务

## 一、核心概念
**静态文件**：不需要后端处理逻辑、直接给浏览器使用的文件
- CSS：页面样式
- JavaScript：页面交互
- 图片、字体、图标等

**静态文件服务**：**后端（Go）提供的功能**
作用：**让浏览器能通过网址，访问到服务器上的静态文件**

---

## 二、关键代码解释
```go
// 1. 绑定服务器本地文件夹 ./static
fs := http.FileServer(http.Dir("./static"))

// 2. 映射网址路径 /static/ 到文件服务
http.Handle("/static/", http.StripPrefix("/static/", fs))
```

### 1. http.Dir("./static")
- 告诉 Go：**静态文件都存在项目里的 `./static` 文件夹**
- 这是**服务器文件夹路径**，不是网址

### 2. http.FileServer
- 创建一个“文件处理器”
- 专门负责：找文件 → 读取文件 → 返回给浏览器

### 3. http.StripPrefix("/static/", fs)
**核心作用：路径翻译**
- 浏览器请求：`/static/style.css`
- 经过 StripPrefix 去掉 `/static/`
- 变成：`style.css`
- 最终去服务器找：`./static/style.css`

**如果不用 StripPrefix，会报 404**
因为会错误去找：`./static/static/style.css`

---


## 三、HTML 中 href 到底是什么？（核心困惑）
```html
<link rel="stylesheet" href="/static/style.css">
```

### 正确理解
1. **href 不是电脑文件夹路径！**
2. **href 是一个网址！**
3. 意思：
   > 浏览器，你去 `http://localhost:8080/static/style.css` 下载这个文件！

### 页面加载真实流程
1. 浏览器请求 `/` → 拿到 HTML（骨架）
2. 浏览器看到 `href` → **自动再发一次请求**
3. 请求 `/static/style.css`
4. Go 静态文件服务找到文件并返回
5. 浏览器用 CSS 美化页面

**一个页面 = 至少两次请求！**

---

## 四、我当时的困惑 + 正确答案
### 我的困惑
1. 静态文件服务到底是前端还是后端？
2. href 是什么？为什么要写这个路径？
3. 页面不是已经有自己的 path 了吗？
4. 后端这段代码到底有什么用？

### 正确答案
1. **静态文件（CSS/JS）属于前端**
   **静态文件服务（Go代码）属于后端**
2. **href = 让浏览器去下载样式文件的网址**
3. **HTML 是骨架，CSS 是皮肤**
   页面 path 只负责返回骨架，**href 负责请求皮肤**
4. **后端代码 = 翻译官**
   把网址路径 `/static/xxx`
   翻译成服务器文件路径 `./static/xxx`

---

## 五、一句话终极总结
**HTML 决定页面有什么，href 告诉浏览器去哪找样式，
Go 静态文件服务负责把文件找出来发给浏览器。
没有后端这段代码，页面就没有样式、没有图片！**

---

## 六、最简单记忆图
```
浏览器请求 /static/style.css  ← 网址路径
          ↓
Go 后端 StripPrefix 翻译路径
          ↓
去服务器 ./static/style.css 找文件 ← 文件路径
          ↓
返回给浏览器 → 页面变美观
```

