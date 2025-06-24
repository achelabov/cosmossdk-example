package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:  DefaultParams(),
		NoteMap: []Note{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	noteIndexMap := make(map[string]struct{})

	for _, elem := range gs.NoteMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := noteIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for note")
		}
		noteIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
