package flex

import "errors"

var (
	ErrNoEntry          = errors.New("no entry for given date")
	ErrNoEntries        = errors.New("no entries for customer")
	ErrNoSuchCustomer   = errors.New("no such customer")
	ErrCustomerExists   = errors.New("customer already exists")
	ErrNilCustomer      = errors.New("customer is nil")
	ErrInvalidJSONInput = errors.New("invalid JSON input")
	ErrEmptyDB          = errors.New("empty flex database")
)
