package enums

type Status string

const (
	StatusActive  Status = "active"
	StatusDeleted Status = "deleted"
)

func (s Status) IsValid() bool {
	return s == StatusActive || s == StatusDeleted
}
