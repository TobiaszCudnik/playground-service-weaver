package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
)

// Reverser component.
type Reverser interface {
	Reverse(context.Context, string) (string, error)
}

// Implementation of the Reverser component.
type reverser struct {
	weaver.Implements[Reverser]
}

func (r *reverser) Reverse(_ context.Context, s string) (string, error) {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes), nil
}

func main() {
	// Get a network listener on address "localhost:12345".
	root := weaver.Init(context.Background())
	opts := weaver.ListenerOptions{LocalAddress: "localhost:12345"}
	lis, err := root.Listener("hello", opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("hello listener available on %v\n", lis)

	// Get a client to the Reverser component.
	reverser, err := weaver.Get[Reverser](root)
	if err != nil {
		log.Fatal(err)
	}

	// Serve the /hello endpoint.
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		reversed, err := reverser.Reverse(r.Context(), r.URL.Query().Get("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "Hello, %s!\n", reversed)
	})
	http.Serve(lis, nil)
}
