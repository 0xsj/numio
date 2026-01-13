// internal/fuzzy/matcher.go

// Package fuzzy provides fuzzy string matching algorithms for typo tolerance.
package fuzzy

// ════════════════════════════════════════════════════════════════
// INTERFACES
// ════════════════════════════════════════════════════════════════

// Matcher is the interface implemented by all fuzzy matching algorithms.
type Matcher interface {
	// Name returns the algorithm name.
	Name() string

	// Distance returns the edit distance between two strings.
	// Lower values indicate more similar strings.
	// Returns -1 if the algorithm doesn't support distance calculation.
	Distance(a, b string) int

	// Similarity returns a similarity score between 0.0 and 1.0.
	// Higher values indicate more similar strings.
	// Returns -1 if the algorithm doesn't support similarity calculation.
	Similarity(a, b string) float64

	// Match checks if two strings are considered a match
	// based on the algorithm's default threshold.
	Match(a, b string) bool
}

// ════════════════════════════════════════════════════════════════
// RESULT TYPES
// ════════════════════════════════════════════════════════════════

// Result represents a fuzzy match result.
type Result struct {
	// Text is the matched string.
	Text string

	// Distance is the edit distance (if applicable).
	Distance int

	// Similarity is the similarity score 0.0-1.0 (if applicable).
	Similarity float64

	// Rank is the overall ranking score (lower is better).
	Rank float64
}

// Results is a slice of Result with sorting capabilities.
type Results []Result

// Len implements sort.Interface.
func (r Results) Len() int {
	return len(r)
}

// Less implements sort.Interface (sorts by Rank ascending).
func (r Results) Less(i, j int) bool {
	return r[i].Rank < r[j].Rank
}

// Swap implements sort.Interface.
func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Best returns the best match (lowest rank), or empty Result if none.
func (r Results) Best() Result {
	if len(r) == 0 {
		return Result{}
	}

	best := r[0]
	for _, result := range r[1:] {
		if result.Rank < best.Rank {
			best = result
		}
	}
	return best
}

// Top returns the top n matches sorted by rank.
func (r Results) Top(n int) Results {
	if n <= 0 || len(r) == 0 {
		return nil
	}

	// Copy and sort
	sorted := make(Results, len(r))
	copy(sorted, r)
	sortResults(sorted)

	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// Filter returns results that pass the filter function.
func (r Results) Filter(fn func(Result) bool) Results {
	var filtered Results
	for _, result := range r {
		if fn(result) {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// WithinDistance returns results with distance <= maxDistance.
func (r Results) WithinDistance(maxDistance int) Results {
	return r.Filter(func(res Result) bool {
		return res.Distance >= 0 && res.Distance <= maxDistance
	})
}

// AboveSimilarity returns results with similarity >= minSimilarity.
func (r Results) AboveSimilarity(minSimilarity float64) Results {
	return r.Filter(func(res Result) bool {
		return res.Similarity >= minSimilarity
	})
}

// sortResults sorts results by rank (simple insertion sort for small slices).
func sortResults(r Results) {
	for i := 1; i < len(r); i++ {
		for j := i; j > 0 && r[j].Rank < r[j-1].Rank; j-- {
			r[j], r[j-1] = r[j-1], r[j]
		}
	}
}

// ════════════════════════════════════════════════════════════════
// OPTIONS
// ════════════════════════════════════════════════════════════════

// Options configures fuzzy matching behavior.
type Options struct {
	// CaseSensitive controls whether matching is case-sensitive.
	// Default: false (case-insensitive).
	CaseSensitive bool

	// MaxDistance is the maximum edit distance to consider a match.
	// Default: 2.
	MaxDistance int

	// MinSimilarity is the minimum similarity score to consider a match.
	// Default: 0.6.
	MinSimilarity float64

	// Limit is the maximum number of results to return.
	// Default: 5. Use 0 for unlimited.
	Limit int
}

// DefaultOptions returns the default options.
func DefaultOptions() Options {
	return Options{
		CaseSensitive: false,
		MaxDistance:   2,
		MinSimilarity: 0.6,
		Limit:         5,
	}
}

// ════════════════════════════════════════════════════════════════
// ALGORITHM ENUM
// ════════════════════════════════════════════════════════════════

// Algorithm represents a fuzzy matching algorithm type.
type Algorithm int

const (
	// AlgorithmLevenshtein uses basic Levenshtein edit distance.
	AlgorithmLevenshtein Algorithm = iota

	// AlgorithmDamerauLevenshtein uses Damerau-Levenshtein (includes transpositions).
	AlgorithmDamerauLevenshtein

	// AlgorithmJaroWinkler uses Jaro-Winkler similarity.
	AlgorithmJaroWinkler

	// AlgorithmSoundex uses Soundex phonetic matching.
	AlgorithmSoundex

	// AlgorithmHybrid combines multiple algorithms.
	AlgorithmHybrid
)

// String returns the algorithm name.
func (a Algorithm) String() string {
	switch a {
	case AlgorithmLevenshtein:
		return "levenshtein"
	case AlgorithmDamerauLevenshtein:
		return "damerau-levenshtein"
	case AlgorithmJaroWinkler:
		return "jaro-winkler"
	case AlgorithmSoundex:
		return "soundex"
	case AlgorithmHybrid:
		return "hybrid"
	default:
		return "unknown"
	}
}

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min3 returns the minimum of three integers.
func min3(a, b, c int) int {
	return min(min(a, b), c)
}

// abs returns the absolute value of an integer.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// toLower converts a string to lowercase without using strings package.
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// normalize prepares a string for comparison based on options.
func normalize(s string, caseSensitive bool) string {
	if caseSensitive {
		return s
	}
	return toLower(s)
}
