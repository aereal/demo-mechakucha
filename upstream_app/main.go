package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aereal/demo-mechakucha/webbase"
)

func main() {
	config := &webbase.Config{HostPort: ":8002"}
	if err := webbase.Run(config, handler()); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

type user struct {
	Name string
}

type listUsersResponse struct {
	Users []*user
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, `{"ok":true}`)
	})
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users := []*user{
			&user{Name: "Kumiko"},
			&user{Name: "Reina"},
			&user{Name: "Haduki"},
			&user{Name: "Sapphire"},
		}
		usersRes := &listUsersResponse{Users: users}
		webbase.RespondJSON(w, usersRes)
	})
	return mux
}
