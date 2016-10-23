// suggest internal testing
package suggest

import (
	"testing"
)

func TestNewSuggest(t *testing.T) {
	opts := &Options{}
	opts.CostDeletion = 40

	s := New(*opts)
	actual := s.Options.getSimilarityMinimum()
	expected := 6
	if actual != expected {
		t.Errorf("New getSimilarityMinimum() = %d, want %d", actual, expected)
	}

	actual = s.Options.getCostDeletion()
	expected = 40
	if actual != expected {
		t.Errorf("New CostDeletion() = %d, want %d", actual, expected)
	}

	opts.CostDeletion = -40
	s = New(*opts)

	actual = s.Options.getCostDeletion()
	expected = defaultCostDeletion
	if actual != expected {
		t.Errorf("New Negative getCostDeletion() = %d, want %d", actual, expected)
	}

}
