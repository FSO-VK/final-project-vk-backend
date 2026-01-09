package medication

// ReleaseForm is a Value Object representing the physical form
// in which the medication is released.
type ReleaseForm string

// NewReleaseForm creates validated medication release form.
func NewReleaseForm(form string) (ReleaseForm, error) {
	return ReleaseForm(form), nil
}

// String returns the string representation of the release form.
func (f ReleaseForm) String() string {
	return string(f)
}
