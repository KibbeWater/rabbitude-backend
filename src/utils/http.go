package utils

import (
	"net/http"
	"strings"
)

func CreatePostRequest(url string, body string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}
