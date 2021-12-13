package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func PostNewCard(trelloListId string, trelloApiKey string, trelloUserToken string) {
	url := fmt.Sprintf("https://api.trello.com/1/cards?idList=%s&key=%s&token=%s", trelloListId, trelloApiKey, trelloUserToken)
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"name":"To this and that","desc":"bla bla description bla","pos":"top"}`)
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

	fmt.Println("response Status:", resp.Status)
}
