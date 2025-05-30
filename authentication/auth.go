package authentication

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pflege-de/go-force/force"
)

func NewForce() (force.ForceApiInterface, error) {
	sfRequest := AuthenticationRequest{
		URL:      os.Getenv("SF_SCAUD"),
		Username: os.Getenv("SF_SCUSER"),
		ClientID: os.Getenv("SF_CLIENT_ID"),
	}

	privateKeyFile, err := os.Open(os.Getenv("SF_SCKEY"))
	if err != nil {
		return nil, err
	}

	authReponse, err := Authenticate(sfRequest, privateKeyFile, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	fapi, err := force.CreateWithAccessToken(
		"v53.0",
		os.Getenv("SF_CLIENT_ID"),
		authReponse.GetToken(),
		os.Getenv("SF_SCINSTANCE"),
		http.DefaultClient,
	)
	return fapi, err
}

func NewForceKeyStringSecret() (force.ForceApiInterface, error) {
	sfRequest := AuthenticationRequest{
		URL:      os.Getenv("SF_SCAUD"),
		Username: os.Getenv("SF_SCUSER"),
		ClientID: os.Getenv("SF_CLIENT_ID"),
	}

	r := io.NopCloser(strings.NewReader(os.Getenv("SF_SCKEY")))

	authReponse, err := Authenticate(sfRequest, r, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	fapi, err := force.CreateWithAccessToken(
		"v53.0",
		os.Getenv("SF_CLIENT_ID"),
		authReponse.GetToken(),
		os.Getenv("SF_SCINSTANCE"),
		http.DefaultClient,
	)
	return fapi, err
}
