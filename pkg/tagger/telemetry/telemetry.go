package telemetry

import "github.com/DataDog/datadog-agent/pkg/telemetry"

var (
	// StoredEntities tracks how many entities are stored in the tagger.
	StoredEntities = telemetry.NewGaugeWithOpts("tagger", "stored_entities",
		[]string{"source", "prefix"}, "Number of entities in the store.",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// UpdatedEntities tracks the number of updates to tagger entities.
	UpdatedEntities = telemetry.NewCounterWithOpts("tagger", "updated_entities",
		[]string{}, "Number of updates made to entities.",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Queries tracks the number of queries made against the tagger.
	Queries = telemetry.NewCounterWithOpts("tagger", "queries",
		[]string{"cardinality"}, "Queries made against the tagger.",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// StreamErrors tracks how many errors were received when streaming
	// tagger events.
	StreamErrors = telemetry.NewCounterWithOpts("tagger", "stream_errors",
		[]string{}, "Errors received when streaming tagger events",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Subscribers tracks how many subscribers the tagger has.
	Subscribers = telemetry.NewGaugeWithOpts("tagger", "subscribers",
		[]string{}, "Number of channels subscribing to tagger events",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Events tracks the number of tagger events being sent out.
	Events = telemetry.NewCounterWithOpts("tagger", "events",
		[]string{"cardinality"}, "Number of tagger events being sent out",
		telemetry.Options{NoDoubleUnderscoreSep: true})

	// Notifications tracks the number of times the tagger has sent a
	// notification with a group of events.
	Notifications = telemetry.NewCounterWithOpts("tagger", "sends",
		[]string{}, "Number of of times the tagger has sent a notification with a group of events",
		telemetry.Options{NoDoubleUnderscoreSep: true})
)
