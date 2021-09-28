package example1

type observer interface {
	ID() string
	handleEvent(event string)
}
