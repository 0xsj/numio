// internal/eval/functions_physics_quantum.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// PHOTONS & PHOTOELECTRIC EFFECT
// ════════════════════════════════════════════════════════════════

// FnPhotonEnergy calculates energy of a photon.
// E = hf = hc/λ
// Args: frequency (Hz) OR wavelength (m), type ("f" or "λ")
func FnPhotonEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("photonenergy requires 2 arguments: value, type")
	}

	val := args[0].AsFloat()
	calcType := args[1].AsString()

	h := physicsConstants["h"].Value
	c := physicsConstants["c"].Value

	var E float64
	if calcType == "f" || calcType == "frequency" {
		if val <= 0 {
			return types.Error("photonenergy: frequency must be positive")
		}
		E = h * val
	} else if calcType == "λ" || calcType == "wavelength" || calcType == "w" {
		if val <= 0 {
			return types.Error("photonenergy: wavelength must be positive")
		}
		E = h * c / val
	} else {
		return types.Error("photonenergy: type must be 'f' or 'λ'")
	}

	return types.Number(E)
}

// FnPhotonMomentum calculates momentum of a photon.
// p = h/λ = E/c
// Args: wavelength (m) OR energy (J), type ("λ" or "E")
func FnPhotonMomentum(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("photonmomentum requires 2 arguments: value, type")
	}

	val := args[0].AsFloat()
	calcType := args[1].AsString()

	h := physicsConstants["h"].Value
	c := physicsConstants["c"].Value

	var p float64
	if calcType == "λ" || calcType == "wavelength" || calcType == "w" {
		if val <= 0 {
			return types.Error("photonmomentum: wavelength must be positive")
		}
		p = h / val
	} else if calcType == "E" || calcType == "energy" {
		if val < 0 {
			return types.Error("photonmomentum: energy cannot be negative")
		}
		p = val / c
	} else {
		return types.Error("photonmomentum: type must be 'λ' or 'E'")
	}

	return types.Number(p)
}

// FnPhotoelectricKE calculates maximum kinetic energy of photoelectrons.
// KE_max = hf - φ = hc/λ - φ
// Args: photon frequency (Hz) or wavelength (m), work function (J), type ("f" or "λ")
func FnPhotoelectricKE(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("photoelectricke requires 3 arguments: value, work_function, type")
	}

	val := args[0].AsFloat()
	phi := args[1].AsFloat()
	calcType := args[2].AsString()

	h := physicsConstants["h"].Value
	c := physicsConstants["c"].Value

	var E float64
	if calcType == "f" || calcType == "frequency" {
		E = h * val
	} else if calcType == "λ" || calcType == "wavelength" || calcType == "w" {
		if val <= 0 {
			return types.Error("photoelectricke: wavelength must be positive")
		}
		E = h * c / val
	} else {
		return types.Error("photoelectricke: type must be 'f' or 'λ'")
	}

	KE := E - phi
	if KE < 0 {
		return types.Error("photoelectricke: photon energy below work function (no emission)")
	}

	return types.Number(KE)
}

// FnWorkFunction calculates work function from threshold frequency.
// φ = hf₀
// Args: threshold frequency (Hz)
func FnWorkFunction(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("workfunction requires 1 argument: threshold_frequency")
	}

	f0 := args[0].AsFloat()

	if f0 <= 0 {
		return types.Error("workfunction: frequency must be positive")
	}

	h := physicsConstants["h"].Value
	return types.Number(h * f0)
}

// FnThresholdFrequency calculates threshold frequency from work function.
// f₀ = φ/h
// Args: work function (J)
func FnThresholdFrequency(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("thresholdfreq requires 1 argument: work_function")
	}

	phi := args[0].AsFloat()

	if phi <= 0 {
		return types.Error("thresholdfreq: work function must be positive")
	}

	h := physicsConstants["h"].Value
	return types.Number(phi / h)
}

// FnThresholdWavelength calculates threshold wavelength from work function.
// λ₀ = hc/φ
// Args: work function (J)
func FnThresholdWavelength(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("thresholdwavelength requires 1 argument: work_function")
	}

	phi := args[0].AsFloat()

	if phi <= 0 {
		return types.Error("thresholdwavelength: work function must be positive")
	}

	h := physicsConstants["h"].Value
	c := physicsConstants["c"].Value

	return types.Number(h * c / phi)
}

// FnStoppingPotential calculates stopping potential.
// V_s = KE_max / e = (hf - φ) / e
// Args: maximum kinetic energy (J)
func FnStoppingPotential(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("stoppingpotential requires 1 argument: KE_max")
	}

	KE := args[0].AsFloat()

	if KE < 0 {
		return types.Error("stoppingpotential: kinetic energy cannot be negative")
	}

	e := physicsConstants["e_charge"].Value
	return types.Number(KE / e)
}

// ════════════════════════════════════════════════════════════════
// DE BROGLIE WAVELENGTH
// ════════════════════════════════════════════════════════════════

// FnDeBroglieWavelength calculates de Broglie wavelength.
// λ = h/p = h/(mv)
// Args: mass (kg), velocity (m/s)
func FnDeBroglieWavelength(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("debrogliewavelength requires 2 arguments: mass, velocity")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()

	if m*v == 0 {
		return types.Error("debrogliewavelength: momentum cannot be zero")
	}

	h := physicsConstants["h"].Value
	return types.Number(h / (m * v))
}

// FnDeBroglieFromKE calculates de Broglie wavelength from kinetic energy.
// λ = h / √(2mKE)
// Args: mass (kg), kinetic energy (J)
func FnDeBroglieFromKE(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("debrogliefromke requires 2 arguments: mass, kinetic_energy")
	}

	m := args[0].AsFloat()
	KE := args[1].AsFloat()

	if m <= 0 {
		return types.Error("debrogliefromke: mass must be positive")
	}
	if KE < 0 {
		return types.Error("debrogliefromke: kinetic energy cannot be negative")
	}
	if KE == 0 {
		return types.Error("debrogliefromke: kinetic energy cannot be zero")
	}

	h := physicsConstants["h"].Value
	lambda := h / math.Sqrt(2*m*KE)

	return types.Number(lambda)
}

// FnElectronWavelength calculates de Broglie wavelength of electron accelerated through voltage.
// λ = h / √(2meV)
// Args: accelerating voltage (V)
func FnElectronWavelength(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("electronwavelength requires 1 argument: voltage")
	}

	V := args[0].AsFloat()

	if V <= 0 {
		return types.Error("electronwavelength: voltage must be positive")
	}

	h := physicsConstants["h"].Value
	me := physicsConstants["me"].Value
	e := physicsConstants["e_charge"].Value

	lambda := h / math.Sqrt(2*me*e*V)
	return types.Number(lambda)
}

// ════════════════════════════════════════════════════════════════
// COMPTON SCATTERING
// ════════════════════════════════════════════════════════════════

// FnComptonShift calculates Compton wavelength shift.
// Δλ = (h/mc)(1 - cos(θ))
// Args: scattering angle (degrees)
func FnComptonShift(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("comptonshift requires 1 argument: angle_deg")
	}

	theta := args[0].AsFloat() * math.Pi / 180

	h := physicsConstants["h"].Value
	me := physicsConstants["me"].Value
	c := physicsConstants["c"].Value

	// Compton wavelength of electron
	lambdaC := h / (me * c)
	deltaLambda := lambdaC * (1 - math.Cos(theta))

	return types.Number(deltaLambda)
}

// FnComptonWavelengthParticle calculates Compton wavelength of any particle.
// λ_C = h/(mc)
// Args: mass (kg)
// Note: FnComptonWavelength for electron is in functions_physics_constants.go
func FnComptonWavelengthParticle(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("comptonwavelengthparticle requires 1 argument: mass")
	}

	m := args[0].AsFloat()

	if m <= 0 {
		return types.Error("comptonwavelengthparticle: mass must be positive")
	}

	h := physicsConstants["h"].Value
	c := physicsConstants["c"].Value

	return types.Number(h / (m * c))
}

// FnComptonScatteredWavelength calculates scattered photon wavelength.
// λ' = λ + (h/mc)(1 - cos(θ))
// Args: initial wavelength (m), scattering angle (degrees)
func FnComptonScatteredWavelength(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("comptonscattered requires 2 arguments: wavelength, angle_deg")
	}

	lambda := args[0].AsFloat()
	theta := args[1].AsFloat() * math.Pi / 180

	if lambda <= 0 {
		return types.Error("comptonscattered: wavelength must be positive")
	}

	h := physicsConstants["h"].Value
	me := physicsConstants["me"].Value
	c := physicsConstants["c"].Value

	lambdaC := h / (me * c)
	lambdaPrime := lambda + lambdaC*(1-math.Cos(theta))

	return types.Number(lambdaPrime)
}

// ════════════════════════════════════════════════════════════════
// HEISENBERG UNCERTAINTY
// ════════════════════════════════════════════════════════════════

// FnUncertaintyPosition calculates minimum position uncertainty.
// Δx ≥ ℏ/(2Δp)
// Args: momentum uncertainty (kg·m/s)
func FnUncertaintyPosition(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("uncertaintyposition requires 1 argument: delta_p")
	}

	deltaP := args[0].AsFloat()

	if deltaP <= 0 {
		return types.Error("uncertaintyposition: momentum uncertainty must be positive")
	}

	hbar := physicsConstants["hbar"].Value
	deltaX := hbar / (2 * deltaP)

	return types.Number(deltaX)
}

// FnUncertaintyMomentum calculates minimum momentum uncertainty.
// Δp ≥ ℏ/(2Δx)
// Args: position uncertainty (m)
func FnUncertaintyMomentum(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("uncertaintymomentum requires 1 argument: delta_x")
	}

	deltaX := args[0].AsFloat()

	if deltaX <= 0 {
		return types.Error("uncertaintymomentum: position uncertainty must be positive")
	}

	hbar := physicsConstants["hbar"].Value
	deltaP := hbar / (2 * deltaX)

	return types.Number(deltaP)
}

// FnUncertaintyEnergy calculates minimum energy uncertainty.
// ΔE ≥ ℏ/(2Δt)
// Args: time uncertainty (s)
func FnUncertaintyEnergy(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("uncertaintyenergy requires 1 argument: delta_t")
	}

	deltaT := args[0].AsFloat()

	if deltaT <= 0 {
		return types.Error("uncertaintyenergy: time uncertainty must be positive")
	}

	hbar := physicsConstants["hbar"].Value
	deltaE := hbar / (2 * deltaT)

	return types.Number(deltaE)
}

// FnUncertaintyTime calculates minimum time uncertainty.
// Δt ≥ ℏ/(2ΔE)
// Args: energy uncertainty (J)
func FnUncertaintyTime(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("uncertaintytime requires 1 argument: delta_E")
	}

	deltaE := args[0].AsFloat()

	if deltaE <= 0 {
		return types.Error("uncertaintytime: energy uncertainty must be positive")
	}

	hbar := physicsConstants["hbar"].Value
	deltaT := hbar / (2 * deltaE)

	return types.Number(deltaT)
}

// ════════════════════════════════════════════════════════════════
// BOHR MODEL
// ════════════════════════════════════════════════════════════════

// FnBohrRadius calculates radius of nth Bohr orbit.
// rₙ = n²a₀/Z where a₀ = 5.29177×10⁻¹¹ m
// Args: principal quantum number (n), [atomic number Z] (defaults to 1 for hydrogen)
func FnBohrRadius(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("bohrradius requires 1-2 arguments: n, [Z]")
	}

	n := args[0].AsFloat()
	Z := 1.0

	if len(args) == 2 {
		Z = args[1].AsFloat()
	}

	if n < 1 {
		return types.Error("bohrradius: n must be >= 1")
	}
	if Z < 1 {
		return types.Error("bohrradius: Z must be >= 1")
	}

	a0 := physicsConstants["a0"].Value
	r := n * n * a0 / Z

	return types.Number(r)
}

// FnBohrEnergy calculates energy of nth Bohr level.
// Eₙ = -13.6 eV × Z²/n²
// Args: principal quantum number (n), [atomic number Z] (defaults to 1)
func FnBohrEnergy(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("bohrenergy requires 1-2 arguments: n, [Z]")
	}

	n := args[0].AsFloat()
	Z := 1.0

	if len(args) == 2 {
		Z = args[1].AsFloat()
	}

	if n < 1 {
		return types.Error("bohrenergy: n must be >= 1")
	}
	if Z < 1 {
		return types.Error("bohrenergy: Z must be >= 1")
	}

	// -13.6 eV in Joules
	E1 := -13.6 * physicsConstants["ev"].Value
	E := E1 * Z * Z / (n * n)

	return types.Number(E)
}

// FnBohrEnergyEV calculates energy of nth Bohr level in eV.
// Eₙ = -13.6 eV × Z²/n²
// Args: principal quantum number (n), [atomic number Z] (defaults to 1)
func FnBohrEnergyEV(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("bohrenergyev requires 1-2 arguments: n, [Z]")
	}

	n := args[0].AsFloat()
	Z := 1.0

	if len(args) == 2 {
		Z = args[1].AsFloat()
	}

	if n < 1 {
		return types.Error("bohrenergyev: n must be >= 1")
	}
	if Z < 1 {
		return types.Error("bohrenergyev: Z must be >= 1")
	}

	E := -13.6 * Z * Z / (n * n)
	return types.Number(E)
}

// FnBohrVelocity calculates velocity of electron in nth Bohr orbit.
// vₙ = (Z/n)(e²/(2ε₀h)) = αc(Z/n)
// Args: principal quantum number (n), [atomic number Z] (defaults to 1)
func FnBohrVelocity(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("bohrvelocity requires 1-2 arguments: n, [Z]")
	}

	n := args[0].AsFloat()
	Z := 1.0

	if len(args) == 2 {
		Z = args[1].AsFloat()
	}

	if n < 1 {
		return types.Error("bohrvelocity: n must be >= 1")
	}
	if Z < 1 {
		return types.Error("bohrvelocity: Z must be >= 1")
	}

	alpha := physicsConstants["alpha"].Value
	c := physicsConstants["c"].Value

	v := alpha * c * Z / n
	return types.Number(v)
}

// FnRydbergWavelength calculates wavelength from Rydberg formula.
// 1/λ = R_H × Z² (1/n₁² - 1/n₂²)
// Args: n1 (lower), n2 (upper), [Z] (defaults to 1)
func FnRydbergWavelength(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("rydbergwavelength requires 2-3 arguments: n1, n2, [Z]")
	}

	n1 := args[0].AsFloat()
	n2 := args[1].AsFloat()
	Z := 1.0

	if len(args) == 3 {
		Z = args[2].AsFloat()
	}

	if n1 < 1 || n2 < 1 {
		return types.Error("rydbergwavelength: quantum numbers must be >= 1")
	}
	if n2 <= n1 {
		return types.Error("rydbergwavelength: n2 must be greater than n1")
	}

	RH := physicsConstants["r_inf"].Value
	invLambda := RH * Z * Z * (1/(n1*n1) - 1/(n2*n2))

	if invLambda == 0 {
		return types.Error("rydbergwavelength: invalid quantum numbers")
	}

	return types.Number(1 / invLambda)
}

// FnRydbergEnergy calculates transition energy.
// ΔE = 13.6 eV × Z² (1/n₁² - 1/n₂²)
// Args: n1 (lower), n2 (upper), [Z] (defaults to 1)
func FnRydbergEnergy(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("rydbergenergy requires 2-3 arguments: n1, n2, [Z]")
	}

	n1 := args[0].AsFloat()
	n2 := args[1].AsFloat()
	Z := 1.0

	if len(args) == 3 {
		Z = args[2].AsFloat()
	}

	if n1 < 1 || n2 < 1 {
		return types.Error("rydbergenergy: quantum numbers must be >= 1")
	}
	if n2 <= n1 {
		return types.Error("rydbergenergy: n2 must be greater than n1")
	}

	E1 := 13.6 * physicsConstants["ev"].Value
	deltaE := E1 * Z * Z * (1/(n1*n1) - 1/(n2*n2))

	return types.Number(deltaE)
}

// ════════════════════════════════════════════════════════════════
// PARTICLE IN A BOX
// ════════════════════════════════════════════════════════════════

// FnParticleBoxEnergy calculates energy levels of particle in 1D infinite well.
// Eₙ = n²h²/(8mL²)
// Args: quantum number (n), mass (kg), box length (m)
func FnParticleBoxEnergy(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("particleboxenergy requires 3 arguments: n, mass, length")
	}

	n := args[0].AsFloat()
	m := args[1].AsFloat()
	L := args[2].AsFloat()

	if n < 1 {
		return types.Error("particleboxenergy: n must be >= 1")
	}
	if m <= 0 {
		return types.Error("particleboxenergy: mass must be positive")
	}
	if L <= 0 {
		return types.Error("particleboxenergy: length must be positive")
	}

	h := physicsConstants["h"].Value
	E := n * n * h * h / (8 * m * L * L)

	return types.Number(E)
}

// FnParticleBoxWavelength calculates wavelength of particle in box.
// λₙ = 2L/n
// Args: quantum number (n), box length (m)
func FnParticleBoxWavelength(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("particleboxwavelength requires 2 arguments: n, length")
	}

	n := args[0].AsFloat()
	L := args[1].AsFloat()

	if n < 1 {
		return types.Error("particleboxwavelength: n must be >= 1")
	}
	if L <= 0 {
		return types.Error("particleboxwavelength: length must be positive")
	}

	return types.Number(2 * L / n)
}

// FnParticleBox3DEnergy calculates energy levels of particle in 3D box.
// E = (h²/8m)(nx²/Lx² + ny²/Ly² + nz²/Lz²)
// Args: nx, ny, nz, mass (kg), Lx (m), Ly (m), Lz (m)
func FnParticleBox3DEnergy(args []types.Value) types.Value {
	if len(args) != 7 {
		return types.Error("particlebox3denergy requires 7 arguments: nx, ny, nz, mass, Lx, Ly, Lz")
	}

	nx := args[0].AsFloat()
	ny := args[1].AsFloat()
	nz := args[2].AsFloat()
	m := args[3].AsFloat()
	Lx := args[4].AsFloat()
	Ly := args[5].AsFloat()
	Lz := args[6].AsFloat()

	if nx < 1 || ny < 1 || nz < 1 {
		return types.Error("particlebox3denergy: quantum numbers must be >= 1")
	}
	if m <= 0 {
		return types.Error("particlebox3denergy: mass must be positive")
	}
	if Lx <= 0 || Ly <= 0 || Lz <= 0 {
		return types.Error("particlebox3denergy: dimensions must be positive")
	}

	h := physicsConstants["h"].Value
	E := (h * h / (8 * m)) * (nx*nx/(Lx*Lx) + ny*ny/(Ly*Ly) + nz*nz/(Lz*Lz))

	return types.Number(E)
}

// ════════════════════════════════════════════════════════════════
// HARMONIC OSCILLATOR
// ════════════════════════════════════════════════════════════════

// FnQuantumHOEnergy calculates energy of quantum harmonic oscillator.
// Eₙ = (n + 1/2)ℏω
// Args: quantum number (n = 0, 1, 2, ...), angular frequency (rad/s)
func FnQuantumHOEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("quantumhoenergy requires 2 arguments: n, omega")
	}

	n := args[0].AsFloat()
	omega := args[1].AsFloat()

	if n < 0 {
		return types.Error("quantumhoenergy: n must be >= 0")
	}
	if omega <= 0 {
		return types.Error("quantumhoenergy: omega must be positive")
	}

	hbar := physicsConstants["hbar"].Value
	E := (n + 0.5) * hbar * omega

	return types.Number(E)
}

// FnQuantumHOFrequency calculates angular frequency from spring constant.
// ω = √(k/m)
// Args: spring constant (N/m), mass (kg)
func FnQuantumHOFrequency(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("quantumhofreq requires 2 arguments: k, mass")
	}

	k := args[0].AsFloat()
	m := args[1].AsFloat()

	if k <= 0 {
		return types.Error("quantumhofreq: spring constant must be positive")
	}
	if m <= 0 {
		return types.Error("quantumhofreq: mass must be positive")
	}

	return types.Number(math.Sqrt(k / m))
}

// FnZeroPointEnergy calculates zero-point energy.
// E₀ = (1/2)ℏω
// Args: angular frequency (rad/s)
func FnZeroPointEnergy(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("zeropointenergy requires 1 argument: omega")
	}

	omega := args[0].AsFloat()

	if omega <= 0 {
		return types.Error("zeropointenergy: omega must be positive")
	}

	hbar := physicsConstants["hbar"].Value
	return types.Number(0.5 * hbar * omega)
}

// ════════════════════════════════════════════════════════════════
// TUNNELING
// ════════════════════════════════════════════════════════════════

// FnTunnelingProbability calculates tunneling probability (approximate).
// T ≈ e^(-2κL) where κ = √(2m(V-E))/ℏ
// Args: mass (kg), barrier height (J), particle energy (J), barrier width (m)
func FnTunnelingProbability(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("tunnelingprob requires 4 arguments: mass, V_barrier, E_particle, width")
	}

	m := args[0].AsFloat()
	V := args[1].AsFloat()
	E := args[2].AsFloat()
	L := args[3].AsFloat()

	if m <= 0 {
		return types.Error("tunnelingprob: mass must be positive")
	}
	if L <= 0 {
		return types.Error("tunnelingprob: width must be positive")
	}
	if E >= V {
		return types.Error("tunnelingprob: particle energy must be less than barrier height")
	}

	hbar := physicsConstants["hbar"].Value
	kappa := math.Sqrt(2*m*(V-E)) / hbar
	T := math.Exp(-2 * kappa * L)

	return types.Number(T)
}

// FnTunnelingDecayConstant calculates decay constant κ for tunneling.
// κ = √(2m(V-E))/ℏ
// Args: mass (kg), barrier height (J), particle energy (J)
func FnTunnelingDecayConstant(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("tunnelingdecay requires 3 arguments: mass, V_barrier, E_particle")
	}

	m := args[0].AsFloat()
	V := args[1].AsFloat()
	E := args[2].AsFloat()

	if m <= 0 {
		return types.Error("tunnelingdecay: mass must be positive")
	}
	if E >= V {
		return types.Error("tunnelingdecay: particle energy must be less than barrier height")
	}

	hbar := physicsConstants["hbar"].Value
	kappa := math.Sqrt(2*m*(V-E)) / hbar

	return types.Number(kappa)
}

// ════════════════════════════════════════════════════════════════
// SPIN & ANGULAR MOMENTUM
// ════════════════════════════════════════════════════════════════

// FnOrbitalAngularMomentum calculates orbital angular momentum magnitude.
// L = ℏ√(l(l+1))
// Args: orbital quantum number (l)
func FnOrbitalAngularMomentum(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("orbitalangmom requires 1 argument: l")
	}

	l := args[0].AsFloat()

	if l < 0 {
		return types.Error("orbitalangmom: l must be >= 0")
	}

	hbar := physicsConstants["hbar"].Value
	L := hbar * math.Sqrt(l*(l+1))

	return types.Number(L)
}

// FnSpinAngularMomentum calculates spin angular momentum magnitude.
// S = ℏ√(s(s+1))
// Args: spin quantum number (s)
func FnSpinAngularMomentum(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("spinangmom requires 1 argument: s")
	}

	s := args[0].AsFloat()

	if s < 0 {
		return types.Error("spinangmom: s must be >= 0")
	}

	hbar := physicsConstants["hbar"].Value
	S := hbar * math.Sqrt(s*(s+1))

	return types.Number(S)
}

// FnMagneticMomentQuantum calculates magnetic moment from quantum numbers.
// μ = g × μ_B × √(l(l+1))
// Args: g-factor, angular momentum quantum number
func FnMagneticMomentQuantum(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("magneticmomentq requires 2 arguments: g_factor, quantum_number")
	}

	g := args[0].AsFloat()
	l := args[1].AsFloat()

	if l < 0 {
		return types.Error("magneticmomentq: quantum number must be >= 0")
	}

	hbar := physicsConstants["hbar"].Value
	me := physicsConstants["me"].Value
	e := physicsConstants["e_charge"].Value

	// Bohr magneton
	muB := e * hbar / (2 * me)

	// Magnetic moment magnitude
	mu := g * muB * math.Sqrt(l*(l+1))

	return types.Number(mu)
}

// FnZeemanSplitting calculates Zeeman energy splitting.
// ΔE = g × μ_B × B × m_l
// Args: g-factor, magnetic field (T), magnetic quantum number (m_l)
func FnZeemanSplitting(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("zeemansplitting requires 3 arguments: g_factor, B, m_l")
	}

	g := args[0].AsFloat()
	B := args[1].AsFloat()
	ml := args[2].AsFloat()

	hbar := physicsConstants["hbar"].Value
	me := physicsConstants["me"].Value
	e := physicsConstants["e_charge"].Value

	muB := e * hbar / (2 * me)
	deltaE := g * muB * B * ml

	return types.Number(deltaE)
}

// ════════════════════════════════════════════════════════════════
// NUCLEAR & PARTICLE PHYSICS
// ════════════════════════════════════════════════════════════════

// FnNuclearBindingEnergy calculates approximate binding energy (semi-empirical).
// B = aᵥA - aₛA^(2/3) - aᶜZ(Z-1)/A^(1/3) - aₐ(A-2Z)²/A + δ(A,Z)
// Args: mass number (A), atomic number (Z)
func FnNuclearBindingEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("nuclearbindingenergy requires 2 arguments: A, Z")
	}

	A := args[0].AsFloat()
	Z := args[1].AsFloat()

	if A < 1 || Z < 1 {
		return types.Error("nuclearbindingenergy: A and Z must be >= 1")
	}
	if Z > A {
		return types.Error("nuclearbindingenergy: Z cannot exceed A")
	}

	// Semi-empirical mass formula coefficients (MeV)
	av := 15.67 // Volume term
	as := 17.23 // Surface term
	ac := 0.75  // Coulomb term
	aa := 93.2  // Asymmetry term
	ap := 11.2  // Pairing term

	N := A - Z

	// Volume term
	Bv := av * A

	// Surface term
	Bs := as * math.Pow(A, 2.0/3.0)

	// Coulomb term
	Bc := ac * Z * (Z - 1) / math.Pow(A, 1.0/3.0)

	// Asymmetry term
	Ba := aa * (N - Z) * (N - Z) / A

	// Pairing term
	var delta float64
	Aint := int(A)
	Zint := int(Z)
	if Aint%2 == 0 && Zint%2 == 0 {
		delta = ap / math.Sqrt(A) // Even-even
	} else if Aint%2 == 1 {
		delta = 0 // Odd A
	} else {
		delta = -ap / math.Sqrt(A) // Odd-odd
	}

	B := Bv - Bs - Bc - Ba + delta

	// Convert to Joules
	MeVtoJ := 1.60218e-13
	return types.Number(B * MeVtoJ)
}

// FnNuclearBindingEnergyMeV same as above but returns MeV.
func FnNuclearBindingEnergyMeV(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("nuclearbindingmev requires 2 arguments: A, Z")
	}

	A := args[0].AsFloat()
	Z := args[1].AsFloat()

	if A < 1 || Z < 1 {
		return types.Error("nuclearbindingmev: A and Z must be >= 1")
	}
	if Z > A {
		return types.Error("nuclearbindingmev: Z cannot exceed A")
	}

	av := 15.67
	as := 17.23
	ac := 0.75
	aa := 93.2
	ap := 11.2

	N := A - Z

	Bv := av * A
	Bs := as * math.Pow(A, 2.0/3.0)
	Bc := ac * Z * (Z - 1) / math.Pow(A, 1.0/3.0)
	Ba := aa * (N - Z) * (N - Z) / A

	var delta float64
	Aint := int(A)
	Zint := int(Z)
	if Aint%2 == 0 && Zint%2 == 0 {
		delta = ap / math.Sqrt(A)
	} else if Aint%2 == 1 {
		delta = 0
	} else {
		delta = -ap / math.Sqrt(A)
	}

	B := Bv - Bs - Bc - Ba + delta
	return types.Number(B)
}

// FnMassDefect calculates mass defect.
// Δm = Zm_p + Nm_n - M_nucleus
// Args: atomic number (Z), neutron number (N), nuclear mass (kg)
func FnMassDefect(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("massdefect requires 3 arguments: Z, N, nuclear_mass")
	}

	Z := args[0].AsFloat()
	N := args[1].AsFloat()
	Mnuc := args[2].AsFloat()

	if Z < 0 || N < 0 {
		return types.Error("massdefect: Z and N must be >= 0")
	}

	mp := physicsConstants["mp"].Value
	mn := physicsConstants["mn"].Value

	deltam := Z*mp + N*mn - Mnuc
	return types.Number(deltam)
}

// FnMassEnergy calculates mass-energy equivalence.
// E = mc²
// Args: mass (kg)
func FnMassEnergy(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("massenergy requires 1 argument: mass")
	}

	m := args[0].AsFloat()
	c := physicsConstants["c"].Value

	return types.Number(m * c * c)
}

// FnMassFromEnergy calculates mass from energy.
// m = E/c²
// Args: energy (J)
func FnMassFromEnergy(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("massfromenergy requires 1 argument: energy")
	}

	E := args[0].AsFloat()
	c := physicsConstants["c"].Value

	return types.Number(E / (c * c))
}

// FnRadioactiveDecay calculates remaining nuclei after decay.
// N(t) = N₀ × e^(-λt)
// Args: initial count (N0), decay constant (λ, s⁻¹), time (s)
func FnRadioactiveDecay(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("radioactivedecay requires 3 arguments: N0, lambda, time")
	}

	N0 := args[0].AsFloat()
	lambda := args[1].AsFloat()
	t := args[2].AsFloat()

	if N0 < 0 {
		return types.Error("radioactivedecay: N0 cannot be negative")
	}
	if lambda < 0 {
		return types.Error("radioactivedecay: decay constant cannot be negative")
	}

	N := N0 * math.Exp(-lambda*t)
	return types.Number(N)
}

// FnHalfLife calculates half-life from decay constant.
// t½ = ln(2)/λ
// Args: decay constant (λ, s⁻¹)
func FnHalfLife(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("halflife requires 1 argument: lambda")
	}

	lambda := args[0].AsFloat()

	if lambda <= 0 {
		return types.Error("halflife: decay constant must be positive")
	}

	return types.Number(math.Ln2 / lambda)
}

// FnDecayConstantFromHalfLife calculates decay constant from half-life.
// λ = ln(2)/t½
// Args: half-life (s)
func FnDecayConstantFromHalfLife(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("decayconstantfromhl requires 1 argument: half_life")
	}

	tHalf := args[0].AsFloat()

	if tHalf <= 0 {
		return types.Error("decayconstantfromhl: half-life must be positive")
	}

	return types.Number(math.Ln2 / tHalf)
}

// FnActivity calculates radioactive activity.
// A = λN = A₀e^(-λt)
// Args: decay constant (λ), number of nuclei (N)
func FnActivity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("activity requires 2 arguments: lambda, N")
	}

	lambda := args[0].AsFloat()
	N := args[1].AsFloat()

	return types.Number(lambda * N)
}
