package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type CallbackServer struct {
	server    *http.Server
	addr      string
	codeChan  chan string
	stateChan chan string
	errChan   chan error
}

func NewCallbackServer(port string) *CallbackServer {
	return &CallbackServer{
		// I'm binding explicitly to 127.0.0.1 since spotify does not accept http://localhost URIs
		addr:      fmt.Sprintf("127.0.0.1:%s", port),
		codeChan:  make(chan string, 1),
		stateChan: make(chan string, 1),
		errChan:   make(chan error, 1),
	}
}

func (s *CallbackServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.handleCallback)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errChan <- err
		}
	}()

	return nil
}

// Stop stops the callback server
func (s *CallbackServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// WaitForCallback waits for the OAuth callback and returns the code and state
func (s *CallbackServer) WaitForCallback(timeout time.Duration) (string, string, error) {
	select {
	case code := <-s.codeChan:
		state := <-s.stateChan
		return code, state, nil
	case err := <-s.errChan:
		return "", "", err
	case <-time.After(timeout):
		return "", "", fmt.Errorf("callback timeout after %v", timeout)
	}
}

// handleCallback handles the OAuth callback
func (s *CallbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	error := r.URL.Query().Get("error")

	if error != "" {
		http.Error(w, fmt.Sprintf("OAuth error: %s", error), http.StatusBadRequest)
		s.errChan <- fmt.Errorf("OAuth error: %s", error)
		return
	}

	if code == "" {
		http.Error(w, "No authorization code received", http.StatusBadRequest)
		s.errChan <- fmt.Errorf("no authorization code received")
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// TODO: add some styling
	fmt.Fprintf(w, `
		<html>
		<head><title>Spotify CLI Authorization</title></head>
		<body>
			<h1>Authorization Successful!</h1>
			<p>You can now close this window and return to the terminal.</p>
		</body>
		</html>
	`)

	// Send the code and state to the channel
	select {
	case s.codeChan <- code:
	case <-time.After(1 * time.Second):
	}

	select {
	case s.stateChan <- state:
	case <-time.After(1 * time.Second):
	}
}
