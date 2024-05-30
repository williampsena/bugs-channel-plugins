package sentry

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
)

const (
	// The sentry kind event
	eventType = "event"
)

// Represents an error when attempting to parse json into a sentry event struct.
var ErrParseSentryEvent = errors.New("an error occurred while attempting to parse JSON to sentry event")

// Represents an error when sentry event is not implemented yet.
var ErrEventNotImplemented = errors.New("an error occurred when sentry event is not implemented yet")

// Represents an sentry untranslated field.
type DynamicSentryField = map[string]interface{}

// Represents an abstract sentry event meta.
type SentryEventMeta struct {
	// The event identifier
	ID string `json:"event_id"`
	// The event timestamp
	SentAt string `json:"sent_at"`
	// The event trace
	Trace SentryEventTrace `json:"trace"`
}

// Represents a sentry event trace
type SentryEventTrace struct {
	// The trace id
	TraceId string `json:"trace_id"`
	// The environment
	Environment string `json:"environment"`
	// The sentry public key
	PublicKey string `json:"public_key"`
	// The event sample rate
	SampleRate float32 `json:"sample_rate"`
}

// Represents an abstract sentry event header.
type SentryEventHeader struct {
	// The event type
	Type string `json:"type"`
	// The event content type
	ContentType string `json:"content_type"`
	// The event size
	Length int `json:"length"`
}

// Represents sentry event item
type SentryEvent struct {
	// The event id
	ID string `json:"event_id"`
	// The event project
	Project int `json:"project"`
	// The event group id
	EventGroupId string `json:"event_group_id"`
	// The environment
	Environment string `json:"environment"`
	// The language
	Platform string `json:"platform"`
	// The release version
	Release string `json:"release"`
	// Server name
	ServerName string `json:"server_name"`
	// Extra fields
	Extra DynamicSentryField `json:"extra"`
	// Error message
	Message string `json:"message"`
	// Timestamp
	Timestamp string `json:"timestamp"`
	// Severity
	Level string `json:"level"`
	// Request information
	Request DynamicSentryField `json:"request"`
	// User information
	User DynamicSentryField `json:"user"`
	// Used modules
	Modules DynamicSentryField `json:"modules"`
	// Bread crumbs
	Breadcrumbs DynamicSentryField `json:"breadcrumbs"`
	// Tag list
	Tags DynamicSentryField `json:"tags"`
	// Error stack trace
	StackTrace StackTrace `json:"exception"`
}

type StackTrace struct {
	// Values of stack trace
	Values []DynamicSentryField `json:"values"`
}

// Parse a JSON into a Sentry event meta to struct.
func NewSentryEventMetaFromJson(rawEvent []byte) (*SentryEventMeta, error) {
	var sentryEventMessage SentryEventMeta

	err := json.Unmarshal(rawEvent, &sentryEventMessage)

	if err != nil {
		return nil, errors.Join(ErrParseSentryEvent, err)
	}

	return &sentryEventMessage, nil
}

// Parse a JSON into a Sentry event header struct.
func NewSentryEventHeaderFromJson(rawEvent []byte) (*SentryEventHeader, error) {
	var sentryEventMessage SentryEventHeader

	err := json.Unmarshal(rawEvent, &sentryEventMessage)

	if err != nil {
		return nil, errors.Join(ErrParseSentryEvent, err)
	}

	return &sentryEventMessage, nil
}

// Parse a JSON into a Sentry event struct.
func NewSentryEventsFromEnvelope(projectId int, eventGroupId string, items [][]string) ([]SentryEvent, error) {
	var events []SentryEvent

	for _, headerAndItem := range items {
		if len(headerAndItem) < 2 {
			continue
		}

		rawHeader := headerAndItem[0]
		item := headerAndItem[1]

		header, err := NewSentryEventHeaderFromJson([]byte(rawHeader))

		if err != nil {
			return nil, err
		}

		if header.Type != eventType {
			log.Warnf("The sentry event type is not implemented yet: %v", header)
			continue
		}

		event, err := NewSentryEventFromJson(projectId, eventGroupId, []byte(item))

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// Parse a JSON into a Sentry event struct.
func NewSentryEventFromJson(projectId int, eventGroupId string, jsonContent []byte) (SentryEvent, error) {
	var sentryEvent SentryEvent

	err := json.Unmarshal(jsonContent, &sentryEvent)

	if err != nil {
		return SentryEvent{}, errors.Join(ErrParseSentryEvent, err)
	}

	sentryEvent.Project = projectId
	sentryEvent.EventGroupId = eventGroupId

	return sentryEvent, nil
}
