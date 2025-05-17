package controller

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vincenttanjaya/go_worker_shift/internal/models"
	"github.com/vincenttanjaya/go_worker_shift/internal/repository"
)

type ShiftRequestController struct {
	Repo       *repository.ShiftRequestRepository
	ShiftRepo  *repository.ShiftRepository
	WorkerRepo *repository.WorkerRepository
}

func (c *ShiftRequestController) RequestShift(workerID, shiftID int) error {
	// Find the shift
	shift, err := c.ShiftRepo.GetShiftByID(shiftID)
	if err != nil {
		return fmt.Errorf("shift not found")
	}

	// Check max per day
	countDay, err := c.Repo.CountShiftsForWorkerOnDate(workerID, shift.Date)
	if err != nil {
		return err
	}
	if countDay >= 1 {
		return fmt.Errorf("worker already has a shift on %s", shift.Date)
	}

	// Compute week as "YYYY-WW"
	shiftDate, err := time.Parse("2006-01-02", shift.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	_, week := shiftDate.ISOWeek()
	weekStr := fmt.Sprintf("%04d-%02d", shiftDate.Year(), week)

	countWeek, err := c.Repo.CountShiftsForWorkerInWeek(workerID, weekStr)
	if err != nil {
		return err
	}
	if countWeek >= 5 {
		return fmt.Errorf("worker cannot have more than 5 shifts in this week")
	}

	// Passed all checks: make the request
	sr := &models.ShiftRequest{WorkerID: workerID, ShiftID: shiftID, Status: "pending"}
	return c.Repo.CreateShiftRequest(sr)
}

func (c *ShiftRequestController) ApproveRequest(requestID int) error {
	req, err := c.Repo.GetShiftRequestByID(requestID)
	if err != nil {
		return err
	}

	shift, err := c.ShiftRepo.GetShiftByID(req.ShiftID)
	if err != nil {
		return err
	}
	if shift.AssignedWorkerID.Int64 != 0 {
		return nil
	}
	shift.AssignedWorkerID = sql.NullInt64{Int64: int64(req.WorkerID), Valid: true}
	if err := c.ShiftRepo.UpdateShift(shift); err != nil {
		return err
	}
	return c.Repo.UpdateShiftRequestStatus(requestID, "approved")
}

func (c *ShiftRequestController) RejectRequest(requestID int) error {
	return c.Repo.UpdateShiftRequestStatus(requestID, "rejected")
}

func (c *ShiftRequestController) ListRequestsByWorker(workerID int) ([]models.ShiftRequest, error) {
	return c.Repo.ListRequestsByWorker(workerID)
}

func (c *ShiftRequestController) ListRequestsByStatus(status string) ([]models.ShiftRequest, error) {
	return c.Repo.ListRequestsByStatus(status)
}
