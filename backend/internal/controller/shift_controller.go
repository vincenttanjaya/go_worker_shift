package controller

import (
	"github.com/vincenttanjaya/go_worker_shift/internal/models"
	"github.com/vincenttanjaya/go_worker_shift/internal/repository"
)

type ShiftController struct {
	Repo *repository.ShiftRepository
}

func (c *ShiftController) ListShifts() ([]models.Shift, error) {
	return c.Repo.GetAllShifts()
}

func (c *ShiftController) ListAvailableShifts() ([]models.Shift, error) {
	return c.Repo.GetUnassignedShifts()
}

func (c *ShiftController) ListAssignedShifts() ([]models.Shift, error) {
	return c.Repo.GetAssignedShifts()
}

func (c *ShiftController) CreateShift(s *models.Shift) error {
	return c.Repo.CreateShift(s)
}

func (c *ShiftController) UpdateShift(s *models.Shift) error {
	return c.Repo.UpdateShift(s)
}

func (c *ShiftController) DeleteShift(id int) error {
	return c.Repo.DeleteShift(id)
}

func (c *ShiftController) GetShiftByID(id int) (*models.Shift, error) {
	return c.Repo.GetShiftByID(id)
}
