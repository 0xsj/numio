// internal/eval/functions_physics_gravity.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// GRAVITATIONAL FORCE & FIELD
// ════════════════════════════════════════════════════════════════

// FnGravity calculates gravitational force between two masses.
// F = Gm₁m₂/r²
// Args: mass1 (kg), mass2 (kg), distance (m)
func FnGravity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("gravity requires 3 arguments: m1, m2, distance")
	}

	m1 := args[0].AsFloat()
	m2 := args[1].AsFloat()
	r := args[2].AsFloat()

	if r == 0 {
		return types.Error("gravity: distance cannot be zero")
	}

	G := physicsConstants["g_gravity"].Value
	F := G * m1 * m2 / (r * r)

	return types.Number(F)
}

// FnGravField calculates gravitational field strength.
// g = GM/r²
// Args: mass (kg), distance (m)
func FnGravField(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("gravfield requires 2 arguments: mass, distance")
	}

	M := args[0].AsFloat()
	r := args[1].AsFloat()

	if r == 0 {
		return types.Error("gravfield: distance cannot be zero")
	}

	G := physicsConstants["g_gravity"].Value
	g := G * M / (r * r)

	return types.Number(g)
}

// FnGravPotential calculates gravitational potential.
// V = -GM/r
// Args: mass (kg), distance (m)
func FnGravPotential(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("gravpotential requires 2 arguments: mass, distance")
	}

	M := args[0].AsFloat()
	r := args[1].AsFloat()

	if r == 0 {
		return types.Error("gravpotential: distance cannot be zero")
	}

	G := physicsConstants["g_gravity"].Value
	V := -G * M / r

	return types.Number(V)
}

// FnGravPotentialEnergy calculates gravitational potential energy.
// U = -Gm₁m₂/r
// Args: mass1 (kg), mass2 (kg), distance (m)
func FnGravPotentialEnergy(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("gravpe requires 3 arguments: m1, m2, distance")
	}

	m1 := args[0].AsFloat()
	m2 := args[1].AsFloat()
	r := args[2].AsFloat()

	if r == 0 {
		return types.Error("gravpe: distance cannot be zero")
	}

	G := physicsConstants["g_gravity"].Value
	U := -G * m1 * m2 / r

	return types.Number(U)
}

// FnSurfaceGravity calculates surface gravity of a body.
// g = GM/R²
// Args: mass (kg), radius (m)
func FnSurfaceGravity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("surfacegravity requires 2 arguments: mass, radius")
	}

	M := args[0].AsFloat()
	R := args[1].AsFloat()

	if R == 0 {
		return types.Error("surfacegravity: radius cannot be zero")
	}

	G := physicsConstants["g_gravity"].Value
	g := G * M / (R * R)

	return types.Number(g)
}

// ════════════════════════════════════════════════════════════════
// ORBITAL MECHANICS
// ════════════════════════════════════════════════════════════════

// FnEscapeVelocity calculates escape velocity.
// v_e = √(2GM/r)
// Args: mass (kg), radius (m)
func FnEscapeVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("escapevel requires 2 arguments: mass, radius")
	}

	M := args[0].AsFloat()
	r := args[1].AsFloat()

	if r <= 0 {
		return types.Error("escapevel: radius must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	v := math.Sqrt(2 * G * M / r)

	return types.Number(v)
}

// FnOrbitalVelocity calculates circular orbital velocity.
// v = √(GM/r)
// Args: central mass (kg), orbital radius (m)
func FnOrbitalVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("orbitalvel requires 2 arguments: mass, radius")
	}

	M := args[0].AsFloat()
	r := args[1].AsFloat()

	if r <= 0 {
		return types.Error("orbitalvel: radius must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	v := math.Sqrt(G * M / r)

	return types.Number(v)
}

// FnOrbitalPeriod calculates orbital period (Kepler's 3rd law).
// T = 2π√(r³/GM)
// Args: central mass (kg), orbital radius (m)
func FnOrbitalPeriod(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("orbitalperiod requires 2 arguments: mass, radius")
	}

	M := args[0].AsFloat()
	r := args[1].AsFloat()

	if r <= 0 {
		return types.Error("orbitalperiod: radius must be positive")
	}
	if M <= 0 {
		return types.Error("orbitalperiod: mass must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	T := 2 * math.Pi * math.Sqrt(r*r*r/(G*M))

	return types.Number(T)
}

// FnOrbitalRadius calculates orbital radius from period.
// r = ∛(GMT²/(4π²))
// Args: central mass (kg), period (s)
func FnOrbitalRadius(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("orbitalradius requires 2 arguments: mass, period")
	}

	M := args[0].AsFloat()
	T := args[1].AsFloat()

	if T <= 0 {
		return types.Error("orbitalradius: period must be positive")
	}
	if M <= 0 {
		return types.Error("orbitalradius: mass must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	r := math.Pow(G*M*T*T/(4*math.Pi*math.Pi), 1.0/3.0)

	return types.Number(r)
}

// FnOrbitalEnergy calculates total orbital energy.
// E = -GMm/(2r) for circular orbit
// Args: central mass (kg), orbiting mass (kg), radius (m)
func FnOrbitalEnergy(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("orbitalenergy requires 3 arguments: M, m, radius")
	}

	M := args[0].AsFloat()
	m := args[1].AsFloat()
	r := args[2].AsFloat()

	if r <= 0 {
		return types.Error("orbitalenergy: radius must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	E := -G * M * m / (2 * r)

	return types.Number(E)
}

// FnSpecificOrbitalEnergy calculates specific orbital energy.
// ε = -GM/(2a) for elliptical orbit
// Args: central mass (kg), semi-major axis (m)
func FnSpecificOrbitalEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("specificorbitalenergy requires 2 arguments: mass, semi_major_axis")
	}

	M := args[0].AsFloat()
	a := args[1].AsFloat()

	if a <= 0 {
		return types.Error("specificorbitalenergy: semi-major axis must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	epsilon := -G * M / (2 * a)

	return types.Number(epsilon)
}

// ════════════════════════════════════════════════════════════════
// ELLIPTICAL ORBITS
// ════════════════════════════════════════════════════════════════

// FnPeriapsis calculates periapsis distance.
// r_p = a(1 - e)
// Args: semi-major axis (m), eccentricity
func FnPeriapsis(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("periapsis requires 2 arguments: semi_major_axis, eccentricity")
	}

	a := args[0].AsFloat()
	e := args[1].AsFloat()

	if e < 0 || e >= 1 {
		return types.Error("periapsis: eccentricity must be in [0, 1) for elliptical orbit")
	}

	return types.Number(a * (1 - e))
}

// FnApoapsis calculates apoapsis distance.
// r_a = a(1 + e)
// Args: semi-major axis (m), eccentricity
func FnApoapsis(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("apoapsis requires 2 arguments: semi_major_axis, eccentricity")
	}

	a := args[0].AsFloat()
	e := args[1].AsFloat()

	if e < 0 || e >= 1 {
		return types.Error("apoapsis: eccentricity must be in [0, 1) for elliptical orbit")
	}

	return types.Number(a * (1 + e))
}

// FnSemiMajorAxis calculates semi-major axis from periapsis and apoapsis.
// a = (r_p + r_a) / 2
// Args: periapsis (m), apoapsis (m)
func FnSemiMajorAxis(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("semimajoraxis requires 2 arguments: periapsis, apoapsis")
	}

	rp := args[0].AsFloat()
	ra := args[1].AsFloat()

	return types.Number((rp + ra) / 2)
}

// FnEccentricity calculates orbital eccentricity.
// e = (r_a - r_p) / (r_a + r_p)
// Args: periapsis (m), apoapsis (m)
func FnEccentricity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("eccentricity requires 2 arguments: periapsis, apoapsis")
	}

	rp := args[0].AsFloat()
	ra := args[1].AsFloat()

	if rp+ra == 0 {
		return types.Error("eccentricity: sum of distances cannot be zero")
	}

	e := (ra - rp) / (ra + rp)
	return types.Number(e)
}

// FnOrbitalVelAtRadius calculates orbital velocity at a given radius (vis-viva equation).
// v = √(GM(2/r - 1/a))
// Args: central mass (kg), radius (m), semi-major axis (m)
func FnOrbitalVelAtRadius(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("orbitalvelatr requires 3 arguments: mass, radius, semi_major_axis")
	}

	M := args[0].AsFloat()
	r := args[1].AsFloat()
	a := args[2].AsFloat()

	if r <= 0 {
		return types.Error("orbitalvelatr: radius must be positive")
	}
	if a <= 0 {
		return types.Error("orbitalvelatr: semi-major axis must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	v := math.Sqrt(G * M * (2/r - 1/a))

	return types.Number(v)
}

// FnPeriapsisVelocity calculates velocity at periapsis.
// v_p = √(GM(1+e)/(a(1-e)))
// Args: central mass (kg), semi-major axis (m), eccentricity
func FnPeriapsisVelocity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("periapsisvel requires 3 arguments: mass, semi_major_axis, eccentricity")
	}

	M := args[0].AsFloat()
	a := args[1].AsFloat()
	e := args[2].AsFloat()

	if a <= 0 {
		return types.Error("periapsisvel: semi-major axis must be positive")
	}
	if e < 0 || e >= 1 {
		return types.Error("periapsisvel: eccentricity must be in [0, 1)")
	}

	G := physicsConstants["g_gravity"].Value
	v := math.Sqrt(G * M * (1 + e) / (a * (1 - e)))

	return types.Number(v)
}

// FnApoapsisVelocity calculates velocity at apoapsis.
// v_a = √(GM(1-e)/(a(1+e)))
// Args: central mass (kg), semi-major axis (m), eccentricity
func FnApoapsisVelocity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("apoapsisvel requires 3 arguments: mass, semi_major_axis, eccentricity")
	}

	M := args[0].AsFloat()
	a := args[1].AsFloat()
	e := args[2].AsFloat()

	if a <= 0 {
		return types.Error("apoapsisvel: semi-major axis must be positive")
	}
	if e < 0 || e >= 1 {
		return types.Error("apoapsisvel: eccentricity must be in [0, 1)")
	}

	G := physicsConstants["g_gravity"].Value
	v := math.Sqrt(G * M * (1 - e) / (a * (1 + e)))

	return types.Number(v)
}

// ════════════════════════════════════════════════════════════════
// ORBITAL MANEUVERS
// ════════════════════════════════════════════════════════════════

// FnHohmannTransfer calculates Δv for Hohmann transfer.
// Returns total Δv for transfer from r1 to r2
// Args: central mass (kg), initial radius (m), final radius (m)
func FnHohmannTransfer(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("hohmanntransfer requires 3 arguments: mass, r1, r2")
	}

	M := args[0].AsFloat()
	r1 := args[1].AsFloat()
	r2 := args[2].AsFloat()

	if r1 <= 0 || r2 <= 0 {
		return types.Error("hohmanntransfer: radii must be positive")
	}

	G := physicsConstants["g_gravity"].Value

	// Initial circular velocity
	v1 := math.Sqrt(G * M / r1)

	// Transfer orbit semi-major axis
	a_transfer := (r1 + r2) / 2

	// Velocity at r1 on transfer orbit
	v_transfer1 := math.Sqrt(G * M * (2/r1 - 1/a_transfer))

	// Velocity at r2 on transfer orbit
	v_transfer2 := math.Sqrt(G * M * (2/r2 - 1/a_transfer))

	// Final circular velocity
	v2 := math.Sqrt(G * M / r2)

	// Total Δv
	dv1 := math.Abs(v_transfer1 - v1)
	dv2 := math.Abs(v2 - v_transfer2)
	totalDv := dv1 + dv2

	return types.Number(totalDv)
}

// FnHohmannTime calculates transfer time for Hohmann transfer.
// t = π√(a³/GM) where a = (r1 + r2)/2
// Args: central mass (kg), initial radius (m), final radius (m)
func FnHohmannTime(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("hohmanntime requires 3 arguments: mass, r1, r2")
	}

	M := args[0].AsFloat()
	r1 := args[1].AsFloat()
	r2 := args[2].AsFloat()

	if r1 <= 0 || r2 <= 0 {
		return types.Error("hohmanntime: radii must be positive")
	}
	if M <= 0 {
		return types.Error("hohmanntime: mass must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	a := (r1 + r2) / 2
	t := math.Pi * math.Sqrt(a*a*a/(G*M))

	return types.Number(t)
}

// FnDeltaVCircularize calculates Δv to circularize orbit at given radius.
// Args: central mass (kg), current velocity (m/s), radius (m)
func FnDeltaVCircularize(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("dvcircularize requires 3 arguments: mass, current_velocity, radius")
	}

	M := args[0].AsFloat()
	v_current := args[1].AsFloat()
	r := args[2].AsFloat()

	if r <= 0 {
		return types.Error("dvcircularize: radius must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	v_circular := math.Sqrt(G * M / r)
	dv := math.Abs(v_circular - v_current)

	return types.Number(dv)
}

// ════════════════════════════════════════════════════════════════
// SPHERE OF INFLUENCE & TIDAL FORCES
// ════════════════════════════════════════════════════════════════

// FnHillSphere calculates Hill sphere radius.
// r_H ≈ a(m/3M)^(1/3)
// Args: orbiting body mass (kg), central body mass (kg), semi-major axis (m)
func FnHillSphere(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("hillsphere requires 3 arguments: m, M, semi_major_axis")
	}

	m := args[0].AsFloat()
	M := args[1].AsFloat()
	a := args[2].AsFloat()

	if M <= 0 {
		return types.Error("hillsphere: central mass must be positive")
	}
	if a <= 0 {
		return types.Error("hillsphere: semi-major axis must be positive")
	}

	rH := a * math.Pow(m/(3*M), 1.0/3.0)
	return types.Number(rH)
}

// FnSphereOfInfluence calculates sphere of influence radius.
// r_SOI ≈ a(m/M)^(2/5)
// Args: orbiting body mass (kg), central body mass (kg), semi-major axis (m)
func FnSphereOfInfluence(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("soi requires 3 arguments: m, M, semi_major_axis")
	}

	m := args[0].AsFloat()
	M := args[1].AsFloat()
	a := args[2].AsFloat()

	if M <= 0 {
		return types.Error("soi: central mass must be positive")
	}
	if a <= 0 {
		return types.Error("soi: semi-major axis must be positive")
	}

	rSOI := a * math.Pow(m/M, 2.0/5.0)
	return types.Number(rSOI)
}

// FnTidalForce calculates tidal force.
// ΔF = 2GMmr/d³ (approximate for small r compared to d)
// Args: central mass (kg), object mass (kg), object size (m), distance (m)
func FnTidalForce(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("tidalforce requires 4 arguments: M, m, object_size, distance")
	}

	M := args[0].AsFloat()
	m := args[1].AsFloat()
	r := args[2].AsFloat()
	d := args[3].AsFloat()

	if d == 0 {
		return types.Error("tidalforce: distance cannot be zero")
	}

	G := physicsConstants["g_gravity"].Value
	F := 2 * G * M * m * r / (d * d * d)

	return types.Number(F)
}

// FnTidalAcceleration calculates tidal acceleration.
// a_tidal = 2GMr/d³
// Args: central mass (kg), object size (m), distance (m)
func FnTidalAcceleration(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("tidalaccel requires 3 arguments: M, object_size, distance")
	}

	M := args[0].AsFloat()
	r := args[1].AsFloat()
	d := args[2].AsFloat()

	if d == 0 {
		return types.Error("tidalaccel: distance cannot be zero")
	}

	G := physicsConstants["g_gravity"].Value
	a := 2 * G * M * r / (d * d * d)

	return types.Number(a)
}

// FnRocheLimit calculates Roche limit.
// d = R(2ρM/ρm)^(1/3) ≈ 2.44 R(ρM/ρm)^(1/3) for fluid body
// Args: primary radius (m), primary density (kg/m³), secondary density (kg/m³)
func FnRocheLimit(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("rochelimit requires 3 arguments: R_primary, density_primary, density_secondary")
	}

	R := args[0].AsFloat()
	rhoM := args[1].AsFloat()
	rhom := args[2].AsFloat()

	if rhom <= 0 {
		return types.Error("rochelimit: secondary density must be positive")
	}

	// Roche limit for fluid body
	d := 2.44 * R * math.Pow(rhoM/rhom, 1.0/3.0)

	return types.Number(d)
}

// ════════════════════════════════════════════════════════════════
// SCHWARZSCHILD & BLACK HOLES
// ════════════════════════════════════════════════════════════════

// FnSchwarzschildRadius calculates Schwarzschild radius.
// r_s = 2GM/c²
// Args: mass (kg)
func FnSchwarzschildRadius(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("schwarzschild requires 1 argument: mass")
	}

	M := args[0].AsFloat()

	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value

	rs := 2 * G * M / (c * c)

	return types.Number(rs)
}

// FnBlackHoleMass calculates black hole mass from Schwarzschild radius.
// M = r_s c² / (2G)
// Args: Schwarzschild radius (m)
func FnBlackHoleMass(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("blackholemass requires 1 argument: schwarzschild_radius")
	}

	rs := args[0].AsFloat()

	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value

	M := rs * c * c / (2 * G)

	return types.Number(M)
}

// FnBlackHoleTemperature calculates Hawking temperature of a black hole.
// T = ℏc³ / (8πGMk_B)
// Args: mass (kg)
func FnBlackHoleTemperature(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("blackholetemp requires 1 argument: mass")
	}

	M := args[0].AsFloat()

	if M <= 0 {
		return types.Error("blackholetemp: mass must be positive")
	}

	hbar := physicsConstants["hbar"].Value
	c := physicsConstants["c"].Value
	G := physicsConstants["g_gravity"].Value
	kB := physicsConstants["kb"].Value

	T := hbar * c * c * c / (8 * math.Pi * G * M * kB)

	return types.Number(T)
}

// FnBlackHoleEntropy calculates Bekenstein-Hawking entropy.
// S = Ak_B c³ / (4Gℏ) where A = 4πr_s²
// Args: mass (kg)
func FnBlackHoleEntropy(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("blackholeentropy requires 1 argument: mass")
	}

	M := args[0].AsFloat()

	if M <= 0 {
		return types.Error("blackholeentropy: mass must be positive")
	}

	hbar := physicsConstants["hbar"].Value
	c := physicsConstants["c"].Value
	G := physicsConstants["g_gravity"].Value
	kB := physicsConstants["kb"].Value

	// Schwarzschild radius
	rs := 2 * G * M / (c * c)

	// Event horizon area
	A := 4 * math.Pi * rs * rs

	// Bekenstein-Hawking entropy
	S := A * kB * c * c * c / (4 * G * hbar)

	return types.Number(S)
}

// ════════════════════════════════════════════════════════════════
// GEOSYNCHRONOUS & SPECIAL ORBITS
// ════════════════════════════════════════════════════════════════

// FnGeosyncRadius calculates geosynchronous orbit radius.
// r = ∛(GMT²/(4π²)) where T = 1 sidereal day
// Args: mass (kg), [rotation period in seconds] (defaults to Earth sidereal day)
func FnGeosyncRadius(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("geosyncradius requires 1-2 arguments: mass, [period]")
	}

	M := args[0].AsFloat()
	T := 86164.0905 // Earth sidereal day in seconds

	if len(args) == 2 {
		T = args[1].AsFloat()
	}

	if M <= 0 {
		return types.Error("geosyncradius: mass must be positive")
	}
	if T <= 0 {
		return types.Error("geosyncradius: period must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	r := math.Pow(G*M*T*T/(4*math.Pi*math.Pi), 1.0/3.0)

	return types.Number(r)
}

// FnGeosyncAltitude calculates geosynchronous altitude above surface.
// Args: mass (kg), surface radius (m), [rotation period]
func FnGeosyncAltitude(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("geosyncalt requires 2-3 arguments: mass, surface_radius, [period]")
	}

	M := args[0].AsFloat()
	R := args[1].AsFloat()
	T := 86164.0905

	if len(args) == 3 {
		T = args[2].AsFloat()
	}

	if M <= 0 {
		return types.Error("geosyncalt: mass must be positive")
	}
	if T <= 0 {
		return types.Error("geosyncalt: period must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	r := math.Pow(G*M*T*T/(4*math.Pi*math.Pi), 1.0/3.0)

	altitude := r - R
	return types.Number(altitude)
}

// FnLagrangeL1 calculates L1 Lagrange point distance from smaller body.
// Approximate: r ≈ R(m/(3M))^(1/3)
// Args: smaller mass (kg), larger mass (kg), separation (m)
func FnLagrangeL1(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("lagrangel1 requires 3 arguments: m, M, separation")
	}

	m := args[0].AsFloat()
	M := args[1].AsFloat()
	R := args[2].AsFloat()

	if M <= 0 {
		return types.Error("lagrangel1: larger mass must be positive")
	}

	r := R * math.Pow(m/(3*M), 1.0/3.0)
	return types.Number(r)
}

// ════════════════════════════════════════════════════════════════
// GRAVITATIONAL WAVES (basic)
// ════════════════════════════════════════════════════════════════

// FnGravWaveLuminosity calculates gravitational wave luminosity for binary.
// L = (32/5)(G⁴/c⁵)(m₁m₂)²(m₁+m₂)/r⁵
// Args: m1 (kg), m2 (kg), separation (m)
func FnGravWaveLuminosity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("gravwaveluminosity requires 3 arguments: m1, m2, separation")
	}

	m1 := args[0].AsFloat()
	m2 := args[1].AsFloat()
	r := args[2].AsFloat()

	if r <= 0 {
		return types.Error("gravwaveluminosity: separation must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	c := physicsConstants["c"].Value

	G4 := G * G * G * G
	c5 := c * c * c * c * c
	r5 := r * r * r * r * r

	L := (32.0 / 5.0) * (G4 / c5) * (m1 * m2) * (m1 * m2) * (m1 + m2) / r5

	return types.Number(L)
}

// FnGravWaveFrequency calculates gravitational wave frequency for circular binary.
// f = (1/π)√(G(m₁+m₂)/r³)
// Args: m1 (kg), m2 (kg), separation (m)
func FnGravWaveFrequency(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("gravwavefreq requires 3 arguments: m1, m2, separation")
	}

	m1 := args[0].AsFloat()
	m2 := args[1].AsFloat()
	r := args[2].AsFloat()

	if r <= 0 {
		return types.Error("gravwavefreq: separation must be positive")
	}

	G := physicsConstants["g_gravity"].Value
	f := (1 / math.Pi) * math.Sqrt(G*(m1+m2)/(r*r*r))

	return types.Number(f)
}
