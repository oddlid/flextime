package flex

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Customer struct {
	Name    string  `json:"customer_name"`
	Entries Entries `json:"flex_entries"`
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

// Print prints a strings representation of the Customer and its Entries to the given
// writer, prefixed by indentString * indentLevel.
// indentLevel is increased by 1 when passed on to the Entries Print function.
//func (customer Customer) Print(writer io.Writer, indentString string, indentLevel int) {
//	prefix := strings.Repeat(indentString, indentLevel)
//	fmt.Fprintf(writer, "%s%s:\n", prefix, customer.Name)
//	if customer.Entries != nil {
//		customer.Entries.Print(writer, indentString, indentLevel+1)
//	}
//}

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

// Sort sorts the Customers slice according to the given criteria
func (customers Customers) Sort(sortOrder CustomerSortOrder) {
	switch sortOrder {
	case CustomerSortByNameAscending:
		sort.Sort(CustomersByName(customers))
	case CustomerSortByNameDescending:
		sort.Sort(sort.Reverse(CustomersByName(customers)))
	}
}

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
