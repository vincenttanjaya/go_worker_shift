package controller

import (
	"github.com/vincenttanjaya/go_worker_shift/internal/models"
	"github.com/vincenttanjaya/go_worker_shift/internal/repository"
)

type WorkerController struct {
	Repo *repository.WorkerRepository
}

// Create (signup) a worker
func (c *WorkerController) SignupWorker(name string) (*models.Worker, error) {
	worker := &models.Worker{Name: name}
	err := c.Repo.CreateWorker(worker)
	if err != nil {
		return nil, err
	}
	return worker, nil
}

// (Optional) Get all workers
func (c *WorkerController) ListWorkers() ([]models.Worker, error) {
	return c.Repo.GetAllWorkers()
}

// (Optional) Get by ID
func (c *WorkerController) GetWorkerByID(id int) (*models.Worker, error) {
	return c.Repo.GetWorkerByID(id)
}

func (c *WorkerController) DeleteWorker(id int) error {
	return c.Repo.DeleteWorker(id)
}
