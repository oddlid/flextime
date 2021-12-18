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

	// TODO: Implement delete for a whole customer
	if customer != nil && date != nil {
		// delete specific date from specific customer
		//log.Debug().
		//	Str("customer_name", customer.Name).
		//	Time("date", *date).
		//	Msg("Delete entry with given date for given customer")
		if err := deleteSpecificDateFromCustomer(customer, *date); err != nil {
			return err
		}
	} else if customer != nil && all {
		// delete all entries from specific customer
		// log.Debug().
		// 	Str("customer_name", customer.Name).
		// 	Bool("all", all).
		// 	Msg("Delete all entries for given customer")
		if err := deleteAllEntriesFromCustomer(customer); err != nil {
			return err
		}
	} else if customer != nil && (from != nil || to != nil) {
		// delete date range from specific customer
		// log.Debug().
		// 	Str("customer_name", customer.Name).
		// 	Msg("Delete entries matching date range for given customer")
		if err := deleteDateRangeFromCustomer(customer, from, to); err != nil {
			return err
		}
	} else if all && customer == nil && date == nil && from == nil && to == nil {
		// delete all entries from all customers
		// log.Debug().
		// 	Bool("all", all).
		// 	Msg("Delete all entries from all customers!")
		if err := deleteAllEntriesFromAllCustomers(db); err != nil {
			return err
		}
	} else if all && customer == nil && date != nil && from == nil && to == nil {
		// delete specific date from all customers
		// log.Debug().
		// 	Bool("all", all).
		// 	Time("date", *date).
		// 	Msg("Delete specific date from all customers")
		if err := deleteSpecificDateFromAllCustomers(db, *date); err != nil {
			return err
		}
	} else if all && customer == nil && date == nil && (from != nil || to != nil) {
		// delete date range from all customers
		// log.Debug().
		// 	Bool("all", all).
		// 	Msg("Delete date range from all customers")
		if err := deleteDateRangeFromAllCustomers(db, from, to); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid combination of options for delete")
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

	switch customer {
	case nil:
		switch all {
		case true:
			switch date {
			case nil:
				if from != nil || to != nil {
					localLog("delete date range from all customers")
					// return deleteDateRangeFromAllCustomers(db, from, to)
				}
				localLog("delete all entries from all customers")
				// return deleteAllEntriesFromAllCustomers(db)
			default: // date is set
				localLog("delete specific date from all customers")
				// return deleteSpecificDateFromAllCustomers(db, *date)
			}
		default: // all is false / not set
			switch date {
			case nil:
			default: // date is set
			}
		}
	default: // customer is set
		switch all {
		case true:
			switch date {
			case nil:
				localLog("delete all entries from specific customer")
				// return deleteAllEntriesFromCustomer(customer)
			default:
			}
		default: // all is false / not set
			switch date {
			case nil:
				if from != nil || to != nil {
					localLog("delete date range from specific customer")
					// return deleteDateRangeFromCustomer(customer, from, to)
				}
				localLog("delete customer itself")
			default: // date is set
				localLog("delete specific date from specific customer")
				// return deleteSpecificDateFromCustomer(customer, *date)
			}
		}
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
