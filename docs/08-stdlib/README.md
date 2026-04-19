# 模块 08：常用标准库：io、json、http、time

Go 标准库是这门语言非常大的优势。很多在其他语言里需要依赖第三方库的事情，在 Go 里标准库已经做得很扎实了。这一章的重点不是把所有标准库都讲完，而是先抓住工作里最常用的几类能力。

## 学习目标

- 理解 `io.Reader` 这类小接口为什么如此重要。
- 掌握 JSON 编解码的基本用法。
- 学会写最小可测试的 HTTP handler。
- 理解标准库如何鼓励“面向接口编程”。

## 一、为什么先学 `io`

因为 `io.Reader` 和 `io.Writer` 是 Go 标准库里最基础、最通用的一组接口。

最经典的是：

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

这意味着只要一个东西实现了 `Read`，它就可以被当作数据输入源：

- 文件
- 字符串
- 网络连接
- 内存缓冲区

所以很多函数根本不关心数据来自哪里，它只关心“你是不是一个 Reader”。

## 二、看一个面向接口的例子

当前模块里的示例：

```go
func ReadAllUpper(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(string(data)), nil
}
```

这段代码很值得学习，因为它完全没有绑定具体输入类型。

它既可以这样用：

```go
ReadAllUpper(strings.NewReader("go"))
```

也可以未来这样用：

```go
ReadAllUpper(file)
ReadAllUpper(response.Body)
```

这就是接口抽象的价值。

## 三、JSON 编解码的基本写法

当前模块的结构体：

```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
```

这里的 `` `json:"name"` `` 是结构体 tag，用来指定 JSON 字段名。

编码：

```go
data, err := json.Marshal(User{Name: "alice", Age: 18})
```

解码：

```go
var user User
err := json.Unmarshal(data, &user)
```

当前模块把它们封装成了：

```go
func EncodeUser(user User) ([]byte, error)
func DecodeUser(data []byte) (User, error)
```

## 四、为什么 JSON 结构体 tag 很重要

如果没有 tag，Go 默认会按字段名生成 JSON 键名，比如：

- `Name`
- `Age`

但实际接口里往往更常见的是：

- `name`
- `age`

所以你会经常看到：

```go
type Payload struct {
	RequestID string `json:"request_id"`
}
```

## 五、为什么 URL 不要手拼

当前模块里的另一个例子：

```go
func BuildQueryURL(base string, params map[string]string) (string, error) {
	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	for key, value := range params {
		query.Set(key, value)
	}

	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
```

很多人会手写：

```go
base + "?q=" + keyword + "&page=" + page
```

这很容易出错，尤其遇到：

- 参数编码
- 特殊字符
- 原始 URL 已经带查询参数

标准库的 `net/url` 会更可靠。

## 六、HTTP handler 为什么也能很好测

当前模块里有这个例子：

```go
func NewHealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})
}
```

这就是一个非常典型的 Go HTTP handler。

它的好处是：

- 逻辑很集中。
- 依赖接口而不是具体服务器。
- 可以配合 `httptest` 完全脱离网络测试。

## 七、`httptest` 为什么值得掌握

看测试中的写法：

```go
request := httptest.NewRequest(http.MethodGet, "/health", nil)
recorder := httptest.NewRecorder()

NewHealthHandler().ServeHTTP(recorder, request)
```

这意味着你不需要真的启动端口、不需要浏览器、不需要 curl，就能验证 handler 行为。

这在工作里非常重要，因为：

- 测试更快。
- 更稳定。
- 更容易覆盖边界条件。

## 八、`time` 虽然这里没单独写示例，但你会不断遇到

这一章标题里带 `time`，是因为在 Go 里它和 HTTP、并发、context 都关系很深。

你后面很常见的代码会是：

```go
time.Now()
time.Sleep(10 * time.Millisecond)
time.After(time.Second)
```

这些能力会在并发和 `context` 模块里继续出现。

## 九、当前模块里的几个示例怎么用

JSON 编解码：

```go
data, _ := EncodeUser(User{Name: "alice", Age: 18})
user, _ := DecodeUser(data)
```

Reader 读取：

```go
upper, _ := ReadAllUpper(strings.NewReader("go"))
```

URL 拼装：

```go
link, _ := BuildQueryURL("https://example.com/search", map[string]string{
	"q": "golang",
})
```

HTTP 测试：

```go
handler := NewHealthHandler()
```

## 十、配套测试怎么读

测试文件在 [stdlib_test.go](/home/wusong/workspace/future-2026/golang-improve/examples/08_stdlib/stdlib_test.go)。

建议重点看：

- `TestEncodeDecodeUser`
- `TestReadAllUpper`
- `TestNewHealthHandler`

这三个测试分别覆盖结构体与 JSON 的映射、接口输入抽象、HTTP handler 的无网络测试。

## 十一、怎么运行和怎么 debug

运行模块：

```bash
go test -v ./examples/08_stdlib
```

只调试 HTTP：

```bash
go test -run TestNewHealthHandler -v ./examples/08_stdlib
```

调试时建议重点看：

- `json.Marshal` 之后输出的数据长什么样。
- `ReadAllUpper` 里的 `reader` 具体类型是什么。
- `ServeHTTP` 执行后，`recorder.Code`、`recorder.Body`、`recorder.Header()` 分别是什么。

## 十二、常见误区

- 把具体类型绑死，不用接口抽象输入输出。
- 手工拼接 URL。
- handler 写得太大，导致难以测试。
- JSON tag 不统一，接口字段命名混乱。

## 十三、工作里的映射

这一章和真实项目的联系非常直接：

- 读文件、读网络响应、读请求体时会反复遇到 `io.Reader`。
- 接口开发里每天都在处理 JSON。
- 服务端和客户端都离不开 `net/http`。
- 测试 HTTP 行为时 `httptest` 很高频。

## 十四、建议练习

你可以自己补两个函数：

```go
func DecodeUsers(data []byte) ([]User, error)
func NewEchoHandler() http.Handler
```

一个练 JSON 数组，一个练 HTTP 请求回显，都是非常贴近工作的小练习。
