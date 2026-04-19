# 模块 10：context 与取消控制

如果说并发章节解决的是“怎么同时做很多事”，那 `context` 解决的就是“这些事什么时候该停”。它是 Go 服务开发里非常核心的一套约定，HTTP、RPC、数据库、任务处理都会大量使用。

## 学习目标

- 理解 `context.Context` 的职责。
- 掌握超时和取消的基本用法。
- 知道怎样传递少量请求级元数据。
- 理解 `context` 的常见误用。

## 一、`context` 到底是干什么的

可以先记一个最实用的定义：

`context` 主要用来传递三类信息：

- 取消信号。
- 截止时间或超时信息。
- 少量请求级数据。

它不是拿来放业务对象、数据库连接、大配置结构体的。

## 二、最常见的创建方式

根上下文：

```go
ctx := context.Background()
```

带取消：

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

带超时：

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
```

这里的 `cancel()` 很重要，哪怕你觉得函数快结束了，也最好保持手动释放的习惯。

## 三、超时是怎么起作用的

看当前模块的例子：

```go
func WorkWithTimeout(ctx context.Context, work time.Duration) error {
	select {
	case <-time.After(work):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

这里的核心逻辑是：

- 如果工作先做完，返回 `nil`。
- 如果上下文先超时或被取消，走 `ctx.Done()` 分支并返回 `ctx.Err()`。

这是 Go 里非常典型的超时控制模板。

## 四、`ctx.Done()` 和 `ctx.Err()` 分别是什么

这两个经常成对出现：

- `ctx.Done()`：一个只读 channel，用来通知“上下文结束了”。
- `ctx.Err()`：结束原因，通常是 `context.Canceled` 或 `context.DeadlineExceeded`。

所以常见模式就是：

```go
select {
case <-ctx.Done():
	return ctx.Err()
}
```

## 五、如何让 goroutine 响应取消

当前模块还有这个例子：

```go
func StreamNumbers(ctx context.Context, limit int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for i := 0; i < limit; i++ {
			select {
			case <-ctx.Done():
				return
			case output <- i:
			}
		}
	}()

	return output
}
```

这个例子非常适合理解 `context` 和 goroutine 的配合：

- 主流程可以从 `output` 读取数据。
- 一旦 `ctx` 被取消，内部 goroutine 会立刻退出。
- `output` 会被关闭，外部 `range` 也能正常结束。

这就是优雅停止的一种最小模型。

## 六、为什么请求级数据可以放进 context

当前模块里还定义了请求 ID：

```go
type requestIDKey string

const requestIDContextKey requestIDKey = "request-id"
```

设置：

```go
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}
```

读取：

```go
func RequestID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(requestIDContextKey).(string)
	return value, ok
}
```

这是一种很常见的模式，用来在调用链上传递：

- request id
- trace id
- user id

但前提是数据要小、要和请求范围绑定。

## 七、为什么自定义 key 类型很重要

很多人会直接写：

```go
context.WithValue(ctx, "request-id", "req-001")
```

虽然也能用，但不推荐，因为字符串 key 容易和别的包冲突。

更推荐像当前模块这样定义一个私有 key 类型：

```go
type requestIDKey string
```

这样更安全。

## 八、`context` 不该做什么

这是工作里很容易踩坑的地方。

不要把这些东西放进 `context`：

- 数据库连接。
- logger 实例的大对象。
- 业务参数集合。
- 可选配置大结构体。

因为 `context` 的设计目标不是“万能参数袋”，而是请求链路控制。

## 九、配套测试怎么读

测试文件在 [context_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/10_context/context_test.go)。

建议重点看：

- `TestWorkWithTimeoutSuccess`
- `TestWorkWithTimeoutDeadlineExceeded`
- `TestStreamNumbers`
- `TestRequestID`

这四个测试基本把 `context` 的核心用法串起来了。

## 十、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/10_context
```

只调试超时：

```bash
go test -run TestWorkWithTimeoutDeadlineExceeded -v ./examples/10_context
```

调试时建议重点看：

- `ctx.Done()` 是在什么时候被关闭的。
- `ctx.Err()` 在超时场景下的具体值。
- `StreamNumbers` 在上下文结束时是如何退出 goroutine 的。

## 十一、常见误区

- 创建了 `WithTimeout` 却忘记 `cancel()`。
- 把 `context` 存进结构体长期持有。
- 用 `context` 传一堆业务参数。
- goroutine 不监听 `ctx.Done()`，导致无法取消。

## 十二、工作里的映射

这一章和真实项目几乎直接对应：

- HTTP 请求超时。
- RPC 调用链取消。
- 数据库查询截止时间。
- 请求链路上的 trace / request id 传递。

## 十三、建议练习

你可以自己补两个函数：

```go
func WaitOrCancel(ctx context.Context, d time.Duration) string
func StreamWithRequestID(ctx context.Context, limit int) <-chan string
```

一个练超时取消，一个练 `context` 数据和 goroutine 配合。
