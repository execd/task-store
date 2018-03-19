package model

import (
	"encoding/json"
	"github.com/satori/go.uuid"
)

// Spec is the specification for a task
type Spec struct {
	ID       *uuid.UUID        `json:"id"`
	Metadata map[string]string `json:"metadata"`
	Image    string            `json:"image"`
	Init     string            `json:"init"`
	InitArgs []string          `json:"initArgs"`
}

// MarshalBinary marshals a Spec
func (s *Spec) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary unmarshals a Spec
func (s *Spec) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// Status : task information
type Status struct {
	ID       *uuid.UUID `json:"id"`
	Type     StatusType
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MarshalBinary marshals
func (i *Status) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

// UnmarshalBinary unmarshals
func (i *Status) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
}

type StatusType string

const (
	SucceededStatus StatusType = "Succeeded"
	FailedStatus    StatusType = "Failed"
	ExecutingStatus StatusType = "Executing"
)
