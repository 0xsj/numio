// internal/fetch/weather.go

package fetch

import (
	"context"
	"strings"
	"sync"
)

// ════════════════════════════════════════════════════════════════
// WEATHER TYPES
// ════════════════════════════════════════════════════════════════

// WeatherResult holds current weather data for a location.
type WeatherResult struct {
	City        string
	Country     string
	Temp        float64 // °C
	FeelsLike   float64 // °C
	Humidity    float64 // %
	WindSpeed   float64 // km/h
	WeatherCode int
	Description string
}

// ════════════════════════════════════════════════════════════════
// GEOCODING CACHE
// ════════════════════════════════════════════════════════════════

type geoCoord struct {
	Lat     float64
	Lon     float64
	Name    string
	Country string
}

var (
	geoCache   = map[string]*geoCoord{}
	geoCacheMu sync.RWMutex
)

// ════════════════════════════════════════════════════════════════
// API RESPONSE TYPES
// ════════════════════════════════════════════════════════════════

type geocodeResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
	} `json:"results"`
}

type weatherResponse struct {
	Current struct {
		Temperature    float64 `json:"temperature_2m"`
		ApparentTemp   float64 `json:"apparent_temperature"`
		Humidity       float64 `json:"relative_humidity_2m"`
		WindSpeed      float64 `json:"wind_speed_10m"`
		WeatherCode    int     `json:"weather_code"`
	} `json:"current"`
}

// ════════════════════════════════════════════════════════════════
// FETCH WEATHER
// ════════════════════════════════════════════════════════════════

// FetchWeather fetches current weather for a city using the Open-Meteo API.
func FetchWeather(ctx context.Context, city string) (*WeatherResult, error) {
	coord, err := geocode(ctx, city)
	if err != nil {
		return nil, err
	}

	url := "https://api.open-meteo.com/v1/forecast?" +
		"latitude=" + formatCoord(coord.Lat) +
		"&longitude=" + formatCoord(coord.Lon) +
		"&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m,apparent_temperature" +
		"&temperature_unit=celsius&timezone=auto"

	var resp weatherResponse
	if err := GetJSON(ctx, url, &resp); err != nil {
		return nil, err
	}

	return &WeatherResult{
		City:        coord.Name,
		Country:     coord.Country,
		Temp:        resp.Current.Temperature,
		FeelsLike:   resp.Current.ApparentTemp,
		Humidity:    resp.Current.Humidity,
		WindSpeed:   resp.Current.WindSpeed,
		WeatherCode: resp.Current.WeatherCode,
		Description: wmoDescription(resp.Current.WeatherCode),
	}, nil
}

// geocode resolves a city name to coordinates using the Open-Meteo geocoding API.
func geocode(ctx context.Context, city string) (*geoCoord, error) {
	key := strings.ToLower(strings.TrimSpace(city))

	// Check cache
	geoCacheMu.RLock()
	if coord, ok := geoCache[key]; ok {
		geoCacheMu.RUnlock()
		return coord, nil
	}
	geoCacheMu.RUnlock()

	// URL-encode the city name (simple: replace spaces with +)
	encoded := strings.ReplaceAll(key, " ", "+")
	url := "https://geocoding-api.open-meteo.com/v1/search?name=" + encoded + "&count=1&language=en"

	var resp geocodeResponse
	if err := GetJSON(ctx, url, &resp); err != nil {
		return nil, err
	}

	if len(resp.Results) == 0 {
		return nil, ErrNotFound
	}

	r := resp.Results[0]
	coord := &geoCoord{
		Lat:     r.Latitude,
		Lon:     r.Longitude,
		Name:    r.Name,
		Country: r.Country,
	}

	// Store in cache
	geoCacheMu.Lock()
	geoCache[key] = coord
	geoCacheMu.Unlock()

	return coord, nil
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

// formatCoord formats a coordinate as a string without fmt.
func formatCoord(f float64) string {
	negative := f < 0
	if negative {
		f = -f
	}

	intPart := int(f)
	fracPart := int((f - float64(intPart)) * 10000)

	result := intToString(intPart) + "." + zeroPad(fracPart, 4)
	if negative {
		return "-" + result
	}
	return result
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func zeroPad(n, width int) string {
	s := intToString(n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

// wmoDescription maps WMO weather interpretation codes to descriptions.
func wmoDescription(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1:
		return "Mainly clear"
	case 2:
		return "Partly cloudy"
	case 3:
		return "Overcast"
	case 45:
		return "Fog"
	case 48:
		return "Depositing rime fog"
	case 51:
		return "Light drizzle"
	case 53:
		return "Moderate drizzle"
	case 55:
		return "Dense drizzle"
	case 56:
		return "Light freezing drizzle"
	case 57:
		return "Dense freezing drizzle"
	case 61:
		return "Slight rain"
	case 63:
		return "Moderate rain"
	case 65:
		return "Heavy rain"
	case 66:
		return "Light freezing rain"
	case 67:
		return "Heavy freezing rain"
	case 71:
		return "Slight snowfall"
	case 73:
		return "Moderate snowfall"
	case 75:
		return "Heavy snowfall"
	case 77:
		return "Snow grains"
	case 80:
		return "Slight rain showers"
	case 81:
		return "Moderate rain showers"
	case 82:
		return "Violent rain showers"
	case 85:
		return "Slight snow showers"
	case 86:
		return "Heavy snow showers"
	case 95:
		return "Thunderstorm"
	case 96:
		return "Thunderstorm with slight hail"
	case 99:
		return "Thunderstorm with heavy hail"
	default:
		return "Unknown"
	}
}
