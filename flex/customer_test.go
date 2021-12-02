package flex

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomerGetTotalFlexWhenEntriesIsNil(t *testing.T) {
	customer := Customer{Name: "Customer1"}
	totalFlex := customer.getTotalFlex()
	assert.Equal(
		t,
		time.Duration(0),
		totalFlex,
	)
}

func TestCustomerGetEntry(t *testing.T) {
	today := time.Now()
	entries := Entries{
		{Date: today.Add(-48 * time.Hour), Amount: 1 * time.Hour},
		{Date: today.Add(-24 * time.Hour), Amount: 30 * time.Minute},
		{Date: today.Add(24 * time.Hour), Amount: -30 * time.Minute},
	}
	customer := Customer{
		Name:    "MyCompany",
		Entries: entries,
	}

	entry, err := customer.getEntry(today)
	assert.Nil(t, entry)
	assert.Error(t, err)

	entry, err = customer.getEntry(today.Add(24 * time.Hour))
	assert.NoError(t, err)
	if assert.NotNil(t, entry) {
		assert.Equal(
			t,
			*entries[2],
			*entry,
		)
	}
	entry, err = customer.getEntry(today.Add(-24 * time.Hour))
	assert.NoError(t, err)
	if assert.NotNil(t, entry) {
		assert.Equal(
			t,
			*entries[1],
			*entry,
		)
	}
	entry, err = customer.getEntry(today.Add(-48 * time.Hour))
	assert.NoError(t, err)
	if assert.NotNil(t, entry) {
		assert.Equal(
			t,
			*entries[0],
			*entry,
		)
	}
}

func TestCustomerSetEntry(t *testing.T) {
	today := time.Now()
	entries := Entries{
		{Date: today.Add(-48 * time.Hour), Amount: 1 * time.Hour},
		{Date: today.Add(-24 * time.Hour), Amount: 30 * time.Minute},
		{Date: today.Add(24 * time.Hour), Amount: -30 * time.Minute},
	}
	customer := &Customer{
		Name:    "MyCompany",
		Entries: entries,
	}

	ok := customer.setEntry(*entries[0], false)
	assert.False(t, ok)

	ok = customer.setEntry(*entries[1], false)
	assert.False(t, ok)

	ok = customer.setEntry(*entries[2], false)
	assert.False(t, ok)

	ok = customer.setEntry(Entry{Date: today, Amount: 30 * time.Minute}, false)
	assert.True(t, ok)

	ok = customer.setEntry(Entry{Date: today, Amount: 15 * time.Minute}, true)
	assert.True(t, ok)

	assert.Equal(
		t,
		75*time.Minute,
		customer.getTotalFlex(),
	)
}

func TestCustomersSortAscending(t *testing.T) {
	c1 := &Customer{Name: "CustomerA"}
	c2 := &Customer{Name: "CustomerB"}
	c3 := &Customer{Name: "CustomerC"}
	c4 := &Customer{Name: "Customerc"}
	customers := Customers{
		c4,
		c3,
		c2,
		c1,
	}
	assert.Equal(t, c4, customers[0])
	assert.Equal(t, c3, customers[1])
	assert.Equal(t, c2, customers[2])
	assert.Equal(t, c1, customers[3])

	customers.Sort(CustomerSortByNameAscending)

	assert.Equal(t, c1, customers[0])
	assert.Equal(t, c2, customers[1])
	assert.Equal(t, c3, customers[2])
	assert.Equal(t, c4, customers[3])
}

func TestCustomersSortDescending(t *testing.T) {
	c1 := &Customer{Name: "CustomerA"}
	c2 := &Customer{Name: "CustomerB"}
	c3 := &Customer{Name: "CustomerC"}
	c4 := &Customer{Name: "Customerc"}
	customers := Customers{
		c1,
		c2,
		c3,
		c4,
	}
	assert.Equal(t, c1, customers[0])
	assert.Equal(t, c2, customers[1])
	assert.Equal(t, c3, customers[2])
	assert.Equal(t, c4, customers[3])

	customers.Sort(CustomerSortByNameDescending)

	assert.Equal(t, c1, customers[3])
	assert.Equal(t, c2, customers[2])
	assert.Equal(t, c3, customers[1])
	assert.Equal(t, c4, customers[0])
}

func TestCustomerPrintWhenEntriesIsNil(t *testing.T) {
	indentString := " "
	indentLevel := 2
	name := "Customer1"
	customer := Customer{Name: name}
	builder := strings.Builder{}
	expected := fmt.Sprintf("%s%s:\n", strings.Repeat(indentString, indentLevel), name)

	customer.Print(&builder, indentString, indentLevel)
	assert.Equal(
		t,
		expected,
		builder.String(),
	)
}

func TestCustomerPrintWithEntries(t *testing.T) {
	indentString := " "
	indentLevel := 2
	name := "Customer1"
	today := time.Now()
	amount := 1 * time.Nanosecond
	customer := Customer{Name: name, Entries: Entries{{Date: today, Amount: amount}}}
	builder := strings.Builder{}
	expected := fmt.Sprintf(
		"%s%s:\n%s%s : %v\n",
		strings.Repeat(indentString, indentLevel),
		name,
		strings.Repeat(indentString, indentLevel+1),
		today.Format(shortDateFormat),
		amount,
	)

	customer.Print(&builder, indentString, indentLevel)
	assert.Equal(
		t,
		expected,
		builder.String(),
	)
}
