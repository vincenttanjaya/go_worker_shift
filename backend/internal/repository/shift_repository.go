package repository

import (
	"database/sql"

	"github.com/vincenttanjaya/go_worker_shift/internal/models"
)

type ShiftRepository struct {
	DB *sql.DB
}

func (r *ShiftRepository) CreateShift(s *models.Shift) error {
	res, err := r.DB.Exec(
		"INSERT INTO shifts (date, start_time, end_time, role, location) VALUES (?, ?, ?, ?, ?)",
		s.Date, s.StartTime, s.EndTime, s.Role, s.Location,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	s.ID = int(id)
	return nil
}

func (r *ShiftRepository) UpdateShift(s *models.Shift) error {
	_, err := r.DB.Exec(
		"UPDATE shifts SET date=?, start_time=?, end_time=?, role=?, location=?, assigned_worker_id=? WHERE id=?",
		s.Date, s.StartTime, s.EndTime, s.Role, s.Location, s.AssignedWorkerID, s.ID,
	)
	return err
}

func (r *ShiftRepository) DeleteShift(id int) error {
	_, err := r.DB.Exec("DELETE FROM shifts WHERE id = ?", id)
	return err
}

func (r *ShiftRepository) GetShiftByID(id int) (*models.Shift, error) {
	row := r.DB.QueryRow("SELECT id, date, start_time, end_time, role, location, assigned_worker_id FROM shifts WHERE id = ?", id)
	var s models.Shift
	err := row.Scan(&s.ID, &s.Date, &s.StartTime, &s.EndTime, &s.Role, &s.Location, &s.AssignedWorkerID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShiftRepository) GetAllShifts() ([]models.Shift, error) {
	rows, err := r.DB.Query("SELECT id, date, start_time, end_time, role, location, assigned_worker_id FROM shifts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Shift
	for rows.Next() {
		var s models.Shift
		rows.Scan(&s.ID, &s.Date, &s.StartTime, &s.EndTime, &s.Role, &s.Location, &s.AssignedWorkerID)
		list = append(list, s)
	}
	return list, nil
}

func (r *ShiftRepository) GetUnassignedShifts() ([]models.Shift, error) {
	rows, err := r.DB.Query("SELECT id, date, start_time, end_time, role, location, assigned_worker_id FROM shifts WHERE assigned_worker_id IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Shift
	for rows.Next() {
		var s models.Shift
		rows.Scan(&s.ID, &s.Date, &s.StartTime, &s.EndTime, &s.Role, &s.Location, &s.AssignedWorkerID)
		list = append(list, s)
	}
	return list, nil
}

func (r *ShiftRepository) GetAssignedShifts() ([]models.Shift, error) {
	rows, err := r.DB.Query("SELECT id, date, start_time, end_time, role, location, assigned_worker_id FROM shifts WHERE assigned_worker_id IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Shift
	for rows.Next() {
		var s models.Shift
		rows.Scan(&s.ID, &s.Date, &s.StartTime, &s.EndTime, &s.Role, &s.Location, &s.AssignedWorkerID)
		list = append(list, s)
	}
	return list, nil
}

func (r *ShiftRequestRepository) ListRequestsByStatus(status string) ([]models.ShiftRequest, error) {
	var rows *sql.Rows
	var err error

	if status == "" {
		rows, err = r.DB.Query("SELECT id, worker_id, shift_id, status FROM shift_requests")
	} else {
		rows, err = r.DB.Query("SELECT id, worker_id, shift_id, status FROM shift_requests WHERE status = ?", status)
	}
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

func (r *ShiftRequestRepository) CountShiftsForWorkerOnDate(workerID int, date string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM shift_requests sr
		JOIN shifts s ON sr.shift_id = s.id
		WHERE sr.worker_id = ?
		  AND s.date = ?
		  AND (sr.status = 'pending' OR sr.status = 'approved')
	`
	err := r.DB.QueryRow(query, workerID, date).Scan(&count)
	return count, err
}

func (r *ShiftRequestRepository) CountShiftsForWorkerInWeek(workerID int, weekStr string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM shift_requests sr
		JOIN shifts s ON sr.shift_id = s.id
		WHERE sr.worker_id = ?
		  AND strftime('%Y-%W', s.date) = ?
		  AND (sr.status = 'pending' OR sr.status = 'approved')
	`
	err := r.DB.QueryRow(query, workerID, weekStr).Scan(&count)
	return count, err
}
