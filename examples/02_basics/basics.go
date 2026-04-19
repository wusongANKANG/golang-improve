package basics

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusActive  Status = "active"
)

type Profile struct {
	Name   string
	Age    int
	Active bool
	Status Status
}

func ZeroValueProfile() Profile {
	return Profile{}
}

func NewProfile(name string, age int) Profile {
	return Profile{
		Name:   name,
		Age:    age,
		Active: true,
		Status: StatusActive,
	}
}

func Swap(a, b int) (int, int) {
	return b, a
}

func TypedAndUntypedConstants() (int64, float64) {
	const maxRetries int64 = 3
	const pi = 3.14

	return maxRetries, pi
}
