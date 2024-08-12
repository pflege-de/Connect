package connect

import (
	"log"

	"github.com/pflege-de/connect/authentication"
	"github.com/pflege-de/go-force/force"
)

func GetForceApi() *force.ForceApi {
	forceApi, err := authentication.NewForce()
	if err != nil {
		log.Printf("Couldn't establish forceApi: %s", err)
	}
	return forceApi
}

func GetForceApiKeyStringSecret() *force.ForceApi {
	forceApi, err := authentication.NewForceKeyStringSecret()
	if err != nil {
		log.Printf("Couldn't establish forceApi: %s", err)
	}
	return forceApi
}

func GetForceApiOAuth() *force.ForceApi {
	forceApi, err := authentication.NewOAuthForce()
	if err != nil {
		log.Printf("Couldn't establish forceApi: %s", err)
	}
	return forceApi
}
