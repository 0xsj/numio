// internal/fuzzy/fuzzy.go

package fuzzy

import "strings"

// ════════════════════════════════════════════════════════════════
// MAIN FUZZY MATCHER
// ════════════════════════════════════════════════════════════════

// Fuzzy provides a unified interface for fuzzy string matching.
// It can use different algorithms and combines their results.
type Fuzzy struct {
	// Algorithm is the primary matching algorithm.
	Algorithm Algorithm

	// Options configures matching behavior.
	Options Options

	// matcher is the underlying matcher implementation.
	matcher Matcher
}

// New creates a new Fuzzy matcher with default settings.
func New() *Fuzzy {
	return &Fuzzy{
		Algorithm: AlgorithmDamerauLevenshtein,
		Options:   DefaultOptions(),
		matcher:   NewDamerauLevenshtein(),
	}
}

// NewWithAlgorithm creates a Fuzzy matcher with a specific algorithm.
func NewWithAlgorithm(alg Algorithm) *Fuzzy {
	f := &Fuzzy{
		Algorithm: alg,
		Options:   DefaultOptions(),
	}
	f.matcher = f.createMatcher(alg)
	return f
}

// NewWithOptions creates a Fuzzy matcher with custom options.
func NewWithOptions(opts Options) *Fuzzy {
	f := New()
	f.Options = opts
	f.applyOptions()
	return f
}

// createMatcher creates a matcher for the given algorithm.
func (f *Fuzzy) createMatcher(alg Algorithm) Matcher {
	switch alg {
	case AlgorithmLevenshtein:
		m := NewLevenshtein()
		m.CaseSensitive = f.Options.CaseSensitive
		m.MaxDistance = f.Options.MaxDistance
		return m
	case AlgorithmDamerauLevenshtein:
		m := NewDamerauLevenshtein()
		m.CaseSensitive = f.Options.CaseSensitive
		m.MaxDistance = f.Options.MaxDistance
		return m
	case AlgorithmJaroWinkler:
		m := NewJaroWinkler()
		m.CaseSensitive = f.Options.CaseSensitive
		m.MinSimilarity = f.Options.MinSimilarity
		return m
	default:
		return NewDamerauLevenshtein()
	}
}

// applyOptions applies current options to the matcher.
func (f *Fuzzy) applyOptions() {
	switch m := f.matcher.(type) {
	case *Levenshtein:
		m.CaseSensitive = f.Options.CaseSensitive
		m.MaxDistance = f.Options.MaxDistance
	case *DamerauLevenshtein:
		m.CaseSensitive = f.Options.CaseSensitive
		m.MaxDistance = f.Options.MaxDistance
	case *JaroWinkler:
		m.CaseSensitive = f.Options.CaseSensitive
		m.MinSimilarity = f.Options.MinSimilarity
	}
}

// SetAlgorithm changes the matching algorithm.
func (f *Fuzzy) SetAlgorithm(alg Algorithm) {
	f.Algorithm = alg
	f.matcher = f.createMatcher(alg)
}

// SetOptions updates the options.
func (f *Fuzzy) SetOptions(opts Options) {
	f.Options = opts
	f.applyOptions()
}

// ════════════════════════════════════════════════════════════════
// CORE MATCHING METHODS
// ════════════════════════════════════════════════════════════════

// Distance returns the edit distance between two strings.
func (f *Fuzzy) Distance(a, b string) int {
	return f.matcher.Distance(a, b)
}

// Similarity returns a similarity score between 0.0 and 1.0.
func (f *Fuzzy) Similarity(a, b string) float64 {
	return f.matcher.Similarity(a, b)
}

// Match checks if two strings are considered a match.
func (f *Fuzzy) Match(a, b string) bool {
	return f.matcher.Match(a, b)
}

// ════════════════════════════════════════════════════════════════
// SEARCH METHODS
// ════════════════════════════════════════════════════════════════

// Find searches for the best match for query in the candidates list.
// Returns the best match or empty string if none found.
func (f *Fuzzy) Find(query string, candidates []string) string {
	results := f.Search(query, candidates)
	if len(results) == 0 {
		return ""
	}
	return results.Best().Text
}

// Search finds all matches for query in candidates, ranked by quality.
func (f *Fuzzy) Search(query string, candidates []string) Results {
	var results Results

	for _, candidate := range candidates {
		dist := f.matcher.Distance(query, candidate)
		sim := f.matcher.Similarity(query, candidate)

		// Check if it passes thresholds
		passesDistance := dist <= f.Options.MaxDistance
		passesSimilarity := sim >= f.Options.MinSimilarity

		if passesDistance || passesSimilarity {
			// Calculate rank (lower is better)
			// Combine distance and similarity into single score
			rank := float64(dist) - sim*10

			results = append(results, Result{
				Text:       candidate,
				Distance:   dist,
				Similarity: sim,
				Rank:       rank,
			})
		}
	}

	// Sort by rank
	sortResults(results)

	// Apply limit
	if f.Options.Limit > 0 && len(results) > f.Options.Limit {
		results = results[:f.Options.Limit]
	}

	return results
}

// FindClosest finds the closest match regardless of thresholds.
// Always returns the best match even if it's not very good.
func (f *Fuzzy) FindClosest(query string, candidates []string) Result {
	if len(candidates) == 0 {
		return Result{}
	}

	best := Result{
		Text:       candidates[0],
		Distance:   f.matcher.Distance(query, candidates[0]),
		Similarity: f.matcher.Similarity(query, candidates[0]),
	}
	best.Rank = float64(best.Distance) - best.Similarity*10

	for _, candidate := range candidates[1:] {
		dist := f.matcher.Distance(query, candidate)
		sim := f.matcher.Similarity(query, candidate)
		rank := float64(dist) - sim*10

		if rank < best.Rank {
			best = Result{
				Text:       candidate,
				Distance:   dist,
				Similarity: sim,
				Rank:       rank,
			}
		}
	}

	return best
}

// ════════════════════════════════════════════════════════════════
// SUGGESTION METHODS
// ════════════════════════════════════════════════════════════════

// Suggest returns suggestions for a potentially misspelled query.
// Useful for "did you mean?" functionality.
func (f *Fuzzy) Suggest(query string, dictionary []string) []string {
	results := f.Search(query, dictionary)

	suggestions := make([]string, len(results))
	for i, r := range results {
		suggestions[i] = r.Text
	}

	return suggestions
}

// SuggestOne returns the single best suggestion or empty string.
func (f *Fuzzy) SuggestOne(query string, dictionary []string) string {
	return f.Find(query, dictionary)
}

// SuggestWithConfidence returns suggestion with confidence score.
func (f *Fuzzy) SuggestWithConfidence(query string, dictionary []string) (string, float64) {
	result := f.FindClosest(query, dictionary)
	if result.Text == "" {
		return "", 0
	}
	return result.Text, result.Similarity
}

// ════════════════════════════════════════════════════════════════
// AUTOCORRECT
// ════════════════════════════════════════════════════════════════

// Autocorrect returns the corrected string if a good match is found,
// otherwise returns the original string.
func (f *Fuzzy) Autocorrect(query string, dictionary []string) string {
	result := f.FindClosest(query, dictionary)

	// Only autocorrect if very confident
	if result.Distance <= 1 || result.Similarity >= 0.9 {
		return result.Text
	}

	return query
}

// AutocorrectWithThreshold corrects if within the specified distance.
func (f *Fuzzy) AutocorrectWithThreshold(query string, dictionary []string, maxDist int) string {
	result := f.FindClosest(query, dictionary)

	if result.Distance <= maxDist {
		return result.Text
	}

	return query
}

// ════════════════════════════════════════════════════════════════
// FACTORY FUNCTIONS
// ════════════════════════════════════════════════════════════════

// ForTypos creates a matcher optimized for catching typos.
// Uses Damerau-Levenshtein with reasonable defaults.
func ForTypos() *Fuzzy {
	return &Fuzzy{
		Algorithm: AlgorithmDamerauLevenshtein,
		Options: Options{
			CaseSensitive: false,
			MaxDistance:   2,
			MinSimilarity: 0.7,
			Limit:         5,
		},
		matcher: NewDamerauLevenshtein(),
	}
}

// ForSimilarity creates a matcher optimized for similarity scoring.
// Uses Jaro-Winkler which produces better similarity scores.
func ForSimilarity() *Fuzzy {
	return &Fuzzy{
		Algorithm: AlgorithmJaroWinkler,
		Options: Options{
			CaseSensitive: false,
			MaxDistance:   3,
			MinSimilarity: 0.8,
			Limit:         5,
		},
		matcher: NewJaroWinkler(),
	}
}

// ForFunctionNames creates a matcher optimized for function name typos.
// Tuned for identifiers like "adddays", "daysbetween", etc.
func ForFunctionNames() *Fuzzy {
	return &Fuzzy{
		Algorithm: AlgorithmDamerauLevenshtein,
		Options: Options{
			CaseSensitive: false,
			MaxDistance:   2,
			MinSimilarity: 0.75,
			Limit:         3,
		},
		matcher: NewDamerauLevenshtein(),
	}
}

// ForStrict creates a matcher with strict matching criteria.
// Reduces false positives at the cost of missing some matches.
func ForStrict() *Fuzzy {
	return &Fuzzy{
		Algorithm: AlgorithmDamerauLevenshtein,
		Options: Options{
			CaseSensitive: false,
			MaxDistance:   1,
			MinSimilarity: 0.9,
			Limit:         3,
		},
		matcher: NewDamerauLevenshtein(),
	}
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS (PACKAGE LEVEL)
// ════════════════════════════════════════════════════════════════

// defaultMatcher is the package-level default matcher.
var defaultMatcher = New()

// Find finds the best match using the default matcher.
func Find(query string, candidates []string) string {
	return defaultMatcher.Find(query, candidates)
}

// Search finds all matches using the default matcher.
func Search(query string, candidates []string) Results {
	return defaultMatcher.Search(query, candidates)
}

// Match checks if two strings match using the default matcher.
func Match(a, b string) bool {
	return defaultMatcher.Match(a, b)
}

// Distance returns edit distance using the default matcher.
func Distance(a, b string) int {
	return defaultMatcher.Distance(a, b)
}

// Similarity returns similarity using the default matcher.
func Similarity(a, b string) float64 {
	return defaultMatcher.Similarity(a, b)
}

// Suggest returns suggestions using the default matcher.
func Suggest(query string, dictionary []string) []string {
	return defaultMatcher.Suggest(query, dictionary)
}

// ════════════════════════════════════════════════════════════════
// HYBRID MATCHER
// ════════════════════════════════════════════════════════════════

// HybridMatcher combines multiple algorithms for better results.
// It uses different algorithms and combines their scores.
type HybridMatcher struct {
	matchers []Matcher
	weights  []float64
	Options  Options
}

// NewHybridMatcher creates a matcher that combines multiple algorithms.
func NewHybridMatcher() *HybridMatcher {
	return &HybridMatcher{
		matchers: []Matcher{
			NewDamerauLevenshtein(),
			NewJaroWinkler(),
		},
		weights: []float64{0.6, 0.4},
		Options: DefaultOptions(),
	}
}

// Name returns the algorithm name.
func (h *HybridMatcher) Name() string {
	return "hybrid"
}

// Distance returns the weighted average distance.
func (h *HybridMatcher) Distance(a, b string) int {
	if len(h.matchers) == 0 {
		return 0
	}

	total := 0.0
	for i, m := range h.matchers {
		weight := 1.0
		if i < len(h.weights) {
			weight = h.weights[i]
		}
		total += float64(m.Distance(a, b)) * weight
	}

	return int(total)
}

// Similarity returns the weighted average similarity.
func (h *HybridMatcher) Similarity(a, b string) float64 {
	if len(h.matchers) == 0 {
		return 0
	}

	total := 0.0
	totalWeight := 0.0

	for i, m := range h.matchers {
		weight := 1.0
		if i < len(h.weights) {
			weight = h.weights[i]
		}
		total += m.Similarity(a, b) * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}

	return total / totalWeight
}

// Match returns true if any matcher considers it a match.
func (h *HybridMatcher) Match(a, b string) bool {
	for _, m := range h.matchers {
		if m.Match(a, b) {
			return true
		}
	}
	return false
}

// Search performs hybrid search across all matchers.
func (h *HybridMatcher) Search(query string, candidates []string) Results {
	var results Results

	for _, candidate := range candidates {
		dist := h.Distance(query, candidate)
		sim := h.Similarity(query, candidate)

		if dist <= h.Options.MaxDistance || sim >= h.Options.MinSimilarity {
			results = append(results, Result{
				Text:       candidate,
				Distance:   dist,
				Similarity: sim,
				Rank:       float64(dist) - sim*10,
			})
		}
	}

	sortResults(results)

	if h.Options.Limit > 0 && len(results) > h.Options.Limit {
		results = results[:h.Options.Limit]
	}

	return results
}

// ════════════════════════════════════════════════════════════════
// UTILITY: DID YOU MEAN
// ════════════════════════════════════════════════════════════════

// DidYouMean generates a "did you mean?" suggestion message.
// Returns empty string if no good suggestion found.
func DidYouMean(query string, dictionary []string) string {
	f := ForFunctionNames()
	result := f.FindClosest(query, dictionary)

	// Only suggest if reasonably close
	if result.Text == "" || result.Distance > 2 {
		return ""
	}

	// Don't suggest if it's the same
	if strings.EqualFold(query, result.Text) {
		return ""
	}

	return result.Text
}

// FormatDidYouMean formats a complete "did you mean?" message.
func FormatDidYouMean(query string, dictionary []string) string {
	suggestion := DidYouMean(query, dictionary)
	if suggestion == "" {
		return ""
	}
	return "Did you mean: " + suggestion + "?"
}
