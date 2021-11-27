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
	FlexEntrySortByDateAscending = iota
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

func (fes FlexEntries) getTotalFlex() time.Duration {
	var total time.Duration
	for _, e := range fes {
		total += e.Amount
	}
	return total
}

func (fes FlexEntries) Len() int {
	return len(fes)
}

func (fe FlexEntry) Print(w io.Writer) {
	fmt.Fprintf(w, "%s : %v", fe.Date.Format(shortDateFormat), fe.Amount)
}

func (fes FlexEntries) Print(w io.Writer, indentString string, indentLevel int, sortOrder FlexEntrySortOrder) {
	switch sortOrder {
	case FlexEntrySortByDateAscending:
		sort.Sort(FlexEntriesByDate(fes))
	case FlexEntrySortByDateDescending:
		sort.Sort(sort.Reverse(FlexEntriesByDate(fes)))
	case FlexEntrySortByAmountAscending:
		sort.Sort(FlexEntriesByAmount(fes))
	case FlexEntrySortByAmountDescending:
		sort.Sort(sort.Reverse(FlexEntriesByAmount(fes)))
	}

	prefix := strings.Repeat(indentString, indentLevel)

	for _, fe := range fes {
		fmt.Fprintf(w, "%s", prefix)
		fe.Print(w)
		fmt.Fprint(w, "\n")
	}
}

func (febd FlexEntriesByDate) Len() int {
	return len(febd)
}

func (febd FlexEntriesByDate) Swap(i, j int) {
	febd[i], febd[j] = febd[j], febd[i]
}

func (febd FlexEntriesByDate) Less(i, j int) bool {
	return febd[i].Date.Before(febd[j].Date)
}

func (feba FlexEntriesByAmount) Len() int {
	return len(feba)
}

func (feba FlexEntriesByAmount) Swap(i, j int) {
	feba[i], feba[j] = feba[j], feba[i]
}

func (feba FlexEntriesByAmount) Less(i, j int) bool {
	return i < j
}
