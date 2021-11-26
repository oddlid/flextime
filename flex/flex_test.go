package flex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFlexDB(t *testing.T) {
	filename := "flex.json"
	db := NewFlexDB(filename)
	if assert.NotNil(t, db, "We should get a *flexDB here") {
		assert.IsType(
			t,
			(*FlexDB)(nil),
			db,
		)
		assert.IsType(
			t,
			Customers{},
			db.Customers,
		)
		assert.Equal(
			t,
			0,
			db.Customers.Len(),
		)
		assert.Equal(
			t,
			filename,
			db.FileName,
		)
	}
}

func TestGetTotalFlexForEntries(t *testing.T) {
	fe := FlexEntries{
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
	fe := FlexEntries{
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

func TestCustomerGetFlexEntry(t *testing.T) {
	today := time.Now()
	fes := FlexEntries{
		{Date: today.Add(-48 * time.Hour), Amount: 1 * time.Hour},
		{Date: today.Add(-24 * time.Hour), Amount: 30 * time.Minute},
		{Date: today.Add(24 * time.Hour), Amount: -30 * time.Minute},
	}
	c := Customer{
		Name:        "MyCompany",
		FlexEntries: fes,
	}

	fe, err := c.getFlexEntry(today)
	assert.Nil(t, fe)
	assert.Error(t, err)

	fe, err = c.getFlexEntry(today.Add(24 * time.Hour))
	assert.NoError(t, err)
	if assert.NotNil(t, fe) {
		assert.Equal(
			t,
			*fes[2],
			*fe,
		)
	}
	fe, err = c.getFlexEntry(today.Add(-24 * time.Hour))
	assert.NoError(t, err)
	if assert.NotNil(t, fe) {
		assert.Equal(
			t,
			*fes[1],
			*fe,
		)
	}
	fe, err = c.getFlexEntry(today.Add(-48 * time.Hour))
	assert.NoError(t, err)
	if assert.NotNil(t, fe) {
		assert.Equal(
			t,
			*fes[0],
			*fe,
		)
	}
}

func TestCustomerSetFlexEntry(t *testing.T) {
	today := time.Now()
	fes := FlexEntries{
		{Date: today.Add(-48 * time.Hour), Amount: 1 * time.Hour},
		{Date: today.Add(-24 * time.Hour), Amount: 30 * time.Minute},
		{Date: today.Add(24 * time.Hour), Amount: -30 * time.Minute},
	}
	c := &Customer{
		Name:        "MyCompany",
		FlexEntries: fes,
	}

	ok := c.setFlexEntry(*fes[0], false)
	assert.False(t, ok)

	ok = c.setFlexEntry(*fes[1], false)
	assert.False(t, ok)

	ok = c.setFlexEntry(*fes[2], false)
	assert.False(t, ok)

	ok = c.setFlexEntry(FlexEntry{Date: today, Amount: 30 * time.Minute}, false)
	assert.True(t, ok)

	ok = c.setFlexEntry(FlexEntry{Date: today, Amount: 15 * time.Minute}, true)
	assert.True(t, ok)

	assert.Equal(
		t,
		75*time.Minute,
		c.getTotalFlex(),
	)
}

func TestFlexDBAddCustomer(t *testing.T) {
	fileName := "flexdb.json"
	db := NewFlexDB(fileName)

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
			"customer Customer1 already exists",
			err3.Error(),
		)
	}

	assert.Equal(
		t,
		2,
		db.Customers.Len(),
	)
}

func TestFlexDBGetTotalFlexForCustomerWhenNoCustomerExists(t *testing.T) {
	db := NewFlexDB("flex.json")
	totalFlex, err := db.GetTotalFlexForCustomer("customer1")
	assert.Equal(t, time.Duration(0), totalFlex)
	assert.Error(t, err)
}

func TestFlexDBGetTotalFlexForCustomer(t *testing.T) {
	today := time.Now()
	fes := FlexEntries{
		{Date: today.Add(-48 * time.Hour), Amount: 1 * time.Hour},
		{Date: today.Add(-24 * time.Hour), Amount: 30 * time.Minute},
		{Date: today.Add(24 * time.Hour), Amount: -30 * time.Minute},
	}
	db := &FlexDB{
		FileName: "flex.json",
		Customers: Customers{
			{Name: "Customer1", FlexEntries: fes},
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

func TestFlexDBGetTotalFlexForAllCustomersWhenNoCustomersExist(t *testing.T) {
	db := NewFlexDB("flex.json")
	totalFlex := db.GetTotalFlexForAllCustomers()
	assert.Equal(t, time.Duration(0), totalFlex)
}

func TestFlexDBGetTotalFlexForAllCustomers(t *testing.T) {
	today := time.Now()
	db := &FlexDB{
		FileName: "flex.json",
		Customers: Customers{
			{
				Name: "Customer1",
				FlexEntries: FlexEntries{
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
				FlexEntries: FlexEntries{
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

func TestFlexDBSetFlexForCustomerNoOverwrite(t *testing.T) {
	today := time.Now()
	overwrite := false
	fes := FlexEntries{
		{Date: today, Amount: 1 * time.Hour},
	}
	db := &FlexDB{
		FileName: "flex.json",
		Customers: Customers{
			{Name: "Customer1", FlexEntries: fes},
		},
	}
	err := db.SetFlexForCustomer("customer1", today, 30*time.Minute, overwrite)
	assert.Error(t, err)
}

func TestFlexDBSetFlexForCustomer(t *testing.T) {
	today := time.Now()
	overwrite := true
	fes := FlexEntries{
		{Date: today.Add(-48 * time.Hour), Amount: 1 * time.Hour},
		{Date: today.Add(-24 * time.Hour), Amount: 30 * time.Minute},
		{Date: today.Add(24 * time.Hour), Amount: -30 * time.Minute},
	}
	db := &FlexDB{
		FileName: "flex.json",
		Customers: Customers{
			{Name: "Customer1", FlexEntries: fes},
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

	//jb, jerr := json.MarshalIndent(db, " ", " ")
	//if jerr == nil {
	//	t.Log(string(jb))
	//}
}
