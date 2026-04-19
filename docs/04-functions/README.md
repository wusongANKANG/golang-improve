# 模块 04：函数、闭包、defer、panic、recover

这一章是 Go 编程风格非常集中的一章。很多 Go 代码之所以看起来“很像 Go”，就是因为函数、多返回值、`defer`、错误处理这些能力被大量使用。

## 学习目标

- 掌握 Go 函数的声明和调用风格。
- 理解闭包为什么可以保存状态。
- 理解 `defer` 的执行时机。
- 分清 `error`、`panic`、`recover` 各自适用的场景。

## 一、Go 函数最常见的样子

Go 函数声明形式非常直接：

```go
func Add(a, b int) int {
	return a + b
}
```

但 Go 最有代表性的地方不是单返回值，而是多返回值：

```go
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}

	return a / b, nil
}
```

这就是当前模块里的核心示例之一。

## 二、为什么 Go 喜欢“结果 + 错误”

Go 不依赖异常来处理普通业务错误，所以很多函数会直接返回：

- 一个真正的结果
- 一个表示是否出错的 `error`

调用方式通常是：

```go
result, err := Divide(10, 2)
if err != nil {
	return err
}
```

这种风格一开始可能会觉得啰嗦，但它有几个明显优点：

- 错误路径非常显式。
- 更容易顺着调用链排查问题。
- 函数行为边界清晰。

## 三、闭包是什么，为什么它能“记住状态”

看当前模块的闭包示例：

```go
func Accumulator(start int) func(delta int) int {
	sum := start

	return func(delta int) int {
		sum += delta
		return sum
	}
}
```

调用：

```go
acc := Accumulator(10)
fmt.Println(acc(5)) // 15
fmt.Println(acc(3)) // 18
```

这里的关键不是语法，而是理解这件事：

- `Accumulator` 返回了一个函数。
- 这个返回的函数捕获了外层的 `sum` 变量。
- 即使 `Accumulator` 已经执行完，`sum` 依然会被保留下来。

这就是闭包。

## 四、闭包在工作里通常用来做什么

闭包常见用途：

- 保存局部状态。
- 封装配置后的处理逻辑。
- 生成函数工厂。
- 在中间件、回调、重试器里封装行为。

例如：

```go
func PrefixLogger(prefix string) func(msg string) string {
	return func(msg string) string {
		return prefix + ": " + msg
	}
}
```

这里返回的函数会一直记住 `prefix`。

## 五、`defer` 到底什么时候执行

`defer` 会在当前函数返回前执行，通常用于“善后动作”。

典型例子：

```go
func ReadFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
```

它最常见的用途：

- 关闭文件。
- 释放锁。
- 打点日志。
- 包装错误。
- 和 `recover` 配合兜底。

## 六、`panic` 和 `error` 有什么区别

要先建立一个非常重要的原则：

- 普通业务错误，用 `error`。
- 程序已经进入异常状态、无法继续安全运行时，才考虑 `panic`。

例如除零这种业务可预期错误，不应该 `panic`，而应该返回：

```go
return 0, ErrDivideByZero
```

而 `panic` 更像是：

- 数组越界
- 明显不该发生的内部状态损坏
- 框架层防御性兜底

## 七、`recover` 只能在 `defer` 里生效

当前模块的示例：

```go
func SafeExecute(fn func() string) (result string, recovered any) {
	defer func() {
		if r := recover(); r != nil {
			recovered = r
			result = fmt.Sprintf("recovered: %v", r)
		}
	}()

	return fn(), nil
}
```

这段代码非常适合理解 `panic` 和 `recover` 的关系：

- `fn()` 可能 panic。
- 外层的 deferred 函数会在函数返回前执行。
- `recover()` 可以拦截 panic，避免程序继续崩掉。

调用示例：

```go
result, recovered := SafeExecute(func() string {
	panic("boom")
})
```

返回会是：

- `result == "recovered: boom"`
- `recovered == "boom"`

## 八、为什么 `recover` 不该被滥用

很多初学者会觉得：既然能 `recover`，那是不是以后都不用好好处理错误了？

不是。

`recover` 更适合：

- HTTP 中间件兜底，避免单个请求把整个服务打崩。
- 消费队列 worker 的保护层。
- 框架边界的异常保护。

它不适合：

- 用来代替普通业务分支。
- 用来吞掉本该上报的错误。
- 用来隐藏程序设计问题。

## 九、配套测试怎么读

测试文件在 [functions_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/04_functions/functions_test.go)。

建议重点看两个测试。

第一个是错误处理：

```go
func TestDivideByZero(t *testing.T) {
	_, err := Divide(10, 0)
	if !errors.Is(err, ErrDivideByZero) {
		t.Fatalf("Divide() error = %v, want %v", err, ErrDivideByZero)
	}
}
```

第二个是 panic 恢复：

```go
func TestSafeExecuteRecover(t *testing.T) {
	result, recovered := SafeExecute(func() string {
		panic("boom")
	})
}
```

这两个测试基本把“普通错误”和“异常崩溃”两条路径都展示出来了。

## 十、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/04_functions
```

只调试闭包：

```bash
go test -run TestAccumulator -v ./examples/04_functions
```

只调试恢复逻辑：

```bash
go test -run TestSafeExecuteRecover -v ./examples/04_functions
```

调试时建议重点看：

- `Accumulator` 返回后，`sum` 为什么还在。
- `Divide(10, 0)` 为什么不崩溃而是返回错误。
- `SafeExecute` 里的 deferred 函数是在什么时候触发的。

## 十一、常见误区

- 把 `panic` 当成普通错误处理。
- 认为闭包只是语法糖，不理解它会保留状态。
- 不清楚 `defer` 的执行时机。
- 误以为 `recover` 在任何位置都能用。

## 十二、工作里的映射

这一章几乎覆盖了 Go 日常函数设计风格：

- 业务函数大量使用“结果 + 错误”。
- 中间件和框架层经常使用 `defer`。
- 闭包在回调、封装、配置注入里非常常见。
- `panic/recover` 常出现在兜底层，而不是业务层。

## 十三、建议练习

可以自己补两个练习：

```go
func Retry(times int, fn func() error) error
func SafeCall(fn func() int) (int, error)
```

第一个用来练习函数作为参数。
第二个用来练习 `defer + recover + error` 的组合。
