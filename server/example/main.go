package main

import (
	"fmt"
	"github.com/akula410/web/server"
	"net/http"
)

func main() {
	var srv = &server.Server{Network:"tcp", Address:":8081"}
	srv.
		HandleFunc(":8080", "/", func (w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprint(w, "Hello World")
			if err != nil {
				panic(err)
			}
		}).
		Start()

	srv.Block()
}
