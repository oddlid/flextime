package flex

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntryMatchDateExpectFalse(t *testing.T) {
	date1 := time.Date(2021, time.December, 3, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2021, time.December, 4, 0, 0, 0, 0, time.UTC)
	entry1 := Entry{Date: date1}
	entry2 := Entry{Date: date2}
	assert.False(t, entry1.MatchDate(entry2))
}

func TestEntryMatchDateExpectTrue(t *testing.T) {
	date1 := time.Date(2021, time.December, 3, 1, 2, 3, 0, time.UTC)
	date2 := time.Date(2021, time.December, 3, 4, 5, 6, 0, time.UTC)
	entry1 := Entry{Date: date1}
	entry2 := Entry{Date: date2}
	assert.True(t, entry1.MatchDate(entry2))
}

func TestEntryWithinDateRange(t *testing.T) {
	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	twoDaysAgo := yesterday.Add(-24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	twoDaysFromNow := tomorrow.Add(24 * time.Hour)

	tests := []struct {
		name     string
		date     time.Time
		from     time.Time
		to       time.Time
		expected bool
	}{
		{
			name:     "today",
			date:     today,
			from:     yesterday,
			to:       tomorrow,
			expected: true,
		},
		{
			name:     "twoDaysAgo",
			date:     twoDaysAgo,
			from:     yesterday,
			to:       tomorrow,
			expected: false,
		},
		{
			name:     "twoDaysFromNow",
			date:     twoDaysFromNow,
			from:     yesterday,
			to:       tomorrow,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := Entry{Date: tt.date}
			result := entry.WithinDateRange(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("Got %t expected %t", result, tt.expected)
			}
		})
	}
}

func TestEntriesFilterByDateRange(t *testing.T) {
	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	twoDaysAgo := yesterday.Add(-24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	twoDaysFromNow := tomorrow.Add(24 * time.Hour)

	inputEntries := Entries{
		{Date: twoDaysAgo},
		{Date: yesterday},
		{Date: today},
		{Date: tomorrow},
		{Date: twoDaysFromNow},
	}

	resultEntries := inputEntries.FilterByDateRange(yesterday, tomorrow)
	assert.NotNil(t, resultEntries)
	assert.Equal(
		t,
		3,
		resultEntries.Len(),
	)
	assert.Equal(
		t,
		inputEntries[1],
		resultEntries[0],
	)
	assert.Equal(
		t,
		inputEntries[2],
		resultEntries[1],
	)
	assert.Equal(
		t,
		inputEntries[3],
		resultEntries[2],
	)
}

func TestEntriesFilterByNotInDateRange(t *testing.T) {
	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	twoDaysAgo := yesterday.Add(-24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	twoDaysFromNow := tomorrow.Add(24 * time.Hour)

	inputEntries := Entries{
		{Date: twoDaysAgo},
		{Date: yesterday},
		{Date: today},
		{Date: tomorrow},
		{Date: twoDaysFromNow},
	}

	resultEntries := inputEntries.FilterByNotInDateRange(yesterday, tomorrow)
	assert.NotNil(t, resultEntries)
	assert.Equal(
		t,
		2,
		resultEntries.Len(),
	)
	assert.Equal(
		t,
		inputEntries[0],
		resultEntries[0],
	)
	assert.Equal(
		t,
		inputEntries[4],
		resultEntries[1],
	)
}

func TestEntriesDelete(t *testing.T) {
}

func TestEntriesDeleteByDate(t *testing.T) {
}

func TestGetTotalFlexForEntries(t *testing.T) {
	entries := Entries{
		{Date: time.Now(), Amount: 1 * time.Hour},
		{Date: time.Now(), Amount: 30 * time.Minute},
		{Date: time.Now(), Amount: -30 * time.Minute},
	}
	assert.Equal(
		t,
		1*time.Hour,
		entries.GetTotalFlex(),
	)
}

func TestEntriesLen(t *testing.T) {
	entries := Entries{
		{Date: time.Now(), Amount: 1 * time.Hour},
		{Date: time.Now(), Amount: 30 * time.Minute},
		{Date: time.Now(), Amount: -30 * time.Minute},
	}
	assert.Equal(
		t,
		3,
		entries.Len(),
	)
}

func TestEntryPrint(t *testing.T) {
	today := time.Now()
	entry := Entry{Date: today, Amount: 1 * time.Hour}
	builder := strings.Builder{}
	entry.Print(&builder)
	assert.Equal(
		t,
		fmt.Sprintf("%s : %v", today.Format(ShortDateFormat), entry.Amount),
		builder.String(),
	)
}

func TestSortEntriesByDateAscending(t *testing.T) {
	today := time.Now()
	entry1 := &Entry{Date: today.Add(48 * time.Hour)}
	entry2 := &Entry{Date: today.Add(24 * time.Hour)}
	entry3 := &Entry{Date: today}
	entries := Entries{
		entry1,
		entry2,
		entry3,
	}
	assert.Equal(t, entry1, entries[0])
	assert.Equal(t, entry2, entries[1])
	assert.Equal(t, entry3, entries[2])

	entries.Sort(EntrySortByDateAscending)

	assert.Equal(t, entry1, entries[2])
	assert.Equal(t, entry2, entries[1])
	assert.Equal(t, entry3, entries[0])
}

func TestSortEntriesByDateDescending(t *testing.T) {
	today := time.Now()
	entry1 := &Entry{Date: today.Add(48 * time.Hour)}
	entry2 := &Entry{Date: today.Add(24 * time.Hour)}
	entry3 := &Entry{Date: today}
	entries := Entries{
		entry3,
		entry2,
		entry1,
	}
	assert.Equal(t, entry1, entries[2])
	assert.Equal(t, entry2, entries[1])
	assert.Equal(t, entry3, entries[0])

	entries.Sort(EntrySortByDateDescending)

	assert.Equal(t, entry1, entries[0])
	assert.Equal(t, entry2, entries[1])
	assert.Equal(t, entry3, entries[2])
}

func TestSortEntriesByAmountAscending(t *testing.T) {
	entry1 := &Entry{Amount: 1 * time.Nanosecond}
	entry2 := &Entry{Amount: 2 * time.Nanosecond}
	entry3 := &Entry{Amount: 3 * time.Nanosecond}
	entries := Entries{
		entry3,
		entry2,
		entry1,
	}
	assert.Equal(t, entry1, entries[2])
	assert.Equal(t, entry2, entries[1])
	assert.Equal(t, entry3, entries[0])

	entries.Sort(EntrySortByAmountAscending)

	assert.Equal(t, entry1, entries[0])
	assert.Equal(t, entry2, entries[1])
	assert.Equal(t, entry3, entries[2])
}

func TestSortEntriesByAmountDescending(t *testing.T) {
	entry1 := &Entry{Amount: 1 * time.Nanosecond}
	entry2 := &Entry{Amount: 2 * time.Nanosecond}
	entry3 := &Entry{Amount: 3 * time.Nanosecond}
	entries := Entries{
		entry1,
		entry2,
		entry3,
	}
	assert.Equal(t, entry1, entries[0])
	assert.Equal(t, entry2, entries[1])
	assert.Equal(t, entry3, entries[2])

	entries.Sort(EntrySortByAmountDescending)

	assert.Equal(t, entry1, entries[2])
	assert.Equal(t, entry2, entries[1])
	assert.Equal(t, entry3, entries[0])
}

func TestEntriesPrint(t *testing.T) {
	today := time.Now()
	entry1 := &Entry{Date: today, Amount: 1 * time.Nanosecond}
	entry2 := &Entry{Date: today.Add(24 * time.Hour), Amount: 2 * time.Nanosecond}
	entry3 := &Entry{Date: today.Add(48 * time.Hour), Amount: 3 * time.Nanosecond}
	entries := Entries{
		entry1,
		entry2,
		entry3,
	}

	entry1Str := fmt.Sprintf("%s : %v", entry1.Date.Format(ShortDateFormat), entry1.Amount)
	entry2Str := fmt.Sprintf("%s : %v", entry2.Date.Format(ShortDateFormat), entry2.Amount)
	entry3Str := fmt.Sprintf("%s : %v", entry3.Date.Format(ShortDateFormat), entry3.Amount)

	indentString := " "
	indentLevel := 2
	prefix := strings.Repeat(indentString, indentLevel)

	expected := fmt.Sprintf("%s%s\n%s%s\n%s%s\n", prefix, entry1Str, prefix, entry2Str, prefix, entry3Str)
	builder := strings.Builder{}

	entries.Print(&builder, indentString, indentLevel)

	assert.Equal(
		t,
		expected,
		builder.String(),
	)
}

func TestEntriesPrintSortedByDateDescending(t *testing.T) {
	today := time.Now()
	entry1 := &Entry{Date: today, Amount: 1 * time.Nanosecond}
	entry2 := &Entry{Date: today.Add(24 * time.Hour), Amount: 2 * time.Nanosecond}
	entry3 := &Entry{Date: today.Add(48 * time.Hour), Amount: 3 * time.Nanosecond}
	entries := Entries{
		entry1,
		entry2,
		entry3,
	}

	entry1Str := fmt.Sprintf("%s : %v", entry1.Date.Format(ShortDateFormat), entry1.Amount)
	entry2Str := fmt.Sprintf("%s : %v", entry2.Date.Format(ShortDateFormat), entry2.Amount)
	entry3Str := fmt.Sprintf("%s : %v", entry3.Date.Format(ShortDateFormat), entry3.Amount)

	indentString := " "
	indentLevel := 2
	prefix := strings.Repeat(indentString, indentLevel)

	expected := fmt.Sprintf("%s%s\n%s%s\n%s%s\n", prefix, entry3Str, prefix, entry2Str, prefix, entry1Str)
	builder := strings.Builder{}

	entries.PrintSorted(&builder, indentString, indentLevel, EntrySortByDateDescending)

	assert.Equal(
		t,
		expected,
		builder.String(),
	)
}
