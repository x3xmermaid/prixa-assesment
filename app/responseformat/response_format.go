package responseformat

import (
	"encoding/json"
	"net/http"
)

// ResponseFormat stands for our default response in API
type ResponseFormat struct {
	Data    interface{} `json:"data,omitempty"`
	Status  interface{} `json:"status,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

// ResponseOK is function to return json OK
func (rf *ResponseFormat) ResponseOK(code int, status string, data interface{}, w http.ResponseWriter) {
	// handle response for non json response
	// please provide data with []byte data type
	// i.e: text/plain, toml, e.t.c..
	if w.Header().Get("Content-Type") != "" {
		w.WriteHeader(code)
		w.Write(data.([]byte))
		return
	}

	// default with json response
	rf.Status = status
	rf.Data = data
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(rf)
	if err != nil {
		resErr := ResponseFormat{
			Message: err,
		}
		jsonErr, _ := json.Marshal(resErr)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErr)
		return
	}
	w.WriteHeader(code)
	w.Write(resp)
}

// ResponseNOK is function to return json NOK
func (rf *ResponseFormat) ResponseNOK(code int, status string, errors interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	rf.Status = status
	rf.Message = errors
	resp, err := json.Marshal(rf)
	if err != nil {
		resErr := ResponseFormat{
			Message: err.Error(),
		}
		jsonErr, _ := json.Marshal(resErr)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonErr)
		return
	}

	w.WriteHeader(code)
	w.Write(resp)
}
