package authentication

import (
	"crypto/x509"
	"io/ioutil"
	"log"
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

	//ssh.ParseRawPrivateKey
	//x509.ParsePKCS1PrivateKey([]byte(os.Getenv("SF_SCKEY")))
	key, err := x509.ParsePKCS8PrivateKey([]byte(os.Getenv("SF_SCKEY")))
	if err != nil {
		return nil, err
	}
	log.Println(key)

	r := ioutil.NopCloser(strings.NewReader(string("test"))) // r type is io.ReadCloser
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
