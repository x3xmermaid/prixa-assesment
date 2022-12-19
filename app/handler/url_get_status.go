package handlers

import (
	"encoding/json"
	"net/http"

	"prixa-assesment/app/model"
	rf "prixa-assesment/app/responseformat"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetShortUrlStatus is method for get short url status
func (h *Handler) GetShortUrlStatus(w http.ResponseWriter, r *http.Request) {
	rf := &rf.ResponseFormat{}
	vars := mux.Vars(r)
	uniqueID := vars["url"]

	data, err := h.redisdb.GetValue(uniqueID)
	if err != nil {
		logrus.Errorf("redis GetValue : %v", err)
		rf.ResponseNOK(http.StatusInternalServerError, statusError, msgInternalServerError, w)
		return
	}

	result := model.ShortUrl{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		logrus.Errorf("unmarshall : %v", err)
		rf.ResponseNOK(http.StatusInternalServerError, statusError, msgInternalServerError, w)
		return
	}

	rf.ResponseOK(http.StatusOK, statusSuccess, result, w)
}
