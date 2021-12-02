package flex

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

type FlexEntrySortOrder uint8

const (
	shortDateFormat = "2006-01-02"
)

const (
	FlexEntrySortByDateAscending FlexEntrySortOrder = iota
	FlexEntrySortByDateDescending
	FlexEntrySortByAmountAscending
	FlexEntrySortByAmountDescending
)

type FlexEntry struct {
	Date   time.Time     `json:"date"`
	Amount time.Duration `json:"amount"`
}

type FlexEntries []*FlexEntry
type FlexEntriesByDate FlexEntries
type FlexEntriesByAmount FlexEntries

func (flexEntries FlexEntries) getTotalFlex() time.Duration {
	var total time.Duration
	for _, e := range flexEntries {
		total += e.Amount
	}
	return total
}

func (flexEntries FlexEntries) Len() int {
	return len(flexEntries)
}

func (flexEntry FlexEntry) Print(w io.Writer) {
	fmt.Fprintf(w, "%s : %v", flexEntry.Date.Format(shortDateFormat), flexEntry.Amount)
}

func (flexEntries FlexEntries) Sort(sortOrder FlexEntrySortOrder) {
	switch sortOrder {
	case FlexEntrySortByDateAscending:
		sort.Sort(FlexEntriesByDate(flexEntries))
	case FlexEntrySortByDateDescending:
		sort.Sort(sort.Reverse(FlexEntriesByDate(flexEntries)))
	case FlexEntrySortByAmountAscending:
		sort.Sort(FlexEntriesByAmount(flexEntries))
	case FlexEntrySortByAmountDescending:
		sort.Sort(sort.Reverse(FlexEntriesByAmount(flexEntries)))
	}
}

func (flexEntries FlexEntries) Print(w io.Writer, indentString string, indentLevel int) {
	prefix := strings.Repeat(indentString, indentLevel)
	for _, fe := range flexEntries {
		fmt.Fprintf(w, "%s", prefix)
		fe.Print(w)
		fmt.Fprint(w, "\n")
	}
}

func (flexEntries FlexEntries) PrintSorted(w io.Writer, indentString string, indentLevel int, sortOrder FlexEntrySortOrder) {
	flexEntries.Sort(sortOrder)
	flexEntries.Print(w, indentString, indentLevel)
}

func (flexEntriesByDate FlexEntriesByDate) Len() int {
	return len(flexEntriesByDate)
}

func (flexEntriesByDate FlexEntriesByDate) Swap(i, j int) {
	flexEntriesByDate[i], flexEntriesByDate[j] = flexEntriesByDate[j], flexEntriesByDate[i]
}

func (flexEntriesByDate FlexEntriesByDate) Less(i, j int) bool {
	return flexEntriesByDate[i].Date.Before(flexEntriesByDate[j].Date)
}

func (flexEntriesByAmount FlexEntriesByAmount) Len() int {
	return len(flexEntriesByAmount)
}

func (flexEntriesByAmount FlexEntriesByAmount) Swap(i, j int) {
	flexEntriesByAmount[i], flexEntriesByAmount[j] = flexEntriesByAmount[j], flexEntriesByAmount[i]
}

func (flexEntriesByAmount FlexEntriesByAmount) Less(i, j int) bool {
	return flexEntriesByAmount[i].Amount < flexEntriesByAmount[j].Amount
}
