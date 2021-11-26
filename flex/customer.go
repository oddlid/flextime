package flex

import (
	"fmt"
	"time"
)

type Customer struct {
	Name        string      `json:"customer_name"`
	FlexEntries FlexEntries `json:"flex_entries"`
}

type Customers []*Customer

func (c Customer) getTotalFlex() time.Duration {
	return c.FlexEntries.getTotalFlex()
}

func (c *Customer) getFlexEntry(date time.Time) (*FlexEntry, error) {
	for _, fe := range c.FlexEntries {
		if fe.Date.Equal(date) {
			return fe, nil
		}
	}
	return nil, fmt.Errorf("no flexentry for date: %v", date)
}

func (c *Customer) setFlexEntry(fe FlexEntry, overwrite bool) bool {
	foundAtIndex := -1
	for idx := range c.FlexEntries {
		if fe.Date.Equal(c.FlexEntries[idx].Date) {
			foundAtIndex = idx
			break
		}
	}
	if foundAtIndex == -1 {
		c.FlexEntries = append(c.FlexEntries, &fe)
		return true
	} else {
		if overwrite {
			c.FlexEntries[foundAtIndex] = &fe
			return true
		}
	}
	return false
}

func (cs Customers) Len() int {
	return len(cs)
}
