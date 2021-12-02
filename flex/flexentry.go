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
	FlexEntrySortByDateAscending EntrySortOrder = iota
	FlexEntrySortByDateDescending
	FlexEntrySortByAmountAscending
	FlexEntrySortByAmountDescending
)

type Entry struct {
	Date   time.Time     `json:"date"`
	Amount time.Duration `json:"amount"`
}

type Entries []*Entry
type EntriesByDate Entries
type EntriesByAmount Entries

func (flexEntries Entries) getTotalFlex() time.Duration {
	var total time.Duration
	for _, e := range flexEntries {
		total += e.Amount
	}
	return total
}

func (flexEntries Entries) Len() int {
	return len(flexEntries)
}

func (flexEntry Entry) Print(w io.Writer) {
	fmt.Fprintf(w, "%s : %v", flexEntry.Date.Format(shortDateFormat), flexEntry.Amount)
}

func (flexEntries Entries) Sort(sortOrder EntrySortOrder) {
	switch sortOrder {
	case FlexEntrySortByDateAscending:
		sort.Sort(EntriesByDate(flexEntries))
	case FlexEntrySortByDateDescending:
		sort.Sort(sort.Reverse(EntriesByDate(flexEntries)))
	case FlexEntrySortByAmountAscending:
		sort.Sort(EntriesByAmount(flexEntries))
	case FlexEntrySortByAmountDescending:
		sort.Sort(sort.Reverse(EntriesByAmount(flexEntries)))
	}
}

func (flexEntries Entries) Print(w io.Writer, indentString string, indentLevel int) {
	prefix := strings.Repeat(indentString, indentLevel)
	for _, fe := range flexEntries {
		fmt.Fprintf(w, "%s", prefix)
		fe.Print(w)
		fmt.Fprint(w, "\n")
	}
}

func (flexEntries Entries) PrintSorted(w io.Writer, indentString string, indentLevel int, sortOrder EntrySortOrder) {
	flexEntries.Sort(sortOrder)
	flexEntries.Print(w, indentString, indentLevel)
}

func (flexEntriesByDate EntriesByDate) Len() int {
	return len(flexEntriesByDate)
}

func (flexEntriesByDate EntriesByDate) Swap(i, j int) {
	flexEntriesByDate[i], flexEntriesByDate[j] = flexEntriesByDate[j], flexEntriesByDate[i]
}

func (flexEntriesByDate EntriesByDate) Less(i, j int) bool {
	return flexEntriesByDate[i].Date.Before(flexEntriesByDate[j].Date)
}

func (flexEntriesByAmount EntriesByAmount) Len() int {
	return len(flexEntriesByAmount)
}

func (flexEntriesByAmount EntriesByAmount) Swap(i, j int) {
	flexEntriesByAmount[i], flexEntriesByAmount[j] = flexEntriesByAmount[j], flexEntriesByAmount[i]
}

func (flexEntriesByAmount EntriesByAmount) Less(i, j int) bool {
	return flexEntriesByAmount[i].Amount < flexEntriesByAmount[j].Amount
}
