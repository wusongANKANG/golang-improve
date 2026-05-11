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

## 三、闭包是什么

先给一个尽量准确但不绕的定义：

闭包可以理解成两部分的组合：

- 一个函数。
- 这个函数创建时所在作用域里的变量环境。

也就是说，闭包不是“单独的函数语法”，而是“函数 + 它能访问并记住的外部变量”。

看当前模块的示例：

```go
func Accumulator(start int) func(delta int) int {
	sum := start

	return func(delta int) int {
		sum += delta
		return sum
	}
}
```

这里返回的匿名函数就是闭包，因为它使用了自己参数之外的变量 `sum`。

## 四、为什么闭包能“记住状态”

先看调用：

```go
acc := Accumulator(10)
fmt.Println(acc(5)) // 15
fmt.Println(acc(3)) // 18
```

很多人第一次看到这里会疑惑：

- `Accumulator(10)` 不是早就执行结束了吗？
- 为什么第二次调用 `acc(3)` 时，还能接着上一次的结果继续算？

关键就在于：闭包捕获的不是“一次性的值替换结果”，而是它所引用的那个外部变量。

可以把这段逻辑拆成下面几步：

第一步，调用 `Accumulator(10)`：

```go
sum := 10
```

第二步，返回一个匿名函数，这个函数内部会访问 `sum`：

```go
func(delta int) int {
	sum += delta
	return sum
}
```

第三步，虽然 `Accumulator` 函数本身返回了，但因为返回出去的闭包后面还要继续使用 `sum`，所以这个 `sum` 不会随着函数结束就“消失掉”。

第四步，每次调用 `acc`，其实都是在继续操作同一个 `sum`：

```go
acc(5) // sum: 10 -> 15
acc(3) // sum: 15 -> 18
```

所以闭包最本质的能力不是“返回函数”，而是“返回一个带着上下文状态的函数”。

## 五、闭包捕获的到底是什么

这是理解闭包最关键的一层。

很多初学者容易把闭包理解成：

- “函数把外部变量复制了一份”

这并不准确。

更接近真实行为的理解是：

- 闭包持有了对外部变量的访问能力。
- 只要这个闭包后面还会用到这个变量，这个变量的生命周期就会被延长。

对于 `Accumulator` 来说，`sum` 不是在第一次调用 `acc` 之后重新生成的，它始终是同一个被闭包持有的状态。

所以这段代码：

```go
acc := Accumulator(10)
fmt.Println(acc(1)) // 11
fmt.Println(acc(1)) // 12
fmt.Println(acc(1)) // 13
```

本质上是在不断修改同一个 `sum`。

## 六、一个闭包实例和多个闭包实例有什么区别

这个点非常重要。

看下面两段代码：

```go
acc := Accumulator(10)
fmt.Println(acc(5)) // 15
fmt.Println(acc(3)) // 18
```

这是“同一个闭包实例”被连续调用，所以它共享同一个 `sum`。

再看：

```go
acc1 := Accumulator(10)
acc2 := Accumulator(10)

fmt.Println(acc1(5)) // 15
fmt.Println(acc1(3)) // 18

fmt.Println(acc2(5)) // 15
fmt.Println(acc2(3)) // 18
```

这里 `acc1` 和 `acc2` 是两次独立调用 `Accumulator` 得到的两个闭包实例。

它们看起来逻辑一样，但内部状态彼此独立：

- `acc1` 有自己的 `sum`
- `acc2` 也有自己的 `sum`

这点很像“同一个模板函数，创建出了两个各自维护状态的小对象”。

## 七、把闭包当成“带状态的函数对象”更容易理解

虽然 Go 没有 class 风格的函数对象写法，但从使用效果上看，闭包常常可以类比为一个“隐藏了内部字段的对象”。

比如 `Accumulator(10)` 返回的闭包，从效果上有点像这样一个概念模型：

```go
type Counter struct {
	sum int
}

func (c *Counter) Add(delta int) int {
	c.sum += delta
	return c.sum
}
```

当然，闭包和结构体方法不是完全同一回事，但这个类比很有帮助：

- 结构体方法通过字段保存状态。
- 闭包通过捕获外部变量保存状态。

如果你能接受“闭包就是一个隐藏了状态存储细节的函数值”，理解会顺很多。

## 八、再看一个更贴近日常的闭包例子

比如你想生成一个带固定前缀的日志函数：

```go
func PrefixLogger(prefix string) func(msg string) string {
	return func(msg string) string {
		return prefix + ": " + msg
	}
}
```

调用：

```go
apiLogger := PrefixLogger("API")
dbLogger := PrefixLogger("DB")

fmt.Println(apiLogger("request start")) // API: request start
fmt.Println(dbLogger("query failed"))   // DB: query failed
```

这里每个返回的函数都记住了各自的 `prefix`。

这个例子说明闭包不一定总是拿来“累加数字”，它也经常用来：

- 预先绑定配置。
- 生成定制化处理函数。
- 简化重复参数传递。

## 九、闭包在工作里通常用来做什么

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

再举几个工作里的典型使用场景：

1. 中间件

```go
func WithTrace(traceID string) func(string) string {
	return func(msg string) string {
		return "[" + traceID + "] " + msg
	}
}
```

2. 重试包装器

```go
func Retry(times int, fn func() error) error {
	var err error
	for i := 0; i < times; i++ {
		if err = fn(); err == nil {
			return nil
		}
	}
	return err
}
```

虽然这里 `Retry` 本身没有返回闭包，但它体现了函数值和闭包经常一起出现的风格。

## 十、闭包在生产中的优势是什么

闭包之所以在真实项目里常见，不是因为它“语法高级”，而是因为它很适合把一小段逻辑和它依赖的上下文绑在一起。

你可以把它理解成：

- 普通函数更像“每次都要把原材料带齐再调用”。
- 闭包更像“先把一部分上下文预装进去，再得到一个可重复调用的函数”。

它在生产里的主要优势通常有这些：

### 1. 减少重复传参

如果某段逻辑总要依赖同一组小配置，比如：

- `prefix`
- `traceID`
- `timeout`
- `retryTimes`

那用闭包先绑定这些上下文，调用点会简洁很多。

例如：

```go
func PrefixLogger(prefix string) func(msg string) string {
	return func(msg string) string {
		return prefix + ": " + msg
	}
}
```

这样后面每次调用都不用重新传 `prefix`。

### 2. 保持局部状态封装

闭包非常适合保存“小而明确的状态”，比如：

- 当前累计值
- 调用次数
- 某次处理链路的上下文信息

这样状态不会暴露成包级变量，也不一定非要专门定义一个结构体来保存。

`Accumulator` 就是最典型的例子：

```go
acc := Accumulator(10)
```

`sum` 只活在这个闭包相关的上下文里，不会污染其他地方。

### 3. 让函数更容易组合

闭包和高阶函数搭配时很自然，尤其适合：

- middleware
- wrapper
- hook
- callback

因为这些场景本质上都在做一件事：

- 先基于外部上下文生成一个函数
- 再把这个函数传来传去、层层包装

这也是为什么 Go 里很多中间件写法都和闭包很像。

### 4. 复用逻辑但保留差异化配置

同一套逻辑经常会有多个“只差一点配置”的版本。

比如：

```go
apiLogger := PrefixLogger("API")
dbLogger := PrefixLogger("DB")
auditLogger := PrefixLogger("AUDIT")
```

三者逻辑一样，但各自保留自己的上下文。闭包非常适合这种“同逻辑，多实例”的模式。

### 5. 对简单场景比定义结构体更轻量

如果你只是想：

- 绑定一点小配置
- 保持一点局部状态
- 生成一个临时函数

那闭包通常比单独定义一个小 `struct + method` 更轻。

当然，前提是状态不要太复杂。状态一旦变大、职责一旦变多，通常还是结构体更清晰。

## 十一、什么时候特别适合用闭包

在 Go 项目里，下面这些场景特别常见：

- 生成带固定前缀的 logger
- 生成带配置的校验函数
- HTTP / RPC middleware 包装
- 重试、熔断、限流包装器
- 测试里快速构造 stub / mock 行为
- 回调函数需要顺手带上一点上下文

这些场景的共同点是：

- 上下文不大
- 状态不复杂
- 逻辑需要被传递或组合

这正是闭包最舒服的工作区间。

## 十二、闭包最容易踩的坑是什么

### 1. 误以为每次调用拿到的是新状态

比如：

```go
acc := Accumulator(10)
```

很多人会潜意识里把 `acc(5)`、`acc(3)` 当成两个互不相关的函数调用，但它们其实共享同一份被捕获的状态。

### 2. 不小心共享了同一个外部变量

例如：

```go
total := 0

add := func(v int) int {
	total += v
	return total
}

sub := func(v int) int {
	total -= v
	return total
}
```

这里 `add` 和 `sub` 看起来是两个不同函数，但它们共同操作的是同一个 `total`。这类代码如果不注意，后面会很难排查状态变化来源。

### 3. 在循环中错误理解闭包和变量关系

循环变量捕获一直是 Go 里的经典话题。现代 Go 版本对 `range` 的变量语义已经做了改进，但你仍然要记住一个更本质的原则：

- 如果多个闭包共享的是同一个外部变量，那它们看到的就是同一份状态。

所以真正该理解的不是“背某个坑”，而是：

- 闭包关心的是它引用了哪个变量。
- 不是你主观上“希望它记住哪个值”。

## 十三、怎么 debug 闭包最有效

如果你想真正把闭包看明白，最推荐直接调试 [functions_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/04_functions/functions_test.go) 里的 `TestAccumulator`。

```go
func TestAccumulator(t *testing.T) {
	acc := Accumulator(10)

	if got := acc(5); got != 15 {
		t.Fatalf("first acc() = %d, want 15", got)
	}

	if got := acc(3); got != 18 {
		t.Fatalf("second acc() = %d, want 18", got)
	}
}
```

调试时重点看：

- `acc := Accumulator(10)` 执行后，返回的到底是什么。
- 第一次执行 `acc(5)` 时，`sum` 的值如何变化。
- 第二次执行 `acc(3)` 时，为什么起点不是 `10`，而是上一次的 `15`。

如果你愿意再进一步，可以自己加一个测试：

```go
func TestAccumulatorInstances(t *testing.T) {
	acc1 := Accumulator(10)
	acc2 := Accumulator(10)

	if got := acc1(5); got != 15 {
		t.Fatalf("acc1 first = %d, want 15", got)
	}

	if got := acc2(5); got != 15 {
		t.Fatalf("acc2 first = %d, want 15", got)
	}
}
```

这个测试特别适合帮助你区分：

- “同一个闭包多次调用”
- “多个闭包实例各自维护状态”

## 十四、闭包这一节你最少要带走什么

如果只记三句话，我建议记这三句：

1. 闭包是“函数 + 它捕获的外部变量环境”。
2. 闭包可以让函数在返回后继续持有并操作外部状态。
3. 同一个闭包实例会共享同一份状态，不同闭包实例通常各自维护自己的状态。

把这三句理解透，后面你再看中间件、回调、工厂函数、带配置的 handler，就会自然很多。

## 十五、`defer` 到底什么时候执行

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

这两个练习现在已经在仓库里实现好了，对应代码在：

- [functions.go](/home/wusong/workspace/future-2026/golang-improve/examples/04_functions/functions.go)
- [functions_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/04_functions/functions_test.go)

建议你优先调试下面几个测试：

- `TestRetryEventuallySuccess`
- `TestRetryExhausted`
- `TestSafeCall`
- `TestSafeCallRecover`

这几组测试很适合观察：

- 函数作为参数被重复调用的过程。
- `Retry` 如何控制重试次数。
- `SafeCall` 如何把 `panic` 转成普通 `error` 返回。
