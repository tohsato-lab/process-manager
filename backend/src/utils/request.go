package utils

import (
	"io"
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
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			return
		}
	}(res.Body)
	return body, nil
}
