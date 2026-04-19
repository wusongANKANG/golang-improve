package stdlibdemo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func EncodeUser(user User) ([]byte, error) {
	return json.Marshal(user)
}

func DecodeUser(data []byte) (User, error) {
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return User{}, err
	}

	return user, nil
}

func ReadAllUpper(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(string(data)), nil
}

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

func NewHealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})
}
