package event

type Event struct {
	Name    string
	Payload []any
}

type Listener struct {
	Name    string
	Handler func(...any) error
}
