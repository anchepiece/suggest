// Package suggest is an implementation of a simple command auto-correct feature.
//
// The main application would be to compare a user entered command against a list of known
// available commands. The suggest library can aim to either autocorrect to a very close match,
// provide a single nearest match, provide a list of possible matches, or if there are no
// similar entries, do nothing.
//
// Usage
// import "github.com/anchepiece/suggest"
//
// Example
//		suggester := suggest.Suggest{}
//
//		suggester.Options.SimilarityMinimum = 6
//		suggester.Options.AutocorrectDisabled = false
//
//		query := "proflie"
//		commands := []string{"perfil", "profiel", "profile", "profil", "account"}
//		suggester.Commands = commands
//
//		if result, err := suggester.Query(query); err == nil {
//			if !result.Success() {
//				fmt.Println("No close matches")
//			} else {
//				fmt.Println("Similar matches:", result.Matches) // [profile profil profiel]
//				fmt.Println("Autocorrect:", result.Autocorrect) // profile
//			}
//		}
package suggest
