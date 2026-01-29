// internal/eval/functions_physics_relativity.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// LORENTZ FACTOR & BASIC RELATIVITY
// ════════════════════════════════════════════════════════════════

// FnLorentzFactor calculates the Lorentz factor γ.
// γ = 1/√(1 - v²/c²)
// Args: velocity (m/s)
func FnLorentzFactor(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("lorentzfactor requires 1 argument: velocity")
	}

	v := args[0].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("lorentzfactor: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	return types.Number(gamma)
}

// FnLorentzFactorFromBeta calculates Lorentz factor from β = v/c.
// γ = 1/√(1 - β²)
// Args: beta (v/c ratio)
func FnLorentzFactorFromBeta(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("lorentzfrombeta requires 1 argument: beta")
	}

	beta := args[0].AsFloat()

	if math.Abs(beta) >= 1 {
		return types.Error("lorentzfrombeta: beta must be less than 1")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	return types.Number(gamma)
}

// FnBetaFromLorentz calculates β from Lorentz factor.
// β = √(1 - 1/γ²)
// Args: gamma (Lorentz factor)
func FnBetaFromLorentz(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("betafromlorentz requires 1 argument: gamma")
	}

	gamma := args[0].AsFloat()

	if gamma < 1 {
		return types.Error("betafromlorentz: gamma must be >= 1")
	}

	beta := math.Sqrt(1 - 1/(gamma*gamma))
	return types.Number(beta)
}

// FnVelocityFromLorentz calculates velocity from Lorentz factor.
// v = c × √(1 - 1/γ²)
// Args: gamma (Lorentz factor)
func FnVelocityFromLorentz(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("velocityfromlorentz requires 1 argument: gamma")
	}

	gamma := args[0].AsFloat()

	if gamma < 1 {
		return types.Error("velocityfromlorentz: gamma must be >= 1")
	}

	c := physicsConstants["c"].Value
	v := c * math.Sqrt(1-1/(gamma*gamma))
	return types.Number(v)
}

// ════════════════════════════════════════════════════════════════
// TIME DILATION
// ════════════════════════════════════════════════════════════════

// FnTimeDilation calculates dilated time.
// Δt = γΔt₀ (moving clock runs slower)
// Args: proper time (s), velocity (m/s)
func FnTimeDilation(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("timedilation requires 2 arguments: proper_time, velocity")
	}

	dt0 := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("timedilation: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	dt := gamma * dt0

	return types.Number(dt)
}

// FnProperTime calculates proper time from dilated time.
// Δt₀ = Δt/γ
// Args: dilated time (s), velocity (m/s)
func FnProperTime(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("propertime requires 2 arguments: dilated_time, velocity")
	}

	dt := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("propertime: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	dt0 := dt / gamma

	return types.Number(dt0)
}

// FnTimeDilationGravity calculates gravitational time dilation.
// Δt = Δt₀ / √(1 - 2GM/(rc²))
// Args: proper time (s), mass (kg), radius (m)
func FnTimeDilationGravity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("timedilationgrav requires 3 arguments: proper_time, mass, radius")
	}

	dt0 := args[0].AsFloat()
	M := args[1].AsFloat()
	r := args[2].AsFloat()

	if r <= 0 {
		return types.Error("timedilationgrav: radius must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value

	rs := 2 * G * M / (c * c) // Schwarzschild radius
	if r <= rs {
		return types.Error("timedilationgrav: radius must be greater than Schwarzschild radius")
	}

	factor := math.Sqrt(1 - rs/r)
	dt := dt0 / factor

	return types.Number(dt)
}

// ════════════════════════════════════════════════════════════════
// LENGTH CONTRACTION
// ════════════════════════════════════════════════════════════════

// FnLengthContraction calculates contracted length.
// L = L₀/γ = L₀√(1 - v²/c²)
// Args: proper length (m), velocity (m/s)
func FnLengthContraction(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("lengthcontraction requires 2 arguments: proper_length, velocity")
	}

	L0 := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("lengthcontraction: velocity must be less than speed of light")
	}

	L := L0 * math.Sqrt(1-beta*beta)
	return types.Number(L)
}

// FnProperLength calculates proper length from contracted length.
// L₀ = γL
// Args: contracted length (m), velocity (m/s)
func FnProperLength(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("properlength requires 2 arguments: contracted_length, velocity")
	}

	L := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("properlength: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	L0 := gamma * L

	return types.Number(L0)
}

// ════════════════════════════════════════════════════════════════
// VELOCITY ADDITION
// ════════════════════════════════════════════════════════════════

// FnRelativisticVelocityAdd calculates relativistic velocity addition.
// u = (v + u')/(1 + vu'/c²)
// Args: frame velocity v (m/s), object velocity in moving frame u' (m/s)
func FnRelativisticVelocityAdd(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relvelocityadd requires 2 arguments: v, u_prime")
	}

	v := args[0].AsFloat()
	uPrime := args[1].AsFloat()
	c := physicsConstants["c"].Value

	if math.Abs(v) >= c || math.Abs(uPrime) >= c {
		return types.Error("relvelocityadd: velocities must be less than speed of light")
	}

	u := (v + uPrime) / (1 + v*uPrime/(c*c))
	return types.Number(u)
}

// FnRelativisticVelocityAddBeta calculates velocity addition using β values.
// β = (β₁ + β₂)/(1 + β₁β₂)
// Args: beta1, beta2
func FnRelativisticVelocityAddBeta(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relvelocityaddbeta requires 2 arguments: beta1, beta2")
	}

	beta1 := args[0].AsFloat()
	beta2 := args[1].AsFloat()

	if math.Abs(beta1) >= 1 || math.Abs(beta2) >= 1 {
		return types.Error("relvelocityaddbeta: betas must be less than 1")
	}

	beta := (beta1 + beta2) / (1 + beta1*beta2)
	return types.Number(beta)
}

// ════════════════════════════════════════════════════════════════
// RELATIVISTIC MOMENTUM & ENERGY
// ════════════════════════════════════════════════════════════════

// FnRelativisticMomentum calculates relativistic momentum.
// p = γmv
// Args: mass (kg), velocity (m/s)
func FnRelativisticMomentum(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relmomentum requires 2 arguments: mass, velocity")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("relmomentum: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	p := gamma * m * v

	return types.Number(p)
}

// FnRelativisticEnergy calculates total relativistic energy.
// E = γmc²
// Args: mass (kg), velocity (m/s)
func FnRelativisticEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relenergy requires 2 arguments: mass, velocity")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("relenergy: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	E := gamma * m * c * c

	return types.Number(E)
}

// FnRestEnergy calculates rest energy.
// E₀ = mc²
// Args: mass (kg)
func FnRestEnergy(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("restenergy requires 1 argument: mass")
	}

	m := args[0].AsFloat()
	c := physicsConstants["c"].Value

	return types.Number(m * c * c)
}

// FnRelativisticKE calculates relativistic kinetic energy.
// KE = (γ - 1)mc²
// Args: mass (kg), velocity (m/s)
func FnRelativisticKE(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relke requires 2 arguments: mass, velocity")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("relke: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	KE := (gamma - 1) * m * c * c

	return types.Number(KE)
}

// FnEnergyMomentumRelation calculates energy from momentum (or vice versa).
// E² = (pc)² + (mc²)²
// Args: momentum (kg·m/s), mass (kg) -> returns energy
func FnEnergyMomentumRelation(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("energymomentum requires 2 arguments: momentum, mass")
	}

	p := args[0].AsFloat()
	m := args[1].AsFloat()
	c := physicsConstants["c"].Value

	E := math.Sqrt(p*p*c*c + m*m*c*c*c*c)
	return types.Number(E)
}

// FnMomentumFromEnergy calculates momentum from total energy and mass.
// p = √(E² - (mc²)²) / c
// Args: total energy (J), mass (kg)
func FnMomentumFromEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("momentumfromenergy requires 2 arguments: energy, mass")
	}

	E := args[0].AsFloat()
	m := args[1].AsFloat()
	c := physicsConstants["c"].Value

	mc2 := m * c * c
	if E < mc2 {
		return types.Error("momentumfromenergy: energy must be >= rest energy")
	}

	p := math.Sqrt(E*E-mc2*mc2) / c
	return types.Number(p)
}

// FnVelocityFromMomentum calculates velocity from relativistic momentum.
// v = pc²/E = p/(γm)
// Args: momentum (kg·m/s), mass (kg)
func FnVelocityFromMomentum(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("velocityfrommom requires 2 arguments: momentum, mass")
	}

	p := args[0].AsFloat()
	m := args[1].AsFloat()
	c := physicsConstants["c"].Value

	if m == 0 {
		// Massless particle
		return types.Number(c)
	}

	E := math.Sqrt(p*p*c*c + m*m*c*c*c*c)
	v := p * c * c / E

	return types.Number(v)
}

// ════════════════════════════════════════════════════════════════
// RELATIVISTIC MASS (for reference, though rest mass is preferred)
// ════════════════════════════════════════════════════════════════

// FnRelativisticMass calculates relativistic mass (historical concept).
// m_rel = γm₀
// Args: rest mass (kg), velocity (m/s)
func FnRelativisticMass(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relmass requires 2 arguments: rest_mass, velocity")
	}

	m0 := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("relmass: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	return types.Number(gamma * m0)
}

// ════════════════════════════════════════════════════════════════
// RELATIVISTIC DOPPLER EFFECT
// ════════════════════════════════════════════════════════════════

// FnRelativisticDopplerApproach calculates Doppler shift for approaching source.
// f_obs = f_source × √((1 + β)/(1 - β))
// Args: source frequency (Hz), velocity (m/s) [positive = approaching]
func FnRelativisticDopplerApproach(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("reldopplerapproach requires 2 arguments: frequency, velocity")
	}

	f := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("reldopplerapproach: velocity must be less than speed of light")
	}

	fObs := f * math.Sqrt((1+beta)/(1-beta))
	return types.Number(fObs)
}

// FnRelativisticDopplerRecede calculates Doppler shift for receding source.
// f_obs = f_source × √((1 - β)/(1 + β))
// Args: source frequency (Hz), velocity (m/s) [positive = receding]
func FnRelativisticDopplerRecede(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("reldopplerrecede requires 2 arguments: frequency, velocity")
	}

	f := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("reldopplerrecede: velocity must be less than speed of light")
	}

	fObs := f * math.Sqrt((1-beta)/(1+beta))
	return types.Number(fObs)
}

// FnRedshift calculates cosmological redshift.
// z = (λ_obs - λ_emit)/λ_emit = √((1 + β)/(1 - β)) - 1
// Args: velocity (m/s) [positive = receding]
func FnRedshift(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("redshift requires 1 argument: velocity")
	}

	v := args[0].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("redshift: velocity must be less than speed of light")
	}

	z := math.Sqrt((1+beta)/(1-beta)) - 1
	return types.Number(z)
}

// FnVelocityFromRedshift calculates velocity from redshift.
// β = ((z+1)² - 1) / ((z+1)² + 1)
// Args: redshift z
func FnVelocityFromRedshift(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("velocityfromredshift requires 1 argument: z")
	}

	z := args[0].AsFloat()

	if z < -1 {
		return types.Error("velocityfromredshift: z must be >= -1")
	}

	c := physicsConstants["c"].Value
	zp1sq := (z + 1) * (z + 1)
	beta := (zp1sq - 1) / (zp1sq + 1)
	v := beta * c

	return types.Number(v)
}

// FnTransverseDoppler calculates transverse Doppler effect.
// f_obs = f_source / γ (time dilation only, perpendicular motion)
// Args: source frequency (Hz), velocity (m/s)
func FnTransverseDoppler(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("transversedoppler requires 2 arguments: frequency, velocity")
	}

	f := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("transversedoppler: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	fObs := f / gamma

	return types.Number(fObs)
}

// ════════════════════════════════════════════════════════════════
// LORENTZ TRANSFORMATIONS
// ════════════════════════════════════════════════════════════════

// FnLorentzTransformX calculates Lorentz-transformed x coordinate.
// x' = γ(x - vt)
// Args: x (m), t (s), velocity (m/s)
func FnLorentzTransformX(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("lorentztransformx requires 3 arguments: x, t, velocity")
	}

	x := args[0].AsFloat()
	t := args[1].AsFloat()
	v := args[2].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("lorentztransformx: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	xPrime := gamma * (x - v*t)

	return types.Number(xPrime)
}

// FnLorentzTransformT calculates Lorentz-transformed time coordinate.
// t' = γ(t - vx/c²)
// Args: x (m), t (s), velocity (m/s)
func FnLorentzTransformT(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("lorentztransformt requires 3 arguments: x, t, velocity")
	}

	x := args[0].AsFloat()
	t := args[1].AsFloat()
	v := args[2].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("lorentztransformt: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	tPrime := gamma * (t - v*x/(c*c))

	return types.Number(tPrime)
}

// FnInverseLorentzX calculates inverse Lorentz transformation for x.
// x = γ(x' + vt')
// Args: x' (m), t' (s), velocity (m/s)
func FnInverseLorentzX(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("inverselorentzx requires 3 arguments: x_prime, t_prime, velocity")
	}

	xPrime := args[0].AsFloat()
	tPrime := args[1].AsFloat()
	v := args[2].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("inverselorentzx: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	x := gamma * (xPrime + v*tPrime)

	return types.Number(x)
}

// FnInverseLorentzT calculates inverse Lorentz transformation for t.
// t = γ(t' + vx'/c²)
// Args: x' (m), t' (s), velocity (m/s)
func FnInverseLorentzT(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("inverselorentzt requires 3 arguments: x_prime, t_prime, velocity")
	}

	xPrime := args[0].AsFloat()
	tPrime := args[1].AsFloat()
	v := args[2].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("inverselorentzt: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	t := gamma * (tPrime + v*xPrime/(c*c))

	return types.Number(t)
}

// ════════════════════════════════════════════════════════════════
// SPACETIME INTERVAL
// ════════════════════════════════════════════════════════════════

// FnSpacetimeInterval calculates spacetime interval squared.
// (Δs)² = (cΔt)² - (Δx)² - (Δy)² - (Δz)² (timelike positive convention)
// Args: Δt (s), Δx (m), [Δy] (m), [Δz] (m)
func FnSpacetimeInterval(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 4 {
		return types.Error("spacetimeinterval requires 2-4 arguments: dt, dx, [dy], [dz]")
	}

	dt := args[0].AsFloat()
	dx := args[1].AsFloat()
	dy := 0.0
	dz := 0.0

	if len(args) >= 3 {
		dy = args[2].AsFloat()
	}
	if len(args) >= 4 {
		dz = args[3].AsFloat()
	}

	c := physicsConstants["c"].Value
	ds2 := c*c*dt*dt - dx*dx - dy*dy - dz*dz

	return types.Number(ds2)
}

// FnProperTimeInterval calculates proper time from spacetime interval.
// Δτ = √((Δs)²) / c for timelike intervals
// Args: spacetime interval squared (m²)
func FnProperTimeInterval(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("propertimeinterval requires 1 argument: ds_squared")
	}

	ds2 := args[0].AsFloat()

	if ds2 < 0 {
		return types.Error("propertimeinterval: interval is spacelike (ds² < 0)")
	}

	c := physicsConstants["c"].Value
	tau := math.Sqrt(ds2) / c

	return types.Number(tau)
}

// FnProperLengthInterval calculates proper length from spacetime interval.
// ΔL = √(-(Δs)²) for spacelike intervals
// Args: spacetime interval squared (m²)
func FnProperLengthInterval(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("properlengthinterval requires 1 argument: ds_squared")
	}

	ds2 := args[0].AsFloat()

	if ds2 > 0 {
		return types.Error("properlengthinterval: interval is timelike (ds² > 0)")
	}

	L := math.Sqrt(-ds2)
	return types.Number(L)
}

// FnIntervalType determines if interval is timelike, spacelike, or lightlike.
// Returns: 1 (timelike), 0 (lightlike), -1 (spacelike)
// Args: spacetime interval squared (m²)
func FnIntervalType(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("intervaltype requires 1 argument: ds_squared")
	}

	ds2 := args[0].AsFloat()

	// Use small epsilon for lightlike comparison
	epsilon := 1e-10

	if ds2 > epsilon {
		return types.Number(1) // Timelike
	} else if ds2 < -epsilon {
		return types.Number(-1) // Spacelike
	}
	return types.Number(0) // Lightlike
}

// ════════════════════════════════════════════════════════════════
// FOUR-VECTORS
// ════════════════════════════════════════════════════════════════

// FnFourMomentumE calculates energy component of four-momentum.
// p⁰ = E/c = γmc
// Args: mass (kg), velocity (m/s)
func FnFourMomentumE(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("fourmomentumt requires 2 arguments: mass, velocity")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("fourmomentumt: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	p0 := gamma * m * c

	return types.Number(p0)
}

// FnFourMomentumP calculates spatial component of four-momentum.
// p = γmv
// Args: mass (kg), velocity (m/s)
func FnFourMomentumP(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("fourmomentumspatial requires 2 arguments: mass, velocity")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("fourmomentumspatial: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	p := gamma * m * v

	return types.Number(p)
}

// FnFourMomentumMagnitude calculates invariant mass from four-momentum.
// m²c² = (E/c)² - p²
// Args: energy (J), momentum magnitude (kg·m/s)
func FnFourMomentumMagnitude(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("fourmomentummag requires 2 arguments: energy, momentum")
	}

	E := args[0].AsFloat()
	p := args[1].AsFloat()
	c := physicsConstants["c"].Value

	m2c2 := E*E/(c*c) - p*p
	if m2c2 < 0 {
		return types.Error("fourmomentummag: invalid four-momentum (tachyonic)")
	}

	m := math.Sqrt(m2c2) / c
	return types.Number(m)
}

// ════════════════════════════════════════════════════════════════
// RELATIVISTIC COLLISIONS
// ════════════════════════════════════════════════════════════════

// FnCenterOfMassEnergy calculates center-of-mass energy for collisions.
// E_cm = √(2E₁E₂ + 2p₁p₂c² + m₁²c⁴ + m₂²c⁴) for head-on collision
// Simplified for one particle at rest: E_cm = √(m₁²c⁴ + m₂²c⁴ + 2E₁m₂c²)
// Args: E1 (J), m1 (kg), m2 (kg) [particle 2 at rest]
func FnCenterOfMassEnergy(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("centerofmassenergy requires 3 arguments: E1, m1, m2")
	}

	E1 := args[0].AsFloat()
	m1 := args[1].AsFloat()
	m2 := args[2].AsFloat()
	c := physicsConstants["c"].Value

	c2 := c * c
	c4 := c2 * c2

	Ecm2 := m1*m1*c4 + m2*m2*c4 + 2*E1*m2*c2
	if Ecm2 < 0 {
		return types.Error("centerofmassenergy: invalid energy configuration")
	}

	return types.Number(math.Sqrt(Ecm2))
}

// FnThresholdEnergy calculates threshold energy for particle production.
// E_threshold = ((M_total)² - (m₁ + m₂)²)c²/(2m₂) for target at rest
// Args: m1 (projectile mass, kg), m2 (target mass, kg), M_total (total product mass, kg)
func FnThresholdEnergy(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("thresholdenergy requires 3 arguments: m1, m2, M_total")
	}

	m1 := args[0].AsFloat()
	m2 := args[1].AsFloat()
	M := args[2].AsFloat()
	c := physicsConstants["c"].Value

	if m2 == 0 {
		return types.Error("thresholdenergy: target mass cannot be zero")
	}

	c2 := c * c
	Eth := (M*M - (m1+m2)*(m1+m2)) * c2 / (2 * m2)

	return types.Number(Eth)
}

// ════════════════════════════════════════════════════════════════
// RELATIVISTIC ACCELERATION
// ════════════════════════════════════════════════════════════════

// FnProperAcceleration calculates proper acceleration.
// a₀ = γ³a (for acceleration parallel to velocity)
// Args: coordinate acceleration (m/s²), velocity (m/s)
func FnProperAcceleration(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("properaccel requires 2 arguments: coordinate_accel, velocity")
	}

	a := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("properaccel: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	a0 := gamma * gamma * gamma * a

	return types.Number(a0)
}

// FnCoordinateAcceleration calculates coordinate acceleration from proper.
// a = a₀/γ³
// Args: proper acceleration (m/s²), velocity (m/s)
func FnCoordinateAcceleration(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("coordaccel requires 2 arguments: proper_accel, velocity")
	}

	a0 := args[0].AsFloat()
	v := args[1].AsFloat()
	c := physicsConstants["c"].Value

	beta := v / c
	if math.Abs(beta) >= 1 {
		return types.Error("coordaccel: velocity must be less than speed of light")
	}

	gamma := 1 / math.Sqrt(1-beta*beta)
	a := a0 / (gamma * gamma * gamma)

	return types.Number(a)
}

// FnRelativisticRocket calculates final velocity for constant proper acceleration.
// v = a₀t / √(1 + (a₀t/c)²) (coordinate time)
// Args: proper acceleration (m/s²), coordinate time (s)
func FnRelativisticRocket(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relrocket requires 2 arguments: proper_accel, time")
	}

	a0 := args[0].AsFloat()
	t := args[1].AsFloat()
	c := physicsConstants["c"].Value

	v := a0 * t / math.Sqrt(1+(a0*t/c)*(a0*t/c))
	return types.Number(v)
}

// FnRelativisticRocketDistance calculates distance traveled under constant proper acceleration.
// x = (c²/a₀)(√(1 + (a₀t/c)²) - 1)
// Args: proper acceleration (m/s²), coordinate time (s)
func FnRelativisticRocketDistance(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relrocketdist requires 2 arguments: proper_accel, time")
	}

	a0 := args[0].AsFloat()
	t := args[1].AsFloat()
	c := physicsConstants["c"].Value

	if a0 == 0 {
		return types.Number(0)
	}

	x := (c * c / a0) * (math.Sqrt(1+(a0*t/c)*(a0*t/c)) - 1)
	return types.Number(x)
}

// FnRelativisticRocketProperTime calculates proper time for constant acceleration.
// τ = (c/a₀) × arcsinh(a₀t/c)
// Args: proper acceleration (m/s²), coordinate time (s)
func FnRelativisticRocketProperTime(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("relrocketpropertime requires 2 arguments: proper_accel, coord_time")
	}

	a0 := args[0].AsFloat()
	t := args[1].AsFloat()
	c := physicsConstants["c"].Value

	if a0 == 0 {
		return types.Number(t)
	}

	tau := (c / a0) * math.Asinh(a0*t/c)
	return types.Number(tau)
}
