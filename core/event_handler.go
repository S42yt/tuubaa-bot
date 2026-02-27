package core

var (
	eventHandlers []interface{}
)

func On(handler interface{}) {
	eventHandlers = append(eventHandlers, handler)
}

func AddEventHandler(handler interface{}) {
	eventHandlers = append(eventHandlers, handler)
}
