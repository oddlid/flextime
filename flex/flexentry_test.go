package flex

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTotalFlexForEntries(t *testing.T) {
	fe := Entries{
		{Date: time.Now(), Amount: 1 * time.Hour},
		{Date: time.Now(), Amount: 30 * time.Minute},
		{Date: time.Now(), Amount: -30 * time.Minute},
	}
	assert.Equal(
		t,
		1*time.Hour,
		fe.getTotalFlex(),
	)
}

func TestFlexEntriesLen(t *testing.T) {
	fe := Entries{
		{Date: time.Now(), Amount: 1 * time.Hour},
		{Date: time.Now(), Amount: 30 * time.Minute},
		{Date: time.Now(), Amount: -30 * time.Minute},
	}
	assert.Equal(
		t,
		3,
		fe.Len(),
	)
}

func TestFlexEntryPrint(t *testing.T) {
	today := time.Now()
	fe := Entry{Date: today, Amount: 1 * time.Hour}
	builder := strings.Builder{}
	fe.Print(&builder)
	assert.Equal(
		t,
		fmt.Sprintf("%s : %v", today.Format(shortDateFormat), fe.Amount),
		builder.String(),
	)
}

func TestSortFlexEntriesByDateAscending(t *testing.T) {
	today := time.Now()
	fe1 := &Entry{Date: today.Add(48 * time.Hour)}
	fe2 := &Entry{Date: today.Add(24 * time.Hour)}
	fe3 := &Entry{Date: today}
	entries := Entries{
		fe1,
		fe2,
		fe3,
	}
	assert.Equal(t, fe1, entries[0])
	assert.Equal(t, fe2, entries[1])
	assert.Equal(t, fe3, entries[2])

	entries.Sort(FlexEntrySortByDateAscending)

	assert.Equal(t, fe1, entries[2])
	assert.Equal(t, fe2, entries[1])
	assert.Equal(t, fe3, entries[0])
}

func TestSortFlexEntriesByDateDescending(t *testing.T) {
	today := time.Now()
	fe1 := &Entry{Date: today.Add(48 * time.Hour)}
	fe2 := &Entry{Date: today.Add(24 * time.Hour)}
	fe3 := &Entry{Date: today}
	entries := Entries{
		fe3,
		fe2,
		fe1,
	}
	assert.Equal(t, fe1, entries[2])
	assert.Equal(t, fe2, entries[1])
	assert.Equal(t, fe3, entries[0])

	entries.Sort(FlexEntrySortByDateDescending)

	assert.Equal(t, fe1, entries[0])
	assert.Equal(t, fe2, entries[1])
	assert.Equal(t, fe3, entries[2])
}

func TestSortFlexEntriesByAmountAscending(t *testing.T) {
	fe1 := &Entry{Amount: 1 * time.Nanosecond}
	fe2 := &Entry{Amount: 2 * time.Nanosecond}
	fe3 := &Entry{Amount: 3 * time.Nanosecond}
	entries := Entries{
		fe3,
		fe2,
		fe1,
	}
	assert.Equal(t, fe1, entries[2])
	assert.Equal(t, fe2, entries[1])
	assert.Equal(t, fe3, entries[0])

	entries.Sort(FlexEntrySortByAmountAscending)

	assert.Equal(t, fe1, entries[0])
	assert.Equal(t, fe2, entries[1])
	assert.Equal(t, fe3, entries[2])
}

func TestSortFlexEntriesByAmountDescending(t *testing.T) {
	fe1 := &Entry{Amount: 1 * time.Nanosecond}
	fe2 := &Entry{Amount: 2 * time.Nanosecond}
	fe3 := &Entry{Amount: 3 * time.Nanosecond}
	entries := Entries{
		fe1,
		fe2,
		fe3,
	}
	assert.Equal(t, fe1, entries[0])
	assert.Equal(t, fe2, entries[1])
	assert.Equal(t, fe3, entries[2])

	entries.Sort(FlexEntrySortByAmountDescending)

	assert.Equal(t, fe1, entries[2])
	assert.Equal(t, fe2, entries[1])
	assert.Equal(t, fe3, entries[0])
}

func TestFlexEntriesPrint(t *testing.T) {
	today := time.Now()
	fe1 := &Entry{Date: today, Amount: 1 * time.Nanosecond}
	fe2 := &Entry{Date: today.Add(24 * time.Hour), Amount: 2 * time.Nanosecond}
	fe3 := &Entry{Date: today.Add(48 * time.Hour), Amount: 3 * time.Nanosecond}
	entries := Entries{
		fe1,
		fe2,
		fe3,
	}

	fe1Str := fmt.Sprintf("%s : %v", fe1.Date.Format(shortDateFormat), fe1.Amount)
	fe2Str := fmt.Sprintf("%s : %v", fe2.Date.Format(shortDateFormat), fe2.Amount)
	fe3Str := fmt.Sprintf("%s : %v", fe3.Date.Format(shortDateFormat), fe3.Amount)

	indentString := " "
	indentLevel := 2
	prefix := strings.Repeat(indentString, indentLevel)

	expected := fmt.Sprintf("%s%s\n%s%s\n%s%s\n", prefix, fe1Str, prefix, fe2Str, prefix, fe3Str)
	builder := strings.Builder{}

	entries.Print(&builder, indentString, indentLevel)

	assert.Equal(
		t,
		expected,
		builder.String(),
	)
}

func TestFlexEntriesPrintSortedByDateDescending(t *testing.T) {
	today := time.Now()
	fe1 := &Entry{Date: today, Amount: 1 * time.Nanosecond}
	fe2 := &Entry{Date: today.Add(24 * time.Hour), Amount: 2 * time.Nanosecond}
	fe3 := &Entry{Date: today.Add(48 * time.Hour), Amount: 3 * time.Nanosecond}
	entries := Entries{
		fe1,
		fe2,
		fe3,
	}

	fe1Str := fmt.Sprintf("%s : %v", fe1.Date.Format(shortDateFormat), fe1.Amount)
	fe2Str := fmt.Sprintf("%s : %v", fe2.Date.Format(shortDateFormat), fe2.Amount)
	fe3Str := fmt.Sprintf("%s : %v", fe3.Date.Format(shortDateFormat), fe3.Amount)

	indentString := " "
	indentLevel := 2
	prefix := strings.Repeat(indentString, indentLevel)

	expected := fmt.Sprintf("%s%s\n%s%s\n%s%s\n", prefix, fe3Str, prefix, fe2Str, prefix, fe1Str)
	builder := strings.Builder{}

	entries.PrintSorted(&builder, indentString, indentLevel, FlexEntrySortByDateDescending)

	assert.Equal(
		t,
		expected,
		builder.String(),
	)
}