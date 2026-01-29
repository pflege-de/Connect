package connect

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pflege-de/connect/authentication"

	"github.com/pflege-de/go-force/force"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func GetForceApi() force.ForceApiInterface {
	forceApi, err := authentication.NewForce()
	if err != nil {
		log.Printf("Couldn't establish forceApi: %s", err)
	}
	return forceApi
}

func GetForceApiKeyStringSecret() force.ForceApiInterface {
	forceApi, err := authentication.NewForceKeyStringSecret()
	if err != nil {
		log.Printf("Couldn't establish forceApi: %s", err)
	}
	return forceApi
}

func GetForceApiOAuth() force.ForceApiInterface {
	forceApi, err := authentication.NewOAuthForce()
	if err != nil {
		log.Printf("Couldn't establish forceApi: %s", err)
	}
	return forceApi
}

type ExternalClientOption interface {
	noop() // TODO: implement option pattern: implement via functions; modify struct pointer
}

func GetForceExternalClient(ctx context.Context, options ...ExternalClientOption) (force.ForceApiInterface, error) {
	// NOTE: in the long term, it would be sensible to finally stop reaching deep
	// into the environment variables and pulling out values for "magic" keys --
	// but since this is just a hotfix, let's keep it compatible with the current
	// deployment configuration as much as possible for now (improve API later).
	//
	// TODO: allow the caller to construct and pass the `clientcredentials.Config`.
	// TODO: allow the caller to pass the `http.Client` (i.e., override default).

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SF_CLIENT_ID"),
		ClientSecret: os.Getenv("SF_CLIENT_SECRET"),

		TokenURL: os.Getenv("SF_TOKEN_URL"),

		Scopes: []string{},

		AuthStyle: oauth2.AuthStyleInParams,
	}

	force, err := authentication.NewForceClientCredentials(ctx, http.DefaultClient, config)

	if err != nil {
		return nil, fmt.Errorf("cannot authenticate via client credentials: %w", err)
	}

	return force, nil
}
