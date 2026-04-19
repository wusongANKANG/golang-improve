# 模块 07：错误处理

错误处理是 Go 工程实践里最重要的基本功之一。很多人学 Go 时最先感受到的就是“为什么总是在写 `if err != nil`”，但真正把这件事理解清楚之后，你会发现它其实是在帮你把错误路径显式化。

## 学习目标

- 理解 Go 的错误是普通值，而不是异常机制。
- 掌握哨兵错误、自定义错误、错误包装的使用方式。
- 学会使用 `errors.Is` 和 `errors.As`。
- 建立更接近工程实践的错误处理习惯。

## 一、Go 的错误为什么是值

Go 的 `error` 本质上是一个接口：

```go
type error interface {
	Error() string
}
```

这意味着错误不是某种特殊语法，而是普通值。函数可以创建它、返回它、包装它、比较它、传递它。

这也是为什么最常见的调用风格是：

```go
value, err := someFunc()
if err != nil {
	return err
}
```

## 二、什么是哨兵错误

当前模块定义了两个哨兵错误：

```go
var (
	ErrInvalidAge  = errors.New("invalid age")
	ErrNonPositive = errors.New("value must be positive")
)
```

所谓哨兵错误，可以先简单理解成“有稳定身份的错误值”。它适合表达那些：

- 含义明确。
- 需要被调用方识别。
- 不依赖太多额外上下文的信息。

例如：

- 参数非法
- 记录不存在
- 权限不足

## 三、自定义错误什么时候有用

看当前模块里的类型：

```go
type FieldError struct {
	Field  string
	Reason string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}
```

如果你只返回一个普通字符串错误，比如：

```go
errors.New("name invalid")
```

那调用方很难知道到底是哪个字段错了、为什么错。

自定义错误的优势在于它可以携带结构化信息，比如：

- 哪个字段失败了。
- 哪一步出错了。
- 是否可重试。

## 四、错误包装为什么重要

当前模块里有这样的代码：

```go
func RegisterUser(name string, age int) error {
	if err := ValidateName(name); err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	if err := ValidateAge(age); err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	return nil
}
```

这段代码里最关键的是 `%w`。

它表示：

- 在保留原始错误身份的同时，增加当前层的上下文信息。

也就是说，最终返回的错误既包含：

- 当前语义：`register user`
- 底层原因：字段为空、年龄非法等

这对于日志排查非常重要。

## 五、`errors.Is` 是干什么的

如果错误被包装了，直接比较字符串就不可靠了。Go 提供了 `errors.Is` 来判断一条错误链中是否包含某个目标错误。

例如：

```go
err := RegisterUser("alice", -1)
if errors.Is(err, ErrInvalidAge) {
	// 年龄非法
}
```

这里即使最外层错误已经被包装成 `register user: invalid age`，`errors.Is` 仍然能识别到底层的 `ErrInvalidAge`。

## 六、`errors.As` 是干什么的

如果你想从错误链里拿出某种具体错误类型，可以用 `errors.As`。

例如当前模块测试里：

```go
var fieldErr *FieldError
if errors.As(err, &fieldErr) {
	fmt.Println(fieldErr.Field)
}
```

这表示：

- 如果错误链上存在 `*FieldError`
- 就把它提取出来，供你继续读取字段信息

这比字符串解析强太多，也更稳。

## 七、`ParsePositiveInt` 这个例子在讲什么

看代码：

```go
func ParsePositiveInt(raw string) (int, error) {
	number, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, fmt.Errorf("parse positive int: %w", err)
	}

	if number <= 0 {
		return 0, fmt.Errorf("parse positive int: %w", ErrNonPositive)
	}

	return number, nil
}
```

这个函数展示了两种错误来源：

- 底层转换错误，比如 `"abc"` 不能转整数。
- 业务语义错误，比如 `0` 不是正数。

这在工作里特别常见，因为一个函数往往既可能遇到系统级错误，也可能遇到业务级错误。

## 八、为什么不要直接比较错误字符串

很多初学者会写：

```go
if err.Error() == "invalid age" {
	...
}
```

这非常脆弱，因为：

- 文案容易变。
- 包装后文本会变化。
- 字符串比较无法表达错误链关系。

更推荐：

```go
errors.Is(err, ErrInvalidAge)
errors.As(err, &fieldErr)
```

## 九、当前模块里的配套测试怎么读

测试文件在 [errors_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/07_errors/errors_test.go)。

建议重点看：

- `TestRegisterUserInvalidAge`
- `TestRegisterUserInvalidName`
- `TestParsePositiveIntSyntax`

它们分别展示：

- 如何判断哨兵错误。
- 如何提取自定义错误。
- 如何识别底层标准库错误类型。

## 十、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/07_errors
```

只调试字段错误：

```bash
go test -run TestRegisterUserInvalidName -v ./examples/07_errors
```

调试时建议重点看：

- `RegisterUser` 返回错误时，最外层和最里层分别是什么。
- `errors.Is` 是如何在包装链里找到 `ErrInvalidAge` 的。
- `errors.As` 成功后，`fieldErr.Field` 和 `fieldErr.Reason` 分别是什么。

## 十一、常见误区

- 把所有错误都写成 `fmt.Errorf("xxx")`，丢掉可判断性。
- 直接比较错误字符串。
- 包装错误时忘记 `%w`。
- 业务错误和底层错误不做区分。

## 十二、工作里的映射

这一章直接影响：

- 接口错误码怎么设计。
- 日志和监控怎么定位问题。
- 调用链里错误语义能不能保留下来。
- 业务层是否能针对具体错误做处理。

## 十三、建议练习

可以自己补两个函数：

```go
var ErrEmptyEmail = errors.New("empty email")

func ValidateEmail(email string) error
func RegisterAccount(name, email string, age int) error
```

练习重点是：

- 什么时候用哨兵错误。
- 什么时候用自定义错误。
- 包装之后还能不能被上层识别。
