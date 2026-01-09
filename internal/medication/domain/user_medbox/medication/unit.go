package medication

// Unit is a VO representing the unit of measurement
// for medication quantities.
type Unit string

// NewUnit creates validated medication unit.
func NewUnit(unit string) (Unit, error) {
	return Unit(unit), nil
}

// String implements stringer interface for a unit
func (u Unit) String() string {
	return string(u)
}
