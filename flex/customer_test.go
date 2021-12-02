package flex

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomerGetTotalFlexWhenFlexEntriesIsNil(t *testing.T) {
	customer := Customer{Name: "Customer1"}
	totalFlex := customer.getTotalFlex()
	assert.Equal(
		t,
		time.Duration(0),
		totalFlex,
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

func TestCustomerPrintWhenFlexEntriesIsNil(t *testing.T) {
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

func TestCustomerPrintWithFlexEntries(t *testing.T) {
	indentString := " "
	indentLevel := 2
	name := "Customer1"
	today := time.Now()
	amount := 1 * time.Nanosecond
	customer := Customer{Name: name, FlexEntries: FlexEntries{{Date: today, Amount: amount}}}
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
