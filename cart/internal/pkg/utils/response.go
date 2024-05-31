package utils

import "net/http"

func SuccessReponse(w http.ResponseWriter) {
	w.Write([]byte("{}"))
}
