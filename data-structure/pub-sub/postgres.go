package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/rs/xid"

	"github.com/mytrix-technology/mylibgo/datastore"
)

type postgreClient struct {
	instance string
	ds       *datastore.DataStore
	//Subscriptions map[string]Subscription
}

type postgreSubscription struct {
	event    chan Event
	channel  string
	ds       *datastore.DataStore
	listener *pq.Listener
}

func NewPostgresPubSubClient(host string, port int, user, password string) (Client, error) {
	config := datastore.DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}

	ds, err := datastore.NewPostgresDatastore(config)
	if err != nil {
		return nil, err
	}

	return &postgreClient{xid.New().String(), ds}, nil
}

func NewPostgresPubSubClientFromDatastore(ds *datastore.DataStore) (Client, error) {
	return &postgreClient{xid.New().String(), ds}, nil
}

func (pg *postgreClient) Instance() string {
	return pg.instance
}

func (pg *postgreClient) Listen(channel string) Observable {
	eventChan := make(chan Event)

	//pg.listeners[channel] = listener
	subs := &postgreSubscription{
		event:   eventChan,
		channel: channel,
		ds:      pg.ds,
	}

	//pg.subscriptions[channel] = subs

	return subs
}

func (pg *postgreClient) Publish(channel string, payload []byte) error {
	msg := Message{
		Source:  pg.instance,
		Payload: payload,
	}

	msggob := new(bytes.Buffer)
	if err := gob.NewEncoder(msggob).Encode(msg); err != nil {
		return fmt.Errorf("failed to serialize payload: %s", err)
	}

	msgStr := hex.EncodeToString(msggob.Bytes())
	if err := pg.ds.SendNotify(channel, msgStr); err != nil {
		return err
	}

	return nil
}

func (sub *postgreSubscription) Subscribe(callback MessageCallback) Subscription {
	listener, err := sub.ds.CreatePGListener(sub.channel, createNewListenerEventCallback(sub.event))
	if err != nil {
		event := Event{
			Type: EVENT_ERROR,
			Message: Message{
				Source:  "ERROR",
				Payload: []byte(err.Error()),
			},
		}

		callback(event)
		return sub
	}

	sub.listener = listener

	go listenForListenerEvent(sub.event, callback)
	go listen(listener, callback)

	return sub
}

func (sub *postgreSubscription) Unsubscribe() error {
	if sub.listener == nil {
		return nil
	}
	return sub.listener.Unlisten(sub.channel)
}

func listen(listener *pq.Listener, callback MessageCallback) {
	for {
		select {
		case notif := <-listener.Notify:
			buf, err := hex.DecodeString(notif.Extra)
			if err != nil {
				callback(Event{
					Type: EVENT_ERROR,
					Message: Message{
						Source:  "ERROR",
						Payload: []byte("failed to decode notif message"),
					},
				})
				continue
			}

			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&msg); err != nil {
				callback(Event{
					Type: EVENT_ERROR,
					Message: Message{
						Source:  "ERROR",
						Payload: []byte("failed to deserialize notif message"),
					},
				})
				continue
			}

			event := Event{
				Type:    EVENT_NOTIFY,
				Message: msg,
			}
			callback(event)
		case <-time.After(60 * time.Second):
			if err := listener.Ping(); err != nil {
				event := Event{
					Type: EVENT_ERROR,
					Message: Message{
						Source:  "ERROR",
						Payload: []byte(err.Error()),
					},
				}
				callback(event)
			}
		}
	}
}

func listenForListenerEvent(eventChan <-chan Event, callback MessageCallback) {
	for ev := range eventChan {
		callback(ev)
	}
}

func createNewListenerEventCallback(eventchan chan<- Event) pq.EventCallbackType {
	return func(ev pq.ListenerEventType, err error) {
		if err != nil {
			event := Event{
				Type: EVENT_ERROR,
				Message: Message{
					Source:  "ERROR",
					Payload: []byte(err.Error()),
				},
			}
			eventchan <- event
		}

		var eventType EventType
		switch ev {
		case pq.ListenerEventConnected:
			eventType = EVENT_LISTENER_CONNECTED
		case pq.ListenerEventConnectionAttemptFailed:
			eventType = EVENT_LISTENER_CONNECTION_ATTEMPT_FAILED
		case pq.ListenerEventDisconnected:
			eventType = EVENT_LISTENER_DISCONNECTED
		case pq.ListenerEventReconnected:
			eventType = EVENT_LISTENER_RECONNECTED
		}

		event := Event{
			Type: eventType,
			Message: Message{
				Source:  "",
				Payload: nil,
			},
		}

		eventchan <- event
	}
}
