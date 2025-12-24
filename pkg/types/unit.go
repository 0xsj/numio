// internal/types/unit.go

package types

import (
	"strings"
)

// UnitType represents a category of units.
type UnitType int

const (
	UnitTypeLength UnitType = iota
	UnitTypeWeight
	UnitTypeTime
	UnitTypeTemperature
	UnitTypeData
	UnitTypeArea
	UnitTypeVolume
	UnitTypeSpeed // Future: compound units
)

// String returns the unit type name.
func (t UnitType) String() string {
	switch t {
	case UnitTypeLength:
		return "length"
	case UnitTypeWeight:
		return "weight"
	case UnitTypeTime:
		return "time"
	case UnitTypeTemperature:
		return "temperature"
	case UnitTypeData:
		return "data"
	case UnitTypeArea:
		return "area"
	case UnitTypeVolume:
		return "volume"
	case UnitTypeSpeed:
		return "speed"
	default:
		return "unknown"
	}
}

// Unit represents a unit of measurement.
type Unit struct {
	Code        string   // Canonical code: "km", "lb", "h"
	Symbol      string   // Display symbol (often same as code)
	Name        string   // Full name: "kilometer", "pound", "hour"
	Plural      string   // Plural form: "kilometers", "pounds", "hours"
	Type        UnitType // Category
	Aliases     []string // Alternative names
	ToBase      float64  // Multiplier to convert to base unit
	FromBaseAdd float64  // Additive offset from base (for temperature)
	IsBase      bool     // True if this is the base unit for its type
}

// String returns the unit code.
func (u Unit) String() string {
	return u.Code
}

// ConvertTo converts a value from this unit to another unit of the same type.
// Returns the converted value and true if successful, or 0 and false if incompatible.
func (u Unit) ConvertTo(value float64, target *Unit) (float64, bool) {
	if u.Type != target.Type {
		return 0, false
	}

	// Special handling for temperature (non-linear conversion)
	if u.Type == UnitTypeTemperature {
		return convertTemperature(value, &u, target), true
	}

	// Linear conversion: value → base → target
	baseValue := value * u.ToBase
	targetValue := baseValue / target.ToBase

	return targetValue, true
}

// convertTemperature handles non-linear temperature conversions.
func convertTemperature(value float64, from, to *Unit) float64 {
	// Convert to Kelvin (base) first
	var kelvin float64
	switch from.Code {
	case "K":
		kelvin = value
	case "C":
		kelvin = value + 273.15
	case "F":
		kelvin = (value-32)*5/9 + 273.15
	default:
		kelvin = value
	}

	// Convert from Kelvin to target
	switch to.Code {
	case "K":
		return kelvin
	case "C":
		return kelvin - 273.15
	case "F":
		return (kelvin-273.15)*9/5 + 32
	default:
		return kelvin
	}
}

// UnitRegistry holds all known units.
type UnitRegistry struct {
	byCode  map[string]*Unit
	byAlias map[string]*Unit
	byType  map[UnitType][]*Unit
}

// Global unit registry.
var units = newUnitRegistry()

// newUnitRegistry creates and populates the unit registry.
func newUnitRegistry() *UnitRegistry {
	r := &UnitRegistry{
		byCode:  make(map[string]*Unit),
		byAlias: make(map[string]*Unit),
		byType:  make(map[UnitType][]*Unit),
	}

	for i := range curatedUnits {
		r.register(&curatedUnits[i])
	}

	return r
}

// register adds a unit to the registry.
func (r *UnitRegistry) register(u *Unit) {
	// By code (case-insensitive for most, but preserve case for symbols)
	r.byCode[u.Code] = u
	r.byCode[strings.ToLower(u.Code)] = u
	r.byCode[strings.ToUpper(u.Code)] = u

	// By aliases (case-insensitive)
	for _, alias := range u.Aliases {
		r.byAlias[strings.ToLower(alias)] = u
	}

	// By type
	r.byType[u.Type] = append(r.byType[u.Type], u)
}

// Lookup finds a unit by code or alias.
func (r *UnitRegistry) Lookup(s string) *Unit {
	// Try exact code match first
	if u, ok := r.byCode[s]; ok {
		return u
	}

	// Try case-insensitive code
	if u, ok := r.byCode[strings.ToLower(s)]; ok {
		return u
	}

	// Try alias (case-insensitive)
	if u, ok := r.byAlias[strings.ToLower(s)]; ok {
		return u
	}

	return nil
}

// curatedUnits contains all supported units.
// Base units have ToBase = 1.0
var curatedUnits = []Unit{
	// ════════════════════════════════════════════════════════════
	// LENGTH (base: meter)
	// ════════════════════════════════════════════════════════════
	{
		Code:    "m",
		Symbol:  "m",
		Name:    "meter",
		Plural:  "meters",
		Type:    UnitTypeLength,
		Aliases: []string{"meter", "meters", "metre", "metres"},
		ToBase:  1.0,
		IsBase:  true,
	},
	{
		Code:    "km",
		Symbol:  "km",
		Name:    "kilometer",
		Plural:  "kilometers",
		Type:    UnitTypeLength,
		Aliases: []string{"kilometer", "kilometers", "kilometre", "kilometres"},
		ToBase:  1000.0,
	},
	{
		Code:    "cm",
		Symbol:  "cm",
		Name:    "centimeter",
		Plural:  "centimeters",
		Type:    UnitTypeLength,
		Aliases: []string{"centimeter", "centimeters", "centimetre", "centimetres"},
		ToBase:  0.01,
	},
	{
		Code:    "mm",
		Symbol:  "mm",
		Name:    "millimeter",
		Plural:  "millimeters",
		Type:    UnitTypeLength,
		Aliases: []string{"millimeter", "millimeters", "millimetre", "millimetres"},
		ToBase:  0.001,
	},
	{
		Code:    "mi",
		Symbol:  "mi",
		Name:    "mile",
		Plural:  "miles",
		Type:    UnitTypeLength,
		Aliases: []string{"mile", "miles"},
		ToBase:  1609.344,
	},
	{
		Code:    "yd",
		Symbol:  "yd",
		Name:    "yard",
		Plural:  "yards",
		Type:    UnitTypeLength,
		Aliases: []string{"yard", "yards"},
		ToBase:  0.9144,
	},
	{
		Code:    "ft",
		Symbol:  "ft",
		Name:    "foot",
		Plural:  "feet",
		Type:    UnitTypeLength,
		Aliases: []string{"foot", "feet"},
		ToBase:  0.3048,
	},
	{
		Code:    "in",
		Symbol:  "in",
		Name:    "inch",
		Plural:  "inches",
		Type:    UnitTypeLength,
		Aliases: []string{"inch", "inches"},
		ToBase:  0.0254,
	},
	{
		Code:    "nm",
		Symbol:  "nm",
		Name:    "nautical mile",
		Plural:  "nautical miles",
		Type:    UnitTypeLength,
		Aliases: []string{"nautical mile", "nautical miles", "nmi"},
		ToBase:  1852.0,
	},

	// ════════════════════════════════════════════════════════════
	// WEIGHT / MASS (base: gram)
	// ════════════════════════════════════════════════════════════
	{
		Code:    "g",
		Symbol:  "g",
		Name:    "gram",
		Plural:  "grams",
		Type:    UnitTypeWeight,
		Aliases: []string{"gram", "grams"},
		ToBase:  1.0,
		IsBase:  true,
	},
	{
		Code:    "kg",
		Symbol:  "kg",
		Name:    "kilogram",
		Plural:  "kilograms",
		Type:    UnitTypeWeight,
		Aliases: []string{"kilogram", "kilograms", "kilo", "kilos"},
		ToBase:  1000.0,
	},
	{
		Code:    "mg",
		Symbol:  "mg",
		Name:    "milligram",
		Plural:  "milligrams",
		Type:    UnitTypeWeight,
		Aliases: []string{"milligram", "milligrams"},
		ToBase:  0.001,
	},
	{
		Code:    "t",
		Symbol:  "t",
		Name:    "metric ton",
		Plural:  "metric tons",
		Type:    UnitTypeWeight,
		Aliases: []string{"ton", "tons", "tonne", "tonnes", "metric ton", "metric tons"},
		ToBase:  1000000.0,
	},
	{
		Code:    "lb",
		Symbol:  "lb",
		Name:    "pound",
		Plural:  "pounds",
		Type:    UnitTypeWeight,
		Aliases: []string{"pound", "pounds", "lbs"},
		ToBase:  453.592,
	},
	{
		Code:    "oz",
		Symbol:  "oz",
		Name:    "ounce",
		Plural:  "ounces",
		Type:    UnitTypeWeight,
		Aliases: []string{"ounce", "ounces"},
		ToBase:  28.3495,
	},
	{
		Code:    "st",
		Symbol:  "st",
		Name:    "stone",
		Plural:  "stones",
		Type:    UnitTypeWeight,
		Aliases: []string{"stone", "stones"},
		ToBase:  6350.29,
	},
	{
		Code:    "ozt",
		Symbol:  "ozt",
		Name:    "troy ounce",
		Plural:  "troy ounces",
		Type:    UnitTypeWeight,
		Aliases: []string{"troy ounce", "troy ounces", "oz t"},
		ToBase:  31.1035,
	},

	// ════════════════════════════════════════════════════════════
	// TIME (base: second)
	// ════════════════════════════════════════════════════════════
	{
		Code:    "s",
		Symbol:  "s",
		Name:    "second",
		Plural:  "seconds",
		Type:    UnitTypeTime,
		Aliases: []string{"second", "seconds", "sec", "secs"},
		ToBase:  1.0,
		IsBase:  true,
	},
	{
		Code:    "ms",
		Symbol:  "ms",
		Name:    "millisecond",
		Plural:  "milliseconds",
		Type:    UnitTypeTime,
		Aliases: []string{"millisecond", "milliseconds"},
		ToBase:  0.001,
	},
	{
		Code:    "min",
		Symbol:  "min",
		Name:    "minute",
		Plural:  "minutes",
		Type:    UnitTypeTime,
		Aliases: []string{"minute", "minutes", "mins"},
		ToBase:  60.0,
	},
	{
		Code:    "h",
		Symbol:  "h",
		Name:    "hour",
		Plural:  "hours",
		Type:    UnitTypeTime,
		Aliases: []string{"hour", "hours", "hr", "hrs"},
		ToBase:  3600.0,
	},
	{
		Code:    "d",
		Symbol:  "d",
		Name:    "day",
		Plural:  "days",
		Type:    UnitTypeTime,
		Aliases: []string{"day", "days"},
		ToBase:  86400.0,
	},
	{
		Code:    "wk",
		Symbol:  "wk",
		Name:    "week",
		Plural:  "weeks",
		Type:    UnitTypeTime,
		Aliases: []string{"week", "weeks"},
		ToBase:  604800.0,
	},
	{
		Code:    "mo",
		Symbol:  "mo",
		Name:    "month",
		Plural:  "months",
		Type:    UnitTypeTime,
		Aliases: []string{"month", "months"},
		ToBase:  2629746.0, // Average month (30.44 days)
	},
	{
		Code:    "y",
		Symbol:  "y",
		Name:    "year",
		Plural:  "years",
		Type:    UnitTypeTime,
		Aliases: []string{"year", "years", "yr", "yrs"},
		ToBase:  31556952.0, // Average year (365.2425 days)
	},

	// ════════════════════════════════════════════════════════════
	// TEMPERATURE (base: Kelvin, but special conversion logic)
	// ════════════════════════════════════════════════════════════
	{
		Code:    "K",
		Symbol:  "K",
		Name:    "Kelvin",
		Plural:  "Kelvin",
		Type:    UnitTypeTemperature,
		Aliases: []string{"kelvin"},
		ToBase:  1.0,
		IsBase:  true,
	},
	{
		Code:    "C",
		Symbol:  "°C",
		Name:    "Celsius",
		Plural:  "Celsius",
		Type:    UnitTypeTemperature,
		Aliases: []string{"celsius", "centigrade"},
		ToBase:  1.0, // Special handling in convertTemperature
	},
	{
		Code:    "F",
		Symbol:  "°F",
		Name:    "Fahrenheit",
		Plural:  "Fahrenheit",
		Type:    UnitTypeTemperature,
		Aliases: []string{"fahrenheit"},
		ToBase:  1.0, // Special handling in convertTemperature
	},

	// ════════════════════════════════════════════════════════════
	// DATA (base: byte)
	// ════════════════════════════════════════════════════════════
	{
		Code:    "B",
		Symbol:  "B",
		Name:    "byte",
		Plural:  "bytes",
		Type:    UnitTypeData,
		Aliases: []string{"byte", "bytes"},
		ToBase:  1.0,
		IsBase:  true,
	},
	{
		Code:    "KB",
		Symbol:  "KB",
		Name:    "kilobyte",
		Plural:  "kilobytes",
		Type:    UnitTypeData,
		Aliases: []string{"kilobyte", "kilobytes", "kb"},
		ToBase:  1024.0,
	},
	{
		Code:    "MB",
		Symbol:  "MB",
		Name:    "megabyte",
		Plural:  "megabytes",
		Type:    UnitTypeData,
		Aliases: []string{"megabyte", "megabytes", "mb"},
		ToBase:  1048576.0, // 1024^2
	},
	{
		Code:    "GB",
		Symbol:  "GB",
		Name:    "gigabyte",
		Plural:  "gigabytes",
		Type:    UnitTypeData,
		Aliases: []string{"gigabyte", "gigabytes", "gb"},
		ToBase:  1073741824.0, // 1024^3
	},
	{
		Code:    "TB",
		Symbol:  "TB",
		Name:    "terabyte",
		Plural:  "terabytes",
		Type:    UnitTypeData,
		Aliases: []string{"terabyte", "terabytes", "tb"},
		ToBase:  1099511627776.0, // 1024^4
	},
	{
		Code:    "PB",
		Symbol:  "PB",
		Name:    "petabyte",
		Plural:  "petabytes",
		Type:    UnitTypeData,
		Aliases: []string{"petabyte", "petabytes", "pb"},
		ToBase:  1125899906842624.0, // 1024^5
	},
	{
		Code:    "bit",
		Symbol:  "bit",
		Name:    "bit",
		Plural:  "bits",
		Type:    UnitTypeData,
		Aliases: []string{"bits"},
		ToBase:  0.125, // 1/8 byte
	},
	{
		Code:    "Kbit",
		Symbol:  "Kbit",
		Name:    "kilobit",
		Plural:  "kilobits",
		Type:    UnitTypeData,
		Aliases: []string{"kilobit", "kilobits", "kbit"},
		ToBase:  128.0, // 1024 bits = 128 bytes
	},
	{
		Code:    "Mbit",
		Symbol:  "Mbit",
		Name:    "megabit",
		Plural:  "megabits",
		Type:    UnitTypeData,
		Aliases: []string{"megabit", "megabits", "mbit"},
		ToBase:  131072.0, // 1024^2 bits
	},
	{
		Code:    "Gbit",
		Symbol:  "Gbit",
		Name:    "gigabit",
		Plural:  "gigabits",
		Type:    UnitTypeData,
		Aliases: []string{"gigabit", "gigabits", "gbit"},
		ToBase:  134217728.0, // 1024^3 bits
	},

	// ════════════════════════════════════════════════════════════
	// AREA (base: square meter)
	// ════════════════════════════════════════════════════════════
	{
		Code:    "sqm",
		Symbol:  "m²",
		Name:    "square meter",
		Plural:  "square meters",
		Type:    UnitTypeArea,
		Aliases: []string{"square meter", "square meters", "sq m", "m2"},
		ToBase:  1.0,
		IsBase:  true,
	},
	{
		Code:    "sqkm",
		Symbol:  "km²",
		Name:    "square kilometer",
		Plural:  "square kilometers",
		Type:    UnitTypeArea,
		Aliases: []string{"square kilometer", "square kilometers", "sq km", "km2"},
		ToBase:  1000000.0,
	},
	{
		Code:    "sqft",
		Symbol:  "ft²",
		Name:    "square foot",
		Plural:  "square feet",
		Type:    UnitTypeArea,
		Aliases: []string{"square foot", "square feet", "sq ft", "ft2"},
		ToBase:  0.092903,
	},
	{
		Code:    "sqmi",
		Symbol:  "mi²",
		Name:    "square mile",
		Plural:  "square miles",
		Type:    UnitTypeArea,
		Aliases: []string{"square mile", "square miles", "sq mi", "mi2"},
		ToBase:  2589988.0,
	},
	{
		Code:    "acre",
		Symbol:  "acre",
		Name:    "acre",
		Plural:  "acres",
		Type:    UnitTypeArea,
		Aliases: []string{"acres"},
		ToBase:  4046.86,
	},
	{
		Code:    "ha",
		Symbol:  "ha",
		Name:    "hectare",
		Plural:  "hectares",
		Type:    UnitTypeArea,
		Aliases: []string{"hectare", "hectares"},
		ToBase:  10000.0,
	},

	// ════════════════════════════════════════════════════════════
	// VOLUME (base: liter)
	// ════════════════════════════════════════════════════════════
	{
		Code:    "L",
		Symbol:  "L",
		Name:    "liter",
		Plural:  "liters",
		Type:    UnitTypeVolume,
		Aliases: []string{"liter", "liters", "litre", "litres", "l"},
		ToBase:  1.0,
		IsBase:  true,
	},
	{
		Code:    "mL",
		Symbol:  "mL",
		Name:    "milliliter",
		Plural:  "milliliters",
		Type:    UnitTypeVolume,
		Aliases: []string{"milliliter", "milliliters", "millilitre", "millilitres", "ml"},
		ToBase:  0.001,
	},
	{
		Code:    "gal",
		Symbol:  "gal",
		Name:    "gallon",
		Plural:  "gallons",
		Type:    UnitTypeVolume,
		Aliases: []string{"gallon", "gallons"},
		ToBase:  3.78541, // US gallon
	},
	{
		Code:    "qt",
		Symbol:  "qt",
		Name:    "quart",
		Plural:  "quarts",
		Type:    UnitTypeVolume,
		Aliases: []string{"quart", "quarts"},
		ToBase:  0.946353,
	},
	{
		Code:    "pt",
		Symbol:  "pt",
		Name:    "pint",
		Plural:  "pints",
		Type:    UnitTypeVolume,
		Aliases: []string{"pint", "pints"},
		ToBase:  0.473176,
	},
	{
		Code:    "cup",
		Symbol:  "cup",
		Name:    "cup",
		Plural:  "cups",
		Type:    UnitTypeVolume,
		Aliases: []string{"cups"},
		ToBase:  0.236588,
	},
	{
		Code:    "floz",
		Symbol:  "fl oz",
		Name:    "fluid ounce",
		Plural:  "fluid ounces",
		Type:    UnitTypeVolume,
		Aliases: []string{"fluid ounce", "fluid ounces", "fl oz"},
		ToBase:  0.0295735,
	},
	{
		Code:    "tbsp",
		Symbol:  "tbsp",
		Name:    "tablespoon",
		Plural:  "tablespoons",
		Type:    UnitTypeVolume,
		Aliases: []string{"tablespoon", "tablespoons"},
		ToBase:  0.0147868,
	},
	{
		Code:    "tsp",
		Symbol:  "tsp",
		Name:    "teaspoon",
		Plural:  "teaspoons",
		Type:    UnitTypeVolume,
		Aliases: []string{"teaspoon", "teaspoons"},
		ToBase:  0.00492892,
	},
	{
		Code:    "m3",
		Symbol:  "m³",
		Name:    "cubic meter",
		Plural:  "cubic meters",
		Type:    UnitTypeVolume,
		Aliases: []string{"cubic meter", "cubic meters", "cubic metre", "cubic metres"},
		ToBase:  1000.0,
	},
}

// ════════════════════════════════════════════════════════════════
// PUBLIC API
// ════════════════════════════════════════════════════════════════

// LookupUnit finds a unit by code or alias.
// Returns nil if not found.
func LookupUnit(s string) *Unit {
	return units.Lookup(s)
}

// ParseUnit parses a string into a unit.
// Accepts codes ("km"), symbols, or natural names ("kilometers").
// Returns nil if not found.
func ParseUnit(s string) *Unit {
	return units.Lookup(strings.TrimSpace(s))
}

// IsUnit checks if a string refers to a known unit.
func IsUnit(s string) bool {
	return units.Lookup(s) != nil
}

// IsUnitCode checks if a string is a unit code.
func IsUnitCode(code string) bool {
	return units.byCode[code] != nil || units.byCode[strings.ToLower(code)] != nil
}

// AllUnits returns all curated units.
func AllUnits() []Unit {
	return curatedUnits
}

// UnitsByType returns all units of a given type.
func UnitsByType(t UnitType) []*Unit {
	return units.byType[t]
}

// UnitCodes returns all unit codes.
func UnitCodes() []string {
	codes := make([]string, len(curatedUnits))
	for i, u := range curatedUnits {
		codes[i] = u.Code
	}
	return codes
}

// BaseUnit returns the base unit for a given unit type.
func BaseUnit(t UnitType) *Unit {
	for _, u := range units.byType[t] {
		if u.IsBase {
			return u
		}
	}
	return nil
}

// ConvertUnit converts a value from one unit to another.
// Returns the converted value and true if successful.
func ConvertUnit(value float64, from, to string) (float64, bool) {
	fromUnit := units.Lookup(from)
	toUnit := units.Lookup(to)

	if fromUnit == nil || toUnit == nil {
		return 0, false
	}

	return fromUnit.ConvertTo(value, toUnit)
}

// CompatibleUnits checks if two units can be converted between each other.
func CompatibleUnits(a, b string) bool {
	unitA := units.Lookup(a)
	unitB := units.Lookup(b)

	if unitA == nil || unitB == nil {
		return false
	}

	return unitA.Type == unitB.Type
}
