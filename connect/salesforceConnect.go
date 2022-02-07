package connect

import (
	"log"

	"github.com/nimajalali/go-force/force"
	"github.com/pflege-de/connect/authentication"
)

func GetForceApi() *force.ForceApi {
	forceApi, err := authentication.NewForce()
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
