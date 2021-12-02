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
	Name        string  `json:"customer_name"`
	FlexEntries Entries `json:"flex_entries"`
}

type Customers []*Customer
type CustomersByName Customers

func (customer Customer) getTotalFlex() time.Duration {
	if customer.FlexEntries == nil {
		return time.Duration(0)
	}
	return customer.FlexEntries.getTotalFlex()
}

func (customer *Customer) getFlexEntry(date time.Time) (*Entry, error) {
	for _, fe := range customer.FlexEntries {
		if fe.Date.Equal(date) {
			return fe, nil
		}
	}
	return nil, fmt.Errorf("no flexentry for date: %v", date)
}

func (customer *Customer) setFlexEntry(fe Entry, overwrite bool) bool {
	foundAtIndex := -1
	for idx := range customer.FlexEntries {
		if fe.Date.Equal(customer.FlexEntries[idx].Date) {
			foundAtIndex = idx
			break
		}
	}
	if foundAtIndex == -1 {
		customer.FlexEntries = append(customer.FlexEntries, &fe)
		return true
	}
	if overwrite {
		customer.FlexEntries[foundAtIndex] = &fe
		return true
	}
	return false
}

func (customer Customer) Print(w io.Writer, indentString string, indentLevel int) {
	prefix := strings.Repeat(indentString, indentLevel)
	fmt.Fprintf(w, "%s%s:\n", prefix, customer.Name)
	if customer.FlexEntries != nil {
		customer.FlexEntries.Print(w, indentString, indentLevel+1)
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
