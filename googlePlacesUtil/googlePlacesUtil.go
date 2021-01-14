package googlePlacesUtil

import (
	"github.com/gyf841010/pz-infra-new/log"

	"github.com/astaxie/beego/config"
)

var defaultApiKey string
var defaultLanguage string = "en"
var defaultMaxResultCount int = 3

func Initialize(conf config.Configer) {
	if conf != nil {
		defaultApiKey = conf.DefaultString("GooglePlacesApiKey", "")
		defaultLanguage = conf.DefaultString("GooglePlacesApiLanguage", "en")
		defaultMaxResultCount = conf.DefaultInt("GooglePlacesApiMaxResultCount", 3)
	} else {
		log.Error("conf is nil")
	}
}

type SearchPlacesQuery struct {
	Query          string
	ApiKey         string
	Language       string
	MaxResultCount int
}

func NewSearchPlacesQuery(query string) *SearchPlacesQuery {
	return &SearchPlacesQuery{
		Query:          query,
		ApiKey:         defaultApiKey,
		Language:       defaultLanguage,
		MaxResultCount: defaultMaxResultCount,
	}
}

func (this *SearchPlacesQuery) WithApiKey(apiKey string) *SearchPlacesQuery {
	this.ApiKey = apiKey
	return this
}

func (this *SearchPlacesQuery) WithLanguage(language string) *SearchPlacesQuery {
	this.Language = language
	return this
}

func (this *SearchPlacesQuery) WithMaxResultCount(maxResultCount int) *SearchPlacesQuery {
	this.MaxResultCount = maxResultCount
	return this
}

func (this *SearchPlacesQuery) SearchPlaces() ([]string, error) {
	if this.Query == "" {
		log.Warn("query is empty")
		return nil, nil
	}
	if this.ApiKey == "" {
		return nil, log.Error("Google places api key is empty")
	}
	return googlePlacesApiInst.SearchPlaces(this)
}
