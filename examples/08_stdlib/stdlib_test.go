package stdlibdemo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncodeDecodeUser(t *testing.T) {
	data, err := EncodeUser(User{Name: "alice", Age: 18})
	if err != nil {
		t.Fatalf("EncodeUser() error = %v", err)
	}

	user, err := DecodeUser(data)
	if err != nil {
		t.Fatalf("DecodeUser() error = %v", err)
	}

	if user.Name != "alice" || user.Age != 18 {
		t.Fatalf("decoded user = %+v", user)
	}
}

func TestReadAllUpper(t *testing.T) {
	got, err := ReadAllUpper(strings.NewReader("go"))
	if err != nil {
		t.Fatalf("ReadAllUpper() error = %v", err)
	}

	if got != "GO" {
		t.Fatalf("ReadAllUpper() = %q, want %q", got, "GO")
	}
}

func TestBuildQueryURL(t *testing.T) {
	got, err := BuildQueryURL("https://example.com/search", map[string]string{
		"q":    "golang",
		"page": "1",
	})
	if err != nil {
		t.Fatalf("BuildQueryURL() error = %v", err)
	}

	if got != "https://example.com/search?page=1&q=golang" {
		t.Fatalf("BuildQueryURL() = %q", got)
	}
}

func TestNewHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	NewHealthHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}

	if body := recorder.Body.String(); !strings.Contains(body, "\"status\":\"ok\"") {
		t.Fatalf("response body = %q", body)
	}
}
