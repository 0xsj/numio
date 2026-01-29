// internal/eval/functions_physics_constants.go

package eval

import (
	"math"

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
}

// GetPhysicsConstant retrieves a physics constant by key.
func GetPhysicsConstant(key string) (PhysicsConstant, bool) {
	c, ok := physicsConstants[key]
	return c, ok
}

// ════════════════════════════════════════════════════════════════
// CONSTANT ACCESSOR FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnSpeedOfLight returns the speed of light.
func FnSpeedOfLight(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("c requires no arguments")
	}
	return types.Number(physicsConstants["c"].Value)
}

// FnPlanckConstant returns Planck's constant.
func FnPlanckConstant(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("planck requires no arguments")
	}
	return types.Number(physicsConstants["h"].Value)
}

// FnReducedPlanck returns the reduced Planck constant.
func FnReducedPlanck(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("hbar requires no arguments")
	}
	return types.Number(physicsConstants["hbar"].Value)
}

// FnGravityEarth returns standard gravity on Earth.
func FnGravityEarth(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("g requires no arguments")
	}
	return types.Number(physicsConstants["g_earth"].Value)
}

// FnGravitationalConstant returns the gravitational constant G.
func FnGravitationalConstant(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("G requires no arguments")
	}
	return types.Number(physicsConstants["g_gravity"].Value)
}

// FnElementaryCharge returns the elementary charge.
func FnElementaryCharge(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("e requires no arguments")
	}
	return types.Number(physicsConstants["e_charge"].Value)
}

// FnElectronMass returns the electron mass.
func FnElectronMass(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("me requires no arguments")
	}
	return types.Number(physicsConstants["me"].Value)
}

// FnProtonMass returns the proton mass.
func FnProtonMass(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("mp requires no arguments")
	}
	return types.Number(physicsConstants["mp"].Value)
}

// FnNeutronMass returns the neutron mass.
func FnNeutronMass(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("mn requires no arguments")
	}
	return types.Number(physicsConstants["mn"].Value)
}

// FnBoltzmannConstant returns the Boltzmann constant.
func FnBoltzmannConstant(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("kb requires no arguments")
	}
	return types.Number(physicsConstants["kb"].Value)
}

// FnAvogadroNumber returns Avogadro's number.
func FnAvogadroNumber(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("na requires no arguments")
	}
	return types.Number(physicsConstants["na"].Value)
}

// FnGasConstant returns the ideal gas constant.
func FnGasConstant(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("rgas requires no arguments")
	}
	return types.Number(physicsConstants["r_gas"].Value)
}

// FnVacuumPermittivity returns the vacuum permittivity.
func FnVacuumPermittivity(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("epsilon0 requires no arguments")
	}
	return types.Number(physicsConstants["epsilon0"].Value)
}

// FnVacuumPermeability returns the vacuum permeability.
func FnVacuumPermeability(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("mu0 requires no arguments")
	}
	return types.Number(physicsConstants["mu0"].Value)
}

// FnCoulombConstant returns Coulomb's constant.
func FnCoulombConstant(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("ke requires no arguments")
	}
	return types.Number(physicsConstants["ke"].Value)
}

// FnStefanBoltzmannConst returns the Stefan-Boltzmann constant.
func FnStefanBoltzmannConst(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("stefanboltz requires no arguments")
	}
	return types.Number(physicsConstants["sigma_sb"].Value)
}

// FnWienConstant returns the Wien displacement constant.
func FnWienConstant(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("wien requires no arguments")
	}
	return types.Number(physicsConstants["wien"].Value)
}

// FnRydbergConstant returns the Rydberg constant.
func FnRydbergConstant(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("rydberg requires no arguments")
	}
	return types.Number(physicsConstants["r_inf"].Value)
}

// FnBohrRadiusConst returns the Bohr radius constant.
func FnBohrRadiusConst(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("a0 requires no arguments")
	}
	return types.Number(physicsConstants["a0"].Value)
}

// FnFineStructure returns the fine-structure constant.
func FnFineStructure(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("alpha requires no arguments")
	}
	return types.Number(physicsConstants["alpha"].Value)
}

// FnElectronVolt returns the electron volt in joules.
func FnElectronVolt(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("ev requires no arguments")
	}
	return types.Number(physicsConstants["ev"].Value)
}

// FnAtomicMassUnit returns the atomic mass unit.
func FnAtomicMassUnit(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("amu requires no arguments")
	}
	return types.Number(physicsConstants["amu"].Value)
}

// ════════════════════════════════════════════════════════════════
// DERIVED CONSTANT FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnComptonWavelength returns the Compton wavelength of the electron.
// λ_C = h/(m_e c)
func FnComptonWavelength(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("comptonwavelength requires no arguments")
	}
	h := physicsConstants["h"].Value
	me := physicsConstants["me"].Value
	c := physicsConstants["c"].Value
	return types.Number(h / (me * c))
}

// FnBohrMagneton returns the Bohr magneton.
// μ_B = eℏ/(2m_e)
func FnBohrMagneton(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("bohrmagneton requires no arguments")
	}
	e := physicsConstants["e_charge"].Value
	hbar := physicsConstants["hbar"].Value
	me := physicsConstants["me"].Value
	return types.Number(e * hbar / (2 * me))
}

// FnNuclearMagneton returns the nuclear magneton.
// μ_N = eℏ/(2m_p)
func FnNuclearMagneton(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("nuclearmagneton requires no arguments")
	}
	e := physicsConstants["e_charge"].Value
	hbar := physicsConstants["hbar"].Value
	mp := physicsConstants["mp"].Value
	return types.Number(e * hbar / (2 * mp))
}

// FnClassicalElectronRadius returns the classical electron radius.
// r_e = k_e e²/(m_e c²)
func FnClassicalElectronRadius(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("classicalelectronradius requires no arguments")
	}
	ke := physicsConstants["ke"].Value
	e := physicsConstants["e_charge"].Value
	me := physicsConstants["me"].Value
	c := physicsConstants["c"].Value
	return types.Number(ke * e * e / (me * c * c))
}

// FnPlanckLength returns the Planck length.
// l_P = √(ℏG/c³)
func FnPlanckLength(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("plancklength requires no arguments")
	}
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	return types.Number(math.Sqrt(hbar * G / (c * c * c)))
}

// FnPlanckTime returns the Planck time.
// t_P = √(ℏG/c⁵)
func FnPlanckTime(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("plancktime requires no arguments")
	}
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	return types.Number(math.Sqrt(hbar * G / (c * c * c * c * c)))
}

// FnPlanckMass returns the Planck mass.
// m_P = √(ℏc/G)
func FnPlanckMass(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("planckmass requires no arguments")
	}
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	return types.Number(math.Sqrt(hbar * c / G))
}

// FnPlanckEnergy returns the Planck energy.
// E_P = √(ℏc⁵/G)
func FnPlanckEnergy(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("planckenergy requires no arguments")
	}
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	return types.Number(math.Sqrt(hbar * c * c * c * c * c / G))
}

// FnPlanckTemperature returns the Planck temperature.
// T_P = √(ℏc⁵/(Gk_B²))
func FnPlanckTemperature(args []types.Value) types.Value {
	if len(args) != 0 {
		return types.Error("plancktemperature requires no arguments")
	}
	hbar := physicsConstants["hbar"].Value
	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value
	kb := physicsConstants["kb"].Value
	return types.Number(math.Sqrt(hbar * c * c * c * c * c / (G * kb * kb)))
}
