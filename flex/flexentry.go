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
	shortDateFormat = "2006-01-02"
)

const (
	EntrySortByDateAscending EntrySortOrder = iota
	EntrySortByDateDescending
	EntrySortByAmountAscending
	EntrySortByAmountDescending
)

type Entry struct {
	Date   time.Time     `json:"date"`
	Amount time.Duration `json:"amount"`
}

type Entries []*Entry
type EntriesByDate Entries
type EntriesByAmount Entries

func (entries Entries) getTotalFlex() time.Duration {
	var total time.Duration
	for _, entry := range entries {
		total += entry.Amount
	}
	return total
}

func (entries Entries) Len() int {
	return len(entries)
}

func (entry Entry) Print(w io.Writer) {
	fmt.Fprintf(w, "%s : %v", entry.Date.Format(shortDateFormat), entry.Amount)
}

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

func (entries Entries) Print(w io.Writer, indentString string, indentLevel int) {
	prefix := strings.Repeat(indentString, indentLevel)
	for _, entry := range entries {
		fmt.Fprintf(w, "%s", prefix)
		entry.Print(w)
		fmt.Fprint(w, "\n")
	}
}

func (entries Entries) PrintSorted(w io.Writer, indentString string, indentLevel int, sortOrder EntrySortOrder) {
	entries.Sort(sortOrder)
	entries.Print(w, indentString, indentLevel)
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
