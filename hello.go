package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/authorize", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		if !query.Has("token") {
			log.Fatalf("no token given: %s", req.URL)
		}

		_, _ = fmt.Fprintf(w, "Hello %s", query.Get("token"))
	})

	_ = http.ListenAndServe(":8080", nil)
}
