package tokenverifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	_GOOGLE_ID_TOKEN_VERIFY_URL = "https://www.googleapis.com/oauth2/v3/tokeninfo"
)

type GoogleToken struct {
	Issuer                string `json:"iss"`
	AccessTokenHash       string `json:"at_hash"`
	Audience              string `json:"aud"`
	Subject               string `json:"sub"`
	EmailVerified         string `json:"email_verified"`
	AuthorizedParty       string `json:"azp"`
	HostedDomain          string `json:"hd"`
	Email                 string `json:"email"`
	IssuedAtTimeSeconds   string `json:"iat"`
	ExpirationTimeSeconds string `json:"exp"`
	Picture               string `json:"picture"`
	Name                  string `json:"name"`
	GivenName             string `json:"given_name"`
	FamilyName            string `json:"family_name"`
	Locale                string `json:"locale"`
	Algorithm             string `json:"alg"`
	keyId                 string `json:"kid"`
}

func (gt GoogleToken) GetId() string {
	return gt.Subject
}

func (gt GoogleToken) GetName() string {
	return gt.Name
}

func (gt GoogleToken) GetPicture() string {
	return gt.Picture
}

func (gt GoogleToken) GetEmail() string {
	return gt.Email
}

type GoogleErrorResponse struct {
	Description string `json:"error_description"`
}

func (ger *GoogleErrorResponse) Error() string {
	return fmt.Sprintf("invalid Google ID token [%+v]", *ger)
}

type googleTokenVerifier struct {
	verifyUrl string
}

func (gtv *googleTokenVerifier) SetVerifyUrl(verifyUrl string) {
	gtv.verifyUrl = verifyUrl
}

func (gtv *googleTokenVerifier) GetVerifyUrl() string {
	if len(gtv.verifyUrl) == 0 {
		gtv.verifyUrl = _GOOGLE_ID_TOKEN_VERIFY_URL
	}
	return gtv.verifyUrl
}

func (gtv *googleTokenVerifier) Verify(tokenString string) (Token, error) {
	postData := make(url.Values)
	postData["id_token"] = []string{tokenString}
	resp, err := http.PostForm(gtv.GetVerifyUrl(), postData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// If the token is properly signed and the iss and exp claims have the expected values,
	// we will get a HTTP 200 response, where the body contains the JSON-formatted ID token claims.
	if resp.StatusCode == 200 {
		gt := &GoogleToken{}
		if err = json.Unmarshal(body, gt); err != nil {
			return nil, err
		}
		return gt, nil
	} else {
		// If the token has expired, has been tampered with, or the permissions revoked,
		// the Google Authorization Server will respond with an error.
		// The error surfaces as a 400 status code, and a JSON body.
		ger := &GoogleErrorResponse{}
		if err = json.Unmarshal(body, ger); err != nil {
			return nil, err
		}
		return nil, ger
	}
}
