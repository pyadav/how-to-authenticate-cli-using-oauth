package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"read:org", "read:user", "read:project", "public_repo", "gist"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:9999/oauth/callback",
	}

	// start server
	ctx := context.Background()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	server := &http.Server{Addr: ":9999"}

	// create a channel to receive the authorization code
	codeChan := make(chan string)

	http.HandleFunc("/oauth/callback", handelOauthCallback(ctx, config, codeChan))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// get the OAuth authorization URL
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above
	fmt.Printf("Your browser has been opened to visit::\n%s\n", url)

	// open user's browser to login page
	if err := browser.OpenURL(url); err != nil {
		panic(fmt.Errorf("failed to open browser for authentication %s", err.Error()))
	}

	// wait for the authorization code to be received
	code := <-codeChan

	// exchange the authorization code for an access token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Failed to exchange authorization code for token: %v", err)
	}

	if !token.Valid() {
		log.Fatalf("Cann't get source information without accessToken: %v", err)
		return
	}

	// write the access token to a file
	if err := writeTokenToFile(token); err != nil {
		log.Fatalf("Failed to write token to file: %v", err)
	}

	// shut down the HTTP server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to shut down server: %v", err)
	}

	log.Println(color.CyanString("Authentication successful"))
}

func handelOauthCallback(ctx context.Context, config *oauth2.Config, codeChan chan string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParts, _ := url.ParseQuery(r.URL.RawQuery)

		// Use the authorization code that is pushed to the redirect URL.
		code := queryParts["code"][0]
		log.Printf("code: %s\n", code)

		// write the authorization code to the channel
		codeChan <- code

		msg := "<p><strong>Authentication successful</strong>. You may now close this tab.</p>"
		// send a success message to the browser
		fmt.Fprint(w, msg)
	}
}

func writeTokenToFile(token *oauth2.Token) error {
	// create file with 0600 permissions
	file, err := os.OpenFile("token.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to create token file: %v", err)
	}
	defer file.Close()

	// encode token as JSON and write to file
	if err := json.NewEncoder(file).Encode(token); err != nil {
		return fmt.Errorf("Unable to write token to file: %v", err)
	}

	return nil
}
