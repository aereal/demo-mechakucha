package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aereal/demo-mechakucha/webbase"
)

var (
	UPSTREAM_ORIGIN string
)

func init() {
	UPSTREAM_ORIGIN = os.Getenv("UPSTREAM_ORIGIN")
}

func main() {
	config := &webbase.Config{HostPort: ":8001"}
	if err := webbase.Run(config, handler()); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

type rootResponse struct {
	Upstream string
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		res := &rootResponse{Upstream: UPSTREAM_ORIGIN}
		json.NewEncoder(w).Encode(&res)
	})
	return mux
}
