package main

import (
	"fmt"
	"github.com/NiklausMaurer/quick-task-creator/trello/authorization"
	"github.com/NiklausMaurer/quick-task-creator/trello/client"
	"github.com/NiklausMaurer/quick-task-creator/userTokenStore"
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

				err = userTokenStore.StoreUserToken(token)
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
	trelloApiKey, present := os.LookupEnv("TRELLO_API_KEY")
	if !present {
		log.Fatalf("TRELLO_API_KEY not set")
	}

	getUserTokenResult, err := userTokenStore.GetUserToken()
	if err != nil {
		return err
	}

	if !getUserTokenResult.TokenFound {
		fmt.Println("quick-task-creator is not authorized yet. Try $quick-task-creator authorize to fix this.")
		return nil
	}

	trelloListId := "5e42613e71e90d4b76228153"

	return client.PostNewCard(taskName, trelloListId, trelloApiKey, getUserTokenResult.Token)
}
