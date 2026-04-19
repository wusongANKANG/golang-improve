# 模块 06：包、模块与泛型

这一章的重点是代码组织能力。写 Go 不是只会写一个文件里的函数，而是要知道怎么拆包、怎么做模块管理、什么时候适合用泛型减少重复。

## 学习目标

- 理解包和模块的职责区别。
- 理解导出标识符与非导出标识符。
- 掌握跨包调用的基本方式。
- 掌握泛型最常见的使用场景和边界。

## 一、包是代码组织单位

Go 里一个目录通常对应一个包。例如这个仓库里：

- `examples/06_packages_generics` 是一个包。
- `examples/06_packages_generics/wordutil` 是另一个包。

包的意义是把职责相近的代码收在一起，并且形成边界。

## 二、模块是依赖管理单位

根目录的 [go.mod](/home/wusong/workspace/future-2026/golang-improve/go.mod) 定义了当前模块：

```go
module golang-improve

go 1.25.0
```

这意味着当你在代码里写：

```go
import "golang-improve/examples/06_packages_generics/wordutil"
```

Go 能知道这是当前模块里的另一个包。

## 三、什么叫导出标识符

Go 通过首字母大小写控制可见性：

- 首字母大写：导出，可被其他包访问。
- 首字母小写：不导出，只在当前包内可见。

例如：

```go
func Sum[T Number](items []T) T
func unique[T comparable](items []T) []T
```

上面如果第二个函数以小写开头，那它就只能在当前包内部使用。

这种规则很简单，但非常常用。

## 四、跨包调用怎么发生

当前模块示例：

```go
import "golang-improve/examples/06_packages_generics/wordutil"

func NormalizeWords(words []string) []string {
	return wordutil.CleanLower(words)
}
```

这里体现了最基本的跨包调用流程：

1. 先 `import` 目标包。
2. 调用目标包中导出的标识符。
3. 通过包前缀明确来源。

例如：

```go
result := wordutil.FirstNonEmpty("", "Go")
```

## 五、为什么拆包很重要

当一个目录里的职责越来越多时，如果不拆包，后面会出现：

- 文件越来越臃肿。
- 依赖关系混乱。
- 测试边界不清晰。

当前示例里把字符串清洗逻辑放到 `wordutil`，就是在体现一种简单但实用的拆包方式：把可复用、职责独立的文本处理逻辑收进小工具包。

## 六、泛型解决的是什么问题

泛型最适合解决“逻辑相同，只是类型不同”的重复代码问题。

看这个例子：

```go
type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](items []T) T {
	var total T
	for _, item := range items {
		total += item
	}

	return total
}
```

如果没有泛型，你可能要分别写：

```go
func SumInt(items []int) int
func SumInt64(items []int64) int64
func SumFloat64(items []float64) float64
```

逻辑完全一样，只是元素类型不同。这个时候泛型就很合适。

## 七、类型约束怎么理解

`Number` 这个接口不是运行时接口，而是类型参数约束。

```go
type Number interface {
	~int | ~int64 | ~float64
}
```

它表达的是：

- 类型 `T` 的底层类型必须是这些数值类型之一。

这样 `Sum[T Number]` 才能安全地执行 `total += item`。

这里的 `~` 表示不仅支持精确类型，也支持底层类型是这些类型的自定义类型。

## 八、`comparable` 约束是干什么的

看另一个例子：

```go
func Unique[T comparable](items []T) []T {
	seen := make(map[T]struct{}, len(items))
	result := make([]T, 0, len(items))

	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}
```

这里要求 `T comparable`，因为 map 的 key 必须是可比较类型。这个约束本质上是在告诉编译器：只有能做 `==` 比较的类型，才能拿来去重。

## 九、泛型不是越多越好

这是很重要的边界意识。

泛型适合：

- 容器工具函数。
- 去重、过滤、映射、求和这类通用算法。
- 数据结构库。

泛型不太适合：

- 业务含义很重、已经很具体的逻辑。
- 为了“显得高级”而做的过度抽象。

如果一个函数只服务于一个具体业务模型，通常直接写具体类型反而更清晰。

## 十、当前模块里的几个示例怎么用

求和：

```go
sum1 := Sum([]int{1, 2, 3})
sum2 := Sum([]float64{1.5, 2.5})
```

去重：

```go
words := Unique([]string{"go", "go", "gopher"})
```

跨包文本清理：

```go
normalized := NormalizeWords([]string{" Go ", "", "Gopher"})
first := FirstKeyword("", "  ", "Go")
```

## 十一、配套测试怎么读

测试文件在 [packages_generics_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/06_packages_generics/packages_generics_test.go)。

建议重点看：

- `TestSum`
- `TestUnique`
- `TestNormalizeWords`

这三个测试分别对应泛型数值计算、泛型容器处理、跨包函数调用。

## 十二、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/06_packages_generics
```

调试时建议重点看：

- `Sum` 在 `[]int` 和 `[]float64` 两种输入下是如何工作的。
- `Unique` 的 `seen` map 是怎样去重的。
- `NormalizeWords` 是怎样调用 [wordutil.go](/home/wusong/workspace/future-2026/golang-improve/examples/06_packages_generics/wordutil/wordutil.go) 里的逻辑的。

## 十三、常见误区

- 把包和模块混为一谈。
- 过早拆出很多包，导致结构碎片化。
- 所有重复代码都想用泛型统一。
- 不理解 `comparable` 和类型约束的边界。

## 十四、工作里的映射

这一章会直接影响：

- 项目目录如何拆分。
- 公共工具函数如何复用。
- 是否需要引入泛型减少重复代码。
- 是否能把“工具抽象”和“业务抽象”区分开。

## 十五、建议练习

可以自己补两个函数：

```go
func LastNonEmpty(words ...string) string
func Contains[T comparable](items []T, target T) bool
```

一个练跨包工具思路，一个练泛型约束和容器遍历。
