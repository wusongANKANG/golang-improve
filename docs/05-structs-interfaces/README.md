# 模块 05：结构体、方法、接口

这一章几乎决定了你能不能顺畅地读懂 Go 业务代码。因为 Go 没有传统面向对象语言里那种“类 + 继承”的强主线，所以很多抽象能力都落在了结构体、方法、接口、组合这几个点上。

## 学习目标

- 理解结构体是 Go 的核心建模方式。
- 掌握值接收者和指针接收者的区别。
- 理解接口是“行为抽象”，不是“继承体系”。
- 理解组合和嵌入在 Go 中的使用方式。

## 一、结构体是 Go 的主要数据模型

在 Go 里，结构体就是“把多个相关字段放在一起”的最常见方式：

```go
type User struct {
	Name  string
	Email string
}
```

这和其他语言里的 class 有点像，但要注意：

- 结构体本身只是数据定义。
- 行为需要通过方法补上。
- Go 不强调继承树，而更强调组合与接口。

## 二、方法是什么

方法本质上就是“绑定在某个类型上的函数”。

当前模块里的示例：

```go
func (u User) Notify() string {
	return fmt.Sprintf("notify %s via %s", u.Name, u.Email)
}
```

这里的 `User` 就是接收者类型。你可以把它理解成“这个函数属于 `User` 这个类型”。

调用方式：

```go
user := User{Name: "alice", Email: "alice@example.com"}
msg := user.Notify()
```

## 三、值接收者和指针接收者怎么区分

这是结构体相关最重要的知识点之一。

当前模块有两个方法：

```go
func (u User) Notify() string
func (u *User) Rename(name string)
```

它们分别代表两种语义：

- 值接收者：方法操作的是对象副本，更适合只读逻辑。
- 指针接收者：方法操作的是原对象，更适合修改状态或避免大对象复制。

看这个例子：

```go
func (u *User) Rename(name string) {
	u.Name = name
}
```

如果这里不用指针接收者，而是写成：

```go
func (u User) Rename(name string) {
	u.Name = name
}
```

那修改的就只是副本，外部对象不会变化。

## 四、为什么 `Rename` 必须用指针接收者

当前测试里：

```go
user := User{Name: "alice", Email: "alice@example.com"}
user.Rename("bob")
```

最后断言 `user.Name == "bob"` 成立，是因为 `Rename` 操作的是原对象。

这里可以顺手建立一个实用判断标准：

- 方法要修改对象状态时，用指针接收者。
- 类型较大时，通常也优先考虑指针接收者。
- 只读、小对象、明确不改状态时，值接收者更自然。

## 五、接口的核心思想是什么

Go 接口不是“父类”，而是“行为描述”。

当前模块里的接口定义：

```go
type Notifier interface {
	Notify() string
}
```

这句的意思非常简单：

- 只要某个类型有 `Notify() string` 方法
- 它就自动实现了 `Notifier`

Go 不需要显式写：

```go
implements Notifier
```

这就是 Go 接口的一个非常重要的特点：隐式实现。

## 六、一个结构体能实现多个接口吗

可以，而且这是 Go 里非常常见的用法。

因为 Go 的接口是隐式实现的，所以一个结构体能实现多少个接口，取决于它的方法集合满足多少个接口。

比如当前模块可以定义两个接口：

```go
type Notifier interface {
	Notify() string
}

type Renamer interface {
	Rename(name string)
}
```

`User` 同时拥有这两个方法：

```go
func (u User) Notify() string {
	return fmt.Sprintf("notify %s via %s", u.Name, u.Email)
}

func (u *User) Rename(name string) {
	u.Name = name
}
```

于是同一个 `*User` 可以同时当作 `Notifier` 和 `Renamer` 使用：

```go
user := &User{Name: "alice", Email: "alice@example.com"}

var notifier Notifier = user
var renamer Renamer = user

renamer.Rename("bob")
fmt.Println(notifier.Notify()) // notify bob via alice@example.com
```

这里有一个非常重要的细节：

- `Notify()` 是值接收者，所以 `User` 和 `*User` 都能实现 `Notifier`。
- `Rename()` 是指针接收者，所以只有 `*User` 实现了 `Renamer`。

也就是说，下面这个写法是可以的：

```go
var notifier Notifier = User{Name: "alice", Email: "alice@example.com"}
```

但下面这个写法不行：

```go
var renamer Renamer = User{Name: "alice", Email: "alice@example.com"}
```

因为 `Rename` 需要修改原对象，它的方法接收者是 `*User`。

## 七、接口为什么强调小而专

Go 标准库里很多经典接口都很小：

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

这背后的设计哲学是：

- 小接口更容易被实现。
- 小接口更容易被复用。
- 小接口更容易解耦。

当前模块的 `Notifier` 也遵循这个思路。它只关心“能不能通知”，而不关心你是不是用户、管理员、机器人。

## 八、组合和嵌入在 Go 里怎么工作

当前模块里还有一个类型：

```go
type Admin struct {
	User
	Level int
}
```

这里不是继承，而是嵌入。

这意味着：

- `Admin` 内部包含一个 `User`
- `Admin` 可以直接访问 `User` 的字段和方法

例如：

```go
admin := Promote(User{Name: "root", Email: "root@example.com"}, 10)
fmt.Println(admin.Name)
fmt.Println(admin.Notify())
```

即使 `Admin` 没有显式写 `Notify()`，它依然能调用，因为这个方法来自嵌入的 `User`。

## 九、为什么 `Admin` 也能当 `Notifier`

这是接口和嵌入一起工作时非常值得理解的点。

因为 `Admin` 通过嵌入 `User` 获得了 `Notify()` 方法，所以它也满足：

```go
type Notifier interface {
	Notify() string
}
```

于是下面的代码可以成立：

```go
messages := SendAll([]Notifier{user, admin})
```

这说明 `SendAll` 完全不需要知道具体类型，它只要求传入的对象具备通知能力。

## 十、面向接口编程是什么意思

看 `SendAll`：

```go
func SendAll(notifiers []Notifier) []string {
	messages := make([]string, 0, len(notifiers))
	for _, notifier := range notifiers {
		messages = append(messages, notifier.Notify())
	}

	return messages
}
```

这段代码的重点不是循环，而是抽象边界：

- `SendAll` 不关心底层具体类型。
- 它只依赖行为约定，也就是 `Notify() string`。
- 这让调用方可以自由传入不同实现。

这就是 Go 里很常见的解耦方式。

## 十一、配套测试怎么读

测试文件在 [structs_interfaces_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/05_structs_interfaces/structs_interfaces_test.go)。

建议重点看两个测试：

```go
func TestRename(t *testing.T)
func TestSendAll(t *testing.T)
func TestUserCanImplementMultipleInterfaces(t *testing.T)
```

它们分别对应：

- 指针接收者修改对象状态。
- 接口屏蔽具体类型差异。
- 一个结构体同时满足多个接口。

## 十二、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/05_structs_interfaces
```

调试时建议重点看：

- `Rename` 执行前后，`user.Name` 如何变化。
- `Admin` 虽然没有显式实现 `Notify`，为什么仍能进入 `SendAll`。
- `[]Notifier{user, admin}` 里的元素在运行时分别是什么具体类型。
- `*User` 为什么能同时赋值给 `Notifier` 和 `Renamer`。

## 十三、常见误区

- 把接口理解成继承层次。
- 动不动就定义大接口。
- 明明要修改对象，却用了值接收者。
- 接口定义放在实现方，导致抽象过早、过重。
- 忽略值接收者和指针接收者对接口实现的影响。

## 十四、工作里的映射

这一章和实际项目高度相关：

- `User`、`Order`、`Config` 这类业务对象通常用结构体建模。
- service、repo、client 经常通过接口解耦。
- 指针接收者和嵌入在业务状态变更中非常常见。

## 十五、建议练习

你可以自己补两个练习：

```go
type SMSNotifier struct {
	Phone string
}

func (s SMSNotifier) Notify() string
func RenameAll(users []*User, prefix string)
```

这样能把“接口多实现”和“指针接收者修改对象”再练一遍。
