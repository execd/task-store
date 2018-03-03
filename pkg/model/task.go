package model

import (
	"encoding/json"
)

type TaskSpec struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Init  string `json:"init"`
}

func (s *TaskSpec) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *TaskSpec) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}
	return nil
}
