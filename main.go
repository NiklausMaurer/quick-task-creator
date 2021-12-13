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

			trelloUserToken := GetUserToken()
			trelloListId := "5e42613e71e90d4b76228153"

			client.PostNewCard(taskName, trelloListId, trelloApiKey, trelloUserToken)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

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
