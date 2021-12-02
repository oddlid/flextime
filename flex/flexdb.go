package flex

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	ErrNoSuchCustomer = errors.New("no such customer")
	ErrCustomerExists = errors.New("customer already exists")
)

type DB struct {
	FileName  string    `json:"-"`
	Customers Customers `json:"customers"`
}

func (fdb *DB) IsEmpty() bool {
	dbStr := fmt.Sprintf("%#v", fdb)
	log.Debug().Str("db", dbStr).Msg("DB passed to IsEmpty")
	if fdb.Customers == nil || fdb.Customers.Len() == 0 {
		return true
	}
	return false
}

func (fdb *DB) getCustomer(name string) (*Customer, error) {
	for _, c := range fdb.Customers {
		if strings.EqualFold(name, c.Name) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("%w: %q", ErrNoSuchCustomer, name)
}

func (fdb *DB) addCustomer(name string) (*Customer, error) {
	if c, err := fdb.getCustomer(name); err == nil {
		return c, fmt.Errorf("%w: %s", ErrCustomerExists, c.Name)
	}
	c := &Customer{
		Name:        name,
		FlexEntries: make(Entries, 0),
	}
	fdb.Customers = append(fdb.Customers, c)
	return c, nil
}

func (fdb *DB) GetTotalFlexForCustomer(customerName string) (time.Duration, error) {
	c, err := fdb.getCustomer(customerName)
	if err != nil {
		return 0, err
	}
	return c.getTotalFlex(), nil
}

func (fdb *DB) GetTotalFlexForAllCustomers() time.Duration {
	if fdb.Customers.Len() == 0 {
		return 0
	}
	var total time.Duration
	for _, c := range fdb.Customers {
		total += c.getTotalFlex()
	}
	return total
}

func (fdb *DB) SetFlexForCustomer(customerName string, date time.Time, amount time.Duration, overwrite bool) error {
	c, err := fdb.addCustomer(customerName)
	if err != nil {
		log.Debug().Msg(err.Error())
	}
	if !c.setFlexEntry(Entry{Date: date, Amount: amount}, overwrite) {
		return fmt.Errorf("failed to add %v flex on %v for customer %s", amount, date, customerName)
	}
	return nil
}
