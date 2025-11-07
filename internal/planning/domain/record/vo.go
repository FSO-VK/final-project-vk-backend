package record

// Status represents the status of the record.
type Status uint

// Enum of statuses.
const (
	StatusDraft Status = iota
	StatusTaken
	StatusMissed
)
