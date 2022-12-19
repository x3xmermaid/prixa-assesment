package handlers

import (
	nconfig "prixa-assesment/app/config"

	"github.com/thedevsaddam/govalidator"
)

const (
	statusSuccess          = "success"
	statusFail             = "fail"
	statusError            = "error"
	msgInternalServerError = "Internal Server Error"
)

// Handler holds the API endpoint's function handler.
type Handler struct {
	config  *nconfig.ServiceConfig
	redisdb RedisDB
}

// NewHandler function to make connection database into handler
func NewHandler(redisdb RedisDB, config *nconfig.ServiceConfig) *Handler {
	return &Handler{
		redisdb: redisdb,
		config:  config,
	}
}

func (h Handler) getRules() govalidator.MapData {
	rules := govalidator.MapData{
		"url": []string{"required", "url"},
	}

	return rules
}

// RedisDB presents the interface for RedisDB instance.
type RedisDB interface {
	Put(key string, value interface{}) error
	GetValue(key string) ([]byte, error)
	IsAvailable(key string) bool
}
