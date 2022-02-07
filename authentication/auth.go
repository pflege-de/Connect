package authentication

import (
	"net/http"
	"os"

	"github.com/nimajalali/go-force/force"
)

func NewForce() (*force.ForceApi, error) {
	sfRequest := AuthenticationRequest{
		URL:      os.Getenv("EVENT_SCAUD"),
		Username: os.Getenv("EVENT_SCUSER"),
		ClientID: os.Getenv("EVENT_CLIENT_ID"),
	}

	privateKeyFile, err := os.Open(os.Getenv("EVENT_SCKEY"))
	if err != nil {
		return nil, err
	}

	authReponse, err := Authenticate(sfRequest, privateKeyFile, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	return force.CreateWithAccessToken(
		"v53.0",
		os.Getenv("EVENT_CLIENT_ID"),
		authReponse.GetToken(),
		os.Getenv("EVENT_SCINSTANCE"),
		http.DefaultClient)
}
