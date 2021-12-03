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

func (db *DB) IsEmpty() bool {
	if db.Customers == nil || db.Customers.Len() == 0 {
		return true
	}
	return false
}

func (db *DB) getCustomer(name string) (*Customer, error) {
	for _, customer := range db.Customers {
		if strings.EqualFold(name, customer.Name) {
			return customer, nil
		}
	}
	return nil, fmt.Errorf("%w: %q", ErrNoSuchCustomer, name)
}

func (db *DB) addCustomer(name string) (*Customer, error) {
	if customer, err := db.getCustomer(name); err == nil {
		return customer, fmt.Errorf("%w: %s", ErrCustomerExists, customer.Name)
	}
	customer := &Customer{
		Name:    name,
		Entries: make(Entries, 0),
	}
	db.Customers = append(db.Customers, customer)
	return customer, nil
}

func (db *DB) GetTotalFlexForCustomer(customerName string) (time.Duration, error) {
	customer, err := db.getCustomer(customerName)
	if err != nil {
		return 0, err
	}
	return customer.GetTotalFlex(), nil
}

func (db *DB) GetTotalFlexForAllCustomers() time.Duration {
	if db.Customers.Len() == 0 {
		return 0
	}
	var total time.Duration
	for _, customer := range db.Customers {
		total += customer.GetTotalFlex()
	}
	return total
}

func (db *DB) SetFlexForCustomer(customerName string, date time.Time, amount time.Duration, overwrite bool) error {
	customer, err := db.addCustomer(customerName)
	if err != nil {
		log.Debug().Msg(err.Error())
	}
	if !customer.SetEntry(Entry{Date: date, Amount: amount}, overwrite) {
		return fmt.Errorf("failed to add %v flex on %v for customer %s", amount, date, customerName)
	}
	return nil
}
