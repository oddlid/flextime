package flex

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type FlexDB struct {
	FileName  string    `json:"-"`
	Customers Customers `json:"customers"`
}

func (db *FlexDB) Save() error {
	if db.FileName == "" {
		return fmt.Errorf("FlexDB has no filename to save to")
	}
	file, err := os.Create(db.FileName)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	err = db.Encode(writer)
	if err != nil {
		return err
	}
	writer.Flush()
	log.Info().Str("filename", db.FileName).Msg("Saved FlexDB to file")
	return nil
}

func (db *FlexDB) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(db)
}

func (db *FlexDB) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(db)
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
