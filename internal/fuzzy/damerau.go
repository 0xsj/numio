// internal/fuzzy/damerau.go

package fuzzy

// ════════════════════════════════════════════════════════════════
// DAMERAU-LEVENSHTEIN MATCHER
// ════════════════════════════════════════════════════════════════

// DamerauLevenshtein implements the Damerau-Levenshtein distance algorithm.
// It extends Levenshtein by also counting transpositions (swapping two
// adjacent characters) as a single edit operation.
//
// This is particularly useful for typos like:
//   - "teh" -> "the" (transposition, distance = 1)
//   - "recieve" -> "receive" (transposition, distance = 1)
//
// With standard Levenshtein, "teh" -> "the" would be distance 2.
//
// Time Complexity: O(m*n)
// Space Complexity: O(m*n) for the full algorithm
type DamerauLevenshtein struct {
	// CaseSensitive controls case sensitivity.
	CaseSensitive bool

	// MaxDistance is the threshold for considering strings a match.
	MaxDistance int
}

// NewDamerauLevenshtein creates a new Damerau-Levenshtein matcher.
func NewDamerauLevenshtein() *DamerauLevenshtein {
	return &DamerauLevenshtein{
		CaseSensitive: false,
		MaxDistance:   2,
	}
}

// Name returns the algorithm name.
func (d *DamerauLevenshtein) Name() string {
	return "damerau-levenshtein"
}

// Distance calculates the Damerau-Levenshtein distance between two strings.
// Includes transpositions as a single edit operation.
func (d *DamerauLevenshtein) Distance(a, b string) int {
	a = normalize(a, d.CaseSensitive)
	b = normalize(b, d.CaseSensitive)

	return damerauLevenshteinDistance(a, b)
}

// Similarity returns a similarity score between 0.0 and 1.0.
func (d *DamerauLevenshtein) Similarity(a, b string) float64 {
	a = normalize(a, d.CaseSensitive)
	b = normalize(b, d.CaseSensitive)

	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}

	dist := damerauLevenshteinDistance(a, b)
	maxLen := max(len(a), len(b))

	return 1.0 - float64(dist)/float64(maxLen)
}

// Match returns true if the distance is within MaxDistance.
func (d *DamerauLevenshtein) Match(a, b string) bool {
	return d.Distance(a, b) <= d.MaxDistance
}

// ════════════════════════════════════════════════════════════════
// CORE ALGORITHM (OPTIMAL STRING ALIGNMENT)
// ════════════════════════════════════════════════════════════════

// damerauLevenshteinDistance computes the Optimal String Alignment distance.
// This is the restricted edit distance that doesn't allow multiple edits
// on the same substring (simpler and faster than true Damerau-Levenshtein).
func damerauLevenshteinDistance(a, b string) int {
	m := len(a)
	n := len(b)

	// Handle edge cases
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Create matrix
	// We need 2 previous rows for transposition check
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	// Initialize first row and column
	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}

			d[i][j] = min3(
				d[i-1][j]+1,      // deletion
				d[i][j-1]+1,      // insertion
				d[i-1][j-1]+cost, // substitution
			)

			// Transposition check
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				d[i][j] = min(d[i][j], d[i-2][j-2]+cost)
			}
		}
	}

	return d[m][n]
}

// ════════════════════════════════════════════════════════════════
// TRUE DAMERAU-LEVENSHTEIN (UNRESTRICTED)
// ════════════════════════════════════════════════════════════════

// TrueDamerauLevenshtein implements the unrestricted Damerau-Levenshtein.
// This allows multiple edits on the same substring, which can produce
// smaller distances in some cases.
//
// Example where true DL differs from OSA:
//   - "ca" -> "abc": OSA = 3, True DL = 2
type TrueDamerauLevenshtein struct {
	CaseSensitive bool
	MaxDistance   int
}

// NewTrueDamerauLevenshtein creates a new true Damerau-Levenshtein matcher.
func NewTrueDamerauLevenshtein() *TrueDamerauLevenshtein {
	return &TrueDamerauLevenshtein{
		CaseSensitive: false,
		MaxDistance:   2,
	}
}

// Name returns the algorithm name.
func (t *TrueDamerauLevenshtein) Name() string {
	return "true-damerau-levenshtein"
}

// Distance calculates the true (unrestricted) Damerau-Levenshtein distance.
func (t *TrueDamerauLevenshtein) Distance(a, b string) int {
	a = normalize(a, t.CaseSensitive)
	b = normalize(b, t.CaseSensitive)

	return trueDamerauLevenshteinDistance(a, b)
}

// Similarity returns a similarity score between 0.0 and 1.0.
func (t *TrueDamerauLevenshtein) Similarity(a, b string) float64 {
	a = normalize(a, t.CaseSensitive)
	b = normalize(b, t.CaseSensitive)

	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}

	dist := trueDamerauLevenshteinDistance(a, b)
	maxLen := max(len(a), len(b))

	return 1.0 - float64(dist)/float64(maxLen)
}

// Match returns true if the distance is within MaxDistance.
func (t *TrueDamerauLevenshtein) Match(a, b string) bool {
	return t.Distance(a, b) <= t.MaxDistance
}

// trueDamerauLevenshteinDistance implements the unrestricted algorithm.
// Uses the approach described by Lowrance and Wagner.
func trueDamerauLevenshteinDistance(a, b string) int {
	m := len(a)
	n := len(b)

	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Build alphabet map for last occurrence
	da := make(map[byte]int)

	// Initialize matrix with extra row/column for boundary
	d := make([][]int, m+2)
	for i := range d {
		d[i] = make([]int, n+2)
	}

	maxDist := m + n
	d[0][0] = maxDist

	for i := 0; i <= m; i++ {
		d[i+1][0] = maxDist
		d[i+1][1] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j+1] = maxDist
		d[1][j+1] = j
	}

	// Fill the matrix
	for i := 1; i <= m; i++ {
		db := 0 // last column with matching character in b

		for j := 1; j <= n; j++ {
			i1 := da[b[j-1]] // last row with matching character in a
			j1 := db

			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
				db = j
			}

			d[i+1][j+1] = min4(
				d[i][j]+cost, // substitution
				d[i+1][j]+1,  // insertion
				d[i][j+1]+1,  // deletion
				d[i1][j1]+(i-i1-1)+1+(j-j1-1), // transposition
			)
		}

		da[a[i-1]] = i
	}

	return d[m+1][n+1]
}

// min4 returns the minimum of four integers.
func min4(a, b, c, d int) int {
	return min(min(a, b), min(c, d))
}

// ════════════════════════════════════════════════════════════════
// BOUNDED DAMERAU-LEVENSHTEIN
// ════════════════════════════════════════════════════════════════

// DamerauLevenshteinBounded computes distance with early termination.
// More efficient when you only care about close matches.
func DamerauLevenshteinBounded(a, b string, maxDist int) int {
	m := len(a)
	n := len(b)

	// Quick length check
	if abs(m-n) > maxDist {
		return maxDist + 1
	}

	if m == 0 {
		if n <= maxDist {
			return n
		}
		return maxDist + 1
	}
	if n == 0 {
		if m <= maxDist {
			return m
		}
		return maxDist + 1
	}

	// Use full matrix for transposition support
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	for i := 1; i <= m; i++ {
		rowMin := d[i][0]

		for j := 1; j <= n; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}

			d[i][j] = min3(
				d[i-1][j]+1,
				d[i][j-1]+1,
				d[i-1][j-1]+cost,
			)

			// Transposition
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				d[i][j] = min(d[i][j], d[i-2][j-2]+cost)
			}

			if d[i][j] < rowMin {
				rowMin = d[i][j]
			}
		}

		// Early termination
		if rowMin > maxDist {
			return maxDist + 1
		}
	}

	if d[m][n] <= maxDist {
		return d[m][n]
	}
	return maxDist + 1
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// TranspositionDistance is a convenience function for Damerau-Levenshtein.
func TranspositionDistance(a, b string) int {
	return damerauLevenshteinDistance(toLower(a), toLower(b))
}

// HasTransposition checks if the difference between strings is a transposition.
func HasTransposition(a, b string) bool {
	a = toLower(a)
	b = toLower(b)

	if len(a) != len(b) {
		return false
	}

	// Count differences
	diffs := 0
	diffPositions := [2]int{}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			if diffs >= 2 {
				return false
			}
			diffPositions[diffs] = i
			diffs++
		}
	}

	// Check if exactly 2 differences that are swapped
	if diffs != 2 {
		return false
	}

	i, j := diffPositions[0], diffPositions[1]
	return a[i] == b[j] && a[j] == b[i]
}

// ════════════════════════════════════════════════════════════════
// COMMON TYPO PATTERNS
// ════════════════════════════════════════════════════════════════

// TypoType represents the type of typo detected.
type TypoType int

const (
	TypoNone TypoType = iota
	TypoInsertion
	TypoDeletion
	TypoSubstitution
	TypoTransposition
	TypoMultiple
)

// String returns the typo type name.
func (t TypoType) String() string {
	switch t {
	case TypoNone:
		return "none"
	case TypoInsertion:
		return "insertion"
	case TypoDeletion:
		return "deletion"
	case TypoSubstitution:
		return "substitution"
	case TypoTransposition:
		return "transposition"
	case TypoMultiple:
		return "multiple"
	default:
		return "unknown"
	}
}

// DetectTypoType identifies the type of typo between two strings.
// Only works reliably for single-edit typos.
func DetectTypoType(typo, correct string) TypoType {
	typo = toLower(typo)
	correct = toLower(correct)

	if typo == correct {
		return TypoNone
	}

	dist := damerauLevenshteinDistance(typo, correct)

	if dist > 1 {
		return TypoMultiple
	}

	// Single edit - determine type
	if len(typo) == len(correct)+1 {
		return TypoInsertion
	}
	if len(typo)+1 == len(correct) {
		return TypoDeletion
	}
	if len(typo) == len(correct) {
		if HasTransposition(typo, correct) {
			return TypoTransposition
		}
		return TypoSubstitution
	}

	return TypoMultiple
}
