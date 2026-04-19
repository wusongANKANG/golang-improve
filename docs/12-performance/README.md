# 模块 12：性能优化与常见陷阱

这一章的目标不是把你变成性能专家，而是帮你建立最基础、最有收益的性能直觉。Go 的很多性能问题，并不需要特别复杂的工具才能发现，往往只是字符串拼接、切片扩容、共享底层数组这些基础行为没想清楚。

## 学习目标

- 理解常见性能热点来自哪里。
- 理解字符串拼接和切片预分配的常见差异。
- 理解什么时候应该复制切片，避免意外持有大内存。
- 学会用 benchmark 验证优化效果。

## 一、性能优化先看什么

先建立一个很重要的原则：

- 先保证正确性。
- 再定位热点。
- 最后做有数据支撑的优化。

不要一上来就把代码写得特别“技巧化”，结果可读性没了，收益却很小。

## 二、为什么字符串频繁拼接会慢

看当前模块的第一个例子：

```go
func JoinWithPlus(parts []string) string {
	joined := ""
	for _, part := range parts {
		joined += part
	}

	return joined
}
```

这段代码看起来很自然，但字符串在 Go 里是不可变的。每次 `joined += part`，通常都会产生新的字符串对象。

如果拼接次数很多，就会带来：

- 更多内存分配。
- 更多复制成本。

## 三、为什么 `strings.Builder` 更适合批量拼接

看另一个版本：

```go
func JoinWithBuilder(parts []string) string {
	var builder strings.Builder

	totalLength := 0
	for _, part := range parts {
		totalLength += len(part)
	}

	builder.Grow(totalLength)
	for _, part := range parts {
		builder.WriteString(part)
	}

	return builder.String()
}
```

这里有两个优化点：

- 用 `strings.Builder` 减少中间字符串创建。
- 用 `Grow(totalLength)` 预先申请容量，减少扩容。

这类优化在：

- 日志拼接
- SQL 组装
- 文本导出
- 模板构造

场景里都很常见。

## 四、切片为什么推荐预分配

看这两个函数：

```go
func BuildNumbersNoPrealloc(n int) []int
func BuildNumbersPrealloc(n int) []int
```

预分配版本：

```go
numbers := make([]int, 0, n)
for i := 0; i < n; i++ {
	numbers = append(numbers, i)
}
```

它的优势在于：

- 已知容量时，能减少多次扩容。
- 更少的扩容意味着更少的内存分配和数据复制。

如果你很清楚结果集大小，预分配通常是低成本高收益的优化。

## 五、什么是“意外持有大数组”

这是切片里一个特别值得警惕的问题。

例如你有一个非常大的切片：

```go
big := make([]int, 1000000)
small := big[:10]
```

虽然你只想保留前 10 个元素，但 `small` 仍然引用着那块巨大的底层数组。只要 `small` 还活着，那块大内存就可能无法释放。

当前模块用这个函数规避问题：

```go
func SafeSubset(items []int, n int) []int {
	if n > len(items) {
		n = len(items)
	}

	return append([]int(nil), items[:n]...)
}
```

这里通过复制创建了一份新的、小的切片，避免继续持有原大数组。

## 六、什么时候该共享，什么时候该复制

可以先记一个实用判断：

适合共享：

- 读多写少。
- 数据量不大。
- 明确知道上下游不会修改。

适合复制：

- 需要隔离修改。
- 要防止下游污染上游数据。
- 要截取小片段但不想保留大对象。

## 七、benchmark 才是优化验证依据

当前模块配了 benchmark：

```go
func BenchmarkJoinWithPlus(b *testing.B)
func BenchmarkJoinWithBuilder(b *testing.B)
func BenchmarkBuildNumbersNoPrealloc(b *testing.B)
func BenchmarkBuildNumbersPrealloc(b *testing.B)
```

你应该用 benchmark 来验证：

- 优化后是不是真的更快。
- 分配次数有没有下降。

命令：

```bash
go test -bench . -benchmem ./examples/12_performance
```

## 八、当前仓库里已经跑出的 benchmark 给了什么信号

之前这个仓库实际跑出的结果里，`JoinWithBuilder` 比 `JoinWithPlus` 分配更少，`BuildNumbersPrealloc` 比不预分配版本更快、分配更低。这正好印证了这章的两个核心结论：

- 批量字符串拼接时，`strings.Builder` 更稳。
- 已知容量时，切片预分配收益明显。

## 九、配套测试怎么读

测试文件在：

- [performance_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/12_performance/performance_test.go)
- [performance_benchmark_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/12_performance/performance_benchmark_test.go)

普通测试保证语义正确，benchmark 负责验证优化收益。

## 十、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/12_performance
```

跑 benchmark：

```bash
go test -bench . -benchmem ./examples/12_performance
```

调试时建议重点看：

- `JoinWithBuilder` 里为什么先统计总长度。
- `BuildNumbersPrealloc` 里 `cap(numbers)` 为什么一开始就是 `n`。
- `SafeSubset` 返回后为什么修改子切片不会影响原切片。

## 十一、常见误区

- 没有 benchmark 就盲目优化。
- 为了优化写出很难维护的代码。
- 不知道切片会意外持有大数组。
- 看到 `append` 就完全不考虑容量和扩容成本。

## 十二、工作里的映射

这一章和工作里的联系很直接：

- 构造日志、响应体、导出文本时会遇到字符串拼接问题。
- 批量收集结果、分页组装时会遇到切片扩容问题。
- 缓存和大数据处理时会遇到子切片引用大数组问题。

## 十三、建议练习

你可以自己补两个 benchmark：

```go
func BenchmarkSafeSubset(b *testing.B)
func BenchmarkJoinLargeParts(b *testing.B)
```

重点不是把数字跑得多漂亮，而是学会“带着问题验证优化”。
