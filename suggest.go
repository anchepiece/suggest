package suggest

import (
	"sort"
	"strings"
)

// Suggest struct wraps a list of Commands with Options
type Suggest struct {
	Options  Options
	Commands []string
}

// Options contains customizable weights to apply to a query
type Options struct {
	CostSwap            int  `json:"costswap,omitempty"`
	CostSubstitution    int  `json:"costsubstitution,omitempty"`
	CostInsertion       int  `json:"costinsertion,omitempty"`
	CostDeletion        int  `json:"costdeletion,omitempty"`
	SimilarityMinimum   int  `json:"similarityminimum,omitempty"`
	AutocorrectDisabled bool `json:"autocorrectdisabled,omitempty"`
}

const (
	defaultCostSwap         = 0
	defaultCostSubstitution = 2
	defaultCostInsertion    = 1
	defaultCostDeletion     = 4
	defaultSimilarity       = 6
)

func (options *Options) getCostSwap() int {
	if options.CostSwap <= 0 {
		options.CostSwap = defaultCostSwap
	}
	return options.CostSwap
}

func (options *Options) getCostSubstitution() int {
	if options.CostSubstitution <= 0 {
		options.CostSubstitution = defaultCostSubstitution
	}
	return options.CostSubstitution
}

func (options *Options) getCostInsertion() int {
	if options.CostInsertion <= 0 {
		options.CostInsertion = defaultCostInsertion
	}
	return options.CostInsertion
}

func (options *Options) getCostDeletion() int {
	if options.CostDeletion <= 0 {
		options.CostDeletion = defaultCostDeletion
	}
	return options.CostDeletion
}

func (options *Options) getSimilarityMinimum() int {
	if options.SimilarityMinimum <= 0 {
		options.SimilarityMinimum = defaultSimilarity
	}
	return options.SimilarityMinimum
}

// Result contains the complete calculated result of a successful query
type Result struct {
	// Autocorrect is the identified best single match when option enabled
	Autocorrect string
	// Matches contains the full listing of similar entries with scores below the similarity minimum
	Matches []string
}

// Success returns true when containing at least one valid result
func (r *Result) Success() bool {
	return len(r.Matches) > 0
}

// New allocates a new Suggest with the given options
func New(options Options) *Suggest {
	return &Suggest{
		Options: options,
	}
}

// Query will calculate the best matching entries including
// an Autocorrect option if requested
func (s *Suggest) Query(query string) (Result, error) {
	return s.QueryAgainst(query, s.Commands)
}

// QueryAgainst will calculate the best matching entries against passed commands
// including an Autocorrect option if requested
func (s *Suggest) QueryAgainst(query string, commands []string) (Result, error) {

	best := 100
	var autocorrect string
	var matches []string

	var scores []int
	scores = make([]int, len(commands))

	scoreboard := make(map[string]int)

	for i, candidate := range commands {

		scores[i] = s.CalculateSimilarity(query, candidate)

		if strings.ToLower(query) == strings.ToLower(candidate) {
			scores[i] = -1
		}

		if query == candidate {
			scores[i] = -2
		}

		if scores[i] <= s.Options.getSimilarityMinimum() {
			//matches = append(matches, candidate)
			scoreboard[candidate] = scores[i]
		}

		//fmt.Println("query/candidate/score: ", query, candidate, scores[i])
		if best >= scores[i] {
			best = scores[i]
		}
	}

	if len(scoreboard) == 0 {
		// no candidates meet the minimum similarity
		return Result{}, nil
	}

	if !s.Options.AutocorrectDisabled {
		for i, score := range scores {
			// TODO: do we need a threshold?
			// if best > 0 {
			// 	break
			// }
			if score == best {
				autocorrect = commands[i]
				break
			}
		}
	}

	//fmt.Println("scoreboard", scoreboard)
	sort.Ints(scores)
	for _, v := range scores {
		for k := range scoreboard {
			if v == scoreboard[k] {
				if contains(matches, k) {
					continue
				}
				matches = append(matches, k)
				break
			}
		}
	}

	//fmt.Println("matches", matches)
	return Result{
		Autocorrect: autocorrect,
		Matches:     matches}, nil
}

// contains checks for the existence of a string in a slice
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// CalculateSimilarity returns a computed distance between two words
// based on weighted costs
func (s *Suggest) CalculateSimilarity(query string, candidate string) int {

	lenQuery := len(query)
	lenCandidate := len(candidate)

	size := lenCandidate + 1
	row0 := make([]int, size)
	row1 := make([]int, size)
	row2 := make([]int, size)

	result := 0

	for j := 0; j <= lenCandidate; j++ {
		row1[j] = j * s.Options.getCostInsertion()
	}

	for i := 0; i < lenQuery; i++ {
		row2[0] = (i + 1) * s.Options.getCostDeletion()

		for j := 0; j < lenCandidate; j++ {

			// Substitution
			row2[j+1] = row1[j]
			if query[i] != candidate[j] {
				row2[j+1] = row1[j] + s.Options.getCostSubstitution()
			}

			// Swap
			if i > 0 && j > 0 && query[i-1] == candidate[j] && query[i] == candidate[j-1] {
				swap := row0[j-1] + s.Options.getCostSwap()
				if row2[j+1] > swap {
					row2[j+1] = swap
				}
			}

			// Deletion
			if row2[j+1] > row1[j+1]+s.Options.getCostDeletion() {
				row2[j+1] = row1[j+1] + s.Options.getCostDeletion()
			}

			// Insertion
			if row2[j+1] > row2[j]+s.Options.getCostInsertion() {
				row2[j+1] = row2[j] + s.Options.getCostInsertion()
			}
		}

		t := row0
		row0 = row1
		row1 = row2
		row2 = t
	}
	result = row1[lenCandidate]
	return result
}

// Autocorrect provides the most likely match falling under the minimum cost
// If there is no match close enough, an empty string is returned
func (s *Suggest) Autocorrect(query string) (string, error) {
	return s.AutocorrectAgainst(query, s.Commands)
}

// AutocorrectAgainst provides the most likely match from passed commands
// If there is no match close enough, an empty string is returned
func (s *Suggest) AutocorrectAgainst(query string, commands []string) (string, error) {

	tmp := s.Options.AutocorrectDisabled
	s.Options.AutocorrectDisabled = false
	q, err := s.QueryAgainst(query, commands)
	// restore option
	s.Options.AutocorrectDisabled = tmp

	if err == nil && q.Success() {
		return q.Autocorrect, nil
	}
	return "", err
}

// ExactMatch returns a valid matching entry from available commands
func (s *Suggest) ExactMatch(query string) string {
	return s.ExactMatchAgainst(query, s.Commands)
}

// ExactMatchAgainst returns a valid matching entry from passed commands
func (s *Suggest) ExactMatchAgainst(query string, commands []string) string {

	for _, arg := range commands {
		if strings.ToLower(arg) == strings.ToLower(query) {
			return arg
		}
	}
	return ""
}
