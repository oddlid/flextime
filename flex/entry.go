package flex

import (
	"sort"
	"time"
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

// WithinDateRange returns true if the Entry is within the two given dates, inclusive, false otherwise.
func (entry Entry) WithinDateRange(from, to time.Time) bool {
	if entry.Date.Before(from) || entry.Date.After(to) {
		return false
	}
	return true
}

// FilterByDateRange returns a new Entries slice with the entries that are within the given range.
func (entries Entries) FilterByDateRange(from, to time.Time) Entries {
	filteredEntries := make(Entries, 0)
	for _, entry := range entries {
		if entry.WithinDateRange(from, to) {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	return filteredEntries
}

// FilterByNotInDateRange returns a new Entries slice with the entries that are not within the given range.
func (entries Entries) FilterByNotInDateRange(from, to time.Time) Entries {
	filteredEntries := make(Entries, 0)
	for _, entry := range entries {
		if !entry.WithinDateRange(from, to) {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	return filteredEntries
}

func (entries Entries) FirstDate() (*time.Time, error) {
	if entries.Len() == 0 {
		return nil, ErrNoEntries
	}
	var date *time.Time = &entries[0].Date
	for _, entry := range entries {
		if entry.Date.Before(*date) {
			date = &entry.Date
		}
	}
	return date, nil
}

func (entries Entries) LastDate() (*time.Time, error) {
	if entries.Len() == 0 {
		return nil, ErrNoEntries
	}
	var date *time.Time = &entries[entries.Len()-1].Date
	for _, entry := range entries {
		if entry.Date.After(*date) {
			date = &entry.Date
		}
	}
	return date, nil
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

// Delete removes a matching entry from the Entries slice.
// Returns true if match found and removed, false if not.
func (entries *Entries) Delete(entry Entry) bool {
	idx := entries.IndexOf(entry)
	if idx == -1 {
		return false
	}
	// use slow delete, preserving order
	copy((*entries)[idx:], (*entries)[idx+1:])
	(*entries)[len(*entries)-1] = nil
	*entries = (*entries)[:len(*entries)-1]

	return true
}

// DeleteByDate removes an entry with a matching date from the Entries slice.
// Returns true if match found and deleted, false if not.
func (entries *Entries) DeleteByDate(date time.Time) bool {
	return entries.Delete(Entry{Date: date})
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
//func (entry Entry) Print(writer io.Writer) {
//	fmt.Fprintf(writer, "%s : %v", entry.Date.Format(ShortDateFormat), entry.Amount)
//}

// PrintSorted first sorts the Entries according to the given criteria, then calls Print with the given parameters
//func (entries Entries) PrintSorted(writer io.Writer, indentString string, indentLevel int, sortOrder EntrySortOrder) {
//	entries.Sort(sortOrder)
//	entries.Print(writer, indentString, indentLevel)
//}

// Print prints each Entry in the Entries slice to the given writer, prefixed by indentString * indentLevel
//func (entries Entries) Print(writer io.Writer, indentString string, indentLevel int) {
//	prefix := strings.Repeat(indentString, indentLevel)
//	for _, entry := range entries {
//		fmt.Fprintf(writer, "%s", prefix)
//		entry.Print(writer)
//		fmt.Fprint(writer, "\n")
//	}
//}

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
	default:
	}
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
