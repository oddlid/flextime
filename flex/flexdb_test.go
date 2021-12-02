package flex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDBIsEmptyExpectTrue(t *testing.T) {
	db := &DB{}
	assert.True(t, db.IsEmpty())
}

func TestDBIsEmptyExpectFalse(t *testing.T) {
	db := &DB{
		Customers: Customers{
			{Name: "Customer1"},
		},
	}
	assert.False(t, db.IsEmpty())
}

func TestDBAddCustomer(t *testing.T) {
	db := NewDB()

	c1, err1 := db.addCustomer("Customer1")
	assert.NotNil(t, c1)
	assert.NoError(t, err1)

	c2, err2 := db.addCustomer("Customer2")
	assert.NotNil(t, c2)
	assert.NoError(t, err2)

	c3, err3 := db.addCustomer("custOMer1")
	if assert.NotNil(t, c3) {
		assert.Equal(
			t,
			c1,
			c3,
		)
	}
	if assert.Error(t, err3) {
		assert.Equal(
			t,
			"customer already exists: Customer1",
			err3.Error(),
		)
	}

	assert.Equal(
		t,
		2,
		db.Customers.Len(),
	)
}

func TestDBGetTotalFlexForCustomerWhenNoCustomerExists(t *testing.T) {
	db := NewDB()
	totalFlex, err := db.GetTotalFlexForCustomer("customer1")
	assert.Equal(t, time.Duration(0), totalFlex)
	assert.Error(t, err)
}

func TestDBGetTotalFlexForCustomer(t *testing.T) {
	today := time.Now()
	entries := Entries{
		{Date: today.Add(-48 * time.Hour), Amount: 1 * time.Hour},
		{Date: today.Add(-24 * time.Hour), Amount: 30 * time.Minute},
		{Date: today.Add(24 * time.Hour), Amount: -30 * time.Minute},
	}
	db := &DB{
		FileName: "flex.json",
		Customers: Customers{
			{Name: "Customer1", Entries: entries},
		},
	}

	totalFlex, err := db.GetTotalFlexForCustomer("customer1")
	assert.NoError(t, err)
	assert.Equal(
		t,
		1*time.Hour,
		totalFlex,
	)
}

func TestDBGetTotalFlexForAllCustomersWhenNoCustomersExist(t *testing.T) {
	db := NewDB()
	totalFlex := db.GetTotalFlexForAllCustomers()
	assert.Equal(t, time.Duration(0), totalFlex)
}

func TestDBGetTotalFlexForAllCustomers(t *testing.T) {
	today := time.Now()
	db := &DB{
		FileName: "flex.json",
		Customers: Customers{
			{
				Name: "Customer1",
				Entries: Entries{
					{
						Date:   today.Add(-24 * time.Hour),
						Amount: 1 * time.Second,
					},
					{
						Date:   today,
						Amount: 1 * time.Second,
					},
				},
			},
			{
				Name: "Customer2",
				Entries: Entries{
					{
						Date:   today.Add(-24 * time.Hour),
						Amount: 1 * time.Second,
					},
					{
						Date:   today,
						Amount: 1 * time.Second,
					},
				},
			},
		},
	}

	total := db.GetTotalFlexForAllCustomers()
	assert.Equal(
		t,
		4*time.Second,
		total,
	)
}

func TestDBSetFlexForCustomerNoOverwrite(t *testing.T) {
	today := time.Now()
	overwrite := false
	entries := Entries{
		{Date: today, Amount: 1 * time.Hour},
	}
	db := &DB{
		FileName: "flex.json",
		Customers: Customers{
			{Name: "Customer1", Entries: entries},
		},
	}
	err := db.SetFlexForCustomer("customer1", today, 30*time.Minute, overwrite)
	assert.Error(t, err)
}

func TestDBSetFlexForCustomer(t *testing.T) {
	today := time.Now()
	overwrite := true
	entries := Entries{
		{Date: today.Add(-48 * time.Hour), Amount: 1 * time.Hour},
		{Date: today.Add(-24 * time.Hour), Amount: 30 * time.Minute},
		{Date: today.Add(24 * time.Hour), Amount: -30 * time.Minute},
	}
	db := &DB{
		FileName: "flex.json",
		Customers: Customers{
			{Name: "Customer1", Entries: entries},
		},
	}
	err := db.SetFlexForCustomer("customer1", today.Add(24*time.Hour), 30*time.Minute, overwrite)
	assert.NoError(t, err)
	totalFlex, err := db.GetTotalFlexForCustomer("custoMEr1")
	assert.NoError(t, err)
	assert.Equal(
		t,
		2*time.Hour,
		totalFlex,
	)

	err = db.SetFlexForCustomer("customer1", today.Add(48*time.Hour), 30*time.Minute, overwrite)
	assert.NoError(t, err)
	totalFlex, err = db.GetTotalFlexForCustomer("custoMEr1")
	assert.NoError(t, err)
	assert.Equal(
		t,
		150*time.Minute,
		totalFlex,
	)

	err = db.SetFlexForCustomer("customer2", today, 1*time.Second, overwrite)
	assert.NoError(t, err)
	totalFlex, err = db.GetTotalFlexForCustomer("customer2")
	assert.NoError(t, err)
	assert.Equal(
		t,
		2,
		db.Customers.Len(),
	)
	assert.Equal(
		t,
		1*time.Second,
		totalFlex,
	)
}
