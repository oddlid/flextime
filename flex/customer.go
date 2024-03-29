package flex

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Customer struct {
	Name    string  `json:"customer_name,omitempty"`
	Entries Entries `json:"flex_entries,omitempty"`
}

type Customers []*Customer
type CustomersByName Customers

// GetTotalFlex returns the sum of Amount for all Entries
func (customer Customer) GetTotalFlex() time.Duration {
	if customer.Entries == nil || customer.Entries.Len() == 0 {
		return time.Duration(0)
	}
	return customer.Entries.GetTotalFlex()
}

// GetEntry returns the entry matching the given date, or nil + error if not found
func (customer *Customer) GetEntry(date time.Time) (*Entry, error) {
	if customer.Entries == nil || customer.Entries.Len() == 0 {
		return nil, ErrNoEntries
	}
	idx := customer.Entries.IndexOf(Entry{Date: date})
	if idx == -1 {
		return nil, fmt.Errorf("%w: %v", ErrNoEntry, date)
	}
	return customer.Entries[idx], nil
}

// SetEntry will, if overwrite is false, add the given Entry to the customers Entries,
// if it does not already exist.
// If overwrite is true, it will replace the entry if already present.
// Returns true if an Entry is set, false if not.
func (customer *Customer) SetEntry(entry Entry, overwrite bool) bool {
	foundAtIndex := -1
	if customer.Entries != nil && customer.Entries.Len() > 0 {
		foundAtIndex = customer.Entries.IndexOf(entry)
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

func (customers Customers) Len() int {
	return len(customers)
}

func (customers Customers) IndexOf(customer Customer) int {
	for idx := range customers {
		if strings.EqualFold(customer.Name, customers[idx].Name) {
			return idx
		}
	}
	return -1
}

func (customers *Customers) Delete(customer Customer) bool {
	idx := customers.IndexOf(customer)
	if idx < 0 {
		return false
	}
	// use slow delete, preserving order
	copy((*customers)[idx:], (*customers)[idx+1:])
	(*customers)[len(*customers)-1] = nil
	*customers = (*customers)[:len(*customers)-1]

	return true
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

// Sort sorts the Customers slice according to the given criteria
func (customers Customers) Sort(sortOrder CustomerSortOrder) {
	switch sortOrder {
	case CustomerSortByNameAscending:
		sort.Sort(CustomersByName(customers))
	case CustomerSortByNameDescending:
		sort.Sort(sort.Reverse(CustomersByName(customers)))
	}
}

// LongestName returns the length of the longest customer name in the collection.
// Useful for alignment when printing.
func (customers Customers) LongestName() int {
	maxLen := 0
	for _, customer := range customers {
		clen := len(customer.Name)
		if clen > maxLen {
			maxLen = clen
		}
	}
	return maxLen
}
