package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/NiklausMaurer/quick-task-creator/authorization"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	trelloApiKey, present := os.LookupEnv("TRELLO_API_KEY")
	if !present {
		log.Fatalf("TRELLO_API_KEY not set")
	}

	// board: 5de620427aa9f3570c298caf
	// list: 5e42613e71e90d4b76228153

	trelloUserToken := GetUserToken()
	trelloListId := "5e42613e71e90d4b76228153"

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

func GetUserToken() string {

	homeDirPath := os.Getenv("HOME")
	tokenFilePath := fmt.Sprintf("%s/.quick-task-creator/token", homeDirPath)
	_, err := os.Stat(tokenFilePath)

	if errors.Is(err, os.ErrNotExist) {
		token, err := authorization.PerformAuthorization()
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Authorization process failed. Reason: ", err)
			os.Exit(1)
		}

		err = os.MkdirAll(filepath.Dir(tokenFilePath), os.ModePerm)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Unable to create config directory. Reason: ", err)
			os.Exit(1)
		}

		tokenFile, err := os.Create(tokenFilePath)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Unable create access token file. Reason: ", err)
			os.Exit(1)
		}

		_, err = tokenFile.WriteString(token)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Unable to write access token to file. Reason: ", err)
			os.Exit(1)
		}

		err = os.Chmod(tokenFilePath, 0600)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Unable to set access token file permissions. Reason: ", err)
			os.Exit(1)
		}
	}

	tokenContent, err := os.ReadFile(tokenFilePath)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Unable to read access token from file. Reason: ", err)
		os.Exit(1)
	}

	token := string(tokenContent)
	return token
}
