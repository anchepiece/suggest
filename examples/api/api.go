package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/anchepiece/suggest"
)

var (
	matches = []string{"foo", "bar", "baz"}
)

type matchResponse struct {
	Changed bool   `json:"changed"`
	Match   string `json:"match,omitempty"`
	Error   string `json:"error,omitempty"`
}

type matchHandler struct {
	ctx context.Context
}

func (mh *matchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	mr := &matchResponse{}
	query := r.URL.Query().Get("q")
	if query == "" {
		mr.Error = "Must supply query parameter 'q' in URL."
		response, err := json.Marshal(mr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	suggester := suggest.Suggest{}

	suggester.Commands = matches
	match, err := suggester.Autocorrect(query)
	if err != nil {
		mr.Error = err.Error()
		response, err := json.Marshal(mr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

	if match == "" {
		mr.Error = "No match was found"
		response, err := json.Marshal(mr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}

	mr.Changed = match != query
	mr.Match = match
	response, err := json.Marshal(mr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func newMatchHandler(ctx context.Context) *matchHandler {
	return &matchHandler{
		ctx: ctx,
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/match", newMatchHandler(context.Background()))
	log.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
