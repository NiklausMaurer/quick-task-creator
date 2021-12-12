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

	trelloApiKey, present := os.LookupEnv("TRELLO_API_KEY")
	if !present {
		log.Fatalf("TRELLO_API_KEY not set")
	}

	openbrowser(fmt.Sprintf("http://localhost:8080/authorize/?token=%s", trelloApiKey))

	http.HandleFunc("/authorize", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		if !query.Has("token") {
			log.Fatalf("no token given: %s", req.URL)
		}

		_, _ = fmt.Fprintf(w, "Hello %s", query.Get("token"))
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
