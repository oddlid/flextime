package flex

import (
	"testing"

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
