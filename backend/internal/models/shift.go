package models

import (
	"database/sql"
	"encoding/json"
)

type Shift struct {
	ID               int           `json:"id"`
	Date             string        `json:"date"`
	StartTime        string        `json:"start_time"`
	EndTime          string        `json:"end_time"`
	Role             string        `json:"role"`
	Location         string        `json:"location,omitempty"`
	AssignedWorkerID sql.NullInt64 `json:"assigned_worker_id,omitempty"`
}

func (s *Shift) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID               int    `json:"id"`
		Date             string `json:"date"`
		StartTime        string `json:"start_time"`
		EndTime          string `json:"end_time"`
		Role             string `json:"role"`
		Location         string `json:"location"`
		AssignedWorkerID *int64 `json:"assigned_worker_id"` // pointer to allow null
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.ID = aux.ID
	s.Date = aux.Date
	s.StartTime = aux.StartTime
	s.EndTime = aux.EndTime
	s.Role = aux.Role
	s.Location = aux.Location

	if aux.AssignedWorkerID != nil {
		s.AssignedWorkerID = sql.NullInt64{Int64: *aux.AssignedWorkerID, Valid: true}
	} else {
		s.AssignedWorkerID = sql.NullInt64{Valid: false}
	}
	return nil
}
