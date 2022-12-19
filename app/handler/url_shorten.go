package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	model "prixa-assesment/app/model"
	rf "prixa-assesment/app/responseformat"

	"github.com/google/uuid"
	"github.com/thedevsaddam/govalidator"

	"github.com/sirupsen/logrus"
)

// ShortenUrl is method for create short url
func (h *Handler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	rf := &rf.ResponseFormat{}

	now := time.Now()
	data := model.ShortUrl{
		TotalRedirect: 0,
		CreatedAt:     &now,
		Updated_at:    &now,
	}
	opts := govalidator.Options{
		Request: r,
		Data:    &data,
		Rules:   h.getRules(),
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) != 0 {
		logrus.Errorf("body validation error : %v", e)
		msg := map[string]interface{}{"validationError": e}
		rf.ResponseNOK(http.StatusBadRequest, statusFail, msg, w)
		return
	}

	isKeyUsed := true
	repeat := 0
	for isKeyUsed {

		uniqueID := generateID(data.URL)
		isKeyUsed := h.redisdb.IsAvailable(uniqueID)

		if !isKeyUsed {
			data.ShortUrl = fmt.Sprintf("%v%v/%v", h.config.ServiceData.LocalDomain, h.config.ServiceData.Address, uniqueID)
			err := h.redisdb.Put(uniqueID, data)
			if err != nil {
				logrus.Errorf("redis put : %v", err)
				rf.ResponseNOK(http.StatusInternalServerError, statusError, msgInternalServerError, w)
				return
			}
		}

		if repeat == 3 {
			break
		}
		repeat++
	}

	rf.ResponseOK(http.StatusCreated, statusSuccess, data, w)
}

func generateID(url string) string {
	newID := uuid.New().String()
	uniqueID := base64.StdEncoding.EncodeToString([]byte(newID))
	return string(uniqueID[0:6])
}
