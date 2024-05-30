// This package includes sentry plugin integration
package sentry

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func BuildSentryEvents(projectId int, envelope string) ([]SentryEvent, error) {
	meta, items := DecodeSentryEnvelope(envelope)

	log.Debugf("meta: %v, items: %v", meta, items)

	eventMeta, err := NewSentryEventMetaFromJson([]byte(meta))

	if err != nil {
		return nil, err
	}

	events, err := NewSentryEventsFromEnvelope(projectId, eventMeta.ID, items)

	log.Debugf("project: %v, events: %v", projectId, events)

	return events, err
}

func DecodeSentryEnvelope(envelope string) (string, [][]string) {
	lines := strings.Split(envelope, "\n")
	meta := lines[0]
	items := lines[1:]

	return meta, chunkBy(items, 2)

}

func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}
