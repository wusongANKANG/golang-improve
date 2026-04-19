# 模块 09：并发：goroutine、channel、select

并发是 Go 最有辨识度的能力之一，但也是最容易“看懂代码、写出 bug”的地方。因为并发问题很多时候不是语法不会，而是对执行顺序、同步边界、生命周期没有建立稳定心智模型。

## 学习目标

- 理解 goroutine 是什么。
- 理解 channel 在并发中的作用。
- 理解 `select` 的基本用法。
- 掌握 `sync.WaitGroup` 这种最常见的同步手段。
- 能通过示例和调试建立最小并发直觉。

## 一、goroutine 是什么

goroutine 可以先简单理解成“由 Go runtime 管理的轻量级并发执行单元”。

最小示例：

```go
go func() {
	fmt.Println("run in goroutine")
}()
```

关键点不是 `go` 关键字本身，而是你要意识到：

- 这段逻辑会异步执行。
- 它和当前函数可能并行推进。
- 如果主流程提前结束，这个 goroutine 可能还没来得及完成。

## 二、为什么不能只会开 goroutine

很多并发 bug 就出在“开了 goroutine，但没有管理它”。

你必须明确：

- 谁负责等待它结束。
- 它的结果往哪里传。
- 它什么时候该停止。

所以 Go 并发通常不是只写一个 `go`，而是要配合：

- `sync.WaitGroup`
- `channel`
- `context`

## 三、`WaitGroup` 的作用是什么

当前模块里的 `SquareAll`：

```go
func SquareAll(nums []int) []int {
	results := make([]int, len(nums))

	var wg sync.WaitGroup
	for index, value := range nums {
		wg.Add(1)
		go func(index, value int) {
			defer wg.Done()
			results[index] = value * value
		}(index, value)
	}

	wg.Wait()
	return results
}
```

这里的 `WaitGroup` 在做一件事：

- 等待所有子 goroutine 都执行完，再返回结果。

如果没有 `wg.Wait()`，主流程可能在某些 goroutine 还没写完结果前就提前返回。

## 四、为什么 `SquareAll` 没有数据竞争

这个问题很重要。

虽然多个 goroutine 同时在写 `results`，但它们写的是不同索引位置：

```go
results[index] = value * value
```

只要每个 goroutine 写入的地址不同，并且主 goroutine 在 `wg.Wait()` 之后才读取整体结果，这个例子就是安全的。

这里还能顺手学到一个经典细节：

```go
go func(index, value int) { ... }(index, value)
```

把循环变量作为参数传进去，是为了避免闭包直接捕获循环变量带来的混乱。

## 五、channel 是什么

channel 是 Go 里 goroutine 之间通信的通道。你可以先把它理解成一种带同步语义的管道。

最小例子：

```go
ch := make(chan int)
go func() {
	ch <- 1
}()

value := <-ch
```

它的核心思想是：

- 不直接共享数据，而是通过通信传递数据。

当然这句话不是绝对规则，但它代表了 Go 并发的典型风格。

## 六、`select` 是怎么工作的

当前模块有个例子：

```go
func RaceMessages(delays map[string]time.Duration) string {
	results := make(chan string, len(delays))

	for label, delay := range delays {
		go func(label string, delay time.Duration) {
			time.Sleep(delay)
			results <- label
		}(label, delay)
	}

	return <-results
}
```

这个函数虽然没有显式写 `select`，但它展示的是“谁先完成，就先拿谁的结果”的思路。

如果写成显式 `select`，常见样子会是：

```go
select {
case v := <-ch1:
	return v
case v := <-ch2:
	return v
}
```

`select` 的作用就是在多个 channel 操作之间做选择。

## 七、worker pool 为什么值得学

当前模块的 `WorkerPool` 是一个非常典型的并发模型：

```go
func WorkerPool(workers int, inputs []int) []int
```

它体现了几件实际项目里很常见的事：

- 任务输入和结果输出通过 channel 传递。
- 有固定数量的 worker 并发消费任务。
- 主流程需要正确关闭 `jobs` 和 `results`。

这是任务消费、批量处理、异步作业里很常见的结构。

## 八、`WorkerPool` 里最值得看的是什么

首先是 worker 启动：

```go
for i := 0; i < workers; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for job := range jobs {
			results <- result{
				index: job.index,
				value: job.value * job.value,
			}
		}
	}()
}
```

这里 `for job := range jobs` 的意思是：

- 只要 `jobs` channel 没关闭，worker 就持续接收任务。
- 一旦 `jobs` 被关闭并且数据读完，循环结束，worker 退出。

然后是结果关闭时机：

```go
close(jobs)
wg.Wait()
close(results)
```

这个顺序很重要，因为：

- 先关闭 `jobs`，告诉 worker 没有新任务了。
- 再等待所有 worker 处理完。
- 最后才能安全关闭 `results`。

## 九、怎么避免 goroutine 泄漏

goroutine 泄漏通常意味着：某个 goroutine 一直在等，但永远等不到结束条件。

常见原因：

- channel 没人接收。
- channel 永远不关闭。
- `select` 没有取消路径。
- 上游提前退出，下游还在阻塞。

这一章先建立意识，下一章的 `context` 会继续解决“如何取消”的问题。

## 十、配套测试怎么读

测试文件在 [concurrency_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/09_concurrency/concurrency_test.go)。

建议重点看：

- `TestSquareAll`
- `TestRaceMessages`
- `TestWorkerPool`

它们分别对应：

- 并发计算但结果保持原顺序。
- 谁快拿谁。
- 固定 worker 数量的任务处理模式。

## 十一、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/09_concurrency
```

调试 worker pool：

```bash
go test -run TestWorkerPool -v ./examples/09_concurrency
```

调试时建议重点看：

- `wg.Add`、`wg.Done`、`wg.Wait` 的调用顺序。
- `jobs` 和 `results` 分别什么时候关闭。
- `output[result.index] = result.value` 为什么能恢复输入顺序。

## 十二、常见误区

- 认为开 goroutine 就等于性能更高。
- 不等 goroutine 完成就返回。
- 循环里直接捕获变量，导致值错乱。
- 错误关闭 channel，或者重复关闭 channel。
- 忽略 goroutine 的退出条件。

## 十三、工作里的映射

这一章在工作里会频繁出现：

- 并发调用多个下游接口。
- 多 worker 消费消息队列。
- 批量处理文件、任务、请求。
- 限定并发度，避免把服务打爆。

## 十四、建议练习

你可以自己补两个函数：

```go
func SumAllConcurrent(nums []int) int
func FirstDone(left, right time.Duration) string
```

一个练 goroutine + 同步，一个练“谁先完成就返回”的选择逻辑。
