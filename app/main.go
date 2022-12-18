package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	config "prixa-assesment/app/config"
	nredis "prixa-assesment/app/db/redis"
	handler "prixa-assesment/app/handler"
)

const (
	// ConfigFileLocation is the file configuration of ths service.
	ConfigFileLocation = "conf/config.yaml"
)

// Handler hold the function handler for API's endpoint.
type Handler interface {
	GetUrl(w http.ResponseWriter, r *http.Request)
	ShortenUrl(w http.ResponseWriter, r *http.Request)
	// GetShortUrlStatus(w http.ResponseWriter, r *http.Request)
}

// NewRouter returns router.
func NewRouter(handler Handler) *mux.Router {
	r := mux.NewRouter()

	// Linux Server Inventory
	r.HandleFunc("/{url}", handler.GetUrl).Methods(http.MethodGet)
	// r.HandleFunc("/{url}/status", handler.GetShortUrlStatus).Methods(http.MethodGet)
	r.HandleFunc("/short-url", handler.ShortenUrl).Methods(http.MethodPost)

	return r
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	if len(arr) >= 2 {
		return fmt.Sprintf("%s/%s", arr[len(arr)-2], arr[len(arr)-1])
	}

	return arr[len(arr)-1]
}

func main() {
	customFormatter := &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf(" %s:%d", formatFilePath(f.File), f.Line)
		},
	}

	logrus.SetFormatter(customFormatter)
	logrus.SetReportCaller(true)
	// Pre-printed text at startup.
	logrus.Printf("prixa Assesment")
	logrus.Println("Start service...")

	// Get Config
	configLoader := config.NewYamlConfigLoader(ConfigFileLocation)
	config, err := configLoader.GetServiceConfig()
	if err != nil {
		logrus.Fatalf("Unable to load configuration: %v", err)
	}

	// redis connection
	redis, err := nredis.NewRedis(config.SourceData.RedisNetwork, config.SourceData.RedisAddress, config.SourceData.RedisPassword,
		config.SourceData.RedisTimeout, config.SourceData.RedisKeyExpireDuration)
	if err != nil {
		logrus.Fatalf("Unable to create RedisDB instance: %v", err)
	}
	logrus.Printf("Redis Connection Test: PASS")

	// initialize handler
	handler := handler.NewHandler(redis, config)
	r := NewRouter(handler)

	// Run Web Server
	logrus.Printf("Starting http server at %v", config.ServiceData.Address)
	err = http.ListenAndServe(config.ServiceData.Address, r)
	if err != nil {
		logrus.Fatalf("Unable to run http server: %v", err)
	}
	logrus.Println("Stopping API Service...")
}
