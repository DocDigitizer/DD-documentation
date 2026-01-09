package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// JSONOutput controls whether output should be in JSON format
var JSONOutput bool

// PrintTable prints data in table format
func PrintTable(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	printRow(headers, widths)

	// Print separator
	sep := make([]string, len(headers))
	for i, w := range widths {
		sep[i] = strings.Repeat("-", w)
	}
	printRow(sep, widths)

	// Print rows
	for _, row := range rows {
		printRow(row, widths)
	}
}

func printRow(cells []string, widths []int) {
	parts := make([]string, len(widths))
	for i := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		parts[i] = fmt.Sprintf("%-*s", widths[i], cell)
	}
	fmt.Println(strings.Join(parts, "  "))
}

// PrintKeyValue prints key-value pairs in a formatted way
func PrintKeyValue(data map[string]string) {
	maxKeyLen := 0
	for k := range data {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	for k, v := range data {
		fmt.Printf("%-*s  %s\n", maxKeyLen, k+":", v)
	}
}

// PrintJSON prints data as JSON
func PrintJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	fmt.Println(msg)
}

// PrintError prints an error message to stderr
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// Print handles printing based on JSONOutput flag
func Print(v interface{}, tableHeaders []string, tableRows [][]string) {
	if JSONOutput {
		PrintJSON(v)
	} else {
		PrintTable(tableHeaders, tableRows)
	}
}

// PrintSingle handles printing a single item based on JSONOutput flag
func PrintSingle(v interface{}, kvData map[string]string) {
	if JSONOutput {
		PrintJSON(v)
	} else {
		PrintKeyValue(kvData)
	}
}

// Truncate truncates a string to maxLen characters
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// PtrString returns the value of a string pointer or a default
func PtrString(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

// BoolString returns "Yes" or "No" for a boolean
func BoolString(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
