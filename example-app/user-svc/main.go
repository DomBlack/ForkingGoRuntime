package main

import (
	"net/http"

	"github.com/DomBlack/ForkingGoRuntime/example-app/pkg/httpsrv"
)

func main() {
	httpsrv.Start("user-svc", 8091, http.HandlerFunc(defaultHandler))
}

func defaultHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}
