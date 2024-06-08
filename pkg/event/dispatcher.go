// This package includes events contracts and behaviors
package event

import log "github.com/sirupsen/logrus"

// This contract represents the service in charge of dispatching events.
type EventsDispatcher interface {
	Dispatch(Event) error
	DispatchMany([]Event) error
}

// A simple implementation of Dispatcher, with the dispatcher target being a database, queue, http, tcp, or service open telemetry.
type EventsLoggerDispatcher struct{}

// Dispatch a event to stdout
func (d *EventsLoggerDispatcher) Dispatch(event Event) error {
	log.Infof("🍔 Ingest Event: %v", event)
	return nil
}

// Dispatch many events to stdout
func (d *EventsLoggerDispatcher) DispatchMany(events []Event) error {
	for _, e := range events {
		d.Dispatch(e)
	}

	return nil
}

// Creates a new logger dispatcher
func NewLoggerDispatcher() *EventsLoggerDispatcher {
	return &EventsLoggerDispatcher{}
}
