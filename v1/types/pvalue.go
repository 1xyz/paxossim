package types

import "fmt"

type PValue struct {
	// Ballot number associated
	BN BallotNumber

	// Slot for this PValue
	Slot Slot

	// Command for this PValue
	Command Command
}

func (pValue PValue) String() string {
	return fmt.Sprintf("(%v, %v, %v)", pValue.BN, pValue.Slot, pValue.Command)
}

// PValueSet
type PValues map[PValue]bool

func (pvalues PValues) Set(value PValue) {
	_, ok := pvalues[value]
	if !ok {
		pvalues[value] = true
	}
}

func (pvalues PValues) Contains(value PValue) bool {
	_, ok := pvalues[value]
	return ok
}

func (pvalues PValues) Update(values PValues) {
	for v, _ := range values {
		pvalues.Set(v)
	}
}
