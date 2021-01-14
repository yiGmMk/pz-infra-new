package tokenverifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	_FACEBOOK_ACCESS_TOKEN_VERIFY_URL = "https://graph.facebook.com/v2.5/me"
	_FACEBOOK_ACCESS_TOKEN_FIELDS     = "id,name,first_name,last_name,age_range,link,gender,locale,picture,timezone,updated_time,verified,email"
)

type FacebookPictureData struct {
	IsSilhouette bool   `json:"is_silhouette"`
	Url          string `json:"url"`
}

type FacebookPicture struct {
	Data FacebookPictureData `json:"data"`
}

type FacebookAgeRange struct {
	Min int `json:"min"`
}

type FacebookPublicProfile struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	AgeRange    FacebookAgeRange `json:"age_range"`
	Link        string           `json:"link"`
	Gender      string           `json:"gender"`
	Locale      string           `json:"locale"`
	Picture     FacebookPicture  `json:"picture"`
	Timezone    int              `json:"timezone"`
	UpdatedTime string           `json:"updated_time"`
	Verified    bool             `json:"verified"`
}

type FacebookEmail struct {
	Email string `json:"email"`
}

type FacebookToken struct {
	FacebookPublicProfile
	FacebookEmail
}

type FacebookError struct {
	Message     string `json:"message"`
	Code        int    `json:"code"`
	Subcode     int    `json:"error_subcode"`
	UserMessage string `json:"error_user_msg"`
	UserTitle   string `json:"error_user_title"`
	TraceId     string `json:"fbtrace_id"`
}

type FacebookErrorResponse struct {
	FacebookError FacebookError `json:"error"`
}

func (fer *FacebookErrorResponse) Error() string {
	return fmt.Sprintf("invalid Facebook access token [%+v]", *fer)
}

func (ft FacebookToken) GetId() string {
	return ft.Id
}

func (ft FacebookToken) GetName() string {
	return ft.Name
}

func (ft FacebookToken) GetPicture() string {
	return ft.Picture.Data.Url
}

func (ft FacebookToken) GetEmail() string {
	return ft.Email
}

type facebookTokenVerifier struct {
	verifyUrl string
	fields    string
}

func (gtv *facebookTokenVerifier) SetVerifyUrl(verifyUrl string) {
	gtv.verifyUrl = verifyUrl
}

func (gtv *facebookTokenVerifier) GetVerifyUrl() string {
	if len(gtv.verifyUrl) == 0 {
		gtv.verifyUrl = _FACEBOOK_ACCESS_TOKEN_VERIFY_URL
	}
	return gtv.verifyUrl
}

func (gtv *facebookTokenVerifier) SetFields(fields string) {
	gtv.fields = fields
}

func (gtv *facebookTokenVerifier) GetFields() string {
	if len(gtv.fields) == 0 {
		gtv.fields = _FACEBOOK_ACCESS_TOKEN_FIELDS
	}
	return gtv.fields
}

func (ftv *facebookTokenVerifier) Verify(tokenString string) (Token, error) {
	// As a best practice, for large requests use a POST request instead of a GET request add a method=GET parameter.
	// the POST will be interpreted as if it were a GET.
	postData := make(url.Values)
	postData["method"] = []string{"GET"}
	postData["fields"] = []string{ftv.GetFields()}
	postData["access_token"] = []string{tokenString}
	resp, err := http.PostForm(ftv.GetVerifyUrl(), postData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		ft := &FacebookToken{}
		if err = json.Unmarshal(body, ft); err != nil {
			return nil, err
		}
		return ft, nil
	} else {
		// If the token has become invalid, the API will return an HTTP 400 status code,
		// a code and a subcode in a JSON body explaining the nature of the error.
		fer := &FacebookErrorResponse{}
		if err = json.Unmarshal(body, fer); err != nil {
			return nil, err
		}
		return nil, fer
	}
}
