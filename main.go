package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	// Simple map of channel_id => redirect URL
	streams map[int]string
	key     string
)

func whepEndpoint(w http.ResponseWriter, r *http.Request) {
	strChannelID := strings.TrimPrefix(r.URL.Path, "/v1/whep/endpoint/")

	channelID, err := strconv.Atoi(strChannelID)
	if err != nil {
		errWrongParams(w, r)
		return
	}

	endpoint, ok := streams[channelID]
	if !ok {
		errNotFound(w, r)
		return
	}

	http.Redirect(w, r, endpoint, http.StatusTemporaryRedirect)
}

func whipEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented")
}

func startStream(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")
	if authKey != key {
		fmt.Printf("have %s want %s", authKey, key)
		errUnauthorized(w, r)
		return
	}

	r.ParseForm()
	strChannelId := r.Form.Get("channel_id")
	endpoint := r.Form.Get("endpoint")
	if strChannelId == "" || endpoint == "" {
		errWrongParams(w, r)
		return
	}
	channelID, err := strconv.Atoi(strChannelId)
	if err != nil {
		errWrongParams(w, r)
		return
	}

	streams[channelID] = endpoint
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("Accepted"))
}

func endStream(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")
	if authKey != key {
		errUnauthorized(w, r)
		return
	}

	r.ParseForm()
	strChannelId := r.Form.Get("channel_id")
	if strChannelId == "" {
		fmt.Printf("channel_id=s=%s", strChannelId)
		errWrongParams(w, r)
		return
	}
	channel_id, err := strconv.Atoi(strChannelId)
	if err != nil {
		fmt.Printf("channel_id=d=%d", channel_id)
		errWrongParams(w, r)
		return
	}

	delete(streams, channel_id)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("OK"))
}

func main() {
	flag.StringVar(&key, "key", "", "a key to prove you are a waveguide server and can register streams")

	flag.Parse()

	key := os.Getenv("RTR_KEY")
	if key == "" {
		panic("A RTR_KEY is required to start the RTRouter")
	}
	httpPort := os.Getenv("RTR_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	streams = make(map[int]string)

	http.HandleFunc("/v1/whep/endpoint", whepEndpoint)
	http.HandleFunc("/v1/whip/endpoint", whipEndpoint)
	http.HandleFunc("/v1/state/start_stream", startStream)
	http.HandleFunc("/v1/state/end_stream", endStream)

	log.Fatal(http.ListenAndServe(":"+httpPort, logRequest(http.DefaultServeMux)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func errUnauthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("Unauthorized"))
}
func errWrongParams(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("Invalid Parameters"))
}
func errNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("Not found"))
}
