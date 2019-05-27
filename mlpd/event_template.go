package mlpd

import "errors"

// EventTemplate alloc event
type EventTemplate struct {
	base *EventBase
}

// IsEventTemplate find out if its an EventTemplate
func IsEventTemplate(base *EventBase) bool {
	return base.Type() == TypeTemplate
}

// ReadEventTemplate reads EventTemplate from reader
func ReadEventTemplate(r *MlpdReader, base *EventBase) (*EventTemplate, error) {
	return nil, errors.New("not implemented")
}
