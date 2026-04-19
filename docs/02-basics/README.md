# 模块 02：变量、常量、类型与零值

这一章是 Go 基础里的基础，但也是工作中最容易“似懂非懂”的部分。因为变量、常量、类型这些语法大家都见过，可一旦遇到零值、自定义类型、类型推导、业务语义表达，就很容易模糊。

## 学习目标

- 重新建立对变量、常量、类型推导的基本感觉。
- 掌握 Go 的零值设计。
- 理解自定义类型在业务代码中的意义。
- 通过小例子熟悉多返回值和类型语义。

## 一、变量声明有哪几种写法

Go 里最常见的变量声明方式有两种：

```go
var age int = 18
name := "alice"
```

一般可以这样理解：

- `var` 更适合需要显式类型、需要零值初始化、或者包级变量场景。
- `:=` 更适合函数内部的简洁声明。

再看一个典型例子：

```go
var count int
fmt.Println(count) // 0
```

这里即使没有显式赋值，`count` 也能直接用，因为 Go 提供了零值。

## 二、什么是零值

零值是 Go 的一个非常重要的设计理念。简单说就是：变量没有显式初始化时，会自动拥有该类型的默认值。

常见零值如下：

```go
var n int        // 0
var s string     // ""
var ok bool      // false
var nums []int   // nil
var m map[string]int // nil
```

这套设计的好处是：

- 很多类型可以直接进入可用状态。
- 初始化代码更少。
- 代码里“默认状态”更清晰。

## 三、结构体的零值为什么很常见

看当前模块的示例：

```go
type Profile struct {
	Name   string
	Age    int
	Active bool
	Status Status
}

func ZeroValueProfile() Profile {
	return Profile{}
}
```

如果你调用 `ZeroValueProfile()`，得到的是：

- `Name == ""`
- `Age == 0`
- `Active == false`
- `Status == ""`

这就是结构体零值展开后的效果。

在工作里，这种设计很常见，因为一个结构体可以先作为“默认空状态”存在，再逐步填充值。

## 四、自定义类型和业务语义

当前模块里有这样一个定义：

```go
type Status string

const (
	StatusUnknown Status = "unknown"
	StatusActive  Status = "active"
)
```

这背后有个非常重要的思路：虽然底层还是 `string`，但 `Status` 已经不是普通字符串，而是“带业务语义的字符串”。

这样做的好处：

- 提升可读性。
- 降低把任意字符串误当状态值使用的概率。
- 后续更容易集中维护状态集合。

例如：

```go
type OrderStatus string
type UserStatus string
```

即使底层都还是 `string`，但语义已经不同了。

## 五、构造函数式写法

当前模块里的另一个示例：

```go
func NewProfile(name string, age int) Profile {
	return Profile{
		Name:   name,
		Age:    age,
		Active: true,
		Status: StatusActive,
	}
}
```

虽然 Go 没有构造函数语法，但我们经常会通过 `NewXxx` 函数表达“推荐初始化方式”。

这个例子传达两个学习点：

- 结构体字面量是 Go 中非常常见的初始化方式。
- `NewXxx` 不一定返回指针，也可以返回值，具体要看是否需要共享和修改。

## 六、多返回值是 Go 的日常语法

看这个函数：

```go
func Swap(a, b int) (int, int) {
	return b, a
}
```

多返回值在 Go 里非常自然，它不仅可以用来返回多个业务结果，也经常被用来返回“结果 + 错误”。

调用方式：

```go
left, right := Swap(1, 2)
```

这也是为什么后面你会频繁看到这种写法：

```go
value, err := someFunc()
```

## 七、常量为什么值得单独理解

再看当前模块的例子：

```go
func TypedAndUntypedConstants() (int64, float64) {
	const maxRetries int64 = 3
	const pi = 3.14

	return maxRetries, pi
}
```

这里展示了两种常量：

- 有类型常量：`maxRetries`
- 无类型常量：`pi`

无类型常量的灵活性更高，因为在很多场景下它可以根据上下文自动适配目标类型。

这也是为什么你经常能直接写：

```go
const timeout = 3
var x int64 = timeout
```

## 八、配套测试应该怎么看

测试文件在 [basics_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/02_basics/basics_test.go)。

例如这个测试：

```go
func TestZeroValueProfile(t *testing.T) {
	profile := ZeroValueProfile()

	if profile.Name != "" || profile.Age != 0 || profile.Active || profile.Status != "" {
		t.Fatalf("unexpected zero value profile: %+v", profile)
	}
}
```

这个测试最适合用来重新建立你对零值的直觉。建议你直接打断点看 `profile` 的每个字段。

## 九、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/02_basics
```

只跑零值相关测试：

```bash
go test -run TestZeroValueProfile -v ./examples/02_basics
```

调试时建议重点看：

- `Profile{}` 创建出来后，每个字段是什么。
- `NewProfile("alice", 18)` 和零值结构体相比多了哪些业务默认值。
- `Swap(1, 2)` 的返回值是怎么落到左值变量里的。

## 十、常见误区

- 误以为“没初始化就不能用”。
- 所有业务字段都直接用基础类型，不表达语义。
- 把 `NewXxx` 机械地理解成必须返回指针。
- 不理解 Go 为什么喜欢多返回值。

## 十一、工作里的映射

你在真实项目里会频繁看到这些模式：

- 用零值结构体作为默认配置。
- 用自定义类型表达状态、枚举、ID 语义。
- 用 `NewXxx` 统一创建对象。
- 用多返回值返回结果和错误。

如果你把这一章吃透，后面读业务代码时会顺畅很多。

## 十二、建议练习

你可以自己加两个练习函数：

```go
func NewGuestProfile() Profile
func IsActive(status Status) bool
```

然后为它们各补一个测试。这样比单纯读文档更容易留下长期记忆。
