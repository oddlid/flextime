package main

import (
	"fmt"
	"time"

	"github.com/oddlid/flextime/flex"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func entryPointDelete(c *cli.Context) error {
	log.Debug().Msg("In entryPointDelete")

	fileName := c.String("file")
	customerName := c.String("customer")
	all := c.Bool("all")
	date := c.Timestamp("date")
	from := c.Timestamp("from")
	to := c.Timestamp("to")

	db, err := openDB(fileName)
	if err != nil {
		if db == nil {
			return err
		}
		if db.IsEmpty() {
			return flex.ErrEmptyDB
		}
	}

	var customer *flex.Customer
	if customerName != "" {
		customer, err = db.GetCustomer(customerName)
		if err != nil {
			return err
		}
	}

	if err := dispatchDeleteAction(all, db, customer, date, from, to); err != nil {
		return err
	}

	err = saveDB(db)
	if err != nil {
		return err
	}

	return nil
}

func dispatchDeleteAction(all bool, db *flex.DB, customer *flex.Customer, date, from, to *time.Time) error {
	fmtCustomer := func(c *flex.Customer) string {
		if c == nil {
			return "<nil>"
		}
		return c.Name
	}

	fmtDate := func(t *time.Time) string {
		if t == nil {
			return "<nil>"
		}
		return t.Format(flex.ShortDateFormat)
	}

	localLog := func(msg string) {
		log.Debug().
			Bool("all", all).
			Str("customer", fmtCustomer(customer)).
			Str("date", fmtDate(date)).
			Str("from", fmtDate(from)).
			Str("to", fmtDate(to)).
			Msg(msg)
	}

	switch all {
	case true:
		switch {
		case customer == nil && date == nil && from == nil && to == nil:
			localLog("delete all entries from all customers")
			return deleteAllEntriesFromAllCustomers(db)
		case customer == nil && date == nil && (from != nil || to != nil):
			localLog("delete date range from all customers")
			return deleteDateRangeFromAllCustomers(db, from, to)
		case customer == nil && date != nil && from == nil && to == nil:
			localLog("delete specific date from all customers")
			return deleteSpecificDateFromAllCustomers(db, *date)
		case customer != nil && date == nil && from == nil && to == nil:
			localLog("delete all entries from specific customer")
			return deleteAllEntriesFromCustomer(customer)
		default:
			localLog("Invalid option combination")
		}
	case false:
		switch {
		case customer != nil && date != nil && from == nil && to == nil:
			localLog("delete specific date from specific customer")
			return deleteSpecificDateFromCustomer(customer, *date)
		case customer != nil && date == nil && (from != nil || to != nil):
			localLog("delete date range from specific customer")
			return deleteDateRangeFromCustomer(customer, from, to)
		case customer != nil && date == nil && from == nil && to == nil:
			localLog("delete customer itself")
			return deleteCustomer(db, customer)
		default:
			localLog("Invalid option combination")
		}
	}

	return ErrInvalidOptionCombination
}

func deleteCustomer(db *flex.DB, customer *flex.Customer) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}
	if customer == nil {
		return flex.ErrNilCustomer
	}
	if !db.Customers.Delete(*customer) {
		return flex.ErrNoSuchCustomer
	}
	return nil
}

func deleteSpecificDateFromCustomer(customer *flex.Customer, date time.Time) error {
	if customer == nil {
		return flex.ErrNilCustomer
	}
	if customer.Entries == nil || customer.Entries.Len() == 0 {
		return flex.ErrNoEntries
	}
	if !customer.Entries.DeleteByDate(date) {
		return fmt.Errorf("%w: %s", flex.ErrNoEntry, date.Format(flex.ShortDateFormat))
	}

	log.Info().
		Str("customer_name", customer.Name).
		Str("date", date.Format(flex.ShortDateFormat)).
		Msg("Deleted entry with given date from customer")

	return nil
}

func deleteAllEntriesFromCustomer(customer *flex.Customer) error {
	if customer == nil {
		return flex.ErrNilCustomer
	}
	customer.Entries = make(flex.Entries, 0)

	log.Info().
		Str("customer_name", customer.Name).
		Msg("Deleted all entries from customer")

	return nil
}

func deleteDateRangeFromCustomer(customer *flex.Customer, from, to *time.Time) error {
	if customer == nil {
		return flex.ErrNilCustomer
	}
	if customer.Entries == nil || customer.Entries.Len() == 0 {
		return flex.ErrNoEntries
	}

	if from == nil {
		firstDate, err := customer.Entries.FirstDate()
		if err != nil {
			return err
		}
		from = firstDate
	}

	if to == nil {
		lastDate, err := customer.Entries.LastDate()
		if err != nil {
			return err
		}
		to = lastDate
	}

	filteredEntries := customer.Entries.FilterByNotInDateRange(*from, *to)
	entriesDeleted := customer.Entries.Len() - filteredEntries.Len()
	customer.Entries = filteredEntries

	log.Info().
		Str("customer_name", customer.Name).
		Int("entries_deleted", entriesDeleted).
		Msg("Deleted entries in date range from customer")

	return nil
}

func deleteAllEntriesFromAllCustomers(db *flex.DB) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}
	for _, customer := range db.Customers {
		customer.Entries = make(flex.Entries, 0)
	}

	log.Info().Msg("Deleted all entries from all customers")

	return nil
}

func deleteSpecificDateFromAllCustomers(db *flex.DB, date time.Time) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}

	entriesDeleted := 0
	for _, customer := range db.Customers {
		if customer.Entries.DeleteByDate(date) {
			entriesDeleted++
		}
	}

	log.Info().
		Time("date", date).
		Int("entries_deleted", entriesDeleted).
		Msg("Deleted entries matching date from all customers")

	return nil
}

func deleteDateRangeFromAllCustomers(db *flex.DB, from, to *time.Time) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}

	var firstDate *time.Time
	var lastDate *time.Time
	var err error
	entriesDeleted := 0
	for _, customer := range db.Customers {
		if customer.Entries == nil || customer.Entries.Len() == 0 {
			continue
		}
		if from == nil {
			firstDate, err = customer.Entries.FirstDate()
			if err != nil {
				continue
			}
		} else {
			firstDate = from
		}
		if to == nil {
			lastDate, err = customer.Entries.LastDate()
			if err != nil {
				continue
			}
		} else {
			lastDate = to
		}
		filteredEntries := customer.Entries.FilterByNotInDateRange(*firstDate, *lastDate)
		entriesDeleted += customer.Entries.Len() - filteredEntries.Len()
		customer.Entries = filteredEntries
	}

	log.Info().
		Int("entries_deleted", entriesDeleted).
		Msg("Deleted entries matching date range from all customers")

	return nil
}
