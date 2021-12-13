package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func PostNewCard(taskName string, trelloListId string, trelloApiKey string, trelloUserToken string) {
	url := fmt.Sprintf("https://api.trello.com/1/cards?idList=%s&key=%s&token=%s", trelloListId, trelloApiKey, trelloUserToken)

	var jsonStr = []byte(fmt.Sprintf(`{"name":"%s","desc":"","pos":"top"}`, taskName))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
}
