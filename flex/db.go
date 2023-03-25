package flex

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type DB struct {
	FileName  string    `json:"-"`
	Customers Customers `json:"customers"`
}

// IsEmpty returns true if its Customers field is nil, or its length i 0, false otherwise.
func (db *DB) IsEmpty() bool {
	if db.Customers == nil || db.Customers.Len() == 0 {
		return true
	}
	return false
}

// GetCustomer returns a pointer to the Customer struct with a matching name (case insensitive),
// or nil and an error if not found.
func (db *DB) GetCustomer(name string) (*Customer, error) {
	for _, customer := range db.Customers {
		if strings.EqualFold(name, customer.Name) {
			return customer, nil
		}
	}
	return nil, fmt.Errorf("%w: %q", ErrNoSuchCustomer, name)
}

// GetDefaultCustomer is a convenience method that will either return
// a customer with the name "default", if found, or the first.
// If no customers exist, a default is created and returned.
func (db *DB) GetDefaultCustomer() *Customer {
	if db.Customers == nil {
		db.Customers = make(Customers, 0)
	}
	defaultCustomer := Customer{Name: DefaultCustomerName}
	if db.Customers.Len() == 0 {
		db.Customers = append(db.Customers, &defaultCustomer)
		return &defaultCustomer
	}
	customer, err := db.GetCustomer(DefaultCustomerName)
	if err == nil {
		return customer
	}
	return db.Customers[0]
}

// AddCustomer will add a new Customer object with the given name to the DB if it doesn't exist.
// It will return a pointer to the Customer struct if found or added. The difference will be that if found,
// it will also return an error stating that the customer already exists.
func (db *DB) AddCustomer(name string) (*Customer, error) {
	if customer, err := db.GetCustomer(name); err == nil {
		return customer, fmt.Errorf("%w: %s", ErrCustomerExists, customer.Name)
	}
	customer := &Customer{
		Name:    name,
		Entries: make(Entries, 0),
	}
	db.Customers = append(db.Customers, customer)
	return customer, nil
}

// GetTotalFlexForCustomer returns the total flex time for the given Customer if found,
// or an error if not found.
func (db *DB) GetTotalFlexForCustomer(customerName string) (time.Duration, error) {
	customer, err := db.GetCustomer(customerName)
	if err != nil {
		return 0, err
	}
	return customer.GetTotalFlex(), nil
}

// GetTotalFlexForAllCustomers returns the som of flex time for all Customers.
// If no customers, it returns time.Duration(0).
func (db *DB) GetTotalFlexForAllCustomers() time.Duration {
	if db.Customers == nil || db.Customers.Len() == 0 {
		return 0
	}
	var total time.Duration
	for _, customer := range db.Customers {
		total += customer.GetTotalFlex()
	}
	return total
}

// SetFlexForCustomer will either retrieve or add a Customer with the given name, depending on whether it exists or not.
// If customerName is blank, it will use default customer.
// It will then add a new flex Entry if no Entry with the same date exists for that customer.
// If overwrite is true, it will replace any Entry with a matching date.
// If overwrite is false, it will return an error if an Entry with a matching date is already present.
func (db *DB) SetFlexForCustomer(customerName string, date time.Time, amount time.Duration, overwrite bool) error {
	var err error
	var customer *Customer
	if customerName == "" {
		customer = db.GetDefaultCustomer()
	} else {
		customer, err = db.AddCustomer(customerName)
		if err != nil {
			log.Debug().Err(err).Send()
		}
	}
	if !customer.SetEntry(Entry{Date: date, Amount: amount}, overwrite) {
		return fmt.Errorf(
			"failed to add %v flex on %s for customer: %s (overwrite: %t)",
			amount,
			date.Format(ShortDateFormat),
			customer.Name,
			overwrite,
		)
	}
	return nil
}
