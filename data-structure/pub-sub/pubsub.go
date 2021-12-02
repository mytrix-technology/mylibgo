package pubsub

type Client interface {
	Listen(channel string) Observable
	Publish(channel string, payload []byte) error
	Instance() string
}

type Observable interface {
	Subscribe(callback MessageCallback) Subscription
}

type Subscription interface {
	Unsubscribe() error
}

type MessageCallback func(Event)

type EventType int
const (
	EVENT_ERROR EventType = iota
	EVENT_LISTENER_CONNECTED
	EVENT_LISTENER_CONNECTION_ATTEMPT_FAILED
	EVENT_LISTENER_DISCONNECTED
	EVENT_LISTENER_RECONNECTED
	EVENT_NOTIFY
)

type Event struct {
	Type EventType
	Message Message
}

type Message struct {
	Source string
	Payload []byte
}
