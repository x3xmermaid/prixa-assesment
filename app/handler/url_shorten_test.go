package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	handler "prixa-assesment/app/handler"
	"testing"

	"github.com/gorilla/mux"
)

func TestShortenUrl(t *testing.T) {
	redisdb := &MockedRedisDB{
		ErrStatement: fmt.Errorf(intentionallyError),
		ErrMap:       map[string]bool{},
	}

	handler := handler.NewHandler(redisdb)

	// Setup Routing
	r := mux.NewRouter()
	r.HandleFunc("/short-link", handler.ShortenUrl).Methods(http.MethodPost)

	// Create httptest Server
	httpServer := setup(r)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	// 	// Hit API Endpoint
	targetPath := fmt.Sprintf("%v%v", serverURL, "/short-link")

	// Insert OK
	t.Run("ShortenUrl OK", func(t *testing.T) {
		var jsonRequest = []byte(`{"url":"www.try.com/satria/amanattullah"}`)
		req, _ := http.NewRequest(http.MethodPost, targetPath, bytes.NewBuffer(jsonRequest))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Unable to get worker status: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Response code should be %v, but have %v", http.StatusCreated, resp.StatusCode)
		}

		resp.Body.Close()
	})
}
