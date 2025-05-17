package repository

import (
	"database/sql"

	"github.com/vincenttanjaya/go_worker_shift/internal/models"
)

type WorkerRepository struct {
	DB *sql.DB
}

func (r *WorkerRepository) CreateWorker(worker *models.Worker) error {
	res, err := r.DB.Exec("INSERT INTO workers (name) VALUES (?)", worker.Name)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	worker.ID = int(id)
	return nil
}

func (r *WorkerRepository) GetAllWorkers() ([]models.Worker, error) {
	rows, err := r.DB.Query("SELECT id, name FROM workers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var workers []models.Worker
	for rows.Next() {
		var w models.Worker
		err = rows.Scan(&w.ID, &w.Name)
		if err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}

func (r *WorkerRepository) GetWorkerByID(id int) (*models.Worker, error) {
	row := r.DB.QueryRow("SELECT id, name FROM workers WHERE id = ?", id)
	var w models.Worker
	err := row.Scan(&w.ID, &w.Name)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WorkerRepository) DeleteWorker(id int) error {
	_, err := r.DB.Exec("DELETE FROM workers WHERE id = ?", id)
	return err
}
