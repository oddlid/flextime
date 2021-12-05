package flex

const (
	ShortDateFormat     = "2006-01-02"
	DefaultCustomerName = "default"
)

type EntrySortOrder uint8

const (
	EntryNoSort EntrySortOrder = iota
	EntrySortByDateAscending
	EntrySortByDateDescending
	EntrySortByAmountAscending
	EntrySortByAmountDescending
)

type CustomerSortOrder uint8

const (
	CustomerNoSort CustomerSortOrder = iota
	CustomerSortByNameAscending
	CustomerSortByNameDescending
)
