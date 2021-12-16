package userTokenStore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type GetUserTokenResult struct {
	TokenFound bool
	Token      string
	Error      error
}

func GetUserToken() GetUserTokenResult {
	tokenFilePath := getTokenFilePath()

	fileExists, err := fileExists(tokenFilePath)
	if err != nil {
		return GetUserTokenResult{false, "", err}
	}

	if !fileExists {
		return GetUserTokenResult{false, "", nil}
	}

	tokenContent, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return GetUserTokenResult{true, "", err}
	}

	token := string(tokenContent)
	return GetUserTokenResult{true, token, nil}
}

func getTokenFilePath() string {
	homeDirPath := os.Getenv("HOME")
	tokenFilePath := fmt.Sprintf("%s/.quick-task-creator/token", homeDirPath)
	return tokenFilePath
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

func StoreUserToken(token string) error {

	tokenFilePath := getTokenFilePath()

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
