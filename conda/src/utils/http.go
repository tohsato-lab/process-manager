package utils

import (
	"net/http"
)

func RespondByte(w http.ResponseWriter, status int, data []byte) {
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusVariantAlsoNegotiates)
		return
	}
}
