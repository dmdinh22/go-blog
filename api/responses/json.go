package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(w httpResponseWriter, statusCode int, data interface()) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Fprint(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error)  {
	if err != nil {
		JSON(w, statusCode, struct {
				Error strint `json:"error"`
		}){
				Error: err.Error(),
		}

		return
	}

	JSON(w, httpStatusBadRequest, nil)
}