package authorization

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"
)

type authorizationResult struct {
	token string
	err   error
}

func PerformAuthorization(apiKey string) (string, error) {

	listener, err := initializeNetworkListener()
	if err != nil {
		return "", err
	}

	ch := make(chan authorizationResult)

	server := initializeWebServer(ch)

	go startWebServer(listener, server, ch)
	go startBrowser(apiKey, listener.Addr().(*net.TCPAddr).Port, ch)
	go sendTimeOutAfter(time.Duration(120)*time.Second, ch)

	result := <-ch
	go stopWebServer(server)

	return result.token, result.err
}

func sendTimeOutAfter(d time.Duration, resultChannel chan authorizationResult) {
	time.Sleep(d)
	resultChannel <- authorizationResult{err: errors.New("timeout expired")}
}

func startBrowser(apiKey string, port int, ch chan<- authorizationResult) {
	err := openBrowser(fmt.Sprintf("https://trello.com/1/authorize?expiration=never&callback_method=fragment&return_url=http://localhost:%d/static/authorize.html&name=quick-task-creator&scope=read,write&response_type=fragment&key=%s", port, apiKey))
	if err != nil {
		ch <- authorizationResult{"", err}
	}
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
