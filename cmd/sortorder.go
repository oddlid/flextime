package main

import (
	"strings"

	"github.com/oddlid/flextime/flex"
)

var customerSortOrder = map[string]flex.CustomerSortOrder{
	"asc":  flex.CustomerSortByNameAscending,
	"desc": flex.CustomerSortByNameDescending,
}

var entrySortOrder = map[string]flex.EntrySortOrder{
	"dateasc":    flex.EntrySortByDateAscending,
	"datedesc":   flex.EntrySortByDateDescending,
	"amountasc":  flex.EntrySortByAmountAscending,
	"amountdesc": flex.EntrySortByAmountDescending,
}

func customerSortOrderOptions() string {
	keys := make([]string, len(customerSortOrder))
	i := 0
	for key := range customerSortOrder {
		keys[i] = key
		i++
	}
	return strings.Join(keys, ", ")
}

func entrySortOrderOptions() string {
	keys := make([]string, len(entrySortOrder))
	i := 0
	for key := range entrySortOrder {
		keys[i] = key
		i++
	}
	return strings.Join(keys, ", ")
}
