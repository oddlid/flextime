package flex

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type FlexEntry struct {
	Date   time.Time     `json:"date"`
	Amount time.Duration `json:"amount"`
}

type FlexEntries []*FlexEntry

type Customer struct {
	Name        string      `json:"customer_name"`
	FlexEntries FlexEntries `json:"flex_entries"`
}

type Customers []*Customer

type FlexDB struct {
	FileName  string    `json:"-"`
	Customers Customers `json:"customers"`
}

func NewFlexDB(fileName string) *FlexDB {
	return &FlexDB{
		FileName:  fileName,
		Customers: make(Customers, 0),
	}
}

//func (fes *flexEntries) add(fe flexEntry) {
//	*fes = append(*fes, fe)
//}

func (fes FlexEntries) getTotalFlex() time.Duration {
	var total time.Duration
	for _, e := range fes {
		total += e.Amount
	}
	return total
}

func (fes FlexEntries) Len() int {
	return len(fes)
}

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

//func (c *customer) addFlexEntry(fe flexEntry) bool {
//	return c.setFlexEntry(fe, false)
//}

func (cs Customers) Len() int {
	return len(cs)
}

func (fdb *FlexDB) getCustomer(name string) (*Customer, error) {
	for _, c := range fdb.Customers {
		if strings.EqualFold(name, c.Name) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("no such customer: %q", name)
}

func (fdb *FlexDB) addCustomer(name string) (*Customer, error) {
	if c, err := fdb.getCustomer(name); err == nil {
		return c, fmt.Errorf("customer %s already exists", c.Name)
	}
	c := &Customer{
		Name:        name,
		FlexEntries: make(FlexEntries, 0),
	}
	fdb.Customers = append(fdb.Customers, c)
	return c, nil
}

func (fdb *FlexDB) GetTotalFlexForCustomer(customerName string) (time.Duration, error) {
	c, err := fdb.getCustomer(customerName)
	if err != nil {
		return 0, err
	}
	return c.getTotalFlex(), nil
}

func (fdb *FlexDB) GetTotalFlexForAllCustomers() time.Duration {
	if fdb.Customers.Len() == 0 {
		return 0
	}
	var total time.Duration
	for _, c := range fdb.Customers {
		total += c.getTotalFlex()
	}
	return total
}

func (fdb *FlexDB) SetFlexForCustomer(customerName string, date time.Time, amount time.Duration, overwrite bool) error {
	c, err := fdb.addCustomer(customerName)
	if err != nil {
		log.Debug().Msg(err.Error())
	}
	if !c.setFlexEntry(FlexEntry{Date: date, Amount: amount}, overwrite) {
		return fmt.Errorf("failed to add %v flex on %v for customer %s", amount, date, customerName)
	}
	return nil
}
