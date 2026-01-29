// internal/eval/functions_physics_em.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// ELECTRIC FORCE & FIELD
// ════════════════════════════════════════════════════════════════

// FnCoulombForce calculates electrostatic force between two charges.
// F = ke × q₁q₂/r²
// Args: charge1 (C), charge2 (C), distance (m)
func FnCoulombForce(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("coulombforce requires 3 arguments: q1, q2, distance")
	}

	q1 := args[0].AsFloat()
	q2 := args[1].AsFloat()
	r := args[2].AsFloat()

	if r == 0 {
		return types.Error("coulombforce: distance cannot be zero")
	}

	ke := physicsConstants["ke"].Value
	F := ke * q1 * q2 / (r * r)

	return types.Number(F)
}

// FnElectricField calculates electric field from a point charge.
// E = ke × q/r²
// Args: charge (C), distance (m)
func FnElectricField(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("electricfield requires 2 arguments: charge, distance")
	}

	q := args[0].AsFloat()
	r := args[1].AsFloat()

	if r == 0 {
		return types.Error("electricfield: distance cannot be zero")
	}

	ke := physicsConstants["ke"].Value
	E := ke * q / (r * r)

	return types.Number(E)
}

// FnElectricFieldPlate calculates electric field from infinite charged plate.
// E = σ/(2ε₀) for single plate, E = σ/ε₀ between parallel plates
// Args: surface charge density (C/m²), [type: "single" or "parallel"]
func FnElectricFieldPlate(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("electricfieldplate requires 1-2 arguments: sigma, [type]")
	}

	sigma := args[0].AsFloat()
	plateType := "single"
	if len(args) == 2 {
		plateType = args[1].AsString()
	}

	epsilon0 := physicsConstants["epsilon0"].Value

	var E float64
	if plateType == "parallel" || plateType == "p" {
		E = sigma / epsilon0
	} else {
		E = sigma / (2 * epsilon0)
	}

	return types.Number(E)
}

// FnElectricPotential calculates electric potential from a point charge.
// V = ke × q/r
// Args: charge (C), distance (m)
func FnElectricPotential(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("electricpotential requires 2 arguments: charge, distance")
	}

	q := args[0].AsFloat()
	r := args[1].AsFloat()

	if r == 0 {
		return types.Error("electricpotential: distance cannot be zero")
	}

	ke := physicsConstants["ke"].Value
	V := ke * q / r

	return types.Number(V)
}

// FnElectricPE calculates electric potential energy.
// U = ke × q₁q₂/r
// Args: charge1 (C), charge2 (C), distance (m)
func FnElectricPE(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("electricpe requires 3 arguments: q1, q2, distance")
	}

	q1 := args[0].AsFloat()
	q2 := args[1].AsFloat()
	r := args[2].AsFloat()

	if r == 0 {
		return types.Error("electricpe: distance cannot be zero")
	}

	ke := physicsConstants["ke"].Value
	U := ke * q1 * q2 / r

	return types.Number(U)
}

// FnElectricFlux calculates electric flux.
// Φ = E·A·cos(θ)
// Args: electric field (N/C), area (m²), [angle in degrees] (defaults to 0)
func FnElectricFlux(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("electricflux requires 2-3 arguments: E, A, [angle_deg]")
	}

	E := args[0].AsFloat()
	A := args[1].AsFloat()
	theta := 0.0

	if len(args) == 3 {
		theta = args[2].AsFloat() * math.Pi / 180
	}

	flux := E * A * math.Cos(theta)
	return types.Number(flux)
}

// FnGaussLaw calculates enclosed charge from electric flux.
// q = ε₀ × Φ
// Args: electric flux (N·m²/C)
func FnGaussLaw(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("gausslaw requires 1 argument: electric_flux")
	}

	flux := args[0].AsFloat()
	epsilon0 := physicsConstants["epsilon0"].Value

	q := epsilon0 * flux
	return types.Number(q)
}

// ════════════════════════════════════════════════════════════════
// CAPACITANCE
// ════════════════════════════════════════════════════════════════

// FnCapacitance calculates capacitance.
// C = Q/V
// Args: charge (C), voltage (V)
func FnCapacitance(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("capacitance requires 2 arguments: charge, voltage")
	}

	Q := args[0].AsFloat()
	V := args[1].AsFloat()

	if V == 0 {
		return types.Error("capacitance: voltage cannot be zero")
	}

	return types.Number(Q / V)
}

// FnParallelPlateCapacitance calculates capacitance of parallel plate capacitor.
// C = ε₀εᵣA/d
// Args: area (m²), separation (m), [relative permittivity] (defaults to 1)
func FnParallelPlateCapacitance(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("parallelplatecap requires 2-3 arguments: area, separation, [epsilon_r]")
	}

	A := args[0].AsFloat()
	d := args[1].AsFloat()
	epsilonR := 1.0

	if len(args) == 3 {
		epsilonR = args[2].AsFloat()
	}

	if d == 0 {
		return types.Error("parallelplatecap: separation cannot be zero")
	}

	epsilon0 := physicsConstants["epsilon0"].Value
	C := epsilon0 * epsilonR * A / d

	return types.Number(C)
}

// FnCylindricalCapacitance calculates capacitance of cylindrical capacitor.
// C = 2πε₀εᵣL / ln(b/a)
// Args: length (m), inner radius (m), outer radius (m), [epsilon_r]
func FnCylindricalCapacitance(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("cylindricalcap requires 3-4 arguments: length, r_inner, r_outer, [epsilon_r]")
	}

	L := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()
	epsilonR := 1.0

	if len(args) == 4 {
		epsilonR = args[3].AsFloat()
	}

	if a <= 0 || b <= 0 {
		return types.Error("cylindricalcap: radii must be positive")
	}
	if b <= a {
		return types.Error("cylindricalcap: outer radius must be greater than inner")
	}

	epsilon0 := physicsConstants["epsilon0"].Value
	C := 2 * math.Pi * epsilon0 * epsilonR * L / math.Log(b/a)

	return types.Number(C)
}

// FnSphericalCapacitance calculates capacitance of spherical capacitor.
// C = 4πε₀εᵣab / (b-a)
// Args: inner radius (m), outer radius (m), [epsilon_r]
func FnSphericalCapacitance(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("sphericalcap requires 2-3 arguments: r_inner, r_outer, [epsilon_r]")
	}

	a := args[0].AsFloat()
	b := args[1].AsFloat()
	epsilonR := 1.0

	if len(args) == 3 {
		epsilonR = args[2].AsFloat()
	}

	if a <= 0 || b <= 0 {
		return types.Error("sphericalcap: radii must be positive")
	}
	if b <= a {
		return types.Error("sphericalcap: outer radius must be greater than inner")
	}

	epsilon0 := physicsConstants["epsilon0"].Value
	C := 4 * math.Pi * epsilon0 * epsilonR * a * b / (b - a)

	return types.Number(C)
}

// FnCapacitorEnergy calculates energy stored in capacitor.
// U = ½CV² = ½QV = Q²/(2C)
// Args: capacitance (F), voltage (V)
func FnCapacitorEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("capacitorenergy requires 2 arguments: capacitance, voltage")
	}

	C := args[0].AsFloat()
	V := args[1].AsFloat()

	U := 0.5 * C * V * V
	return types.Number(U)
}

// FnCapacitorSeries calculates equivalent capacitance in series.
// 1/C_eq = 1/C₁ + 1/C₂ + ...
// Args: C1, C2, ...
func FnCapacitorSeries(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("capacitorseries requires at least 2 arguments")
	}

	sum := 0.0
	for _, arg := range args {
		C := arg.AsFloat()
		if C == 0 {
			return types.Error("capacitorseries: capacitance cannot be zero")
		}
		sum += 1 / C
	}

	return types.Number(1 / sum)
}

// FnCapacitorParallel calculates equivalent capacitance in parallel.
// C_eq = C₁ + C₂ + ...
// Args: C1, C2, ...
func FnCapacitorParallel(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("capacitorparallel requires at least 2 arguments")
	}

	sum := 0.0
	for _, arg := range args {
		sum += arg.AsFloat()
	}

	return types.Number(sum)
}

// ════════════════════════════════════════════════════════════════
// CURRENT & RESISTANCE
// ════════════════════════════════════════════════════════════════

// FnOhmsLawV calculates voltage.
// V = IR
// Args: current (A), resistance (Ω)
func FnOhmsLawV(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("ohmslawv requires 2 arguments: current, resistance")
	}

	I := args[0].AsFloat()
	R := args[1].AsFloat()

	return types.Number(I * R)
}

// FnOhmsLawI calculates current.
// I = V/R
// Args: voltage (V), resistance (Ω)
func FnOhmsLawI(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("ohmslawi requires 2 arguments: voltage, resistance")
	}

	V := args[0].AsFloat()
	R := args[1].AsFloat()

	if R == 0 {
		return types.Error("ohmslawi: resistance cannot be zero")
	}

	return types.Number(V / R)
}

// FnOhmsLawR calculates resistance.
// R = V/I
// Args: voltage (V), current (A)
func FnOhmsLawR(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("ohmslawv requires 2 arguments: voltage, current")
	}

	V := args[0].AsFloat()
	I := args[1].AsFloat()

	if I == 0 {
		return types.Error("ohmslawv: current cannot be zero")
	}

	return types.Number(V / I)
}

// FnResistivity calculates resistance from resistivity.
// R = ρL/A
// Args: resistivity (Ω·m), length (m), cross-sectional area (m²)
func FnResistivity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("resistivity requires 3 arguments: rho, length, area")
	}

	rho := args[0].AsFloat()
	L := args[1].AsFloat()
	A := args[2].AsFloat()

	if A == 0 {
		return types.Error("resistivity: area cannot be zero")
	}

	return types.Number(rho * L / A)
}

// FnResistorSeries calculates equivalent resistance in series.
// R_eq = R₁ + R₂ + ...
// Args: R1, R2, ...
func FnResistorSeries(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("resistorseries requires at least 2 arguments")
	}

	sum := 0.0
	for _, arg := range args {
		sum += arg.AsFloat()
	}

	return types.Number(sum)
}

// FnResistorParallel calculates equivalent resistance in parallel.
// 1/R_eq = 1/R₁ + 1/R₂ + ...
// Args: R1, R2, ...
func FnResistorParallel(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("resistorparallel requires at least 2 arguments")
	}

	sum := 0.0
	for _, arg := range args {
		R := arg.AsFloat()
		if R == 0 {
			return types.Error("resistorparallel: resistance cannot be zero")
		}
		sum += 1 / R
	}

	return types.Number(1 / sum)
}

// FnConductance calculates conductance.
// G = 1/R
// Args: resistance (Ω)
func FnConductance(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("conductance requires 1 argument: resistance")
	}

	R := args[0].AsFloat()

	if R == 0 {
		return types.Error("conductance: resistance cannot be zero")
	}

	return types.Number(1 / R)
}

// FnCurrentDensity calculates current density.
// J = I/A
// Args: current (A), cross-sectional area (m²)
func FnCurrentDensity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("currentdensity requires 2 arguments: current, area")
	}

	I := args[0].AsFloat()
	A := args[1].AsFloat()

	if A == 0 {
		return types.Error("currentdensity: area cannot be zero")
	}

	return types.Number(I / A)
}

// FnDriftVelocity calculates drift velocity of electrons.
// v_d = I/(nAe)
// Args: current (A), carrier density (m⁻³), area (m²)
func FnDriftVelocity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("driftvelocity requires 3 arguments: current, carrier_density, area")
	}

	I := args[0].AsFloat()
	n := args[1].AsFloat()
	A := args[2].AsFloat()

	if n*A == 0 {
		return types.Error("driftvelocity: n × A cannot be zero")
	}

	e := physicsConstants["e_charge"].Value
	vd := I / (n * A * e)

	return types.Number(vd)
}

// ════════════════════════════════════════════════════════════════
// ELECTRIC POWER
// ════════════════════════════════════════════════════════════════

// FnElectricPower calculates electric power.
// P = IV = I²R = V²/R
// Args: voltage (V), current (A)
func FnElectricPower(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("electricpower requires 2 arguments: voltage, current")
	}

	V := args[0].AsFloat()
	I := args[1].AsFloat()

	return types.Number(V * I)
}

// FnPowerFromResistance calculates power from current and resistance.
// P = I²R
// Args: current (A), resistance (Ω)
func FnPowerFromResistance(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("powerfromresistance requires 2 arguments: current, resistance")
	}

	I := args[0].AsFloat()
	R := args[1].AsFloat()

	return types.Number(I * I * R)
}

// FnPowerFromVoltage calculates power from voltage and resistance.
// P = V²/R
// Args: voltage (V), resistance (Ω)
func FnPowerFromVoltage(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("powerfromvoltage requires 2 arguments: voltage, resistance")
	}

	V := args[0].AsFloat()
	R := args[1].AsFloat()

	if R == 0 {
		return types.Error("powerfromvoltage: resistance cannot be zero")
	}

	return types.Number(V * V / R)
}

// FnElectricEnergy calculates electric energy.
// E = Pt
// Args: power (W), time (s)
func FnElectricEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("electricenergy requires 2 arguments: power, time")
	}

	P := args[0].AsFloat()
	t := args[1].AsFloat()

	return types.Number(P * t)
}

// ════════════════════════════════════════════════════════════════
// RC CIRCUITS
// ════════════════════════════════════════════════════════════════

// FnRCTimeConstant calculates RC time constant.
// τ = RC
// Args: resistance (Ω), capacitance (F)
func FnRCTimeConstant(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rctimeconstant requires 2 arguments: resistance, capacitance")
	}

	R := args[0].AsFloat()
	C := args[1].AsFloat()

	return types.Number(R * C)
}

// FnRCCharging calculates voltage during capacitor charging.
// V(t) = V₀(1 - e^(-t/RC))
// Args: V0 (V), R (Ω), C (F), time (s)
func FnRCCharging(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("rccharging requires 4 arguments: V0, R, C, time")
	}

	V0 := args[0].AsFloat()
	R := args[1].AsFloat()
	C := args[2].AsFloat()
	t := args[3].AsFloat()

	if R*C == 0 {
		return types.Error("rccharging: RC cannot be zero")
	}

	V := V0 * (1 - math.Exp(-t/(R*C)))
	return types.Number(V)
}

// FnRCDischarging calculates voltage during capacitor discharging.
// V(t) = V₀ × e^(-t/RC)
// Args: V0 (V), R (Ω), C (F), time (s)
func FnRCDischarging(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("rcdischarging requires 4 arguments: V0, R, C, time")
	}

	V0 := args[0].AsFloat()
	R := args[1].AsFloat()
	C := args[2].AsFloat()
	t := args[3].AsFloat()

	if R*C == 0 {
		return types.Error("rcdischarging: RC cannot be zero")
	}

	V := V0 * math.Exp(-t/(R*C))
	return types.Number(V)
}

// ════════════════════════════════════════════════════════════════
// MAGNETIC FIELD
// ════════════════════════════════════════════════════════════════

// FnMagneticForceCharge calculates magnetic force on a moving charge.
// F = qvB sin(θ)
// Args: charge (C), velocity (m/s), magnetic field (T), [angle in degrees] (defaults to 90)
func FnMagneticForceCharge(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("magneticforcecharge requires 3-4 arguments: charge, velocity, B, [angle_deg]")
	}

	q := args[0].AsFloat()
	v := args[1].AsFloat()
	B := args[2].AsFloat()
	theta := 90.0

	if len(args) == 4 {
		theta = args[3].AsFloat()
	}

	thetaRad := theta * math.Pi / 180
	F := math.Abs(q) * v * B * math.Sin(thetaRad)

	return types.Number(F)
}

// FnMagneticForceWire calculates magnetic force on a current-carrying wire.
// F = BIL sin(θ)
// Args: magnetic field (T), current (A), length (m), [angle in degrees] (defaults to 90)
func FnMagneticForceWire(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("magneticforcewire requires 3-4 arguments: B, I, length, [angle_deg]")
	}

	B := args[0].AsFloat()
	I := args[1].AsFloat()
	L := args[2].AsFloat()
	theta := 90.0

	if len(args) == 4 {
		theta = args[3].AsFloat()
	}

	thetaRad := theta * math.Pi / 180
	F := B * I * L * math.Sin(thetaRad)

	return types.Number(F)
}

// FnCyclotronRadius calculates cyclotron radius.
// r = mv/(qB)
// Args: mass (kg), velocity (m/s), charge (C), magnetic field (T)
func FnCyclotronRadius(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("cyclotronradius requires 4 arguments: mass, velocity, charge, B")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()
	q := args[2].AsFloat()
	B := args[3].AsFloat()

	if q*B == 0 {
		return types.Error("cyclotronradius: charge × B cannot be zero")
	}

	r := m * v / (math.Abs(q) * B)
	return types.Number(r)
}

// FnCyclotronFrequency calculates cyclotron frequency.
// f = qB/(2πm)
// Args: charge (C), magnetic field (T), mass (kg)
func FnCyclotronFrequency(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("cyclotronfreq requires 3 arguments: charge, B, mass")
	}

	q := args[0].AsFloat()
	B := args[1].AsFloat()
	m := args[2].AsFloat()

	if m == 0 {
		return types.Error("cyclotronfreq: mass cannot be zero")
	}

	f := math.Abs(q) * B / (2 * math.Pi * m)
	return types.Number(f)
}

// FnMagneticFieldWire calculates magnetic field from a long straight wire.
// B = μ₀I/(2πr)
// Args: current (A), distance (m)
func FnMagneticFieldWire(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("magneticfieldwire requires 2 arguments: current, distance")
	}

	I := args[0].AsFloat()
	r := args[1].AsFloat()

	if r == 0 {
		return types.Error("magneticfieldwire: distance cannot be zero")
	}

	mu0 := physicsConstants["mu0"].Value
	B := mu0 * I / (2 * math.Pi * r)

	return types.Number(B)
}

// FnMagneticFieldLoop calculates magnetic field at center of circular loop.
// B = μ₀I/(2R)
// Args: current (A), radius (m)
func FnMagneticFieldLoop(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("magneticfieldloop requires 2 arguments: current, radius")
	}

	I := args[0].AsFloat()
	R := args[1].AsFloat()

	if R == 0 {
		return types.Error("magneticfieldloop: radius cannot be zero")
	}

	mu0 := physicsConstants["mu0"].Value
	B := mu0 * I / (2 * R)

	return types.Number(B)
}

// FnMagneticFieldSolenoid calculates magnetic field inside a solenoid.
// B = μ₀nI where n = N/L
// Args: number of turns, length (m), current (A)
func FnMagneticFieldSolenoid(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("magneticfieldsolenoid requires 3 arguments: turns, length, current")
	}

	N := args[0].AsFloat()
	L := args[1].AsFloat()
	I := args[2].AsFloat()

	if L == 0 {
		return types.Error("magneticfieldsolenoid: length cannot be zero")
	}

	mu0 := physicsConstants["mu0"].Value
	n := N / L
	B := mu0 * n * I

	return types.Number(B)
}

// FnMagneticFieldToroid calculates magnetic field inside a toroid.
// B = μ₀NI/(2πr)
// Args: number of turns, current (A), radius (m)
func FnMagneticFieldToroid(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("magneticfieldtoroid requires 3 arguments: turns, current, radius")
	}

	N := args[0].AsFloat()
	I := args[1].AsFloat()
	r := args[2].AsFloat()

	if r == 0 {
		return types.Error("magneticfieldtoroid: radius cannot be zero")
	}

	mu0 := physicsConstants["mu0"].Value
	B := mu0 * N * I / (2 * math.Pi * r)

	return types.Number(B)
}

// ════════════════════════════════════════════════════════════════
// MAGNETIC FLUX & INDUCTANCE
// ════════════════════════════════════════════════════════════════

// FnMagneticFlux calculates magnetic flux.
// Φ = BA cos(θ)
// Args: magnetic field (T), area (m²), [angle in degrees] (defaults to 0)
func FnMagneticFlux(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("magneticflux requires 2-3 arguments: B, area, [angle_deg]")
	}

	B := args[0].AsFloat()
	A := args[1].AsFloat()
	theta := 0.0

	if len(args) == 3 {
		theta = args[2].AsFloat() * math.Pi / 180
	}

	flux := B * A * math.Cos(theta)
	return types.Number(flux)
}

// FnInducedEMF calculates induced EMF (Faraday's law).
// ε = -N × dΦ/dt
// Args: number of turns, change in flux (Wb), change in time (s)
func FnInducedEMF(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("inducedemf requires 3 arguments: turns, delta_flux, delta_time")
	}

	N := args[0].AsFloat()
	dPhi := args[1].AsFloat()
	dt := args[2].AsFloat()

	if dt == 0 {
		return types.Error("inducedemf: delta_time cannot be zero")
	}

	emf := -N * dPhi / dt
	return types.Number(emf)
}

// FnMotionalEMF calculates motional EMF.
// ε = BLv
// Args: magnetic field (T), length (m), velocity (m/s)
func FnMotionalEMF(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("motionalemf requires 3 arguments: B, length, velocity")
	}

	B := args[0].AsFloat()
	L := args[1].AsFloat()
	v := args[2].AsFloat()

	return types.Number(B * L * v)
}

// FnInductance calculates inductance.
// L = NΦ/I
// Args: number of turns, magnetic flux (Wb), current (A)
func FnInductance(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("inductance requires 3 arguments: turns, flux, current")
	}

	N := args[0].AsFloat()
	Phi := args[1].AsFloat()
	I := args[2].AsFloat()

	if I == 0 {
		return types.Error("inductance: current cannot be zero")
	}

	return types.Number(N * Phi / I)
}

// FnSolenoidInductance calculates inductance of a solenoid.
// L = μ₀N²A/L
// Args: number of turns, area (m²), length (m)
func FnSolenoidInductance(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("solenoidinductance requires 3 arguments: turns, area, length")
	}

	N := args[0].AsFloat()
	A := args[1].AsFloat()
	length := args[2].AsFloat()

	if length == 0 {
		return types.Error("solenoidinductance: length cannot be zero")
	}

	mu0 := physicsConstants["mu0"].Value
	L := mu0 * N * N * A / length

	return types.Number(L)
}

// FnInductorEnergy calculates energy stored in an inductor.
// U = ½LI²
// Args: inductance (H), current (A)
func FnInductorEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("inductorenergy requires 2 arguments: inductance, current")
	}

	L := args[0].AsFloat()
	I := args[1].AsFloat()

	return types.Number(0.5 * L * I * I)
}

// ════════════════════════════════════════════════════════════════
// RL CIRCUITS
// ════════════════════════════════════════════════════════════════

// FnRLTimeConstant calculates RL time constant.
// τ = L/R
// Args: inductance (H), resistance (Ω)
func FnRLTimeConstant(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rltimeconstant requires 2 arguments: inductance, resistance")
	}

	L := args[0].AsFloat()
	R := args[1].AsFloat()

	if R == 0 {
		return types.Error("rltimeconstant: resistance cannot be zero")
	}

	return types.Number(L / R)
}

// FnRLCurrentGrowth calculates current during inductor energizing.
// I(t) = (V/R)(1 - e^(-Rt/L))
// Args: voltage (V), resistance (Ω), inductance (H), time (s)
func FnRLCurrentGrowth(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("rlcurrentgrowth requires 4 arguments: voltage, R, L, time")
	}

	V := args[0].AsFloat()
	R := args[1].AsFloat()
	L := args[2].AsFloat()
	t := args[3].AsFloat()

	if R == 0 {
		return types.Error("rlcurrentgrowth: resistance cannot be zero")
	}
	if L == 0 {
		return types.Error("rlcurrentgrowth: inductance cannot be zero")
	}

	I := (V / R) * (1 - math.Exp(-R*t/L))
	return types.Number(I)
}

// FnRLCurrentDecay calculates current during inductor de-energizing.
// I(t) = I₀ × e^(-Rt/L)
// Args: initial current (A), resistance (Ω), inductance (H), time (s)
func FnRLCurrentDecay(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("rlcurrentdecay requires 4 arguments: I0, R, L, time")
	}

	I0 := args[0].AsFloat()
	R := args[1].AsFloat()
	L := args[2].AsFloat()
	t := args[3].AsFloat()

	if L == 0 {
		return types.Error("rlcurrentdecay: inductance cannot be zero")
	}

	I := I0 * math.Exp(-R*t/L)
	return types.Number(I)
}

// ════════════════════════════════════════════════════════════════
// AC CIRCUITS
// ════════════════════════════════════════════════════════════════

// FnCapacitiveReactance calculates capacitive reactance.
// Xc = 1/(2πfC)
// Args: frequency (Hz), capacitance (F)
func FnCapacitiveReactance(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("capacitivereactance requires 2 arguments: frequency, capacitance")
	}

	f := args[0].AsFloat()
	C := args[1].AsFloat()

	if f*C == 0 {
		return types.Error("capacitivereactance: frequency × capacitance cannot be zero")
	}

	Xc := 1 / (2 * math.Pi * f * C)
	return types.Number(Xc)
}

// FnInductiveReactance calculates inductive reactance.
// XL = 2πfL
// Args: frequency (Hz), inductance (H)
func FnInductiveReactance(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("inductivereactance requires 2 arguments: frequency, inductance")
	}

	f := args[0].AsFloat()
	L := args[1].AsFloat()

	XL := 2 * math.Pi * f * L
	return types.Number(XL)
}

// FnImpedance calculates impedance of RLC circuit.
// Z = √(R² + (XL - Xc)²)
// Args: resistance (Ω), inductive reactance (Ω), capacitive reactance (Ω)
func FnImpedance(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("impedance requires 3 arguments: R, XL, Xc")
	}

	R := args[0].AsFloat()
	XL := args[1].AsFloat()
	Xc := args[2].AsFloat()

	Z := math.Sqrt(R*R + (XL-Xc)*(XL-Xc))
	return types.Number(Z)
}

// FnPhaseAngle calculates phase angle in AC circuit.
// φ = arctan((XL - Xc)/R)
// Args: resistance (Ω), inductive reactance (Ω), capacitive reactance (Ω)
func FnPhaseAngle(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("phaseangle requires 3 arguments: R, XL, Xc")
	}

	R := args[0].AsFloat()
	XL := args[1].AsFloat()
	Xc := args[2].AsFloat()

	if R == 0 {
		if XL > Xc {
			return types.Number(90)
		} else if XL < Xc {
			return types.Number(-90)
		}
		return types.Number(0)
	}

	phi := math.Atan((XL-Xc)/R) * 180 / math.Pi
	return types.Number(phi)
}

// FnResonantFrequency calculates resonant frequency of LC circuit.
// f₀ = 1/(2π√(LC))
// Args: inductance (H), capacitance (F)
func FnResonantFrequency(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("resonantfreq requires 2 arguments: inductance, capacitance")
	}

	L := args[0].AsFloat()
	C := args[1].AsFloat()

	if L <= 0 || C <= 0 {
		return types.Error("resonantfreq: L and C must be positive")
	}

	f0 := 1 / (2 * math.Pi * math.Sqrt(L*C))
	return types.Number(f0)
}

// FnQualityFactor calculates Q factor of RLC circuit.
// Q = (1/R)√(L/C)
// Args: resistance (Ω), inductance (H), capacitance (F)
func FnQualityFactor(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("qualityfactor requires 3 arguments: R, L, C")
	}

	R := args[0].AsFloat()
	L := args[1].AsFloat()
	C := args[2].AsFloat()

	if R == 0 {
		return types.Error("qualityfactor: resistance cannot be zero")
	}
	if L <= 0 || C <= 0 {
		return types.Error("qualityfactor: L and C must be positive")
	}

	Q := (1 / R) * math.Sqrt(L/C)
	return types.Number(Q)
}

// FnRMSVoltage calculates RMS voltage from peak voltage.
// V_rms = V_peak / √2
// Args: peak voltage (V)
func FnRMSVoltage(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("rmsvoltage requires 1 argument: peak_voltage")
	}

	Vpeak := args[0].AsFloat()
	return types.Number(Vpeak / math.Sqrt2)
}

// FnPeakVoltage calculates peak voltage from RMS voltage.
// V_peak = V_rms × √2
// Args: RMS voltage (V)
func FnPeakVoltage(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("peakvoltage requires 1 argument: rms_voltage")
	}

	Vrms := args[0].AsFloat()
	return types.Number(Vrms * math.Sqrt2)
}

// FnPowerFactor calculates power factor.
// PF = cos(φ) = R/Z
// Args: resistance (Ω), impedance (Ω)
func FnPowerFactor(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("powerfactor requires 2 arguments: resistance, impedance")
	}

	R := args[0].AsFloat()
	Z := args[1].AsFloat()

	if Z == 0 {
		return types.Error("powerfactor: impedance cannot be zero")
	}

	return types.Number(R / Z)
}

// FnAveragePowerAC calculates average power in AC circuit.
// P_avg = V_rms × I_rms × cos(φ)
// Args: V_rms (V), I_rms (A), phase angle (degrees)
func FnAveragePowerAC(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("avgpowerac requires 3 arguments: V_rms, I_rms, phase_angle_deg")
	}

	Vrms := args[0].AsFloat()
	Irms := args[1].AsFloat()
	phi := args[2].AsFloat() * math.Pi / 180

	P := Vrms * Irms * math.Cos(phi)
	return types.Number(P)
}

// ════════════════════════════════════════════════════════════════
// TRANSFORMERS
// ════════════════════════════════════════════════════════════════

// FnTransformerVoltage calculates secondary voltage.
// V₂/V₁ = N₂/N₁
// Args: primary voltage (V), primary turns, secondary turns
func FnTransformerVoltage(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("transformervoltage requires 3 arguments: V1, N1, N2")
	}

	V1 := args[0].AsFloat()
	N1 := args[1].AsFloat()
	N2 := args[2].AsFloat()

	if N1 == 0 {
		return types.Error("transformervoltage: primary turns cannot be zero")
	}

	V2 := V1 * N2 / N1
	return types.Number(V2)
}

// FnTransformerCurrent calculates secondary current (ideal transformer).
// I₂/I₁ = N₁/N₂
// Args: primary current (A), primary turns, secondary turns
func FnTransformerCurrent(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("transformercurrent requires 3 arguments: I1, N1, N2")
	}

	I1 := args[0].AsFloat()
	N1 := args[1].AsFloat()
	N2 := args[2].AsFloat()

	if N2 == 0 {
		return types.Error("transformercurrent: secondary turns cannot be zero")
	}

	I2 := I1 * N1 / N2
	return types.Number(I2)
}

// FnTransformerRatio calculates turns ratio.
// n = N₂/N₁ = V₂/V₁
// Args: primary turns, secondary turns
func FnTransformerRatio(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("transformerratio requires 2 arguments: N1, N2")
	}

	N1 := args[0].AsFloat()
	N2 := args[1].AsFloat()

	if N1 == 0 {
		return types.Error("transformerratio: primary turns cannot be zero")
	}

	return types.Number(N2 / N1)
}
