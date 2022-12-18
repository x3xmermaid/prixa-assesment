package handlers

// Handler holds the API endpoint's function handler.
type Handler struct {
	redisdb RedisDB
}

// NewHandler function to make connection database into handler
func NewHandler(redisdb RedisDB) *Handler {
	return &Handler{
		redisdb: redisdb,
	}
}

// RedisDB presents the interface for RedisDB instance.
type RedisDB interface {
}
