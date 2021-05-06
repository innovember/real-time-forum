package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	State    bool        `json:"state"`
	HTTPCode int         `json:"httpCode"`
	Message  interface{} `json:"message"`
	Data     interface{} `json:"data"`
}

var (
	err    error
	output []byte
)

func Respond(w http.ResponseWriter, responseStatus bool, httpCode int, message, data interface{}) {
	output, err = json.Marshal(Response{responseStatus, httpCode, message, data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		output = []byte(`{"state":false, "httpCode":500,"message":"failed to marshal JSON in response.JSON()","data":null}`)
		w.Write(output)
		return
	}
	w.WriteHeader(httpCode)
	w.Write(output)
}

func JSON(w http.ResponseWriter, responseStatus bool, httpCode int, message string, data interface{}) {
	Respond(w, responseStatus, httpCode, message, data)
}
