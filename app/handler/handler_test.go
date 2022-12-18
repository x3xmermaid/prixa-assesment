package handlers_test

import (
	"net/http/httptest"

	"github.com/gorilla/mux"
)

const (
	// IntentionallyError is the intentional error message for testing.
	intentionallyError  = "Error created intentionally."
	authorizationHeader = "Authorization"
)

type MockedRedisDB struct {
	ErrMap       map[string]bool
	CondMap      map[string]bool
	ErrStatement error
}

// Setup sets up a test HTTP server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
// It is inspired by go-octokit.
func setup(r *mux.Router) *httptest.Server {
	// test server
	server := httptest.NewServer(r)
	return server
}

func (m *MockedRedisDB) Put(key string, value interface{}) error {
	if m.ErrMap["Put"] {
		return m.ErrStatement
	}
	return nil
}

func (m *MockedRedisDB) GetValue(key string) ([]byte, error) {
	return []byte{}, nil
}

func (m *MockedRedisDB) IsAvailable(key string) bool {
	return true
}
