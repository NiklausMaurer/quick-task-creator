package main

import (
	"errors"
	"fmt"
	"github.com/NiklausMaurer/quick-task-creator/secretStore"
	"github.com/NiklausMaurer/quick-task-creator/trello/authorization"
	"github.com/NiklausMaurer/quick-task-creator/trello/client"
	"github.com/urfave/cli/v2"
	"log"
	"os"
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

			command := c.Args().First()
			if command == "authorize" {
				token, err := authorization.PerformAuthorization()
				if err != nil {
					return err
				}

				err = secretStore.StoreSecret("token", token)
				if err != nil {
					return err
				}

				return nil
			}

			taskName := c.String("t")
			if len(taskName) > 0 {
				return addCardToDefaultList(taskName)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func addCardToDefaultList(taskName string) error {
	getTrelloApiKeyResult := secretStore.GetSecret("trello-api-key")
	if getTrelloApiKeyResult.Error != nil {
		return getTrelloApiKeyResult.Error
	}
	if !getTrelloApiKeyResult.SecretFound {
		return errors.New("trello api key not set")
	}

	getTokenResult := secretStore.GetSecret("token")
	if getTokenResult.Error != nil {
		return getTokenResult.Error
	}

	if !getTokenResult.SecretFound {
		fmt.Println("quick-task-creator is not authorized yet. Try $quick-task-creator authorize to fix this.")
		return nil
	}

	trelloListId := "5e42613e71e90d4b76228153"

	return client.PostNewCard(taskName, trelloListId, getTrelloApiKeyResult.Secret, getTokenResult.Secret)
}
