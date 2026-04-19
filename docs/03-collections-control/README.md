# 模块 03：数组、切片、Map 与流程控制

这一章是 Go 日常编码里使用频率最高的一章。尤其是切片和 map，几乎每个项目都会反复出现。如果这部分理解不扎实，后面很多 bug 会看起来“很玄学”，实际上只是底层行为没有想清楚。

## 学习目标

- 掌握数组、切片、map 的区别和用法。
- 理解切片和底层数组共享内存的本质。
- 熟悉 `if`、`for`、`switch`、`range` 的惯用写法。
- 能通过调试观察切片共享和复制行为。

## 一、数组和切片的核心区别

数组是固定长度的，长度属于类型的一部分：

```go
var a [3]int
var b [4]int
```

这里 `[3]int` 和 `[4]int` 是两种不同类型。

切片不是数组本身，而是对数组一段区间的描述：

```go
nums := []int{1, 2, 3}
```

切片更灵活，所以 Go 业务代码里大多数时候你看到的都是切片，而不是数组。

## 二、切片到底是什么

可以先把切片想象成一个小结构，它大致记录了三件事：

- 指向底层数组的指针。
- 当前长度 `len`。
- 当前容量 `cap`。

这也是为什么切片经常会出现“我改了 A，B 也跟着变”的情况。因为两个切片可能并没有复制数据，只是指向了同一块底层数组。

## 三、先看一个最简单的 `range`

当前模块里的求和例子：

```go
func Sum(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}

	return total
}
```

这里可以顺手记住 `range` 的两个特点：

- 第一个返回值通常是索引。
- 第二个返回值通常是元素副本。

如果你不需要索引，可以用 `_` 忽略。

## 四、map 的最常见使用场景

词频统计是理解 map 很好的例子：

```go
func WordFrequency(words []string) map[string]int {
	result := make(map[string]int, len(words))
	for _, word := range words {
		normalized := strings.ToLower(strings.TrimSpace(word))
		if normalized == "" {
			continue
		}

		result[normalized]++
	}

	return result
}
```

这个例子值得学习的点很多：

- `make(map[string]int, len(words))` 是初始化 map 的常见方式。
- `result[normalized]++` 之所以成立，是因为不存在的 key 会先得到值类型零值，这里就是 `0`。
- 对字符串先做标准化，是业务处理中非常常见的预处理思路。

## 五、`switch` 在 Go 里经常比一连串 `if` 更清晰

例如评分逻辑：

```go
func Grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	default:
		return "D"
	}
}
```

这里的 `switch` 没有表达式，等价于一组按顺序判断的条件分支。这是 Go 里非常常见的写法，尤其适合区间判断和状态分流。

## 六、这一章最关键的知识点：切片共享底层数组

看当前模块最重要的示例：

```go
func SliceSharingDemo() (base []int, shared []int, safeCopy []int) {
	base = make([]int, 3, 4)
	copy(base, []int{1, 2, 3})

	shared = base[:2]
	shared = append(shared, 9)

	safeCopy = append([]int(nil), base...)
	safeCopy[0] = 100

	return base, shared, safeCopy
}
```

这段代码是理解切片最值得反复调试的一段。

先按步骤拆开：

第一步：

```go
base = make([]int, 3, 4)
copy(base, []int{1, 2, 3})
```

此时：

- `base = [1 2 3]`
- `len(base) = 3`
- `cap(base) = 4`

第二步：

```go
shared = base[:2]
```

此时：

- `shared = [1 2]`
- `shared` 和 `base` 指向同一个底层数组

第三步：

```go
shared = append(shared, 9)
```

因为 `base` 还有剩余容量，所以这次 `append` 很可能直接复用原底层数组，于是 `base` 也会变成：

```go
base == []int{1, 2, 9}
```

这就是很多人工作里第一次遇到时会很困惑的地方：我明明改的是 `shared`，为什么 `base` 也变了？

答案就是：它们共享底层数组。

## 七、如何安全复制切片

示例里用了这句：

```go
safeCopy = append([]int(nil), base...)
```

它的作用是新建一份数据，而不是继续共用原底层数组。

这在工作里很常见，比如：

- 你不想让下游修改上游传来的切片。
- 你只想截取一部分数据，但不想长期引用整个大数组。
- 你要把数据缓存起来，避免后续原数据被覆盖。

## 八、配套测试怎么读

测试文件在 [collections_control_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/03_collections_control/collections_control_test.go)。

最关键的是这个测试：

```go
func TestSliceSharingDemo(t *testing.T) {
	base, shared, safeCopy := SliceSharingDemo()

	if base[2] != 9 {
		t.Fatalf("base[2] = %d, want 9", base[2])
	}
}
```

这个测试不是在考你记忆，而是在告诉你一个行为事实：`shared` 的 `append` 的确影响了 `base`。

## 九、怎么运行和怎么 debug

运行整个模块：

```bash
go test -v ./examples/03_collections_control
```

只跑切片示例：

```bash
go test -run TestSliceSharingDemo -v ./examples/03_collections_control
```

调试时建议观察：

- `base` 在 `append` 前后的内容变化。
- `shared` 的长度和容量。
- `safeCopy[0] = 100` 为什么不会反过来影响 `base[0]`。

## 十、常见误区

- 把切片理解成“动态数组值拷贝”。
- 忽略 `append` 可能复用原底层数组。
- 忘记 `range` 拿到的是值副本。
- 读取 map 时不考虑 key 不存在的情况。

## 十一、工作里的映射

这一章几乎每天都会遇到：

- 收集数据库结果时用切片。
- 做聚合统计时用 map。
- 做状态分类时用 `switch`。
- 切片在多个函数间传递时，必须注意是否共享数据。

## 十二、建议练习

可以自己再写两个函数：

```go
func Top2(nums []int) []int
func MergeCount(left, right map[string]int) map[string]int
```

重点不是功能本身，而是刻意练习切片复制和 map 累加。
