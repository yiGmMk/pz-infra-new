package staticMapUtil

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/yiGmMk/pz-infra-new/geoutil"
	. "github.com/yiGmMk/pz-infra-new/logging"

	"github.com/astaxie/beego"
	"github.com/twpayne/go-polyline"
)

const (
	IMAGE_TYPE_PNG = "png"
	IMAGE_TYPE_JPG = "jpg"

	DefaultRoadMapType = "roadmap"
	DefaultImageSize   = "480x620"
	DefaultScaleValue  = 2
	BriefImageSize     = "644x200"
	BriefScaleValue    = 1

	StartPointColor = "0x34C981FF"
	EndPointColor   = "0xF64444FF"
	PathColor       = "0x24B4E7FF"
	PathWeight      = 6

	GoogleStaticMapApiPath = "https://maps.googleapis.com/maps/api/staticmap"
)

type GoogleStaticMapClient struct {
	GoogleApiKey string
	ImageType    string
	Center       string // Center & ZoomLevel are no need to set if Markers or Path are provided
	ZoomLevel    int
	ImageSize    string      // Size with format as width x height, such as "480x620"
	Scale        int         // Scale of dpi density, supported values 1, 2, default using 2 as more clear
	MapType      string      // Map Type, supporting values are roadmap(default), satellite, terrain, hybrid
	Markers      []mapMarker // Map Markers
	Paths        []mapPath   // Map Path
}

func NewGoogleStaticMapClient() *GoogleStaticMapClient {
	client := GoogleStaticMapClient{
		GoogleApiKey: beego.AppConfig.String("google.apiKey"),
		ImageType:    IMAGE_TYPE_PNG,
		Scale:        DefaultScaleValue,
		ImageSize:    DefaultImageSize,
		MapType:      DefaultRoadMapType,
		Markers:      make([]mapMarker, 0),
		Paths:        make([]mapPath, 0),
	}
	return &client
}

type mapMarker struct {
	IsCustom bool    // whether using custom icons, if true, icon field should be set
	Color    string  // RGB Color for the Marker, in format of 0xF64444FF
	Icon     string  // For Custom Marker Only, the url of Icon
	Lat      float64 // Latitude of the marker
	Lng      float64 // Longitude of the marker
}

func (m *mapMarker) BuildQuery() string {
	if m.IsCustom {
		// Custom Marker: icon:http://goo.gl/TgkC8N|37.76567,-122.47544
		return fmt.Sprintf("markers=icon:%s|%s", url.QueryEscape(m.Icon), geoutil.FormatLatLng(m.Lat, m.Lng))
	}
	// Default Marker : color:0x34C981FF|37.765542,-122.477998
	return fmt.Sprintf("markers=color:%s|%s", m.Color, geoutil.FormatLatLng(m.Lat, m.Lng))
}

type mapPath struct {
	Color          string // RGB Color for the Path, in format of 0xF64444FF
	Weight         int    // Weight of the Path
	PolylineEncode string // Google Polyline Encoding String
}

func (p *mapPath) BuildQuery() string {
	// path = color:0x24b400AA|weight:18|geodesic:true|enc:wboeFrpojVHzLX~N
	return fmt.Sprintf("path=color:%s|weight:%d|geodesic:true|enc:%s", p.Color, p.Weight, url.QueryEscape(p.PolylineEncode))
}

func (c *GoogleStaticMapClient) SetApiKey(apiKey string) {
	c.GoogleApiKey = apiKey
}

func (c *GoogleStaticMapClient) SetImageType(imgType string) {
	c.ImageType = imgType
}

func (c *GoogleStaticMapClient) SetImageSize(imgSize string) {
	c.ImageSize = imgSize
}

func (c *GoogleStaticMapClient) SetCenter(lat, lng float64) {
	c.Center = geoutil.FormatLatLng(lat, lng)
}

func (c *GoogleStaticMapClient) SetZoomLevel(zoomLevel int) {
	c.ZoomLevel = zoomLevel
}

func (c *GoogleStaticMapClient) SetScale(scale int) {
	c.Scale = scale
}

func (c *GoogleStaticMapClient) SetMapType(mapType string) {
	c.MapType = mapType
}

func (c *GoogleStaticMapClient) AddDefaultMarker(color string, lat, lng float64) {
	defaultMarker := mapMarker{
		IsCustom: false,
		Color:    color,
		Lat:      lat,
		Lng:      lng,
	}
	c.Markers = append(c.Markers, defaultMarker)
}

func (c *GoogleStaticMapClient) AddCustomMarker(iconUrl string, lat, lng float64) {
	customMarker := mapMarker{
		IsCustom: true,
		Icon:     iconUrl,
		Lat:      lat,
		Lng:      lng,
	}
	c.Markers = append(c.Markers, customMarker)
}

// Adding Coordinates of one path into the map.
//
// Notice: If the length of coordinates larger than 1000, will truncate to 1000 by scattering out some points inner them
func (c *GoogleStaticMapClient) AddPath(color string, weight int, coordinates [][]float64) {
	path := mapPath{
		Color:  color,
		Weight: weight,
	}
	// To Filter out the points greater than 1000
	totalPoints := len(coordinates)
	filteredCoordinates := make([][]float64, 0)
	if totalPoints > 1000 {
		numToReduce := totalPoints - 1000
		numHasReduce := 0
		step := int(math.Ceil(float64(totalPoints) / float64(1000)))
		Log.Debug("Reduce Coordinates should has step no less than 2", With("numToReduce", numToReduce), With("step", step))
		for i := 0; i < totalPoints; i = i + step {
			filteredCoordinates = append(filteredCoordinates, coordinates[i])
			numHasReduce = numHasReduce + step - 1
			if numHasReduce >= numToReduce {
				Log.Debug("Already met the numToReduce", With("numToReduce", numToReduce), With("numHasReduce", numHasReduce))
				if totalPoints-1 > i+step {
					filteredCoordinates = append(filteredCoordinates, coordinates[i+step:totalPoints-1]...)
				} else {
					filteredCoordinates = append(filteredCoordinates, coordinates[totalPoints-1])
				}
				break
			}
		}
	} else {
		filteredCoordinates = coordinates
	}

	path.PolylineEncode = string(polyline.EncodeCoords(filteredCoordinates))
	c.Paths = append(c.Paths, path)
}

// build the query API URL
func (c *GoogleStaticMapClient) buildApiUrl() string {
	baseUrl := GoogleStaticMapApiPath
	// API Key
	apiUrl := baseUrl + "?key=" + c.GoogleApiKey
	// Image Scale
	if c.Scale > 0 {
		apiUrl = fmt.Sprintf("%s&scale=%d", apiUrl, c.Scale)
	}
	// Image Format Type
	if c.ImageType != "" {
		apiUrl = apiUrl + "&format=" + c.ImageType
	}
	// Image Size
	if c.ImageSize != "" {
		apiUrl = apiUrl + "&size=" + c.ImageSize
	}
	// Center
	if c.Center != "" {
		apiUrl = apiUrl + "&center=" + c.Center
	}
	// Zoom Level
	if c.ZoomLevel > 0 {
		apiUrl = apiUrl + "&zoom=" + strconv.Itoa(c.ZoomLevel)
	}
	// Map Markers
	if len(c.Markers) > 0 {
		for _, marker := range c.Markers {
			apiUrl = fmt.Sprintf("%s&%s", apiUrl, marker.BuildQuery())
		}
	}
	// Map Paths
	if len(c.Paths) > 0 {
		for _, path := range c.Paths {
			apiUrl = fmt.Sprintf("%s&%s", apiUrl, path.BuildQuery())
		}
	}
	return apiUrl
}

// Send the request to Static Map API, return the Content of the Image
func (c *GoogleStaticMapClient) Send() ([]byte, error) {
	requestUrl := c.buildApiUrl()
	resp, err := http.Get(requestUrl)
	Log.Debug("Calling Google Static Map API", With("requestUrl", requestUrl))
	if err != nil {
		return nil, Log.Error("Failed to call Google Static Map API", With("requestUrl", requestUrl), WithError(err))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, Log.Error("Failed to read response body from Google Static Map", With("requestUrl", requestUrl), WithError(err))
	}
	return respBody, nil
}
