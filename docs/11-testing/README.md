# 模块 11：测试、基准测试、竞态检测

这一章不是单纯讲测试语法，而是把“怎么验证你真的掌握了知识点”这个问题讲清楚。对于学习 Go 来说，测试不是附属品，而是最好的练习和 debug 入口。

## 学习目标

- 掌握 Go 单元测试的基础写法。
- 理解表驱动测试和子测试。
- 学会跑 benchmark。
- 理解 `-race` 的作用和限制。

## 一、Go 测试的基本结构

一个最基本的测试函数长这样：

```go
func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Fatalf("Add() = %d, want 5", got)
	}
}
```

这里需要记住几个最小规则：

- 文件名通常是 `*_test.go`。
- 函数名以 `Test` 开头。
- 参数必须是 `t *testing.T`。

## 二、为什么 `got / want` 这种写法很常见

因为它简洁、统一、利于阅读：

```go
got := Add(2, 3)
want := 5
```

这在团队协作里非常有价值。测试越统一，大家越容易快速看懂失败原因。

## 三、表驱动测试是什么

Go 社区非常常见的测试风格就是表驱动测试。看当前模块的 `TestDivide`：

```go
testCases := []struct {
	name    string
	a       float64
	b       float64
	want    float64
	wantErr error
}{
	{name: "normal", a: 10, b: 2, want: 5},
	{name: "zero divisor", a: 10, b: 0, wantErr: ErrDivideByZero},
}
```

它的好处是：

- 多个测试场景集中管理。
- 扩展 case 很方便。
- 每个 case 的输入输出边界清晰。

## 四、子测试为什么有用

继续看 `TestDivide`：

```go
for _, testCase := range testCases {
	t.Run(testCase.name, func(t *testing.T) {
		got, err := Divide(testCase.a, testCase.b)
		...
	})
}
```

`t.Run` 会把每个 case 作为独立子测试运行，这样好处是：

- 失败报告更清楚。
- 可以单独定位某个 case。
- 复杂场景下结构更清晰。

例如输出里你会直接看到：

- `TestDivide/normal`
- `TestDivide/zero_divisor`

## 五、基准测试是怎么写的

当前模块里有两个 benchmark：

```go
func BenchmarkAdd(b *testing.B)
func BenchmarkFibonacci(b *testing.B)
```

最常见的写法是：

```go
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Add(10, 20)
	}
}
```

这里的重点是：

- 框架会自动决定 `b.N` 的大小。
- 你只需要把待测逻辑放进循环里。

## 六、benchmark 不只是“谁快”

很多人第一次接触 benchmark 时只看 `ns/op`，但更有价值的往往还有：

- `B/op`
- `allocs/op`

所以建议经常用：

```bash
go test -bench . -benchmem ./examples/11_testing
```

这样你不仅看到速度，还能看到内存分配情况。

## 七、`-race` 是做什么的

`-race` 用来检测数据竞争。也就是多个 goroutine 在没有正确同步的情况下并发读写共享数据。

常用命令：

```bash
go test -race ./examples/09_concurrency
```

这对并发代码非常有帮助，因为很多数据竞争肉眼很难看出来。

## 八、为什么学习仓库里每章都配测试

因为测试在这里有三层意义：

- 它是行为验证。
- 它是最小例子。
- 它是最好的断点入口。

你完全可以把测试文件当成“带断言的练习题”。

## 九、当前模块里的示例怎么读

`Add` 的测试很基础：

```go
func TestAdd(t *testing.T)
```

`Divide` 的测试更有代表性，因为它展示了：

- 正常路径。
- 错误路径。
- 表驱动结构。
- 子测试组织方式。

而 `BenchmarkFibonacci` 则展示了如何对一个纯计算函数做最小性能测试。

## 十、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/11_testing
```

只跑一个测试：

```bash
go test -run TestDivide -v ./examples/11_testing
```

跑 benchmark：

```bash
go test -bench . -benchmem ./examples/11_testing
```

调试时建议重点看：

- `testCases` 在循环里是如何逐个执行的。
- 子测试失败时名字是怎么展示的。
- `BenchmarkAdd` 在 IDE 里虽然不常断点，但它的结构很适合作为性能验证模板记住。

## 十一、常见误区

- 只写 happy path，不写错误路径。
- 不做表驱动，导致重复测试代码很多。
- benchmark 只看快慢，不看内存分配。
- 并发代码不跑 `-race`。

## 十二、工作里的映射

这一章会直接影响：

- 你改代码时有没有安全感。
- 你能不能快速复现线上边界问题。
- 你能不能用 benchmark 支撑优化决策。
- 团队是否能稳定重构。

## 十三、建议练习

可以自己补两个内容：

```go
func TestFibonacci(t *testing.T)
func BenchmarkDivide(b *testing.B)
```

一个练表驱动正常函数测试，一个练 benchmark 写法。
