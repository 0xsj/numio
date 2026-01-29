// internal/eval/functions_physics_mechanics.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// FORCE & MOTION
// ════════════════════════════════════════════════════════════════

// FnForce calculates force using Newton's second law.
// F = ma
// Args: mass (kg), acceleration (m/s²)
func FnForce(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("force requires 2 arguments: mass, acceleration")
	}

	m := args[0].AsFloat()
	a := args[1].AsFloat()

	return types.Number(m * a)
}

// FnMass calculates mass from force and acceleration.
// m = F/a
// Args: force (N), acceleration (m/s²)
func FnMass(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("mass requires 2 arguments: force, acceleration")
	}

	F := args[0].AsFloat()
	a := args[1].AsFloat()

	if a == 0 {
		return types.Error("mass: acceleration cannot be zero")
	}

	return types.Number(F / a)
}

// FnAcceleration calculates acceleration from force and mass.
// a = F/m
// Args: force (N), mass (kg)
func FnAcceleration(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("acceleration requires 2 arguments: force, mass")
	}

	F := args[0].AsFloat()
	m := args[1].AsFloat()

	if m == 0 {
		return types.Error("acceleration: mass cannot be zero")
	}

	return types.Number(F / m)
}

// FnWeight calculates weight (gravitational force).
// W = mg
// Args: mass (kg), [g] (defaults to Earth gravity 9.80665 m/s²)
func FnWeight(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("weight requires 1-2 arguments: mass, [g]")
	}

	m := args[0].AsFloat()
	g := physicsConstants["g_earth"].Value

	if len(args) == 2 {
		g = args[1].AsFloat()
	}

	return types.Number(m * g)
}

// ════════════════════════════════════════════════════════════════
// MOMENTUM & IMPULSE
// ════════════════════════════════════════════════════════════════

// FnMomentum calculates linear momentum.
// p = mv
// Args: mass (kg), velocity (m/s)
func FnMomentum(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("momentum requires 2 arguments: mass, velocity")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()

	return types.Number(m * v)
}

// FnImpulse calculates impulse.
// J = FΔt or J = Δp
// Args: force (N), time (s) OR mass (kg), Δv (m/s)
func FnImpulse(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("impulse requires 2 arguments: force, time OR mass, delta_v")
	}

	a := args[0].AsFloat()
	b := args[1].AsFloat()

	return types.Number(a * b)
}

// FnMomentumConservation calculates final velocity after collision.
// m1*v1 + m2*v2 = (m1+m2)*vf (perfectly inelastic)
// Args: m1, v1, m2, v2
func FnMomentumConservation(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("momentumcons requires 4 arguments: m1, v1, m2, v2")
	}

	m1 := args[0].AsFloat()
	v1 := args[1].AsFloat()
	m2 := args[2].AsFloat()
	v2 := args[3].AsFloat()

	if m1+m2 == 0 {
		return types.Error("momentumcons: total mass cannot be zero")
	}

	vf := (m1*v1 + m2*v2) / (m1 + m2)
	return types.Number(vf)
}

// ════════════════════════════════════════════════════════════════
// ENERGY
// ════════════════════════════════════════════════════════════════

// FnKineticEnergy calculates kinetic energy.
// KE = ½mv²
// Args: mass (kg), velocity (m/s)
func FnKineticEnergy(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("kineticenergy requires 2 arguments: mass, velocity")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()

	return types.Number(0.5 * m * v * v)
}

// FnPotentialEnergy calculates gravitational potential energy.
// PE = mgh
// Args: mass (kg), height (m), [g] (defaults to Earth gravity)
func FnPotentialEnergy(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("potentialenergy requires 2-3 arguments: mass, height, [g]")
	}

	m := args[0].AsFloat()
	h := args[1].AsFloat()
	g := physicsConstants["g_earth"].Value

	if len(args) == 3 {
		g = args[2].AsFloat()
	}

	return types.Number(m * g * h)
}

// FnElasticPE calculates elastic potential energy.
// PE = ½kx²
// Args: spring constant (N/m), displacement (m)
func FnElasticPE(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("elasticpe requires 2 arguments: spring_constant, displacement")
	}

	k := args[0].AsFloat()
	x := args[1].AsFloat()

	return types.Number(0.5 * k * x * x)
}

// FnVelocityFromKE calculates velocity from kinetic energy.
// v = √(2KE/m)
// Args: kinetic energy (J), mass (kg)
func FnVelocityFromKE(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("velocityfromke requires 2 arguments: kinetic_energy, mass")
	}

	KE := args[0].AsFloat()
	m := args[1].AsFloat()

	if m == 0 {
		return types.Error("velocityfromke: mass cannot be zero")
	}
	if KE < 0 {
		return types.Error("velocityfromke: kinetic energy cannot be negative")
	}

	return types.Number(math.Sqrt(2 * KE / m))
}

// ════════════════════════════════════════════════════════════════
// WORK & POWER
// ════════════════════════════════════════════════════════════════

// FnWork calculates work done by a force.
// W = Fd·cos(θ)
// Args: force (N), distance (m), [angle in degrees] (defaults to 0)
func FnWork(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("work requires 2-3 arguments: force, distance, [angle_deg]")
	}

	F := args[0].AsFloat()
	d := args[1].AsFloat()
	theta := 0.0

	if len(args) == 3 {
		theta = args[2].AsFloat() * math.Pi / 180 // Convert to radians
	}

	return types.Number(F * d * math.Cos(theta))
}

// FnPowerMech calculates mechanical power.
// P = W/t or P = Fv
// Args: work (J), time (s) OR force (N), velocity (m/s)
func FnPowerMech(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("powermech requires 2 arguments: work, time OR force, velocity")
	}

	a := args[0].AsFloat()
	b := args[1].AsFloat()

	if b == 0 {
		return types.Error("powermech: second argument cannot be zero")
	}

	return types.Number(a / b)
}

// FnPowerFromFV calculates power from force and velocity.
// P = Fv
// Args: force (N), velocity (m/s)
func FnPowerFromFV(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("powerfv requires 2 arguments: force, velocity")
	}

	F := args[0].AsFloat()
	v := args[1].AsFloat()

	return types.Number(F * v)
}

// FnEfficiency calculates mechanical efficiency.
// η = (useful output / total input) × 100%
// Args: useful output, total input
func FnEfficiency(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("efficiency requires 2 arguments: output, input")
	}

	output := args[0].AsFloat()
	input := args[1].AsFloat()

	if input == 0 {
		return types.Error("efficiency: input cannot be zero")
	}

	return types.Number((output / input) * 100)
}

// ════════════════════════════════════════════════════════════════
// KINEMATICS
// ════════════════════════════════════════════════════════════════

// FnDisplacement calculates displacement using kinematic equations.
// s = v₀t + ½at²
// Args: initial velocity (m/s), time (s), acceleration (m/s²)
func FnDisplacement(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("displacement requires 3 arguments: v0, time, acceleration")
	}

	v0 := args[0].AsFloat()
	t := args[1].AsFloat()
	a := args[2].AsFloat()

	return types.Number(v0*t + 0.5*a*t*t)
}

// FnFinalVelocity calculates final velocity.
// v = v₀ + at
// Args: initial velocity (m/s), acceleration (m/s²), time (s)
func FnFinalVelocity(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("finalvelocity requires 3 arguments: v0, acceleration, time")
	}

	v0 := args[0].AsFloat()
	a := args[1].AsFloat()
	t := args[2].AsFloat()

	return types.Number(v0 + a*t)
}

// FnFinalVelocitySq calculates final velocity squared.
// v² = v₀² + 2as
// Args: initial velocity (m/s), acceleration (m/s²), displacement (m)
func FnFinalVelocitySq(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("finalvelocitysq requires 3 arguments: v0, acceleration, displacement")
	}

	v0 := args[0].AsFloat()
	a := args[1].AsFloat()
	s := args[2].AsFloat()

	vSq := v0*v0 + 2*a*s
	if vSq < 0 {
		return types.Error("finalvelocitysq: result would be imaginary")
	}

	return types.Number(math.Sqrt(vSq))
}

// FnTimeFromKinematics calculates time from initial velocity, final velocity, and acceleration.
// t = (v - v₀) / a
// Args: v0, v, a
func FnTimeFromKinematics(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("timefromkin requires 3 arguments: v0, v, acceleration")
	}

	v0 := args[0].AsFloat()
	v := args[1].AsFloat()
	a := args[2].AsFloat()

	if a == 0 {
		return types.Error("timefromkin: acceleration cannot be zero")
	}

	return types.Number((v - v0) / a)
}

// FnFreeFall calculates free fall parameters.
// h = ½gt² (dropping from rest)
// Args: time (s), [g] (defaults to Earth gravity)
func FnFreeFall(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("freefall requires 1-2 arguments: time, [g]")
	}

	t := args[0].AsFloat()
	g := physicsConstants["g_earth"].Value

	if len(args) == 2 {
		g = args[1].AsFloat()
	}

	return types.Number(0.5 * g * t * t)
}

// FnFreeFallTime calculates time to fall a given height.
// t = √(2h/g)
// Args: height (m), [g] (defaults to Earth gravity)
func FnFreeFallTime(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("freefalltime requires 1-2 arguments: height, [g]")
	}

	h := args[0].AsFloat()
	g := physicsConstants["g_earth"].Value

	if len(args) == 2 {
		g = args[1].AsFloat()
	}

	if h < 0 {
		return types.Error("freefalltime: height cannot be negative")
	}
	if g <= 0 {
		return types.Error("freefalltime: g must be positive")
	}

	return types.Number(math.Sqrt(2 * h / g))
}

// FnFreeFallVelocity calculates velocity after falling a given height.
// v = √(2gh)
// Args: height (m), [g] (defaults to Earth gravity)
func FnFreeFallVelocity(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("freefallvel requires 1-2 arguments: height, [g]")
	}

	h := args[0].AsFloat()
	g := physicsConstants["g_earth"].Value

	if len(args) == 2 {
		g = args[1].AsFloat()
	}

	if h < 0 {
		return types.Error("freefallvel: height cannot be negative")
	}

	return types.Number(math.Sqrt(2 * g * h))
}

// ════════════════════════════════════════════════════════════════
// PROJECTILE MOTION
// ════════════════════════════════════════════════════════════════

// FnProjectileRange calculates horizontal range of a projectile.
// R = (v₀² sin(2θ)) / g
// Args: initial velocity (m/s), angle (degrees), [g]
func FnProjectileRange(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("projectilerange requires 2-3 arguments: v0, angle_deg, [g]")
	}

	v0 := args[0].AsFloat()
	theta := args[1].AsFloat() * math.Pi / 180 // Convert to radians
	g := physicsConstants["g_earth"].Value

	if len(args) == 3 {
		g = args[2].AsFloat()
	}

	if g <= 0 {
		return types.Error("projectilerange: g must be positive")
	}

	R := (v0 * v0 * math.Sin(2*theta)) / g
	return types.Number(R)
}

// FnProjectileHeight calculates maximum height of a projectile.
// H = (v₀² sin²(θ)) / (2g)
// Args: initial velocity (m/s), angle (degrees), [g]
func FnProjectileHeight(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("projectileheight requires 2-3 arguments: v0, angle_deg, [g]")
	}

	v0 := args[0].AsFloat()
	theta := args[1].AsFloat() * math.Pi / 180
	g := physicsConstants["g_earth"].Value

	if len(args) == 3 {
		g = args[2].AsFloat()
	}

	if g <= 0 {
		return types.Error("projectileheight: g must be positive")
	}

	sinTheta := math.Sin(theta)
	H := (v0 * v0 * sinTheta * sinTheta) / (2 * g)
	return types.Number(H)
}

// FnProjectileTime calculates total flight time of a projectile.
// T = (2v₀ sin(θ)) / g
// Args: initial velocity (m/s), angle (degrees), [g]
func FnProjectileTime(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("projectiletime requires 2-3 arguments: v0, angle_deg, [g]")
	}

	v0 := args[0].AsFloat()
	theta := args[1].AsFloat() * math.Pi / 180
	g := physicsConstants["g_earth"].Value

	if len(args) == 3 {
		g = args[2].AsFloat()
	}

	if g <= 0 {
		return types.Error("projectiletime: g must be positive")
	}

	T := (2 * v0 * math.Sin(theta)) / g
	return types.Number(T)
}

// ════════════════════════════════════════════════════════════════
// CIRCULAR MOTION
// ════════════════════════════════════════════════════════════════

// FnCentripetal calculates centripetal force.
// F = mv²/r
// Args: mass (kg), velocity (m/s), radius (m)
func FnCentripetal(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("centripetal requires 3 arguments: mass, velocity, radius")
	}

	m := args[0].AsFloat()
	v := args[1].AsFloat()
	r := args[2].AsFloat()

	if r == 0 {
		return types.Error("centripetal: radius cannot be zero")
	}

	return types.Number(m * v * v / r)
}

// FnCentripetalAccel calculates centripetal acceleration.
// a = v²/r
// Args: velocity (m/s), radius (m)
func FnCentripetalAccel(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("centripetalaccel requires 2 arguments: velocity, radius")
	}

	v := args[0].AsFloat()
	r := args[1].AsFloat()

	if r == 0 {
		return types.Error("centripetalaccel: radius cannot be zero")
	}

	return types.Number(v * v / r)
}

// FnAngularVelocity calculates angular velocity.
// ω = v/r or ω = 2πf
// Args: linear velocity (m/s), radius (m) OR frequency (Hz)
func FnAngularVelocity(args []types.Value) types.Value {
	if len(args) == 1 {
		// From frequency
		f := args[0].AsFloat()
		return types.Number(2 * math.Pi * f)
	}

	if len(args) == 2 {
		// From linear velocity and radius
		v := args[0].AsFloat()
		r := args[1].AsFloat()

		if r == 0 {
			return types.Error("angularvel: radius cannot be zero")
		}

		return types.Number(v / r)
	}

	return types.Error("angularvel requires 1-2 arguments: frequency OR velocity, radius")
}

// FnLinearVelocity calculates linear velocity from angular.
// v = ωr
// Args: angular velocity (rad/s), radius (m)
func FnLinearVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("linearvel requires 2 arguments: angular_velocity, radius")
	}

	omega := args[0].AsFloat()
	r := args[1].AsFloat()

	return types.Number(omega * r)
}

// FnPeriod calculates period from frequency or angular velocity.
// T = 1/f or T = 2π/ω
// Args: frequency (Hz) or angular velocity (rad/s), [type: "f" or "omega"]
func FnPeriod(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("period requires 1-2 arguments: frequency, [type]")
	}

	val := args[0].AsFloat()
	if val == 0 {
		return types.Error("period: input cannot be zero")
	}

	// Default: assume frequency
	if len(args) == 1 {
		return types.Number(1 / val)
	}

	// Check type
	typeStr := args[1].AsString()
	if typeStr == "omega" || typeStr == "w" {
		return types.Number(2 * math.Pi / val)
	}

	return types.Number(1 / val)
}

// FnFrequency calculates frequency from period.
// f = 1/T
// Args: period (s)
func FnFrequencyMech(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("frequencymech requires 1 argument: period")
	}

	T := args[0].AsFloat()
	if T == 0 {
		return types.Error("frequencymech: period cannot be zero")
	}

	return types.Number(1 / T)
}

// ════════════════════════════════════════════════════════════════
// ROTATIONAL MECHANICS
// ════════════════════════════════════════════════════════════════

// FnTorque calculates torque.
// τ = rF·sin(θ)
// Args: radius/lever arm (m), force (N), [angle in degrees] (defaults to 90)
func FnTorque(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("torque requires 2-3 arguments: radius, force, [angle_deg]")
	}

	r := args[0].AsFloat()
	F := args[1].AsFloat()
	theta := 90.0 // Default perpendicular

	if len(args) == 3 {
		theta = args[2].AsFloat()
	}

	thetaRad := theta * math.Pi / 180
	return types.Number(r * F * math.Sin(thetaRad))
}

// FnMomentOfInertia calculates moment of inertia for a point mass.
// I = mr²
// Args: mass (kg), radius (m)
func FnMomentOfInertia(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("momentofinertia requires 2 arguments: mass, radius")
	}

	m := args[0].AsFloat()
	r := args[1].AsFloat()

	return types.Number(m * r * r)
}

// FnMomentOfInertiaSphere calculates moment of inertia for a solid sphere.
// I = (2/5)mr²
// Args: mass (kg), radius (m)
func FnMomentOfInertiaSphere(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("moisphere requires 2 arguments: mass, radius")
	}

	m := args[0].AsFloat()
	r := args[1].AsFloat()

	return types.Number((2.0 / 5.0) * m * r * r)
}

// FnMomentOfInertiaCylinder calculates moment of inertia for a solid cylinder.
// I = (1/2)mr²
// Args: mass (kg), radius (m)
func FnMomentOfInertiaCylinder(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("moicylinder requires 2 arguments: mass, radius")
	}

	m := args[0].AsFloat()
	r := args[1].AsFloat()

	return types.Number(0.5 * m * r * r)
}

// FnMomentOfInertiaRod calculates moment of inertia for a rod (about center).
// I = (1/12)mL²
// Args: mass (kg), length (m)
func FnMomentOfInertiaRod(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("moirod requires 2 arguments: mass, length")
	}

	m := args[0].AsFloat()
	L := args[1].AsFloat()

	return types.Number((1.0 / 12.0) * m * L * L)
}

// FnAngularMomentum calculates angular momentum.
// L = Iω
// Args: moment of inertia (kg·m²), angular velocity (rad/s)
func FnAngularMomentum(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("angularmomentum requires 2 arguments: moment_of_inertia, angular_velocity")
	}

	I := args[0].AsFloat()
	omega := args[1].AsFloat()

	return types.Number(I * omega)
}

// FnRotationalKE calculates rotational kinetic energy.
// KE = ½Iω²
// Args: moment of inertia (kg·m²), angular velocity (rad/s)
func FnRotationalKE(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rotationalke requires 2 arguments: moment_of_inertia, angular_velocity")
	}

	I := args[0].AsFloat()
	omega := args[1].AsFloat()

	return types.Number(0.5 * I * omega * omega)
}

// FnAngularAccel calculates angular acceleration.
// α = τ/I
// Args: torque (N·m), moment of inertia (kg·m²)
func FnAngularAccel(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("angularaccel requires 2 arguments: torque, moment_of_inertia")
	}

	tau := args[0].AsFloat()
	I := args[1].AsFloat()

	if I == 0 {
		return types.Error("angularaccel: moment of inertia cannot be zero")
	}

	return types.Number(tau / I)
}

// ════════════════════════════════════════════════════════════════
// SIMPLE HARMONIC MOTION
// ════════════════════════════════════════════════════════════════

// FnSHMPeriodSpring calculates period of a mass-spring system.
// T = 2π√(m/k)
// Args: mass (kg), spring constant (N/m)
func FnSHMPeriodSpring(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("shmperiodspring requires 2 arguments: mass, spring_constant")
	}

	m := args[0].AsFloat()
	k := args[1].AsFloat()

	if m <= 0 {
		return types.Error("shmperiodspring: mass must be positive")
	}
	if k <= 0 {
		return types.Error("shmperiodspring: spring constant must be positive")
	}

	return types.Number(2 * math.Pi * math.Sqrt(m/k))
}

// FnSHMPeriodPendulum calculates period of a simple pendulum.
// T = 2π√(L/g)
// Args: length (m), [g] (defaults to Earth gravity)
func FnSHMPeriodPendulum(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("shmperiodpendulum requires 1-2 arguments: length, [g]")
	}

	L := args[0].AsFloat()
	g := physicsConstants["g_earth"].Value

	if len(args) == 2 {
		g = args[1].AsFloat()
	}

	if L <= 0 {
		return types.Error("shmperiodpendulum: length must be positive")
	}
	if g <= 0 {
		return types.Error("shmperiodpendulum: g must be positive")
	}

	return types.Number(2 * math.Pi * math.Sqrt(L/g))
}

// FnSHMFrequency calculates frequency of SHM.
// f = 1/(2π)√(k/m)
// Args: spring constant (N/m), mass (kg)
func FnSHMFrequency(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("shmfrequency requires 2 arguments: spring_constant, mass")
	}

	k := args[0].AsFloat()
	m := args[1].AsFloat()

	if k <= 0 {
		return types.Error("shmfrequency: spring constant must be positive")
	}
	if m <= 0 {
		return types.Error("shmfrequency: mass must be positive")
	}

	return types.Number(math.Sqrt(k/m) / (2 * math.Pi))
}

// FnSHMDisplacement calculates displacement in SHM.
// x = A·cos(ωt + φ)
// Args: amplitude (m), angular frequency (rad/s), time (s), [phase] (defaults to 0)
func FnSHMDisplacement(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("shmdisplacement requires 3-4 arguments: amplitude, omega, time, [phase]")
	}

	A := args[0].AsFloat()
	omega := args[1].AsFloat()
	t := args[2].AsFloat()
	phi := 0.0

	if len(args) == 4 {
		phi = args[3].AsFloat()
	}

	return types.Number(A * math.Cos(omega*t+phi))
}

// FnSHMVelocity calculates velocity in SHM.
// v = -Aω·sin(ωt + φ)
// Args: amplitude (m), angular frequency (rad/s), time (s), [phase] (defaults to 0)
func FnSHMVelocity(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("shmvelocity requires 3-4 arguments: amplitude, omega, time, [phase]")
	}

	A := args[0].AsFloat()
	omega := args[1].AsFloat()
	t := args[2].AsFloat()
	phi := 0.0

	if len(args) == 4 {
		phi = args[3].AsFloat()
	}

	return types.Number(-A * omega * math.Sin(omega*t+phi))
}

// FnSHMMaxVelocity calculates maximum velocity in SHM.
// v_max = Aω
// Args: amplitude (m), angular frequency (rad/s)
func FnSHMMaxVelocity(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("shmmaxvel requires 2 arguments: amplitude, omega")
	}

	A := args[0].AsFloat()
	omega := args[1].AsFloat()

	return types.Number(math.Abs(A * omega))
}

// ════════════════════════════════════════════════════════════════
// FRICTION & DRAG
// ════════════════════════════════════════════════════════════════

// FnFriction calculates friction force.
// f = μN
// Args: coefficient of friction, normal force (N)
func FnFriction(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("friction requires 2 arguments: mu, normal_force")
	}

	mu := args[0].AsFloat()
	N := args[1].AsFloat()

	return types.Number(mu * N)
}

// FnDragForce calculates drag force.
// F_d = ½ρv²C_dA
// Args: density (kg/m³), velocity (m/s), drag coefficient, area (m²)
func FnDragForce(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("dragforce requires 4 arguments: density, velocity, drag_coeff, area")
	}

	rho := args[0].AsFloat()
	v := args[1].AsFloat()
	Cd := args[2].AsFloat()
	A := args[3].AsFloat()

	return types.Number(0.5 * rho * v * v * Cd * A)
}

// FnTerminalVelocity calculates terminal velocity.
// v_t = √(2mg / (ρC_dA))
// Args: mass (kg), drag coefficient, area (m²), [density] (defaults to air at sea level)
func FnTerminalVelocity(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("terminalvel requires 3-4 arguments: mass, drag_coeff, area, [density]")
	}

	m := args[0].AsFloat()
	Cd := args[1].AsFloat()
	A := args[2].AsFloat()
	rho := 1.225 // Air density at sea level (kg/m³)

	if len(args) == 4 {
		rho = args[3].AsFloat()
	}

	g := physicsConstants["g_earth"].Value

	denominator := rho * Cd * A
	if denominator == 0 {
		return types.Error("terminalvel: denominator cannot be zero")
	}

	return types.Number(math.Sqrt(2 * m * g / denominator))
}
