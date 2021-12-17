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
	go stopWebServer(server)

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
		ch <- authorizationResult{"", errors.New("the environment variable TRELLO_API_KEY is not set")}
	}

	err := openBrowser(fmt.Sprintf("https://trello.com/1/authorize?expiration=never&callback_method=fragment&return_url=http://localhost:%d/static/authorize.html&name=quick-task-creator&scope=read,write&response_type=fragment&key=%s", port, trelloApiKey))
	if err != nil {
		ch <- authorizationResult{"", err}
	}
}

func initializeWebServer(token chan authorizationResult) *http.Server {

	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/static/authorize.html", func(w http.ResponseWriter, req *http.Request) {

		_, err := w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>Authorize Quick-Task-Creator</title></head>
	<body>
		<script type="text/javascript">
			const hash = window.top.location.hash.substr(1);
			const xhr = new XMLHttpRequest();
			xhr.onreadystatechange = function() {
				if (this.readyState === 4) {
					if(this.status === 200) {
						document.body.innerHTML = "<h1>Authorization successful</h1><p>Have fun using Quick-Task-Creator</p>";
					}
					else {
						document.body.innerHTML = "<h1>Whoops, something went wrong</h1>".concat("<p>", this.response, "</p>");
					}
				}
			}
			xhr.open("POST", "/authorize", true);
			xhr.setRequestHeader('Content-Type', 'application/json');
			xhr.send(hash);
		</script>
	</body>
</html>`))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

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

			token <- authorizationResult{strings.TrimPrefix(tokenWithPrefix, "token="), nil}
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
