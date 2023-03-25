package main

import (
	"fmt"
	"time"

	"github.com/oddlid/flextime/flex"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

/*

* No input filename:
	- Create new DB
	- Set its filename to "-" so that it will be encoded to stdout on save
* "-" as input filename:
	- Load DB from stdin
	- Set DB filename to "-" so it will be decoded to stdout on save
* Any other input filename:
	- Try to load a DB from the file
	- If invalid content or nonexistent file, create new DB
	- Set its filename to the input filename, so it will be saved back to where it came from

* No customer name:
	- Create and return default customer
* Nonexistent customer name:
	- Create and return a customer with the given name
* Existing customer name:
	- Return existing customer

* No date given:
	- Use todays date
* Nonexistent date given:
	- Add entry for given date, regardless of overwrite param
* Existing date given:
	- Only set new entry of overwrite is true

* No amount given
	- Do nothing, but log or return error

*/

func entryPointAdd(c *cli.Context) error {
	log.Debug().Msg("In entryPointAdd")

	fileName := c.String("file")
	customerName := c.String("customer")
	date := c.Timestamp("date")
	amount := c.Duration("amount")
	overwrite := c.Bool("overwrite")

	fmtDate := func(t *time.Time) string {
		if t == nil {
			return "<nil>"
		}
		return t.Format(flex.ShortDateFormat)
	}

	log.Debug().
		Str("file", fileName).
		Str("CustomerName", customerName).
		Str("Date", fmtDate(date)).
		Dur("Amount", amount).
		Bool("Overwrite", overwrite).
		Send()

	db, err := openDB(fileName)
	if err != nil {
		if db == nil {
			return err
		}
		log.Error().Err(err).Send()
	}

	if date == nil && customerName != "" {
		_, err := db.AddCustomer(customerName)
		if err != nil {
			return err
		}
	} else {
		if amount == time.Duration(0) {
			return fmt.Errorf("refusing to add entry with 0 flex amount")
		}
		err = db.SetFlexForCustomer(customerName, *date, amount, overwrite)
		if err != nil {
			return err
		}
	}

	err = saveDB(db)
	if err != nil {
		return err
	}

	return nil
}
