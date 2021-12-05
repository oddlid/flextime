package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/oddlid/flextime/flex"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func entryPointList(c *cli.Context) error {
	log.Debug().Msg("In entryPointList")

	tfmt := func(t *time.Time) string {
		if t == nil {
			return "nil"
		}
		return t.Format(flex.ShortDateFormat)
	}

	fileName := c.String("file")
	customerName := c.String("customer")
	verbose := c.Bool("verbose")
	all := c.Bool("all")
	customerSort := c.String("customer-sort")
	entrySort := c.String("entry-sort")
	date := c.Timestamp("date")
	from := c.Timestamp("from")
	to := c.Timestamp("to")

	log.Debug().
		Str("FileName", fileName).
		Str("CustomerName", customerName).
		Bool("Verbose", verbose).
		Bool("All", all).
		Str("CustomerSort", customerSort).
		Str("EntrySort", entrySort).
		Str("Date", tfmt(date)).
		Str("From", tfmt(from)).
		Str("To", tfmt(to)).
		Send()

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

	customerSortValue := flex.CustomerNoSort
	entrySortValue := flex.EntryNoSort
	if value, ok := customerSortOrder[customerSort]; ok {
		customerSortValue = value
	}
	if value, ok := entrySortOrder[entrySort]; ok {
		entrySortValue = value
	}

	builder := strings.Builder{}

	if customer != nil && date != nil {
		if err := listSpecificDateForCustomer(&builder, customer, *date); err != nil {
			return err
		}
	} else if customer != nil && date == nil && from == nil && to == nil && !verbose {
		if err := listShortSummaryForCustomer(&builder, customer); err != nil {
			return err
		}
	} else if customer != nil && date == nil && from == nil && to == nil && verbose {
		if err := listAllForCustomer(&builder, customer, entrySortValue); err != nil {
			return err
		}
	} else if customer != nil && date == nil && (from != nil || to != nil) && !verbose {
		if err := listDateRangeSummaryForCustomer(&builder, customer, from, to); err != nil {
			return err
		}
	} else if customer != nil && date == nil && (from != nil || to != nil) && verbose {
		if err := listDateRangeForCustomer(&builder, customer, entrySortValue, from, to); err != nil {
			return err
		}
	} else if all && date != nil {
		if err := listSpecificDateForAllCustomers(&builder, db, *date, customerSortValue); err != nil {
			return err
		}
	} else if all && date == nil && from == nil && to == nil && !verbose {
		if err := listShortSummaryForAllCustomers(&builder, db, customerSortValue); err != nil {
			return err
		}
	} else if all && date == nil && from == nil && to == nil && verbose {
		if err := listAllForAllCustomers(&builder, db, customerSortValue, entrySortValue); err != nil {
			return err
		}
	} else if all && date == nil && (from != nil || to != nil) && !verbose {
		if err := listDateRangeSummaryForAllCustomers(&builder, db, customerSortValue, from, to); err != nil {
			return err
		}
	} else if all && date == nil && (from != nil || to != nil) && verbose {
		if err := listDateRangeForAllCustomers(&builder, db, customerSortValue, entrySortValue, from, to); err != nil {
			return err
		}
	}

	fmt.Print(builder.String())

	return nil
}

func listSpecificDateForCustomer(writer io.Writer, customer *flex.Customer, date time.Time) error {
	if customer == nil {
		return flex.ErrNilCustomer
	}
	entry, err := customer.GetEntry(date)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		writer,
		"%s:\n\t* %s: %v\n",
		customer.Name,
		entry.Date.Format(flex.ShortDateFormat),
		entry.Amount,
	)

	return nil
}

func listSpecificDateForAllCustomers(writer io.Writer, db *flex.DB, date time.Time, sortOrder flex.CustomerSortOrder) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}
	db.Customers.Sort(sortOrder)
	for _, customer := range db.Customers {
		entry, err := customer.GetEntry(date)
		if err != nil {
			continue
		}
		fmt.Fprintf(
			writer,
			"%s:\n\t* %s: %v\n",
			customer.Name,
			entry.Date.Format(flex.ShortDateFormat),
			entry.Amount,
		)
	}
	return nil
}

func listShortSummaryForCustomer(writer io.Writer, customer *flex.Customer) error {
	if customer == nil {
		return flex.ErrNilCustomer
	}

	fmt.Fprintf(
		writer,
		"%s: %v\n",
		customer.Name,
		customer.GetTotalFlex(),
	)
	return nil
}

func listShortSummaryForAllCustomers(writer io.Writer, db *flex.DB, sortOrder flex.CustomerSortOrder) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}
	db.Customers.Sort(sortOrder)
	formatStr := fmt.Sprintf("%s%d%s", "%-", db.Customers.LongestName(), "s : %v\n")
	for _, customer := range db.Customers {
		fmt.Fprintf(
			writer,
			formatStr,
			customer.Name,
			customer.GetTotalFlex(),
		)
	}
	return nil
}

func listAllForAllCustomers(writer io.Writer, db *flex.DB, customerSortOrder flex.CustomerSortOrder, entrySortOrder flex.EntrySortOrder) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}
	customerFormat := "%s: %v\n"
	entryFormat := "\t* %s: %v\n"
	db.Customers.Sort(customerSortOrder)
	for _, customer := range db.Customers {
		fmt.Fprintf(
			writer,
			customerFormat,
			customer.Name,
			customer.GetTotalFlex(),
		)
		customer.Entries.Sort(entrySortOrder)
		for _, entry := range customer.Entries {
			fmt.Fprintf(
				writer,
				entryFormat,
				entry.Date.Format(flex.ShortDateFormat),
				entry.Amount,
			)
		}
	}
	return nil
}

func listAllForCustomer(writer io.Writer, customer *flex.Customer, sortOrder flex.EntrySortOrder) error {
	if customer == nil {
		return flex.ErrNilCustomer
	}
	if customer.Entries == nil || customer.Entries.Len() == 0 {
		return flex.ErrNoEntries
	}

	customer.Entries.Sort(sortOrder)

	fmt.Fprintf(writer, "%s: %v\n", customer.Name, customer.GetTotalFlex())
	for _, entry := range customer.Entries {
		fmt.Fprintf(
			writer,
			"\t* %s: %v\n",
			entry.Date.Format(flex.ShortDateFormat),
			entry.Amount,
		)
	}

	return nil
}

func listDateRangeSummaryForCustomer(writer io.Writer, customer *flex.Customer, from, to *time.Time) error {
	if customer == nil {
		return flex.ErrNilCustomer
	}
	if customer.Entries == nil {
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

	filteredEntries := customer.Entries.FilterByDateRange(*from, *to)
	fmt.Fprintf(
		writer,
		"%s: %v\n",
		customer.Name,
		filteredEntries.GetTotalFlex(),
	)

	return nil
}

func listDateRangeSummaryForAllCustomers(writer io.Writer, db *flex.DB, sortOrder flex.CustomerSortOrder, from, to *time.Time) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}

	var firstDate *time.Time
	var lastDate *time.Time
	var err error
	formatStr := fmt.Sprintf("%s%d%s", "%-", db.Customers.LongestName(), "s : %v\n")
	db.Customers.Sort(sortOrder)
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
		filteredEntries := customer.Entries.FilterByDateRange(*firstDate, *lastDate)
		fmt.Fprintf(
			writer,
			formatStr,
			customer.Name,
			filteredEntries.GetTotalFlex(),
		)
	}
	return nil
}

func listDateRangeForCustomer(writer io.Writer, customer *flex.Customer, sortOrder flex.EntrySortOrder, from, to *time.Time) error {
	if customer == nil {
		return flex.ErrNilCustomer
	}
	if customer.Entries == nil {
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

	filteredEntries := customer.Entries.FilterByDateRange(*from, *to)
	filteredEntries.Sort(sortOrder)

	fmt.Fprintf(
		writer,
		"%s: %v\n",
		customer.Name,
		filteredEntries.GetTotalFlex(),
	)
	for _, entry := range filteredEntries {
		fmt.Fprintf(
			writer,
			"\t* %s: %v\n",
			entry.Date.Format(flex.ShortDateFormat),
			entry.Amount,
		)
	}

	return nil
}

func listDateRangeForAllCustomers(writer io.Writer, db *flex.DB, customerSortOrder flex.CustomerSortOrder, entrySortOrder flex.EntrySortOrder, from, to *time.Time) error {
	if db == nil || db.IsEmpty() {
		return flex.ErrEmptyDB
	}

	var firstDate *time.Time
	var lastDate *time.Time
	var err error
	customerFormat := "%s: %v\n"
	entryFormat := "\t* %s: %v\n"
	db.Customers.Sort(customerSortOrder)

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
		filteredEntries := customer.Entries.FilterByDateRange(*firstDate, *lastDate)
		fmt.Fprintf(
			writer,
			customerFormat,
			customer.Name,
			filteredEntries.GetTotalFlex(),
		)
		filteredEntries.Sort(entrySortOrder)
		for _, entry := range filteredEntries {
			fmt.Fprintf(
				writer,
				entryFormat,
				entry.Date.Format(flex.ShortDateFormat),
				entry.Amount,
			)
		}
	}

	return nil
}
