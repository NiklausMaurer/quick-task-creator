package main

import (
	"bufio"
	"fmt"
	"github.com/NiklausMaurer/quick-task-creator/secretStore"
	"github.com/NiklausMaurer/quick-task-creator/trello/authorization"
	"github.com/NiklausMaurer/quick-task-creator/trello/client"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
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
			{
				Name:   "configure",
				Usage:  "configure this installation",
				Action: executeConfigureCommand,
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

	config, err := GetConfig()
	if err != nil {
		return err
	}

	return client.PostNewCard(taskName, config.DefaultListId, getTokenResult.Secret)
}

func executeConfigureCommand(*cli.Context) error {

	trelloListId, err := requestInput("What's the list id of the trello list you'd like to add tasks to?")
	if err != nil {
		return err
	}

	config := Config{
		DefaultListId: trelloListId,
	}

	err = SetConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func requestInput(caption string) (string, error) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(caption)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		if len(input) > 0 {
			return strings.TrimSuffix(input, "\n"), nil
		}
	}
}
