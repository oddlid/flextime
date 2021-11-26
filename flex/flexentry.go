package flex

import "time"

type FlexEntry struct {
	Date   time.Time     `json:"date"`
	Amount time.Duration `json:"amount"`
}

type FlexEntries []*FlexEntry

func (fes FlexEntries) getTotalFlex() time.Duration {
	var total time.Duration
	for _, e := range fes {
		total += e.Amount
	}
	return total
}

func (fes FlexEntries) Len() int {
	return len(fes)
}
