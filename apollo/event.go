package apollo

type ChangeType int

const (
	ADD ChangeType = iota

	MODIFY

	DELETE
)

type Listener interface {
	OnChange(event *ChangeEvent)
}

type Change struct {
	Key        string
	OldValue   string
	NewValue   string
	ChangeType ChangeType
}

type ChangeEvent struct {
	Namespace      string
	NotificationID int
	Changes        map[string]*Change
}

func onDelete(key, value string) *Change {
	return &Change{
		Key:        key,
		ChangeType: DELETE,
		OldValue:   value,
	}
}

func onModify(key, oldValue, newValue string) *Change {
	return &Change{
		Key:        key,
		ChangeType: MODIFY,
		OldValue:   oldValue,
		NewValue:   newValue,
	}
}

func onAdd(key, value string) *Change {
	return &Change{
		Key:        key,
		ChangeType: ADD,
		NewValue:   value,
	}
}
