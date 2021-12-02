package flex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlexDBAddCustomer(t *testing.T) {
	db := NewFlexDB()

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

func TestFlexDBGetTotalFlexForCustomerWhenNoCustomerExists(t *testing.T) {
	db := NewFlexDB()
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
	db := NewFlexDB()
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
}

//func TestFlexDBSaveWithNoFilename(t *testing.T) {
//	db := &FlexDB{}
//	err := db.Save()
//	if assert.Error(t, err) {
//		assert.Equal(
//			t,
//			"FlexDB has no filename to save to",
//			err.Error(),
//		)
//	}
//}
//
//func TestFlexDBSaveWithInvalidFilename(t *testing.T) {
//	filename := "/invalid/path/to/file.json"
//	db := &FlexDB{FileName: filename}
//	err := db.Save()
//	if assert.Error(t, err) {
//		assert.Equal(
//			t,
//			fmt.Sprintf("open %s: no such file or directory", filename),
//			err.Error(),
//		)
//	}
//}
//
//func TestFlexDBSave(t *testing.T) {
//	file, err := os.CreateTemp("", "flextime")
//	if err != nil {
//		t.Errorf("Failed to create temp file: %v", err)
//		return
//	}
//	if err = file.Close(); err != nil {
//		t.Errorf("%v", err)
//	}
//	defer os.Remove(file.Name())
//	db := &FlexDB{FileName: file.Name()}
//	err = db.Save()
//	assert.NoError(t, err)
//}
