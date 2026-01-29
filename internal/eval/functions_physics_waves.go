// internal/eval/functions_physics_waves.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// WAVE FUNDAMENTALS
// ════════════════════════════════════════════════════════════════

// FnWavelength calculates wavelength from frequency and velocity.
// λ = v/f
// Args: velocity (m/s), frequency (Hz)
func FnWavelength(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("wavelength requires 2 arguments: velocity, frequency")
	}

	v := args[0].AsFloat()
	f := args[1].AsFloat()

	if f == 0 {
		return types.Error("wavelength: frequency cannot be zero")
	}

	return types.Number(v / f)
}

// FnFrequencyWave calculates frequency from wavelength and velocity.
// f = v/λ
// Args: velocity (m/s), wavelength (m)
func FnFrequencyWave(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("frequencywave requires 2 arguments: velocity, wavelength")
	}

	v := args[0].AsFloat()
	lambda := args[1].AsFloat()

	if lambda == 0 {
		return types.Error("frequencywave: wavelength cannot be zero")
	}

	return types.Number(v / lambda)
}

// FnWaveVelocity calculates wave velocity.
// v = fλ
// Args: frequency (Hz), wavelength (m)
func FnWaveVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("wavevelocity requires 2 arguments: frequency, wavelength")
	}

	f := args[0].AsFloat()
	lambda := args[1].AsFloat()

	return types.Number(f * lambda)
}

// FnWaveNumber calculates wave number.
// k = 2π/λ
// Args: wavelength (m)
func FnWaveNumber(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("wavenumber requires 1 argument: wavelength")
	}

	lambda := args[0].AsFloat()

	if lambda == 0 {
		return types.Error("wavenumber: wavelength cannot be zero")
	}

	return types.Number(2 * math.Pi / lambda)
}

// FnAngularFrequency calculates angular frequency.
// ω = 2πf
// Args: frequency (Hz)
func FnAngularFrequency(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("angularfreq requires 1 argument: frequency")
	}

	f := args[0].AsFloat()
	return types.Number(2 * math.Pi * f)
}

// FnWavePeriod calculates wave period.
// T = 1/f
// Args: frequency (Hz)
func FnWavePeriod(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("waveperiod requires 1 argument: frequency")
	}

	f := args[0].AsFloat()

	if f == 0 {
		return types.Error("waveperiod: frequency cannot be zero")
	}

	return types.Number(1 / f)
}

// ════════════════════════════════════════════════════════════════
// WAVE EQUATION
// ════════════════════════════════════════════════════════════════

// FnWaveDisplacement calculates wave displacement at position and time.
// y = A sin(kx - ωt + φ)
// Args: amplitude, wave number (k), position (x), angular freq (ω), time (t), [phase]
func FnWaveDisplacement(args []types.Value) types.Value {
	if len(args) < 5 || len(args) > 6 {
		return types.Error("wavedisplacement requires 5-6 arguments: A, k, x, omega, t, [phase]")
	}

	A := args[0].AsFloat()
	k := args[1].AsFloat()
	x := args[2].AsFloat()
	omega := args[3].AsFloat()
	t := args[4].AsFloat()
	phi := 0.0

	if len(args) == 6 {
		phi = args[5].AsFloat()
	}

	y := A * math.Sin(k*x-omega*t+phi)
	return types.Number(y)
}

// FnWaveVelocityParticle calculates particle velocity in wave.
// v_y = -Aω cos(kx - ωt + φ)
// Args: amplitude, wave number (k), position (x), angular freq (ω), time (t), [phase]
func FnWaveVelocityParticle(args []types.Value) types.Value {
	if len(args) < 5 || len(args) > 6 {
		return types.Error("wavevelparticle requires 5-6 arguments: A, k, x, omega, t, [phase]")
	}

	A := args[0].AsFloat()
	k := args[1].AsFloat()
	x := args[2].AsFloat()
	omega := args[3].AsFloat()
	t := args[4].AsFloat()
	phi := 0.0

	if len(args) == 6 {
		phi = args[5].AsFloat()
	}

	vy := -A * omega * math.Cos(k*x-omega*t+phi)
	return types.Number(vy)
}

// FnWaveIntensity calculates wave intensity.
// I = (1/2)ρvω²A²
// Args: density (kg/m³), velocity (m/s), angular frequency (rad/s), amplitude (m)
func FnWaveIntensity(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("waveintensity requires 4 arguments: density, velocity, omega, amplitude")
	}

	rho := args[0].AsFloat()
	v := args[1].AsFloat()
	omega := args[2].AsFloat()
	A := args[3].AsFloat()

	I := 0.5 * rho * v * omega * omega * A * A
	return types.Number(I)
}

// FnWaveEnergy calculates wave energy per unit length.
// E/L = (1/2)μω²A² where μ = linear mass density
// Args: linear density (kg/m), angular frequency (rad/s), amplitude (m)
func FnWaveEnergy(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("waveenergy requires 3 arguments: linear_density, omega, amplitude")
	}

	mu := args[0].AsFloat()
	omega := args[1].AsFloat()
	A := args[2].AsFloat()

	E := 0.5 * mu * omega * omega * A * A
	return types.Number(E)
}

// FnWavePower calculates power transmitted by wave.
// P = (1/2)μω²A²v
// Args: linear density (kg/m), angular frequency (rad/s), amplitude (m), velocity (m/s)
func FnWavePower(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("wavepower requires 4 arguments: linear_density, omega, amplitude, velocity")
	}

	mu := args[0].AsFloat()
	omega := args[1].AsFloat()
	A := args[2].AsFloat()
	v := args[3].AsFloat()

	P := 0.5 * mu * omega * omega * A * A * v
	return types.Number(P)
}

// ════════════════════════════════════════════════════════════════
// STRING WAVES
// ════════════════════════════════════════════════════════════════

// FnStringWaveVelocity calculates wave velocity on a string.
// v = √(T/μ)
// Args: tension (N), linear density (kg/m)
func FnStringWaveVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("stringwavevel requires 2 arguments: tension, linear_density")
	}

	T := args[0].AsFloat()
	mu := args[1].AsFloat()

	if mu <= 0 {
		return types.Error("stringwavevel: linear density must be positive")
	}
	if T < 0 {
		return types.Error("stringwavevel: tension cannot be negative")
	}

	return types.Number(math.Sqrt(T / mu))
}

// FnStringFundamental calculates fundamental frequency of a string.
// f₁ = (1/2L)√(T/μ)
// Args: length (m), tension (N), linear density (kg/m)
func FnStringFundamental(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("stringfundamental requires 3 arguments: length, tension, linear_density")
	}

	L := args[0].AsFloat()
	T := args[1].AsFloat()
	mu := args[2].AsFloat()

	if L <= 0 {
		return types.Error("stringfundamental: length must be positive")
	}
	if mu <= 0 {
		return types.Error("stringfundamental: linear density must be positive")
	}

	f := (1 / (2 * L)) * math.Sqrt(T/mu)
	return types.Number(f)
}

// FnStringHarmonic calculates nth harmonic frequency of a string.
// fₙ = n × f₁ = (n/2L)√(T/μ)
// Args: harmonic number (n), length (m), tension (N), linear density (kg/m)
func FnStringHarmonic(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("stringharmonic requires 4 arguments: n, length, tension, linear_density")
	}

	n := args[0].AsFloat()
	L := args[1].AsFloat()
	T := args[2].AsFloat()
	mu := args[3].AsFloat()

	if n < 1 {
		return types.Error("stringharmonic: harmonic number must be >= 1")
	}
	if L <= 0 {
		return types.Error("stringharmonic: length must be positive")
	}
	if mu <= 0 {
		return types.Error("stringharmonic: linear density must be positive")
	}

	f := (n / (2 * L)) * math.Sqrt(T/mu)
	return types.Number(f)
}

// ════════════════════════════════════════════════════════════════
// SOUND WAVES
// ════════════════════════════════════════════════════════════════

// FnSoundVelocity calculates speed of sound in a medium.
// v = √(B/ρ) for fluids, v = √(E/ρ) for solids
// Args: bulk modulus or Young's modulus (Pa), density (kg/m³)
func FnSoundVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("soundvelocity requires 2 arguments: modulus, density")
	}

	B := args[0].AsFloat()
	rho := args[1].AsFloat()

	if rho <= 0 {
		return types.Error("soundvelocity: density must be positive")
	}
	if B < 0 {
		return types.Error("soundvelocity: modulus cannot be negative")
	}

	return types.Number(math.Sqrt(B / rho))
}

// FnSoundVelocityGas calculates speed of sound in ideal gas.
// v = √(γRT/M) = √(γP/ρ)
// Args: gamma (Cp/Cv), temperature (K), molar mass (kg/mol)
func FnSoundVelocityGas(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("soundvelocitygas requires 3 arguments: gamma, T, molar_mass")
	}

	gamma := args[0].AsFloat()
	T := args[1].AsFloat()
	M := args[2].AsFloat()

	if T < 0 {
		return types.Error("soundvelocitygas: temperature cannot be negative")
	}
	if M <= 0 {
		return types.Error("soundvelocitygas: molar mass must be positive")
	}

	R := physicsConstants["r_gas"].Value
	v := math.Sqrt(gamma * R * T / M)

	return types.Number(v)
}

// FnSoundIntensity calculates sound intensity.
// I = P/(4πr²) for spherical wave
// Args: power (W), distance (m)
func FnSoundIntensity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("soundintensity requires 2 arguments: power, distance")
	}

	P := args[0].AsFloat()
	r := args[1].AsFloat()

	if r <= 0 {
		return types.Error("soundintensity: distance must be positive")
	}

	I := P / (4 * math.Pi * r * r)
	return types.Number(I)
}

// FnDecibelLevel calculates sound level in decibels.
// β = 10 log₁₀(I/I₀) where I₀ = 10⁻¹² W/m²
// Args: intensity (W/m²)
func FnDecibelLevel(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("decibellevel requires 1 argument: intensity")
	}

	I := args[0].AsFloat()

	if I <= 0 {
		return types.Error("decibellevel: intensity must be positive")
	}

	I0 := 1e-12 // Reference intensity
	beta := 10 * math.Log10(I/I0)

	return types.Number(beta)
}

// FnIntensityFromDecibel calculates intensity from decibel level.
// I = I₀ × 10^(β/10)
// Args: decibel level (dB)
func FnIntensityFromDecibel(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("intensityfromdb requires 1 argument: decibel_level")
	}

	beta := args[0].AsFloat()

	I0 := 1e-12
	I := I0 * math.Pow(10, beta/10)

	return types.Number(I)
}

// FnDecibelAdd calculates combined decibel level.
// β_total = 10 log₁₀(10^(β₁/10) + 10^(β₂/10) + ...)
// Args: dB1, dB2, ...
func FnDecibelAdd(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("decibleadd requires at least 2 arguments")
	}

	sum := 0.0
	for _, arg := range args {
		beta := arg.AsFloat()
		sum += math.Pow(10, beta/10)
	}

	total := 10 * math.Log10(sum)
	return types.Number(total)
}

// ════════════════════════════════════════════════════════════════
// DOPPLER EFFECT
// ════════════════════════════════════════════════════════════════

// FnDopplerFrequency calculates observed frequency with Doppler effect.
// f' = f × (v ± v_o) / (v ∓ v_s)
// Args: source freq (Hz), wave velocity (m/s), observer velocity (m/s), source velocity (m/s)
// Positive velocities: toward each other
func FnDopplerFrequency(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("dopplerfreq requires 4 arguments: f_source, v_wave, v_observer, v_source")
	}

	f := args[0].AsFloat()
	v := args[1].AsFloat()
	vo := args[2].AsFloat() // Observer velocity (positive = toward source)
	vs := args[3].AsFloat() // Source velocity (positive = toward observer)

	denom := v - vs
	if denom == 0 {
		return types.Error("dopplerfreq: source at wave velocity (sonic boom)")
	}

	fPrime := f * (v + vo) / denom
	return types.Number(fPrime)
}

// FnDopplerShift calculates Doppler shift ratio.
// Δf/f = v_relative / v (for v << c)
// Args: relative velocity (m/s), wave velocity (m/s)
func FnDopplerShift(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("dopplershift requires 2 arguments: v_relative, v_wave")
	}

	vRel := args[0].AsFloat()
	v := args[1].AsFloat()

	if v == 0 {
		return types.Error("dopplershift: wave velocity cannot be zero")
	}

	return types.Number(vRel / v)
}

// FnRelativisticDoppler calculates relativistic Doppler shift for light.
// f' = f × √((1 - β)/(1 + β)) where β = v/c (approaching)
// Args: source frequency (Hz), relative velocity (m/s)
func FnRelativisticDoppler(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("reldoppler requires 2 arguments: f_source, velocity")
	}

	f := args[0].AsFloat()
	v := args[1].AsFloat() // Positive = approaching

	c := physicsConstants["c"].Value
	beta := v / c

	if math.Abs(beta) >= 1 {
		return types.Error("reldoppler: velocity must be less than speed of light")
	}

	// For approaching: blueshift
	fPrime := f * math.Sqrt((1+beta)/(1-beta))
	return types.Number(fPrime)
}

// ════════════════════════════════════════════════════════════════
// BEATS & INTERFERENCE
// ════════════════════════════════════════════════════════════════

// FnBeatFrequency calculates beat frequency.
// f_beat = |f₁ - f₂|
// Args: frequency1 (Hz), frequency2 (Hz)
func FnBeatFrequency(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("beatfreq requires 2 arguments: f1, f2")
	}

	f1 := args[0].AsFloat()
	f2 := args[1].AsFloat()

	return types.Number(math.Abs(f1 - f2))
}

// FnPathDifference calculates path difference for interference.
// Δ = d sin(θ)
// Args: slit separation (m), angle (degrees)
func FnPathDifference(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("pathdiff requires 2 arguments: separation, angle_deg")
	}

	d := args[0].AsFloat()
	theta := args[1].AsFloat() * math.Pi / 180

	return types.Number(d * math.Sin(theta))
}

// FnConstructiveInterference calculates angles for constructive interference.
// d sin(θ) = mλ
// Args: slit separation (m), wavelength (m), order (m)
func FnConstructiveAngle(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("constructiveangle requires 3 arguments: separation, wavelength, order")
	}

	d := args[0].AsFloat()
	lambda := args[1].AsFloat()
	m := args[2].AsFloat()

	if d == 0 {
		return types.Error("constructiveangle: separation cannot be zero")
	}

	sinTheta := m * lambda / d
	if math.Abs(sinTheta) > 1 {
		return types.Error("constructiveangle: no solution exists for this order")
	}

	theta := math.Asin(sinTheta) * 180 / math.Pi
	return types.Number(theta)
}

// FnDestructiveAngle calculates angles for destructive interference.
// d sin(θ) = (m + 1/2)λ
// Args: slit separation (m), wavelength (m), order (m)
func FnDestructiveAngle(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("destructiveangle requires 3 arguments: separation, wavelength, order")
	}

	d := args[0].AsFloat()
	lambda := args[1].AsFloat()
	m := args[2].AsFloat()

	if d == 0 {
		return types.Error("destructiveangle: separation cannot be zero")
	}

	sinTheta := (m + 0.5) * lambda / d
	if math.Abs(sinTheta) > 1 {
		return types.Error("destructiveangle: no solution exists for this order")
	}

	theta := math.Asin(sinTheta) * 180 / math.Pi
	return types.Number(theta)
}

// ════════════════════════════════════════════════════════════════
// STANDING WAVES & RESONANCE
// ════════════════════════════════════════════════════════════════

// FnStandingWaveNodes calculates number of nodes.
// For string fixed at both ends: n + 1 nodes for nth harmonic
// Args: harmonic number
func FnStandingWaveNodes(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("standingwavenodes requires 1 argument: harmonic_number")
	}

	n := int(args[0].AsFloat())
	if n < 1 {
		return types.Error("standingwavenodes: harmonic number must be >= 1")
	}

	return types.Number(float64(n + 1))
}

// FnStandingWaveAntinodes calculates number of antinodes.
// For string fixed at both ends: n antinodes for nth harmonic
// Args: harmonic number
func FnStandingWaveAntinodes(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("standingwaveantinodes requires 1 argument: harmonic_number")
	}

	n := int(args[0].AsFloat())
	if n < 1 {
		return types.Error("standingwaveantinodes: harmonic number must be >= 1")
	}

	return types.Number(float64(n))
}

// FnOpenPipeFundamental calculates fundamental frequency of open pipe.
// f₁ = v/(2L)
// Args: velocity (m/s), length (m)
func FnOpenPipeFundamental(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("openpipefund requires 2 arguments: velocity, length")
	}

	v := args[0].AsFloat()
	L := args[1].AsFloat()

	if L <= 0 {
		return types.Error("openpipefund: length must be positive")
	}

	return types.Number(v / (2 * L))
}

// FnClosedPipeFundamental calculates fundamental frequency of closed pipe.
// f₁ = v/(4L)
// Args: velocity (m/s), length (m)
func FnClosedPipeFundamental(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("closedpipefund requires 2 arguments: velocity, length")
	}

	v := args[0].AsFloat()
	L := args[1].AsFloat()

	if L <= 0 {
		return types.Error("closedpipefund: length must be positive")
	}

	return types.Number(v / (4 * L))
}

// FnOpenPipeHarmonic calculates nth harmonic of open pipe.
// fₙ = nv/(2L) (all harmonics)
// Args: harmonic number, velocity (m/s), length (m)
func FnOpenPipeHarmonic(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("openpipeharmonic requires 3 arguments: n, velocity, length")
	}

	n := args[0].AsFloat()
	v := args[1].AsFloat()
	L := args[2].AsFloat()

	if n < 1 {
		return types.Error("openpipeharmonic: harmonic number must be >= 1")
	}
	if L <= 0 {
		return types.Error("openpipeharmonic: length must be positive")
	}

	return types.Number(n * v / (2 * L))
}

// FnClosedPipeHarmonic calculates nth harmonic of closed pipe.
// fₙ = nv/(4L) (odd harmonics only: n = 1, 3, 5, ...)
// Args: odd harmonic number, velocity (m/s), length (m)
func FnClosedPipeHarmonic(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("closedpipeharmonic requires 3 arguments: n_odd, velocity, length")
	}

	n := int(args[0].AsFloat())
	v := args[1].AsFloat()
	L := args[2].AsFloat()

	if n < 1 || n%2 == 0 {
		return types.Error("closedpipeharmonic: n must be odd (1, 3, 5, ...)")
	}
	if L <= 0 {
		return types.Error("closedpipeharmonic: length must be positive")
	}

	return types.Number(float64(n) * v / (4 * L))
}

// ════════════════════════════════════════════════════════════════
// OPTICS - REFLECTION & REFRACTION
// ════════════════════════════════════════════════════════════════

// FnSnellsLaw calculates refracted angle using Snell's law.
// n₁ sin(θ₁) = n₂ sin(θ₂)
// Args: n1, angle1 (degrees), n2
func FnSnellsLaw(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("snellslaw requires 3 arguments: n1, angle1_deg, n2")
	}

	n1 := args[0].AsFloat()
	theta1 := args[1].AsFloat() * math.Pi / 180
	n2 := args[2].AsFloat()

	if n2 == 0 {
		return types.Error("snellslaw: n2 cannot be zero")
	}

	sinTheta2 := n1 * math.Sin(theta1) / n2

	if math.Abs(sinTheta2) > 1 {
		return types.Error("snellslaw: total internal reflection occurs")
	}

	theta2 := math.Asin(sinTheta2) * 180 / math.Pi
	return types.Number(theta2)
}

// FnCriticalAngle calculates critical angle for total internal reflection.
// θc = arcsin(n₂/n₁)
// Args: n1 (higher index), n2 (lower index)
func FnCriticalAngle(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("criticalangle requires 2 arguments: n1, n2")
	}

	n1 := args[0].AsFloat()
	n2 := args[1].AsFloat()

	if n1 <= 0 {
		return types.Error("criticalangle: n1 must be positive")
	}
	if n2 >= n1 {
		return types.Error("criticalangle: n1 must be greater than n2 for TIR")
	}

	thetaC := math.Asin(n2/n1) * 180 / math.Pi
	return types.Number(thetaC)
}

// FnBrewsterAngle calculates Brewster's angle.
// θB = arctan(n₂/n₁)
// Args: n1, n2
func FnBrewsterAngle(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("brewsterangle requires 2 arguments: n1, n2")
	}

	n1 := args[0].AsFloat()
	n2 := args[1].AsFloat()

	if n1 == 0 {
		return types.Error("brewsterangle: n1 cannot be zero")
	}

	thetaB := math.Atan(n2/n1) * 180 / math.Pi
	return types.Number(thetaB)
}

// FnRefractiveIndex calculates refractive index from velocities.
// n = c/v
// Args: velocity in medium (m/s)
func FnRefractiveIndex(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("refractiveindex requires 1 argument: velocity_in_medium")
	}

	v := args[0].AsFloat()

	if v <= 0 {
		return types.Error("refractiveindex: velocity must be positive")
	}

	c := physicsConstants["c"].Value
	return types.Number(c / v)
}

// ════════════════════════════════════════════════════════════════
// OPTICS - LENSES & MIRRORS
// ════════════════════════════════════════════════════════════════

// FnThinLensEquation calculates image distance using thin lens equation.
// 1/f = 1/do + 1/di
// Args: focal length (m), object distance (m)
func FnThinLensImage(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("thinlensimage requires 2 arguments: focal_length, object_distance")
	}

	f := args[0].AsFloat()
	do := args[1].AsFloat()

	if f == 0 {
		return types.Error("thinlensimage: focal length cannot be zero")
	}

	denom := 1/f - 1/do
	if denom == 0 {
		return types.Error("thinlensimage: image at infinity")
	}

	di := 1 / denom
	return types.Number(di)
}

// FnThinLensFocal calculates focal length from object and image distances.
// 1/f = 1/do + 1/di
// Args: object distance (m), image distance (m)
func FnThinLensFocal(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("thinlensfocal requires 2 arguments: object_distance, image_distance")
	}

	do := args[0].AsFloat()
	di := args[1].AsFloat()

	if do == 0 || di == 0 {
		return types.Error("thinlensfocal: distances cannot be zero")
	}

	f := 1 / (1/do + 1/di)
	return types.Number(f)
}

// FnMagnification calculates magnification.
// M = -di/do = hi/ho
// Args: object distance (m), image distance (m)
func FnMagnification(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("magnification requires 2 arguments: object_distance, image_distance")
	}

	do := args[0].AsFloat()
	di := args[1].AsFloat()

	if do == 0 {
		return types.Error("magnification: object distance cannot be zero")
	}

	M := -di / do
	return types.Number(M)
}

// FnLensMakerEquation calculates focal length using lensmaker's equation.
// 1/f = (n-1)(1/R₁ - 1/R₂)
// Args: refractive index, R1 (m), R2 (m)
func FnLensMaker(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("lensmaker requires 3 arguments: n, R1, R2")
	}

	n := args[0].AsFloat()
	R1 := args[1].AsFloat()
	R2 := args[2].AsFloat()

	// Handle infinite radii (flat surface)
	term1 := 0.0
	term2 := 0.0

	if R1 != 0 {
		term1 = 1 / R1
	}
	if R2 != 0 {
		term2 = 1 / R2
	}

	invF := (n - 1) * (term1 - term2)
	if invF == 0 {
		return types.Error("lensmaker: focal length is infinite")
	}

	return types.Number(1 / invF)
}

// FnMirrorEquation calculates image distance for spherical mirror.
// 1/f = 1/do + 1/di, where f = R/2
// Args: radius of curvature (m), object distance (m)
func FnMirrorImage(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("mirrorimage requires 2 arguments: radius_of_curvature, object_distance")
	}

	R := args[0].AsFloat()
	do := args[1].AsFloat()

	f := R / 2

	if f == 0 {
		return types.Error("mirrorimage: flat mirror (R = 0)")
	}

	denom := 1/f - 1/do
	if denom == 0 {
		return types.Error("mirrorimage: image at infinity")
	}

	di := 1 / denom
	return types.Number(di)
}

// FnDioptricPower calculates power of lens in diopters.
// P = 1/f (where f in meters)
// Args: focal length (m)
func FnDioptricPower(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("dioptricpower requires 1 argument: focal_length")
	}

	f := args[0].AsFloat()

	if f == 0 {
		return types.Error("dioptricpower: focal length cannot be zero")
	}

	return types.Number(1 / f)
}

// ════════════════════════════════════════════════════════════════
// OPTICS - DIFFRACTION
// ════════════════════════════════════════════════════════════════

// FnSingleSlitMinima calculates angle for single slit diffraction minima.
// a sin(θ) = mλ
// Args: slit width (m), wavelength (m), order (m)
func FnSingleSlitMinima(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("singleslitmin requires 3 arguments: slit_width, wavelength, order")
	}

	a := args[0].AsFloat()
	lambda := args[1].AsFloat()
	m := args[2].AsFloat()

	if a == 0 {
		return types.Error("singleslitmin: slit width cannot be zero")
	}
	if m == 0 {
		return types.Error("singleslitmin: order cannot be zero (central maximum)")
	}

	sinTheta := m * lambda / a
	if math.Abs(sinTheta) > 1 {
		return types.Error("singleslitmin: no solution for this order")
	}

	theta := math.Asin(sinTheta) * 180 / math.Pi
	return types.Number(theta)
}

// FnDiffractionGrating calculates angle for diffraction grating maxima.
// d sin(θ) = mλ
// Args: grating spacing (m), wavelength (m), order (m)
func FnDiffractionGrating(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("diffractiongrating requires 3 arguments: spacing, wavelength, order")
	}

	d := args[0].AsFloat()
	lambda := args[1].AsFloat()
	m := args[2].AsFloat()

	if d == 0 {
		return types.Error("diffractiongrating: spacing cannot be zero")
	}

	sinTheta := m * lambda / d
	if math.Abs(sinTheta) > 1 {
		return types.Error("diffractiongrating: no solution for this order")
	}

	theta := math.Asin(sinTheta) * 180 / math.Pi
	return types.Number(theta)
}

// FnRayleighCriterion calculates minimum resolvable angle (Rayleigh criterion).
// θ = 1.22 λ/D
// Args: wavelength (m), aperture diameter (m)
func FnRayleighCriterion(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rayleighcriterion requires 2 arguments: wavelength, diameter")
	}

	lambda := args[0].AsFloat()
	D := args[1].AsFloat()

	if D <= 0 {
		return types.Error("rayleighcriterion: diameter must be positive")
	}

	theta := 1.22 * lambda / D
	return types.Number(theta) // In radians
}

// FnAiryDiskRadius calculates radius of Airy disk.
// r = 1.22 λf/D
// Args: wavelength (m), focal length (m), aperture diameter (m)
func FnAiryDiskRadius(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("airydiskradius requires 3 arguments: wavelength, focal_length, diameter")
	}

	lambda := args[0].AsFloat()
	f := args[1].AsFloat()
	D := args[2].AsFloat()

	if D <= 0 {
		return types.Error("airydiskradius: diameter must be positive")
	}

	r := 1.22 * lambda * f / D
	return types.Number(r)
}

// ════════════════════════════════════════════════════════════════
// ELECTROMAGNETIC SPECTRUM
// ════════════════════════════════════════════════════════════════

// FnLightWavelength calculates wavelength of light from frequency.
// λ = c/f
// Args: frequency (Hz)
func FnLightWavelength(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("lightwavelength requires 1 argument: frequency")
	}

	f := args[0].AsFloat()

	if f == 0 {
		return types.Error("lightwavelength: frequency cannot be zero")
	}

	c := physicsConstants["c"].Value
	return types.Number(c / f)
}

// FnLightFrequency calculates frequency of light from wavelength.
// f = c/λ
// Args: wavelength (m)
func FnLightFrequency(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("lightfrequency requires 1 argument: wavelength")
	}

	lambda := args[0].AsFloat()

	if lambda == 0 {
		return types.Error("lightfrequency: wavelength cannot be zero")
	}

	c := physicsConstants["c"].Value
	return types.Number(c / lambda)
}

// FnLightEnergy calculates energy of a photon.
// E = hf = hc/λ
// Args: wavelength (m) or frequency (Hz), type ("wavelength" or "frequency")
func FnLightEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("lightenergy requires 2 arguments: value, type")
	}

	val := args[0].AsFloat()
	calcType := args[1].AsString()

	h := physicsConstants["h"].Value
	c := physicsConstants["c"].Value

	var E float64
	if calcType == "wavelength" || calcType == "w" || calcType == "λ" {
		if val <= 0 {
			return types.Error("lightenergy: wavelength must be positive")
		}
		E = h * c / val
	} else if calcType == "frequency" || calcType == "f" {
		if val <= 0 {
			return types.Error("lightenergy: frequency must be positive")
		}
		E = h * val
	} else {
		return types.Error("lightenergy: type must be 'wavelength' or 'frequency'")
	}

	return types.Number(E)
}
