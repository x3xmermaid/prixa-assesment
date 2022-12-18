package handlers_test

import (
	"fmt"
	"net/http"
	"net/url"
	handler "prixa-assesment/app/handler"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetUrl(t *testing.T) {
	redisdb := &MockedRedisDB{
		ErrStatement: fmt.Errorf(intentionallyError),
		ErrMap:       map[string]bool{},
	}

	handler := handler.NewHandler(redisdb, nil)

	// Setup Routing
	r := mux.NewRouter()
	r.HandleFunc("/{url}", handler.GetUrl).Methods(http.MethodGet)

	// Create httptest Server
	httpServer := setup(r)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	// 	// Hit API Endpoint
	targetPath := fmt.Sprintf("%v%v", serverURL, "/h36bKa")

	// Insert OK
	t.Run("GetUrl OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, targetPath, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Unable to get worker status: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Response code should be %v, but have %v", http.StatusOK, resp.StatusCode)
		}

		resp.Body.Close()
	})
}
