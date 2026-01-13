// internal/fuzzy/jarowinkler.go

package fuzzy

// ════════════════════════════════════════════════════════════════
// JARO SIMILARITY
// ════════════════════════════════════════════════════════════════

// Jaro implements the Jaro similarity algorithm.
// It produces a similarity score between 0.0 and 1.0 based on:
//   - Number of matching characters
//   - Number of transpositions
//
// Two characters are considered matching if they are the same and
// not farther than floor(max(|s1|,|s2|)/2) - 1.
//
// Formula: jaro = (m/|s1| + m/|s2| + (m-t)/m) / 3
// Where: m = matches, t = transpositions/2
//
// Time Complexity: O(m*n)
// Space Complexity: O(m+n)
type Jaro struct {
	// CaseSensitive controls case sensitivity.
	CaseSensitive bool

	// MinSimilarity is the threshold for considering strings a match.
	MinSimilarity float64
}

// NewJaro creates a new Jaro matcher.
func NewJaro() *Jaro {
	return &Jaro{
		CaseSensitive: false,
		MinSimilarity: 0.7,
	}
}

// Name returns the algorithm name.
func (j *Jaro) Name() string {
	return "jaro"
}

// Distance returns an approximate edit distance derived from similarity.
// Not a true edit distance, but useful for ranking.
func (j *Jaro) Distance(a, b string) int {
	sim := j.Similarity(a, b)
	if sim == 1.0 {
		return 0
	}
	// Convert similarity to pseudo-distance
	maxLen := max(len(a), len(b))
	return int(float64(maxLen) * (1.0 - sim))
}

// Similarity calculates the Jaro similarity between two strings.
func (j *Jaro) Similarity(a, b string) float64 {
	a = normalize(a, j.CaseSensitive)
	b = normalize(b, j.CaseSensitive)

	return jaroSimilarity(a, b)
}

// Match returns true if similarity is above MinSimilarity.
func (j *Jaro) Match(a, b string) bool {
	return j.Similarity(a, b) >= j.MinSimilarity
}

// jaroSimilarity computes the Jaro similarity score.
func jaroSimilarity(a, b string) float64 {
	lenA := len(a)
	lenB := len(b)

	// Handle edge cases
	if lenA == 0 && lenB == 0 {
		return 1.0
	}
	if lenA == 0 || lenB == 0 {
		return 0.0
	}
	if a == b {
		return 1.0
	}

	// Calculate match window
	matchWindow := max(lenA, lenB)/2 - 1
	if matchWindow < 0 {
		matchWindow = 0
	}

	// Track matches
	aMatched := make([]bool, lenA)
	bMatched := make([]bool, lenB)

	matches := 0
	transpositions := 0

	// Find matches
	for i := 0; i < lenA; i++ {
		start := max(0, i-matchWindow)
		end := min(lenB, i+matchWindow+1)

		for j := start; j < end; j++ {
			if bMatched[j] || a[i] != b[j] {
				continue
			}
			aMatched[i] = true
			bMatched[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	// Count transpositions
	k := 0
	for i := 0; i < lenA; i++ {
		if !aMatched[i] {
			continue
		}
		for !bMatched[k] {
			k++
		}
		if a[i] != b[k] {
			transpositions++
		}
		k++
	}

	// Calculate Jaro similarity
	m := float64(matches)
	t := float64(transpositions) / 2.0

	return (m/float64(lenA) + m/float64(lenB) + (m-t)/m) / 3.0
}

// ════════════════════════════════════════════════════════════════
// JARO-WINKLER SIMILARITY
// ════════════════════════════════════════════════════════════════

// JaroWinkler implements the Jaro-Winkler similarity algorithm.
// It extends Jaro by giving more weight to strings that share a
// common prefix, which is useful because typos are less likely
// at the beginning of words.
//
// Formula: jw = jaro + (prefixLen * p * (1 - jaro))
// Where: p = scaling factor (typically 0.1), prefixLen <= 4
//
// This works well for:
//   - Names (people often get beginnings right)
//   - Function names (common prefixes like "add", "get", "set")
//   - Short strings where prefixes matter more
type JaroWinkler struct {
	// CaseSensitive controls case sensitivity.
	CaseSensitive bool

	// MinSimilarity is the threshold for considering strings a match.
	MinSimilarity float64

	// PrefixScale is the scaling factor for common prefix bonus.
	// Standard value is 0.1. Should not exceed 0.25.
	PrefixScale float64

	// MaxPrefixLength is the maximum prefix length to consider.
	// Standard value is 4.
	MaxPrefixLength int

	// BoostThreshold is the Jaro similarity above which prefix boost applies.
	// Standard value is 0.7 (only boost already-similar strings).
	BoostThreshold float64
}

// NewJaroWinkler creates a new Jaro-Winkler matcher with standard settings.
func NewJaroWinkler() *JaroWinkler {
	return &JaroWinkler{
		CaseSensitive:   false,
		MinSimilarity:   0.8,
		PrefixScale:     0.1,
		MaxPrefixLength: 4,
		BoostThreshold:  0.7,
	}
}

// Name returns the algorithm name.
func (jw *JaroWinkler) Name() string {
	return "jaro-winkler"
}

// Distance returns an approximate edit distance derived from similarity.
func (jw *JaroWinkler) Distance(a, b string) int {
	sim := jw.Similarity(a, b)
	if sim == 1.0 {
		return 0
	}
	maxLen := max(len(a), len(b))
	return int(float64(maxLen) * (1.0 - sim))
}

// Similarity calculates the Jaro-Winkler similarity between two strings.
func (jw *JaroWinkler) Similarity(a, b string) float64 {
	a = normalize(a, jw.CaseSensitive)
	b = normalize(b, jw.CaseSensitive)

	return jw.jaroWinklerSimilarity(a, b)
}

// Match returns true if similarity is above MinSimilarity.
func (jw *JaroWinkler) Match(a, b string) bool {
	return jw.Similarity(a, b) >= jw.MinSimilarity
}

func (jw *JaroWinkler) jaroWinklerSimilarity(a, b string) float64 {
	// Get base Jaro similarity
	jaroSim := jaroSimilarity(a, b)

	// Only apply prefix boost if above threshold
	if jaroSim < jw.BoostThreshold {
		return jaroSim
	}

	// Calculate common prefix length
	prefixLen := commonPrefixLength(a, b, jw.MaxPrefixLength)

	// Apply Winkler modification
	return jaroSim + float64(prefixLen)*jw.PrefixScale*(1.0-jaroSim)
}

// commonPrefixLength returns the length of common prefix up to maxLen.
func commonPrefixLength(a, b string, maxLen int) int {
	limit := min(min(len(a), len(b)), maxLen)

	for i := 0; i < limit; i++ {
		if a[i] != b[i] {
			return i
		}
	}

	return limit
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// JaroSimilarity is a convenience function for quick Jaro similarity.
func JaroSimilarity(a, b string) float64 {
	return jaroSimilarity(toLower(a), toLower(b))
}

// JaroWinklerSimilarity is a convenience function for quick JW similarity.
func JaroWinklerSimilarity(a, b string) float64 {
	jw := NewJaroWinkler()
	return jw.Similarity(a, b)
}

// ════════════════════════════════════════════════════════════════
// TUNED VARIANTS
// ════════════════════════════════════════════════════════════════

// JaroWinklerStrict creates a matcher with stricter matching.
// Useful for reducing false positives.
func JaroWinklerStrict() *JaroWinkler {
	return &JaroWinkler{
		CaseSensitive:   false,
		MinSimilarity:   0.9,
		PrefixScale:     0.1,
		MaxPrefixLength: 4,
		BoostThreshold:  0.8,
	}
}

// JaroWinklerLoose creates a matcher with looser matching.
// Useful for catching more potential matches.
func JaroWinklerLoose() *JaroWinkler {
	return &JaroWinkler{
		CaseSensitive:   false,
		MinSimilarity:   0.7,
		PrefixScale:     0.1,
		MaxPrefixLength: 4,
		BoostThreshold:  0.6,
	}
}

// JaroWinklerForNames creates a matcher optimized for person names.
// Names often have correct first letters but variations in spelling.
func JaroWinklerForNames() *JaroWinkler {
	return &JaroWinkler{
		CaseSensitive:   false,
		MinSimilarity:   0.85,
		PrefixScale:     0.15, // Higher prefix weight for names
		MaxPrefixLength: 4,
		BoostThreshold:  0.7,
	}
}

// JaroWinklerForCode creates a matcher optimized for code identifiers.
// Function names, variables, etc. often share prefixes.
func JaroWinklerForCode() *JaroWinkler {
	return &JaroWinkler{
		CaseSensitive:   false,
		MinSimilarity:   0.8,
		PrefixScale:     0.12,
		MaxPrefixLength: 6, // Longer prefixes for function names
		BoostThreshold:  0.7,
	}
}

// ════════════════════════════════════════════════════════════════
// NORMALIZED VARIANTS
// ════════════════════════════════════════════════════════════════

// JaroWinklerNormalized adjusts similarity based on string lengths.
// Short strings get slightly penalized to reduce false positives
// from coincidental matches.
type JaroWinklerNormalized struct {
	*JaroWinkler

	// MinLength is the minimum length below which penalty applies.
	MinLength int

	// LengthPenalty is the penalty factor for short strings.
	LengthPenalty float64
}

// NewJaroWinklerNormalized creates a length-normalized matcher.
func NewJaroWinklerNormalized() *JaroWinklerNormalized {
	return &JaroWinklerNormalized{
		JaroWinkler:   NewJaroWinkler(),
		MinLength:     4,
		LengthPenalty: 0.05,
	}
}

// Name returns the algorithm name.
func (jwn *JaroWinklerNormalized) Name() string {
	return "jaro-winkler-normalized"
}

// Similarity calculates length-normalized Jaro-Winkler similarity.
func (jwn *JaroWinklerNormalized) Similarity(a, b string) float64 {
	baseSim := jwn.JaroWinkler.Similarity(a, b)

	// Apply length penalty for short strings
	minLen := min(len(a), len(b))
	if minLen < jwn.MinLength {
		penalty := float64(jwn.MinLength-minLen) * jwn.LengthPenalty
		baseSim = baseSim * (1.0 - penalty)
		if baseSim < 0 {
			baseSim = 0
		}
	}

	return baseSim
}

// ════════════════════════════════════════════════════════════════
// SORTED JARO-WINKLER
// ════════════════════════════════════════════════════════════════

// SortedJaroWinkler handles multi-word strings by sorting words first.
// This helps match strings like "John Smith" with "Smith, John".
type SortedJaroWinkler struct {
	*JaroWinkler
}

// NewSortedJaroWinkler creates a sorted word matcher.
func NewSortedJaroWinkler() *SortedJaroWinkler {
	return &SortedJaroWinkler{
		JaroWinkler: NewJaroWinkler(),
	}
}

// Name returns the algorithm name.
func (sjw *SortedJaroWinkler) Name() string {
	return "sorted-jaro-winkler"
}

// Similarity calculates similarity after sorting words in both strings.
func (sjw *SortedJaroWinkler) Similarity(a, b string) float64 {
	a = normalize(a, sjw.CaseSensitive)
	b = normalize(b, sjw.CaseSensitive)

	// Split, sort, and rejoin
	aSorted := sortWords(a)
	bSorted := sortWords(b)

	return sjw.JaroWinkler.jaroWinklerSimilarity(aSorted, bSorted)
}

// sortWords splits a string into words, sorts them, and rejoins.
func sortWords(s string) string {
	// Split into words
	words := splitWords(s)
	if len(words) <= 1 {
		return s
	}

	// Simple insertion sort for small word lists
	for i := 1; i < len(words); i++ {
		for j := i; j > 0 && words[j] < words[j-1]; j-- {
			words[j], words[j-1] = words[j-1], words[j]
		}
	}

	// Rejoin
	result := words[0]
	for i := 1; i < len(words); i++ {
		result += " " + words[i]
	}
	return result
}

// splitWords splits a string into words.
func splitWords(s string) []string {
	var words []string
	var current []byte

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == ',' || c == ';' {
			if len(current) > 0 {
				words = append(words, string(current))
				current = current[:0]
			}
		} else {
			current = append(current, c)
		}
	}

	if len(current) > 0 {
		words = append(words, string(current))
	}

	return words
}
