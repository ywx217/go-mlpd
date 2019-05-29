package mlpd

// EventIter iterate function for reading events
type EventIter func(*BufferHeader, *Event) error

// StopEventIterError stops iteration error
type StopEventIterError struct{}

// Error error interface
func (e *StopEventIterError) Error() string {
	return "stops event iteration"
}
