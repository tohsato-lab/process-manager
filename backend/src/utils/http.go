package utils

import (
	"io/ioutil"
	"net/http"
	"time"
)

func RequestHTTP(method string, endpoint string, timeOut time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: timeOut,
	}
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := res.Body.Close(); err != nil {
		return nil, err
	}
	return body, nil
}

func RespondByte(w http.ResponseWriter, status int, data []byte) {
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusVariantAlsoNegotiates)
		return
	}
}
