package client

import (
	"bytes"
	"fmt"
	"net/http"
)

const apiKey = "e6b342d7e5d3c98eb4cd2b14c6d7f599"

type RequestError struct {
	StatusCode int
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("Request failed with status %d", e.StatusCode)
}

func PostNewCard(taskName string, trelloListId string, trelloUserToken string, trelloApiUrl string) error {

	url := fmt.Sprintf("%s/1/cards?idList=%s&key=%s&token=%s", trelloApiUrl, trelloListId, apiKey, trelloUserToken)

	var jsonStr = []byte(fmt.Sprintf(`{"name":"%s","desc":"","pos":"top"}`, taskName))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return &RequestError{StatusCode: resp.StatusCode}
	}

	return nil
}
