package authorization

import (
	"errors"
	"fmt"
	"io"
	"log"
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

func PerformAuthorization() string {

	listener := initializeNetworkListener()
	authorizationResultChannel := make(chan authorizationResult)
	server := initializeWebServer(authorizationResultChannel)

	go startWebServer(listener, server)
	go startBrowser(listener.Addr().(*net.TCPAddr).Port)
	go sendTimeOutAfter(time.Duration(120)*time.Second, authorizationResultChannel)

	result := <-authorizationResultChannel
	stopWebServer(server)

	if result.err == nil {
		log.Printf("Aaaand the token iiis...: %s, I'm done here.", result.token)
	} else {
		log.Fatal("An error occurred during the authorization process: ", result.err)
	}

	return result.token
}

func sendTimeOutAfter(d time.Duration, resultChannel chan authorizationResult) {
	time.Sleep(d)
	resultChannel <- authorizationResult{err: errors.New("timeout expired")}
}

func stopWebServer(server *http.Server) {
	err := server.Close()
	if err != nil {
		log.Print("Stopping the web server failed: ", err)
	}
}

func initializeNetworkListener() net.Listener {

	listener, err := net.Listen("tcp", ":42671")

	if err == nil {
		log.Print("Web server initialized, listening to ", listener.Addr())
	} else {
		log.Fatal("Web server initialisation went wrong: ", err)
	}

	return listener
}

func startWebServer(listener net.Listener, server *http.Server) {
	func() {
		err := server.Serve(listener)
		if err != nil {
			log.Fatal("Web server crashed: ", err)
		}
	}()
}

func startBrowser(port int) {
	trelloApiKey, present := os.LookupEnv("TRELLO_API_KEY")
	if !present {
		log.Fatalf("TRELLO_API_KEY not set")
	}

	openBrowser(fmt.Sprintf("https://trello.com/1/authorize?expiration=never&callback_method=fragment&return_url=http://localhost:%d/static/authorize.html&name=quick-task-creator&scope=read,write&response_type=fragment&key=%s", port, trelloApiKey))
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

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}