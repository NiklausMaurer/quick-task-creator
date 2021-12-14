package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type RequestError struct {
	StatusCode int
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("Request failed with status %d", e.StatusCode)
}

func PostNewCard(taskName string, trelloListId string, trelloApiKey string, trelloUserToken string) error {
	url := fmt.Sprintf("https://api.trello.com/1/cards?idList=%s&key=%s&token=%s", trelloListId, trelloApiKey, trelloUserToken)

	var jsonStr = []byte(fmt.Sprintf(`{"name":"%s","desc":"","pos":"top"}`, taskName))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode >= 400 {
		return &RequestError{StatusCode: resp.StatusCode}
	}

	return nil
}
