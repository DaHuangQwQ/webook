package events

type InconsistentEvent struct {
	Type      string
	ID        int64
	Direction string
}

const (
	// InconsistentEventTypeTargetMissing INSERT
	InconsistentEventTypeTargetMissing = "target_missing"
	// InconsistentEventTypeBaseMissing UPDATE
	InconsistentEventTypeBaseMissing = "base_missing"
	// InconsistentEventTypeNotEqual DELETE
	InconsistentEventTypeNotEqual = "neq"
)
