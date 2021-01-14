package googlePlacesUtil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"

	"github.com/yiGmMk/pz-infra-new/log"
)

type googlePlacesResultItem struct {
	Name string `json:"name"`
}

type googlePlacesResponse struct {
	Status  string                   `json:"status"`
	Results []googlePlacesResultItem `json:"results"`
}

var googlePlacesApiInst googlePlacesApi = &googlePlacesApiImpl{}

type googlePlacesApi interface {
	SearchPlaces(*SearchPlacesQuery) ([]string, error)
}

type googlePlacesApiImpl struct {
}

func (*googlePlacesApiImpl) SearchPlaces(query *SearchPlacesQuery) ([]string, error) {
	googlePlacesApiUrl := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	queryParameters := url.Values{}
	queryParameters.Set("query", query.Query)
	queryParameters.Set("key", query.ApiKey)
	queryParameters.Set("language", query.Language)
	finalUrl := fmt.Sprintf("%s?%s", googlePlacesApiUrl, queryParameters.Encode())
	resp, err := http.Get(finalUrl)
	if err != nil {
		return nil, log.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, log.Error(fmt.Sprintf("Google places api response %d", resp.StatusCode))
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, log.Error(err)
	}
	if len(bodyBytes) <= 0 {
		return nil, log.Error("empty Google place Api response")
	}
	var placesResp googlePlacesResponse
	if err = json.Unmarshal(bodyBytes, &placesResp); err != nil {
		return nil, log.Error(err)
	}
	if placesResp.Status == "ZERO_RESULTS" {
		log.Info("Google places returns ZERO_RESULTS")
		return nil, nil
	}
	if placesResp.Status != "OK" {
		return nil, log.Error("Google places response status: %s", placesResp.Status)
	}
	result := make([]string, 0, int(math.Min(float64(query.MaxResultCount), float64(len(placesResp.Results)))))
	index := 0
	for index < query.MaxResultCount && index < len(placesResp.Results) {
		name := placesResp.Results[index].Name
		upperName := strings.ToUpper(name)
		found := false
		for _, item := range result {
			if strings.ToUpper(item) == upperName {
				found = true
				break
			}
		}
		if !found {
			result = append(result, name)
		}
		index += 1
	}
	return result, nil
}
