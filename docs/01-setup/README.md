# 模块 01：开发环境与程序结构

这个模块解决的是一个很基础、但很容易被忽略的问题：Go 程序到底是怎么组织起来的。很多人工作里会写 Go 代码，但如果对“包、目录、模块、入口、测试”的关系理解不清楚，后面学习切片、接口、并发时很容易总觉得代码是散的。

## 学习目标

- 理解 Go 程序最小可执行结构。
- 理解目录、包、模块之间的关系。
- 知道 `go run`、`go build`、`go test` 的区别。
- 建立“知识点也应该放进可运行工程里”的学习习惯。

## 一、Go 程序的最小骨架

一个最小 Go 程序通常长这样：

```go
package main

import "fmt"

func main() {
	fmt.Println("hello, gopher")
}
```

这里面有三个关键点：

- `package main` 表示这是一个可执行程序包。
- `func main()` 是程序入口。
- `import` 用来引入标准库或其他包。

如果一个目录里的代码不是 `package main`，那它通常不是一个直接执行的程序，而是一个供别人调用的包。

## 二、目录、包、模块是什么关系

这三个词很容易混：

- 目录：磁盘上的文件夹。
- 包：Go 的代码组织单位，通常一个目录对应一个包。
- 模块：依赖管理单位，通常由 `go.mod` 定义。

可以先用一句非常实用的话记忆：

- 写代码时，主要和“包”打交道。
- 管依赖时，主要和“模块”打交道。

在这个仓库里：

- 根目录有一个 [go.mod](/home/wusong/workspace/future-2026/golang-improve/go.mod)，表示整个仓库是一个模块。
- `examples/01_setup` 是一个包。
- `examples/06_packages_generics/wordutil` 也是一个包。

## 三、为什么学习 Go 时推荐优先用 `go test`

很多人学习新知识点时第一反应是 `go run`。这当然没错，但如果你的目标是“学会并记住”，`go test` 往往更适合。

原因很简单：

- 测试天然就是最小示例。
- 测试可以精确验证输出。
- 测试更容易单步调试。
- 后面你改动示例时，测试能帮你防止理解跑偏。

常用命令区别如下：

```bash
go run ./cmd/app
go build ./...
go test ./...
```

它们分别更适合：

- `go run`：快速执行一个程序入口。
- `go build`：确认代码是否能编译通过。
- `go test`：验证行为是否正确，也是最适合学习和 debug 的方式。

## 四、这个模块里的示例在讲什么

这个模块对应的代码在 [setup.go](/home/wusong/workspace/future-2026/golang-improve/examples/01_setup/setup.go)。

先看第一个例子：

```go
func Greeting(name string) string {
	if name == "" {
		name = "gopher"
	}

	return fmt.Sprintf("hello, %s", name)
}
```

这个例子很小，但它已经体现了几个 Go 的基础写法：

- 函数声明格式是 `func 名称(参数) 返回值类型`。
- `if` 不需要括号。
- 字符串格式化常用 `fmt.Sprintf`。
- 返回值写在参数列表后面。

再看第二个例子：

```go
func ProgramShape() []string {
	return []string{
		"package main",
		"import (...)",
		"func main() { ... }",
	}
}
```

这个函数不是为了业务逻辑，而是为了把“程序骨架”转成可测试的数据。这里很重要的一点是：学习代码也可以被写成函数、被测试、被断点调试，而不是只能靠看笔记。

## 五、配套测试怎么看

测试代码在 [setup_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/01_setup/setup_test.go)。

例如：

```go
func TestGreeting(t *testing.T) {
	if got := Greeting("alice"); got != "hello, alice" {
		t.Fatalf("Greeting() = %q, want %q", got, "hello, alice")
	}
}
```

这里可以顺手认识 Go 测试最基本的套路：

- 测试函数名必须以 `Test` 开头。
- 参数必须是 `t *testing.T`。
- 常见写法是 `got` / `want`。
- 判断不符合预期时用 `t.Fatalf` 直接终止当前测试。

## 六、怎么运行和怎么 debug

运行当前模块：

```bash
go test -v ./examples/01_setup
```

只跑一个测试：

```bash
go test -run TestGreeting -v ./examples/01_setup
```

如果你用 VS Code，可以直接在 [setup_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/01_setup/setup_test.go) 的 `TestGreeting` 上打断点，然后使用仓库里的 [launch.json](/home/wusong/workspace/future-2026/golang-improve/.vscode/launch.json) 进行调试。

调试时重点观察：

- `name == ""` 条件是否成立。
- `name` 在进入 `if` 前后如何变化。
- `Greeting("")` 和 `Greeting("alice")` 的返回结果有什么不同。

## 七、工作里的映射

这个模块虽然基础，但它对应的是非常实际的工程能力：

- 你能不能快速看懂一个 Go 项目的结构。
- 你是否知道一个目录里的代码为什么能被测试但不能直接运行。
- 你会不会用最小测试来验证一个函数，而不是靠肉眼猜结果。

如果这一步打牢，后面所有模块都会学得更顺。

## 八、你可以自己动手改的练习

建议你现在就试三个小改动：

1. 给 `Greeting` 增加一个默认前缀，比如返回 `welcome, gopher`。
2. 给 `ProgramShape` 再加一项，比如 `go test ./...`。
3. 自己补一个 `TestGreetingDefaultName`，只测试空字符串场景。

这类改动很小，但非常适合把“看懂”变成“会写”。
