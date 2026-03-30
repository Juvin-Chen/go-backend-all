# GORM 核心常用方法大全
基于**学生(Student) + 借阅记录(BorrowRecord)** 场景编写
覆盖：**基础操作、条件查询、关联查询、原生SQL** 所有高频 `db.()` 方法


---

## 前置准备（共用结构体）
以下所有示例，都基于这两个结构体（你写的借伞业务）：
```go
// 学生表（一方）
type Student struct {
	gorm.Model
	Name    string
	Records []BorrowRecord // 一对多：收纳盒（假字段）
}

// 借阅记录表（多方）
type BorrowRecord struct {
	gorm.Model
	StudentID  int    // 外键（真字段）
	UmbrellaSN string // 伞编号
	Status     string // 状态：borrowing/returned
}
```

---

# 一、数据库初始化（必用）
## 1. `db.AutoMigrate()` 自动建表/更新表
**作用**：根据结构体，自动创建数据库表，新增字段时自动更新
```go
// 自动创建 students、borrow_records 两张表
db.AutoMigrate(&Student{}, &BorrowRecord{})
```

---

# 二、基础增删改查（CRUD）
## 🔹 新增（Create）
### 1. `db.Create()` 单条新增
**作用**：插入一条数据到数据库
```go
// 新增单个学生
student := Student{Name: "lisi"}
db.Create(&student)

// 【关联新增】新增学生+同时新增借阅记录（一对多）
student := Student{
	Name: "zhangsan",
	Records: []BorrowRecord{
		{UmbrellaSN: "UMB-001", Status: "borrowing"},
	},
}
db.Create(&student) // GORM自动保存两张表数据
```

---

## 🔹 查询（Find/First/Take）
### 1. `db.First()` 查询单条数据（按主键ID）
**作用**：查询**第一条**数据，找不到返回报错
```go
var student Student
db.First(&student, 1) // 查询ID=1的学生
```

### 2. `db.Take()` 查询单条数据（无排序）
**作用**：随机查一条，比 First 更轻量
```go
var student Student
db.Take(&student)
```

### 3. `db.Find()` 查询多条数据
**作用**：查询**所有匹配**的数据（列表）
```go
var students []Student
db.Find(&students) // 查询所有学生
```

---

## 🔹 更新（Save/Update）
### 1. `db.Save()` 全量更新
**作用**：更新整条数据（所有字段覆盖）
```go
var student Student
db.First(&student, 1)
student.Name = "wangwu" // 修改姓名
db.Save(&student)       // 保存更新
```

### 2. `db.Update()` 单字段更新
```go
db.Model(&Student{}).Where("id = ?", 1).Update("name", "zhaoliu")
```

### 3. `db.Updates()` 多字段更新
```go
// 批量更新字段
db.Model(&Student{}).Where("id = ?", 1).Updates(map[string]any{
	"name": "sunqi",
})
```

---

## 🔹 删除（Delete）
### 1. `db.Delete()` 删除数据
**作用**：删除数据（软删除，gorm.Model自带）
```go
// 删除ID=1的学生
db.Delete(&Student{}, 1)
```

---

# 三、条件查询（Where/Or/Not 最常用）
## 1. `db.Where()` 条件筛选（核心）
**作用**：指定查询条件，几乎所有查询都要用
```go
var student Student
// 1. 精确匹配：查询姓名为 zhangsan 的学生
db.Where("name = ?", "zhangsan").First(&student)

var records []BorrowRecord
// 2. 条件查询：查询未归还的借伞记录
db.Where("status = ?", "borrowing").Find(&records)

// 3. 多条件：查询学生1的已归还记录
db.Where("student_id = ? AND status = ?", 1, "returned").Find(&records)
```

## 2. `db.Or()` 或条件
```go
// 查询 未归还 或 伞编号为UMB-001 的记录
db.Where("status = ?", "borrowing").Or("umbrella_sn = ?", "UMB-001").Find(&records)
```

## 3. `db.Not()` 非条件
```go
// 查询 不是未归还 的记录
db.Not("status = ?", "borrowing").Find(&records)
```

## 4. `db.Order()` 排序
```go
// 按ID倒序查询所有记录
db.Order("id desc").Find(&records)
```

## 5. `db.Limit()` 限制条数 / `db.Offset()` 分页
```go
var records []BorrowRecord
// 查询前5条记录
db.Limit(5).Find(&records)

// 分页：跳过前2条，查3条
db.Offset(2).Limit(3).Find(&records)
```

## 6. `db.Count()` 统计总数
```go
var total int64
// 统计所有借伞记录数量
db.Model(&BorrowRecord{}).Count(&total)
```

---

# 四、🔥 关联查询
专门解决：**查学生 → 带出他的所有借伞记录**

## 1. `db.Preload()` 预加载（一对多必用）
**作用**：查询主数据时，**顺带查询关联数据**（自动装满收纳盒）
```go
var student Student
// 查询张三 + 他的所有借阅记录
db.Preload("Records").Where("name = ?", "zhangsan").First(&student)

// 使用：直接取关联数据
fmt.Println(student.Records) // 输出该学生所有借伞记录
```

### ✅ 关键说明
- 不加 `Preload`：`student.Records` 是空的
- 加 `Preload`：GORM自动查询关联表，把数据装进切片

## 2. 多层预加载（进阶）
```go
// 如果有更深的关联，直接点语法
db.Preload("Records.Order").First(&student)
```

---

# 五、原生SQL操作（Exec/Raw）
## 1. `db.Exec()` 执行原生SQL（无返回数据）
**作用**：直接让数据库执行SQL语句（增删改）
```go
// 原生SQL：把学生1的所有借伞记录改为已归还
db.Exec("UPDATE borrow_records SET status = ? WHERE student_id = ?", "returned", 1)
```

## 2. `db.Raw()` 执行原生SQL查询（有返回数据）
```go
var records []BorrowRecord
// 原生SQL查询
db.Raw("SELECT * FROM borrow_records WHERE student_id = ?", 1).Scan(&records)
```

---

# 六、常用辅助方法
## 1. `db.Model()` 指定操作表
**作用**：告诉GORM当前操作哪张表（更新/删除必用）
```go
// 指定操作 Student 表
db.Model(&Student{}).Where("id=?", 1).Update("name", "test")
```

## 2. `db.Error` 获取执行错误
```go
err := db.First(&student, 999).Error
if err != nil {
	fmt.Println("查询失败：", err)
}
```

---

# 七、GORM 方法速查表（闭眼查）
| 方法                | 作用                                   | 适用场景                     |
| ------------------- | -------------------------------------- | ---------------------------- |
| `AutoMigrate`       | 自动建表/更新表                        | 项目初始化                   |
| `Create`            | 新增数据                               | 注册、添加记录               |
| `First`             | 查询单条（主键）| 查详情                       |
| `Find`              | 查询多条/列表                          | 列表展示                     |
| `Save`/`Updates`    | 更新数据                               | 修改信息                     |
| `Delete`            | 删除数据                               | 删除记录                     |
| `Where`             | 条件筛选                               | 精准查询                     |
| `Order`/`Limit`     | 排序/分页                              | 列表页                       |
| `Preload`           | 关联查询（一对多）| 查主数据+带出关联数据（核心） |
| `Exec`              | 执行原生SQL                            | 复杂更新/删除                |
| `Raw`               | 原生SQL查询                            | 复杂查询                     |

---

# 八、3条核心
1. **基础操作**：增(Create)、删(Delete)、改(Updates)、查(Find/First)
2. **条件查询**：所有筛选都用 `Where`
3. **关联查询**：一对多必须用 `Preload`，否则关联数据为空

