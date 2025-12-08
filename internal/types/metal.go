// internal/types/metal.go

package types

import (
	"strings"
)

// Metal represents a precious metal.
type Metal struct {
	Code      string   // ISO 4217 code: "XAU", "XAG"
	Symbol    string   // Display symbol
	Name      string   // Full name: "Gold", "Silver"
	Aliases   []string // Natural language aliases
	UnitName  string   // Standard unit: "oz", "gram"
	UnitLabel string   // Display label: "per troy oz", "per gram"
}

// String returns the metal code.
func (m Metal) String() string {
	return m.Code
}

// MetalRegistry holds all known metals.
type MetalRegistry struct {
	byCode  map[string]*Metal
	byAlias map[string]*Metal
}

// Global metal registry.
var metals = newMetalRegistry()

// newMetalRegistry creates and populates the metal registry.
func newMetalRegistry() *MetalRegistry {
	r := &MetalRegistry{
		byCode:  make(map[string]*Metal),
		byAlias: make(map[string]*Metal),
	}

	for i := range curatedMetals {
		r.register(&curatedMetals[i])
	}

	return r
}

// register adds a metal to the registry.
func (r *MetalRegistry) register(m *Metal) {
	// By code (case-insensitive)
	r.byCode[strings.ToUpper(m.Code)] = m
	r.byCode[strings.ToLower(m.Code)] = m

	// By aliases (case-insensitive)
	for _, alias := range m.Aliases {
		r.byAlias[strings.ToLower(alias)] = m
	}
}

// Lookup finds a metal by code or alias.
func (r *MetalRegistry) Lookup(s string) *Metal {
	// Try code first (case-insensitive)
	if m, ok := r.byCode[strings.ToUpper(s)]; ok {
		return m
	}

	// Try alias (case-insensitive)
	if m, ok := r.byAlias[strings.ToLower(s)]; ok {
		return m
	}

	return nil
}

// curatedMetals contains precious metals with ISO 4217 codes.
// These codes are used by financial APIs for metal prices.
var curatedMetals = []Metal{
	// ════════════════════════════════════════════════════════════
	// PRECIOUS METALS (ISO 4217)
	// ════════════════════════════════════════════════════════════
	{
		Code:      "XAU",
		Symbol:    "Au",
		Name:      "Gold",
		Aliases:   []string{"gold", "au", "xau"},
		UnitName:  "oz",
		UnitLabel: "per troy oz",
	},
	{
		Code:      "XAG",
		Symbol:    "Ag",
		Name:      "Silver",
		Aliases:   []string{"silver", "ag", "xag"},
		UnitName:  "oz",
		UnitLabel: "per troy oz",
	},
	{
		Code:      "XPT",
		Symbol:    "Pt",
		Name:      "Platinum",
		Aliases:   []string{"platinum", "pt", "xpt"},
		UnitName:  "oz",
		UnitLabel: "per troy oz",
	},
	{
		Code:      "XPD",
		Symbol:    "Pd",
		Name:      "Palladium",
		Aliases:   []string{"palladium", "pd", "xpd"},
		UnitName:  "oz",
		UnitLabel: "per troy oz",
	},

	// ════════════════════════════════════════════════════════════
	// OTHER METALS (commonly traded)
	// ════════════════════════════════════════════════════════════
	{
		Code:      "XCU",
		Symbol:    "Cu",
		Name:      "Copper",
		Aliases:   []string{"copper", "cu", "xcu"},
		UnitName:  "lb",
		UnitLabel: "per pound",
	},
	{
		Code:      "XAL",
		Symbol:    "Al",
		Name:      "Aluminum",
		Aliases:   []string{"aluminum", "aluminium", "al", "xal"},
		UnitName:  "lb",
		UnitLabel: "per pound",
	},
	{
		Code:      "XNI",
		Symbol:    "Ni",
		Name:      "Nickel",
		Aliases:   []string{"nickel", "ni", "xni"},
		UnitName:  "lb",
		UnitLabel: "per pound",
	},
	{
		Code:      "XZN",
		Symbol:    "Zn",
		Name:      "Zinc",
		Aliases:   []string{"zinc", "zn", "xzn"},
		UnitName:  "lb",
		UnitLabel: "per pound",
	},
	{
		Code:      "XPB",
		Symbol:    "Pb",
		Name:      "Lead",
		Aliases:   []string{"lead", "pb", "xpb"},
		UnitName:  "lb",
		UnitLabel: "per pound",
	},
	{
		Code:      "XSN",
		Symbol:    "Sn",
		Name:      "Tin",
		Aliases:   []string{"tin", "sn", "xsn"},
		UnitName:  "lb",
		UnitLabel: "per pound",
	},
}

// ════════════════════════════════════════════════════════════════
// PUBLIC API
// ════════════════════════════════════════════════════════════════

// LookupMetal finds a metal by code or alias.
// Returns nil if not found.
func LookupMetal(s string) *Metal {
	return metals.Lookup(s)
}

// ParseMetal parses a string into a metal.
// Accepts codes ("XAU"), symbols ("Au"), or natural names ("gold").
// Returns nil if not found.
func ParseMetal(s string) *Metal {
	return metals.Lookup(strings.TrimSpace(s))
}

// IsMetal checks if a string refers to a known metal.
func IsMetal(s string) bool {
	return metals.Lookup(s) != nil
}

// IsMetalCode checks if a string is a metal ISO code.
func IsMetalCode(code string) bool {
	return metals.byCode[strings.ToUpper(code)] != nil
}

// AllMetals returns all curated metals.
func AllMetals() []Metal {
	return curatedMetals
}

// MetalCodes returns all metal codes.
func MetalCodes() []string {
	codes := make([]string, len(curatedMetals))
	for i, m := range curatedMetals {
		codes[i] = m.Code
	}
	return codes
}

// PreciousMetals returns only the precious metals (XAU, XAG, XPT, XPD).
func PreciousMetals() []Metal {
	precious := make([]Metal, 0, 4)
	for _, m := range curatedMetals {
		switch m.Code {
		case "XAU", "XAG", "XPT", "XPD":
			precious = append(precious, m)
		}
	}
	return precious
}