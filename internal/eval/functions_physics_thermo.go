// internal/eval/functions_physics_thermo.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// IDEAL GAS LAW
// ════════════════════════════════════════════════════════════════

// FnIdealGasPressure calculates pressure using ideal gas law.
// P = nRT/V
// Args: moles (n), temperature (K), volume (m³)
func FnIdealGasPressure(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("idealgaspressure requires 3 arguments: n, T, V")
	}

	n := args[0].AsFloat()
	T := args[1].AsFloat()
	V := args[2].AsFloat()

	if V == 0 {
		return types.Error("idealgaspressure: volume cannot be zero")
	}

	R := physicsConstants["r_gas"].Value
	P := n * R * T / V

	return types.Number(P)
}

// FnIdealGasVolume calculates volume using ideal gas law.
// V = nRT/P
// Args: moles (n), temperature (K), pressure (Pa)
func FnIdealGasVolume(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("idealgasvolume requires 3 arguments: n, T, P")
	}

	n := args[0].AsFloat()
	T := args[1].AsFloat()
	P := args[2].AsFloat()

	if P == 0 {
		return types.Error("idealgasvolume: pressure cannot be zero")
	}

	R := physicsConstants["r_gas"].Value
	V := n * R * T / P

	return types.Number(V)
}

// FnIdealGasTemp calculates temperature using ideal gas law.
// T = PV/(nR)
// Args: pressure (Pa), volume (m³), moles (n)
func FnIdealGasTemp(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("idealgastemp requires 3 arguments: P, V, n")
	}

	P := args[0].AsFloat()
	V := args[1].AsFloat()
	n := args[2].AsFloat()

	if n == 0 {
		return types.Error("idealgastemp: moles cannot be zero")
	}

	R := physicsConstants["r_gas"].Value
	T := P * V / (n * R)

	return types.Number(T)
}

// FnIdealGasMoles calculates moles using ideal gas law.
// n = PV/(RT)
// Args: pressure (Pa), volume (m³), temperature (K)
func FnIdealGasMoles(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("idealgasmoles requires 3 arguments: P, V, T")
	}

	P := args[0].AsFloat()
	V := args[1].AsFloat()
	T := args[2].AsFloat()

	if T == 0 {
		return types.Error("idealgasmoles: temperature cannot be zero")
	}

	R := physicsConstants["r_gas"].Value
	n := P * V / (R * T)

	return types.Number(n)
}

// FnIdealGasDensity calculates density of ideal gas.
// ρ = PM/(RT)
// Args: pressure (Pa), molar mass (kg/mol), temperature (K)
func FnIdealGasDensity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("idealgasdensity requires 3 arguments: P, M, T")
	}

	P := args[0].AsFloat()
	M := args[1].AsFloat()
	T := args[2].AsFloat()

	if T == 0 {
		return types.Error("idealgasdensity: temperature cannot be zero")
	}

	R := physicsConstants["r_gas"].Value
	rho := P * M / (R * T)

	return types.Number(rho)
}

// ════════════════════════════════════════════════════════════════
// HEAT TRANSFER
// ════════════════════════════════════════════════════════════════

// FnHeat calculates heat transfer.
// Q = mcΔT
// Args: mass (kg), specific heat (J/(kg·K)), temperature change (K)
func FnHeat(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("heat requires 3 arguments: mass, specific_heat, delta_T")
	}

	m := args[0].AsFloat()
	c := args[1].AsFloat()
	deltaT := args[2].AsFloat()

	return types.Number(m * c * deltaT)
}

// FnHeatMoles calculates heat transfer using molar heat capacity.
// Q = nCΔT
// Args: moles (n), molar heat capacity (J/(mol·K)), temperature change (K)
func FnHeatMoles(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("heatmoles requires 3 arguments: n, C, delta_T")
	}

	n := args[0].AsFloat()
	C := args[1].AsFloat()
	deltaT := args[2].AsFloat()

	return types.Number(n * C * deltaT)
}

// FnFinalTemp calculates final temperature after heat transfer.
// T_f = T_i + Q/(mc)
// Args: initial temp (K), heat added (J), mass (kg), specific heat (J/(kg·K))
func FnFinalTemp(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("finaltemp requires 4 arguments: T_initial, Q, mass, specific_heat")
	}

	Ti := args[0].AsFloat()
	Q := args[1].AsFloat()
	m := args[2].AsFloat()
	c := args[3].AsFloat()

	if m*c == 0 {
		return types.Error("finaltemp: mass * specific_heat cannot be zero")
	}

	Tf := Ti + Q/(m*c)
	return types.Number(Tf)
}

// FnThermalEquilibrium calculates equilibrium temperature of two objects.
// T_f = (m₁c₁T₁ + m₂c₂T₂) / (m₁c₁ + m₂c₂)
// Args: m1, c1, T1, m2, c2, T2
func FnThermalEquilibrium(args []types.Value) types.Value {
	if len(args) != 6 {
		return types.Error("thermalequilibrium requires 6 arguments: m1, c1, T1, m2, c2, T2")
	}

	m1 := args[0].AsFloat()
	c1 := args[1].AsFloat()
	T1 := args[2].AsFloat()
	m2 := args[3].AsFloat()
	c2 := args[4].AsFloat()
	T2 := args[5].AsFloat()

	denom := m1*c1 + m2*c2
	if denom == 0 {
		return types.Error("thermalequilibrium: total heat capacity cannot be zero")
	}

	Tf := (m1*c1*T1 + m2*c2*T2) / denom
	return types.Number(Tf)
}

// FnLatentHeat calculates heat for phase change.
// Q = mL
// Args: mass (kg), latent heat (J/kg)
func FnLatentHeat(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("latentheat requires 2 arguments: mass, latent_heat")
	}

	m := args[0].AsFloat()
	L := args[1].AsFloat()

	return types.Number(m * L)
}

// ════════════════════════════════════════════════════════════════
// HEAT CONDUCTION
// ════════════════════════════════════════════════════════════════

// FnHeatConduction calculates heat conduction rate (Fourier's law).
// P = kA(ΔT/Δx)
// Args: thermal conductivity (W/(m·K)), area (m²), temp difference (K), thickness (m)
func FnHeatConduction(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("heatconduction requires 4 arguments: k, A, delta_T, thickness")
	}

	k := args[0].AsFloat()
	A := args[1].AsFloat()
	deltaT := args[2].AsFloat()
	dx := args[3].AsFloat()

	if dx == 0 {
		return types.Error("heatconduction: thickness cannot be zero")
	}

	P := k * A * deltaT / dx
	return types.Number(P)
}

// FnThermalResistance calculates thermal resistance.
// R = Δx / (kA)
// Args: thickness (m), thermal conductivity (W/(m·K)), area (m²)
func FnThermalResistance(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("thermalresistance requires 3 arguments: thickness, k, A")
	}

	dx := args[0].AsFloat()
	k := args[1].AsFloat()
	A := args[2].AsFloat()

	if k*A == 0 {
		return types.Error("thermalresistance: k * A cannot be zero")
	}

	R := dx / (k * A)
	return types.Number(R)
}

// ════════════════════════════════════════════════════════════════
// RADIATION
// ════════════════════════════════════════════════════════════════

// FnStefanBoltzmann calculates radiated power (Stefan-Boltzmann law).
// P = εσAT⁴
// Args: emissivity, area (m²), temperature (K)
func FnStefanBoltzmann(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("stefanboltzmann requires 3 arguments: emissivity, A, T")
	}

	epsilon := args[0].AsFloat()
	A := args[1].AsFloat()
	T := args[2].AsFloat()

	sigma := physicsConstants["sigma_sb"].Value
	P := epsilon * sigma * A * math.Pow(T, 4)

	return types.Number(P)
}

// FnBlackBodyPower calculates power radiated by a black body.
// P = σAT⁴
// Args: area (m²), temperature (K)
func FnBlackBodyPower(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("blackbodypower requires 2 arguments: A, T")
	}

	A := args[0].AsFloat()
	T := args[1].AsFloat()

	sigma := physicsConstants["sigma_sb"].Value
	P := sigma * A * math.Pow(T, 4)

	return types.Number(P)
}

// FnNetRadiation calculates net radiation between object and surroundings.
// P_net = εσA(T⁴ - T_surr⁴)
// Args: emissivity, area (m²), object temp (K), surroundings temp (K)
func FnNetRadiation(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("netradiation requires 4 arguments: emissivity, A, T_object, T_surroundings")
	}

	epsilon := args[0].AsFloat()
	A := args[1].AsFloat()
	T := args[2].AsFloat()
	Tsurr := args[3].AsFloat()

	sigma := physicsConstants["sigma_sb"].Value
	P := epsilon * sigma * A * (math.Pow(T, 4) - math.Pow(Tsurr, 4))

	return types.Number(P)
}

// FnWienDisplacement calculates peak wavelength (Wien's law).
// λ_max = b/T where b = 2.898×10⁻³ m·K
// Args: temperature (K)
func FnWienDisplacement(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("wiendisplacement requires 1 argument: T")
	}

	T := args[0].AsFloat()

	if T <= 0 {
		return types.Error("wiendisplacement: temperature must be positive")
	}

	b := physicsConstants["wien"].Value
	lambda := b / T

	return types.Number(lambda)
}

// FnTempFromWien calculates temperature from peak wavelength.
// T = b/λ_max
// Args: wavelength (m)
func FnTempFromWien(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tempfromwien requires 1 argument: wavelength")
	}

	lambda := args[0].AsFloat()

	if lambda <= 0 {
		return types.Error("tempfromwien: wavelength must be positive")
	}

	b := physicsConstants["wien"].Value
	T := b / lambda

	return types.Number(T)
}

// ════════════════════════════════════════════════════════════════
// THERMODYNAMIC PROCESSES
// ════════════════════════════════════════════════════════════════

// FnIsothermalWork calculates work done in isothermal process.
// W = nRT ln(V₂/V₁)
// Args: moles (n), temperature (K), V1 (m³), V2 (m³)
func FnIsothermalWork(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("isothermalwork requires 4 arguments: n, T, V1, V2")
	}

	n := args[0].AsFloat()
	T := args[1].AsFloat()
	V1 := args[2].AsFloat()
	V2 := args[3].AsFloat()

	if V1 <= 0 || V2 <= 0 {
		return types.Error("isothermalwork: volumes must be positive")
	}

	R := physicsConstants["r_gas"].Value
	W := n * R * T * math.Log(V2/V1)

	return types.Number(W)
}

// FnIsobaricWork calculates work done in isobaric (constant pressure) process.
// W = PΔV
// Args: pressure (Pa), V1 (m³), V2 (m³)
func FnIsobaricWork(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("isobaricwork requires 3 arguments: P, V1, V2")
	}

	P := args[0].AsFloat()
	V1 := args[1].AsFloat()
	V2 := args[2].AsFloat()

	W := P * (V2 - V1)
	return types.Number(W)
}

// FnAdiabaticWork calculates work done in adiabatic process.
// W = (P₁V₁ - P₂V₂) / (γ - 1)
// Args: P1 (Pa), V1 (m³), P2 (Pa), V2 (m³), gamma
func FnAdiabaticWork(args []types.Value) types.Value {
	if len(args) != 5 {
		return types.Error("adiabaticwork requires 5 arguments: P1, V1, P2, V2, gamma")
	}

	P1 := args[0].AsFloat()
	V1 := args[1].AsFloat()
	P2 := args[2].AsFloat()
	V2 := args[3].AsFloat()
	gamma := args[4].AsFloat()

	if gamma == 1 {
		return types.Error("adiabaticwork: gamma cannot equal 1")
	}

	W := (P1*V1 - P2*V2) / (gamma - 1)
	return types.Number(W)
}

// FnAdiabaticRelation calculates final state in adiabatic process.
// P₁V₁^γ = P₂V₂^γ or T₁V₁^(γ-1) = T₂V₂^(γ-1)
// Args: P1 or T1, V1, V2, gamma, type ("P" or "T")
func FnAdiabaticFinal(args []types.Value) types.Value {
	if len(args) != 5 {
		return types.Error("adiabaticfinal requires 5 arguments: initial, V1, V2, gamma, type")
	}

	initial := args[0].AsFloat()
	V1 := args[1].AsFloat()
	V2 := args[2].AsFloat()
	gamma := args[3].AsFloat()
	calcType := args[4].AsString()

	if V1 <= 0 || V2 <= 0 {
		return types.Error("adiabaticfinal: volumes must be positive")
	}

	var final float64
	if calcType == "P" || calcType == "p" {
		// P₂ = P₁(V₁/V₂)^γ
		final = initial * math.Pow(V1/V2, gamma)
	} else if calcType == "T" || calcType == "t" {
		// T₂ = T₁(V₁/V₂)^(γ-1)
		final = initial * math.Pow(V1/V2, gamma-1)
	} else {
		return types.Error("adiabaticfinal: type must be 'P' or 'T'")
	}

	return types.Number(final)
}

// ════════════════════════════════════════════════════════════════
// ENTROPY
// ════════════════════════════════════════════════════════════════

// FnEntropy calculates entropy change for reversible heat transfer.
// ΔS = Q/T
// Args: heat (J), temperature (K)
func FnEntropy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("entropy requires 2 arguments: Q, T")
	}

	Q := args[0].AsFloat()
	T := args[1].AsFloat()

	if T == 0 {
		return types.Error("entropy: temperature cannot be zero")
	}

	return types.Number(Q / T)
}

// FnEntropyIdealGas calculates entropy change for ideal gas.
// ΔS = nCᵥln(T₂/T₁) + nRln(V₂/V₁)
// Args: moles (n), Cv (J/(mol·K)), T1 (K), T2 (K), V1 (m³), V2 (m³)
func FnEntropyIdealGas(args []types.Value) types.Value {
	if len(args) != 6 {
		return types.Error("entropyidealgas requires 6 arguments: n, Cv, T1, T2, V1, V2")
	}

	n := args[0].AsFloat()
	Cv := args[1].AsFloat()
	T1 := args[2].AsFloat()
	T2 := args[3].AsFloat()
	V1 := args[4].AsFloat()
	V2 := args[5].AsFloat()

	if T1 <= 0 || T2 <= 0 || V1 <= 0 || V2 <= 0 {
		return types.Error("entropyidealgas: temperatures and volumes must be positive")
	}

	R := physicsConstants["r_gas"].Value
	deltaS := n*Cv*math.Log(T2/T1) + n*R*math.Log(V2/V1)

	return types.Number(deltaS)
}

// FnEntropyMixing calculates entropy of mixing for ideal gases.
// ΔS_mix = -R Σ(nᵢ ln(xᵢ))
// Args: n1, n2, ... (moles of each component)
func FnEntropyMixing(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("entropymixing requires at least 2 arguments: n1, n2, ...")
	}

	// Calculate total moles
	total := 0.0
	for _, arg := range args {
		n := arg.AsFloat()
		if n < 0 {
			return types.Error("entropymixing: moles cannot be negative")
		}
		total += n
	}

	if total == 0 {
		return types.Error("entropymixing: total moles cannot be zero")
	}

	// Calculate entropy of mixing
	R := physicsConstants["r_gas"].Value
	deltaS := 0.0
	for _, arg := range args {
		n := arg.AsFloat()
		if n > 0 {
			x := n / total
			deltaS -= n * math.Log(x)
		}
	}
	deltaS *= R

	return types.Number(deltaS)
}

// ════════════════════════════════════════════════════════════════
// HEAT ENGINES & EFFICIENCY
// ════════════════════════════════════════════════════════════════

// FnCarnotEfficiency calculates Carnot efficiency.
// η = 1 - T_cold/T_hot
// Args: T_hot (K), T_cold (K)
func FnCarnotEfficiency(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("carnotefficiency requires 2 arguments: T_hot, T_cold")
	}

	Th := args[0].AsFloat()
	Tc := args[1].AsFloat()

	if Th <= 0 {
		return types.Error("carnotefficiency: T_hot must be positive")
	}
	if Tc < 0 {
		return types.Error("carnotefficiency: T_cold cannot be negative")
	}
	if Tc >= Th {
		return types.Error("carnotefficiency: T_cold must be less than T_hot")
	}

	eta := 1 - Tc/Th
	return types.Number(eta)
}

// FnHeatEngineEfficiency calculates heat engine efficiency.
// η = W/Q_h = (Q_h - Q_c)/Q_h
// Args: Q_hot (J), Q_cold (J)
func FnHeatEngineEfficiency(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("heatengineeff requires 2 arguments: Q_hot, Q_cold")
	}

	Qh := args[0].AsFloat()
	Qc := args[1].AsFloat()

	if Qh == 0 {
		return types.Error("heatengineeff: Q_hot cannot be zero")
	}

	eta := (Qh - Qc) / Qh
	return types.Number(eta)
}

// FnHeatEngineWork calculates work output of heat engine.
// W = Q_h - Q_c
// Args: Q_hot (J), Q_cold (J)
func FnHeatEngineWork(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("heatengwork requires 2 arguments: Q_hot, Q_cold")
	}

	Qh := args[0].AsFloat()
	Qc := args[1].AsFloat()

	return types.Number(Qh - Qc)
}

// FnCOP calculates coefficient of performance for refrigerator/heat pump.
// COP_refrigerator = Q_c / W = Q_c / (Q_h - Q_c)
// COP_heat_pump = Q_h / W = Q_h / (Q_h - Q_c)
// Args: Q_hot (J), Q_cold (J), type ("refrigerator" or "heatpump")
func FnCOP(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("cop requires 3 arguments: Q_hot, Q_cold, type")
	}

	Qh := args[0].AsFloat()
	Qc := args[1].AsFloat()
	copType := args[2].AsString()

	W := Qh - Qc
	if W == 0 {
		return types.Error("cop: work cannot be zero")
	}

	var cop float64
	if copType == "refrigerator" || copType == "r" {
		cop = Qc / W
	} else if copType == "heatpump" || copType == "hp" {
		cop = Qh / W
	} else {
		return types.Error("cop: type must be 'refrigerator' or 'heatpump'")
	}

	return types.Number(cop)
}

// FnCarnotCOP calculates Carnot COP.
// COP_refrigerator = T_c / (T_h - T_c)
// COP_heat_pump = T_h / (T_h - T_c)
// Args: T_hot (K), T_cold (K), type
func FnCarnotCOP(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("carnotcop requires 3 arguments: T_hot, T_cold, type")
	}

	Th := args[0].AsFloat()
	Tc := args[1].AsFloat()
	copType := args[2].AsString()

	deltaT := Th - Tc
	if deltaT == 0 {
		return types.Error("carnotcop: temperature difference cannot be zero")
	}

	var cop float64
	if copType == "refrigerator" || copType == "r" {
		cop = Tc / deltaT
	} else if copType == "heatpump" || copType == "hp" {
		cop = Th / deltaT
	} else {
		return types.Error("carnotcop: type must be 'refrigerator' or 'heatpump'")
	}

	return types.Number(cop)
}

// ════════════════════════════════════════════════════════════════
// KINETIC THEORY OF GASES
// ════════════════════════════════════════════════════════════════

// FnRMSVelocity calculates root-mean-square velocity of gas molecules.
// v_rms = √(3RT/M) = √(3kT/m)
// Args: temperature (K), molar mass (kg/mol)
func FnRMSVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rmsvelocity requires 2 arguments: T, molar_mass")
	}

	T := args[0].AsFloat()
	M := args[1].AsFloat()

	if T < 0 {
		return types.Error("rmsvelocity: temperature cannot be negative")
	}
	if M <= 0 {
		return types.Error("rmsvelocity: molar mass must be positive")
	}

	R := physicsConstants["r_gas"].Value
	vrms := math.Sqrt(3 * R * T / M)

	return types.Number(vrms)
}

// FnMeanVelocity calculates mean velocity of gas molecules.
// v_mean = √(8RT/(πM))
// Args: temperature (K), molar mass (kg/mol)
func FnMeanVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("meanvelocity requires 2 arguments: T, molar_mass")
	}

	T := args[0].AsFloat()
	M := args[1].AsFloat()

	if T < 0 {
		return types.Error("meanvelocity: temperature cannot be negative")
	}
	if M <= 0 {
		return types.Error("meanvelocity: molar mass must be positive")
	}

	R := physicsConstants["r_gas"].Value
	vmean := math.Sqrt(8 * R * T / (math.Pi * M))

	return types.Number(vmean)
}

// FnMostProbableVelocity calculates most probable velocity of gas molecules.
// v_p = √(2RT/M)
// Args: temperature (K), molar mass (kg/mol)
func FnMostProbableVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("mostprobablevel requires 2 arguments: T, molar_mass")
	}

	T := args[0].AsFloat()
	M := args[1].AsFloat()

	if T < 0 {
		return types.Error("mostprobablevel: temperature cannot be negative")
	}
	if M <= 0 {
		return types.Error("mostprobablevel: molar mass must be positive")
	}

	R := physicsConstants["r_gas"].Value
	vp := math.Sqrt(2 * R * T / M)

	return types.Number(vp)
}

// FnMeanFreePath calculates mean free path.
// λ = 1/(√2 · n · σ) = kT/(√2 · P · σ)
// Args: temperature (K), pressure (Pa), collision diameter (m)
func FnMeanFreePath(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("meanfreepath requires 3 arguments: T, P, collision_diameter")
	}

	T := args[0].AsFloat()
	P := args[1].AsFloat()
	d := args[2].AsFloat()

	if P <= 0 {
		return types.Error("meanfreepath: pressure must be positive")
	}
	if d <= 0 {
		return types.Error("meanfreepath: collision diameter must be positive")
	}

	kB := physicsConstants["kb"].Value
	sigma := math.Pi * d * d // Collision cross-section
	lambda := kB * T / (math.Sqrt2 * P * sigma)

	return types.Number(lambda)
}

// FnCollisionFrequency calculates collision frequency.
// z = v_mean / λ
// Args: temperature (K), pressure (Pa), molar mass (kg/mol), collision diameter (m)
func FnCollisionFrequency(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("collisionfreq requires 4 arguments: T, P, molar_mass, collision_diameter")
	}

	T := args[0].AsFloat()
	P := args[1].AsFloat()
	M := args[2].AsFloat()
	d := args[3].AsFloat()

	if T <= 0 || P <= 0 || M <= 0 || d <= 0 {
		return types.Error("collisionfreq: all arguments must be positive")
	}

	R := physicsConstants["r_gas"].Value
	kB := physicsConstants["kb"].Value

	// Mean velocity
	vmean := math.Sqrt(8 * R * T / (math.Pi * M))

	// Mean free path
	sigma := math.Pi * d * d
	lambda := kB * T / (math.Sqrt2 * P * sigma)

	// Collision frequency
	z := vmean / lambda

	return types.Number(z)
}

// FnAvgKineticEnergy calculates average kinetic energy per molecule.
// KE_avg = (3/2)kT
// Args: temperature (K)
func FnAvgKineticEnergy(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("avgkineticenergy requires 1 argument: T")
	}

	T := args[0].AsFloat()

	if T < 0 {
		return types.Error("avgkineticenergy: temperature cannot be negative")
	}

	kB := physicsConstants["kb"].Value
	KE := 1.5 * kB * T

	return types.Number(KE)
}

// FnInternalEnergy calculates internal energy of ideal gas.
// U = (f/2)nRT where f = degrees of freedom
// Args: moles (n), temperature (K), degrees of freedom (f)
func FnInternalEnergy(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("internalenergy requires 3 arguments: n, T, degrees_of_freedom")
	}

	n := args[0].AsFloat()
	T := args[1].AsFloat()
	f := args[2].AsFloat()

	R := physicsConstants["r_gas"].Value
	U := (f / 2) * n * R * T

	return types.Number(U)
}

// ════════════════════════════════════════════════════════════════
// TEMPERATURE CONVERSIONS
// ════════════════════════════════════════════════════════════════

// FnCelsiusToKelvin converts Celsius to Kelvin.
func FnCelsiusToKelvin(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("celsiustokelvin requires 1 argument: celsius")
	}

	C := args[0].AsFloat()
	return types.Number(C + 273.15)
}

// FnKelvinToCelsius converts Kelvin to Celsius.
func FnKelvinToCelsius(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("kelvintocelsius requires 1 argument: kelvin")
	}

	K := args[0].AsFloat()
	return types.Number(K - 273.15)
}

// FnFahrenheitToCelsius converts Fahrenheit to Celsius.
func FnFahrenheitToCelsius(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("fahrenheittocelsius requires 1 argument: fahrenheit")
	}

	F := args[0].AsFloat()
	return types.Number((F - 32) * 5 / 9)
}

// FnCelsiusToFahrenheit converts Celsius to Fahrenheit.
func FnCelsiusToFahrenheit(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("celsiustofahrenheit requires 1 argument: celsius")
	}

	C := args[0].AsFloat()
	return types.Number(C*9/5 + 32)
}

// FnFahrenheitToKelvin converts Fahrenheit to Kelvin.
func FnFahrenheitToKelvin(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("fahrenheittokelvin requires 1 argument: fahrenheit")
	}

	F := args[0].AsFloat()
	return types.Number((F-32)*5/9 + 273.15)
}

// FnKelvinToFahrenheit converts Kelvin to Fahrenheit.
func FnKelvinToFahrenheit(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("kelvintofahrenheit requires 1 argument: kelvin")
	}

	K := args[0].AsFloat()
	return types.Number((K-273.15)*9/5 + 32)
}
