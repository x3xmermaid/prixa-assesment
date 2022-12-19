package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"prixa-assesment/app/config"
	handler "prixa-assesment/app/handler"
	"testing"

	"github.com/gorilla/mux"
)

func TestShortenUrl(t *testing.T) {
	redisdb := &MockedRedisDB{
		ErrStatement: fmt.Errorf(intentionallyError),
		ErrMap:       map[string]bool{},
	}

	nconfig := config.ServiceConfig{
		ServiceData: config.ServiceDataConfig{
			Address:     ":8000",
			LocalDomain: "localhost",
		},
	}

	handler := handler.NewHandler(redisdb, &nconfig)

	// Setup Routing
	r := mux.NewRouter()
	r.HandleFunc("/short-url", handler.ShortenUrl).Methods(http.MethodPost)

	// Create httptest Server
	httpServer := setup(r)
	defer httpServer.Close()
	serverURL, _ := url.Parse(httpServer.URL)

	// 	// Hit API Endpoint
	targetPath := fmt.Sprintf("%v%v", serverURL, "/short-url")

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

	t.Run("ShortenUrl Put Error", func(t *testing.T) {
		redisdb.ErrMap["IsAvailable"] = true
		redisdb.ErrMap["Put"] = true
		var jsonRequest = []byte(`{"url":"www.try.com/satria/amanattullah"}`)
		req, _ := http.NewRequest(http.MethodPost, targetPath, bytes.NewBuffer(jsonRequest))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Unable to get worker status: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Response code should be %v, but have %v", http.StatusInternalServerError, resp.StatusCode)
		}

		resp.Body.Close()
	})

	t.Run("ShortenUrl Request Body Error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, targetPath, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Unable to get worker status: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Response code should be %v, but have %v", http.StatusBadRequest, resp.StatusCode)
		}

		resp.Body.Close()
	})
}
