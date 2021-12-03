package flex

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

type EntrySortOrder uint8

const (
	ShortDateFormat = "2006-01-02"
)

const (
	EntrySortByDateAscending EntrySortOrder = iota
	EntrySortByDateDescending
	EntrySortByAmountAscending
	EntrySortByAmountDescending
)

// An Entry is the unit for recording flex time +/- for a given date
type Entry struct {
	Date   time.Time     `json:"date"`
	Amount time.Duration `json:"amount"`
}

type Entries []*Entry
type EntriesByDate Entries
type EntriesByAmount Entries

// MatchDate returns true of the date for the two Entries match on year, month and day, false otherwise
func (entry Entry) MatchDate(otherEntry Entry) bool {
	if entry.Date.Year() == otherEntry.Date.Year() &&
		entry.Date.Month() == otherEntry.Date.Month() &&
		entry.Date.Day() == otherEntry.Date.Day() {
		return true
	}
	return false
}

// IndexOf returns the index of the matching entry, if found,
// or -1 if not found.
func (entries Entries) IndexOf(entry Entry) int {
	for idx := range entries {
		if entry.MatchDate(*entries[idx]) {
			return idx
		}
	}
	return -1
}

// GetTotalFlex returns the sum of the Amount fields in all Entries
func (entries Entries) GetTotalFlex() time.Duration {
	var total time.Duration
	for _, entry := range entries {
		total += entry.Amount
	}
	return total
}

// Len returns how many elements in the Entries slice
func (entries Entries) Len() int {
	return len(entries)
}

// Print will print the content of the Entry formatted to the given writer
func (entry Entry) Print(writer io.Writer) {
	fmt.Fprintf(writer, "%s : %v", entry.Date.Format(ShortDateFormat), entry.Amount)
}

// Sort will sort the Entries slice according to the given criteria
func (entries Entries) Sort(sortOrder EntrySortOrder) {
	switch sortOrder {
	case EntrySortByDateAscending:
		sort.Sort(EntriesByDate(entries))
	case EntrySortByDateDescending:
		sort.Sort(sort.Reverse(EntriesByDate(entries)))
	case EntrySortByAmountAscending:
		sort.Sort(EntriesByAmount(entries))
	case EntrySortByAmountDescending:
		sort.Sort(sort.Reverse(EntriesByAmount(entries)))
	}
}

// Print prints each Entry in the Entries slice to the given writer, prefixed by indentString * indentLevel
func (entries Entries) Print(writer io.Writer, indentString string, indentLevel int) {
	prefix := strings.Repeat(indentString, indentLevel)
	for _, entry := range entries {
		fmt.Fprintf(writer, "%s", prefix)
		entry.Print(writer)
		fmt.Fprint(writer, "\n")
	}
}

// PrintSorted first sorts the Entries according to the given criteria, then calls Print with the given parameters
func (entries Entries) PrintSorted(writer io.Writer, indentString string, indentLevel int, sortOrder EntrySortOrder) {
	entries.Sort(sortOrder)
	entries.Print(writer, indentString, indentLevel)
}

func (entriesByDate EntriesByDate) Len() int {
	return len(entriesByDate)
}

func (entriesByDate EntriesByDate) Swap(i, j int) {
	entriesByDate[i], entriesByDate[j] = entriesByDate[j], entriesByDate[i]
}

func (entriesByDate EntriesByDate) Less(i, j int) bool {
	return entriesByDate[i].Date.Before(entriesByDate[j].Date)
}

func (entriesByAmount EntriesByAmount) Len() int {
	return len(entriesByAmount)
}

func (entriesByAmount EntriesByAmount) Swap(i, j int) {
	entriesByAmount[i], entriesByAmount[j] = entriesByAmount[j], entriesByAmount[i]
}

func (entriesByAmount EntriesByAmount) Less(i, j int) bool {
	return entriesByAmount[i].Amount < entriesByAmount[j].Amount
}
