package suggest_test

import (
	"reflect"
	"testing"

	"github.com/anchepiece/suggest"
)

var writeTests = []struct {
	query    string
	commands []string
	expected string
}{
	{
		"test",
		[]string{"key", "value", "Test"},
		"Test",
	},
	{
		"unique",
		[]string{"key", "value", "Test"},
		"",
	},
}

func TestExactMatch(t *testing.T) {

	s := suggest.Suggest{}

	for _, tt := range writeTests {
		query := tt.query
		s.Commands = tt.commands
		actual := s.ExactMatch(query)

		if actual != tt.expected {
			t.Errorf("ExactMatch(%q, %v) = %q, want %q", tt.query, tt.commands, actual, tt.expected)
		}
	}
}

func TestExactMatchAgainst(t *testing.T) {

	s := suggest.Suggest{}

	for _, tt := range writeTests {
		query := tt.query
		commands := tt.commands
		actual := s.ExactMatchAgainst(query, commands)

		if actual != tt.expected {
			t.Errorf("ExactMatchAgainst(%q, %v) = %q, want %q", tt.query, tt.commands, actual, tt.expected)
		}
	}
}

var distanceTests = []struct {
	words    []string
	expected int
}{
	{
		[]string{"", ""},
		0,
	},
	{
		[]string{"", "a"},
		1,
	},
	{
		[]string{"fgrep", "fgerp"}, // swap cost is 0
		0,
	},
	{
		[]string{"fgrep", "fgerps"}, // insertion cost is 1
		1,
	},
	{
		[]string{"fgrep", "fgre"}, // deletion cost is 4
		4,
	},
	{
		[]string{"fgrep", "fgreP"}, // substitution cost is 2
		2,
	},
	{
		[]string{"word", "wrdo"},
		5,
	},
	{
		[]string{"kissing", "sitting"},
		6,
	},
	{
		[]string{"two words", "one two three"},
		12,
	},
}

// TestCalculateSimilarity tests the expected distance using the default parameters
func TestCalculateSimilarity(t *testing.T) {

	s := suggest.Suggest{}

	for _, tt := range distanceTests {
		actual := s.CalculateSimilarity(tt.words[0], tt.words[1])
		if actual != tt.expected {
			t.Errorf("CalculateSimilarity(%v) = %d, want %d", tt.words, actual, tt.expected)
		}
	}
}

var autocorrectTests = []struct {
	query    string
	commands []string
	expected string
}{
	{
		"test",
		[]string{"key", "value", "Test", "tEst"},
		"Test",
	},
	{
		"unique",
		[]string{"key", "value", "Test"},
		"",
	},

	{
		"install",
		[]string{"insatll", "install", "Install"},
		"install",
	},
}

func TestAutocorrect(t *testing.T) {

	s := suggest.Suggest{}

	if s.Options.AutocorrectDisabled != false {
		t.Errorf("AutocorrectDisabled should be initialized to false")
	}

	for _, tt := range autocorrectTests {
		query := tt.query
		s.Commands = tt.commands
		actual, err := s.Autocorrect(query)
		if err != nil {
			t.Errorf("Autocorrect(%v) returned error %v", tt.query, err)
			continue
		}
		if actual != tt.expected {
			t.Errorf("Autocorrect(%q, %v) = %q, want %q", tt.query, tt.commands, actual, tt.expected)
		}
	}
}

func TestAutocorrectAgainst(t *testing.T) {

	s := suggest.Suggest{}

	if s.Options.AutocorrectDisabled != false {
		t.Errorf("AutocorrectDisabled should be initialized to false")
	}

	for _, tt := range autocorrectTests {
		query := tt.query
		s.Commands = tt.commands
		actual, err := s.AutocorrectAgainst(query, tt.commands)
		if err != nil {
			t.Errorf("AutocorrectAgainst(%v) returned error %v", tt.query, err)
			continue
		}
		if actual != tt.expected {
			t.Errorf("AutocorrectAgainst(%q, %v) = %q, want %q", tt.query, tt.commands, actual, tt.expected)
		}
	}
}

func TestAutocorrectDisabled(t *testing.T) {

	var autocorrectDisabledTests = []struct {
		query    string
		commands []string
		expected string
	}{
		{
			"test",
			[]string{"test"},
			"",
		},
		{
			"unique",
			[]string{"key", "value", "test"},
			"",
		},

		{
			"isntall",
			[]string{"help", "branch", "install"},
			"",
		},
	}

	s := suggest.New(suggest.Options{AutocorrectDisabled: true})

	for _, tt := range autocorrectDisabledTests {
		query := tt.query
		s.Commands = tt.commands
		result, err := s.Query(query)
		actual := result.Autocorrect

		if err != nil {
			t.Errorf("Autocorrect(%v) returned error %v", tt.query, err)
			continue
		}

		if actual != tt.expected {
			t.Errorf("Autocorrect(%q, %v) = %q, want %q", tt.query, tt.commands, actual, tt.expected)
		}
	}
}

var queryTests = []struct {
	query    string
	commands []string
	expected []string
}{
	{
		"kittens",
		[]string{"tests", "sitting", "mittens"},
		[]string{"mittens", "sitting"},
	},
	{
		"isntall",
		[]string{"help", "branch", "install"},
		[]string{"install"},
	},
	{
		"Arnold Swarzeneger",
		[]string{"Arnold Schwarzenegger", "skip", "list"},
		[]string{"Arnold Schwarzenegger"},
	},
}

func TestQueryAgainst(t *testing.T) {

	s := suggest.Suggest{}

	for _, tt := range queryTests {
		query := tt.query
		result, err := s.QueryAgainst(query, tt.commands)
		actual := result.Matches

		if err != nil {
			t.Errorf("TestQueryAgainst(%v) returned error %v", tt.query, err)
			continue
		}

		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("TestQueryAgainst(%q, %v) = %q, want %q", tt.query, tt.commands, actual, tt.expected)
		}
	}
}

func TestQueryMatchOrder(t *testing.T) {

	s := suggest.Suggest{}

	for _, tt := range queryTests {
		query := tt.query
		s.Commands = tt.commands
		result, err := s.Query(query)
		actual := result.Matches

		if err != nil {
			t.Errorf("TestQueryMatchOrder(%v) returned error %v", tt.query, err)
			continue
		}

		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("TestQueryMatchOrder(%q, %v) = %q, want %q", tt.query, tt.commands, actual, tt.expected)
		}
	}
}
