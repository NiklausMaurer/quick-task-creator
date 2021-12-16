package authorization

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type authorizationResult struct {
	token string
	err   error
}

func successfulAuthorizationResult(token string) authorizationResult {
	return authorizationResult{token: token}
}

func PerformAuthorization() (string, error) {

	listener, err := initializeNetworkListener()
	if err != nil {
		return "", err
	}

	ch := make(chan authorizationResult)

	server := initializeWebServer(ch)

	go startWebServer(listener, server, ch)
	go startBrowser(listener.Addr().(*net.TCPAddr).Port, ch)
	go sendTimeOutAfter(time.Duration(120)*time.Second, ch)

	result := <-ch
	stopWebServer(server)

	return result.token, result.err
}

func sendTimeOutAfter(d time.Duration, resultChannel chan authorizationResult) {
	time.Sleep(d)
	resultChannel <- authorizationResult{err: errors.New("timeout expired")}
}

func stopWebServer(server *http.Server) {
	_ = server.Close()
}

func initializeNetworkListener() (net.Listener, error) {

	listener, err := net.Listen("tcp", ":42671")

	if err != nil {
		return nil, err
	}

	return listener, nil
}

func startWebServer(listener net.Listener, server *http.Server, ch chan<- authorizationResult) {
	err := server.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		ch <- authorizationResult{"", err}
	}
}

func startBrowser(port int, ch chan<- authorizationResult) {
	trelloApiKey, present := os.LookupEnv("TRELLO_API_KEY")
	if !present {
		ch <- authorizationResult{"", errors.New("The environment variable TRELLO_API_KEY is not set")}
	}

	err := openBrowser(fmt.Sprintf("https://trello.com/1/authorize?expiration=never&callback_method=fragment&return_url=http://localhost:%d/static/authorize.html&name=quick-task-creator&scope=read,write&response_type=fragment&key=%s", port, trelloApiKey))
	if err != nil {
		ch <- authorizationResult{"", err}
	}
}

func initializeWebServer(token chan authorizationResult) *http.Server {

	serverMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	serverMux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	serverMux.HandleFunc("/authorize", func(w http.ResponseWriter, req *http.Request) {

		if req.Method == http.MethodPost {
			data, err := io.ReadAll(req.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			tokenWithPrefix := string(data)

			if !strings.HasPrefix(tokenWithPrefix, "token=") {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			token <- successfulAuthorizationResult(strings.TrimPrefix(tokenWithPrefix, "token="))
		}
	})

	server := http.Server{Handler: serverMux}

	return &server
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
