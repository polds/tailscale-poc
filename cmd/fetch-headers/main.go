// Fetch Headers attempts to fetch the headers of a service routed through Tailscale.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var client = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	},
}

func main() {
	log.Print("starting server...")
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		log.Printf("defaulting to port %s", port)
	}

	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	endpoint := os.Getenv("TAILSCALE_ENDPOINT")

	fmt.Fprintf(w, "fetching %s...\n\n", endpoint)
	hdrFetch(ctx, w, endpoint)

	if err := ctx.Err(); err != nil {
		fmt.Fprintf(w, "context error: %v\n", err)
	}
}

func hdrFetch(ctx context.Context, w http.ResponseWriter, addr string) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		fmt.Fprintf(w, "unable to create new request: %v\n", err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "unable to issue http request: %v\n", err)
		return
	}
	defer res.Body.Close()
	var sb strings.Builder
	for k, v := range res.Header {
		fmt.Fprintf(&sb, "%s => %q\n", k, v)
	}
	fmt.Fprintln(w, sb.String())
}
