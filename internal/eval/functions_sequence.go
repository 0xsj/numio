// internal/eval/functions_sequence.go

package eval

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// FIBONACCI
// ════════════════════════════════════════════════════════════════

// FnFib returns the nth Fibonacci number.
// Args: n (0-indexed)
func FnFib(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("fib requires exactly 1 argument: n")
	}

	n := int(args[0].AsFloat())
	if n < 0 {
		return types.Error("fib: n must be non-negative")
	}
	if n > 1000 {
		return types.Error("fib: n too large (max 1000)")
	}

	return types.Number(fibonacci(n))
}

// FnFibSeq returns a sequence of Fibonacci numbers.
// Args: count, [start index]
func FnFibSeq(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("fibseq requires 1-2 arguments: count, [start]")
	}

	count := int(args[0].AsFloat())
	start := 0
	if len(args) == 2 {
		start = int(args[1].AsFloat())
	}

	if count < 1 {
		return types.StringValue("[]")
	}
	if count > 100 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	results := make([]string, count)
	for i := 0; i < count; i++ {
		results[i] = formatNumber(fibonacci(start + i))
	}

	return types.StringValue("[" + strings.Join(results, ", ") + "]")
}

// FnIsFib checks if a number is a Fibonacci number.
func FnIsFib(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isfib requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	if n < 0 || n != math.Floor(n) {
		return types.Number(0)
	}

	// A number is Fibonacci if 5n² + 4 or 5n² - 4 is a perfect square
	n2 := n * n
	if isPerfectSquare(5*n2+4) || isPerfectSquare(5*n2-4) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnFibIndex returns the index of a Fibonacci number (or -1 if not Fibonacci).
func FnFibIndex(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("fibindex requires exactly 1 argument")
	}

	target := args[0].AsFloat()
	if target < 0 || target != math.Floor(target) {
		return types.Number(-1)
	}

	// Search for the index
	a, b := 0.0, 1.0
	index := 0
	for a <= target {
		if a == target {
			return types.Number(float64(index))
		}
		a, b = b, a+b
		index++
		if index > 1000 {
			break
		}
	}

	return types.Number(-1)
}

// fibonacci calculates the nth Fibonacci number using fast doubling.
func fibonacci(n int) float64 {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	// Use iterative approach for simplicity
	a, b := 0.0, 1.0
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// ════════════════════════════════════════════════════════════════
// PRIME NUMBERS
// ════════════════════════════════════════════════════════════════

// FnPrime returns the nth prime number (1-indexed).
func FnPrime(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("prime requires exactly 1 argument: n")
	}

	n := int(args[0].AsFloat())
	if n < 1 {
		return types.Error("prime: n must be positive")
	}
	if n > 10000 {
		return types.Error("prime: n too large (max 10000)")
	}

	return types.Number(float64(nthPrime(n)))
}

// FnIsPrime checks if a number is prime.
func FnIsPrime(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isprime requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	if n < 2 || n != math.Floor(n) {
		return types.Number(0)
	}

	if isPrime(int64(n)) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnNextPrime returns the next prime >= n.
func FnNextPrime(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("nextprime requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 2 {
		return types.Number(2)
	}

	for !isPrime(n) {
		n++
		if n > 1e12 {
			return types.Error("nextprime: value too large")
		}
	}

	return types.Number(float64(n))
}

// FnPrevPrime returns the previous prime <= n.
func FnPrevPrime(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("prevprime requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 2 {
		return types.Error("prevprime: no prime <= n")
	}

	for !isPrime(n) {
		n--
		if n < 2 {
			return types.Error("prevprime: no prime <= n")
		}
	}

	return types.Number(float64(n))
}

// FnPrimeSeq returns a sequence of primes.
// Args: count, [start value]
func FnPrimeSeq(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("primeseq requires 1-2 arguments: count, [start]")
	}

	count := int(args[0].AsFloat())
	start := int64(2)
	if len(args) == 2 {
		start = int64(args[1].AsFloat())
		if start < 2 {
			start = 2
		}
	}

	if count < 1 {
		return types.StringValue("[]")
	}
	if count > 100 {
		count = 100
	}

	primes := make([]string, 0, count)
	n := start
	for len(primes) < count {
		if isPrime(n) {
			primes = append(primes, formatNumber(float64(n)))
		}
		n++
		if n > 1e10 {
			break
		}
	}

	return types.StringValue("[" + strings.Join(primes, ", ") + "]")
}

// FnPrimesBelow returns all primes below n.
func FnPrimesBelow(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("primesbelow requires exactly 1 argument")
	}

	limit := int(args[0].AsFloat())
	if limit < 2 {
		return types.StringValue("[]")
	}
	if limit > 100000 {
		return types.Error("primesbelow: limit too large (max 100000)")
	}

	primes := sieveOfEratosthenes(limit)
	results := make([]string, len(primes))
	for i, p := range primes {
		results[i] = fmt.Sprintf("%d", p)
	}

	return types.StringValue("[" + strings.Join(results, ", ") + "]")
}

// FnPrimeCount counts primes up to n (prime counting function π(n)).
func FnPrimeCount(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("primecount requires exactly 1 argument")
	}

	n := int(args[0].AsFloat())
	if n < 2 {
		return types.Number(0)
	}
	if n > 10000000 {
		return types.Error("primecount: n too large (max 10000000)")
	}

	primes := sieveOfEratosthenes(n)
	return types.Number(float64(len(primes)))
}

// FnPrimePi approximates the prime counting function using Li(x).
func FnPrimePi(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("primepi requires exactly 1 argument")
	}

	x := args[0].AsFloat()
	if x < 2 {
		return types.Number(0)
	}

	// Use logarithmic integral approximation: π(x) ≈ x / ln(x)
	approx := x / math.Log(x)
	return types.Number(math.Round(approx))
}

// isPrime checks if n is prime using trial division and Miller-Rabin for large n.
func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	// Trial division for small numbers
	if n < 1000000 {
		for i := int64(5); i*i <= n; i += 6 {
			if n%i == 0 || n%(i+2) == 0 {
				return false
			}
		}
		return true
	}

	// Miller-Rabin for larger numbers
	return millerRabin(n)
}

// millerRabin performs the Miller-Rabin primality test.
func millerRabin(n int64) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	// Write n-1 as 2^r * d
	d := n - 1
	r := 0
	for d%2 == 0 {
		d /= 2
		r++
	}

	// Witnesses to test (sufficient for n < 3,317,044,064,679,887,385,961,981)
	witnesses := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}

	for _, a := range witnesses {
		if a >= n {
			continue
		}
		if !millerRabinWitness(a, d, n, r) {
			return false
		}
	}
	return true
}

func millerRabinWitness(a, d, n int64, r int) bool {
	x := modPow(a, d, n)
	if x == 1 || x == n-1 {
		return true
	}
	for i := 0; i < r-1; i++ {
		x = modPow(x, 2, n)
		if x == n-1 {
			return true
		}
	}
	return false
}

func modPow(base, exp, mod int64) int64 {
	result := int64(1)
	base = base % mod
	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) % mod
		}
		exp /= 2
		base = (base * base) % mod
	}
	return result
}

// nthPrime returns the nth prime number.
func nthPrime(n int) int {
	if n == 1 {
		return 2
	}

	count := 1
	candidate := 1
	for count < n {
		candidate += 2
		if isPrime(int64(candidate)) {
			count++
		}
	}
	return candidate
}

// sieveOfEratosthenes returns all primes up to limit.
func sieveOfEratosthenes(limit int) []int {
	if limit < 2 {
		return []int{}
	}

	sieve := make([]bool, limit+1)
	for i := range sieve {
		sieve[i] = true
	}
	sieve[0], sieve[1] = false, false

	for i := 2; i*i <= limit; i++ {
		if sieve[i] {
			for j := i * i; j <= limit; j += i {
				sieve[j] = false
			}
		}
	}

	primes := make([]int, 0)
	for i, isPrime := range sieve {
		if isPrime {
			primes = append(primes, i)
		}
	}
	return primes
}

// ════════════════════════════════════════════════════════════════
// FACTORIZATION
// ════════════════════════════════════════════════════════════════

// FnFactors returns the prime factorization of a number.
func FnFactors(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("factors requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("factors: n must be positive")
	}
	if n == 1 {
		return types.StringValue("[]")
	}

	factors := primeFactors(n)
	results := make([]string, len(factors))
	for i, f := range factors {
		results[i] = fmt.Sprintf("%d", f)
	}

	return types.StringValue("[" + strings.Join(results, ", ") + "]")
}

// FnFactorization returns prime factorization with exponents.
func FnFactorization(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("factorization requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("factorization: n must be positive")
	}
	if n == 1 {
		return types.StringValue("1")
	}

	factors := primeFactors(n)

	// Count occurrences
	counts := make(map[int64]int)
	for _, f := range factors {
		counts[f]++
	}

	// Sort primes
	primes := make([]int64, 0, len(counts))
	for p := range counts {
		primes = append(primes, p)
	}
	sort.Slice(primes, func(i, j int) bool { return primes[i] < primes[j] })

	// Format
	parts := make([]string, 0, len(primes))
	for _, p := range primes {
		exp := counts[p]
		if exp == 1 {
			parts = append(parts, fmt.Sprintf("%d", p))
		} else {
			parts = append(parts, fmt.Sprintf("%d^%d", p, exp))
		}
	}

	return types.StringValue(strings.Join(parts, " × "))
}

// FnDivisors returns all divisors of a number.
func FnDivisors(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("divisors requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("divisors: n must be positive")
	}

	divs := divisors(n)
	results := make([]string, len(divs))
	for i, d := range divs {
		results[i] = fmt.Sprintf("%d", d)
	}

	return types.StringValue("[" + strings.Join(results, ", ") + "]")
}

// FnDivisorCount returns the number of divisors.
func FnDivisorCount(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("divisorcount requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("divisorcount: n must be positive")
	}

	divs := divisors(n)
	return types.Number(float64(len(divs)))
}

// FnDivisorSum returns the sum of divisors (σ(n)).
func FnDivisorSum(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("divisorsum requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("divisorsum: n must be positive")
	}

	divs := divisors(n)
	sum := int64(0)
	for _, d := range divs {
		sum += d
	}

	return types.Number(float64(sum))
}

// FnProperDivisorSum returns the sum of proper divisors (excluding n).
func FnProperDivisorSum(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("properdivisorsum requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("properdivisorsum: n must be positive")
	}

	divs := divisors(n)
	sum := int64(0)
	for _, d := range divs {
		if d != n {
			sum += d
		}
	}

	return types.Number(float64(sum))
}

// primeFactors returns the prime factors of n (with repetition).
func primeFactors(n int64) []int64 {
	factors := make([]int64, 0)

	// Factor out 2s
	for n%2 == 0 {
		factors = append(factors, 2)
		n /= 2
	}

	// Factor out odd numbers
	for i := int64(3); i*i <= n; i += 2 {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}

	// If n is still > 1, it's prime
	if n > 1 {
		factors = append(factors, n)
	}

	return factors
}

// divisors returns all divisors of n in sorted order.
func divisors(n int64) []int64 {
	divs := make([]int64, 0)

	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			divs = append(divs, i)
			if i != n/i {
				divs = append(divs, n/i)
			}
		}
	}

	sort.Slice(divs, func(i, j int) bool { return divs[i] < divs[j] })
	return divs
}

// ════════════════════════════════════════════════════════════════
// SPECIAL NUMBER CLASSIFICATIONS
// ════════════════════════════════════════════════════════════════

// FnIsPerfect checks if a number is perfect (equals sum of proper divisors).
func FnIsPerfect(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isperfect requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Number(0)
	}

	divs := divisors(n)
	sum := int64(0)
	for _, d := range divs {
		if d != n {
			sum += d
		}
	}

	if sum == n {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsAbundant checks if a number is abundant (proper divisor sum > n).
func FnIsAbundant(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isabundant requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Number(0)
	}

	divs := divisors(n)
	sum := int64(0)
	for _, d := range divs {
		if d != n {
			sum += d
		}
	}

	if sum > n {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsDeficient checks if a number is deficient (proper divisor sum < n).
func FnIsDeficient(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isdeficient requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Number(0)
	}

	divs := divisors(n)
	sum := int64(0)
	for _, d := range divs {
		if d != n {
			sum += d
		}
	}

	if sum < n {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsSquare checks if a number is a perfect square.
func FnIsSquare(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("issquare requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	if n < 0 {
		return types.Number(0)
	}

	if isPerfectSquare(n) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsCube checks if a number is a perfect cube.
func FnIsCube(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("iscube requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	root := math.Cbrt(n)
	rounded := math.Round(root)

	if math.Abs(rounded*rounded*rounded-n) < 1e-9 {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsPowerOf checks if a number is a power of another.
func FnIsPowerOf(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("ispowerof requires 2 arguments: n, base")
	}

	n := args[0].AsFloat()
	base := args[1].AsFloat()

	if base <= 0 || base == 1 || n <= 0 {
		return types.Number(0)
	}

	// Check if log_base(n) is an integer
	logVal := math.Log(n) / math.Log(base)
	rounded := math.Round(logVal)

	if math.Abs(logVal-rounded) < 1e-9 && rounded >= 0 {
		// Verify by computing base^rounded
		if math.Abs(math.Pow(base, rounded)-n) < 1e-9 {
			return types.Number(1)
		}
	}
	return types.Number(0)
}

// isPerfectSquare checks if n is a perfect square.
func isPerfectSquare(n float64) bool {
	if n < 0 {
		return false
	}
	root := math.Sqrt(n)
	rounded := math.Round(root)
	return math.Abs(rounded*rounded-n) < 1e-9
}

// ════════════════════════════════════════════════════════════════
// OTHER SEQUENCES
// ════════════════════════════════════════════════════════════════

// FnTriangular returns the nth triangular number.
func FnTriangular(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("triangular requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	if n < 0 {
		return types.Error("triangular: n must be non-negative")
	}

	// T(n) = n(n+1)/2
	return types.Number(n * (n + 1) / 2)
}

// FnIsTriangular checks if a number is triangular.
func FnIsTriangular(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("istriangular requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	if n < 0 {
		return types.Number(0)
	}

	// n is triangular if 8n+1 is a perfect square
	if isPerfectSquare(8*n + 1) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnSquareNum returns the nth square number.
func FnSquareNum(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("squarenum requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	if n < 0 {
		return types.Error("squarenum: n must be non-negative")
	}

	return types.Number(n * n)
}

// FnCubeNum returns the nth cube number.
func FnCubeNum(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("cubenum requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	return types.Number(n * n * n)
}

// FnPentagonal returns the nth pentagonal number.
func FnPentagonal(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("pentagonal requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	if n < 0 {
		return types.Error("pentagonal: n must be non-negative")
	}

	// P(n) = n(3n-1)/2
	return types.Number(n * (3*n - 1) / 2)
}

// FnHexagonal returns the nth hexagonal number.
func FnHexagonal(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hexagonal requires exactly 1 argument")
	}

	n := args[0].AsFloat()
	if n < 0 {
		return types.Error("hexagonal: n must be non-negative")
	}

	// H(n) = n(2n-1)
	return types.Number(n * (2*n - 1))
}

// FnCatalan returns the nth Catalan number.
func FnCatalan(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("catalan requires exactly 1 argument")
	}

	n := int(args[0].AsFloat())
	if n < 0 {
		return types.Error("catalan: n must be non-negative")
	}
	if n > 30 {
		return types.Error("catalan: n too large (max 30)")
	}

	// C(n) = C(2n, n) / (n+1)
	result := binomialCoeff(2*n, n) / float64(n+1)
	return types.Number(result)
}

// FnLucas returns the nth Lucas number.
func FnLucas(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("lucas requires exactly 1 argument")
	}

	n := int(args[0].AsFloat())
	if n < 0 {
		return types.Error("lucas: n must be non-negative")
	}
	if n > 1000 {
		return types.Error("lucas: n too large (max 1000)")
	}

	// L(0) = 2, L(1) = 1, L(n) = L(n-1) + L(n-2)
	if n == 0 {
		return types.Number(2)
	}
	if n == 1 {
		return types.Number(1)
	}

	a, b := 2.0, 1.0
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return types.Number(b)
}

// FnCollatz returns the Collatz sequence starting from n.
func FnCollatz(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("collatz requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("collatz: n must be positive")
	}

	sequence := make([]string, 0)
	maxSteps := 1000

	for n != 1 && len(sequence) < maxSteps {
		sequence = append(sequence, fmt.Sprintf("%d", n))
		if n%2 == 0 {
			n = n / 2
		} else {
			n = 3*n + 1
		}
	}
	sequence = append(sequence, "1")

	return types.StringValue("[" + strings.Join(sequence, ", ") + "]")
}

// FnCollatzLength returns the length of the Collatz sequence.
func FnCollatzLen(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("collatzlen requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("collatzlen: n must be positive")
	}

	steps := 0
	maxSteps := 10000

	for n != 1 && steps < maxSteps {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = 3*n + 1
		}
		steps++
	}

	return types.Number(float64(steps))
}

// ════════════════════════════════════════════════════════════════
// NUMBER THEORY UTILITIES
// ════════════════════════════════════════════════════════════════

// FnTotient returns Euler's totient function φ(n).
func FnTotient(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("totient requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("totient: n must be positive")
	}

	result := n
	temp := n

	// Factor out 2
	if temp%2 == 0 {
		result -= result / 2
		for temp%2 == 0 {
			temp /= 2
		}
	}

	// Factor out odd primes
	for i := int64(3); i*i <= temp; i += 2 {
		if temp%i == 0 {
			result -= result / i
			for temp%i == 0 {
				temp /= i
			}
		}
	}

	// If temp is still > 1, it's a prime factor
	if temp > 1 {
		result -= result / temp
	}

	return types.Number(float64(result))
}

// FnMobius returns the Möbius function μ(n).
func FnMobius(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("mobius requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("mobius: n must be positive")
	}
	if n == 1 {
		return types.Number(1)
	}

	factors := primeFactors(n)

	// Check for square factors
	counts := make(map[int64]int)
	for _, f := range factors {
		counts[f]++
		if counts[f] > 1 {
			return types.Number(0) // Has a square factor
		}
	}

	// μ(n) = (-1)^k where k is the number of distinct prime factors
	if len(factors)%2 == 0 {
		return types.Number(1)
	}
	return types.Number(-1)
}

// FnRadical returns the radical of n (product of distinct prime factors).
func FnRadical(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("radical requires exactly 1 argument")
	}

	n := int64(args[0].AsFloat())
	if n < 1 {
		return types.Error("radical: n must be positive")
	}
	if n == 1 {
		return types.Number(1)
	}

	factors := primeFactors(n)

	// Get distinct factors
	seen := make(map[int64]bool)
	radical := int64(1)
	for _, f := range factors {
		if !seen[f] {
			seen[f] = true
			radical *= f
		}
	}

	return types.Number(float64(radical))
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

// formatNumber formats a float for display in sequences.
func formatNumber(n float64) string {
	if n == math.Trunc(n) {
		return fmt.Sprintf("%.0f", n)
	}
	return fmt.Sprintf("%g", n)
}
