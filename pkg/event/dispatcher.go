// This package includes events contracts and behaviors
package event

import log "github.com/sirupsen/logrus"

const (
	SentryOrigin = "sentry"
)

// Contains extra event info
type EventExtra = map[string]interface{}
type StackTrace = []map[string]interface{}

// Represents a Bugs Channel Event
type Event struct {
	// The event id
	ID string `json:"event_id"`
	// The service
	ServiceId string `json:"service_id"`
	// The meta event id
	MetaId string `json:"meta_id"`
	// The language
	Platform string `json:"platform"`
	// The environment
	Environment string `json:"environment"`
	// The release version
	Release string `json:"release"`
	// Server name
	ServerName string `json:"server_name"`
	// The title message
	Title string `json:"title"`
	// The body message
	Body string `json:"body"`
	// The stack trace
	StackTrace StackTrace `json:"stack_trace"`
	// The event severity (info, warn, error)
	Level string `json:"level"`
	// The event kind (event or error)
	Kind string `json:"kind"`
	// Error origin (home, sentry, ...)
	Origin string `json:"origin"`
	// Tags list
	Tags []string `json:"tags"`
	// Extra fields
	Extra EventExtra `json:"extra"`
	// Timestamp
	Timestamp string `json:"timestamp"`
}

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
