package flex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
