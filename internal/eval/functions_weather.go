// internal/eval/functions_weather.go

package eval

import (
	"context"
	"time"

	"github.com/0xsj/numio/internal/fetch"
	"github.com/0xsj/numio/pkg/types"
)

func registerWeatherFunctions() {
	register("weather", 1, 1, false, FnWeather)
	register("temp", 1, 1, false, FnTemp)
	register("humidity", 1, 1, false, FnHumidity)
	register("wind", 1, 1, false, FnWind)
	register("feelslike", 1, 1, false, FnFeelsLike)
}

// fetchWeather is a shared helper that fetches weather for a city string arg.
func fetchWeather(args []types.Value) (*fetch.WeatherResult, types.Value) {
	city := args[0].AsString()
	if city == "" {
		return nil, types.Errorf("weather: city name required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := fetch.FetchWeather(ctx, city)
	if err != nil {
		return nil, types.Errorf("weather: %s", err.Error())
	}

	return result, types.Empty()
}

// FnWeather returns a formatted weather summary string.
// Usage: weather("istanbul")
func FnWeather(args []types.Value) types.Value {
	w, errVal := fetchWeather(args)
	if w == nil {
		return errVal
	}

	summary := w.City + ": " +
		formatTemp(w.Temp) + ", " +
		w.Description + ", " +
		"Humidity: " + formatPct(w.Humidity) + ", " +
		"Wind: " + formatSpeed(w.WindSpeed)

	return types.StringValue(summary)
}

// FnTemp returns the current temperature as a numeric value (°C).
// Usage: temp("istanbul") → 8.5
func FnTemp(args []types.Value) types.Value {
	w, errVal := fetchWeather(args)
	if w == nil {
		return errVal
	}
	return types.Number(w.Temp)
}

// FnHumidity returns the current humidity as a numeric value (%).
// Usage: humidity("istanbul") → 72
func FnHumidity(args []types.Value) types.Value {
	w, errVal := fetchWeather(args)
	if w == nil {
		return errVal
	}
	return types.Number(w.Humidity)
}

// FnWind returns the current wind speed as a numeric value (km/h).
// Usage: wind("istanbul") → 15
func FnWind(args []types.Value) types.Value {
	w, errVal := fetchWeather(args)
	if w == nil {
		return errVal
	}
	return types.Number(w.WindSpeed)
}

// FnFeelsLike returns the apparent temperature as a numeric value (°C).
// Usage: feelslike("istanbul") → 6.2
func FnFeelsLike(args []types.Value) types.Value {
	w, errVal := fetchWeather(args)
	if w == nil {
		return errVal
	}
	return types.Number(w.FeelsLike)
}

// ════════════════════════════════════════════════════════════════
// FORMATTING HELPERS
// ════════════════════════════════════════════════════════════════

func formatTemp(t float64) string {
	return formatWeatherFloat(t) + "°C"
}

func formatPct(p float64) string {
	return formatWeatherFloat(p) + "%"
}

func formatSpeed(s float64) string {
	return formatWeatherFloat(s) + " km/h"
}

func formatWeatherFloat(f float64) string {
	negative := f < 0
	if negative {
		f = -f
	}

	intPart := int(f)
	fracPart := int((f - float64(intPart)) * 10)

	result := weatherIntToStr(intPart)
	if fracPart > 0 {
		result += "." + weatherIntToStr(fracPart)
	}

	if negative {
		return "-" + result
	}
	return result
}

func weatherIntToStr(n int) string {
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
