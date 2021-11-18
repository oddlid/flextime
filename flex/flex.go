package flex

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type flexEntry struct {
	Date   time.Time     `json:"date"`
	Amount time.Duration `json:"amount"`
}

type flexEntries []*flexEntry

type customer struct {
	Name        string      `json:"customer_name"`
	FlexEntries flexEntries `json:"flex_entries"`
}

type customers []*customer

type flexDB struct {
	FileName  string    `json:"-"`
	Customers customers `json:"customers"`
}

func NewFlexDB(fileName string) *flexDB {
	return &flexDB{
		FileName:  fileName,
		Customers: make(customers, 0),
	}
}

//func (fes *flexEntries) add(fe flexEntry) {
//	*fes = append(*fes, fe)
//}

func (fes flexEntries) getTotalFlex() time.Duration {
	var total time.Duration
	for _, e := range fes {
		total += e.Amount
	}
	return total
}

func (fes flexEntries) Len() int {
	return len(fes)
}

func (c customer) getTotalFlex() time.Duration {
	return c.FlexEntries.getTotalFlex()
}

func (c *customer) getFlexEntry(date time.Time) (*flexEntry, error) {
	for _, fe := range c.FlexEntries {
		if fe.Date.Equal(date) {
			return fe, nil
		}
	}
	return nil, fmt.Errorf("no flexentry for date: %v", date)
}

func (c *customer) setFlexEntry(fe flexEntry, overwrite bool) bool {
	foundAtIndex := -1
	for idx := range c.FlexEntries {
		if fe.Date.Equal(c.FlexEntries[idx].Date) {
			foundAtIndex = idx
			break
		}
	}
	if foundAtIndex == -1 {
		//log.Debug().
		//	Str("customer_name", c.Name).
		//	Time("date", fe.Date).
		//	Dur("amount", fe.Amount).
		//	Msg("Adding new FlexEntry")
		c.FlexEntries = append(c.FlexEntries, &fe)
		return true
	} else {
		if overwrite {
			//log.Debug().
			//	Str("customer_name", c.Name).
			//	Time("date", fe.Date).
			//	Dur("amount", fe.Amount).
			//	Int("index", foundAtIndex).
			//	Msg("Overwriting existing FlexEntry")
			c.FlexEntries[foundAtIndex] = &fe
			return true
		}
	}
	return false
}

//func (c *customer) addFlexEntry(fe flexEntry) bool {
//	return c.setFlexEntry(fe, false)
//}

func (cs customers) Len() int {
	return len(cs)
}

func (fdb *flexDB) getCustomer(name string) (*customer, error) {
	for _, c := range fdb.Customers {
		if strings.EqualFold(name, c.Name) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("no such customer: %q", name)
}

func (fdb *flexDB) addCustomer(name string) *customer {
	if _, err := fdb.getCustomer(name); err == nil {
		//log.Debug().
		//	Str("customer_name", name).
		//	Msg("Customer already exists")
		return nil // customer exists
	}
	c := &customer{
		Name:        name,
		FlexEntries: make(flexEntries, 0),
	}
	fdb.Customers = append(fdb.Customers, c)
	return c
}

func (fdb *flexDB) GetTotalFlexForCustomer(customerName string) (time.Duration, error) {
	c, err := fdb.getCustomer(customerName)
	if err != nil {
		return 0, err
	}
	return c.getTotalFlex(), nil
}

func (fdb *flexDB) GetTotalFlexForAllCustomers() time.Duration {
	if fdb.Customers.Len() == 0 {
		return 0
	}
	var total time.Duration
	for _, c := range fdb.Customers {
		total += c.getTotalFlex()
	}
	return total
}

func (fdb *flexDB) SetFlexForCustomer(customerName string, date time.Time, amount time.Duration) {
	c, err := fdb.getCustomer(customerName)
	if err != nil {
		//log.Debug().
		//	Str("customer_name", customerName).
		//	Msg("No such customer, creating...")
		c = fdb.addCustomer(customerName)
	}
	if !c.setFlexEntry(flexEntry{Date: date, Amount: amount}, true) {
		log.Error().
			Str("customer_name", customerName).
			Msg("Failed to add flex")
	}
}
