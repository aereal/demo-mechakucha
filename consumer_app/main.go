package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aereal/demo-mechakucha/webbase"
)

var (
	UPSTREAM_ORIGIN string
)

func init() {
	UPSTREAM_ORIGIN = os.Getenv("UPSTREAM_ORIGIN")
	if UPSTREAM_ORIGIN == "" {
		panic(fmt.Errorf("UPSTREAM_ORIGIN is empty"))
	}
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

type favoritesResponse struct {
	FavoriteUsers []*user
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		res := &rootResponse{Upstream: UPSTREAM_ORIGIN}
		json.NewEncoder(w).Encode(&res)
	})
	mux.HandleFunc("/favorites", func(w http.ResponseWriter, r *http.Request) {
		users, err := getUsers()
		if err != nil {
			webbase.RespondErrorJSON(w, http.StatusServiceUnavailable, fmt.Errorf("Failed to request: %s", err))
			return
		}
		res := &favoritesResponse{
			FavoriteUsers: users.Users,
		}
		webbase.RespondJSON(w, res)
	})
	return mux
}

type user struct {
	Name string
}

type usersRes struct {
	Users []*user
}

func getUsers() (*usersRes, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/users", UPSTREAM_ORIGIN))
	if err != nil {
		return nil, fmt.Errorf("Failed to request: %s", err)
	}
	var body usersRes
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %s", err)
	}
	return &body, nil
}
