// internal/fuzzy/levenshtein.go

package fuzzy

// ════════════════════════════════════════════════════════════════
// LEVENSHTEIN MATCHER
// ════════════════════════════════════════════════════════════════

// Levenshtein implements the classic Levenshtein edit distance algorithm.
// It counts the minimum number of single-character edits (insertions,
// deletions, or substitutions) required to transform one string into another.
//
// Time Complexity: O(m*n) where m and n are string lengths
// Space Complexity: O(min(m,n)) using optimized single-row approach
type Levenshtein struct {
	// CaseSensitive controls case sensitivity.
	CaseSensitive bool

	// MaxDistance is the threshold for considering strings a match.
	MaxDistance int
}

// NewLevenshtein creates a new Levenshtein matcher with default settings.
func NewLevenshtein() *Levenshtein {
	return &Levenshtein{
		CaseSensitive: false,
		MaxDistance:   2,
	}
}

// Name returns the algorithm name.
func (l *Levenshtein) Name() string {
	return "levenshtein"
}

// Distance calculates the Levenshtein distance between two strings.
// Returns the minimum number of edits to transform a into b.
func (l *Levenshtein) Distance(a, b string) int {
	a = normalize(a, l.CaseSensitive)
	b = normalize(b, l.CaseSensitive)

	return levenshteinDistance(a, b)
}

// Similarity returns a similarity score between 0.0 and 1.0.
// Calculated as: 1 - (distance / max(len(a), len(b)))
func (l *Levenshtein) Similarity(a, b string) float64 {
	a = normalize(a, l.CaseSensitive)
	b = normalize(b, l.CaseSensitive)

	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}

	dist := levenshteinDistance(a, b)
	maxLen := max(len(a), len(b))

	return 1.0 - float64(dist)/float64(maxLen)
}

// Match returns true if the distance is within MaxDistance.
func (l *Levenshtein) Match(a, b string) bool {
	return l.Distance(a, b) <= l.MaxDistance
}

// ════════════════════════════════════════════════════════════════
// CORE ALGORITHM
// ════════════════════════════════════════════════════════════════

// levenshteinDistance computes the edit distance using Wagner-Fischer algorithm.
// Uses space-optimized approach with single row.
func levenshteinDistance(a, b string) int {
	// Ensure a is the shorter string for space optimization
	if len(a) > len(b) {
		a, b = b, a
	}

	m := len(a)
	n := len(b)

	// Handle edge cases
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Single row for space optimization
	// prev[j] represents distance for a[0:i-1] to b[0:j]
	prev := make([]int, m+1)
	curr := make([]int, m+1)

	// Initialize first row
	for j := 0; j <= m; j++ {
		prev[j] = j
	}

	// Fill the matrix row by row
	for i := 1; i <= n; i++ {
		curr[0] = i

		for j := 1; j <= m; j++ {
			cost := 1
			if b[i-1] == a[j-1] {
				cost = 0
			}

			curr[j] = min3(
				prev[j]+1,      // deletion
				curr[j-1]+1,    // insertion
				prev[j-1]+cost, // substitution
			)
		}

		// Swap rows
		prev, curr = curr, prev
	}

	return prev[m]
}

// ════════════════════════════════════════════════════════════════
// BOUNDED DISTANCE (EARLY TERMINATION)
// ════════════════════════════════════════════════════════════════

// LevenshteinBounded computes distance with early termination.
// Returns the distance if <= maxDist, otherwise returns maxDist+1.
// This is more efficient when you only care about close matches.
func LevenshteinBounded(a, b string, maxDist int) int {
	// Ensure a is the shorter string
	if len(a) > len(b) {
		a, b = b, a
	}

	m := len(a)
	n := len(b)

	// Quick length check
	if n-m > maxDist {
		return maxDist + 1
	}

	// Handle edge cases
	if m == 0 {
		if n <= maxDist {
			return n
		}
		return maxDist + 1
	}

	prev := make([]int, m+1)
	curr := make([]int, m+1)

	// Initialize first row
	for j := 0; j <= m; j++ {
		prev[j] = j
	}

	for i := 1; i <= n; i++ {
		curr[0] = i
		minInRow := curr[0]

		for j := 1; j <= m; j++ {
			cost := 1
			if b[i-1] == a[j-1] {
				cost = 0
			}

			curr[j] = min3(
				prev[j]+1,
				curr[j-1]+1,
				prev[j-1]+cost,
			)

			if curr[j] < minInRow {
				minInRow = curr[j]
			}
		}

		// Early termination: if minimum in row exceeds maxDist,
		// the final result will too
		if minInRow > maxDist {
			return maxDist + 1
		}

		prev, curr = curr, prev
	}

	if prev[m] <= maxDist {
		return prev[m]
	}
	return maxDist + 1
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// EditDistance is a convenience function for quick Levenshtein distance.
func EditDistance(a, b string) int {
	return levenshteinDistance(toLower(a), toLower(b))
}

// IsCloseMatch checks if two strings are within n edits.
func IsCloseMatch(a, b string, maxEdits int) bool {
	return LevenshteinBounded(toLower(a), toLower(b), maxEdits) <= maxEdits
}

// ════════════════════════════════════════════════════════════════
// WEIGHTED LEVENSHTEIN
// ════════════════════════════════════════════════════════════════

// WeightedLevenshtein allows custom costs for different operations.
type WeightedLevenshtein struct {
	CaseSensitive bool
	MaxDistance   int

	// Costs for different operations
	InsertCost     int
	DeleteCost     int
	SubstituteCost int
}

// NewWeightedLevenshtein creates a weighted matcher with default costs.
func NewWeightedLevenshtein() *WeightedLevenshtein {
	return &WeightedLevenshtein{
		CaseSensitive:  false,
		MaxDistance:    2,
		InsertCost:     1,
		DeleteCost:     1,
		SubstituteCost: 1,
	}
}

// Name returns the algorithm name.
func (w *WeightedLevenshtein) Name() string {
	return "weighted-levenshtein"
}

// Distance calculates the weighted edit distance.
func (w *WeightedLevenshtein) Distance(a, b string) int {
	a = normalize(a, w.CaseSensitive)
	b = normalize(b, w.CaseSensitive)

	return w.weightedDistance(a, b)
}

// Similarity returns a similarity score between 0.0 and 1.0.
func (w *WeightedLevenshtein) Similarity(a, b string) float64 {
	a = normalize(a, w.CaseSensitive)
	b = normalize(b, w.CaseSensitive)

	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}

	dist := w.weightedDistance(a, b)
	maxLen := max(len(a), len(b))
	maxCost := max(max(w.InsertCost, w.DeleteCost), w.SubstituteCost)

	return 1.0 - float64(dist)/float64(maxLen*maxCost)
}

// Match returns true if the distance is within MaxDistance.
func (w *WeightedLevenshtein) Match(a, b string) bool {
	return w.Distance(a, b) <= w.MaxDistance
}

func (w *WeightedLevenshtein) weightedDistance(a, b string) int {
	if len(a) > len(b) {
		a, b = b, a
	}

	m := len(a)
	n := len(b)

	if m == 0 {
		return n * w.InsertCost
	}
	if n == 0 {
		return m * w.DeleteCost
	}

	prev := make([]int, m+1)
	curr := make([]int, m+1)

	for j := 0; j <= m; j++ {
		prev[j] = j * w.DeleteCost
	}

	for i := 1; i <= n; i++ {
		curr[0] = i * w.InsertCost

		for j := 1; j <= m; j++ {
			if b[i-1] == a[j-1] {
				curr[j] = prev[j-1]
			} else {
				curr[j] = min3(
					prev[j]+w.DeleteCost,
					curr[j-1]+w.InsertCost,
					prev[j-1]+w.SubstituteCost,
				)
			}
		}

		prev, curr = curr, prev
	}

	return prev[m]
}

// ════════════════════════════════════════════════════════════════
// KEYBOARD-AWARE LEVENSHTEIN
// ════════════════════════════════════════════════════════════════

// KeyboardDistance provides reduced cost for adjacent key typos.
// Based on QWERTY keyboard layout.
var keyboardAdjacent = map[byte]string{
	'q': "wa",
	'w': "qeas",
	'e': "wrds",
	'r': "etfd",
	't': "rygf",
	'y': "tuhg",
	'u': "yijh",
	'i': "uokj",
	'o': "iplk",
	'p': "ol",
	'a': "qwsz",
	's': "weadzx",
	'd': "erfcxs",
	'f': "rtgvcd",
	'g': "tyhbvf",
	'h': "yujnbg",
	'j': "uikmnh",
	'k': "iolmj",
	'l': "opk",
	'z': "asx",
	'x': "sdzc",
	'c': "dfxv",
	'v': "fgcb",
	'b': "ghnv",
	'n': "hjbm",
	'm': "jkn",
}

// IsAdjacentKey checks if two keys are adjacent on a QWERTY keyboard.
func IsAdjacentKey(a, b byte) bool {
	// Normalize to lowercase
	if a >= 'A' && a <= 'Z' {
		a += 'a' - 'A'
	}
	if b >= 'A' && b <= 'Z' {
		b += 'a' - 'A'
	}

	adjacent, ok := keyboardAdjacent[a]
	if !ok {
		return false
	}

	for i := 0; i < len(adjacent); i++ {
		if adjacent[i] == b {
			return true
		}
	}
	return false
}

// KeyboardLevenshtein uses reduced cost for adjacent key substitutions.
type KeyboardLevenshtein struct {
	CaseSensitive   bool
	MaxDistance     int
	AdjacentKeyCost float64 // Cost multiplier for adjacent key typos (0.5 = half cost)
}

// NewKeyboardLevenshtein creates a keyboard-aware matcher.
func NewKeyboardLevenshtein() *KeyboardLevenshtein {
	return &KeyboardLevenshtein{
		CaseSensitive:   false,
		MaxDistance:     2,
		AdjacentKeyCost: 0.5,
	}
}

// Name returns the algorithm name.
func (k *KeyboardLevenshtein) Name() string {
	return "keyboard-levenshtein"
}

// Distance calculates keyboard-aware edit distance.
// Returns distance * 100 to maintain integer precision with fractional costs.
func (k *KeyboardLevenshtein) Distance(a, b string) int {
	a = normalize(a, k.CaseSensitive)
	b = normalize(b, k.CaseSensitive)

	return k.keyboardDistance(a, b)
}

// Similarity returns a similarity score between 0.0 and 1.0.
func (k *KeyboardLevenshtein) Similarity(a, b string) float64 {
	a = normalize(a, k.CaseSensitive)
	b = normalize(b, k.CaseSensitive)

	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}

	dist := k.keyboardDistance(a, b)
	maxLen := max(len(a), len(b))

	// Distance is scaled by 100, so normalize
	return 1.0 - float64(dist)/float64(maxLen*100)
}

// Match returns true if the distance is within MaxDistance.
func (k *KeyboardLevenshtein) Match(a, b string) bool {
	// MaxDistance is scaled by 100
	return k.Distance(a, b) <= k.MaxDistance*100
}

func (k *KeyboardLevenshtein) keyboardDistance(a, b string) int {
	if len(a) > len(b) {
		a, b = b, a
	}

	m := len(a)
	n := len(b)

	if m == 0 {
		return n * 100
	}
	if n == 0 {
		return m * 100
	}

	prev := make([]int, m+1)
	curr := make([]int, m+1)

	for j := 0; j <= m; j++ {
		prev[j] = j * 100
	}

	for i := 1; i <= n; i++ {
		curr[0] = i * 100

		for j := 1; j <= m; j++ {
			if b[i-1] == a[j-1] {
				curr[j] = prev[j-1]
			} else {
				subCost := 100
				if IsAdjacentKey(a[j-1], b[i-1]) {
					subCost = int(100 * k.AdjacentKeyCost)
				}

				curr[j] = min3(
					prev[j]+100,       // deletion
					curr[j-1]+100,     // insertion
					prev[j-1]+subCost, // substitution
				)
			}
		}

		prev, curr = curr, prev
	}

	return prev[m]
}
