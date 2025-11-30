package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient() (*http.Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %v", err)
	}

	configDir := filepath.Join(home, ".config", "personalcli")
	credsPath := filepath.Join(configDir, "credentials.json")
	tokenPath := filepath.Join(configDir, "token.json")

	b, err := os.ReadFile(credsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file at %s: %v", credsPath, err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		log.Println("No token found, starting web authentication flow...")
		tok, err = getTokenFromWeb(config, tokenPath)
		if err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), tok), nil
}

// getTokenFromWeb starts a local server to handle the OAuth2 callback.
func getTokenFromWeb(config *oauth2.Config, tokenPath string) (*oauth2.Token, error) {
	// Create a channel to receive the token
	tokenChan := make(chan *oauth2.Token)
	errChan := make(chan error)

	// State token to prevent CSRF
	state := "random-state-string"
	config.RedirectURL = "http://localhost:8080"

	// Define the callback handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "State token does not match", http.StatusBadRequest)
			errChan <- fmt.Errorf("state token mismatch")
			return
		}

		code := r.URL.Query().Get("code")
		tok, err := config.Exchange(context.TODO(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			errChan <- fmt.Errorf("unable to retrieve token from web: %v", err)
			return
		}

		fmt.Fprintln(w, "Authentication successful! You can close this tab.")
		tokenChan <- tok
	})

	// Start the server in a goroutine
	server := &http.Server{Addr: ":8080"}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- fmt.Errorf("could not start server: %v", err)
		}
	}()

	// Open the browser for authentication
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Println("Your browser should open for authentication automatically.")
	fmt.Printf("If it doesn't, please visit this link: %s\n", authURL)
	err := browser.OpenURL(authURL)
	if err != nil {
		log.Printf("Could not open browser: %v. Please open the URL manually.", err)
	}

	// Wait for the token or an error
	select {
	case tok := <-tokenChan:
		saveToken(tokenPath, tok)
		// Shutdown the server
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
		return tok, nil
	case err := <-errChan:
		// Shutdown the server
		if errS := server.Shutdown(context.Background()); errS != nil {
			log.Printf("Error shutting down server: %v", errS)
		}
		return nil, err
	}
}

// tokenFromFile retrieves a Token from a given file path.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
