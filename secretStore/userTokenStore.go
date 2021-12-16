package secretStore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type GetSecretResult struct {
	SecretFound bool
	Secret      string
	Error       error
}

func GetSecret(key string) GetSecretResult {
	secretFilePath := getSecretFilePath(key)

	fileExists, err := fileExists(secretFilePath)
	if err != nil {
		return GetSecretResult{false, "", err}
	}

	if !fileExists {
		return GetSecretResult{false, "", nil}
	}

	secretContent, err := os.ReadFile(secretFilePath)
	if err != nil {
		return GetSecretResult{true, "", err}
	}

	secret := string(secretContent)

	return GetSecretResult{true, secret, nil}
}

func getSecretFilePath(key string) string {
	homeDirPath := os.Getenv("HOME")
	secretFilePath := fmt.Sprintf("%s/.quick-task-creator/%s", homeDirPath, key)
	return secretFilePath
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

func StoreSecret(key string, secret string) error {

	secretFilePath := getSecretFilePath(key)

	err := os.MkdirAll(filepath.Dir(secretFilePath), os.ModePerm)
	if err != nil {
		return err
	}

	secretFile, err := os.Create(secretFilePath)
	if err != nil {
		return err
	}

	_, err = secretFile.WriteString(secret)
	if err != nil {
		return err
	}

	err = os.Chmod(secretFilePath, 0600)
	if err != nil {
		return err
	}

	return nil
}
