package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pflege-de/go-force/force"
)

const (
	getTokenPath string = "services/oauth2/token" //nolint:gosec
	// (Wrongly recognized - G101: Potential hardcoded credentials)
	authorizePath string = "services/oauth2/authorize"
)

// postAuthorizationCodeResponse is the response we get after posting the initial authentication code to salesforce
type postAuthorizationCodeReponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Signature    string `json:"signature"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
	InstanceURL  string `json:"instance_url"`
	ID           string `json:"id"`
	TokenType    string `json:"token_type"`
	IssuedAt     string `json:"issued_at"`
}

// OAuthHandler implements the http.Handler interface and shares the token with our code
type OAuthHandler struct {
	shouldCloseServerChannel chan bool // we will send a single value to the channel after our request to shutdown the server
	token                    postAuthorizationCodeReponse
}

// ServeHTTP waits for the user of this application to finish logging in with Salesforce on the site printed by promptUser()
// and accepts the request from the redirected user. The authentication code is sent as a query parameter.
func (s *OAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, ok := r.URL.Query()["code"]
	if !ok || len(code[0]) < 1 {
		fmt.Println("url parameter code was not sent. This is probably the second request automatically sent by the browser")
		// For some reason the browser send two requests to our server
		// one with code as query parameter, one without any query params
		// this is only a real error if no authorizationCode was parsed previously
		return
	}
	client := &http.Client{}
	// send the authentication code to salesforce
	tokenResponse, err := getToken(client, code[0])
	if err != nil {
		fmt.Println("Error sending authentication token to salesforce: %w", err)
		return
	}
	s.token = tokenResponse
	fmt.Fprintf(w, "Got code: %s\nGot Token: %s\nGot Refresh Token: %s\n", code[0], tokenResponse.AccessToken, tokenResponse.RefreshToken)
	s.shouldCloseServerChannel <- true
}

// promptUsers prints the user a link to our salesforce instance. The user must authenticate this app with salesforce, and is then redirected to redirectURL
// redirectURL must be set when creating a new salesforce connected app. For development edit /etc/hosts with an entry redirecting to 127.0.0.1 for redirectURL
// later on, we can use localhorst.dev.p4e.io which redirects to 127.0.0.1 in AWS Route53.
func promptUser() {
	u := url.URL{
		Scheme: "https",
		Host:   os.Getenv("SF_OAUTH_HOST"),
		Path:   authorizePath,
	}
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", os.Getenv("SF_OAUTH_CLIENT_ID"))
	q.Set("redirect_uri", os.Getenv("SF_OAUTH_REDIRECT_URL"))
	u.RawQuery = q.Encode()

	fmt.Printf("\033[32mOpen this link in your browser to authenticate with salesforce OAuth\n\n\033[0m%s\n\n", u.String())
}

// getToken exchanges the authorizationCode for an access and refresh token
func getToken(client *http.Client, authorizationCode string) (postAuthorizationCodeReponse, error) {
	tokenResponse := postAuthorizationCodeReponse{}
	u := url.URL{
		Scheme: "https",
		Host:   os.Getenv("SF_OAUTH_HOST"),
		Path:   getTokenPath,
	}

	q := u.Query()
	q.Set("grant_type", "authorization_code")
	q.Set("code", authorizationCode)
	q.Set("client_id", os.Getenv("SF_OAUTH_CLIENT_ID"))
	q.Set("client_secret", os.Getenv("SF_OAUTH_CLIENT_SECRET"))
	q.Set("redirect_uri", os.Getenv("SF_OAUTH_REDIRECT_URL"))
	u.RawQuery = q.Encode()

	res, err := client.Post(u.String(), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return tokenResponse, fmt.Errorf("error sending the request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return tokenResponse, fmt.Errorf("error reading the body: %w", err)
	}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return tokenResponse, fmt.Errorf("error unmarshaling the json: %w", err)
	}
	return tokenResponse, nil
}

func NewOAuthForce() (force.ForceApiInterface, error) {
	// prompt the user by printing the url to a salesforce login page
	// launch goroutine accepting the redirect
	// block until token is returned
	// quit server after response was received
	// create a force with access token
	promptUser()

	stopChannel := make(chan bool)
	authHandler := OAuthHandler{shouldCloseServerChannel: stopChannel}
	srv := http.Server{Addr: ":443", Handler: &authHandler, ReadHeaderTimeout: 30 * time.Second}
	ctx := context.Background()

	go func(ctx context.Context, srv *http.Server, c chan bool) {
		<-c // block until a value is send to the channel
		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Println("error closing the server")
		}
	}(ctx, &srv, stopChannel)

	err := srv.ListenAndServeTLS("./tls.crt", "./tls.key")
	// erstell das zertifikat aus cert-manager
	//
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error starting server: %v\n", err)
	}

	fapi, err :=
		force.CreateWithAccessToken(
			"v53.0",
			os.Getenv("EVENT_CLIENT_ID"),
			authHandler.token.AccessToken,
			os.Getenv("EVENT_SCINSTANCE"),
			http.DefaultClient,
		)

	return fapi, err
}
