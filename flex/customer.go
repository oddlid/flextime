package flex

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

type CustomerSortOrder uint8

const (
	CustomerSortByNameAscending CustomerSortOrder = iota
	CustomerSortByNameDescending
)

type Customer struct {
	Name    string  `json:"customer_name"`
	Entries Entries `json:"flex_entries"`
}

type Customers []*Customer
type CustomersByName Customers

func (customer Customer) GetTotalFlex() time.Duration {
	if customer.Entries == nil {
		return time.Duration(0)
	}
	return customer.Entries.GetTotalFlex()
}

func (customer *Customer) GetEntry(date time.Time) (*Entry, error) {
	for _, entry := range customer.Entries {
		if entry.Date.Equal(date) {
			return entry, nil
		}
	}
	return nil, fmt.Errorf("no entry for date: %v", date)
}

func (customer *Customer) SetEntry(entry Entry, overwrite bool) bool {
	foundAtIndex := -1
	for idx := range customer.Entries {
		if entry.Date.Equal(customer.Entries[idx].Date) {
			foundAtIndex = idx
			break
		}
	}
	if foundAtIndex == -1 {
		customer.Entries = append(customer.Entries, &entry)
		return true
	}
	if overwrite {
		customer.Entries[foundAtIndex] = &entry
		return true
	}
	return false
}

func (customer Customer) Print(w io.Writer, indentString string, indentLevel int) {
	prefix := strings.Repeat(indentString, indentLevel)
	fmt.Fprintf(w, "%s%s:\n", prefix, customer.Name)
	if customer.Entries != nil {
		customer.Entries.Print(w, indentString, indentLevel+1)
	}
}

func (customers Customers) Len() int {
	return len(customers)
}

func (customersByName CustomersByName) Len() int {
	return len(customersByName)
}

func (customersByName CustomersByName) Swap(i, j int) {
	customersByName[i], customersByName[j] = customersByName[j], customersByName[i]
}

func (customersByName CustomersByName) Less(i, j int) bool {
	iName := customersByName[i].Name
	jName := customersByName[j].Name
	iNameLower := strings.ToLower(iName)
	jNameLower := strings.ToLower(jName)
	if iNameLower == jNameLower {
		return iName < jName
	}
	return iNameLower < jNameLower
}

func (customers Customers) Sort(sortOrder CustomerSortOrder) {
	switch sortOrder {
	case CustomerSortByNameAscending:
		sort.Sort(CustomersByName(customers))
	case CustomerSortByNameDescending:
		sort.Sort(sort.Reverse(CustomersByName(customers)))
	}
}
