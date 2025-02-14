package api

import (
	"bearer-issuer/bearer"
	"net/http"
)

func process(w http.ResponseWriter, r *http.Request) {
	requestType := r.Form.Get("requestType")

	if requestType == "getBearer" {

		result, err := bearer.Register()

		WriteOut(result, err, w)

		return
	}
}
