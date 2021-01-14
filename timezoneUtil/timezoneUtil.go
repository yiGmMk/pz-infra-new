package timezoneUtil

import (
	"encoding/json"
	"fmt"

	"github.com/gyf841010/pz-infra-new/geoutil"
	"github.com/gyf841010/pz-infra-new/httpUtil"
	. "github.com/gyf841010/pz-infra-new/logging"
	timeutil "github.com/gyf841010/pz-infra-new/timeUtil"

	"github.com/astaxie/beego"
)

const (
	GoogleTimezoneApiPath = "https://maps.googleapis.com/maps/api/timezone/json?key=%s&location=%s&timestamp=%d&language=%s"
	DEFAULT_LANGUAGE      = "en"
	STATUS_OK             = "OK"
)

type TimezoneInfo struct {
	TimeZoneId   string `json:"timeZoneId"`
	TimeZoneName string `json:"timeZoneName"`
	DstOffset    int64  `json:"dstOffset" description:"以秒数表示的夏令时偏移。"`
	RawOffset    int64  `json:"rawOffset" description:"给定位置与协调世界时的偏移（单位：秒）。该值未将夏令时考虑在内。"`
}

type TimezoneResp struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	TimezoneInfo
}

type GoogleTimezoneClient struct {
	GoogleApiKey string
	Lat          float64 // Latitude of the marker
	Lng          float64 // Longitude of the marker
	Language     string
}

func NewGoogleTimezoneClient() *GoogleTimezoneClient {
	client := GoogleTimezoneClient{
		GoogleApiKey: beego.AppConfig.String("google.apiKey"),
		Language:     DEFAULT_LANGUAGE,
	}
	return &client
}

func (c *GoogleTimezoneClient) SetLocation(lat, lng float64) {
	c.Lat = lat
	c.Lng = lng
}

func (c *GoogleTimezoneClient) SetLanguage(language string) {
	c.Language = language
}

// build the query API URL
func (c *GoogleTimezoneClient) buildApiUrl() string {
	location := geoutil.FormatLatLng(c.Lat, c.Lng)
	timestamp := timeutil.CurrentUnixInt()
	// API Key
	apiUrl := fmt.Sprintf(GoogleTimezoneApiPath, c.GoogleApiKey, location, timestamp, c.Language)
	return apiUrl
}

// Send the request to Static Map API, return the Content of the Image
func (c *GoogleTimezoneClient) Send() (*TimezoneInfo, error) {
	requestUrl := c.buildApiUrl()
	respContent, err := httpUtil.GetJson(requestUrl, nil)
	if err != nil {
		Log.Error("Failed to read response body from Google Timezone Client", With("requestUrl", requestUrl), WithError(err))
		return nil, err
	}
	var searchResponse TimezoneResp
	err = json.Unmarshal(respContent, &searchResponse)
	if err != nil {
		Log.Error("Failed to Unmarshall Google Timezone Search Content", With("requestUrl", requestUrl), With("respContent", string(respContent)), WithError(err))
		return nil, err
	}
	Log.Debug("Find Google Timezone Search Response", With("respContent", string(respContent)))

	if searchResponse.Status != STATUS_OK {
		return nil, Log.Error("No Timezone Info Found", With("lat", c.Lat), With("lng", c.Lng), With("status", searchResponse.Status), With("errorMessage", searchResponse.ErrorMessage))
	}
	timezoneInfo := searchResponse.TimezoneInfo
	Log.Debug("Google Timezone Search", With("timeZone", timezoneInfo.TimeZoneId))
	return &timezoneInfo, nil
}
