package event

type Event struct {
	Name    string
	Payload []any
}

type Listener = func(...any) error
