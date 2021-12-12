package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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

	fmt.Println(token)
}

func startWebserver(token chan<- string) {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/authorize", func(w http.ResponseWriter, req *http.Request) {

		log.Println(req.URL)

		query := req.URL.Query()
		if query.Has("token") {
			token <- query.Get("token")
			_, _ = fmt.Fprintf(w, "Success!")
			return
		}

		_, _ = fmt.Fprintf(w, "Fail!")
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
