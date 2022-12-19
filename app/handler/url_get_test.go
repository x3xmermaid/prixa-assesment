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

	t.Run("GetUrl Redis Put Error", func(t *testing.T) {
		redisdb.ErrMap["Put"] = true
		req, _ := http.NewRequest(http.MethodGet, targetPath, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Unable to get worker status: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Response code should be %v, but have %v", http.StatusInternalServerError, resp.StatusCode)
		}

		resp.Body.Close()
	})

	t.Run("GetUrl Redis GetValue-ErrorValue Error", func(t *testing.T) {
		redisdb.ErrMap["GetValue-ErrorValue"] = true
		req, _ := http.NewRequest(http.MethodGet, targetPath, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Unable to get worker status: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Response code should be %v, but have %v", http.StatusInternalServerError, resp.StatusCode)
		}

		resp.Body.Close()
	})

	t.Run("GetUrl Redis GetValue Error", func(t *testing.T) {
		redisdb.ErrMap["GetValue"] = true
		req, _ := http.NewRequest(http.MethodGet, targetPath, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Unable to get worker status: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Response code should be %v, but have %v", http.StatusInternalServerError, resp.StatusCode)
		}

		resp.Body.Close()
	})

}
