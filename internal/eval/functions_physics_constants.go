// internal/eval/functions_physics_constants.go

package eval

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/0xsj/numio/pkg/types"
)

// PhysicsConstant represents a physical constant with metadata.
type PhysicsConstant struct {
	Value       float64
	Unit        string
	Description string
	Symbol      string
}

// physicsConstants holds all physical constants.
var physicsConstants = map[string]PhysicsConstant{
	// Fundamental constants
	"c": {
		Value:       299792458,
		Unit:        "m/s",
		Description: "Speed of light in vacuum",
		Symbol:      "c",
	},
	"h": {
		Value:       6.62607015e-34,
		Unit:        "J·s",
		Description: "Planck constant",
		Symbol:      "h",
	},
	"hbar": {
		Value:       1.054571817e-34,
		Unit:        "J·s",
		Description: "Reduced Planck constant",
		Symbol:      "ℏ",
	},
	"g_earth": {
		Value:       9.80665,
		Unit:        "m/s²",
		Description: "Standard gravity on Earth",
		Symbol:      "g",
	},
	"g_gravity": {
		Value:       6.67430e-11,
		Unit:        "N·m²/kg²",
		Description: "Gravitational constant",
		Symbol:      "G",
	},
	"e_charge": {
		Value:       1.602176634e-19,
		Unit:        "C",
		Description: "Elementary charge",
		Symbol:      "e",
	},
	"me": {
		Value:       9.1093837015e-31,
		Unit:        "kg",
		Description: "Electron mass",
		Symbol:      "mₑ",
	},
	"mp": {
		Value:       1.67262192369e-27,
		Unit:        "kg",
		Description: "Proton mass",
		Symbol:      "mₚ",
	},
	"mn": {
		Value:       1.67492749804e-27,
		Unit:        "kg",
		Description: "Neutron mass",
		Symbol:      "mₙ",
	},
	"kb": {
		Value:       1.380649e-23,
		Unit:        "J/K",
		Description: "Boltzmann constant",
		Symbol:      "kB",
	},
	"na": {
		Value:       6.02214076e23,
		Unit:        "mol⁻¹",
		Description: "Avogadro number",
		Symbol:      "Nₐ",
	},
	"r_gas": {
		Value:       8.314462618,
		Unit:        "J/(mol·K)",
		Description: "Ideal gas constant",
		Symbol:      "R",
	},
	"epsilon0": {
		Value:       8.8541878128e-12,
		Unit:        "F/m",
		Description: "Vacuum permittivity",
		Symbol:      "ε₀",
	},
	"mu0": {
		Value:       1.25663706212e-6,
		Unit:        "H/m",
		Description: "Vacuum permeability",
		Symbol:      "μ₀",
	},
	"ke": {
		Value:       8.9875517923e9,
		Unit:        "N·m²/C²",
		Description: "Coulomb constant",
		Symbol:      "kₑ",
	},
	"sigma_sb": {
		Value:       5.670374419e-8,
		Unit:        "W/(m²·K⁴)",
		Description: "Stefan-Boltzmann constant",
		Symbol:      "σ",
	},
	"wien": {
		Value:       2.897771955e-3,
		Unit:        "m·K",
		Description: "Wien displacement constant",
		Symbol:      "b",
	},
	"r_inf": {
		Value:       1.0973731568160e7,
		Unit:        "m⁻¹",
		Description: "Rydberg constant",
		Symbol:      "R∞",
	},
	"a0": {
		Value:       5.29177210903e-11,
		Unit:        "m",
		Description: "Bohr radius",
		Symbol:      "a₀",
	},
	"alpha": {
		Value:       7.2973525693e-3,
		Unit:        "",
		Description: "Fine-structure constant",
		Symbol:      "α",
	},
	"ev": {
		Value:       1.602176634e-19,
		Unit:        "J",
		Description: "Electron volt",
		Symbol:      "eV",
	},
	"amu": {
		Value:       1.66053906660e-27,
		Unit:        "kg",
		Description: "Atomic mass unit",
		Symbol:      "u",
	},
	// Astronomical constants
	"ly": {
		Value:       9.4607304725808e15,
		Unit:        "m",
		Description: "Light year",
		Symbol:      "ly",
	},
	"pc": {
		Value:       3.0856775814913673e16,
		Unit:        "m",
		Description: "Parsec",
		Symbol:      "pc",
	},
	"au": {
		Value:       1.495978707e11,
		Unit:        "m",
		Description: "Astronomical unit",
		Symbol:      "AU",
	},
}

// GetPhysicsConstant retrieves a physics constant by key.
func GetPhysicsConstant(key string) (PhysicsConstant, bool) {
	c, ok := physicsConstants[key]
	return c, ok
}

// IsPhysicsConstant checks if a name is a physics constant.
func IsPhysicsConstant(name string) bool {
	_, ok := physicsConstants[strings.ToLower(name)]
	return ok
}

// ════════════════════════════════════════════════════════════════
// CONSTANT LOOKUP FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnPhysConst looks up a physics constant by name.
func FnPhysConst(args []types.Value) types.Value {
	name := strings.ToLower(args[0].AsString())

	c, ok := physicsConstants[name]
	if !ok {
		return types.Errorf("unknown physics constant: %s", name)
	}

	return types.Number(c.Value)
}

// FnPhysConstInfo returns info about a physics constant.
func FnPhysConstInfo(args []types.Value) types.Value {
	name := strings.ToLower(args[0].AsString())

	c, ok := physicsConstants[name]
	if !ok {
		return types.Errorf("unknown physics constant: %s", name)
	}

	info := fmt.Sprintf("%s (%s) = %g %s - %s", name, c.Symbol, c.Value, c.Unit, c.Description)
	return types.StringValue(info)
}

// FnListPhysConsts lists all available physics constants.
func FnListPhysConsts(args []types.Value) types.Value {
	names := make([]string, 0, len(physicsConstants))
	for name := range physicsConstants {
		names = append(names, name)
	}
	sort.Strings(names)

	var sb strings.Builder
	for _, name := range names {
		c := physicsConstants[name]
		sb.WriteString(fmt.Sprintf("%s (%s): %g %s\n", name, c.Symbol, c.Value, c.Unit))
	}

	return types.StringValue(sb.String())
}

// ════════════════════════════════════════════════════════════════
// CONSTANT ACCESSOR FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnSpeedOfLight returns the speed of light.
func FnSpeedOfLight(args []types.Value) types.Value {
	return types.Number(physicsConstants["c"].Value)
}

// FnPlanckConstant returns Planck's constant.
func FnPlanckConstant(args []types.Value) types.Value {
	return types.Number(physicsConstants["h"].Value)
}

// FnReducedPlanck returns the reduced Planck constant.
func FnReducedPlanck(args []types.Value) types.Value {
	return types.Number(physicsConstants["hbar"].Value)
}

// FnGravityEarth returns standard gravity on Earth.
func FnGravityEarth(args []types.Value) types.Value {
	return types.Number(physicsConstants["g_earth"].Value)
}

// FnGravitationalConstant returns the gravitational constant G.
func FnGravitationalConstant(args []types.Value) types.Value {
	return types.Number(physicsConstants["g_gravity"].Value)
}

// FnElementaryCharge returns the elementary charge.
func FnElementaryCharge(args []types.Value) types.Value {
	return types.Number(physicsConstants["e_charge"].Value)
}

// FnElectronMass returns the electron mass.
func FnElectronMass(args []types.Value) types.Value {
	return types.Number(physicsConstants["me"].Value)
}

// FnProtonMass returns the proton mass.
func FnProtonMass(args []types.Value) types.Value {
	return types.Number(physicsConstants["mp"].Value)
}

// FnNeutronMass returns the neutron mass.
func FnNeutronMass(args []types.Value) types.Value {
	return types.Number(physicsConstants["mn"].Value)
}

// FnBoltzmannConstant returns the Boltzmann constant.
func FnBoltzmannConstant(args []types.Value) types.Value {
	return types.Number(physicsConstants["kb"].Value)
}

// FnAvogadroNumber returns Avogadro's number.
func FnAvogadroNumber(args []types.Value) types.Value {
	return types.Number(physicsConstants["na"].Value)
}

// FnGasConstant returns the ideal gas constant.
func FnGasConstant(args []types.Value) types.Value {
	return types.Number(physicsConstants["r_gas"].Value)
}

// FnVacuumPermittivity returns the vacuum permittivity.
func FnVacuumPermittivity(args []types.Value) types.Value {
	return types.Number(physicsConstants["epsilon0"].Value)
}

// FnVacuumPermeability returns the vacuum permeability.
func FnVacuumPermeability(args []types.Value) types.Value {
	return types.Number(physicsConstants["mu0"].Value)
}

// FnCoulombConstant returns Coulomb's constant.
func FnCoulombConstant(args []types.Value) types.Value {
	return types.Number(physicsConstants["ke"].Value)
}

// FnStefanBoltzmannConst returns the Stefan-Boltzmann constant.
func FnStefanBoltzmannConst(args []types.Value) types.Value {
	return types.Number(physicsConstants["sigma_sb"].Value)
}

// FnWienConstant returns the Wien displacement constant.
func FnWienConstant(args []types.Value) types.Value {
	return types.Number(physicsConstants["wien"].Value)
}

// FnRydbergConstant returns the Rydberg constant.
func FnRydbergConstant(args []types.Value) types.Value {
	return types.Number(physicsConstants["r_inf"].Value)
}

// FnBohrRadiusConst returns the Bohr radius constant.
func FnBohrRadiusConst(args []types.Value) types.Value {
	return types.Number(physicsConstants["a0"].Value)
}

// FnFineStructure returns the fine-structure constant.
func FnFineStructure(args []types.Value) types.Value {
	return types.Number(physicsConstants["alpha"].Value)
}

// FnElectronVolt returns the electron volt in joules.
func FnElectronVolt(args []types.Value) types.Value {
	return types.Number(physicsConstants["ev"].Value)
}

// FnAtomicMassUnit returns the atomic mass unit.
func FnAtomicMassUnit(args []types.Value) types.Value {
	return types.Number(physicsConstants["amu"].Value)
}

// ════════════════════════════════════════════════════════════════
// DERIVED CONSTANT FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnComptonWavelength returns the Compton wavelength of the electron.
// λ_C = h/(m_e c)
func FnComptonWavelength(args []types.Value) types.Value {
	h := physicsConstants["h"].Value
	me := physicsConstants["me"].Value
	c := physicsConstants["c"].Value
	return types.Number(h / (me * c))
}

// FnBohrMagneton returns the Bohr magneton.
// μ_B = eℏ/(2m_e)
func FnBohrMagneton(args []types.Value) types.Value {
	e := physicsConstants["e_charge"].Value
	hbar := physicsConstants["hbar"].Value
	me := physicsConstants["me"].Value
	return types.Number(e * hbar / (2 * me))
}

// FnNuclearMagneton returns the nuclear magneton.
// μ_N = eℏ/(2m_p)
func FnNuclearMagneton(args []types.Value) types.Value {
	e := physicsConstants["e_charge"].Value
	hbar := physicsConstants["hbar"].Value
	mp := physicsConstants["mp"].Value
	return types.Number(e * hbar / (2 * mp))
}

// FnClassicalElectronRadius returns the classical electron radius.
// r_e = k_e e²/(m_e c²)
func FnClassicalElectronRadius(args []types.Value) types.Value {
	ke := physicsConstants["ke"].Value
	e := physicsConstants["e_charge"].Value
	me := physicsConstants["me"].Value
	c := physicsConstants["c"].Value
	return types.Number(ke * e * e / (me * c * c))
}

// FnPlanckLength returns the Planck length.
// l_P = √(ℏG/c³)
func FnPlanckLength(args []types.Value) types.Value {
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	return types.Number(math.Sqrt(hbar * G / (c * c * c)))
}

// FnPlanckTime returns the Planck time.
// t_P = √(ℏG/c⁵)
func FnPlanckTime(args []types.Value) types.Value {
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	return types.Number(math.Sqrt(hbar * G / (c * c * c * c * c)))
}

// FnPlanckMass returns the Planck mass.
// m_P = √(ℏc/G)
func FnPlanckMass(args []types.Value) types.Value {
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	return types.Number(math.Sqrt(hbar * c / G))
}

// FnPlanckEnergy returns the Planck energy.
// E_P = √(ℏc⁵/G)
func FnPlanckEnergy(args []types.Value) types.Value {
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	return types.Number(math.Sqrt(hbar * c * c * c * c * c / G))
}

// FnPlanckTemperature returns the Planck temperature.
// T_P = √(ℏc⁵/(Gk_B²))
func FnPlanckTemperature(args []types.Value) types.Value {
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	kb := physicsConstants["kb"].Value
	return types.Number(math.Sqrt(hbar * c * c * c * c * c / (G * kb * kb)))
}

// FnImpedanceOfFreeSpace returns the impedance of free space.
// Z_0 = μ_0 c ≈ 376.73 Ω
func FnImpedanceOfFreeSpace(args []types.Value) types.Value {
	mu0 := physicsConstants["mu0"].Value
	c := physicsConstants["c"].Value
	return types.Number(mu0 * c)
}

// FnMagneticFluxQuantum returns the magnetic flux quantum.
// Φ_0 = h/(2e)
func FnMagneticFluxQuantum(args []types.Value) types.Value {
	h := physicsConstants["h"].Value
	e := physicsConstants["e_charge"].Value
	return types.Number(h / (2 * e))
}

// FnConductanceQuantum returns the conductance quantum.
// G_0 = 2e²/h
func FnConductanceQuantum(args []types.Value) types.Value {
	h := physicsConstants["h"].Value
	e := physicsConstants["e_charge"].Value
	return types.Number(2 * e * e / h)
}

// FnThomsonCrossSection returns the Thomson cross section.
// σ_T = (8π/3) r_e²
func FnThomsonCrossSection(args []types.Value) types.Value {
	ke := physicsConstants["ke"].Value
	e := physicsConstants["e_charge"].Value
	me := physicsConstants["me"].Value
	c := physicsConstants["c"].Value

	re := ke * e * e / (me * c * c) // Classical electron radius
	return types.Number((8.0 * math.Pi / 3.0) * re * re)
}

// ════════════════════════════════════════════════════════════════
// UNIT CONVERSION FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnEvToJoules converts electron volts to joules.
func FnEvToJoules(args []types.Value) types.Value {
	ev := args[0].AsFloat()
	return types.Number(ev * physicsConstants["ev"].Value)
}

// FnJoulesToEv converts joules to electron volts.
func FnJoulesToEv(args []types.Value) types.Value {
	j := args[0].AsFloat()
	return types.Number(j / physicsConstants["ev"].Value)
}

// FnEvToKelvin converts electron volts to Kelvin.
// T = E / k_B
func FnEvToKelvin(args []types.Value) types.Value {
	ev := args[0].AsFloat()
	eJ := ev * physicsConstants["ev"].Value
	kb := physicsConstants["kb"].Value
	return types.Number(eJ / kb)
}

// FnKelvinToEv converts Kelvin to electron volts.
// E = k_B T
func FnKelvinToEv(args []types.Value) types.Value {
	T := args[0].AsFloat()
	kb := physicsConstants["kb"].Value
	ev := physicsConstants["ev"].Value
	return types.Number(kb * T / ev)
}

// FnAmuToKg converts atomic mass units to kilograms.
func FnAmuToKg(args []types.Value) types.Value {
	amu := args[0].AsFloat()
	return types.Number(amu * physicsConstants["amu"].Value)
}

// FnKgToAmu converts kilograms to atomic mass units.
func FnKgToAmu(args []types.Value) types.Value {
	kg := args[0].AsFloat()
	return types.Number(kg / physicsConstants["amu"].Value)
}

// FnLyToM converts light years to meters.
func FnLyToM(args []types.Value) types.Value {
	ly := args[0].AsFloat()
	return types.Number(ly * physicsConstants["ly"].Value)
}

// FnMToLy converts meters to light years.
func FnMToLy(args []types.Value) types.Value {
	m := args[0].AsFloat()
	return types.Number(m / physicsConstants["ly"].Value)
}

// FnPcToM converts parsecs to meters.
func FnPcToM(args []types.Value) types.Value {
	pc := args[0].AsFloat()
	return types.Number(pc * physicsConstants["pc"].Value)
}

// FnMToPc converts meters to parsecs.
func FnMToPc(args []types.Value) types.Value {
	m := args[0].AsFloat()
	return types.Number(m / physicsConstants["pc"].Value)
}

// FnAuToM converts astronomical units to meters.
func FnAuToM(args []types.Value) types.Value {
	au := args[0].AsFloat()
	return types.Number(au * physicsConstants["au"].Value)
}

// FnMToAu converts meters to astronomical units.
func FnMToAu(args []types.Value) types.Value {
	m := args[0].AsFloat()
	return types.Number(m / physicsConstants["au"].Value)
}
