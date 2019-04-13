package main

import (
	"net/http"

	"github.com/jiripetrlik/handy-dns/internal/app/rest"
)

func main() {
	rest.HandleRestAPI()
	http.ListenAndServe(":8080", nil)
}
