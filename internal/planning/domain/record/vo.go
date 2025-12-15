package record

// Status represents the status of the record.
type Status uint

// Enum of statuses.
const (
	StatusDraft Status = iota
	StatusTaken
	StatusMissed
)

func (s Status) String() string {
	switch s {
	case StatusDraft:
		return "Запланировано"
	case StatusTaken:
		return "Принято"
	case StatusMissed:
		return "Пропущено"
	}
	return "Невозможный статус"
}
