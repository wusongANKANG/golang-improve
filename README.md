# Golang Improve

这是一个面向 Go 开发者的系统化学习仓库，目标不是“看完概念”，而是做到：

- 能按模块重新捡回 Go 基础知识。
- 每个模块都有对应代码用例，可以直接运行和断点调试。
- 学完之后能把知识映射回真实工作场景。

## 仓库结构

```text
.
├── docs/                 # 按知识模块整理的中文文档
├── examples/             # 每个模块对应的 Go 示例代码与测试
├── .vscode/launch.json   # VS Code 调试配置
├── go.mod                # 根模块，方便统一运行所有示例
└── Makefile              # 常用命令入口
```

## 推荐学习顺序

1. [学习路线图](docs/00-roadmap/README.md)
2. [开发环境与程序结构](docs/01-setup/README.md)
3. [变量、常量、类型与零值](docs/02-basics/README.md)
4. [数组、切片、Map 与流程控制](docs/03-collections-control/README.md)
5. [函数、闭包、defer、panic、recover](docs/04-functions/README.md)
6. [结构体、方法、接口](docs/05-structs-interfaces/README.md)
7. [包、模块与泛型](docs/06-packages-generics/README.md)
8. [错误处理](docs/07-errors/README.md)
9. [常用标准库：io、json、http、time](docs/08-stdlib/README.md)
10. [并发：goroutine、channel、select](docs/09-concurrency/README.md)
11. [context 与取消控制](docs/10-context/README.md)
12. [测试、基准测试、竞态检测](docs/11-testing/README.md)
13. [性能优化与常见陷阱](docs/12-performance/README.md)

## 快速开始

先跑通整个仓库：

```bash
go test ./...
```

只学某个模块：

```bash
go test -v ./examples/02_basics
go test -v ./examples/03_collections_control
```

只看某个测试：

```bash
go test -run TestSliceSharingDemo -v ./examples/03_collections_control
```

跑 benchmark：

```bash
go test -bench . -benchmem ./examples/11_testing
go test -bench . -benchmem ./examples/12_performance
```

做竞态检测：

```bash
go test -race ./examples/09_concurrency
```

## 用 Makefile 跑

```bash
make test
make module MODULE=./examples/09_concurrency
make run TEST=TestWorkerPool MODULE=./examples/09_concurrency
make bench MODULE=./examples/12_performance
make race MODULE=./examples/09_concurrency
```

## 如何快速 Debug

如果你用 VS Code：

1. 安装官方 Go 插件。
2. 打开任意 `*_test.go` 文件。
3. 使用仓库里的 `.vscode/launch.json`。
4. 选择 “Debug current package tests” 或 “Debug current file”。

推荐的学习方式：

- 先读对应模块文档。
- 再看 `examples/<module>` 里的代码。
- 直接对测试打断点，单步观察变量、切片、接口、channel 的变化。

## 学习建议

- 第一轮先追求“能讲明白”，不要追求背语法。
- 第二轮开始把每个示例改写一遍，比如自己再补一个测试。
- 第三轮结合工作代码，把模块里的知识点映射到实际项目。

如果你想继续扩展这个仓库，后面可以加：

- 数据库与 `database/sql`
- `gin` / `grpc` / `protobuf`
- 项目工程化布局
- 依赖注入、日志、配置、可观测性
- 面试专题与高频陷阱
