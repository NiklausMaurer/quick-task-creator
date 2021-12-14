package main

import (
	"errors"
	"fmt"
	"github.com/NiklausMaurer/quick-task-creator/trello/authorization"
	"github.com/NiklausMaurer/quick-task-creator/trello/client"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Required: true,
				Name:     "t",
				Value:    "",
				Usage:    "Trello task name",
			},
		},
		Action: func(c *cli.Context) error {

			taskName := c.String("t")
			if len(taskName) == 0 {
				return nil
			}

			trelloApiKey, present := os.LookupEnv("TRELLO_API_KEY")
			if !present {
				log.Fatalf("TRELLO_API_KEY not set")
			}

			trelloUserToken, err := GetUserToken()
			if err != nil {
				return err
			}

			trelloListId := "5e42613e71e90d4b76228153"

			return client.PostNewCard(taskName, trelloListId, trelloApiKey, trelloUserToken)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func fileExists(filePath string) (bool, error) {

	_, err := os.Stat(filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func GetUserToken() (string, error) {

	homeDirPath := os.Getenv("HOME")
	tokenFilePath := fmt.Sprintf("%s/.quick-task-creator/token", homeDirPath)

	fileExists, err := fileExists(tokenFilePath)
	if err != nil {
		return "", err
	}

	if !fileExists {
		token, err := authorization.PerformAuthorization()
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Authorization process failed. Reason: ", err)
			os.Exit(1)
		}

		err = writeTokenToFile(token, tokenFilePath)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "There was an issue while saving the user token: ", err)
			os.Exit(1)
		}
	}

	tokenContent, err := os.ReadFile(tokenFilePath)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Unable to read access token from file. Reason: ", err)
		os.Exit(1)
	}

	token := string(tokenContent)
	return token, nil
}

func writeTokenToFile(token string, tokenFilePath string) error {
	err := os.MkdirAll(filepath.Dir(tokenFilePath), os.ModePerm)
	if err != nil {
		return err
	}

	tokenFile, err := os.Create(tokenFilePath)
	if err != nil {
		return err
	}

	_, err = tokenFile.WriteString(token)
	if err != nil {
		return err
	}

	err = os.Chmod(tokenFilePath, 0600)
	if err != nil {
		return err
	}

	return nil
}
