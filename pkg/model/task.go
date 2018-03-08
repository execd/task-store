package model

import (
	"encoding/json"
	"github.com/satori/go.uuid"
)

// Spec is the specification for a task
type Spec struct {
	ID       *uuid.UUID `json:"id"`
	Name     string     `json:"name"`
	Image    string     `json:"image"`
	Init     string     `json:"init"`
	InitArgs []string   `json:"initArgs"`
}

// MarshalBinary marshals a Spec
func (s *Spec) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary unmarshals a Spec
func (s *Spec) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// Info : task information
type Info struct {
	ID           *uuid.UUID     `json:"id"`
	Metadata     interface{}    `json:"metadata"`
	Succeeded    bool           `json:"succeeded"`
	FailureStats *FailureStatus `json:"failureStats,omitempty"`
}

// MarshalBinary marshals a Spec
func (i *Info) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

// UnmarshalBinary unmarshals a Spec
func (i *Info) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
}

// FailureStatus : contains status information for the result of an execution
type FailureStatus struct {
	Type        string          `json:"type"`               // The type of the "thing" that failed
	Name        string          `json:"name"`               // A name for this status
	Reason      string          `json:"reason,omitempty"`   // A general cause for the failure
	Message     string          `json:"message,omitempty"`  // A more detailed failure message
	ChildStatus []FailureStatus `json:"children,omitempty"` // In the case failures have a hierarchy - i.e. pod -> containers
}

// MarshalBinary marshals a Spec
func (s *FailureStatus) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary unmarshals a Spec
func (s *FailureStatus) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
