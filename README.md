# suggest [![Doc Status](https://godoc.org/github.com/anchepiece/suggest?status.png)](https://godoc.org/github.com/anchepiece/suggest)

## Introduction

Go implementation of a simple command auto-correct  feature. Inspired by "Did you mean?" search suggestions.

The main application would be to compare a user entered command against a list of known
available commands. The suggest library can aim to either auto-correct to a very close match, 
provide a single nearest match, provide a list of possible matches, or if there are no 
similar entries, do nothing.


## Design

Suggest should provide reasonable defaults to return predictable behavior.


## Example

	go get github.com/anchepiece/suggest

```go
import "github.com/anchepiece/suggest"

func main() {

	suggester := suggest.Suggest{}

	query := "fgerp"
	commands := []string{"cat", "mkdir", "fgrep", "history"}

	suggester.Commands = commands
	if match, err := suggester.Autocorrect("mkdri"); err == nil {
		fmt.Println("Autocorrected to:", match) // "mkdir"
	}

	// Alternate autocorrect usage pattern
	match, _ := suggester.AutocorrectAgainst(query, commands)
	if match != "" {
		fmt.Println("Autocorrected to:", match) // "fgrep"
	}

	// Alternate usage pattern
	query = "println"
	commands = []string{"Fprint", "Fprintf", "Fprintln", "Sprintf", "Print", "Printf", "Println"}
	suggester.Options.SimilarityMinimum = 8
	
	fmt.Printf("Searching %v for %s\n", query, commands)

	if result, err := suggester.QueryAgainst(query, commands); err == nil {
		if !result.Success() {
			fmt.Println("No close matches")

		} else {
			fmt.Println("Similar matches:", result.Matches) 
			// [Println Fprintln]

			fmt.Println("Autocorrect:", result.Autocorrect) 
			// Println
		}
	}
}
```

## GoDoc

[GoDoc](https://godoc.org/github.com/anchepiece/suggest)

## TODO
- [ ] Boost prefix
- [ ] Complete testing 
- [ ] Web backend service
- [ ] Command-line help usage example 
- [ ] Ideas as to other applications

## License

This library is under the [MIT License](http://opensource.org/licenses/MIT)

