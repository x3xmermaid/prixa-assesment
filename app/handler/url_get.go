package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"prixa-assesment/app/model"
	rf "prixa-assesment/app/responseformat"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetUrl is method for get url
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) {
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

	result.TotalRedirect++
	now := time.Now()
	result.Updated_at = &now

	err = h.redisdb.Put(uniqueID, result)
	if err != nil {
		logrus.Errorf("redis Put : %v", err)
		rf.ResponseNOK(http.StatusInternalServerError, statusError, msgInternalServerError, w)
		return
	}

	url := result.URL
	if !strings.Contains(url, "http://") || !strings.Contains(url, "https://") {
		url = fmt.Sprintf("https://%v", url)
	}

	logrus.Infof("short_url:%v, total_redirect: %v, last_redirect:%v", result.ShortUrl, result.TotalRedirect, result.Updated_at)
	http.Redirect(w, r, url, http.StatusFound)
}
