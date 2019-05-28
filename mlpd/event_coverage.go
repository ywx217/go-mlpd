package mlpd

import "errors"

// EventCoverage coverage event
type EventCoverage struct {
	base *EventBase
	// type coverage format
	// type: TYPE_COVERAGE
	// exinfo: one of TYPE_COVERAGE_METHOD, TYPE_COVERAGE_STATEMENT, TYPE_COVERAGE_ASSEMBLY, TYPE_COVERAGE_CLASS
	// if exinfo == TYPE_COVERAGE_METHOD
	//  [assembly: string] name of assembly
	//  [class: string] name of the class
	//  [name: string] name of the method
	//  [signature: string] the signature of the method
	//  [filename: string] the file path of the file that contains this method
	//  [token: uleb128] the method token
	//  [method_id: uleb128] an ID for this data to associate with the buffers of TYPE_COVERAGE_STATEMENTS
	//  [len: uleb128] the number of TYPE_COVERAGE_BUFFERS associated with this method
	// if exinfo == TYPE_COVERAGE_STATEMENTS
	//  [method_id: uleb128] an the TYPE_COVERAGE_METHOD buffer to associate this with
	//  [offset: uleb128] the il offset relative to the previous offset
	//  [counter: uleb128] the counter for this instruction
	//  [line: uleb128] the line of filename containing this instruction
	//  [column: uleb128] the column containing this instruction
	// if exinfo == TYPE_COVERAGE_ASSEMBLY
	//  [name: string] assembly name
	//  [guid: string] assembly GUID
	//  [filename: string] assembly filename
	//  [number_of_methods: uleb128] the number of methods in this assembly
	//  [fully_covered: uleb128] the number of fully covered methods
	//  [partially_covered: uleb128] the number of partially covered methods
	//    currently partially_covered will always be 0, and fully_covered is the
	//    number of methods that are fully and partially covered.
	// if exinfo == TYPE_COVERAGE_CLASS
	//  [name: string] assembly name
	//  [class: string] class name
	//  [number_of_methods: uleb128] the number of methods in this class
	//  [fully_covered: uleb128] the number of fully covered methods
	//  [partially_covered: uleb128] the number of partially covered methods
	//    currently partially_covered will always be 0, and fully_covered is the
	//    number of methods that are fully and partially covered.
}

// IsEventCoverage find out if its an EventCoverage
func IsEventCoverage(base *EventBase) bool {
	return base.Type() == TypeCoverage
}

// ReadEventCoverage reads EventCoverage from reader
func ReadEventCoverage(r *MlpdReader, base *EventBase) (*EventCoverage, error) {
	return nil, errors.New("not implemented")
}
