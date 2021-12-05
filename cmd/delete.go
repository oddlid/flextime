package main

import (
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

	if customer != nil && date != nil {
		// delete specific date from specific customer
		log.Debug().
			Str("CustomerName", customer.Name).
			Time("Date", *date).
			Msg("Delete entry with given date for given customer")
	} else if customer != nil && all {
		// delete all entries from specific customer
		log.Debug().
			Str("CustomerName", customer.Name).
			Bool("All", all).
			Msg("Delete all entries for given customer")
	} else if customer != nil && (from != nil || to != nil) {
		// delete date range from specific customer
		log.Debug().
			Str("CustomerName", customer.Name).
			Msg("Delete entries matching date range for given customer")
	} else if all && customer == nil && date == nil && from == nil && to == nil {
		// delete all entries from all customers
		log.Debug().
			Bool("All", all).
			Msg("Delete all entries from all customers!")
	} else if all && customer == nil && date != nil && from == nil && to == nil {
		// delete specific date from all customers
		log.Debug().
			Bool("All", all).
			Time("Date", *date).
			Msg("Delete specific date from all customers")
	} else if all && customer == nil && date == nil && (from != nil || to != nil) {
		// delete date range from all customers
		log.Debug().
			Bool("All", all).
			Msg("Delete date range from all customers")
	} else {
		log.Info().Msg("Invalid combination of options for delete")
	}

	return nil
}
