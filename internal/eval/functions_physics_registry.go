// internal/eval/functions_physics_registry.go

package eval

import "github.com/0xsj/numio/pkg/types"

// PhysicsFunctionInfo holds metadata about a physics function.
type PhysicsFunctionInfo struct {
	Fn          func([]types.Value) types.Value
	Description string
	Args        string
	Category    string
	Formula     string
}

// PhysicsFunctionRegistry maps function names to their implementations and metadata.
var PhysicsFunctionRegistry = map[string]PhysicsFunctionInfo{

	// ════════════════════════════════════════════════════════════════
	// CONSTANTS (functions_physics_constants.go)
	// ════════════════════════════════════════════════════════════════

	// Constant accessors
	"c":           {Fn: FnSpeedOfLight, Description: "Speed of light in vacuum", Args: "", Category: "constants", Formula: "299792458 m/s"},
	"planck":      {Fn: FnPlanckConstant, Description: "Planck constant", Args: "", Category: "constants", Formula: "6.62607015×10⁻³⁴ J·s"},
	"hbar":        {Fn: FnReducedPlanck, Description: "Reduced Planck constant", Args: "", Category: "constants", Formula: "ℏ = h/(2π)"},
	"g":           {Fn: FnGravityEarth, Description: "Standard gravity on Earth", Args: "", Category: "constants", Formula: "9.80665 m/s²"},
	"G":           {Fn: FnGravitationalConstant, Description: "Gravitational constant", Args: "", Category: "constants", Formula: "6.67430×10⁻¹¹ N·m²/kg²"},
	"e":           {Fn: FnElementaryCharge, Description: "Elementary charge", Args: "", Category: "constants", Formula: "1.602176634×10⁻¹⁹ C"},
	"me":          {Fn: FnElectronMass, Description: "Electron mass", Args: "", Category: "constants", Formula: "9.1093837015×10⁻³¹ kg"},
	"mp":          {Fn: FnProtonMass, Description: "Proton mass", Args: "", Category: "constants", Formula: "1.67262192369×10⁻²⁷ kg"},
	"mn":          {Fn: FnNeutronMass, Description: "Neutron mass", Args: "", Category: "constants", Formula: "1.67492749804×10⁻²⁷ kg"},
	"kb":          {Fn: FnBoltzmannConstant, Description: "Boltzmann constant", Args: "", Category: "constants", Formula: "1.380649×10⁻²³ J/K"},
	"na":          {Fn: FnAvogadroNumber, Description: "Avogadro number", Args: "", Category: "constants", Formula: "6.02214076×10²³ mol⁻¹"},
	"rgas":        {Fn: FnGasConstant, Description: "Ideal gas constant", Args: "", Category: "constants", Formula: "8.314462618 J/(mol·K)"},
	"epsilon0":    {Fn: FnVacuumPermittivity, Description: "Vacuum permittivity", Args: "", Category: "constants", Formula: "8.8541878128×10⁻¹² F/m"},
	"mu0":         {Fn: FnVacuumPermeability, Description: "Vacuum permeability", Args: "", Category: "constants", Formula: "1.25663706212×10⁻⁶ H/m"},
	"ke":          {Fn: FnCoulombConstant, Description: "Coulomb constant", Args: "", Category: "constants", Formula: "8.9875517923×10⁹ N·m²/C²"},
	"stefanboltz": {Fn: FnStefanBoltzmannConst, Description: "Stefan-Boltzmann constant", Args: "", Category: "constants", Formula: "5.670374419×10⁻⁸ W/(m²·K⁴)"},
	"wien":        {Fn: FnWienConstant, Description: "Wien displacement constant", Args: "", Category: "constants", Formula: "2.897771955×10⁻³ m·K"},
	"rydberg":     {Fn: FnRydbergConstant, Description: "Rydberg constant", Args: "", Category: "constants", Formula: "1.0973731568160×10⁷ m⁻¹"},
	"a0":          {Fn: FnBohrRadiusConst, Description: "Bohr radius", Args: "", Category: "constants", Formula: "5.29177210903×10⁻¹¹ m"},
	"alpha":       {Fn: FnFineStructure, Description: "Fine-structure constant", Args: "", Category: "constants", Formula: "≈ 1/137"},
	"ev":          {Fn: FnElectronVolt, Description: "Electron volt in joules", Args: "", Category: "constants", Formula: "1.602176634×10⁻¹⁹ J"},
	"amu":         {Fn: FnAtomicMassUnit, Description: "Atomic mass unit", Args: "", Category: "constants", Formula: "1.66053906660×10⁻²⁷ kg"},

	// Derived constant functions
	"comptonwavelength":       {Fn: FnComptonWavelength, Description: "Compton wavelength of electron", Args: "", Category: "constants", Formula: "λ_C = h/(m_e c)"},
	"bohrmagneton":            {Fn: FnBohrMagneton, Description: "Bohr magneton", Args: "", Category: "constants", Formula: "μ_B = eℏ/(2m_e)"},
	"nuclearmagneton":         {Fn: FnNuclearMagneton, Description: "Nuclear magneton", Args: "", Category: "constants", Formula: "μ_N = eℏ/(2m_p)"},
	"classicalelectronradius": {Fn: FnClassicalElectronRadius, Description: "Classical electron radius", Args: "", Category: "constants", Formula: "r_e = k_e e²/(m_e c²)"},
	"plancklengthn":           {Fn: FnPlanckLength, Description: "Planck length", Args: "", Category: "constants", Formula: "l_P = √(ℏG/c³)"},
	"plancktime":              {Fn: FnPlanckTime, Description: "Planck time", Args: "", Category: "constants", Formula: "t_P = √(ℏG/c⁵)"},
	"planckmass":              {Fn: FnPlanckMass, Description: "Planck mass", Args: "", Category: "constants", Formula: "m_P = √(ℏc/G)"},
	"planckenergy":            {Fn: FnPlanckEnergy, Description: "Planck energy", Args: "", Category: "constants", Formula: "E_P = √(ℏc⁵/G)"},
	"plancktemperature":       {Fn: FnPlanckTemperature, Description: "Planck temperature", Args: "", Category: "constants", Formula: "T_P = √(ℏc⁵/(Gk_B²))"},

	// ════════════════════════════════════════════════════════════════
	// MECHANICS (functions_physics_mechanics.go)
	// ════════════════════════════════════════════════════════════════

	// Force & Motion
	"force":        {Fn: FnForce, Description: "Newton's second law", Args: "mass, acceleration", Category: "mechanics", Formula: "F = ma"},
	"mass":         {Fn: FnMass, Description: "Mass from force and acceleration", Args: "force, acceleration", Category: "mechanics", Formula: "m = F/a"},
	"acceleration": {Fn: FnAcceleration, Description: "Acceleration from force and mass", Args: "force, mass", Category: "mechanics", Formula: "a = F/m"},
	"weight":       {Fn: FnWeight, Description: "Weight (gravitational force)", Args: "mass, [g]", Category: "mechanics", Formula: "W = mg"},

	// Momentum & Impulse
	"momentum":     {Fn: FnMomentum, Description: "Linear momentum", Args: "mass, velocity", Category: "mechanics", Formula: "p = mv"},
	"impulse":      {Fn: FnImpulse, Description: "Impulse", Args: "force, time", Category: "mechanics", Formula: "J = Ft"},
	"momentumcons": {Fn: FnMomentumConservation, Description: "Final velocity after inelastic collision", Args: "m1, v1, m2, v2", Category: "mechanics", Formula: "v_f = (m₁v₁+m₂v₂)/(m₁+m₂)"},

	// Energy
	"kineticenergy":   {Fn: FnKineticEnergy, Description: "Kinetic energy", Args: "mass, velocity", Category: "mechanics", Formula: "KE = ½mv²"},
	"potentialenergy": {Fn: FnPotentialEnergy, Description: "Gravitational potential energy", Args: "mass, height, [g]", Category: "mechanics", Formula: "PE = mgh"},
	"elasticpe":       {Fn: FnElasticPE, Description: "Elastic potential energy", Args: "spring_constant, displacement", Category: "mechanics", Formula: "PE = ½kx²"},
	"velocityfromke":  {Fn: FnVelocityFromKE, Description: "Velocity from kinetic energy", Args: "KE, mass", Category: "mechanics", Formula: "v = √(2KE/m)"},

	// Work & Power
	"work":       {Fn: FnWork, Description: "Work done by force", Args: "force, distance, [angle_deg]", Category: "mechanics", Formula: "W = Fd·cos(θ)"},
	"powermech":  {Fn: FnPowerMech, Description: "Mechanical power", Args: "work, time", Category: "mechanics", Formula: "P = W/t"},
	"powerfv":    {Fn: FnPowerFromFV, Description: "Power from force and velocity", Args: "force, velocity", Category: "mechanics", Formula: "P = Fv"},
	"efficiency": {Fn: FnEfficiency, Description: "Mechanical efficiency", Args: "output, input", Category: "mechanics", Formula: "η = (out/in)×100%"},

	// Kinematics
	"displacement":    {Fn: FnDisplacement, Description: "Kinematic displacement", Args: "v0, time, acceleration", Category: "mechanics", Formula: "s = v₀t + ½at²"},
	"finalvelocity":   {Fn: FnFinalVelocity, Description: "Final velocity", Args: "v0, acceleration, time", Category: "mechanics", Formula: "v = v₀ + at"},
	"finalvelocitysq": {Fn: FnFinalVelocitySq, Description: "Final velocity from displacement", Args: "v0, acceleration, displacement", Category: "mechanics", Formula: "v = √(v₀² + 2as)"},
	"timefromkin":     {Fn: FnTimeFromKinematics, Description: "Time from velocities", Args: "v0, v, acceleration", Category: "mechanics", Formula: "t = (v-v₀)/a"},
	"freefall":        {Fn: FnFreeFall, Description: "Free fall distance", Args: "time, [g]", Category: "mechanics", Formula: "h = ½gt²"},
	"freefalltime":    {Fn: FnFreeFallTime, Description: "Free fall time", Args: "height, [g]", Category: "mechanics", Formula: "t = √(2h/g)"},
	"freefallvel":     {Fn: FnFreeFallVelocity, Description: "Free fall velocity", Args: "height, [g]", Category: "mechanics", Formula: "v = √(2gh)"},

	// Projectile Motion
	"projectilerange":  {Fn: FnProjectileRange, Description: "Projectile range", Args: "v0, angle_deg, [g]", Category: "mechanics", Formula: "R = v₀²sin(2θ)/g"},
	"projectileheight": {Fn: FnProjectileHeight, Description: "Maximum projectile height", Args: "v0, angle_deg, [g]", Category: "mechanics", Formula: "H = v₀²sin²(θ)/(2g)"},
	"projectiletime":   {Fn: FnProjectileTime, Description: "Projectile flight time", Args: "v0, angle_deg, [g]", Category: "mechanics", Formula: "T = 2v₀sin(θ)/g"},

	// Circular Motion
	"centripetal":      {Fn: FnCentripetal, Description: "Centripetal force", Args: "mass, velocity, radius", Category: "mechanics", Formula: "F = mv²/r"},
	"centripetalaccel": {Fn: FnCentripetalAccel, Description: "Centripetal acceleration", Args: "velocity, radius", Category: "mechanics", Formula: "a = v²/r"},
	"angularvel":       {Fn: FnAngularVelocity, Description: "Angular velocity", Args: "velocity, radius OR frequency", Category: "mechanics", Formula: "ω = v/r or 2πf"},
	"linearvel":        {Fn: FnLinearVelocity, Description: "Linear velocity from angular", Args: "omega, radius", Category: "mechanics", Formula: "v = ωr"},
	"period":           {Fn: FnPeriod, Description: "Period from frequency", Args: "frequency, [type]", Category: "mechanics", Formula: "T = 1/f"},
	"frequencymech":    {Fn: FnFrequencyMech, Description: "Frequency from period", Args: "period", Category: "mechanics", Formula: "f = 1/T"},

	// Rotational Mechanics
	"torque":          {Fn: FnTorque, Description: "Torque", Args: "radius, force, [angle_deg]", Category: "mechanics", Formula: "τ = rF·sin(θ)"},
	"momentofinertia": {Fn: FnMomentOfInertia, Description: "Moment of inertia (point mass)", Args: "mass, radius", Category: "mechanics", Formula: "I = mr²"},
	"moisphere":       {Fn: FnMomentOfInertiaSphere, Description: "Moment of inertia (solid sphere)", Args: "mass, radius", Category: "mechanics", Formula: "I = (2/5)mr²"},
	"moicylinder":     {Fn: FnMomentOfInertiaCylinder, Description: "Moment of inertia (solid cylinder)", Args: "mass, radius", Category: "mechanics", Formula: "I = (1/2)mr²"},
	"moirod":          {Fn: FnMomentOfInertiaRod, Description: "Moment of inertia (rod, center)", Args: "mass, length", Category: "mechanics", Formula: "I = (1/12)mL²"},
	"angularmomentum": {Fn: FnAngularMomentum, Description: "Angular momentum", Args: "moment_of_inertia, omega", Category: "mechanics", Formula: "L = Iω"},
	"rotationalke":    {Fn: FnRotationalKE, Description: "Rotational kinetic energy", Args: "moment_of_inertia, omega", Category: "mechanics", Formula: "KE = ½Iω²"},
	"angularaccel":    {Fn: FnAngularAccel, Description: "Angular acceleration", Args: "torque, moment_of_inertia", Category: "mechanics", Formula: "α = τ/I"},

	// Simple Harmonic Motion
	"shmperiodspring":   {Fn: FnSHMPeriodSpring, Description: "SHM period (spring)", Args: "mass, spring_constant", Category: "mechanics", Formula: "T = 2π√(m/k)"},
	"shmperiodpendulum": {Fn: FnSHMPeriodPendulum, Description: "SHM period (pendulum)", Args: "length, [g]", Category: "mechanics", Formula: "T = 2π√(L/g)"},
	"shmfrequency":      {Fn: FnSHMFrequency, Description: "SHM frequency", Args: "spring_constant, mass", Category: "mechanics", Formula: "f = √(k/m)/(2π)"},
	"shmdisplacement":   {Fn: FnSHMDisplacement, Description: "SHM displacement", Args: "amplitude, omega, time, [phase]", Category: "mechanics", Formula: "x = A·cos(ωt+φ)"},
	"shmvelocity":       {Fn: FnSHMVelocity, Description: "SHM velocity", Args: "amplitude, omega, time, [phase]", Category: "mechanics", Formula: "v = -Aω·sin(ωt+φ)"},
	"shmmaxvel":         {Fn: FnSHMMaxVelocity, Description: "SHM maximum velocity", Args: "amplitude, omega", Category: "mechanics", Formula: "v_max = Aω"},

	// Friction & Drag
	"friction":    {Fn: FnFriction, Description: "Friction force", Args: "mu, normal_force", Category: "mechanics", Formula: "f = μN"},
	"dragforce":   {Fn: FnDragForce, Description: "Drag force", Args: "density, velocity, drag_coeff, area", Category: "mechanics", Formula: "F = ½ρv²CdA"},
	"terminalvel": {Fn: FnTerminalVelocity, Description: "Terminal velocity", Args: "mass, drag_coeff, area, [density]", Category: "mechanics", Formula: "v_t = √(2mg/(ρCdA))"},

	// ════════════════════════════════════════════════════════════════
	// GRAVITY (functions_physics_gravity.go)
	// ════════════════════════════════════════════════════════════════

	// Gravitational Force & Field
	"gravity":        {Fn: FnGravity, Description: "Gravitational force", Args: "m1, m2, distance", Category: "gravity", Formula: "F = Gm₁m₂/r²"},
	"gravfield":      {Fn: FnGravField, Description: "Gravitational field strength", Args: "mass, distance", Category: "gravity", Formula: "g = GM/r²"},
	"gravpotential":  {Fn: FnGravPotential, Description: "Gravitational potential", Args: "mass, distance", Category: "gravity", Formula: "V = -GM/r"},
	"gravpe":         {Fn: FnGravPotentialEnergy, Description: "Gravitational potential energy", Args: "m1, m2, distance", Category: "gravity", Formula: "U = -Gm₁m₂/r"},
	"surfacegravity": {Fn: FnSurfaceGravity, Description: "Surface gravity", Args: "mass, radius", Category: "gravity", Formula: "g = GM/R²"},

	// Orbital Mechanics
	"escapevel":             {Fn: FnEscapeVelocity, Description: "Escape velocity", Args: "mass, radius", Category: "gravity", Formula: "v = √(2GM/r)"},
	"orbitalvel":            {Fn: FnOrbitalVelocity, Description: "Orbital velocity", Args: "mass, radius", Category: "gravity", Formula: "v = √(GM/r)"},
	"orbitalperiod":         {Fn: FnOrbitalPeriod, Description: "Orbital period", Args: "mass, radius", Category: "gravity", Formula: "T = 2π√(r³/GM)"},
	"orbitalradius":         {Fn: FnOrbitalRadius, Description: "Orbital radius from period", Args: "mass, period", Category: "gravity", Formula: "r = ∛(GMT²/4π²)"},
	"orbitalenergy":         {Fn: FnOrbitalEnergy, Description: "Orbital energy", Args: "M, m, radius", Category: "gravity", Formula: "E = -GMm/(2r)"},
	"specificorbitalenergy": {Fn: FnSpecificOrbitalEnergy, Description: "Specific orbital energy", Args: "mass, semi_major_axis", Category: "gravity", Formula: "ε = -GM/(2a)"},

	// Elliptical Orbits
	"periapsis":     {Fn: FnPeriapsis, Description: "Periapsis distance", Args: "semi_major_axis, eccentricity", Category: "gravity", Formula: "r_p = a(1-e)"},
	"apoapsis":      {Fn: FnApoapsis, Description: "Apoapsis distance", Args: "semi_major_axis, eccentricity", Category: "gravity", Formula: "r_a = a(1+e)"},
	"semimajoraxis": {Fn: FnSemiMajorAxis, Description: "Semi-major axis", Args: "periapsis, apoapsis", Category: "gravity", Formula: "a = (r_p+r_a)/2"},
	"eccentricity":  {Fn: FnEccentricity, Description: "Orbital eccentricity", Args: "periapsis, apoapsis", Category: "gravity", Formula: "e = (r_a-r_p)/(r_a+r_p)"},
	"orbitalvelatr": {Fn: FnOrbitalVelAtRadius, Description: "Velocity at radius (vis-viva)", Args: "mass, radius, semi_major_axis", Category: "gravity", Formula: "v = √(GM(2/r-1/a))"},
	"periapsisvel":  {Fn: FnPeriapsisVelocity, Description: "Velocity at periapsis", Args: "mass, semi_major_axis, eccentricity", Category: "gravity", Formula: "vis-viva at periapsis"},
	"apoapsisvel":   {Fn: FnApoapsisVelocity, Description: "Velocity at apoapsis", Args: "mass, semi_major_axis, eccentricity", Category: "gravity", Formula: "vis-viva at apoapsis"},

	// Orbital Maneuvers
	"hohmanntransfer": {Fn: FnHohmannTransfer, Description: "Hohmann transfer Δv", Args: "mass, r1, r2", Category: "gravity", Formula: "Total Δv for transfer"},
	"hohmanntime":     {Fn: FnHohmannTime, Description: "Hohmann transfer time", Args: "mass, r1, r2", Category: "gravity", Formula: "t = π√(a³/GM)"},
	"dvcircularize":   {Fn: FnDeltaVCircularize, Description: "Δv to circularize", Args: "mass, current_velocity, radius", Category: "gravity", Formula: "Δv = |v_circ - v|"},

	// Sphere of Influence
	"hillsphere": {Fn: FnHillSphere, Description: "Hill sphere radius", Args: "m, M, semi_major_axis", Category: "gravity", Formula: "r_H = a(m/3M)^⅓"},
	"soi":        {Fn: FnSphereOfInfluence, Description: "Sphere of influence", Args: "m, M, semi_major_axis", Category: "gravity", Formula: "r_SOI = a(m/M)^⅖"},
	"tidalforce": {Fn: FnTidalForce, Description: "Tidal force", Args: "M, m, object_size, distance", Category: "gravity", Formula: "ΔF = 2GMmr/d³"},
	"tidalaccel": {Fn: FnTidalAcceleration, Description: "Tidal acceleration", Args: "M, object_size, distance", Category: "gravity", Formula: "a = 2GMr/d³"},
	"rochelimit": {Fn: FnRocheLimit, Description: "Roche limit", Args: "R_primary, density_primary, density_secondary", Category: "gravity", Formula: "d = 2.44R(ρM/ρm)^⅓"},

	// Black Holes
	"schwarzschild":    {Fn: FnSchwarzschildRadius, Description: "Schwarzschild radius", Args: "mass", Category: "gravity", Formula: "r_s = 2GM/c²"},
	"blackholemass":    {Fn: FnBlackHoleMass, Description: "Black hole mass from radius", Args: "schwarzschild_radius", Category: "gravity", Formula: "M = r_s c²/(2G)"},
	"blackholetemp":    {Fn: FnBlackHoleTemperature, Description: "Hawking temperature", Args: "mass", Category: "gravity", Formula: "T = ℏc³/(8πGMk_B)"},
	"blackholeentropy": {Fn: FnBlackHoleEntropy, Description: "Bekenstein-Hawking entropy", Args: "mass", Category: "gravity", Formula: "S = Ak_B c³/(4Gℏ)"},

	// Special Orbits
	"geosyncradius": {Fn: FnGeosyncRadius, Description: "Geosynchronous radius", Args: "mass, [period]", Category: "gravity", Formula: "r = ∛(GMT²/4π²)"},
	"geosyncalt":    {Fn: FnGeosyncAltitude, Description: "Geosynchronous altitude", Args: "mass, surface_radius, [period]", Category: "gravity", Formula: "alt = r - R"},
	"lagrangel1":    {Fn: FnLagrangeL1, Description: "L1 Lagrange point distance", Args: "m, M, separation", Category: "gravity", Formula: "r ≈ R(m/3M)^⅓"},

	// Gravitational Waves
	"gravwaveluminosity": {Fn: FnGravWaveLuminosity, Description: "Gravitational wave luminosity", Args: "m1, m2, separation", Category: "gravity", Formula: "Binary GW power"},
	"gravwavefreq":       {Fn: FnGravWaveFrequency, Description: "Gravitational wave frequency", Args: "m1, m2, separation", Category: "gravity", Formula: "f = (1/π)√(G(m₁+m₂)/r³)"},

	// ════════════════════════════════════════════════════════════════
	// THERMODYNAMICS (functions_physics_thermo.go)
	// ════════════════════════════════════════════════════════════════

	// Ideal Gas Law
	"idealgaspressure": {Fn: FnIdealGasPressure, Description: "Ideal gas pressure", Args: "n, T, V", Category: "thermo", Formula: "P = nRT/V"},
	"idealgasvolume":   {Fn: FnIdealGasVolume, Description: "Ideal gas volume", Args: "n, T, P", Category: "thermo", Formula: "V = nRT/P"},
	"idealgastemp":     {Fn: FnIdealGasTemp, Description: "Ideal gas temperature", Args: "P, V, n", Category: "thermo", Formula: "T = PV/(nR)"},
	"idealgasmoles":    {Fn: FnIdealGasMoles, Description: "Ideal gas moles", Args: "P, V, T", Category: "thermo", Formula: "n = PV/(RT)"},
	"idealgasdensity":  {Fn: FnIdealGasDensity, Description: "Ideal gas density", Args: "P, M, T", Category: "thermo", Formula: "ρ = PM/(RT)"},

	// Heat Transfer
	"heat":               {Fn: FnHeat, Description: "Heat transfer", Args: "mass, specific_heat, delta_T", Category: "thermo", Formula: "Q = mcΔT"},
	"heatmoles":          {Fn: FnHeatMoles, Description: "Heat transfer (molar)", Args: "n, C, delta_T", Category: "thermo", Formula: "Q = nCΔT"},
	"finaltemp":          {Fn: FnFinalTemp, Description: "Final temperature", Args: "T_initial, Q, mass, specific_heat", Category: "thermo", Formula: "T_f = T_i + Q/(mc)"},
	"thermalequilibrium": {Fn: FnThermalEquilibrium, Description: "Thermal equilibrium temp", Args: "m1, c1, T1, m2, c2, T2", Category: "thermo", Formula: "T_f = Σ(mcT)/Σ(mc)"},
	"latentheat":         {Fn: FnLatentHeat, Description: "Latent heat", Args: "mass, latent_heat", Category: "thermo", Formula: "Q = mL"},

	// Conduction
	"heatconduction":    {Fn: FnHeatConduction, Description: "Heat conduction rate", Args: "k, A, delta_T, thickness", Category: "thermo", Formula: "P = kAΔT/Δx"},
	"thermalresistance": {Fn: FnThermalResistance, Description: "Thermal resistance", Args: "thickness, k, A", Category: "thermo", Formula: "R = Δx/(kA)"},

	// Radiation
	"stefanboltzmann":  {Fn: FnStefanBoltzmann, Description: "Stefan-Boltzmann radiation", Args: "emissivity, A, T", Category: "thermo", Formula: "P = εσAT⁴"},
	"blackbodypower":   {Fn: FnBlackBodyPower, Description: "Black body power", Args: "A, T", Category: "thermo", Formula: "P = σAT⁴"},
	"netradiation":     {Fn: FnNetRadiation, Description: "Net radiation", Args: "emissivity, A, T_object, T_surroundings", Category: "thermo", Formula: "P = εσA(T⁴-T_surr⁴)"},
	"wiendisplacement": {Fn: FnWienDisplacement, Description: "Wien peak wavelength", Args: "T", Category: "thermo", Formula: "λ = b/T"},
	"tempfromwien":     {Fn: FnTempFromWien, Description: "Temperature from Wien", Args: "wavelength", Category: "thermo", Formula: "T = b/λ"},

	// Thermodynamic Processes
	"isothermalwork": {Fn: FnIsothermalWork, Description: "Isothermal work", Args: "n, T, V1, V2", Category: "thermo", Formula: "W = nRT ln(V₂/V₁)"},
	"isobaricwork":   {Fn: FnIsobaricWork, Description: "Isobaric work", Args: "P, V1, V2", Category: "thermo", Formula: "W = PΔV"},
	"adiabaticwork":  {Fn: FnAdiabaticWork, Description: "Adiabatic work", Args: "P1, V1, P2, V2, gamma", Category: "thermo", Formula: "W = (P₁V₁-P₂V₂)/(γ-1)"},
	"adiabaticfinal": {Fn: FnAdiabaticFinal, Description: "Adiabatic final state", Args: "initial, V1, V2, gamma, type", Category: "thermo", Formula: "P or T adiabatic relation"},

	// Entropy
	"entropy":         {Fn: FnEntropy, Description: "Entropy change", Args: "Q, T", Category: "thermo", Formula: "ΔS = Q/T"},
	"entropyidealgas": {Fn: FnEntropyIdealGas, Description: "Entropy change (ideal gas)", Args: "n, Cv, T1, T2, V1, V2", Category: "thermo", Formula: "ΔS = nCv ln(T₂/T₁) + nR ln(V₂/V₁)"},
	"entropymixing":   {Fn: FnEntropyMixing, Description: "Entropy of mixing", Args: "n1, n2, ...", Category: "thermo", Formula: "ΔS = -R Σ(nᵢ ln xᵢ)"},

	// Heat Engines
	"carnotefficiency": {Fn: FnCarnotEfficiency, Description: "Carnot efficiency", Args: "T_hot, T_cold", Category: "thermo", Formula: "η = 1 - T_c/T_h"},
	"heatengineeff":    {Fn: FnHeatEngineEfficiency, Description: "Heat engine efficiency", Args: "Q_hot, Q_cold", Category: "thermo", Formula: "η = (Q_h-Q_c)/Q_h"},
	"heatengwork":      {Fn: FnHeatEngineWork, Description: "Heat engine work", Args: "Q_hot, Q_cold", Category: "thermo", Formula: "W = Q_h - Q_c"},
	"cop":              {Fn: FnCOP, Description: "Coefficient of performance", Args: "Q_hot, Q_cold, type", Category: "thermo", Formula: "COP = Q/W"},
	"carnotcop":        {Fn: FnCarnotCOP, Description: "Carnot COP", Args: "T_hot, T_cold, type", Category: "thermo", Formula: "COP_Carnot"},

	// Kinetic Theory
	"rmsvelocity":      {Fn: FnRMSVelocity, Description: "RMS velocity", Args: "T, molar_mass", Category: "thermo", Formula: "v = √(3RT/M)"},
	"meanvelocity":     {Fn: FnMeanVelocity, Description: "Mean velocity", Args: "T, molar_mass", Category: "thermo", Formula: "v = √(8RT/πM)"},
	"mostprobablevel":  {Fn: FnMostProbableVelocity, Description: "Most probable velocity", Args: "T, molar_mass", Category: "thermo", Formula: "v = √(2RT/M)"},
	"meanfreepath":     {Fn: FnMeanFreePath, Description: "Mean free path", Args: "T, P, collision_diameter", Category: "thermo", Formula: "λ = kT/(√2·P·σ)"},
	"collisionfreq":    {Fn: FnCollisionFrequency, Description: "Collision frequency", Args: "T, P, molar_mass, collision_diameter", Category: "thermo", Formula: "z = v/λ"},
	"avgkineticenergy": {Fn: FnAvgKineticEnergy, Description: "Average kinetic energy", Args: "T", Category: "thermo", Formula: "KE = (3/2)kT"},
	"internalenergy":   {Fn: FnInternalEnergy, Description: "Internal energy", Args: "n, T, degrees_of_freedom", Category: "thermo", Formula: "U = (f/2)nRT"},

	// Temperature Conversions
	"celsiustokelvin":     {Fn: FnCelsiusToKelvin, Description: "Celsius to Kelvin", Args: "celsius", Category: "thermo", Formula: "K = C + 273.15"},
	"kelvintocelsius":     {Fn: FnKelvinToCelsius, Description: "Kelvin to Celsius", Args: "kelvin", Category: "thermo", Formula: "C = K - 273.15"},
	"fahrenheittocelsius": {Fn: FnFahrenheitToCelsius, Description: "Fahrenheit to Celsius", Args: "fahrenheit", Category: "thermo", Formula: "C = (F-32)×5/9"},
	"celsiustofahrenheit": {Fn: FnCelsiusToFahrenheit, Description: "Celsius to Fahrenheit", Args: "celsius", Category: "thermo", Formula: "F = C×9/5 + 32"},
	"fahrenheittokelvin":  {Fn: FnFahrenheitToKelvin, Description: "Fahrenheit to Kelvin", Args: "fahrenheit", Category: "thermo", Formula: "K = (F-32)×5/9 + 273.15"},
	"kelvintofahrenheit":  {Fn: FnKelvinToFahrenheit, Description: "Kelvin to Fahrenheit", Args: "kelvin", Category: "thermo", Formula: "F = (K-273.15)×9/5 + 32"},

	// ════════════════════════════════════════════════════════════════
	// WAVES & OPTICS (functions_physics_waves.go)
	// ════════════════════════════════════════════════════════════════

	// Wave Fundamentals
	"wavelength":    {Fn: FnWavelength, Description: "Wavelength", Args: "velocity, frequency", Category: "waves", Formula: "λ = v/f"},
	"frequencywave": {Fn: FnFrequencyWave, Description: "Wave frequency", Args: "velocity, wavelength", Category: "waves", Formula: "f = v/λ"},
	"wavevelocity":  {Fn: FnWaveVelocity, Description: "Wave velocity", Args: "frequency, wavelength", Category: "waves", Formula: "v = fλ"},
	"wavenumber":    {Fn: FnWaveNumber, Description: "Wave number", Args: "wavelength", Category: "waves", Formula: "k = 2π/λ"},
	"angularfreq":   {Fn: FnAngularFrequency, Description: "Angular frequency", Args: "frequency", Category: "waves", Formula: "ω = 2πf"},
	"waveperiod":    {Fn: FnWavePeriod, Description: "Wave period", Args: "frequency", Category: "waves", Formula: "T = 1/f"},

	// Wave Equation
	"wavedisplacement": {Fn: FnWaveDisplacement, Description: "Wave displacement", Args: "A, k, x, omega, t, [phase]", Category: "waves", Formula: "y = A sin(kx-ωt+φ)"},
	"wavevelparticle":  {Fn: FnWaveVelocityParticle, Description: "Particle velocity in wave", Args: "A, k, x, omega, t, [phase]", Category: "waves", Formula: "v = -Aω cos(kx-ωt+φ)"},
	"waveintensity":    {Fn: FnWaveIntensity, Description: "Wave intensity", Args: "density, velocity, omega, amplitude", Category: "waves", Formula: "I = ½ρvω²A²"},
	"waveenergy":       {Fn: FnWaveEnergy, Description: "Wave energy per length", Args: "linear_density, omega, amplitude", Category: "waves", Formula: "E/L = ½μω²A²"},
	"wavepower":        {Fn: FnWavePower, Description: "Wave power", Args: "linear_density, omega, amplitude, velocity", Category: "waves", Formula: "P = ½μω²A²v"},

	// String Waves
	"stringwavevel":     {Fn: FnStringWaveVelocity, Description: "String wave velocity", Args: "tension, linear_density", Category: "waves", Formula: "v = √(T/μ)"},
	"stringfundamental": {Fn: FnStringFundamental, Description: "String fundamental frequency", Args: "length, tension, linear_density", Category: "waves", Formula: "f₁ = (1/2L)√(T/μ)"},
	"stringharmonic":    {Fn: FnStringHarmonic, Description: "String harmonic", Args: "n, length, tension, linear_density", Category: "waves", Formula: "fₙ = nf₁"},

	// Sound
	"soundvelocity":    {Fn: FnSoundVelocity, Description: "Speed of sound", Args: "modulus, density", Category: "waves", Formula: "v = √(B/ρ)"},
	"soundvelocitygas": {Fn: FnSoundVelocityGas, Description: "Speed of sound in gas", Args: "gamma, T, molar_mass", Category: "waves", Formula: "v = √(γRT/M)"},
	"soundintensity":   {Fn: FnSoundIntensity, Description: "Sound intensity", Args: "power, distance", Category: "waves", Formula: "I = P/(4πr²)"},
	"decibellevel":     {Fn: FnDecibelLevel, Description: "Decibel level", Args: "intensity", Category: "waves", Formula: "β = 10 log(I/I₀)"},
	"intensityfromdb":  {Fn: FnIntensityFromDecibel, Description: "Intensity from decibels", Args: "decibel_level", Category: "waves", Formula: "I = I₀×10^(β/10)"},
	"decibleadd":       {Fn: FnDecibelAdd, Description: "Add decibel levels", Args: "dB1, dB2, ...", Category: "waves", Formula: "β_total = 10 log(Σ10^(βᵢ/10))"},

	// Doppler Effect
	"dopplerfreq":  {Fn: FnDopplerFrequency, Description: "Doppler frequency", Args: "f_source, v_wave, v_observer, v_source", Category: "waves", Formula: "f' = f(v+v_o)/(v-v_s)"},
	"dopplershift": {Fn: FnDopplerShift, Description: "Doppler shift ratio", Args: "v_relative, v_wave", Category: "waves", Formula: "Δf/f = v_rel/v"},
	"reldoppler":   {Fn: FnRelativisticDoppler, Description: "Relativistic Doppler", Args: "f_source, velocity", Category: "waves", Formula: "f' = f√((1+β)/(1-β))"},

	// Interference
	"beatfreq":          {Fn: FnBeatFrequency, Description: "Beat frequency", Args: "f1, f2", Category: "waves", Formula: "f_beat = |f₁-f₂|"},
	"pathdiff":          {Fn: FnPathDifference, Description: "Path difference", Args: "separation, angle_deg", Category: "waves", Formula: "Δ = d sin(θ)"},
	"constructiveangle": {Fn: FnConstructiveAngle, Description: "Constructive interference angle", Args: "separation, wavelength, order", Category: "waves", Formula: "d sin(θ) = mλ"},
	"destructiveangle":  {Fn: FnDestructiveAngle, Description: "Destructive interference angle", Args: "separation, wavelength, order", Category: "waves", Formula: "d sin(θ) = (m+½)λ"},

	// Standing Waves
	"standingwavenodes":     {Fn: FnStandingWaveNodes, Description: "Number of nodes", Args: "harmonic_number", Category: "waves", Formula: "n+1 nodes"},
	"standingwaveantinodes": {Fn: FnStandingWaveAntinodes, Description: "Number of antinodes", Args: "harmonic_number", Category: "waves", Formula: "n antinodes"},
	"openpipefund":          {Fn: FnOpenPipeFundamental, Description: "Open pipe fundamental", Args: "velocity, length", Category: "waves", Formula: "f = v/(2L)"},
	"closedpipefund":        {Fn: FnClosedPipeFundamental, Description: "Closed pipe fundamental", Args: "velocity, length", Category: "waves", Formula: "f = v/(4L)"},
	"openpipeharmonic":      {Fn: FnOpenPipeHarmonic, Description: "Open pipe harmonic", Args: "n, velocity, length", Category: "waves", Formula: "f = nv/(2L)"},
	"closedpipeharmonic":    {Fn: FnClosedPipeHarmonic, Description: "Closed pipe harmonic", Args: "n_odd, velocity, length", Category: "waves", Formula: "f = nv/(4L)"},

	// Optics - Refraction
	"snellslaw":       {Fn: FnSnellsLaw, Description: "Snell's law", Args: "n1, angle1_deg, n2", Category: "optics", Formula: "n₁sinθ₁ = n₂sinθ₂"},
	"criticalangle":   {Fn: FnCriticalAngle, Description: "Critical angle", Args: "n1, n2", Category: "optics", Formula: "θ_c = arcsin(n₂/n₁)"},
	"brewsterangle":   {Fn: FnBrewsterAngle, Description: "Brewster's angle", Args: "n1, n2", Category: "optics", Formula: "θ_B = arctan(n₂/n₁)"},
	"refractiveindex": {Fn: FnRefractiveIndex, Description: "Refractive index", Args: "velocity_in_medium", Category: "optics", Formula: "n = c/v"},

	// Optics - Lenses & Mirrors
	"thinlensimage": {Fn: FnThinLensImage, Description: "Thin lens image distance", Args: "focal_length, object_distance", Category: "optics", Formula: "1/f = 1/d_o + 1/d_i"},
	"thinlensfocal": {Fn: FnThinLensFocal, Description: "Thin lens focal length", Args: "object_distance, image_distance", Category: "optics", Formula: "1/f = 1/d_o + 1/d_i"},
	"magnification": {Fn: FnMagnification, Description: "Magnification", Args: "object_distance, image_distance", Category: "optics", Formula: "M = -d_i/d_o"},
	"lensmaker":     {Fn: FnLensMaker, Description: "Lensmaker equation", Args: "n, R1, R2", Category: "optics", Formula: "1/f = (n-1)(1/R₁-1/R₂)"},
	"mirrorimage":   {Fn: FnMirrorImage, Description: "Mirror image distance", Args: "radius_of_curvature, object_distance", Category: "optics", Formula: "1/f = 1/d_o + 1/d_i"},
	"dioptricpower": {Fn: FnDioptricPower, Description: "Dioptic power", Args: "focal_length", Category: "optics", Formula: "P = 1/f"},

	// Optics - Diffraction
	"singleslitmin":      {Fn: FnSingleSlitMinima, Description: "Single slit minima", Args: "slit_width, wavelength, order", Category: "optics", Formula: "a sin(θ) = mλ"},
	"diffractiongrating": {Fn: FnDiffractionGrating, Description: "Diffraction grating maxima", Args: "spacing, wavelength, order", Category: "optics", Formula: "d sin(θ) = mλ"},
	"rayleighcriterion":  {Fn: FnRayleighCriterion, Description: "Rayleigh criterion", Args: "wavelength, diameter", Category: "optics", Formula: "θ = 1.22λ/D"},
	"airydiskradius":     {Fn: FnAiryDiskRadius, Description: "Airy disk radius", Args: "wavelength, focal_length, diameter", Category: "optics", Formula: "r = 1.22λf/D"},

	// EM Spectrum
	"lightwavelength": {Fn: FnLightWavelength, Description: "Light wavelength", Args: "frequency", Category: "optics", Formula: "λ = c/f"},
	"lightfrequency":  {Fn: FnLightFrequency, Description: "Light frequency", Args: "wavelength", Category: "optics", Formula: "f = c/λ"},
	"lightenergy":     {Fn: FnLightEnergy, Description: "Photon energy", Args: "value, type", Category: "optics", Formula: "E = hf or hc/λ"},

	// ════════════════════════════════════════════════════════════════
	// ELECTROMAGNETISM (functions_physics_em.go)
	// ════════════════════════════════════════════════════════════════

	// Electric Force & Field
	"coulombforce":       {Fn: FnCoulombForce, Description: "Coulomb force", Args: "q1, q2, distance", Category: "em", Formula: "F = k_e q₁q₂/r²"},
	"electricfield":      {Fn: FnElectricField, Description: "Electric field", Args: "charge, distance", Category: "em", Formula: "E = k_e q/r²"},
	"electricfieldplate": {Fn: FnElectricFieldPlate, Description: "Electric field (plate)", Args: "sigma, [type]", Category: "em", Formula: "E = σ/ε₀"},
	"electricpotential":  {Fn: FnElectricPotential, Description: "Electric potential", Args: "charge, distance", Category: "em", Formula: "V = k_e q/r"},
	"electricpe":         {Fn: FnElectricPE, Description: "Electric potential energy", Args: "q1, q2, distance", Category: "em", Formula: "U = k_e q₁q₂/r"},
	"electricflux":       {Fn: FnElectricFlux, Description: "Electric flux", Args: "E, A, [angle_deg]", Category: "em", Formula: "Φ = EA cos(θ)"},
	"gausslaw":           {Fn: FnGaussLaw, Description: "Gauss's law", Args: "electric_flux", Category: "em", Formula: "q = ε₀Φ"},

	// Capacitance
	"capacitance":       {Fn: FnCapacitance, Description: "Capacitance", Args: "charge, voltage", Category: "em", Formula: "C = Q/V"},
	"parallelplatecap":  {Fn: FnParallelPlateCapacitance, Description: "Parallel plate capacitance", Args: "area, separation, [epsilon_r]", Category: "em", Formula: "C = ε₀εᵣA/d"},
	"cylindricalcap":    {Fn: FnCylindricalCapacitance, Description: "Cylindrical capacitance", Args: "length, r_inner, r_outer, [epsilon_r]", Category: "em", Formula: "C = 2πε₀εᵣL/ln(b/a)"},
	"sphericalcap":      {Fn: FnSphericalCapacitance, Description: "Spherical capacitance", Args: "r_inner, r_outer, [epsilon_r]", Category: "em", Formula: "C = 4πε₀εᵣab/(b-a)"},
	"capacitorenergy":   {Fn: FnCapacitorEnergy, Description: "Capacitor energy", Args: "capacitance, voltage", Category: "em", Formula: "U = ½CV²"},
	"capacitorseries":   {Fn: FnCapacitorSeries, Description: "Capacitors in series", Args: "C1, C2, ...", Category: "em", Formula: "1/C = Σ(1/Cᵢ)"},
	"capacitorparallel": {Fn: FnCapacitorParallel, Description: "Capacitors in parallel", Args: "C1, C2, ...", Category: "em", Formula: "C = ΣCᵢ"},

	// Resistance
	"ohmslawv":         {Fn: FnOhmsLawV, Description: "Ohm's law (voltage)", Args: "current, resistance", Category: "em", Formula: "V = IR"},
	"ohmslawi":         {Fn: FnOhmsLawI, Description: "Ohm's law (current)", Args: "voltage, resistance", Category: "em", Formula: "I = V/R"},
	"resistivity":      {Fn: FnResistivity, Description: "Resistance from resistivity", Args: "rho, length, area", Category: "em", Formula: "R = ρL/A"},
	"resistorseries":   {Fn: FnResistorSeries, Description: "Resistors in series", Args: "R1, R2, ...", Category: "em", Formula: "R = ΣRᵢ"},
	"resistorparallel": {Fn: FnResistorParallel, Description: "Resistors in parallel", Args: "R1, R2, ...", Category: "em", Formula: "1/R = Σ(1/Rᵢ)"},
	"conductance":      {Fn: FnConductance, Description: "Conductance", Args: "resistance", Category: "em", Formula: "G = 1/R"},
	"currentdensity":   {Fn: FnCurrentDensity, Description: "Current density", Args: "current, area", Category: "em", Formula: "J = I/A"},
	"driftvelocity":    {Fn: FnDriftVelocity, Description: "Drift velocity", Args: "current, carrier_density, area", Category: "em", Formula: "v_d = I/(nAe)"},

	// Electric Power
	"electricpower":       {Fn: FnElectricPower, Description: "Electric power", Args: "voltage, current", Category: "em", Formula: "P = VI"},
	"powerfromresistance": {Fn: FnPowerFromResistance, Description: "Power from I and R", Args: "current, resistance", Category: "em", Formula: "P = I²R"},
	"powerfromvoltage":    {Fn: FnPowerFromVoltage, Description: "Power from V and R", Args: "voltage, resistance", Category: "em", Formula: "P = V²/R"},
	"electricenergy":      {Fn: FnElectricEnergy, Description: "Electric energy", Args: "power, time", Category: "em", Formula: "E = Pt"},

	// RC Circuits
	"rctimeconstant": {Fn: FnRCTimeConstant, Description: "RC time constant", Args: "resistance, capacitance", Category: "em", Formula: "τ = RC"},
	"rccharging":     {Fn: FnRCCharging, Description: "RC charging voltage", Args: "V0, R, C, time", Category: "em", Formula: "V = V₀(1-e^(-t/RC))"},
	"rcdischarging":  {Fn: FnRCDischarging, Description: "RC discharging voltage", Args: "V0, R, C, time", Category: "em", Formula: "V = V₀e^(-t/RC)"},

	// Magnetic Field
	"magneticforcecharge":   {Fn: FnMagneticForceCharge, Description: "Magnetic force on charge", Args: "charge, velocity, B, [angle_deg]", Category: "em", Formula: "F = qvB sin(θ)"},
	"magneticforcewire":     {Fn: FnMagneticForceWire, Description: "Magnetic force on wire", Args: "B, I, length, [angle_deg]", Category: "em", Formula: "F = BIL sin(θ)"},
	"cyclotronradius":       {Fn: FnCyclotronRadius, Description: "Cyclotron radius", Args: "mass, velocity, charge, B", Category: "em", Formula: "r = mv/(qB)"},
	"cyclotronfreq":         {Fn: FnCyclotronFrequency, Description: "Cyclotron frequency", Args: "charge, B, mass", Category: "em", Formula: "f = qB/(2πm)"},
	"magneticfieldwire":     {Fn: FnMagneticFieldWire, Description: "Magnetic field (wire)", Args: "current, distance", Category: "em", Formula: "B = μ₀I/(2πr)"},
	"magneticfieldloop":     {Fn: FnMagneticFieldLoop, Description: "Magnetic field (loop center)", Args: "current, radius", Category: "em", Formula: "B = μ₀I/(2R)"},
	"magneticfieldsolenoid": {Fn: FnMagneticFieldSolenoid, Description: "Magnetic field (solenoid)", Args: "turns, length, current", Category: "em", Formula: "B = μ₀nI"},
	"magneticfieldtoroid":   {Fn: FnMagneticFieldToroid, Description: "Magnetic field (toroid)", Args: "turns, current, radius", Category: "em", Formula: "B = μ₀NI/(2πr)"},

	// Inductance
	"magneticflux":       {Fn: FnMagneticFlux, Description: "Magnetic flux", Args: "B, area, [angle_deg]", Category: "em", Formula: "Φ = BA cos(θ)"},
	"inducedemf":         {Fn: FnInducedEMF, Description: "Induced EMF", Args: "turns, delta_flux, delta_time", Category: "em", Formula: "ε = -NdΦ/dt"},
	"motionalemf":        {Fn: FnMotionalEMF, Description: "Motional EMF", Args: "B, length, velocity", Category: "em", Formula: "ε = BLv"},
	"inductance":         {Fn: FnInductance, Description: "Inductance", Args: "turns, flux, current", Category: "em", Formula: "L = NΦ/I"},
	"solenoidinductance": {Fn: FnSolenoidInductance, Description: "Solenoid inductance", Args: "turns, area, length", Category: "em", Formula: "L = μ₀N²A/L"},
	"inductorenergy":     {Fn: FnInductorEnergy, Description: "Inductor energy", Args: "inductance, current", Category: "em", Formula: "U = ½LI²"},

	// RL Circuits
	"rltimeconstant":  {Fn: FnRLTimeConstant, Description: "RL time constant", Args: "inductance, resistance", Category: "em", Formula: "τ = L/R"},
	"rlcurrentgrowth": {Fn: FnRLCurrentGrowth, Description: "RL current growth", Args: "voltage, R, L, time", Category: "em", Formula: "I = (V/R)(1-e^(-Rt/L))"},
	"rlcurrentdecay":  {Fn: FnRLCurrentDecay, Description: "RL current decay", Args: "I0, R, L, time", Category: "em", Formula: "I = I₀e^(-Rt/L)"},

	// AC Circuits
	"capacitivereactance": {Fn: FnCapacitiveReactance, Description: "Capacitive reactance", Args: "frequency, capacitance", Category: "em", Formula: "X_c = 1/(2πfC)"},
	"inductivereactance":  {Fn: FnInductiveReactance, Description: "Inductive reactance", Args: "frequency, inductance", Category: "em", Formula: "X_L = 2πfL"},
	"impedance":           {Fn: FnImpedance, Description: "Impedance", Args: "R, XL, Xc", Category: "em", Formula: "Z = √(R²+(X_L-X_c)²)"},
	"phaseangle":          {Fn: FnPhaseAngle, Description: "Phase angle", Args: "R, XL, Xc", Category: "em", Formula: "φ = atan((X_L-X_c)/R)"},
	"resonantfreq":        {Fn: FnResonantFrequency, Description: "Resonant frequency", Args: "inductance, capacitance", Category: "em", Formula: "f₀ = 1/(2π√LC)"},
	"qualityfactor":       {Fn: FnQualityFactor, Description: "Q factor", Args: "R, L, C", Category: "em", Formula: "Q = (1/R)√(L/C)"},
	"rmsvoltage":          {Fn: FnRMSVoltage, Description: "RMS voltage", Args: "peak_voltage", Category: "em", Formula: "V_rms = V_peak/√2"},
	"peakvoltage":         {Fn: FnPeakVoltage, Description: "Peak voltage", Args: "rms_voltage", Category: "em", Formula: "V_peak = V_rms×√2"},
	"powerfactor":         {Fn: FnPowerFactor, Description: "Power factor", Args: "resistance, impedance", Category: "em", Formula: "PF = R/Z"},
	"avgpowerac":          {Fn: FnAveragePowerAC, Description: "Average AC power", Args: "V_rms, I_rms, phase_angle_deg", Category: "em", Formula: "P = VI cos(φ)"},

	// Transformers
	"transformervoltage": {Fn: FnTransformerVoltage, Description: "Transformer voltage", Args: "V1, N1, N2", Category: "em", Formula: "V₂ = V₁N₂/N₁"},
	"transformercurrent": {Fn: FnTransformerCurrent, Description: "Transformer current", Args: "I1, N1, N2", Category: "em", Formula: "I₂ = I₁N₁/N₂"},
	"transformerratio":   {Fn: FnTransformerRatio, Description: "Transformer turns ratio", Args: "N1, N2", Category: "em", Formula: "n = N₂/N₁"},
	// ════════════════════════════════════════════════════════════════
	// QUANTUM MECHANICS (functions_physics_quantum.go)
	// ════════════════════════════════════════════════════════════════

	// Photons & Photoelectric
	"photonenergy":        {Fn: FnPhotonEnergy, Description: "Photon energy", Args: "value, type", Category: "quantum", Formula: "E = hf or hc/λ"},
	"photonmomentum":      {Fn: FnPhotonMomentum, Description: "Photon momentum", Args: "value, type", Category: "quantum", Formula: "p = h/λ or E/c"},
	"photoelectricke":     {Fn: FnPhotoelectricKE, Description: "Photoelectric KE", Args: "value, work_function, type", Category: "quantum", Formula: "KE = hf - φ"},
	"workfunction":        {Fn: FnWorkFunction, Description: "Work function", Args: "threshold_frequency", Category: "quantum", Formula: "φ = hf₀"},
	"thresholdfreq":       {Fn: FnThresholdFrequency, Description: "Threshold frequency", Args: "work_function", Category: "quantum", Formula: "f₀ = φ/h"},
	"thresholdwavelength": {Fn: FnThresholdWavelength, Description: "Threshold wavelength", Args: "work_function", Category: "quantum", Formula: "λ₀ = hc/φ"},
	"stoppingpotential":   {Fn: FnStoppingPotential, Description: "Stopping potential", Args: "KE_max", Category: "quantum", Formula: "V_s = KE/e"},

	// De Broglie
	"debrogliewavelength": {Fn: FnDeBroglieWavelength, Description: "De Broglie wavelength", Args: "mass, velocity", Category: "quantum", Formula: "λ = h/(mv)"},
	"debrogliefromke":     {Fn: FnDeBroglieFromKE, Description: "De Broglie from KE", Args: "mass, kinetic_energy", Category: "quantum", Formula: "λ = h/√(2mKE)"},
	"electronwavelength":  {Fn: FnElectronWavelength, Description: "Electron wavelength", Args: "voltage", Category: "quantum", Formula: "λ = h/√(2meV)"},

	// Compton Scattering
	"comptonshift":              {Fn: FnComptonShift, Description: "Compton shift", Args: "angle_deg", Category: "quantum", Formula: "Δλ = λ_c(1-cosθ)"},
	"comptonwavelengthparticle": {Fn: FnComptonWavelengthParticle, Description: "Compton wavelength (particle)", Args: "mass", Category: "quantum", Formula: "λ_c = h/(mc)"},
	"comptonscattered":          {Fn: FnComptonScatteredWavelength, Description: "Scattered wavelength", Args: "wavelength, angle_deg", Category: "quantum", Formula: "λ' = λ + Δλ"},

	// Uncertainty
	"uncertaintyposition": {Fn: FnUncertaintyPosition, Description: "Position uncertainty", Args: "delta_p", Category: "quantum", Formula: "Δx ≥ ℏ/(2Δp)"},
	"uncertaintymomentum": {Fn: FnUncertaintyMomentum, Description: "Momentum uncertainty", Args: "delta_x", Category: "quantum", Formula: "Δp ≥ ℏ/(2Δx)"},
	"uncertaintyenergy":   {Fn: FnUncertaintyEnergy, Description: "Energy uncertainty", Args: "delta_t", Category: "quantum", Formula: "ΔE ≥ ℏ/(2Δt)"},
	"uncertaintytime":     {Fn: FnUncertaintyTime, Description: "Time uncertainty", Args: "delta_E", Category: "quantum", Formula: "Δt ≥ ℏ/(2ΔE)"},

	// Bohr Model
	"bohrradius":        {Fn: FnBohrRadius, Description: "Bohr orbit radius", Args: "n, [Z]", Category: "quantum", Formula: "r = n²a₀/Z"},
	"bohrenergy":        {Fn: FnBohrEnergy, Description: "Bohr energy (J)", Args: "n, [Z]", Category: "quantum", Formula: "E = -13.6Z²/n² eV"},
	"bohrenergyev":      {Fn: FnBohrEnergyEV, Description: "Bohr energy (eV)", Args: "n, [Z]", Category: "quantum", Formula: "E = -13.6Z²/n² eV"},
	"bohrvelocity":      {Fn: FnBohrVelocity, Description: "Bohr orbit velocity", Args: "n, [Z]", Category: "quantum", Formula: "v = αcZ/n"},
	"rydbergwavelength": {Fn: FnRydbergWavelength, Description: "Rydberg wavelength", Args: "n1, n2, [Z]", Category: "quantum", Formula: "1/λ = R_H Z²(1/n₁²-1/n₂²)"},
	"rydbergenergy":     {Fn: FnRydbergEnergy, Description: "Rydberg energy", Args: "n1, n2, [Z]", Category: "quantum", Formula: "ΔE = 13.6Z²(1/n₁²-1/n₂²)"},

	// Particle in Box
	"particleboxenergy":     {Fn: FnParticleBoxEnergy, Description: "Particle in box energy", Args: "n, mass, length", Category: "quantum", Formula: "E = n²h²/(8mL²)"},
	"particleboxwavelength": {Fn: FnParticleBoxWavelength, Description: "Particle in box wavelength", Args: "n, length", Category: "quantum", Formula: "λ = 2L/n"},
	"particlebox3denergy":   {Fn: FnParticleBox3DEnergy, Description: "3D box energy", Args: "nx, ny, nz, mass, Lx, Ly, Lz", Category: "quantum", Formula: "E = (h²/8m)Σ(nᵢ²/Lᵢ²)"},

	// Harmonic Oscillator
	"quantumhoenergy": {Fn: FnQuantumHOEnergy, Description: "QHO energy", Args: "n, omega", Category: "quantum", Formula: "E = (n+½)ℏω"},
	"quantumhofreq":   {Fn: FnQuantumHOFrequency, Description: "QHO frequency", Args: "k, mass", Category: "quantum", Formula: "ω = √(k/m)"},
	"zeropointenergy": {Fn: FnZeroPointEnergy, Description: "Zero-point energy", Args: "omega", Category: "quantum", Formula: "E₀ = ½ℏω"},

	// Tunneling
	"tunnelingprob":  {Fn: FnTunnelingProbability, Description: "Tunneling probability", Args: "mass, V_barrier, E_particle, width", Category: "quantum", Formula: "T ≈ e^(-2κL)"},
	"tunnelingdecay": {Fn: FnTunnelingDecayConstant, Description: "Tunneling decay constant", Args: "mass, V_barrier, E_particle", Category: "quantum", Formula: "κ = √(2m(V-E))/ℏ"},

	// Angular Momentum
	"orbitalangmom":   {Fn: FnOrbitalAngularMomentum, Description: "Orbital angular momentum", Args: "l", Category: "quantum", Formula: "L = ℏ√(l(l+1))"},
	"spinangmom":      {Fn: FnSpinAngularMomentum, Description: "Spin angular momentum", Args: "s", Category: "quantum", Formula: "S = ℏ√(s(s+1))"},
	"magneticmomentq": {Fn: FnMagneticMomentQuantum, Description: "Magnetic moment", Args: "g_factor, quantum_number", Category: "quantum", Formula: "μ = gμ_B√(l(l+1))"},
	"zeemansplitting": {Fn: FnZeemanSplitting, Description: "Zeeman splitting", Args: "g_factor, B, m_l", Category: "quantum", Formula: "ΔE = gμ_B B m_l"},

	// Nuclear
	"nuclearbindingenergy": {Fn: FnNuclearBindingEnergy, Description: "Nuclear binding energy (J)", Args: "A, Z", Category: "quantum", Formula: "Semi-empirical formula"},
	"nuclearbindingmev":    {Fn: FnNuclearBindingEnergyMeV, Description: "Nuclear binding energy (MeV)", Args: "A, Z", Category: "quantum", Formula: "Semi-empirical formula"},
	"massdefect":           {Fn: FnMassDefect, Description: "Mass defect", Args: "Z, N, nuclear_mass", Category: "quantum", Formula: "Δm = Zm_p + Nm_n - M"},
	"massenergy":           {Fn: FnMassEnergy, Description: "Mass-energy", Args: "mass", Category: "quantum", Formula: "E = mc²"},
	"massfromenergy":       {Fn: FnMassFromEnergy, Description: "Mass from energy", Args: "energy", Category: "quantum", Formula: "m = E/c²"},
	"radioactivedecay":     {Fn: FnRadioactiveDecay, Description: "Radioactive decay", Args: "N0, lambda, time", Category: "quantum", Formula: "N = N₀e^(-λt)"},
	"halflife":             {Fn: FnHalfLife, Description: "Half-life", Args: "lambda", Category: "quantum", Formula: "t½ = ln2/λ"},
	"decayconstantfromhl":  {Fn: FnDecayConstantFromHalfLife, Description: "Decay constant", Args: "half_life", Category: "quantum", Formula: "λ = ln2/t½"},
	"activity":             {Fn: FnActivity, Description: "Radioactive activity", Args: "lambda, N", Category: "quantum", Formula: "A = λN"},

	// ════════════════════════════════════════════════════════════════
	// RELATIVITY (functions_physics_relativity.go)
	// ════════════════════════════════════════════════════════════════

	// Lorentz Factor
	"lorentzfactor":       {Fn: FnLorentzFactor, Description: "Lorentz factor", Args: "velocity", Category: "relativity", Formula: "γ = 1/√(1-v²/c²)"},
	"lorentzfrombeta":     {Fn: FnLorentzFactorFromBeta, Description: "Lorentz factor from β", Args: "beta", Category: "relativity", Formula: "γ = 1/√(1-β²)"},
	"betafromlorentz":     {Fn: FnBetaFromLorentz, Description: "β from γ", Args: "gamma", Category: "relativity", Formula: "β = √(1-1/γ²)"},
	"velocityfromlorentz": {Fn: FnVelocityFromLorentz, Description: "Velocity from γ", Args: "gamma", Category: "relativity", Formula: "v = c√(1-1/γ²)"},

	// Time Dilation
	"timedilation":     {Fn: FnTimeDilation, Description: "Time dilation", Args: "proper_time, velocity", Category: "relativity", Formula: "Δt = γΔt₀"},
	"propertime":       {Fn: FnProperTime, Description: "Proper time", Args: "dilated_time, velocity", Category: "relativity", Formula: "Δt₀ = Δt/γ"},
	"timedilationgrav": {Fn: FnTimeDilationGravity, Description: "Gravitational time dilation", Args: "proper_time, mass, radius", Category: "relativity", Formula: "Δt = Δt₀/√(1-r_s/r)"},

	// Length Contraction
	"lengthcontraction": {Fn: FnLengthContraction, Description: "Length contraction", Args: "proper_length, velocity", Category: "relativity", Formula: "L = L₀/γ"},
	"properlength":      {Fn: FnProperLength, Description: "Proper length", Args: "contracted_length, velocity", Category: "relativity", Formula: "L₀ = γL"},

	// Velocity Addition
	"relvelocityadd":     {Fn: FnRelativisticVelocityAdd, Description: "Relativistic velocity addition", Args: "v, u_prime", Category: "relativity", Formula: "u = (v+u')/(1+vu'/c²)"},
	"relvelocityaddbeta": {Fn: FnRelativisticVelocityAddBeta, Description: "Velocity addition (β)", Args: "beta1, beta2", Category: "relativity", Formula: "β = (β₁+β₂)/(1+β₁β₂)"},

	// Momentum & Energy
	"relmomentum":        {Fn: FnRelativisticMomentum, Description: "Relativistic momentum", Args: "mass, velocity", Category: "relativity", Formula: "p = γmv"},
	"relenergy":          {Fn: FnRelativisticEnergy, Description: "Relativistic energy", Args: "mass, velocity", Category: "relativity", Formula: "E = γmc²"},
	"restenergy":         {Fn: FnRestEnergy, Description: "Rest energy", Args: "mass", Category: "relativity", Formula: "E₀ = mc²"},
	"relke":              {Fn: FnRelativisticKE, Description: "Relativistic KE", Args: "mass, velocity", Category: "relativity", Formula: "KE = (γ-1)mc²"},
	"energymomentum":     {Fn: FnEnergyMomentumRelation, Description: "Energy-momentum relation", Args: "momentum, mass", Category: "relativity", Formula: "E = √((pc)²+(mc²)²)"},
	"momentumfromenergy": {Fn: FnMomentumFromEnergy, Description: "Momentum from energy", Args: "energy, mass", Category: "relativity", Formula: "p = √(E²-(mc²)²)/c"},
	"velocityfrommom":    {Fn: FnVelocityFromMomentum, Description: "Velocity from momentum", Args: "momentum, mass", Category: "relativity", Formula: "v = pc²/E"},
	"relmass":            {Fn: FnRelativisticMass, Description: "Relativistic mass", Args: "rest_mass, velocity", Category: "relativity", Formula: "m_rel = γm₀"},

	// Doppler Effect
	"reldopplerapproach":   {Fn: FnRelativisticDopplerApproach, Description: "Relativistic Doppler (approach)", Args: "frequency, velocity", Category: "relativity", Formula: "f' = f√((1+β)/(1-β))"},
	"reldopplerrecede":     {Fn: FnRelativisticDopplerRecede, Description: "Relativistic Doppler (recede)", Args: "frequency, velocity", Category: "relativity", Formula: "f' = f√((1-β)/(1+β))"},
	"redshift":             {Fn: FnRedshift, Description: "Redshift", Args: "velocity", Category: "relativity", Formula: "z = √((1+β)/(1-β))-1"},
	"velocityfromredshift": {Fn: FnVelocityFromRedshift, Description: "Velocity from redshift", Args: "z", Category: "relativity", Formula: "β = ((z+1)²-1)/((z+1)²+1)"},
	"transversedoppler":    {Fn: FnTransverseDoppler, Description: "Transverse Doppler", Args: "frequency, velocity", Category: "relativity", Formula: "f' = f/γ"},

	// Lorentz Transformations
	"lorentztransformx": {Fn: FnLorentzTransformX, Description: "Lorentz transform x", Args: "x, t, velocity", Category: "relativity", Formula: "x' = γ(x-vt)"},
	"lorentztransformt": {Fn: FnLorentzTransformT, Description: "Lorentz transform t", Args: "x, t, velocity", Category: "relativity", Formula: "t' = γ(t-vx/c²)"},
	"inverselorentzx":   {Fn: FnInverseLorentzX, Description: "Inverse Lorentz x", Args: "x_prime, t_prime, velocity", Category: "relativity", Formula: "x = γ(x'+vt')"},
	"inverselorentzt":   {Fn: FnInverseLorentzT, Description: "Inverse Lorentz t", Args: "x_prime, t_prime, velocity", Category: "relativity", Formula: "t = γ(t'+vx'/c²)"},

	// Spacetime Interval
	"spacetimeinterval":    {Fn: FnSpacetimeInterval, Description: "Spacetime interval squared", Args: "dt, dx, [dy], [dz]", Category: "relativity", Formula: "(Δs)² = (cΔt)²-Δr²"},
	"propertimeinterval":   {Fn: FnProperTimeInterval, Description: "Proper time from interval", Args: "ds_squared", Category: "relativity", Formula: "τ = √(ds²)/c"},
	"properlengthinterval": {Fn: FnProperLengthInterval, Description: "Proper length from interval", Args: "ds_squared", Category: "relativity", Formula: "L = √(-ds²)"},
	"intervaltype":         {Fn: FnIntervalType, Description: "Interval type", Args: "ds_squared", Category: "relativity", Formula: "1=timelike, 0=null, -1=spacelike"},

	// Four-Vectors
	"fourmomentumt":       {Fn: FnFourMomentumE, Description: "Four-momentum E/c", Args: "mass, velocity", Category: "relativity", Formula: "p⁰ = γmc"},
	"fourmomentumspatial": {Fn: FnFourMomentumP, Description: "Four-momentum spatial", Args: "mass, velocity", Category: "relativity", Formula: "p = γmv"},
	"fourmomentummag":     {Fn: FnFourMomentumMagnitude, Description: "Invariant mass", Args: "energy, momentum", Category: "relativity", Formula: "m = √((E/c)²-p²)/c"},

	// Collisions
	"centerofmassenergy": {Fn: FnCenterOfMassEnergy, Description: "Center-of-mass energy", Args: "E1, m1, m2", Category: "relativity", Formula: "E_cm"},
	"thresholdenergy":    {Fn: FnThresholdEnergy, Description: "Threshold energy", Args: "m1, m2, M_total", Category: "relativity", Formula: "E_threshold"},

	// Acceleration
	"properaccel":         {Fn: FnProperAcceleration, Description: "Proper acceleration", Args: "coordinate_accel, velocity", Category: "relativity", Formula: "a₀ = γ³a"},
	"coordaccel":          {Fn: FnCoordinateAcceleration, Description: "Coordinate acceleration", Args: "proper_accel, velocity", Category: "relativity", Formula: "a = a₀/γ³"},
	"relrocket":           {Fn: FnRelativisticRocket, Description: "Relativistic rocket velocity", Args: "proper_accel, time", Category: "relativity", Formula: "v = a₀t/√(1+(a₀t/c)²)"},
	"relrocketdist":       {Fn: FnRelativisticRocketDistance, Description: "Relativistic rocket distance", Args: "proper_accel, time", Category: "relativity", Formula: "x = (c²/a₀)(√(1+(a₀t/c)²)-1)"},
	"relrocketpropertime": {Fn: FnRelativisticRocketProperTime, Description: "Relativistic rocket proper time", Args: "proper_accel, coord_time", Category: "relativity", Formula: "τ = (c/a₀)arcsinh(a₀t/c)"},
}

// GetPhysicsFunction retrieves a physics function by name.
func GetPhysicsFunction(name string) (func([]types.Value) types.Value, bool) {
	info, ok := PhysicsFunctionRegistry[name]
	if !ok {
		return nil, false
	}
	return info.Fn, true
}

// GetPhysicsFunctionInfo retrieves full function info by name.
func GetPhysicsFunctionInfo(name string) (PhysicsFunctionInfo, bool) {
	info, ok := PhysicsFunctionRegistry[name]
	return info, ok
}

// GetPhysicsFunctionsByCategory returns all functions in a category.
func GetPhysicsFunctionsByCategory(category string) map[string]PhysicsFunctionInfo {
	result := make(map[string]PhysicsFunctionInfo)
	for name, info := range PhysicsFunctionRegistry {
		if info.Category == category {
			result[name] = info
		}
	}
	return result
}

// GetAllPhysicsCategories returns all available categories.
func GetAllPhysicsCategories() []string {
	categorySet := make(map[string]bool)
	for _, info := range PhysicsFunctionRegistry {
		categorySet[info.Category] = true
	}
	categories := make([]string, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, cat)
	}
	return categories
}

// PhysicsCategories provides organized category descriptions.
var PhysicsCategories = map[string]string{
	"constants":  "Physical constants (c, h, G, e, etc.)",
	"mechanics":  "Classical mechanics (force, energy, momentum, rotation)",
	"gravity":    "Gravitation and orbital mechanics",
	"thermo":     "Thermodynamics and heat transfer",
	"waves":      "Waves, sound, and acoustics",
	"optics":     "Optics, refraction, lenses, diffraction",
	"em":         "Electromagnetism (electric, magnetic, circuits)",
	"quantum":    "Quantum mechanics and atomic physics",
	"relativity": "Special and general relativity",
}
