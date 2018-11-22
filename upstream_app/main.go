package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aereal/demo-mechakucha/webbase"
)

var (
	Up      = true
	ErrDown = fmt.Errorf("server is down")
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

type updateStatus struct {
	Up bool
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleWithStatus(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, `{"ok":true}`)
	}))
	mux.HandleFunc("/users", handleWithStatus(func(w http.ResponseWriter, r *http.Request) {
		users := []*user{
			&user{Name: "Kumiko"},
			&user{Name: "Reina"},
			&user{Name: "Haduki"},
			&user{Name: "Sapphire"},
		}
		usersRes := &listUsersResponse{Users: users}
		webbase.RespondJSON(w, usersRes)
	}))
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var body updateStatus
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				webbase.RespondErrorJSON(w, http.StatusBadRequest, err)
				return
			}
			Up = body.Up
		default:
			status := &updateStatus{Up: Up}
			webbase.RespondJSON(w, status)
		}
	})
	return mux
}

func handleWithStatus(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !Up {
			unavailableHandler(w, r)
		} else {
			f(w, r)
		}
	}
}

func unavailableHandler(w http.ResponseWriter, r *http.Request) {
	webbase.RespondErrorJSON(w, http.StatusServiceUnavailable, ErrDown)
}
