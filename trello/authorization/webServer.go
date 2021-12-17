package authorization

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

const authorizeHtml = `<!DOCTYPE html>
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
</html>`

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

func initializeWebServer(token chan authorizationResult) *http.Server {

	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/static/authorize.html", func(w http.ResponseWriter, req *http.Request) {

		_, err := w.Write([]byte(authorizeHtml))
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
