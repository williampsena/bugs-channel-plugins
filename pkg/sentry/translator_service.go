package sentry

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/williampsena/bugs-channel-plugins/pkg/event"
	"github.com/williampsena/bugs-channel-plugins/pkg/gzip"
	"github.com/williampsena/bugs-channel-plugins/pkg/text"
)

const (
	titleMaxLen = 100
)

type eventHeader struct {
	title string
	body  string
	kind  string
}

// Represents an error while decompressing the sentry envelope.
var ErrUnzipSentryEnvelope = errors.New("an error occurred while attempting to unzip the Sentry envelope")

// This function translates raw sentry events to recognized intermediate structures.
func TranslateEvents(project string, body io.Reader) ([]event.Event, error) {
	projectId, _ := strconv.Atoi(project)
	content, err := gzip.UnzipReader(body)

	if err != nil {
		return nil, errors.Join(ErrUnzipSentryEnvelope, err)
	}

	events, err := BuildSentryEvents(projectId, content)

	if err != nil {
		return nil, ErrParseSentryEvent
	}

	return parseToLocalEvents(events), nil
}

func parseToLocalEvents(sentryEvents []SentryEvent) []event.Event {
	var events []event.Event

	for _, e := range sentryEvents {
		events = append(events, parseToLocalEvent(&e))
	}

	return events
}

func parseToLocalEvent(sentryEvent *SentryEvent) event.Event {
	titleAndBody := buildTitleAndBody(sentryEvent)

	return event.Event{
		ID:          sentryEvent.ID,
		ServiceId:   strconv.Itoa(sentryEvent.Project),
		MetaId:      sentryEvent.EventGroupId,
		Platform:    sentryEvent.Platform,
		Environment: sentryEvent.Environment,
		Release:     sentryEvent.Release,
		ServerName:  sentryEvent.ServerName,
		Title:       titleAndBody.title,
		Body:        titleAndBody.body,
		StackTrace:  sentryEvent.StackTrace.Values,
		Kind:        titleAndBody.kind,
		Level:       sentryEvent.Level,
		Origin:      event.SentryOrigin,
		Tags:        buildEventTags(sentryEvent.Tags),
		Extra:       sentryEvent.Extra,
		Timestamp:   sentryEvent.Timestamp,
	}
}

func buildTitleAndBody(event *SentryEvent) eventHeader {
	if len(event.StackTrace.Values) > 0 {
		stackTrace := event.StackTrace.Values[0]

		return eventHeader{
			title: text.Truncate(fmt.Sprintf("%v: %v", stackTrace["type"], stackTrace["value"]), 100),
			body:  stackTrace["value"].(string),
			kind:  "error"}
	} else {
		return eventHeader{
			title: text.Truncate(event.Message, titleMaxLen),
			body:  event.Message,
			kind:  "event"}
	}
}

func buildEventTags(rawTags map[string]interface{}) []string {
	var tags []string

	for k, v := range rawTags {
		tags = append(tags, fmt.Sprintf("%v:%v", k, v))
	}

	return tags
}
