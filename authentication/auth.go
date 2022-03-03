package authentication

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/nimajalali/go-force/force"
)

func NewForce() (*force.ForceApi, error) {
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

	return force.CreateWithAccessToken(
		"v53.0",
		os.Getenv("SF_CLIENT_ID"),
		authReponse.GetToken(),
		os.Getenv("SF_SCINSTANCE"))
}

func NewForceKeyStringSecret() (*force.ForceApi, error) {
	sfRequest := AuthenticationRequest{
		URL:      os.Getenv("SF_SCAUD"),
		Username: os.Getenv("SF_SCUSER"),
		ClientID: os.Getenv("SF_CLIENT_ID"),
	}

	key, _ := x509.ParsePKCS1PrivateKey([]byte(os.Getenv("SF_SCKEY")))

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	r := ioutil.NopCloser(strings.NewReader(string(pemdata))) // r type is io.ReadCloser
	defer r.Close()

	authReponse, err := Authenticate(sfRequest, r, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	return force.CreateWithAccessToken(
		"v53.0",
		os.Getenv("SF_CLIENT_ID"),
		authReponse.GetToken(),
		os.Getenv("SF_SCINSTANCE"))
}
