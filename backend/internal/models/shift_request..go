package models

type ShiftRequest struct {
	ID       int    `json:"id"`
	WorkerID int    `json:"worker_id"`
	ShiftID  int    `json:"shift_id"`
	Status   string `json:"status"` // pending, approved, rejected
}
