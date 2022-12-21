package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func HttpRequest(method string, url string, body map[string]interface{}, queryString map[string]interface{}, headers map[string]interface{}) ([]byte, error) {
	var requestBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer(jsonBody)
	}
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, err
	}
	if queryString != nil {
		q := req.URL.Query()
		for k, v := range queryString {
			q.Add(k, v.(string))
		}
		req.URL.RawQuery = q.Encode()
	}
	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
