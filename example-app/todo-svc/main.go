package main

import (
	"net/http"

	"github.com/DomBlack/ForkingGoRuntime/example-app/pkg/httpsrv"
)

func main() {
	httpsrv.Start("todo-svc", 8090, http.HandlerFunc(defaultHandler))
}

func defaultHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
