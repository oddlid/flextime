package flex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomerGetTotalFlexWhenEntriesIsNil(t *testing.T) {
	customer := Customer{Name: "Customer1"}
	totalFlex := customer.GetTotalFlex()
	assert.Equal(
		t,
		time.Duration(0),
		totalFlex,
	)
}

func TestCustomerGetEntryWhenEntriesIsNil(t *testing.T) {
	customer := Customer{Name: "Customer1"}
	entry, err := customer.GetEntry(time.Now())
	assert.Nil(t, entry)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNoEntries)
	}
}

func TestCustomerGetEntryWhenEntriesLenIs0(t *testing.T) {
	customer := Customer{
		Name:    "Customer1",
		Entries: make(Entries, 0),
	}
	entry, err := customer.GetEntry(time.Now())
	assert.Nil(t, entry)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNoEntries)
	}
}

func TestCustomerGetEntryWhenEntryNotFound(t *testing.T) {
	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	twoDaysAgo := today.Add(-48 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	customer := Customer{
		Name: "Customer1",
		Entries: Entries{
			{Date: tomorrow},
			{Date: yesterday},
			{Date: twoDaysAgo},
		},
	}
	entry, err := customer.GetEntry(today)
	assert.Nil(t, entry)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNoEntry)
	}
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

	entry, err := customer.GetEntry(today)
	assert.Nil(t, entry)
	assert.Error(t, err)

	entry, err = customer.GetEntry(today.Add(24 * time.Hour))
	assert.NoError(t, err)
	if assert.NotNil(t, entry) {
		assert.Equal(
			t,
			*entries[2],
			*entry,
		)
	}
	entry, err = customer.GetEntry(today.Add(-24 * time.Hour))
	assert.NoError(t, err)
	if assert.NotNil(t, entry) {
		assert.Equal(
			t,
			*entries[1],
			*entry,
		)
	}
	entry, err = customer.GetEntry(today.Add(-48 * time.Hour))
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

	ok := customer.SetEntry(*entries[0], false)
	assert.False(t, ok)

	ok = customer.SetEntry(*entries[1], false)
	assert.False(t, ok)

	ok = customer.SetEntry(*entries[2], false)
	assert.False(t, ok)

	ok = customer.SetEntry(Entry{Date: today, Amount: 30 * time.Minute}, false)
	assert.True(t, ok)

	ok = customer.SetEntry(Entry{Date: today, Amount: 15 * time.Minute}, true)
	assert.True(t, ok)

	assert.Equal(
		t,
		75*time.Minute,
		customer.GetTotalFlex(),
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

func TestCustomersLongestName(t *testing.T) {
	customers := Customers{
		{Name: "123"},
		{Name: "1234"},
		{Name: "12345"},
	}
	assert.Equal(t, 5, customers.LongestName())
}

func TestCustomersIndexOf(t *testing.T) {
	customers := Customers{
		{Name: "Customer1"},
		{Name: "Customer2"},
		{Name: "Customer3"},
	}

	assert.Equal(t, 1, customers.IndexOf(*customers[1]))
	assert.Equal(t, -1, customers.IndexOf(Customer{Name: "NoSuchCustomer"}))
}

func TestCustomersDeleteExpectFalse(t *testing.T) {
	customers := Customers{
		{Name: "Customer1"},
		{Name: "Customer2"},
		{Name: "Customer3"},
	}
	customerNotInList := Customer{Name: "CustomerNotInList"}

	assert.False(t, customers.Delete(customerNotInList))
	assert.Equal(t, 3, len(customers))
}

func TestCustomersDeleteExpectTrue(t *testing.T) {
	customers := Customers{
		{Name: "Customer1"},
		{Name: "Customer2"},
		{Name: "Customer3"},
	}

	assert.True(t, customers.Delete(*customers[1]))
	assert.Equal(t, 2, len(customers))
}
