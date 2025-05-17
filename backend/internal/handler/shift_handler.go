package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vincenttanjaya/go_worker_shift/internal/controller"
	"github.com/vincenttanjaya/go_worker_shift/internal/models"
)

type ShiftHandler struct {
	Ctrl *controller.ShiftController
}

func (h *ShiftHandler) ListShifts(w http.ResponseWriter, r *http.Request) {
	shifts, err := h.Ctrl.ListShifts()
	if err != nil {
		http.Error(w, "Cannot fetch shifts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shifts)
}

func (h *ShiftHandler) ListAvailableShifts(w http.ResponseWriter, r *http.Request) {
	shifts, err := h.Ctrl.ListAvailableShifts()
	if err != nil {
		http.Error(w, "Cannot fetch available shifts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shifts)
}

func (h *ShiftHandler) ListAssignedShifts(w http.ResponseWriter, r *http.Request) {
	shifts, err := h.Ctrl.ListAssignedShifts()
	if err != nil {
		http.Error(w, "Cannot fetch assigned shifts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shifts)
}

func (h *ShiftHandler) CreateShift(w http.ResponseWriter, r *http.Request) {
	var s models.Shift
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Println("Failed to decode JSON:", err)
		return
	}
	if err := h.Ctrl.CreateShift(&s); err != nil {
		http.Error(w, "Could not create shift", http.StatusInternalServerError)
		fmt.Println("Failed to create shift:", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *ShiftHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var s models.Shift
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	s.ID = id
	if err := h.Ctrl.UpdateShift(&s); err != nil {
		http.Error(w, "Could not update shift", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

func (h *ShiftHandler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.Ctrl.DeleteShift(id); err != nil {
		http.Error(w, "Could not delete shift", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
