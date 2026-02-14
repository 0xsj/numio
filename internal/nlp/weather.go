// internal/nlp/weather.go

package nlp

// ════════════════════════════════════════════════════════════════
// WEATHER PATTERNS
// ════════════════════════════════════════════════════════════════

// RegisterWeatherPatterns registers all weather-related NLP patterns.
func RegisterWeatherPatterns(r *PatternRegistry) {
	r.RegisterAll([]*Pattern{
		patternWeatherCity(),
		patternTempCity(),
		patternHumidityCity(),
		patternWindCity(),
	})
}

// patternWeatherCity: "weather istanbul", "weather in tokyo", "weather new york"
func patternWeatherCity() *Pattern {
	return NewPattern("weather_city").
		Regex(`weather\s+(?:in\s+)?(.+)`).
		Keywords("weather").
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}
			city := matches[1]
			if city == "" {
				return "", false
			}
			return `weather("` + city + `")`, true
		}).
		Build()
}

// patternTempCity: "temperature istanbul", "temp tokyo", "temp in london"
func patternTempCity() *Pattern {
	return NewPattern("temp_city").
		Regex(`(?:temperature|temp)\s+(?:in\s+)?(.+)`).
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}
			city := matches[1]
			if city == "" {
				return "", false
			}
			return `temp("` + city + `")`, true
		}).
		Build()
}

// patternHumidityCity: "humidity istanbul", "humidity in tokyo"
func patternHumidityCity() *Pattern {
	return NewPattern("humidity_city").
		Regex(`humidity\s+(?:in\s+)?(.+)`).
		Keywords("humidity").
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}
			city := matches[1]
			if city == "" {
				return "", false
			}
			return `humidity("` + city + `")`, true
		}).
		Build()
}

// patternWindCity: "wind istanbul", "wind speed tokyo", "wind in london"
func patternWindCity() *Pattern {
	return NewPattern("wind_city").
		Regex(`wind\s+(?:speed\s+)?(?:in\s+)?(.+)`).
		Keywords("wind").
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}
			city := matches[1]
			if city == "" {
				return "", false
			}
			return `wind("` + city + `")`, true
		}).
		Build()
}
