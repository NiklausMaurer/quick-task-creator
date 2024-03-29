package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiKey = "e6b342d7e5d3c98eb4cd2b14c6d7f599"

type RequestError struct {
	StatusCode int
}

type trelloCard struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"desc,omitempty"`
	Position    string `json:"pos,omitempty"`
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("Request failed with status %d", e.StatusCode)
}

func PostNewCard(taskName string, taskDescription string, trelloListId string, trelloUserToken string, trelloApiUrl string) error {

	url := fmt.Sprintf("%s/1/cards?idList=%s&key=%s&token=%s", trelloApiUrl, trelloListId, apiKey, trelloUserToken)

	client := http.Client{Timeout: 2 * time.Second}

	var card = trelloCard{
		Name:        taskName,
		Description: taskDescription,
		Position:    "top",
	}

	var jsonStr, err = json.Marshal(card)
	if err != nil {
		return err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return &RequestError{StatusCode: resp.StatusCode}
	}

	return nil
}
