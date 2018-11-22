package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aereal/demo-mechakucha/webbase"
)

func main() {
	// TODO: upstream origin
	config := &webbase.Config{HostPort: ":8001"}
	if err := webbase.Run(config, handler()); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, `{"ok":true}`)
	})
	return mux
}
