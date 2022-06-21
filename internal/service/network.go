package service

import (
	"bytes"
	"io"
	"net/http"
)

func SendPostRequest(url string, buffer *bytes.Buffer) (error error, status int,
	response string) {
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return err, 0, ""
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, 0, ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, 0, ""
	}
	return nil, resp.StatusCode, string(body)
}
