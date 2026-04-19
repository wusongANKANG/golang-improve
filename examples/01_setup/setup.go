package setup

import "fmt"

func Greeting(name string) string {
	if name == "" {
		name = "gopher"
	}

	return fmt.Sprintf("hello, %s", name)
}

func ProgramShape() []string {
	return []string{
		"package main",
		"import (...)",
		"func main() { ... }",
	}
}
