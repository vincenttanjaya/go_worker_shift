package repository

import (
	"database/sql"

	"github.com/vincenttanjaya/go_worker_shift/internal/models"
)

type ShiftRequestRepository struct {
	DB *sql.DB
}

func (r *ShiftRequestRepository) CreateShiftRequest(sr *models.ShiftRequest) error {
	_, err := r.DB.Exec(
		"INSERT INTO shift_requests (worker_id, shift_id, status) VALUES (?, ?, ?)",
		sr.WorkerID, sr.ShiftID, sr.Status,
	)
	return err
}

func (r *ShiftRequestRepository) UpdateShiftRequestStatus(requestID int, status string) error {
	_, err := r.DB.Exec(
		"UPDATE shift_requests SET status = ? WHERE id = ?",
		status, requestID,
	)
	return err
}

func (r *ShiftRequestRepository) GetShiftRequestByID(id int) (*models.ShiftRequest, error) {
	row := r.DB.QueryRow("SELECT id, worker_id, shift_id, status FROM shift_requests WHERE id = ?", id)
	var sr models.ShiftRequest
	err := row.Scan(&sr.ID, &sr.WorkerID, &sr.ShiftID, &sr.Status)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}

func (r *ShiftRequestRepository) ListRequestsByWorker(workerID int) ([]models.ShiftRequest, error) {
	rows, err := r.DB.Query("SELECT id, worker_id, shift_id, status FROM shift_requests WHERE worker_id = ?", workerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ShiftRequest
	for rows.Next() {
		var sr models.ShiftRequest
		rows.Scan(&sr.ID, &sr.WorkerID, &sr.ShiftID, &sr.Status)
		list = append(list, sr)
	}
	return list, nil
}
