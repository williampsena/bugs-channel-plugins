// This package includes events contracts and behaviors
package event

import (
	"encoding/json"
	"errors"
)

// Represents an error when attempting to parse event to json
var ErrParseEventToJson = errors.New("an error occurred while attempting to event to JSON")

const (
	// The event origin from plugin/integration
	SentryOrigin = "sentry"
)

// Contains extra event info
type EventExtra = map[string]interface{}
type StackTrace = []map[string]interface{}

// Represents a Bugs Channel Event
type Event struct {
	// The event id
	ID string `json:"event_id"`
	// The raw event
	RawEvent string `json:"raw_event"`
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

func (e *Event) Json() (string, error) {
	b, err := json.Marshal(e)

	if err != nil {
		return "", errors.Join(ErrParseEventToJson, err)
	}

	return string(b), err
}
