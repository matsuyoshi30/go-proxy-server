package main

import (
	"io"
	"log"
	"net/http"
)

type proxyHandler struct {
	transport *http.Transport
}

func newProxyHandler(ts *http.Transport) *proxyHandler {
	return &proxyHandler{ts}
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Host == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Request.Host:%s, Request.URL: %s\n", r.Host, r.URL)
	cr := new(http.Request)
	*cr = *r
	cr.URL.Scheme = "http"
	cr.URL.Host = r.Host
	resp, err := h.transport.RoundTrip(cr)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// Usage
// Terminal1:
//     go run main.go
// Terminal2:
//     curl -H 'Host: <want to access Host without scheme>' localhost:8081/<want to access endpoint>
//   example)
//      curl -H 'Host: example.com' localhost:8081

func main() {
	ph := newProxyHandler(http.DefaultTransport.(*http.Transport))
	server := &http.Server{
		Addr:    ":8081",
		Handler: ph,
	}
	log.Fatal(server.ListenAndServe())
}
