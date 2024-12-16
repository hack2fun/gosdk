package types

func (gs GenesisState) Validate() error {
	if gs.Owner == "" {
		return ErrInvalidOwner
	}
	if gs.Denom == "" {
		return ErrInvalidDenom
	}
	if gs.ExchangeRate == 0 {
		return ErrInvalidRate
	}
	return nil
}
