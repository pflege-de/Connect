package authentication

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
)

const (
	grantType     string = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	tokenEndpoint string = "/services/oauth2/token" //nolint:gosec (Wrongly recognized - G101: Potential hardcoded credentials)
)

type AuthenticationRequest struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	ClientID string `json:"-"`
}

type authenticationResponse struct {
	Token       string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

type AuthenticationResponse interface {
	GetToken() string
	GetInstanceURL() string
	GetID() string
	GetTokenType() string
	GetIssuedAt() string
	GetSignature() string
}

// DeviceClaims represents the needed claims for the server-to-server oauth flow
// https://help.salesforce.com/s/articleView?id=sf.remoteaccess_oauth_jwt_flow.htm&type=5
type DeviceClaims struct {
	Audience string `json:"aud"`
	// Audience aus dem RegisteredClaims wird standardkonform in ein Array umgewandelt
	// Salesforce akzeptiert hier jedoch nur Strings
	jwt.RegisteredClaims
}

// GetToken returns the authenication token.
func (response authenticationResponse) GetToken() string { return response.Token }

// GetInstanceURL returns the Salesforce instance URL to use with the authenication information.
func (response authenticationResponse) GetInstanceURL() string { return response.InstanceURL }

// GetID returns the Salesforce ID of the authenication.
func (response authenticationResponse) GetID() string { return response.ID }

// GetTokenType returns the authenication token type.
func (response authenticationResponse) GetTokenType() string { return response.TokenType }

// GetIssuedAt returns the time when the token was issued.
func (response authenticationResponse) GetIssuedAt() string { return response.IssuedAt }

// GetSignature returns the signature of the authenication.
func (response authenticationResponse) GetSignature() string { return response.Signature }

// Authenicate will exchange the JWT signed request for access token.
func Authenticate(request AuthenticationRequest, privateKey io.ReadCloser, client *http.Client) (AuthenticationResponse, error) {

	pemData, err := io.ReadAll(privateKey)
	if err != nil {
		return nil, err
	}
	err = privateKey.Close()
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, DeviceClaims{
		request.URL,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Minute)),
			Issuer:    request.ClientID,
			Subject:   request.Username,
			ID:        uuid.New().String(),
		}})

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemData)
	if err != nil {
		return nil, err
	}

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("grant_type", grantType)
	form.Add("assertion", tokenString)

	tokenURL := request.URL + tokenEndpoint
	httpRequest, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, respErr := client.Do(httpRequest)

	if respErr != nil {
		return nil, errors.Wrap(respErr, "response is bad")
	}

	body, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return nil, bodyErr
	}

	if response.StatusCode >= 300 {
		fmt.Println(string(body))
		return nil, errors.New(response.Status)
	}

	var jsonResponse authenticationResponse

	unMarshallErr := json.Unmarshal(body, &jsonResponse)

	if unMarshallErr != nil {
		return nil, unMarshallErr
	}

	return jsonResponse, nil

}
