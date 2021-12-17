package main

import (
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
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a new card to the default list",
				Action:  executeAddCommand,
			},
			{
				Name:   "authorize",
				Usage:  "authorize this installation with trello",
				Action: executeAuthorizeCommand,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func executeAuthorizeCommand(*cli.Context) error {
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

func executeAddCommand(c *cli.Context) error {

	taskName := c.Args().First()

	if len(taskName) == 0 {
		return nil
	}

	getTokenResult := secretStore.GetSecret("token")
	if getTokenResult.Error != nil {
		return getTokenResult.Error
	}

	if !getTokenResult.SecretFound {
		fmt.Println("quick-task-creator is not authorized yet. Try $quick-task-creator authorize to fix this.")
		return nil
	}

	homeDirPath := os.Getenv("HOME")
	configFilePath := fmt.Sprintf("%s/.quick-task-creator/%s", homeDirPath, "config.json")
	config, err := GetConfig(configFilePath)
	if err != nil {
		return err
	}

	return client.PostNewCard(taskName, config.DefaultListId, getTokenResult.Secret)
}
