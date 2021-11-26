package flex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
