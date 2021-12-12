package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	_ "strings"
)

func main() {

	tokenChannel := make(chan string)
	go startWebserver(tokenChannel)

	trelloApiKey, present := os.LookupEnv("TRELLO_API_KEY")
	if !present {
		log.Fatalf("TRELLO_API_KEY not set")
	}

	openbrowser(fmt.Sprintf("https://trello.com/1/authorize?expiration=never&callback_method=fragment&return_url=http://localhost:8080/static/authorize.html&name=quick-task-creator&scope=read,write&response_type=fragment&key=%s", trelloApiKey))

	token := <-tokenChannel

	log.Printf("Aaaand the token iiis...: %s, I'm done here.", token)
}

func startWebserver(token chan<- string) {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/authorize", func(w http.ResponseWriter, req *http.Request) {

		if req.Method == http.MethodPost {
			data, err := io.ReadAll(req.Body)
			if err != nil {
				return
			}

			tokenWithPrefix := string(data)

			if strings.HasPrefix(tokenWithPrefix, "token=") {
				token <- strings.TrimPrefix(tokenWithPrefix, "token=")
			}
		}

	})

	_ = http.ListenAndServe(":8080", nil)

}

func openbrowser(url string) {
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
